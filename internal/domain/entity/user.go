/*
 * This file was last modified at 2024-08-03 17:21 by Victor N. Skurikhin.
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
	"time"

	"github.com/goccy/go-json"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/tool"
)

type User struct {
	TAttributes
	upk     string
	version int64
}

type user struct {
	UPK       string
	Version   int64
	Deleted   JsonNullBool `json:",omitempty"`
	CreatedAt time.Time
	UpdatedAt JsonNullTime `json:",omitempty"`
}

var _ domain.Entity = (*User)(nil)

func GetUser(ctx context.Context, repo domain.Repo[*User], upk string) (User, error) {

	var err error
	result := &User{upk: upk}

	result, er0 := repo.Get(ctx, result, func(scanner domain.Scanner) {
		err = scanner.Scan(
			&result.upk,
			&result.version,
			&result.deleted,
			&result.createdAt,
			&result.updatedAt,
		)
	})
	if er0 != nil {
		return User{}, er0
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

func MakeUserWithVersion(upk string, version int64, a TAttributes) User {
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
		upk:     upk,
		version: version,
	}
}

func IsUserNotFound(u User, err error) bool {
	return tool.NoRowsInResultSet(err) || u == User{}
}

func (u User) Upk() string {
	return u.upk
}

func (u User) Version() int64 {
	return u.version
}

func (u User) Deleted() sql.NullBool {
	return u.deleted
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() sql.NullTime {
	return u.updatedAt
}

func (u *User) Copy() domain.Entity {
	c := *u
	return &c
}

func (u *User) Delete(ctx context.Context, repo domain.Repo[*User]) (err error) {

	_, e := repo.Delete(ctx, u, func(s domain.Scanner) {
		t := *u
		err = s.Scan(&t.version, &t.deleted, &t.updatedAt)
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
	return `UPDATE users
	SET version = version + 1, deleted = true
	WHERE upk = $1
	RETURNING version, deleted, updated_at`
}

func (u *User) FromJSON(data []byte) (err error) {

	var t user
	err = json.Unmarshal(data, &t)

	if err != nil {
		return err
	}
	u.upk = t.UPK
	u.version = t.Version
	u.deleted = t.Deleted.ToNullBool()
	u.createdAt = t.CreatedAt
	u.updatedAt = t.UpdatedAt.ToNullTime()

	return nil
}

func (u *User) GetArgs() []any {
	return []any{u.upk}
}

func (u *User) GetByFilterArgs() []any {
	return []any{}
}

func (u *User) GetByFilterSQL() string {
	return `SELECT upk, version, deleted, created_at, updated_at FROM users WHERE deleted IS NOT TRUE`
}

func (u *User) GetSQL() string {
	return `SELECT upk, version, deleted, created_at, updated_at FROM users WHERE upk = $1`
}

func (u *User) Insert(ctx context.Context, repo domain.Repo[*User]) (err error) {

	_, e := repo.Insert(ctx, u, func(s domain.Scanner) {
		t := *u
		err = s.Scan(&t.upk, &t.version, &t.createdAt)
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
	return []any{u.upk, u.version, u.createdAt}
}

func (u *User) InsertSQL() string {
	return `INSERT INTO users
	(upk, version, created_at)
	VALUES ($1, $2, $3)
	RETURNING upk, version, created_at`
}

func (u *User) Key() string {
	return u.upk
}

func (u *User) String() string {
	return fmt.Sprintf(
		"{%s %d %v %v %v}\n",
		u.upk,
		u.version,
		u.deleted,
		u.createdAt,
		u.updatedAt,
	)
}

func (u *User) ToJSON() ([]byte, error) {

	result, err := json.Marshal(user{
		UPK:       u.upk,
		Version:   u.version,
		Deleted:   FromNullBool(u.deleted),
		CreatedAt: u.createdAt,
		UpdatedAt: FromNullTime(u.updatedAt),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *User) Update(ctx context.Context, repo domain.Repo[*User]) (err error) {

	_, er0 := repo.Update(ctx, u, func(s domain.Scanner) {
		t := *u
		err = s.Scan(&t.version, &t.updatedAt)
		if err == nil {
			*u = t
		}
	})
	if er0 != nil {
		return er0
	}
	return err
}

func (u *User) UpdateArgs() []any {
	return []any{u.upk, u.version, u.updatedAt}
}

func (u *User) UpdateSQL() string {
	return `UPDATE users SET version = $2, updated_at = $3 WHERE upk = $1 RETURNING version, updated_at`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
