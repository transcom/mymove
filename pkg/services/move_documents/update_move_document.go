package movedocument

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
	"log"
	"time"
)

// createForm is a service object to create a form with data
type moveDocument struct {
	db *pop.Connection
}

func NewMoveDocumentUpdater(db *pop.Connection) services.MoveDocumentUpdater {
	return &moveDocument{db: db}
}

func (m moveDocument) Update(params movedocop.UpdateMoveDocumentParams, moveDocId uuid.UUID, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	//payload := params.UpdateMoveDocument
	//newType := models.MoveDocumentType(payload.MoveDocumentType)
	// fetch move document from move id
	originalMoveDocument, err := models.FetchMoveDocument(m.db, session, moveDocId)
	if err != nil {
		return nil, returnVerrs, models.ErrFetchNotFound
	}
	updatedMoveDocument, returnVerrs, err := m.UpdateMoveDocumentStatus(params, originalMoveDocument, session)
	updatedMoveDocument, returnVerrs, err = m.UpdateMoveDocumentPPM(params, updatedMoveDocument, session)
	return m.Commit(params, updatedMoveDocument, session)
}

func (m moveDocument) UpdateMoveDocumentStatus(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	if moveDoc == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("updateMoveDocumentStatus: missing move document")
	}
	payload := params.UpdateMoveDocument
	newStatus := models.MoveDocumentStatus(payload.Status)
	origStatus := moveDoc.Status
	// only modify the move document status if there is a change in status
	if newStatus == origStatus {
		return moveDoc, returnVerrs, nil
	}
	err := moveDoc.AttemptTransition(newStatus)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.Wrap(err, "updateMoveDocumentStatus: transition error")
	}
	return moveDoc, returnVerrs, nil
}

func (m moveDocument) Commit(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	moveDoc.Title = *payload.Title
	moveDoc.Notes = payload.Notes
	moveDoc.MoveDocumentType = newType

	// If we are an expense type, we need to either delete, create, or update a MovingExpenseType
	// depending on which type of document already exists
	if newType == models.MoveDocumentTypeEXPENSE {
		return m.updateMovingExpense(params, moveDoc)
	}
	// If we are a weight ticket set type, we need to either delete, create, or update a WeightTicketSetDocument
	// depending on which type of document already exists
	if newType == models.MoveDocumentTypeWEIGHTTICKETSET {
		return m.updateWeightTicket(params, moveDoc)
	}
	if newType == models.MoveDocumentTypeSHIPMENTSUMMARY {
		verrs, err := models.SaveMoveDocument(m.db, moveDoc, "", "")
		return moveDoc, verrs, err
	}
	return moveDoc, returnVerrs, nil
}

func (m moveDocument) UpdateMoveDocumentPPM(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	newStatus := models.MoveDocumentStatus(payload.Status)
	if moveDoc.PersonallyProcuredMoveID == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("No PPM loaded for approved move doc")
	}
	switch {
	case newType == models.MoveDocumentTypeWEIGHTTICKETSET:
		return m.updatePPMNetWeight(moveDoc, session)
	case newStatus == models.MoveDocumentStatusOK && moveDoc.MoveDocumentType == models.MoveDocumentTypeSHIPMENTSUMMARY:
		return m.updatePPMStatus(moveDoc)
	case newType == models.MoveDocumentTypeSTORAGEEXPENSE:
		return m.updatePPMSIT(moveDoc, session)
	default:
		return moveDoc, returnVerrs, nil
	}
}

func (m moveDocument) updateMovingExpense(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	// We should have a MovingExpenseDocument model
	payload := params.UpdateMoveDocument
	requestedAmt := unit.Cents(payload.RequestedAmountCents)
	paymentMethod := payload.PaymentMethod

	var storageStartDate *time.Time
	if payload.StorageStartDate != nil {
		storageStartDate = handlers.FmtDatePtrToPopPtr(payload.StorageStartDate)
	}
	var storageEndDate *time.Time
	if payload.StorageEndDate != nil {
		storageEndDate = handlers.FmtDatePtrToPopPtr(payload.StorageEndDate)
	}
	if moveDoc.MovingExpenseDocument == nil {
		// But we don't have one, so create it to be saved later
		moveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
			MoveDocumentID:       moveDoc.ID,
			MoveDocument:         *moveDoc,
			MovingExpenseType:    models.MovingExpenseType(payload.MovingExpenseType),
			RequestedAmountCents: requestedAmt,
			PaymentMethod:        paymentMethod,
			StorageStartDate:     storageStartDate,
			StorageEndDate:       storageEndDate,
		}
	} else {
		// We have one already, so update the fields
		moveDoc.MovingExpenseDocument.MovingExpenseType = models.MovingExpenseType(payload.MovingExpenseType)
		moveDoc.MovingExpenseDocument.RequestedAmountCents = requestedAmt
		moveDoc.MovingExpenseDocument.PaymentMethod = paymentMethod
		moveDoc.MovingExpenseDocument.StorageStartDate = storageStartDate
		moveDoc.MovingExpenseDocument.StorageEndDate = storageEndDate
	}
	// TODO pull this out and move to very end maybe inside a transaction....
	var saveWeightTicketAction models.MoveWeightTicketSetDocumentSaveAction
	if moveDoc.MovingExpenseDocument != nil {
		saveWeightTicketAction = models.MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL
	}
	returnVerrs, err := models.SaveMoveDocument(m.db, moveDoc, models.MoveDocumentSaveActionSAVEEXPENSEMODEL, saveWeightTicketAction)
	if err != nil || returnVerrs.HasAny() {
		return &models.MoveDocument{}, returnVerrs, err
	}
	return moveDoc, returnVerrs, nil
}

