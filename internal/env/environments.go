/*
 * This file was last modified at 2024-07-29 16:42 by Victor N. Skurikhin.
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
	"time"
)

type environments struct {
	CacheExpireMs                  int           `env:"CACHE_EXPIRE_MS"`
	CacheGCIntervalSec             int           `env:"CACHE_GC_INTERVAL_SEC"`
	DataBaseDSN                    string        `env:"DATABASE_DSN"`
	ExternalAssetGRPCAddress       []string      `env:"ASSET_GRPC_ADDRESS" envSeparator:":"`
	ExternalAuthGRPCAddress        []string      `env:"AUTH_GRPC_ADDRESS" envSeparator:":"`
	ExternalRequestTimeoutInterval int           `env:"REQUEST_TIMEOUT_INTERVAL_MS"`
	GRPCAddress                    []string      `env:"GRPC_ADDRESS" envSeparator:":"`
	GRPCCAFile                     string        `env:"GRPC_CA_FILE"`
	GRPCCertFile                   string        `env:"GRPC_CERT_FILE"`
	GRPCKeyFile                    string        `env:"GRPC_KEY_FILE"`
	HTTPAddress                    []string      `env:"HTTP_ADDRESS" envSeparator:":"`
	HTTPCAFile                     string        `env:"HTTP_CA_FILE"`
	HTTPCertFile                   string        `env:"HTTP_CERT_FILE"`
	HTTPKeyFile                    string        `env:"HTTP_KEY_FILE"`
	JwtExpiresIn                   time.Duration `env:"JWT_EXPIRED_IN"`
	JwtMaxAge                      int           `env:"JWT_MAX_AGE_SEC"`
	JwtSecret                      string        `env:"JWT_SECRET"`
	UpkPrivateKeyFile              string        `env:"UPK_PRIVATE_KEY_FILE"`
	UpkPublicKeyFile               string        `env:"UPK_PUBLIC_KEY_FILE"`
	UpkSecret                      string        `env:"UPK_SECRET"`
}

func getEnvironments() (env *environments, err error) {

	env = new(environments)
	err = c0env.Parse(env)
	return
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
