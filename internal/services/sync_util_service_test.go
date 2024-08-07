/*
 * This file was last modified at 2024-08-03 17:39 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * sync_util_service.go
 * $Id$
 */
//!+

// Package services сервисы бизнес логики.
package services

import (
	"context"
	"database/sql"
	"log/slog"
	"testing"

	"github.com/vskurikhin/gofavorites/internal/domain/repo"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/batch"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/mongo"
	"github.com/vskurikhin/gofavorites/internal/env"
	"go.uber.org/mock/gomock"
)

func TestSyncUtilService(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "test #0 positive Sync util Service GetSyncUtilService",
			fRun: testGetSyncUtilService,
		},
		{
			name: "test #1 positive #1 Sync util Service Sync",
			fRun: testFavoritesServiceSyncPositive1,
		},
		{
			name: "test #2 positive #2 Sync util Service Sync",
			fRun: testFavoritesServiceSyncPositive2,
		},
		{
			name: "test #3 positive #3 Sync util Service Sync",
			fRun: testFavoritesServiceSyncPositive3,
		},
		{
			name: "test #4 positive #4 Sync util Service Sync",
			fRun: testFavoritesServiceSyncPositive4,
		},
		{
			name: "test #5 positive #5 Sync util Service Sync",
			fRun: testFavoritesServiceSyncPositive5,
		},
		{
			name: "test #6 positive #6 Sync util Service Sync",
			fRun: testFavoritesServiceSyncPositive6,
		},
		{
			name: "test #7 positive #7 Sync util Service Sync",
			fRun: testFavoritesServiceSyncPositive7,
		},
		{
			name: "test #8 positive #8 Sync util Service Sync",
			fRun: testFavoritesServiceSyncPositive8,
		},
		{
			name: "test #9 negative #1 Sync util Service Sync",
			fRun: testFavoritesServiceSyncNegative1,
		},
		{
			name: "test #10 negative #2 Sync util Service Sync",
			fRun: testFavoritesServiceSyncNegative2,
		},
	}
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testFavoritesServiceSyncNegative1(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	mockMongo.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(repo.ErrNotFound).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, repo.ErrBadPool).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, repo.ErrBadPool).
		AnyTimes()
	repoFavoritesDeleted.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*entity.FavoritesDeleted{{}}, repo.ErrBadPool).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	favorites1 := entity.Favorites{}
	favorites2 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 2, entity.DefaultTAttributes()),
		sql.NullInt64{},
		entity.DefaultTAttributes(),
	)
	favorites3 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 3, entity.DefaultTAttributes()),
		sql.NullInt64{},
		entity.DefaultTAttributes(),
	)
	favorites0 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 0, entity.DefaultTAttributes()),
		sql.NullInt64{},
		entity.DefaultTAttributes(),
	)
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{favorites1}, []entity.Favorites{favorites2, favorites3, favorites0})
	assert.Nil(t, err)
	assert.NotNil(t, resp)

}

func testFavoritesServiceSyncNegative2(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesInsertsBatch.
		EXPECT().
		Do(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, repo.ErrBadPool).
		AnyTimes()
	repoFavoritesDeleted.
		EXPECT().
		Delete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.FavoritesDeleted{}, repo.ErrBadPool).
		AnyTimes()
	repoUser.
		EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.User{}, repo.ErrBadPool).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	favorites0 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 0, entity.DefaultTAttributes()),
		sql.NullInt64{Int64: 0, Valid: true},
		entity.DefaultTAttributes(),
	)
	favorites1 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 1, entity.DefaultTAttributes()),
		sql.NullInt64{Int64: 1, Valid: true},
		entity.DefaultTAttributes(),
	)
	favorites2 := entity.Favorites{}
	favorites3 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 3, entity.DefaultTAttributes()),
		sql.NullInt64{Int64: 1, Valid: true},
		entity.DefaultTAttributes(),
	)
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{favorites0, favorites3, favorites1}, []entity.Favorites{favorites2})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testGetSyncUtilService(t *testing.T) {
	prop := env.GetProperties()
	got := GetSyncUtilService(prop)
	assert.NotNil(t, got)
}

