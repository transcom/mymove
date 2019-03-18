package models

import (
	"crypto/sha256"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// MoveStatus represents the status of an order record's lifecycle
type MoveStatus string

const (
	// MoveStatusDRAFT captures enum value "DRAFT"
	MoveStatusDRAFT MoveStatus = "DRAFT"
	// MoveStatusSUBMITTED captures enum value "SUBMITTED"
	MoveStatusSUBMITTED MoveStatus = "SUBMITTED"
	// MoveStatusAPPROVED captures enum value "APPROVED"
	MoveStatusAPPROVED MoveStatus = "APPROVED"
	// MoveStatusCOMPLETED captures enum value "COMPLETED"
	MoveStatusCOMPLETED MoveStatus = "COMPLETED"
	// MoveStatusCANCELED captures enum value "CANCELED"
	MoveStatusCANCELED MoveStatus = "CANCELED"
)

// SelectedMoveType represents the type of move being represented
type SelectedMoveType string

// This lists available move types in the system
// Combination move types like HHG+PPM should be added as an underscore separated list
// The list should be lexigraphically sorted. Ex: UB + PPM will always be 'PPM_UB'
const (
	// MoveStatusHHG captures enum value "HHG" for House Hold Goods
	SelectedMoveTypeHHG SelectedMoveType = "HHG"
	// MoveStatusPPM captures enum value "PPM" for Personally Procured Move
	SelectedMoveTypePPM SelectedMoveType = "PPM"
	// MoveStatusUB captures enum value "UB" for Unaccompanied Baggage
	SelectedMoveTypeUB SelectedMoveType = "UB"
	// MoveStatusPOV captures enum value "POV" for Privately-Owned Vehicle
	SelectedMoveTypePOV SelectedMoveType = "POV"
	// MoveStatusNTS captures enum value "NTS" for Non-Temporary Storage
	SelectedMoveTypeNTS SelectedMoveType = "NTS"
	// MoveStatusHHGPPM captures enum value "HHG_PPM" for combination move HHG + PPM
	SelectedMoveTypeHHGPPM SelectedMoveType = "HHG_PPM"
)

const maxLocatorAttempts = 3
const locatorLength = 6

// This set of letters should produce 'non-word' type strings
var locatorLetters = []rune("346789BCDFGHJKMPQRTVWXY")

// Move is an object representing a move
type Move struct {
	ID                      uuid.UUID               `json:"id" db:"id"`
	Locator                 string                  `json:"locator" db:"locator"`
	CreatedAt               time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time               `json:"updated_at" db:"updated_at"`
	OrdersID                uuid.UUID               `json:"orders_id" db:"orders_id"`
	Orders                  Order                   `belongs_to:"orders"`
	SelectedMoveType        *SelectedMoveType       `json:"selected_move_type" db:"selected_move_type"`
	PersonallyProcuredMoves PersonallyProcuredMoves `has_many:"personally_procured_moves" order_by:"created_at desc"`
	Shipments               Shipments               `has_many:"shipments"`
	MoveDocuments           MoveDocuments           `has_many:"move_documents" order_by:"created_at desc"`
	Status                  MoveStatus              `json:"status" db:"status"`
	SignedCertifications    SignedCertifications    `has_many:"signed_certifications" order_by:"created_at desc"`
	CancelReason            *string                 `json:"cancel_reason" db:"cancel_reason"`
}

// Moves is not required by pop and may be deleted
type Moves []Move

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Move) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: m.OrdersID, Name: "OrdersID"},
		&validators.StringIsPresent{Field: string(m.Status), Name: "Status"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *Move) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *Move) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// State Machine
// Avoid calling Move.Status = ... ever. Use these methods to change the state.

