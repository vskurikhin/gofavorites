/*
 * Copyright text:
 * This file was last modified at 2024-07-10 20:19 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * environments_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEnvironments(t *testing.T) {
	type want struct {
		env *environments
		err error
	}
	var tests = []struct {
		name string
		fRun func() (env *environments, err error)
		want want
	}{
		{
			name: `positive test #0 nil environments`,
			fRun: getEnvironments,
			want: want{&environments{}, nil},
		},
	}
	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun()
			assert.Equal(t, test.want.env, got)
			assert.Equal(t, test.want.err, err)
		})
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
