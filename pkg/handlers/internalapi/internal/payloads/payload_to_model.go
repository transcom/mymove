package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// AddressModel model
func AddressModel(address *internalmessages.Address) *models.Address {
	if address == nil {
		return nil
	}
	if address.County == nil {
		address.County = models.StringPointer("")
	}

	usprcID := uuid.FromStringOrNil(address.UsprcID.String())

	return &models.Address{
		ID:                 uuid.FromStringOrNil(address.ID.String()),
		StreetAddress1:     *address.StreetAddress1,
		StreetAddress2:     address.StreetAddress2,
		StreetAddress3:     address.StreetAddress3,
		City:               *address.City,
		State:              *address.State,
		PostalCode:         *address.PostalCode,
		County:             *address.County,
		UsPostRegionCityId: &usprcID,
	}
}

func VLocationModel(vLocation *internalmessages.VLocation) *models.VLocation {
	if vLocation == nil {
		return nil
	}

	usprcID := uuid.FromStringOrNil(vLocation.UsPostRegionCitiesID.String())

	return &models.VLocation{
		CityName:      vLocation.City,
		StateName:     vLocation.State,
		UsprZipID:     vLocation.PostalCode,
		UsprcCountyNm: *vLocation.County,
		UprcId:        &usprcID,
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

	model := &models.MTOShipment{
		MoveTaskOrderID: uuid.FromStringOrNil(mtoShipment.MoveTaskOrderID.String()),
		CustomerRemarks: mtoShipment.CustomerRemarks,
		ShipmentType:    models.MTOShipmentType(*mtoShipment.ShipmentType),
	}
	isBoatShipment := model.ShipmentType == models.MTOShipmentTypeBoatHaulAway || model.ShipmentType == models.MTOShipmentTypeBoatTowAway
	isMobileHomeShipment := model.ShipmentType == models.MTOShipmentTypeMobileHome
	// PPM and Boat type shipment begins in DRAFT because it requires a multi-page series to complete.
	// After move submission a the status will change to SUBMITTED
	if model.ShipmentType == models.MTOShipmentTypePPM || isBoatShipment || isMobileHomeShipment {
		model.Status = models.MTOShipmentStatusDraft
	} else {
		model.Status = models.MTOShipmentStatusSubmitted
	}

	requestedPickupDate := time.Time(mtoShipment.RequestedPickupDate)
	if !requestedPickupDate.IsZero() {
		model.RequestedPickupDate = &requestedPickupDate
	}

	requestedDeliveryDate := time.Time(mtoShipment.RequestedDeliveryDate)
	if !requestedDeliveryDate.IsZero() {
		model.RequestedDeliveryDate = &requestedDeliveryDate
	}

	model.PickupAddress = AddressModel(mtoShipment.PickupAddress)

	model.SecondaryPickupAddress = AddressModel(mtoShipment.SecondaryPickupAddress)
	model.HasSecondaryPickupAddress = handlers.FmtBool(mtoShipment.SecondaryPickupAddress != nil)

	model.TertiaryPickupAddress = AddressModel(mtoShipment.TertiaryPickupAddress)
	model.HasTertiaryPickupAddress = handlers.FmtBool(mtoShipment.TertiaryPickupAddress != nil)

	model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)

	model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)
	model.HasSecondaryDeliveryAddress = handlers.FmtBool(mtoShipment.SecondaryDeliveryAddress != nil)

	model.TertiaryDeliveryAddress = AddressModel(mtoShipment.TertiaryDeliveryAddress)
	model.HasTertiaryDeliveryAddress = handlers.FmtBool(mtoShipment.TertiaryDeliveryAddress != nil)

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
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
func PPMShipmentModelFromCreate(ppmShipment *internalmessages.CreatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	model := &models.PPMShipment{
		SITExpected:           ppmShipment.SitExpected,
		ExpectedDepartureDate: handlers.FmtDatePtrToPop(ppmShipment.ExpectedDepartureDate),
	}

	if ppmShipment.PickupAddress != nil {
		model.PickupAddress = AddressModel(ppmShipment.PickupAddress)
	}

	model.HasSecondaryPickupAddress = handlers.FmtBool(ppmShipment.SecondaryPickupAddress != nil)
	if ppmShipment.SecondaryPickupAddress != nil {
		model.SecondaryPickupAddress = AddressModel(ppmShipment.SecondaryPickupAddress)
	}

	model.HasTertiaryPickupAddress = handlers.FmtBool(ppmShipment.TertiaryPickupAddress != nil)
	if ppmShipment.TertiaryPickupAddress != nil {
		model.TertiaryPickupAddress = AddressModel(ppmShipment.TertiaryPickupAddress)
	}

	if ppmShipment.DestinationAddress != nil {
		model.DestinationAddress = AddressModel(ppmShipment.DestinationAddress)
	}

	model.HasSecondaryDestinationAddress = handlers.FmtBool(ppmShipment.SecondaryDestinationAddress != nil)
	if ppmShipment.SecondaryDestinationAddress != nil {
		model.SecondaryDestinationAddress = AddressModel(ppmShipment.SecondaryDestinationAddress)
	}

	model.HasTertiaryDestinationAddress = handlers.FmtBool(ppmShipment.TertiaryDestinationAddress != nil)
	if ppmShipment.TertiaryDestinationAddress != nil {
		model.TertiaryDestinationAddress = AddressModel(ppmShipment.TertiaryDestinationAddress)
	}

	if ppmShipment.IsActualExpenseReimbursement != nil {
		model.IsActualExpenseReimbursement = ppmShipment.IsActualExpenseReimbursement
	}

	return model
}

