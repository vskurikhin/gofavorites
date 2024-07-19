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
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	pb "github.com/vskurikhin/gofavorites/proto"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFavoritesService(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #1 Favorites Service Get",
			fRun: testFavoritesServiceGetPositive,
		},
		{
			name: "negative test #2 Favorites Service Get #1",
			fRun: testFavoritesServiceGetNegative1,
		},
		{
			name: "positive test #3 Favorites Service GetForUser",
			fRun: testFavoritesServiceGetForUserPositive1,
		},
		{
			name: "positive test #4 Favorites Service Set",
			fRun: testFavoritesServiceSetPositive,
		},
		{
			name: "negative test #5 Favorites Service Set #1",
			fRun: testFavoritesServiceSetNegative1,
		},
		{
			name: "negative test #6 Favorites Service Set #2",
			fRun: testFavoritesServiceSetNegative2,
		},
		{
			name: "negative test #7 Favorites Service Set #3",
			fRun: testFavoritesServiceSetNegative3,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testFavoritesServiceGetPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		Get(context.TODO(), gomock.Any(), gomock.Any()).Return(&entity.Favorites{}, nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, repoFavorites, userLookup)
	resp, err := favoritesService.Get(context.TODO(), &pb.FavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
	assert.Equal(t, "", resp.GetFavorites().GetAsset().GetIsin())
	assert.Equal(t, "", resp.GetFavorites().GetAsset().GetAssetType().GetName())
	assert.Equal(t, "", resp.GetFavorites().GetUser().GetUpk())
}

func testFavoritesServiceGetNegative1(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		Get(context.TODO(), gomock.Any(), gomock.Any()).Return(nil, pgx.ErrTxCommitRollback).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, repoFavorites, userLookup)
	resp, err := favoritesService.Get(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceGetForUserPositive1(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	repoFavorites.
		EXPECT().
		GetByFilter(context.TODO(), gomock.Any(), gomock.Any()).
		Return(make([]*entity.Favorites, 0), nil).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, repoFavorites, userLookup)
	resp, err := favoritesService.GetForUser(context.TODO(), &pb.UserFavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
	assert.Equal(t, 0, len(resp.GetFavorites()))
}

// TODO
func testFavoritesServiceGetForUserPositive2(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, repoFavorites, userLookup)
	resp, err := favoritesService.GetForUser(context.TODO(), &pb.UserFavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
	assert.Equal(t, 0, len(resp.GetFavorites()))
}

func testFavoritesServiceGetForUserNegative(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, repoFavorites, userLookup)
	resp, err := favoritesService.GetForUser(context.TODO(), &pb.UserFavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
	assert.Equal(t, 0, len(resp.GetFavorites()))
}

func testFavoritesServiceSetPositive(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(context.TODO(), gomock.Any()).Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(context.TODO(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(context.TODO(), gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, repoFavorites, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.Nil(t, err)
	assert.Equal(t, pb.Status_OK, resp.GetStatus())
	assert.Equal(t, "", resp.GetFavorites().GetAsset().GetIsin())
	assert.Equal(t, "", resp.GetFavorites().GetAsset().GetAssetType().GetName())
	assert.Equal(t, "", resp.GetFavorites().GetUser().GetUpk())
}

func testFavoritesServiceSetNegative1(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(context.TODO(), gomock.Any()).Return(false).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(context.TODO(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(context.TODO(), gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, repoFavorites, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceSetNegative2(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(context.TODO(), gomock.Any()).Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(context.TODO(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(context.TODO(), gomock.Any(), gomock.Any()).
		Return(false).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, repoFavorites, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func testFavoritesServiceSetNegative3(t *testing.T) {
	ctrl := gomock.NewController(t)
	assetLookup := NewMockAssetSearchService(ctrl)
	dftFavorites := NewMockDft[*entity.Favorites](ctrl)
	repoFavorites := NewMockRepo[*entity.Favorites](ctrl)
	userLookup := NewMockUserSearchService(ctrl)
	assetLookup.
		EXPECT().
		Lookup(context.TODO(), gomock.Any()).Return(true).
		AnyTimes()
	dftFavorites.
		EXPECT().
		DoUpsert(context.TODO(), gomock.Any(), gomock.Any()).
		Return(pgx.ErrTxCommitRollback).
		AnyTimes()
	userLookup.
		EXPECT().
		Lookup(context.TODO(), gomock.Any(), gomock.Any()).
		Return(true).
		AnyTimes()
	favoritesService := getTestFavoritesService(assetLookup, dftFavorites, repoFavorites, userLookup)
	resp, err := favoritesService.Set(context.TODO(), &pb.FavoritesRequest{})
	assert.NotNil(t, err)
	assert.Equal(t, pb.Status_FAIL, resp.GetStatus())
}

func getTestFavoritesService(
	assetLookup AssetSearchService,
	dftFavorites domain.Dft[*entity.Favorites],
	repoFavorites domain.Repo[*entity.Favorites],
	userLookup UserSearchService,
) FavoritesService {
	favoritesSrv = new(favoritesService)
	favoritesSrv.assetLookup = assetLookup
	favoritesSrv.dftFavorites = dftFavorites
	favoritesSrv.repoFavorites = repoFavorites
	favoritesSrv.userLookup = userLookup
	return favoritesSrv
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
