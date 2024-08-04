/*
 * This file was last modified at 2024-08-03 14:56 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * batch_test.go
 * $Id$
 */

// Package batch TODO.
package batch

import (
	"context"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/env"
)

func TestFavoritesInsertsBatch(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "test #1 negative",
			fRun: testBatchDoNegative,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testBatchDoNegative(t *testing.T) {
	defer func() { _ = recover() }() // TODO
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("DATABASE_DSN", "")

	prop := env.GetProperties()
	batch := GetBatchPostgres(prop)
	err := batch.Do(context.Background(), []entity.Favorites{{}}, "")
	assert.NotNil(t, err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
