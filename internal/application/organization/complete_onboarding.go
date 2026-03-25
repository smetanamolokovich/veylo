package organization

import (
	"context"
	"fmt"

	domainorg "github.com/smetanamolokovich/veylo/internal/domain/organization"
)

type CompleteOnboardingUseCase struct {
	orgRepo domainorg.Repository
}

type CompleteOnboardingRequest struct {
	OrganizationID string
}

type CompleteOnboardingResponse struct {
	OrganizationID      string `json:"organization_id"`
	OnboardingCompleted bool   `json:"onboarding_completed"`
}

func NewCompleteOnboardingUseCase(orgRepo domainorg.Repository) *CompleteOnboardingUseCase {
	return &CompleteOnboardingUseCase{orgRepo: orgRepo}
}

func (uc *CompleteOnboardingUseCase) Execute(ctx context.Context, req CompleteOnboardingRequest) (*CompleteOnboardingResponse, error) {
	org, err := uc.orgRepo.FindByID(ctx, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("CompleteOnboardingUseCase.Execute: find org: %w", err)
	}

	org.CompleteOnboarding()

	if err := uc.orgRepo.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("CompleteOnboardingUseCase.Execute: save org: %w", err)
	}

	return &CompleteOnboardingResponse{
		OrganizationID:      org.ID(),
		OnboardingCompleted: true,
	}, nil
}
