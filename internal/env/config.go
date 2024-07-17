/*
 * This file was last modified at 2024-07-17 10:29 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * config.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"encoding/base64"
	"fmt"
)

var _ Config = (*config)(nil)

type Config interface {
	fmt.Stringer
	CacheEnabled() bool
	CacheExpire() int
	CacheGCInterval() int
	DBEnabled() bool
	DBHost() string
	DBName() string
	DBPort() int
	DBUserName() string
	DBUserPassword() string
	Enabled() bool
	GRPCAddress() string
	GRPCEnabled() bool
	GRPCPort() int
	GRPCProto() string
	GRPCTLSCAFile() string
	GRPCTLSCertFile() string
	GRPCTLSEnabled() bool
	GRPCTLSKeyFile() string
	HTTPAddress() string
	HTTPEnabled() bool
	HTTPPort() int
	HTTPTLSCAFile() string
	HTTPTLSCertFile() string
	HTTPTLSEnabled() bool
	HTTPTLSKeyFile() string
	Token() string
}

type config struct {
	Favorites struct {
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
	}
}

type cacheConfig struct {
	Expire     int `mapstructure:"expire_ms"`
	GCInterval int `mapstructure:"gc_interval_sec"`
}

type dbConfig struct {
	Name         string
	Host         string
	Port         int16
	UserName     string
	UserPassword string `mapstructure:"password"`
}

type grpcConfig struct {
	Address string
	Port    int16
	Proto   string
}

type httpConfig struct {
	Address string
	Port    int16
}

type tlsConfig struct {
	CAFile   string `mapstructure:"ca_file"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// moduleConfig could be in a module specific package
type goFavoritesConfig struct {
	Token string
}

func (y *config) CacheEnabled() bool {

	if y != nil {
		return y.Favorites.Cache.Enabled
	}
	return false
}

func (y *config) CacheExpire() int {

	if y != nil {
		return y.Favorites.Cache.Expire
	}
	return 0
}

func (y *config) CacheGCInterval() int {

	if y != nil {
		return y.Favorites.Cache.GCInterval
	}
	return 0
}

func (y *config) DBEnabled() bool {

	if y != nil {
		return y.Favorites.DB.Enabled
	}
	return false
}

func (y *config) DBHost() string {

	if y != nil {
		return y.Favorites.DB.Host
	}
	return ""
}

func (y *config) DBName() string {

	if y != nil {
		return y.Favorites.DB.Name
	}
	return ""
}

func (y *config) DBPort() int {

	if y != nil {
		return int(y.Favorites.DB.Port)
	}
	return 0
}

func (y *config) DBUserName() string {

	if y != nil {
		return y.Favorites.DB.UserName
	}
	return ""
}

func (y *config) DBUserPassword() string {

	if y != nil {
		return y.Favorites.DB.UserPassword
	}
	return ""
}

func (y *config) Enabled() bool {

	if y != nil {
		return y.Favorites.Enabled
	}
	return false
}

func (y *config) GRPCAddress() string {

	if y != nil {
		return y.Favorites.GRPC.Address
	}
	return ""
}

func (y *config) GRPCEnabled() bool {

	if y != nil {
		return y.Favorites.GRPC.Enabled
	}
	return false
}

func (y *config) GRPCPort() int {

	if y != nil {
		return int(y.Favorites.GRPC.Port)
	}
	return 0
}

func (y *config) GRPCProto() string {

	if y != nil {
		return y.Favorites.GRPC.Proto
	}
	return ""
}

func (y *config) GRPCTLSCAFile() string {

	if y != nil {
		return y.Favorites.GRPC.TLS.CAFile
	}
	return ""
}

func (y *config) GRPCTLSCertFile() string {

	if y != nil {
		return y.Favorites.GRPC.TLS.CertFile
	}
	return ""
}

func (y *config) GRPCTLSKeyFile() string {

	if y != nil {
		return y.Favorites.GRPC.TLS.KeyFile
	}
	return ""
}

func (y *config) GRPCTLSEnabled() bool {

	if y != nil {
		return y.Favorites.GRPC.TLS.Enabled
	}
	return false
}

func (y *config) HTTPAddress() string {

	if y != nil {
		return y.Favorites.HTTP.Address
	}
	return ""
}

func (y *config) HTTPEnabled() bool {

	if y != nil {
		return y.Favorites.HTTP.Enabled
	}
	return false
}

func (y *config) HTTPPort() int {

	if y != nil {
		return int(y.Favorites.HTTP.Port)
	}
	return 0
}

func (y *config) HTTPTLSCAFile() string {

	if y != nil {
		return y.Favorites.HTTP.TLS.CAFile
	}
	return ""
}

func (y *config) HTTPTLSCertFile() string {

	if y != nil {
		return y.Favorites.HTTP.TLS.CertFile
	}
	return ""
}

func (y *config) HTTPTLSKeyFile() string {

	if y != nil {
		return y.Favorites.HTTP.TLS.KeyFile
	}
	return ""
}

func (y *config) HTTPTLSEnabled() bool {

	if y != nil {
		return y.Favorites.HTTP.TLS.Enabled
	}
	return false
}

func (y *config) Token() string {

	if y != nil {
		return y.Favorites.Token
	}
	return ""
}

func (y *config) String() string {
	return fmt.Sprintf(
		`CacheEnabled: %v
CacheExpire: %d
CacheGCInterval: %d
DBHost: %s
DBName: %s
DBEnabled: %v
DBPort: %d
DBUserName: %s
DBUserPassword: %s
Enabled: %v
GRPCAddress: %s
GRPCEnabled: %v
GRPCPort: %d
GRPCProto: %s
GRPCTLSCAFile: %s
GRPCTLSCertFile: %s
GRPCTLSKeyFile: %s
GRPCTLSEnabled: %v
HTTPAddress: %s
HTTPEnabled: %v
HTTPPort: %d
HTTPTLSCAFile: %s
HTTPTLSCertFile: %s
HTTPTLSEnabled: %v
HTTPTLSKeyFile: %s
Token: %s`,
		y.CacheEnabled(),
		y.CacheExpire(),
		y.CacheGCInterval(),
		y.DBHost(),
		y.DBName(),
		y.DBEnabled(),
		y.DBPort(),
		y.DBUserName(),
		base64.StdEncoding.EncodeToString([]byte(y.DBUserPassword())),
		y.Enabled(),
		y.GRPCAddress(),
		y.GRPCEnabled(),
		y.GRPCPort(),
		y.GRPCProto(),
		y.GRPCTLSCAFile(),
		y.GRPCTLSCertFile(),
		y.GRPCTLSKeyFile(),
		y.GRPCTLSEnabled(),
		y.HTTPAddress(),
		y.HTTPEnabled(),
		y.HTTPPort(),
		y.HTTPTLSCAFile(),
		y.HTTPTLSCertFile(),
		y.HTTPTLSEnabled(),
		y.HTTPTLSKeyFile(),
		y.Token(),
	)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
