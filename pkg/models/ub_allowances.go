package models

// UBAllowances represents the UB weight allowance for a branch, grade, order type, dependenths authorized, and accompanied tour variables on an order
type UBAllowances struct {
	BranchOfService string `db:"branch"`
	OrderPayGrade   string `db:"grade"`
	OrdersType      string `db:"orders_type"`
	HasDependents   bool   `db:"dependents_authorized"`
	AccompaniedTour bool   `db:"accompanied_tour"`
	BaseUBAllowance int    `db:"ub_weight_allowance"`
}

// TableName overrides the table name used by Pop.
func (m UBAllowances) TableName() string {
	return "ub_allowances"
}
