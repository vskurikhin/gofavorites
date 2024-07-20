/*
 * This file was last modified at 2024-07-20 13:37 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * properties.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"net"
	"sync"
	"time"
)

const (
	propertyCacheExpire                    = "cache-expire"
	propertyConfig                         = "config"
	propertyDBPool                         = "db-pool"
	propertyEnvironments                   = "environments"
	propertyExternalAssetGRPCAddress       = "external-asset-grpc-address"
	propertyExternalAuthGRPCAddress        = "external-auth-grpc-address"
	propertyExternalRequestTimeoutInterval = "external-request-timeout-interval"
	propertyFlags                          = "flags"
	propertyGRPCAddress                    = "grpc-address"
	propertyGRPCTransportCredentials       = "grpc-transport-credentials"
)

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
	OutboundIP() net.IP
}

type mapProperties struct {
	mp sync.Map
}

var properties Properties = (*mapProperties)(nil)
var once = new(sync.Once)

// GetProperties — свойства преобразованные из конфигурации и окружения.
func GetProperties() Properties {

	once.Do(func() {
		yml, err := LoadConfig(".")
		tool.IfErrorThenPanic(err)
		env, err := getEnvironments()
		tool.IfErrorThenPanic(err)
		flm := makeFlagsParse()

		cacheExpire, err := getCacheExpire(flm, env, yml)
		slog.Warn(MSG+" GetProperties", "cacheExpire", cacheExpire, "err", err)
		cacheGCInterval, err := getCacheGCInterval(flm, env, yml)
		slog.Warn(MSG+" GetProperties", "cacheGCInterval", cacheGCInterval, "err", err)
		dbPool, err := makeDBPool(flm, env, yml)
		slog.Warn(MSG+" GetProperties", "dbDisable", err)
		grpcAddress, err := getGRPCAddress(flm, env, yml)
		slog.Warn(MSG+" GetProperties", "grpcAddress", grpcAddress, "err", err)
		tCredentials, err := getGRPCTransportCredentials(flm, env, yml)
		slog.Warn(MSG+" GetProperties", "grpcTransportCredentials", tCredentials, "err", err)
		assetGRPCAddress, err := getExternalAssetGRPCAddress(flm, env, yml)
		slog.Warn(MSG+" GetProperties", "assetGRPCAddress", assetGRPCAddress, "err", err)
		authGRPCAddress, err := getExternalAuthGRPCAddress(flm, env, yml)
		slog.Warn(MSG+" GetProperties", "authGRPCAddress", authGRPCAddress, "err", err)
		requestTimeoutInterval, err := getExternalRequestTimeoutInterval(flm, env, yml)
		slog.Warn(MSG+" GetProperties", "requestTimeoutInterval", requestTimeoutInterval, "err", err)

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
			WithGRPCTransportCredentials(tCredentials),
		)
	})
	return properties
}

// WithCacheExpire — TODO.
func WithCacheExpire(cacheExpire time.Duration) func(*mapProperties) {
	return func(p *mapProperties) {
		if cacheExpire > 0 {
			p.mp.Store(propertyCacheExpire, cacheExpire)
		}
	}
}

// CacheExpire — TODO.
func (p *mapProperties) CacheExpire() time.Duration {
	if a, ok := p.mp.Load(propertyCacheExpire); ok {
		if cacheExpire, ok := a.(time.Duration); ok {
			return cacheExpire
		}
	}
	return 0
}

// WithCacheGCInterval — TODO.
func WithCacheGCInterval(cacheGCInterval time.Duration) func(*mapProperties) {
	return func(p *mapProperties) {
		if cacheGCInterval > 0 {
			p.mp.Store(propertyCacheExpire, cacheGCInterval)
		}
	}
}

// CacheGCInterval — TODO.
func (p *mapProperties) CacheGCInterval() time.Duration {
	if a, ok := p.mp.Load(propertyCacheExpire); ok {
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

// Environments — флаги командной строки.
func (p *mapProperties) Environments() environments {
	if f, ok := p.mp.Load(propertyEnvironments); ok {
		if env, ok := f.(environments); ok {
			return env
		}
	}
	return environments{}
}

// WithExternalAssetGRPCAddress — Окружение.
func WithExternalAssetGRPCAddress(address string) func(*mapProperties) {
	return func(p *mapProperties) {
		p.mp.Store(propertyExternalAssetGRPCAddress, address)
	}
}

// ExternalAssetGRPCAddress — флаги командной строки.
func (p *mapProperties) ExternalAssetGRPCAddress() string {
	if a, ok := p.mp.Load(propertyExternalAssetGRPCAddress); ok {
		if address, ok := a.(string); ok {
			return address
		}
	}
	return ""
}

// WithExternalAuthGRPCAddress — Окружение.
func WithExternalAuthGRPCAddress(address string) func(*mapProperties) {
	return func(p *mapProperties) {
		p.mp.Store(propertyExternalAuthGRPCAddress, address)
	}
}

// ExternalAuthGRPCAddress — флаги командной строки.
func (p *mapProperties) ExternalAuthGRPCAddress() string {
	if a, ok := p.mp.Load(propertyExternalAuthGRPCAddress); ok {
		if address, ok := a.(string); ok {
			return address
		}
	}
	return ""
}

// WithExternalRequestTimeoutInterval — Окружение.
func WithExternalRequestTimeoutInterval(timeoutInterval time.Duration) func(*mapProperties) {
	return func(p *mapProperties) {
		p.mp.Store(propertyExternalRequestTimeoutInterval, timeoutInterval)
	}
}

// ExternalRequestTimeoutInterval — флаги командной строки.
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

// GRPCAddress — геттер адреса gRPC сервера.
func (p *mapProperties) GRPCAddress() string {
	if a, ok := p.mp.Load(propertyGRPCAddress); ok {
		if address, ok := a.(string); ok {
			return address
		}
	}
	return ""
}

// WithGRPCTransportCredentials — TODO.
func WithGRPCTransportCredentials(tCredentials credentials.TransportCredentials) func(*mapProperties) {
	return func(p *mapProperties) {
		if tCredentials != nil {
			p.mp.Store(propertyGRPCTransportCredentials, tCredentials)
		}
	}
}

// GRPCTransportCredentials — геттер TODO.
func (p *mapProperties) GRPCTransportCredentials() credentials.TransportCredentials {
	if c, ok := p.mp.Load(propertyGRPCTransportCredentials); ok {
		if tCredentials, ok := c.(credentials.TransportCredentials); ok {
			return tCredentials
		}
	}
	return nil
}

// withDBPool — Флаги.
func withDBPool(pool *pgxpool.Pool) func(*mapProperties) {
	return func(p *mapProperties) {
		if pool != nil {
			p.mp.Store(propertyDBPool, pool)
		}
	}
}

func (p *mapProperties) DBPool() *pgxpool.Pool {
	if p, ok := p.mp.Load(propertyDBPool); ok {
		if pool, ok := p.(*pgxpool.Pool); ok {
			return pool
		}
	}
	return nil
}

func (p *mapProperties) OutboundIP() net.IP {
	return nil
}

func (p *mapProperties) String() string {
	format := `
%s
DBPool: %v
Environments: %v
ExternalAssetGRPCAddress: %s
ExternalAuthGRPCAddress: %s
ExternalRequestTimeoutInterval: %d
Flags: %v
GRPCAddress: %s
GRPCTransportCredentials: %v
OutboundIP: %v
`
	return fmt.Sprintf(format,
		p.Config(),
		p.DBPool(),
		p.Environments(),
		p.ExternalAssetGRPCAddress(),
		p.ExternalAuthGRPCAddress(),
		p.ExternalRequestTimeoutInterval(),
		p.Flags(),
		p.GRPCAddress(),
		p.GRPCTransportCredentials(),
		p.OutboundIP(),
	)
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
