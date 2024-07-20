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

	assert.NotNil(t, t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fRun(t)
		})
	}
}

func testGetAssetSearchService(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("ASSET_GRPC_ADDRESS", fmt.Sprintf("localhost:%d", 65500+rnd.Intn(34)))
	t.Setenv("REQUEST_TIMEOUT_INTERVAL_MS", "500")
	prop := env.GetProperties()
	got := GetAssetSearchService(prop)
	assert.NotNil(t, got)
}

func testAssetSearchServiceLookupPositiveCase1(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("ASSET_GRPC_ADDRESS", fmt.Sprintf("localhost:%d", 65500+rnd.Intn(34)))
	t.Setenv("REQUEST_TIMEOUT_INTERVAL_MS", "500")
	prop := env.GetProperties()
	ctrl := gomock.NewController(t)
	repoAsset := NewMockRepo[*entity.Asset](ctrl)
	asset := entity.MakeAsset("test", entity.AssetType{}, entity.DefaultTAttributes())
	repoAsset.
		EXPECT().
		Get(context.TODO(), gomock.Any(), gomock.Any()).
		Return(&asset, nil).
		AnyTimes()
	assetSearchService := getAssetSearchService(prop, repoAsset)
	got := assetSearchService.Lookup(context.TODO(), asset.Isin())
	assert.True(t, got)
}

func testAssetSearchServiceLookupPositiveCase2(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("ASSET_GRPC_ADDRESS", fmt.Sprintf("localhost:%d", 65500+rnd.Intn(34)))
	t.Setenv("REQUEST_TIMEOUT_INTERVAL_MS", "500")
	prop := env.GetProperties()
	ctrl := gomock.NewController(t)
	repoAsset := NewMockRepo[*entity.Asset](ctrl)
	asset := entity.MakeAsset("test", entity.AssetType{}, entity.DefaultTAttributes())
	repoAsset.
		EXPECT().
		Get(context.TODO(), gomock.Any(), gomock.Any()).
		Return(&asset, repo.ErrNotFound).
		AnyTimes()
	ctx, cancel := context.WithTimeout(context.TODO(), 500*time.Millisecond)
	defer func() {
		cancel()
		time.Sleep(100 * time.Millisecond)
	}()
	go grpcServeAssetServiceServer(ctx, prop, assetServicePositive{})
	assetSearchService := getAssetSearchService(prop, repoAsset)
	got := assetSearchService.Lookup(context.TODO(), asset.Isin())
	assert.True(t, got)
}

func testAssetSearchServiceLookupNegativeCase2(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	t.Setenv("GO_FAVORITES_SKIP_LOAD_CONFIG", "True")
	t.Setenv("ASSET_GRPC_ADDRESS", fmt.Sprintf("localhost:%d", 65500+rnd.Intn(34)))
	t.Setenv("REQUEST_TIMEOUT_INTERVAL_MS", "500")
	prop := env.GetProperties()
	ctrl := gomock.NewController(t)
	repoAsset := NewMockRepo[*entity.Asset](ctrl)
	asset := entity.MakeAsset("test", entity.AssetType{}, entity.DefaultTAttributes())
	repoAsset.
		EXPECT().
		Get(context.TODO(), gomock.Any(), gomock.Any()).
		Return(nil, repo.ErrNotFound).
		AnyTimes()
	assetSearchService := getAssetSearchService(prop, repoAsset)
	got := assetSearchService.Lookup(context.TODO(), asset.Isin())
	assert.False(t, got)
}

func getAssetSearchService(prop env.Properties, repoAsset domain.Repo[*entity.Asset]) AssetSearchService {
	assetSearchSrv = new(assetSearchService)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	assetSearchSrv.assetGRPCAddress = prop.ExternalAssetGRPCAddress()
	assetSearchSrv.opts = opts
	assetSearchSrv.repoAsset = repoAsset
	assetSearchSrv.requestInterval = prop.ExternalRequestTimeoutInterval()
	return assetSearchSrv
}

type assetServicePositive struct {
	pb.UnimplementedAssetServiceServer
}

func (a assetServicePositive) Get(_ context.Context, request *pb.AssetRequest) (*pb.AssetResponse, error) {
	return &pb.AssetResponse{Asset: &pb.Asset{Isin: request.GetAsset().GetIsin()}, Status: pb.Status_OK}, nil
}

func grpcServeAssetServiceServer(ctx context.Context, prop env.Properties, srv pb.AssetServiceServer) {
	listen, err := net.Listen("tcp", prop.ExternalAssetGRPCAddress())
	tool.IfErrorThenPanic(err)
	opts := []grpc.ServerOption{grpc.Creds(local.NewCredentials())}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAssetServiceServer(grpcServer, srv)
	go func() {
		for {
			select {
			case <-ctx.Done():
				grpcServer.Stop()
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
