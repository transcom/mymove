package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
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
	S3ID        uuid.UUID `db:"s3_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// String is not required by pop and may be deleted
func (u Upload) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Uploads is not required by pop and may be deleted
type Uploads []Upload

// String is not required by pop and may be deleted
func (u Uploads) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

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

// BeforeCreate sets an Upload's S3ID before it is saved the first time
func (u *Upload) BeforeCreate(tx *pop.Connection) error {
	id, err := uuid.NewV4()
	if err != nil {
		return errors.WithStack(err)
	}
	u.S3ID = id
	return nil
}
