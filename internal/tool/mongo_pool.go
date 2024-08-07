/*
 * This file was last modified at 2024-08-05 23:27 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * mongo_pool.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoPool struct {
	pool        chan *mongo.Client
	timeout     time.Duration
	uri         string
	connections int
	poolSize    int
}

var (
	mongoPool     *MongoPool
	onceMongoPool sync.Once
)

func MongodbConnect(dsn string) *MongoPool {

	onceMongoPool.Do(func() {
		mongoPool = &MongoPool{
			pool:        make(chan *mongo.Client, 10),
			connections: 0,
			timeout:     500 * time.Millisecond,
			uri:         dsn,
			poolSize:    10,
		}
	})
	return mongoPool
}

//goland:noinspection GoVetLostCancel
func (mp *MongoPool) getContextTimeOut() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), mp.timeout)
	return ctx
}

func (mp *MongoPool) createToChan() {

	var client *mongo.Client
	client, err := mongo.Connect(mp.getContextTimeOut(), options.Client().ApplyURI(mp.uri))

	if err != nil {
		sLog.Error(MSG+"Create the MongoPool failed", "err", err)
	}
	mp.pool <- client
	mp.connections++
}

func (mp *MongoPool) CloseConnection(conn *mongo.Client) error {
	select {
	case mp.pool <- conn:
		return nil
	default:
		if err := conn.Disconnect(context.TODO()); err != nil {
			sLog.Error(MSG+"Close the MongoPool failed", "err", err)
			return err
		}
		mp.connections--
		return nil
	}
}

func (mp *MongoPool) GetConnection() (*mongo.Client, error) {
	for {
		select {
		case conn := <-mp.pool:
			err := conn.Ping(mp.getContextTimeOut(), readpref.Primary())
			if err != nil {
				sLog.Error(MSG+"GetConnection: Failed to obtain connection mongoPool connection", "err", err)
				return nil, err
			}
			return conn, nil
		default:
			if mp.connections < mp.poolSize {
				mp.createToChan()
			}
		}
	}
}

func GetCollection(conn *mongo.Client, dbname, collection string) *mongo.Collection {
	return conn.Database(dbname).Collection(collection)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
