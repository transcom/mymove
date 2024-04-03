package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// AddressModel model
func AddressModel(address *primemessages.Address) *models.Address {
	// To check if the model is intended to be blank, we'll look at both ID and StreetAddress1
	// We should always have ID if the user intends to update an Address,
	// and StreetAddress1 is a required field on creation. If both are blank, it should be treated as nil.
	var blankSwaggerID strfmt.UUID
	if address == nil || (address.ID == blankSwaggerID && address.StreetAddress1 == nil || address.County == nil) {
		return nil
	}
	modelAddress := &models.Address{
		ID:             uuid.FromStringOrNil(address.ID.String()),
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		Country:        address.Country,
		County:         *address.County,
	}
	if address.StreetAddress1 != nil {
		modelAddress.StreetAddress1 = *address.StreetAddress1
	}
	if address.City != nil {
		modelAddress.City = *address.City
	}
	if address.State != nil {
		modelAddress.State = *address.State
	}
	if address.PostalCode != nil {
		modelAddress.PostalCode = *address.PostalCode
	}
	return modelAddress
}

// ReweighModelFromUpdate model
func ReweighModelFromUpdate(reweigh *primemessages.UpdateReweigh, reweighID strfmt.UUID, mtoShipmentID strfmt.UUID) *models.Reweigh {
	if reweigh == nil {
		return nil
	}

	model := &models.Reweigh{
		ID:         uuid.FromStringOrNil(reweighID.String()),
		ShipmentID: uuid.FromStringOrNil(mtoShipmentID.String()),
	}

	if reweigh.Weight != nil {
		model.Weight = handlers.PoundPtrFromInt64Ptr(reweigh.Weight)
	}

	if reweigh.VerificationReason != nil {
		model.VerificationReason = reweigh.VerificationReason
	}

	return model
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
		MoveTaskOrderID:             uuid.FromStringOrNil(mtoShipment.MoveTaskOrderID.String()),
		CustomerRemarks:             mtoShipment.CustomerRemarks,
		Diversion:                   mtoShipment.Diversion,
		CounselorRemarks:            mtoShipment.CounselorRemarks,
		HasSecondaryPickupAddress:   handlers.FmtBool(false),
		HasSecondaryDeliveryAddress: handlers.FmtBool(false),
	}

	if mtoShipment.ShipmentType != nil {
		model.ShipmentType = models.MTOShipmentType(*mtoShipment.ShipmentType)
	}

	if mtoShipment.PrimeEstimatedWeight != nil {
		estimatedWeight := unit.Pound(*mtoShipment.PrimeEstimatedWeight)
		model.PrimeEstimatedWeight = &estimatedWeight
		recordedDate := time.Now()
		model.PrimeEstimatedWeightRecordedDate = &recordedDate
	}

	if mtoShipment.RequestedPickupDate != nil {
		model.RequestedPickupDate = models.TimePointer(time.Time(*mtoShipment.RequestedPickupDate))
	}

	// Set up address models
	var addressModel *models.Address

	addressModel = AddressModel(&mtoShipment.PickupAddress.Address)
	if addressModel != nil {
		model.PickupAddress = addressModel
	}

	addressModel = AddressModel(&mtoShipment.DestinationAddress.Address)
	if addressModel != nil {
		model.DestinationAddress = addressModel
	}

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	if mtoShipment.PpmShipment != nil {
		model.PPMShipment = PPMShipmentModelFromCreate(mtoShipment.PpmShipment)
		model.PPMShipment.Shipment = *model
	}

	return model
}

// Non SIT Address update Model
func ShipmentAddressUpdateModel(nonSITAddressUpdate *primemessages.UpdateShipmentDestinationAddress, MtoShipmentID uuid.UUID) *models.ShipmentAddressUpdate {
	if nonSITAddressUpdate == nil {
		return nil
	}

	model := &models.ShipmentAddressUpdate{
		ContractorRemarks: *nonSITAddressUpdate.ContractorRemarks,
		ShipmentID:        MtoShipmentID,
	}

	addressModel := AddressModel(nonSITAddressUpdate.NewAddress)
	if addressModel != nil {
		model.NewAddress = *addressModel
	}

	return model

}

