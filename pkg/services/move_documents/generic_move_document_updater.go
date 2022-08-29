package movedocument

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// GenericUpdater is a genertic document updater
type GenericUpdater struct {
	moveDocumentStatusUpdater
}

// Update updates the generic (non-special case) move documents
func (gu GenericUpdater) Update(appCtx appcontext.AppContext, moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
	if moveDocumentPayload.MoveDocumentType == nil {
		return nil, nil, errors.New("missing required field: MoveDocumentType")
	}

	newType := models.MoveDocumentType(*moveDocumentPayload.MoveDocumentType)
	updatedMoveDoc, returnVerrs, err := gu.UpdateMoveDocumentStatus(appCtx, moveDocumentPayload, moveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "update: error updating move document status")
	}
	var title string
	if moveDocumentPayload.Title != nil {
		title = *moveDocumentPayload.Title
	}
	var receiptMissing bool
	if moveDocumentPayload.ReceiptMissing != nil {
		receiptMissing = *moveDocumentPayload.ReceiptMissing
	}
	updatedMoveDoc.Title = title
	updatedMoveDoc.Notes = moveDocumentPayload.Notes
	updatedMoveDoc.MoveDocumentType = newType
	if newType == models.MoveDocumentTypeEXPENSE {
		if updatedMoveDoc.MovingExpenseDocument == nil {
			updatedMoveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
				MoveDocumentID: moveDoc.ID,
				MoveDocument:   *moveDoc,
			}
		}
		updatedMoveDoc.MovingExpenseDocument.MovingExpenseType = models.MovingExpenseType(moveDocumentPayload.MovingExpenseType)
		updatedMoveDoc.MovingExpenseDocument.RequestedAmountCents = unit.Cents(moveDocumentPayload.RequestedAmountCents)
		updatedMoveDoc.MovingExpenseDocument.PaymentMethod = moveDocumentPayload.PaymentMethod
		updatedMoveDoc.MovingExpenseDocument.ReceiptMissing = receiptMissing
		// Storage expenses have their own updater StorageExpenseUpdater
		updatedMoveDoc.MovingExpenseDocument.StorageStartDate = nil
		updatedMoveDoc.MovingExpenseDocument.StorageEndDate = nil
	}

	var saveExpenseAction models.MoveExpenseDocumentSaveAction
	if newType == models.MoveDocumentTypeEXPENSE {
		saveExpenseAction = models.MoveDocumentSaveActionSAVEEXPENSEMODEL
	}
	if moveDoc.MovingExpenseDocument != nil && newType != models.MoveDocumentTypeEXPENSE {
		saveExpenseAction = models.MoveDocumentSaveActionDELETEEXPENSEMODEL
	}
	var saveWeightTicketAction models.MoveWeightTicketSetDocumentSaveAction
	if moveDoc.WeightTicketSetDocument != nil {
		saveWeightTicketAction = models.MoveDocumentSaveActionDELETEWEIGHTTICKETSETMODEL
	}
	returnVerrs, err = models.SaveMoveDocument(appCtx.DB(), updatedMoveDoc, saveExpenseAction, saveWeightTicketAction)
	if err != nil || returnVerrs.HasAny() {
		return &models.MoveDocument{}, returnVerrs, err
	}
	return moveDoc, returnVerrs, nil
}
