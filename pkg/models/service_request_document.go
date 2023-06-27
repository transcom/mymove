package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// ServiceRequestDocument represents a document for a service item requst by the Prime
type ServiceRequestDocument struct {
	ID               uuid.UUID `json:"id" db:"id"`
	MTOServiceItemID uuid.UUID `json:"mto_service_item_id" db:"mto_service_item_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`

	//Associations
	MTOServiceItem                MTOServiceItem                `belongs_to:"mto_service_item" fk_id:"mto_service_item_id"`
	ServiceRequestDocumentUploads ServiceRequestDocumentUploads `has_many:"service_request_document_uploads" fk_id:"service_request_documents_id" order_by:"created_at asc"`
}

// TableName overrides the table name used by Pop.
func (s ServiceRequestDocument) TableName() string {
	return "service_request_documents"
}

type ServiceRequestDocuments []ServiceRequestDocument

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (s *ServiceRequestDocument) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.MTOServiceItemID, Name: "MTOServiceItemID"},
	), nil
}
