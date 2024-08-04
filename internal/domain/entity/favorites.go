/*
 * This file was last modified at 2024-08-04 20:14 by Victor N. Skurikhin.
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
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
)

const (
	KeyFormat = "%32s%s"

	FavoritesSelectSQL = `SELECT
	f.id, f.version, f.deleted, f.created_at, f.updated_at,
    a.isin, a.deleted, a.created_at, a.updated_at,
    t.name, t.deleted, t.created_at, t.updated_at,
    u.upk, u.version, u.deleted, u.created_at, u.updated_at
    FROM favorites f
    JOIN assets a ON f.isin = a.isin
    JOIN asset_types t ON a.asset_type = t.name 
    JOIN users u ON f.user_upk = u.upk 
    WHERE f.isin = $1 AND f.user_upk = $2`

	FavoritesSelectForUserSQL = `SELECT
	f.id, f.version, f.deleted, f.created_at, f.updated_at,
    a.isin, a.deleted, a.created_at, a.updated_at,
    t.name, t.deleted, t.created_at, t.updated_at,
    u.upk, u.version, u.deleted, u.created_at, u.updated_at
    FROM favorites f
    JOIN assets a ON f.isin = a.isin
    JOIN asset_types t ON a.asset_type = t.name 
    JOIN users u ON f.user_upk = u.upk
    WHERE f.user_upk = $1
	AND f.deleted IS NOT TRUE
	AND a.deleted IS NOT TRUE
	AND t.deleted IS NOT TRUE
	AND u.deleted IS NOT TRUE`

	FavoritesDeleteSQL = `UPDATE favorites
	SET deleted = true, version = (SELECT u.version FROM users u WHERE u.upk = $3)
	WHERE isin = $1 AND user_upk = $2
	RETURNING id, isin, user_upk, version, deleted, created_at, updated_at`

	FavoritesInsertSQL = `INSERT INTO favorites
    (isin, user_upk, version, created_at)
    VALUES ($1, $2, (SELECT u.version FROM users u WHERE u.upk = $3), $4)
    RETURNING id, version, created_at`

	FavoritesUpdateSQL = `UPDATE favorites
    SET updated_at = $3, version = (SELECT u.version FROM users u WHERE u.upk = $4)
    WHERE isin = $1 AND user_upk = $2
    RETURNING version, updated_at`

	FavoritesDeleteTxSQL = `UPDATE favorites
	SET version = version + 1, deleted = true
	WHERE isin = $1 AND user_upk = $2
	RETURNING id, isin, user_upk, version, deleted, created_at, updated_at`

	FavoritesUpsertTxAssetSQL = `INSERT INTO assets
	(isin, asset_type, created_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (isin)
	DO UPDATE SET asset_type = $2, updated_at = $4`

	FavoritesUpsertTxAssetTypeSQL = `INSERT INTO asset_types
	(name, created_at)
	VALUES ($1, $2)
	ON CONFLICT (name)
	DO NOTHING`

	FavoritesUpsertTxUserSQL = `INSERT INTO users
	(upk, version, created_at)
	VALUES ($1, 1, $2)
	ON CONFLICT (upk)
	DO UPDATE SET version = users.version + 1`

	FavoritesUpsertTxFavoritesSQL = `INSERT INTO favorites
    (isin, user_upk, created_at)
    VALUES ($1, $2, $3)
	ON CONFLICT (isin, user_upk)
	DO UPDATE SET updated_at = $4
    RETURNING id, isin, user_upk, version, deleted, created_at, updated_at,
	(SELECT created_at FROM asset_types WHERE name = $5),
	(SELECT created_at FROM assets WHERE isin = $1),
	(SELECT updated_at FROM assets WHERE isin = $1),
	(SELECT created_at FROM users WHERE upk = $2)`
)

type Favorites struct {
	TAttributes
	id      uuid.UUID
	asset   Asset
	user    User
	version sql.NullInt64
}

type favorites struct {
	ID        uuid.UUID
	Asset     asset
	User      user
	Version   int64
	Deleted   JsonNullBool `json:",omitempty"`
	CreatedAt time.Time
	UpdatedAt JsonNullTime `json:",omitempty"`
}

var _ domain.Entity = (*Favorites)(nil)

func GetFavorites(ctx context.Context, repo domain.Repo[*Favorites], isin, upk string) (Favorites, error) {

	var err error
	result := &Favorites{asset: Asset{isin: isin}, user: User{upk: upk}}

	_, er0 := repo.Get(ctx, result, func(scanner domain.Scanner) {
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
	})
	if er0 != nil {
		return Favorites{}, er0
	}
	if err != nil {
		return Favorites{}, err
	}
	return *result, nil
}

func GetFavoritesForUser(ctx context.Context, repo domain.Repo[*Favorites], upk string) ([]Favorites, error) {

	var err error
	results := make([]Favorites, 0)
	_, er0 := repo.GetByFilter(ctx, &Favorites{user: User{upk: upk}}, func(scanner domain.Scanner) *Favorites {
		result := Favorites{}
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
		return results, er0
	}
	return results, err
}

func MakeFavorites(id uuid.UUID, asset Asset, user User, version sql.NullInt64, a TAttributes) Favorites {
	return Favorites{
		TAttributes: struct {
			deleted   sql.NullBool
			createdAt time.Time
			updatedAt sql.NullTime
		}{
			deleted:   a.deleted,
			createdAt: a.createdAt,
			updatedAt: a.updatedAt,
		},
		id:      id,
		asset:   asset,
		user:    user,
		version: version,
	}
}

func IsFavoritesNotFound(f Favorites, err error) bool {
	return tool.NoRowsInResultSet(err) || f == Favorites{}
}

func (f Favorites) ID() uuid.UUID {
	return f.id
}

func (f Favorites) Asset() Asset {
	return f.asset
}

func (f Favorites) User() User {
	return f.user
}

func (f Favorites) Version() sql.NullInt64 {
	return f.version
}

func (f Favorites) Deleted() sql.NullBool {
	return f.deleted
}

func (f Favorites) CreatedAt() time.Time {
	return f.createdAt
}

func (f Favorites) UpdatedAt() sql.NullTime {
	return f.updatedAt
}

func (f *Favorites) Copy() domain.Entity {
	c := *f
	return &c
}

func (f *Favorites) Delete(ctx context.Context, dtf domain.Dft[*Favorites], inTransaction func()) (err error) {

	var fv Favorites

	er0 := dtf.DoDelete(ctx, f, func(scanner domain.Scanner) {
		err = scanner.Scan(
			&fv.id, &fv.asset.isin, &fv.user.upk, &fv.version, &fv.deleted, &fv.createdAt, &fv.updatedAt,
		)
		if err != nil {
			slog.ErrorContext(ctx, env.MSG+"Favorites.Delete", "err", err)
		} else {
			f.id = fv.id
			f.version = fv.version
			f.deleted = fv.deleted
			f.createdAt = fv.createdAt
			f.updatedAt = fv.updatedAt
			f.asset.isin = fv.asset.isin
			f.user.upk = fv.user.upk
			inTransaction()
		}
	})
	if er0 != nil {
		return er0
	}
	return err
}

func (f *Favorites) DeleteArgs() []any {
	return []any{f.asset.isin, f.user.upk, f.user.upk}
}

func (f *Favorites) DeleteSQL() string {
	return FavoritesDeleteSQL
}

func (f *Favorites) DeleteTxArgs() domain.TxArgs {
	return domain.TxArgs{
		SQLs: []string{
			FavoritesDeleteTxSQL,
		},
		Args: [][]any{
			{f.asset.isin, f.user.upk},
		},
	}
}

func (f *Favorites) FromJSON(data []byte) (err error) {

	var t favorites
	err = json.Unmarshal(data, &t)

	if err != nil {
		return err
	}
	f.id = t.ID
	f.deleted = t.Deleted.ToNullBool()
	f.createdAt = t.CreatedAt
	f.updatedAt = t.UpdatedAt.ToNullTime()

	f.asset.isin = t.Asset.Isin
	f.asset.deleted = t.Asset.Deleted.ToNullBool()
	f.asset.createdAt = t.Asset.CreatedAt
	f.asset.updatedAt = t.Asset.UpdatedAt.ToNullTime()

	f.asset.assetType.name = t.Asset.AssetType.Name
	f.asset.assetType.deleted = t.Asset.AssetType.Deleted.ToNullBool()
	f.asset.assetType.createdAt = t.Asset.AssetType.CreatedAt
	f.asset.assetType.updatedAt = t.Asset.AssetType.UpdatedAt.ToNullTime()

	f.user.upk = t.User.UPK
	f.user.version = t.User.Version
	f.user.deleted = t.User.Deleted.ToNullBool()
	f.user.createdAt = t.User.CreatedAt
	f.user.updatedAt = t.User.UpdatedAt.ToNullTime()

	return nil
}

func (f *Favorites) GetArgs() []any {
	return []any{f.asset.isin, f.user.upk}
}

func (f *Favorites) GetByFilterArgs() []any {
	return []any{f.user.upk}
}

func (f *Favorites) GetByFilterSQL() string {
	return FavoritesSelectForUserSQL
}

func (f *Favorites) GetSQL() string {
	return FavoritesSelectSQL
}

func (f *Favorites) InsertArgs() []any {
	return []any{f.asset.isin, f.user.upk, f.user.upk, f.createdAt}
}

func (f *Favorites) InsertSQL() string {
	return FavoritesInsertSQL
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

func (f Favorites) ToJSON() ([]byte, error) {

	result, err := json.Marshal(favorites{
		ID: f.id,
		Asset: asset{
			Isin: f.asset.isin,
			AssetType: assetType{
				Name:      f.asset.assetType.name,
				Deleted:   FromNullBool(f.asset.assetType.deleted),
				CreatedAt: f.asset.assetType.createdAt,
				UpdatedAt: FromNullTime(f.asset.assetType.updatedAt),
			},
			Deleted:   FromNullBool(f.asset.deleted),
			CreatedAt: f.asset.createdAt,
			UpdatedAt: FromNullTime(f.asset.updatedAt),
		},
		User: user{
			UPK:       f.user.upk,
			Version:   f.user.version,
			Deleted:   FromNullBool(f.user.deleted),
			CreatedAt: f.user.createdAt,
			UpdatedAt: FromNullTime(f.user.updatedAt),
		},
		Version:   f.version.Int64,
		Deleted:   FromNullBool(f.deleted),
		CreatedAt: f.createdAt,
		UpdatedAt: FromNullTime(f.updatedAt),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (f *Favorites) Update(ctx context.Context, repo domain.Repo[*Favorites]) (err error) {

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

func (f *Favorites) UpdateArgs() []any {
	return []any{f.asset.isin, f.user.upk, f.updatedAt, f.user.upk}
}

func (f *Favorites) UpdateSQL() string {
	return FavoritesUpdateSQL
}

func (f *Favorites) Upsert(ctx context.Context, dtf domain.Dft[*Favorites], inTransaction func()) (err error) {

	var fv Favorites
	er0 := dtf.DoUpsert(ctx, f, func(scanner domain.Scanner) {
		err = scanner.Scan(
			&fv.id, &fv.asset.isin, &fv.user.upk, &fv.version, &fv.deleted, &fv.createdAt, &fv.updatedAt,
			&fv.asset.assetType.createdAt, &fv.asset.createdAt, &fv.asset.updatedAt, &fv.asset.createdAt,
		)
		if err != nil {
			slog.ErrorContext(ctx, env.MSG+"Favorites.Upsert", "err", err)
		} else {
			f.id = fv.id
			f.version = fv.version
			f.deleted = fv.deleted
			f.createdAt = fv.createdAt
			f.updatedAt = fv.updatedAt
			f.asset.isin = fv.asset.isin
			f.asset.createdAt = fv.asset.createdAt
			f.asset.updatedAt = fv.asset.updatedAt
			f.asset.assetType.createdAt = fv.asset.assetType.createdAt
			f.user.upk = fv.user.upk
			f.user.createdAt = fv.asset.createdAt
			inTransaction()
		}
	})
	if er0 != nil {
		return er0
	}
	return err
}

func (f *Favorites) UpsertTxArgs() domain.TxArgs {
	return domain.TxArgs{
		SQLs: []string{
			FavoritesUpsertTxAssetTypeSQL,
			FavoritesUpsertTxAssetSQL,
			FavoritesUpsertTxUserSQL,
			FavoritesUpsertTxFavoritesSQL,
		},
		Args: [][]any{
			{f.asset.assetType.name, f.asset.assetType.createdAt},
			{f.asset.isin, f.asset.assetType.name, f.asset.createdAt, f.asset.updatedAt},
			{f.user.upk, f.user.createdAt},
			{f.asset.isin, f.user.upk, f.createdAt, f.updatedAt, f.asset.assetType.name},
		},
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
