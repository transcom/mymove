package movedocument

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

type GenericUpdater struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

//Update updates the generic (non-special case) move documents
func (gu GenericUpdater) Update(moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	newType := models.MoveDocumentType(moveDocumentPayload.MoveDocumentType)
	updatedMoveDoc, returnVerrs, err := gu.UpdateMoveDocumentStatus(moveDocumentPayload, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "update: error updating move document status")
	}
	var title string
	if moveDocumentPayload.Title != nil {
		title = *moveDocumentPayload.Title
	}
	var recieptMissing bool
	if moveDocumentPayload.ReceiptMissing != nil {
		recieptMissing = *moveDocumentPayload.ReceiptMissing
	}
	updatedMoveDoc.Title = title
	updatedMoveDoc.Notes = moveDocumentPayload.Notes
	updatedMoveDoc.MoveDocumentType = newType
	if newType == models.MoveDocumentTypeEXPENSE {
		if updatedMoveDoc.MovingExpenseDocument == nil {
			updatedMoveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
				MoveDocumentID: updatedMoveDoc.ID,
				MoveDocument:   *updatedMoveDoc,
			}
		}
		updatedMoveDoc.MovingExpenseDocument.MovingExpenseType = models.MovingExpenseType(moveDocumentPayload.MovingExpenseType)
		updatedMoveDoc.MovingExpenseDocument.RequestedAmountCents = unit.Cents(moveDocumentPayload.RequestedAmountCents)
		updatedMoveDoc.MovingExpenseDocument.PaymentMethod = moveDocumentPayload.PaymentMethod
		updatedMoveDoc.MovingExpenseDocument.ReceiptMissing = recieptMissing
	}
	updatedMoveDoc.MovingExpenseDocument.StorageStartDate = nil
	updatedMoveDoc.MovingExpenseDocument.StorageEndDate = nil

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
	returnVerrs, err = models.SaveMoveDocument(gu.db, updatedMoveDoc, saveExpenseAction, saveWeightTicketAction)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "Update: error updating move document ppm")
	}
	return moveDoc, returnVerrs, nil
}
