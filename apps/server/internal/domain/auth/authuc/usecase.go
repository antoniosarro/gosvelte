package authuc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/antoniosarro/gosvelte/apps/server/config"
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/account"
	"github.com/antoniosarro/gosvelte/apps/server/internal/domain/auth"
	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/httperrors"
	"github.com/antoniosarro/gosvelte/apps/server/internal/sdk/token"
	"github.com/antoniosarro/gosvelte/apps/server/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type iAuthCacheRepo interface {
	AddAccessTokenToBlacklist(ctx context.Context, token string, exp time.Duration) error
	AddRefreshTokenToBlacklist(ctx context.Context, accountID uuid.UUID, token string, exp time.Duration) error
	RemoveRefreshTokenFromBlacklist(ctx context.Context, accountID uuid.UUID) error
}

type iAccountDBRepo interface {
	GetOneByEmail(ctx context.Context, email string) (*account.AccountDTO, error)
	GetOne(ctx context.Context, id uuid.UUID) (*account.AccountDTO, error)
}

type Usecase struct {
	authCacheRepo iAuthCacheRepo
	conf          *config.Config
	log           *logger.Log
	accountDBRepo iAccountDBRepo
}

func New(conf *config.Config, log *logger.Log, authCacheRepo iAuthCacheRepo, accountDBRepo iAccountDBRepo) *Usecase {
	return &Usecase{
		accountDBRepo: accountDBRepo,
		authCacheRepo: authCacheRepo,
		conf:          conf,
		log:           log,
	}
}

func (uc *Usecase) Login(ctx context.Context, email, password string) (*auth.AuthDTO, error) {
	a, err := uc.accountDBRepo.GetOneByEmail(ctx, email)

	if err != nil {
		return nil, httperrors.New(httperrors.InvalidCredentials, "Credentials are invalid, either email and/or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		return nil, httperrors.New(httperrors.InvalidCredentials, "Credentials are invalid, either email and/or password")
	}

	if err := uc.authCacheRepo.RemoveRefreshTokenFromBlacklist(ctx, a.ID); err != nil {
		e := httperrors.New(httperrors.Internal, "Failed to refresh token from blacklist")
		e.AddDetail(fmt.Sprintf("data: %v", err))
		return nil, e
	}

	wg := new(sync.WaitGroup)
	atCh, rtCh := make(chan *string, 1), make(chan *string, 1) //access token channel & refresh token channel
	wg.Add(2)

	go func(ch chan *string) {
		defer wg.Done()

		at, err := token.GenerateAccess(uc.conf, token.AccessTokenPayload{
			AccountID: a.ID,
			Email:     a.Email,
			Role:      a.Role,
		})
		if err != nil {
			e := httperrors.New(httperrors.Internal, fmt.Sprintf("Failed to generate access_token, %v", err))
			uc.log.Error(e.LogForDebug())
			ch <- nil // store nil pointer to channel
			return
		}
		ch <- &at
	}(atCh)

	go func(ch chan *string) {
		defer wg.Done()

		rt, err := token.GenerateRefresh(uc.conf, token.RefreshTokenPayload{
			AccountID: a.ID,
		})
		if err != nil {
			e := httperrors.New(httperrors.Internal, fmt.Sprintf("Failed to generate refresh_token, %v", err))
			uc.log.Error(e.LogForDebug())
			ch <- nil // store nil pointer to channel
			return
		}
		ch <- &rt
	}(rtCh)

	wg.Wait()
	at, rt := <-atCh, <-rtCh

	if at == nil || rt == nil {
		return nil, httperrors.New(httperrors.Internal, "Something went wrong")
	}

	return &auth.AuthDTO{
		Account:      a,
		AccessToken:  *at,
		RefreshToken: *rt,
	}, nil
}

func (uc *Usecase) Logout(ctx context.Context, accessToken, refreshToken string, atc *token.AccessTokenClaims, rtc *token.RefreshTokenClaims) error {
	atTime := token.RemainingTime(&atc.RegisteredClaims)
	rtTime := token.RemainingTime(&rtc.RegisteredClaims)

	chanErrs := make(chan error, 2)
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func(e chan error) {
		defer wg.Done()

		err := uc.authCacheRepo.AddAccessTokenToBlacklist(ctx, accessToken, atTime)
		if err != nil {
			e := httperrors.New(httperrors.Internal, "Failed to add access token from blacklist")
			e.AddDetail(fmt.Sprintf("data: %v", err))
			chanErrs <- e
		}
		chanErrs <- nil
	}(chanErrs)

	go func(e chan error) {
		defer wg.Done()

		err := uc.authCacheRepo.AddRefreshTokenToBlacklist(ctx, rtc.AccountID, refreshToken, rtTime)
		if err != nil {
			e := httperrors.New(httperrors.Internal, "Failed to add refresh token from blacklist")
			e.AddDetail(fmt.Sprintf("data: %v", err))
			chanErrs <- e
		}
		chanErrs <- nil
	}(chanErrs)

	wg.Wait()

	if e := <-chanErrs; e != nil {
		return e
	}

	return nil
}

func (uc *Usecase) Refresh(ctx context.Context, refreshToken string, accountID uuid.UUID) (*auth.AuthDTO, error) {
	a, err := uc.accountDBRepo.GetOne(ctx, accountID)
	if err != nil {
		return nil, httperrors.New(httperrors.NotFound, "Something went wrong")
	}

	if a == nil {
		e := httperrors.New(httperrors.NotFound, "Account could not be found")
		e.AddDetail(fmt.Sprintf("data: account with id %s not found", accountID))
		return nil, e
	}

	at, err := token.GenerateAccess(uc.conf, token.AccessTokenPayload{
		AccountID: a.ID,
		Email:     a.Email,
		Role:      a.Role,
	})

	if err != nil {
		return nil, httperrors.New(httperrors.NotFound, "Something went wrong")
	}

	return &auth.AuthDTO{
		Account:      a,
		AccessToken:  at,
		RefreshToken: refreshToken,
	}, nil
}