// Submit submits the Move
func (m *Move) Submit() error {
	if m.Status != MoveStatusDRAFT {
		return errors.Wrap(ErrInvalidTransition, "Submit")
	}

	m.Status = MoveStatusSUBMITTED

	// Update PPM status too
	for i := range m.PersonallyProcuredMoves {
		err := m.PersonallyProcuredMoves[i].Submit()
		if err != nil {
			return err
		}
	}

	// Update HHG (Shipment) status too
	for i := range m.Shipments {
		err := m.Shipments[i].Submit()
		if err != nil {
			return err
		}
	}

	for _, ppm := range m.PersonallyProcuredMoves {
		if ppm.Advance != nil {
			err := ppm.Advance.Request()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Approve approves the Move
func (m *Move) Approve() error {
	if m.Status != MoveStatusSUBMITTED {
		return errors.Wrap(ErrInvalidTransition, "Approve")
	}

	m.Status = MoveStatusAPPROVED
	return nil
}

// Complete Completes the Move
func (m *Move) Complete() error {
	if m.Status != MoveStatusAPPROVED {
		return errors.Wrap(ErrInvalidTransition, "Complete")
	}

	m.Status = MoveStatusCOMPLETED
	return nil
}

// Cancel cancels the Move and its associated PPMs
func (m *Move) Cancel(reason string) error {
	// We can cancel any move that isn't already complete.
	if m.Status == MoveStatusCOMPLETED || m.Status == MoveStatusCANCELED {
		return errors.Wrap(ErrInvalidTransition, "Cancel")
	}

	m.Status = MoveStatusCANCELED

	// If a reason was submitted, add it to the move record.
	if reason != "" {
		m.CancelReason = &reason
	}

	// This will work only if you use the PPM in question rather than a var representing it
	// i.e. you can't use _, ppm := range PPMs, has to be PPMS[i] as below
	for i := range m.PersonallyProcuredMoves {
		err := m.PersonallyProcuredMoves[i].Cancel()
		if err != nil {
			return err
		}
	}

	// TODO: Orders can exist after related moves are canceled
	err := m.Orders.Cancel()
	if err != nil {
		return err
	}

	return nil
}

// FetchMove fetches and validates a Move for this User
func FetchMove(db *pop.Connection, session *auth.Session, id uuid.UUID) (*Move, error) {
	var move Move
	err := db.Q().Eager("PersonallyProcuredMoves.Advance",
		"SignedCertifications",
		"Orders",
		"MoveDocuments.Document",
		"Shipments.TrafficDistributionList",
		"Shipments.ServiceAgents").Find(&move, id)

	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	// Eager loading of nested has_many associations is broken
	for i, moveDoc := range move.MoveDocuments {
		db.Load(&moveDoc.Document, "Uploads")
		move.MoveDocuments[i] = moveDoc
	}

	// Eager loading of nested has_many associations is broken
	for i, shipment := range move.Shipments {
		if shipment.PickupAddressID != nil {
			pickupAddress := Address{}
			if err = db.Find(&pickupAddress, shipment.PickupAddressID); err == nil {
				move.Shipments[i].PickupAddress = &pickupAddress
			}
		}

		if shipment.HasSecondaryPickupAddress && shipment.SecondaryPickupAddressID != nil {
			secondaryPickupAddress := Address{}
			if err = db.Find(&secondaryPickupAddress, shipment.SecondaryPickupAddressID); err == nil {
				move.Shipments[i].SecondaryPickupAddress = &secondaryPickupAddress
			}
		}

		if shipment.HasDeliveryAddress && shipment.DeliveryAddressID != nil {
			deliveryAddress := Address{}
			if err = db.Find(&deliveryAddress, shipment.DeliveryAddressID); err == nil {
				move.Shipments[i].DeliveryAddress = &deliveryAddress
			}
		}
	}

	// Ensure that the logged-in user is authorized to access this move
	_, authErr := FetchOrderForUser(db, session, move.OrdersID)
	if authErr != nil {
		return nil, authErr
	}

	return &move, nil
}

func (m Move) createMoveDocumentWithoutTransaction(
	db *pop.Connection,
	uploads Uploads,
	modelID *uuid.UUID,
	moveDocumentType MoveDocumentType,
	title string,
	notes *string,
	moveType SelectedMoveType) (*MoveDocument, *validate.Errors, error) {

	var responseError error
	responseVErrors := validate.NewErrors()

	// Make a generic Document
	newDoc := Document{
		ServiceMemberID: m.Orders.ServiceMemberID,
		Uploads:         uploads,
	}
	verrs, err := db.ValidateAndCreate(&newDoc)
	if err != nil || verrs.HasAny() {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error creating document for move document")
		return nil, responseVErrors, responseError
	}

	// Associate uploads to the new document
	for _, upload := range uploads {
		upload.DocumentID = &newDoc.ID
		verrs, err := db.ValidateAndUpdate(&upload)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error updating upload")
			return nil, responseVErrors, responseError
		}
	}

	var newMoveDocument *MoveDocument
	if moveType == SelectedMoveTypeHHG || moveType == SelectedMoveTypeHHGPPM {
		newMoveDocument = &MoveDocument{
			Move:             m,
			MoveID:           m.ID,
			Document:         newDoc,
			DocumentID:       newDoc.ID,
			ShipmentID:       modelID,
			MoveDocumentType: moveDocumentType,
			Title:            title,
			Status:           MoveDocumentStatusAWAITINGREVIEW,
		}
	} else {
		// Finally create the MoveDocument to tie it to the Move
		newMoveDocument = &MoveDocument{
			Move:                     m,
			MoveID:                   m.ID,
			Document:                 newDoc,
			DocumentID:               newDoc.ID,
			PersonallyProcuredMoveID: modelID,
			MoveDocumentType:         moveDocumentType,
			Title:                    title,
			Status:                   MoveDocumentStatusAWAITINGREVIEW,
			Notes:                    notes,
		}
	}

	verrs, err = db.ValidateAndCreate(newMoveDocument)
	if err != nil || verrs.HasAny() {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error creating move document")
		return nil, responseVErrors, responseError
	}

	return newMoveDocument, responseVErrors, nil
}

// CreateMoveDocument creates a move document associated to a move & ppm or shipment
func (m Move) CreateMoveDocument(
	db *pop.Connection,
	uploads Uploads,
	modelID *uuid.UUID,
	moveDocumentType MoveDocumentType,
	title string,
	notes *string,
	moveType SelectedMoveType) (*MoveDocument, *validate.Errors, error) {

	var newMoveDocument *MoveDocument
	var responseError error
	responseVErrors := validate.NewErrors()

	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		newMoveDocument, responseVErrors, responseError = m.createMoveDocumentWithoutTransaction(
			db,
			uploads,
			modelID,
			moveDocumentType,
			title,
			notes,
			moveType)

		if responseVErrors.HasAny() || responseError != nil {
			return transactionError
		}

		return nil

	})

	return newMoveDocument, responseVErrors, responseError
}

