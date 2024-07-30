/*
 * This file was last modified at 2024-07-29 13:57 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * upk_util_service_test.go
 * $Id$
 */

package services

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
)

func TestUpkUtilService(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #0 User Service RSA",
			fRun: testUpkUtilServiceRSAPositiveCase,
		},
		{
			name: "positive test #1 User Service AES",
			fRun: testUpkUtilServiceAESPositiveCase,
		},
		{
			name: "positive test #2 User Service AES case #1",
			fRun: testUpkUtilServiceEncryptPersonalKeyPositiveCase,
		},
		{
			name: "positive test #3 User Service AES case #2",
			fRun: testUpkUtilServiceEncryptPersonalKeyPositiveCase2,
		},
	}
	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testUpkUtilServiceRSAPositiveCase(t *testing.T) {
	expected := "рентгеноэлектрокардиографический"
	srv := getTestUpkUtilService(
		tool.LoadPrivateKey("test_private-key.pem"),
		tool.LoadPublicKey("test_public-key.pem"),
		make([]byte, 32),
	)
	encrypt, err := srv.EncryptRSA([]byte(expected))
	assert.Nil(t, err)
	got, err := srv.DecryptRSA(encrypt)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(got))
}

func testUpkUtilServiceAESPositiveCase(t *testing.T) {
	secretKey := make([]byte, 32)
	if _, err := rand.Reader.Read(secretKey); err != nil {
		t.Fail()
	}
	expected := "рентгено"
	srv := getTestUpkUtilService(nil, nil, secretKey)
	encrypt, err := srv.EncryptAES([]byte(expected))
	assert.Nil(t, err)
	got, err := srv.DecryptAES(encrypt)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(got))
}

func testUpkUtilServiceEncryptPersonalKeyPositiveCase(t *testing.T) {
	test := "test"
	expected := "CdTLuDHCHE1DrSaSm2WZtgAAAAAAAAAAAAAAAAAAAAA="
	secretKey := []byte{48, 17, 60, 87, 186, 101, 173, 89, 205, 24, 23, 245, 219, 42, 222, 100}
	srv := getTestUpkUtilService(nil, nil, secretKey)
	encrypt, err := srv.EncryptPersonalKey(test)
	assert.Nil(t, err)
	assert.Equal(t, expected, encrypt)
}

func testUpkUtilServiceEncryptPersonalKeyPositiveCase2(t *testing.T) {

	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("UPK_PRIVATE_KEY_FILE", "test_private-key.pem")
	t.Setenv("UPK_PUBLIC_KEY_FILE", "test_public-key.pem")
	t.Setenv("UPK_SECRET", "qYhaPtg+PIQtBhAU5fHCeQw7XIF3WLKoLPZnJgq1H//DDOB8o2qrP9goVCUZldOdwqLAHxWOGHuvXcwaIFRrD8I3Hz5tRCgCeI+cEZD9h4c4h6ADSjkcrPXg5eRwnANasBkKKZQz8noYwvt9Z9p7HdOtrBmQOi7OVjTfY0T2SnI=")

	expected := "ZRiw9fJPhGsLDByaA9eQDAAAAAAAAAAAAAAAAAAAAAA="
	prop := env.GetProperties()
	srv := GetUpkUtilService(prop)
	test := "test"
	encrypt, err := srv.EncryptPersonalKey(test)
	assert.Nil(t, err)
	assert.Equal(t, expected, encrypt)
}

func getTestUpkUtilService(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, secretKey []byte) UpkUtilService {
	u := new(upkUtilService)
	u.rsaPrivateKey = privateKey
	u.rsaPublicKey = publicKey
	u.secretKey = secretKey
	return u
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