// PPMShipmentModelFromCreate model
func PPMShipmentModelFromCreate(ppmShipment *primemessages.CreatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	model := &models.PPMShipment{
		Status:                         models.PPMShipmentStatusSubmitted,
		SITExpected:                    ppmShipment.SitExpected,
		SecondaryPickupPostalCode:      ppmShipment.SecondaryPickupPostalCode,
		SecondaryDestinationPostalCode: ppmShipment.SecondaryDestinationPostalCode,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
	}

	expectedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.ExpectedDepartureDate)
	if expectedDepartureDate != nil && !expectedDepartureDate.IsZero() {
		model.ExpectedDepartureDate = *expectedDepartureDate
	}

	if ppmShipment.PickupPostalCode != nil {
		model.PickupPostalCode = *ppmShipment.PickupPostalCode
	}
	if ppmShipment.DestinationPostalCode != nil {
		model.DestinationPostalCode = *ppmShipment.DestinationPostalCode
	}

	if model.SITExpected != nil && *model.SITExpected {
		if ppmShipment.SitLocation != nil {
			sitLocation := models.SITLocationType(*ppmShipment.SitLocation)
			model.SITLocation = &sitLocation
		}

		model.SITEstimatedWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.SitEstimatedWeight)

		sitEstimatedEntryDate := handlers.FmtDatePtrToPopPtr(ppmShipment.SitEstimatedEntryDate)
		if sitEstimatedEntryDate != nil && !sitEstimatedEntryDate.IsZero() {
			model.SITEstimatedEntryDate = sitEstimatedEntryDate
		}
		sitEstimatedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.SitEstimatedDepartureDate)
		if sitEstimatedDepartureDate != nil && !sitEstimatedDepartureDate.IsZero() {
			model.SITEstimatedDepartureDate = sitEstimatedDepartureDate
		}
	}

	if model.HasProGear != nil && *model.HasProGear {
		model.ProGearWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight)
		model.SpouseProGearWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight)
	}

	return model
}

// MTOShipmentModelFromUpdate model
func MTOShipmentModelFromUpdate(mtoShipment *primemessages.UpdateMTOShipment, mtoShipmentID strfmt.UUID) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	model := &models.MTOShipment{
		ID:                         uuid.FromStringOrNil(mtoShipmentID.String()),
		ActualPickupDate:           handlers.FmtDatePtrToPopPtr(mtoShipment.ActualPickupDate),
		FirstAvailableDeliveryDate: handlers.FmtDatePtrToPopPtr(mtoShipment.FirstAvailableDeliveryDate),
		ScheduledPickupDate:        handlers.FmtDatePtrToPopPtr(mtoShipment.ScheduledPickupDate),
		ActualDeliveryDate:         handlers.FmtDatePtrToPopPtr(mtoShipment.ActualDeliveryDate),
		ScheduledDeliveryDate:      handlers.FmtDatePtrToPopPtr(mtoShipment.ScheduledDeliveryDate),
		ShipmentType:               models.MTOShipmentType(mtoShipment.ShipmentType),
		Diversion:                  mtoShipment.Diversion,
		CounselorRemarks:           mtoShipment.CounselorRemarks,
	}

	if mtoShipment.PrimeActualWeight != nil {
		actualWeight := unit.Pound(*mtoShipment.PrimeActualWeight)
		model.PrimeActualWeight = &actualWeight
	}

	if mtoShipment.NtsRecordedWeight != nil && *mtoShipment.NtsRecordedWeight > 0 {
		ntsRecordedWeight := unit.Pound(*mtoShipment.NtsRecordedWeight)
		model.NTSRecordedWeight = &ntsRecordedWeight
	}

	if mtoShipment.PrimeEstimatedWeight != nil {
		estimatedWeight := unit.Pound(*mtoShipment.PrimeEstimatedWeight)
		model.PrimeEstimatedWeight = &estimatedWeight
	}

	// Set up address models
	var addressModel *models.Address

	addressModel = AddressModel(&mtoShipment.PickupAddress.Address)
	if addressModel != nil {
		model.PickupAddress = addressModel
	}

	addressModel = AddressModel(&mtoShipment.DestinationAddress.Address)
	if addressModel != nil {
		model.DestinationAddress = addressModel
	}

	if mtoShipment.DestinationType != nil {
		destinationType := models.DestinationType(*mtoShipment.DestinationType)
		model.DestinationType = &destinationType
	}

	addressModel = AddressModel(&mtoShipment.SecondaryPickupAddress.Address)
	if addressModel != nil {
		model.SecondaryPickupAddress = addressModel
		secondaryPickupAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.SecondaryPickupAddressID = &secondaryPickupAddressID
		model.HasSecondaryPickupAddress = handlers.FmtBool(true)
	}

	addressModel = AddressModel(&mtoShipment.SecondaryDeliveryAddress.Address)
	if addressModel != nil {
		model.SecondaryDeliveryAddress = addressModel
		secondaryDeliveryAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.SecondaryDeliveryAddressID = &secondaryDeliveryAddressID
		model.HasSecondaryDeliveryAddress = handlers.FmtBool(true)
	}

	if mtoShipment.PpmShipment != nil {
		model.PPMShipment = PPMShipmentModelFromUpdate(mtoShipment.PpmShipment)
		model.PPMShipment.Shipment = *model
	}

	return model
}

