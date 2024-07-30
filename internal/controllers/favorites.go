/*
 * This file was last modified at 2024-07-31 14:52 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites.go
 * $Id$
 */
//!+

// Package controllers TODO.
package controllers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vskurikhin/gofavorites/internal/controllers/dto"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/models"
	"github.com/vskurikhin/gofavorites/internal/services"
)

type Favorites struct {
	favoritesServ services.ApiFavoritesService
	jwtExpiresIn  time.Duration
	jwtMaxAge     int
	jwtSecret     string
	token         string
}

var (
	onceFavorites = new(sync.Once)
	favoritesCont *Favorites
)

func GetFavoritesController(prop env.Properties) *Favorites {

	onceFavorites.Do(func() {
		favoritesCont = new(Favorites)
		favoritesCont.favoritesServ = services.GetFavoritesService(prop)
		favoritesCont.jwtExpiresIn = prop.JwtExpiresIn()
		favoritesCont.jwtMaxAge = prop.JwtMaxAgeSec()
		favoritesCont.jwtSecret = prop.JwtSecret()
		favoritesCont.token = prop.Config().Token()
	})
	return favoritesCont
}

func (f *Favorites) Get(c *fiber.Ctx) error {

	var payload dto.Favorites

	requestId := c.Locals("requestid")
	user, ok := c.Locals("user").(string)

	if !ok {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "fail", "requestId": requestId, "message": "user failed"})
	}
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
	model := models.FavoritesFromDto(payload, user, "")
	ctx := context.WithValue(c.Context(), "request-id", requestId)
	favorites, err := f.favoritesServ.ApiFavoritesGet(ctx, model)

	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  "fail",
				"message": fmt.Sprintf("error: %v", err),
			})
	}
	response := favorites.ToDto()

	return c.
		Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status":    "success",
			"requestId": requestId,
			"data":      fiber.Map{"favorites": response, "user": user},
		})
}

func (f *Favorites) GetForUser(c *fiber.Ctx) error {

	requestId := c.Locals("requestid")
	user, ok := c.Locals("user").(string)

	if !ok {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "fail", "requestId": requestId, "message": "user failed"})
	}
	model := models.MakeUser(user, "")
	ctx := context.WithValue(c.Context(), "request-id", requestId)
	favorites, err := f.favoritesServ.ApiFavoritesGetForUser(ctx, model)

	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  "fail",
				"message": fmt.Sprintf("error: %v", err),
			})
	}
	response := models.FavoritesSliceToDto(favorites)

	return c.
		Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status":    "success",
			"requestId": requestId,
			"data":      fiber.Map{"favorites": response, "user": user},
		})
}

func (f *Favorites) Set(c *fiber.Ctx) error {

	var payload dto.Favorites

	requestId := c.Locals("requestid")
	user, ok := c.Locals("user").(string)

	if !ok {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "fail", "requestId": requestId, "message": "user failed"})
	}
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
	model := models.FavoritesFromDto(payload, user, "")
	ctx := context.WithValue(c.Context(), "request-id", requestId)
	favorites, err := f.favoritesServ.ApiFavoritesSet(ctx, model)
	response := favorites.ToDto()

	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  "fail",
				"message": fmt.Sprintf("error: %v", err),
			})
	}
	return c.
		Status(fiber.StatusOK).
		JSON(fiber.Map{
			"status":    "success",
			"requestId": requestId,
			"data":      fiber.Map{"favorites": response, "user": user},
		})
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
