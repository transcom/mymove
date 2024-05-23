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
	return &models.Address{
		ID:             uuid.FromStringOrNil(address.ID.String()),
		StreetAddress1: *address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           *address.City,
		State:          *address.State,
		PostalCode:     *address.PostalCode,
		Country:        address.Country,
		County:         *address.County,
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

	// A PPM type shipment begins in DRAFT because it requires a multi-page series to complete.
	// After move submission a PPM's status will change to SUBMITTED
	if model.ShipmentType == models.MTOShipmentTypePPM {
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
	model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)
	model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)
	model.HasSecondaryDeliveryAddress = handlers.FmtBool(mtoShipment.SecondaryDeliveryAddress != nil)

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	if mtoShipment.PpmShipment != nil {
		model.PPMShipment = PPMShipmentModelFromCreate(mtoShipment.PpmShipment)
		model.PPMShipment.Shipment = *model
	}

	return model
}

// PPMShipmentModelFromCreate model
func PPMShipmentModelFromCreate(ppmShipment *internalmessages.CreatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	model := &models.PPMShipment{
		PickupPostalCode:               *ppmShipment.PickupPostalCode,
		SecondaryPickupPostalCode:      handlers.FmtNullableStringToStringPtrNilToNil(ppmShipment.SecondaryPickupPostalCode),
		DestinationPostalCode:          *ppmShipment.DestinationPostalCode,
		SecondaryDestinationPostalCode: handlers.FmtNullableStringToStringPtrNilToNil(ppmShipment.SecondaryDestinationPostalCode),
		SITExpected:                    ppmShipment.SitExpected,
		ExpectedDepartureDate:          handlers.FmtDatePtrToPop(ppmShipment.ExpectedDepartureDate),
	}

	if ppmShipment.PickupAddress != nil {
		model.PickupAddress = AddressModel(ppmShipment.PickupAddress)
	}

	model.HasSecondaryPickupAddress = handlers.FmtBool(ppmShipment.SecondaryPickupAddress != nil)
	if ppmShipment.SecondaryPickupAddress != nil {
		model.SecondaryPickupAddress = AddressModel(ppmShipment.SecondaryPickupAddress)
	}

	if ppmShipment.DestinationAddress != nil {
		model.DestinationAddress = AddressModel(ppmShipment.DestinationAddress)
	}

	model.HasSecondaryDestinationAddress = handlers.FmtBool(ppmShipment.SecondaryDestinationAddress != nil)
	if ppmShipment.SecondaryDestinationAddress != nil {
		model.SecondaryDestinationAddress = AddressModel(ppmShipment.SecondaryDestinationAddress)
	}

	return model
}

func UpdatePPMShipmentModel(ppmShipment *internalmessages.UpdatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	ppmModel := &models.PPMShipment{
		ActualMoveDate:                 (*time.Time)(ppmShipment.ActualMoveDate),
		SecondaryPickupPostalCode:      handlers.FmtNullableStringToStringPtrNilToBlankString(ppmShipment.SecondaryPickupPostalCode),
		ActualPickupPostalCode:         ppmShipment.ActualPickupPostalCode,
		SecondaryDestinationPostalCode: handlers.FmtNullableStringToStringPtrNilToBlankString(ppmShipment.SecondaryDestinationPostalCode),
		ActualDestinationPostalCode:    ppmShipment.ActualDestinationPostalCode,
		SITExpected:                    ppmShipment.SitExpected,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtInt64PtrToPopPtr(ppmShipment.AdvanceAmountRequested),
		HasReceivedAdvance:             ppmShipment.HasReceivedAdvance,
		AdvanceAmountReceived:          handlers.FmtInt64PtrToPopPtr(ppmShipment.AdvanceAmountReceived),
		FinalIncentive:                 handlers.FmtInt64PtrToPopPtr(ppmShipment.FinalIncentive),
		HasSecondaryPickupAddress:      ppmShipment.HasSecondaryPickupAddress,
		HasSecondaryDestinationAddress: ppmShipment.HasSecondaryDestinationAddress,
	}

	ppmModel.W2Address = AddressModel(ppmShipment.W2Address)
	if ppmShipment.ExpectedDepartureDate != nil {
		ppmModel.ExpectedDepartureDate = *handlers.FmtDatePtrToPopPtr(ppmShipment.ExpectedDepartureDate)
	}

	if ppmShipment.PickupPostalCode != nil {
		ppmModel.PickupPostalCode = *ppmShipment.PickupPostalCode
	}

	if ppmShipment.PickupAddress != nil {
		ppmModel.PickupAddress = AddressModel(ppmShipment.PickupAddress)
	}

	if ppmShipment.SecondaryPickupAddress != nil {
		ppmModel.SecondaryPickupAddress = AddressModel(ppmShipment.SecondaryPickupAddress)
	}

	if ppmShipment.DestinationPostalCode != nil {
		ppmModel.DestinationPostalCode = *ppmShipment.DestinationPostalCode
	}

	if ppmShipment.DestinationAddress != nil {
		ppmModel.DestinationAddress = AddressModel(ppmShipment.DestinationAddress)
	}

	if ppmShipment.SecondaryDestinationAddress != nil {
		ppmModel.SecondaryDestinationAddress = AddressModel(ppmShipment.SecondaryDestinationAddress)
	}

	if ppmShipment.FinalIncentive != nil {
		ppmModel.FinalIncentive = handlers.FmtInt64PtrToPopPtr(ppmShipment.FinalIncentive)
	}

	return ppmModel
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
		ActualProGearWeight:         handlers.PoundPtrFromInt64Ptr(mtoShipment.ActualProGearWeight),
		ActualSpouseProGearWeight:   handlers.PoundPtrFromInt64Ptr(mtoShipment.ActualSpouseProGearWeight),
	}

	model.PickupAddress = AddressModel(mtoShipment.PickupAddress)
	if mtoShipment.HasSecondaryPickupAddress != nil {
		if *mtoShipment.HasSecondaryPickupAddress {
			model.SecondaryPickupAddress = AddressModel(mtoShipment.SecondaryPickupAddress)
		}
	}
	if mtoShipment.HasSecondaryDeliveryAddress != nil {
		if *mtoShipment.HasSecondaryDeliveryAddress {
			model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)
		}
	}
	model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)
	model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	model.PPMShipment = UpdatePPMShipmentModel(mtoShipment.PpmShipment)

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
