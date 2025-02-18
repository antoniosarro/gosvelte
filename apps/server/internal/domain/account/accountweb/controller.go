package accountweb

import (
	"context"
	"net/http"

	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/account"
	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/httperrors"
	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/validate"
	"github.com/antoniosarro/gosvelte/apps/server/internal/web/webcontext"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type iUsecase interface {
	Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.AccountDTO, error)
	ChangePassword(ctx context.Context, cp *account.ChangePasswordDTO, email string) error
	Me(ctx context.Context, accountID uuid.UUID) (*account.AccountDTO, error)
}

type controller struct {
	accountUC iUsecase
	log       *logger.Log
}

func newController(accountUC iUsecase, log *logger.Log) *controller {
	return &controller{accountUC, log}
}

func (con *controller) register(c echo.Context) error {
	dto := new(account.NewAccouuntDTO)

	if err := c.Bind(dto); err != nil {
		return httperrors.New(httperrors.InvalidArgument, "Invalid given JSON")
	}

	if err := dto.Validate(); err != nil {
		e := httperrors.New(httperrors.InvalidArgument, "Not valid given JSON")

		validateErrs := validate.SplitErrors(err)
		for _, s := range validateErrs {
			e.AddDetail(s)
		}

		return e
	}

	a, err := con.accountUC.Register(c.Request().Context(), dto)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, a)
}

func (con *controller) me(c echo.Context) error {
	claims := webcontext.GetAccessTokenClaims(c.Request().Context())
	if claims == nil {
		e := httperrors.New(httperrors.Unauthenticated, "Could not give access to this resource")
		e.AddDetail("token: access_token is missing")
		return e
	}

	a, err := con.accountUC.Me(c.Request().Context(), claims.AccountID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, a)
}
