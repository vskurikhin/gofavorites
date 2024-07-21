/*
 * This file was last modified at 2024-07-21 11:42 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * transactional.go
 * $Id$
 */

package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/env"
	"log/slog"
	"sync"
	"time"
)

type TxPostgres[S domain.Suite] struct {
	pool  *pgxpool.Pool
	cache cache[S]
}

var _ domain.Dft[domain.Suite] = (*TxPostgres[domain.Suite])(nil)
var (
	onceAssetDft     = new(sync.Once)
	assetDft         *TxPostgres[*entity.Asset]
	onceFavoritesDft = new(sync.Once)
	favoritesDft     *TxPostgres[*entity.Favorites]
)

func GetAssetTxPostgres(prop env.Properties) domain.Dft[*entity.Asset] {
	onceAssetDft.Do(func() {
		assetDft = new(TxPostgres[*entity.Asset])
		assetDft.pool = prop.DBPool()
		assetDft.cache = getAssetCache(prop)
	})
	return assetDft
}

func GetFavoritesTxPostgres(prop env.Properties) domain.Dft[*entity.Favorites] {
	onceFavoritesDft.Do(func() {
		favoritesDft = new(TxPostgres[*entity.Favorites])
		favoritesDft.pool = prop.DBPool()
		favoritesDft.cache = getFavoritesCache(prop)
	})
	return favoritesDft
}

func (p *TxPostgres[S]) DoDelete(ctx context.Context, entity S, scan func(domain.Scanner)) (err error) {

	err = p.cache.delete(entity)

	if err != nil {
		slog.Warn(env.MSG+" DoDelete", "err", err)
	}
	return scanPostgreTxArgs(ctx, p.pool, entity.DeleteTxArgs(), scan)
}

func (p *TxPostgres[S]) DoUpsert(ctx context.Context, entity S, scan func(domain.Scanner)) (err error) {

	err = p.cache.delete(entity)

	if err != nil {
		slog.Warn(env.MSG+" DoUpsert", "err", err)
	}
	err = scanPostgreTxArgs(ctx, p.pool, entity.UpsertTxArgs(), scan)

	if err == nil {
		_, e := p.cache.set(entity)
		if e != nil {
			slog.Warn(env.MSG+" DoUpsert", "err", err)
		}
	}
	return err
}

func scanPostgreTxArgs(ctx context.Context, pool *pgxpool.Pool, txArgs domain.TxArgs, scan func(domain.Scanner)) (err error) {

	if pool == nil {
		return ErrBadPool
	}
	conn, err := pool.Acquire(ctx)

	for i := 1; err != nil && i < tries*increase; i += increase {
		time.Sleep(time.Duration(i) * time.Second)
		slog.Warn(env.MSG+" retry pool acquire", "err", err)
		conn, err = pool.Acquire(ctx)
	}
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()
	for i := 0; i < len(txArgs.SQLs)-1; i++ {
		args := getTxArgs(txArgs, i)
		ct, err := tx.Exec(ctx, txArgs.SQLs[i], args...)

		if err != nil {
			return err
		}
		slog.Debug(
			env.MSG+" read",
			"commandTag", ct.String(),
			"commandTag.RowsAffected()", ct.RowsAffected(),
		)
	}
	query := txArgs.SQLs[len(txArgs.SQLs)-1]
	args := getTxArgs(txArgs, len(txArgs.SQLs)-1)
	row := tx.QueryRow(ctx, query, args...)
	scan(row)

	return err
}

func getTxArgs(txArgs domain.TxArgs, i int) []any {

	var args []any

	if i < len(txArgs.Args) {
		args = txArgs.Args[i]
	} else {
		args = []any{}
	}
	return args
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