// CreateMovingExpenseDocument creates a moving expense document associated to a move and move document
func (m Move) CreateMovingExpenseDocument(
	db *pop.Connection,
	uploads Uploads,
	personallyProcuredMoveID *uuid.UUID,
	moveDocumentType MoveDocumentType,
	title string,
	notes *string,
	requestedAmountCents unit.Cents,
	paymentMethod string,
	movingExpenseType MovingExpenseType,
	moveType SelectedMoveType) (*MovingExpenseDocument, *validate.Errors, error) {

	var newMovingExpenseDocument *MovingExpenseDocument
	var responseError error
	responseVErrors := validate.NewErrors()

	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		var newMoveDocument *MoveDocument
		newMoveDocument, responseVErrors, responseError = m.createMoveDocumentWithoutTransaction(
			db,
			uploads,
			personallyProcuredMoveID,
			moveDocumentType,
			title,
			notes,
			moveType)
		if responseVErrors.HasAny() || responseError != nil {
			return transactionError
		}

		// Finally, create the MovingExpenseDocument
		newMovingExpenseDocument = &MovingExpenseDocument{
			MoveDocumentID:       newMoveDocument.ID,
			MoveDocument:         *newMoveDocument,
			MovingExpenseType:    movingExpenseType,
			RequestedAmountCents: requestedAmountCents,
			PaymentMethod:        paymentMethod,
		}
		verrs, err := db.ValidateAndCreate(newMovingExpenseDocument)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error creating moving expense document")
			newMovingExpenseDocument = nil
			return transactionError
		}

		return nil

	})

	return newMovingExpenseDocument, responseVErrors, responseError
}

