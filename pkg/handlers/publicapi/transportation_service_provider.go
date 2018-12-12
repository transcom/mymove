package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	// "github.com/gofrs/uuid"
	// "go.uber.org/zap"

	// "github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	transportationserviceproviderop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/transportation_service_provider"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForTransportationServieAgentModel(t models.TransportationServiceProvider) *apimessages.TransportationServiceProvider {
	transporationServiceProviderPayload := &apimessages.TransportationServiceProvider{
		ID:                       *handlers.FmtUUID(t.ID),
		StandardCarrierAlphaCode: handlers.FmtString(t.StandardCarrierAlphaCode),
		CreatedAt:                strfmt.DateTime(t.CreatedAt),
		UpdatedAt:                strfmt.DateTime(t.UpdatedAt),
		Enrolled:                 t.Enrolled,
		Name:                     t.Name,
		PocGeneralName:           t.PocGeneralName,
		PocGeneralEmail:          t.PocGeneralEmail,
		PocGeneralPhone:          t.PocGeneralPhone,
		PocClaimsName:            t.PocClaimsName,
		PocClaimsEmail:           t.PocClaimsEmail,
		PocClaimsPhone:           t.PocClaimsPhone,
	}
	return transporationServiceProviderPayload
}

// ShipmentGetTSPHandler returns a TSP for a shipment
type ShipmentGetTSPHandler struct {
	handlers.HandlerContext
}

// Handle getting the tsp for a shipment
func (h ShipmentGetTSPHandler) Handle(params transportationserviceproviderop.GetTransportationServiceProviderParams) middleware.Responder {
	return middleware.NotImplemented("operation .getTSP has not yet been implemented")
}
