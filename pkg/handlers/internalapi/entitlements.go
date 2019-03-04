package internalapi

import (
	"fmt"
	"reflect"

	"github.com/dustin/go-humanize"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/honeycombio/beeline-go"

	"github.com/transcom/mymove/pkg/auth"
	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// ValidateEntitlementHandler validates a weight estimate based on entitlement for a PPM move
type ValidateEntitlementHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h ValidateEntitlementHandler) Handle(params entitlementop.ValidateEntitlementParams) middleware.Responder {

	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Fetch move, orders, serviceMember and PPM
	move, err := models.FetchMove(h.DB(), session, moveID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	orders, err := models.FetchOrderForUser(h.DB(), session, move.OrdersID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	serviceMember, err := models.FetchServiceMemberForUser(ctx, h.DB(), session, orders.ServiceMemberID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	// Return 404 if there's no PPM or Rank, or this is an HHG
	// TODO: Handle COMBO moves
	if len(move.PersonallyProcuredMoves) < 1 || len(move.Shipments) >= 1 || serviceMember.Rank == nil {
		return entitlementop.NewValidateEntitlementNotFound()
	}
	// PPMs are in descending order - this is the last one created
	weightEstimate := *move.PersonallyProcuredMoves[0].WeightEstimate

	smEntitlement, err := models.GetEntitlement(*serviceMember.Rank, orders.HasDependents, orders.SpouseHasProGear)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	if int(weightEstimate) > smEntitlement {
		return handlers.ResponseForConflictErrors(h.Logger(), fmt.Errorf("your estimated weight of %s lbs is above your weight entitlement of %s lbs. \n You will only be paid for the weight you move up to your weight entitlement", humanize.Comma(weightEstimate), humanize.Comma(int64(smEntitlement))))
	}

	return entitlementop.NewValidateEntitlementOK()
}
