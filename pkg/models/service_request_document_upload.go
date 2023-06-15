package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/db/utilities"
)

// A ServiceRequestDocumentUpload represents an user uploaded file, such as an image or PDF.
type ServiceRequestDocumentUpload struct {
	ID                       uuid.UUID              `db:"id"`
	ServiceRequestDocumentID uuid.UUID              `db:"service_request_documents_id"`
	ServiceRequestDocument   ServiceRequestDocument `belongs_to:"service_request_documents" fk_id:"service_request_documents_id"`
	ContractorID             uuid.UUID              `db:"contractor_id"`
	Contractor               Contractor             `belongs_to:"contractors" fk_id:"contractor_id"`
	UploadID                 uuid.UUID              `db:"upload_id"`
	Upload                   Upload                 `belongs_to:"uploads" fk_id:"upload_id"`
	CreatedAt                time.Time              `db:"created_at"`
	UpdatedAt                time.Time              `db:"updated_at"`
	DeletedAt                *time.Time             `db:"deleted_at"`
}

// TableName overrides the table name used by Pop.
func (u ServiceRequestDocumentUpload) TableName() string {
	return "service_request_document_uploads"
}

type ServiceRequestDocumentUploads []ServiceRequestDocumentUpload

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (u *ServiceRequestDocumentUpload) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: u.ContractorID, Name: "ContractorID"},
		&validators.UUIDIsPresent{Field: u.ServiceRequestDocumentID, Name: "ServiceRequestDocumentID"},
	), nil
}

// DeleteServiceRequestDocumentUpload deletes an upload from the database
func DeleteServiceRequestDocumentUpload(dbConn *pop.Connection, serviceRequestUpload *ServiceRequestDocumentUpload) error {
	if dbConn.TX != nil {
		err := utilities.SoftDestroy(dbConn, serviceRequestUpload)
		if err != nil {
			return err
		}
	} else {
		return dbConn.Transaction(func(db *pop.Connection) error {
			err := utilities.SoftDestroy(db, serviceRequestUpload)
			if err != nil {
				return err
			}
			return nil
		})
	}
	return nil
}
