/*
 * This file was last modified at 2024-08-03 12:36 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * properties_tool_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPropertiesToolNegative(t *testing.T) {
	for _, test := range []struct {
		name string
		fRun func(*testing.T) (interface{}, error)
	}{
		{
			"test #1 negative for getUpkSecretKey",
			testGetUpkSecretKey,
		},
		{
			"test #2 negative for intPrepareProperty",
			testIntPrepareProperty,
		},
		{
			"test #3 negative for getFileName",
			testGetFileName,
		},
		{
			"test #4 negative #1 for serverAddressPrepareProperty",
			testServerAddressPrepareProperty1,
		},
		{
			"test #5 negative #2 for serverAddressPrepareProperty",
			testServerAddressPrepareProperty2,
		},
		{
			"test #6 negative for stringsAddressPrepareProperty",
			testStringsAddressPrepareProperty,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(t)
			assert.NotNil(t, err)
			assert.Nil(t, got)
		})
	}
}

func testGetUpkSecretKey(_ *testing.T) (interface{}, error) {
	return getUpkSecretKey(make(map[string]interface{}), &environments{}, &config{}, nil)
}

func testIntPrepareProperty(t *testing.T) (interface{}, error) {
	f := false
	got, err := intPrepareProperty("", &f, 0, 0)
	assert.Equal(t, 0, got)
	return nil, err
}

func testGetFileName(t *testing.T) (interface{}, error) {
	f := false
	got, err := getFileName("", &f, "", "")
	assert.Equal(t, "", got)
	return nil, err
}

func testServerAddressPrepareProperty1(t *testing.T) (interface{}, error) {
	s := ""
	m := map[string]interface{}{"": &s}
	got, err := serverAddressPrepareProperty("", m, []string{}, "", 0)
	assert.Equal(t, "", got)
	return nil, err
}

func testServerAddressPrepareProperty2(t *testing.T) (interface{}, error) {
	m := map[string]interface{}{"": false}
	got, err := serverAddressPrepareProperty("", m, []string{}, "", 0)
	assert.Equal(t, ":0", got)
	return nil, err
}

func testStringsAddressPrepareProperty(t *testing.T) (interface{}, error) {
	f := false
	got, err := stringsAddressPrepareProperty("", &f, []string{}, "")
	assert.Equal(t, "", got)
	return nil, err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
