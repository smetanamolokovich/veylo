package organization

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	domainorg "github.com/smetanamolokovich/veylo/internal/domain/organization"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
)

// CreateOrganizationUseCase is onboarding step 2: creates an organization and
// default workflow, links the authenticated user as admin, and issues a new JWT
// with the org_id populated.
type CreateOrganizationUseCase struct {
	orgRepo      domainorg.Repository
	workflowRepo workflow.Repository
	userRepo     user.Repository
	jwtManager   JWTManager
}

type CreateOrganizationRequest struct {
	// UserID is taken from the JWT claims (authenticated user without org yet).
	UserID   string
	OrgName  string
	Vertical string
}

type CreateOrganizationResponse struct {
	OrganizationID string `json:"organization_id"`
	AccessToken    string `json:"access_token"`
}

func NewCreateOrganizationUseCase(
	orgRepo domainorg.Repository,
	workflowRepo workflow.Repository,
	userRepo user.Repository,
	jwtManager JWTManager,
) *CreateOrganizationUseCase {
	return &CreateOrganizationUseCase{
		orgRepo:      orgRepo,
		workflowRepo: workflowRepo,
		userRepo:     userRepo,
		jwtManager:   jwtManager,
	}
}

func (uc *CreateOrganizationUseCase) Execute(ctx context.Context, req CreateOrganizationRequest) (*CreateOrganizationResponse, error) {
	vertical := domainorg.Vertical(req.Vertical)
	if !vertical.IsValid() {
		return nil, fmt.Errorf("CreateOrganizationUseCase.Execute: invalid vertical %q", req.Vertical)
	}

	orgID := ulid.Make().String()
	org, err := domainorg.NewOrganization(orgID, req.OrgName, vertical)
	if err != nil {
		return nil, fmt.Errorf("CreateOrganizationUseCase.Execute: %w", err)
	}
	if err := uc.orgRepo.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("CreateOrganizationUseCase.Execute: save org: %w", err)
	}

	var wf *workflow.Workflow
	switch vertical {
	case domainorg.VerticalVehicle:
		wf = workflow.DefaultVehicleWorkflow(ulid.Make().String(), orgID)
	default:
		wf, err = workflow.NewWorkflow(ulid.Make().String(), orgID)
		if err != nil {
			return nil, fmt.Errorf("CreateOrganizationUseCase.Execute: create workflow: %w", err)
		}
	}
	if err := uc.workflowRepo.Save(ctx, wf); err != nil {
		return nil, fmt.Errorf("CreateOrganizationUseCase.Execute: save workflow: %w", err)
	}

	// Link the user to the new organization.
	u, err := uc.userRepo.FindByIDOnly(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("CreateOrganizationUseCase.Execute: find user: %w", err)
	}
	u.SetOrganizationID(orgID)
	if err := uc.userRepo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("CreateOrganizationUseCase.Execute: update user org: %w", err)
	}

	// Issue a new JWT with org_id populated.
	accessToken, err := uc.jwtManager.Generate(req.UserID, orgID, string(user.RoleAdmin))
	if err != nil {
		return nil, fmt.Errorf("CreateOrganizationUseCase.Execute: generate token: %w", err)
	}

	return &CreateOrganizationResponse{
		OrganizationID: orgID,
		AccessToken:    accessToken,
	}, nil
}
