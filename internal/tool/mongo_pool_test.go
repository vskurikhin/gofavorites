/*
 * This file was last modified at 2024-07-31 17:21 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * mongo_pool_test.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMongodb(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "test #1 positive ",
			fRun: testMongodbConnect,
		},
		{
			name: "test #2 negative ",
			fRun: testMongodbGetConnectionNegative,
		},
		{
			name: "test #3 negative ",
			fRun: testMongodbCloseConnectionNegative,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testMongodbConnect(t *testing.T) {
	got := MongodbConnect("")
	assert.NotNil(t, got)
}

func testMongodbGetConnectionNegative(t *testing.T) {
	defer func() { _ = recover() }() // TODO
	pool := MongodbConnect("")
	assert.NotNil(t, pool)
	got, err := pool.GetConnection()
	assert.Nil(t, got)
	assert.NotNil(t, err)
}

func testMongodbCloseConnectionNegative(t *testing.T) {
	defer func() { _ = recover() }() // TODO
	pool := MongodbConnect("")
	assert.NotNil(t, pool)
	err := pool.CloseConnection(nil)
	assert.Nil(t, err)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
