/*
 * Copyright text:
 * This file was last modified at 2024-07-10 21:53 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * main.go
 * $Id$
 */
//!+

package main

import (
	"context"
	"crypto/tls"
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/vskurikhin/gofavorites/internal/alog"
	"github.com/vskurikhin/gofavorites/internal/controllers"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/interceptors"
	"github.com/vskurikhin/gofavorites/internal/middleware"
	"github.com/vskurikhin/gofavorites/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/vskurikhin/gofavorites/docs"
	pb "github.com/vskurikhin/gofavorites/proto"
	_ "google.golang.org/grpc/encoding/gzip"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	sLog         *slog.Logger
)

//go:embed migrations/*.sql
var embedMigrations embed.FS //

// @title GoFavorites API
// @Security     Bearer
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8443
// @BasePath /
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization

func main() {
	run(context.Background())
}

func run(ctx context.Context) {
	slog.Info(env.MSG+"meta info",
		"build_version", buildVersion,
		"build_date", buildDate,
		"build_commit", buildCommit,
	)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	prop := env.GetProperties()
	sLog = alog.GetLogger()
	dbMigrations(prop)
	serve(ctx, prop)
}

func dbMigrations(prop env.Properties) {

	pool := prop.DBPool()
	if pool == nil {
		sLog.Warn(env.MSG+"dbMigrations", "pool", pool)
		return
	}
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	db := stdlib.OpenDBFromPool(pool)

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
	if err := goose.Version(db, "migrations"); err != nil {
		log.Fatal(err)
	}
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func serve(ctx context.Context, prop env.Properties) {

	listen, err := net.Listen("tcp", prop.GRPCAddress())

	if err != nil {
		log.Fatal(err)
	}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	httpServer := makeHTTP(prop)
	grpcServer := makeGRPC(prop)

	go func() {
		<-sigint
		grpcServer.GracefulStop()
		sLog.Info(env.MSG+"graceful stop", "msg", "Выключение сервера gRPC")
		if err := httpServer.Shutdown(); err != nil {
			sLog.Error(env.MSG+"graceful stop", "msg", "Ошибка при выключение сервера HTTP", "err", err)
		}
		sLog.Info(env.MSG+"graceful stop", "msg", "Выключение сервера HTTP")
		close(idleConnsClosed)
	}()
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		sLog.Info(env.MSG+"graceful stop", "msg", "Выключение сервера gRPC")
		if err := httpServer.Shutdown(); err != nil {
			sLog.Error(env.MSG+"graceful stop", "msg", "Ошибка при выключение сервера HTTP", "err", err)
		}
		sLog.Info(env.MSG+"graceful stop", "msg", "Выключение сервера HTTP")
		close(idleConnsClosed)
	}()
	go func() {
		sLog.Info(env.MSG+"start app", "msg", "Сервер gRPC начал работу")
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()
	sLog.Info(env.MSG+"start app", "msg", "Сервер HTTP начал работу")
	if prop.Config().HTTPTLSEnabled() {

		ln, err := tls.Listen("tcp", prop.HTTPAddress(), prop.HTTPTLSConfig())
		if err != nil {
			panic(err)
		}
		if err := httpServer.Listener(ln); err != nil {
			sLog.Error(env.MSG+"start app", "msg", "Ошибка при выключение сервера HTTP", "err", err)
		}
	} else if err := httpServer.Listen(prop.HTTPAddress()); err != nil {
		sLog.Error(env.MSG+"start app", "msg", "Ошибка при выключение сервера HTTP", "err", err)
	}
	<-idleConnsClosed
	sLog.Info(env.MSG+"shutdown app", "msg", "Корректное завершение работы сервера")
}

func makeHTTP(prop env.Properties) *fiber.App {

	logHandler := logger.New(logger.Config{
		Format:       "${pid} | ${time} | ${status} | ${locals:requestid} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat:   "15:04:05.000000",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	})
	slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := fiber.New()
	micro := fiber.New()
	app.Mount("/api", micro)
	app.Use(requestid.New())

	micro.Use(requestid.New())

	if prop.SlogJSON() {
		app.Use(alog.New(slogLogger))
		micro.Use(alog.New(slogLogger))
	} else {
		app.Use(logHandler)
		micro.Use(logHandler)
	}
	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/login", controllers.GetAuthController(prop).SignInUser)
	})
	app.Get("/swagger/*", swagger.New(swagger.Config{PreauthorizeApiKey: "Bearer"}))
	micro.Get(
		"/favorites/get",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		controllers.GetFavoritesController(prop).GetForUser,
	)
	micro.Post(
		"/favorites/get",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		controllers.GetFavoritesController(prop).Get,
	)
	micro.Post(
		"/favorites/set",
		middleware.GetUserJwtHandler(prop).DeserializeUser,
		controllers.GetFavoritesController(prop).Set,
	)
	micro.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "JWT Authentication with Golang and Fiber",
		})
	})
	micro.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	})
	return app
}

func makeGRPC(prop env.Properties) *grpc.Server {

	var opts []grpc.ServerOption
	authInterceptor := interceptors.GetAuthInterceptor(prop)

	if prop.Config().GRPCTLSEnabled() {
		opts = []grpc.ServerOption{
			grpc.Creds(prop.GRPCTransportCredentials()),
			grpc.StreamInterceptor(authInterceptor.Stream()),
			grpc.UnaryInterceptor(authInterceptor.Unary()),
		}
	} else {
		opts = []grpc.ServerOption{
			grpc.Creds(insecure.NewCredentials()),
			grpc.StreamInterceptor(authInterceptor.Stream()),
			grpc.UnaryInterceptor(authInterceptor.Unary()),
		}
	}
	grpcServer := grpc.NewServer(opts...)
	favoritesService := services.GetFavoritesService(prop)
	pb.RegisterFavoritesServiceServer(grpcServer, favoritesService)
	reflection.Register(grpcServer)

	return grpcServer
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
