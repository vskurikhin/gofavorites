/*
 * This file was last modified at 2024-07-23 14:43 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * auth.go
 * $Id$
 */

package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/models"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

type Auth struct {
	token        string
	jwtExpiresIn time.Duration
	jwtMaxAge    int
	jwtSecret    string
}

var (
	onceAuth = new(sync.Once)
	authCont *Auth
)

func GetAuthController(prop env.Properties) *Auth {

	onceAuth.Do(func() {
		authCont = new(Auth)
		authCont.token = prop.Config().Token()
		authCont.jwtExpiresIn = prop.JwtExpiresIn()
		authCont.jwtMaxAge = prop.JwtMaxAgeSec()
		authCont.jwtSecret = prop.JwtSecret()
	})
	return authCont
}

func (a *Auth) SignInUser(c *fiber.Ctx) error {

	var payload models.SignInRequest
	requestId := c.Locals("requestid")

	if err := c.BodyParser(&payload); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "fail", "message": err.Error(), "requestId": requestId})
	}
	errors := models.ValidateStruct(payload)

	if errors != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(errors)
	}
	err := bcrypt.CompareHashAndPassword([]byte(a.token), []byte(payload.Password))

	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid email or Password", "requestId": requestId,
			})
	}
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = payload.UserName
	claims["exp"] = now.Add(a.jwtExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(a.jwtSecret))

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   a.jwtMaxAge,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	if err != nil {
		return c.
			Status(fiber.StatusBadGateway).
			JSON(fiber.Map{
				"status":  "fail",
				"message": fmt.Sprintf("generating JWT Token failed: %v", err),
			})
	}
	return c.
		Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "requestId": requestId, "token": tokenString})
}
