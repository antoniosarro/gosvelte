package authweb

import (
	"context"
	"net/http"

	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/auth"
	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/httperrors"
	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/token"
	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/validate"
	"github.com/antoniosarro/gosvelte/apps/server/internal/web/webcontext"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type iUsecase interface {
	Login(ctx context.Context, email, password string) (*auth.AuthDTO, error)
	Logout(ctx context.Context, accessToken, refreshToken string, atc *token.AccessTokenClaims, rtc *token.RefreshTokenClaims) error
	Refresh(ctx context.Context, refreshToken string, accountID uuid.UUID) (*auth.AuthDTO, error)
}

type controller struct {
	authUC iUsecase
	log    *logger.Log
}

func newController(authUC iUsecase, log *logger.Log) *controller {
	return &controller{authUC, log}
}

// Login godoc
//
//	@Summary		Authenticate a user
//	@Description	Perform user login
//	@ID				auth-login
//	@Tags			Auth Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body		auth.LoginDTO	true	"User's credentials"
//	@Success		200		{object}	auth.AuthDTO
//	@Failure		401		{object}	error
//	@Router			/login [post]
func (con *controller) login(c echo.Context) error {
	dto := new(auth.LoginDTO)

	if err := c.Bind(dto); err != nil {
		return httperrors.New(httperrors.InvalidArgument, "Given JSON is invalid")
	}

	if err := dto.Validate(); err != nil {
		e := httperrors.New(httperrors.InvalidArgument, "Given JSON is out of validation rules")

		validationErrs := validate.SplitErrors(err)
		for _, s := range validationErrs {
			e.AddDetail(s)
		}

		return e
	}

	a, err := con.authUC.Login(c.Request().Context(), dto.Email, dto.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, a)
}

// Logout godoc
//
//	@Summary		Disconnect a user
//	@Description	Perform user logout
//	@ID				auth-logout
//	@Tags			Auth Actions
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		401		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/logout [post]
func (con *controller) logout(c echo.Context) error {
	at := webcontext.GetAccessToken(c.Request().Context())
	rt := webcontext.GetRefreshToken(c.Request().Context())

	if at == "" || rt == "" {
		e := httperrors.New(httperrors.Unauthenticated, "Could not give access to this resource")
		if at == "" {
			e.AddDetail("token: access_token is missing")
		}

		if rt == "" {
			e.AddDetail("token: refresh_token is missing")
		}

		return e
	}

	atc := webcontext.GetAccessTokenClaims(c.Request().Context())
	rtc := webcontext.GetRefreshTokenClaims(c.Request().Context())

	if atc == nil || rtc == nil {
		e := httperrors.New(httperrors.Unauthenticated, "Could not give access to this resource")
		e.AddDetail("data: claims are not found")
		return e
	}

	err := con.authUC.Logout(c.Request().Context(), at, rt, atc, rtc)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// Refresh godoc
//
//	@Summary		Refresh user token
//	@Description	Refresh user token
//	@ID				auth-refresh
//	@Tags			Auth Actions
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	auth.AuthDTO
//	@Failure		401		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/refresh [post]
func (con *controller) refreshToken(c echo.Context) error {
	rt := webcontext.GetRefreshToken(c.Request().Context())

	if rt == "" {
		e := httperrors.New(httperrors.Unauthenticated, "Could not give access to this resource")
		e.AddDetail("token: refresh_token is missing")
		return e
	}

	claims := webcontext.GetRefreshTokenClaims(c.Request().Context())

	if claims == nil {
		e := httperrors.New(httperrors.Unauthenticated, "Could not give access to this resource")
		e.AddDetail("data: claims are not found")
		return e
	}

	a, err := con.authUC.Refresh(c.Request().Context(), rt, claims.AccountID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, a)
}
