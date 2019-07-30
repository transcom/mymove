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
	"time"
)

// createForm is a service object to create a form with data
type moveDocumentUpdater struct {
	db                    *pop.Connection
	weightTicketUpdater   Updater
	storageExpenseUpdater Updater
	ppmCompleter          Updater
	genericUpdater        Updater
	moveDocumentStatusUpdater
}

type Updater interface {
	Update(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error)
}

func NewMoveDocumentUpdater(db *pop.Connection) services.MoveDocumentUpdater {
	mdsu := moveDocumentStatusUpdater{}
	wtu := WeightTicketUpdater{db, mdsu}
	su := StorageExpenseUpdater{db, mdsu}
	ppmc := PPMCompleter{db, mdsu}
	gu := GenericUpdater{db, mdsu}
	return &moveDocumentUpdater{
		db:                        db,
		weightTicketUpdater:       wtu,
		storageExpenseUpdater:     su,
		ppmCompleter:              ppmc,
		genericUpdater:            gu,
		moveDocumentStatusUpdater: mdsu,
	}
}

func (m moveDocumentUpdater) Update(params movedocop.UpdateMoveDocumentParams, moveDocId uuid.UUID, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newStatus := models.MoveDocumentStatus(payload.Status)
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	newExpenseType := models.MovingExpenseType(payload.MovingExpenseType)
	originalMoveDocument, err := models.FetchMoveDocument(m.db, session, moveDocId)
	if err != nil {
		return nil, returnVerrs, models.ErrFetchNotFound
	}
	switch {
	case newType == models.MoveDocumentTypeEXPENSE && newExpenseType == models.MovingExpenseTypeSTORAGE:
		return m.storageExpenseUpdater.Update(params, originalMoveDocument, session)
	case newType == models.MoveDocumentTypeWEIGHTTICKETSET:
		return m.weightTicketUpdater.Update(params, originalMoveDocument, session)
	case newType == models.MoveDocumentTypeSHIPMENTSUMMARY && newStatus == models.MoveDocumentStatusOK:
		return m.ppmCompleter.Update(params, originalMoveDocument, session)
	default:
		return m.genericUpdater.Update(params, originalMoveDocument, session)
	}
}

type moveDocumentStatusUpdater struct {
}

func (mds moveDocumentStatusUpdater) UpdateMoveDocumentStatus(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newStatus := models.MoveDocumentStatus(payload.Status)
	if moveDoc == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("updateMoveDocumentStatus: missing move document")
	}
	origStatus := moveDoc.Status
	// only modify the move document status if its changed
	if newStatus == origStatus {
		return moveDoc, returnVerrs, nil
	}
	err := moveDoc.AttemptTransition(newStatus)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.Wrap(err, "updateMoveDocumentStatus: transition error")
	}
	return moveDoc, returnVerrs, nil
}

type GenericUpdater struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

func (gu GenericUpdater) Update(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	updatedMoveDoc, returnVerrs, err := gu.UpdateMoveDocumentStatus(params, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "update: error updating move document status")
	}
	var title string
	if payload.Title != nil {
		title = *payload.Title
	}
	updatedMoveDoc.Title = title
	updatedMoveDoc.Notes = payload.Notes
	updatedMoveDoc.MoveDocumentType = newType

	var saveExpenseAction models.MoveExpenseDocumentSaveAction
	if moveDoc.MovingExpenseDocument != nil {
		saveExpenseAction = models.MoveDocumentSaveActionDELETEEXPENSEMODEL
	}
	var saveWeightTicketAction models.MoveWeightTicketSetDocumentSaveAction
	if moveDoc.WeightTicketSetDocument != nil {
		saveWeightTicketAction = models.MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL
	}
	returnVerrs, err = models.SaveMoveDocument(gu.db, updatedMoveDoc, saveExpenseAction, saveWeightTicketAction)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "Update: error updating move document ppm")
	}
	return moveDoc, returnVerrs, nil
}

type WeightTicketUpdater struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

