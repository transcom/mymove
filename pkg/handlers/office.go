package handlers

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	officeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/office"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// ShowAccountingHandler creates a new service member via GET /move/{moveId}/accounting
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
}

// PatchAccountingHandler patches a move's accounting information via PATCH /move/{moveId}/accounting
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
}
