package category

import (
	"context"

	"github.com/google/uuid"
)

// CategoryRepository defines the interface for accessing categories.
type CategoryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Category, error)
}

// In a real app implementation, this would be a Postgres repository.
