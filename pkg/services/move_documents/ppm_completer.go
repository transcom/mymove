package movedocument

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/models"
)

type PPMCompleter struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

//Update moves ppm status to complete when ssw is uploaded
func (ppmc PPMCompleter) Update(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	moveDoc.Title = *payload.Title
	moveDoc.Notes = payload.Notes
	moveDoc.MoveDocumentType = newType
	updatedMoveDoc, returnVerrs, err := ppmc.UpdateMoveDocumentStatus(params, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "ppmcompleter.update: error updating move document status")
	}
	if moveDoc.PersonallyProcuredMoveID == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("ppmcompleter.update: no PPM loaded for move doc")
	}
	updatedMoveDoc, returnVerrs, err = ppmc.completePPM(updatedMoveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "ppmcompleter.update: error updating move document ppm")
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
			return &models.MoveDocument{}, returnVerrs, errors.Wrap(err, "ppmcompleter.completeppm: error completing ppm")
		}
	}
	return moveDoc, returnVerrs, nil
}
