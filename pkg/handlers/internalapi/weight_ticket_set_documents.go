package internalapi

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	movedocop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/move_docs"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/unit"
)

func payloadForWeightTicketSetMoveDocumentModel(storer storage.FileStorer, weightTicketSet models.WeightTicketSetDocument) (*internalmessages.MoveDocumentPayload, error) {

	documentPayload, err := payloadForDocumentModel(storer, weightTicketSet.MoveDocument.Document)
	if err != nil {
		return nil, err
	}

	var ppmID *strfmt.UUID
	if weightTicketSet.MoveDocument.PersonallyProcuredMoveID != nil {
		ppmID = handlers.FmtUUID(*weightTicketSet.MoveDocument.PersonallyProcuredMoveID)
	}
	var vehicleNickname *string
	if weightTicketSet.VehicleNickname != nil {
		vehicleNickname = weightTicketSet.VehicleNickname
	}
	var vehicleMake *string
	if weightTicketSet.VehicleMake != nil {
		vehicleMake = weightTicketSet.VehicleMake
	}
	var vehicleModel *string
	if weightTicketSet.VehicleModel != nil {
		vehicleModel = weightTicketSet.VehicleModel
	}
	var emptyWeight *int64
	if weightTicketSet.EmptyWeight != nil {
		ew := int64(*weightTicketSet.EmptyWeight)
		emptyWeight = &ew
	}
	var fullWeight *int64
	if weightTicketSet.FullWeight != nil {
		fw := int64(*weightTicketSet.FullWeight)
		fullWeight = &fw
	}
	var weighTicketDate *strfmt.Date
	if weightTicketSet.WeightTicketDate != nil {
		weighTicketDate = handlers.FmtDate(*weightTicketSet.WeightTicketDate)
	}
	status := internalmessages.MoveDocumentStatus(weightTicketSet.MoveDocument.Status)
	moveDocumentType := internalmessages.MoveDocumentType(weightTicketSet.MoveDocument.MoveDocumentType)
	weightTicketSetType := internalmessages.WeightTicketSetType(weightTicketSet.WeightTicketSetType)
	genericMoveDocumentPayload := internalmessages.MoveDocumentPayload{
		ID:                       handlers.FmtUUID(weightTicketSet.MoveDocument.ID),
		MoveID:                   handlers.FmtUUID(weightTicketSet.MoveDocument.MoveID),
		Document:                 documentPayload,
		Title:                    &weightTicketSet.MoveDocument.Title,
		MoveDocumentType:         &moveDocumentType,
		VehicleNickname:          vehicleNickname,
		VehicleMake:              vehicleMake,
		VehicleModel:             vehicleModel,
		WeightTicketSetType:      &weightTicketSetType,
		PersonallyProcuredMoveID: ppmID,
		EmptyWeight:              emptyWeight,
		EmptyWeightTicketMissing: handlers.FmtBool(weightTicketSet.EmptyWeightTicketMissing),
		FullWeight:               fullWeight,
		FullWeightTicketMissing:  handlers.FmtBool(weightTicketSet.FullWeightTicketMissing),
		TrailerOwnershipMissing:  handlers.FmtBool(weightTicketSet.TrailerOwnershipMissing),
		WeightTicketDate:         weighTicketDate,
		Status:                   &status,
		Notes:                    weightTicketSet.MoveDocument.Notes,
	}

	return &genericMoveDocumentPayload, nil
}

// CreateWeightTicketSetDocumentHandler creates a WeightTicketSetDocument
type CreateWeightTicketSetDocumentHandler struct {
	handlers.HandlerContext
}

