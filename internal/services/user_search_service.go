/*
 * This file was last modified at 2024-08-06 20:32 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * user_search_service.go
 * $Id$
 */
//!+

// Package services сервисы бизнес логики.
package services

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/repo"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/models"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	pb "github.com/vskurikhin/gofavorites/proto"
)

type UserSearchService interface {
	Lookup(ctx context.Context, user models.User) bool
}

type userSearchService struct {
	authGRPCAddress string
	creds           credentials.TransportCredentials
	opts            []grpc.DialOption
	repoUser        domain.Repo[*entity.User]
	requestInterval time.Duration
	sLog            *slog.Logger
}

var _ UserSearchService = (*userSearchService)(nil)
var (
	onceUserSearch = new(sync.Once)
	userSearchServ *userSearchService
)

// GetUserSearchService — потокобезопасное (thread-safe) создание
// сервиса поиска пользователей в базе данных или в системе аутентификации.
func GetUserSearchService(prop env.Properties) UserSearchService {

	onceUserSearch.Do(func() {
		userSearchServ = new(userSearchService)
		opts := []grpc.DialOption{
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		}
		tlsCredentials, err := tool.LoadClientTLSCredentials(prop.Config().GRPCTLSCAFile())
		if err != nil {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
		}
		userSearchServ.authGRPCAddress = prop.ExternalAuthGRPCAddress()
		userSearchServ.opts = opts
		userSearchServ.repoUser = repo.GetUserPostgresCachedRepo(prop)
		userSearchServ.requestInterval = prop.ExternalRequestTimeoutInterval()
		userSearchServ.sLog = prop.Logger()
	})
	return userSearchServ
}

const cntUserSearchLookupJobs = 3

// Lookup поиск пользователя в базе данных или в системе аутентификации.
func (u *userSearchService) Lookup(ctx context.Context, user models.User) bool {

	var wg sync.WaitGroup
	wg.Add(cntUserSearchLookupJobs)

	quit := make(chan struct{})
	results := make(chan bool, cntUserSearchLookupJobs)

	go func() {
		defer wg.Done()
		results <- u.dbLookup(ctx, user.Upk())
	}()
	go func() {
		defer wg.Done()
		results <- u.grpcLookupUpk(ctx, user.Upk())
	}()
	go func() {
		defer wg.Done()
		results <- u.grpcLookupPersonalKey(ctx, user.PersonalKey())
	}()
	go func() {
		wg.Wait()
		close(results)
		close(quit)
	}()
	for {
		select {
		case result := <-results:
			if result {
				return result
			}
		case <-quit:
			return false
		}
	}
}

func (u *userSearchService) dbLookup(ctx context.Context, upk string) bool {

	asset, err := entity.GetUser(ctx, u.repoUser, upk)

	if entity.IsUserNotFound(asset, err) {
		return false
	}
	return true
}

func (u *userSearchService) grpcLookupPersonalKey(ctx context.Context, personalKey string) bool {

	conn, err := grpc.NewClient(u.authGRPCAddress, u.opts...)
	if err != nil {
		return false
	}
	defer func() { _ = conn.Close() }()
	c := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer func() {
		cancel()
	}()
	var request pb.UserRequest
	request.User = &pb.User{PersonalKey: personalKey}
	resp, err := c.Get(ctx, &request)

	for i := 1; err != nil && tool.IsUpperBoundWithSleep(i, 300, u.requestInterval); i++ {
		u.sLog.ErrorContext(ctx,
			env.MSG+"UserSearchService.grpcLookupPersonalKey",
			"msg", "user search service gRPC lookup personal key",
			"err", err,
		)
		time.Sleep(300 * time.Millisecond * time.Duration(i))
		resp, err = c.Get(ctx, &request)
	}
	if resp != nil && resp.Status == pb.Status_OK {
		return true
	}
	return false
}

func (u *userSearchService) grpcLookupUpk(ctx context.Context, upk string) bool {

	conn, err := grpc.NewClient(u.authGRPCAddress, u.opts...)
	if err != nil {
		return false
	}
	defer func() { _ = conn.Close() }()
	c := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer func() {
		cancel()
	}()
	var request pb.UserRequest
	request.User = &pb.User{Upk: upk}
	resp, err := c.Get(ctx, &request)

	for i := 1; err != nil && tool.IsUpperBoundWithSleep(i, 300, u.requestInterval); i++ {
		u.sLog.ErrorContext(ctx,
			env.MSG+"UserSearchService.grpcLookupUpk",
			"msg", "user search service gRPC lookup upk",
			"err", err,
		)
		time.Sleep(300 * time.Millisecond * time.Duration(i))
		resp, err = c.Get(ctx, &request)
	}
	if err == nil && resp.Status == pb.Status_OK {
		return true
	}
	return false
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
