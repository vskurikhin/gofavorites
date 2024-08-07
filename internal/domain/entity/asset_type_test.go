/*
 * This file was last modified at 2024-07-15 17:10 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * asset_type_test.go
 * $Id$
 */
//!+

// Package entity TODO.
package entity

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/tool"
)

func TestAssetType(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{name: "positive test #0 AssetType Cloneable", fRun: testAssetTypeCloneable},
		{name: "positive test #1 AssetType FromJSON and ToJSON", fRun: testAssetTypeJSON},
		{name: "positive test #2 AssetType stubRepoOk", fRun: testAssetTypeRepoOk},
		{name: "negative test #3 AssetType stubRepoErr", fRun: testAssetTypeRepoErr},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testAssetTypeCloneable(t *testing.T) {
	ta := MakeTAttributes(sql.NullBool{Bool: true, Valid: true}, time.Time{}, sql.NullTime{})
	expected := MakeAssetType(tool.RandStringBytes(32), ta)
	got := expected.Copy()
	assert.NotNil(t, got)
	assert.Equal(t, &expected, got)
}

func testAssetTypeJSON(t *testing.T) {
	name := tool.RandStringBytes(32)
	ta := MakeTAttributes(sql.NullBool{Bool: true, Valid: true}, time.Time{}, sql.NullTime{})
	expected := MakeAssetType(name, ta)
	j, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.NotNil(t, j)
	got := AssetType{}
	err = (&got).FromJSON(j)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.String(), got.String())
	assert.Equal(t, name, got.Key())
}

func testAssetTypeRepoOk(t *testing.T) {
	name := tool.RandStringBytes(32)
	assetType := MakeAssetType(name, DefaultTAttributes())
	err := assetType.Insert(context.TODO(), &stubRepoOk[*AssetType]{})
	assert.Nil(t, err)
	assert.False(t, assetType.Deleted().Valid)
	got, err := GetAssetType(context.TODO(), &stubRepoOk[*AssetType]{}, name)
	assert.Nil(t, err)
	assert.Equal(t, assetType, got)
	assert.Equal(t, assetType.CreatedAt(), got.CreatedAt())
	assert.Equal(t, assetType.Name(), got.Name())
	assert.Equal(t, assetType.Deleted(), got.Deleted())
	assert.False(t, assetType.UpdatedAt().Valid)
	err = assetType.Update(context.TODO(), &stubRepoOk[*AssetType]{})
	assert.Nil(t, err)
	assert.False(t, assetType.Deleted().Valid)
	assert.False(t, assetType.UpdatedAt().Valid)
	err = assetType.Delete(context.TODO(), &stubRepoOk[*AssetType]{})
	assert.Nil(t, err)
	assert.False(t, assetType.Deleted().Valid)
}

func testAssetTypeRepoErr(t *testing.T) {
	name := tool.RandStringBytes(32)
	assetType := MakeAssetType(name, DefaultTAttributes())
	err := assetType.Insert(context.TODO(), &stubRepoErr[*AssetType]{})
	assert.NotNil(t, err)
	_, err = GetAssetType(context.TODO(), &stubRepoErr[*AssetType]{}, name)
	assert.NotNil(t, err)
	err = assetType.Update(context.TODO(), &stubRepoErr[*AssetType]{})
	assert.NotNil(t, err)
	err = assetType.Delete(context.TODO(), &stubRepoErr[*AssetType]{})
	assert.NotNil(t, err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
