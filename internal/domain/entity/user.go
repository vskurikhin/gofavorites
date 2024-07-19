/*
 * This file was last modified at 2024-07-20 11:01 by Victor N. Skurikhin.
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
	"fmt"
	"github.com/goccy/go-json"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"time"
)

type User struct {
	TAttributes
	upk string
}

var _ domain.Entity = (*User)(nil)

func GetUser(ctx context.Context, repo domain.Repo[*User], upk string) (User, error) {

	var e error
	result := &User{upk: upk}

	result, err := repo.Get(ctx, result, func(scanner domain.Scanner) {
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
	return *result, nil
}

func MakeUser(upk string, a TAttributes) User {
	return User{
		TAttributes: struct {
			deleted   sql.NullBool
			createdAt time.Time
			updatedAt sql.NullTime
		}{
			deleted:   a.deleted,
			createdAt: a.createdAt,
			updatedAt: a.updatedAt,
		},
		upk: upk,
	}
}

func IsUserNotFound(u User, err error) bool {
	return tool.NoRowsInResultSet(err) || u == User{}
}

func (u *User) Upk() string {
	return u.upk
}

func (u *User) Deleted() sql.NullBool {
	return u.deleted
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() sql.NullTime {
	return u.updatedAt
}

func (u *User) Copy() domain.Entity {
	c := *u
	return &c
}

func (u *User) Delete(ctx context.Context, repo domain.Repo[*User]) (err error) {

	_, e := repo.Delete(ctx, u, func(s domain.Scanner) {
		t := *u
		err = s.Scan(&t.deleted, &t.updatedAt)
		if err == nil {
			*u = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (u *User) DeleteArgs() []any {
	return []any{u.upk}
}

func (u *User) DeleteSQL() string {
	return `UPDATE users SET deleted = true WHERE upk = $1 RETURNING deleted, updated_at`
}

type userJSON struct {
	UPK       string
	Deleted   *bool `json:",omitempty"`
	CreatedAt time.Time
	UpdatedAt *time.Time `json:",omitempty"`
}

func (u *User) FromJSON(data []byte) (err error) {

	var t userJSON
	err = json.Unmarshal(data, &t)

	if err != nil {
		return err
	}
	u.upk = t.UPK
	u.deleted = tool.ConvertBoolPointerToNullBool(t.Deleted)
	u.createdAt = t.CreatedAt
	u.updatedAt = tool.ConvertTimePointerToNullTime(t.UpdatedAt)

	return nil
}

func (u *User) GetArgs() []any {
	return []any{u.upk}
}

func (u *User) GetByFilterArgs() []any {
	return []any{}
}

func (u *User) GetByFilterSQL() string {
	return `SELECT upk, deleted, created_at, updated_at FROM users WHERE deleted IS NOT TRUE`
}

func (u *User) GetSQL() string {
	return `SELECT upk, deleted, created_at, updated_at FROM users WHERE upk = $1`
}

func (u *User) Insert(ctx context.Context, repo domain.Repo[*User]) (err error) {

	_, e := repo.Insert(ctx, u, func(s domain.Scanner) {
		t := *u
		err = s.Scan(&t.upk, &t.createdAt)
		if err == nil {
			*u = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (u *User) InsertArgs() []any {
	return []any{u.upk, u.createdAt}
}

func (u *User) InsertSQL() string {
	return `INSERT INTO users (upk, created_at) VALUES ($1, $2) RETURNING upk, created_at`
}

func (u *User) Key() string {
	return u.upk
}

func (u *User) String() string {
	return fmt.Sprintf(
		"{%s %v %v %v}\n",
		u.upk,
		u.deleted,
		u.createdAt,
		u.updatedAt,
	)
}

func (u *User) ToJSON() ([]byte, error) {

	result, err := json.Marshal(userJSON{
		UPK:       u.upk,
		Deleted:   tool.ConvertNullBoolToBoolPointer(u.deleted),
		CreatedAt: u.createdAt,
		UpdatedAt: tool.ConvertNullTimeToTimePointer(u.updatedAt),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *User) Update(ctx context.Context, repo domain.Repo[*User]) (err error) {

	_, e := repo.Update(ctx, u, func(s domain.Scanner) {
		t := *u
		err = s.Scan(&t.updatedAt)
		if err == nil {
			*u = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (u *User) UpdateArgs() []any {
	return []any{u.upk, u.updatedAt}
}

func (u *User) UpdateSQL() string {
	return `UPDATE users SET updated_at = $2 WHERE upk = $1 RETURNING updated_at`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
