/*
 * This file was last modified at 2024-07-11 09:38 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * load_tls_test.go
 * $Id$
 */
//!+

package tool

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/credentials"
	"testing"
)

func TestLoadServerTLSCredentials(t *testing.T) {
	type input struct {
		certFile string
		keyFile  string
	}
	var tests = []struct {
		name  string
		input input
		fRun  func(*testing.T, string, string)
	}{
		{
			name: "positive test #0 LoadServerTLSCredentials",
			input: input{
				certFile: "server-cert.pem",
				keyFile:  "server-key.pem",
			},
			fRun: positiveLoadServerTLSCredentials,
		},
		{
			name: "negative test #1 LoadServerTLSCredentials",
			input: input{
				certFile: "",
				keyFile:  "",
			},
			fRun: negativeLoadServerTLSCredentials,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t, test.input.certFile, test.input.keyFile)
		})
	}
}

func negativeLoadServerTLSCredentials(t *testing.T, certFile, keyFile string) {
	_, err := LoadServerTLSCredentials(certFile, keyFile)
	assert.NotNil(t, err)
}

func positiveLoadServerTLSCredentials(t *testing.T, certFile, keyFile string) {
	expectedInfo := credentials.ProtocolInfo(
		credentials.ProtocolInfo{
			ProtocolVersion:  "",
			SecurityProtocol: "tls",
			SecurityVersion:  "1.2",
			ServerName:       "",
		})

	got, err := LoadServerTLSCredentials(certFile, keyFile)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, expectedInfo, got.Info())
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
