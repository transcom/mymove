package models

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"
	"path"
)

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

// DeleteUpload deletes an upload from the database
func DeleteUpload(db *pop.Connection, upload *Upload) error {
	return db.Destroy(upload)
}

func (pdb *popDocumentDB) FetchUpload(id uuid.UUID) (*Upload, error) {
	var upload Upload
	err := pdb.db.Q().Eager().Find(&upload, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	return &upload, nil
}
