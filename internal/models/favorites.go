/*
 * This file was last modified at 2024-07-26 11:20 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites.go
 * $Id$
 */

package models

import (
	"github.com/google/uuid"
	"github.com/vskurikhin/gofavorites/internal/controllers/dto"
	pb "github.com/vskurikhin/gofavorites/proto"
	"math"
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

type Favorites struct {
	id      uuid.UUID
	asset   Asset
	user    User
	version int64
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

func (f Favorites) ToDto() dto.Favorites {
	return dto.Favorites{
		ID:        f.id.String(),
		Isin:      f.asset.isin,
		AssetType: f.asset.assetType,
	}
}

func FavoritesSliceToDto(favorites []Favorites) (result []dto.Favorites) {

	result = make([]dto.Favorites, 0, len(favorites))

	for _, fav := range favorites {
		item := fav.ToDto()
		result = append(result, item)
	}
	return result
}

type User struct {
	personalKey string
	upk         string
}

func (u User) PersonalKey() string {
	return u.personalKey
}

func (u User) Upk() string {
	return u.upk
}

func FavoritesFromDto(dto dto.Favorites, personalKey, upk string) Favorites {

	assetType := dto.AssetType
	isin := dto.Isin
	asset := MakeAsset(isin, assetType)
	user := MakeUser(personalKey, upk)

	return MakeFavorites(uuid.Max, asset, user, math.MaxInt64)
}

func FavoritesFromProto(proto *pb.Favorites) Favorites {

	if proto == nil {
		return Favorites{id: uuid.Max, version: math.MaxInt64}
	}
	assetType := proto.GetAsset().GetAssetType().GetName()
	isin := proto.GetAsset().GetIsin()
	asset := MakeAsset(isin, assetType)
	user := MakeUser(proto.GetUser().GetPersonalKey(), proto.GetUser().GetUpk())

	return MakeFavorites(uuid.Max, asset, user, math.MaxInt64)
}

func UserFromProto(proto *pb.User) User {

	if proto == nil {
		return User{}
	}
	return MakeUser(proto.GetPersonalKey(), proto.GetUpk())
}

func MakeAsset(isin, assetType string) Asset {
	return Asset{
		isin:      isin,
		assetType: assetType,
	}
}

func MakeFavorites(id uuid.UUID, asset Asset, user User, version int64) Favorites {
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
