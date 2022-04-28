package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/utilities"
)

// An UserUpload represents an user uploaded file, such as an image or PDF.
type UserUpload struct {
	ID         uuid.UUID  `db:"id"`
	DocumentID *uuid.UUID `db:"document_id"`
	Document   Document   `belongs_to:"documents" fk_id:"document_id"`
	UploaderID uuid.UUID  `db:"uploader_id"`
	UploadID   uuid.UUID  `db:"upload_id"`
	Upload     Upload     `belongs_to:"uploads" fk_id:"upload_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

// UserUploads is not required by pop and may be deleted
type UserUploads []UserUpload

// UploadsFromUserUploads returns a slice of Uploads given a slice of UserUploads
func UploadsFromUserUploads(db *pop.Connection, userUploads UserUploads) (Uploads, error) {
	var uploads Uploads
	for _, userUpload := range userUploads {
		var upload Upload
		err := db.Q().Where("uploads.deleted_at is null").Find(&upload, userUpload.UploadID)
		if err != nil {
			if errors.Cause(err).Error() == RecordNotFoundErrorString {
				return Uploads{}, errors.Wrap(ErrFetchNotFound, "error fetching upload")
			}
			// Otherwise, it's an unexpected err so we return that.
			return Uploads{}, err
		}
		uploads = append(uploads, upload)
	}
	return uploads, nil
}

// UploadsFromUserUploadsNoDatabase returns a slice of Uploads given a slice of UserUploads
func UploadsFromUserUploadsNoDatabase(userUploads UserUploads) (Uploads, error) {
	var uploads Uploads
	for _, userUpload := range userUploads {
		if userUpload.UploadID != uuid.Nil {
			uploads = append(uploads, userUpload.Upload)
		} else {
			return Uploads{}, errors.New("error invalid UploadID in UserUpload")
		}
	}
	return uploads, nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *UserUpload) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: u.UploaderID, Name: "UploaderID"},
	), nil
}

// FetchUserUpload returns an UserUpload if the user has access to that upload
func FetchUserUpload(db *pop.Connection, session *auth.Session, id uuid.UUID) (UserUpload, error) {
	var userUpload UserUpload
	err := db.Q().
		Where("deleted_at is null").Eager("Document", "Upload").Find(&userUpload, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return UserUpload{}, errors.Wrap(ErrFetchNotFound, "error fetching user_uploads")
		}
		// Otherwise, it's an unexpected err so we return that.
		return UserUpload{}, err
	}

	// If there's a document, check permissions. Otherwise user must
	// have been the uploader
	if userUpload.DocumentID != nil {
		_, docErr := FetchDocument(db, session, *userUpload.DocumentID, false)
		if docErr != nil {
			return UserUpload{}, docErr
		}
	} else if userUpload.UploaderID != session.UserID {
		return UserUpload{}, errors.Wrap(ErrFetchNotFound, "user ID doesn't match uploader ID")
	}
	return userUpload, nil
}

// FetchUserUploadFromUploadID returns an UserUpload if the user has access to that upload
func FetchUserUploadFromUploadID(db *pop.Connection, session *auth.Session, uploadID uuid.UUID) (UserUpload, error) {
	var userUpload UserUpload
	err := db.Q().
		Join("uploads AS ups", "ups.id = user_uploads.upload_id").
		Where("ups.ID = $1 and user_uploads.deleted_at is null", uploadID).Eager("Document", "Upload").First(&userUpload)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return UserUpload{}, errors.Wrap(ErrFetchNotFound, "error fetching user_uploads")
		}
		// Otherwise, it's an unexpected err so we return that.
		return UserUpload{}, err
	}

	// If there's a document, check permissions. Otherwise user must
	// have been the uploader
	if userUpload.DocumentID != nil {
		_, docErr := FetchDocument(db, session, *userUpload.DocumentID, false)
		if docErr != nil {
			return UserUpload{}, docErr
		}
	} else if userUpload.UploaderID != session.UserID {
		return UserUpload{}, errors.Wrap(ErrFetchNotFound, "user ID doesn't match uploader ID")
	}
	return userUpload, nil
}

// DeleteUserUpload deletes an upload from the database
func DeleteUserUpload(dbConn *pop.Connection, userUpload *UserUpload) error {
	if dbConn.TX != nil {
		err := utilities.SoftDestroy(dbConn, userUpload)
		if err != nil {
			return err
		}
	} else {
		return dbConn.Transaction(func(db *pop.Connection) error {
			err := utilities.SoftDestroy(db, userUpload)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return nil
}
