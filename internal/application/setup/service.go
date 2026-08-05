package setup

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/seisram/FinApp/internal/domain/household"
	"github.com/seisram/FinApp/internal/domain/user"
)

// SetupService handles the initial user and household creation.
type SetupService struct {
	// In a real app, we'd inject repositories here.
	// For this demo/foundation, we'll use simplified in-memory or mock logic.
}

func NewSetupService() *SetupService {
	return &SetupService{}
}

// Execute performs the multi-tenant setup: Create User -> Create Household.
func (s *SetupService) Execute(ctx context.Context, params user.CreateUserParams) (*user.User, error) {
	// 1. Logic to create user in DB would go here.
	newUser := &user.User{
		ID:        uuid.New(),
		Email:     params.Email,
		FullName:  params.FullName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 2. Logic to create household linked to the user.
	mode := household.HouseholdMode(params.HouseholdMode)
	if mode != household.SingleMode && mode != household.HouseholdModeType {
		return nil, fmt.Errorf("invalid household mode: %s", params.HouseholdMode)
	}

	// In reality, we'd use a transaction to ensure both or neither are created.
	newHousehold := &household.Household{
		ID:        uuid.New(),
		Name:      fmt.Sprintf("%s's Household", newUser.FullName),
		OwnerID:   newUser.ID,
		Mode:      mode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// For the sake of this-step-foundation, we return them so tests can verify logic flow.
	_ = newHousehold 

	return newUser, nil
}
