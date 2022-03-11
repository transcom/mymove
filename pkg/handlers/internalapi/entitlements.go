package internalapi

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	entitlementop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/entitlements"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForEntitlementModel(e models.WeightAllotment) internalmessages.WeightAllotment {
	// Type Conversion
	TotalWeightSelf := int64(e.TotalWeightSelf)
	TotalWeightSelfPlusDependents := int64(e.TotalWeightSelfPlusDependents)
	ProGearWeight := int64(e.ProGearWeight)
	ProGearWeightSpouse := int64(e.ProGearWeightSpouse)

	return internalmessages.WeightAllotment{
		TotalWeightSelf:               &TotalWeightSelf,
		TotalWeightSelfPlusDependents: &TotalWeightSelfPlusDependents,
		ProGearWeight:                 &ProGearWeight,
		ProGearWeightSpouse:           &ProGearWeightSpouse,
	}
}

// IndexEntitlementsHandler indexes entitlements
type IndexEntitlementsHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h IndexEntitlementsHandler) Handle(params entitlementop.IndexEntitlementsParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {
			entitlements := models.AllWeightAllotments()
			payload := make(map[string]internalmessages.WeightAllotment)
			for k, v := range entitlements {
				rank := string(k)
				allotment := payloadForEntitlementModel(v)
				payload[rank] = allotment
			}
			return entitlementop.NewIndexEntitlementsOK().WithPayload(payload)
		})
}

// ValidateEntitlementHandler validates a weight estimate based on entitlement for a PPM move
type ValidateEntitlementHandler struct {
	handlers.HandlerContext
}

// Handle is the handler
func (h ValidateEntitlementHandler) Handle(params entitlementop.ValidateEntitlementParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	moveID, _ := uuid.FromString(params.MoveID.String())

	// Fetch move, orders, serviceMember and PPM
	move, err := models.FetchMove(appCtx.DB(), appCtx.Session(), moveID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	orders, err := models.FetchOrderForUser(appCtx.DB(), appCtx.Session(), move.OrdersID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), orders.ServiceMemberID)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	// Return 404 if there's no PPM or Shipment,  or if there is no Rank
	if (len(move.PersonallyProcuredMoves) < 1) || serviceMember.Rank == nil {
		return entitlementop.NewValidateEntitlementNotFound()
	}
	var weightEstimate int64
	if len(move.PersonallyProcuredMoves) >= 1 {
		// PPMs are in descending order - this is the last one created
		ppm := move.PersonallyProcuredMoves[0]
		if ppm.WeightEstimate != nil {
			weightEstimate = int64(*ppm.WeightEstimate)
		} else {
			weightEstimate = int64(0)
		}

	}

	smEntitlement, err := models.GetEntitlement(*serviceMember.Rank, orders.HasDependents)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	if weightEstimate > int64(smEntitlement) {
		return handlers.ResponseForConflictErrors(appCtx.Logger(), fmt.Errorf("your estimated weight of %s lbs is above your weight entitlement of %s lbs. \n You will only be paid for the weight you move up to your weight entitlement", humanize.Comma(weightEstimate), humanize.Comma(int64(smEntitlement))))
	}

	return entitlementop.NewValidateEntitlementOK()
}