func (wtu WeightTicketUpdater) Update(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	var emptyWeight, fullWeight *unit.Pound
	updatedMoveDoc, returnVerrs, err := wtu.UpdateMoveDocumentStatus(params, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "update: error updating move document status")
	}
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
	var title string
	if payload.Title != nil {
		title = *payload.Title
	}
	var emptyWeightTicketMissing bool
	if payload.EmptyWeightTicketMissing != nil {
		emptyWeightTicketMissing = *payload.EmptyWeightTicketMissing
	}
	var fullWeightTicketMissing bool
	if payload.EmptyWeightTicketMissing != nil {
		emptyWeightTicketMissing = *payload.FullWeightTicketMissing
	}
	updatedMoveDoc.Title = title
	updatedMoveDoc.Notes = payload.Notes
	updatedMoveDoc.MoveDocumentType = newType
	if updatedMoveDoc.WeightTicketSetDocument == nil {
		updatedMoveDoc.WeightTicketSetDocument = &models.WeightTicketSetDocument{}
	}
	updatedMoveDoc.WeightTicketSetDocument.EmptyWeight = emptyWeight
	updatedMoveDoc.WeightTicketSetDocument.EmptyWeightTicketMissing = emptyWeightTicketMissing
	updatedMoveDoc.WeightTicketSetDocument.FullWeight = fullWeight
	updatedMoveDoc.WeightTicketSetDocument.FullWeightTicketMissing = fullWeightTicketMissing
	updatedMoveDoc.WeightTicketSetDocument.VehicleNickname = payload.VehicleNickname
	updatedMoveDoc.WeightTicketSetDocument.VehicleOptions = payload.VehicleOptions
	updatedMoveDoc.WeightTicketSetDocument.WeightTicketDate = weightTicketDate
	updatedMoveDoc.WeightTicketSetDocument.TrailerOwnershipMissing = trailerOwnershipMissing
	updatedMoveDoc, returnVerrs, err = wtu.updatePPMNetWeight(params, updatedMoveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "update: error updating weight ticket ppm")
	}
	updatedMoveDoc, returnVerrs, err = wtu.updateWeightTicket(params, updatedMoveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "update: error updating weight ticket")
	}
	return moveDoc, returnVerrs, nil
}

func (wtu WeightTicketUpdater) updatePPMNetWeight(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	// weight tickets require that we save the ppm again to
	// reflect updated net weight derived from the updated weight tickets
	returnVerrs := validate.NewErrors()
	status := moveDoc.Status
	ppm := &moveDoc.PersonallyProcuredMove
	if moveDoc.PersonallyProcuredMoveID == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("updatePPMNetWeight: no PPM loaded for move doc")
	}
	origMovDoc, err := models.FetchMoveDocument(wtu.db, session, moveDoc.ID)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.Wrap(err, "updatePPMNetWeight: error fetching original weight ticket set")
	}
	var origNetWeight unit.Pound
	if ppm.NetWeight != nil {
		origNetWeight = *ppm.NetWeight
	}
	var origFullWeight, origEmptyWeight unit.Pound
	if origMovDoc.WeightTicketSetDocument != nil && origMovDoc.WeightTicketSetDocument.FullWeight != nil {
		origFullWeight = *origMovDoc.WeightTicketSetDocument.FullWeight
	}
	if origMovDoc.WeightTicketSetDocument != nil && origMovDoc.WeightTicketSetDocument.EmptyWeight != nil {
		origEmptyWeight = *origMovDoc.WeightTicketSetDocument.EmptyWeight
	}
	var updatedFullWeight, updatedEmptyWeight unit.Pound
	if moveDoc.WeightTicketSetDocument.FullWeight != nil {
		updatedFullWeight = *moveDoc.WeightTicketSetDocument.FullWeight
	}
	if moveDoc.WeightTicketSetDocument.EmptyWeight != nil {
		updatedEmptyWeight = *moveDoc.WeightTicketSetDocument.EmptyWeight
	}
	prevWeightTicketNetWeight := origFullWeight - origEmptyWeight
	updatedNetWeight := origNetWeight - prevWeightTicketNetWeight
	if status == models.MoveDocumentStatusOK {
		updatedNetWeight += updatedFullWeight - updatedEmptyWeight
	}
	ppm.NetWeight = &updatedNetWeight
	return moveDoc, returnVerrs, nil
}

func (wtu WeightTicketUpdater) updateWeightTicket(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	if moveDoc.WeightTicketSetDocument == nil {
		// create new weight ticket set
		moveDoc.WeightTicketSetDocument = &models.WeightTicketSetDocument{
			ID:                       uuid.UUID{},
			MoveDocumentID:           moveDoc.ID,
			MoveDocument:             *moveDoc,
			EmptyWeight:              moveDoc.WeightTicketSetDocument.EmptyWeight,
			EmptyWeightTicketMissing: moveDoc.WeightTicketSetDocument.EmptyWeightTicketMissing,
			FullWeight:               moveDoc.WeightTicketSetDocument.FullWeight,
			FullWeightTicketMissing:  moveDoc.WeightTicketSetDocument.FullWeightTicketMissing,
			VehicleNickname:          moveDoc.WeightTicketSetDocument.VehicleNickname,
			VehicleOptions:           moveDoc.WeightTicketSetDocument.VehicleOptions,
			WeightTicketDate:         moveDoc.WeightTicketSetDocument.WeightTicketDate,
			TrailerOwnershipMissing:  moveDoc.WeightTicketSetDocument.TrailerOwnershipMissing,
		}
	}
	var saveExpenseAction models.MoveExpenseDocumentSaveAction
	if moveDoc.MovingExpenseDocument != nil {
		saveExpenseAction = models.MoveDocumentSaveActionDELETEEXPENSEMODEL
	}
	returnVerrs, err := models.SaveMoveDocument(wtu.db, moveDoc, saveExpenseAction, models.MoveDocumentSaveActionSAVEWEIGHTTICKETSETMODEL)
	if err != nil || returnVerrs.HasAny() {
		return &models.MoveDocument{}, returnVerrs, err
	}
	return moveDoc, returnVerrs, nil
}

