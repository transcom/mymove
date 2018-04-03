package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// An Upload represents an uploaded file, such as an image or PDF.
type Upload struct {
	ID          uuid.UUID `db:"id"`
	DocumentID  uuid.UUID `db:"document_id"`
	UploaderID  uuid.UUID `db:"uploader_id"`
	Filename    string    `db:"filename"`
	Bytes       int64     `db:"bytes"`
	ContentType string    `db:"content_type"`
	Checksum    string    `db:"checksum"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Uploads is not required by pop and may be deleted
type Uploads []Upload

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *Upload) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: u.DocumentID, Name: "DocumentID"},
		&validators.UUIDIsPresent{Field: u.UploaderID, Name: "UploaderID"},
		&validators.StringIsPresent{Field: u.Filename, Name: "Filename"},
		&Int64IsPresent{Field: u.Bytes, Name: "Bytes"},
		&validators.StringIsPresent{Field: u.ContentType, Name: "ContentType"},
		&validators.StringIsPresent{Field: u.Checksum, Name: "Checksum"},
	), nil
}