// Handle is the handler for CreateWeightTicketSetDocumentHandler
func (h CreateWeightTicketSetDocumentHandler) Handle(params movedocop.CreateWeightTicketDocumentParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			moveID, err := uuid.FromString(params.MoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			// Validate that this move belongs to the current user
			move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}

			payload := params.CreateWeightTicketDocument
			uploadIds := payload.UploadIds
			userUploads := models.UserUploads{}
			for _, id := range uploadIds {
				convertedUploadID := uuid.Must(uuid.FromString(id.String()))
				userUpload, fetchUploadErr := models.FetchUserUploadFromUploadID(appCtx.DB(), appCtx.Session(), convertedUploadID)
				if fetchUploadErr != nil {
					return handlers.ResponseForError(appCtx.Logger(), fetchUploadErr)
				}
				userUploads = append(userUploads, userUpload)
			}

			ppmID := uuid.Must(uuid.FromString(payload.PersonallyProcuredMoveID.String()))

			// Enforce that the ppm's move_id matches our move
			ppm, fetchPPMErr := models.FetchPersonallyProcuredMove(appCtx.DB(), appCtx.Session(), ppmID)
			if fetchPPMErr != nil {
				return handlers.ResponseForError(appCtx.Logger(), fetchPPMErr)
			}
			if ppm.MoveID != moveID {
				return movedocop.NewCreateWeightTicketDocumentBadRequest()
			}

			if (string(*payload.WeightTicketSetType) == string(models.WeightTicketSetTypeCAR)) ||
				(string(*payload.WeightTicketSetType) == string(models.WeightTicketSetTypeCARTRAILER)) {
				if payload.VehicleMake == nil || payload.VehicleModel == nil {
					wtTypeErr := fmt.Errorf("weight ticket set for type %s must have values for vehicle make and model", string(*payload.WeightTicketSetType))
					return handlers.ResponseForCustomErrors(appCtx.Logger(), wtTypeErr, 422)
				}
			}

			if (string(*payload.WeightTicketSetType) == string(models.WeightTicketSetTypeBOXTRUCK)) ||
				(string(*payload.WeightTicketSetType) == string(models.WeightTicketSetTypePROGEAR)) {
				if payload.VehicleNickname == nil {
					wtTypeErr := fmt.Errorf("weight ticket set for type %s must have value for vehicle nickname", string(*payload.WeightTicketSetType))
					return handlers.ResponseForCustomErrors(appCtx.Logger(), wtTypeErr, 422)
				}
			}

			var emptyWeight *unit.Pound
			if payload.EmptyWeight != nil {
				pound := unit.Pound(*payload.EmptyWeight)
				emptyWeight = &pound
			}
			var fullWeight *unit.Pound
			if payload.FullWeight != nil {
				pound := unit.Pound(*payload.FullWeight)
				fullWeight = &pound
			}
			var weighTicketDate *time.Time
			if payload.WeightTicketDate != nil {
				weighTicketDate = (*time.Time)(payload.WeightTicketDate)
			}

			wtsd := models.WeightTicketSetDocument{
				EmptyWeight:              emptyWeight,
				EmptyWeightTicketMissing: *payload.EmptyWeightTicketMissing,
				FullWeight:               fullWeight,
				FullWeightTicketMissing:  *payload.FullWeightTicketMissing,
				VehicleNickname:          payload.VehicleNickname,
				VehicleMake:              payload.VehicleMake,
				VehicleModel:             payload.VehicleModel,
				WeightTicketSetType:      models.WeightTicketSetType(*payload.WeightTicketSetType),
				WeightTicketDate:         weighTicketDate,
				TrailerOwnershipMissing:  *payload.TrailerOwnershipMissing,
			}
			newWeightTicketSetDocument, verrs, err := move.CreateWeightTicketSetDocument(
				appCtx.DB(),
				userUploads,
				&ppmID,
				&wtsd,
				*move.SelectedMoveType,
			)

			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
			}

			newPayload, err := payloadForWeightTicketSetMoveDocumentModel(h.FileStorer(), *newWeightTicketSetDocument)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err)
			}
			return movedocop.NewCreateWeightTicketDocumentOK().WithPayload(newPayload)
		})
}
