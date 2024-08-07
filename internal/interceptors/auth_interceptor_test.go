/*
 * This file was last modified at 2024-08-04 22:01 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * auth_interceptor_test.go
 * $Id$
 */
//!+

// Package interceptors TODO.
package interceptors

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gofavorites/internal/controllers/dto"
	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/jwt"
	"github.com/vskurikhin/gofavorites/internal/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/local"

	pb "github.com/vskurikhin/gofavorites/proto"
)

func TestAuthInterceptor(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "test #1 positive",
			fRun: testPositive1,
		},
		{
			name: "test #2 negative #1",
			fRun: testNegative1,
		},
		{
			name: "test #3 negative #2",
			fRun: testNegative2,
		},
		{
			name: "test #4 negative #3",
			fRun: testNegative3,
		},
		{
			name: "test #5 negative #4",
			fRun: testNegative4,
		},
	}

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testPositive1(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65285+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx, address, favoritesServicePositive{}, up)
	<-up
	manager := jwt.GetJWTManager(env.GetProperties())
	token, err := manager.Generate(dto.SignInRequest{
		UserName: "test",
		Password: "password",
	})
	if err != nil {
		t.Fail()
	}
	conn, client, err := makeFavoritesServiceClient(t, address, token)
	defer func() { _ = conn.Close() }()

	var request pb.FavoritesRequest
	resp, err := client.Get(ctx, &request)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testNegative1(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65285+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeTestServiceServer(ctx, address, &testServicePositive{}, up)
	<-up
	manager := jwt.GetJWTManager(env.GetProperties())
	token, err := manager.Generate(dto.SignInRequest{
		UserName: "test",
		Password: "password",
	})
	if err != nil {
		t.Fail()
	}
	conn, client, err := makeTestServiceClient(t, address, token)
	defer func() { _ = conn.Close() }()
	var request testpb.PingRequest
	resp, err := client.Ping(ctx, &request)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testNegative2(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65285+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx, address, favoritesServicePositive{}, up)
	<-up
	conn, client, err := makeFavoritesServiceClient(t, address, "")
	defer func() { _ = conn.Close() }()
	var request pb.FavoritesRequest
	resp, err := client.Get(ctx, &request)

	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

type mdIncomingKey struct{}

func testNegative3(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65285+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	ctx0 := context.WithValue(context.Background(), mdIncomingKey{}, nil)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx0, address, favoritesServicePositive{}, up)
	<-up
	manager := jwt.GetJWTManager(env.GetProperties())
	token, err := manager.Generate(dto.SignInRequest{
		UserName: "test",
		Password: "password",
	})
	if err != nil {
		t.Fail()
	}
	conn, client, err := makeFavoritesServiceClient(t, address, token)
	defer func() { _ = conn.Close() }()

	var request pb.FavoritesRequest
	ctx1 := context.WithValue(context.Background(), mdIncomingKey{}, nil)
	resp, err := client.Get(ctx1, &request)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func testNegative4(t *testing.T) {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	address := fmt.Sprintf("127.0.0.1:%d", 65285+rnd.Intn(34))

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	ctx0 := context.WithValue(context.Background(), mdIncomingKey{}, nil)
	defer func() {
		cancel()
		ctx.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeFavoritesServiceServer(ctx0, address, favoritesServicePositive{}, up)
	<-up
	conn, client, err := makeFavoritesServiceClientWithoutToken(t, address)
	defer func() { _ = conn.Close() }()

	var request pb.FavoritesRequest
	ctx1 := context.WithValue(context.Background(), mdIncomingKey{}, nil)
	resp, err := client.Get(ctx1, &request)

	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func grpcServeFavoritesServiceServer(ctx context.Context, address string, srv pb.FavoritesServiceServer, up chan struct{}) {
	listen, err := net.Listen("tcp", address)
	tool.IfErrorThenPanic(err)
	authInterceptor := GetAuthInterceptor(env.GetProperties())
	opts := []grpc.ServerOption{
		grpc.Creds(local.NewCredentials()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterFavoritesServiceServer(grpcServer, srv)
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

func grpcServeTestServiceServer(ctx context.Context, address string, srv testpb.TestServiceServer, up chan struct{}) {
	listen, err := net.Listen("tcp", address)
	tool.IfErrorThenPanic(err)
	authInterceptor := GetAuthInterceptor(env.GetProperties())
	opts := []grpc.ServerOption{
		grpc.Creds(local.NewCredentials()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	}
	grpcServer := grpc.NewServer(opts...)
	testpb.RegisterTestServiceServer(grpcServer, srv)
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

type jwtAuthorization struct {
	Token string
}

func (a jwtAuthorization) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"authorization": a.Token}, nil
}

func (a jwtAuthorization) RequireTransportSecurity() bool {
	return false
}

var _ credentials.PerRPCCredentials = (*jwtAuthorization)(nil)

func makeFavoritesServiceClient(t *testing.T, address, token string) (*grpc.ClientConn, pb.FavoritesServiceClient, error) {

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(jwtAuthorization{Token: token}),
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		t.Fail()
	}
	client := pb.NewFavoritesServiceClient(conn)

	return conn, client, err
}

func makeFavoritesServiceClientWithoutToken(t *testing.T, address string) (*grpc.ClientConn, pb.FavoritesServiceClient, error) {

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		t.Fail()
	}
	client := pb.NewFavoritesServiceClient(conn)

	return conn, client, err
}

func makeTestServiceClient(t *testing.T, address, token string) (*grpc.ClientConn, testpb.TestServiceClient, error) {

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(jwtAuthorization{Token: token}),
	}

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		t.Fail()
	}
	client := testpb.NewTestServiceClient(conn)

	return conn, client, err
}

type FavoritesServiceTest interface {
	pb.FavoritesServiceServer
}

type favoritesServicePositive struct {
	pb.UnimplementedFavoritesServiceServer
}

type testServicePositive struct {
	testpb.UnimplementedTestServiceServer
}

var _ FavoritesServiceTest = (*favoritesServicePositive)(nil)

func (f favoritesServicePositive) Get(_ context.Context, _ *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {
	return &pb.FavoritesResponse{}, nil
}

func (f favoritesServicePositive) GetForUser(_ context.Context, _ *pb.UserFavoritesRequest) (*pb.UserFavoritesResponse, error) {
	return &pb.UserFavoritesResponse{}, nil
}

func (f favoritesServicePositive) Set(_ context.Context, _ *pb.FavoritesRequest) (*pb.FavoritesResponse, error) {
	return &pb.FavoritesResponse{}, nil
}

func (t testServicePositive) PingEmpty(ctx context.Context, request *testpb.PingEmptyRequest) (*testpb.PingEmptyResponse, error) {
	return &testpb.PingEmptyResponse{}, nil
}

func (t testServicePositive) Ping(ctx context.Context, request *testpb.PingRequest) (*testpb.PingResponse, error) {
	return &testpb.PingResponse{}, nil
}

func (t testServicePositive) PingError(ctx context.Context, request *testpb.PingErrorRequest) (*testpb.PingErrorResponse, error) {
	return &testpb.PingErrorResponse{}, nil
}

func (t testServicePositive) PingList(request *testpb.PingListRequest, server testpb.TestService_PingListServer) error {
	return nil
}

func (t testServicePositive) PingStream(server testpb.TestService_PingStreamServer) error {
	return nil
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
