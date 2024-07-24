/*
 * This file was last modified at 2024-07-24 09:10 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * user.go
 * $Id$
 */
//!+

// Package dto TODO.
package dto

type SignInRequest struct {
	UserName string `json:"user_name"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
