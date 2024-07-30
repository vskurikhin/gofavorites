/*
 * This file was last modified at 2024-08-03 12:36 by Victor N. Skurikhin.
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
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/vskurikhin/gofavorites/internal/alog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"google.golang.org/grpc/credentials"
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

func getExternalAssetGRPCAddress(flm map[string]interface{}, env *environments, yml Config) (string, error) {
	return stringsAddressPrepareProperty(
		flagExternalAssetGRPCAddress,
		flm[flagExternalAssetGRPCAddress],
		env.ExternalAssetGRPCAddress,
		yml.ExternalAssetGRPCAddress(),
	)
}

func getExternalAuthGRPCAddress(flm map[string]interface{}, env *environments, yml Config) (string, error) {
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

func getGRPCAddress(flm map[string]interface{}, env *environments, yml Config) (string, error) {
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
) (credentials.TransportCredentials, error) {
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

func getHTTPAddress(flm map[string]interface{}, env *environments, yml Config) (string, error) {
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

func getHTTPTLSConfig(flm map[string]interface{}, env *environments, yml Config) (*tls.Config, error) {
	if yml.HTTPTLSEnabled() {
		return serverTLSConfigPrepareProperty(
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

func getRSAPrivateKey(flm map[string]interface{}, env *environments, yml Config) (*rsa.PrivateKey, error) {
	return loadRSAPrivateKey(
		flagUpkPrivateKeyFile,
		flm[flagUpkPrivateKeyFile],
		env.UpkPrivateKeyFile,
		yml.UpkRSAPrivateKeyFile(),
	)
}

func getRSAPublicKey(flm map[string]interface{}, env *environments, yml Config) (*rsa.PublicKey, error) {
	return loadRSAPublicKey(
		flagUpkPublicKeyFile,
		flm[flagUpkPublicKeyFile],
		env.UpkPublicKeyFile,
		yml.UpkRSAPublicKeyFile(),
	)
}

func getUpkSecret(flm map[string]interface{}, env *environments, yml Config) (string, error) {
	return stringPrepareProperty(
		flagUpkSecret,
		flm[flagUpkSecret],
		env.UpkSecret,
		yml.UpkSecret(),
	)
}

func getUpkSecretKey(
	flm map[string]interface{},
	env *environments,
	yml Config,
	rsaPrivateKey *rsa.PrivateKey,
) ([]byte, error) {

	secret, err := getUpkSecret(flm, env, yml)

	if err != nil {
		return nil, err
	}
	encrypt, err := base64.StdEncoding.DecodeString(secret)

	if err != nil {
		return nil, err
	}
	secretKey, err := tool.DecryptRSA(rsaPrivateKey, encrypt)

	if err != nil {
		return nil, err
	}
	result := make([]byte, 32)
	copy(result, secretKey)

	return result, nil
}

func intPrepareProperty(name string, flag interface{}, env int, yaml int) (int, error) {

	var result int
	var err error

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

func loadRSAPrivateKey(name string, flag interface{}, env string, yaml string) (*rsa.PrivateKey, error) {

	err, fileName := getFileName(name, flag, env, yaml)

	if err != nil {
		return nil, err
	}
	return tool.LoadPrivateKey(fileName), err
}

func loadRSAPublicKey(name string, flag interface{}, env string, yaml string) (*rsa.PublicKey, error) {

	err, fileName := getFileName(name, flag, env, yaml)

	if err != nil {
		return nil, err
	}
	return tool.LoadPublicKey(fileName), err
}

func getFileName(name string, flag interface{}, env string, yaml string) (error, string) {
	var err error
	var fileName string
	getFlag := func() {
		if a, ok := flag.(*string); !ok {
			err = fmt.Errorf("bad value")
		} else {
			fileName = *a
		}
	}
	if yaml != "" {
		fileName = yaml
	}
	if env != "" {
		fileName = env
	} else if fileName == "" {
		getFlag()
	}
	setIfFlagChanged(name, getFlag)
	return err, fileName
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
		slog.Debug(MSG+"makeDBPool", "DatabaseDSN", dsn)

		return tool.DBConnect(dsn), nil
	}
	return nil, fmt.Errorf("connect to DataBase disabled")
}

func makeMongodbPool(flm map[string]interface{}, env *environments, yml Config) (*tool.MongoPool, error) {
	if yml.MongoEnabled() {

		dsn := fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/%s?authSource=admin",
			yml.MongoUserName(), yml.MongoUserPassword(), yml.MongoHost(), yml.MongoPort(), yml.MongoName(),
		)
		getFlagDatabaseDSN := func() {
			dsn = *(flm[flagMongodbDSN].(*string))
		}
		if env.MongodbDSN != "" {
			dsn = env.MongodbDSN
		} else if dsn == "mongodb://:@:/?authSource=admin" {
			getFlagDatabaseDSN()
		}
		setIfFlagChanged(flagMongodbDSN, getFlagDatabaseDSN)
		slog.Debug(MSG+"makeMongodbPool", "MongodbDSN", dsn)

		return tool.MongodbConnect(dsn), nil
	}
	return nil, fmt.Errorf("connect to MongoDB disabled")
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

func serverAddressPrepareProperty(
	name string,
	flm map[string]interface{},
	envAddress []string,
	ymlAddress string,
	ymlPort int,
) (string, error) {

	var address string
	var err error

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

func setupLogger(slogJSON bool) *slog.Logger {
	if slogJSON {
		alog.NewLogger(alog.NewHandlerJSON(os.Stdout, nil))
	} else {
		opts := alog.PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		alog.NewLogger(alog.NewPrettyHandlerText(os.Stdout, opts))
	}
	return alog.GetLogger()
}

func stringsAddressPrepareProperty(name string, flag interface{}, env []string, yaml string) (string, error) {

	var result string
	var err error

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

func stringPrepareProperty(name string, flag interface{}, env string, yaml string) (string, error) {

	var result string
	var err error

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

func serverTLSConfigPrepareProperty(
	nameCertFile string,
	nameKeyFile string,
	flm map[string]interface{},
	envTLSCertFile string,
	envTLSKeyFile string,
	ymlTLSCertFile string,
	ymlTLSKeyFile string,
) (tConfig *tls.Config, err error) {
	certFile, keyFile := ymlTLSCertFile, ymlTLSKeyFile
	getFlagCertFile := func() {
		if cf, ok := flm[nameCertFile].(*string); !ok {
			err = fmt.Errorf("bad value of %s : %v", flagGRPCCertFile, flm[flagGRPCCertFile])
		} else {
			certFile = *cf
		}
	}
	getFlagKeyFile := func() {
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
		getFlagCertFile()
	}
	if keyFile == "" {
		getFlagKeyFile()
	}
	setIfFlagChanged(nameCertFile, getFlagCertFile)
	setIfFlagChanged(nameKeyFile, getFlagKeyFile)
	if err != nil {
		return nil, err
	}
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cer}}, nil
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

func timePrepareProperty(name string, flag interface{}, env time.Duration, yaml time.Duration) (time.Duration, error) {

	var result time.Duration
	var err error

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

func toTimePrepareProperty(name string, flag interface{}, env int, yaml int, scale time.Duration) (time.Duration, error) {

	var result time.Duration
	var err error

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
