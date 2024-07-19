/*
 * This file was last modified at 2024-07-20 11:01 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * asset.go
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
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"log/slog"
	"time"
)

const (
	AssetSelectSQL = `SELECT
	a.isin, a.deleted, a.created_at, a.updated_at,
	t.name, t.deleted, t.created_at, t.updated_at
	FROM assets a
	JOIN asset_types t ON a.asset_type = t.name
	WHERE a.isin = $1`

	AssetSelectByAssetTypeSQL = `SELECT
	a.isin, a.deleted, a.created_at, a.updated_at,
	t.name, t.deleted, t.created_at, t.updated_at
	FROM assets a
	JOIN asset_types t ON a.asset_type = t.name
	WHERE a.asset_type = $1 AND a.deleted IS NOT TRUE AND t.deleted IS NOT TRUE`

	AssetDeleteSQL = `UPDATE assets
	SET deleted = true
	WHERE isin = $1
	RETURNING deleted, updated_at`

	AssetDeleteTxSQL = `UPDATE assets
	SET deleted = true WHERE isin = $1
	RETURNING isin, asset_type, deleted, created_at, updated_at`

	AssetInsertSQL = `INSERT INTO assets
	(isin, asset_type, created_at)
	VALUES ($1, $2, $3)
	RETURNING isin, asset_type, created_at`

	AssetUpdateSQL = `UPDATE assets
	SET asset_type = $2, updated_at = $3
	WHERE isin = $1
	RETURNING asset_type, updated_at`

	AssetUpsertTxAssetTypeSQL = `INSERT INTO asset_types
	(name, created_at)
	VALUES ($1, $2)
	ON CONFLICT (name)
	DO NOTHING`
)

type Asset struct {
	TAttributes
	isin      string
	assetType AssetType
}

var _ domain.Entity = (*Asset)(nil)

func GetAsset(ctx context.Context, repo domain.Repo[*Asset], isin string) (Asset, error) {

	var e error
	result := &Asset{isin: isin}

	result, err := repo.Get(ctx, result, func(scanner domain.Scanner) {
		e = scanner.Scan(
			&result.isin,
			&result.deleted,
			&result.createdAt,
			&result.updatedAt,
			&result.assetType.name,
			&result.assetType.deleted,
			&result.assetType.createdAt,
			&result.assetType.updatedAt,
		)
	})
	if e != nil {
		return Asset{}, e
	}
	if err != nil {
		return Asset{}, err
	}
	return *result, nil
}

func MakeAsset(isin string, at AssetType, a TAttributes) Asset {
	return Asset{
		TAttributes: struct {
			deleted   sql.NullBool
			createdAt time.Time
			updatedAt sql.NullTime
		}{
			deleted:   a.deleted,
			createdAt: a.createdAt,
			updatedAt: a.updatedAt,
		},
		isin:      isin,
		assetType: at,
	}
}

func IsAssetNotFound(a Asset, err error) bool {
	return tool.NoRowsInResultSet(err) || a == Asset{}
}

func (a *Asset) Isin() string {
	return a.isin
}

func (a *Asset) AssetType() AssetType {
	return a.assetType
}

func (a *Asset) Deleted() sql.NullBool {
	return a.deleted
}

func (a *Asset) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Asset) UpdatedAt() sql.NullTime {
	return a.updatedAt
}

// Copy shallow copy and return same type
func (a *Asset) Copy() domain.Entity {
	c := *a
	return &c
}

func (a *Asset) Delete(ctx context.Context, dtf domain.Dft[*Asset], inTransaction func()) (err error) {

	var isin, assetType string
	var deleted sql.NullBool
	var createdAt time.Time
	var updatedAt sql.NullTime

	e := dtf.DoDelete(ctx, a, func(scanner domain.Scanner) {
		err = scanner.Scan(
			&isin, &assetType, &deleted, &createdAt, &updatedAt,
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
		at := a.assetType
		*a = MakeAsset(isin, at, MakeTAttributes(deleted, createdAt, updatedAt))
	}
	return err
}

func (a *Asset) DeleteArgs() []any {
	return []any{a.isin}
}

func (a *Asset) DeleteSQL() string {
	return AssetDeleteSQL
}

func (a *Asset) DeleteTxArgs() domain.TxArgs {
	return domain.TxArgs{
		SQLs: []string{
			AssetDeleteTxSQL,
		},
		Args: [][]any{
			{a.isin},
		},
	}
}

type assetJSON struct {
	Isin      string
	AssetType assetTypeJSON
	Deleted   *bool      `json:",omitempty"`
	CreatedAt time.Time  `json:",omitempty"`
	UpdatedAt *time.Time `json:",omitempty"`
}

func (a *Asset) FromJSON(data []byte) (err error) {

	var t assetJSON
	err = json.Unmarshal(data, &t)

	if err != nil {
		return err
	}
	a.isin = t.Isin
	a.deleted = tool.ConvertBoolPointerToNullBool(t.Deleted)
	a.createdAt = t.CreatedAt
	a.updatedAt = tool.ConvertTimePointerToNullTime(t.UpdatedAt)

	a.assetType.name = t.AssetType.Name
	a.assetType.deleted = tool.ConvertBoolPointerToNullBool(t.AssetType.Deleted)
	a.assetType.createdAt = t.AssetType.CreatedAt
	a.assetType.updatedAt = tool.ConvertTimePointerToNullTime(t.AssetType.UpdatedAt)

	return nil
}

func (a *Asset) GetArgs() []any {
	return []any{a.isin}
}

func (a *Asset) GetByFilterArgs() []any {
	return []any{a.assetType.name}
}

func (a *Asset) GetByFilterSQL() string {
	return AssetSelectByAssetTypeSQL
}

func (a *Asset) GetSQL() string {
	return AssetSelectSQL
}

func (a *Asset) InsertArgs() []any {
	return []any{a.isin, a.assetType.name, a.createdAt}
}

func (a *Asset) InsertSQL() string {
	return AssetInsertSQL
}

func (a *Asset) Key() string {
	return a.isin
}

func (a *Asset) String() string {
	return fmt.Sprintf(
		"{%s {%s %v %v %v} %v %v %v}\n",
		a.isin,
		a.assetType.name,
		a.assetType.deleted,
		a.assetType.createdAt,
		a.assetType.updatedAt,
		a.deleted,
		a.createdAt,
		a.updatedAt,
	)
}

func (a *Asset) ToJSON() ([]byte, error) {

	result, err := json.Marshal(assetJSON{
		Isin: a.isin,
		AssetType: assetTypeJSON{
			Name:      a.assetType.name,
			Deleted:   tool.ConvertNullBoolToBoolPointer(a.assetType.deleted),
			CreatedAt: a.assetType.createdAt,
			UpdatedAt: tool.ConvertNullTimeToTimePointer(a.assetType.updatedAt),
		},
		Deleted:   tool.ConvertNullBoolToBoolPointer(a.deleted),
		CreatedAt: a.createdAt,
		UpdatedAt: tool.ConvertNullTimeToTimePointer(a.updatedAt),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *Asset) UpdateArgs() []any {
	return []any{a.isin, a.assetType.name, a.updatedAt}
}

func (a *Asset) UpdateSQL() string {
	return AssetUpdateSQL
}

const AssetUpsertTxAssetSQL = `INSERT INTO assets
	(isin, asset_type, created_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (isin)
	DO UPDATE SET asset_type = $2, updated_at = $4
	RETURNING isin, asset_type, deleted, created_at, updated_at,
	(SELECT created_at FROM asset_types WHERE name = $2)`

func (a *Asset) Upsert(ctx context.Context, dtf domain.Dft[*Asset], inTransaction func()) (err error) {

	var isin, assetType string
	var deleted sql.NullBool
	var createdAt, atCreatedAt time.Time
	var updatedAt sql.NullTime

	e := dtf.DoUpsert(ctx, a, func(scanner domain.Scanner) {
		err = scanner.Scan(
			&isin, &assetType, &deleted, &createdAt, &updatedAt, &atCreatedAt,
		)
		if err == nil {
			at := a.assetType
			*a = MakeAsset(isin, at, MakeTAttributes(deleted, createdAt, updatedAt))
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

func (a *Asset) UpsertTxArgs() domain.TxArgs {
	return domain.TxArgs{
		SQLs: []string{
			AssetUpsertTxAssetTypeSQL,
			AssetUpsertTxAssetSQL,
		},
		Args: [][]any{
			{a.assetType.name, a.assetType.createdAt},
			{a.isin, a.assetType.name, a.createdAt, a.updatedAt},
		},
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
