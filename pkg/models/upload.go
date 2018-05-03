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
		&AllowedFiletype{Field: u.ContentType, Name: "ContentType"},
		&validators.StringIsPresent{Field: u.Checksum, Name: "Checksum"},
	), nil
}

// FetchUpload returns an Upload if the user has access to that upload
func FetchUpload(db *pop.Connection, user User, id uuid.UUID) (Upload, error) {
	var upload Upload
	err := db.Q().Eager().Find(&upload, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return Upload{}, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return Upload{}, err
	}

	_, docErr := FetchDocument(db, user, upload.DocumentID)
	if docErr != nil {
		return Upload{}, docErr
	}
	return upload, nil
}
