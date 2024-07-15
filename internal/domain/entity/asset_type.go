/*
 * This file was last modified at 2024-07-15 16:57 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * asset_type.go
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

type AssetType struct {
	name      string
	deleted   sql.NullBool
	createdAt time.Time
	updatedAt sql.NullTime
}

type assetTypeJSON struct {
	Name      string
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

var _ domain.Entity = (*AssetType)(nil)

func GetAssetTypes(ctx context.Context, repo domain.Repo, name string) (AssetType, error) {

	var e error
	var result AssetType

	err := repo.Get(ctx, &AssetType{name: name}, func(scanner domain.Scanner) {
		e = scanner.Scan(
			&result.name,
			&result.deleted,
			&result.createdAt,
			&result.updatedAt,
		)
	})
	if e != nil {
		return AssetType{}, e
	}
	if err != nil {
		return AssetType{}, err
	}
	return result, nil
}

func NewAssetTypes(name string, createdAt time.Time) AssetType {
	return AssetType{
		name:      name,
		createdAt: createdAt,
	}
}

func (a *AssetType) Delete(ctx context.Context, repo domain.Repo) (err error) {
	return repo.Delete(ctx, a)
}

func (a *AssetType) DeleteArgs() []any {
	return []any{a.name}
}

func (a *AssetType) DeleteSQL() string {
	return `UPDATE asset_types SET deleted = true WHERE name = $1`
}

func (a *AssetType) GetArgs() []any {
	return []any{a.name}
}

func (a *AssetType) GetSQL() string {
	return `SELECT name, deleted, created_at, updated_at FROM asset_types WHERE name = $1`
}

func (a *AssetType) Insert(ctx context.Context, repo domain.Repo) (err error) {
	return repo.Insert(ctx, a)
}

func (a *AssetType) InsertArgs() []any {
	return []any{a.name, a.createdAt}
}

func (a *AssetType) InsertSQL() string {
	return `INSERT INTO asset_types (name, created_at) VALUES ($1, $2)`
}

func (a *AssetType) JSON() ([]byte, error) {

	result, err := json.Marshal(assetTypeJSON{
		Name:      a.name,
		Deleted:   a.deleted.Bool,
		CreatedAt: a.createdAt,
		UpdatedAt: a.updatedAt,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *AssetType) Key() string {
	return a.name
}

func (a *AssetType) UpdateArgs() []any {
	return []any{a.name, a.updatedAt}
}

func (a *AssetType) UpdateSQL() string {
	return `UPDATE asset_types SET updatedAt = $2 WHERE name = $1`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
