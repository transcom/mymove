package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
)

// MoveDocumentStatus represents the status of a move document record's lifecycle
type MoveDocumentStatus string

const (
	// MoveDocumentStatusAWAITINGREVIEW captures enum value "AWAITING_REVIEW"
	MoveDocumentStatusAWAITINGREVIEW MoveDocumentStatus = "AWAITING_REVIEW"
	// MoveDocumentStatusOK captures enum value "OK"
	MoveDocumentStatusOK MoveDocumentStatus = "OK"
	// MoveDocumentStatusHASISSUE captures enum value "HAS_ISSUE"
	MoveDocumentStatusHASISSUE MoveDocumentStatus = "HAS_ISSUE"
)

// MoveDocumentType represents types of different move documents
type MoveDocumentType string

const (
	// Shared Doc Types

	// MoveDocumentTypeOTHER captures enum value "OTHER"
	MoveDocumentTypeOTHER MoveDocumentType = "OTHER"
	// MoveDocumentTypeWEIGHTTICKET captures enum value "WEIGHT_TICKET"
	MoveDocumentTypeWEIGHTTICKET MoveDocumentType = "WEIGHT_TICKET"
	// MoveDocumentTypeWEIGHTTICKETREWEIGH captures enum value "WEIGHT_TICKET_REWEIGH"
	MoveDocumentTypeWEIGHTTICKETREWEIGH = "WEIGHT_TICKET_REWEIGH"

	// PPM Doc Types

	// MoveDocumentTypeSTORAGEEXPENSE captures enum value "STORAGE_EXPENSE"
	MoveDocumentTypeSTORAGEEXPENSE MoveDocumentType = "STORAGE_EXPENSE"
	// MoveDocumentTypeSHIPMENTSUMMARY captures enum value "SHIPMENT_SUMMARY"
	MoveDocumentTypeSHIPMENTSUMMARY MoveDocumentType = "SHIPMENT_SUMMARY"
	// MoveDocumentTypeEXPENSE captures enum value "EXPENSE"
	MoveDocumentTypeEXPENSE MoveDocumentType = "EXPENSE"

	// HHG Doc Types

	// MoveDocumentTypeGOVBILLOFLADING captures enum value "GOV_BILL_OF_LADING"
	MoveDocumentTypeGOVBILLOFLADING MoveDocumentType = "GOV_BILL_OF_LADING"
	// MoveDocumentTypeORIGINPACKET captures enum value "ORIGIN_PACKET"
	MoveDocumentTypeORIGINPACKET = "ORIGIN_PACKET"
	// MoveDocumentTypeORIGIN619 captures enum value "ORIGIN_619"
	MoveDocumentTypeORIGIN619 = "ORIGIN_619"
	// MoveDocumentTypeORIGININVENTORY captures enum value "ORIGIN_INVENTORY"
	MoveDocumentTypeORIGININVENTORY = "ORIGIN_INVENTORY"
	// MoveDocumentTypeDESTINATIONPACKET captures enum value "DESTINATION_PACKET"
	MoveDocumentTypeDESTINATIONPACKET = "DESTINATION_PACKET"
	// MoveDocumentTypeDESTINATION619 captures enum value "DESTINATION_619"
	MoveDocumentTypeDESTINATION619 = "DESTINATION_619"
	// MoveDocumentTypeDESTINATION6191 captures enum value "DESTINATION_619_1"
	MoveDocumentTypeDESTINATION6191 = "DESTINATION_619_1"
	// MoveDocumentTypeDESTINATIONINVENTORY captures enum value "DESTINATION_INVENTORY"
	MoveDocumentTypeDESTINATIONINVENTORY = "DESTINATION_INVENTORY"
	// MoveDocumentTypeTHIRDPARTYINVOICE captures enum value "THIRD_PARTY_INVOICE"
	MoveDocumentTypeTHIRDPARTYINVOICE = "THIRD_PARTY_INVOICE"
	// MoveDocumentTypeTHIRDPARTYESTIMATE captures enum value "THIRD_PARTY_ESTIMATE"
	MoveDocumentTypeTHIRDPARTYESTIMATE = "THIRD_PARTY_ESTIMATE"
	// MoveDocumentTypeNOTICEOFLOSSORDAMAGE captures enum value "NOTICE_OF_LOSS_OR_DAMAGE"
	MoveDocumentTypeNOTICEOFLOSSORDAMAGE = "NOTICE_OF_LOSS_OR_DAMAGE"
	// MoveDocumentTypeFIREARMSCHAINOFCUSTODY captures enum value "FIREARMS_CHAIN_OF_CUSTODY"
	MoveDocumentTypeFIREARMSCHAINOFCUSTODY = "FIREARMS_CHAIN_OF_CUSTODY"
	// MoveDocumentTypePHOTO captures enum value "PHOTO"
	MoveDocumentTypePHOTO = "PHOTO"
)

