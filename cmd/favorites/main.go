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
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/vskurikhin/gofavorites/internal/env"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS //

func main() {
	run(context.Background())
}

func run(ctx context.Context) {
	_ = fiber.New()
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
	dbMigrations(prop)
	for {
		select {
		case <-ctx.Done():
			return
		case <-sigint:
			return
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func dbMigrations(prop env.Properties) {
	pool := prop.DBPool()
	if pool == nil {
		slog.Warn(env.MSG+" dbMigrations", "pool", pool)
		return
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
	if err := goose.Version(db, "migrations"); err != nil {
		log.Fatal(err)
	}
	if err := db.Close(); err != nil {
		panic(err)
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