// CreatePPM creates a new PPM associated with this move
func (m Move) CreatePPM(db *pop.Connection,
	size *internalmessages.TShirtSize,
	weightEstimate *int64,
	originalMoveDate *time.Time,
	pickupPostalCode *string,
	hasAdditionalPostalCode *bool,
	additionalPickupPostalCode *string,
	destinationPostalCode *string,
	hasSit *bool,
	daysInStorage *int64,
	estimatedStorageReimbursement *string,
	hasRequestedAdvance bool,
	advance *Reimbursement) (*PersonallyProcuredMove, *validate.Errors, error) {

	newPPM := PersonallyProcuredMove{
		MoveID:                        m.ID,
		Move:                          m,
		Size:                          size,
		WeightEstimate:                weightEstimate,
		OriginalMoveDate:              originalMoveDate,
		PickupPostalCode:              pickupPostalCode,
		HasAdditionalPostalCode:       hasAdditionalPostalCode,
		AdditionalPickupPostalCode:    additionalPickupPostalCode,
		DestinationPostalCode:         destinationPostalCode,
		HasSit:                        hasSit,
		DaysInStorage:                 daysInStorage,
		Status:                        PPMStatusDRAFT,
		HasRequestedAdvance:           hasRequestedAdvance,
		Advance:                       advance,
		EstimatedStorageReimbursement: estimatedStorageReimbursement,
	}

	verrs, err := SavePersonallyProcuredMove(db, &newPPM)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return &newPPM, verrs, nil
}

// CreateSignedCertification creates a new SignedCertification associated with this move
func (m Move) CreateSignedCertification(db *pop.Connection,
	submittingUserID uuid.UUID,
	certificationText string,
	signature string,
	date time.Time) (*SignedCertification, *validate.Errors, error) {

	newSignedCertification := SignedCertification{
		MoveID:            m.ID,
		SubmittingUserID:  submittingUserID,
		CertificationText: certificationText,
		Signature:         signature,
		Date:              date,
	}

	verrs, err := db.ValidateAndCreate(&newSignedCertification)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return &newSignedCertification, verrs, nil
}

// GetMovesForUserID gets all move models for a given user ID
func GetMovesForUserID(db *pop.Connection, userID uuid.UUID) (Moves, error) {
	var moves Moves
	query := db.Where("user_id = $1", userID)
	err := query.All(&moves)
	return moves, err
}

// GenerateLocator constructs a record locator - a unique 6 character alphanumeric string
func GenerateLocator() string {
	// Get a UUID as a source of (almost certainly) unique bytes
	seed, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	// Scramble them via SHA256 in case UUID has structure
	scrambledBytes := sha256.Sum256(seed.Bytes())
	// Now convert bytes to letters
	locatorRunes := make([]rune, locatorLength)
	for idx := 0; idx < locatorLength; idx++ {
		j := int(scrambledBytes[idx]) % len(locatorLetters)
		locatorRunes[idx] = locatorLetters[j]
	}
	return string(locatorRunes)
}

// createNewMove adds a new Move record into the DB. In the (unlikely) event that we have a clash on Locators we
// retry with a new record locator.
func createNewMove(db *pop.Connection,
	orders Order,
	selectedType *SelectedMoveType) (*Move, *validate.Errors, error) {

	var stringSelectedType SelectedMoveType
	if selectedType != nil {
		stringSelectedType = SelectedMoveType(*selectedType)
	}
	for i := 0; i < maxLocatorAttempts; i++ {
		move := Move{
			Orders:           orders,
			OrdersID:         orders.ID,
			Locator:          GenerateLocator(),
			SelectedMoveType: &stringSelectedType,
			Status:           MoveStatusDRAFT,
		}
		verrs, err := db.ValidateAndCreate(&move)
		if verrs.HasAny() {
			return nil, verrs, nil
		}
		if err != nil {
			if strings.HasPrefix(errors.Cause(err).Error(), uniqueConstraintViolationErrorPrefix) {
				// If we have a collision, try again for maxLocatorAttempts
				continue
			}
			return nil, verrs, err
		}

		return &move, verrs, nil
	}
	// the only way we get here is if we got a unique constraint error maxLocatorAttempts times.
	verrs := validate.NewErrors()
	return nil, verrs, ErrLocatorGeneration
}

