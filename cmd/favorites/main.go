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
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vskurikhin/gofavorites/internal/controllers"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/services"

	pb "github.com/vskurikhin/gofavorites/proto"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS //

func main() {
	run(context.Background())
}

func run(ctx context.Context) {
	slog.Info(env.MSG,
		"build_version", buildVersion,
		"build_date", buildDate,
		"build_commit", buildCommit,
	)
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	// запускаем горутину обработки пойманных прерываний
	prop := env.GetProperties()
	_, _ = fmt.Fprintf(os.Stderr, "PROPERTIES%s\n", prop)
	dbMigrations(prop)
	serve(ctx, prop)
}

func dbMigrations(prop env.Properties) {

	pool := prop.DBPool()
	if pool == nil {
		slog.Warn(env.MSG+" dbMigrations", "pool", pool)
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
		log.Println("Выключение сервера gRPC")
		if err := httpServer.Shutdown(); err != nil {
			log.Printf("Ошибка при выключение сервера HTTP: %v\n", err)
		}
		log.Println("Выключение сервера HTTP")
		close(idleConnsClosed)
	}()
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		log.Println("Выключение сервера gRPC")
		if err := httpServer.Shutdown(); err != nil {
			log.Printf("Ошибка при выключение сервера HTTP: %v\n", err)
		}
		log.Println("Выключение сервера HTTP")
		close(idleConnsClosed)
	}()
	go func() {
		log.Println("Сервер gRPC начал работу")
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Сервер HTTP начал работу")
	if err := httpServer.Listen(":8000"); err != nil {
		log.Printf("Ошибка при выключение сервера HTTP: %v\n", err)
	}
	<-idleConnsClosed
	log.Println("Корректное завершение работы сервера")
}

func makeHTTP(prop env.Properties) *fiber.App {

	logHandler := logger.New(logger.Config{
		Format:       "${pid} | ${time} | ${status} | ${locals:requestid} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat:   "15:04:05.999",
		TimeZone:     "Local",
		TimeInterval: 500 * time.Millisecond,
		Output:       os.Stdout,
	})

	app := fiber.New()
	micro := fiber.New()
	app.Mount("/api", micro)
	app.Use(logHandler)
	app.Use(requestid.New())
	micro.Use(logHandler)
	micro.Use(requestid.New())

	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/login", controllers.GetAuthController(prop).SignInUser)
	})

	// micro.Get("/users/me", middleware.DeserializeUser, controllers.GetMe)

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

	opts := []grpc.ServerOption{grpc.Creds(prop.GRPCTransportCredentials())}
	grpcServer := grpc.NewServer(opts...)
	favoritesService := services.GetFavoritesService(prop)
	pb.RegisterFavoritesServiceServer(grpcServer, favoritesService)

	return grpcServer
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
