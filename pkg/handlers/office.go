package handlers

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// ShowAccountingHandler creates a new service member via GET /moves/{moveId}/accounting
type ShowAccountingHandler HandlerContext

// Handle ... creates a new ServiceMember from a request payload
func (h ShowAccountingHandler) Handle(params officeop.ShowAccountingParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	accountingInfo, err := FetchAccountingInfo(h.db, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	return officeop.NewShowAccountingInfoOK().WithPayload(&ppmEstimate)
}

// PatchAccountingHandler patches a move's accounting information via PATCH /moves/{moveId}/accounting
type PatchAccountingHandler HandlerContext

// Handle ... patches a new ServiceMember from a request payload
func (h PatchAccountingHandler) Handle(params officeop.PatchAccountingParams) middleware.Responder {
	// User should always be populated by middleware
	user, _ := auth.GetUser(params.HTTPRequest.Context())
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Validate that this move belongs to the current user
	move, err := models.FetchMove(h.db, user, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	accountingInfo, err := FetchAccountingInfo(h.db, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	payload := params.PatchAccountingParams
	newTAC := payload.Tac
	newDeptIndicator := payload.DeptIndicator

	if newTAC != nil {
		// TODO: Set TAC here
	}

	if newDeptIndicator != nil {
		// TODO: Set DeptIndicator here
	}

	// TODO: Validate and update whatever obj holds these values
	// verrs, err := h.db.ValidateAndUpdate(move)
	// if err != nil || verrs.HasAny() {
	// 	return responseForVErrors(h.logger, verrs, err)
	// }

	accountingInfo = internalmessages.AccountingInfo{
		Tac:           newTAC,
		DeptIndicator: newDeptIndicator,
	}

	return officeop.NewShowAccountingOK().WithPayload(&accountingInfo)
}
