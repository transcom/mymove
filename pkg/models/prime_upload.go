package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/db/utilities"
)

// An PrimeUpload represents an user uploaded file, such as an image or PDF.
type PrimeUpload struct {
	ID                  uuid.UUID         `db:"id"`
	ProofOfServiceDocID uuid.UUID         `db:"proof_of_service_docs_id"`
	ProofOfServiceDoc   ProofOfServiceDoc `belongs_to:"proof_of_service_docs" fk_id:"proof_of_service_docs_id"`
	ContractorID        uuid.UUID         `db:"contractor_id"`
	Contractor          Contractor        `belongs_to:"contractors" fk_id:"contractor_id"`
	UploadID            uuid.UUID         `db:"upload_id"`
	Upload              Upload            `belongs_to:"uploads" fk_id:"upload_id"`
	CreatedAt           time.Time         `db:"created_at"`
	UpdatedAt           time.Time         `db:"updated_at"`
	DeletedAt           *time.Time        `db:"deleted_at"`
}

// PrimeUploads is not required by pop and may be deleted
type PrimeUploads []PrimeUpload

// UploadsFromPrimeUploads return a slice of uploads given a slice of prime uploads
func UploadsFromPrimeUploads(db *pop.Connection, primeUploads PrimeUploads) (Uploads, error) {
	var uploads Uploads
	for _, PrimeUpload := range primeUploads {
		var upload Upload
		err := db.Q().Where("uploads.deleted_at is null").Eager("ProofOfServiceDoc", "Contractor", "Upload").Find(&upload, PrimeUpload.UploadID)
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

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *PrimeUpload) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: u.ContractorID, Name: "ContractorID"},
		&validators.UUIDIsPresent{Field: u.ProofOfServiceDocID, Name: "ProofOfServiceDocID"},
	), nil
}

// FetchPrimeUpload returns an PrimeUpload if the contractor has access to that upload
func FetchPrimeUpload(db *pop.Connection, contractorID uuid.UUID, id uuid.UUID) (PrimeUpload, error) {
	var primeUpload PrimeUpload
	err := db.Q().
		Join("uploads AS ups", "ups.id = prime_uploads.upload_id").
		Where("ups.deleted_at is null and prime_uploads.deleted_at is null").Eager("ProofOfServiceDoc", "Contractor", "Upload").Find(&primeUpload, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return PrimeUpload{}, errors.Wrap(ErrFetchNotFound, "error fetching prime_uploads")
		}
		// Otherwise, it's an unexpected err so we return that.
		return PrimeUpload{}, err
	}

	// If there's a proof of service doc, check permissions.
	if primeUpload.ContractorID != contractorID {
		return PrimeUpload{}, errors.Wrap(ErrFetchNotFound, "contractor ID doesn't match primeUpload.ContractorID")
	}
	return primeUpload, nil
}

// DeletePrimeUpload deletes an upload from the database
func DeletePrimeUpload(dbConn *pop.Connection, primeUpload *PrimeUpload) error {
	if dbConn.TX != nil {
		err := utilities.SoftDestroy(dbConn, primeUpload)
		if err != nil {
			return err
		}
	} else {
		return dbConn.Transaction(func(db *pop.Connection) error {
			err := utilities.SoftDestroy(db, primeUpload)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return nil
}
