/*
 * This file was last modified at 2024-07-16 20:57 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * entity.go
 * $Id$
 */
//!+

// Package domain TODO.
package domain

type Cloneable interface {
	Copy() Entity
}

type Serializable interface {
	FromJSON(data []byte) (err error)
	Key() string
	ToJSON() ([]byte, error)
}

type Entity interface {
	Cloneable
	Serializable
	DeleteArgs() []any
	DeleteSQL() string
	GetArgs() []any
	GetSQL() string
	InsertArgs() []any
	InsertSQL() string
	UpdateArgs() []any
	UpdateSQL() string
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