func (m moveDocument) updateWeightTicket(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	payload := params.UpdateMoveDocument
	var emptyWeight, fullWeight *unit.Pound
	if payload.EmptyWeight != nil {
		ew := unit.Pound(*payload.EmptyWeight)
		emptyWeight = &ew
	}
	if payload.FullWeight != nil {
		fw := unit.Pound(*payload.FullWeight)
		fullWeight = &fw
	}
	var weightTicketDate *time.Time
	if payload.WeightTicketDate != nil {
		weightTicketDate = (*time.Time)(payload.WeightTicketDate)
	}
	var trailerOwnershipMissing bool
	if payload.TrailerOwnershipMissing != nil {
		trailerOwnershipMissing = *payload.TrailerOwnershipMissing
	}

	if moveDoc.WeightTicketSetDocument == nil {
		// create new weight ticket set
		moveDoc.WeightTicketSetDocument = &models.WeightTicketSetDocument{
			MoveDocumentID:          moveDoc.ID,
			MoveDocument:            *moveDoc,
			EmptyWeight:             emptyWeight,
			FullWeight:              fullWeight,
			VehicleNickname:         payload.VehicleNickname,
			VehicleOptions:          payload.VehicleOptions,
			WeightTicketDate:        weightTicketDate,
			TrailerOwnershipMissing: trailerOwnershipMissing,
		}
	} else {
		// update existing weight ticket set
		moveDoc.WeightTicketSetDocument.EmptyWeight = emptyWeight
		moveDoc.WeightTicketSetDocument.FullWeight = fullWeight

	}
	var saveExpenseAction models.MoveExpenseDocumentSaveAction
	if moveDoc.MovingExpenseDocument != nil {
		saveExpenseAction = models.MoveDocumentSaveActionDELETEEXPENSEMODEL
	}
	returnVerrs, err := models.SaveMoveDocument(m.db, moveDoc, saveExpenseAction, models.MoveDocumentSaveActionSAVEWEIGHTTICKETSETMODEL)
	if err != nil || returnVerrs.HasAny() {
		return &models.MoveDocument{}, returnVerrs, err
	}
	return moveDoc, returnVerrs, nil
}

func (m moveDocument) updatePPMSIT(moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	status := moveDoc.Status
	ppm := moveDoc.PersonallyProcuredMove
	origMovDoc, err := models.FetchMoveDocument(m.db, session, moveDoc.ID)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("unable to fetch original move document")
	}
	daysInStorage, err := origMovDoc.MovingExpenseDocument.DaysInStorage()
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("unable to calculate days in storage")
	}
	updatedDaysInSit := *ppm.DaysInStorage - int64(daysInStorage)
	updatedTotalSit := *ppm.TotalSITCost - origMovDoc.MovingExpenseDocument.RequestedAmountCents
	newdaysInStorage, err := moveDoc.MovingExpenseDocument.DaysInStorage()
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("unable to calculate days in storage")
	}
	if status == models.MoveDocumentStatusOK {
		updatedDaysInSit = updatedDaysInSit + int64(newdaysInStorage)
		updatedTotalSit = updatedTotalSit + moveDoc.MovingExpenseDocument.RequestedAmountCents
	}
	ppm.DaysInStorage = &updatedDaysInSit
	ppm.TotalSITCost = &updatedTotalSit
	return moveDoc, returnVerrs, nil
}

func (m moveDocument) updatePPMStatus(moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	// If the status has already been completed
	// (because the document has been toggled between OK and HAS_ISSUE and back)
	// then don't complete it again.
	returnVerrs := validate.NewErrors()
	ppm := &moveDoc.PersonallyProcuredMove
	if ppm.Status != models.PPMStatusCOMPLETED {
		err := ppm.Complete()
		if err != nil {
			return &models.MoveDocument{}, returnVerrs, errors.Wrap(err, "updatePPMStatus: error completing ppm")
		}
	}
	log.Println("updatePPMStatus", ppm.Status)
	return moveDoc, returnVerrs, nil
}

func (m moveDocument) updatePPMNetWeight(moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	// weight tickets require that we save the ppm again to reflect updated net weight derived from the
	// updated weight tickets
	returnVerrs := validate.NewErrors()
	status := moveDoc.Status
	ppm := moveDoc.PersonallyProcuredMove
	origMovDoc, err := models.FetchMoveDocument(m.db, session, moveDoc.ID)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.Wrap(err, "updatePPMNetWeight: error fetching original weight ticket set")
	}
	var origNetWeight unit.Pound
	if ppm.NetWeight != nil {
		origNetWeight = *ppm.NetWeight
	}
	delta := *origMovDoc.WeightTicketSetDocument.FullWeight - *origMovDoc.WeightTicketSetDocument.EmptyWeight
	updatedNetWeight := origNetWeight - delta
	if status == models.MoveDocumentStatusOK {
		updatedNetWeight += *moveDoc.WeightTicketSetDocument.FullWeight - *moveDoc.WeightTicketSetDocument.EmptyWeight
	}
	ppm.NetWeight = &updatedNetWeight
	return moveDoc, returnVerrs, nil
}
