package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	publicserviceagentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/service_agents"
	"github.com/transcom/mymove/pkg/models"
)

/*
 * ------------------------------------------
 * The code below is for the PUBLIC REST API.
 * ------------------------------------------
 */

func publicPayloadForServiceAgentModel(s models.ServiceAgent) *apimessages.ServiceAgent {
	serviceAgentPayload := &apimessages.ServiceAgent{
		ID:               *fmtUUID(s.ID),
		ShipmentID:       *fmtUUID(s.ShipmentID),
		CreatedAt:        strfmt.DateTime(s.CreatedAt),
		UpdatedAt:        strfmt.DateTime(s.UpdatedAt),
		Role:             apimessages.ServiceAgentRole(s.Role),
		PointOfContact:   fmtString(s.PointOfContact),
		Email:            s.Email,
		PhoneNumber:      s.PhoneNumber,
		FaxNumber:        s.FaxNumber,
		EmailIsPreferred: s.EmailIsPreferred,
		PhoneIsPreferred: s.PhoneIsPreferred,
		Notes:            s.Notes,
	}
	return serviceAgentPayload
}

// PublicCreateServiceAgentHandler ... creates a new service agent on a shipment via POST /shipment/{shipmentId}/serviceAgent
type PublicCreateServiceAgentHandler HandlerContext

// Handle ... creates a new ServiceAgent from a request payload - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PublicCreateServiceAgentHandler) Handle(params publicserviceagentop.CreateServiceAgentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	shipment, err := models.FetchShipment(h.db, session, shipmentID)
	if err != nil {
		return responseForError(h.logger, err)
	}

	payload := params.ServiceAgent

	serviceAgentRole := models.Role(payload.Role)
	newServiceAgent, verrs, err := models.CreateServiceAgent(h.db, shipment.ID, serviceAgentRole, *payload.PointOfContact, payload.Email, payload.PhoneNumber)
	if err != nil || verrs.HasAny() {
		return responseForVErrors(h.logger, verrs, err)
	}

	serviceAgentPayload := publicPayloadForServiceAgentModel(newServiceAgent)
	return publicserviceagentop.NewCreateServiceAgentOK().WithPayload(serviceAgentPayload)
}
