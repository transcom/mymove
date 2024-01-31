package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcv2messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MTOAgentModel model
func MTOAgentModel(mtoAgent *ghcv2messages.MTOAgent) *models.MTOAgent {
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
func MTOAgentsModel(mtoAgents *ghcv2messages.MTOAgents) *models.MTOAgents {
	if mtoAgents == nil {
		return nil
	}

	agents := make(models.MTOAgents, len(*mtoAgents))

	for i, m := range *mtoAgents {
		agents[i] = *MTOAgentModel(m)
	}

	return &agents
}

// AddressModel model
func AddressModel(address *ghcv2messages.Address) *models.Address {
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
		Country:        address.Country,
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

// StorageFacilityModel model
func StorageFacilityModel(storageFacility *ghcv2messages.StorageFacility) *models.StorageFacility {
	// To check if the model is intended to be blank, we'll look at both ID and FacilityName
	// We should always have ID if the user intends to update a Storage Facility,
	// and FacilityName is a required field on creation. If both are blank, it should be treated as nil.
	var blankSwaggerID strfmt.UUID
	if storageFacility == nil || (storageFacility.ID == blankSwaggerID && storageFacility.FacilityName == "") {
		return nil
	}

	modelStorageFacility := &models.StorageFacility{
		ID:           uuid.FromStringOrNil(storageFacility.ID.String()),
		FacilityName: storageFacility.FacilityName,
		LotNumber:    storageFacility.LotNumber,
		Phone:        storageFacility.Phone,
		Email:        storageFacility.Email,
	}

	addressModel := AddressModel(storageFacility.Address)
	if addressModel != nil {
		modelStorageFacility.Address = *addressModel
	}

	return modelStorageFacility
}

// MTOShipmentModelFromUpdate model
func MTOShipmentModelFromUpdate(mtoShipment *ghcv2messages.UpdateShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	var requestedPickupDate *time.Time
	if mtoShipment.RequestedPickupDate != nil {
		rpd := time.Time(*mtoShipment.RequestedPickupDate)
		requestedPickupDate = &rpd
	}
	var requestedDeliveryDate *time.Time
	if mtoShipment.RequestedDeliveryDate != nil {
		rdd := time.Time(*mtoShipment.RequestedDeliveryDate)
		requestedDeliveryDate = &rdd
	}
	var billableWeightCap *unit.Pound
	if mtoShipment.BillableWeightCap != nil {
		bwc := unit.Pound(*mtoShipment.BillableWeightCap)
		billableWeightCap = &bwc
	}

	var tacType *models.LOAType
	if mtoShipment.TacType.Present {
		tt := models.LOAType(*mtoShipment.TacType.Value)
		tacType = &tt
	}

	var sacType *models.LOAType
	if mtoShipment.SacType.Present {
		tt := models.LOAType(*mtoShipment.SacType.Value)
		sacType = &tt
	}

	var usesExternalVendor bool
	if mtoShipment.UsesExternalVendor != nil {
		usesExternalVendor = *mtoShipment.UsesExternalVendor
	}

	model := &models.MTOShipment{
		BillableWeightCap:           billableWeightCap,
		BillableWeightJustification: mtoShipment.BillableWeightJustification,
		ShipmentType:                models.MTOShipmentType(mtoShipment.ShipmentType),
		RequestedPickupDate:         requestedPickupDate,
		RequestedDeliveryDate:       requestedDeliveryDate,
		CustomerRemarks:             mtoShipment.CustomerRemarks,
		CounselorRemarks:            mtoShipment.CounselorRemarks,
		TACType:                     tacType,
		SACType:                     sacType,
		UsesExternalVendor:          usesExternalVendor,
		ServiceOrderNumber:          mtoShipment.ServiceOrderNumber,
		HasSecondaryPickupAddress:   mtoShipment.HasSecondaryPickupAddress,
		HasSecondaryDeliveryAddress: mtoShipment.HasSecondaryDeliveryAddress,
	}

	model.PickupAddress = AddressModel(&mtoShipment.PickupAddress.Address)
	model.DestinationAddress = AddressModel(&mtoShipment.DestinationAddress.Address)
	if mtoShipment.HasSecondaryPickupAddress != nil {
		if *mtoShipment.HasSecondaryPickupAddress {
			model.SecondaryPickupAddress = AddressModel(&mtoShipment.SecondaryPickupAddress.Address)
		}
	}
	if mtoShipment.HasSecondaryDeliveryAddress != nil {
		if *mtoShipment.HasSecondaryDeliveryAddress {
			model.SecondaryDeliveryAddress = AddressModel(&mtoShipment.SecondaryDeliveryAddress.Address)
		}
	}

	if mtoShipment.DestinationType != nil {
		valDestinationType := models.DestinationType(*mtoShipment.DestinationType)
		model.DestinationType = &valDestinationType
	}

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	if mtoShipment.NtsRecordedWeight != nil {
		ntsRecordedWeight := handlers.PoundPtrFromInt64Ptr(mtoShipment.NtsRecordedWeight)
		model.NTSRecordedWeight = ntsRecordedWeight
	}

	storageFacilityModel := StorageFacilityModel(mtoShipment.StorageFacility)
	if storageFacilityModel != nil {
		model.StorageFacility = storageFacilityModel
	}

	if mtoShipment.PpmShipment != nil {
		model.PPMShipment = PPMShipmentModelFromUpdate(mtoShipment.PpmShipment)
		model.PPMShipment.Shipment = *model
	}

	return model
}

// PPMShipmentModelFromUpdate model
func PPMShipmentModelFromUpdate(ppmShipment *ghcv2messages.UpdatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}
	model := &models.PPMShipment{
		ActualMoveDate:                 (*time.Time)(ppmShipment.ActualMoveDate),
		SecondaryPickupPostalCode:      ppmShipment.SecondaryPickupPostalCode,
		SecondaryDestinationPostalCode: ppmShipment.SecondaryDestinationPostalCode,
		SITExpected:                    ppmShipment.SitExpected,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtInt64PtrToPopPtr(ppmShipment.AdvanceAmountRequested),
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

	var addressModel *models.Address

	if ppmShipment.W2Address != nil {
		addressModel = AddressModel(ppmShipment.W2Address)
		model.W2Address = addressModel
	}

	if ppmShipment.SitLocation != nil {
		sitLocation := models.SITLocationType(*ppmShipment.SitLocation)
		model.SITLocation = &sitLocation
	}

	if ppmShipment.AdvanceStatus != nil {
		advanceStatus := models.PPMAdvanceStatus(*ppmShipment.AdvanceStatus)
		model.AdvanceStatus = &advanceStatus
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
