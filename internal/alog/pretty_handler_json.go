/*
 * This file was last modified at 2024-07-31 13:53 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * pretty_handler_json.go
 * $Id$
 */
//!+

// Package tool TODO.
package alog

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/hokaccha/go-prettyjson"

	"github.com/google/uuid"
)

type PrettyHandlerJSON struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandlerJSON) Handle(ctx context.Context, r slog.Record) error {

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		switch x := a.Value.Any().(type) {
		case error:
			fields[a.Key] = x.Error()
		case fmt.Stringer:
			fields[a.Key] = x.String()
		default:
			fields[a.Key] = a.Value.Any()
		}
		return true
	})

	fields["level"] = r.Level.String()
	fields["pid"] = os.Getpid()
	fields["time"] = r.Time.Format("15:05:05.000000")
	fields["msg"] = r.Message
	id := uuid.Max.String()
	if ri := ctx.Value("request-id"); ri != nil {
		if requestId, ok := ri.(string); ok {
			id = requestId
		}
	}
	fields["requestId"] = id

	b, err := prettyjson.Marshal(fields)
	if err != nil {
		return err
	}
	h.l.Println(string(b))

	return nil
}

func NewPrettyHandlerJSON(out io.Writer, opts PrettyHandlerOptions) *PrettyHandlerJSON {
	h := &PrettyHandlerJSON{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
