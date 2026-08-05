package household

import (
	"time"

	"github.com/google/uuid"
)

type HouseholdMode string

const (
	SingleMode    HouseholdMode = "SINGLE"
	HouseholdModeType HouseholdMode = "HOUSEHOLD"
)

type Household struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	OwnerID   uuid.UUID     `json:"owner_id"`
	Mode      HouseholdMode `json:"mode"`
	CreatedAt time.Time     `json("created_at")`
	UpdatedAt time.Time     `json:"updated_at"`
}

type CreateHouseholdParams struct {
	Name    string        `json:"name"`
	OwnerID uuid.UUID     `json:"owner_id"`
	Mode    HouseholdMode `json:"mode"`
}
