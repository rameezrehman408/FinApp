package setup

import (
	"context"
	"testing"

	"github.com/seisram/FinApp/internal/domain/household"
	"github.com/seisram/FinApp/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestSetupService_Execute(t *testing.T) {
	service := NewSetupService()
	ctx := context.Background()

	tests := []struct {
		name         string
		params       user.CreateUserParams
		wantErr      bool
		expectedMode household.HouseholdMode
	}{
		{
			name: "Successful Single User Setup",
			params: user.CreateUserParams{
				Email:         "test@example.com",
				Password:      "password123",
				FullName:      "Test User",
				HouseholdMode: "SINGLE",
			},
			wantErr:      false,
			expectedMode: household.SingleMode,
		},
		{
			name: "Successful Household Setup",
			params: user.CreateUserParams{
				Email:         "family@example.com",
				Password:      "password123",
				FullName:      "The Family",
				HouseholdMode: "HOUSEHOLD",
			},
			wantErr:      false,
			expectedMode: household.HouseholdModeType,
		},
		{
			name: "Failure - Invalid Mode",
			params: user.CreateUserParams{
				Email:         "bad@example.com",
				Password:      "password123",
				FullName:      "Bad User",
				HouseholdMode: "INVALID_MODE",
			},
			wantErr:      true,
			expectedMode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := service.Execute(ctx, tt.params)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotUser)
				assert.Equal(t, tt.params.Email, gotUser.Email)
			}
		})
	}
}
