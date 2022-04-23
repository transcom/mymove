package internalapi

import (
	"github.com/go-openapi/runtime/middleware"

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
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			entitlements := models.AllWeightAllotments()
			payload := make(map[string]internalmessages.WeightAllotment)
			for k, v := range entitlements {
				rank := string(k)
				allotment := payloadForEntitlementModel(v)
				payload[rank] = allotment
			}
			return entitlementop.NewIndexEntitlementsOK().WithPayload(payload), nil
		})
}
