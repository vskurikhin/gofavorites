/*
 * This file was last modified at 2024-07-21 11:42 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * postgres.go
 * $Id$
 */
//!+

// Package repo TODO.
package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/env"
	"log/slog"
	"sync"
	"time"
)

const (
	increase = 1
	tries    = 3
)

type Postgres[E domain.Entity] struct {
	pool *pgxpool.Pool
}

var _ domain.Repo[domain.Entity] = (*Postgres[domain.Entity])(nil)
var (
	ErrBadPool        = fmt.Errorf("bad Database pool")
	onceAssetTypeRepo = new(sync.Once)
	assetTypeRepo     *Postgres[*entity.AssetType]
	onceUserRepo      = new(sync.Once)
	userRepo          *Postgres[*entity.User]
)

func GetAssetTypePostgresRepo(prop env.Properties) domain.Repo[*entity.AssetType] {
	onceAssetTypeRepo.Do(func() {
		assetTypeRepo = new(Postgres[*entity.AssetType])
		assetTypeRepo.pool = prop.DBPool()
	})
	return assetTypeRepo
}

func GetUserPostgresRepo(prop env.Properties) domain.Repo[*entity.User] {
	onceUserRepo.Do(func() {
		userRepo = new(Postgres[*entity.User])
		userRepo.pool = prop.DBPool()
	})
	return userRepo
}

func (p *Postgres[E]) Delete(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	err := scanPostgreSQL(ctx, p.pool, scan, entity.DeleteSQL(), entity.DeleteArgs()...)
	return entity, err
}

func (p *Postgres[E]) Get(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	err := scanPostgreSQL(ctx, p.pool, scan, entity.GetSQL(), entity.GetArgs()...)
	return entity, err
}

func (p *Postgres[E]) GetByFilter(ctx context.Context, entity E, scan func(domain.Scanner) E) ([]E, error) {

	result := make([]E, 0)
	rows, err := rowsPostgreSQL(ctx, p.pool, entity.GetByFilterSQL(), entity.GetByFilterArgs()...)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		e := scan(rows)
		result = append(result, e)
	}
	return result, err
}

func (p *Postgres[E]) Insert(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	err := scanPostgreSQL(ctx, p.pool, scan, entity.InsertSQL(), entity.InsertArgs()...)
	return entity, err
}

func (p *Postgres[E]) Update(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	err := scanPostgreSQL(ctx, p.pool, scan, entity.UpdateSQL(), entity.UpdateArgs()...)
	return entity, err
}

func scanPostgreSQL(ctx context.Context, pool *pgxpool.Pool, scan func(domain.Scanner), sql string, args ...any) error {

	row, err := rowPostgreSQL(ctx, pool, sql, args...)

	if err != nil {
		return err
	}
	scan(row)

	return nil
}

func rowPostgreSQL(ctx context.Context, pool *pgxpool.Pool, sql string, args ...any) (pgx.Row, error) {

	if pool == nil {
		return nil, ErrBadPool
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
	if conn == nil || err != nil {
		return nil, fmt.Errorf(" while connecting %v", err)
	}
	return conn.QueryRow(ctx, sql, args...), nil
}

func rowsPostgreSQL(ctx context.Context, pool *pgxpool.Pool, sql string, args ...any) (pgx.Rows, error) {

	if pool == nil {
		return nil, ErrBadPool
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
	if conn == nil || err != nil {
		return nil, fmt.Errorf(" while connecting %v", err)
	}
	return conn.Query(ctx, sql, args...)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
