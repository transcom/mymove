package models

import "github.com/gofrs/uuid"

// MoveToGBLOC represents the view that associates each move ID with a GBLOC based on the postal code of its first shipment.
// This view is used to encapsulate query logic that was impossible to express with Pop.
// It will be used for the TOO and TIO queues, but not for Services Counseling.
type MoveToGBLOC struct {
	MoveID uuid.UUID `db:"move_id"`
	Move   Move      `belongs_to:"moves" fk_id:"move_id"`
	GBLOC  *string   `db:"gbloc"`
}

// MoveToGBLOCs is not required by pop and may be deleted
type MoveToGBLOCs []MoveToGBLOC

// TableName overrides the table name used by Pop.
func (m MoveToGBLOC) TableName() string {
	return "move_to_gbloc"
}
