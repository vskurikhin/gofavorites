/*
 * This file was last modified at 2024-07-17 11:20 by Victor N. Skurikhin.
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
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/memory"
	"github.com/vskurikhin/gofavorites/internal/env"
	"sync"
	"time"
)

type Cache interface {
	Reset() error
}

type CachedPostgres[E domain.Entity] struct {
	cache *memory.Storage
	pool  *pgxpool.Pool
	exp   time.Duration
}

var _ Cache = (*CachedPostgres[domain.Entity])(nil)
var _ domain.Repo[domain.Entity] = (*CachedPostgres[domain.Entity])(nil)
var (
	onceAssetCachedRepo     = new(sync.Once)
	assetCachedRepo         *CachedPostgres[*entity.Asset]
	onceAssetTypeCachedRepo = new(sync.Once)
	assetTypeCachedRepo     *CachedPostgres[*entity.AssetType]
	onceFavoritesCachedRepo = new(sync.Once)
	favoritesCachedRepo     *CachedPostgres[*entity.Favorites]
	onceUserCachedRepo      = new(sync.Once)
	userCachedRepo          *CachedPostgres[*entity.User]
)

func GetAssetCachedPostgresRepo(prop env.Properties) domain.Repo[*entity.Asset] {
	onceAssetCachedRepo.Do(func() {
		assetCachedRepo = new(CachedPostgres[*entity.Asset])
		assetCachedRepo.cache = memory.New(memory.Config{GCInterval: prop.CacheGCInterval()})
		assetCachedRepo.pool = prop.DBPool()
		assetCachedRepo.exp = prop.CacheExpire()
	})
	return assetCachedRepo
}

func GetAssetTypeCachedPostgresRepo(prop env.Properties) domain.Repo[*entity.AssetType] {
	onceAssetTypeCachedRepo.Do(func() {
		assetTypeCachedRepo = new(CachedPostgres[*entity.AssetType])
		assetTypeCachedRepo.cache = memory.New(memory.Config{GCInterval: prop.CacheGCInterval()})
		assetTypeCachedRepo.pool = prop.DBPool()
		assetTypeCachedRepo.exp = prop.CacheExpire()
	})
	return assetTypeCachedRepo
}

func GetFavoritesCachedPostgresRepo(prop env.Properties) domain.Repo[*entity.Favorites] {
	onceFavoritesCachedRepo.Do(func() {
		favoritesCachedRepo = new(CachedPostgres[*entity.Favorites])
		favoritesCachedRepo.cache = memory.New(memory.Config{GCInterval: prop.CacheGCInterval()})
		favoritesCachedRepo.pool = prop.DBPool()
		favoritesCachedRepo.exp = prop.CacheExpire()
	})
	return favoritesCachedRepo
}

func GetUserCachedPostgresRepo(prop env.Properties) domain.Repo[*entity.User] {
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
	err := rowScanPostgreSQL(ctx, p.pool, scan, entity.DeleteSQL(), entity.DeleteArgs()...)

	return entity, err
}

func (p *CachedPostgres[E]) Get(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {

	t := entity.Copy()
	data, err := p.cache.Get(entity.Key())

	if err == nil && data != nil {
		err = t.FromJSON(data)
		if err == nil {
			_ = entity.FromJSON(data)
			return entity, nil
		}
	}
	err = rowScanPostgreSQL(ctx, p.pool, scan, entity.GetSQL(), entity.GetArgs()...)

	return entity, err
}

func (p *CachedPostgres[E]) Insert(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {

	err := rowScanPostgreSQL(ctx, p.pool, scan, entity.InsertSQL(), entity.InsertArgs()...)

	if data, e := entity.ToJSON(); e == nil {
		_ = p.cache.Set(entity.Key(), data, time.Second)
	}
	return entity, err
}

func (p *CachedPostgres[E]) Update(ctx context.Context, entity E, scan func(domain.Scanner)) (E, error) {

	err := rowScanPostgreSQL(ctx, p.pool, scan, entity.UpdateSQL(), entity.UpdateArgs()...)

	if data, e := entity.ToJSON(); e == nil {
		_ = p.cache.Set(entity.Key(), data, time.Second)
	}
	return entity, err
}

func (p *CachedPostgres[E]) Reset() error {
	return p.cache.Reset()
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
