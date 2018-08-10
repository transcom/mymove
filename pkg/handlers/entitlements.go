package handlers

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/models"
)

// ValidateEntitlementHandler validates a weight estimate based on entitlement
type ValidateEntitlementHandler HandlerContext

// Handle is the handler
func (h ValidateEntitlementHandler) Handle(params entitlementop.ValidateEntitlementParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Fetch move, orders, serviceMember and PPM
	move, err := models.FetchMove(h.db, session, moveID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	orders, err := models.FetchOrderForUser(h.db, session, move.OrdersID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	serviceMember, err := models.FetchServiceMemberForUser(h.db, session, orders.ServiceMemberID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	// Return 404 if there's no PPM or Rank
	if len(move.PersonallyProcuredMoves) < 1 || serviceMember.Rank == nil {
		return entitlementop.NewValidateEntitlementNotFound()
	}
	// PPMs are in descending order - this is the last one created
	weightEstimate := *move.PersonallyProcuredMoves[0].WeightEstimate

	smEntitlement := models.GetEntitlement(*serviceMember.Rank, orders.HasDependents, orders.SpouseHasProGear)
	if int(weightEstimate) > smEntitlement {
		return responseForConflictErrors(h.logger, fmt.Errorf("your estimated weight of %s lbs is above your weight entitlement of %s lbs. \n You will only be paid for the weight you move up to your weight entitlement", humanize.Comma(weightEstimate), humanize.Comma(int64(smEntitlement))))
	}

	return entitlementop.NewValidateEntitlementOK()
}
