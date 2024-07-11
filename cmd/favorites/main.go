/*
 * Copyright text:
 * This file was last modified at 2024-07-10 21:53 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * main.go
 * $Id$
 */
//!+

package main

import (
	"context"
	"fmt"
	"github.com/vskurikhin/gofavorites/internal/env"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	run(context.Background())
}

func run(ctx context.Context) {

	slog.Info(env.MSG,
		"build_version", buildVersion,
		"build_date", buildDate,
		"build_commit", buildCommit,
	)
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	// запускаем горутину обработки пойманных прерываний
	prop := env.GetProperties()
	_, _ = fmt.Fprintf(os.Stderr, "PROPERTIES%s\n", prop)
	for {
		select {
		case <-ctx.Done():
			return
		case <-sigint:
			return
		}
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
