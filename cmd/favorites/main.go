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
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/services"
	"google.golang.org/grpc"

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
	_ = fiber.New()
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

	// определяем порт для gRPC сервера
	listen, err := net.Listen("tcp", prop.GRPCAddress())
	fmt.Println(prop.GRPCAddress())
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.ServerOption{}
	tlsCredentials := prop.GRPCTransportCredentials()

	if err != nil {
		log.Println("Не удалось загрузить сертификаты для сервера gRPC")
	} else {
		opts = append(opts, grpc.Creds(tlsCredentials))
	}
	grpcServer := grpc.NewServer(opts...)
	favoritesService := services.GetFavoritesService(prop)
	pb.RegisterFavoritesServiceServer(grpcServer, favoritesService)

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-sigint
		grpcServer.Stop()
		log.Println("Выключение сервера gRPC")
		close(idleConnsClosed)
	}()
	go func() {
		<-ctx.Done()
		grpcServer.Stop()
		log.Println("Выключение сервера gRPC")
		close(idleConnsClosed)
	}()

	log.Println("Сервер gRPC начал работу")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}
	<-idleConnsClosed
	log.Println("Корректное завершение работы сервера")
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
