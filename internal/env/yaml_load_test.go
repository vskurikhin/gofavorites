/*
 * Copyright text:
 * This file was last modified at 2024-07-10 20:32 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * yaml_load_test.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	assert.NotNil(t, t)
	type want struct {
		yamlConfig Config
		err        error
	}
	var tests = []struct {
		name  string
		input string
		fRun  func(string) (Config, error)
		want  want
	}{
		{
			name:  `positive test #0 LoadConfig(".")`,
			input: ".",
			fRun:  LoadConfig,
			want: want{
				yamlConfig: &config{Favorites: struct {
					Cache struct {
						Enabled     bool
						cacheConfig `mapstructure:",squash"`
					}
					DB struct {
						Enabled  bool
						dbConfig `mapstructure:",squash"`
					}
					Enabled  bool
					External struct {
						externalConfig `mapstructure:",squash"`
					}
					GRPC struct {
						Enabled    bool
						grpcConfig `mapstructure:",squash"`
						TLS        struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}
					}
					HTTP struct {
						Enabled    bool
						httpConfig `mapstructure:",squash"`
						TLS        struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}
					}
					JWT struct {
						jwtConfig `mapstructure:",squash"`
					}
					goFavoritesConfig `mapstructure:",squash"`
				}{
					Cache: struct {
						Enabled     bool
						cacheConfig `mapstructure:",squash"`
					}{
						Enabled: true,
						cacheConfig: cacheConfig{
							ExpireMs:      1000,
							GCIntervalSec: 10,
						},
					},
					DB: struct {
						Enabled  bool
						dbConfig `mapstructure:",squash"`
					}{
						Enabled: false,
						dbConfig: dbConfig{
							Name:         "db",
							Host:         "localhost",
							Port:         5432,
							UserName:     "dbuser",
							UserPassword: "password",
						},
					},
					Enabled: true,
					External: struct {
						externalConfig `mapstructure:",squash"`
					}{
						externalConfig: externalConfig{
							AssetGRPCAddress:       "localhost:8444",
							AuthGRPCAddress:        "localhost:8444",
							RequestTimeoutInterval: 3333,
						},
					},
					GRPC: struct {
						Enabled    bool
						grpcConfig `mapstructure:",squash"`
						TLS        struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}
					}{
						Enabled: true,
						grpcConfig: grpcConfig{
							Address: "localhost",
							Port:    8443,
							Proto:   "tcp",
						},
						TLS: struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}{
							Enabled: true,
							tlsConfig: tlsConfig{
								CAFile:   "cert/grpc-ca-cert.pem",
								CertFile: "cert/grpc-server-cert.pem",
								KeyFile:  "cert/grpc-server-key.pem",
							},
						},
					},
					HTTP: struct {
						Enabled    bool
						httpConfig `mapstructure:",squash"`
						TLS        struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}
					}{
						Enabled: true,
						httpConfig: httpConfig{
							Address: "localhost",
							Port:    443,
						},
						TLS: struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}{
							Enabled: true,
							tlsConfig: tlsConfig{
								CAFile:   "cert/http-ca-cert.pem",
								CertFile: "cert/http-server-cert.pem",
								KeyFile:  "cert/http-server-key.pem",
							},
						},
					},
					JWT: struct {
						jwtConfig `mapstructure:",squash"`
					}{
						jwtConfig: jwtConfig{
							JwtSecret:    "HyZPFEaRf5he4zezLWy5QdSvAdOBWoAgJq5wTvoUH06TYVucOnSGPhSRPp7mkFF",
							JwtExpiresIn: time.Duration(60) * time.Minute,
							JwtMaxAgeSec: 60,
						},
					},
					goFavoritesConfig: goFavoritesConfig{
						Token: "$2a$11$ZTzzVGdLUJGcYKJws9UoUug3Q3kCMELVziajBSJPY3k0pNu2XWHBy",
					},
				}},
				err: nil,
			},
		},
	}
	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.fRun(test.input)
			assert.Equal(t, test.want.yamlConfig, got)
			assert.Equal(t, test.want.err, err)
		})
	}
}
