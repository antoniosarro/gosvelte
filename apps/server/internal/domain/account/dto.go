package account

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type AccountDTO struct {
	ID        uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewAccouuntDTO struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

func (na NewAccouuntDTO) Validate() error {
	return validation.ValidateStruct(
		&na,
		validation.Field(&na.Firstname, validation.Required, validation.Length(1, 30)),
		validation.Field(&na.Lastname, validation.Required, validation.Length(1, 30)),
		validation.Field(&na.Email, validation.Required, is.Email),
		validation.Field(&na.Password, validation.Required),
	)
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (d ChangePasswordDTO) Validate() error {
	return validation.ValidateStruct(
		&d,
		validation.Field(&d.OldPassword, validation.Required),
		validation.Field(&d.NewPassword, validation.Required),
	)
}
