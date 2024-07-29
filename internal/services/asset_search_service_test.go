/*
 * This file was last modified at 2024-07-20 19:34 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * asset_search_service_test.go
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

func TestAssetSearchService(t *testing.T) {
	var tests = []struct {
		name string
		fRun func(*testing.T)
	}{
		{
			name: "positive test #1 Asset Service Lookup case #1",
			fRun: testAssetSearchServiceLookupPositiveCase1,
		},
		{
			name: "positive test #2 Asset Service Lookup case #2",
			fRun: testAssetSearchServiceLookupPositiveCase2,
		},
		{
			name: "negative test #3 Asset Service Lookup case #2",
			fRun: testAssetSearchServiceLookupNegativeCase2,
		},
		{
			name: "negative test #4 GetAssetSearchService",
			fRun: testGetAssetSearchService,
		},
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("ASSET_GRPC_ADDRESS", fmt.Sprintf("127.0.0.1:%d", 65501+rnd.Intn(34)))
	t.Setenv("AUTH_GRPC_ADDRESS", fmt.Sprintf("127.0.0.1:%d", 65501+rnd.Intn(34)))
	t.Setenv("REQUEST_TIMEOUT_INTERVAL_MS", "500")
	t.Setenv("UPK_PRIVATE_KEY_FILE", "test_private-key.pem")
	t.Setenv("UPK_PUBLIC_KEY_FILE", "test_public-key.pem")
	t.Setenv("UPK_SECRET", "qYhaPtg+PIQtBhAU5fHCeQw7XIF3WLKoLPZnJgq1H//DDOB8o2qrP9goVCUZldOdwqLAHxWOGHuvXcwaIFRrD8I3Hz5tRCgCeI+cEZD9h4c4h6ADSjkcrPXg5eRwnANasBkKKZQz8noYwvt9Z9p7HdOtrBmQOi7OVjTfY0T2SnI=")

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testGetAssetSearchService(t *testing.T) {
	prop := env.GetProperties()
	got := GetAssetSearchService(prop)
	assert.NotNil(t, got)
}

func testAssetSearchServiceLookupPositiveCase1(t *testing.T) {
	prop := env.GetProperties()
	ctrl := gomock.NewController(t)
	repoAsset := NewMockRepo[*entity.Asset](ctrl)
	asset := entity.MakeAsset("test", entity.AssetType{}, entity.DefaultTAttributes())
	repoAsset.
		EXPECT().
		Get(context.Background(), gomock.Any(), gomock.Any()).
		Return(&asset, nil).
		AnyTimes()
	assetSearchService := getAssetSearchService(prop, repoAsset)
	got := assetSearchService.Lookup(context.Background(), asset.Isin())
	assert.True(t, got)
}

func testAssetSearchServiceLookupPositiveCase2(t *testing.T) {
	prop := env.GetProperties()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer func() {
		cancel()
		time.Sleep(500 * time.Millisecond)
	}()
	up := make(chan struct{})
	go grpcServeAssetServiceServer(ctx, prop, assetServicePositive{}, up)
	<-up

	ctrl := gomock.NewController(t)
	repoAsset := NewMockRepo[*entity.Asset](ctrl)
	asset := entity.MakeAsset("test", entity.AssetType{}, entity.DefaultTAttributes())
	repoAsset.
		EXPECT().
		Get(context.Background(), gomock.Any(), gomock.Any()).
		Return(&asset, repo.ErrNotFound).
		AnyTimes()

	assetSearchService := getAssetSearchService(prop, repoAsset)
	got := assetSearchService.Lookup(context.Background(), asset.Isin())
	assert.True(t, got)
}

func testAssetSearchServiceLookupNegativeCase2(t *testing.T) {
	prop := env.GetProperties()
	ctrl := gomock.NewController(t)
	repoAsset := NewMockRepo[*entity.Asset](ctrl)
	asset := entity.MakeAsset("test", entity.AssetType{}, entity.DefaultTAttributes())
	repoAsset.
		EXPECT().
		Get(context.Background(), gomock.Any(), gomock.Any()).
		Return(nil, repo.ErrNotFound).
		AnyTimes()
	assetSearchService := getAssetSearchService(prop, repoAsset)
	got := assetSearchService.Lookup(context.Background(), asset.Isin())
	assert.False(t, got)
}

func getAssetSearchService(prop env.Properties, repoAsset domain.Repo[*entity.Asset]) AssetSearchService {
	assetSearchServ = new(assetSearchService)
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithNoProxy(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	assetSearchServ.assetGRPCAddress = prop.ExternalAssetGRPCAddress()
	assetSearchServ.opts = opts
	assetSearchServ.repoAsset = repoAsset
	assetSearchServ.requestInterval = prop.ExternalRequestTimeoutInterval()
	return assetSearchServ
}

type assetServicePositive struct {
	pb.UnimplementedAssetServiceServer
}

func (a assetServicePositive) Get(_ context.Context, request *pb.AssetRequest) (*pb.AssetResponse, error) {
	return &pb.AssetResponse{Asset: &pb.Asset{Isin: request.GetAsset().GetIsin()}, Status: pb.Status_OK}, nil
}

func grpcServeAssetServiceServer(ctx context.Context, prop env.Properties, srv pb.AssetServiceServer, up chan struct{}) {
	listen, err := net.Listen("tcp", prop.ExternalAssetGRPCAddress())
	tool.IfErrorThenPanic(err)
	opts := []grpc.ServerOption{grpc.Creds(local.NewCredentials())}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAssetServiceServer(grpcServer, srv)
	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()
	if up != nil {
		close(up)
	}
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
