/*
 * This file was last modified at 2024-07-19 15:41 by Victor N. Skurikhin.
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
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
	pb "github.com/vskurikhin/gofavorites/proto"
	"log/slog"
	"time"
)

const (
	KeyFormat = "%32s%s"

	FavoritesSelectSQL = `SELECT
	f.id, f.version, f.deleted, f.created_at, f.updated_at,
    a.isin, a.deleted, a.created_at, a.updated_at,
    t.name, t.deleted, t.created_at, t.updated_at,
    u.upk, u.deleted, u.created_at, u.updated_at
    FROM favorites f
    JOIN assets a ON f.isin = a.isin
    JOIN asset_types t ON a.asset_type = t.name 
    JOIN users u ON f.user_upk = u.upk 
    WHERE f.isin = $1 AND f.user_upk = $2`

	FavoritesDeleteSQL = `UPDATE favorites
	SET deleted = true
	WHERE isin = $1 AND user_upk = $2
	RETURNING id, isin, user_upk, version, deleted, created_at, updated_at`

	FavoritesInsertSQL = `INSERT INTO favorites
    (isin, user_upk, version, created_at)
    VALUES ($1, $2, $3, $4)
    RETURNING id, version, created_at`

	FavoritesUpdateSQL = `UPDATE favorites
    SET version = $3, updated_at = $4
    WHERE isin = $1 AND user_upk = $2
    RETURNING version, updated_at`

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
	(upk, created_at)
	VALUES ($1, $2)
	ON CONFLICT (upk)
	DO NOTHING`
)

type Favorites struct {
	TAttributes
	id      uuid.UUID
	asset   Asset
	user    User
	version sql.NullInt64
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

func FavoritesFromProto(proto *pb.Favorites, upk string) Favorites {

	assetType := proto.GetAsset().GetAssetType().GetName()
	isin := proto.GetAsset().GetIsin()
	id := uuid.New()
	at := MakeAssetType(assetType, DefaultTAttributes())
	asset := MakeAsset(isin, at, DefaultTAttributes())
	user := MakeUser(upk, DefaultTAttributes())

	return MakeFavorites(id, asset, user, sql.NullInt64{}, DefaultTAttributes())
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

const FavoritesDeleteTxSQL = `UPDATE favorites
	SET deleted = true
	WHERE isin = $1 AND user_upk = $2
	RETURNING id, isin, user_upk, version, deleted, created_at, updated_at`

func (f *Favorites) Delete(ctx context.Context, dtf domain.Dft[*Favorites], inTransaction func()) (err error) {

	var id uuid.UUID
	var isin, userUpk string
	var version sql.NullInt64
	var deleted sql.NullBool
	var createdAt time.Time
	var updatedAt sql.NullTime

	e := dtf.DoDelete(ctx, f, func(scanner domain.Scanner) {
		err = scanner.Scan(
			&id, &isin, &userUpk, &version, &deleted, &createdAt, &updatedAt,
		)
		if err == nil {
			inTransaction()
		} else {
			slog.Error(env.MSG+" Delete", "err", err)
		}
	})
	if e != nil {
		return e
	}
	if err == nil {
		asset := f.asset
		user := f.user
		*f = MakeFavorites(id, asset, user, version, MakeTAttributes(deleted, createdAt, updatedAt))
	}
	return err
}

func (f *Favorites) DeleteArgs() []any {
	return []any{f.asset.isin, f.user.upk}
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

	f.user.upk = t.User.UPK
	f.user.deleted = tool.ConvertBoolPointerToNullBool(t.User.Deleted)
	f.user.createdAt = t.User.CreatedAt
	f.user.updatedAt = tool.ConvertTimePointerToNullTime(t.User.UpdatedAt)

	return nil
}

func (f *Favorites) GetArgs() []any {
	return []any{f.asset.isin, f.user.upk}
}

func (f *Favorites) GetSQL() string {
	return FavoritesSelectSQL
}

func (f *Favorites) InsertArgs() []any {
	return []any{f.asset.isin, f.user.upk, f.version, f.createdAt}
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

func (f *Favorites) ToProto() pb.Favorites {
	return pb.Favorites{
		Asset: &pb.Asset{
			Isin: f.asset.isin,
			AssetType: &pb.AssetType{
				Name: f.asset.assetType.name,
			},
		},
		User: &pb.User{
			Upk: f.user.upk,
		},
	}
}

func (f *Favorites) UpdateArgs() []any {
	return []any{f.asset.isin, f.user.upk, f.version, f.updatedAt}
}

func (f *Favorites) UpdateSQL() string {
	return FavoritesUpdateSQL
}

const FavoritesUpsertTxFavoritesSQL = `INSERT INTO favorites
    (isin, user_upk, version, created_at)
    VALUES ($1, $2, $3, $4)
	ON CONFLICT (isin, user_upk)
	DO UPDATE SET version = $3, updated_at = $5
    RETURNING id, isin, user_upk, version, deleted, created_at, updated_at,
	(SELECT created_at FROM asset_types WHERE name = $6),
	(SELECT created_at FROM assets WHERE isin = $1),
	(SELECT updated_at FROM assets WHERE isin = $1),
	(SELECT created_at FROM users WHERE upk = $2)`

func (f *Favorites) Upsert(ctx context.Context, dtf domain.Dft[*Favorites], inTransaction func()) (err error) {

	var id uuid.UUID
	var isin, userUpk string
	var version sql.NullInt64
	var deleted sql.NullBool
	var createdAt, assetsCreatedAt, assetTypesCreatedAt, usersCreatedAt time.Time
	var updatedAt, assetsUpdatedAt sql.NullTime

	e := dtf.DoUpsert(ctx, f, func(scanner domain.Scanner) {
		err = scanner.Scan(
			&id, &isin, &userUpk, &version, &deleted, &createdAt, &updatedAt,
			&assetTypesCreatedAt, &assetsCreatedAt, &assetsUpdatedAt, &usersCreatedAt,
		)
		if err == nil {
			at := MakeAssetType(
				f.asset.assetType.name,
				MakeTAttributes(f.asset.assetType.deleted, assetTypesCreatedAt, f.asset.assetType.updatedAt),
			)
			asset := MakeAsset(isin, at, MakeTAttributes(f.asset.deleted, assetsCreatedAt, assetsUpdatedAt))
			user := MakeUser(userUpk, MakeTAttributes(f.user.deleted, usersCreatedAt, f.user.updatedAt))
			*f = MakeFavorites(id, asset, user, version, MakeTAttributes(deleted, createdAt, updatedAt))
			inTransaction()
		} else {
			slog.Error(env.MSG+" Upsert", "err", err)
		}
	})
	if e != nil {
		return e
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
			{f.asset.isin, f.user.upk, f.version, f.createdAt, f.updatedAt, f.asset.assetType.name},
		},
	}
}

func IsFavoritesNotFound(f Favorites, err error) bool {
	return f == Favorites{} || tool.NoRowsInResultSet(err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
