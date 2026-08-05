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

type Category struct {
	ID            uuid.UUID    `json:"id"`
	HouseholdID   uuid.UUID    `json:"household_id"`
	Name          string       `json:"name"`
	Type          string       `json:"type"` // INCOME, EXPENSE, etc.
	BudgetLimit   float64      `json:"budget_limit"`
	BudgetPeriod  BudgetPeriod `json:"budget_period"`
	ParentID      *uuid.UUID   `json:"parent_id,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
