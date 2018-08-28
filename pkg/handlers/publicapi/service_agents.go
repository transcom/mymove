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
		ShipmentID:       *handlers.FmtUUID(s.ShipmentID),
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

// IndexServiceAgentsHandler returns a list of service agents via GET /shipments/{shipmentId}/service_agents
type IndexServiceAgentsHandler struct {
	handlers.HandlerContext
}

// Handle returns a list of service agents - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h IndexServiceAgentsHandler) Handle(params serviceagentop.IndexServiceAgentsParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.Logger().Error("Missing TSP User ID")
		return serviceagentop.NewIndexServiceAgentsForbidden()
	}

	// TODO (2018_08_27 cgilmer): Find a way to check Shipment belongs to TSP without 2 queries
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return serviceagentop.NewIndexServiceAgentsForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return serviceagentop.NewIndexServiceAgentsBadRequest()
	}

	serviceAgents, err := models.FetchServiceAgentsByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipment.ID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	serviceAgentPayloadList := make(apimessages.IndexServiceAgents, len(serviceAgents))
	for i, serviceAgent := range serviceAgents {
		serviceAgentPayloadList[i] = payloadForServiceAgentModel(serviceAgent)
	}
	return serviceagentop.NewIndexServiceAgentsOK().WithPayload(serviceAgentPayloadList)
}

// CreateServiceAgentHandler creates a new service agent on a shipment via POST /shipments/{shipmentId}/service_agents
type CreateServiceAgentHandler struct {
	handlers.HandlerContext
}

// Handle creates a new ServiceAgent from a request payload - checks that currently logged in user is authorized to act for the TSP assigned the shipment
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

// PatchServiceAgentHandler allows a user to update a service agent
type PatchServiceAgentHandler struct {
	handlers.HandlerContext
}

// Handle updates the service agent - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PatchServiceAgentHandler) Handle(params serviceagentop.PatchServiceAgentParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	serviceAgentID, _ := uuid.FromString(params.ServiceAgentID.String())

	// Possible they are coming from the wrong endpoint and thus the session is missing the
	// TspUserID
	if session.TspUserID == uuid.Nil {
		h.Logger().Error("Missing TSP User ID")
		return serviceagentop.NewPatchServiceAgentForbidden()
	}

	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return serviceagentop.NewPatchServiceAgentForbidden()
	}

	serviceAgent, err := models.FetchServiceAgentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID, serviceAgentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return serviceagentop.NewPatchServiceAgentBadRequest()
	}

	// Update the Service Agent
	payload := params.Update
	serviceAgent.PointOfContact = *payload.PointOfContact
	serviceAgent.Email = payload.Email
	serviceAgent.PhoneNumber = payload.PhoneNumber
	serviceAgent.Notes = payload.Notes

	verrs, err := h.DB().ValidateAndSave(serviceAgent)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	serviceAgentPayload := payloadForServiceAgentModel(*serviceAgent)
	return serviceagentop.NewPatchServiceAgentOK().WithPayload(serviceAgentPayload)
}
