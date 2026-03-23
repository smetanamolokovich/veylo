package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/smetanamolokovich/veylo/internal/domain/organization"
	"github.com/smetanamolokovich/veylo/internal/domain/refreshtoken"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
)

type SignupUseCase struct {
	orgRepo          organization.Repository
	workflowRepo     workflow.Repository
	userRepo         user.Repository
	refreshTokenRepo refreshtoken.Repository
	pwdHasher        PasswordHasher
	jwtManager       JWTManager
}

type SignupRequest struct {
	OrgName   string
	Vertical  string
	Email     string
	Password  string
	FullName  string
}

type SignupResponse struct {
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
}

func NewSignupUseCase(
	orgRepo organization.Repository,
	workflowRepo workflow.Repository,
	userRepo user.Repository,
	refreshTokenRepo refreshtoken.Repository,
	pwdHasher PasswordHasher,
	jwtManager JWTManager,
) *SignupUseCase {
	return &SignupUseCase{
		orgRepo:          orgRepo,
		workflowRepo:     workflowRepo,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		pwdHasher:        pwdHasher,
		jwtManager:       jwtManager,
	}
}

func (uc *SignupUseCase) Execute(ctx context.Context, req SignupRequest) (*SignupResponse, error) {
	vertical := organization.Vertical(req.Vertical)
	if !vertical.IsValid() {
		return nil, fmt.Errorf("auth.Signup: invalid vertical %q", req.Vertical)
	}

	orgID := ulid.Make().String()
	org, err := organization.NewOrganization(orgID, req.OrgName, vertical)
	if err != nil {
		return nil, fmt.Errorf("auth.Signup: %w", err)
	}
	if err := uc.orgRepo.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("auth.Signup: save org: %w", err)
	}

	var wf *workflow.Workflow
	switch vertical {
	case organization.VerticalVehicle:
		wf = workflow.DefaultVehicleWorkflow(ulid.Make().String(), orgID)
	default:
		wf, err = workflow.NewWorkflow(ulid.Make().String(), orgID)
		if err != nil {
			return nil, fmt.Errorf("auth.Signup: create workflow: %w", err)
		}
	}
	if err := uc.workflowRepo.Save(ctx, wf); err != nil {
		return nil, fmt.Errorf("auth.Signup: save workflow: %w", err)
	}

	hash, err := uc.pwdHasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("auth.Signup: hash password: %w", err)
	}
	userID := ulid.Make().String()
	admin, err := user.NewUser(userID, orgID, req.Email, hash, req.FullName, user.RoleAdmin)
	if err != nil {
		return nil, fmt.Errorf("auth.Signup: create user: %w", err)
	}
	if err := uc.userRepo.Save(ctx, admin); err != nil {
		return nil, fmt.Errorf("auth.Signup: save user: %w", err)
	}

	accessToken, err := uc.jwtManager.Generate(userID, orgID, string(user.RoleAdmin))
	if err != nil {
		return nil, fmt.Errorf("auth.Signup: generate access token: %w", err)
	}

	rawRefresh, err := uc.jwtManager.GenerateRefresh()
	if err != nil {
		return nil, fmt.Errorf("auth.Signup: generate refresh token: %w", err)
	}
	hashedRefresh, err := uc.pwdHasher.Hash(rawRefresh)
	if err != nil {
		return nil, fmt.Errorf("auth.Signup: hash refresh token: %w", err)
	}

	rt, err := refreshtoken.NewRefreshToken(
		ulid.Make().String(),
		userID,
		orgID,
		hashedRefresh,
		time.Now().UTC().Add(7*24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("auth.Signup: create refresh token: %w", err)
	}
	if err := uc.refreshTokenRepo.Save(ctx, rt); err != nil {
		return nil, fmt.Errorf("auth.Signup: save refresh token: %w", err)
	}

	return &SignupResponse{
		OrganizationID: orgID,
		UserID:         userID,
		AccessToken:    accessToken,
		RefreshToken:   rawRefresh,
	}, nil
}
