package movedocument

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// StorageExpenseUpdater updates storage expenses
type StorageExpenseUpdater struct {
	moveDocumentStatusUpdater
}

// Update updates the storage expense documents
func (seu StorageExpenseUpdater) Update(appCtx appcontext.AppContext, moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	if moveDocumentPayload.MoveDocumentType == nil {
		return nil, nil, errors.New("missing required field: MoveDocumentType")
	}
	newType := models.MoveDocumentType(*moveDocumentPayload.MoveDocumentType)
	updatedMoveDoc, returnVerrs, err := seu.UpdateMoveDocumentStatus(appCtx, moveDocumentPayload, moveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "storageexpenseupdater.update: error updating move document status")
	}
	var storageStartDate *time.Time
	if moveDocumentPayload.StorageStartDate != nil {
		storageStartDate = handlers.FmtDatePtrToPopPtr(moveDocumentPayload.StorageStartDate)
	}
	var storageEndDate *time.Time
	if moveDocumentPayload.StorageEndDate != nil {
		storageEndDate = handlers.FmtDatePtrToPopPtr(moveDocumentPayload.StorageEndDate)
	}
	var receiptMissing bool
	if moveDocumentPayload.ReceiptMissing != nil {
		receiptMissing = *moveDocumentPayload.ReceiptMissing
	}
	updatedMoveDoc.Title = *moveDocumentPayload.Title
	updatedMoveDoc.Notes = moveDocumentPayload.Notes
	updatedMoveDoc.MoveDocumentType = newType
	if updatedMoveDoc.MovingExpenseDocument == nil {
		updatedMoveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
			MoveDocumentID: moveDoc.ID,
			MoveDocument:   *moveDoc,
		}
	}
	updatedMoveDoc.MovingExpenseDocument.MovingExpenseType = models.MovingExpenseTypeSTORAGE
	updatedMoveDoc.MovingExpenseDocument.RequestedAmountCents = unit.Cents(moveDocumentPayload.RequestedAmountCents)
	updatedMoveDoc.MovingExpenseDocument.PaymentMethod = moveDocumentPayload.PaymentMethod
	updatedMoveDoc.MovingExpenseDocument.StorageStartDate = storageStartDate
	updatedMoveDoc.MovingExpenseDocument.StorageEndDate = storageEndDate
	updatedMoveDoc.MovingExpenseDocument.ReceiptMissing = receiptMissing

	updatedMoveDoc, returnVerrs, err = seu.updatePPMSIT(appCtx, updatedMoveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "storageexpenseupdater.update: error updating move document ppm")
	}
	updatedMoveDoc, returnVerrs, err = seu.updateMovingExpense(appCtx, updatedMoveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "storageexpenseupdater.update: error updating move document")
	}
	return updatedMoveDoc, returnVerrs, nil
}

func (seu StorageExpenseUpdater) updatePPMSIT(appCtx appcontext.AppContext, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	ppm := &moveDoc.PersonallyProcuredMove
	if moveDoc.PersonallyProcuredMoveID == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("storageexpenseupdater.updateppmsit: no PPM loaded for move doc")
	}
	okStatus := models.MoveDocumentStatusOK
	mergedMoveDocuments, err := mergeMoveDocuments(appCtx, ppm.ID, moveDoc, models.MoveDocumentTypeEXPENSE, okStatus)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("storageexpenseupdater.updateppmsit: unable to merge move documents")
	}
	movingExpenseDocuments := models.FilterMovingExpenseDocuments(mergedMoveDocuments)
	sitExpenses := models.FilterSITExpenses(movingExpenseDocuments)
	var updatedDaysInSit int64
	var updatedTotalSit unit.Cents
	for _, sitExpense := range sitExpenses {
		days, err := sitExpense.DaysInStorage()
		if err == nil {
			updatedDaysInSit += int64(days)
		}
		updatedTotalSit += sitExpense.RequestedAmountCents
	}
	ppm.DaysInStorage = &updatedDaysInSit
	ppm.TotalSITCost = &updatedTotalSit
	return moveDoc, returnVerrs, nil
}

func (seu StorageExpenseUpdater) updateMovingExpense(appCtx appcontext.AppContext, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	var saveWeightTicketAction models.MoveWeightTicketSetDocumentSaveAction
	if moveDoc.WeightTicketSetDocument != nil {
		saveWeightTicketAction = models.MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL
	}
	returnVerrs, err := models.SaveMoveDocument(appCtx.DB(), moveDoc, models.MoveDocumentSaveActionSAVEEXPENSEMODEL, saveWeightTicketAction)
	if err != nil || returnVerrs.HasAny() {
		return &models.MoveDocument{}, returnVerrs, err
	}
	return moveDoc, returnVerrs, nil
}
