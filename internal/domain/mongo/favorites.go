/*
 * This file was last modified at 2024-07-30 12:07 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites.go
 * $Id$
 */
//!+

// Package mongo TODO.
package mongo

import (
	"context"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"sync"
)

const (
	Collection = "favorites"
	UPK        = "upk"
	ISIN       = "isin"
	VERSION    = "version"
)

type Mongo interface {
	MaxVersion(ctx context.Context, entity entity.User) int64
	Save(ctx context.Context, entity entity.Favorites) error
}

type repo struct {
	dbName      string
	mongodbPool *tool.MongoPool
}

var _ Mongo = (*repo)(nil)
var (
	onceMongo = new(sync.Once)
	mongoRepo *repo
)

func GetMongoRepo(prop env.Properties) Mongo {
	onceMongo.Do(func() {
		mongoRepo = new(repo)
		mongoRepo.dbName = prop.Config().MongoName()
		mongoRepo.mongodbPool = prop.MongodbPool()

	})
	return mongoRepo
}

func (r *repo) MaxVersion(ctx context.Context, entity entity.User) int64 {

	conn, err := r.mongodbPool.GetConnection()

	if err != nil {
		return 0
	}
	defer func() { _ = r.mongodbPool.CloseConnection(conn) }()

	collection := tool.GetCollection(conn, r.dbName, Collection)
	opts := options.FindOne().SetSort(bson.D{{VERSION, -1}})
	res := collection.FindOne(ctx, bson.D{
		{UPK, entity.Upk()},
	}, opts)

	var result favorites
	err = res.Decode(&result)

	if err != nil {
		return 0
	}
	return result.Version
}

func (r *repo) Save(ctx context.Context, entity entity.Favorites) error {

	conn, err := r.mongodbPool.GetConnection()

	if err != nil {
		return err
	}
	defer func() { _ = r.mongodbPool.CloseConnection(conn) }()

	collection := tool.GetCollection(conn, r.dbName, Collection)
	cur, err := collection.Find(ctx, bson.D{
		{UPK, entity.User().Upk()},
		{ISIN, entity.Asset().Isin()},
	})

	if err != nil {
		return err
	}
	defer func() { _ = cur.Close(ctx) }()

	if cur.RemainingBatchLength() == 0 {
		res, err := collection.InsertOne(ctx, bson.D{
			{UPK, entity.User().Upk()},
			{ISIN, entity.Asset().Isin()},
			{VERSION, entity.Version().Int64},
		})
		slog.Debug(env.MSG+" repo.Save", "res.InsertedID", res.InsertedID, "err", err)
	} else {
		for cur.Next(ctx) {
			// To decode into a struct, use cursor.Decode()
			var result favorites
			err := cur.Decode(&result)

			if err != nil {
				slog.Debug(env.MSG+" repo.Save", "cur.Decode(&result)", err)
				return err
			}
			if result.Version < entity.Version().Int64 {
				filter := bson.D{{"_id", result.ID}}
				update := bson.D{{"$set",
					bson.D{
						{UPK, entity.User().Upk()},
						{ISIN, entity.Asset().Isin()},
						{VERSION, entity.Version().Int64},
					}},
				}
				res, err := collection.UpdateOne(ctx, filter, update)
				if err != nil {
					slog.Debug(env.MSG+" repo.Save", "collection.UpdateOne", err)
					return err
				}
				slog.Debug(env.MSG+" FavoritesService.Set", "res.UpsertedID", res.UpsertedID, "err", err)
			}
		}
	}
	return nil
}

type favorites struct {
	ID      primitive.ObjectID `bson:"_id"`
	Upk     string             `bson:"upk"`
	Isin    string             `bson:"isin"`
	Version int64              `bson:"version"`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
