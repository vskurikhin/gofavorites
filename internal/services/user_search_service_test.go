/*
 * This file was last modified at 2024-07-20 19:34 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * user_search_service_test.go
 * $Id$
 */
//!+

// Package services TODO.
package services

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/domain"
	"github.com/vskurikhin/gofavorites/internal/domain/entity"
	"github.com/vskurikhin/gofavorites/internal/domain/repo"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/tool"
	pb "github.com/vskurikhin/gofavorites/proto"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/local"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestUserSearchService(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #1 User Service Lookup case #1",
			fRun: testUserSearchServiceLookupPositiveCase1,
		},
		{
			name: "positive test #2 User Service Lookup case #2",
			fRun: testUserSearchServiceLookupPositiveCase2,
		},
		{
			name: "positive test #3 User Service Lookup case #3",
			fRun: testUserSearchServiceLookupPositiveCase3,
		},
		{
			name: "negative test #4 User Service Lookup case #1",
			fRun: testUserSearchServiceLookupNegativeCase1,
		},
		{
			name: "negative test #5 User Service Lookup case #2",
			fRun: testUserSearchServiceLookupNegativeCase2,
		},
		{
			name: "negative test #6 GetUserSearchService",
			fRun: testGetUserSearchService,
		},
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("ASSET_GRPC_ADDRESS", fmt.Sprintf("127.0.0.1:%d", 65501+rnd.Intn(34)))
	t.Setenv("AUTH_GRPC_ADDRESS", fmt.Sprintf("127.0.0.1:%d", 65501+rnd.Intn(34)))
	t.Setenv("REQUEST_TIMEOUT_INTERVAL_MS", "500")

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testGetUserSearchService(t *testing.T) {
	prop := env.GetProperties()
	got := GetUserSearchService(prop)
	assert.NotNil(t, got)
}

func testUserSearchServiceLookupPositiveCase1(t *testing.T) {
	prop := env.GetProperties()
	ctrl := gomock.NewController(t)
	repoUser := NewMockRepo[*entity.User](ctrl)
	user1 := entity.MakeUser("test", entity.DefaultTAttributes())
	repoUser.
		EXPECT().
		Get(context.Background(), gomock.Any(), gomock.Any()).
		Return(&user1, nil).
		AnyTimes()
	userSearchService := getUserSearchService(prop, repoUser)
	got := userSearchService.Lookup(context.Background(), "", user1.Upk())
	assert.True(t, got)
}

