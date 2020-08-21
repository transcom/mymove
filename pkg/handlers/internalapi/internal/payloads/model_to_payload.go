package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// Address payload
func Address(address *models.Address) *internalmessages.Address {
	if address == nil {
		return nil
	}
	return &internalmessages.Address{
		ID:             strfmt.UUID(address.ID.String()),
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           &address.City,
		State:          &address.State,
		PostalCode:     &address.PostalCode,
		Country:        address.Country,
	}
}

// MTOAgent payload
func MTOAgent(mtoAgent *models.MTOAgent) *internalmessages.MTOAgent {
	if mtoAgent == nil {
		return nil
	}

	return &internalmessages.MTOAgent{
		AgentType:     internalmessages.MTOAgentType(mtoAgent.MTOAgentType),
		FirstName:     mtoAgent.FirstName,
		LastName:      mtoAgent.LastName,
		Phone:         mtoAgent.Phone,
		Email:         mtoAgent.Email,
		ID:            strfmt.UUID(mtoAgent.ID.String()),
		MtoShipmentID: strfmt.UUID(mtoAgent.MTOShipmentID.String()),
		CreatedAt:     strfmt.DateTime(mtoAgent.CreatedAt),
		UpdatedAt:     strfmt.DateTime(mtoAgent.UpdatedAt),
	}
}

// MTOAgents payload
func MTOAgents(mtoAgents *models.MTOAgents) *internalmessages.MTOAgents {
	if mtoAgents == nil {
		return nil
	}

	agents := make(internalmessages.MTOAgents, len(*mtoAgents))

	for i, m := range *mtoAgents {
		agents[i] = MTOAgent(&m)
	}

	return &agents
}

// MTOShipment payload
func MTOShipment(mtoShipment *models.MTOShipment) *internalmessages.MTOShipment {
	payload := &internalmessages.MTOShipment{
		ID:                 strfmt.UUID(mtoShipment.ID.String()),
		Agents:             *MTOAgents(&mtoShipment.MTOAgents),
		MoveTaskOrderID:    strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:       internalmessages.MTOShipmentType(mtoShipment.ShipmentType),
		CustomerRemarks:    mtoShipment.CustomerRemarks,
		PickupAddress:      Address(mtoShipment.PickupAddress),
		DestinationAddress: Address(mtoShipment.DestinationAddress),
		CreatedAt:          strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:          strfmt.DateTime(mtoShipment.UpdatedAt),
		Status:             internalmessages.MTOShipmentStatus(mtoShipment.Status),
		ETag:               etag.GenerateEtag(mtoShipment.UpdatedAt),
	}

	if mtoShipment.RequestedPickupDate != nil && !mtoShipment.RequestedPickupDate.IsZero() {
		payload.RequestedPickupDate = strfmt.Date(*mtoShipment.RequestedPickupDate)
	}

	if mtoShipment.RequestedDeliveryDate != nil && !mtoShipment.RequestedDeliveryDate.IsZero() {
		payload.RequestedDeliveryDate = strfmt.Date(*mtoShipment.RequestedDeliveryDate)
	}

	return payload
}

// MTOShipments payload
func MTOShipments(mtoShipments *models.MTOShipments) *internalmessages.MTOShipments {
	payload := make(internalmessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		payload[i] = MTOShipment(&m)
	}
	return &payload
}

// InternalServerError describes errors in a standard structure to be returned in the payload.
// If detail is nil, string defaults to "An internal server error has occurred."
func InternalServerError(detail *string, traceID uuid.UUID) *internalmessages.Error {
	payload := internalmessages.Error{
		Title:    handlers.FmtString(handlers.InternalServerErrMessage),
		Detail:   handlers.FmtString(handlers.InternalServerErrDetail),
		Instance: strfmt.UUID(traceID.String()),
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}

// ValidationError describes validation errors from the model or properties
func ValidationError(detail string, instance uuid.UUID, validationErrors *validate.Errors) *internalmessages.ValidationError {
	payload := &internalmessages.ValidationError{
		ClientError: *ClientError(handlers.ValidationErrMessage, detail, instance),
	}
	if validationErrors != nil {
		payload.InvalidFields = handlers.NewValidationErrorListResponse(validationErrors).Errors
	}
	return payload
}

// ClientError describes errors in a standard structure to be returned in the payload
func ClientError(title string, detail string, instance uuid.UUID) *internalmessages.ClientError {
	return &internalmessages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
}
