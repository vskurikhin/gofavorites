/*
 * This file was last modified at 2024-08-04 22:13 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * mongo.go
 * $Id$
 */
//!+

// Package mongo TODO.
package mongo

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/vskurikhin/gofavorites/internal/tool"

	"github.com/google/uuid"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	Collection = "favorites"
	AssetType  = "asset-type"
	UPK        = "upk"
	ISIN       = "isin"
	Version    = "version"
)

type Mongo interface {
	Delete(ctx context.Context, entity entity.Favorites) error
	Load(ctx context.Context, upk string) ([]entity.Favorites, error)
	Save(ctx context.Context, entity entity.Favorites) error
}

type repo struct {
	dbName      string
	mongodbPool *tool.MongoPool
	sLog        *slog.Logger
}

type favorites struct {
	ID        primitive.ObjectID `bson:"_id"`
	Upk       string             `bson:"upk"`
	Isin      string             `bson:"isin"`
	AssetType string             `bson:"asset-type"`
	Version   int64              `bson:"version"`
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
		mongoRepo.sLog = prop.Logger()
	})
	return mongoRepo
}

func (r *repo) Delete(ctx context.Context, entity entity.Favorites) error {

	conn, err := r.mongodbPool.GetConnection()

	if err != nil {
		return err
	}
	defer func() { _ = r.mongodbPool.CloseConnection(conn) }()

	collection := tool.GetCollection(conn, r.dbName, Collection)
	res, err := collection.DeleteOne(ctx, bson.D{
		{Key: UPK, Value: entity.User().Upk()},
		{Key: ISIN, Value: entity.Asset().Isin()},
	})
	fmt.Printf("Number of documents deleted: %d\n", res.DeletedCount)

	return err
}

func (r *repo) Load(ctx context.Context, upk string) ([]entity.Favorites, error) {

	conn, err := r.mongodbPool.GetConnection()
	result := make([]entity.Favorites, 0)

	if err != nil {
		return result, err
	}
	defer func() { _ = r.mongodbPool.CloseConnection(conn) }()

	collection := tool.GetCollection(conn, r.dbName, Collection)
	cur, err := collection.Find(ctx, bson.D{
		{Key: UPK, Value: upk},
	})
	if err != nil {
		return result, err
	}
	for cur.Next(ctx) {
		// To decode into a struct, use cursor.Decode()
		var fav favorites
		err = cur.Decode(&fav)

		if err != nil {
			r.sLog.DebugContext(ctx, env.MSG+"MongoRepo.Save", "cur.Decode(&result)", err)
			return result, err
		}
		at := entity.MakeAssetType(fav.AssetType, entity.DefaultTAttributes())
		us := entity.MakeUserWithVersion(fav.Upk, fav.Version, entity.DefaultTAttributes())
		as := entity.MakeAsset(fav.Isin, at, entity.DefaultTAttributes())
		vn := sql.NullInt64{Int64: fav.Version, Valid: true}
		fv := entity.MakeFavorites(uuid.Max, as, us, vn, entity.DefaultTAttributes())
		result = append(result, fv)
	}
	return result, nil
}

func (r *repo) Save(ctx context.Context, entity entity.Favorites) error {

	conn, err := r.mongodbPool.GetConnection()

	if err != nil {
		return err
	}
	defer func() { _ = r.mongodbPool.CloseConnection(conn) }()

	collection := tool.GetCollection(conn, r.dbName, Collection)
	cur, err := collection.Find(ctx, bson.D{
		{Key: UPK, Value: entity.User().Upk()},
		{Key: ISIN, Value: entity.Asset().Isin()},
	})

	if err != nil {
		return err
	}
	defer func() { _ = cur.Close(ctx) }()

	if cur.RemainingBatchLength() == 0 {
		res, err := collection.InsertOne(ctx, bson.D{
			{Key: UPK, Value: entity.User().Upk()},
			{Key: ISIN, Value: entity.Asset().Isin()},
			{Key: AssetType, Value: entity.Asset().AssetType().Name()},
			{Key: Version, Value: entity.Version().Int64},
		})
		r.sLog.DebugContext(ctx, env.MSG+"MongoRepo.Save", "res.InsertedID", res.InsertedID, "err", err)
	} else {
		for cur.Next(ctx) {
			// To decode into a struct, use cursor.Decode()
			var result favorites
			err := cur.Decode(&result)

			if err != nil {
				r.sLog.DebugContext(ctx, env.MSG+"MongoRepo.Save cur.Decode", "err", err)
				return err
			}
			if result.Version < entity.Version().Int64 {
				filter := bson.D{{Key: "_id", Value: result.ID}}
				update := bson.D{{Key: "$set",
					Value: bson.D{
						{Key: UPK, Value: entity.User().Upk()},
						{Key: ISIN, Value: entity.Asset().Isin()},
						{Key: AssetType, Value: entity.Asset().AssetType().Name()},
						{Key: Version, Value: entity.Version().Int64},
					}},
				}
				res, err := collection.UpdateOne(ctx, filter, update)
				if err != nil {
					r.sLog.DebugContext(ctx, env.MSG+"MongoRepo.Save collection.UpdateOne", "err", err)
					return err
				}
				r.sLog.DebugContext(ctx, env.MSG+"MongoRepo.Save", "res.UpsertedID", res.UpsertedID, "err", err)
			}
		}
	}
	return nil
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
