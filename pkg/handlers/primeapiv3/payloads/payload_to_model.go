package payloads

import (
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/primev3messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// CountryModel model
func CountryModel(country *string) *models.Country {
	// The prime doesn't know the uuids of our countries, so for now we are going to just populate the name so we can query that
	// when creating the address IF it is provided - else this will be nil and a US country will be created
	if country == nil {
		return nil
	}

	modelCountry := &models.Country{
		Country: *country,
	}
	return modelCountry
}

// AddressModel model
func AddressModel(address *primev3messages.Address) *models.Address {
	// To check if the model is intended to be blank, we'll look at both ID and StreetAddress1
	// We should always have ID if the user intends to update an Address,
	// and StreetAddress1 is a required field on creation. If both are blank, it should be treated as nil.
	var blankSwaggerID strfmt.UUID
	if address == nil || (address.ID == blankSwaggerID && address.StreetAddress1 == nil) {
		return nil
	}
	modelAddress := &models.Address{
		ID:             uuid.FromStringOrNil(address.ID.String()),
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
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
	if address.Country != nil {
		modelAddress.Country = CountryModel(address.Country)
	}
	return modelAddress
}

func PPMDestinationAddressModel(address *primev3messages.PPMDestinationAddress) *models.Address {
	// To check if the model is intended to be blank, we'll look at ID and City, State, PostalCode
	// We should always have ID if the user intends to update an Address,
	// and City, State, PostalCode is a required field on creation. If both are blank, it should be treated as nil.
	var blankSwaggerID strfmt.UUID
	// unlike other addresses PPM destination address can be created without StreetAddress1
	if address == nil || (address.ID == blankSwaggerID && address.City == nil && address.State == nil && address.PostalCode == nil) {
		return nil
	}
	modelAddress := &models.Address{
		ID:             uuid.FromStringOrNil(address.ID.String()),
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}

	if address.StreetAddress1 != nil && len(strings.Trim(*address.StreetAddress1, " ")) > 0 {
		modelAddress.StreetAddress1 = *address.StreetAddress1
	} else {
		// Street address 1 is optional for certain business context but not nullable on the database level.
		// Use place holder text to represent NULL.
		modelAddress.StreetAddress1 = models.STREET_ADDRESS_1_NOT_PROVIDED
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
	if address.Country != nil {
		modelAddress.Country = CountryModel(address.Country)
	}
	return modelAddress
}

// ReweighModelFromUpdate model
func ReweighModelFromUpdate(reweigh *primev3messages.UpdateReweigh, reweighID strfmt.UUID, mtoShipmentID strfmt.UUID) *models.Reweigh {
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
func MTOAgentModel(mtoAgent *primev3messages.MTOAgent) *models.MTOAgent {
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
func MTOAgentsModel(mtoAgents *primev3messages.MTOAgents) *models.MTOAgents {
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
func MTOServiceItemModelListFromCreate(mtoShipment *primev3messages.CreateMTOShipment) (models.MTOServiceItems, *validate.Errors) {

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
func MTOShipmentModelFromCreate(mtoShipment *primev3messages.CreateMTOShipment) (*models.MTOShipment, *validate.Errors) {
	verrs := validate.NewErrors()

	if mtoShipment == nil {
		verrs.Add("mtoShipment", "mtoShipment object is nil.")
		return nil, verrs
	}

	if mtoShipment.MoveTaskOrderID == nil {
		verrs.Add("mtoShipment", "MoveTaskOrderID is nil.")
		return nil, verrs
	}

	var divertedFromShipmentID *uuid.UUID
	if mtoShipment.DivertedFromShipmentID != "" {
		// Create the UUID in memory so it can be referenced
		uuid := uuid.FromStringOrNil(mtoShipment.DivertedFromShipmentID.String())
		divertedFromShipmentID = &uuid
	}

	model := &models.MTOShipment{
		MoveTaskOrderID:             uuid.FromStringOrNil(mtoShipment.MoveTaskOrderID.String()),
		CustomerRemarks:             mtoShipment.CustomerRemarks,
		Diversion:                   mtoShipment.Diversion,
		DivertedFromShipmentID:      divertedFromShipmentID,
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

	addressModel = AddressModel(&mtoShipment.SecondaryPickupAddress.Address)
	if addressModel != nil {
		model.SecondaryPickupAddress = addressModel
		model.HasSecondaryPickupAddress = handlers.FmtBool(true)
	}

	addressModel = AddressModel(&mtoShipment.TertiaryPickupAddress.Address)
	if addressModel != nil {
		model.TertiaryPickupAddress = addressModel
		model.HasTertiaryPickupAddress = handlers.FmtBool(true)
	}

	addressModel = AddressModel(&mtoShipment.DestinationAddress.Address)
	if addressModel != nil {
		model.DestinationAddress = addressModel
	}

	addressModel = AddressModel(&mtoShipment.SecondaryDestinationAddress.Address)
	if addressModel != nil {
		model.SecondaryDeliveryAddress = addressModel
		model.HasSecondaryDeliveryAddress = handlers.FmtBool(true)
	}

	addressModel = AddressModel(&mtoShipment.TertiaryDestinationAddress.Address)
	if addressModel != nil {
		model.TertiaryDeliveryAddress = addressModel
		model.HasTertiaryDeliveryAddress = handlers.FmtBool(true)
	}

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	if mtoShipment.PpmShipment != nil {
		ppmShipment, err := PPMShipmentModelFromCreate(mtoShipment.PpmShipment)
		if err != nil {
			verrs.Add("mtoShipment", err.Error())
			return nil, verrs
		} else {
			model.PPMShipment = ppmShipment
			model.PPMShipment.Shipment = *model
		}
	}

	if mtoShipment.BoatShipment != nil {
		model.BoatShipment, verrs = BoatShipmentModelFromCreate(mtoShipment)
		if verrs != nil && verrs.HasAny() {
			return nil, verrs
		}
		model.BoatShipment.Shipment = *model
	}

	if mtoShipment.MobileHomeShipment != nil {
		model.MobileHome, verrs = MobileHomeShipmentModelFromCreate(mtoShipment)
		if verrs != nil && verrs.HasAny() {
			return nil, verrs
		}
		model.MobileHome.Shipment = *model
	}

	return model, nil
}

// Non SIT Address update Model
func ShipmentAddressUpdateModel(nonSITAddressUpdate *primev3messages.UpdateShipmentDestinationAddress, MtoShipmentID uuid.UUID) *models.ShipmentAddressUpdate {
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
func PPMShipmentModelFromCreate(ppmShipment *primev3messages.CreatePPMShipment) (*models.PPMShipment, *validate.Errors) {
	verrs := validate.NewErrors()
	if ppmShipment == nil {
		return nil, nil
	}

	model := &models.PPMShipment{
		Status:          models.PPMShipmentStatusSubmitted,
		SITExpected:     ppmShipment.SitExpected,
		EstimatedWeight: handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		HasProGear:      ppmShipment.HasProGear,
	}

	expectedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.ExpectedDepartureDate)
	if expectedDepartureDate != nil && !expectedDepartureDate.IsZero() {
		model.ExpectedDepartureDate = *expectedDepartureDate
	}

	// Set up address models
	var addressModel *models.Address

	addressModel = AddressModel(&ppmShipment.PickupAddress.Address)
	if addressModel != nil {
		model.PickupAddress = addressModel
	}

	addressModel = AddressModel(&ppmShipment.SecondaryPickupAddress.Address)
	if addressModel != nil {
		model.SecondaryPickupAddress = addressModel
		model.HasSecondaryPickupAddress = handlers.FmtBool(true)
	}

	addressModel = AddressModel(&ppmShipment.TertiaryPickupAddress.Address)
	if addressModel != nil {
		if AddressModel(&ppmShipment.SecondaryPickupAddress.Address) == nil {
			verrs.Add("ppmShipment", "Shipment cannot have a third pickup address without a second pickup address present")
			return nil, verrs
		}
		model.TertiaryPickupAddress = addressModel
		model.HasTertiaryPickupAddress = handlers.FmtBool(true)
	}

	addressModel = PPMDestinationAddressModel(&ppmShipment.DestinationAddress.PPMDestinationAddress)
	if addressModel != nil {
		model.DestinationAddress = addressModel
	}

	addressModel = AddressModel(&ppmShipment.SecondaryDestinationAddress.Address)
	if addressModel != nil {
		model.SecondaryDestinationAddress = addressModel
		model.HasSecondaryDestinationAddress = handlers.FmtBool(true)
	}

	addressModel = AddressModel(&ppmShipment.TertiaryDestinationAddress.Address)
	if addressModel != nil {
		if AddressModel(&ppmShipment.SecondaryDestinationAddress.Address) == nil {
			verrs.Add("ppmShipment", "Shipment cannot have a third destination address without a second destination address present")
			return nil, verrs
		}
		model.TertiaryDestinationAddress = addressModel
		model.HasTertiaryDestinationAddress = handlers.FmtBool(true)
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

	if ppmShipment.IsActualExpenseReimbursement != nil {
		model.IsActualExpenseReimbursement = ppmShipment.IsActualExpenseReimbursement
	}

	return model, nil
}

// BoatShipmentModelFromCreate model
func BoatShipmentModelFromCreate(mtoShipment *primev3messages.CreateMTOShipment) (*models.BoatShipment, *validate.Errors) {
	reasonVerrs := validateBoatShipmentType(*mtoShipment.ShipmentType)
	if reasonVerrs.HasAny() {
		return nil, reasonVerrs
	}

	var shipmentType models.BoatShipmentType

	if *mtoShipment.ShipmentType == primev3messages.MTOShipmentTypeBOATHAULAWAY {
		shipmentType = models.BoatShipmentTypeHaulAway
	} else if *mtoShipment.ShipmentType == primev3messages.MTOShipmentTypeBOATTOWAWAY {
		shipmentType = models.BoatShipmentTypeTowAway
	}

	year := int(*mtoShipment.BoatShipment.Year)
	lengthInInches := int(*mtoShipment.BoatShipment.LengthInInches)
	widthInInches := int(*mtoShipment.BoatShipment.WidthInInches)
	heightInInches := int(*mtoShipment.BoatShipment.HeightInInches)
	model := &models.BoatShipment{
		Type:           shipmentType,
		Year:           &year,
		Make:           mtoShipment.BoatShipment.Make,
		Model:          mtoShipment.BoatShipment.Model,
		LengthInInches: &lengthInInches,
		WidthInInches:  &widthInInches,
		HeightInInches: &heightInInches,
		HasTrailer:     mtoShipment.BoatShipment.HasTrailer,
		IsRoadworthy:   mtoShipment.BoatShipment.IsRoadworthy,
	}

	return model, nil
}

// MobileHomeShipmentModelFromCreate model
func MobileHomeShipmentModelFromCreate(mtoShipment *primev3messages.CreateMTOShipment) (*models.MobileHome, *validate.Errors) {
	year := int(*mtoShipment.MobileHomeShipment.Year)
	lengthInInches := int(*mtoShipment.MobileHomeShipment.LengthInInches)
	widthInInches := int(*mtoShipment.MobileHomeShipment.WidthInInches)
	heightInInches := int(*mtoShipment.MobileHomeShipment.HeightInInches)
	model := &models.MobileHome{
		Year:           &year,
		Make:           mtoShipment.MobileHomeShipment.Make,
		Model:          mtoShipment.MobileHomeShipment.Model,
		LengthInInches: &lengthInInches,
		WidthInInches:  &widthInInches,
		HeightInInches: &heightInInches,
	}

	return model, nil
}

// MTOShipmentModelFromUpdate model
func MTOShipmentModelFromUpdate(mtoShipment *primev3messages.UpdateMTOShipment, mtoShipmentID strfmt.UUID) *models.MTOShipment {
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

	if mtoShipment.ActualProGearWeight != nil {
		actualProGearWeight := unit.Pound(*mtoShipment.ActualProGearWeight)
		model.ActualProGearWeight = &actualProGearWeight
	}

	if mtoShipment.ActualSpouseProGearWeight != nil {
		actualSpouseProGearWeight := unit.Pound(*mtoShipment.ActualSpouseProGearWeight)
		model.ActualSpouseProGearWeight = &actualSpouseProGearWeight
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

	addressModel = AddressModel(&mtoShipment.TertiaryPickupAddress.Address)
	if addressModel != nil {
		model.TertiaryPickupAddress = addressModel
		tertiaryPickupAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.TertiaryPickupAddressID = &tertiaryPickupAddressID
		model.HasTertiaryPickupAddress = handlers.FmtBool(true)
	}

	addressModel = AddressModel(&mtoShipment.SecondaryDeliveryAddress.Address)
	if addressModel != nil {
		model.SecondaryDeliveryAddress = addressModel
		secondaryDeliveryAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.SecondaryDeliveryAddressID = &secondaryDeliveryAddressID
		model.HasSecondaryDeliveryAddress = handlers.FmtBool(true)
	}

	addressModel = AddressModel(&mtoShipment.TertiaryDeliveryAddress.Address)
	if addressModel != nil {
		model.TertiaryDeliveryAddress = addressModel
		tertiaryDeliveryAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.TertiaryDeliveryAddressID = &tertiaryDeliveryAddressID
		model.HasTertiaryDeliveryAddress = handlers.FmtBool(true)
	}

	if mtoShipment.PpmShipment != nil {
		model.PPMShipment = PPMShipmentModelFromUpdate(mtoShipment.PpmShipment)
		model.PPMShipment.Shipment = *model
	}

	return model
}

// PPMShipmentModelFromUpdate model
func PPMShipmentModelFromUpdate(ppmShipment *primev3messages.UpdatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	model := &models.PPMShipment{
		SITExpected:                    ppmShipment.SitExpected,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight),
		HasSecondaryPickupAddress:      ppmShipment.HasSecondaryPickupAddress,
		HasSecondaryDestinationAddress: ppmShipment.HasSecondaryDestinationAddress,
		HasTertiaryPickupAddress:       ppmShipment.HasTertiaryPickupAddress,
		HasTertiaryDestinationAddress:  ppmShipment.HasTertiaryDestinationAddress,
	}

	// Set up address models
	var addressModel *models.Address

	addressModel = AddressModel(&ppmShipment.PickupAddress.Address)
	if addressModel != nil {
		model.PickupAddress = addressModel
	}

	// only assume secondary address is added if has flag is set to true. If false the address will be ignored.
	if ppmShipment.HasSecondaryPickupAddress != nil && *ppmShipment.HasSecondaryPickupAddress {
		addressModel = AddressModel(&ppmShipment.SecondaryPickupAddress.Address)
		if addressModel != nil {
			model.SecondaryPickupAddress = addressModel
			secondaryPickupAddressID := uuid.FromStringOrNil(addressModel.ID.String())
			model.SecondaryPickupAddressID = &secondaryPickupAddressID
		}
	}

	if ppmShipment.HasTertiaryPickupAddress != nil && *ppmShipment.HasTertiaryPickupAddress {
		addressModel = AddressModel(&ppmShipment.TertiaryPickupAddress.Address)
		if addressModel != nil {
			model.TertiaryPickupAddress = addressModel
			tertiaryPickupAddressID := uuid.FromStringOrNil(addressModel.ID.String())
			model.TertiaryPickupAddressID = &tertiaryPickupAddressID
		}
	}

	addressModel = PPMDestinationAddressModel(&ppmShipment.DestinationAddress.PPMDestinationAddress)
	if addressModel != nil {
		model.DestinationAddress = addressModel
	}

	// only assume secondary address is added if has flag is set to true. If false the address will be ignored.
	if ppmShipment.HasSecondaryDestinationAddress != nil && *ppmShipment.HasSecondaryDestinationAddress {
		addressModel = AddressModel(&ppmShipment.SecondaryDestinationAddress.Address)
		if addressModel != nil {
			model.SecondaryDestinationAddress = addressModel
			secondaryDestinationAddressID := uuid.FromStringOrNil(addressModel.ID.String())
			model.SecondaryDestinationAddressID = &secondaryDestinationAddressID
		}
	}

	if ppmShipment.HasTertiaryDestinationAddress != nil && *ppmShipment.HasTertiaryDestinationAddress {
		addressModel = AddressModel(&ppmShipment.TertiaryDestinationAddress.Address)
		if addressModel != nil {
			model.TertiaryDestinationAddress = addressModel
			tertiaryDestinationAddressID := uuid.FromStringOrNil(addressModel.ID.String())
			model.TertiaryDestinationAddressID = &tertiaryDestinationAddressID
		}
	}

	expectedDepartureDate := handlers.FmtDatePtrToPopPtr(ppmShipment.ExpectedDepartureDate)
	if expectedDepartureDate != nil && !expectedDepartureDate.IsZero() {
		model.ExpectedDepartureDate = *expectedDepartureDate
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

	if ppmShipment.IsActualExpenseReimbursement != nil {
		model.IsActualExpenseReimbursement = ppmShipment.IsActualExpenseReimbursement
	}

	return model
}

// MTOServiceItemModel model
func MTOServiceItemModel(mtoServiceItem primev3messages.MTOServiceItem) (*models.MTOServiceItem, *validate.Errors) {
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
	case primev3messages.MTOServiceItemModelTypeMTOServiceItemOriginSIT:

		originsit := mtoServiceItem.(*primev3messages.MTOServiceItemOriginSIT)

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

		model.SITPostalCode = originsit.SitPostalCode

		model.SITOriginHHGActualAddress = AddressModel(originsit.SitHHGActualOrigin)
		if model.SITOriginHHGActualAddress != nil {
			model.SITOriginHHGActualAddressID = &model.SITOriginHHGActualAddress.ID
		}

	case primev3messages.MTOServiceItemModelTypeMTOServiceItemDestSIT:
		destsit := mtoServiceItem.(*primev3messages.MTOServiceItemDestSIT)

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

	case primev3messages.MTOServiceItemModelTypeMTOServiceItemShuttle:
		shuttleService := mtoServiceItem.(*primev3messages.MTOServiceItemShuttle)
		// values to get from payload
		model.ReService.Code = models.ReServiceCode(*shuttleService.ReServiceCode)
		model.Reason = shuttleService.Reason
		model.EstimatedWeight = handlers.PoundPtrFromInt64Ptr(shuttleService.EstimatedWeight)
		model.ActualWeight = handlers.PoundPtrFromInt64Ptr(shuttleService.ActualWeight)

	case primev3messages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating:
		domesticCrating := mtoServiceItem.(*primev3messages.MTOServiceItemDomesticCrating)

		// additional validation for this specific service item type
		verrs := validateDomesticCrating(*domesticCrating)
		if verrs.HasAny() {
			return nil, verrs
		}

		// have to get code from payload
		model.ReService.Code = models.ReServiceCode(*domesticCrating.ReServiceCode)
		model.Description = domesticCrating.Description
		model.Reason = domesticCrating.Reason
		model.StandaloneCrate = domesticCrating.StandaloneCrate
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
	case primev3messages.MTOServiceItemModelTypeMTOServiceItemInternationalCrating:
		internationalCrating := mtoServiceItem.(*primev3messages.MTOServiceItemInternationalCrating)

		// additional validation for this specific service item type
		verrs := validateInternationalCrating(*internationalCrating)
		if verrs.HasAny() {
			return nil, verrs
		}

		// have to get code from payload
		model.ReService.Code = models.ReServiceCode(*internationalCrating.ReServiceCode)
		model.Description = internationalCrating.Description
		model.Reason = internationalCrating.Reason
		model.StandaloneCrate = internationalCrating.StandaloneCrate
		model.ExternalCrate = internationalCrating.ExternalCrate

		if model.ReService.Code == models.ReServiceCodeICRT {
			if internationalCrating.StandaloneCrate == nil {
				model.StandaloneCrate = models.BoolPointer(false)
			}
			if internationalCrating.ExternalCrate == nil {
				model.ExternalCrate = models.BoolPointer(false)
			}
		}
		model.Dimensions = models.MTOServiceItemDimensions{
			models.MTOServiceItemDimension{
				Type:   models.DimensionTypeItem,
				Length: unit.ThousandthInches(*internationalCrating.Item.Length),
				Height: unit.ThousandthInches(*internationalCrating.Item.Height),
				Width:  unit.ThousandthInches(*internationalCrating.Item.Width),
			},
			models.MTOServiceItemDimension{
				Type:   models.DimensionTypeCrate,
				Length: unit.ThousandthInches(*internationalCrating.Crate.Length),
				Height: unit.ThousandthInches(*internationalCrating.Crate.Height),
				Width:  unit.ThousandthInches(*internationalCrating.Crate.Width),
			},
		}
	default:
		// assume basic service item, take in provided re service code
		basic := mtoServiceItem.(*primev3messages.MTOServiceItemBasic)
		if basic.ReServiceCode != nil {
			model.ReService.Code = models.ReServiceCode(*basic.ReServiceCode)
		}
	}

	return model, nil
}

// MTOServiceItemModelFromUpdate converts the payload from UpdateMTOServiceItem to a normal MTOServiceItem model.
// The payload for this is different than the one for create.
func MTOServiceItemModelFromUpdate(mtoServiceItemID string, mtoServiceItem primev3messages.UpdateMTOServiceItem) (*models.MTOServiceItem, *validate.Errors) {
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
	case primev3messages.UpdateMTOServiceItemModelTypeUpdateMTOServiceItemSIT:
		sit := mtoServiceItem.(*primev3messages.UpdateMTOServiceItemSIT)
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

	case primev3messages.UpdateMTOServiceItemModelTypeUpdateMTOServiceItemShuttle:
		shuttle := mtoServiceItem.(*primev3messages.UpdateMTOServiceItemShuttle)
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

func ServiceRequestDocumentUploadModel(u models.Upload) *primev3messages.UploadWithOmissions {
	return &primev3messages.UploadWithOmissions{
		Bytes:       &u.Bytes,
		ContentType: &u.ContentType,
		Filename:    &u.Filename,
		CreatedAt:   (strfmt.DateTime)(u.CreatedAt),
		UpdatedAt:   (strfmt.DateTime)(u.UpdatedAt),
	}
}

// SITExtensionModel transform the request data the sitExtension model
func SITExtensionModel(sitExtension *primev3messages.CreateSITExtension, mtoShipmentID strfmt.UUID) *models.SITDurationUpdate {
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

func validateDomesticCrating(m primev3messages.MTOServiceItemDomesticCrating) *validate.Errors {
	return validate.Validate(
		&models.ItemCanFitInsideCrateV3{
			Name:         "Item",
			NameCompared: "Crate",
			Item:         &m.Item.MTOServiceItemDimension,
			Crate:        &m.Crate.MTOServiceItemDimension,
		},
	)
}

// validateInternationalCrating validates this mto service item international crating
func validateInternationalCrating(m primev3messages.MTOServiceItemInternationalCrating) *validate.Errors {
	return validate.Validate(
		&models.ItemCanFitInsideCrateV3{
			Name:         "Item",
			NameCompared: "Crate",
			Item:         &m.Item.MTOServiceItemDimension,
			Crate:        &m.Crate.MTOServiceItemDimension,
		},
	)
}

// validateDDFSITForCreate validates DDFSIT service item has all required fields
func validateDDFSITForCreate(m primev3messages.MTOServiceItemDestSIT) *validate.Errors {
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
func validateDestSITForUpdate(m primev3messages.UpdateMTOServiceItemSIT) *validate.Errors {
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
func validateReasonDestSIT(m primev3messages.MTOServiceItemDestSIT) *validate.Errors {
	verrs := validate.NewErrors()

	if m.Reason == nil || m.Reason == models.StringPointer("") {
		verrs.Add("reason", "reason is required in body.")
	}
	return verrs
}

// validateReasonOriginSIT validates that Origin SIT service items have required Reason field
func validateReasonOriginSIT(m primev3messages.MTOServiceItemOriginSIT) *validate.Errors {
	verrs := validate.NewErrors()

	if m.Reason == nil || m.Reason == models.StringPointer("") {
		verrs.Add("reason", "reason is required in body.")
	}
	return verrs
}

// validateBoatShipmentType validates that the shipment type is a valid boat type, and is not nil.
func validateBoatShipmentType(s primev3messages.MTOShipmentType) *validate.Errors {
	verrs := validate.NewErrors()

	if s != primev3messages.MTOShipmentTypeBOATHAULAWAY && s != primev3messages.MTOShipmentTypeBOATTOWAWAY {
		verrs.Add("Boat Shipment Type (mtoShipment.shipmentType)", "shipmentType must be either "+string(primev3messages.BoatShipmentTypeTOWAWAY)+" or "+string(primev3messages.BoatShipmentTypeHAULAWAY))
	}

	return verrs
}
