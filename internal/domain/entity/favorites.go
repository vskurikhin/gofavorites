/*
 * This file was last modified at 2024-07-15 17:59 by Victor N. Skurikhin.
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

type favoritesJSON struct {
	Id        uuid.UUID
	Asset     assetJSON
	User      userJSON
	Version   int64
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

var _ domain.Entity = (*Favorites)(nil)

func GetFavorites(ctx context.Context, repo domain.Repo, isin, upk string) (Favorites, error) {

	var e error
	var result Favorites

	err := repo.Get(ctx, &Favorites{asset: Asset{isin: isin}, user: User{upk: upk}}, func(scanner domain.Scanner) {
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
	return result, nil
}

func NewFavorites(isin, assetType string, user User, createdAt time.Time) Favorites {
	return Favorites{
		id:        uuid.New(),
		asset:     NewAsset(isin, assetType, createdAt),
		user:      user,
		createdAt: createdAt,
	}
}

func (f *Favorites) Delete(ctx context.Context, repo domain.Repo) (err error) {
	return repo.Delete(ctx, f)
}

func (f *Favorites) DeleteArgs() []any {
	return []any{f.asset.isin, f.user.upk}
}

func (f *Favorites) DeleteSQL() string {
	return `UPDATE favorites SET deleted = true WHERE isin = $1 AND user_upk = $2`
}

func (f *Favorites) GetArgs() []any {
	return []any{f.asset.isin, f.user.upk}
}

func (f *Favorites) GetSQL() string {
	return `SELECT f.id, f.version, f.deleted, f.created_at, f.updated_at,
                   a.isin, a.deleted, a.created_at, a.updated_at,
                   t.name, t.deleted, t.created_at, t.updated_at,
                   u.upk, u.deleted, u.created_at, u.updated_at,
    FROM favorites f
    JOIN assets a ON f.isin = a.isin
    JOIN asset_types t ON a.asset_type = t.name 
    JOIN users u ON f.user_upk = u.upk 
    WHERE f.isin = $1 AND u.upk = $2`
}

func (f *Favorites) Insert(ctx context.Context, repo domain.Repo) (err error) {
	return repo.Insert(ctx, f)
}

func (f *Favorites) InsertArgs() []any {
	return []any{f.asset.isin, f.user.upk, f.version, f.createdAt}
}

func (f *Favorites) InsertSQL() string {
	return `INSERT INTO favorites (isin, user_upk, version, created_at) VALUES ($1, $2, $3, $4)`
}

func (f *Favorites) JSON() ([]byte, error) {

	result, err := json.Marshal(favoritesJSON{
		Id: f.id,
		Asset: assetJSON{
			Isin: f.asset.isin,
			AssetType: assetTypeJSON{
				Name:      f.asset.assetType.name,
				Deleted:   f.asset.assetType.deleted.Bool,
				CreatedAt: f.asset.assetType.createdAt,
				UpdatedAt: f.asset.assetType.updatedAt,
			},
			Deleted:   f.asset.deleted.Bool,
			CreatedAt: f.asset.createdAt,
			UpdatedAt: f.asset.updatedAt,
		},
		User: userJSON{
			UPK:       f.user.upk,
			Deleted:   f.user.deleted.Bool,
			CreatedAt: f.user.createdAt,
			UpdatedAt: f.user.updatedAt,
		},
		Version:   f.version.Int64,
		Deleted:   f.deleted.Bool,
		CreatedAt: f.createdAt,
		UpdatedAt: f.updatedAt,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *Favorites) Key() string {
	return fmt.Sprintf(KeyFormat, f.asset.isin, f.user.upk)
}

func (f *Favorites) UpdateArgs() []any {
	return []any{f.asset.isin, f.user.upk, f.updatedAt, f.version}
}

func (f *Favorites) UpdateSQL() string {
	return `UPDATE favorites SET updatedAt = $3, version = $4 WHERE isin = $1 AND user_upk = $2`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