// MoveDocumentSaveAction represents actions that can be taken during save
type MoveDocumentSaveAction string

const (
	// MoveDocumentSaveActionDELETEEXPENSEMODEL encodes an action to delete a linked expense model
	MoveDocumentSaveActionDELETEEXPENSEMODEL MoveDocumentSaveAction = "DELETE_EXPENSE_MODEL"
	// MoveDocumentSaveActionSAVEEXPENSEMODEL encodes an action to save a linked expense model
	MoveDocumentSaveActionSAVEEXPENSEMODEL MoveDocumentSaveAction = "SAVE_EXPENSE_MODEL"
)

// MoveDocument is an object representing a move document
type MoveDocument struct {
	ID                       uuid.UUID              `json:"id" db:"id"`
	DocumentID               uuid.UUID              `json:"document_id" db:"document_id"`
	Document                 Document               `belongs_to:"documents"`
	MoveID                   uuid.UUID              `json:"move_id" db:"move_id"`
	Move                     Move                   `belongs_to:"moves"`
	PersonallyProcuredMoveID *uuid.UUID             `json:"personally_procured_move_id" db:"personally_procured_move_id"`
	PersonallyProcuredMove   PersonallyProcuredMove `belongs_to:"personally_procured_moves"`
	ShipmentID               *uuid.UUID             `json:"shipment_id" db:"shipment_id"`
	Shipment                 Shipment               `belongs_to:"shipments"`
	Title                    string                 `json:"title" db:"title"`
	Status                   MoveDocumentStatus     `json:"status" db:"status"`
	MoveDocumentType         MoveDocumentType       `json:"move_document_type" db:"move_document_type"`
	MovingExpenseDocument    *MovingExpenseDocument `has_one:"moving_expense_document"`
	Notes                    *string                `json:"notes" db:"notes"`
	CreatedAt                time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time              `json:"updated_at" db:"updated_at"`
}

// MoveDocuments is not required by pop and may be deleted
type MoveDocuments []MoveDocument

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MoveDocument) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: m.DocumentID, Name: "DocumentID"},
		&validators.UUIDIsPresent{Field: m.MoveID, Name: "MoveID"},
		&validators.StringIsPresent{Field: string(m.Title), Name: "Title"},
		&validators.StringIsPresent{Field: string(m.Status), Name: "Status"},
		&validators.StringIsPresent{Field: string(m.MoveDocumentType), Name: "MoveDocumentType"},
	), nil
}

// State Machinery

// AttemptTransition is glue for when you are modifying the status of a model
// via a PUT rather than an action url. This translates the target status into an action method.
func (m *MoveDocument) AttemptTransition(targetStatus MoveDocumentStatus) error {
	// If it's the same it's not a transition
	if targetStatus == m.Status {
		return nil
	}

	switch targetStatus {
	case MoveDocumentStatusOK:
		return m.Approve()
	case MoveDocumentStatusHASISSUE:
		return m.Reject()
	}

	return errors.Wrap(ErrInvalidTransition, string(targetStatus))
}

// Approve marks the Document as OK
func (m *MoveDocument) Approve() error {
	if m.Status == MoveDocumentStatusOK {
		return errors.Wrap(ErrInvalidTransition, "Approve")
	}

	m.Status = MoveDocumentStatusOK
	return nil
}

// Reject marks the Document as HAS_ISSUE
func (m *MoveDocument) Reject() error {
	if m.Status == MoveDocumentStatusHASISSUE {
		return errors.Wrap(ErrInvalidTransition, "Reject")
	}

	m.Status = MoveDocumentStatusHASISSUE
	return nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MoveDocument) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MoveDocument) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// FetchMoveDocument fetches a MoveDocument model
func FetchMoveDocument(db *pop.Connection, session *auth.Session, id uuid.UUID) (*MoveDocument, error) {
	// Allow all office users to fetch move doc
	if session.IsOfficeApp() && session.OfficeUserID == uuid.Nil {
		return &MoveDocument{}, ErrFetchForbidden
	}

	var moveDoc MoveDocument
	err := db.Q().Eager("Document.Uploads", "Move", "PersonallyProcuredMove", "Shipment").Find(&moveDoc, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		return nil, err
	}

	// Pointer associations are buggy, so we manually load expense document things
	movingExpenseDocument := MovingExpenseDocument{}
	moveDoc.MovingExpenseDocument = nil
	err = db.Where("move_document_id = $1", moveDoc.ID.String()).Eager().First(&movingExpenseDocument)
	if err != nil {
		if errors.Cause(err).Error() != recordNotFoundErrorString {
			return nil, err
		}
	} else {
		moveDoc.MovingExpenseDocument = &movingExpenseDocument
	}

	// Check that the logged-in service member is associated to the document
	if session.IsMilApp() && moveDoc.Document.ServiceMemberID != session.ServiceMemberID {
		return &MoveDocument{}, ErrFetchForbidden
	}

	return &moveDoc, nil
}

