package accountrepo

import (
	"context"

	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/account"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (dbrepo *repository) GetOneByEmail(ctx context.Context, email string) (*account.AccountDTO, error) {
	model := new(Model)

	query := `
		SELECT * FROM users
		WHERE email = $1
		LIMIT 1
	`

	if err := dbrepo.QueryRowxContext(ctx, query, email).StructScan(model); err != nil {
		return nil, err
	}

	return model.intoDTO(), nil
}

func (dbrepo *repository) GetOne(ctx context.Context, id uuid.UUID) (*account.AccountDTO, error) {
	model := new(Model)

	query := `
		SELECT * FROM users
		WHERE id = $1
		LIMIT 1
	`

	if err := dbrepo.QueryRowxContext(ctx, query, id).StructScan(model); err != nil {
		return nil, err
	}

	return model.intoDTO(), nil
}

func (dbrepo *repository) Create(ctx context.Context, a *account.AccountDTO) error {
	query := `
		INSERT INTO users
			(id, firstname, lastname, email, password, role, created_at, updated_at)
		VALUES
			(:id, :firstname, :lastname, :email, :password, :role, :created_at, :updated_at)
	`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(a))
	return err
}

func (repo *repository) ChangePassword(ctx context.Context, email, password string) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE email = $2
	`

	_, err := repo.ExecContext(ctx, query, password, email)
	return err
}
