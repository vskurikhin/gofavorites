/*
 * This file was last modified at 2024-07-15 16:32 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * entity.go
 * $Id$
 */
//!+

// Package domain TODO.
package domain

type Entity interface {
	DeleteArgs() []any
	DeleteSQL() string
	GetArgs() []any
	GetSQL() string
	InsertArgs() []any
	InsertSQL() string
	JSON() ([]byte, error)
	Key() string
	UpdateArgs() []any
	UpdateSQL() string
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
