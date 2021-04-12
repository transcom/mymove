package models

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/random"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/dberr"
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
	// MoveStatusCANCELED captures enum value "CANCELED"
	MoveStatusCANCELED MoveStatus = "CANCELED"
	// MoveStatusAPPROVALSREQUESTED captures enum value "APPROVALS REQUESTED"
	MoveStatusAPPROVALSREQUESTED MoveStatus = "APPROVALS REQUESTED"
	// MoveStatusNeedsServiceCounseling captures enum value "NEEDS SERVICE COUNSELING"
	MoveStatusNeedsServiceCounseling MoveStatus = "NEEDS SERVICE COUNSELING"
	// MoveStatusServiceCounselingCompleted captures enum value "SERVICE COUNSELING COMPLETED"
	MoveStatusServiceCounselingCompleted MoveStatus = "SERVICE COUNSELING COMPLETED"
)

// SelectedMoveType represents the type of move being represented
type SelectedMoveType string

func (s SelectedMoveType) String() string {
	return string(s)
}

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
	SelectedMoveTypeNTS SelectedMoveType = NTSRaw
	// MoveStatusNTS captures enum value "NTS" for Non-Temporary Storage Release
	SelectedMoveTypeNTSR SelectedMoveType = NTSrRaw
	// MoveStatusHHGPPM captures enum value "HHG_PPM" for combination move HHG + PPM
	SelectedMoveTypeHHGPPM SelectedMoveType = "HHG_PPM"
)

const maxLocatorAttempts = 3
const locatorLength = 6

// This set of letters should produce 'non-word' type strings
var locatorLetters = []rune("346789BCDFGHJKMPQRTVWXY")

// Move is an object representing a move
type Move struct {
	ID                           uuid.UUID               `json:"id" db:"id"`
	Locator                      string                  `json:"locator" db:"locator"`
	CreatedAt                    time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time               `json:"updated_at" db:"updated_at"`
	SubmittedAt                  *time.Time              `json:"submitted_at" db:"submitted_at"`
	OrdersID                     uuid.UUID               `json:"orders_id" db:"orders_id"`
	Orders                       Order                   `belongs_to:"orders"`
	SelectedMoveType             *SelectedMoveType       `json:"selected_move_type" db:"selected_move_type"`
	PersonallyProcuredMoves      PersonallyProcuredMoves `has_many:"personally_procured_moves" order_by:"created_at desc"`
	MoveDocuments                MoveDocuments           `has_many:"move_documents" order_by:"created_at desc"`
	Status                       MoveStatus              `json:"status" db:"status"`
	SignedCertifications         SignedCertifications    `has_many:"signed_certifications" order_by:"created_at desc"`
	CancelReason                 *string                 `json:"cancel_reason" db:"cancel_reason"`
	Show                         *bool                   `json:"show" db:"show"`
	AvailableToPrimeAt           *time.Time              `db:"available_to_prime_at"`
	ContractorID                 *uuid.UUID              `db:"contractor_id"`
	Contractor                   *Contractor             `belongs_to:"contractors"`
	PPMEstimatedWeight           *unit.Pound             `db:"ppm_estimated_weight"`
	PPMType                      *string                 `db:"ppm_type"`
	MTOServiceItems              MTOServiceItems         `has_many:"mto_service_items"`
	PaymentRequests              PaymentRequests         `has_many:"payment_requests"`
	MTOShipments                 MTOShipments            `has_many:"mto_shipments"`
	ReferenceID                  *string                 `db:"reference_id"`
	ServiceCounselingCompletedAt *time.Time              `db:"service_counseling_completed_at"`
}

// MoveOptions is used when creating new moves based on parameters
type MoveOptions struct {
	SelectedType *SelectedMoveType
	Show         *bool
}

