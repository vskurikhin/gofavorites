/*
 * This file was last modified at 2024-08-04 22:01 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * deserialize-user.go
 * $Id$
 */

package middleware

import (
	"fmt"
	"strings"
	"sync"

	"github.com/vskurikhin/gofavorites/internal/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/vskurikhin/gofavorites/internal/env"
)

type UserJwtHandler struct {
	jwtManager jwt.Manager
}

var (
	onceUserJwt = new(sync.Once)
	userJwtMidl *UserJwtHandler
)

func GetUserJwtHandler(prop env.Properties) *UserJwtHandler {

	onceUserJwt.Do(func() {
		userJwtMidl = new(UserJwtHandler)
		userJwtMidl.jwtManager = jwt.GetJWTManager(prop)
	})
	return userJwtMidl
}

func (u *UserJwtHandler) DeserializeUser(c *fiber.Ctx) error {

	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}
	if tokenString == "" {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}
	claims, err := u.jwtManager.Verify(tokenString)

	if err != nil || claims == nil {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{
				"status":  "fail",
				"message": fmt.Sprintf("invalidate token: %v", err),
			})
	}
	c.Locals("user", fmt.Sprint(claims.UserName()))

	return c.Next()
}