type StorageExpenseUpdater struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

func (seu StorageExpenseUpdater) Update(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	updatedMoveDoc, returnVerrs, err := seu.UpdateMoveDocumentStatus(params, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "update: error updating move document status")
	}
	var storageStartDate *time.Time
	if payload.StorageStartDate != nil {
		storageStartDate = handlers.FmtDatePtrToPopPtr(payload.StorageStartDate)
	}
	var storageEndDate *time.Time
	if payload.StorageEndDate != nil {
		storageEndDate = handlers.FmtDatePtrToPopPtr(payload.StorageEndDate)
	}
	var recieptMissing bool
	if payload.ReceiptMissing != nil {
		recieptMissing = *payload.ReceiptMissing
	}
	updatedMoveDoc.Title = *payload.Title
	updatedMoveDoc.Notes = payload.Notes
	updatedMoveDoc.MoveDocumentType = newType
	if updatedMoveDoc.MovingExpenseDocument == nil {
		updatedMoveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{}
	}
	updatedMoveDoc.MovingExpenseDocument.RequestedAmountCents = unit.Cents(payload.RequestedAmountCents)
	updatedMoveDoc.MovingExpenseDocument.PaymentMethod = payload.PaymentMethod
	updatedMoveDoc.MovingExpenseDocument.StorageStartDate = storageStartDate
	updatedMoveDoc.MovingExpenseDocument.StorageEndDate = storageEndDate
	updatedMoveDoc.MovingExpenseDocument.ReceiptMissing = recieptMissing

	updatedMoveDoc, returnVerrs, err = seu.updatePPMSIT(updatedMoveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "Update: error updating move document ppm")
	}
	updatedMoveDoc, returnVerrs, err = seu.updateMovingExpense(params, updatedMoveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "Update: error updating move document")
	}
	return updatedMoveDoc, returnVerrs, nil
}

func (seu StorageExpenseUpdater) updatePPMSIT(moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	ppm := &moveDoc.PersonallyProcuredMove
	if moveDoc.PersonallyProcuredMoveID == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("updatePPMSIT: no PPM loaded for move doc")
	}
	origMovDoc, err := models.FetchMoveDocument(seu.db, session, moveDoc.ID)
	moveDocuments, err := models.FetchMoveDocuments(seu.db, session, ppm.ID, nil, models.MoveDocumentTypeEXPENSE)
	//sitExpenses := fetchSITExpenses(moveDocuments)
	// TODO This is not going to work since users can add there own expenses and override ours
	// TODO Need to get all docs and recalc like ssw
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("updatePPMSIT: unable to fetch original move document")
	}
	var daysInStorage int
	if origMovDoc.MovingExpenseDocument != nil &&
		origMovDoc.MovingExpenseDocument.StorageStartDate != nil &&
		origMovDoc.MovingExpenseDocument.StorageEndDate != nil {
		if daysInStorage, err = origMovDoc.MovingExpenseDocument.DaysInStorage(); err != nil {
			return &models.MoveDocument{}, returnVerrs, errors.New("updatePPMSIT: unable to calculate days in storage")
		}
	}
	var prevDaysInStorage int64
	if ppm.DaysInStorage != nil {
		prevDaysInStorage = *ppm.DaysInStorage
	}
	var prevSitCost unit.Cents
	if ppm.TotalSITCost != nil {
		prevSitCost = *ppm.TotalSITCost
	}
	updatedDaysInSit := prevDaysInStorage - int64(daysInStorage)
	updatedTotalSit := prevSitCost - origMovDoc.MovingExpenseDocument.RequestedAmountCents
	newdaysInStorage, err := moveDoc.MovingExpenseDocument.DaysInStorage()
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("updatePPMSIT: unable to calculate days in storage")
	}
	if moveDoc.Status == models.MoveDocumentStatusOK {
		updatedDaysInSit = updatedDaysInSit + int64(newdaysInStorage)
		updatedTotalSit = updatedTotalSit + moveDoc.MovingExpenseDocument.RequestedAmountCents
	}
	ppm.DaysInStorage = &updatedDaysInSit
	ppm.TotalSITCost = &updatedTotalSit
	return moveDoc, returnVerrs, nil
}

