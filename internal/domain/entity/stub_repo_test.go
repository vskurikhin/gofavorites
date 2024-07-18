/*
 * This file was last modified at 2024-07-16 21:11 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * stub_repo_test.go
 * $Id$
 */
//!+

// Package entity TODO.
package entity

import (
	"context"
	"fmt"
	"github.com/vskurikhin/gofavorites/internal/domain"
)

type stubRepoOk[E domain.Entity] struct {
}

var _ domain.Repo[domain.Entity] = (*stubRepoOk[domain.Entity])(nil)

func (p *stubRepoOk[E]) Delete(_ context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	if entity.DeleteSQL() == "" {
		panic(`entity.DeleteSQL() == ""`)
	}
	if len(entity.DeleteArgs()) < 1 {
		panic(`len(entity.DeleteArgs()) < 1`)
	}
	scan(&stubScannerOk{})
	return entity, nil
}

func (p *stubRepoOk[E]) Get(_ context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	if entity.GetSQL() == "" {
		panic(`entity.GetSQL() == ""`)
	}
	if len(entity.GetArgs()) < 1 {
		panic(`len(entity.GetArgs()) < 1`)
	}
	scan(&stubScannerOk{})
	return entity, nil
}

func (p *stubRepoOk[E]) Insert(_ context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	if entity.InsertSQL() == "" {
		panic(`entity.InsertSQL() == ""`)
	}
	if len(entity.InsertArgs()) < 1 {
		panic(`len(entity.DeleteArgs()) < 1`)
	}
	scan(&stubScannerOk{})
	return entity, nil
}

func (p *stubRepoOk[E]) Update(_ context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	if entity.UpdateSQL() == "" {
		panic(`entity.UpdateSQL() == ""`)
	}
	if len(entity.UpdateArgs()) < 1 {
		panic(`len(entity.UpdateArgs()) < 1`)
	}
	scan(&stubScannerOk{})
	return entity, nil
}

type stubTxRepoOk[S domain.Suite] struct {
}

var _ domain.Dft[domain.Suite] = (*stubTxRepoOk[domain.Suite])(nil)

func (p *stubTxRepoOk[S]) DoDelete(ctx context.Context, entity S, scan func(domain.Scanner)) error {
	if len(entity.DeleteTxArgs().Args) < 1 {
		panic(`len(entity.DeleteTxArgs().Args) < 1`)
	}
	if len(entity.DeleteTxArgs().SQLs) < 1 {
		panic(`len(entity.DeleteTxArgs().SQLs) < 1`)
	}
	scan(&stubScannerOk{})
	return nil
}

func (p *stubTxRepoOk[S]) DoUpsert(ctx context.Context, entity S, scan func(domain.Scanner)) error {
	if len(entity.UpsertTxArgs().Args) < 1 {
		panic(`len(entity.UpsertTxArgs().Args) < 1`)
	}
	if len(entity.UpsertTxArgs().SQLs) < 1 {
		panic(`len(entity.UpsertTxArgs().SQLs) < 1`)
	}
	scan(&stubScannerOk{})
	return nil
}

type stubRepoErr[E domain.Entity] struct {
}

var _ domain.Repo[domain.Entity] = (*stubRepoErr[domain.Entity])(nil)

func (p *stubRepoErr[E]) Delete(_ context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	scan(&stubScannerErr{})
	return entity, nil
}

func (p *stubRepoErr[E]) Get(_ context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	scan(&stubScannerErr{})
	return entity, nil
}

func (p *stubRepoErr[E]) Insert(_ context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	scan(&stubScannerErr{})
	return entity, nil
}

func (p *stubRepoErr[E]) Update(_ context.Context, entity E, scan func(domain.Scanner)) (E, error) {
	scan(&stubScannerErr{})
	return entity, nil
}

type stubTxRepoErr[S domain.Suite] struct {
}

var _ domain.Dft[domain.Suite] = (*stubTxRepoErr[domain.Suite])(nil)

func (p *stubTxRepoErr[S]) DoDelete(_ context.Context, _ S, scan func(domain.Scanner)) error {
	scan(&stubScannerErr{})
	return nil
}

func (p *stubTxRepoErr[S]) DoUpsert(_ context.Context, _ S, scan func(domain.Scanner)) error {
	scan(&stubScannerErr{})
	return nil
}

type stubScannerOk struct {
}

func (s *stubScannerOk) Scan(_ ...any) error {
	return nil
}

type stubScannerErr struct {
}

func (s *stubScannerErr) Scan(_ ...any) error {
	return fmt.Errorf("%s error", "stub")
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
