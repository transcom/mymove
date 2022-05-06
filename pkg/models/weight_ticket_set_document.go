package models

import (
	"time"

	"github.com/transcom/mymove/pkg/auth"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// WeightTicketSetType represents types of weight ticket sets
type WeightTicketSetType string

const (
	// WeightTicketSetTypeCAR captures enum value "CAR"
	WeightTicketSetTypeCAR WeightTicketSetType = "CAR"

	// WeightTicketSetTypeCARTRAILER captures enum value "CAR_TRAILER"
	WeightTicketSetTypeCARTRAILER WeightTicketSetType = "CAR_TRAILER"

	// WeightTicketSetTypeBOXTRUCK captures enum value "BOX_TRUCK"
	WeightTicketSetTypeBOXTRUCK WeightTicketSetType = "BOX_TRUCK"

	// WeightTicketSetTypePROGEAR captures enum value "PRO_GEAR"
	WeightTicketSetTypePROGEAR WeightTicketSetType = "PRO_GEAR"
)

// WeightTicketSetDocument weight ticket documents payload
type WeightTicketSetDocument struct {
	ID                       uuid.UUID           `json:"id" db:"id"`
	MoveDocumentID           uuid.UUID           `json:"move_document_id" db:"move_document_id"`
	MoveDocument             MoveDocument        `belongs_to:"move_documents" fk_id:"move_document_id"`
	EmptyWeight              *unit.Pound         `json:"empty_weight,omitempty" db:"empty_weight"`
	EmptyWeightTicketMissing bool                `json:"empty_weight_ticket_missing,omitempty" db:"empty_weight_ticket_missing"`
	FullWeight               *unit.Pound         `json:"full_weight,omitempty" db:"full_weight"`
	FullWeightTicketMissing  bool                `json:"full_weight_ticket_missing,omitempty" db:"full_weight_ticket_missing"`
	VehicleNickname          *string             `json:"vehicle_nickname,omitempty" db:"vehicle_nickname"`
	VehicleMake              *string             `json:"vehicle_make,omitempty" db:"vehicle_make"`
	VehicleModel             *string             `json:"vehicle_model,omitempty" db:"vehicle_model"`
	WeightTicketSetType      WeightTicketSetType `json:"weight_ticket_set_type,omitempty" db:"weight_ticket_set_type"`
	WeightTicketDate         *time.Time          `json:"weight_ticket_date,omitempty" db:"weight_ticket_date"`
	TrailerOwnershipMissing  bool                `json:"trailer_ownership_missing,omitempty" db:"trailer_ownership_missing"`
	CreatedAt                time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time           `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time          `db:"deleted_at"`
}

// WeightTicketSetDocuments slice of WeightTicketSetDocuments
type WeightTicketSetDocuments []WeightTicketSetDocuments

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *WeightTicketSetDocument) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: m.MoveDocumentID, Name: "MoveDocumentID"},
		&validators.StringIsPresent{Field: string(m.WeightTicketSetType), Name: "WeightTicketSetType"},
		&MustBeBothNilOrBothHaveValue{
			FieldName1:  "VehicleMake",
			FieldValue1: m.VehicleMake,
			FieldName2:  "VehicleModel",
			FieldValue2: m.VehicleModel,
		},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *WeightTicketSetDocument) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *WeightTicketSetDocument) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// SumWeightTicketSetsForPPM iterates through move documents that are weight ticket sets and accumulates
// the net weight if the weight ticket set has an OK status
func SumWeightTicketSetsForPPM(db *pop.Connection, session *auth.Session, ppmID uuid.UUID) (*unit.Pound, error) {
	status := MoveDocumentStatusOK
	var totalWeight unit.Pound
	weightTicketSets, err := FetchMoveDocuments(db, session, ppmID, &status, MoveDocumentTypeWEIGHTTICKETSET, false)

	if err != nil {
		return &totalWeight, err
	}

	for _, weightTicketSet := range weightTicketSets {
		wt := weightTicketSet.WeightTicketSetDocument
		if wt != nil && wt.FullWeight != nil && wt.EmptyWeight != nil {
			totalWeight += *wt.FullWeight - *wt.EmptyWeight
		}
	}
	return &totalWeight, nil
}

// CreateWeightTicketSetDocument creates a moving weight ticket document associated to a move and move document
func (m Move) CreateWeightTicketSetDocument(
	db *pop.Connection,
	userUploads UserUploads,
	personallyProcuredMoveID *uuid.UUID,
	weightTicketSetDocument *WeightTicketSetDocument,
	moveType SelectedMoveType) (*WeightTicketSetDocument, *validate.Errors, error) {

	weightTicketSetTitle := "vehicle_weight"
	if weightTicketSetDocument.WeightTicketSetType == WeightTicketSetTypePROGEAR {
		weightTicketSetTitle = "pro_gear_weight"
	}

	var responseError error
	responseVErrors := validate.NewErrors()

	transactionErr := db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		var newMoveDocument *MoveDocument
		newMoveDocument, responseVErrors, responseError = m.createMoveDocumentWithoutTransaction(
			db,
			userUploads,
			personallyProcuredMoveID,
			MoveDocumentTypeWEIGHTTICKETSET,
			weightTicketSetTitle,
			weightTicketSetDocument.VehicleNickname,
			moveType)
		responseError = errors.Wrap(responseError, "Error creating move document")
		if responseVErrors.HasAny() || responseError != nil {
			return transactionError
		}

		weightTicketSetDocument.MoveDocument = *newMoveDocument
		weightTicketSetDocument.MoveDocumentID = newMoveDocument.ID

		verrs, err := db.ValidateAndCreate(weightTicketSetDocument)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error creating moving expense document")
			weightTicketSetDocument = nil
			return transactionError
		}

		return nil
	})

	if transactionErr != nil {
		return weightTicketSetDocument, responseVErrors, responseError
	}

	return weightTicketSetDocument, responseVErrors, responseError
}
