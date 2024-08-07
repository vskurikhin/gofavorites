/*
 * This file was last modified at 2024-08-06 20:17 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * properties.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"google.golang.org/grpc/credentials"
)

const (
	propertyCacheExpireMs                  = "cache-expire"
	propertyCacheGCIntervalSec             = "cache-gc-interval"
	propertyConfig                         = "config"
	propertyDBPool                         = "db-pool"
	propertyEnvironments                   = "environments"
	propertyExternalAssetGRPCAddress       = "external-asset-grpc-address"
	propertyExternalAuthGRPCAddress        = "external-auth-grpc-address"
	propertyExternalRequestTimeoutInterval = "external-request-timeout-interval"
	propertyFlags                          = "flags"
	propertyGRPCAddress                    = "grpc-address"
	propertyGRPCTransportCredentials       = "grpc-transport-credentials"
	propertyHTTPAddress                    = "http-address"
	propertyHTTPHTTPTLSConfig              = "http-tls-config"
	propertyJwtExpiresIn                   = "jwt-expires-in"
	propertyJwtMaxAgeSec                   = "jwt-max-age-sec"
	propertyJwtSecret                      = "jwt-secret"
	propertyLogger                         = "logger"
	propertyMongodbPool                    = "mongodb-pool"
	propertyUpkRSAPrivateKey               = "upk-rsa-private-key"
	propertyUpkRSAPublicKey                = "upk-rsa-public-key"
	propertyUpkSecretKey                   = "upk-secret-key"
)

// Properties конфигурация собранная из Yaml-файла, переменных окружения и флагов командной строки.
type Properties interface {
	fmt.Stringer
	CacheExpire() time.Duration
	CacheGCInterval() time.Duration
	Config() Config
	DBPool() *pgxpool.Pool
	Environments() environments
	ExternalAssetGRPCAddress() string
	ExternalAuthGRPCAddress() string
	ExternalRequestTimeoutInterval() time.Duration
	Flags() map[string]interface{}
	GRPCAddress() string
	GRPCTransportCredentials() credentials.TransportCredentials
	HTTPAddress() string
	HTTPTLSConfig() *tls.Config
	JwtExpiresIn() time.Duration
	JwtMaxAgeSec() int
	JwtSecret() string
	Logger() *slog.Logger
	MongodbPool() *tool.MongoPool
	SlogJSON() bool
	OutboundIP() net.IP
	UpkRSAPrivateKey() *rsa.PrivateKey
	UpkRSAPublicKey() *rsa.PublicKey
	UpkSecretKey() []byte
}

type mapProperties struct {
	mp sync.Map
}

var properties Properties = (*mapProperties)(nil)
var once = new(sync.Once)

// GetProperties — свойства преобразованные из конфигурации и окружения.
// потокобезопасное (thread-safe) создание.
func GetProperties() Properties {

	once.Do(func() {
		yml, err := LoadConfig(".")
		tool.IfErrorThenPanic(err)
		env, err := getEnvironments()
		tool.IfErrorThenPanic(err)
		flm := makeFlagsParse()

		cacheExpire, err := getCacheExpire(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "cacheExpire", cacheExpire, "err", err)
		cacheGCInterval, err := getCacheGCInterval(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "cacheGCInterval", cacheGCInterval, "err", err)

		dbPool, err := makeDBPool(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "dbDisable", err)

		grpcAddress, err := getGRPCAddress(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "grpcAddress", grpcAddress, "err", err)
		tgRPCCredentials, err := getGRPCTransportCredentials(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "grpcTransportCredentials", tgRPCCredentials, "err", err)

		httpAddress, err := getHTTPAddress(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "httpAddress", httpAddress, "err", err)
		tHTTPConfig, err := getHTTPTLSConfig(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "tHTTPConfig", tHTTPConfig, "err", err)

		assetGRPCAddress, err := getExternalAssetGRPCAddress(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "assetGRPCAddress", assetGRPCAddress, "err", err)
		authGRPCAddress, err := getExternalAuthGRPCAddress(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "authGRPCAddress", authGRPCAddress, "err", err)
		requestTimeoutInterval, err := getExternalRequestTimeoutInterval(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "requestTimeoutInterval", requestTimeoutInterval, "err", err)

		jwtExpiresIn, err := getJwtExpiresIn(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "jwtExpiresIn", jwtExpiresIn, "err", err)
		jwtMaxAgeSec, err := getJwtMaxAgeSec(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "jwtMaxAgeSec", jwtMaxAgeSec, "err", err)
		jwtSecret, err := getJwtSecret(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "jwtSecret", jwtSecret, "err", err)

		mongodbPool, err := makeMongodbPool(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "mongodbDisable", err)

		upkRSAPrivateKey, err := getRSAPrivateKey(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "upkRSAPrivateKey", upkRSAPrivateKey, "err", err)
		upkRSAPublicKey, err := getRSAPublicKey(flm, env, yml)
		slog.Debug(MSG+"GetProperties", "upkRSAPublicKey", upkRSAPublicKey, "err", err)
		upkSecretKey, err := getUpkSecretKey(flm, env, yml, upkRSAPrivateKey)
		slog.Debug(MSG+"GetProperties", "upkSecretKey", base64.StdEncoding.EncodeToString(upkSecretKey), "err", err)

		properties = getProperties(
			WithCacheExpire(cacheExpire),
			WithCacheGCInterval(cacheGCInterval),
			WithConfig(yml),
			WithEnvironments(*env),
			WithExternalAssetGRPCAddress(assetGRPCAddress),
			WithExternalAuthGRPCAddress(authGRPCAddress),
			WithExternalRequestTimeoutInterval(requestTimeoutInterval),
			WithFlags(flm),
			withDBPool(dbPool),
			WithGRPCAddress(grpcAddress),
			WithGRPCTransportCredentials(tgRPCCredentials),
			WithHTTPAddress(httpAddress),
			WithHTTPTLSConfig(tHTTPConfig),
			WithJwtExpiresIn(jwtExpiresIn),
			WithJwtMaxAgeSec(jwtMaxAgeSec),
			WithJwtSecret(jwtSecret),
			WithLogger(setupLogger(slogJSON(flm))),
			withMongodbPool(mongodbPool),
			WithUpkRSAPrivateKey(upkRSAPrivateKey),
			WithUpkRSAPublicKey(upkRSAPublicKey),
			WithUpkSecretKey(upkSecretKey),
		)
	})

	return properties
}

// WithCacheExpire — срок действия записи в кэше.
func WithCacheExpire(cacheExpire time.Duration) func(*mapProperties) {
	return func(p *mapProperties) {
		if cacheExpire > 0 {
			p.mp.Store(propertyCacheExpireMs, cacheExpire)
		}
	}
}

// CacheExpire геттер срока действия записи в кэше.
func (p *mapProperties) CacheExpire() time.Duration {
	if a, ok := p.mp.Load(propertyCacheExpireMs); ok {
		if cacheExpire, ok := a.(time.Duration); ok {
			return cacheExpire
		}
	}
	return 0
}

// WithCacheGCInterval — интервал очистки кэша.
func WithCacheGCInterval(cacheGCInterval time.Duration) func(*mapProperties) {
	return func(p *mapProperties) {
		if cacheGCInterval > 0 {
			p.mp.Store(propertyCacheGCIntervalSec, cacheGCInterval)
		}
	}
}

// CacheGCInterval геттер интервала очистки кэша.
func (p *mapProperties) CacheGCInterval() time.Duration {
	if a, ok := p.mp.Load(propertyCacheGCIntervalSec); ok {
		if cacheGCInterval, ok := a.(time.Duration); ok {
			return cacheGCInterval
		}
	}
	return 0
}

// WithConfig — Конфигурация.
func WithConfig(config Config) func(*mapProperties) {
	return func(p *mapProperties) {
		if config != nil {
			p.mp.Store(propertyConfig, config)
		}
	}
}

// Config — текущая конфигурация.
func (p *mapProperties) Config() Config {
	if c, ok := p.mp.Load(propertyConfig); ok {
		if cfg, ok := c.(Config); ok {
			return cfg
		}
	}
	return nil
}

// WithEnvironments — Окружение.
func WithEnvironments(env environments) func(*mapProperties) {
	return func(p *mapProperties) {
		p.mp.Store(propertyEnvironments, env)
	}
}

// Environments геттер Окружения.
func (p *mapProperties) Environments() environments {
	if f, ok := p.mp.Load(propertyEnvironments); ok {
		if env, ok := f.(environments); ok {
			return env
		}
	}
	return environments{}
}

// WithExternalAssetGRPCAddress — внешний адрес gRPC-сервиса по биржевым инструментам.
func WithExternalAssetGRPCAddress(address string) func(*mapProperties) {
	return func(p *mapProperties) {
		p.mp.Store(propertyExternalAssetGRPCAddress, address)
	}
}

// ExternalAssetGRPCAddress геттер внешнего адреса gRPC-сервиса по биржевым инструментам.
func (p *mapProperties) ExternalAssetGRPCAddress() string {
	if a, ok := p.mp.Load(propertyExternalAssetGRPCAddress); ok {
		if address, ok := a.(string); ok {
			return address
		}
	}
	return ""
}

// WithExternalAuthGRPCAddress — внешний адрес gRPC-сервиса аутентификации пользователей.
func WithExternalAuthGRPCAddress(address string) func(*mapProperties) {
	return func(p *mapProperties) {
		p.mp.Store(propertyExternalAuthGRPCAddress, address)
	}
}

// ExternalAuthGRPCAddress геттер внешнего адреса gRPC-сервиса аутентификации пользователей.
func (p *mapProperties) ExternalAuthGRPCAddress() string {
	if a, ok := p.mp.Load(propertyExternalAuthGRPCAddress); ok {
		if address, ok := a.(string); ok {
			return address
		}
	}
	return ""
}

// WithExternalRequestTimeoutInterval — интервал ожидания ответа от внешних gRPC-сервисов.
func WithExternalRequestTimeoutInterval(timeoutInterval time.Duration) func(*mapProperties) {
	return func(p *mapProperties) {
		p.mp.Store(propertyExternalRequestTimeoutInterval, timeoutInterval)
	}
}

// ExternalRequestTimeoutInterval геттер интервала ожидания ответа от внешних gRPC-сервисов.
func (p *mapProperties) ExternalRequestTimeoutInterval() time.Duration {
	if a, ok := p.mp.Load(propertyExternalRequestTimeoutInterval); ok {
		if timeoutInterval, ok := a.(time.Duration); ok {
			return timeoutInterval
		}
	}
	return 0
}

// WithFlags — Флаги.
func WithFlags(flags map[string]interface{}) func(*mapProperties) {
	return func(p *mapProperties) {
		if flags != nil {
			p.mp.Store(propertyFlags, flags)
		}
	}
}

// Flags — флаги командной строки.
func (p *mapProperties) Flags() map[string]interface{} {
	if f, ok := p.mp.Load(propertyFlags); ok {
		if flags, ok := f.(map[string]interface{}); ok {
			return flags
		}
	}
	return nil
}

// WithGRPCAddress — адрес gRPC сервера.
func WithGRPCAddress(address string) func(*mapProperties) {
	return func(p *mapProperties) {
		if address != "" {
			p.mp.Store(propertyGRPCAddress, address)
		}
	}
}

// GRPCAddress геттер адреса gRPC сервера.
func (p *mapProperties) GRPCAddress() string {
	if a, ok := p.mp.Load(propertyGRPCAddress); ok {
		if address, ok := a.(string); ok {
			return address
		}
	}
	return ""
}

// WithGRPCTransportCredentials — TLS реквизиты для gRPC-сервера.
func WithGRPCTransportCredentials(tCredentials credentials.TransportCredentials) func(*mapProperties) {
	return func(p *mapProperties) {
		if tCredentials != nil {
			p.mp.Store(propertyGRPCTransportCredentials, tCredentials)
		}
	}
}

// GRPCTransportCredentials геттер TLS реквизитов для gRPC-сервера.
func (p *mapProperties) GRPCTransportCredentials() credentials.TransportCredentials {
	if c, ok := p.mp.Load(propertyGRPCTransportCredentials); ok {
		if tCredentials, ok := c.(credentials.TransportCredentials); ok {
			return tCredentials
		}
	}
	return nil
}

// WithHTTPAddress — адрес HTTP сервера.
func WithHTTPAddress(address string) func(*mapProperties) {
	return func(p *mapProperties) {
		if address != "" {
			p.mp.Store(propertyHTTPAddress, address)
		}
	}
}

// HTTPAddress геттер адреса HTTP сервера.
func (p *mapProperties) HTTPAddress() string {
	if a, ok := p.mp.Load(propertyHTTPAddress); ok {
		if address, ok := a.(string); ok {
			return address
		}
	}
	return ""
}

// WithHTTPTLSConfig — TLS конфигурация для HTTP-сервера.
func WithHTTPTLSConfig(tCredentials *tls.Config) func(*mapProperties) {
	return func(p *mapProperties) {
		if tCredentials != nil {
			p.mp.Store(propertyHTTPHTTPTLSConfig, tCredentials)
		}
	}
}

// HTTPTLSConfig геттер TLS конфигурации для HTTP-сервера.
func (p *mapProperties) HTTPTLSConfig() *tls.Config {
	if c, ok := p.mp.Load(propertyHTTPHTTPTLSConfig); ok {
		if tCredentials, ok := c.(*tls.Config); ok {
			return tCredentials
		}
	}
	return nil
}

// WithJwtExpiresIn — Утверждение «exp» (время истечения срока действия)
// определяет время истечения срока действия или после чего JWT
// НЕ ДОЛЖЕН приниматься в обработку.
func WithJwtExpiresIn(jwtExpiresIn time.Duration) func(*mapProperties) {
	return func(p *mapProperties) {
		if jwtExpiresIn > 0 {
			p.mp.Store(propertyJwtExpiresIn, jwtExpiresIn)
		}
	}
}

// JwtExpiresIn геттер времени истечения срока действия JWT.
func (p *mapProperties) JwtExpiresIn() time.Duration {
	if a, ok := p.mp.Load(propertyJwtExpiresIn); ok {
		if jwtExpiresIn, ok := a.(time.Duration); ok {
			return jwtExpiresIn
		}
	}
	return 0
}

// WithJwtMaxAgeSec — определяет время жизни куки в секундах.
func WithJwtMaxAgeSec(maxAge int) func(*mapProperties) {
	return func(p *mapProperties) {
		if maxAge > 0 {
			p.mp.Store(propertyJwtMaxAgeSec, maxAge)
		}
	}
}

// JwtMaxAgeSec геттер времени жизни куки в секундах.
func (p *mapProperties) JwtMaxAgeSec() int {
	if a, ok := p.mp.Load(propertyJwtMaxAgeSec); ok {
		if maxAge, ok := a.(int); ok {
			return maxAge
		}
	}
	return 0
}

// WithJwtSecret — секрет для подписи JWТокена.
func WithJwtSecret(secret string) func(*mapProperties) {
	return func(p *mapProperties) {
		if secret != "" {
			p.mp.Store(propertyJwtSecret, secret)
		}
	}
}

// JwtSecret геттер секрета для подписи JWТокена.
func (p *mapProperties) JwtSecret() string {
	if a, ok := p.mp.Load(propertyJwtSecret); ok {
		if secret, ok := a.(string); ok {
			return secret
		}
	}
	return ""
}

// WithLogger — логгер приложения.
func WithLogger(logger *slog.Logger) func(*mapProperties) {
	return func(p *mapProperties) {
		if logger != nil {
			p.mp.Store(propertyLogger, logger)
		}
	}
}

// Logger получение логгера приложения.
func (p *mapProperties) Logger() *slog.Logger {
	if a, ok := p.mp.Load(propertyLogger); ok {
		if logger, ok := a.(*slog.Logger); ok {
			return logger
		}
	}
	return slog.Default()
}

// MongodbPool пул подключения к базе данных MongoDB.
func (p *mapProperties) MongodbPool() *tool.MongoPool {
	if p, ok := p.mp.Load(propertyMongodbPool); ok {
		if pool, ok := p.(*tool.MongoPool); ok {
			return pool
		}
	}
	return nil
}

func (p *mapProperties) SlogJSON() bool {
	return slogJSON(p.Flags())
}

// WithUpkRSAPrivateKey — RSA ключ для дешифрации секрета
// который применяется в симметричном шифровании UPK (User Personal Key).
func WithUpkRSAPrivateKey(privateKey *rsa.PrivateKey) func(*mapProperties) {
	return func(p *mapProperties) {
		if privateKey != nil {
			p.mp.Store(propertyUpkRSAPrivateKey, privateKey)
		}
	}
}

// UpkRSAPrivateKey геттер RSA ключа для дешифрации секрета.
func (p *mapProperties) UpkRSAPrivateKey() *rsa.PrivateKey {
	if a, ok := p.mp.Load(propertyUpkRSAPrivateKey); ok {
		if PrivateKey, ok := a.(*rsa.PrivateKey); ok {
			return PrivateKey
		}
	}
	return nil
}

// WithUpkRSAPublicKey — RSA ключ для шифрования секрета
// который применяется в симметричном шифровании UPK (User Personal Key).
func WithUpkRSAPublicKey(publicKey *rsa.PublicKey) func(*mapProperties) {
	return func(p *mapProperties) {
		if publicKey != nil {
			p.mp.Store(propertyUpkRSAPublicKey, publicKey)
		}
	}
}

// UpkRSAPublicKey геттер RSA ключа для шифрования секрета.
func (p *mapProperties) UpkRSAPublicKey() *rsa.PublicKey {
	if a, ok := p.mp.Load(propertyUpkRSAPublicKey); ok {
		if publicKey, ok := a.(*rsa.PublicKey); ok {
			return publicKey
		}
	}
	return nil
}

// WithUpkSecretKey — секрет симметричного шифрования UPK (User Personal Key).
func WithUpkSecretKey(secretKey []byte) func(*mapProperties) {
	return func(p *mapProperties) {
		if secretKey != nil {
			p.mp.Store(propertyUpkSecretKey, secretKey)
		}
	}
}

// UpkSecretKey геттер секрета симметричного шифрования UPK.
func (p *mapProperties) UpkSecretKey() []byte {
	if a, ok := p.mp.Load(propertyUpkSecretKey); ok {
		if secretKey, ok := a.([]byte); ok {
			return secretKey
		}
	}
	return nil
}

// withDBPool — пул подключения к базе данных PostgreSQL.
func withDBPool(pool *pgxpool.Pool) func(*mapProperties) {
	return func(p *mapProperties) {
		if pool != nil {
			p.mp.Store(propertyDBPool, pool)
		}
	}
}

// DBPool геттер пула подключения к базе данных PostgreSQL.
func (p *mapProperties) DBPool() *pgxpool.Pool {
	if p, ok := p.mp.Load(propertyDBPool); ok {
		if pool, ok := p.(*pgxpool.Pool); ok {
			return pool
		}
	}
	return nil
}

// withMongodbPool — пул подключения к базе данных MongoDB.
func withMongodbPool(pool *tool.MongoPool) func(*mapProperties) {
	return func(p *mapProperties) {
		if pool != nil {
			p.mp.Store(propertyMongodbPool, pool)
		}
	}
}

func (p *mapProperties) OutboundIP() net.IP {
	return nil
}

func (p *mapProperties) String() string {
	format := `
%s
CacheExpire: %v
CacheGCInterval: %v
DBPool: %v
Environments: %v
ExternalAssetGRPCAddress: %s
ExternalAuthGRPCAddress: %s
ExternalRequestTimeoutInterval: %d
Flags: %v
GRPCAddress: %s
GRPCTransportCredentials: %v
HTTPAddress: %s
HTTPTransportCredentials: %v
JwtExpiresIn: %v
JwtMaxAgeSec: %d
JwtSecret: %s
MongodbPool: %v
OutboundIP: %v
UpkRSAPrivateKey: %v
UpkRSAPublicKey: %v
UpkSecretKey: %v
`
	return fmt.Sprintf(format,
		p.Config(),
		p.CacheExpire(),
		p.CacheGCInterval(),
		p.DBPool(),
		p.Environments(),
		p.ExternalAssetGRPCAddress(),
		p.ExternalAuthGRPCAddress(),
		p.ExternalRequestTimeoutInterval(),
		p.Flags(),
		p.GRPCAddress(),
		p.GRPCTransportCredentials(),
		p.HTTPAddress(),
		p.HTTPTLSConfig(),
		p.JwtExpiresIn(),
		p.JwtMaxAgeSec(),
		p.JwtSecret(),
		p.MongodbPool(),
		p.OutboundIP(),
		p.UpkRSAPrivateKey(),
		p.UpkRSAPublicKey(),
		p.UpkSecretKey(),
	)
}

func slogJSON(flags map[string]interface{}) bool {
	if sj, ok := flags[flagSlogJson]; ok {
		if slogJSON, ok := sj.(*bool); ok {
			return *slogJSON
		}
	}
	return false
}

func getProperties(opts ...func(*mapProperties)) *mapProperties {

	var property = new(mapProperties)

	// вызываем все указанные функции для установки параметров
	for _, opt := range opts {
		opt(property)
	}

	return property
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
