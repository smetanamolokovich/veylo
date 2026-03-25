package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/smetanamolokovich/veylo/internal/domain/refreshtoken"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

type LoginUseCase struct {
	userRepo         user.Repository
	refreshTokenRepo refreshtoken.Repository
	pwdHasher        PasswordHasher
	jwtManager       JWTManager
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
}

func NewLoginUseCase(userRepo user.Repository, refreshTokenRepo refreshtoken.Repository, pwdHasher PasswordHasher, jwtManager JWTManager) *LoginUseCase {
	return &LoginUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		pwdHasher:        pwdHasher,
		jwtManager:       jwtManager,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	existing, err := uc.userRepo.FindByEmailNoOrg(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing.Status() == user.StatusBlocked {
		return nil, user.ErrBlocked
	}
	if ok := uc.pwdHasher.Compare(req.Password, existing.PasswordHash()); !ok {
		return nil, user.ErrInvalidCredentials
	}

	accessToken, err := uc.jwtManager.Generate(existing.ID(), existing.OrganizationID(), string(existing.Role()))
	if err != nil {
		return nil, fmt.Errorf("auth.Login: generate access token: %w", err)
	}

	rawRefreshToken, err := uc.jwtManager.GenerateRefresh()
	if err != nil {
		return nil, fmt.Errorf("auth.Login: generate refresh token: %w", err)
	}

	hashedRefreshToken, err := uc.pwdHasher.Hash(rawRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("auth.Login: hash refresh token: %w", err)
	}

	rt, err := refreshtoken.NewRefreshToken(
		ulid.Make().String(),
		existing.ID(),
		existing.OrganizationID(),
		hashedRefreshToken,
		time.Now().UTC().Add(7*24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("auth.Login: create refresh token: %w", err)
	}

	if err := uc.refreshTokenRepo.DeleteByUserID(ctx, existing.ID(), existing.OrganizationID()); err != nil {
		return nil, fmt.Errorf("auth.Login: delete old refresh token: %w", err)
	}

	if err := uc.refreshTokenRepo.Save(ctx, rt); err != nil {
		return nil, fmt.Errorf("auth.Login: save refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		UserID:       existing.ID(),
		Role:         string(existing.Role()),
	}, nil
}
