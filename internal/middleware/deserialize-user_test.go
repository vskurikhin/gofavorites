/*
 * This file was last modified at 2024-07-23 15:18 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * deserialize-user_test.go
 * $Id$
 */
//!+

// Package tool TODO.
package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vskurikhin/gofavorites/internal/env"
	"io"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"
)

// go test -run TestDeserializeUser
func TestDeserializeUser(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	httpPort := 65500 + rnd.Intn(34)
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("HTTP_ADDRESS", fmt.Sprintf("127.0.0.1:%d", httpPort))

	prop := env.GetProperties()
	testHandler := GetUserJwtHandler(prop).DeserializeUser
	t.Run("test", func(t *testing.T) {
		app := fiber.New()

		app.Get("/", testHandler, func(c *fiber.Ctx) error {
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Send([]byte("ok"))
		})

		tokenByte := jwt.New(jwt.SigningMethodHS256)
		now := time.Now().UTC()
		claims := tokenByte.Claims.(jwt.MapClaims)

		claims["sub"] = "test"
		claims["exp"] = now.Add(prop.JwtExpiresIn()).Unix()
		claims["iat"] = now.Unix()
		claims["nbf"] = now.Unix()

		tokenString, err := tokenByte.SignedString([]byte(prop.JwtSecret()))

		req := httptest.NewRequest(fiber.MethodGet, "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 200, resp.StatusCode, "Status code")

		// Validate body
		body, err := io.ReadAll(resp.Body)
		utils.AssertEqual(t, nil, err)
		utils.AssertEqual(t, []byte("ok"), body)
	})
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
