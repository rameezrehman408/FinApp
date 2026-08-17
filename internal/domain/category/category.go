package category

import (
	"time"

	"github.com/google/uuid"
)

type BudgetPeriod string

const (
	Weekly  BudgetPeriod = "WEEKLY"
	Monthly BudgetPeriod = "MONTHLY"
	Yearly  BudgetPeriod = "YEARLY"
)

type CategoryType string

const (
	Income  CategoryType = "INCOME"
	Expense CategoryType = "EXPENSE"
)

type Category struct {
	ID           uuid.UUID    `json:"id"`
	HouseholdID  uuid.UUID    `json:"household_id"`
	Name         string       `json:"name"`
	Type         CategoryType `json:"type"` // INCOME, EXPENSE, etc.
	BudgetLimit  float64      `json:"budget_limit"`
	BudgetPeriod BudgetPeriod `json:"budget_period"`
	ParentID     *uuid.UUID   `json:"parent_id,omitempty"`
	Description  string       `json:"description"`
	Icon         string       `json:"icon"`
	Color        string       `json:"color"`
	IsSystem     bool         `json:"is_system"`
	Version      int          `json:"version"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func NewCategory(id, householdID uuid.UUID, name string, cType CategoryType) *Category {
	now := time.Now()
	return &Category{
		ID:          id,
		HouseholdID: householdID,
		Name:        name,
		Type:        cType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

