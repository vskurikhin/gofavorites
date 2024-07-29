/*
 * This file was last modified at 2024-07-29 13:42 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * trace.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func TraceInOut(ctx context.Context, name, format string, values ...any) func() {

	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		return func() {}
	}
	start := time.Now()
	f := fmt.Sprintf(" in: %s(", name) + format + ")"
	i := fmt.Sprintf(f, values...)
	slog.Debug(fmt.Sprintf("%s %s", MSG, name), "trace", i)

	return func() {
		o := fmt.Sprintf("out: %s [%s]", name, time.Since(start))
		slog.Debug(fmt.Sprintf("%s %s", MSG, name), "trace", o)
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
