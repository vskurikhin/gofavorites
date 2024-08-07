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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
					MONGO struct {
						Enabled  bool
						dbConfig `mapstructure:",squash"`
					}
					goFavoritesConfig `mapstructure:",squash"`
					UPK               struct {
						upkConfig `mapstructure:",squash"`
					}
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
							Port:    8442,
							Proto:   "tcp",
						},
						TLS: struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}{
							Enabled: true,
							tlsConfig: tlsConfig{
								CAFile:   "cert/grpc-test_ca-cert.pem",
								CertFile: "cert/grpc-test_server-cert.pem",
								KeyFile:  "cert/grpc-test_server-key.pem",
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
							Port:    8443,
						},
						TLS: struct {
							Enabled   bool
							tlsConfig `mapstructure:",squash"`
						}{
							Enabled: true,
							tlsConfig: tlsConfig{
								CAFile:   "cert/http-test_ca-cert.pem",
								CertFile: "cert/http-test_server-cert.pem",
								KeyFile:  "cert/http-test_server-key.pem",
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
					MONGO: struct {
						Enabled  bool
						dbConfig `mapstructure:",squash"`
					}{
						Enabled: false,
						dbConfig: dbConfig{
							Name:         "db",
							Host:         "localhost",
							Port:         27017,
							UserName:     "mongouser",
							UserPassword: "password",
						},
					},
					goFavoritesConfig: goFavoritesConfig{
						Token: "$2a$11$ZTzzVGdLUJGcYKJws9UoUug3Q3kCMELVziajBSJPY3k0pNu2XWHBy",
					},
					UPK: struct {
						upkConfig `mapstructure:",squash"`
					}{
						upkConfig: upkConfig{
							RSAPrivateKeyFile: "cert/upk-private-key.pem",
							RSAPublicKeyFile:  "cert/upk-public-key.pem",
							Secret:            "qYhaPtg+PIQtBhAU5fHCeQw7XIF3WLKoLPZnJgq1H//DDOB8o2qrP9goVCUZldOdwqLAHxWOGHuvXcwaIFRrD8I3Hz5tRCgCeI+cEZD9h4c4h6ADSjkcrPXg5eRwnANasBkKKZQz8noYwvt9Z9p7HdOtrBmQOi7OVjTfY0T2SnI=",
						},
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
