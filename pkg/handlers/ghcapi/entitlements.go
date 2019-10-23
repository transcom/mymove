package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	entitlementscodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/entitlements"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

// GetEntitlementsHandler fetches the entitlements for a move task order
type GetEntitlementsHandler struct {
	handlers.HandlerContext
}

// Handle getting the entitlements for a move task order
func (h GetEntitlementsHandler) Handle(params entitlementscodeop.GetEntitlementsParams) middleware.Responder {
	// for now just return static data
	entitlements := &ghcmessages.Entitlements{
		ID:                    "571008b1-b0de-454d-b843-d71be9f02c04",
		DependentsAuthorized:  false,
		NonTemporaryStorage:   false,
		PrivatelyOwnedVehicle: true,
		ProGearWeight:         200,
		ProGearWeightSpouse:   100,
		StorageInTransit:      90,
		TotalDependents:       3,
		TotalWeightSelf:       1300,
	}
	return entitlementscodeop.NewGetEntitlementsOK().WithPayload(entitlements)
}
