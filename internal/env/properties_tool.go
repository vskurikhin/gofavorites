/*
 * This file was last modified at 2024-07-11 11:30 by Victor N. Skurikhin.
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
)

func getGRPCAddress(flm map[string]interface{}, env *environments, yml Config) (address string, err error) {
	if yml.GRPCEnabled() {

		getFlagGRPCAddress := func() {
			if a, ok := flm[flagGRPCAddress].(*string); !ok {
				err = fmt.Errorf("bad value of %s : %v", flagGRPCAddress, flm[flagGRPCAddress])
			} else {
				address = *a
			}
		}
		address = fmt.Sprintf("%s:%s:%d", yml.GRPCProto(), yml.GRPCAddress(), yml.GRPCPort())

		if len(env.GRPCAddress) > 0 {
			address = parseEnvAddress(env.GRPCAddress)
		} else if yml.GRPCProto() == "" && yml.GRPCAddress() == "" && yml.GRPCPort() == 0 {
			getFlagGRPCAddress()
		}
		setIfFlagChanged(flagGRPCAddress, getFlagGRPCAddress)

		if address == "" {
			err = fmt.Errorf("can't configure gRPC address : %s", address)
		}
		return address, err
	}
	return "", fmt.Errorf("gRPC server disabled")
}

func getGRPCTransportCredentials(
	flm map[string]interface{},
	env *environments,
	yml Config,
) (tCredentials credentials.TransportCredentials, err error) {
	if yml.GRPCEnabled() {

		certFile, keyFile := yml.GRPCTLSCertFile(), yml.GRPCTLSKeyFile()
		getFlagGRPCCertFile := func() {
			if cf, ok := flm[flagGRPCCertFile].(*string); !ok {
				err = fmt.Errorf("bad value of %s : %v", flagGRPCCertFile, flm[flagGRPCCertFile])
			} else {
				certFile = *cf
			}
		}
		getFlagGRPCKeyFile := func() {
			if kf, ok := flm[flagGRPCKeyFile].(*string); !ok {
				err = fmt.Errorf("bad value of %s : %v", flagGRPCKeyFile, flm[flagGRPCKeyFile])
			} else {
				keyFile = *kf
			}
		}
		if env.GRPCCertFile != "" {
			certFile = env.GRPCCertFile
		}
		if env.GRPCKeyFile != "" {
			keyFile = env.GRPCKeyFile
		}
		if certFile == "" {
			getFlagGRPCCertFile()
		}
		if keyFile == "" {
			getFlagGRPCKeyFile()
		}
		setIfFlagChanged(flagGRPCCertFile, getFlagGRPCCertFile)
		setIfFlagChanged(flagGRPCKeyFile, getFlagGRPCKeyFile)
		if err != nil {
			return nil, err
		}
		return tool.LoadServerTLSCredentials(certFile, keyFile)
	}
	return nil, fmt.Errorf("gRPC server disabled")
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

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
