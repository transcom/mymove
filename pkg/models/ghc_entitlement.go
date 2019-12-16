package models

import (
	"time"

	"github.com/gofrs/uuid"
)

//TODO not sure if this should exist but current Entitlement model is not db table and differs a bit in structure
//TODO so just going to break out into separate object for now.
type GHCEntitlement struct {
	ID                    uuid.UUID      `db:"id"`
	DependentsAuthorized  bool           `db:"dependents_authorized"`
	TotalDependents       int            `db:"total_dependents"`
	NonTemporaryStorage   bool           `db:"non_temporary_storage"`
	PrivatelyOwnedVehicle bool           `db:"privately_owned_vehicle"`
	ProGearWeight         int            `db:"pro_gear_weight"`
	ProGearWeightSpouse   int            `db:"pro_gear_weight_spouse"`
	StorageInTransit      int            `db:"storage_in_transit"`
	CreatedAt             time.Time      `db:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at"`
	MoveTaskOrder         *MoveTaskOrder `belongs_to:"move_task_orders"`
	MoveTaskOrderID       uuid.UUID      `db:"move_task_order_id"`
}

// TableName overrides the table name used by Pop.
func (ghce GHCEntitlement) TableName() string {
	return "ghc_entitlements"
}
