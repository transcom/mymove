package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// A Document represents a physical artifact such as a multipage form that was
// filled out by hand. A Document can have many associated Uploads, which allows
// for handling multiple files that belong to the same document.
type Document struct {
	ID         uuid.UUID `db:"id"`
	UploaderID uuid.UUID `db:"uploader_id"`
	MoveID     uuid.UUID `db:"move_id"`
	Name       string    `db:"name"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// Documents is not required by pop and may be deleted
type Documents []Document

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (d *Document) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: d.UploaderID, Name: "UploaderID"},
		&validators.UUIDIsPresent{Field: d.MoveID, Name: "MoveID"},
	), nil
}

// ValidateDocumentOwnership validates that a user owns the move that contains a document and that move and document both exist
func ValidateDocumentOwnership(db *pop.Connection, userID uuid.UUID, moveID uuid.UUID, documentID uuid.UUID) (bool, bool) {
	exists := false
	userOwns := false
	var move Move
	var document Document
	docErr := db.Find(&document, documentID)
	moveErr := db.Find(&move, moveID)
	if docErr == nil && moveErr == nil {
		exists = true
		// TODO: Handle case where more than one user is authorized to modify move
		if uuid.Equal(move.UserID, userID) && uuid.Equal(document.MoveID, moveID) {
			userOwns = true
		}
	}
	return exists, userOwns
}