// PPMShipmentModelFromUpdate model
func PPMShipmentModelFromUpdate(ppmShipment *primemessages.UpdatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	model := &models.PPMShipment{
		SecondaryPickupPostalCode:      ppmShipment.SecondaryPickupPostalCode,
		SecondaryDestinationPostalCode: ppmShipment.SecondaryDestinationPostalCode,
		SITExpected:                    ppmShipment.SitExpected,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight),
	}

	expectedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.ExpectedDepartureDate)
	if expectedDepartureDate != nil && !expectedDepartureDate.IsZero() {
		model.ExpectedDepartureDate = *expectedDepartureDate
	}

	if ppmShipment.PickupPostalCode != nil {
		model.PickupPostalCode = *ppmShipment.PickupPostalCode
	}
	if ppmShipment.DestinationPostalCode != nil {
		model.DestinationPostalCode = *ppmShipment.DestinationPostalCode
	}

	if ppmShipment.SitLocation != nil {
		sitLocation := models.SITLocationType(*ppmShipment.SitLocation)
		model.SITLocation = &sitLocation
	}

	model.SITEstimatedWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.SitEstimatedWeight)

	sitEstimatedEntryDate := handlers.FmtDatePtrToPopPtr(ppmShipment.SitEstimatedEntryDate)
	if sitEstimatedEntryDate != nil && !sitEstimatedEntryDate.IsZero() {
		model.SITEstimatedEntryDate = sitEstimatedEntryDate
	}
	sitEstimatedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.SitEstimatedDepartureDate)
	if sitEstimatedDepartureDate != nil && !sitEstimatedDepartureDate.IsZero() {
		model.SITEstimatedDepartureDate = sitEstimatedDepartureDate
	}

	return model
}

