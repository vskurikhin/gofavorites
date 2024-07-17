/*
 * This file was last modified at 2024-07-15 18:40 by Victor N. Skurikhin.
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
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"testing"
	"time"
)

func TestCachedPostgresRepos(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #0 AssetType Postgres Repo",
			fRun: testAssetTypeCachedPostgresRepo,
		},
		{
			name: "positive test #1 Asset Postgres Repo",
			fRun: testAssetCachedPostgresRepo,
		},
		{
			name: "positive test #2 User Cached Postgres Repo",
			fRun: testUserCachedPostgresRepo,
		},
		{
			name: "positive test #3 userRepo",
			fRun: testFavoritesCachedPostgresRepo,
		},
		{
			name: "positive test #999 clear",
			fRun: testClear,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testAssetTypeCachedPostgresRepo(t *testing.T) {
	prop := env.GetProperties()
	assetTypePostgresRepo := GetAssetTypeCachedPostgresRepo(prop)
	cache := assetTypePostgresRepo.(*CachedPostgres[*entity.AssetType])
	assetType = tool.RandStringBytes(32)
	expected := entity.NewAssetType(assetType, time.Time{})
	err := expected.Insert(context.TODO(), assetTypePostgresRepo)
	assert.Nil(t, err)
	assert.False(t, expected.Deleted().Valid)

	data1, err := cache.cache.Get(assetType)
	assert.Nil(t, err)
	assert.NotNil(t, data1)
	data2, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, data1, data2)

	got, err := entity.GetAssetType(context.TODO(), assetTypePostgresRepo, assetType)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.CreatedAt(), got.CreatedAt())
	assert.Equal(t, expected.Name(), got.Name())
	assert.Equal(t, expected.Deleted(), got.Deleted())
	assert.False(t, expected.UpdatedAt().Valid)
	err = expected.Update(context.TODO(), assetTypePostgresRepo)
	assert.Nil(t, err)
	assert.False(t, expected.Deleted().Valid)
	assert.True(t, expected.UpdatedAt().Valid)
	err = expected.Delete(context.TODO(), assetTypePostgresRepo)
	assert.True(t, expected.Deleted().Valid)
	assert.True(t, expected.Deleted().Bool)
	assert.Nil(t, err)
	conn, err := prop.DBPool().Acquire(context.TODO())
	defer func() { conn.Release() }()
	assert.Nil(t, err)
}

func testAssetCachedPostgresRepo(t *testing.T) {
	prop := env.GetProperties()
	assetTypePostgresRepo := GetAssetTypePostgresRepo(prop)
	assetPostgresRepo := GetAssetCachedPostgresRepo(prop)
	cache := assetPostgresRepo.(*CachedPostgres[*entity.Asset])
	at, err := entity.GetAssetType(context.TODO(), assetTypePostgresRepo, assetType)
	assert.Nil(t, err)
	isin = tool.RandStringBytes(32)
	expected := entity.NewAsset(isin, at, time.Time{})
	err = expected.Insert(context.TODO(), assetPostgresRepo)
	assert.Nil(t, err)
	assert.False(t, expected.Deleted().Valid)

	data1, err := cache.cache.Get(isin)
	assert.Nil(t, err)
	assert.NotNil(t, data1)
	data2, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, data1, data2)

	got, err := entity.GetAsset(context.TODO(), assetPostgresRepo, isin)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.CreatedAt(), got.CreatedAt())
	assert.Equal(t, expected.Isin(), got.Isin())
	assert.Equal(t, expected.Deleted(), got.Deleted())
	assert.False(t, expected.UpdatedAt().Valid)
	err = expected.Update(context.TODO(), assetPostgresRepo)
	assert.Nil(t, err)
	assert.False(t, expected.Deleted().Valid)
	assert.True(t, expected.UpdatedAt().Valid)
	err = expected.Delete(context.TODO(), assetPostgresRepo)
	assert.True(t, expected.Deleted().Valid)
	assert.True(t, expected.Deleted().Bool)
	assert.Nil(t, err)
	conn, err := prop.DBPool().Acquire(context.TODO())
	defer func() { conn.Release() }()
	assert.Nil(t, err)
	tu, err := entity.GetAsset(context.TODO(), assetPostgresRepo, "bla_bla_bla_test")
	assert.True(t, entity.IsAssetNotFound(tu, err))
}

func testUserCachedPostgresRepo(t *testing.T) {
	prop := env.GetProperties()
	userPostgresRepo := GetUserCachedPostgresRepo(prop)
	cache := userPostgresRepo.(*CachedPostgres[*entity.User])
	upk = tool.RandStringBytes(32)
	user := entity.NewUser(upk, time.Time{})
	err := user.Insert(context.TODO(), userPostgresRepo)
	assert.Nil(t, err)
	assert.False(t, user.Deleted().Valid)
	data1, err := cache.cache.Get(upk)
	assert.Nil(t, err)
	assert.NotNil(t, data1)
	data2, err := user.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, data1, data2)
	got, err := entity.GetUser(context.TODO(), userPostgresRepo, upk)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, user, got)
	assert.Equal(t, user.CreatedAt(), got.CreatedAt())
	assert.Equal(t, user.Upk(), got.Upk())
	assert.Equal(t, user.Deleted(), got.Deleted())
	assert.False(t, user.UpdatedAt().Valid)
	err = user.Update(context.TODO(), userPostgresRepo)
	assert.Nil(t, err)
	assert.False(t, user.Deleted().Valid)
	assert.True(t, user.UpdatedAt().Valid)
	err = user.Delete(context.TODO(), userPostgresRepo)
	assert.True(t, user.Deleted().Valid)
	assert.True(t, user.Deleted().Bool)
	assert.Nil(t, err)
}

func testFavoritesCachedPostgresRepo(t *testing.T) {
	prop := env.GetProperties()
	favoritesPostgresRepo := GetFavoritesCachedPostgresRepo(prop)
	userPostgresRepo := GetUserPostgresRepo(prop)
	user, err := entity.GetUser(context.TODO(), userPostgresRepo, upk)
	assert.Nil(t, err)
	assetPostgresRepo := GetAssetPostgresRepo(prop)
	asset, err := entity.GetAsset(context.TODO(), assetPostgresRepo, isin)
	expected := entity.NewFavorites(asset, user, time.Time{})
	cache := favoritesPostgresRepo.(*CachedPostgres[*entity.Favorites])
	err = expected.Insert(context.TODO(), favoritesPostgresRepo)
	assert.Nil(t, err)
	assert.False(t, expected.Deleted().Valid)

	data1, err := cache.cache.Get(expected.Key())
	assert.Nil(t, err)
	assert.NotNil(t, data1)
	data2, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, data1, data2)

	id = expected.ID()
	got, err := entity.GetFavorites(context.TODO(), favoritesPostgresRepo, isin, upk)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.CreatedAt(), got.CreatedAt())
	assert.Equal(t, expected.ID(), got.ID())
	assert.Equal(t, expected.Deleted(), got.Deleted())
	assert.False(t, expected.UpdatedAt().Valid)
	err = expected.Update(context.TODO(), favoritesPostgresRepo)
	assert.Nil(t, err)
	assert.False(t, expected.Deleted().Valid)
	assert.True(t, expected.UpdatedAt().Valid)
	err = expected.Delete(context.TODO(), favoritesPostgresRepo)
	assert.True(t, expected.Deleted().Valid)
	assert.True(t, expected.Deleted().Bool)
	assert.Nil(t, err)
	tu, err := entity.GetFavorites(context.TODO(), favoritesPostgresRepo, "bla_bla_bla_test", "bla_bla_bla_test")
	assert.True(t, entity.IsFavoritesNotFound(tu, err))
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
