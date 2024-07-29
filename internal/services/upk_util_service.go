/*
 * This file was last modified at 2024-07-29 19:13 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * upk_util_service.go
 * $Id$
 */

package services

import (
	"crypto/rsa"
	"encoding/base64"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"sync"
)

type UpkUtilService interface {
	EncryptAES(plain []byte) ([]byte, error)
	DecryptAES(bytes []byte) ([]byte, error)
	EncryptPersonalKey(personalKey string) (string, error)
	EncryptRSA(plain []byte) ([]byte, error)
	DecryptRSA(bytes []byte) ([]byte, error)
}

type upkUtilService struct {
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
	secretKey     []byte
}

var _ UpkUtilService = (*upkUtilService)(nil)
var (
	onceUpkUtil = new(sync.Once)
	upkUtilServ *upkUtilService
)

func GetUpkUtilService(prop env.Properties) UpkUtilService {

	onceUpkUtil.Do(func() {
		upkUtilServ = new(upkUtilService)
		upkUtilServ.rsaPrivateKey = prop.UpkRSAPrivateKey()
		upkUtilServ.rsaPublicKey = prop.UpkRSAPublicKey()
		upkUtilServ.secretKey = prop.UpkSecretKey()
	})
	return upkUtilServ
}

func (u *upkUtilService) EncryptAES(plain []byte) ([]byte, error) {
	return tool.EncryptAES(u.secretKey, plain)
}

func (u *upkUtilService) DecryptAES(bytes []byte) ([]byte, error) {
	return tool.DecryptAES(u.secretKey, bytes)
}

func (u *upkUtilService) EncryptPersonalKey(personalKey string) (string, error) {

	bytes := make([]byte, 32)
	copy(bytes, personalKey)
	encrypted, err := u.EncryptAES(bytes)

	if err != nil {
		return "", err
	}
	upk := base64.StdEncoding.EncodeToString(encrypted)

	return upk, nil
}

func (u *upkUtilService) EncryptRSA(plain []byte) ([]byte, error) {
	return tool.EncryptRSA(u.rsaPublicKey, plain)
}

func (u *upkUtilService) DecryptRSA(bytes []byte) ([]byte, error) {
	return tool.DecryptRSA(u.rsaPrivateKey, bytes)
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
