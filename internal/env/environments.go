/*
 * This file was last modified at 2024-07-21 11:24 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * environments.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	c0env "github.com/caarlos0/env"
)

type environments struct {
	CacheExpire                    int      `env:"CACHE_EXPIRE_MS"`
	CacheGCInterval                int      `env:"CACHE_GC_INTERVAL_SEC"`
	DataBaseDSN                    string   `env:"DATABASE_DSN"`
	ExternalAssetGRPCAddress       []string `env:"ASSET_GRPC_ADDRESS" envSeparator:":"`
	ExternalAuthGRPCAddress        []string `env:"AUTH_GRPC_ADDRESS" envSeparator:":"`
	ExternalRequestTimeoutInterval int      `env:"REQUEST_TIMEOUT_INTERVAL_MS"`
	GRPCAddress                    []string `env:"GRPC_ADDRESS" envSeparator:":"`
	GRPCCAFile                     string   `env:"GRPC_CA_FILE"`
	GRPCCertFile                   string   `env:"GRPC_CERT_FILE"`
	GRPCKeyFile                    string   `env:"GRPC_KEY_FILE"`
	HTTPAddress                    []string `env:"HTTP_ADDRESS" envSeparator:":"`
	HTTPCAFile                     string   `env:"HTTP_CA_FILE"`
	HTTPCertFile                   string   `env:"HTTP_CERT_FILE"`
	HTTPKeyFile                    string   `env:"HTTP_KEY_FILE"`
}

func getEnvironments() (env *environments, err error) {

	env = new(environments)
	err = c0env.Parse(env)
	return
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
