package movedocument

import (
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

type WeightTicketUpdater struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

//Update updates the weight ticket documents
func (wtu WeightTicketUpdater) Update(moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	newType := models.MoveDocumentType(moveDocumentPayload.MoveDocumentType)
	var emptyWeight, fullWeight *unit.Pound
	updatedMoveDoc, returnVerrs, err := wtu.UpdateMoveDocumentStatus(moveDocumentPayload, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "weightticketupdater.update: error updating move document status")
	}
	if moveDocumentPayload.EmptyWeight != nil {
		ew := unit.Pound(*moveDocumentPayload.EmptyWeight)
		emptyWeight = &ew
	}
	if moveDocumentPayload.FullWeight != nil {
		fw := unit.Pound(*moveDocumentPayload.FullWeight)
		fullWeight = &fw
	}
	var weightTicketDate *time.Time
	if moveDocumentPayload.WeightTicketDate != nil {
		weightTicketDate = (*time.Time)(moveDocumentPayload.WeightTicketDate)
	}
	var trailerOwnershipMissing bool
	if moveDocumentPayload.TrailerOwnershipMissing != nil {
		trailerOwnershipMissing = *moveDocumentPayload.TrailerOwnershipMissing
	}
	var title string
	if moveDocumentPayload.Title != nil {
		title = *moveDocumentPayload.Title
	}
	var emptyWeightTicketMissing bool
	if moveDocumentPayload.EmptyWeightTicketMissing != nil {
		emptyWeightTicketMissing = *moveDocumentPayload.EmptyWeightTicketMissing
	}
	var fullWeightTicketMissing bool
	if moveDocumentPayload.EmptyWeightTicketMissing != nil {
		emptyWeightTicketMissing = *moveDocumentPayload.FullWeightTicketMissing
	}
	updatedMoveDoc.Title = title
	updatedMoveDoc.Notes = moveDocumentPayload.Notes
	updatedMoveDoc.MoveDocumentType = newType
	if updatedMoveDoc.WeightTicketSetDocument == nil {
		updatedMoveDoc.WeightTicketSetDocument = &models.WeightTicketSetDocument{
			MoveDocumentID: updatedMoveDoc.ID,
			MoveDocument:   *updatedMoveDoc,
		}
	}
	updatedMoveDoc.WeightTicketSetDocument.EmptyWeight = emptyWeight
	updatedMoveDoc.WeightTicketSetDocument.EmptyWeightTicketMissing = emptyWeightTicketMissing
	updatedMoveDoc.WeightTicketSetDocument.FullWeight = fullWeight
	updatedMoveDoc.WeightTicketSetDocument.FullWeightTicketMissing = fullWeightTicketMissing
	updatedMoveDoc.WeightTicketSetDocument.VehicleNickname = moveDocumentPayload.VehicleNickname
	updatedMoveDoc.WeightTicketSetDocument.VehicleOptions = moveDocumentPayload.VehicleOptions
	updatedMoveDoc.WeightTicketSetDocument.WeightTicketDate = weightTicketDate
	updatedMoveDoc.WeightTicketSetDocument.TrailerOwnershipMissing = trailerOwnershipMissing
	updatedMoveDoc, returnVerrs, err = wtu.updatePPMNetWeight(updatedMoveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "weightticketupdater.update: error updating weight ticket ppm")
	}
	updatedMoveDoc, returnVerrs, err = wtu.updateWeightTicket(updatedMoveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "weightticketupdater.update: error updating weight ticket")
	}
	return moveDoc, returnVerrs, nil
}

func (wtu WeightTicketUpdater) updatePPMNetWeight(moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	// weight tickets require that we save the ppm again to
	// reflect updated net weight derived from the updated weight tickets
	returnVerrs := validate.NewErrors()
	ppm := &moveDoc.PersonallyProcuredMove
	if moveDoc.PersonallyProcuredMoveID == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("weightticketupdater.updateppmnetweight: no PPM loaded for move doc")
	}
	okStatus := models.MoveDocumentStatusOK
	mergedMoveDocuments, err := mergeMoveDocuments(wtu.db, session, ppm.ID, moveDoc, okStatus)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("weightticketupdater.updateppmnetweight: unable to merge move documents")
	}
	var updatedNetWeight unit.Pound
	for _, weightTicket := range mergedMoveDocuments {
		wts := weightTicket.WeightTicketSetDocument
		if wts != nil && wts.EmptyWeight != nil && wts.FullWeight != nil {
			updatedNetWeight += *wts.FullWeight - *wts.EmptyWeight
		}
	}
	ppm.NetWeight = &updatedNetWeight
	return moveDoc, returnVerrs, nil
}

func (wtu WeightTicketUpdater) updateWeightTicket(moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
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
