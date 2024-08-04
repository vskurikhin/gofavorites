/*
 * This file was last modified at 2024-07-19 17:50 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites_service_test.go
 * $Id$
 */
//!+

package services

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/vskurikhin/gofavorites/internal/domain/repo"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/mongo"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/local"

	pb "github.com/vskurikhin/gofavorites/proto"
)

func TestFavoritesService(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "test #0 positive Favorites Service GetFavoritesService",
			fRun: testGetFavoritesService,
		},
		{
			name: "test #1 positive Favorites Service ApiFavoritesGet",
			fRun: testFavoritesServiceApiFavoritesGetPositive,
		},
		{
			name: "test #2 positive Favorites Service Get",
			fRun: testFavoritesServiceGetPositive,
		},
		{
			name: "test #3 negative #0 Favorites Service Get",
			fRun: testFavoritesServiceGetNegative0,
		},
		{
			name: "test #4 negative #1 Favorites Service Get",
			fRun: testFavoritesServiceGetNegative1,
		},
		{
			name: "test #5 negative #2 Favorites Service Get",
			fRun: testFavoritesServiceGetNegative2,
		},
		{
			name: "test #6 positive Favorites Service ApiFavoritesGetForUser",
			fRun: testFavoritesServiceApiFavoritesGetForUserPositive,
		},
		{
			name: "test #7 positive #1 Favorites Service GetForUser",
			fRun: testFavoritesServiceGetForUserPositive1,
		},
		{
			name: "test #8 positive #1 Favorites Service GetForUser",
			fRun: testFavoritesServiceGetForUserPositive1,
		},
		{
			name: "test #9 positive #2 Favorites Service GetForUser",
			fRun: testFavoritesServiceGetForUserPositive2,
		},
		{
			name: "test #10 negative #0 Favorites Service GetForUser",
			fRun: testFavoritesServiceGetForUserNegative0,
		},
		{
			name: "test #11 negative #1 Favorites Service GetForUser",
			fRun: testFavoritesServiceGetForUserNegative1,
		},
		{
			name: "test #12 negative #1 Favorites Service GetForUser",
			fRun: testFavoritesServiceGetForUserNegative2,
		},
		{
			name: "test #13 positive Favorites Service ApiFavoritesSet",
			fRun: testFavoritesServiceApiFavoritesSetPositive,
		},
		{
			name: "test #14 positive Favorites Service Set",
			fRun: testFavoritesServiceSetPositive,
		},
		{
			name: "test #15 negative #0 Favorites Service Set",
			fRun: testFavoritesServiceSetNegative0,
		},
		{
			name: "test #16 negative #1 Favorites Service Set",
			fRun: testFavoritesServiceSetNegative1,
		},
		{
			name: "test #17 negative #2 Favorites Service Set",
			fRun: testFavoritesServiceSetNegative2,
		},
		{
			name: "test #18 negative #3 Favorites Service Set",
			fRun: testFavoritesServiceSetNegative3,
		},
		{
			name: "test #19 negative #4 Favorites Service Set",
			fRun: testFavoritesServiceSetNegative4,
		},
		{
			name: "test #20 negative #5 Favorites Service Set",
			fRun: testFavoritesServiceSetNegative5,
		},
		{
			name: "test #21 negative #6 Favorites Service Set",
			fRun: testFavoritesServiceSetNegative6,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testGetFavoritesService(t *testing.T) {
	prop := env.GetProperties()
	got := GetFavoritesService(prop)
	assert.NotNil(t, got)
}

func testFavoritesServiceApiFavoritesGetPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.ApiFavoritesGet(context.TODO(), models.Favorites{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceGetPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Get(context.TODO(), &pb.FavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
	assert.Equal(t, "", resp.GetFavorites().GetAsset().GetIsin())
	assert.Equal(t, "", resp.GetFavorites().GetAsset().GetAssetType().GetName())
	assert.Equal(t, "", resp.GetFavorites().GetUser().GetUpk())
}

func testFavoritesServiceGetNegative0(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, pgx.ErrTxCommitRollback).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Get(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceGetNegative1(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, pgx.ErrTxCommitRollback).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Get(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceGetNegative2(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", tool.ErrEncryptAES).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Get(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceApiFavoritesGetForUserPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	mockMongo.
		EXPECT().
		Load(gomock.Any(), gomock.Any()).
		Return(make([]entity.Favorites, 0), nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(make([]*entity.Favorites, 0), nil).
		AnyTimes()
	syncUtil.
		EXPECT().
		Sync(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(make([]entity.Favorites, 0), nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.ApiFavoritesGetForUser(context.TODO(), models.User{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceGetForUserPositive1(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	mockMongo.
		EXPECT().
		Load(gomock.Any(), gomock.Any()).
		Return(make([]entity.Favorites, 0), nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(make([]*entity.Favorites, 0), nil).
		AnyTimes()
	syncUtil.
		EXPECT().
		Sync(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(make([]entity.Favorites, 0), nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.GetForUser(context.TODO(), &pb.UserFavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
	assert.Equal(t, 0, len(resp.GetFavorites()))
}

func testFavoritesServiceGetForUserPositive2(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	favorites := entity.Favorites{}
	mockMongo.
		EXPECT().
		Load(gomock.Any(), gomock.Any()).
		Return([]entity.Favorites{favorites}, nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]*entity.Favorites{&favorites}, nil).
		AnyTimes()
	syncUtil.
		EXPECT().
		Sync(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]entity.Favorites{favorites}, nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.GetForUser(context.TODO(), &pb.UserFavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
	assert.Equal(t, 1, len(resp.GetFavorites()))
}

func testFavoritesServiceGetForUserNegative0(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(make([]*entity.Favorites, 0), nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", tool.ErrEncryptAES).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.GetForUser(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
	assert.Equal(t, 0, len(resp.GetFavorites()))
}

func testFavoritesServiceGetForUserNegative1(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	mockMongo.
		EXPECT().
		Load(gomock.Any(), gomock.Any()).
		Return(make([]entity.Favorites, 0), repo.ErrNotFound).
		AnyTimes()
	repoFavorites.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(make([]*entity.Favorites, 0), repo.ErrNotFound).
		AnyTimes()
	syncUtil.
		EXPECT().
		Sync(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(make([]entity.Favorites, 0), repo.ErrNotFound).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.GetForUser(context.TODO(), &pb.UserFavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
	assert.Equal(t, 0, len(resp.GetFavorites()))
}

func testFavoritesServiceGetForUserNegative2(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		GetByFilter(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(make([]*entity.Favorites, 0), nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", tool.ErrEncryptAES).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.GetForUser(context.TODO(), &pb.UserFavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
	assert.Equal(t, 0, len(resp.GetFavorites()))
}

func testFavoritesServiceApiFavoritesSetPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.ApiFavoritesSet(context.TODO(), models.Favorites{})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testFavoritesServiceSetPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
	assert.Equal(t, "", resp.GetFavorites().GetAsset().GetIsin())
	assert.Equal(t, "", resp.GetFavorites().GetAsset().GetAssetType().GetName())
	assert.Equal(t, "", resp.GetFavorites().GetUser().GetUpk())
}

func testFavoritesServiceSetNegative0(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).Return(false).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Set(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceSetNegative1(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(false).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceSetNegative2(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(false).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceSetNegative3(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(pgx.ErrTxCommitRollback).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceSetNegative4(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", tool.ErrEncryptAES).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceSetNegative5(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(repo.ErrNotFound).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceSetNegative6(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	mockMongo := NewMockMongo(ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	syncUtil := NewMockSyncUtilService(ctrl)
	upkUtil := NewMockUpkUtilService(ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	repoFavorites.
		EXPECT().
		Update(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&entity.Favorites{}, repo.ErrNotFound).
		AnyTimes()
	upkUtil.
		EXPECT().
		EncryptPersonalKey(gomock.Any()).
		Return("", nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, mockMongo, repoFavorites, syncUtil, upkUtil, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
}

func getTestFavoritesService(
	assetLookup AssetSearchService,
	dftFavorites domain.Dft[*entity.Favorites],
	mongo mongo.Mongo,
	repoFavorites domain.Repo[*entity.Favorites],
	syncService SyncUtilService,
	upkUtil UpkUtilService,
	userLookup UserSearchService,
) FavoritesService {
	favoritesServ = new(favoritesService)
	favoritesServ.assetLookup = assetLookup
	favoritesServ.dftFavorites = dftFavorites
	favoritesServ.mongo = mongo
	favoritesServ.repoFavorites = repoFavorites
	favoritesServ.upkUtil = upkUtil
	favoritesServ.userLookup = userLookup
	favoritesServ.sLog = slog.Default()
	favoritesServ.syncService = syncService
	return favoritesServ
}

func TestGRPCFavoritesService(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #1 Favorites Service Get",
			fRun: testGRPCFavoritesServiceGetPositive,
		},
		{
			name: "positive test #2 Favorites Service GetForUser",
			fRun: testGRPCFavoritesServiceGetForUserPositive,
		},
		{
			name: "positive test #3 Favorites Service Set",
			fRun: testGRPCFavoritesServiceSetPositive,
		},
		{
			name: "negative test #4 Favorites Service Get",
			fRun: testGRPCFavoritesServiceGetNegative,
		},
		{
			name: "negative test #5 Favorites Service Get",
			fRun: testGRPCFavoritesServiceGetForUserNegative,
		},
		{
			name: "negative test #6 Favorites Service Get",
			fRun: testGRPCFavoritesServiceSetNegative,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testGRPCFavoritesServiceGetPositive(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65285+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx, address, favoritesServicePositive{}, up)
	<-up

	conn, client, err := makeFavoritesServiceClient(t, address)
	defer func() { _ = conn.Close() }()

	var request pb.FavoritesRequest
	resp, err := client.Get(ctx, &request)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testGRPCFavoritesServiceGetForUserPositive(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65321+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx, address, favoritesServicePositive{}, up)
	<-up

	conn, client, err := makeFavoritesServiceClient(t, address)
	defer func() { _ = conn.Close() }()

	var request pb.UserFavoritesRequest
	resp, err := client.GetForUser(ctx, &request)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testGRPCFavoritesServiceSetPositive(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65357+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx, address, favoritesServicePositive{}, up)
	<-up

	conn, client, err := makeFavoritesServiceClient(t, address)
	defer func() { _ = conn.Close() }()

	var request pb.FavoritesRequest
	resp, err := client.Set(ctx, &request)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testGRPCFavoritesServiceGetNegative(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65393+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx, address, favoritesServiceNegative{}, up)
	<-up

	conn, client, err := makeFavoritesServiceClient(t, address)
	defer func() { _ = conn.Close() }()

	var request pb.FavoritesRequest
	resp, err := client.Get(ctx, &request)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func testGRPCFavoritesServiceGetForUserNegative(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65429+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx, address, favoritesServiceNegative{}, up)
	<-up

	conn, client, err := makeFavoritesServiceClient(t, address)
	defer func() { _ = conn.Close() }()

	var request pb.UserFavoritesRequest
	resp, err := client.GetForUser(ctx, &request)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func testGRPCFavoritesServiceSetNegative(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65465+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer func() {
		cancel()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx, address, favoritesServiceNegative{}, up)
	<-up

	conn, client, err := makeFavoritesServiceClient(t, address)
	defer func() { _ = conn.Close() }()

	var request pb.FavoritesRequest
	resp, err := client.Set(ctx, &request)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func makeFavoritesServiceClient(t *testing.T, address string) (*grpc.ClientConn, pb.FavoritesServiceClient, error) {

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		t.Fail()
	}
	client := pb.NewFavoritesServiceClient(conn)

	return conn, client, err
}

func grpcServeFavoritesServiceServer(ctx context.Context, address string, srv pb.FavoritesServiceServer, up chan struct{}) {
	listen, err := net.Listen("tcp", address)
	tool.IfErrorThenPanic(err)
	opts := []grpc.ServerOption{grpc.Creds(local.NewCredentials())}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterFavoritesServiceServer(grpcServer, srv)
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()
	if up != nil {
		close(up)
	}
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

type FavoritesServicePositive interface {
	pb.FavoritesServiceServer
}

type favoritesServicePositive struct {
	pb.UnimplementedFavoritesServiceServer
}

var _ FavoritesServicePositive = (*favoritesServicePositive)(nil)

func (f favoritesServicePositive) Get(_ context.Context, _ *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {
	return &pb.FavoritesResponse{}, nil
}

func (f favoritesServicePositive) GetForUser(_ context.Context, _ *pb.UserFavoritesRequest) (*pb.UserFavoritesResponse, error) {
	return &pb.UserFavoritesResponse{}, nil
}

func (f favoritesServicePositive) Set(_ context.Context, _ *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {
	return &pb.FavoritesResponse{}, nil
}

type FavoritesServiceNegative interface {
	pb.FavoritesServiceServer
}

type favoritesServiceNegative struct {
	pb.UnimplementedFavoritesServiceServer
}

var _ FavoritesServiceNegative = (*favoritesServiceNegative)(nil)

func (f favoritesServiceNegative) Get(_ context.Context, request *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {
	return nil, fmt.Errorf("test")
}

func (f favoritesServiceNegative) GetForUser(_ context.Context, request *pb.UserFavoritesRequest) (*pb.UserFavoritesResponse, error) {
	return nil, fmt.Errorf("test")
}

func (f favoritesServiceNegative) Set(_ context.Context, request *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {
	return nil, fmt.Errorf("test")
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
