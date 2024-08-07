/*
 * This file was last modified at 2024-08-03 14:56 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * batch.go
 * $Id$
 */

// Package batch TODO.
package batch

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/env"
)

const (
	increase = 1
	tries    = 3
)

type FavoritesInsertsBatch interface {
	Do(ctx context.Context, favorites []entity.Favorites, upk string) error
}

type postgres struct {
	pool *pgxpool.Pool
	sLog *slog.Logger
}

var (
	ErrBadPool    = fmt.Errorf("bad Database pool")
	oncePostgres  = new(sync.Once)
	batchPostgres *postgres
)
var _ FavoritesInsertsBatch = (*postgres)(nil)

func GetBatchPostgres(prop env.Properties) FavoritesInsertsBatch {
	oncePostgres.Do(func() {
		batchPostgres = new(postgres)
		batchPostgres.pool = prop.DBPool()
		batchPostgres.sLog = prop.Logger()
	})
	return batchPostgres
}

func (p *postgres) Do(ctx context.Context, favorites []entity.Favorites, upk string) error {

	if len(favorites) < 1 {
		return nil
	}
	maxVersion := slices.MaxFunc[[]entity.Favorites, entity.Favorites](favorites, func(x, y entity.Favorites) int {
		if x.User().Version() > y.User().Version() {
			return 1
		} else if x.User().Version() > y.User().Version() {
			return -1
		}
		return 0
	})
	sqls := make([]string, 0)
	args := make([][]any, 0)
	sqls = append(sqls, `
		INSERT INTO users
		(upk, version, created_at) VALUES ($1, $2, $3)
		ON CONFLICT (upk)
		DO UPDATE SET version = users.version + 1
	`)
	args = append(args, []any{upk, maxVersion.User().Version(), maxVersion.User().CreatedAt()})

	for _, f := range favorites {
		sqls = append(sqls, `
			INSERT INTO asset_types
			(name, created_at)
			VALUES ($1, $2)
			ON CONFLICT (name)
			DO NOTHING
		`)
		args = append(args, []any{f.Asset().AssetType().Name(), f.Asset().AssetType().CreatedAt()})
		sqls = append(sqls, `
			INSERT INTO assets
			(isin, asset_type, created_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (isin)
			DO UPDATE SET asset_type = $2, updated_at = $4
		`)
		args = append(args, []any{f.Asset().Isin(), f.Asset().AssetType().Name(), f.Asset().CreatedAt(), f.Asset().UpdatedAt()})
		sqls = append(sqls, `
			INSERT INTO favorites
    		(isin, user_upk, version, deleted, created_at)
    		VALUES ($1, $2, $3, NULL, $4)
			ON CONFLICT (isin, user_upk)
			DO UPDATE SET version = $3, deleted = NULL, updated_at = $5
		`)
		args = append(args, []any{f.Asset().Isin(), f.User().Upk(), f.Version(), f.CreatedAt(), f.UpdatedAt()})
	}
	return rowsPostgreSQL(ctx, p.sLog, p.pool, sqls, args)
}

func rowsPostgreSQL(
	ctx context.Context,
	log *slog.Logger,
	pool *pgxpool.Pool,
	sqls []string,
	args [][]any,
) error {

	if pool == nil {
		return ErrBadPool
	}
	conn, err := pool.Acquire(ctx)

	for i := 1; err != nil && i < tries*increase; i += increase {
		time.Sleep(time.Duration(i) * time.Second)
		log.WarnContext(ctx, env.MSG+"batch.rowsPostgreSQL", "msg", "retry pool acquire rows", "err", err)
		conn, err = pool.Acquire(ctx)
	}
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()
	if conn == nil || err != nil {
		return fmt.Errorf(" while connecting %v", err)
	}
	batch := &pgx.Batch{}

	for i, sql := range sqls {
		if i < len(args) && len(args[i]) > 0 {
			batch.Queue(sql, args[i]...)
		} else {
			batch.Queue(sql)
		}
	}
	br := conn.SendBatch(ctx, batch)

	for i := 0; i < len(sqls); i++ {
		ct, err := br.Exec()
		log.DebugContext(ctx, env.MSG+"batch.rowsPostgreSQL", "commandTag", ct.String(), "rowsAffected", ct.RowsAffected())
		if err != nil {
			return err
		}
	}
	if err = br.Close(); err != nil {
		return err
	}
	return nil
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
