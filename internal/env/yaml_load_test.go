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
					Enabled bool
					GRPC    struct {
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
					goFavoritesConfig `mapstructure:",squash"`
				}{

					Cache: struct {
						Enabled     bool
						cacheConfig `mapstructure:",squash"`
					}{
						Enabled: true,
						cacheConfig: cacheConfig{
							Expire:     1000,
							GCInterval: 10,
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
					goFavoritesConfig: goFavoritesConfig{
						Token: "89h3f98hbwf987h3f98wenf89ehf",
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
