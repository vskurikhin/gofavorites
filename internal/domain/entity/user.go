/*
 * This file was last modified at 2024-07-15 17:10 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * user.go
 * $Id$
 */

//!+

// Package entity TODO.
package entity

import (
	"context"
	"database/sql"
	"github.com/goccy/go-json"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"time"
)

type User struct {
	upk       string
	deleted   sql.NullBool
	createdAt time.Time
	updatedAt sql.NullTime
}

type userJSON struct {
	UPK       string
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

var _ domain.Entity = (*User)(nil)

func GetUser(ctx context.Context, repo domain.Repo, upk string) (User, error) {

	var e error
	var result User

	err := repo.Get(ctx, &User{upk: upk}, func(scanner domain.Scanner) {
		e = scanner.Scan(
			&result.upk,
			&result.deleted,
			&result.createdAt,
			&result.updatedAt,
		)
	})
	if e != nil {
		return User{}, e
	}
	if err != nil {
		return User{}, err
	}
	return result, nil
}

func NewUser(upk string, createdAt time.Time) User {
	return User{
		upk:       upk,
		createdAt: createdAt,
	}
}

func (u *User) Delete(ctx context.Context, repo domain.Repo) (err error) {
	return repo.Delete(ctx, u)
}

func (u *User) DeleteArgs() []any {
	return []any{u.upk}
}

func (u *User) DeleteSQL() string {
	return `UPDATE users SET deleted = true WHERE upk = $1`
}

func (u *User) GetArgs() []any {
	return []any{u.upk}
}

func (u *User) GetSQL() string {
	return `SELECT upk, deleted, created_at, updated_at FROM users WHERE upk = $1`
}

func (u *User) Insert(ctx context.Context, repo domain.Repo) (err error) {
	return repo.Insert(ctx, u)
}

func (u *User) InsertArgs() []any {
	return []any{u.upk, u.createdAt}
}

func (u *User) InsertSQL() string {
	return `INSERT INTO users (upk, created_at) VALUES ($1, $2)`
}

func (u *User) JSON() ([]byte, error) {

	result, err := json.Marshal(userJSON{
		UPK:       u.upk,
		Deleted:   u.deleted.Bool,
		CreatedAt: u.createdAt,
		UpdatedAt: u.updatedAt,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *User) Key() string {
	return u.upk
}

func (u *User) UpdateArgs() []any {
	return []any{u.upk, u.updatedAt}
}

func (u *User) UpdateSQL() string {
	return `UPDATE users SET updatedAt = $2 WHERE upk = $1`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
