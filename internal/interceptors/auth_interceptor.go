/*
 * This file was last modified at 2024-08-06 20:17 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * auth_interceptor.go
 * $Id$
 */
//!+

// Package interceptors TODO.
package interceptors

import (
	"context"
	"log/slog"
	"sync"

	"github.com/vskurikhin/gofavorites/internal/env"
	"github.com/vskurikhin/gofavorites/internal/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor interface {
	Stream() grpc.StreamServerInterceptor
	Unary() grpc.UnaryServerInterceptor
}

type authInterceptor struct {
	accessibleRoles map[string][]string
	jwtManager      jwt.Manager
	sLog            *slog.Logger
}

var _ AuthInterceptor = (*authInterceptor)(nil)
var (
	ErrAuthorizationTokenIsNotProvided = status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	ErrMetadataIsNotProvided           = status.Errorf(codes.Unauthenticated, "metadata is not provided")
	ErrNoPermissionToAccessThisRPC     = status.Error(codes.PermissionDenied, "no permission to access this RPC")
	authInter                          *authInterceptor
	onceAuth                           = new(sync.Once)
)

func GetAuthInterceptor(prop env.Properties) AuthInterceptor {

	onceAuth.Do(func() {
		authInter = new(authInterceptor)
		authInter.accessibleRoles = authMethods()
		authInter.jwtManager = jwt.GetJWTManager(prop)
		authInter.sLog = prop.Logger()
	})
	return authInter
}

func (a *authInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		a.sLog.DebugContext(stream.Context(), "--> stream interceptor: ", "FullMethod", info.FullMethod)

		err := a.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (a *authInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		a.sLog.DebugContext(ctx, "--> unary interceptor: ", "FullMethod", info.FullMethod)

		err := a.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (a *authInterceptor) authorize(ctx context.Context, method string) error {

	accessibleRoles, ok := a.accessibleRoles[method]

	if !ok {
		return nil
	}
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return ErrMetadataIsNotProvided
	}
	values := md["authorization"]

	if len(values) == 0 {
		return ErrAuthorizationTokenIsNotProvided
	}
	accessToken := values[0]
	claims, err := a.jwtManager.Verify(accessToken)

	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	for _, role := range accessibleRoles {
		if role == claims.Role() {
			return nil
		}
	}
	return ErrNoPermissionToAccessThisRPC
}

func authMethods() map[string][]string {
	const methodServicePath = "/proto.FavoritesService/"

	return map[string][]string{
		methodServicePath + "Get":        {"USER"},
		methodServicePath + "GetForUser": {"USER"},
		methodServicePath + "Set":        {"USER"},
	}
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
