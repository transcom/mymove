package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	authctx "github.com/transcom/mymove/pkg/auth/context"
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
	var response middleware.Responder
	userID, ok := authctx.GetUserID(params.HTTPRequest.Context())
	if !ok {
		h.logger.Fatal("No User ID, this should never happen.")
	}
	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		h.logger.Fatal("Invalid MoveID, this should never happen.")
	}

	// Validate that this move belongs to the current user
	moveResult, err := models.GetMoveForUser(h.db, userID, moveID)
	if err != nil {
		h.logger.Error("DB Error checking on move validity", zap.Error(err))
		response = ppmop.NewCreatePersonallyProcuredMoveInternalServerError()
	} else if !moveResult.IsValid() {
		switch errCode := moveResult.ErrorCode(); errCode {
		case models.FetchErrorNotFound: // this won't work yet...
			response = ppmop.NewCreatePersonallyProcuredMoveNotFound()
		case models.FetchErrorForbidden:
			response = ppmop.NewCreatePersonallyProcuredMoveForbidden()
		default:
			h.logger.Fatal("This case statement is no longer exhaustive!")
		}
	} else { // The given move does belong to the current user.
		newPersonallyProcuredMove := models.PersonallyProcuredMove{
			MoveID:         moveID,
			Size:           params.CreatePersonallyProcuredMovePayload.Size,
			WeightEstimate: params.CreatePersonallyProcuredMovePayload.WeightEstimate,
		}

		if verrs, err := h.db.ValidateAndCreate(&newPersonallyProcuredMove); err != nil {
			h.logger.Error("DB Insertion", zap.Error(err))
			response = ppmop.NewCreatePersonallyProcuredMoveBadRequest()
		} else if verrs.HasAny() {
			h.logger.Error("We got verrs!", zap.String("verrs", verrs.String()))
			response = ppmop.NewCreatePersonallyProcuredMoveBadRequest()
		} else {
			ppmPayload := payloadForPPMModel(newPersonallyProcuredMove)
			response = ppmop.NewCreatePersonallyProcuredMoveCreated().WithPayload(&ppmPayload)
		}
	}
	return response
}

// IndexPersonallyProcuredMovesHandler returns a list of all the PPMs associated with this move.
type IndexPersonallyProcuredMovesHandler HandlerContext

// Handle handles the request
func (h IndexPersonallyProcuredMovesHandler) Handle(params ppmop.IndexPersonallyProcuredMovesParams) middleware.Responder {
	var response middleware.Responder
	userID, ok := authctx.GetUserID(params.HTTPRequest.Context())
	if !ok {
		h.logger.Fatal("No User ID, this should never happen.")
	}
	moveID, err := uuid.FromString(params.MoveID.String())
	if err != nil {
		h.logger.Fatal("Invalid MoveID, this should never happen.")
	}

	// Validate that this move belongs to the current user
	moveResult, err := models.GetMoveForUser(h.db, userID, moveID)
	if err != nil {
		h.logger.Error("DB Error checking on move validity", zap.Error(err))
		response = ppmop.NewCreatePersonallyProcuredMoveInternalServerError()
	} else if !moveResult.IsValid() {
		switch errCode := moveResult.ErrorCode(); errCode {
		case models.FetchErrorNotFound:
			response = ppmop.NewCreatePersonallyProcuredMoveNotFound()
		case models.FetchErrorForbidden:
			response = ppmop.NewCreatePersonallyProcuredMoveForbidden()
		default:
			h.logger.Fatal("This case statement is no longer exhaustive!")
		}
	} else { // The given move does belong to the current user.
		ppms, err := models.GetPersonallyProcuredMovesForMoveID(h.db, moveID)
		if err != nil {
			h.logger.Error("DB Error checking on move validity", zap.Error(err))
			response = ppmop.NewCreatePersonallyProcuredMoveInternalServerError()
		} else {
			ppmsPayload := make(internalmessages.IndexPersonallyProcuredMovePayload, len(ppms))
			for i, ppm := range ppms {
				ppmPayload := payloadForPPMModel(ppm)
				ppmsPayload[i] = &ppmPayload
			}
			response = ppmop.NewIndexPersonallyProcuredMovesOK().WithPayload(ppmsPayload)
		}
	}

	return response
}
