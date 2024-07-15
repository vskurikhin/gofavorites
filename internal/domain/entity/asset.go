/*
 * This file was last modified at 2024-07-15 17:53 by Victor N. Skurikhin.
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
	"time"
)

type Asset struct {
	isin      string
	assetType AssetType
	deleted   sql.NullBool
	createdAt time.Time
	updatedAt sql.NullTime
}

type assetJSON struct {
	Isin      string
	AssetType assetTypeJSON
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

var _ domain.Entity = (*Asset)(nil)

func GetAsset(ctx context.Context, repo domain.Repo, isin string) (Asset, error) {

	var e error
	var result Asset

	err := repo.Get(ctx, &Asset{isin: isin}, func(scanner domain.Scanner) {
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
	return result, nil
}

func NewAsset(isin, assetType string, createdAt time.Time) Asset {
	return Asset{
		isin:      isin,
		assetType: NewAssetTypes(assetType, createdAt),
		createdAt: createdAt,
	}
}

func (a *Asset) Delete(ctx context.Context, repo domain.Repo) (err error) {
	return repo.Delete(ctx, a)
}

func (a *Asset) DeleteArgs() []any {
	return []any{a.isin}
}

func (a *Asset) DeleteSQL() string {
	return `UPDATE assets SET deleted = true WHERE name = $1`
}

func (a *Asset) GetArgs() []any {
	return []any{a.isin}
}

func (a *Asset) GetSQL() string {
	return `SELECT a.isin, a.deleted, a.created_at, a.updated_at,
                   t.name, t.deleted, t.created_at, t.updated_at
	FROM assets a JOIN asset_types t ON a.asset_type = t.name WHERE a.isin = $1`
}

func (a *Asset) Insert(ctx context.Context, repo domain.Repo) (err error) {
	return repo.Insert(ctx, a)
}

func (a *Asset) InsertArgs() []any {
	return []any{a.isin, a.assetType.name, a.createdAt}
}

func (a *Asset) InsertSQL() string {
	return `INSERT INTO asset (isin, asset_type, created_at) VALUES ($1, $2, $3)`
}

func (a *Asset) JSON() ([]byte, error) {

	result, err := json.Marshal(assetJSON{
		Isin: a.isin,
		AssetType: assetTypeJSON{
			Name:      a.assetType.name,
			Deleted:   a.assetType.deleted.Bool,
			CreatedAt: a.assetType.createdAt,
			UpdatedAt: a.assetType.updatedAt,
		},
		Deleted:   a.deleted.Bool,
		CreatedAt: a.createdAt,
		UpdatedAt: a.updatedAt,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *Asset) Key() string {
	return a.isin
}

func (a *Asset) UpdateArgs() []any {
	return []any{a.isin, a.updatedAt, a.assetType.name}
}

func (a *Asset) UpdateSQL() string {
	return `UPDATE assets SET updatedAt = $2, asset_type = $3 WHERE name = $1`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
