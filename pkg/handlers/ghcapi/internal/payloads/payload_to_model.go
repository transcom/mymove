package payloads

import (
	"time"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

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

	var address models.Address
	if payload.CurrentAddress != nil {
		address = models.Address{
			ID:             uuid.FromStringOrNil(payload.CurrentAddress.ID.String()),
			StreetAddress1: *payload.CurrentAddress.StreetAddress1,
			StreetAddress2: payload.CurrentAddress.StreetAddress2,
			StreetAddress3: payload.CurrentAddress.StreetAddress3,
			City:           *payload.CurrentAddress.City,
			State:          *payload.CurrentAddress.State,
			PostalCode:     *payload.CurrentAddress.PostalCode,
			Country:        payload.CurrentAddress.Country,
		}
	}

	var backupContacts []models.BackupContact
	if payload.BackupContact != nil {
		backupContacts = []models.BackupContact{{
			Email: *payload.BackupContact.Email,
			Name:  *payload.BackupContact.Name,
			Phone: payload.BackupContact.Phone,
		}}
	}

	return models.ServiceMember{
		ResidentialAddress: &address,
		BackupContacts:     backupContacts,
		FirstName:          &payload.FirstName,
		LastName:           &payload.LastName,
		Suffix:             payload.Suffix,
		MiddleName:         payload.MiddleName,
		PersonalEmail:      payload.Email,
		Telephone:          payload.Phone,
	}
}

// AddressModel model
func AddressModel(address *ghcmessages.Address) *models.Address {
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
func StorageFacilityModel(storageFacility *ghcmessages.StorageFacility) *models.StorageFacility {
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

// ApprovedSITExtensionFromCreate model
func ApprovedSITExtensionFromCreate(sitExtension *ghcmessages.CreateSITExtensionAsTOO, shipmentID strfmt.UUID) *models.SITExtension {
	if sitExtension == nil {
		return nil
	}
	now := time.Now()
	ad := int(*sitExtension.ApprovedDays)
	model := &models.SITExtension{
		MTOShipmentID: uuid.FromStringOrNil(shipmentID.String()),
		RequestReason: models.SITExtensionRequestReason(*sitExtension.RequestReason),
		RequestedDays: int(*sitExtension.ApprovedDays),
		Status:        models.SITExtensionStatusApproved,
		ApprovedDays:  &ad,
		OfficeRemarks: sitExtension.OfficeRemarks,
		DecisionDate:  &now,
	}

	return model
}

// MTOShipmentModelFromCreate model
func MTOShipmentModelFromCreate(mtoShipment *ghcmessages.CreateMTOShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	var tacType *models.LOAType
	if mtoShipment.TacType != nil {
		tt := models.LOAType(*mtoShipment.TacType)
		tacType = &tt
	}

	var sacType *models.LOAType
	if mtoShipment.SacType != nil {
		st := models.LOAType(*mtoShipment.SacType)
		sacType = &st
	}

	var usesExternalVendor bool
	if mtoShipment.UsesExternalVendor != nil {
		usesExternalVendor = *mtoShipment.UsesExternalVendor
	}

	model := &models.MTOShipment{
		MoveTaskOrderID:    uuid.FromStringOrNil(mtoShipment.MoveTaskOrderID.String()),
		Status:             models.MTOShipmentStatusSubmitted,
		CustomerRemarks:    mtoShipment.CustomerRemarks,
		CounselorRemarks:   mtoShipment.CounselorRemarks,
		TACType:            tacType,
		SACType:            sacType,
		UsesExternalVendor: usesExternalVendor,
		ServiceOrderNumber: mtoShipment.ServiceOrderNumber,
	}

	if mtoShipment.ShipmentType != nil {
		model.ShipmentType = models.MTOShipmentType(*mtoShipment.ShipmentType)
	}

	if mtoShipment.RequestedPickupDate != nil {
		model.RequestedPickupDate = swag.Time(time.Time(*mtoShipment.RequestedPickupDate))
	}

	if mtoShipment.RequestedDeliveryDate != nil {
		model.RequestedDeliveryDate = swag.Time(time.Time(*mtoShipment.RequestedDeliveryDate))
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
		valDestinationType := models.DestinationType(*mtoShipment.DestinationType)
		model.DestinationType = &valDestinationType
	}

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	if mtoShipment.NtsRecordedWeight != nil {
		ntsRecordedWeight := unit.Pound(*mtoShipment.NtsRecordedWeight)
		model.NTSRecordedWeight = &ntsRecordedWeight
	}

	storageFacilityModel := StorageFacilityModel(mtoShipment.StorageFacility)
	if storageFacilityModel != nil {
		model.StorageFacility = storageFacilityModel
	}

	if mtoShipment.PpmShipment != nil {
		model.PPMShipment = PPMShipmentModelFromCreate(mtoShipment.PpmShipment)
		model.PPMShipment.Shipment = *model
	}

	return model
}

// PPMShipmentModelFromCreate model
func PPMShipmentModelFromCreate(ppmShipment *ghcmessages.CreatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	model := &models.PPMShipment{
		Status:                         models.PPMShipmentStatusSubmitted,
		SitExpected:                    ppmShipment.SitExpected,
		SecondaryPickupPostalCode:      ppmShipment.SecondaryPickupPostalCode,
		SecondaryDestinationPostalCode: ppmShipment.SecondaryDestinationPostalCode,
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

	if model.SitExpected != nil && *model.SitExpected {
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

	model.EstimatedWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight)

	if model.HasProGear != nil && *model.HasProGear {
		model.ProGearWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight)
		model.SpouseProGearWeight = handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight)
	}

	return model
}

// MTOShipmentModelFromUpdate model
func MTOShipmentModelFromUpdate(mtoShipment *ghcmessages.UpdateShipment) *models.MTOShipment {
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
	}

	model.PickupAddress = AddressModel(&mtoShipment.PickupAddress.Address)
	model.DestinationAddress = AddressModel(&mtoShipment.DestinationAddress.Address)

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

	return model
}
