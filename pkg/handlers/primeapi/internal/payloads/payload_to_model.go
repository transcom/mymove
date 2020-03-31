package payloads

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// AddressModel model
func AddressModel(address *primemessages.Address) *models.Address {
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
func MTOAgentModel(mtoAgent *primemessages.MTOAgent) *models.MTOAgent {
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
func MTOAgentsModel(mtoAgents *primemessages.MTOAgents) *models.MTOAgents {
	if mtoAgents == nil {
		return nil
	}

	agents := make(models.MTOAgents, len(*mtoAgents))

	for i, m := range *mtoAgents {
		agents[i] = *MTOAgentModel(m)
	}

	return &agents
}

// MTOShipmentModel model
func MTOShipmentModel(mtoShipment *primemessages.MTOShipment) *models.MTOShipment {
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

	firstAvailableDeliveryDate := time.Time(mtoShipment.FirstAvailableDeliveryDate)
	if !firstAvailableDeliveryDate.IsZero() {
		model.FirstAvailableDeliveryDate = &firstAvailableDeliveryDate
	}

	requestedPickupDate := time.Time(mtoShipment.RequestedPickupDate)
	if !requestedPickupDate.IsZero() {
		model.RequestedPickupDate = &requestedPickupDate
	}

	actualPickupDate := time.Time(mtoShipment.ActualPickupDate)
	if !actualPickupDate.IsZero() {
		model.ActualPickupDate = &actualPickupDate
	}

	if mtoShipment.PickupAddress != nil {
		model.PickupAddress = AddressModel(mtoShipment.PickupAddress)
	}

	if mtoShipment.DestinationAddress != nil {
		model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)
	}

	if mtoShipment.PrimeActualWeight > 0 {
		actualWeight := unit.Pound(mtoShipment.PrimeActualWeight)
		model.PrimeActualWeight = &actualWeight
	}

	if mtoShipment.PrimeEstimatedWeight > 0 {
		estimatedWeight := unit.Pound(mtoShipment.PrimeEstimatedWeight)
		model.PrimeEstimatedWeight = &estimatedWeight
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

// MTOServiceItemModel model
func MTOServiceItemModel(mtoServiceItem primemessages.MTOServiceItem) *models.MTOServiceItem {
	if mtoServiceItem == nil {
		return nil
	}

	shipmentID := uuid.FromStringOrNil(mtoServiceItem.MtoShipmentID().String())

	// basic service item
	model := &models.MTOServiceItem{
		ID:              uuid.FromStringOrNil(mtoServiceItem.ID().String()),
		MoveTaskOrderID: uuid.FromStringOrNil(mtoServiceItem.MoveTaskOrderID().String()),
		MTOShipmentID:   &shipmentID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// here we initialize more fields below for other service item types. Eg. MTOServiceItemDOFSIT
	switch mtoServiceItem.ModelType() {
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDOFSIT:
		dofsit := mtoServiceItem.(*primemessages.MTOServiceItemDOFSIT)
		model.ReService.Code = models.ReServiceCodeDOFSIT
		model.Reason = dofsit.Reason
		model.PickupPostalCode = dofsit.PickupPostalCode
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating:
		domesticCrating := mtoServiceItem.(*primemessages.MTOServiceItemDomesticCrating)
		// have to get code from payload
		model.ReService.Code = models.ReServiceCode(*domesticCrating.ReServiceCode)
		model.Description = domesticCrating.Description
		model.Dimensions = models.MTOServiceItemDimensions{
			models.MTOServiceItemDimension{
				Type:   models.DimensionTypeItem,
				Length: unit.ThousandthInches(*domesticCrating.Item.Length),
				Height: unit.ThousandthInches(*domesticCrating.Item.Height),
				Width:  unit.ThousandthInches(*domesticCrating.Item.Width),
			},
			models.MTOServiceItemDimension{
				Type:   models.DimensionTypeCrate,
				Length: unit.ThousandthInches(*domesticCrating.Crate.Length),
				Height: unit.ThousandthInches(*domesticCrating.Crate.Height),
				Width:  unit.ThousandthInches(*domesticCrating.Crate.Width),
			},
		}
	default:
		// assume basic service item, take in provided re service code
		basic := mtoServiceItem.(*primemessages.MTOServiceItemBasic)
		model.ReService.Code = models.ReServiceCode(basic.ReServiceCode)
	}

	return model
}
