package models

import "github.com/gofrs/uuid"

// OriginDutyLocationToGBLOC represents the view that associates each move ID with a GBLOC based on the postal code of its origin duty location.
// This view is used to encapsulate query logic that was impossible to express with Pop.
// It will be used for the for Services Counseling queue.
type OriginDutyLocationToGBLOC struct {
	ID     uuid.UUID `db:"id"`
	MoveID uuid.UUID `db:"move_id"`
	GBLOC  string    `db:"gbloc"`
}

// TableName overrides the table name used by Pop.
func (m OriginDutyLocationToGBLOC) TableName() string {
	return "origin_duty_location_to_gbloc"
}
