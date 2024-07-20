/*
 * This file was last modified at 2024-07-20 19:34 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * cached_postgres.go
 * $Id$
 */
//!+

// Package repo TODO.
package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/memory"
	"github.com/vskurikhin/gofavorites/internal/env"
	"sync"
	"time"
)

type cache[E domain.Entity] interface {
	delete(entity E) error
	get(entity E) (E, error)
	invalidate() error
	set(entity E) (E, error)
}

type CachedPostgres[E domain.Entity] struct {
	cache *memory.Storage
	pool  *pgxpool.Pool
	exp   time.Duration
}

var _ cache[domain.Entity] = (*CachedPostgres[domain.Entity])(nil)
var _ domain.Repo[domain.Entity] = (*CachedPostgres[domain.Entity])(nil)
var (
	ErrNotFound             = fmt.Errorf("not found in cache")
	onceAssetCachedRepo     = new(sync.Once)
	assetCachedRepo         *CachedPostgres[*entity.Asset]
	onceAssetTypeCachedRepo = new(sync.Once)
	assetTypeCachedRepo     *CachedPostgres[*entity.AssetType]
	onceFavoritesCachedRepo = new(sync.Once)
	favoritesCachedRepo     *CachedPostgres[*entity.Favorites]
	onceUserCachedRepo      = new(sync.Once)
	userCachedRepo          *CachedPostgres[*entity.User]
)

func getAssetCache(prop env.Properties) cache[*entity.Asset] {
	return getAssetCachedPostgresRepo(prop)
}

func GetAssetPostgresCachedRepo(prop env.Properties) domain.Repo[*entity.Asset] {
	return getAssetCachedPostgresRepo(prop)
}

func getAssetCachedPostgresRepo(prop env.Properties) *CachedPostgres[*entity.Asset] {
	onceAssetCachedRepo.Do(func() {
		assetCachedRepo = new(CachedPostgres[*entity.Asset])
		assetCachedRepo.cache = memory.New(memory.Config{GCInterval: prop.CacheGCInterval()})
		assetCachedRepo.pool = prop.DBPool()
		assetCachedRepo.exp = prop.CacheExpire()
	})
	return assetCachedRepo
}

func GetAssetTypePostgresCachedRepo(prop env.Properties) domain.Repo[*entity.AssetType] {
	return getAssetTypeCachedPostgresRepo(prop)
}

func getAssetTypeCachedPostgresRepo(prop env.Properties) *CachedPostgres[*entity.AssetType] {
	onceAssetTypeCachedRepo.Do(func() {
		assetTypeCachedRepo = new(CachedPostgres[*entity.AssetType])
		assetTypeCachedRepo.cache = memory.New(memory.Config{GCInterval: prop.CacheGCInterval()})
		assetTypeCachedRepo.pool = prop.DBPool()
		assetTypeCachedRepo.exp = prop.CacheExpire()
	})
	return assetTypeCachedRepo
}

func getFavoritesCache(prop env.Properties) cache[*entity.Favorites] {
	return getFavoritesCachedPostgresRepo(prop)
}

func GetFavoritesPostgresCachedRepo(prop env.Properties) domain.Repo[*entity.Favorites] {
	return getFavoritesCachedPostgresRepo(prop)
}

func getFavoritesCachedPostgresRepo(prop env.Properties) *CachedPostgres[*entity.Favorites] {
	onceFavoritesCachedRepo.Do(func() {
		favoritesCachedRepo = new(CachedPostgres[*entity.Favorites])
		favoritesCachedRepo.cache = memory.New(memory.Config{GCInterval: prop.CacheGCInterval()})
		favoritesCachedRepo.pool = prop.DBPool()
		favoritesCachedRepo.exp = prop.CacheExpire()
	})
	return favoritesCachedRepo
}

func GetUserPostgresCachedRepo(prop env.Properties) domain.Repo[*entity.User] {
	return getUserCachedPostgresRepo(prop)
}

func getUserCachedPostgresRepo(prop env.Properties) *CachedPostgres[*entity.User] {
	onceUserCachedRepo.Do(func() {
		userCachedRepo = new(CachedPostgres[*entity.User])
		userCachedRepo.cache = memory.New(memory.Config{GCInterval: prop.CacheGCInterval()})
		userCachedRepo.pool = prop.DBPool()
		userCachedRepo.exp = prop.CacheExpire()
	})
	return userCachedRepo
}

func (p *CachedPostgres[E]) Delete(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {

	_ = p.cache.Delete(entity.Key())
	err := scanPostgreSQL(ctx, p.pool, scan, entity.DeleteSQL(), entity.DeleteArgs()...)

	return entity, err
}

func (p *CachedPostgres[E]) Get(ctx context.Context, entity E, scan func(domain.Scanner)) (e E, err error) {

	entity, err = p.get(entity)

	if err == nil {
		return entity, nil
	}
	err = scanPostgreSQL(ctx, p.pool, scan, entity.GetSQL(), entity.GetArgs()...)

	return entity, err
}

func (p *CachedPostgres[E]) GetByFilter(ctx context.Context, entity E, scan func(domain.Scanner) E) ([]E, error) {

	result := make([]E, 0)
	rows, err := rowsPostgreSQL(ctx, p.pool, entity.GetByFilterSQL(), entity.GetByFilterArgs()...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		e := scan(rows)
		result = append(result, e)
		if data, err := e.ToJSON(); err == nil {
			_ = p.cache.Set(e.Key(), data, p.exp)
		}
	}
	return result, err
}

func (p *CachedPostgres[E]) Insert(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {

	err := scanPostgreSQL(ctx, p.pool, scan, entity.InsertSQL(), entity.InsertArgs()...)

	if err == nil {
		return p.set(entity)
	}
	return entity, err
}

func (p *CachedPostgres[E]) Update(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {

	err := scanPostgreSQL(ctx, p.pool, scan, entity.UpdateSQL(), entity.UpdateArgs()...)

	if data, e := entity.ToJSON(); e == nil {
		_ = p.cache.Set(entity.Key(), data, p.exp)
	}
	return entity, err
}

func (p *CachedPostgres[E]) delete(entity E) error {
	return p.cache.Delete(entity.Key())
}

func (p *CachedPostgres[E]) get(entity E) (E, error) {

	t := entity.Copy()
	data, err := p.cache.Get(entity.Key())

	if err == nil && data != nil && len(data) <= 0 {
		if err = t.FromJSON(data); err != nil {
			return entity, err
		}
		_ = entity.FromJSON(data)
		return entity, nil
	}
	return entity, ErrNotFound
}

func (p *CachedPostgres[E]) invalidate() error {
	return p.cache.Invalidate()
}

func (p *CachedPostgres[E]) set(entity E) (e E, err error) {

	if data, e := entity.ToJSON(); e == nil {
		err = p.cache.Set(entity.Key(), data, p.exp)
	}
	return entity, err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
