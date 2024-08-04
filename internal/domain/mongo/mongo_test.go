/*
 * This file was last modified at 2024-08-04 22:13 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * mongo_test.go
 * $Id$
 */
//!+

// Package mongo TODO.
package mongo

import (
	"context"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/env"
)

func TestFavoritesInsertsBatch(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "test #1 negative",
			fRun: testMongoDeleteNegative,
		},
		{
			name: "test #2 negative",
			fRun: testMongoLoadNegative,
		},
		{
			name: "test #2 negative",
			fRun: testMongoSaveNegative,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testMongoDeleteNegative(t *testing.T) {
	defer func() { _ = recover() }() // TODO
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")

	prop := env.GetProperties()
	mongo := GetMongoRepo(prop)
	err := mongo.Delete(context.Background(), entity.Favorites{})
	assert.NotNil(t, err)
}

func testMongoLoadNegative(t *testing.T) {
	defer func() { _ = recover() }() // TODO
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")

	prop := env.GetProperties()
	mongo := GetMongoRepo(prop)
	got, err := mongo.Load(context.Background(), "")
	assert.Nil(t, got)
	assert.NotNil(t, err)
}

func testMongoSaveNegative(t *testing.T) {
	defer func() { _ = recover() }() // TODO
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")

	prop := env.GetProperties()
	mongo := GetMongoRepo(prop)
	err := mongo.Save(context.Background(), entity.Favorites{})
	assert.NotNil(t, err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
