/*
 * This file was last modified at 2024-07-16 17:09 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * repo.go
 * $Id$
 */
//!+

// Package domain TODO.
package domain

import "context"

type Repo[E Entity] interface {
	Delete(ctx context.Context, entity E, scan func(Scanner)) (E, error)
	Get(ctx context.Context, entity E, scan func(Scanner)) (E, error)
	Insert(ctx context.Context, entity E, scan func(Scanner)) (E, error)
	Update(ctx context.Context, entity E, scan func(Scanner)) (E, error)
}

type Scanner interface {
	Scan(dest ...any) error
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
