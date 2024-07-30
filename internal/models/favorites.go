/*
 * This file was last modified at 2024-07-31 14:52 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites.go
 * $Id$
 */

package models

import (
	"database/sql"
	"math"

	"github.com/google/uuid"
	"github.com/ssoroka/slice"
	"github.com/vskurikhin/gofavorites/internal/controllers/dto"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	pb "github.com/vskurikhin/gofavorites/proto"
)

type Asset struct {
	isin      string
	assetType string
}

func (a Asset) Isin() string {
	return a.isin
}

func (a Asset) AssetType() string {
	return a.assetType
}

func (a Asset) ToEntity() entity.Asset {
	at := entity.MakeAssetType(a.assetType, entity.DefaultTAttributes())
	return entity.MakeAsset(a.isin, at, entity.DefaultTAttributes())
}

type Favorites struct {
	id      uuid.UUID
	asset   Asset
	user    User
	version int64
}

func AssetFromEntity(entity entity.Asset) Asset {
	return Asset{
		isin:      entity.Isin(),
		assetType: entity.AssetType().Name(),
	}
}

func (f Favorites) Id() uuid.UUID {
	return f.id
}

func (f Favorites) Asset() Asset {
	return f.asset
}

func (f Favorites) User() User {
	return f.user
}

func (f Favorites) Version() int64 {
	return f.version
}

func (f Favorites) WithUpk(upk string) Favorites {
	t := f
	t.user.upk = upk
	return t
}

func (f Favorites) ToDto() dto.Favorites {
	return dto.Favorites{
		ID:        f.id.String(),
		Isin:      f.asset.isin,
		AssetType: f.asset.assetType,
	}
}

func (f Favorites) ToEntity() entity.Favorites {

	var version sql.NullInt64

	if f.version > 0 {
		version.Int64 = f.version
		version.Valid = true
	}
	return entity.MakeFavorites(
		f.id, f.asset.ToEntity(), f.user.ToEntity(), version,
		entity.DefaultTAttributes(),
	)
}

func (f Favorites) ToProto() *pb.Favorites {
	return &pb.Favorites{
		Asset: &pb.Asset{
			Isin: f.asset.isin,
			AssetType: &pb.AssetType{
				Name: f.asset.assetType,
			},
		},
		User: &pb.User{
			PersonalKey: f.user.personalKey,
			Upk:         f.user.upk,
		},
	}
}

func FavoritesSliceToDto(favorites []Favorites) (result []dto.Favorites) {

	result = make([]dto.Favorites, 0, len(favorites))
	result = slice.Map[Favorites, dto.Favorites](
		favorites,
		func(i int, fav Favorites) dto.Favorites {
			return fav.ToDto()
		})
	return result
}

type User struct {
	personalKey string
	upk         string
	version     int64
}

func (u User) PersonalKey() string {
	return u.personalKey
}

func (u User) Upk() string {
	return u.upk
}

func (u User) Version() int64 {
	return u.version
}

func (u User) ToEntity() entity.User {
	return entity.MakeUser(u.upk, entity.DefaultTAttributes())
}

func FavoritesFromDto(dto dto.Favorites, personalKey, upk string) Favorites {

	assetType := dto.AssetType
	isin := dto.Isin
	asset := makeAsset(isin, assetType)
	user := MakeUser(personalKey, upk)

	return makeFavorites(uuid.Max, asset, user, math.MinInt64)
}

func FavoritesFromEntity(entity entity.Favorites) Favorites {

	asset := AssetFromEntity(entity.Asset())
	user := UserFromEntity(entity.User())

	return makeFavorites(entity.ID(), asset, user, entity.Version().Int64)
}

func FavoritesFromProto(proto *pb.Favorites) Favorites {

	if proto == nil {
		return Favorites{id: uuid.Max, version: math.MinInt64}
	}
	assetType := proto.GetAsset().GetAssetType().GetName()
	isin := proto.GetAsset().GetIsin()
	asset := makeAsset(isin, assetType)
	user := MakeUser(proto.GetUser().GetPersonalKey(), proto.GetUser().GetUpk())

	return makeFavorites(uuid.Max, asset, user, math.MinInt64)
}

func UserFromEntity(entity entity.User) User {
	return User{upk: entity.Upk(), version: entity.Version()}
}

func UserFromProto(proto *pb.User) User {

	if proto == nil {
		return User{}
	}
	return MakeUser(proto.GetPersonalKey(), proto.GetUpk())
}

func makeAsset(isin, assetType string) Asset {
	return Asset{
		isin:      isin,
		assetType: assetType,
	}
}

func makeFavorites(id uuid.UUID, asset Asset, user User, version int64) Favorites {
	return Favorites{
		id:      id,
		asset:   asset,
		user:    user,
		version: version,
	}
}

func MakeUser(personalKey, upk string) User {
	return User{
		personalKey: personalKey,
		upk:         upk,
	}
}
