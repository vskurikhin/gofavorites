/*
 * This file was last modified at 2024-07-30 12:03 by Victor N. Skurikhin.
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
	"fmt"
	"github.com/ssoroka/slice"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/mongo"
	"github.com/vskurikhin/gofavorites/internal/domain/repo"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/models"
	"github.com/vskurikhin/gofavorites/internal/tool"
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
	mongo         mongo.Mongo
	repoFavorites domain.Repo[*entity.Favorites]
	upkUtil       UpkUtilService
	userLookup    UserSearchService
}

var _ FavoritesService = (*favoritesService)(nil)
var (
	ErrRequestNil = fmt.Errorf("request is nil")
	onceFavorites = new(sync.Once)
	favoritesServ *favoritesService
)

func GetFavoritesService(prop env.Properties) FavoritesService {

	onceFavorites.Do(func() {
		favoritesServ = new(favoritesService)
		favoritesServ.assetLookup = GetAssetSearchService(prop)
		favoritesServ.dftFavorites = repo.GetFavoritesTxPostgres(prop)
		favoritesServ.mongo = mongo.GetMongoRepo(prop)
		favoritesServ.repoFavorites = repo.GetFavoritesPostgresCachedRepo(prop)
		favoritesServ.upkUtil = GetUpkUtilService(prop)
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
		return &response, err
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

func (f *favoritesService) encrypt(personalKey string) (string, error) {

	upk, err := f.upkUtil.EncryptPersonalKey(personalKey)

	if err != nil {
		slog.Error(env.MSG+" FavoritesService.encrypt", "err", err)
		return "", tool.ErrEncryptAES
	}
	return upk, nil
}

func (f *favoritesService) get(ctx context.Context, model models.Favorites) (models.Favorites, error) {

	var (
		err      error
		response models.Favorites
		upk      string
	)
	isin := model.Asset().Isin()
	personalKey := model.User().PersonalKey()
	upk = model.User().Upk()

	if upk == "" {
		if upk, err = f.encrypt(personalKey); err != nil {
			return models.Favorites{}, err
		}
	}
	favorites, err := entity.GetFavorites(ctx, f.repoFavorites, isin, upk)

	if err != nil {
		slog.Error(env.MSG+" FavoritesService.Get", "err", err)
	} else {
		response = models.FavoritesFromEntity(favorites)
	}
	return response, err
}

func (f *favoritesService) getForUser(ctx context.Context, model models.User) ([]models.Favorites, error) {

	var err error
	var upk string

	personalKey := model.PersonalKey()
	upk = model.Upk()

	if upk == "" {
		if upk, err = f.encrypt(personalKey); err != nil {
			return nil, err
		}
	}
	favorites, err := entity.GetFavoritesForUser(ctx, f.repoFavorites, upk)

	if err != nil {
		return nil, err
	}
	response := make([]models.Favorites, 0, len(favorites))
	response = slice.Map[entity.Favorites, models.Favorites](
		favorites,
		func(i int, fav entity.Favorites) models.Favorites {
			return models.FavoritesFromEntity(fav)
		})
	return response, err
}

func (f *favoritesService) set(ctx context.Context, model models.Favorites) (models.Favorites, error) {

	var (
		err      error
		response models.Favorites
		upk      string
	)
	personalKey := model.User().PersonalKey()
	upk = model.User().Upk()

	if upk == "" {
		if upk, err = f.encrypt(personalKey); err != nil {
			return models.Favorites{}, err
		}
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
		slog.Error(env.MSG+" FavoritesService.set", "err", err)
	} else {
		model = model.WithUpk(upk)
		favorites := model.ToEntity()
		err = favorites.Upsert(ctx, f.dftFavorites, func() {
			err := f.mongo.Save(ctx, favorites)
			if err != nil {
				slog.Error(env.MSG+" FavoritesService.set in transaction", "err", err)
				return
			}
		})

		if err != nil {
			slog.Error(env.MSG+" FavoritesService.set", "err", err)
		} else {
			response = models.FavoritesFromEntity(favorites)
		}
	}
	return response, err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
