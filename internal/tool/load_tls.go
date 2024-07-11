/*
 * This file was last modified at 2024-07-11 11:30 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * load_tls.go
 * $Id$
 */
//!+

// Package tool TODO.
package tool

import (
	"crypto/tls"

	"google.golang.org/grpc/credentials"
)

func LoadServerTLSCredentials(certFile, keyFile string) (credentials.TransportCredentials, error) {

	// Загрузка серверного сертификата и закрытого ключа.
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Создание учётных данных для конфигурации TLS.
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
