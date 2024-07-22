/*
 * This file was last modified at 2024-07-22 23:58 by Victor N. Skurikhin.
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
	"time"
)

var _ Config = (*config)(nil)

type Config interface {
	fmt.Stringer
	CacheEnabled() bool
	CacheExpireMs() int
	CacheGCIntervalSec() int
	DBEnabled() bool
	DBHost() string
	DBName() string
	DBPort() int
	DBUserName() string
	DBUserPassword() string
	Enabled() bool
	ExternalAssetGRPCAddress() string
	ExternalAuthGRPCAddress() string
	ExternalRequestTimeoutInterval() int
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
	JwtExpiresIn() time.Duration
	JwtMaxAgeSec() int
	JwtSecret() string
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
	}
}

type cacheConfig struct {
	ExpireMs      int `mapstructure:"expire_ms"`
	GCIntervalSec int `mapstructure:"gc_interval_sec"`
}

type dbConfig struct {
	Name         string
	Host         string
	Port         int16
	UserName     string
	UserPassword string `mapstructure:"password"`
}

type externalConfig struct {
	AssetGRPCAddress       string `mapstructure:"asset_grpc_address"`
	AuthGRPCAddress        string `mapstructure:"auth_grpc_address"`
	RequestTimeoutInterval int    `mapstructure:"request_timeout_interval_ms"`
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

type jwtConfig struct {
	JwtSecret    string        `mapstructure:"jwt_secret"`
	JwtExpiresIn time.Duration `mapstructure:"jwt_expired_in"`
	JwtMaxAgeSec int           `mapstructure:"jwt_max_age_sec"`
}

type tlsConfig struct {
	CAFile   string `mapstructure:"ca_file"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// moduleConfig could be in a module specific package
type goFavoritesConfig struct {
	Token string `mapstructure:"token"`
}

func (y *config) CacheEnabled() bool {

	if y != nil {
		return y.Favorites.Cache.Enabled
	}
	return false
}

func (y *config) CacheExpireMs() int {

	if y != nil {
		return y.Favorites.Cache.ExpireMs
	}
	return 0
}

func (y *config) CacheGCIntervalSec() int {

	if y != nil {
		return y.Favorites.Cache.GCIntervalSec
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

func (y *config) ExternalAssetGRPCAddress() string {
	if y != nil {
		return y.Favorites.External.AssetGRPCAddress
	}
	return ""
}

func (y *config) ExternalAuthGRPCAddress() string {
	if y != nil {
		return y.Favorites.External.AuthGRPCAddress
	}
	return ""
}

func (y *config) ExternalRequestTimeoutInterval() int {

	if y != nil {
		return y.Favorites.External.RequestTimeoutInterval
	}
	return 0
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

func (y *config) JwtExpiresIn() time.Duration {

	if y != nil {
		return y.Favorites.JWT.JwtExpiresIn
	}
	return 0
}

func (y *config) JwtMaxAgeSec() int {

	if y != nil {
		return y.Favorites.JWT.JwtMaxAgeSec
	}
	return 0
}

func (y *config) JwtSecret() string {

	if y != nil {
		return y.Favorites.JWT.JwtSecret
	}
	return ""
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
ExternalAssetGRPCAddress: %s
ExternalAuthGRPCAddress: %s
ExternalRequestTimeoutInterval: %d
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
		y.CacheExpireMs(),
		y.CacheGCIntervalSec(),
		y.DBHost(),
		y.DBName(),
		y.DBEnabled(),
		y.DBPort(),
		y.DBUserName(),
		base64.StdEncoding.EncodeToString([]byte(y.DBUserPassword())),
		y.Enabled(),
		y.ExternalAssetGRPCAddress(),
		y.ExternalAuthGRPCAddress(),
		y.ExternalRequestTimeoutInterval(),
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
