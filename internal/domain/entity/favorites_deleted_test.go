/*
 * This file was last modified at 2024-07-15 17:10 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites_deleted_test.go
 * $Id$
 */
//!+

// Package entity TODO.
package entity

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFavoritesDeleted(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{name: "positive test #0 FavoritesDeleted Cloneable", fRun: testFavoritesDeletedCloneable},
		{name: "positive test #1 FavoritesDeleted FromJSON and ToJSON", fRun: testFavoritesDeletedJSON},
		{name: "positive test #2 FavoritesDeleted stubRepoOk", fRun: testFavoritesDeletedRepoOk},
		{name: "negative test #3 FavoritesDeleted stubRepoErr", fRun: testFavoritesDeletedRepoErr},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testFavoritesDeletedCloneable(t *testing.T) {
	expected := MakeFavoritesDeletedUser("")
	got := expected.Copy()
	assert.NotNil(t, got)
	assert.Equal(t, &expected, got)
}

func testFavoritesDeletedJSON(t *testing.T) {
	expected := MakeFavoritesDeletedUser("")
	j, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.NotNil(t, j)
	got := FavoritesDeleted{}
	err = (&got).FromJSON(j)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.String(), got.String())
	assert.Equal(t, "                                ", got.Key())
}

func testFavoritesDeletedRepoOk(t *testing.T) {
	favorites := MakeFavoritesDeletedUser("")
	err := favorites.Delete(context.TODO(), &stubRepoOk[*FavoritesDeleted]{})
	assert.Nil(t, err)
	got, err := GetFavoritesDeletedForUser(context.TODO(), &stubRepoOk[*FavoritesDeleted]{}, "")
	assert.Nil(t, err)
	assert.Equal(t, []FavoritesDeleted{favorites}, got)
	err = favorites.Update(context.TODO(), &stubRepoOk[*FavoritesDeleted]{})
	assert.Nil(t, err)
}

func testFavoritesDeletedRepoErr(t *testing.T) {
	favorites := MakeFavoritesDeletedUser("")
	err := favorites.Delete(context.TODO(), &stubRepoErr[*FavoritesDeleted]{})
	assert.NotNil(t, err)
	_, err = GetFavoritesDeletedForUser(context.TODO(), &stubRepoErr[*FavoritesDeleted]{}, "")
	assert.NotNil(t, err)
	err = favorites.Update(context.TODO(), &stubRepoErr[*FavoritesDeleted]{})
	assert.NotNil(t, err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
