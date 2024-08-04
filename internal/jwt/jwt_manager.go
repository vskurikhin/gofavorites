/*
 * This file was last modified at 2024-08-04 14:16 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * jwt_manager.go
 * $Id$
 */
//!+

// Package jwt TODO.
package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vskurikhin/gofavorites/internal/controllers/dto"
	"github.com/vskurikhin/gofavorites/internal/env"
	"sync"
	"time"
)

type Manager interface {
	Generate(user dto.SignInRequest) (string, error)
	Verify(accessToken string) (*Claims, error)
}

type Claims struct {
	jwt.RegisteredClaims
	userName string
	role     string
}

func (c Claims) UserName() string {
	return c.userName
}

func (c Claims) Role() string {
	return c.role
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

type manager struct {
	jwtSecret    string
	jwtExpiresIn time.Duration
}

var _ Manager = (*manager)(nil)
var (
	ErrInvalidTokenClaims           = fmt.Errorf("invalid token claims")
	ErrUnexpectedTokenSigningMethod = fmt.Errorf("unexpected token signing method")
	onceManager                     = new(sync.Once)
	jwtManager                      *manager
)

func GetJWTManager(prop env.Properties) Manager {

	onceManager.Do(func() {
		jwtManager = new(manager)
		jwtManager.jwtSecret = prop.JwtSecret()
		jwtManager.jwtExpiresIn = prop.JwtExpiresIn()
	})
	return jwtManager
}

func (m *manager) Generate(user dto.SignInRequest) (string, error) {

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.jwtExpiresIn)),
		},
		Username: user.UserName,
		Role:     "USER",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.jwtSecret))
}

// Verify verifies the access token string and return a user claim if the token is valid
func (m *manager) Verify(accessToken string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrUnexpectedTokenSigningMethod
			}

			return []byte(m.jwtSecret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	userClaims, ok := token.Claims.(*UserClaims)

	if !ok {
		return nil, ErrInvalidTokenClaims
	}
	claims := &Claims{
		userName: userClaims.Username,
		role:     userClaims.Role,
	}
	return claims, nil
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
