package budget

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBudgetService_GetUtilization(t *testing.T) {
	service := NewBudgetService()
	ctx := context.Background()

	tests := []struct {
		name        string
		spent       float64
		limit       float64
		expected    float64
		wantErr     bool
	}{
		{
			name:     "50 percent utilization",
			spent:    50.0,
			limit:    100.0,
			expected: 50.0,
			wantErr:  false,
		},
		{
			name:     "100 percent utilization",
			spent:    100.0,
			limit:    100.0,
			expected: 100.0,
			wantErr:  false,
		},
		{
			name:     "Over budget",
			spent:    120.0,
			limit:    100.0,
			expected: 120.0,
			wantErr:  false,
		},
		{
			name:     "Zero limit error",
			spent:    50.0,
			limit:    0,
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetUtilization(ctx, tt.spent, tt.limit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
