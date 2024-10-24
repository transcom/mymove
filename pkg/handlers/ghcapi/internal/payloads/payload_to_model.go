package payloads

import (
	"errors"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	timeHHMMFormat = "15:04"
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
	address := AddressModel(&payload.CurrentAddress.Address)
	backupAddress := AddressModel(&payload.BackupAddress.Address)

	var backupContacts []models.BackupContact
	if payload.BackupContact != nil {
		backupContacts = []models.BackupContact{{
			Email: *payload.BackupContact.Email,
			Name:  *payload.BackupContact.Name,
			Phone: payload.BackupContact.Phone,
		}}
	}

	return models.ServiceMember{
		ResidentialAddress:   address,
		BackupContacts:       backupContacts,
		FirstName:            &payload.FirstName,
		LastName:             &payload.LastName,
		Suffix:               payload.Suffix,
		MiddleName:           payload.MiddleName,
		PersonalEmail:        payload.Email,
		Telephone:            payload.Phone,
		SecondaryTelephone:   payload.SecondaryTelephone,
		PhoneIsPreferred:     &payload.PhoneIsPreferred,
		EmailIsPreferred:     &payload.EmailIsPreferred,
		BackupMailingAddress: backupAddress,
		CacValidated:         payload.CacValidated,
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

// AddressModel model
func PPMDestinationAddressModel(address *ghcmessages.PPMDestinationAddress) *models.Address {
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
func ApprovedSITExtensionFromCreate(sitExtension *ghcmessages.CreateApprovedSITDurationUpdate, shipmentID strfmt.UUID) *models.SITDurationUpdate {
	if sitExtension == nil {
		return nil
	}
	now := time.Now()
	ad := int(*sitExtension.ApprovedDays)
	model := &models.SITDurationUpdate{
		MTOShipmentID: uuid.FromStringOrNil(shipmentID.String()),
		RequestReason: models.SITDurationUpdateRequestReason(*sitExtension.RequestReason),
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
		model.RequestedPickupDate = models.TimePointer(time.Time(*mtoShipment.RequestedPickupDate))
	}

	if mtoShipment.RequestedDeliveryDate != nil {
		model.RequestedDeliveryDate = models.TimePointer(time.Time(*mtoShipment.RequestedDeliveryDate))
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

	addressModel = AddressModel(&mtoShipment.SecondaryPickupAddress.Address)
	if addressModel != nil {
		model.SecondaryPickupAddress = addressModel
	}

	addressModel = AddressModel(&mtoShipment.SecondaryDeliveryAddress.Address)
	if addressModel != nil {
		model.SecondaryDeliveryAddress = addressModel
	}

	addressModel = AddressModel(&mtoShipment.TertiaryPickupAddress.Address)
	if addressModel != nil {
		model.TertiaryPickupAddress = addressModel
	}

	addressModel = AddressModel(&mtoShipment.TertiaryDeliveryAddress.Address)
	if addressModel != nil {
		model.TertiaryDeliveryAddress = addressModel
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
	} else if mtoShipment.BoatShipment != nil {
		model.BoatShipment = BoatShipmentModelFromCreate(mtoShipment.BoatShipment)
		model.BoatShipment.Shipment = *model
	} else if mtoShipment.MobileHomeShipment != nil {
		model.MobileHome = MobileHomeShipmentModelFromCreate(mtoShipment.MobileHomeShipment)
		model.MobileHome.Shipment = *model
	}

	return model
}

// PPMShipmentModelFromCreate model
func PPMShipmentModelFromCreate(ppmShipment *ghcmessages.CreatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
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
		model.TertiaryDestinationAddress = addressModel
		model.HasTertiaryDestinationAddress = handlers.FmtBool(true)
	}

	if ppmShipment.IsActualExpenseReimbursement != nil {
		model.IsActualExpenseReimbursement = ppmShipment.IsActualExpenseReimbursement
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

// BoatShipmentModelFromCreate model
func BoatShipmentModelFromCreate(boatShipment *ghcmessages.CreateBoatShipment) *models.BoatShipment {
	if boatShipment == nil {
		return nil
	}
	var year *int
	if boatShipment.Year != nil {
		val := int(*boatShipment.Year)
		year = &val
	}
	var lengthInInches *int
	if boatShipment.LengthInInches != nil {
		val := int(*boatShipment.LengthInInches)
		lengthInInches = &val
	}
	var widthInInches *int
	if boatShipment.WidthInInches != nil {
		val := int(*boatShipment.WidthInInches)
		widthInInches = &val
	}
	var heightInInches *int
	if boatShipment.HeightInInches != nil {
		val := int(*boatShipment.HeightInInches)
		heightInInches = &val
	}

	model := &models.BoatShipment{
		Type:           models.BoatShipmentType(*boatShipment.Type),
		Year:           year,
		Make:           boatShipment.Make,
		Model:          boatShipment.Model,
		LengthInInches: lengthInInches,
		WidthInInches:  widthInInches,
		HeightInInches: heightInInches,
		HasTrailer:     boatShipment.HasTrailer,
		IsRoadworthy:   boatShipment.IsRoadworthy,
	}

	if model.HasTrailer == models.BoolPointer(false) {
		model.IsRoadworthy = nil
	}

	return model
}

// MobileHomeShipmentModelFromCreate model
func MobileHomeShipmentModelFromCreate(mobileHomeShipment *ghcmessages.CreateMobileHomeShipment) *models.MobileHome {
	if mobileHomeShipment == nil {
		return nil
	}
	var year *int
	if mobileHomeShipment.Year != nil {
		val := int(*mobileHomeShipment.Year)
		year = &val
	}
	var lengthInInches *int
	if mobileHomeShipment.LengthInInches != nil {
		val := int(*mobileHomeShipment.LengthInInches)
		lengthInInches = &val
	}
	var heightInInches *int
	if mobileHomeShipment.HeightInInches != nil {
		val := int(*mobileHomeShipment.HeightInInches)
		heightInInches = &val
	}
	var widthInInches *int
	if mobileHomeShipment.WidthInInches != nil {
		val := int(*mobileHomeShipment.WidthInInches)
		widthInInches = &val
	}

	model := &models.MobileHome{
		Make:           mobileHomeShipment.Make,
		Model:          mobileHomeShipment.Model,
		Year:           year,
		LengthInInches: lengthInInches,
		HeightInInches: heightInInches,
		WidthInInches:  widthInInches,
	}

	return model
}

func CustomerSupportRemarkModelFromCreate(remark *ghcmessages.CreateCustomerSupportRemark) *models.CustomerSupportRemark {
	if remark == nil {
		return nil
	}

	model := &models.CustomerSupportRemark{
		Content:      *remark.Content,
		OfficeUserID: uuid.FromStringOrNil(remark.OfficeUserID.String()),
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
		HasSecondaryPickupAddress:   mtoShipment.HasSecondaryPickupAddress,
		HasSecondaryDeliveryAddress: mtoShipment.HasSecondaryDeliveryAddress,
		HasTertiaryPickupAddress:    mtoShipment.HasTertiaryPickupAddress,
		HasTertiaryDeliveryAddress:  mtoShipment.HasTertiaryDeliveryAddress,
		ActualProGearWeight:         handlers.PoundPtrFromInt64Ptr(mtoShipment.ActualProGearWeight),
		ActualSpouseProGearWeight:   handlers.PoundPtrFromInt64Ptr(mtoShipment.ActualSpouseProGearWeight),
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

	if mtoShipment.HasTertiaryPickupAddress != nil {
		if *mtoShipment.HasTertiaryPickupAddress {
			model.TertiaryPickupAddress = AddressModel(&mtoShipment.TertiaryPickupAddress.Address)
		}
	}
	if mtoShipment.HasTertiaryDeliveryAddress != nil {
		if *mtoShipment.HasTertiaryDeliveryAddress {
			model.TertiaryDeliveryAddress = AddressModel(&mtoShipment.TertiaryDeliveryAddress.Address)
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

	// making sure both shipmentType and boatShipment.Type match
	if mtoShipment.BoatShipment != nil && mtoShipment.BoatShipment.Type != nil {
		if *mtoShipment.BoatShipment.Type == string(models.BoatShipmentTypeHaulAway) {
			model.ShipmentType = models.MTOShipmentTypeBoatHaulAway
		} else {
			model.ShipmentType = models.MTOShipmentTypeBoatTowAway
		}
		model.BoatShipment = BoatShipmentModelFromUpdate(mtoShipment.BoatShipment)
		model.BoatShipment.Shipment = *model
	}

	if mtoShipment.MobileHomeShipment != nil {
		model.MobileHome = MobileHomeShipmentModelFromUpdate(mtoShipment.MobileHomeShipment)
		model.MobileHome.Shipment = *model
	}

	return model
}

// PPMShipmentModelFromUpdate model
func PPMShipmentModelFromUpdate(ppmShipment *ghcmessages.UpdatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}
	model := &models.PPMShipment{
		ActualMoveDate:                 (*time.Time)(ppmShipment.ActualMoveDate),
		SITExpected:                    ppmShipment.SitExpected,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		AllowableWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.AllowableWeight),
		HasProGear:                     ppmShipment.HasProGear,
		IsActualExpenseReimbursement:   ppmShipment.IsActualExpenseReimbursement,
		ProGearWeight:                  handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtInt64PtrToPopPtr(ppmShipment.AdvanceAmountRequested),
		HasSecondaryPickupAddress:      ppmShipment.HasSecondaryPickupAddress,
		HasSecondaryDestinationAddress: ppmShipment.HasSecondaryDestinationAddress,
		HasTertiaryPickupAddress:       ppmShipment.HasTertiaryPickupAddress,
		HasTertiaryDestinationAddress:  ppmShipment.HasTertiaryDestinationAddress,
		AdvanceAmountReceived:          handlers.FmtInt64PtrToPopPtr(ppmShipment.AdvanceAmountReceived),
		HasReceivedAdvance:             ppmShipment.HasReceivedAdvance,
		ActualPickupPostalCode:         ppmShipment.ActualPickupPostalCode,
		ActualDestinationPostalCode:    ppmShipment.ActualDestinationPostalCode,
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
		secondaryPickupAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.SecondaryPickupAddressID = &secondaryPickupAddressID
	}

	addressModel = AddressModel(&ppmShipment.TertiaryPickupAddress.Address)
	if addressModel != nil {
		model.TertiaryPickupAddress = addressModel
		tertiaryPickupAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.TertiaryPickupAddressID = &tertiaryPickupAddressID
	}

	addressModel = PPMDestinationAddressModel(&ppmShipment.DestinationAddress.PPMDestinationAddress)
	if addressModel != nil {
		model.DestinationAddress = addressModel
	}

	addressModel = AddressModel(&ppmShipment.SecondaryDestinationAddress.Address)
	if addressModel != nil {
		model.SecondaryDestinationAddress = addressModel
		secondaryDestinationAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.SecondaryDestinationAddressID = &secondaryDestinationAddressID
	}

	addressModel = AddressModel(&ppmShipment.TertiaryDestinationAddress.Address)
	if addressModel != nil {
		model.TertiaryDestinationAddress = addressModel
		tertiaryDestinationAddressID := uuid.FromStringOrNil(addressModel.ID.String())
		model.TertiaryDestinationAddressID = &tertiaryDestinationAddressID
	}

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

// BoatShipmentModelFromUpdate model
func BoatShipmentModelFromUpdate(boatShipment *ghcmessages.UpdateBoatShipment) *models.BoatShipment {
	if boatShipment == nil {
		return nil
	}
	var year *int
	if boatShipment.Year != nil {
		val := int(*boatShipment.Year)
		year = &val
	}
	var lengthInInches *int
	if boatShipment.LengthInInches != nil {
		val := int(*boatShipment.LengthInInches)
		lengthInInches = &val
	}
	var widthInInches *int
	if boatShipment.WidthInInches != nil {
		val := int(*boatShipment.WidthInInches)
		widthInInches = &val
	}
	var heightInInches *int
	if boatShipment.HeightInInches != nil {
		val := int(*boatShipment.HeightInInches)
		heightInInches = &val
	}

	boatModel := &models.BoatShipment{
		Year:           year,
		Make:           boatShipment.Make,
		Model:          boatShipment.Model,
		LengthInInches: lengthInInches,
		WidthInInches:  widthInInches,
		HeightInInches: heightInInches,
		HasTrailer:     boatShipment.HasTrailer,
		IsRoadworthy:   boatShipment.IsRoadworthy,
	}

	if boatShipment.Type != nil {
		boatModel.Type = models.BoatShipmentType(*boatShipment.Type)
	}

	if boatShipment.HasTrailer == models.BoolPointer(false) {
		boatModel.IsRoadworthy = nil
	}

	return boatModel
}

func MobileHomeShipmentModelFromUpdate(mobileHomeShipment *ghcmessages.UpdateMobileHomeShipment) *models.MobileHome {
	if mobileHomeShipment == nil {
		return nil
	}
	var year *int
	if mobileHomeShipment.Year != nil {
		val := int(*mobileHomeShipment.Year)
		year = &val
	}
	var lengthInInches *int
	if mobileHomeShipment.LengthInInches != nil {
		val := int(*mobileHomeShipment.LengthInInches)
		lengthInInches = &val
	}
	var heightInInches *int
	if mobileHomeShipment.HeightInInches != nil {
		val := int(*mobileHomeShipment.HeightInInches)
		heightInInches = &val
	}

	var widthInInches *int
	if mobileHomeShipment.WidthInInches != nil {
		val := int(*mobileHomeShipment.WidthInInches)
		widthInInches = &val
	}

	mobileHomeModel := &models.MobileHome{
		Make:           mobileHomeShipment.Make,
		Model:          mobileHomeShipment.Model,
		Year:           year,
		LengthInInches: lengthInInches,
		HeightInInches: heightInInches,
		WidthInInches:  widthInInches,
	}

	return mobileHomeModel
}

// ProgearWeightTicketModelFromUpdate model
func ProgearWeightTicketModelFromUpdate(progearWeightTicket *ghcmessages.UpdateProGearWeightTicket) *models.ProgearWeightTicket {
	if progearWeightTicket == nil {
		return nil
	}

	model := &models.ProgearWeightTicket{
		Weight:           handlers.PoundPtrFromInt64Ptr(progearWeightTicket.Weight),
		HasWeightTickets: handlers.FmtBool(progearWeightTicket.HasWeightTickets),
		BelongsToSelf:    handlers.FmtBool(progearWeightTicket.BelongsToSelf),
		Status:           (*models.PPMDocumentStatus)(handlers.FmtString(string(progearWeightTicket.Status))),
		Reason:           handlers.FmtString(progearWeightTicket.Reason),
	}
	return model
}

// WeightTicketModelFromUpdate
func WeightTicketModelFromUpdate(weightTicket *ghcmessages.UpdateWeightTicket) *models.WeightTicket {
	if weightTicket == nil {
		return nil
	}
	model := &models.WeightTicket{
		EmptyWeight:          handlers.PoundPtrFromInt64Ptr(weightTicket.EmptyWeight),
		FullWeight:           handlers.PoundPtrFromInt64Ptr(weightTicket.FullWeight),
		OwnsTrailer:          handlers.FmtBool(weightTicket.OwnsTrailer),
		TrailerMeetsCriteria: handlers.FmtBool(weightTicket.TrailerMeetsCriteria),
		Status:               (*models.PPMDocumentStatus)(handlers.FmtString(string(weightTicket.Status))),
		Reason:               handlers.FmtString(weightTicket.Reason),
		AdjustedNetWeight:    handlers.PoundPtrFromInt64Ptr(weightTicket.AdjustedNetWeight),
		NetWeightRemarks:     handlers.FmtString(weightTicket.NetWeightRemarks),
	}
	return model
}

// MovingExpenseModelFromUpdate
func MovingExpenseModelFromUpdate(movingExpense *ghcmessages.UpdateMovingExpense) *models.MovingExpense {
	var model models.MovingExpense

	if movingExpense == nil {
		return nil
	}

	var expenseType models.MovingExpenseReceiptType
	if movingExpense.MovingExpenseType != nil {
		expenseType = models.MovingExpenseReceiptType(*movingExpense.MovingExpenseType.Pointer())
		model.MovingExpenseType = &expenseType
	}

	if movingExpense.Description != nil {
		model.Description = movingExpense.Description
	}

	if movingExpense.SitLocation != nil {
		model.SITLocation = (*models.SITLocationType)(handlers.FmtString(string(*movingExpense.SitLocation)))
	}

	model.Amount = handlers.FmtInt64PtrToPopPtr(&movingExpense.Amount)
	model.SITStartDate = handlers.FmtDatePtrToPopPtr(&movingExpense.SitStartDate)
	model.SITEndDate = handlers.FmtDatePtrToPopPtr(&movingExpense.SitEndDate)
	model.Status = (*models.PPMDocumentStatus)(handlers.FmtString(string(movingExpense.Status)))
	model.Reason = handlers.FmtString(movingExpense.Reason)
	model.WeightStored = handlers.PoundPtrFromInt64Ptr(&movingExpense.WeightStored)
	model.SITEstimatedCost = handlers.FmtInt64PtrToPopPtr(movingExpense.SitEstimatedCost)
	model.SITReimburseableAmount = handlers.FmtInt64PtrToPopPtr(movingExpense.SitReimburseableAmount)

	return &model
}

func EvaluationReportFromUpdate(evaluationReport *ghcmessages.EvaluationReport) (*models.EvaluationReport, error) {
	if evaluationReport == nil {
		err := apperror.NewPreconditionFailedError(uuid.UUID{}, errors.New("cannot update empty report"))
		return nil, err
	}

	var inspectionType *models.EvaluationReportInspectionType
	if evaluationReport.InspectionType != nil {
		tempInspectionType := models.EvaluationReportInspectionType(*evaluationReport.InspectionType)
		inspectionType = &tempInspectionType
	}

	var location *models.EvaluationReportLocationType
	if evaluationReport.Location != nil {
		tempLocation := models.EvaluationReportLocationType(*evaluationReport.Location)
		location = &tempLocation
	}

	var timeDepart *time.Time
	if evaluationReport.TimeDepart != nil {
		td, err := time.Parse(timeHHMMFormat, *evaluationReport.TimeDepart)

		if err != nil {
			return nil, apperror.NewPreconditionFailedError(uuid.UUID{}, err)
		}

		timeDepart = &td
	}

	var evalStart *time.Time
	if evaluationReport.EvalStart != nil {
		es, err := time.Parse(timeHHMMFormat, *evaluationReport.EvalStart)

		if err != nil {
			return nil, apperror.NewPreconditionFailedError(uuid.UUID{}, err)
		}

		evalStart = &es
	}

	var evalEnd *time.Time
	if evaluationReport.EvalEnd != nil {
		ee, err := time.Parse(timeHHMMFormat, *evaluationReport.EvalEnd)

		if err != nil {
			return nil, apperror.NewPreconditionFailedError(uuid.UUID{}, err)
		}

		evalEnd = &ee
	}

	model := models.EvaluationReport{
		ID:                                 uuid.FromStringOrNil(evaluationReport.ID.String()),
		OfficeUser:                         models.OfficeUser{},
		OfficeUserID:                       uuid.Nil,
		Move:                               models.Move{},
		MoveID:                             uuid.Nil,
		Shipment:                           nil,
		ShipmentID:                         nil,
		Type:                               models.EvaluationReportType(evaluationReport.Type),
		InspectionDate:                     (*time.Time)(evaluationReport.InspectionDate),
		InspectionType:                     inspectionType,
		TimeDepart:                         timeDepart,
		EvalStart:                          evalStart,
		EvalEnd:                            evalEnd,
		Location:                           location,
		LocationDescription:                evaluationReport.LocationDescription,
		ObservedShipmentDeliveryDate:       (*time.Time)(evaluationReport.ObservedShipmentDeliveryDate),
		ObservedShipmentPhysicalPickupDate: (*time.Time)(evaluationReport.ObservedShipmentPhysicalPickupDate),
		ViolationsObserved:                 evaluationReport.ViolationsObserved,
		Remarks:                            evaluationReport.Remarks,
		SeriousIncident:                    evaluationReport.SeriousIncident,
		SeriousIncidentDesc:                evaluationReport.SeriousIncidentDesc,
		ObservedClaimsResponseDate:         (*time.Time)(evaluationReport.ObservedClaimsResponseDate),
		ObservedPickupDate:                 (*time.Time)(evaluationReport.ObservedPickupDate),
		ObservedPickupSpreadStartDate:      (*time.Time)(evaluationReport.ObservedPickupSpreadStartDate),
		ObservedPickupSpreadEndDate:        (*time.Time)(evaluationReport.ObservedPickupSpreadEndDate),
		ObservedDeliveryDate:               (*time.Time)(evaluationReport.ObservedDeliveryDate),
		SubmittedAt:                        handlers.FmtDateTimePtrToPopPtr(evaluationReport.SubmittedAt),
	}
	return &model, nil
}
