package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// CustomerSupportRemark struct represents the shape of an office move remark
type CustomerSupportRemark struct {
	ID           uuid.UUID  `db:"id"`
	Content      string     `db:"content"`
	OfficeUser   OfficeUser `belongs_to:"office_users" fk_id:"office_user_id"`
	OfficeUserID uuid.UUID  `db:"office_user_id"`
	Move         Move       `belongs_to:"moves" fk_id:"move_id"`
	MoveID       uuid.UUID  `db:"move_id"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

type CustomerSupportRemarks []CustomerSupportRemark

// TableName overrides the table name used by Pop.
func (o CustomerSupportRemark) TableName() string {
	return "customer_support_remarks"
}