// FetchApprovedMovingExpenseDocuments fetches all approved move expense document for a ppm
func FetchApprovedMovingExpenseDocuments(db *pop.Connection, session *auth.Session, ppmID uuid.UUID) (MoveDocuments, error) {
	// Allow all logged in office users to fetch move docs
	if session.IsOfficeApp() && session.OfficeUserID == uuid.Nil {
		return nil, ErrFetchForbidden
	}
	// Validate the move is associated to the logged-in service member
	_, fetchErr := FetchPersonallyProcuredMove(db, session, ppmID)
	if fetchErr != nil {
		return nil, ErrFetchForbidden
	}

	var moveDocuments MoveDocuments
	err := db.Where("move_document_type = $1", string(MoveDocumentTypeEXPENSE)).Where("status = $2", string(MoveDocumentStatusOK)).Where("personally_procured_move_id = $3", ppmID.String()).All(&moveDocuments)
	if err != nil {
		if errors.Cause(err).Error() != recordNotFoundErrorString {
			return nil, err
		}
	}

	for i, moveDoc := range moveDocuments {
		movingExpenseDocument := MovingExpenseDocument{}
		moveDoc.MovingExpenseDocument = nil
		err = db.Where("move_document_id = $1", moveDoc.ID.String()).Eager().First(&movingExpenseDocument)
		if err != nil {
			if errors.Cause(err).Error() != recordNotFoundErrorString {
				return nil, err
			}
		} else {
			moveDocuments[i].MovingExpenseDocument = &movingExpenseDocument
		}
	}

	return moveDocuments, nil
}

// FetchMoveDocumentsByTypeForShipment fetches move documents for shipment and move document type
func FetchMoveDocumentsByTypeForShipment(db *pop.Connection, session *auth.Session, moveDocumentType MoveDocumentType, shipmentID uuid.UUID) (MoveDocuments, error) {

	// Verify that the logged-in TSP user is authorized to generate GBL
	// Does this need to be checked here if already checked in create gbl handler?
	if session.IsTspApp() {
		if session.TspUserID == uuid.Nil {
			return nil, ErrFetchForbidden
		}
		tspUser, err := FetchTspUserByID(db, session.TspUserID)
		if err != nil {
			return nil, ErrFetchNotFound
		}
		shipment, err := FetchShipmentByTSP(db, tspUser.TransportationServiceProviderID, shipmentID)
		if err != nil {
			return nil, ErrFetchForbidden
		}
		if shipment.ID != shipmentID {
			return nil, ErrFetchForbidden
		}
	}

	// Allow all logged in office users to fetch move docs
	if session.IsOfficeApp() && session.OfficeUserID == uuid.Nil {
		return nil, ErrFetchForbidden
	}

	var moveDocuments MoveDocuments
	err := db.Where("move_document_type = $1", string(moveDocumentType)).Where("shipment_id = $2", shipmentID.String()).All(&moveDocuments)
	if err != nil {
		if errors.Cause(err).Error() != recordNotFoundErrorString {
			return nil, err
		}
	}
	return moveDocuments, nil
}

// SaveMoveDocument saves a move document
func SaveMoveDocument(db *pop.Connection, moveDocument *MoveDocument, saveAction MoveDocumentSaveAction) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()

	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		if saveAction == MoveDocumentSaveActionSAVEEXPENSEMODEL {
			// Save expense document
			expenseDocument := moveDocument.MovingExpenseDocument
			if verrs, err := db.ValidateAndSave(expenseDocument); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Creating Moving Expense Document")
				return transactionError
			}
		} else if saveAction == MoveDocumentSaveActionDELETEEXPENSEMODEL {
			// destroy expense document
			expenseDocument := moveDocument.MovingExpenseDocument
			if err := db.Destroy(expenseDocument); err != nil {
				responseError = errors.Wrap(err, "Error Deleting Moving Expense Document")
				return transactionError
			}
			moveDocument.MovingExpenseDocument = nil
		}

		// Updating the move document can cause the PPM to be updated
		if moveDocument.PersonallyProcuredMoveID != nil {
			ppm := moveDocument.PersonallyProcuredMove

			if verrs, err := db.ValidateAndSave(&ppm); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Saving Move Document's PPM")
				return transactionError
			}
		}

		// Updating the move document can cause the Shipment to be updated
		if moveDocument.ShipmentID != nil {
			shipment := moveDocument.Shipment

			if verrs, err := db.ValidateAndSave(&shipment); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Saving Move Document's Shipment")
				return transactionError
			}
		}

		// Finally, save the MoveDocument
		if verrs, err := db.ValidateAndSave(moveDocument); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving Move Document")
			return transactionError
		}

		return nil
	})

	return responseVErrors, responseError
}
