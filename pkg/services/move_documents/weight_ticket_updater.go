package movedocument

import (
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// WeightTicketUpdater updates weight tickets
type WeightTicketUpdater struct {
	db *pop.Connection
	moveDocumentStatusUpdater
}

// Update updates the weight ticket documents
func (wtu WeightTicketUpdater) Update(moveDocumentPayload *internalmessages.MoveDocumentPayload, moveDoc *models.MoveDocument, session *auth.Session) (*models.MoveDocument, *validate.Errors, error) {
	newType := models.MoveDocumentType(moveDocumentPayload.MoveDocumentType)
	var emptyWeight, fullWeight *unit.Pound
	updatedMoveDoc, returnVerrs, err := wtu.UpdateMoveDocumentStatus(moveDocumentPayload, moveDoc, session)
	if err != nil || returnVerrs.HasAny() {
		return nil, returnVerrs, errors.Wrap(err, "weightticketupdater.update: error updating move document status")
	}

	var vehicleNickname *string
	if moveDocumentPayload.VehicleNickname != nil {
		vehicleNickname = moveDocumentPayload.VehicleNickname
	}

	var vehicleMake *string
	if moveDocumentPayload.VehicleMake != nil {
		vehicleMake = moveDocumentPayload.VehicleMake
	}

	var vehicleModel *string
	if moveDocumentPayload.VehicleModel != nil {
		vehicleModel = moveDocumentPayload.VehicleModel
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
			MoveDocumentID: moveDoc.ID,
			MoveDocument:   *moveDoc,
		}
	}
	updatedMoveDoc.WeightTicketSetDocument.EmptyWeight = emptyWeight
	updatedMoveDoc.WeightTicketSetDocument.EmptyWeightTicketMissing = emptyWeightTicketMissing
	updatedMoveDoc.WeightTicketSetDocument.FullWeight = fullWeight
	updatedMoveDoc.WeightTicketSetDocument.FullWeightTicketMissing = fullWeightTicketMissing
	updatedMoveDoc.WeightTicketSetDocument.VehicleNickname = vehicleNickname
	updatedMoveDoc.WeightTicketSetDocument.VehicleMake = vehicleMake
	updatedMoveDoc.WeightTicketSetDocument.VehicleModel = vehicleModel
	updatedMoveDoc.WeightTicketSetDocument.WeightTicketSetType = models.WeightTicketSetType(*moveDocumentPayload.WeightTicketSetType)
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
	return updatedMoveDoc, returnVerrs, nil
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
	mergedMoveDocuments, err := mergeMoveDocuments(wtu.db, session, ppm.ID, moveDoc, models.MoveDocumentTypeWEIGHTTICKETSET, okStatus)
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
