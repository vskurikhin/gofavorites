/*
 * This file was last modified at 2024-07-29 11:38 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites_service.go
 * $Id$
 */
//!+

// Package services TODO.
package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/ssoroka/slice"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/repo"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/models"
	pb "github.com/vskurikhin/gofavorites/proto"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"sync"
)

type FavoritesService interface {
	ApiFavoritesService
	pb.FavoritesServiceServer
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
	ErrRequestNil = fmt.Errorf("request is nil")
	onceFavorites = new(sync.Once)
	favoritesServ *favoritesService
	inTransaction = func() {
		// TODO
	}
)

func GetFavoritesService(prop env.Properties) FavoritesService {

	onceFavorites.Do(func() {
		favoritesServ = new(favoritesService)
		favoritesServ.assetLookup = GetAssetSearchService(prop)
		favoritesServ.dftFavorites = repo.GetFavoritesTxPostgres(prop)
		favoritesServ.repoFavorites = repo.GetFavoritesPostgresCachedRepo(prop)
		favoritesServ.userLookup = GetUserSearchService(prop)
	})
	return favoritesServ
}

func (f *favoritesService) ApiFavoritesGet(
	ctx context.Context,
	favorites models.Favorites,
) (response models.Favorites, err error) {
	return f.get(ctx, favorites)
}

func (f *favoritesService) ApiFavoritesGetForUser(
	ctx context.Context,
	user models.User,
) (response []models.Favorites, err error) {
	return f.getForUser(ctx, user)
}

func (f *favoritesService) ApiFavoritesSet(
	ctx context.Context,
	favorites models.Favorites,
) (response models.Favorites, err error) {
	return f.set(ctx, favorites)
}

func (f *favoritesService) Get(ctx context.Context, request *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {

	var response pb.FavoritesResponse

	if request == nil {
		response.Status = pb.Status_FAIL
		return &response, ErrRequestNil
	}
	favorites := models.FavoritesFromProto(request.GetFavorites())
	model, err := f.get(ctx, favorites)

	if err != nil {
		slog.Error(env.MSG+" FavoritesService.Get", "err", err)
		response.Status = pb.Status_FAIL
	} else {
		response.Favorites = model.ToProto()
		response.Status = pb.Status_OK
	}
	return &response, err
}

func (f *favoritesService) GetForUser(ctx context.Context, request *pb.UserFavoritesRequest) (*pb.UserFavoritesResponse, error) {

	var response pb.UserFavoritesResponse

	if request == nil {
		response.Status = pb.Status_FAIL
		return &response, ErrRequestNil
	}
	user := models.UserFromProto(request.GetUser())
	favorites, err := f.getForUser(ctx, user)

	if err != nil {
		response.Status = pb.Status_FAIL
		return nil, err
	}
	response.Favorites = make([]*pb.Favorites, len(favorites))

	for _, fav := range favorites {
		item := fav.ToProto()
		response.Favorites = append(response.Favorites, item)
	}
	response.Status = pb.Status_OK

	return &response, err
}

func (f *favoritesService) Set(ctx context.Context, request *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {

	var response pb.FavoritesResponse

	if request == nil {
		response.Status = pb.Status_FAIL
		return &response, ErrRequestNil
	}
	favorites := models.FavoritesFromProto(request.GetFavorites())
	model, err := f.set(ctx, favorites)

	if err != nil {
		response.Status = pb.Status_FAIL
	} else {
		response.Favorites = model.ToProto()
		response.Status = pb.Status_OK
	}
	return &response, err
}

func (f *favoritesService) get(ctx context.Context, model models.Favorites) (response models.Favorites, err error) {

	var upk string
	isin := model.Asset().Isin()
	personalKey := model.User().PersonalKey()
	upk = model.User().Upk()

	if upk == "" {
		upk = base64.StdEncoding.EncodeToString([]byte(personalKey)) // TODO RSA Encrypt
	}
	favorites, err := entity.GetFavorites(ctx, f.repoFavorites, isin, upk)

	if err != nil {
		slog.Error(env.MSG+" FavoritesService.Get", "err", err)
	} else {
		response = favorites.ToModel()
	}
	return response, err
}

func (f *favoritesService) getForUser(ctx context.Context, model models.User) ([]models.Favorites, error) {

	var upk string
	personalKey := model.PersonalKey()
	upk = model.Upk()

	if upk == "" {
		upk = base64.StdEncoding.EncodeToString([]byte(personalKey)) // TODO RSA Encrypt
	}
	favorites, err := entity.GetFavoritesForUser(ctx, f.repoFavorites, upk)

	if err != nil {
		return nil, err
	}
	response := make([]models.Favorites, 0, len(favorites))
	response = slice.Map[entity.Favorites, models.Favorites](favorites,
		func(i int, fav entity.Favorites) models.Favorites {
			return fav.ToModel()
		})
	return response, err
}

func (f *favoritesService) set(ctx context.Context, model models.Favorites) (response models.Favorites, err error) {

	var upk string
	personalKey := model.User().PersonalKey()
	upk = model.User().Upk()

	if upk == "" {
		upk = base64.StdEncoding.EncodeToString([]byte(personalKey)) // TODO RSA Encrypt
	}
	user := models.MakeUser(personalKey, upk)
	g, c := errgroup.WithContext(ctx)

	g.Go(func() error {
		if f.userLookup.Lookup(c, user) {
			return nil
		}
		return fmt.Errorf("user by upk: %s not found", model.User().Upk())
	})
	g.Go(func() error {
		if f.assetLookup.Lookup(c, model.Asset().Isin()) {
			return nil
		}
		return fmt.Errorf("asset by isin: %s not found", model.Asset().Isin())
	})
	if err = g.Wait(); err != nil {
		slog.Error(env.MSG+" FavoritesService.Set", "err", err)
	} else {
		favorites := entity.FavoritesFromModel(
			models.MakeFavorites(model.Id(), model.Asset(), user, model.Version()),
		)
		err = favorites.Upsert(ctx, f.dftFavorites, inTransaction)

		if err != nil {
			slog.Error(env.MSG+" FavoritesService.Set", "err", err)
		} else {
			response = favorites.ToModel()
		}
	}
	return response, err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
