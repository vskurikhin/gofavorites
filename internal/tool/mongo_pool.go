/*
 * This file was last modified at 2024-07-29 22:31 by Victor N. Skurikhin.
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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"sync"
	"time"
)

type MongoPool struct {
	pool        chan *mongo.Client
	timeout     time.Duration
	uri         string
	connections int
	poolSize    int
}

var (
	mongoPool *MongoPool
	onceMongo sync.Once
)

func MongodbConnect(dsn string) *MongoPool {

	onceMongo.Do(func() {
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

func (mp *MongoPool) getContextTimeOut() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), mp.timeout)
	return ctx
}

func (mp *MongoPool) createToChan() {

	var client *mongo.Client
	client, err := mongo.Connect(mp.getContextTimeOut(), options.Client().ApplyURI(mp.uri))

	if err != nil {
		log.Fatalf("Create the Pool failed，err=%v", err)
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
			log.Fatalf("Close the Pool failed，err=%v", err)
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
				log.Fatalf("Failed to obtain connection pool connection，err=%v", err)
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
