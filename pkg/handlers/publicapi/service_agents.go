package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
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
		Email:            s.Email,
		PhoneNumber:      s.PhoneNumber,
		FaxNumber:        s.FaxNumber,
		EmailIsPreferred: s.EmailIsPreferred,
		PhoneIsPreferred: s.PhoneIsPreferred,
		Notes:            s.Notes,
		Company:          handlers.FmtString(s.Company),
	}
	return serviceAgentPayload
}

// IndexServiceAgentsHandler returns a list of service agents via GET /shipments/{shipmentId}/service_agents
type IndexServiceAgentsHandler struct {
	handlers.HandlerContext
}

// Handle returns a list of service agents - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h IndexServiceAgentsHandler) Handle(params serviceagentop.IndexServiceAgentsParams) middleware.Responder {
	var serviceAgents []models.ServiceAgent
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	if session.IsTspUser() {
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

		serviceAgents, err = models.FetchServiceAgentsByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipment.ID)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else if session.IsOfficeUser() {
		shipment, err := models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			h.Logger().Error("DB Query", zap.Error(err))
			return serviceagentop.NewIndexServiceAgentsBadRequest()
		}
		serviceAgents, err = models.FetchServiceAgentsOnShipment(h.DB(), shipment.ID)
		if err != nil {
			return handlers.ResponseForError(h.Logger(), err)
		}
	} else {
		return serviceagentop.NewIndexServiceAgentsForbidden()
	}

	serviceAgentPayloadList := make(apimessages.ServiceAgents, len(serviceAgents))
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
		payload.Company,
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
	var serviceAgent *models.ServiceAgent
	var err error

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	serviceAgentID, _ := uuid.FromString(params.ServiceAgentID.String())

	if session.IsTspUser() {
		tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
		if err != nil {
			h.Logger().Error("DB Query", zap.Error(err))
			return serviceagentop.NewPatchServiceAgentForbidden()
		}

		serviceAgent, err = models.FetchServiceAgentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID, serviceAgentID)
		if err != nil {
			h.Logger().Error("DB Query", zap.Error(err))
			return serviceagentop.NewPatchServiceAgentBadRequest()
		}
	} else if session.IsOfficeUser() {
		serviceAgent, err = models.FetchServiceAgentForOffice(h.DB(), shipmentID, serviceAgentID)
		if err != nil {
			h.Logger().Error("DB Query", zap.Error(err))
			return serviceagentop.NewPatchServiceAgentBadRequest()
		}
	} else {
		h.Logger().Error("Non office or TSP user attempted to patch service agent")
		return serviceagentop.NewPatchServiceAgentForbidden()
	}
	// Update the Service Agent
	payload := params.ServiceAgent
	if payload.Company != nil {
		serviceAgent.Company = *payload.Company
	}
	if payload.Email != nil {
		serviceAgent.Email = payload.Email
	}
	if payload.PhoneNumber != nil {
		serviceAgent.PhoneNumber = payload.PhoneNumber
	}
	if payload.Notes != nil {
		serviceAgent.Notes = payload.Notes
	}

	verrs, err := h.DB().ValidateAndSave(serviceAgent)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	serviceAgentPayload := payloadForServiceAgentModel(*serviceAgent)
	return serviceagentop.NewPatchServiceAgentOK().WithPayload(serviceAgentPayload)
}
