/*
 * This file was last modified at 2024-07-21 09:58 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * user_search_service.go
 * $Id$
 */
//!+

// Package services TODO.
package services

import (
	"context"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/repo"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
	pb "github.com/vskurikhin/gofavorites/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"log/slog"
	"sync"
	"time"
)

type UserSearchService interface {
	Lookup(ctx context.Context, personalKey, upk string) bool
}

type userSearchService struct {
	authGRPCAddress string
	creds           credentials.TransportCredentials
	opts            []grpc.DialOption
	repoUser        domain.Repo[*entity.User]
	requestInterval time.Duration
}

var _ UserSearchService = (*userSearchService)(nil)
var (
	onceUserSearch = new(sync.Once)
	userSearchSrv  *userSearchService
)

func GetUserSearchService(prop env.Properties) UserSearchService {

	onceUserSearch.Do(func() {
		userSearchSrv = new(userSearchService)
		opts := []grpc.DialOption{
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		}
		tlsCredentials, err := tool.LoadAgentTLSCredentials(prop.Config().GRPCTLSCAFile())
		if err != nil {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
		}
		userSearchSrv.authGRPCAddress = prop.ExternalAuthGRPCAddress()
		userSearchSrv.opts = opts
		userSearchSrv.repoUser = repo.GetUserPostgresCachedRepo(prop)
		userSearchSrv.requestInterval = prop.ExternalRequestTimeoutInterval()
	})
	return userSearchSrv
}

func (u *userSearchService) Lookup(ctx context.Context, personalKey, upk string) bool {

	if u.dbLookup(ctx, upk) {
		return true
	}
	if u.grpcLookupUpk(ctx, upk) {
		return true
	}
	if u.grpcLookupPersonalKey(ctx, personalKey) {
		return true
	}
	return false
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
	ctx, cancel := context.WithTimeout(ctx, u.requestInterval)
	defer func() {
		cancel()
		ctx.Done()
	}()
	var request pb.UserRequest
	request.User = &pb.User{PersonalKey: personalKey}
	resp, err := c.Get(ctx, &request)

	for i := 1; err != nil && tool.IsUpperBound(i, u.requestInterval); i++ {
		time.Sleep(100 * time.Millisecond * time.Duration(i))
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
	ctx, cancel := context.WithTimeout(ctx, u.requestInterval)
	defer func() {
		cancel()
		ctx.Done()
	}()
	var request pb.UserRequest
	request.User = &pb.User{Upk: upk}
	resp, err := c.Get(ctx, &request)

	for i := 1; err != nil && tool.IsUpperBound(i, u.requestInterval); i++ {
		slog.Warn(env.MSG+" UserSearchService.grpcLookup", "err", err)
		time.Sleep(100 * time.Millisecond * time.Duration(i))
		resp, err = c.Get(ctx, &request)
	}
	if err == nil && resp.Status == pb.Status_OK {
		return true
	}
	return false
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
