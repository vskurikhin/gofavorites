/*
 * This file was last modified at 2024-07-25 13:49 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites.go
 * $Id$
 */
//!+

// Package dto TODO.
package dto

type Favorites struct {
	ID        string `json:"id"`
	Isin      string `json:"isin" validate:"required"`
	AssetType string `json:"asset_type" validate:"required"`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
