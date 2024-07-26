/*
 * This file was last modified at 2024-07-27 10:54 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * deserialize-user.go
 * $Id$
 */

package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vskurikhin/gofavorites/internal/env"
	"strings"
	"sync"
)

type UserJwtHandler struct {
	jwtSecret string
}

var (
	onceUserJwt = new(sync.Once)
	userJwtMidl *UserJwtHandler
)

func GetUserJwtHandler(prop env.Properties) *UserJwtHandler {

	onceUserJwt.Do(func() {
		userJwtMidl = new(UserJwtHandler)
		userJwtMidl.jwtSecret = prop.JwtSecret()
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
	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(u.jwtSecret), nil
	})
	if err != nil {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{
				"status":  "fail",
				"message": fmt.Sprintf("invalidate token: %v", err),
			})
	}
	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

	}
	//var user models.User
	//initializers.DB.First(&user, "id = ?", fmt.Sprint(claims["sub"]))
	//if user.ID.String() != claims["sub"] {
	//	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	//}

	c.Locals("user", fmt.Sprint(claims["sub"]))

	return c.Next()
}
