/*
 * This file was last modified at 2024-07-15 17:10 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * user_test.go
 * $Id$
 */
//!+

// Package entity TODO.
package entity

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{name: "positive test #0 User Cloneable", fRun: testUserCloneable},
		{name: "positive test #1 User FromJSON and ToJSON", fRun: testUserJSON},
		{name: "positive test #2 User IsUserNotFound", fRun: testIsUserNotFound},
		{name: "positive test #3 User stubRepoOk", fRun: testUserRepoOk},
		{name: "negative test #4 User stubRepoErr", fRun: testUserRepoErr},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testUserCloneable(t *testing.T) {
	expected := MakeUser(tool.RandStringBytes(32), MakeTAttributes(
		sql.NullBool{true, true},
		time.Time{},
		sql.NullTime{time.Time{}, false},
	))
	got := expected.Copy()
	assert.NotNil(t, got)
	assert.Equal(t, &expected, got)
}

func testUserJSON(t *testing.T) {
	upk := tool.RandStringBytes(32)
	expected := MakeUser(upk, MakeTAttributes(
		sql.NullBool{true, true},
		time.Time{},
		sql.NullTime{time.Time{}, false},
	))
	j, err := expected.ToJSON()
	assert.Nil(t, err)
	assert.NotNil(t, j)
	got := User{}
	err = (&got).FromJSON(j)
	assert.Nil(t, err)
	assert.Equal(t, expected, got)
	assert.Equal(t, expected.String(), got.String())
	assert.Equal(t, upk, got.Key())
}

func testIsUserNotFound(t *testing.T) {
	assert.True(t, IsUserNotFound(User{upk: "test"}, pgx.ErrNoRows))
	assert.True(t, IsUserNotFound(User{}, errors.New("")))
	assert.False(t, IsUserNotFound(User{upk: "test"}, errors.New("")))
}

func testUserRepoOk(t *testing.T) {
	upk := tool.RandStringBytes(32)
	user := MakeUser(upk, DefaultTAttributes())
	err := user.Insert(context.TODO(), &stubRepoOk[*User]{})
	assert.Nil(t, err)
	assert.False(t, user.Deleted().Valid)
	got, err := GetUser(context.TODO(), &stubRepoOk[*User]{}, upk)
	assert.Nil(t, err)
	assert.Equal(t, user, got)
	assert.Equal(t, user.CreatedAt(), got.CreatedAt())
	assert.Equal(t, user.Upk(), got.Upk())
	assert.Equal(t, user.Deleted(), got.Deleted())
	assert.False(t, user.UpdatedAt().Valid)
	err = user.Update(context.TODO(), &stubRepoOk[*User]{})
	assert.Nil(t, err)
	assert.False(t, user.Deleted().Valid)
	assert.False(t, user.UpdatedAt().Valid)
	err = user.Delete(context.TODO(), &stubRepoOk[*User]{})
	assert.Nil(t, err)
	assert.False(t, user.Deleted().Valid)
}

func testUserRepoErr(t *testing.T) {
	upk := tool.RandStringBytes(32)
	user := MakeUser(upk, DefaultTAttributes())
	err := user.Insert(context.TODO(), &stubRepoErr[*User]{})
	assert.NotNil(t, err)
	_, err = GetUser(context.TODO(), &stubRepoErr[*User]{}, upk)
	assert.NotNil(t, err)
	err = user.Update(context.TODO(), &stubRepoErr[*User]{})
	assert.NotNil(t, err)
	err = user.Delete(context.TODO(), &stubRepoErr[*User]{})
	assert.NotNil(t, err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
