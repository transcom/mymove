package movedocument

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

type GenericUpdater struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

//Update updates the generic (non-special case) move documents
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
	var recieptMissing bool
	if payload.ReceiptMissing != nil {
		recieptMissing = *payload.ReceiptMissing
	}
	updatedMoveDoc.Title = title
	updatedMoveDoc.Notes = payload.Notes
	updatedMoveDoc.MoveDocumentType = newType
	if newType == models.MoveDocumentTypeEXPENSE {
		if updatedMoveDoc.MovingExpenseDocument == nil {
			updatedMoveDoc.MovingExpenseDocument = &models.MovingExpenseDocument{
				MoveDocumentID: updatedMoveDoc.ID,
				MoveDocument:   *updatedMoveDoc,
			}
		}
		updatedMoveDoc.MovingExpenseDocument.MovingExpenseType = models.MovingExpenseType(payload.MovingExpenseType)
		updatedMoveDoc.MovingExpenseDocument.RequestedAmountCents = unit.Cents(payload.RequestedAmountCents)
		updatedMoveDoc.MovingExpenseDocument.PaymentMethod = payload.PaymentMethod
		updatedMoveDoc.MovingExpenseDocument.ReceiptMissing = recieptMissing
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
	returnVerrs, err = models.SaveMoveDocument(gu.db, updatedMoveDoc, saveExpenseAction, saveWeightTicketAction)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "Update: error updating move document ppm")
	}
	return moveDoc, returnVerrs, nil
}
