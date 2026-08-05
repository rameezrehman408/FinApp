package budget

import (
	"context"
	"errors"
)

// BudgetService calculates budget utilization for a household.
type BudgetService struct {
	// In real life, we'd inject a repository here.
}

func NewBudgetService() *BudgetService {
	return &BudgetService{}
}

// GetUtilization returns the usage percentage of a category budget.
func (s *BudgetService) GetUtilization(ctx context.Context, spent float64, limit float64) (float64, error) {
	if limit <= 0 {
		return 0, errors.New("budget limit must be greater than zero")
	}
	return (spent / limit) * 100, nil
}

// CheckBudgetOverrun returns true if the spent amount exceeds the limit.
func (s *BudgetService) CheckBudgetOverrun(ctx context.Context, spent float64, limit float64) (bool, error) {
	utilization, err := s.GetUtilization(ctx, spent, limit)
	if err != nil {
		return false, err
	}
	return utilization > 100.0, nil
}