func testFavoritesServiceSyncPositive1(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{}, []entity.Favorites{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceSyncPositive2(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	favorites := entity.Favorites{}
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{}, []entity.Favorites{favorites, favorites})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceSyncPositive3(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	userLookup.EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	favorites := entity.Favorites{}
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{favorites, favorites}, []entity.Favorites{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceSyncPositive4(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	favorites := entity.Favorites{}
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{favorites}, []entity.Favorites{favorites})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceSyncPositive5(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesInsertsBatch.
		EXPECT().
		Do(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	repoFavoritesDeleted.
		EXPECT().
		Delete(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.FavoritesDeleted{}, nil).
		AnyTimes()
	repoUser.
		EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.User{}, nil).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	favorites0 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 0, entity.DefaultTAttributes()),
		sql.NullInt64{Int64: 0, Valid: true},
		entity.DefaultTAttributes(),
	)
	favorites1 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 1, entity.DefaultTAttributes()),
		sql.NullInt64{Int64: 1, Valid: true},
		entity.DefaultTAttributes(),
	)
	favorites2 := entity.Favorites{}
	favorites3 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 3, entity.DefaultTAttributes()),
		sql.NullInt64{Int64: 1, Valid: true},
		entity.DefaultTAttributes(),
	)
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{favorites0, favorites3, favorites1}, []entity.Favorites{favorites2})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceSyncPositive6(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	mockMongo.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).Return(nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	repoFavoritesDeleted.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*entity.FavoritesDeleted{{}}, nil).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	favorites1 := entity.Favorites{}
	favorites2 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 1, entity.DefaultTAttributes()),
		sql.NullInt64{Int64: 1, Valid: true},
		entity.DefaultTAttributes(),
	)
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{favorites1}, []entity.Favorites{favorites2})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceSyncPositive7(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	mockMongo.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).Return(nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	repoFavoritesDeleted.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*entity.FavoritesDeleted{{}}, nil).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	favorites1 := entity.Favorites{}
	favorites2 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 2, entity.DefaultTAttributes()),
		sql.NullInt64{},
		entity.DefaultTAttributes(),
	)
	favorites3 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 3, entity.DefaultTAttributes()),
		sql.NullInt64{},
		entity.DefaultTAttributes(),
	)
	favorites0 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 0, entity.DefaultTAttributes()),
		sql.NullInt64{},
		entity.DefaultTAttributes(),
	)
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{favorites1}, []entity.Favorites{favorites2, favorites3, favorites0})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceSyncPositive8(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	favoritesInsertsBatch := NewMockFavoritesInsertsBatch(ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	repoFavoritesDeleted := NewMockRepo[*entity.FavoritesDeleted](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoUser := NewMockRepo[*entity.User](ctrl)
	mockMongo.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).Return(nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	fg := entity.FavoritesDeleted{}
	fdl := make([]*entity.FavoritesDeleted, 0)
	fdl = append(fdl, &fg)
	repoFavoritesDeleted.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(fdl, nil).
		AnyTimes()
	syncUtilService := getTestSyncUtilService(
		assetLookup,
		favoritesInsertsBatch,
		mockMongo,
		repoFavorites,
		repoFavoritesDeleted,
		userLookup,
		repoUser,
	)
	favorites1 := entity.Favorites{}
	favorites2 := entity.MakeFavorites(
		uuid.New(),
		entity.Asset{},
		entity.MakeUserWithVersion("", 1, entity.DefaultTAttributes()),
		sql.NullInt64{Int64: 1, Valid: true},
		entity.DefaultTAttributes(),
	)
	resp, err := syncUtilService.Sync(context.TODO(), []entity.Favorites{favorites1}, []entity.Favorites{favorites2})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func getTestSyncUtilService(
	assetLookup AssetSearchService,
	batch batch.FavoritesInsertsBatch,
	mongo mongo.Mongo,
	repoFavorites domain.Repo[*entity.Favorites],
	repoFavoritesDeleted domain.Repo[*entity.FavoritesDeleted],
	userLookup UserSearchService,
	userRepo domain.Repo[*entity.User],
) SyncUtilService {
	syncUtilServ = new(syncUtilService)
	syncUtilServ.assetLookup = assetLookup
	syncUtilServ.batch = batch
	syncUtilServ.mongo = mongo
	syncUtilServ.repoFavorites = repoFavorites
	syncUtilServ.repoFavoritesDeleted = repoFavoritesDeleted
	syncUtilServ.sLog = slog.Default()
	syncUtilServ.userLookup = userLookup
	syncUtilServ.userRepo = userRepo
	return syncUtilServ
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