func UpdatePPMShipmentModel(ppmShipment *internalmessages.UpdatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	ppmModel := &models.PPMShipment{
		ActualMoveDate:                 (*time.Time)(ppmShipment.ActualMoveDate),
		ActualPickupPostalCode:         ppmShipment.ActualPickupPostalCode,
		ActualDestinationPostalCode:    ppmShipment.ActualDestinationPostalCode,
		SITExpected:                    ppmShipment.SitExpected,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
		IsActualExpenseReimbursement:   ppmShipment.IsActualExpenseReimbursement,
		ProGearWeight:                  handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtInt64PtrToPopPtr(ppmShipment.AdvanceAmountRequested),
		HasReceivedAdvance:             ppmShipment.HasReceivedAdvance,
		AdvanceAmountReceived:          handlers.FmtInt64PtrToPopPtr(ppmShipment.AdvanceAmountReceived),
		FinalIncentive:                 handlers.FmtInt64PtrToPopPtr(ppmShipment.FinalIncentive),
		HasSecondaryPickupAddress:      ppmShipment.HasSecondaryPickupAddress,
		HasSecondaryDestinationAddress: ppmShipment.HasSecondaryDestinationAddress,
		HasTertiaryPickupAddress:       ppmShipment.HasTertiaryPickupAddress,
		HasTertiaryDestinationAddress:  ppmShipment.HasTertiaryDestinationAddress,
	}

	ppmModel.W2Address = AddressModel(ppmShipment.W2Address)
	if ppmShipment.ExpectedDepartureDate != nil {
		ppmModel.ExpectedDepartureDate = *handlers.FmtDatePtrToPopPtr(ppmShipment.ExpectedDepartureDate)
	}

	if ppmShipment.PickupAddress != nil {
		ppmModel.PickupAddress = AddressModel(ppmShipment.PickupAddress)
	}

	if ppmShipment.SecondaryPickupAddress != nil {
		ppmModel.SecondaryPickupAddress = AddressModel(ppmShipment.SecondaryPickupAddress)
	}

	if ppmShipment.TertiaryPickupAddress != nil {
		ppmModel.TertiaryPickupAddress = AddressModel(ppmShipment.TertiaryPickupAddress)
	}

	if ppmShipment.DestinationAddress != nil {
		ppmModel.DestinationAddress = AddressModel(ppmShipment.DestinationAddress)
	}

	if ppmShipment.SecondaryDestinationAddress != nil {
		ppmModel.SecondaryDestinationAddress = AddressModel(ppmShipment.SecondaryDestinationAddress)
	}

	if ppmShipment.TertiaryDestinationAddress != nil {
		ppmModel.TertiaryDestinationAddress = AddressModel(ppmShipment.TertiaryDestinationAddress)
	}

	if ppmShipment.FinalIncentive != nil {
		ppmModel.FinalIncentive = handlers.FmtInt64PtrToPopPtr(ppmShipment.FinalIncentive)
	}

	return ppmModel
}

