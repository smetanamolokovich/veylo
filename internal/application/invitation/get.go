package invitation

import (
	"context"
	"errors"
	"fmt"

	doaminvitation "github.com/smetanamolokovich/veylo/internal/domain/invitation"
	"github.com/smetanamolokovich/veylo/internal/domain/organization"
)

type GetInvitationUseCase struct {
	invitationRepo doaminvitation.Repository
	orgRepo        organization.Repository
}

type GetInvitationRequest struct {
	Token string
}

type GetInvitationResponse struct {
	Email            string `json:"email"`
	OrganizationName string `json:"organization_name"`
	Role             string `json:"role"`
	ExpiresAt        string `json:"expires_at"`
	IsExpired        bool   `json:"is_expired"`
}

func NewGetInvitationUseCase(
	invitationRepo doaminvitation.Repository,
	orgRepo organization.Repository,
) *GetInvitationUseCase {
	return &GetInvitationUseCase{
		invitationRepo: invitationRepo,
		orgRepo:        orgRepo,
	}
}

func (uc *GetInvitationUseCase) Execute(ctx context.Context, req GetInvitationRequest) (*GetInvitationResponse, error) {
	inv, err := uc.invitationRepo.FindByToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, doaminvitation.ErrNotFound) {
			return nil, doaminvitation.ErrNotFound
		}
		return nil, fmt.Errorf("GetInvitationUseCase.Execute: find invitation: %w", err)
	}

	org, err := uc.orgRepo.FindByID(ctx, inv.OrganizationID())
	if err != nil {
		return nil, fmt.Errorf("GetInvitationUseCase.Execute: find org: %w", err)
	}

	return &GetInvitationResponse{
		Email:            inv.Email(),
		OrganizationName: org.Name(),
		Role:             string(inv.Role()),
		ExpiresAt:        inv.ExpiresAt().Format("2006-01-02T15:04:05Z07:00"),
		IsExpired:        inv.IsExpired(),
	}, nil
}
