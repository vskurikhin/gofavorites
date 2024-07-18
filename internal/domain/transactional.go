/*
 * This file was last modified at 2024-07-18 22:20 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * transactional.go
 * $Id$
 */

package domain

import "context"

type TxArgs struct {
	Args [][]any
	SQLs []string
}

type Suite interface {
	Entity
	DeleteTxArgs() TxArgs
	UpsertTxArgs() TxArgs
}

type Dft[S Suite] interface {
	DoDelete(ctx context.Context, entity S, scan func(Scanner)) error
	DoUpsert(ctx context.Context, entity S, scan func(Scanner)) error
}
