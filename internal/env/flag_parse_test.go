/*
 * Copyright text:
 * This file was last modified at 2024-07-10 21:11 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * flag_parse_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMakeFlagsParse(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #0 makeFlagsParse",
			fRun: tryMakeFlagsParse,
		},
	}
	oldCommandLine := pflag.CommandLine
	defer func() { pflag.CommandLine = oldCommandLine }()
	ResetForTesting(func() {})

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func tryMakeFlagsParse(t *testing.T) {
	tests := []string{
		flagDatabaseDSN,
		flagGRPCAddress,
		flagGRPCCAFile,
		flagGRPCCertFile,
		flagGRPCKeyFile,
		flagHTTPAddress,
		flagHTTPCAFile,
		flagHTTPCertFile,
		flagHTTPKeyFile,
		flagJwtSecret,
		flagUpkPrivateKeyFile,
		flagUpkPublicKeyFile,
		flagUpkSecret,
	}
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd"}
	for _, in := range tests {
		os.Args = append(os.Args, "--"+in)
		os.Args = append(os.Args, "test")
	}
	m := makeFlagsParse()
	for _, in := range tests {
		ps, ok := m[in].(*string)
		assert.True(t, ok)
		assert.Equal(t, "test", *ps)
	}
}

func ResetForTesting(usage func()) {
	pflag.CommandLine = &pflag.FlagSet{}
	pflag.Usage = usage
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
