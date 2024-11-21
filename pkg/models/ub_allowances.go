package models

import "github.com/gofrs/uuid"

// UBAllowances represents the UB weight allowance for a branch, grade, order type, dependents authorized, and accompanied tour variables on an order
type UBAllowances struct {
	ID              uuid.UUID `db:"id" rw:"r"`
	BranchOfService *string   `db:"branch" rw:"r"`
	OrderPayGrade   *string   `db:"grade" rw:"r"`
	OrdersType      *string   `db:"orders_type" rw:"r"`
	HasDependents   *bool     `db:"dependents_authorized" rw:"r"`
	AccompaniedTour *bool     `db:"accompanied_tour" rw:"r"`
	UBAllowance     *int      `db:"ub_weight_allowance" rw:"r"`
}

// TableName overrides the table name used by Pop.
func (m UBAllowances) TableName() string {
	return "ub_allowances"
}
