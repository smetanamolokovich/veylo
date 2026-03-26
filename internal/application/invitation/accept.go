package invitation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	doaminvitation "github.com/smetanamolokovich/veylo/internal/domain/invitation"
	"github.com/smetanamolokovich/veylo/internal/domain/refreshtoken"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

type AcceptInvitationUseCase struct {
	invitationRepo   doaminvitation.Repository
	userRepo         user.Repository
	refreshTokenRepo refreshtoken.Repository
	pwdHasher        PasswordHasher
	jwtManager       JWTManager
}

type AcceptInvitationRequest struct {
	Token    string
	FullName string
	Password string
}

type AcceptInvitationResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAcceptInvitationUseCase(
	invitationRepo doaminvitation.Repository,
	userRepo user.Repository,
	refreshTokenRepo refreshtoken.Repository,
	pwdHasher PasswordHasher,
	jwtManager JWTManager,
) *AcceptInvitationUseCase {
	return &AcceptInvitationUseCase{
		invitationRepo:   invitationRepo,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		pwdHasher:        pwdHasher,
		jwtManager:       jwtManager,
	}
}

func (uc *AcceptInvitationUseCase) Execute(ctx context.Context, req AcceptInvitationRequest) (*AcceptInvitationResponse, error) {
	inv, err := uc.invitationRepo.FindByToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, doaminvitation.ErrNotFound) {
			return nil, doaminvitation.ErrNotFound
		}
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: find invitation: %w", err)
	}

	if err := inv.Accept(); err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: %w", err)
	}

	// Check if this email is already registered in any org
	existing, err := uc.userRepo.FindByEmailNoOrg(ctx, inv.Email())
	if err != nil && !errors.Is(err, user.ErrNotFound) {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: check existing user: %w", err)
	}
	if existing != nil {
		return nil, user.ErrAlreadyExists
	}

	hash, err := uc.pwdHasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: hash password: %w", err)
	}

	userID := ulid.Make().String()
	newUser, err := user.NewUser(userID, inv.OrganizationID(), inv.Email(), hash, req.FullName, inv.Role())
	if err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: create user: %w", err)
	}

	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: save user: %w", err)
	}

	if err := uc.invitationRepo.Save(ctx, inv); err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: save invitation: %w", err)
	}

	accessToken, err := uc.jwtManager.Generate(userID, inv.OrganizationID(), string(inv.Role()))
	if err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: generate access token: %w", err)
	}

	rawRefresh, err := uc.jwtManager.GenerateRefresh()
	if err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: generate refresh token: %w", err)
	}
	hashedRefresh, err := uc.pwdHasher.Hash(rawRefresh)
	if err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: hash refresh token: %w", err)
	}

	rt, err := refreshtoken.NewRefreshToken(
		ulid.Make().String(),
		userID,
		inv.OrganizationID(),
		hashedRefresh,
		time.Now().UTC().Add(7*24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: create refresh token: %w", err)
	}
	if err := uc.refreshTokenRepo.Save(ctx, rt); err != nil {
		return nil, fmt.Errorf("AcceptInvitationUseCase.Execute: save refresh token: %w", err)
	}

	return &AcceptInvitationResponse{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}
