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

// MoveDocumentStatus represents the status of a move document record's lifecycle
type MoveDocumentStatus string

const (
	// MoveDocumentStatusAWAITINGREVIEW captures enum value "AWAITING_REVIEW"
	MoveDocumentStatusAWAITINGREVIEW MoveDocumentStatus = "AWAITING_REVIEW"
	// MoveDocumentStatusOK captures enum value "OK"
	MoveDocumentStatusOK MoveDocumentStatus = "OK"
	// MoveDocumentStatusHASISSUE captures enum value "HAS_ISSUE"
	MoveDocumentStatusHASISSUE MoveDocumentStatus = "HAS_ISSUE"
	// MoveDocumentStatusEXCLUDEFROMCALCULATION captures enum value "EXCLUDE_FROM_CALCULATION"
	MoveDocumentStatusEXCLUDEFROMCALCULATION MoveDocumentStatus = "EXCLUDE_FROM_CALCULATION"
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
	MoveDocumentTypeWEIGHTTICKETREWEIGH = "WEIGHT_TICKET_REWEIGH" // TODO: remove reweigh type?

	// PPM Doc Types

	// MoveDocumentTypeSTORAGEEXPENSE captures enum value "STORAGE_EXPENSE"
	MoveDocumentTypeSTORAGEEXPENSE MoveDocumentType = "STORAGE_EXPENSE"
	// MoveDocumentTypeSHIPMENTSUMMARY captures enum value "SHIPMENT_SUMMARY"
	MoveDocumentTypeSHIPMENTSUMMARY MoveDocumentType = "SHIPMENT_SUMMARY"
	// MoveDocumentTypeEXPENSE captures enum value "EXPENSE"
	MoveDocumentTypeEXPENSE MoveDocumentType = "EXPENSE"
	// MoveDocumentTypeWEIGHTTICKETSET captures enum value "WEIGHT_TICKET_SET"
	MoveDocumentTypeWEIGHTTICKETSET MoveDocumentType = "WEIGHT_TICKET_SET"

	// TODO: remove HHG doc types
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

// MoveExpenseDocumentSaveAction represents actions that can be taken during save for expense documents
type MoveExpenseDocumentSaveAction string

const (
	// MoveDocumentSaveActionDELETEEXPENSEMODEL encodes an action to delete a linked expense model
	MoveDocumentSaveActionDELETEEXPENSEMODEL MoveExpenseDocumentSaveAction = "DELETE_EXPENSE_MODEL"
	// MoveDocumentSaveActionSAVEEXPENSEMODEL encodes an action to save a linked expense model
	MoveDocumentSaveActionSAVEEXPENSEMODEL MoveExpenseDocumentSaveAction = "SAVE_EXPENSE_MODEL"
)

// MoveWeightTicketSetDocumentSaveAction represents actions that can be taken during save for weight ticket set documents
type MoveWeightTicketSetDocumentSaveAction string

const (
	// MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL encodes an action to delete a linked expense model
	MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL MoveWeightTicketSetDocumentSaveAction = "DELETE_WEIGHT_TICKET_SET_MODEL"
	// MoveDocumentSaveActionSAVEWEIGHTTICKETSETMODEL encodes an action to save a linked expense model
	MoveDocumentSaveActionSAVEWEIGHTTICKETSETMODEL MoveWeightTicketSetDocumentSaveAction = "SAVE_WEIGHT_TICKET_SET_MODEL"
)

// MoveDocument is an object representing a move document
type MoveDocument struct {
	ID                       uuid.UUID                `json:"id" db:"id"`
	DocumentID               uuid.UUID                `json:"document_id" db:"document_id"`
	Document                 Document                 `belongs_to:"documents" fk_id:"document_id"`
	MoveID                   uuid.UUID                `json:"move_id" db:"move_id"`
	Move                     Move                     `belongs_to:"moves" fk_id:"move_id"`
	PersonallyProcuredMoveID *uuid.UUID               `json:"personally_procured_move_id" db:"personally_procured_move_id"`
	PersonallyProcuredMove   PersonallyProcuredMove   `belongs_to:"personally_procured_moves" fk_id:"personally_procured_move_id"`
	Title                    string                   `json:"title" db:"title"`
	Status                   MoveDocumentStatus       `json:"status" db:"status"`
	MoveDocumentType         MoveDocumentType         `json:"move_document_type" db:"move_document_type"`
	MovingExpenseDocument    *MovingExpenseDocument   `has_one:"moving_expense_document"`
	Notes                    *string                  `json:"notes" db:"notes"`
	CreatedAt                time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time                `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time               `db:"deleted_at"`
	WeightTicketSetDocument  *WeightTicketSetDocument `has_one:"weight_ticket_set_document"`
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
	case MoveDocumentStatusEXCLUDEFROMCALCULATION:
		return m.Exclude()
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

// Exclude marks the Document as HAS_ISSUE
func (m *MoveDocument) Exclude() error {
	if m.Status == MoveDocumentStatusEXCLUDEFROMCALCULATION {
		return errors.Wrap(ErrInvalidTransition, "Exclude")
	}

	m.Status = MoveDocumentStatusEXCLUDEFROMCALCULATION
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

// DeleteMoveDocument deletes a MoveDocument model
func DeleteMoveDocument(db *pop.Connection, moveDoc *MoveDocument) error {
	docType := moveDoc.MoveDocumentType

	// only delete weight tickets, weight ticket sets, and expense documents at this time
	if !(docType == MoveDocumentTypeEXPENSE || docType == MoveDocumentTypeWEIGHTTICKETSET || docType == MoveDocumentTypeWEIGHTTICKET) {
		return errors.New("Can only delete weight ticket set and expense documents")
	}

	return db.Transaction(func(db *pop.Connection) error {
		return utilities.SoftDestroy(db, moveDoc)
	})
}

// FetchMoveDocument fetches a MoveDocument model
func FetchMoveDocument(db *pop.Connection, session *auth.Session, id uuid.UUID, includedDeletedMoveDocuments bool) (*MoveDocument, error) {
	// Allow all office users to fetch move doc
	if session.IsOfficeApp() && session.OfficeUserID == uuid.Nil {
		return &MoveDocument{}, ErrFetchForbidden
	}

	var moveDoc MoveDocument
	query := db.Q()

	if !includedDeletedMoveDocuments {
		query = query.Where("deleted_at is null")
	}

	err := query.Eager("Document.UserUploads.Upload", "Move", "PersonallyProcuredMove").Find(&moveDoc, id)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		return nil, err
	}

	// Pointer associations are buggy, so we manually load expense document things
	q := db.Where("move_document_id = $1", moveDoc.ID.String())
	movingExpenseDocument := &MovingExpenseDocument{}
	var movingDocumentErr error
	if movingDocumentErr = q.Eager().First(movingExpenseDocument); movingDocumentErr == nil {
		moveDoc.MovingExpenseDocument = movingExpenseDocument
	}
	if movingDocumentErr != nil && errors.Cause(movingDocumentErr).Error() != RecordNotFoundErrorString {
		return nil, err
	}

	weightTicketSetDocument := &WeightTicketSetDocument{}
	var weightTicketSetDocumentErr error
	if weightTicketSetDocumentErr = q.Eager().First(weightTicketSetDocument); weightTicketSetDocumentErr == nil {
		moveDoc.WeightTicketSetDocument = weightTicketSetDocument
	}
	if weightTicketSetDocumentErr != nil && errors.Cause(weightTicketSetDocumentErr).Error() != RecordNotFoundErrorString {
		return nil, err
	}

	// Check that the logged-in service member is associated to the document
	if session.IsMilApp() && moveDoc.Document.ServiceMemberID != session.ServiceMemberID {
		return &MoveDocument{}, ErrFetchForbidden
	}

	return &moveDoc, nil
}

// FetchMoveDocuments fetches all move expense and weight ticket set documents for a ppm
// the optional status parameter can be used for restricting to a subset of statuses.
func FetchMoveDocuments(db *pop.Connection, session *auth.Session, ppmID uuid.UUID, status *MoveDocumentStatus, moveDocumentType MoveDocumentType, includedDeletedMoveDocuments bool) (MoveDocuments, error) {
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
	query := db.Q()

	if !includedDeletedMoveDocuments {
		query = query.Where("deleted_at is null")
	}

	query = query.Where("move_document_type = $1", string(moveDocumentType)).Where("personally_procured_move_id = $2", ppmID.String())
	if status != nil {
		query = query.Where("status = $3", string(*status))
	}
	err := query.All(&moveDocuments)
	if err != nil {
		if errors.Cause(err).Error() != RecordNotFoundErrorString {
			return nil, err
		}
	}

	for i, moveDoc := range moveDocuments {
		movingExpenseDocument := MovingExpenseDocument{}
		moveDoc.MovingExpenseDocument = nil
		err = db.Where("move_document_id = $1", moveDoc.ID.String()).Eager().First(&movingExpenseDocument)
		if err != nil {
			if errors.Cause(err).Error() != RecordNotFoundErrorString {
				return nil, err
			}
		} else {
			moveDocuments[i].MovingExpenseDocument = &movingExpenseDocument
		}
	}

	for i, moveDoc := range moveDocuments {
		weightTicketSet := WeightTicketSetDocument{}
		moveDoc.WeightTicketSetDocument = nil
		err = db.Where("move_document_id = $1", moveDoc.ID.String()).Eager().First(&weightTicketSet)
		if err != nil {
			if errors.Cause(err).Error() != RecordNotFoundErrorString {
				return nil, err
			}
		} else {
			moveDocuments[i].WeightTicketSetDocument = &weightTicketSet
		}
	}

	return moveDocuments, nil
}

// SaveMoveDocument saves a move document
func SaveMoveDocument(db *pop.Connection, moveDocument *MoveDocument, saveExpenseAction MoveExpenseDocumentSaveAction, saveWeightTicketSetAction MoveWeightTicketSetDocumentSaveAction) (*validate.Errors, error) {
	var responseError error
	responseVErrors := validate.NewErrors()

	err := db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback the transaction")

		if saveExpenseAction == MoveDocumentSaveActionSAVEEXPENSEMODEL {
			// Save expense document
			expenseDocument := moveDocument.MovingExpenseDocument
			if verrs, err := db.ValidateAndSave(expenseDocument); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Creating Moving Expense Document")
				return transactionError
			}
		} else if saveExpenseAction == MoveDocumentSaveActionDELETEEXPENSEMODEL {
			// destroy expense document
			expenseDocument := moveDocument.MovingExpenseDocument
			if err := utilities.SoftDestroy(db, expenseDocument); err != nil {
				responseError = errors.Wrap(err, "Error Deleting Moving Expense Document")
				return transactionError
			}
			moveDocument.MovingExpenseDocument = nil
		}

		if saveWeightTicketSetAction == MoveDocumentSaveActionSAVEWEIGHTTICKETSETMODEL {
			// save weight ticket set document
			weightTicketSetDocument := moveDocument.WeightTicketSetDocument
			if verrs, err := db.ValidateAndSave(weightTicketSetDocument); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Creating Moving Expense Document")
				return transactionError
			}
		} else if saveWeightTicketSetAction == MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL {
			// destroy weight ticket set document
			weightTicketSetDocument := moveDocument.WeightTicketSetDocument
			if err := utilities.SoftDestroy(db, weightTicketSetDocument); err != nil {
				responseError = errors.Wrap(err, "Error Deleting Moving Expense Document")
				return transactionError
			}
			moveDocument.WeightTicketSetDocument = nil
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

		// Finally, save the MoveDocument
		if verrs, err := db.ValidateAndSave(moveDocument); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving Move Document")
			return transactionError
		}

		return nil
	})

	if err != nil {
		return responseVErrors, responseError
	}

	return responseVErrors, responseError
}
