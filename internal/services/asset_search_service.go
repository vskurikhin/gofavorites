/*
 * This file was last modified at 2024-08-06 20:32 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * asset_search_service.go
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
	"github.com/vskurikhin/gofavorites/internal/tool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	pb "github.com/vskurikhin/gofavorites/proto"
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
	sLog             *slog.Logger
}

var _ AssetSearchService = (*assetSearchService)(nil)
var (
	onceAssetSearch = new(sync.Once)
	assetSearchServ *assetSearchService
)

// GetAssetSearchService — потокобезопасное (thread-safe) создание
// сервиса поиска биржевых инструментов в базе данных или во внешней системе.
func GetAssetSearchService(prop env.Properties) AssetSearchService {

	onceAssetSearch.Do(func() {
		assetSearchServ = new(assetSearchService)
		opts := []grpc.DialOption{
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		}
		tlsCredentials, err := tool.LoadClientTLSCredentials(prop.Config().GRPCTLSCAFile())
		if err != nil {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
		}
		assetSearchServ.assetGRPCAddress = prop.ExternalAssetGRPCAddress()
		assetSearchServ.opts = opts
		assetSearchServ.repoAsset = repo.GetAssetPostgresCachedRepo(prop)
		assetSearchServ.requestInterval = prop.ExternalRequestTimeoutInterval()
		assetSearchServ.sLog = prop.Logger()
	})
	return assetSearchServ
}

const cntAssetSearchLookupJobs = 2

// Lookup поиск биржевого инструмента в базе данных или во внешней системе.
func (a *assetSearchService) Lookup(ctx context.Context, isin string) bool {

	var wg sync.WaitGroup
	wg.Add(cntAssetSearchLookupJobs)

	quit := make(chan struct{})
	results := make(chan bool, cntAssetSearchLookupJobs)

	go func() {
		defer wg.Done()
		results <- a.dbLookup(ctx, isin)
	}()
	go func() {
		defer wg.Done()
		results <- a.grpcLookup(ctx, isin)
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
	var request pb.AssetRequest
	request.Asset = &pb.Asset{Isin: isin}
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer func() {
		cancel()
	}()
	resp, err := c.Get(ctx, &request)

	for i := 1; err != nil && tool.IsUpperBoundWithSleep(i, 300, a.requestInterval); i++ {
		a.sLog.ErrorContext(ctx,
			env.MSG+"AssetSearchService.grpcLookup",
			"msg", "asset search service gRPC lookup",
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
