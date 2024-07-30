/*
 * This file was last modified at 2024-07-24 10:23 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * auth_test.go
 * $Id$
 */
//!+

// Package controllers TODO.
package controllers

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/vskurikhin/gofavorites/internal/env"
)

// go test -run TestDeserializeUser
func Test_Auth_SignInUser_Positive_1(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	httpPort := 65500 + rnd.Intn(34)
	t.Setenv("HTTP_ADDRESS", fmt.Sprintf("127.0.0.1:%d", httpPort))

	prop := env.GetProperties()
	t.Run("Auth SignInUser Positive #1", func(t *testing.T) {
		app := fiber.New()

		app.Post("/", GetAuthController(prop).SignInUser)

		req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(`{"user_name":"test","password":"password"}`))
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := app.Test(req)
		utils.AssertEqual(t, nil, err, "app.Test(req)")
		utils.AssertEqual(t, 200, resp.StatusCode, "Status code")

		_, err = io.ReadAll(resp.Body)
		utils.AssertEqual(t, nil, err)
	})
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
