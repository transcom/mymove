package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/app"
	"github.com/transcom/mymove/pkg/auth"
	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/models"
)

// ApproveMoveHandler approves a move via POST /moves/{moveId}/approve
type ApproveMoveHandler HandlerContext

// Handle ... approves a Move from a request payload
func (h ApproveMoveHandler) Handle(params officeop.ApproveMoveParams) middleware.Responder {
	// #nosec User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	// #nosec UUID is pattern matched by swagger and will be ok
	moveID, _ := uuid.FromString(params.MoveID.String())
	reqApp := app.GetAppFromContext(params.HTTPRequest)

	move, err := models.FetchMove(h.db, user, reqApp, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	move.Status = models.MoveStatusAPPROVED

	verrs, err := h.db.ValidateAndUpdate(move)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	movePayload := payloadForMoveModel(move.Orders, *move)
	return officeop.NewApproveMoveOK().WithPayload(&movePayload)
}

// ApprovePPMHandler approves a move via POST /moves/{moveId}/approve
type ApprovePPMHandler HandlerContext

// Handle ... approves a Personally Procured Move from a request payload
func (h ApprovePPMHandler) Handle(params officeop.ApprovePPMParams) middleware.Responder {
	// #nosec User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	// #nosec UUID is pattern matched by swagger and will be ok
	ppmID, _ := uuid.FromString(params.PersonallyProcuredMoveID.String())
	reqApp := app.GetAppFromContext(params.HTTPRequest)

	ppm, err := models.FetchPersonallyProcuredMove(h.db, user, reqApp, ppmID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	ppm.Status = models.PPMStatusAPPROVED

	verrs, err := h.db.ValidateAndUpdate(ppm)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	ppmPayload := payloadForPPMModel(*ppm)
	return officeop.NewApprovePPMOK().WithPayload(&ppmPayload)
}
