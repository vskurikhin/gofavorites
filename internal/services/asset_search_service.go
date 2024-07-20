/*
 * This file was last modified at 2024-07-21 08:50 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * asset_search_service.go
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

type AssetSearchService interface {
	Lookup(ctx context.Context, isin string) bool
}

type assetSearchService struct {
	assetGRPCAddress string
	creds            credentials.TransportCredentials
	opts             []grpc.DialOption
	repoAsset        domain.Repo[*entity.Asset]
	requestInterval  time.Duration
}

var _ AssetSearchService = (*assetSearchService)(nil)
var (
	onceAssetSearch = new(sync.Once)
	assetSearchSrv  *assetSearchService
)

func GetAssetSearchService(prop env.Properties) AssetSearchService {

	onceAssetSearch.Do(func() {
		assetSearchSrv = new(assetSearchService)
		opts := []grpc.DialOption{
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		}
		tlsCredentials, err := tool.LoadAgentTLSCredentials(prop.Config().GRPCTLSCAFile())
		if err != nil {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
		}
		assetSearchSrv.assetGRPCAddress = prop.ExternalAssetGRPCAddress()
		assetSearchSrv.opts = opts
		assetSearchSrv.repoAsset = repo.GetAssetPostgresCachedRepo(prop)
		assetSearchSrv.requestInterval = prop.ExternalRequestTimeoutInterval()
	})
	return assetSearchSrv
}

func (a *assetSearchService) Lookup(ctx context.Context, isin string) bool {

	if a.dbLookup(ctx, isin) {
		return true
	}
	if a.grpcLookup(ctx, isin) {
		return true
	}
	return false
}

func (a *assetSearchService) dbLookup(ctx context.Context, isin string) bool {

	asset, err := entity.GetAsset(ctx, a.repoAsset, isin)

	if entity.IsAssetNotFound(asset, err) {
		return false
	}
	return true
}

func (a *assetSearchService) grpcLookup(ctx context.Context, isin string) bool {

	conn, err := grpc.NewClient(a.assetGRPCAddress, a.opts...)
	if err != nil {
		return false
	}
	defer func() { _ = conn.Close() }()
	c := pb.NewAssetServiceClient(conn)
	ctx, cancel := context.WithTimeout(ctx, a.requestInterval)
	defer func() {
		cancel()
		ctx.Done()
	}()
	var request pb.AssetRequest
	request.Asset = &pb.Asset{Isin: isin}
	resp, err := c.Get(ctx, &request)

	for i := 1; err != nil && tool.IsUpperBound(i, a.requestInterval); i++ {
		slog.Warn(env.MSG+" AssetSearchService.grpcLookup", "err", err)
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
