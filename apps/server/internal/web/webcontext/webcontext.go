package webcontext

import (
	"context"

	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/token"
)

type contextKey string

const (
	accessClaimsKey  contextKey = "access_claims_key"
	refreshClaimsKey contextKey = "refresh_claims_key"
	accessKey        contextKey = "access_key"
	refreshKey       contextKey = "refresh_key"
)

func SetAccessTokenClaims(ctx context.Context, cl *token.AccessTokenClaims) context.Context {
	return context.WithValue(ctx, accessClaimsKey, cl)
}

func GetAccessTokenClaims(ctx context.Context) *token.AccessTokenClaims {
	val, ok := ctx.Value(accessClaimsKey).(*token.AccessTokenClaims)

	if !ok {
		return &token.AccessTokenClaims{}
	}

	return val
}

func SetAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, accessKey, token)
}

func GetAccessToken(ctx context.Context) string {
	val, ok := ctx.Value(accessKey).(string)

	if !ok {
		return ""
	}

	return val
}

func SetRefreshToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, refreshKey, token)
}

func GetRefreshToken(ctx context.Context) string {
	val, ok := ctx.Value(refreshKey).(string)

	if !ok {
		return ""
	}

	return val
}

func SetRefreshTokenClaims(ctx context.Context, cl *token.RefreshTokenClaims) context.Context {
	return context.WithValue(ctx, refreshClaimsKey, cl)
}

func GetRefreshTokenClaims(ctx context.Context) *token.RefreshTokenClaims {
	val, ok := ctx.Value(refreshClaimsKey).(*token.RefreshTokenClaims)

	if !ok {
		return &token.RefreshTokenClaims{}
	}

	return val
}
