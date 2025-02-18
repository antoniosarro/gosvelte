package middlewares

import (
	"errors"
	"fmt"

	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/httperrors"
	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/token"
	"github.com/antoniosarro/gosvelte/apps/server/internal/web/webcontext"
	"github.com/labstack/echo/v4"
)

func (mid *Middleware) Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		t, err := token.ExtractBearerToken(authHeader)
		if err != nil {
			e := httperrors.New(httperrors.Unauthenticated, "Authorization header missing")
			e.AddDetail(fmt.Sprintf("data: %s", err))
			return e
		}

		// whenever token exist in blacklist, it will return error
		if val, _ := mid.cache.Get(c.Request().Context(), t).Result(); val != "" {
			e := httperrors.New(httperrors.Unauthenticated, "User already logged out")
			e.AddDetail(fmt.Sprintf("data: %s", err))
			return e
		}

		claims, err := token.ValidateAccess(mid.conf, t)
		if err != nil {
			if errors.Is(err, token.ErrInvalidToken) {
				e := httperrors.New(httperrors.Unauthenticated, "Invalid access token")
				e.AddDetail("token: access_token on bearer authentication header is invalid")
				return e
			}

			return httperrors.New(httperrors.Internal, "Something went wrong")
		}

		ctx := webcontext.SetAccessTokenClaims(c.Request().Context(), claims)
		ctx = webcontext.SetAccessToken(ctx, t)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (mid *Middleware) RefreshAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		refreshToken := c.Request().Header.Get("RF-Token")
		if refreshToken == "" {
			e := httperrors.New(httperrors.Unauthenticated, "RT-Token header missing")
			e.AddDetail("token: expected RT-Token header with refresh token as its value")
		}

		claims, err := token.ValidateRefresh(mid.conf, refreshToken)
		if err != nil {
			e := httperrors.New(httperrors.Unauthenticated, "Invalid refresh token")
			e.AddDetail("token: refresh_token on RT-Token header is invalid")
			return e
		}

		// whenever token exist in blacklist, it will return error
		if val, _ := mid.cache.Get(c.Request().Context(), claims.AccountID.String()).Result(); val != "" {
			e := httperrors.New(httperrors.Unauthenticated, "User already logged out")
			e.AddDetail(fmt.Sprintf("data: %s", err))
			return e
		}

		ctx := webcontext.SetRefreshTokenClaims(c.Request().Context(), claims)
		ctx = webcontext.SetRefreshToken(ctx, refreshToken)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
