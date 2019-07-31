package movedocument

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

type StorageExpenseUpdater struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

//Update updates the storage expense documents
func (seu StorageExpenseUpdater) Update(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	updatedMoveDoc, returnVerrs, err := seu.UpdateMoveDocumentStatus(params, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "storageexpenseupdater.update: error updating move document status")
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
		updatedMoveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
			MoveDocumentID:    moveDoc.ID,
			MoveDocument:      *moveDoc,
			MovingExpenseType: models.MovingExpenseTypeSTORAGE,
		}
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
		return &models.MoveDocument{}, returnVerrs, errors.New("storageexpenseupdater.updateppmsit: no PPM loaded for move doc")
	}
	okStatus := models.MoveDocumentStatusOK
	moveDocuments, err := models.FetchMoveDocuments(seu.db, session, ppm.ID, &okStatus, models.MoveDocumentTypeEXPENSE)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("storageexpenseupdater.updateppmsit: unable to fetch move documents")
	}
	mergedMoveDocuments := mergeMoveDocuments(moveDocuments, *moveDoc)
	sitExpenses := filterSitExpenses(mergedMoveDocuments)
	var updatedDaysInSit int64
	var updatedTotalSit unit.Cents
	for _, v := range sitExpenses {
		days, err := v.DaysInStorage()
		if err == nil {
			updatedDaysInSit += int64(days)
		}
		updatedTotalSit += v.RequestedAmountCents
	}
	ppm.DaysInStorage = &updatedDaysInSit
	ppm.TotalSITCost = &updatedTotalSit
	return moveDoc, returnVerrs, nil
}

func (seu StorageExpenseUpdater) updateMovingExpense(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
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

func filterSitExpenses(moveDocuments models.MoveDocuments) []models.MovingExpenseDocument {
	var sitExpenses []models.MovingExpenseDocument
	for _, moveDocument := range moveDocuments {
		if moveDocument.MovingExpenseDocument != nil &&
			moveDocument.MovingExpenseDocument.MovingExpenseType == models.MovingExpenseTypeSTORAGE {
			sitExpenses = append(sitExpenses, *moveDocument.MovingExpenseDocument)
		}
	}
	return sitExpenses
}
