package models

import "github.com/gofrs/uuid"

// UBAllowances represents the UB weight allowance for a branch, grade, order type, dependents authorized, and accompanied tour variables on an order
type UBAllowances struct {
	ID              uuid.UUID `db:"id"`
	BranchOfService string    `db:"branch"`
	OrderPayGrade   string    `db:"grade"`
	OrdersType      string    `db:"orders_type"`
	HasDependents   bool      `db:"dependents_authorized"`
	AccompaniedTour bool      `db:"accompanied_tour"`
	UBAllowance     int       `db:"ub_weight_allowance"`
}

// TableName overrides the table name used by Pop.
func (m UBAllowances) TableName() string {
	return "ub_allowances"
}
