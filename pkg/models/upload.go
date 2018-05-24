package models

import (
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
	Document    Document  `belongs_to:"documents"`
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
		NewAllowedFileTypeValidator(u.ContentType, "ContentType"),
		&validators.StringIsPresent{Field: u.Checksum, Name: "Checksum"},
	), nil
}

// FetchUpload returns an Upload if the user has access to that upload
func FetchUpload(db *pop.Connection, user User, reqApp string, id uuid.UUID) (Upload, error) {
	var upload Upload
	err := db.Q().Eager().Find(&upload, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return Upload{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Upload{}, err
	}

	_, docErr := FetchDocument(db, user, reqApp, upload.DocumentID)
	if docErr != nil {
		return Upload{}, docErr
	}
	return upload, nil
}

// DeleteUpload deletes an upload from the database
func DeleteUpload(db *pop.Connection, upload *Upload) error {
	return db.Destroy(upload)
}