func testUserSearchServiceLookupPositiveCase2(t *testing.T) {
	prop := env.GetProperties()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer func() {
		cancel()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeUserServiceServer(ctx, prop, userServicePositiveCase1{}, up)
	<-up

	ctrl := gomock.NewController(t)
	repoUser := NewMockRepo[*entity.User](ctrl)
	user1 := entity.MakeUser("test", entity.DefaultTAttributes())
	repoUser.
		EXPECT().
		Get(context.Background(), gomock.Any(), gomock.Any()).
		Return(nil, repo.ErrNotFound).
		AnyTimes()
	defer func() {
		cancel()
		time.Sleep(100 * time.Millisecond)
	}()

	userSearchService := getUserSearchService(prop, repoUser)
	got := userSearchService.Lookup(context.Background(), "", user1.Upk())
	assert.True(t, got)
}

func testUserSearchServiceLookupPositiveCase3(t *testing.T) {
	prop := env.GetProperties()
	ctx, cancel := context.WithTimeout(context.Background(), 900*time.Millisecond)
	defer func() {
		cancel()
	}()
	up := make(chan struct{})
	go grpcServeUserServiceServer(ctx, prop, userServicePositiveCase2{}, up)
	<-up

	ctrl := gomock.NewController(t)
	repoUser := NewMockRepo[*entity.User](ctrl)
	repoUser.
		EXPECT().
		Get(context.Background(), gomock.Any(), gomock.Any()).
		Return(nil, repo.ErrNotFound).
		AnyTimes()
	defer func() {
		cancel()
		time.Sleep(100 * time.Millisecond)
	}()

	userSearchService := getUserSearchService(prop, repoUser)
	got := userSearchService.Lookup(context.Background(), "test", "")
	assert.True(t, got)
}

func testUserSearchServiceLookupNegativeCase1(t *testing.T) {
	prop := env.GetProperties()
	ctrl := gomock.NewController(t)
	repoUser := NewMockRepo[*entity.User](ctrl)
	repoUser.
		EXPECT().
		Get(context.Background(), gomock.Any(), gomock.Any()).
		Return(nil, repo.ErrNotFound).
		AnyTimes()
	userSearchService := getUserSearchService(prop, repoUser)
	got := userSearchService.Lookup(context.Background(), "test", "")
	assert.False(t, got)
}

func testUserSearchServiceLookupNegativeCase2(t *testing.T) {
	prop := env.GetProperties()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer func() {
		cancel()
	}()
	up := make(chan struct{})
	go grpcServeUserServiceServer(ctx, prop, userServiceNegativeCase1{}, up)
	<-up

	ctrl := gomock.NewController(t)
	repoUser := NewMockRepo[*entity.User](ctrl)
	repoUser.
		EXPECT().
		Get(context.Background(), gomock.Any(), gomock.Any()).
		Return(nil, repo.ErrNotFound).
		AnyTimes()
	defer func() {
		cancel()
		time.Sleep(100 * time.Millisecond)
	}()

	userSearchService := getUserSearchService(prop, repoUser)
	got := userSearchService.Lookup(context.Background(), "test", "")
	assert.False(t, got)
}

func getUserSearchService(prop env.Properties, repoUser domain.Repo[*entity.User]) UserSearchService {
	userSearchServ = new(userSearchService)
	opts := []grpc.DialOption{
		grpc.WithNoProxy(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	userSearchServ.authGRPCAddress = prop.ExternalAuthGRPCAddress()
	userSearchServ.opts = opts
	userSearchServ.repoUser = repoUser
	userSearchServ.requestInterval = prop.ExternalRequestTimeoutInterval()
	return userSearchServ
}

type userServicePositiveCase1 struct {
	pb.UnimplementedUserServiceServer
}

func (a userServicePositiveCase1) Get(_ context.Context, request *pb.UserRequest) (*pb.UserResponse, error) {
	if request.GetUser().GetUpk() != "" {
		return &pb.UserResponse{User: &pb.User{Upk: request.GetUser().GetUpk()}, Status: pb.Status_OK}, nil
	}
	return &pb.UserResponse{User: &pb.User{}, Status: pb.Status_FAIL}, nil
}

type userServicePositiveCase2 struct {
	pb.UnimplementedUserServiceServer
}

func (a userServicePositiveCase2) Get(_ context.Context, request *pb.UserRequest) (*pb.UserResponse, error) {
	if request.GetUser().GetUpk() == "" && request.GetUser().GetPersonalKey() != "" {
		return &pb.UserResponse{User: &pb.User{Upk: request.GetUser().GetUpk()}, Status: pb.Status_OK}, nil
	}
	return &pb.UserResponse{User: &pb.User{}, Status: pb.Status_FAIL}, nil
}

type userServiceNegativeCase1 struct {
	pb.UnimplementedUserServiceServer
}

func (a userServiceNegativeCase1) Get(_ context.Context, _ *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{User: &pb.User{}, Status: pb.Status_FAIL}, nil
}

func grpcServeUserServiceServer(ctx context.Context, prop env.Properties, srv pb.UserServiceServer, up chan struct{}) {
	listen, err := net.Listen("tcp", prop.ExternalAuthGRPCAddress())
	tool.IfErrorThenPanic(err)
	opts := []grpc.ServerOption{grpc.Creds(local.NewCredentials())}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUserServiceServer(grpcServer, srv)
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()
	if up != nil {
		close(up)
	}
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
