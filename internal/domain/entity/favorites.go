/*
 * This file was last modified at 2024-07-16 23:18 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites.go
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
	"github.com/google/uuid"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"time"
)

const (
	KeyFormat = "%32s%s"
)

type Favorites struct {
	id        uuid.UUID
	asset     Asset
	user      User
	version   sql.NullInt64
	deleted   sql.NullBool
	createdAt time.Time
	updatedAt sql.NullTime
}

var _ domain.Entity = (*Favorites)(nil)

func GetFavorites(ctx context.Context, repo domain.Repo[*Favorites], isin, upk string) (Favorites, error) {

	var e error
	result := &Favorites{asset: Asset{isin: isin}, user: User{upk: upk}}

	_, err := repo.Get(ctx, result, func(scanner domain.Scanner) {
		e = scanner.Scan(
			&result.id,
			&result.version,
			&result.deleted,
			&result.createdAt,
			&result.updatedAt,

			&result.asset.isin,
			&result.asset.deleted,
			&result.asset.createdAt,
			&result.asset.updatedAt,

			&result.asset.assetType.name,
			&result.asset.assetType.deleted,
			&result.asset.assetType.createdAt,
			&result.asset.assetType.updatedAt,

			&result.user.upk,
			&result.user.deleted,
			&result.user.createdAt,
			&result.user.updatedAt,
		)
	})
	if e != nil {
		return Favorites{}, e
	}
	if err != nil {
		return Favorites{}, err
	}
	return *result, nil
}

func NewFavorites(asset Asset, user User, createdAt time.Time) Favorites {
	return Favorites{
		id:        uuid.New(),
		asset:     asset,
		user:      user,
		createdAt: createdAt,
	}
}

func (f *Favorites) ID() uuid.UUID {
	return f.id
}

func (f *Favorites) Asset() Asset {
	return f.asset
}

func (f *Favorites) User() User {
	return f.user
}

func (f *Favorites) Version() sql.NullInt64 {
	return f.version
}

func (f *Favorites) Deleted() sql.NullBool {
	return f.deleted
}

func (f *Favorites) CreatedAt() time.Time {
	return f.createdAt
}

func (f *Favorites) UpdatedAt() sql.NullTime {
	return f.updatedAt
}

func (f *Favorites) Copy() domain.Entity {
	c := *f
	return &c
}

func (f *Favorites) Delete(ctx context.Context, repo domain.Repo[*Favorites]) (err error) {

	_, e := repo.Delete(ctx, f, func(s domain.Scanner) {
		t := *f
		err = s.Scan(&t.deleted, &t.updatedAt)
		if err == nil {
			*f = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (f *Favorites) DeleteArgs() []any {
	return []any{f.asset.isin, f.user.upk}
}

func (f *Favorites) DeleteSQL() string {
	return `UPDATE favorites
	SET deleted = true
	WHERE isin = $1 AND user_upk = $2
	RETURNING deleted, updated_at`
}

type favoritesJSON struct {
	ID        uuid.UUID
	Asset     assetJSON
	User      userJSON
	Version   int64
	Deleted   *bool `json:",omitempty"`
	CreatedAt time.Time
	UpdatedAt *time.Time `json:",omitempty"`
}

func (f *Favorites) FromJSON(data []byte) (err error) {

	var t favoritesJSON
	err = json.Unmarshal(data, &t)

	if err != nil {
		return err
	}
	f.id = t.ID
	f.deleted = tool.ConvertBoolPointerToNullBool(t.Deleted)
	f.createdAt = t.CreatedAt
	f.updatedAt = tool.ConvertTimePointerToNullTime(t.UpdatedAt)

	f.asset.isin = t.Asset.Isin
	f.asset.deleted = tool.ConvertBoolPointerToNullBool(t.Asset.Deleted)
	f.asset.createdAt = t.Asset.CreatedAt
	f.asset.updatedAt = tool.ConvertTimePointerToNullTime(t.Asset.UpdatedAt)

	f.asset.assetType.name = t.Asset.AssetType.Name
	f.asset.assetType.deleted = tool.ConvertBoolPointerToNullBool(t.Asset.AssetType.Deleted)
	f.asset.assetType.createdAt = t.Asset.AssetType.CreatedAt
	f.asset.assetType.updatedAt = tool.ConvertTimePointerToNullTime(t.Asset.AssetType.UpdatedAt)

	return nil
}

func (f *Favorites) GetArgs() []any {
	return []any{f.asset.isin, f.user.upk}
}

func (f *Favorites) GetSQL() string {
	return `SELECT f.id, f.version, f.deleted, f.created_at, f.updated_at,
                   a.isin, a.deleted, a.created_at, a.updated_at,
                   t.name, t.deleted, t.created_at, t.updated_at,
                   u.upk, u.deleted, u.created_at, u.updated_at
    FROM favorites f
    JOIN assets a ON f.isin = a.isin
    JOIN asset_types t ON a.asset_type = t.name 
    JOIN users u ON f.user_upk = u.upk 
    WHERE f.isin = $1 AND f.user_upk = $2`
}

func (f *Favorites) Insert(ctx context.Context, repo domain.Repo[*Favorites]) (err error) {

	_, e := repo.Insert(ctx, f, func(s domain.Scanner) {
		t := *f
		err = s.Scan(&t.id, &t.version, &t.createdAt)
		if err == nil {
			*f = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (f *Favorites) InsertArgs() []any {
	return []any{f.asset.isin, f.user.upk, f.version, f.createdAt}
}

func (f *Favorites) InsertSQL() string {
	return `INSERT INTO favorites
    (isin, user_upk, version, created_at)
    VALUES ($1, $2, $3, $4)
    RETURNING id, version, created_at`
}

func (f *Favorites) Key() string {
	return fmt.Sprintf(KeyFormat, f.asset.isin, f.user.upk)
}

func (f *Favorites) String() string {
	return fmt.Sprintf(
		"{%v {%s {%s %v %v %v} %v %v %v} {%s %v %v %v} %v %v %v %v}\n",
		f.id,
		f.asset.isin,
		f.asset.assetType.name,
		f.asset.assetType.deleted,
		f.asset.assetType.createdAt,
		f.asset.assetType.updatedAt,
		f.asset.deleted,
		f.asset.createdAt,
		f.asset.updatedAt,
		f.user.upk,
		f.user.deleted,
		f.user.createdAt,
		f.user.updatedAt,
		f.version,
		f.deleted,
		f.createdAt,
		f.updatedAt,
	)
}

func (f *Favorites) ToJSON() ([]byte, error) {

	result, err := json.Marshal(favoritesJSON{
		ID: f.id,
		Asset: assetJSON{
			Isin: f.asset.isin,
			AssetType: assetTypeJSON{
				Name:      f.asset.assetType.name,
				Deleted:   tool.ConvertNullBoolToBoolPointer(f.asset.assetType.deleted),
				CreatedAt: f.asset.assetType.createdAt,
				UpdatedAt: tool.ConvertNullTimeToTimePointer(f.asset.assetType.updatedAt),
			},
			Deleted:   tool.ConvertNullBoolToBoolPointer(f.asset.deleted),
			CreatedAt: f.asset.createdAt,
			UpdatedAt: tool.ConvertNullTimeToTimePointer(f.asset.updatedAt),
		},
		User: userJSON{
			UPK:       f.user.upk,
			Deleted:   tool.ConvertNullBoolToBoolPointer(f.user.deleted),
			CreatedAt: f.user.createdAt,
			UpdatedAt: tool.ConvertNullTimeToTimePointer(f.user.updatedAt),
		},
		Version:   f.version.Int64,
		Deleted:   tool.ConvertNullBoolToBoolPointer(f.deleted),
		CreatedAt: f.createdAt,
		UpdatedAt: tool.ConvertNullTimeToTimePointer(f.updatedAt),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *Favorites) Update(ctx context.Context, repo domain.Repo[*Favorites]) (err error) {

	_, e := repo.Update(ctx, f, func(s domain.Scanner) {
		t := *f
		err = s.Scan(&t.version, &t.updatedAt)
		if err == nil {
			*f = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (f *Favorites) UpdateArgs() []any {
	return []any{f.asset.isin, f.user.upk, f.version, f.updatedAt}
}

func (f *Favorites) UpdateSQL() string {
	return `UPDATE favorites
    SET version = $3, updated_at = $4
    WHERE isin = $1 AND user_upk = $2
    RETURNING version, updated_at`
}

func IsFavoritesNotFound(f Favorites, err error) bool {
	return f == Favorites{} || tool.NoRowsInResultSet(err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