// Moves is not required by pop and may be deleted
type Moves []Move

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Move) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Locator, Name: "Locator"},
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
	submitDate := swag.Time(time.Now())
	m.SubmittedAt = submitDate

	// Update PPM status too
	for i := range m.PersonallyProcuredMoves {
		ppm := &m.PersonallyProcuredMoves[i]
		err := ppm.Submit(*submitDate)
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

// SendToServiceCounseling sends the move to needs service counseling
func (m *Move) SendToServiceCounseling() error {
	if m.Status != MoveStatusDRAFT {
		return errors.Wrap(ErrInvalidTransition, fmt.Sprintf("Cannot move to NeedsServiceCounseling state when the Move is not in Draft status. Its current status is %s", m.Status))
	}
	m.Status = MoveStatusNeedsServiceCounseling
	submitDate := swag.Time(time.Now())
	m.SubmittedAt = submitDate

	return nil
}

var validStatusesBeforeApproval = []MoveStatus{
	MoveStatusSUBMITTED,
	MoveStatusAPPROVALSREQUESTED,
	MoveStatusServiceCounselingCompleted,
}

func statusSliceContains(statusSlice []MoveStatus, status MoveStatus) bool {
	for _, validStatus := range statusSlice {
		if status == validStatus {
			return true
		}
	}
	return false
}

// Approve approves the Move
func (m *Move) Approve() error {
	if m.approvable() {
		m.Status = MoveStatusAPPROVED
		return nil
	}
	if m.alreadyApproved() {
		return nil
	}
	return errors.Wrap(
		ErrInvalidTransition, fmt.Sprintf(
			"A move can only be approved if it's in one of these states: %q. However, its current status is: %s",
			validStatusesBeforeApproval, m.Status,
		),
	)
}

func (m *Move) alreadyApproved() bool {
	return m.Status == MoveStatusAPPROVED
}

func (m *Move) approvable() bool {
	return statusSliceContains(validStatusesBeforeApproval, m.Status)
}

// SetApprovalsRequested sets the move to approvals requested
func (m *Move) SetApprovalsRequested() error {
	// Do nothing if it's already in the desired state
	if m.Status == MoveStatusAPPROVALSREQUESTED {
		return nil
	}
	if m.Status != MoveStatusAPPROVED {
		return errors.Wrap(ErrInvalidTransition, fmt.Sprintf("The status for the Move with ID %s can only be set to 'Approvals Requested' from the 'Approved' status, but its current status is %s.", m.ID, m.Status))
	}
	m.Status = MoveStatusAPPROVALSREQUESTED
	return nil
}

// Cancel cancels the Move and its associated PPMs
func (m *Move) Cancel(reason string) error {
	// We can cancel any move that isn't already complete.
	if m.Status == MoveStatusCANCELED {
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
		"MTOShipments.MTOAgents",
		"MTOShipments.PickupAddress",
		"MTOShipments.DestinationAddress",
		"SignedCertifications",
		"Orders",
		"MoveDocuments.Document",
	).Where("show = TRUE").Find(&move, id)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	// Eager loading of nested has_many associations is broken
	for i, moveDoc := range move.MoveDocuments {
		err := db.Load(&moveDoc.Document, "UserUploads.Upload")
		if err != nil {
			return nil, err
		}
		move.MoveDocuments[i] = moveDoc
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
	userUploads UserUploads,
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
		UserUploads:     userUploads,
	}
	newDocVerrs, newDocErr := db.ValidateAndCreate(&newDoc)
	if newDocErr != nil || newDocVerrs.HasAny() {
		responseVErrors.Append(newDocVerrs)
		responseError = errors.Wrap(newDocErr, "Error creating document for move document")
		return nil, responseVErrors, responseError
	}

	// Associate uploads to the new document
	for _, upload := range userUploads {
		copyOfUpload := upload // Make copy to avoid implicit memory aliasing of items from a range statement.
		copyOfUpload.DocumentID = &newDoc.ID
		verrs, err := db.ValidateAndUpdate(&copyOfUpload)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error updating upload")
			return nil, responseVErrors, responseError
		}
	}

	var newMoveDocument *MoveDocument
	if moveType == SelectedMoveTypeHHG {
		newMoveDocument = &MoveDocument{
			Move:             m,
			MoveID:           m.ID,
			Document:         newDoc,
			DocumentID:       newDoc.ID,
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

	verrs, err := db.ValidateAndCreate(newMoveDocument)
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
	userUploads UserUploads,
	modelID *uuid.UUID,
	moveDocumentType MoveDocumentType,
	title string,
	notes *string,
	moveType SelectedMoveType) (*MoveDocument, *validate.Errors, error) {

	var newMoveDocument *MoveDocument
	var responseError error
	responseVErrors := validate.NewErrors()

	transactionErr := db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		newMoveDocument, responseVErrors, responseError = m.createMoveDocumentWithoutTransaction(
			db,
			userUploads,
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

	if transactionErr != nil {
		return nil, responseVErrors, responseError
	}

	return newMoveDocument, responseVErrors, responseError
}

// CreateMovingExpenseDocument creates a moving expense document associated to a move and move document
func (m Move) CreateMovingExpenseDocument(
	db *pop.Connection,
	userUploads UserUploads,
	personallyProcuredMoveID *uuid.UUID,
	moveDocumentType MoveDocumentType,
	title string,
	notes *string,
	expenseDocument MovingExpenseDocument,
	moveType SelectedMoveType,
) (*MovingExpenseDocument, *validate.Errors, error) {

	var newMovingExpenseDocument *MovingExpenseDocument
	var responseError error
	responseVErrors := validate.NewErrors()

	transactionErr := db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		var newMoveDocument *MoveDocument
		newMoveDocument, responseVErrors, responseError = m.createMoveDocumentWithoutTransaction(
			db,
			userUploads,
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
			MovingExpenseType:    expenseDocument.MovingExpenseType,
			RequestedAmountCents: expenseDocument.RequestedAmountCents,
			PaymentMethod:        expenseDocument.PaymentMethod,
			ReceiptMissing:       expenseDocument.ReceiptMissing,
			StorageStartDate:     expenseDocument.StorageStartDate,
			StorageEndDate:       expenseDocument.StorageEndDate,
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

	if transactionErr != nil {
		return nil, responseVErrors, transactionErr
	}

	return newMovingExpenseDocument, responseVErrors, responseError
}

// CreatePPM creates a new PPM associated with this move
func (m Move) CreatePPM(db *pop.Connection,
	weightEstimate *unit.Pound,
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
	date time.Time,
	ppmID *uuid.UUID,
	certificationType *SignedCertificationType) (*SignedCertification, *validate.Errors, error) {

	newSignedCertification := SignedCertification{
		MoveID:                   m.ID,
		PersonallyProcuredMoveID: ppmID,
		CertificationType:        certificationType,
		SubmittingUserID:         submittingUserID,
		CertificationText:        certificationText,
		Signature:                signature,
		Date:                     date,
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
	moveOptions MoveOptions) (*Move, *validate.Errors, error) {

	var stringSelectedType SelectedMoveType
	if moveOptions.SelectedType != nil {
		stringSelectedType = SelectedMoveType(*moveOptions.SelectedType)
	}

	show := swag.Bool(true)
	if moveOptions.Show != nil {
		show = moveOptions.Show
	}

	var contractor Contractor
	err := db.Where("contract_number = ?", "HTC111-11-1-1111").First(&contractor)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not find contractor: %w", err)
	}

	referenceID, err := GenerateReferenceID(db)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not generate a unique ReferenceID: %w", err)
	}

	for i := 0; i < maxLocatorAttempts; i++ {
		move := Move{
			Orders:           orders,
			OrdersID:         orders.ID,
			Locator:          GenerateLocator(),
			SelectedMoveType: &stringSelectedType,
			Status:           MoveStatusDRAFT,
			Show:             show,
			ContractorID:     &contractor.ID,
			ReferenceID:      &referenceID,
		}
		verrs, err := db.ValidateAndCreate(&move)
		if verrs.HasAny() {
			return nil, verrs, nil
		}
		if err != nil {
			if dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "moves_locator_idx") {
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

// GenerateReferenceID generates a reference ID for the MTO
func GenerateReferenceID(db *pop.Connection) (string, error) {
	const maxAttempts = 10
	var referenceID string
	var err error
	for i := 0; i < maxAttempts; i++ {
		referenceID, err = generateReferenceIDHelper(db)
		if err == nil {
			return referenceID, nil
		}
	}
	return "", fmt.Errorf("move: failed to generate reference id; %w", err)
}

// GenerateReferenceID creates a random ID for an MTO. Format (xxxx-xxxx) with X being a number 0-9 (ex. 0009-1234. 4321-4444)
func generateReferenceIDHelper(db *pop.Connection) (string, error) {
	min := 0
	max := 10000
	firstNum, err := random.GetRandomIntAddend(min, max)
	if err != nil {
		return "", err
	}

	secondNum, err := random.GetRandomIntAddend(min, max)
	if err != nil {
		return "", err
	}

	newReferenceID := fmt.Sprintf("%04d-%04d", firstNum, secondNum)

	count, err := db.Where(`reference_id= $1`, newReferenceID).Count(&Move{})
	if err != nil {
		return "", err
	} else if count > 0 {
		return "", errors.New("move: reference_id already exists")
	}

	return newReferenceID, nil
}

// SaveMoveDependencies safely saves a Move status, ppms' advances' statuses, orders statuses,
// and shipment GBLOCs.
func SaveMoveDependencies(db *pop.Connection, move *Move) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	transactionErr := db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		for _, ppm := range move.PersonallyProcuredMoves {
			copyOfPpm := ppm // Make copy to avoid implicit memory aliasing of items from a range statement.
			if copyOfPpm.Advance != nil {
				if verrs, err := db.ValidateAndSave(copyOfPpm.Advance); verrs.HasAny() || err != nil {
					responseVErrors.Append(verrs)
					responseError = errors.Wrap(err, "Error Saving Advance")
					return transactionError
				}
			}

			if verrs, err := db.ValidateAndSave(&copyOfPpm); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Saving PPM")
				return transactionError
			}
		}

		order := &move.Orders
		err := db.Load(order, "Moves")
		if err != nil {
			responseError = errors.Wrap(err, "Error Loading Order")
			return transactionError
		}
		if verrs, err := db.ValidateAndSave(order); verrs.HasAny() || err != nil {
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

	if transactionErr != nil {
		return responseVErrors, transactionErr
	}

	return responseVErrors, responseError
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

// FetchMoveByOrderID returns a station for a given id
func FetchMoveByOrderID(db *pop.Connection, orderID uuid.UUID) (Move, error) {
	var move Move
	err := db.Where("orders_id = ?", orderID).First(&move)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Move{}, ErrFetchNotFound
		}
		return Move{}, err
	}
	return move, nil
}

// FetchMoveByMoveID returns a station for a given id
func FetchMoveByMoveID(db *pop.Connection, moveID uuid.UUID) (Move, error) {
	var move Move
	err := db.Q().Find(&move, moveID)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Move{}, ErrFetchNotFound
		}
		return Move{}, err
	}
	return move, nil
}

// IsCanceled returns true if the Move's status is `CANCELED`, false otherwise
func (m Move) IsCanceled() *bool {
	if m.Status == MoveStatusCANCELED {
		return swag.Bool(true)
	}
	return swag.Bool(false)
}
