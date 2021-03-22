package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
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

// MTOServiceItemModelListFromCreate model
func MTOServiceItemModelListFromCreate(mtoShipment *primemessages.CreateMTOShipment) (models.MTOServiceItems, *validate.Errors) {

	if mtoShipment == nil {
		return nil, nil
	}

	serviceItemsListFromPayload := mtoShipment.MtoServiceItems()

	serviceItemsList := make(models.MTOServiceItems, len(serviceItemsListFromPayload))

	for i, m := range serviceItemsListFromPayload {
		serviceItem, verrs := MTOServiceItemModel(m)
		if verrs != nil && verrs.HasAny() {
			return nil, verrs
		}

		serviceItemsList[i] = *serviceItem
	}

	return serviceItemsList, nil
}

// MTOShipmentModelFromCreate model
func MTOShipmentModelFromCreate(mtoShipment *primemessages.CreateMTOShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	model := &models.MTOShipment{
		MoveTaskOrderID: uuid.FromStringOrNil(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:    models.MTOShipmentType(mtoShipment.ShipmentType),
		CustomerRemarks: mtoShipment.CustomerRemarks,
	}

	if mtoShipment.PrimeEstimatedWeight > 0 {
		estimatedWeight := unit.Pound(mtoShipment.PrimeEstimatedWeight)
		model.PrimeEstimatedWeight = &estimatedWeight
		recordedDate := time.Now()
		model.PrimeEstimatedWeightRecordedDate = &recordedDate
	}

	if mtoShipment.RequestedPickupDate != nil {
		model.RequestedPickupDate = swag.Time(time.Time(*mtoShipment.RequestedPickupDate))
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
func MTOServiceItemModel(mtoServiceItem primemessages.MTOServiceItem) (*models.MTOServiceItem, *validate.Errors) {
	if mtoServiceItem == nil {
		return nil, nil
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
	case primemessages.MTOServiceItemModelTypeMTOServiceItemOriginSIT:

		originsit := mtoServiceItem.(*primemessages.MTOServiceItemOriginSIT)

		if originsit.ReServiceCode != nil {
			model.ReService.Code = models.ReServiceCode(*originsit.ReServiceCode)
		}

		model.Reason = originsit.Reason
		sitEntryDate := handlers.FmtDatePtrToPopPtr(originsit.SitEntryDate)

		if sitEntryDate != nil {
			model.SITEntryDate = sitEntryDate
		}

		model.SITPostalCode = originsit.SitPostalCode

		model.SITOriginHHGActualAddress = AddressModel(originsit.SitHHGActualOrigin)
		if model.SITOriginHHGActualAddress != nil {
			model.SITOriginHHGActualAddressID = &model.SITOriginHHGActualAddress.ID
		}

	case primemessages.MTOServiceItemModelTypeMTOServiceItemDestSIT:
		destsit := mtoServiceItem.(*primemessages.MTOServiceItemDestSIT)

		if destsit.ReServiceCode != nil {
			model.ReService.Code = models.ReServiceCode(*destsit.ReServiceCode)
		}

		model.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeFirst,
				TimeMilitary:               *destsit.TimeMilitary1,
				FirstAvailableDeliveryDate: time.Time(*destsit.FirstAvailableDeliveryDate1),
			},
			models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeSecond,
				TimeMilitary:               *destsit.TimeMilitary2,
				FirstAvailableDeliveryDate: time.Time(*destsit.FirstAvailableDeliveryDate2),
			},
		}

		sitEntryDate := handlers.FmtDatePtrToPopPtr(destsit.SitEntryDate)

		if sitEntryDate != nil {
			model.SITEntryDate = sitEntryDate
		}

		if destsit.SitDepartureDate != nil {
			model.SITDepartureDate = handlers.FmtDatePtrToPopPtr(destsit.SitDepartureDate)
		}

	case primemessages.MTOServiceItemModelTypeMTOServiceItemShuttle:
		shuttleService := mtoServiceItem.(*primemessages.MTOServiceItemShuttle)
		// values to get from payload
		model.ReService.Code = models.ReServiceCode(*shuttleService.ReServiceCode)
		model.Reason = shuttleService.Reason
		model.Description = shuttleService.Description
	case primemessages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating:
		domesticCrating := mtoServiceItem.(*primemessages.MTOServiceItemDomesticCrating)

		// additional validation for this specific service item type
		verrs := validateDomesticCrating(*domesticCrating)
		if verrs.HasAny() {
			return nil, verrs
		}

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

	return model, nil
}

// MTOServiceItemModelFromUpdate converts the payload from UpdateMTOServiceItem to a normal MTOServiceItem model.
// The payload for this is different than the one for create.
func MTOServiceItemModelFromUpdate(mtoServiceItemID string, mtoServiceItem primemessages.UpdateMTOServiceItem) (*models.MTOServiceItem, *validate.Errors) {
	verrs := validate.NewErrors()
	if mtoServiceItem == nil {
		verrs.Add("mtoServiceItem", "was nil")
		return nil, verrs
	}

	nilUUID := strfmt.UUID(uuid.Nil.String())

	if mtoServiceItem.ID().String() != "" && mtoServiceItem.ID() != nilUUID && mtoServiceItem.ID().String() != mtoServiceItemID {
		verrs.Add("id", "value does not agree with mtoServiceItemID in path - omit from body or correct")
	}

	// Create the service item model
	model := &models.MTOServiceItem{
		ID: uuid.FromStringOrNil(mtoServiceItemID),
	}

	// Here we initialize more fields below for the specific model types.
	// Currently only UpdateMTOServiceItemSIT is supported, more to be expected
	modelType := mtoServiceItem.ModelType()
	if modelType == primemessages.UpdateMTOServiceItemModelTypeUpdateMTOServiceItemSIT {
		sit := mtoServiceItem.(*primemessages.UpdateMTOServiceItemSIT)
		model.SITDepartureDate = swag.Time(time.Time(sit.SitDepartureDate))
		model.ReService.Code = models.ReServiceCode(sit.ReServiceCode)
		model.SITDestinationFinalAddress = AddressModel(sit.SitDestinationFinalAddress)
		if model.SITDestinationFinalAddress != nil {
			model.SITDestinationFinalAddressID = &model.SITDestinationFinalAddress.ID
		}

		if verrs != nil && verrs.HasAny() {
			return nil, verrs
		}

		return model, nil
	}

	verrs.Add("mtoServiceItem", "The model type of the service item is not allowed")
	return nil, verrs

}

// validateDomesticCrating validates this mto service item domestic crating
func validateDomesticCrating(m primemessages.MTOServiceItemDomesticCrating) *validate.Errors {
	return validate.Validate(
		&models.ItemCanFitInsideCrate{Name: "Item", NameCompared: "Crate", Item: m.Item, Crate: m.Crate},
	)
}
