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
	ID              uuid.UUID `db:"id"`
	UploaderID      uuid.UUID `db:"uploader_id"`
	ServiceMemberID uuid.UUID `db:"service_member_id"`
	Name            string    `db:"name"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// String is not required by pop and may be deleted
func (d Document) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Documents is not required by pop and may be deleted
type Documents []Document

// String is not required by pop and may be deleted
func (d Documents) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (d *Document) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: d.UploaderID, Name: "UploaderID"},
		&validators.UUIDIsPresent{Field: d.ServiceMemberID, Name: "ServiceMemberID"},
	), nil
}

// ValidateDocumentAccess validates that a user has access to document
func ValidateDocumentAccess(db *pop.Connection, userID uuid.UUID, documentID uuid.UUID) (bool, bool) {
	exists := false
	userHasAccess := false
	var document Document
	docErr := db.Find(&document, documentID)
	if docErr == nil {
		exists = true
		userHasAccess = ValidateServiceMemberAccess(db, userID, document.ServiceMemberID)
	}
	return exists, userHasAccess
}
