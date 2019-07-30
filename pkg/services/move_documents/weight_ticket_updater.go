package movedocument

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
	"time"
)

type WeightTicketUpdater struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

//Update updates the weight ticket documents
func (wtu WeightTicketUpdater) Update(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	returnVerrs := validate.NewErrors()
	payload := params.UpdateMoveDocument
	newType := models.MoveDocumentType(payload.MoveDocumentType)
	var emptyWeight, fullWeight *unit.Pound
	updatedMoveDoc, returnVerrs, err := wtu.UpdateMoveDocumentStatus(params, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "weightticketupdater.update: error updating move document status")
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
		return nil, returnVerrs, errors.Wrap(err, "weightticketupdater.update: error updating weight ticket ppm")
	}
	updatedMoveDoc, returnVerrs, err = wtu.updateWeightTicket(params, updatedMoveDoc)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "weightticketupdater.update: error updating weight ticket")
	}
	return moveDoc, returnVerrs, nil
}

func (wtu WeightTicketUpdater) updatePPMNetWeight(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	// weight tickets require that we save the ppm again to
	// reflect updated net weight derived from the updated weight tickets
	returnVerrs := validate.NewErrors()
	ppm := &moveDoc.PersonallyProcuredMove
	if moveDoc.PersonallyProcuredMoveID == nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("weightticketupdater.updateppmnetweight: no PPM loaded for move doc")
	}
	okStatus := models.MoveDocumentStatusOK
	moveDocuments, err := models.FetchMoveDocuments(wtu.db, session, ppm.ID, &okStatus, models.MoveDocumentTypeWEIGHTTICKETSET)
	if err != nil {
		return &models.MoveDocument{}, returnVerrs, errors.New("weightticketupdater.updateppmnetweight: unable to fetch move documents")
	}
	mergedMoveDocuments := mergeMoveDocuments(moveDocuments, *moveDoc)
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

func (wtu WeightTicketUpdater) updateWeightTicket(params movedocop.UpdateMoveDocumentParams, moveDoc *models.MoveDocument) (*models.MoveDocument, *validate.Errors, error) {
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