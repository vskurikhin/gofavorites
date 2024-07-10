/*
 * Copyright text:
 * This file was last modified at 2024-07-10 21:51 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * properties_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

func TestGetProperties(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #0 GetProperties",
			fRun: tryDefaultGetProperties,
		},
		{
			name: "positive test #1 GetProperties",
			fRun: tryEnvGRPCAddressGetProperties,
		},
		{
			name: "positive test #2 GetProperties",
			fRun: tryFlagGRPCAddressGetProperties,
		},
		{
			name: "positive test #3 GetProperties",
			fRun: tryEnvDatabaseDSNGetProperties,
		},
		{
			name: "positive test #4 GetProperties",
			fRun: tryFlagDatabaseDSNGetProperties,
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

func tryDefaultGetProperties(t *testing.T) {
	defer func() { _ = recover() }()
	once = new(sync.Once)
	got := GetProperties()
	expected := fmt.Sprintf("%p", got)
	assert.NotNil(t, got)
	assert.NotNil(t, got.String())
	assert.NotEqual(t, "", got.String())
	assert.Equal(t, expected, fmt.Sprintf("%p", GetProperties()))
	assert.Equal(t, "db", got.DBPool().Config().ConnConfig.Database)
}

func tryEnvGRPCAddressGetProperties(t *testing.T) {
	once = new(sync.Once)
	t.Setenv("GRPC_ADDRESS", "udp:127.0.0.1:0")
	got := GetProperties()
	expected := fmt.Sprintf("%p", got)
	assert.NotNil(t, got)
	assert.NotNil(t, got.String())
	assert.NotEqual(t, "", got.String())
	assert.Equal(t, expected, fmt.Sprintf("%p", GetProperties()))
	assert.Equal(t, "udp:127.0.0.1:0", got.GRPCAddress())
}

func tryFlagGRPCAddressGetProperties(t *testing.T) {
	once = new(sync.Once)
	oldCommandLine := pflag.CommandLine
	defer func() { pflag.CommandLine = oldCommandLine }()
	ResetForTesting(func() {})
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	t.Setenv("GRPC_ADDRESS", "udp:127.0.0.1:1")
	os.Args = []string{"cmd"}
	os.Args = append(os.Args, "--"+flagGRPCAddress)
	os.Args = append(os.Args, "udp:127.0.0.1:0")
	got := GetProperties()
	expected := fmt.Sprintf("%p", got)
	assert.NotNil(t, got)
	assert.NotNil(t, got.String())
	assert.NotEqual(t, "", got.String())
	assert.Equal(t, expected, fmt.Sprintf("%p", GetProperties()))
	assert.Equal(t, "udp:127.0.0.1:0", got.GRPCAddress())
}

func tryEnvDatabaseDSNGetProperties(t *testing.T) {
	defer func() { _ = recover() }()
	once = new(sync.Once)
	t.Setenv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable")
	got := GetProperties()
	expected := fmt.Sprintf("%p", got)
	assert.NotNil(t, got)
	assert.NotNil(t, got.String())
	assert.NotEqual(t, "", got.String())
	assert.Equal(t, expected, fmt.Sprintf("%p", GetProperties()))
	assert.Equal(t, "praktikum", got.DBPool().Config().ConnConfig.Database)
}

func tryFlagDatabaseDSNGetProperties(t *testing.T) {
	defer func() { _ = recover() }()
	once = new(sync.Once)
	oldCommandLine := pflag.CommandLine
	defer func() { pflag.CommandLine = oldCommandLine }()
	ResetForTesting(func() {})
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd"}
	os.Args = append(os.Args, "--"+flagDatabaseDSN)
	os.Args = append(os.Args, "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable")
	got := GetProperties()
	expected := fmt.Sprintf("%p", got)
	assert.NotNil(t, got)
	assert.NotNil(t, got.String())
	assert.NotEqual(t, "", got.String())
	assert.Equal(t, expected, fmt.Sprintf("%p", GetProperties()))
	assert.Equal(t, "praktikum", got.DBPool().Config().ConnConfig.Database)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
