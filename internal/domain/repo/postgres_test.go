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
)

var (
	assetType string
	id        uuid.UUID
	isin      string
	upk       string
)

func TestPostgresRepos(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "negative test #1 AssetType Postgres Repo",
			fRun: testAssetTypePostgresRepoNegative,
		},
		{
			name: "negative test #2 User Postgres Repo",
			fRun: testUserPostgresRepoNegative,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testUserPostgresRepoNegative(t *testing.T) {
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")
	prop := env.GetProperties()
	userPostgresRepo := GetUserPostgresRepo(prop)
	upk = tool.RandStringBytes(32)
	user := entity.MakeUser(upk, entity.DefaultTAttributes())
	err := user.Insert(context.TODO(), userPostgresRepo)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)
	err = user.Update(context.TODO(), userPostgresRepo)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)
	err = user.Delete(context.TODO(), userPostgresRepo)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)
	_ = GetUserPostgresRepo(prop)
}

func testAssetTypePostgresRepoNegative(t *testing.T) {
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")
	prop := env.GetProperties()
	assetTypePostgresRepo := GetAssetTypePostgresRepo(prop)
	assetType = tool.RandStringBytes(32)
	expected := entity.MakeAssetType(assetType, entity.DefaultTAttributes())
	err := expected.Insert(context.TODO(), assetTypePostgresRepo)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)

	assert.False(t, expected.Deleted().Valid)
	_, err = entity.GetAssetType(context.TODO(), assetTypePostgresRepo, assetType)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)

	err = expected.Update(context.TODO(), assetTypePostgresRepo)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)

	err = expected.Delete(context.TODO(), assetTypePostgresRepo)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)
}

func testAssetTypePostgresRepoPositive(t *testing.T) {
	defer func() { _ = recover() }()
	prop := env.GetProperties()
	assetTypePostgresRepo := GetAssetTypePostgresRepo(prop)
	assetType = tool.RandStringBytes(32)
	expected := entity.MakeAssetType(assetType, entity.DefaultTAttributes())
	err := expected.Insert(context.TODO(), assetTypePostgresRepo)
	assert.Nil(t, err)
	defer testClearAssetTypes(t)
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

func testUserPostgresRepoPositive(t *testing.T) {
	defer func() { _ = recover() }()
	prop := env.GetProperties()
	userPostgresRepo := GetUserPostgresRepo(prop)
	upk = tool.RandStringBytes(32)
	user := entity.MakeUser(upk, entity.DefaultTAttributes())
	err := user.Insert(context.TODO(), userPostgresRepo)
	assert.Nil(t, err)
	defer testClearUsers(t)
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

func testClearAssetTypes(t *testing.T) {
	defer func() { _ = recover() }()
	prop := env.GetProperties()
	conn, err := prop.DBPool().Acquire(context.TODO())
	defer func() { conn.Release() }()
	assert.Nil(t, err)

	row := conn.QueryRow(context.TODO(), "DELETE FROM asset_types WHERE name = $1 RETURNING name", assetType)
	assert.NotNil(t, row)
	var deletedName string
	err = row.Scan(&deletedName)
	assert.Nil(t, err)
	assert.Equal(t, assetType, deletedName)
}

func testClearAssets(t *testing.T) {
	defer func() { _ = recover() }()
	prop := env.GetProperties()
	conn, err := prop.DBPool().Acquire(context.TODO())
	defer func() { conn.Release() }()
	assert.Nil(t, err)

	row := conn.QueryRow(context.TODO(), "DELETE FROM assets WHERE isin = $1 RETURNING isin", isin)
	assert.NotNil(t, row)
	var deletedIsin string
	err = row.Scan(&deletedIsin)
	assert.Nil(t, err)
	assert.Equal(t, isin, deletedIsin)
}

func testClearUsers(t *testing.T) {
	defer func() { _ = recover() }()
	prop := env.GetProperties()
	conn, err := prop.DBPool().Acquire(context.TODO())
	defer func() { conn.Release() }()
	assert.Nil(t, err)

	row := conn.QueryRow(context.TODO(), "DELETE FROM users WHERE upk = $1 RETURNING upk", upk)
	assert.NotNil(t, row)
	var deletedUPK string
	err = row.Scan(&deletedUPK)
	assert.Nil(t, err)
	assert.Equal(t, upk, deletedUPK)
}

func testClearFavorites(t *testing.T) {
	defer func() { _ = recover() }()
	prop := env.GetProperties()
	conn, err := prop.DBPool().Acquire(context.TODO())
	defer func() { conn.Release() }()
	assert.Nil(t, err)

	row := conn.QueryRow(context.TODO(), "DELETE FROM favorites WHERE isin = $1 AND user_upk = $2 RETURNING isin, user_upk", isin, upk)
	assert.NotNil(t, row)
	var deletedIsin, deletedUPK string
	err = row.Scan(&deletedIsin, &deletedUPK)
	assert.Nil(t, err)
	assert.Equal(t, isin, deletedIsin)
	assert.Equal(t, upk, deletedUPK)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
