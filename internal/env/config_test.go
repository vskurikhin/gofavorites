/*
 * Copyright text:
 * This file was last modified at 2024-07-10 20:02 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * config_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	var tests = []struct {
		name  string
		fRun  func() Config
		isNil bool
		want  string
	}{
		{
			name:  `positive test #0 nil config`,
			fRun:  nilConfig,
			isNil: true,
			want: `CacheEnabled: false
CacheExpire: 0
CacheGCInterval: 0
DBHost: 
DBName: 
DBEnabled: false
DBPort: 0
DBUserName: 
DBUserPassword: 
Enabled: false
ExternalAssetGRPCAddress: 
ExternalAuthGRPCAddress: 
ExternalRequestTimeoutInterval: 0
GRPCAddress: 
GRPCEnabled: false
GRPCPort: 0
GRPCProto: 
GRPCTLSCAFile: 
GRPCTLSCertFile: 
GRPCTLSKeyFile: 
GRPCTLSEnabled: false
HTTPAddress: 
HTTPEnabled: false
HTTPPort: 0
HTTPTLSCAFile: 
HTTPTLSCertFile: 
HTTPTLSEnabled: false
HTTPTLSKeyFile: 
MongoHost: 
MongoName: 
MongoEnabled: false
MongoPort: 0
MongoUserName: 
MongoUserPassword: 
Token: 
UpkRSAPrivateKeyFile: 
UpkRSAPublicKeyFile: 
UpkSecretKey: `,
		},
		{
			name: `positive test #1 zero config`,
			fRun: zeroConfig,
			want: `CacheEnabled: false
CacheExpire: 0
CacheGCInterval: 0
DBHost: 
DBName: 
DBEnabled: false
DBPort: 0
DBUserName: 
DBUserPassword: 
Enabled: false
ExternalAssetGRPCAddress: 
ExternalAuthGRPCAddress: 
ExternalRequestTimeoutInterval: 0
GRPCAddress: 
GRPCEnabled: false
GRPCPort: 0
GRPCProto: 
GRPCTLSCAFile: 
GRPCTLSCertFile: 
GRPCTLSKeyFile: 
GRPCTLSEnabled: false
HTTPAddress: 
HTTPEnabled: false
HTTPPort: 0
HTTPTLSCAFile: 
HTTPTLSCertFile: 
HTTPTLSEnabled: false
HTTPTLSKeyFile: 
MongoHost: 
MongoName: 
MongoEnabled: false
MongoPort: 0
MongoUserName: 
MongoUserPassword: 
Token: 
UpkRSAPrivateKeyFile: 
UpkRSAPublicKeyFile: 
UpkSecretKey: `,
		},
	}
	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.fRun()
			if !test.isNil {
				assert.Equal(t, test.want, got.String())
			} else {
				assert.Equal(t, test.want, (*config)(nil).String())
			}
		})
	}
}

func nilConfig() Config {
	return nil
}

func zeroConfig() Config {
	return &config{}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
