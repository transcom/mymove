package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForPPMModel(personallyProcuredMove models.PersonallyProcuredMove) internalmessages.PersonallyProcuredMovePayload {
	ppmPayload := internalmessages.PersonallyProcuredMovePayload{
		ID:             fmtUUID(personallyProcuredMove.ID),
		CreatedAt:      fmtDateTime(personallyProcuredMove.CreatedAt),
		UpdatedAt:      fmtDateTime(personallyProcuredMove.UpdatedAt),
		Size:           personallyProcuredMove.Size,
		WeightEstimate: personallyProcuredMove.WeightEstimate,
	}
	return ppmPayload
}

// CreatePersonallyProcuredMoveHandler creates a PPM
type CreatePersonallyProcuredMoveHandler HandlerContext

// Handle is the handler
func (h CreatePersonallyProcuredMoveHandler) Handle(params ppmop.CreatePersonallyProcuredMoveParams) middleware.Responder {
	// get the user, and validate that this move belongs to them.
	var response middleware.Responder
	user, err := models.GetUserFromRequest(h.db, params.HTTPRequest)
	if err != nil {
		h.logger.Error("Bad User: ", zap.Error(err))
		response = ppmop.NewCreatePersonallyProcuredMoveUnauthorized()
	} else {
		moveID, err := uuid.FromString(params.MoveID.String())
		if err != nil {
			response = ppmop.NewCreatePersonallyProcuredMoveBadRequest()
		} else {

			_, err := models.GetMoveForUser(h.db, user.ID, moveID)
			if err != nil {
				switch errMsg := err.Error(); errMsg {
				case models.ModelFetchErrorNotFound:
					response = ppmop.NewCreatePersonallyProcuredMoveNotFound()
				case models.ModelFetchErrorNotAuthorized:
					response = ppmop.NewCreatePersonallyProcuredMoveUnauthorized()
				default:
					h.logger.Error("Unexpected DB error: ", zap.Error(err))
					response = ppmop.NewCreatePersonallyProcuredMoveInternalServerError()
				}
			} else {
				newPersonallyProcuredMove := models.PersonallyProcuredMove{
					MoveID:         moveID,
					Size:           params.CreatePersonallyProcuredMovePayload.Size,
					WeightEstimate: params.CreatePersonallyProcuredMovePayload.WeightEstimate,
				}

				if verrs, err := h.db.ValidateAndCreate(&newPersonallyProcuredMove); err != nil {
					h.logger.Error("DB Insertion", zap.Error(err))
					response = ppmop.NewCreatePersonallyProcuredMoveBadRequest()
				} else if verrs != nil {
					h.logger.Error("We've got verrrers!")
					response = ppmop.NewCreatePersonallyProcuredMoveBadRequest()
				} else {
					ppmPayload := payloadForPPMModel(newPersonallyProcuredMove)
					response = ppmop.NewCreatePersonallyProcuredMoveCreated().WithPayload(&ppmPayload)
				}
			}
		}
	}
	return response
}

// // IndexPersonallyProcuredMoveHandler returns a list of all the PPMs associated with this move.
// type IndexPersonallyProcuredMoveHandler HandlerContext

// func (h IndexPersonallyProcuredMoveHandler) Handle(params ppmop.IndexPersonallyProcuredMoveParams) middleware.Responder {

// }
