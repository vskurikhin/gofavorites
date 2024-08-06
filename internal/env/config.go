/*
 * This file was last modified at 2024-08-06 18:20 by Victor N. Skurikhin.
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

// Config статичная конфигурация собранная из Yaml-файла.
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
	MongoEnabled() bool
	MongoHost() string
	MongoName() string
	MongoPort() int
	MongoUserName() string
	MongoUserPassword() string
	Token() string
	UpkRSAPrivateKeyFile() string
	UpkRSAPublicKeyFile() string
	UpkSecret() string
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
		MONGO struct {
			Enabled  bool
			dbConfig `mapstructure:",squash"`
		}
		goFavoritesConfig `mapstructure:",squash"`
		UPK               struct {
			upkConfig `mapstructure:",squash"`
		}
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

type goFavoritesConfig struct {
	Token string `mapstructure:"token"`
}

type upkConfig struct {
	RSAPrivateKeyFile string `mapstructure:"rsa_private_key_file"`
	RSAPublicKeyFile  string `mapstructure:"rsa_public_key_file"`
	Secret            string `mapstructure:"secret"`
}

// CacheEnabled тумблер включения локального кэша.
func (y *config) CacheEnabled() bool {

	if y != nil {
		return y.Favorites.Cache.Enabled
	}
	return false
}

// CacheExpireMs срок действия записи в кэше истекает в миллисекундах.
func (y *config) CacheExpireMs() int {

	if y != nil {
		return y.Favorites.Cache.ExpireMs
	}
	return 0
}

// CacheGCIntervalSec интервал очистки кэша в секундах.
func (y *config) CacheGCIntervalSec() int {

	if y != nil {
		return y.Favorites.Cache.GCIntervalSec
	}
	return 0
}

// DBEnabled тумблер подключения к базе данных PostgreSQL.
func (y *config) DBEnabled() bool {

	if y != nil {
		return y.Favorites.DB.Enabled
	}
	return false
}

// DBHost хост базы данных PostgreSQL.
func (y *config) DBHost() string {

	if y != nil {
		return y.Favorites.DB.Host
	}
	return ""
}

// DBName имя базы данных PostgreSQL.
func (y *config) DBName() string {

	if y != nil {
		return y.Favorites.DB.Name
	}
	return ""
}

// DBPort порт базы данных PostgreSQL.
func (y *config) DBPort() int {

	if y != nil {
		return int(y.Favorites.DB.Port)
	}
	return 0
}

// DBUserName имя пользователя базы данных PostgreSQL.
func (y *config) DBUserName() string {

	if y != nil {
		return y.Favorites.DB.UserName
	}
	return ""
}

// DBUserPassword пароль пользователя базы данных PostgreSQL.
func (y *config) DBUserPassword() string {

	if y != nil {
		return y.Favorites.DB.UserPassword
	}
	return ""
}

// Enabled тумблер старта приложения.
func (y *config) Enabled() bool {

	if y != nil {
		return y.Favorites.Enabled
	}
	return false
}

// ExternalAssetGRPCAddress внешний адрес gRPC-сервиса по биржевым инструментам.
func (y *config) ExternalAssetGRPCAddress() string {
	if y != nil {
		return y.Favorites.External.AssetGRPCAddress
	}
	return ""
}

// ExternalAuthGRPCAddress внешний адрес gRPC-сервиса аутентификации пользователей.
func (y *config) ExternalAuthGRPCAddress() string {
	if y != nil {
		return y.Favorites.External.AuthGRPCAddress
	}
	return ""
}

// ExternalRequestTimeoutInterval интервал ожидания ответа от внешних gRPC-сервисов.
func (y *config) ExternalRequestTimeoutInterval() int {

	if y != nil {
		return y.Favorites.External.RequestTimeoutInterval
	}
	return 0
}

// GRPCAddress адрес для выставления конечных точек gRPC-сервера.
func (y *config) GRPCAddress() string {

	if y != nil {
		return y.Favorites.GRPC.Address
	}
	return ""
}

// GRPCEnabled тумблер включения gRPC-сервера.
func (y *config) GRPCEnabled() bool {

	if y != nil {
		return y.Favorites.GRPC.Enabled
	}
	return false
}

// GRPCPort порт для gRPC-сервера.
func (y *config) GRPCPort() int {

	if y != nil {
		return int(y.Favorites.GRPC.Port)
	}
	return 0
}

// GRPCProto протокол для gRPC-сервера.
func (y *config) GRPCProto() string {

	if y != nil {
		return y.Favorites.GRPC.Proto
	}
	return ""
}

// GRPCTLSCAFile корневой сертификат центра сертификации который выдал TLS сертификат для gRPC-сервера.
func (y *config) GRPCTLSCAFile() string {

	if y != nil {
		return y.Favorites.GRPC.TLS.CAFile
	}
	return ""
}

// GRPCTLSCertFile TLS сертификат для gRPC-сервера.
func (y *config) GRPCTLSCertFile() string {

	if y != nil {
		return y.Favorites.GRPC.TLS.CertFile
	}
	return ""
}

// GRPCTLSKeyFile TLS ключ для gRPC-сервера.
func (y *config) GRPCTLSKeyFile() string {

	if y != nil {
		return y.Favorites.GRPC.TLS.KeyFile
	}
	return ""
}

// GRPCTLSEnabled тумблер включения на gRPC-сервере TLS шифрования.
func (y *config) GRPCTLSEnabled() bool {

	if y != nil {
		return y.Favorites.GRPC.TLS.Enabled
	}
	return false
}

// HTTPAddress адрес для выставления конечных точек HTTP-сервера.
func (y *config) HTTPAddress() string {

	if y != nil {
		return y.Favorites.HTTP.Address
	}
	return ""
}

// HTTPEnabled тумблер включения HTTP-сервера.
func (y *config) HTTPEnabled() bool {

	if y != nil {
		return y.Favorites.HTTP.Enabled
	}
	return false
}

// HTTPPort порт для HTTP-сервера.
func (y *config) HTTPPort() int {

	if y != nil {
		return int(y.Favorites.HTTP.Port)
	}
	return 0
}

// HTTPTLSCAFile корневой сертификат центра сертификации который выдал TLS сертификат для HTTP-сервера.
func (y *config) HTTPTLSCAFile() string {

	if y != nil {
		return y.Favorites.HTTP.TLS.CAFile
	}
	return ""
}

// HTTPTLSCertFile TLS сертификат для HTTP-сервера.
func (y *config) HTTPTLSCertFile() string {

	if y != nil {
		return y.Favorites.HTTP.TLS.CertFile
	}
	return ""
}

// HTTPTLSKeyFile TLS ключ для HTTP-сервера.
func (y *config) HTTPTLSKeyFile() string {

	if y != nil {
		return y.Favorites.HTTP.TLS.KeyFile
	}
	return ""
}

// HTTPTLSEnabled тумблер включения на HTTP-сервере TLS шифрования.
func (y *config) HTTPTLSEnabled() bool {

	if y != nil {
		return y.Favorites.HTTP.TLS.Enabled
	}
	return false
}

// JwtExpiresIn Утверждение «exp» (время истечения срока действия)
// определяет время истечения срока действия или после чего JWT
// НЕ ДОЛЖЕН приниматься в обработку.
func (y *config) JwtExpiresIn() time.Duration {

	if y != nil {
		return y.Favorites.JWT.JwtExpiresIn
	}
	return 0
}

// JwtMaxAgeSec определяет время жизни куки в секундах.
func (y *config) JwtMaxAgeSec() int {

	if y != nil {
		return y.Favorites.JWT.JwtMaxAgeSec
	}
	return 0
}

// JwtSecret секрет для подписи JWТокена.
func (y *config) JwtSecret() string {

	if y != nil {
		return y.Favorites.JWT.JwtSecret
	}
	return ""
}

// MongoEnabled тумблер подключения к MongoDB.
func (y *config) MongoEnabled() bool {

	if y != nil {
		return y.Favorites.MONGO.Enabled
	}
	return false
}

// MongoHost хост MongoDB.
func (y *config) MongoHost() string {

	if y != nil {
		return y.Favorites.MONGO.Host
	}
	return ""
}

// MongoName имя базы данных MongoDB.
func (y *config) MongoName() string {

	if y != nil {
		return y.Favorites.MONGO.Name
	}
	return ""
}

// MongoPort порт базы данных MongoDB.
func (y *config) MongoPort() int {

	if y != nil {
		return int(y.Favorites.MONGO.Port)
	}
	return 0
}

// MongoUserName имя пользователя базы данных MongoDB.
func (y *config) MongoUserName() string {

	if y != nil {
		return y.Favorites.MONGO.UserName
	}
	return ""
}

// MongoUserPassword пароль пользователя базы данных MongoDB.
func (y *config) MongoUserPassword() string {

	if y != nil {
		return y.Favorites.MONGO.UserPassword
	}
	return ""
}

func (y *config) Token() string {

	if y != nil {
		return y.Favorites.Token
	}
	return ""
}

// UpkRSAPrivateKeyFile RSA ключ для дешифрации секрета
// который применяется в симметричном шифровании UPK (User Personal Key).
func (y *config) UpkRSAPrivateKeyFile() string {

	if y != nil {
		return y.Favorites.UPK.RSAPrivateKeyFile
	}
	return ""
}

// UpkRSAPublicKeyFile RSA ключ для шифрования секрета
// который применяется в симметричном шифровании UPK (User Personal Key).
func (y *config) UpkRSAPublicKeyFile() string {

	if y != nil {
		return y.Favorites.UPK.RSAPublicKeyFile
	}
	return ""
}

// UpkSecret секрет который применяется в симметричном шифровании UPK (User Personal Key).
func (y *config) UpkSecret() string {

	if y != nil {
		return y.Favorites.UPK.Secret
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
MongoHost: %s
MongoName: %s
MongoEnabled: %v
MongoPort: %d
MongoUserName: %s
MongoUserPassword: %s
Token: %s
UpkRSAPrivateKeyFile: %s
UpkRSAPublicKeyFile: %s
UpkSecretKey: %s`,
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
		y.MongoHost(),
		y.MongoName(),
		y.MongoEnabled(),
		y.MongoPort(),
		y.MongoUserName(),
		base64.StdEncoding.EncodeToString([]byte(y.MongoUserPassword())),
		y.Token(),
		y.UpkRSAPrivateKeyFile(),
		y.UpkRSAPublicKeyFile(),
		base64.StdEncoding.EncodeToString([]byte(y.UpkSecret())),
	)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
