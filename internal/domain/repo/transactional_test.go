package repo

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
)

func TestTxPostgres(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "negative test #1 Asset TxPostgres Repo",
			fRun: testAssetTxPostgresNegative,
		},
		{
			name: "negative test #3 Favorites TxPostgres Repo",
			fRun: testFavoritesTxPostgresNegative,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testAssetTxPostgresNegative(t *testing.T) {
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")

	prop := env.GetProperties()
	assetPostgresDft := GetAssetTxPostgres(prop)

	assetType = "stocks"
	isin = tool.RandStringBytes(32)
	at := entity.MakeAssetType(assetType, entity.DefaultTAttributes())
	asset := entity.MakeAsset(isin, at, entity.DefaultTAttributes())
	// TODO var ok bool
	err := asset.Upsert(context.TODO(), assetPostgresDft, func() {
		// TODO ok = true
	})
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)
	// TODO assert.True(t, ok)
}

func testAssetTxPostgresPositive(t *testing.T) {
	defer func() { _ = recover() }() // TODO

	prop := env.GetProperties()
	assetPostgresDft := GetAssetTxPostgres(prop)

	assetType = "stocks"
	isin = tool.RandStringBytes(32)
	at := entity.MakeAssetType(assetType, entity.DefaultTAttributes())
	asset := entity.MakeAsset(isin, at, entity.DefaultTAttributes())
	// TODO var ok bool
	_ = asset.Upsert(context.TODO(), assetPostgresDft, func() {
		// TODO ok = true
	})
	//assert.Nil(t, err) TODO
	// TODO assert.True(t, ok)

	defer testClearAssetTypes(t)
	defer testClearAssets(t)
}

func testFavoritesTxPostgresNegative(t *testing.T) {
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")

	prop := env.GetProperties()
	dft := GetFavoritesTxPostgres(prop)
	repo := getFavoritesCachedPostgresRepo(prop)

	assetType = "stocks"
	id = uuid.New()
	isin = tool.RandStringBytes(16)
	upk = tool.RandStringBytes(32)

	at := entity.MakeAssetType(assetType, entity.DefaultTAttributes())
	asset := entity.MakeAsset(isin, at, entity.DefaultTAttributes())
	user := entity.MakeUser(upk, entity.DefaultTAttributes())

	expected := entity.MakeFavorites(id, asset, user, sql.NullInt64{}, entity.DefaultTAttributes())

	var ok bool
	inTransaction := func() {
		ok = !ok
	}
	err := expected.Upsert(context.TODO(), dft, inTransaction)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)

	_, err = entity.GetFavorites(context.TODO(), repo, isin, upk)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)

	_, err = entity.GetFavoritesForUser(context.TODO(), repo, upk)
	assert.Nil(t, err)

	err = expected.Delete(context.TODO(), dft, inTransaction)
	assert.NotNil(t, err)
	assert.Equal(t, ErrBadPool, err)
}

func testFavoritesTxPostgresPositive(t *testing.T) {
	defer func() { _ = recover() }()

	prop := env.GetProperties()
	dft := GetFavoritesTxPostgres(prop)
	repo := getFavoritesCachedPostgresRepo(prop)
	txPostgres := dft.(*TxPostgres[*entity.Favorites])
	cache := txPostgres.cache.(*CachedPostgres[*entity.Favorites])

	assetType = "stocks"
	id = uuid.New()
	isin = tool.RandStringBytes(16)
	upk = tool.RandStringBytes(32)

	at := entity.MakeAssetType(assetType, entity.DefaultTAttributes())
	asset := entity.MakeAsset(isin, at, entity.DefaultTAttributes())
	user := entity.MakeUser(upk, entity.DefaultTAttributes())

	expected := entity.MakeFavorites(id, asset, user, sql.NullInt64{}, entity.DefaultTAttributes())

	var ok bool
	inTransaction := func() {
		ok = !ok
	}
	err := expected.Upsert(context.TODO(), dft, inTransaction)
	assert.Nil(t, err)
	assert.True(t, ok)

	defer testClearAssetTypes(t)
	defer testClearAssets(t)
	defer testClearUsers(t)
	defer testClearFavorites(t)

	assert.False(t, expected.Deleted().Valid)
	assert.NotEqual(t, expected.CreatedAt(), sql.NullTime{})
	assert.False(t, expected.UpdatedAt().Valid)

	asset = expected.Asset()
	assert.False(t, asset.Deleted().Valid)
	assert.NotEqual(t, asset.CreatedAt(), sql.NullTime{})
	assert.False(t, asset.UpdatedAt().Valid)

	at = asset.AssetType()
	assert.False(t, at.Deleted().Valid)
	assert.NotEqual(t, at.CreatedAt(), sql.NullTime{})
	assert.False(t, at.UpdatedAt().Valid)

	user = expected.User()
	assert.False(t, user.Deleted().Valid)
	assert.NotEqual(t, user.CreatedAt(), sql.NullTime{})
	assert.False(t, user.UpdatedAt().Valid)

	data1, err := cache.cache.Get(expected.Key())
	assert.Nil(t, err)
	assert.NotNil(t, data1)
	data2, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, data1, data2)

	got, err := entity.GetFavorites(context.TODO(), repo, isin, upk)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)

	list, err := entity.GetFavoritesForUser(context.TODO(), repo, upk)
	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, expected, list[0])

	err = expected.Delete(context.TODO(), dft, inTransaction)
	assert.Nil(t, err)
	assert.False(t, ok)
	d, err := cache.cache.Get(expected.Key())
	assert.Nil(t, d)
	assert.True(t, expected.Deleted().Bool)
	assert.True(t, expected.Deleted().Valid)

}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
