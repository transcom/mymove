package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	serviceagentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/service_agents"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForServiceAgentModel(s models.ServiceAgent) *apimessages.ServiceAgent {
	serviceAgentPayload := &apimessages.ServiceAgent{
		ID:               *handlers.FmtUUID(s.ID),
		CreatedAt:        strfmt.DateTime(s.CreatedAt),
		UpdatedAt:        strfmt.DateTime(s.UpdatedAt),
		Role:             apimessages.ServiceAgentRole(s.Role),
		PointOfContact:   handlers.FmtString(s.PointOfContact),
		Email:            s.Email,
		PhoneNumber:      s.PhoneNumber,
		FaxNumber:        s.FaxNumber,
		EmailIsPreferred: s.EmailIsPreferred,
		PhoneIsPreferred: s.PhoneIsPreferred,
		Notes:            s.Notes,
	}
	return serviceAgentPayload
}

// CreateServiceAgentHandler ... creates a new service agent on a shipment via POST /shipment/{shipmentId}/serviceAgent
type CreateServiceAgentHandler struct {
	handlers.HandlerContext
}

// Handle ... creates a new ServiceAgent from a request payload - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CreateServiceAgentHandler) Handle(params serviceagentop.CreateServiceAgentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO (rebecca): Find a way to check Shipment belongs to TSP without 2 queries
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return serviceagentop.NewCreateServiceAgentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return serviceagentop.NewCreateServiceAgentBadRequest()
	}

	payload := params.ServiceAgent

	serviceAgentRole := models.Role(payload.Role)
	newServiceAgent, verrs, err := models.CreateServiceAgent(
		h.DB(),
		shipment.ID,
		serviceAgentRole,
		payload.PointOfContact,
		payload.Email,
		payload.PhoneNumber,
		payload.EmailIsPreferred,
		payload.PhoneIsPreferred,
		payload.Notes)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	serviceAgentPayload := payloadForServiceAgentModel(newServiceAgent)
	return serviceagentop.NewCreateServiceAgentOK().WithPayload(serviceAgentPayload)
}
