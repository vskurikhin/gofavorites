/*
 * This file was last modified at 2024-08-03 12:13 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * transactional.go
 * $Id$
 */

package repo

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/env"
)

type TxPostgres[S domain.Suite] struct {
	cache cache[S]
	pool  *pgxpool.Pool
	sLog  *slog.Logger
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
		assetDft.cache = getAssetCache(prop)
		assetDft.pool = prop.DBPool()
		assetDft.sLog = prop.Logger()
	})
	return assetDft
}

func GetFavoritesTxPostgres(prop env.Properties) domain.Dft[*entity.Favorites] {
	onceFavoritesDft.Do(func() {
		favoritesDft = new(TxPostgres[*entity.Favorites])
		favoritesDft.cache = getFavoritesCache(prop)
		favoritesDft.pool = prop.DBPool()
		favoritesDft.sLog = prop.Logger()
	})
	return favoritesDft
}

func (p *TxPostgres[S]) DoDelete(ctx context.Context, entity S, scan func(domain.Scanner)) (err error) {

	err = p.cache.delete(entity)

	if err != nil {
		p.sLog.ErrorContext(ctx, env.MSG+"TxPostgres.DoDelete", "err", err)
	}
	return scanPostgreTxArgs(ctx, p.sLog, p.pool, entity.DeleteTxArgs(), scan)
}

func (p *TxPostgres[S]) DoUpsert(ctx context.Context, entity S, scan func(domain.Scanner)) (err error) {

	err = p.cache.delete(entity)

	if err != nil {
		p.sLog.ErrorContext(ctx, env.MSG+"TxPostgres.DoUpsert", "err", err)
	}
	err = scanPostgreTxArgs(ctx, p.sLog, p.pool, entity.UpsertTxArgs(), scan)

	if err == nil {
		_, e := p.cache.set(entity)
		if e != nil {
			p.sLog.ErrorContext(ctx, env.MSG+"TxPostgres.DoUpsert", "err", err)
		}
	}
	return err
}

func scanPostgreTxArgs(
	ctx context.Context,
	log *slog.Logger,
	pool *pgxpool.Pool,
	txArgs domain.TxArgs,
	scan func(domain.Scanner),
) (err error) {

	if pool == nil {
		return ErrBadPool
	}
	conn, err := pool.Acquire(ctx)

	for i := 1; err != nil && i < tries*increase; i += increase {
		time.Sleep(time.Duration(i) * time.Second)
		log.WarnContext(ctx, env.MSG+"TxPostgres.scanPostgreTxArgs", "msg", "retry pool acquire tx", "err", err)
		conn, err = pool.Acquire(ctx)
	}
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()

	if err != nil {
		return err
	}
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})

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
			log.ErrorContext(ctx, env.MSG+"TxPostgres.scanPostgreTxArgs", "err", err)
			return err
		}
		log.DebugContext(ctx, env.MSG+"TxPostgres.scanPostgreTxArgs", "commandTag", ct.String(), "rowsAffected", ct.RowsAffected())
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