// MTOServiceItemModel model
func MTOServiceItemModel(mtoServiceItem primemessages.MTOServiceItem) (*models.MTOServiceItem, *validate.Errors) {
	if mtoServiceItem == nil {
		return nil, nil
	}

	shipmentID := uuid.FromStringOrNil(mtoServiceItem.MtoShipmentID().String())

	// Default requested approvals value when an MTOServiceItem is created
	requestedApprovalsRequestedStatus := false

	// basic service item
	model := &models.MTOServiceItem{
		ID:                                uuid.FromStringOrNil(mtoServiceItem.ID().String()),
		MoveTaskOrderID:                   uuid.FromStringOrNil(mtoServiceItem.MoveTaskOrderID().String()),
		MTOShipmentID:                     &shipmentID,
		CreatedAt:                         time.Now(),
		UpdatedAt:                         time.Now(),
		RequestedApprovalsRequestedStatus: &requestedApprovalsRequestedStatus,
	}

	// here we initialize more fields below for other service item types. Eg. MTOServiceItemDOFSIT
	switch mtoServiceItem.ModelType() {
	case primemessages.MTOServiceItemModelTypeMTOServiceItemOriginSIT:

		originsit := mtoServiceItem.(*primemessages.MTOServiceItemOriginSIT)

		if originsit.ReServiceCode != nil {
			model.ReService.Code = models.ReServiceCode(*originsit.ReServiceCode)
		}

		model.Reason = originsit.Reason
		// Check for reason required field on a DDFSIT
		if model.ReService.Code == models.ReServiceCodeDOASIT {
			reasonVerrs := validateReasonOriginSIT(*originsit)

			if reasonVerrs.HasAny() {
				return nil, reasonVerrs
			}
		}

		if model.ReService.Code == models.ReServiceCodeDOFSIT {
			reasonVerrs := validateReasonOriginSIT(*originsit)

			if reasonVerrs.HasAny() {
				return nil, reasonVerrs
			}
		}

		sitEntryDate := handlers.FmtDatePtrToPopPtr(originsit.SitEntryDate)

		if sitEntryDate != nil {
			model.SITEntryDate = sitEntryDate
		}

		if originsit.SitDepartureDate != nil {
			model.SITDepartureDate = handlers.FmtDatePtrToPopPtr(originsit.SitDepartureDate)
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

		model.Reason = destsit.Reason
		sitEntryDate := handlers.FmtDatePtrToPopPtr(destsit.SitEntryDate)

		// Check for required fields on a DDFSIT
		if model.ReService.Code == models.ReServiceCodeDDFSIT {
			verrs := validateDDFSITForCreate(*destsit)
			reasonVerrs := validateReasonDestSIT(*destsit)

			if verrs.HasAny() {
				return nil, verrs
			}

			if reasonVerrs.HasAny() {
				return nil, reasonVerrs
			}
		}

		var customerContacts models.MTOServiceItemCustomerContacts

		if destsit.TimeMilitary1 != nil && destsit.FirstAvailableDeliveryDate1 != nil && destsit.DateOfContact1 != nil {
			customerContacts = append(customerContacts, models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeFirst,
				DateOfContact:              time.Time(*destsit.DateOfContact1),
				TimeMilitary:               *destsit.TimeMilitary1,
				FirstAvailableDeliveryDate: time.Time(*destsit.FirstAvailableDeliveryDate1),
			})
		}
		if destsit.TimeMilitary2 != nil && destsit.FirstAvailableDeliveryDate2 != nil && destsit.DateOfContact2 != nil {
			customerContacts = append(customerContacts, models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeSecond,
				DateOfContact:              time.Time(*destsit.DateOfContact2),
				TimeMilitary:               *destsit.TimeMilitary2,
				FirstAvailableDeliveryDate: time.Time(*destsit.FirstAvailableDeliveryDate2),
			})
		}

		model.CustomerContacts = customerContacts

		if sitEntryDate != nil {
			model.SITEntryDate = sitEntryDate
		}

		if destsit.SitDepartureDate != nil {
			model.SITDepartureDate = handlers.FmtDatePtrToPopPtr(destsit.SitDepartureDate)
		}

		model.SITDestinationFinalAddress = AddressModel(destsit.SitDestinationFinalAddress)
		if model.SITDestinationFinalAddress != nil {
			model.SITDestinationFinalAddressID = &model.SITDestinationFinalAddress.ID
		}

	case primemessages.MTOServiceItemModelTypeMTOServiceItemShuttle:
		shuttleService := mtoServiceItem.(*primemessages.MTOServiceItemShuttle)
		// values to get from payload
		model.ReService.Code = models.ReServiceCode(*shuttleService.ReServiceCode)
		model.Reason = shuttleService.Reason
		model.EstimatedWeight = handlers.PoundPtrFromInt64Ptr(shuttleService.EstimatedWeight)
		model.ActualWeight = handlers.PoundPtrFromInt64Ptr(shuttleService.ActualWeight)

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
		model.Reason = domesticCrating.Reason
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
		if basic.ReServiceCode != nil {
			model.ReService.Code = models.ReServiceCode(*basic.ReServiceCode)
		}
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
	switch mtoServiceItem.ModelType() {
	case primemessages.UpdateMTOServiceItemModelTypeUpdateMTOServiceItemSIT:
		sit := mtoServiceItem.(*primemessages.UpdateMTOServiceItemSIT)
		model.ReService.Code = models.ReServiceCode(sit.ReServiceCode)
		model.SITDestinationFinalAddress = AddressModel(sit.SitDestinationFinalAddress)
		model.SITRequestedDelivery = (*time.Time)(sit.SitRequestedDelivery)
		model.Status = models.MTOServiceItemStatusSubmitted
		model.Reason = sit.UpdateReason

		var zeroDate strfmt.Date
		if sit.SitDepartureDate != zeroDate {
			model.SITDepartureDate = models.TimePointer(time.Time(sit.SitDepartureDate))
		}

		if sit.SitEntryDate != nil {
			model.SITEntryDate = (*time.Time)(sit.SitEntryDate)
		}

		if sit.SitPostalCode != nil {
			newPostalCode := sit.SitPostalCode
			model.SITPostalCode = newPostalCode
		}

		if model.SITDestinationFinalAddress != nil {
			model.SITDestinationFinalAddressID = &model.SITDestinationFinalAddress.ID
		}

		if sit.SitCustomerContacted != nil {
			model.SITCustomerContacted = handlers.FmtDatePtrToPopPtr(sit.SitCustomerContacted)
		}

		if sit.SitRequestedDelivery != nil {
			model.SITRequestedDelivery = handlers.FmtDatePtrToPopPtr(sit.SitRequestedDelivery)
		}

		// If the request params have a have the RequestApprovalsRequestedStatus set the model RequestApprovalsRequestedStatus value to the incoming value
		if sit.RequestApprovalsRequestedStatus != nil {
			pointerValue := *sit.RequestApprovalsRequestedStatus
			model.RequestedApprovalsRequestedStatus = &pointerValue
		}

		if sit.ReServiceCode == string(models.ReServiceCodeDDDSIT) ||
			sit.ReServiceCode == string(models.ReServiceCodeDDASIT) ||
			sit.ReServiceCode == string(models.ReServiceCodeDDFSIT) ||
			sit.ReServiceCode == string(models.ReServiceCodeDDSFSC) {
			destSitVerrs := validateDestSITForUpdate(*sit)

			if destSitVerrs.HasAny() {
				return nil, destSitVerrs
			}
			var customerContacts models.MTOServiceItemCustomerContacts
			if sit.TimeMilitary1 != nil && sit.FirstAvailableDeliveryDate1 != nil && sit.DateOfContact1 != nil {
				contact1 := models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeFirst,
					DateOfContact:              time.Time(*sit.DateOfContact1),
					TimeMilitary:               *sit.TimeMilitary1,
					FirstAvailableDeliveryDate: time.Time(*sit.FirstAvailableDeliveryDate1),
				}
				customerContacts = append(customerContacts, contact1)
			}
			if sit.TimeMilitary2 != nil && sit.FirstAvailableDeliveryDate2 != nil && sit.DateOfContact2 != nil {
				contact2 := models.MTOServiceItemCustomerContact{
					Type:                       models.CustomerContactTypeSecond,
					DateOfContact:              time.Time(*sit.DateOfContact2),
					TimeMilitary:               *sit.TimeMilitary2,
					FirstAvailableDeliveryDate: time.Time(*sit.FirstAvailableDeliveryDate2),
				}
				customerContacts = append(customerContacts, contact2)
			}
			if len(customerContacts) > 0 {
				model.CustomerContacts = customerContacts
			}

			model.SITCustomerContacted = handlers.FmtDatePtrToPopPtr(sit.SitCustomerContacted)
			model.SITRequestedDelivery = handlers.FmtDatePtrToPopPtr(sit.SitRequestedDelivery)
		}

		if verrs != nil && verrs.HasAny() {
			return nil, verrs
		}

	case primemessages.UpdateMTOServiceItemModelTypeUpdateMTOServiceItemShuttle:
		shuttle := mtoServiceItem.(*primemessages.UpdateMTOServiceItemShuttle)
		model.EstimatedWeight = handlers.PoundPtrFromInt64Ptr(shuttle.EstimatedWeight)
		model.ActualWeight = handlers.PoundPtrFromInt64Ptr(shuttle.ActualWeight)

		if verrs != nil && verrs.HasAny() {
			return nil, verrs
		}
	default:
		// assume basic service item
		if verrs != nil && verrs.HasAny() {
			return nil, verrs
		}
	}

	return model, nil
}

