package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/smetanamolokovich/veylo/internal/domain/refreshtoken"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

type RefreshTokenUseCase struct {
	refreshTokenRepo refreshtoken.Repository
	userRepo         user.Repository
	jwtManager       JWTManager
	hasher           PasswordHasher
}

func NewRefreshTokenUseCase(refreshTokenRepo refreshtoken.Repository, userRepo user.Repository, jwtManager JWTManager, hasher PasswordHasher) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		refreshTokenRepo: refreshTokenRepo,
		userRepo:         userRepo,
		jwtManager:       jwtManager,
		hasher:           hasher,
	}
}

type RefreshRequest struct {
	RefreshToken   string
	UserID         string
	OrganizationID string
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, req RefreshRequest) (*RefreshResponse, error) {
	rt, err := uc.refreshTokenRepo.FindByUserID(ctx, req.UserID, req.OrganizationID)
	if err != nil {
		return nil, err
	}

	if rt.IsExpired() {
		return nil, refreshtoken.ErrExpiredRefreshToken
	}

	if !uc.hasher.Compare(req.RefreshToken, rt.TokenHash()) {
		return nil, refreshtoken.ErrInvalidRefreshToken
	}

	existing, err := uc.userRepo.FindByID(ctx, req.UserID, req.OrganizationID)
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.jwtManager.Generate(existing.ID(), existing.OrganizationID(), string(existing.Role()))
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: generate access token: %w", err)
	}

	newRawToken, err := uc.jwtManager.GenerateRefresh()
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: generate refresh token: %w", err)
	}

	hashedToken, err := uc.hasher.Hash(newRawToken)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: hash refresh token: %w", err)
	}

	newRT, err := refreshtoken.NewRefreshToken(
		ulid.Make().String(),
		req.UserID,
		req.OrganizationID,
		hashedToken,
		time.Now().UTC().Add(7*24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: create refresh token: %w", err)
	}

	if err := uc.refreshTokenRepo.DeleteByUserID(ctx, req.UserID, req.OrganizationID); err != nil {
		return nil, fmt.Errorf("auth.Refresh: delete old token: %w", err)
	}

	if err := uc.refreshTokenRepo.Save(ctx, newRT); err != nil {
		return nil, fmt.Errorf("auth.Refresh: save new token: %w", err)
	}

	return &RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRawToken,
	}, nil
}