// FetchMovingExpensesShipmentSummaryWorksheet fetches moving expenses for the Shipment Summary Worksheet
// TODO add merge method for new expense and old one
func fetchSITExpenses(move models.Move, db *pop.Connection, session *auth.Session) ([]models.MovingExpenseDocument, error) {
	var sitExpenses []models.MovingExpenseDocument
	if len(move.PersonallyProcuredMoves) > 0 {
		ppm := move.PersonallyProcuredMoves[0]
		moveDocuments, err := models.FetchMoveDocuments(db, session, ppm.ID, nil, models.MoveDocumentTypeEXPENSE)
		if err != nil {
			return sitExpenses, err
		}
		sitExpenses, err = models.FetchMovingExpenses(moveDocuments)
		if err != nil {
			return sitExpenses, err
		}
		sitExpenses = models.FilterSITExpenses(sitExpenses)
	}
	return sitExpenses, nil
}

func (seu StorageExpenseUpdater) updateMovingExpense(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument

	if moveDoc.MovingExpenseDocument == nil {
		// But we don't have one, so create it to be saved later
		moveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
			MoveDocumentID:       moveDoc.ID,
			MoveDocument:         *moveDoc,
			MovingExpenseType:    models.MovingExpenseType(payload.MovingExpenseType),
			RequestedAmountCents: moveDoc.MovingExpenseDocument.RequestedAmountCents,
			PaymentMethod:        moveDoc.MovingExpenseDocument.PaymentMethod,
			ReceiptMissing:       moveDoc.MovingExpenseDocument.ReceiptMissing,
			StorageStartDate:     moveDoc.MovingExpenseDocument.StorageStartDate,
			StorageEndDate:       moveDoc.MovingExpenseDocument.StorageEndDate,
		}
	}
	//TODO check w/ ren that got this right, but understanding as if move document wasn't nil
	//TODO i.e. we're in update situation, want to clear weight ticket since this is an expense
	//TODO not sure how this situation would arise, but would be like prev was a wt, but now an expense
	var saveWeightTicketAction models.MoveWeightTicketSetDocumentSaveAction
	if moveDoc.WeightTicketSetDocument != nil {
		saveWeightTicketAction = models.MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL
	}
	returnVerrs, err := models.SaveMoveDocument(seu.db, moveDoc, models.MoveDocumentSaveActionSAVEEXPENSEMODEL, saveWeightTicketAction)
	if err != nil || returnVerrs.HasAny() {
		return &models.MoveDocument{}, returnVerrs, err
	}
	return moveDoc, returnVerrs, nil
}

type PPMCompleter struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

func (ppmc PPMCompleter) Update(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	moveDoc.Title = *payload.Title
	moveDoc.Notes = payload.Notes
	moveDoc.MoveDocumentType = newType
	updatedMoveDoc, returnVerrs, err := ppmc.UpdateMoveDocumentStatus(params, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "update: error updating move document status")
	}
	if moveDoc.PersonallyProcuredMoveID == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("updatePPMNetWeight: no PPM loaded for move doc")
	}
	updatedMoveDoc, returnVerrs, err = ppmc.completePPM(updatedMoveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "Update: error updating move document ppm")
	}
	var saveExpenseAction models.MoveExpenseDocumentSaveAction
	if moveDoc.MovingExpenseDocument != nil {
		saveExpenseAction = models.MoveDocumentSaveActionDELETEEXPENSEMODEL
	}
	var saveWeightTicketAction models.MoveWeightTicketSetDocumentSaveAction
	if moveDoc.WeightTicketSetDocument != nil {
		saveWeightTicketAction = models.MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL
	}
	returnVerrs, err = models.SaveMoveDocument(ppmc.db, updatedMoveDoc, saveExpenseAction, saveWeightTicketAction)
	return moveDoc, returnVerrs, err
}

func (ppmc PPMCompleter) completePPM(moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	// If the status has already been completed
	// (because the document has been toggled between OK and HAS_ISSUE and back)
	// then don't complete it again.
	returnVerrs := validate.NewErrors()
	ppm := &moveDoc.PersonallyProcuredMove
	if ppm.Status != models.PPMStatusCOMPLETED {
		err := ppm.Complete()
		if err != nil {
			return &models.MoveDocument{}, returnVerrs, errors.Wrap(err, "completePPM: error completing ppm")
		}
	}
	return moveDoc, returnVerrs, nil
}
