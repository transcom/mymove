package models

import (
	"path"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth"
)

// An Upload represents an uploaded file, such as an image or PDF.
type Upload struct {
	ID          uuid.UUID  `db:"id"`
	DocumentID  *uuid.UUID `db:"document_id"`
	Document    Document   `belongs_to:"documents"`
	UploaderID  uuid.UUID  `db:"uploader_id"`
	Filename    string     `db:"filename"`
	Bytes       int64      `db:"bytes"`
	ContentType string     `db:"content_type"`
	Checksum    string     `db:"checksum"`
	StorageKey  string     `db:"storage_key"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// Uploads is not required by pop and may be deleted
type Uploads []Upload

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *Upload) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: u.UploaderID, Name: "UploaderID"},
		&validators.StringIsPresent{Field: u.Filename, Name: "Filename"},
		&Int64IsPresent{Field: u.Bytes, Name: "Bytes"},
		NewAllowedFileTypeValidator(u.ContentType, "ContentType"),
		&validators.StringIsPresent{Field: u.Checksum, Name: "Checksum"},
	), nil
}

// BeforeCreate populates the StorageKey on a newly created Upload
func (u *Upload) BeforeCreate(tx *pop.Connection) error {
	// Populate ID if not exists
	if uuid.Equal(u.ID, uuid.UUID{}) {
		u.ID = uuid.Must(uuid.NewV4())
	}

	if u.StorageKey == "" {
		u.StorageKey = path.Join("user", u.UploaderID.String(), "uploads", u.ID.String())
	}

	return nil
}

// FetchUpload returns an Upload if the user has access to that upload
func FetchUpload(db *pop.Connection, session *auth.Session, id uuid.UUID) (Upload, error) {
	var upload Upload
	err := db.Q().Eager().Find(&upload, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return Upload{}, errors.Wrap(ErrFetchNotFound, "error fetching upload")
		}
		// Otherwise, it's an unexpected err so we return that.
		return Upload{}, err
	}

	// If there's a document, check permissions. Otherwise user must
	// have been the uploader
	if upload.DocumentID != nil {
		_, docErr := FetchDocument(db, session, *upload.DocumentID)
		if docErr != nil {
			return Upload{}, docErr
		}
	} else if upload.UploaderID != session.UserID {
		return Upload{}, errors.Wrap(ErrFetchNotFound, "user ID doesn't match uploader ID")
	}
	return upload, nil
}

// DeleteUpload deletes an upload from the database
func DeleteUpload(db *pop.Connection, upload *Upload) error {
	return db.Destroy(upload)
}