// BoatShipmentModelFromCreate model
func BoatShipmentModelFromCreate(boatShipment *internalmessages.CreateBoatShipment) *models.BoatShipment {
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

func UpdateBoatShipmentModel(boatShipment *internalmessages.UpdateBoatShipment) *models.BoatShipment {
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

// MobileHomeShipmentModelFromCreate model
func MobileHomeShipmentModelFromCreate(mobileHomeShipment *internalmessages.CreateMobileHomeShipment) *models.MobileHome {
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

func UpdateMobileHomeShipmentModel(mobileHomeShipment *internalmessages.UpdateMobileHomeShipment) *models.MobileHome {
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

// MTOShipmentModelFromUpdate model
func MTOShipmentModelFromUpdate(mtoShipment *internalmessages.UpdateShipment) *models.MTOShipment {
	if mtoShipment == nil {
		return nil
	}

	var requestedPickupDate, requestedDeliveryDate *time.Time
	if mtoShipment.RequestedPickupDate != nil {
		date := time.Time(*mtoShipment.RequestedPickupDate)
		requestedPickupDate = &date
	}

	if mtoShipment.RequestedDeliveryDate != nil {
		date := time.Time(*mtoShipment.RequestedDeliveryDate)
		requestedDeliveryDate = &date
	}

	model := &models.MTOShipment{
		ShipmentType:                models.MTOShipmentType(mtoShipment.ShipmentType),
		RequestedPickupDate:         requestedPickupDate,
		RequestedDeliveryDate:       requestedDeliveryDate,
		CustomerRemarks:             mtoShipment.CustomerRemarks,
		Status:                      models.MTOShipmentStatus(mtoShipment.Status),
		HasSecondaryPickupAddress:   mtoShipment.HasSecondaryPickupAddress,
		HasSecondaryDeliveryAddress: mtoShipment.HasSecondaryDeliveryAddress,
		HasTertiaryPickupAddress:    mtoShipment.HasTertiaryPickupAddress,
		HasTertiaryDeliveryAddress:  mtoShipment.HasTertiaryDeliveryAddress,
		ActualProGearWeight:         handlers.PoundPtrFromInt64Ptr(mtoShipment.ActualProGearWeight),
		ActualSpouseProGearWeight:   handlers.PoundPtrFromInt64Ptr(mtoShipment.ActualSpouseProGearWeight),
	}

	model.PickupAddress = AddressModel(mtoShipment.PickupAddress)
	if mtoShipment.HasSecondaryPickupAddress != nil {
		if *mtoShipment.HasSecondaryPickupAddress {
			model.SecondaryPickupAddress = AddressModel(mtoShipment.SecondaryPickupAddress)
		}
	}
	if mtoShipment.HasTertiaryPickupAddress != nil {
		if *mtoShipment.HasTertiaryPickupAddress {
			model.TertiaryPickupAddress = AddressModel(mtoShipment.TertiaryPickupAddress)
		}
	}
	if mtoShipment.HasSecondaryDeliveryAddress != nil {
		if *mtoShipment.HasSecondaryDeliveryAddress {
			model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)
		}
	}
	if mtoShipment.HasTertiaryDeliveryAddress != nil {
		if *mtoShipment.HasTertiaryDeliveryAddress {
			model.TertiaryDeliveryAddress = AddressModel(mtoShipment.TertiaryDeliveryAddress)
		}
	}
	model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)
	model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	model.PPMShipment = UpdatePPMShipmentModel(mtoShipment.PpmShipment)

	// making sure both shipmentType and boatShipment.Type match
	if mtoShipment.BoatShipment != nil && mtoShipment.BoatShipment.Type != nil {
		if *mtoShipment.BoatShipment.Type == string(models.BoatShipmentTypeHaulAway) {
			model.ShipmentType = models.MTOShipmentTypeBoatHaulAway
		} else {
			model.ShipmentType = models.MTOShipmentTypeBoatTowAway
		}
	}
	model.BoatShipment = UpdateBoatShipmentModel(mtoShipment.BoatShipment)

	model.MobileHome = UpdateMobileHomeShipmentModel(mtoShipment.MobileHomeShipment)

	return model
}

// MovingExpenseModelFromUpdate
func MovingExpenseModelFromUpdate(movingExpense *internalmessages.UpdateMovingExpense) *models.MovingExpense {
	if movingExpense == nil {
		return nil
	}
	model := &models.MovingExpense{
		MovingExpenseType: (*models.MovingExpenseReceiptType)(movingExpense.MovingExpenseType),
		Description:       handlers.FmtStringPtr(movingExpense.Description),
		Amount:            handlers.FmtInt64PtrToPopPtr(movingExpense.Amount),
		SITStartDate:      handlers.FmtDatePtrToPopPtr(&movingExpense.SitStartDate),
		SITEndDate:        handlers.FmtDatePtrToPopPtr(&movingExpense.SitEndDate),
		WeightStored:      handlers.PoundPtrFromInt64Ptr(&movingExpense.WeightStored),
	}

	if movingExpense.PaidWithGTCC != nil {
		model.PaidWithGTCC = handlers.FmtBool(*movingExpense.PaidWithGTCC)
	}

	if movingExpense.MissingReceipt != nil {
		model.MissingReceipt = handlers.FmtBool(*movingExpense.MissingReceipt)
	}

	if movingExpense.SitLocation != nil {
		model.SITLocation = (*models.SITLocationType)(handlers.FmtString(string(*movingExpense.SitLocation)))
	}

	if movingExpense.SitReimburseableAmount != nil {
		model.SITReimburseableAmount = handlers.FmtInt64PtrToPopPtr(movingExpense.SitReimburseableAmount)
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

	requestedPickupDate := time.Time(*mtoShipment.RequestedPickupDate)
	if !requestedPickupDate.IsZero() {
		model.RequestedPickupDate = &requestedPickupDate
	}

	requestedDeliveryDate := time.Time(*mtoShipment.RequestedDeliveryDate)
	if !requestedDeliveryDate.IsZero() {
		model.RequestedDeliveryDate = &requestedDeliveryDate
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

// WeightTicketModelFromUpdate
func WeightTicketModelFromUpdate(weightTicket *internalmessages.UpdateWeightTicket) *models.WeightTicket {
	if weightTicket == nil {
		return nil
	}
	model := &models.WeightTicket{
		VehicleDescription:       &weightTicket.VehicleDescription,
		EmptyWeight:              handlers.PoundPtrFromInt64Ptr(weightTicket.EmptyWeight),
		MissingEmptyWeightTicket: handlers.FmtBool(weightTicket.MissingEmptyWeightTicket),
		FullWeight:               handlers.PoundPtrFromInt64Ptr(weightTicket.FullWeight),
		MissingFullWeightTicket:  handlers.FmtBool(weightTicket.MissingFullWeightTicket),
		OwnsTrailer:              handlers.FmtBool(weightTicket.OwnsTrailer),
		TrailerMeetsCriteria:     handlers.FmtBool(weightTicket.TrailerMeetsCriteria),
		AdjustedNetWeight:        handlers.PoundPtrFromInt64Ptr(weightTicket.AdjustedNetWeight),
		NetWeightRemarks:         handlers.FmtString(weightTicket.NetWeightRemarks),
	}
	return model
}

// ProgearWeightTicketModelFromUpdate
func ProgearWeightTicketModelFromUpdate(progearWeightTicket *internalmessages.UpdateProGearWeightTicket) *models.ProgearWeightTicket {
	if progearWeightTicket == nil {
		return nil
	}
	model := &models.ProgearWeightTicket{
		Description:      &progearWeightTicket.Description,
		Weight:           handlers.PoundPtrFromInt64Ptr(progearWeightTicket.Weight),
		HasWeightTickets: handlers.FmtBool(progearWeightTicket.HasWeightTickets),
		BelongsToSelf:    handlers.FmtBool(progearWeightTicket.BelongsToSelf),
	}
	return model
}

// SavePPMShipmentSignedCertification converts from the SavePPMShipmentSignedCertification payload and the
// SignedCertification model
func SavePPMShipmentSignedCertification(ppmShipmentID uuid.UUID, signedCertification internalmessages.SavePPMShipmentSignedCertification) models.SignedCertification {
	model := models.SignedCertification{
		PpmID: &ppmShipmentID,
		Date:  handlers.FmtDatePtrToPop(signedCertification.Date),
	}

	if signedCertification.CertificationText != nil {
		model.CertificationText = *signedCertification.CertificationText
	}

	if signedCertification.Signature != nil {
		model.Signature = *signedCertification.Signature
	}

	return model
}

// ReSavePPMShipmentSignedCertification converts from the SavePPMShipmentSignedCertification payload and the
// SignedCertification model, taking into account an existing ID
func ReSavePPMShipmentSignedCertification(ppmShipmentID uuid.UUID, signedCertificationID uuid.UUID, signedCertification internalmessages.SavePPMShipmentSignedCertification) models.SignedCertification {
	model := SavePPMShipmentSignedCertification(ppmShipmentID, signedCertification)

	model.ID = signedCertificationID

	return model
}

// SignedCertificationFromSubmit
func SignedCertificationFromSubmit(payload *internalmessages.SubmitMoveForApprovalPayload, userID uuid.UUID, moveID strfmt.UUID) *models.SignedCertification {
	if payload == nil {
		return nil
	}
	date := time.Time(*payload.Certificate.Date)
	certType := models.SignedCertificationType(*payload.Certificate.CertificationType)
	newSignedCertification := models.SignedCertification{
		MoveID:            uuid.FromStringOrNil(moveID.String()),
		CertificationType: &certType,
		SubmittingUserID:  userID,
		CertificationText: *payload.Certificate.CertificationText,
		Signature:         *payload.Certificate.Signature,
		Date:              date,
	}

	return &newSignedCertification
}
