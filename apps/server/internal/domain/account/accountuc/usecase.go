package accountuc

import (
	"context"
	"fmt"
	"time"

	"github.com/antoniosarro/gosvelte/apps/server/config"
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/account"
	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/httperrors"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type iDBRepository interface {
	GetOne(ctx context.Context, id uuid.UUID) (*account.AccountDTO, error)
	GetOneByEmail(ctx context.Context, email string) (*account.AccountDTO, error)
	Create(ctx context.Context, a *account.AccountDTO) error
	ChangePassword(ctx context.Context, email, password string) error
}

type iCacheRepository interface {
	SetMe(ctx context.Context, accountPayload *account.AccountDTO) error
	GetMe(ctx context.Context, id uuid.UUID) (*account.AccountDTO, error)
}

type Usecase struct {
	conf      *config.Config
	log       *logger.Log
	dbRepo    iDBRepository
	cacheRepo iCacheRepository
}

func New(conf *config.Config, log *logger.Log, dbRepo iDBRepository, cacheRepo iCacheRepository) *Usecase {
	return &Usecase{
		conf:      conf,
		log:       log,
		dbRepo:    dbRepo,
		cacheRepo: cacheRepo,
	}
}

func (uc *Usecase) Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.AccountDTO, error) {
	exist, err := uc.dbRepo.GetOneByEmail(ctx, na.Email)
	if exist != nil && err == nil {
		err := httperrors.New(httperrors.AlreadyExists, "Email not available")
		err.AddDetail(fmt.Sprintf("data: email %s not available", na.Email))
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(na.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, httperrors.New(httperrors.Internal, "Failed hashing password")
	}

	now := time.Now()
	a := &account.AccountDTO{
		ID:        uuid.New(),
		Firstname: na.Firstname,
		Lastname:  na.Lastname,
		Email:     na.Email,
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.dbRepo.Create(ctx, a); err != nil {
		return nil, httperrors.New(httperrors.Internal, "Error inserting new user")
	}

	return a, nil

}

func (uc *Usecase) ChangePassword(ctx context.Context, cp *account.ChangePasswordDTO, email string) error {
	a, err := uc.dbRepo.GetOneByEmail(ctx, email)
	if err != nil {
		httperrors.New(httperrors.Internal, "Error retriving email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(cp.OldPassword)); err != nil {
		return httperrors.New(httperrors.InvalidCredentials, "Credentials are invalid, either email and/or password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cp.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return httperrors.New(httperrors.Internal, "Failed hashing password")
	}

	if err := uc.dbRepo.ChangePassword(ctx, email, string(hashedPassword)); err != nil {
		httperrors.New(httperrors.Internal, "Something went wrong")
	}

	return nil
}

func (uc *Usecase) Me(ctx context.Context, accountID uuid.UUID) (*account.AccountDTO, error) {
	cached, err := uc.cacheRepo.GetMe(ctx, accountID)
	if err != nil {
		return nil, httperrors.New(httperrors.Internal, "Could not retrieve data from cache")
	}

	if cached != nil {
		return cached, nil
	}

	res, err := uc.dbRepo.GetOne(ctx, accountID)
	if err != nil {
		return nil, httperrors.New(httperrors.Internal, "Could not retrieve data from db")
	}

	if err := uc.cacheRepo.SetMe(ctx, res); err != nil {
		return nil, httperrors.New(httperrors.Internal, "Could not set data into cache")
	}

	return res, nil
}
