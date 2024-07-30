/*
 * This file was last modified at 2024-07-31 00:15 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * api_favorites_service.go
 * $Id$
 */
//!+

// Package services TODO.
package services

import (
	"context"

	"github.com/vskurikhin/gofavorites/internal/models"
)

type ApiFavoritesService interface {
	ApiFavoritesGet(ctx context.Context, model models.Favorites) (models.Favorites, error)
	ApiFavoritesGetForUser(ctx context.Context, model models.User) ([]models.Favorites, error)
	ApiFavoritesSet(ctx context.Context, model models.Favorites) (models.Favorites, error)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
