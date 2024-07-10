/*
 * This file was last modified at 2024-07-11 11:30 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * db.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

func DBConnect(dsn string) *pgxpool.Pool {

	config, err := pgxpool.ParseConfig(dsn)
	IfErrorThenPanic(err)
	slog.Debug(MSG, "dbConnect", "config parsed")

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		slog.Debug(MSG, "dbConnect", "Acquire connect ping...")
		if err = conn.Ping(ctx); err != nil {
			panic(err)
		}
		slog.Debug(MSG, "dbConnect", "Acquire connect Ok")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), config)
	IfErrorThenPanic(err)
	_, err = pool.Acquire(context.TODO())
	IfErrorThenPanic(err)
	slog.Debug(MSG, "dbConnect", "Acquire pool Ok")

	return pool
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
