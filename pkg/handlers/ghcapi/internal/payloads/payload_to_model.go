package payloads

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
)

// AddressModel model
func AddressModel(address *ghcmessages.Address) *models.Address {
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
func MTOAgentModel(mtoAgent *ghcmessages.MTOAgent) *models.MTOAgent {
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
func MTOAgentsModel(mtoAgents *ghcmessages.MTOAgents) *models.MTOAgents {
	if mtoAgents == nil {
		return nil
	}

	agents := make(models.MTOAgents, len(*mtoAgents))

	for i, m := range *mtoAgents {
		agents[i] = *MTOAgentModel(m)
	}

	return &agents
}

// CustomerToServiceMember transforms UpdateCustomerPayload to ServiceMember model
func CustomerToServiceMember(payload ghcmessages.UpdateCustomerPayload) models.ServiceMember {

	var address = models.Address{
		ID:             uuid.FromStringOrNil(payload.CurrentAddress.ID.String()),
		StreetAddress1: *payload.CurrentAddress.StreetAddress1,
		StreetAddress2: payload.CurrentAddress.StreetAddress2,
		StreetAddress3: payload.CurrentAddress.StreetAddress3,
		City:           *payload.CurrentAddress.City,
		State:          *payload.CurrentAddress.State,
		PostalCode:     *payload.CurrentAddress.PostalCode,
		Country:        payload.CurrentAddress.Country,
	}

	var backupContact = models.BackupContact{
		Email: *payload.BackupContact.Email,
		Name:  *payload.BackupContact.Name,
		Phone: payload.BackupContact.Phone,
	}

	var backupContacts []models.BackupContact
	backupContacts = append(backupContacts, backupContact)

	return models.ServiceMember{
		ResidentialAddress: &address,
		BackupContacts:     backupContacts,
		FirstName:          &payload.FirstName,
		LastName:           &payload.LastName,
		PersonalEmail:      payload.Email,
		Telephone:          payload.Phone,
	}
}

// MTOShipmentModelFromUpdate model
func MTOShipmentModelFromUpdate(mtoShipment *ghcmessages.UpdateShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	requestedPickupDate := time.Time(mtoShipment.RequestedPickupDate)
	requestedDeliveryDate := time.Time(mtoShipment.RequestedDeliveryDate)

	model := &models.MTOShipment{
		ShipmentType:          models.MTOShipmentType(mtoShipment.ShipmentType),
		RequestedPickupDate:   &requestedPickupDate,
		RequestedDeliveryDate: &requestedDeliveryDate,
		CustomerRemarks:       mtoShipment.CustomerRemarks,
		Status:                models.MTOShipmentStatus(mtoShipment.Status),
	}

	model.PickupAddress = AddressModel(mtoShipment.PickupAddress)
	model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	return model
}
