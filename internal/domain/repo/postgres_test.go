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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"testing"
	"time"
)

func TestPostgresRepos(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #0 AssetType Postgres Repo",
			fRun: testAssetTypePostgresRepo,
		},
		{
			name: "positive test #1 Asset Postgres Repo",
			fRun: testAssetPostgresRepo,
		},
		{
			name: "positive test #2 User Postgres Repo",
			fRun: testUserPostgresRepo,
		},
		{
			name: "positive test #3 Favorites Postgres Repo",
			fRun: testFavoritesPostgresRepo,
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

var (
	assetType string
	id        uuid.UUID
	isin      string
	upk       string
)

func testAssetTypePostgresRepo(t *testing.T) {
	prop := env.GetProperties()
	assetTypePostgresRepo := GetAssetTypePostgresRepo(prop)
	assetType = tool.RandStringBytes(32)
	expected := entity.NewAssetType(assetType, time.Time{})
	err := expected.Insert(context.TODO(), assetTypePostgresRepo)
	assert.Nil(t, err)
	assert.False(t, expected.Deleted().Valid)
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

func testAssetPostgresRepo(t *testing.T) {
	prop := env.GetProperties()
	assetTypePostgresRepo := GetAssetTypePostgresRepo(prop)
	assetPostgresRepo := GetAssetPostgresRepo(prop)
	at, err := entity.GetAssetType(context.TODO(), assetTypePostgresRepo, assetType)
	assert.Nil(t, err)
	isin = tool.RandStringBytes(32)
	expected := entity.NewAsset(isin, at, time.Time{})
	err = expected.Insert(context.TODO(), assetPostgresRepo)
	assert.Nil(t, err)
	assert.False(t, expected.Deleted().Valid)
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

func testUserPostgresRepo(t *testing.T) {
	prop := env.GetProperties()
	userPostgresRepo := GetUserPostgresRepo(prop)
	upk = tool.RandStringBytes(32)
	user := entity.NewUser(upk, time.Time{})
	err := user.Insert(context.TODO(), userPostgresRepo)
	assert.Nil(t, err)
	assert.False(t, user.Deleted().Valid)
	got, err := entity.GetUser(context.TODO(), userPostgresRepo, upk)
	assert.Nil(t, err)
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
	tu, err := entity.GetUser(context.TODO(), userPostgresRepo, "bla_bla_bla_test")
	assert.True(t, entity.IsUserNotFound(tu, err))
}

func testFavoritesPostgresRepo(t *testing.T) {
	prop := env.GetProperties()
	favoritesPostgresRepo := GetFavoritesPostgresRepo(prop)
	userPostgresRepo := GetUserPostgresRepo(prop)
	user, err := entity.GetUser(context.TODO(), userPostgresRepo, upk)
	assert.Nil(t, err)
	assetPostgresRepo := GetAssetPostgresRepo(prop)
	asset, err := entity.GetAsset(context.TODO(), assetPostgresRepo, isin)
	expected := entity.NewFavorites(asset, user, time.Time{})
	err = expected.Insert(context.TODO(), favoritesPostgresRepo)
	assert.Nil(t, err)
	assert.False(t, expected.Deleted().Valid)
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

func testClear(t *testing.T) {
	prop := env.GetProperties()
	conn, err := prop.DBPool().Acquire(context.TODO())
	defer func() { conn.Release() }()
	assert.Nil(t, err)

	row := conn.QueryRow(context.TODO(), "DELETE FROM favorites WHERE id = $1 RETURNING id", id)
	assert.NotNil(t, row)
	var deletedID uuid.UUID
	err = row.Scan(&deletedID)
	assert.Nil(t, err)
	assert.Equal(t, id, deletedID)

	row = conn.QueryRow(context.TODO(), "DELETE FROM assets WHERE isin = $1 RETURNING isin", isin)
	assert.NotNil(t, row)
	var deletedIsin string
	err = row.Scan(&deletedIsin)
	assert.Nil(t, err)
	assert.Equal(t, isin, deletedIsin)

	row = conn.QueryRow(context.TODO(), "DELETE FROM users WHERE upk = $1 RETURNING upk", upk)
	assert.NotNil(t, row)
	var deletedUPK string
	err = row.Scan(&deletedUPK)
	assert.Nil(t, err)
	assert.Equal(t, upk, deletedUPK)

	row = conn.QueryRow(context.TODO(), "DELETE FROM asset_types WHERE name = $1 RETURNING name", assetType)
	assert.NotNil(t, row)
	var deletedName string
	err = row.Scan(&deletedName)
	assert.Nil(t, err)
	assert.Equal(t, assetType, deletedName)
}

/*

 */
//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
