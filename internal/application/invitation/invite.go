package invitation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	doaminvitation "github.com/smetanamolokovich/veylo/internal/domain/invitation"
	"github.com/smetanamolokovich/veylo/internal/domain/organization"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

type InviteUserUseCase struct {
	invitationRepo doaminvitation.Repository
	orgRepo        organization.Repository
}

type InviteUserRequest struct {
	OrganizationID string
	InviterUserID  string
	Email          string
	Role           string
}

type InviteUserResponse struct {
	InvitationID string `json:"invitation_id"`
	InviteToken  string `json:"invite_token"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	ExpiresAt    string `json:"expires_at"`
}

func NewInviteUserUseCase(
	invitationRepo doaminvitation.Repository,
	orgRepo organization.Repository,
) *InviteUserUseCase {
	return &InviteUserUseCase{
		invitationRepo: invitationRepo,
		orgRepo:        orgRepo,
	}
}

func (uc *InviteUserUseCase) Execute(ctx context.Context, req InviteUserRequest) (*InviteUserResponse, error) {
	role := user.Role(req.Role)
	if role != user.RoleInspector && role != user.RoleEvaluator && role != user.RoleManager {
		return nil, fmt.Errorf("InviteUserUseCase.Execute: invalid role: %s", req.Role)
	}

	_, err := uc.orgRepo.FindByID(ctx, req.OrganizationID)
	if err != nil {
		if errors.Is(err, organization.ErrNotFound) {
			return nil, fmt.Errorf("InviteUserUseCase.Execute: %w", organization.ErrNotFound)
		}
		return nil, fmt.Errorf("InviteUserUseCase.Execute: find org: %w", err)
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("InviteUserUseCase.Execute: generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	id := ulid.Make().String()
	inv, err := doaminvitation.NewInvitation(id, req.OrganizationID, req.Email, role, token, req.InviterUserID)
	if err != nil {
		return nil, fmt.Errorf("InviteUserUseCase.Execute: create invitation: %w", err)
	}

	if err := uc.invitationRepo.Save(ctx, inv); err != nil {
		if errors.Is(err, doaminvitation.ErrDuplicate) {
			return nil, doaminvitation.ErrDuplicate
		}
		return nil, fmt.Errorf("InviteUserUseCase.Execute: save invitation: %w", err)
	}

	return &InviteUserResponse{
		InvitationID: inv.ID(),
		InviteToken:  token,
		Email:        inv.Email(),
		Role:         string(inv.Role()),
		ExpiresAt:    inv.ExpiresAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