// SaveMoveDependencies safely saves a Move status, ppms' advances' statuses, orders statuses,
// and shipment GBLOCs.
func SaveMoveDependencies(db *pop.Connection, move *Move) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		for _, ppm := range move.PersonallyProcuredMoves {
			if ppm.Advance != nil {
				if verrs, err := db.ValidateAndSave(ppm.Advance); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = errors.Wrap(err, "Error Saving Advance")
					return transactionError
				}
			}

			if verrs, err := db.ValidateAndSave(&ppm); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Saving PPM")
				return transactionError
			}
		}

		if move.Status == MoveStatusSUBMITTED {

			// Save Shipment GBLOCs
			orders, err := FetchOrder(db, move.OrdersID)
			if err != nil {
				responseError = errors.Wrap(err, "Error fetching orders")
				return transactionError
			}

			for _, shipment := range move.Shipments {
				serviceMember, err := FetchServiceMember(db, shipment.ServiceMemberID)
				if err != nil {
					responseError = errors.Wrap(err, "Error fetching service member")
					return transactionError
				}

				destinationGbloc, err := getGbloc(db, orders.NewDutyStationID)
				if err != nil {
					responseError = errors.Wrap(err, "Error getting shipment destination GBLOC")
					return transactionError
				}
				shipment.DestinationGBLOC = &destinationGbloc

				sourceGbloc, err := getGbloc(db, *serviceMember.DutyStationID)
				if err != nil {
					responseError = errors.Wrap(err, "Error getting shipment destination GBLOC")
					return transactionError
				}
				shipment.SourceGBLOC = &sourceGbloc

				// Assign a new unique GBL number using source GBLOC
				err = shipment.AssignGBLNumber(db)
				if err != nil {
					responseError = errors.Wrap(err, "Error assigning GBL number for shipment")
					return transactionError
				}

				if verrs, err := db.ValidateAndSave(&shipment); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = errors.Wrap(err, "Error Saving Shipment")
					return transactionError
				}
			}
		}

		if verrs, err := db.ValidateAndSave(&move.Orders); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving Orders")
			return transactionError
		}

		if verrs, err := db.ValidateAndSave(move); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving Move")
			return transactionError
		}
		return nil
	})

	return responseVErrors, responseError
}

func getGbloc(db *pop.Connection, dutyStationID uuid.UUID) (gbloc string, err error) {
	transportationOffice, err := FetchDutyStationTransportationOffice(db, dutyStationID)
	if err != nil {
		return "", errors.Wrap(err, "could not load transportation office for duty station")
	}
	return transportationOffice.Gbloc, nil
}

// FetchMoveForAdvancePaperwork returns a Move with all of the associations required
// to generate the Advance paperwork.
func FetchMoveForAdvancePaperwork(db *pop.Connection, moveID uuid.UUID) (Move, error) {
	var move Move
	if err := db.Q().Eager("Orders.NewDutyStation", "Orders.ServiceMember.BackupContacts", "Orders.ServiceMember.ResidentialAddress", "PersonallyProcuredMoves.Advance").Find(&move, moveID); err != nil {
		return move, errors.Wrap(err, "could not load move")
	}
	return move, nil
}

// FetchMoveForMoveDates returns a Move along with all the associations needed to determine
// the move dates summary information.
func FetchMoveForMoveDates(db *pop.Connection, moveID uuid.UUID) (Move, error) {
	var move Move
	err := db.
		Eager(
			"Orders.ServiceMember.DutyStation.Address",
			"Orders.NewDutyStation.Address",
			"Orders.ServiceMember",
		).
		Find(&move, moveID)

	return move, err
}
