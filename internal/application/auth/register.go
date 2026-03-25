package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/smetanamolokovich/veylo/internal/domain/refreshtoken"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

// RegisterUseCase is onboarding step 1: creates a user without an organization.
// The returned JWT has an empty org_id claim. The client must then call
// POST /api/v1/organizations to complete registration.
type RegisterUseCase struct {
	userRepo         user.Repository
	refreshTokenRepo refreshtoken.Repository
	pwdHasher        PasswordHasher
	jwtManager       JWTManager
}

type RegisterRequest struct {
	Email    string
	Password string
	FullName string
}

type RegisterResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewRegisterUseCase(
	userRepo user.Repository,
	refreshTokenRepo refreshtoken.Repository,
	pwdHasher PasswordHasher,
	jwtManager JWTManager,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		pwdHasher:        pwdHasher,
		jwtManager:       jwtManager,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	existing, err := uc.userRepo.FindByEmailNoOrg(ctx, req.Email)
	if err != nil && !errors.Is(err, user.ErrNotFound) {
		return nil, fmt.Errorf("RegisterUseCase.Execute: %w", err)
	}
	if existing != nil {
		return nil, user.ErrAlreadyExists
	}

	hash, err := uc.pwdHasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("RegisterUseCase.Execute: hash password: %w", err)
	}

	userID := ulid.Make().String()
	newUser, err := user.NewUserWithoutOrg(userID, req.Email, hash, req.FullName)
	if err != nil {
		return nil, fmt.Errorf("RegisterUseCase.Execute: create user: %w", err)
	}

	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("RegisterUseCase.Execute: save user: %w", err)
	}

	// Issue access token with empty org_id — user has no org yet.
	accessToken, err := uc.jwtManager.Generate(userID, "", string(user.RoleAdmin))
	if err != nil {
		return nil, fmt.Errorf("RegisterUseCase.Execute: generate access token: %w", err)
	}

	rawRefresh, err := uc.jwtManager.GenerateRefresh()
	if err != nil {
		return nil, fmt.Errorf("RegisterUseCase.Execute: generate refresh token: %w", err)
	}
	hashedRefresh, err := uc.pwdHasher.Hash(rawRefresh)
	if err != nil {
		return nil, fmt.Errorf("RegisterUseCase.Execute: hash refresh token: %w", err)
	}

	rt, err := refreshtoken.NewRefreshToken(
		ulid.Make().String(),
		userID,
		"", // no org yet
		hashedRefresh,
		time.Now().UTC().Add(7*24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("RegisterUseCase.Execute: create refresh token: %w", err)
	}
	if err := uc.refreshTokenRepo.Save(ctx, rt); err != nil {
		return nil, fmt.Errorf("RegisterUseCase.Execute: save refresh token: %w", err)
	}

	return &RegisterResponse{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}
