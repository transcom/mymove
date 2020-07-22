package payloads

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// AddressModel model
func AddressModel(address *internalmessages.Address) *models.Address {
	if address == nil {
		return nil
	}
	return &models.Address{
		ID:             uuid.FromStringOrNil(address.ID.String()),
		StreetAddress1: *address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           *address.City,
		State:          *address.State,
		PostalCode:     *address.PostalCode,
		Country:        address.Country,
	}
}

// MTOAgentModel model
func MTOAgentModel(mtoAgent *internalmessages.MTOAgent) *models.MTOAgent {
	if mtoAgent == nil {
		return nil
	}

	return &models.MTOAgent{
		ID:            uuid.FromStringOrNil(mtoAgent.ID.String()),
		MTOShipmentID: uuid.FromStringOrNil(mtoAgent.MtoShipmentID.String()),
		FirstName:     mtoAgent.FirstName,
		LastName:      mtoAgent.LastName,
		Email:         mtoAgent.Email,
		Phone:         mtoAgent.Phone,
		MTOAgentType:  models.MTOAgentType(mtoAgent.AgentType),
	}
}

// MTOAgentsModel model
func MTOAgentsModel(mtoAgents *internalmessages.MTOAgents) *models.MTOAgents {
	if mtoAgents == nil {
		return nil
	}

	agents := make(models.MTOAgents, len(*mtoAgents))

	for i, m := range *mtoAgents {
		agents[i] = *MTOAgentModel(m)
	}

	return &agents
}

// MTOShipmentModelFromCreate model
func MTOShipmentModelFromCreate(mtoShipment *internalmessages.CreateShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	requestedPickupDate := time.Time(mtoShipment.RequestedPickupDate)
	model := &models.MTOShipment{
		MoveTaskOrderID:     uuid.FromStringOrNil(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:        models.MTOShipmentType(mtoShipment.ShipmentType),
		RequestedPickupDate: &requestedPickupDate,
		CustomerRemarks:     mtoShipment.CustomerRemarks,
	}

	if mtoShipment.PickupAddress != nil {
		model.PickupAddress = AddressModel(mtoShipment.PickupAddress)
	}

	if mtoShipment.DestinationAddress != nil {
		model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)
	}

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	return model
}

// MTOShipmentModel model
func MTOShipmentModel(mtoShipment *internalmessages.MTOShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	model := &models.MTOShipment{
		ID:           uuid.FromStringOrNil(mtoShipment.ID.String()),
		ShipmentType: models.MTOShipmentType(mtoShipment.ShipmentType),
	}

	scheduledPickupDate := time.Time(mtoShipment.ScheduledPickupDate)
	if !scheduledPickupDate.IsZero() {
		model.ScheduledPickupDate = &scheduledPickupDate
	}

	requestedPickupDate := time.Time(mtoShipment.RequestedPickupDate)
	if !requestedPickupDate.IsZero() {
		model.RequestedPickupDate = &requestedPickupDate
	}

	actualPickupDate := time.Time(mtoShipment.ActualPickupDate)
	if !actualPickupDate.IsZero() {
		model.ActualPickupDate = &actualPickupDate
	}

	requiredDeliveryDate := time.Time(mtoShipment.RequiredDeliveryDate)
	if !requiredDeliveryDate.IsZero() {
		model.RequiredDeliveryDate = &requiredDeliveryDate
	}

	if mtoShipment.PickupAddress != nil {
		model.PickupAddress = AddressModel(mtoShipment.PickupAddress)
	}

	if mtoShipment.DestinationAddress != nil {
		model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)
	}

	if mtoShipment.SecondaryPickupAddress != nil {
		model.SecondaryPickupAddress = AddressModel(mtoShipment.SecondaryPickupAddress)
		secondaryPickupAddressID := uuid.FromStringOrNil(mtoShipment.SecondaryPickupAddress.ID.String())
		model.SecondaryPickupAddressID = &secondaryPickupAddressID
	}

	if mtoShipment.SecondaryDeliveryAddress != nil {
		model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)
		secondaryDeliveryAddressID := uuid.FromStringOrNil(mtoShipment.SecondaryDeliveryAddress.ID.String())
		model.SecondaryDeliveryAddressID = &secondaryDeliveryAddressID
	}

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	return model
}
