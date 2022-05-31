package payloads

import (
	"time"

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
	model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)
	model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)

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
		SITExpected: ppmShipment.SitExpected,
	}

	expectedDepartureDate := time.Time(*ppmShipment.ExpectedDepartureDate)
	if !expectedDepartureDate.IsZero() {
		model.ExpectedDepartureDate = expectedDepartureDate
	}

	if ppmShipment.PickupPostalCode != nil {
		model.PickupPostalCode = *ppmShipment.PickupPostalCode
	}

	if ppmShipment.DestinationPostalCode != nil {
		model.DestinationPostalCode = *ppmShipment.DestinationPostalCode
	}

	return model
}

func UpdatePPMShipmentModel(ppmShipment *internalmessages.UpdatePPMShipment) *models.PPMShipment {
	if ppmShipment == nil {
		return nil
	}

	ppmModel := &models.PPMShipment{
		ActualMoveDate:                 (*time.Time)(ppmShipment.ActualMoveDate),
		SecondaryPickupPostalCode:      handlers.FmtNullableStringToStringPtr(ppmShipment.SecondaryPickupPostalCode),
		ActualPickupPostalCode:         ppmShipment.ActualPickupPostalCode,
		SecondaryDestinationPostalCode: handlers.FmtNullableStringToStringPtr(ppmShipment.SecondaryDestinationPostalCode),
		ActualDestinationPostalCode:    ppmShipment.ActualDestinationPostalCode,
		SITExpected:                    ppmShipment.SitExpected,
		EstimatedWeight:                handlers.PoundPtrFromInt64Ptr(ppmShipment.EstimatedWeight),
		NetWeight:                      handlers.PoundPtrFromInt64Ptr(ppmShipment.NetWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.PoundPtrFromInt64Ptr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.PoundPtrFromInt64Ptr(ppmShipment.SpouseProGearWeight),
		AdvanceRequested:               ppmShipment.AdvanceRequested,
		HasRequestedAdvance:            ppmShipment.AdvanceRequested,
		Advance:                        handlers.FmtInt64PtrToPopPtr(ppmShipment.Advance),
		AdvanceAmountRequested:         handlers.FmtInt64PtrToPopPtr(ppmShipment.Advance),
		HasReceivedAdvance:             ppmShipment.HasReceivedAdvance,
		AdvanceAmountReceived:          handlers.FmtInt64PtrToPopPtr(ppmShipment.AdvanceAmountReceived),
	}

	if ppmShipment.ExpectedDepartureDate != nil {
		ppmModel.ExpectedDepartureDate = *handlers.FmtDatePtrToPopPtr(ppmShipment.ExpectedDepartureDate)
	}

	if ppmShipment.DestinationPostalCode != nil {
		ppmModel.DestinationPostalCode = *ppmShipment.DestinationPostalCode
	}

	if ppmShipment.PickupPostalCode != nil {
		ppmModel.PickupPostalCode = *ppmShipment.PickupPostalCode
	}

	return ppmModel
}

// MTOShipmentModelFromUpdate model
func MTOShipmentModelFromUpdate(mtoShipment *internalmessages.UpdateShipment) *models.MTOShipment {
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
	model.SecondaryPickupAddress = AddressModel(mtoShipment.SecondaryPickupAddress)
	model.DestinationAddress = AddressModel(mtoShipment.DestinationAddress)
	model.SecondaryDeliveryAddress = AddressModel(mtoShipment.SecondaryDeliveryAddress)

	if mtoShipment.Agents != nil {
		model.MTOAgents = *MTOAgentsModel(&mtoShipment.Agents)
	}

	model.PPMShipment = UpdatePPMShipmentModel(mtoShipment.PpmShipment)

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
