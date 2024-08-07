/*
 * This file was last modified at 2024-08-03 13:01 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites_deleted.go
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

	"github.com/google/uuid"
	"github.com/vskurikhin/gofavorites/internal/domain"
)

const (
	FavoritesDeletedDeleteForUserSQL = `UPDATE favorites
	SET deleted = true, version = (SELECT u.version FROM users u WHERE u.upk = $2)
	WHERE user_upk = $1`

	FavoritesDeletedSelectForUserSQL = `SELECT
	f.id, f.version, f.deleted, f.created_at, f.updated_at,
    a.isin, a.deleted, a.created_at, a.updated_at,
    t.name, t.deleted, t.created_at, t.updated_at,
    u.upk, u.version, u.deleted, u.created_at, u.updated_at
    FROM favorites f
    JOIN assets a ON f.isin = a.isin
    JOIN asset_types t ON a.asset_type = t.name 
    JOIN users u ON f.user_upk = u.upk
    WHERE f.user_upk = $1 AND f.version IS NULL
	AND f.deleted IS TRUE`
)

type FavoritesDeleted struct {
	TAttributes
	id      uuid.UUID
	asset   Asset
	user    User
	version sql.NullInt64
}

var _ domain.Entity = (*FavoritesDeleted)(nil)

func GetFavoritesDeletedForUser(ctx context.Context, repo domain.Repo[*FavoritesDeleted], upk string) ([]FavoritesDeleted, error) {

	var err error
	results := make([]FavoritesDeleted, 0)
	_, er0 := repo.GetByFilter(ctx, &FavoritesDeleted{user: User{upk: upk}}, func(scanner domain.Scanner) *FavoritesDeleted {
		result := FavoritesDeleted{}
		err = scanner.Scan(
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
			&result.user.version,
			&result.user.deleted,
			&result.user.createdAt,
			&result.user.updatedAt,
		)
		results = append(results, result)
		return &result
	})
	if er0 != nil {
		return results, err
	}
	return results, err
}

func MakeFavoritesDeletedUser(upk string) FavoritesDeleted {
	return FavoritesDeleted{
		user: User{upk: upk},
	}
}

func (f *FavoritesDeleted) Copy() domain.Entity {
	c := *f
	return &c
}

func (f *FavoritesDeleted) Delete(ctx context.Context, repo domain.Repo[*FavoritesDeleted]) (err error) {

	_, err = repo.Delete(ctx, f, func(s domain.Scanner) {})
	return err
}

func (f *FavoritesDeleted) DeleteArgs() []any {
	return []any{f.user.upk, f.user.upk}
}

func (f *FavoritesDeleted) DeleteSQL() string {
	return FavoritesDeletedDeleteForUserSQL
}

func (f *FavoritesDeleted) GetArgs() []any {
	return f.ToFavorites().GetArgs()
}

func (f *FavoritesDeleted) GetByFilterArgs() []any {
	return []any{f.user.upk}
}

func (f *FavoritesDeleted) GetByFilterSQL() string {
	return FavoritesDeletedSelectForUserSQL
}

func (f *FavoritesDeleted) GetSQL() string {
	return f.ToFavorites().GetSQL()
}

func (f *FavoritesDeleted) InsertArgs() []any {
	return f.ToFavorites().InsertArgs()
}

func (f *FavoritesDeleted) InsertSQL() string {
	return f.ToFavorites().InsertSQL()
}

func (f *FavoritesDeleted) FromJSON(data []byte) (err error) {
	return f.ToFavorites().FromJSON(data)
}

func (f *FavoritesDeleted) Key() string {
	return f.ToFavorites().Key()
}

func (f *FavoritesDeleted) ToJSON() ([]byte, error) {
	return f.ToFavorites().ToJSON()
}

func (f *FavoritesDeleted) Update(ctx context.Context, repo domain.Repo[*FavoritesDeleted]) (err error) {

	_, er0 := repo.Update(ctx, f, func(s domain.Scanner) {
		t := *f
		err = s.Scan(&f.version, &f.updatedAt)
		if err == nil {
			*f = t
		}
	})
	if er0 != nil {
		return er0
	}
	return err

}

func (f *FavoritesDeleted) UpdateArgs() []any {
	return f.ToFavorites().UpdateArgs()
}

func (f *FavoritesDeleted) UpdateSQL() string {
	return f.ToFavorites().UpdateSQL()
}

func (f *FavoritesDeleted) String() string {
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

func (f FavoritesDeleted) ToFavorites() *Favorites {

	return &Favorites{
		TAttributes: struct {
			deleted   sql.NullBool
			createdAt time.Time
			updatedAt sql.NullTime
		}{
			deleted:   f.deleted,
			createdAt: f.createdAt,
			updatedAt: f.updatedAt,
		},
		id:      f.id,
		asset:   f.asset,
		user:    f.user,
		version: f.version,
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
