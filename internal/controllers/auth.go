/*
 * This file was last modified at 2024-08-06 20:17 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * auth.go
 * $Id$
 */
//!+

// Package controllers REST-ful (endpoints) конечные точки REST веб-сервиса.
package controllers

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/vskurikhin/gofavorites/internal/controllers/dto"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	token      string
	jwtManager jwt.Manager
	jwtMaxAge  int
}

var (
	onceAuth = new(sync.Once)
	authCont *Auth
)

// GetAuthController — потокобезопасное (thread-safe) создание
// REST веб-сервиса аутентификации.
func GetAuthController(prop env.Properties) *Auth {

	onceAuth.Do(func() {
		authCont = new(Auth)
		authCont.token = prop.Config().Token()
		authCont.jwtManager = jwt.GetJWTManager(prop)
		authCont.jwtMaxAge = prop.JwtMaxAgeSec()
	})
	return authCont
}

// SignInUser handler
//
//	@Summary		аутентификация
//	@Description	аутентификация пользователя
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		none
//	@Param			request			body		dto.SignInRequest	true	"Формат запроса JSON (body)"
//	@Success		200				{object}	dto.SignInRequest	"пользователь успешно аутентифицирован"
//	@Failure		400				{object}	dto.SignInRequest	"неверный формат запроса"
//	@Failure		401				{object}	dto.SignInRequest	"неверная пара логин/пароль"
//	@Failure		500				{string}	string			"Internal Server Error"
//	@Router			/api/auth/login	[post]
func (a *Auth) SignInUser(c *fiber.Ctx) error {

	var payload dto.SignInRequest
	requestId := c.Locals("requestid")

	if err := c.BodyParser(&payload); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "fail", "message": err.Error(), "requestId": requestId})
	}
	errors := dto.ValidateStruct(payload)

	if errors != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(errors)
	}
	err := bcrypt.CompareHashAndPassword([]byte(a.token), []byte(payload.Password))

	if err != nil {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{
				"status":  "fail",
				"message": "Invalid email or Password", "requestId": requestId,
			})
	}

	tokenString, err := a.jwtManager.Generate(payload)

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

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
