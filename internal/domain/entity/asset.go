/*
 * This file was last modified at 2024-07-16 20:57 by Victor N. Skurikhin.
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
	"github.com/goccy/go-json"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"time"
)

type Asset struct {
	isin      string
	assetType AssetType
	deleted   sql.NullBool
	createdAt time.Time
	updatedAt sql.NullTime
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

func NewAsset(isin string, assetType AssetType, createdAt time.Time) Asset {
	return Asset{
		isin:      isin,
		assetType: assetType,
		createdAt: createdAt,
	}
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

// shallow copy and return same type
func (a *Asset) Copy() domain.Entity {
	c := *a
	return &c
}

func (a *Asset) Delete(ctx context.Context, repo domain.Repo[*Asset]) (err error) {

	_, e := repo.Delete(ctx, a, func(s domain.Scanner) {
		t := *a
		err = s.Scan(&t.deleted, &t.updatedAt)
		if err == nil {
			*a = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (a *Asset) DeleteArgs() []any {
	return []any{a.isin}
}

func (a *Asset) DeleteSQL() string {
	return `UPDATE assets SET deleted = true WHERE isin = $1 RETURNING deleted, updated_at`
}

type assetJSON struct {
	Isin      string
	AssetType assetTypeJSON
	Deleted   *bool
	CreatedAt time.Time
	UpdatedAt *time.Time
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
	a.assetType.deleted = tool.ConvertBoolPointerToNullBool(t.Deleted)
	a.assetType.createdAt = t.CreatedAt
	a.assetType.updatedAt = tool.ConvertTimePointerToNullTime(t.UpdatedAt)

	return nil
}

func (a *Asset) GetArgs() []any {
	return []any{a.isin}
}

func (a *Asset) GetSQL() string {
	return `SELECT a.isin, a.deleted, a.created_at, a.updated_at,
                   t.name, t.deleted, t.created_at, t.updated_at
	FROM assets a JOIN asset_types t ON a.asset_type = t.name WHERE a.isin = $1`
}

func (a *Asset) Insert(ctx context.Context, repo domain.Repo[*Asset]) (err error) {

	_, e := repo.Insert(ctx, a, func(s domain.Scanner) {
		t := *a
		err = s.Scan(&t.isin, &t.assetType.name, &t.createdAt)
		if err == nil {
			*a = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (a *Asset) InsertArgs() []any {
	return []any{a.isin, a.assetType.name, a.createdAt}
}

func (a *Asset) InsertSQL() string {
	return `INSERT INTO assets
	(isin, asset_type, created_at)
	VALUES ($1, $2, $3)
	RETURNING isin, asset_type, created_at`
}

func (a *Asset) Key() string {
	return a.isin
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

func (a *Asset) Update(ctx context.Context, repo domain.Repo[*Asset]) (err error) {

	_, e := repo.Update(ctx, a, func(s domain.Scanner) {
		t := *a
		err = s.Scan(&t.assetType.name, &t.updatedAt)
		if err == nil {
			*a = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (a *Asset) UpdateArgs() []any {
	return []any{a.isin, a.assetType.name, a.updatedAt}
}

func (a *Asset) UpdateSQL() string {
	return `UPDATE assets
	SET asset_type = $2, updated_at = $3
	WHERE isin = $1
	RETURNING asset_type, updated_at`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
