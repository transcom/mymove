package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Entitlement is an object representing entitlements for move orders
type Entitlement struct {
	ID                    uuid.UUID `db:"id"`
	DependentsAuthorized  *bool     `db:"dependents_authorized"`
	TotalDependents       *int      `db:"total_dependents"`
	NonTemporaryStorage   *bool     `db:"non_temporary_storage"`
	PrivatelyOwnedVehicle *bool     `db:"privately_owned_vehicle"`
	ProGearWeight         *int      `db:"pro_gear_weight"`
	ProGearWeightSpouse   *int      `db:"pro_gear_weight_spouse"`
	StorageInTransit      *int      `db:"storage_in_transit"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
}
