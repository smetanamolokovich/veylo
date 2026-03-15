package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

type RegisterUseCase struct {
	userRepo  user.Repository
	pwdHasher PasswordHasher
}

type RegisterRequest struct {
	Email          string
	Password       string
	OrganizationID string
	FullName       string
	Role           user.Role
}

type RegisterResponse struct {
	ID    string    `json:"id"`
	Email string    `json:"email"`
	Role  user.Role `json:"role"`
}

func NewRegisterUseCase(userRepo user.Repository, pwdHasher PasswordHasher) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:  userRepo,
		pwdHasher: pwdHasher,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	existing, err := uc.userRepo.FindByEmail(ctx, req.Email, req.OrganizationID)
	if err != nil && !errors.Is(err, user.ErrNotFound) {
		return nil, fmt.Errorf("auth.Register: %w", err)
	}
	if existing != nil {
		return nil, user.ErrAlreadyExists
	}

	hash, err := uc.pwdHasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	id := ulid.Make().String()

	newUser, err := user.NewUser(
		id,
		req.OrganizationID,
		req.Email,
		hash,
		req.FullName,
		req.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &RegisterResponse{
		ID:    newUser.ID(),
		Email: newUser.Email(),
		Role:  newUser.Role(),
	}, nil
}
