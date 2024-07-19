/*
 * This file was last modified at 2024-07-15 17:10 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites_test.go
 * $Id$
 */
//!+

// Package entity TODO.
package entity

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFavorites(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{name: "positive test #0 Favorites Cloneable", fRun: testFavoritesCloneable},
		{name: "positive test #1 Favorites FromJSON and ToJSON", fRun: testFavoritesJSON},
		{name: "positive test #2 Favorites IsFavoritesNotFound", fRun: testIsFavoritesNotFound},
		{name: "positive test #3 Favorites stubRepoOk", fRun: testFavoritesRepoOk},
		{name: "negative test #4 Favorites stubRepoErr", fRun: testFavoritesRepoErr},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testFavoritesCloneable(t *testing.T) {
	expected := MakeFavorites(
		uuid.New(),
		Asset{},
		User{},
		sql.NullInt64{},
		MakeTAttributes(
			sql.NullBool{true, true},
			time.Time{},
			sql.NullTime{time.Time{}, false},
		))
	got := expected.Copy()
	assert.NotNil(t, got)
	assert.Equal(t, &expected, got)
}

func testFavoritesJSON(t *testing.T) {
	expected := MakeFavorites(
		uuid.New(),
		Asset{},
		User{},
		sql.NullInt64{},
		MakeTAttributes(
			sql.NullBool{true, true},
			time.Time{},
			sql.NullTime{time.Time{}, false},
		))
	j, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.NotNil(t, j)
	got := Favorites{}
	err = (&got).FromJSON(j)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.String(), got.String())
	assert.Equal(t, "                                ", got.Key())
}

func testIsFavoritesNotFound(t *testing.T) {
	assert.True(t, IsFavoritesNotFound(Favorites{id: uuid.New()}, pgx.ErrNoRows))
	assert.True(t, IsFavoritesNotFound(Favorites{}, errors.New("")))
	assert.False(t, IsFavoritesNotFound(Favorites{id: uuid.New()}, errors.New("")))
}

func testFavoritesRepoOk(t *testing.T) {
	favorites := MakeFavorites(uuid.New(), Asset{}, User{}, sql.NullInt64{}, DefaultTAttributes())
	err := favorites.Upsert(context.TODO(), &stubTxRepoOk[*Favorites]{}, func() {})
	assert.Nil(t, err)
	assert.False(t, favorites.Deleted().Valid)
	got, err := GetFavorites(context.TODO(), &stubRepoOk[*Favorites]{}, "", "")
	assert.Nil(t, err)
	assert.Equal(t, favorites.CreatedAt(), got.CreatedAt())
	assert.Equal(t, favorites.Deleted(), got.Deleted())
	assert.False(t, favorites.UpdatedAt().Valid)
	err = favorites.Delete(context.TODO(), &stubTxRepoOk[*Favorites]{}, func() {})
	assert.Nil(t, err)
	assert.False(t, favorites.Deleted().Valid)
}

func testFavoritesRepoErr(t *testing.T) {
	favorites := MakeFavorites(uuid.New(), Asset{}, User{}, sql.NullInt64{}, DefaultTAttributes())
	err := favorites.Upsert(context.TODO(), &stubTxRepoErr[*Favorites]{}, func() {})
	assert.NotNil(t, err)
	_, err = GetFavorites(context.TODO(), &stubRepoErr[*Favorites]{}, "", "")
	assert.NotNil(t, err)
	err = favorites.Delete(context.TODO(), &stubTxRepoErr[*Favorites]{}, func() {})
	assert.NotNil(t, err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
