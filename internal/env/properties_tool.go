/*
 * This file was last modified at 2024-07-22 23:58 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * properties_tool.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"bytes"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"strconv"
	"time"
)

var ErrEmptyAddress = fmt.Errorf("can't configure epmty address")

func getCacheExpire(flm map[string]interface{}, env *environments, yml Config) (time.Duration, error) {
	return toTimePrepareProperty(
		flagCacheExpireMs,
		flm[flagCacheExpireMs],
		env.CacheExpireMs,
		yml.CacheExpireMs(),
		time.Millisecond,
	)
}

func getCacheGCInterval(flm map[string]interface{}, env *environments, yml Config) (time.Duration, error) {
	return toTimePrepareProperty(
		flagCacheGCIntervalSec,
		flm[flagCacheGCIntervalSec],
		env.CacheGCIntervalSec,
		yml.CacheGCIntervalSec(),
		time.Second,
	)
}

func getExternalAssetGRPCAddress(flm map[string]interface{}, env *environments, yml Config) (result string, err error) {
	return stringsAddressPrepareProperty(
		flagExternalAssetGRPCAddress,
		flm[flagExternalAssetGRPCAddress],
		env.ExternalAssetGRPCAddress,
		yml.ExternalAssetGRPCAddress(),
	)
}

func getExternalAuthGRPCAddress(flm map[string]interface{}, env *environments, yml Config) (result string, err error) {
	return stringsAddressPrepareProperty(
		flagExternalAuthGRPCAddress,
		flm[flagExternalAuthGRPCAddress],
		env.ExternalAuthGRPCAddress,
		yml.ExternalAuthGRPCAddress(),
	)
}

func getExternalRequestTimeoutInterval(flm map[string]interface{}, env *environments, yml Config) (time.Duration, error) {
	return toTimePrepareProperty(
		flagExternalRequestTimeoutInterval,
		flm[flagExternalRequestTimeoutInterval],
		env.ExternalRequestTimeoutInterval,
		yml.ExternalRequestTimeoutInterval(),
		time.Millisecond,
	)
}

func getGRPCAddress(flm map[string]interface{}, env *environments, yml Config) (address string, err error) {
	if yml.GRPCEnabled() {
		return serverAddressPrepareProperty(
			flagGRPCAddress, flm,
			env.GRPCAddress,
			yml.GRPCAddress(),
			yml.GRPCPort())
	}
	return "", fmt.Errorf("gRPC server disabled")
}

func getGRPCTransportCredentials(
	flm map[string]interface{},
	env *environments,
	yml Config,
) (tCredentials credentials.TransportCredentials, err error) {
	if yml.GRPCEnabled() {
		return serverTransportCredentialsPrepareProperty(
			flagGRPCCertFile,
			flagGRPCKeyFile, flm,
			env.GRPCCertFile,
			env.GRPCKeyFile,
			yml.GRPCTLSCertFile(),
			yml.GRPCTLSKeyFile(),
		)
	}
	return nil, fmt.Errorf("gRPC server disabled")
}

func getHTTPAddress(flm map[string]interface{}, env *environments, yml Config) (address string, err error) {
	if yml.GRPCEnabled() {
		return serverAddressPrepareProperty(
			flagHTTPAddress, flm,
			env.HTTPAddress,
			yml.HTTPAddress(),
			yml.HTTPPort(),
		)
	}
	return "", fmt.Errorf("HTTP server disabled")
}

func getHTTPTransportCredentials(
	flm map[string]interface{},
	env *environments,
	yml Config,
) (tCredentials credentials.TransportCredentials, err error) {
	if yml.GRPCEnabled() {
		return serverTransportCredentialsPrepareProperty(
			flagHTTPCertFile,
			flagHTTPKeyFile, flm,
			env.HTTPCertFile,
			env.HTTPKeyFile,
			yml.HTTPTLSCertFile(),
			yml.HTTPTLSKeyFile(),
		)
	}
	return nil, fmt.Errorf("HTTP server disabled")
}

func getJwtExpiresIn(flm map[string]interface{}, env *environments, yml Config) (time.Duration, error) {
	return timePrepareProperty(
		flagJwtExpiresIn,
		flm[flagJwtExpiresIn],
		env.JwtExpiresIn,
		yml.JwtExpiresIn(),
	)
}

func getJwtMaxAgeSec(flm map[string]interface{}, env *environments, yml Config) (int, error) {
	return intPrepareProperty(
		flagJwtMaxAgeSec,
		flm[flagJwtMaxAgeSec],
		env.JwtMaxAge,
		yml.JwtMaxAgeSec(),
	)
}

func getJwtSecret(flm map[string]interface{}, env *environments, yml Config) (string, error) {
	return stringPrepareProperty(
		flagJwtSecret,
		flm[flagJwtSecret],
		env.JwtSecret,
		yml.JwtSecret(),
	)
}

func makeDBPool(flm map[string]interface{}, env *environments, yml Config) (*pgxpool.Pool, error) {
	if yml.DBEnabled() {

		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			yml.DBUserName(), yml.DBUserPassword(), yml.DBHost(), yml.DBPort(), yml.DBName(),
		)
		getFlagDatabaseDSN := func() {
			dsn = *(flm[flagDatabaseDSN].(*string))
		}
		if env.DataBaseDSN != "" {
			dsn = env.DataBaseDSN
		} else if dsn == "postgres://:@:/?sslmode=disable" {
			getFlagDatabaseDSN()
		}
		setIfFlagChanged(flagDatabaseDSN, getFlagDatabaseDSN)
		slog.Warn(MSG, "DatabaseDSN", dsn)

		return tool.DBConnect(dsn), nil
	}
	return nil, fmt.Errorf("connect to DataBase disabled")
}

func parseEnvAddress(address []string) string {

	port, err := strconv.Atoi(address[len(address)-1])
	tool.IfErrorThenPanic(err)
	var bb bytes.Buffer

	if len(address) > 1 {
		for i := 0; i < len(address)-1; i++ {
			bb.WriteString(address[i])
			bb.WriteRune(':')
		}
	} else {
		bb.WriteRune(':')
	}
	return fmt.Sprintf("%s%d", bb.String(), port)
}

func intPrepareProperty(
	name string,
	flag interface{},
	env int,
	yaml int,
) (result int, err error) {

	getFlag := func() {
		if a, ok := flag.(*int); !ok {
			err = fmt.Errorf("bad value")
		} else {
			result = *a
		}
	}
	if yaml > 0 {
		result = yaml
	}
	if env > 0 {
		result = env
	} else if result == 0 {
		getFlag()
	}
	setIfFlagChanged(name, getFlag)

	return result, err
}

func serverAddressPrepareProperty(
	name string,
	flm map[string]interface{},
	envAddress []string,
	ymlAddress string,
	ymlPort int,
) (address string, err error) {

	getFlagAddress := func() {
		if a, ok := flm[name].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", name, flm[name])
		} else {
			address = *a
		}
	}
	address = fmt.Sprintf("%s:%d", ymlAddress, ymlPort)

	if len(envAddress) > 0 {
		address = parseEnvAddress(envAddress)
	} else if ymlAddress == "" && ymlPort == 0 {
		getFlagAddress()
	}
	setIfFlagChanged(name, getFlagAddress)

	if address == "" {
		err = ErrEmptyAddress
	}
	return address, err
}

func stringsAddressPrepareProperty(
	name string,
	flag interface{},
	env []string,
	yaml string,
) (result string, err error) {

	getFlag := func() {
		if a, ok := flag.(*string); !ok {
			err = fmt.Errorf("bad value")
		} else {
			result = *a
		}
	}
	if yaml != "" {
		result = yaml
	}
	if len(env) > 0 {
		result = parseEnvAddress(env)
	} else if result == "" {
		getFlag()
	}
	setIfFlagChanged(name, getFlag)

	return result, err
}

func stringPrepareProperty(
	name string,
	flag interface{},
	env string,
	yaml string,
) (result string, err error) {

	getFlag := func() {
		if a, ok := flag.(*string); !ok {
			err = fmt.Errorf("bad value")
		} else {
			result = *a
		}
	}
	if yaml != "" {
		result = yaml
	}
	if env != "" {
		result = env
	} else if result == "" {
		getFlag()
	}
	setIfFlagChanged(name, getFlag)

	return result, err
}

func serverTransportCredentialsPrepareProperty(
	nameCertFile string,
	nameKeyFile string,
	flm map[string]interface{},
	envTLSCertFile string,
	envTLSKeyFile string,
	ymlTLSCertFile string,
	ymlTLSKeyFile string,
) (tCredentials credentials.TransportCredentials, err error) {

	certFile, keyFile := ymlTLSCertFile, ymlTLSKeyFile
	getFlagGRPCCertFile := func() {
		if cf, ok := flm[nameCertFile].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", flagGRPCCertFile, flm[flagGRPCCertFile])
		} else {
			certFile = *cf
		}
	}
	getFlagGRPCKeyFile := func() {
		if kf, ok := flm[nameKeyFile].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", flagGRPCKeyFile, flm[flagGRPCKeyFile])
		} else {
			keyFile = *kf
		}
	}
	if envTLSCertFile != "" {
		certFile = envTLSCertFile
	}
	if envTLSKeyFile != "" {
		keyFile = envTLSKeyFile
	}
	if certFile == "" {
		getFlagGRPCCertFile()
	}
	if keyFile == "" {
		getFlagGRPCKeyFile()
	}
	setIfFlagChanged(nameCertFile, getFlagGRPCCertFile)
	setIfFlagChanged(nameKeyFile, getFlagGRPCKeyFile)
	if err != nil {
		return nil, err
	}
	return tool.LoadServerTLSCredentials(certFile, keyFile)
}

func timePrepareProperty(
	name string,
	flag interface{},
	env time.Duration,
	yaml time.Duration,
) (result time.Duration, err error) {

	getFlag := func() {
		if a, ok := flag.(*time.Duration); !ok {
			err = fmt.Errorf("bad value")
		} else {
			result = *a
		}
	}
	if yaml > 0 {
		result = yaml
	}
	if env > 0 {
		result = env
	} else if result == 0 {
		getFlag()
	}
	setIfFlagChanged(name, getFlag)

	return result, err
}

func toTimePrepareProperty(
	name string,
	flag interface{},
	env int,
	yaml int,
	scale time.Duration,
) (result time.Duration, err error) {

	getFlag := func() {
		if a, ok := flag.(*int); !ok {
			err = fmt.Errorf("bad value")
		} else {
			result = time.Duration(*a) * scale
		}
	}
	if yaml > 0 {
		result = time.Duration(yaml) * scale
	}
	if env > 0 {
		result = time.Duration(env) * scale
	} else if result == 0 {
		getFlag()
	}
	setIfFlagChanged(name, getFlag)

	return result, err
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
