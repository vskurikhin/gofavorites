/*
 * This file was last modified at 2024-08-03 14:47 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * sync_util_service.go
 * $Id$
 */
//!+

// Package services TODO.
package services

import (
	"context"
	"database/sql"
	"log/slog"
	"slices"
	"sync"

	"github.com/vskurikhin/gofavorites/internal/domain/batch"
	"github.com/vskurikhin/gofavorites/internal/models"

	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/mongo"
	"github.com/vskurikhin/gofavorites/internal/domain/repo"

	"github.com/vskurikhin/gofavorites/internal/env"
)

type SyncUtilService interface {
	Sync(ctx context.Context, mongodbFavorites, pgDBFavorites []entity.Favorites) ([]entity.Favorites, error)
}

type syncUtilService struct {
	assetLookup          AssetSearchService
	batch                batch.FavoritesInsertsBatch
	mongo                mongo.Mongo
	repoFavorites        domain.Repo[*entity.Favorites]
	repoFavoritesDeleted domain.Repo[*entity.FavoritesDeleted]
	sLog                 *slog.Logger
	userLookup           UserSearchService
}

var _ SyncUtilService = (*syncUtilService)(nil)
var (
	onceSyncUtil = new(sync.Once)
	syncUtilServ *syncUtilService
)

func GetSyncUtilService(prop env.Properties) SyncUtilService {

	onceSyncUtil.Do(func() {
		syncUtilServ = new(syncUtilService)
		syncUtilServ.assetLookup = GetAssetSearchService(prop)
		syncUtilServ.batch = batch.GetBatchPostgres(prop)
		syncUtilServ.mongo = mongo.GetMongoRepo(prop)
		syncUtilServ.repoFavorites = repo.GetFavoritesPostgresCachedRepo(prop)
		syncUtilServ.repoFavoritesDeleted = repo.GetFavoritesDeletedPostgresRepo(prop)
		syncUtilServ.sLog = prop.Logger()
		syncUtilServ.userLookup = GetUserSearchService(prop)
	})
	return syncUtilServ
}

func (s syncUtilService) Sync(
	ctx context.Context,
	mongodbFavorites, pgDBFavorites []entity.Favorites,
) ([]entity.Favorites, error) {

	if len(mongodbFavorites) < 1 {
		return pgDBFavorites, nil
	}
	maxMongodbUser := slices.MaxFunc[[]entity.Favorites, entity.Favorites](mongodbFavorites, func(x, y entity.Favorites) int {
		if x.Version().Int64 > y.Version().Int64 {
			return 1
		} else if x.Version().Int64 > y.Version().Int64 {
			return -1
		}
		return 0
	}).User()
	s.sLog.InfoContext(ctx, env.MSG+"Sync", "maxMongodb", maxMongodbUser.Version())

	var maxPostgresUser entity.User

	if len(pgDBFavorites) < 1 {
		if !s.userLookup.Lookup(ctx, models.UserFromEntity(maxMongodbUser)) {
			return mongodbFavorites, nil
		}
		maxPostgresUser = entity.MakeUser(maxMongodbUser.Upk(), entity.DefaultTAttributes())
	} else {
		maxPostgresUser = slices.
			MaxFunc[[]entity.Favorites, entity.Favorites](pgDBFavorites, func(x, y entity.Favorites) int {
			if x.User().Version() > y.User().Version() {
				return 1
			} else if x.User().Version() > y.User().Version() {
				return -1
			}
			return 0
		}).User()
	}
	s.sLog.InfoContext(ctx, env.MSG+"Sync", "maxPostgresUser", maxPostgresUser.Version())
	result := pgDBFavorites

	if maxPostgresUser.Version() > maxMongodbUser.Version() {
		s.sLog.InfoContext(ctx, env.MSG+"Sync", "maxPostgresUser", maxPostgresUser)
		go s.syncToMongoDB(ctx, pgDBFavorites, maxPostgresUser.Upk())
		result = pgDBFavorites
	} else if maxPostgresUser.Version() < maxMongodbUser.Version() {
		s.sLog.InfoContext(ctx, env.MSG+"Sync", "maxMongodb", maxMongodbUser)
		go s.syncToPostgreSQL(ctx, mongodbFavorites, maxPostgresUser.Upk())
		result = mongodbFavorites
	}
	return result, nil
}

func (s *syncUtilService) syncToMongoDB(ctx context.Context, favorites []entity.Favorites, upk string) {

	for _, fav := range favorites {

		var f = fav

		if !fav.Version().Valid {

			a := entity.MakeTAttributes(fav.Deleted(), fav.CreatedAt(), fav.UpdatedAt())
			v := sql.NullInt64{Int64: fav.User().Version(), Valid: true}
			f = entity.MakeFavorites(fav.ID(), fav.Asset(), fav.User(), v, a)
			err := f.Update(ctx, s.repoFavorites)

			if err != nil {
				s.sLog.ErrorContext(ctx, env.MSG+"syncToMongoDB", "err", err)
			}
		}
		if (!f.Deleted().Valid && !f.Deleted().Bool) || !f.Deleted().Bool {
			err := s.mongo.Save(ctx, f)
			if err != nil {
				s.sLog.ErrorContext(ctx, env.MSG+"syncToMongoDB save", "err", err)
			}
		}
	}
	favoritesDeleted, err := entity.GetFavoritesDeletedForUser(ctx, s.repoFavoritesDeleted, upk)
	if err != nil {
		s.sLog.ErrorContext(ctx, env.MSG+"syncToMongoDB get favorites deleted", "err", err)
	}
	for _, deleted := range favoritesDeleted {
		if err := s.mongo.Delete(ctx, *deleted.ToFavorites()); err != nil {
			s.sLog.ErrorContext(ctx, env.MSG+"syncToMongoDB delete in MongoDB", "err", err)
		} else {
			if err := deleted.Update(ctx, s.repoFavoritesDeleted); err != nil {
				s.sLog.ErrorContext(ctx, env.MSG+"syncToMongoDB delete in PostgreSQL", "err", err)
			}
		}
	}
}

func (s *syncUtilService) syncToPostgreSQL(ctx context.Context, favorites []entity.Favorites, upk string) {

	deleted := entity.MakeFavoritesDeletedUser(upk)

	if err := deleted.Delete(ctx, s.repoFavoritesDeleted); err != nil {
		s.sLog.ErrorContext(ctx, env.MSG+"syncToPostgreSQL delete in PostgreSQL", "err", err)
	}
	batchFavorites := make([]entity.Favorites, 0, len(favorites))

	for _, favorite := range favorites {
		if s.assetLookup.Lookup(ctx, favorite.Asset().Isin()) {
			batchFavorites = append(batchFavorites, favorite)
			s.sLog.DebugContext(ctx, env.MSG+"syncToPostgreSQL favorite add to batchFavorites", "favorite", favorite)
		}
	}
	if err := s.batch.Do(ctx, batchFavorites, upk); err != nil {
		s.sLog.ErrorContext(ctx, env.MSG+"syncToPostgreSQL batch to PostgreSQL", "err", err)
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
