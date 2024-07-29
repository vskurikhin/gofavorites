/*
 * This file was last modified at 2024-07-29 16:43 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * flag_parse.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"github.com/spf13/pflag"
	"time"
)

const (
	flagCacheExpireMs                  = "cache-expire-ms"
	flagCacheGCIntervalSec             = "cache-gc-interval-sec"
	flagDatabaseDSN                    = "database-dsn"
	flagExternalAssetGRPCAddress       = "asset-grpc-address"
	flagExternalAuthGRPCAddress        = "auth-grpc-address"
	flagExternalRequestTimeoutInterval = "request-timeout-interval"
	flagGRPCAddress                    = "grpc-address"
	flagGRPCCAFile                     = "grpc-ca-file"
	flagGRPCCertFile                   = "grpc-cert-file"
	flagGRPCKeyFile                    = "grpc-key-file"
	flagHTTPAddress                    = "http-address"
	flagHTTPCAFile                     = "http-ca-file"
	flagHTTPCertFile                   = "http-cert-file"
	flagHTTPKeyFile                    = "http-key-file"
	flagJwtExpiresIn                   = "jwt-expires-in"
	flagJwtMaxAgeSec                   = "Jwt-max-age-sec"
	flagJwtSecret                      = "jwt-secret"
	flagUpkPrivateKeyFile              = "upk-private-key-file"
	flagUpkPublicKeyFile               = "upk-public-key-file"
	flagUpkSecret                      = "upk-secret"
)

func makeFlagsParse() map[string]interface{} {

	var flagsMap = make(map[string]interface{})

	if !pflag.Parsed() {
		flagsMap[flagCacheExpireMs] = pflag.Int(
			flagCacheExpireMs,
			1000,
			"time to expire key in millisecond",
		)
		flagsMap[flagCacheGCIntervalSec] = pflag.Int(
			flagCacheGCIntervalSec,
			10,
			"time before deleting expired keys in second",
		)
		flagsMap[flagDatabaseDSN] = pflag.StringP(
			flagDatabaseDSN,
			"d",
			"postgres://dbuser:password@localhost:5432/db?sslmode=disable",
			"database DSN",
		)
		flagsMap[flagExternalAssetGRPCAddress] = pflag.String(
			flagExternalAssetGRPCAddress,
			"localhost:8444",
			"asset gRPC server host and port",
		)
		flagsMap[flagExternalAuthGRPCAddress] = pflag.String(
			flagExternalAuthGRPCAddress,
			"localhost:8444",
			"asset gRPC server host and port",
		)
		flagsMap[flagExternalRequestTimeoutInterval] = pflag.IntP(
			flagExternalRequestTimeoutInterval,
			"t",
			5000,
			"timeout in millisecond via gRPC clients for external services",
		)
		flagsMap[flagGRPCAddress] = pflag.StringP(
			flagGRPCAddress,
			"g",
			"localhost:8443",
			"gRPC server host and port",
		)
		flagsMap[flagGRPCCAFile] = pflag.String(
			flagGRPCCAFile,
			"cert/grpc-test_ca-cert.pem",
			"gRPC CA file",
		)
		flagsMap[flagGRPCCertFile] = pflag.String(
			flagGRPCCertFile,
			"cert/grpc-test_server-cert.pem",
			"gRPC server certificate file",
		)
		flagsMap[flagGRPCKeyFile] = pflag.String(
			flagGRPCKeyFile,
			"cert/grpc-test_server-key.pem",
			"gRPC server key file",
		)
		flagsMap[flagHTTPAddress] = pflag.StringP(
			flagHTTPAddress,
			"h",
			"localhost:443",
			"HTTP server host and port",
		)
		flagsMap[flagHTTPCAFile] = pflag.String(
			flagHTTPCAFile,
			"cert/http-test_ca-cert.pem",
			"HTTP CA file",
		)
		flagsMap[flagHTTPCertFile] = pflag.String(
			flagHTTPCertFile,
			"cert/http-test_server-cert.pem",
			"HTTP server certificate file",
		)
		flagsMap[flagHTTPKeyFile] = pflag.String(
			flagHTTPKeyFile,
			"cert/http-test_server-key.pem",
			"HTTP server key file",
		)
		flagsMap[flagJwtExpiresIn] = pflag.Duration(
			flagJwtExpiresIn,
			time.Duration(60)*time.Second,
			"JWT expires in in second",
		)
		flagsMap[flagJwtMaxAgeSec] = pflag.Int(
			flagJwtMaxAgeSec,
			60,
			"JWT max age",
		)
		flagsMap[flagJwtSecret] = pflag.String(
			flagJwtSecret,
			"",
			"JWT secret",
		)
		flagsMap[flagUpkPrivateKeyFile] = pflag.String(
			flagUpkPrivateKeyFile,
			"cert/upk-private-key.pem",
			"Private key file for decrypt UPK secret",
		)
		flagsMap[flagUpkPublicKeyFile] = pflag.String(
			flagUpkPublicKeyFile,
			"cert/upk-public-key.pem",
			"Public key for encrypt UPK secret",
		)
		flagsMap[flagUpkSecret] = pflag.String(
			flagUpkSecret,
			"",
			"UPK secret",
		)
		pflag.Parse()
	}
	return flagsMap
}

func setIfFlagChanged(name string, set func()) {
	pflag.VisitAll(func(f *pflag.Flag) {
		if f.Changed && f.Name == name {
			set()
		}
	})
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
