/*
 * This file was last modified at 2024-07-26 11:26 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * favorites_test.go
 * $Id$
 */
//!+

// Package controllers REST-ful (endpoints) конечные точки REST веб-сервиса.
package controllers

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/middleware"
	"github.com/vskurikhin/gofavorites/internal/models"
	"github.com/vskurikhin/gofavorites/internal/services"
	"go.uber.org/mock/gomock"
)

func TestFavoritesGet(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #0 Favorites.Get",
			fRun: positiveFavoritesGet,
		},
		{
			name: "positive test #1 Favorites.GetForUser",
			fRun: positiveFavoritesGetForUser,
		},
		{
			name: "positive test #2 Favorites.Set",
			fRun: positiveFavoritesSet,
		},
		{
			name: "negative test #3 Favorites.Get #0",
			fRun: negativeFavoritesGet0,
		},
		{
			name: "negative test #4 Favorites.Get #1",
			fRun: negativeFavoritesGet1,
		},
		{
			name: "negative test #5 Favorites.Get #2",
			fRun: negativeFavoritesGet2,
		},
		{
			name: "negative test #6 Favorites.Get #3",
			fRun: negativeFavoritesGet3,
		},
		{
			name: "negative test #7 Favorites.GetForUser #0",
			fRun: negativeFavoritesGetForUser0,
		},
		{
			name: "negative test #8 Favorites.GetForUser #1",
			fRun: negativeFavoritesGetForUser1,
		},
		{
			name: "negative test #9 Favorites.Set #0",
			fRun: negativeFavoritesSet0,
		},
		{
			name: "negative test #10 Favorites.Set #1",
			fRun: negativeFavoritesSet1,
		},
		{
			name: "negative test #11 Favorites.Set #2",
			fRun: negativeFavoritesSet2,
		},
		{
			name: "negative test #12 Favorites.Set #3",
			fRun: negativeFavoritesSet3,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func positiveFavoritesGet(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesGet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, nil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Get,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(`{"isin":"test","asset_type":"test"}`))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}

func positiveFavoritesGetForUser(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	var stubResult = make([]models.Favorites, 0)
	favoritesServ.
		EXPECT().
		ApiFavoritesGetForUser(gomock.Any(), gomock.Any()).
		Return(stubResult, nil).
		AnyTimes()
	app.Get("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).GetForUser,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}

func positiveFavoritesSet(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesSet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, nil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Set,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(`{"isin":"test","asset_type":"test"}`))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}

func negativeFavoritesGet0(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesGet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, nil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Get,
	)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(``))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 401, resp.StatusCode, "Status code")
}

func negativeFavoritesGet1(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesGet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, nil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Get,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(``))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}

func negativeFavoritesGet2(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesGet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, nil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Get,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}

func negativeFavoritesGet3(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesGet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, services.ErrRequestNil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Get,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(`{"isin":"test","asset_type":"test"}`))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}

func negativeFavoritesGetForUser0(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	var stubResult = make([]models.Favorites, 0)
	favoritesServ.
		EXPECT().
		ApiFavoritesGetForUser(gomock.Any(), gomock.Any()).
		Return(stubResult, nil).
		AnyTimes()
	app.Get("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).GetForUser,
	)

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 401, resp.StatusCode, "Status code")
}

func negativeFavoritesGetForUser1(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	var stubResult = make([]models.Favorites, 0)
	favoritesServ.
		EXPECT().
		ApiFavoritesGetForUser(gomock.Any(), gomock.Any()).
		Return(stubResult, services.ErrRequestNil).
		AnyTimes()
	app.Get("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).GetForUser,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}

func negativeFavoritesSet0(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesSet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, nil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Set,
	)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(``))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 401, resp.StatusCode, "Status code")
}

func negativeFavoritesSet1(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesSet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, nil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Set,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(``))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}

func negativeFavoritesSet2(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesSet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, nil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Set,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}

func negativeFavoritesSet3(t *testing.T) {

	ctrl := gomock.NewController(t)
	prop := env.GetProperties()
	app := fiber.New()
	app.Use(requestid.New())

	var favoritesServ = NewMockApiFavoritesService(ctrl)
	favoritesServ.
		EXPECT().
		ApiFavoritesSet(gomock.Any(), gomock.Any()).
		Return(models.Favorites{}, services.ErrRequestNil).
		AnyTimes()
	app.Post("/",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		getTestFavoritesController(prop, favoritesServ).Set,
	)
	tokenString, err := getTokenString(prop)
	assert.Nil(t, err)

	req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBufferString(`{"isin":"test","asset_type":"test"}`))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := app.Test(req)
	utils.AssertEqual(t, nil, err, "app.Test(req)")
	utils.AssertEqual(t, 400, resp.StatusCode, "Status code")
}

func getTokenString(prop env.Properties) (string, error) {
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = "test"
	claims["exp"] = now.Add(prop.JwtExpiresIn()).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(prop.JwtSecret()))
	return tokenString, err
}

func getTestFavoritesController(prop env.Properties, favoritesServ services.ApiFavoritesService) *Favorites {

	favoritesCont = new(Favorites)
	favoritesCont.favoritesServ = favoritesServ
	favoritesCont.jwtExpiresIn = prop.JwtExpiresIn()
	favoritesCont.jwtMaxAge = prop.JwtMaxAgeSec()
	favoritesCont.jwtSecret = prop.JwtSecret()
	favoritesCont.token = prop.Config().Token()

	return favoritesCont
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