func ServiceRequestDocumentUploadModel(u models.Upload) *primemessages.UploadWithOmissions {
	return &primemessages.UploadWithOmissions{
		Bytes:       &u.Bytes,
		ContentType: &u.ContentType,
		Filename:    &u.Filename,
		CreatedAt:   (strfmt.DateTime)(u.CreatedAt),
		UpdatedAt:   (strfmt.DateTime)(u.UpdatedAt),
	}
}

// SITExtensionModel transform the request data the sitExtension model
func SITExtensionModel(sitExtension *primemessages.CreateSITExtension, mtoShipmentID strfmt.UUID) *models.SITDurationUpdate {
	if sitExtension == nil {
		return nil
	}

	model := &models.SITDurationUpdate{
		MTOShipmentID:     uuid.FromStringOrNil(mtoShipmentID.String()),
		RequestedDays:     int(*sitExtension.RequestedDays),
		ContractorRemarks: sitExtension.ContractorRemarks,
		RequestReason:     models.SITDurationUpdateRequestReason(*sitExtension.RequestReason),
	}

	return model
}

// SITAddressUpdateModel
func SITAddressUpdateModel(sitAddressUpdate *primemessages.CreateSITAddressUpdateRequest) *models.SITAddressUpdate {
	if sitAddressUpdate == nil {
		return nil
	}

	model := &models.SITAddressUpdate{
		ContractorRemarks: sitAddressUpdate.ContractorRemarks,
		MTOServiceItemID:  uuid.FromStringOrNil(sitAddressUpdate.MtoServiceItemID.String()),
	}

	addressModel := AddressModel(sitAddressUpdate.NewAddress)
	if addressModel != nil {
		model.NewAddress = *addressModel
		newAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.NewAddressID = newAddressID
	}

	return model
}

