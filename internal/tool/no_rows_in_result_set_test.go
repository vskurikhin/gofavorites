/*
 * This file was last modified at 2024-07-16 23:24 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * no_rows_in_result_set_test.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoRowsInResultSet(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{name: "positive test #0 NoRowsInResultSet", fRun: testNoRowsInResultSet},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testNoRowsInResultSet(t *testing.T) {
	assert.True(t, NoRowsInResultSet(errors.New("no rows in result set")))
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
