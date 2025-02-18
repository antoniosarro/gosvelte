package auth

import (
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/account"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type AuthDTO struct {
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	Account      *account.AccountDTO `json:"account"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d LoginDTO) Validate() error {
	return validation.ValidateStruct(
		&d,
		validation.Field(&d.Email, validation.Required, is.Email),
		validation.Field(&d.Password, validation.Required),
	)
}
