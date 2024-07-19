/*
 * This file was last modified at 2024-07-15 17:10 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * asset_test.go
 * $Id$
 */
//!+

// Package entity TODO.
package entity

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"testing"
	"time"
)

func TestAsset(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{name: "positive test #0 Asset Cloneable", fRun: testAssetCloneable},
		{name: "positive test #1 Asset FromJSON and ToJSON", fRun: testAssetJSON},
		{name: "positive test #2 Asset IsAssetNotFound", fRun: testIsAssetNotFound},
		{name: "positive test #3 Asset stubRepoOk", fRun: testAssetRepoOk},
		{name: "negative test #4 Asset stubRepoErr", fRun: testAssetRepoErr},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testAssetCloneable(t *testing.T) {
	expected := MakeAsset(
		tool.RandStringBytes(32),
		AssetType{},
		MakeTAttributes(
			sql.NullBool{true, true},
			time.Time{},
			sql.NullTime{time.Time{}, false},
		))
	got := expected.Copy()
	assert.NotNil(t, got)
	assert.Equal(t, &expected, got)
}

func testAssetJSON(t *testing.T) {
	isin := tool.RandStringBytes(32)
	expected := MakeAsset(
		isin,
		AssetType{},
		MakeTAttributes(
			sql.NullBool{true, true},
			time.Time{},
			sql.NullTime{time.Time{}, false},
		))
	j, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.NotNil(t, j)
	got := Asset{}
	err = (&got).FromJSON(j)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.String(), got.String())
	assert.Equal(t, isin, got.Key())
}

func testIsAssetNotFound(t *testing.T) {
	assert.True(t, IsAssetNotFound(Asset{isin: "test"}, pgx.ErrNoRows))
	assert.True(t, IsAssetNotFound(Asset{}, errors.New("")))
	assert.False(t, IsAssetNotFound(Asset{isin: "test"}, errors.New("")))
}

func testAssetRepoOk(t *testing.T) {
	asset := MakeAsset("", AssetType{}, MakeTAttributes(sql.NullBool{}, time.Time{}, sql.NullTime{}))
	err := asset.Upsert(context.TODO(), &stubTxRepoOk[*Asset]{}, func() {})
	assert.Nil(t, err)
	assert.False(t, asset.Deleted().Valid)
	got, err := GetAsset(context.TODO(), &stubRepoOk[*Asset]{}, "")
	assert.Nil(t, err)
	assert.Equal(t, asset, got)
	assert.Equal(t, asset.CreatedAt(), got.CreatedAt())
	assert.Equal(t, asset.Isin(), got.Isin())
	assert.Equal(t, asset.Deleted(), got.Deleted())
	assert.False(t, asset.UpdatedAt().Valid)
	err = asset.Delete(context.TODO(), &stubTxRepoOk[*Asset]{}, func() {})
	assert.Nil(t, err)
	assert.False(t, asset.Deleted().Valid)
}

func testAssetRepoErr(t *testing.T) {
	isin := tool.RandStringBytes(32)
	asset := MakeAsset(isin, AssetType{}, MakeTAttributes(sql.NullBool{}, time.Time{}, sql.NullTime{}))
	err := asset.Upsert(context.TODO(), &stubTxRepoErr[*Asset]{}, func() {})
	assert.NotNil(t, err)
	_, err = GetAsset(context.TODO(), &stubRepoErr[*Asset]{}, isin)
	assert.NotNil(t, err)
	err = asset.Delete(context.TODO(), &stubTxRepoErr[*Asset]{}, func() {})
	assert.NotNil(t, err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
