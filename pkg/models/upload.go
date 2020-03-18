package models

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/utilities"
)

// UploadType represents the type of upload this is, whether is it uploaded for a User or for the Prime
type UploadType string

const (
	// UploadTypeUSER string USER
	UploadTypeUSER UploadType = "USER"
	// UploadTypePRIME string PRIME
	UploadTypePRIME UploadType = "PRIME"
)

// An Upload represents an uploaded file, such as an image or PDF.
type Upload struct {
	ID          uuid.UUID  `db:"id"`
	Filename    string     `db:"filename"`
	Bytes       int64      `db:"bytes"`
	ContentType string     `db:"content_type"`
	Checksum    string     `db:"checksum"`
	StorageKey  string     `db:"storage_key"`
	UploadType  UploadType `db:"upload_type"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

// Uploads is not required by pop and may be deleted
type Uploads []Upload

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *Upload) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	vs = append(vs, &validators.StringInclusion{Field: string(u.UploadType), Name: "UploadType", List: []string{
		string(UploadTypeUSER),
		string(UploadTypePRIME),
	}})
	vs = append(vs,
		&validators.StringIsPresent{Field: u.Filename, Name: "Filename"},
		&Int64IsPresent{Field: u.Bytes, Name: "Bytes"},
		&validators.StringIsPresent{Field: u.ContentType, Name: "ContentType"},
		&validators.StringIsPresent{Field: u.Checksum, Name: "Checksum"},
	)
	return validate.Validate(vs...), nil
}

// BeforeCreate populates the StorageKey on a newly created UserUpload
func (u *Upload) BeforeCreate(tx *pop.Connection) error {
	// Populate ID if not exists
	if u.ID == uuid.Nil {
		u.ID = uuid.Must(uuid.NewV4())
	}

	if u.StorageKey == "" {
		u.StorageKey = path.Join(string(u.UploadType), "uploads", u.ID.String())
	}

	return nil
}

// FetchUpload returns an UserUpload if the user has access to that upload
func FetchUpload(ctx context.Context, db *pop.Connection, session *auth.Session, id uuid.UUID) (Upload, error) {
	var upload Upload
	err := db.Q().Where("uploads.deleted_at is null").Eager().Find(&upload, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Upload{}, errors.Wrap(ErrFetchNotFound, "error fetching upload")
		}
		// Otherwise, it's an unexpected err so we return that.
		return Upload{}, err
	}
	// In order to check permissions, need to fetch upload type
	if upload.UploadType == UploadTypeUSER {
		_, err := FetchUserUploadFromUploadID(ctx, db, session, upload.ID)
		if err != nil {
			return Upload{}, fmt.Errorf("error fetching user upload errors: %w", err)
		}
	} else if upload.UploadType == UploadTypePRIME {
		_, err := FetchPrimeUploadFromUploadID(ctx, db, session, upload.ID)
		if err != nil {
			return Upload{}, fmt.Errorf("error fetching prime upload errors: %w", err)
		}
	}
	return upload, nil
}

// DeleteUpload deletes an upload from the database
func DeleteUpload(dbConn *pop.Connection, upload *Upload) error {
	if dbConn.TX != nil {
		err := utilities.SoftDestroy(dbConn, upload)
		if err != nil {
			return err
		}
	} else {
		return dbConn.Transaction(func(db *pop.Connection) error {
			err := utilities.SoftDestroy(db, upload)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return nil
}
