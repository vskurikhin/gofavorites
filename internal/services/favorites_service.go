/*
 * This file was last modified at 2024-07-20 10:30 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites_service.go
 * $Id$
 */

package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/repo"
	"github.com/vskurikhin/gofavorites/internal/env"
	pb "github.com/vskurikhin/gofavorites/proto"
	"log/slog"
	"sync"
)

type FavoritesService interface {
	pb.FavoritesServiceServer
}

type AssetSearchService interface {
	Lookup(ctx context.Context, isin string) bool
}

type UserSearchService interface {
	Lookup(ctx context.Context, personalKey, upk string) bool
}

type favoritesService struct {
	pb.UnimplementedFavoritesServiceServer
	assetLookup   AssetSearchService
	dftFavorites  domain.Dft[*entity.Favorites]
	repoFavorites domain.Repo[*entity.Favorites]
	userLookup    UserSearchService
}

var _ FavoritesService = (*favoritesService)(nil)
var (
	onceFavoritesService = new(sync.Once)
	favoritesSrv         *favoritesService
	inTransaction        = func() {
		// TODO
	}
)

func GetFavoritesService(prop env.Properties) FavoritesService {

	onceFavoritesService.Do(func() {
		favoritesSrv = new(favoritesService)
		favoritesSrv.assetLookup = GetAssetSearchService(prop)
		favoritesSrv.dftFavorites = repo.GetFavoritesTxPostgres(prop)
		favoritesSrv.repoFavorites = repo.GetFavoritesPostgresCachedRepo(prop)
		favoritesSrv.userLookup = GetUserSearchService(prop)
	})
	return favoritesSrv
}

func (f *favoritesService) Get(ctx context.Context, request *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {

	var response pb.FavoritesResponse
	favoritesProto := request.GetFavorites()
	isin := favoritesProto.GetAsset().GetIsin()
	personalKey := favoritesProto.GetUser().GetPersonalKey()
	upk := base64.StdEncoding.EncodeToString([]byte(personalKey)) // TODO RSA Encrypt
	favorites, err := entity.GetFavorites(ctx, f.repoFavorites, isin, upk)

	if err != nil {
		slog.Error(env.MSG+" FavoritesService.Get", "err", err)
		response.Status = pb.Status_FAIL
	} else {
		favoritesProto := favorites.ToProto()
		response.Status = pb.Status_OK
		response.Favorites = &favoritesProto
	}
	return &response, err
}

func (f *favoritesService) GetForUser(ctx context.Context, request *pb.UserFavoritesRequest) (*pb.UserFavoritesResponse, error) {

	var response pb.UserFavoritesResponse
	userProto := request.GetUser()
	personalKey := userProto.GetPersonalKey()
	upk := base64.StdEncoding.EncodeToString([]byte(personalKey)) // TODO RSA Encrypt
	favorites, err := entity.GetFavoritesForUser(ctx, f.repoFavorites, upk)

	if err != nil {
		slog.Error(env.MSG+" FavoritesService.Get", "err", err)
		response.Status = pb.Status_FAIL
	} else {
		response.Favorites = make([]*pb.Favorites, 0)

		for _, f := range favorites {
			p := f.ToProto()
			response.Favorites = append(response.Favorites, &p)
		}
		response.Status = pb.Status_OK
	}
	return &response, err
}

func (f *favoritesService) Set(ctx context.Context, request *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {

	var err error
	var response pb.FavoritesResponse
	favoritesProto := request.GetFavorites()
	favoritesProto.GetAsset()
	personalKey := favoritesProto.GetUser().GetPersonalKey()
	upk := base64.StdEncoding.EncodeToString([]byte(personalKey)) // TODO RSA Encrypt
	isin := favoritesProto.GetAsset().GetIsin()

	if !f.assetLookup.Lookup(ctx, isin) {
		response.Status = pb.Status_FAIL
		err = fmt.Errorf("asset by isin: %s not found", isin)
		slog.Error(env.MSG+" FavoritesService.Set", "err", err)
	} else if !f.userLookup.Lookup(ctx, personalKey, upk) {
		response.Status = pb.Status_FAIL
		err = fmt.Errorf("user by upk: %s not found", upk)
		slog.Error(env.MSG+" FavoritesService.Set", "err", err)
	} else {
		favorites := entity.FavoritesFromProto(favoritesProto, upk)
		err = favorites.Upsert(ctx, f.dftFavorites, inTransaction)

		if err != nil {
			response.Status = pb.Status_FAIL
			slog.Error(env.MSG+" FavoritesService.Set", "err", err)
		} else {
			favoritesProto := favorites.ToProto()
			response.Status = pb.Status_OK
			response.Favorites = &favoritesProto
		}
	}
	return &response, err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
