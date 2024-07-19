/*
 * This file was last modified at 2024-07-20 11:01 by Victor N. Skurikhin.
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
	"fmt"
	"github.com/goccy/go-json"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"time"
)

type AssetType struct {
	TAttributes
	name string
}

var _ domain.Entity = (*AssetType)(nil)

func GetAssetType(ctx context.Context, repo domain.Repo[*AssetType], name string) (AssetType, error) {

	var e error
	result := &AssetType{name: name}

	result, err := repo.Get(ctx, result, func(scanner domain.Scanner) {
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
	return *result, nil
}

func MakeAssetType(name string, a TAttributes) AssetType {
	return AssetType{
		TAttributes: struct {
			deleted   sql.NullBool
			createdAt time.Time
			updatedAt sql.NullTime
		}{
			deleted:   a.deleted,
			createdAt: a.createdAt,
			updatedAt: a.updatedAt,
		},
		name: name,
	}
}

func (a *AssetType) Name() string {
	return a.name
}

func (a *AssetType) Deleted() sql.NullBool {
	return a.deleted
}

func (a *AssetType) CreatedAt() time.Time {
	return a.createdAt
}

func (a *AssetType) UpdatedAt() sql.NullTime {
	return a.updatedAt
}

func (a *AssetType) Copy() domain.Entity {
	c := *a
	return &c
}

func (a *AssetType) Delete(ctx context.Context, repo domain.Repo[*AssetType]) (err error) {

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

func (a *AssetType) DeleteArgs() []any {
	return []any{a.name}
}

func (a *AssetType) DeleteSQL() string {
	return `UPDATE asset_types SET deleted = true WHERE name = $1 RETURNING deleted, updated_at`
}

type assetTypeJSON struct {
	Name      string
	Deleted   *bool `json:",omitempty"`
	CreatedAt time.Time
	UpdatedAt *time.Time `json:",omitempty"`
}

func (a *AssetType) FromJSON(data []byte) (err error) {

	var t assetTypeJSON
	err = json.Unmarshal(data, &t)

	if err != nil {
		return err
	}
	a.name = t.Name
	a.deleted = tool.ConvertBoolPointerToNullBool(t.Deleted)
	a.createdAt = t.CreatedAt
	a.updatedAt = tool.ConvertTimePointerToNullTime(t.UpdatedAt)

	return nil
}

func (a *AssetType) GetArgs() []any {
	return []any{a.name}
}

func (a *AssetType) GetByFilterArgs() []any {
	return []any{}
}

func (a *AssetType) GetByFilterSQL() string {
	return `SELECT name, deleted, created_at, updated_at FROM asset_types WHERE deleted IS NOT TRUE`
}

func (a *AssetType) GetSQL() string {
	return `SELECT name, deleted, created_at, updated_at FROM asset_types WHERE name = $1`
}

func (a *AssetType) Insert(ctx context.Context, repo domain.Repo[*AssetType]) (err error) {

	_, e := repo.Insert(ctx, a, func(s domain.Scanner) {
		t := *a
		err = s.Scan(&t.name, &t.createdAt)
		if err == nil {
			*a = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (a *AssetType) InsertArgs() []any {
	return []any{a.name, a.createdAt}
}

func (a *AssetType) InsertSQL() string {
	return `INSERT INTO asset_types (name, created_at) VALUES ($1, $2) RETURNING name, created_at`
}

func (a *AssetType) Key() string {
	return a.name
}

func (a *AssetType) String() string {
	return fmt.Sprintf(
		"{%s %v %v %v}\n",
		a.name,
		a.deleted,
		a.createdAt,
		a.updatedAt,
	)
}

func (a *AssetType) ToJSON() ([]byte, error) {

	result, err := json.Marshal(assetTypeJSON{
		Name:      a.name,
		Deleted:   tool.ConvertNullBoolToBoolPointer(a.deleted),
		CreatedAt: a.createdAt,
		UpdatedAt: tool.ConvertNullTimeToTimePointer(a.updatedAt),
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *AssetType) Update(ctx context.Context, repo domain.Repo[*AssetType]) (err error) {

	_, e := repo.Update(ctx, a, func(s domain.Scanner) {
		t := *a
		err = s.Scan(&t.updatedAt)
		if err == nil {
			*a = t
		}
	})
	if e != nil {
		return e
	}
	return
}

func (a *AssetType) UpdateArgs() []any {
	return []any{a.name, a.updatedAt}
}

func (a *AssetType) UpdateSQL() string {
	return `UPDATE asset_types SET updated_at = $2 WHERE name = $1 RETURNING updated_at`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