// validateDomesticCrating validates this mto service item domestic crating
func validateDomesticCrating(m primemessages.MTOServiceItemDomesticCrating) *validate.Errors {
	return validate.Validate(
		&models.ItemCanFitInsideCrate{
			Name:         "Item",
			NameCompared: "Crate",
			Item:         &m.Item.MTOServiceItemDimension,
			Crate:        &m.Crate.MTOServiceItemDimension,
		},
	)
}

// validateDDFSITForCreate validates DDFSIT service item has all required fields
func validateDDFSITForCreate(m primemessages.MTOServiceItemDestSIT) *validate.Errors {
	verrs := validate.NewErrors()

	if m.FirstAvailableDeliveryDate1 == nil && m.DateOfContact1 != nil && m.TimeMilitary1 != nil {
		verrs.Add("firstAvailableDeliveryDate1", "firstAvailableDeliveryDate1, dateOfContact1, and timeMilitary1 must be provided together in body.")
	}
	if m.DateOfContact1 == nil && m.TimeMilitary1 != nil && m.FirstAvailableDeliveryDate1 != nil {
		verrs.Add("DateOfContact1", "dateOfContact1, timeMilitary1, and firstAvailableDeliveryDate1 must be provided together in body.")
	}
	if m.TimeMilitary1 == nil && m.DateOfContact1 != nil && m.FirstAvailableDeliveryDate1 != nil {
		verrs.Add("timeMilitary1", "timeMilitary1, dateOfContact1, and firstAvailableDeliveryDate1 must be provided together in body.")
	}
	if m.FirstAvailableDeliveryDate2 == nil && m.DateOfContact2 != nil && m.TimeMilitary2 != nil {
		verrs.Add("firstAvailableDeliveryDate2", "firstAvailableDeliveryDate2, dateOfContact2, and timeMilitary2 must be provided together in body.")
	}
	if m.DateOfContact2 == nil && m.TimeMilitary2 != nil && m.FirstAvailableDeliveryDate2 != nil {
		verrs.Add("DateOfContact1", "dateOfContact2, firstAvailableDeliveryDate2, and timeMilitary2 must be provided together in body.")
	}
	if m.TimeMilitary2 == nil && m.DateOfContact2 != nil && m.FirstAvailableDeliveryDate2 != nil {
		verrs.Add("timeMilitary2", "timeMilitary2, firstAvailableDeliveryDate2, and dateOfContact2 must be provided together in body.")
	}
	return verrs
}

// validateDestSITForUpdate validates DDDSIT service item has all required fields
func validateDestSITForUpdate(m primemessages.UpdateMTOServiceItemSIT) *validate.Errors {
	verrs := validate.NewErrors()

	if m.FirstAvailableDeliveryDate1 == nil && m.DateOfContact1 != nil && m.TimeMilitary1 != nil {
		verrs.Add("firstAvailableDeliveryDate1", "firstAvailableDeliveryDate1, dateOfContact1, and timeMilitary1 must be provided together in body.")
	}
	if m.DateOfContact1 == nil && m.TimeMilitary1 != nil && m.FirstAvailableDeliveryDate1 != nil {
		verrs.Add("DateOfContact1", "dateOfContact1, timeMilitary1, and firstAvailableDeliveryDate1 must be provided together in body.")
	}
	if m.TimeMilitary1 == nil && m.DateOfContact1 != nil && m.FirstAvailableDeliveryDate1 != nil {
		verrs.Add("timeMilitary1", "timeMilitary1, dateOfContact1, and firstAvailableDeliveryDate1 must be provided together in body.")
	}
	if m.FirstAvailableDeliveryDate2 == nil && m.DateOfContact2 != nil && m.TimeMilitary2 != nil {
		verrs.Add("firstAvailableDeliveryDate2", "firstAvailableDeliveryDate2, dateOfContact2, and timeMilitary2 must be provided together in body.")
	}
	if m.DateOfContact2 == nil && m.TimeMilitary2 != nil && m.FirstAvailableDeliveryDate2 != nil {
		verrs.Add("DateOfContact1", "dateOfContact2, firstAvailableDeliveryDate2, and timeMilitary2 must be provided together in body.")
	}
	if m.TimeMilitary2 == nil && m.DateOfContact2 != nil && m.FirstAvailableDeliveryDate2 != nil {
		verrs.Add("timeMilitary2", "timeMilitary2, firstAvailableDeliveryDate2, and dateOfContact2 must be provided together in body.")
	}
	return verrs
}

// validateReasonDestSIT validates that Destination SIT service items have required Reason field
func validateReasonDestSIT(m primemessages.MTOServiceItemDestSIT) *validate.Errors {
	verrs := validate.NewErrors()

	if m.Reason == nil || m.Reason == models.StringPointer("") {
		verrs.Add("reason", "reason is required in body.")
	}
	return verrs
}

// validateReasonOriginSIT validates that Origin SIT service items have required Reason field
func validateReasonOriginSIT(m primemessages.MTOServiceItemOriginSIT) *validate.Errors {
	verrs := validate.NewErrors()

	if m.Reason == nil || m.Reason == models.StringPointer("") {
		verrs.Add("reason", "reason is required in body.")
	}
	return verrs
}
