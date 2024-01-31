package payloads

import (
	"errors"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/ghcv2messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/storage"
)

// Address payload
func Address(address *models.Address) *ghcv2messages.Address {
	if address == nil {
		return nil
	}
	return &ghcv2messages.Address{
		ID:             strfmt.UUID(address.ID.String()),
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           &address.City,
		State:          &address.State,
		PostalCode:     &address.PostalCode,
		Country:        address.Country,
		ETag:           etag.GenerateEtag(address.UpdatedAt),
	}
}

// StorageFacility payload
func StorageFacility(storageFacility *models.StorageFacility) *ghcv2messages.StorageFacility {
	if storageFacility == nil {
		return nil
	}

	payload := ghcv2messages.StorageFacility{
		ID:           strfmt.UUID(storageFacility.ID.String()),
		FacilityName: storageFacility.FacilityName,
		Address:      Address(&storageFacility.Address),
		LotNumber:    storageFacility.LotNumber,
		Phone:        storageFacility.Phone,
		Email:        storageFacility.Email,
		ETag:         etag.GenerateEtag(storageFacility.UpdatedAt),
	}

	return &payload
}

// SITDurationUpdate payload
func SITDurationUpdate(sitDurationUpdate *models.SITDurationUpdate) *ghcv2messages.SITExtension {
	if sitDurationUpdate == nil {
		return nil
	}
	payload := &ghcv2messages.SITExtension{
		ID:                strfmt.UUID(sitDurationUpdate.ID.String()),
		ETag:              etag.GenerateEtag(sitDurationUpdate.UpdatedAt),
		MtoShipmentID:     strfmt.UUID(sitDurationUpdate.MTOShipmentID.String()),
		RequestReason:     string(sitDurationUpdate.RequestReason),
		RequestedDays:     int64(sitDurationUpdate.RequestedDays),
		Status:            string(sitDurationUpdate.Status),
		CreatedAt:         strfmt.DateTime(sitDurationUpdate.CreatedAt),
		UpdatedAt:         strfmt.DateTime(sitDurationUpdate.UpdatedAt),
		ApprovedDays:      handlers.FmtIntPtrToInt64(sitDurationUpdate.ApprovedDays),
		ContractorRemarks: handlers.FmtStringPtr(sitDurationUpdate.ContractorRemarks),
		DecisionDate:      handlers.FmtDateTimePtr(sitDurationUpdate.DecisionDate),
		OfficeRemarks:     handlers.FmtStringPtr(sitDurationUpdate.OfficeRemarks),
	}

	return payload
}

// SITDurationUpdates payload
func SITDurationUpdates(sitDurationUpdates *models.SITDurationUpdates) *ghcv2messages.SITExtensions {
	payload := make(ghcv2messages.SITExtensions, len(*sitDurationUpdates))

	if len(*sitDurationUpdates) > 0 {
		for i, m := range *sitDurationUpdates {
			copyOfSITDurationUpdate := m // Make copy to avoid implicit memory aliasing of items from a range statement.
			payload[i] = SITDurationUpdate(&copyOfSITDurationUpdate)
		}
		// Reversing the SIT duration updates as they are saved in the order
		// they are created and we want to always display them in the reverse
		// order.
		for i, j := 0, len(payload)-1; i < j; i, j = i+1, j-1 {
			payload[i], payload[j] = payload[j], payload[i]
		}
	}
	return &payload
}

func currentSIT(currentSIT *services.CurrentSIT) *ghcv2messages.SITStatusCurrentSIT {
	if currentSIT == nil {
		return nil
	}
	return &ghcv2messages.SITStatusCurrentSIT{
		ServiceItemID:        *handlers.FmtUUID(currentSIT.ServiceItemID),
		Location:             currentSIT.Location,
		DaysInSIT:            handlers.FmtIntPtrToInt64(&currentSIT.DaysInSIT),
		SitEntryDate:         handlers.FmtDate(currentSIT.SITEntryDate),
		SitDepartureDate:     handlers.FmtDatePtr(currentSIT.SITDepartureDate),
		SitAllowanceEndDate:  handlers.FmtDate(currentSIT.SITAllowanceEndDate),
		SitCustomerContacted: handlers.FmtDatePtr(currentSIT.SITCustomerContacted),
		SitRequestedDelivery: handlers.FmtDatePtr(currentSIT.SITRequestedDelivery),
	}
}

// SITStatus payload
func SITStatus(shipmentSITStatuses *services.SITStatus, storer storage.FileStorer) *ghcv2messages.SITStatus {
	if shipmentSITStatuses == nil {
		return nil
	}
	payload := &ghcv2messages.SITStatus{
		PastSITServiceItems:      MTOServiceItemModels(shipmentSITStatuses.PastSITs, storer),
		TotalSITDaysUsed:         handlers.FmtIntPtrToInt64(&shipmentSITStatuses.TotalSITDaysUsed),
		TotalDaysRemaining:       handlers.FmtIntPtrToInt64(&shipmentSITStatuses.TotalDaysRemaining),
		CalculatedTotalDaysInSIT: handlers.FmtIntPtrToInt64(&shipmentSITStatuses.CalculatedTotalDaysInSIT),
		CurrentSIT:               currentSIT(shipmentSITStatuses.CurrentSIT),
	}

	return payload
}

// SITStatuses payload
func SITStatuses(shipmentSITStatuses map[string]services.SITStatus, storer storage.FileStorer) map[string]*ghcv2messages.SITStatus {
	sitStatuses := map[string]*ghcv2messages.SITStatus{}
	if len(shipmentSITStatuses) == 0 {
		return sitStatuses
	}

	for _, sitStatus := range shipmentSITStatuses {
		copyOfSITStatus := sitStatus
		sitStatuses[sitStatus.ShipmentID.String()] = SITStatus(&copyOfSITStatus, storer)
	}

	return sitStatuses
}

// PPMShipment payload
func PPMShipment(_ storage.FileStorer, ppmShipment *models.PPMShipment) *ghcv2messages.PPMShipment {
	if ppmShipment == nil || ppmShipment.ID.IsNil() {
		return nil
	}

	payloadPPMShipment := &ghcv2messages.PPMShipment{
		ID:                             *handlers.FmtUUID(ppmShipment.ID),
		ShipmentID:                     *handlers.FmtUUID(ppmShipment.ShipmentID),
		CreatedAt:                      strfmt.DateTime(ppmShipment.CreatedAt),
		UpdatedAt:                      strfmt.DateTime(ppmShipment.UpdatedAt),
		Status:                         ghcv2messages.PPMShipmentStatus(ppmShipment.Status),
		ExpectedDepartureDate:          handlers.FmtDate(ppmShipment.ExpectedDepartureDate),
		ActualMoveDate:                 handlers.FmtDatePtr(ppmShipment.ActualMoveDate),
		SubmittedAt:                    handlers.FmtDateTimePtr(ppmShipment.SubmittedAt),
		ReviewedAt:                     handlers.FmtDateTimePtr(ppmShipment.ReviewedAt),
		ApprovedAt:                     handlers.FmtDateTimePtr(ppmShipment.ApprovedAt),
		PickupPostalCode:               &ppmShipment.PickupPostalCode,
		SecondaryPickupPostalCode:      ppmShipment.SecondaryPickupPostalCode,
		ActualPickupPostalCode:         ppmShipment.ActualPickupPostalCode,
		DestinationPostalCode:          &ppmShipment.DestinationPostalCode,
		SecondaryDestinationPostalCode: ppmShipment.SecondaryDestinationPostalCode,
		ActualDestinationPostalCode:    ppmShipment.ActualDestinationPostalCode,
		SitExpected:                    ppmShipment.SITExpected,
		EstimatedWeight:                handlers.FmtPoundPtr(ppmShipment.EstimatedWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.FmtPoundPtr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.FmtPoundPtr(ppmShipment.SpouseProGearWeight),
		EstimatedIncentive:             handlers.FmtCost(ppmShipment.EstimatedIncentive),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtCost(ppmShipment.AdvanceAmountRequested),
		HasReceivedAdvance:             ppmShipment.HasReceivedAdvance,
		AdvanceAmountReceived:          handlers.FmtCost(ppmShipment.AdvanceAmountReceived),
		SitEstimatedWeight:             handlers.FmtPoundPtr(ppmShipment.SITEstimatedWeight),
		SitEstimatedEntryDate:          handlers.FmtDatePtr(ppmShipment.SITEstimatedEntryDate),
		SitEstimatedDepartureDate:      handlers.FmtDatePtr(ppmShipment.SITEstimatedDepartureDate),
		SitEstimatedCost:               handlers.FmtCost(ppmShipment.SITEstimatedCost),
		ETag:                           etag.GenerateEtag(ppmShipment.UpdatedAt),
	}

	if ppmShipment.SITLocation != nil {
		sitLocation := ghcv2messages.SITLocationType(*ppmShipment.SITLocation)
		payloadPPMShipment.SitLocation = &sitLocation
	}

	if ppmShipment.AdvanceStatus != nil {
		advanceStatus := ghcv2messages.PPMAdvanceStatus(*ppmShipment.AdvanceStatus)
		payloadPPMShipment.AdvanceStatus = &advanceStatus
	}

	if ppmShipment.W2Address != nil {
		payloadPPMShipment.W2Address = Address(ppmShipment.W2Address)
	}

	return payloadPPMShipment
}

// ProGearWeightTickets sets up a ProGearWeightTicket slice for the api using model data.
func ProGearWeightTickets(storer storage.FileStorer, proGearWeightTickets models.ProgearWeightTickets) []*ghcv2messages.ProGearWeightTicket {
	payload := make([]*ghcv2messages.ProGearWeightTicket, len(proGearWeightTickets))
	for i, proGearWeightTicket := range proGearWeightTickets {
		copyOfProGearWeightTicket := proGearWeightTicket
		proGearWeightTicketPayload := ProGearWeightTicket(storer, &copyOfProGearWeightTicket)
		payload[i] = proGearWeightTicketPayload
	}
	return payload
}

// ProGearWeightTicket payload
func ProGearWeightTicket(storer storage.FileStorer, progear *models.ProgearWeightTicket) *ghcv2messages.ProGearWeightTicket {
	ppmShipmentID := strfmt.UUID(progear.PPMShipmentID.String())

	document, err := PayloadForDocumentModel(storer, progear.Document)
	if err != nil {
		return nil
	}

	payload := &ghcv2messages.ProGearWeightTicket{
		ID:               strfmt.UUID(progear.ID.String()),
		PpmShipmentID:    ppmShipmentID,
		CreatedAt:        *handlers.FmtDateTime(progear.CreatedAt),
		UpdatedAt:        *handlers.FmtDateTime(progear.UpdatedAt),
		DocumentID:       *handlers.FmtUUID(progear.DocumentID),
		Document:         document,
		Weight:           handlers.FmtPoundPtr(progear.Weight),
		BelongsToSelf:    progear.BelongsToSelf,
		HasWeightTickets: progear.HasWeightTickets,
		Description:      progear.Description,
		ETag:             etag.GenerateEtag(progear.UpdatedAt),
	}

	if progear.Status != nil {
		status := ghcv2messages.OmittablePPMDocumentStatus(*progear.Status)
		payload.Status = &status
	}

	if progear.Reason != nil {
		reason := ghcv2messages.PPMDocumentStatusReason(*progear.Reason)
		payload.Reason = &reason
	}

	return payload
}

// MovingExpense payload
func MovingExpense(storer storage.FileStorer, movingExpense *models.MovingExpense) *ghcv2messages.MovingExpense {

	document, err := PayloadForDocumentModel(storer, movingExpense.Document)
	if err != nil {
		return nil
	}

	payload := &ghcv2messages.MovingExpense{
		ID:             *handlers.FmtUUID(movingExpense.ID),
		PpmShipmentID:  *handlers.FmtUUID(movingExpense.PPMShipmentID),
		DocumentID:     *handlers.FmtUUID(movingExpense.DocumentID),
		Document:       document,
		CreatedAt:      strfmt.DateTime(movingExpense.CreatedAt),
		UpdatedAt:      strfmt.DateTime(movingExpense.UpdatedAt),
		Description:    movingExpense.Description,
		PaidWithGtcc:   movingExpense.PaidWithGTCC,
		Amount:         handlers.FmtCost(movingExpense.Amount),
		MissingReceipt: movingExpense.MissingReceipt,
		ETag:           etag.GenerateEtag(movingExpense.UpdatedAt),
	}
	if movingExpense.MovingExpenseType != nil {
		movingExpenseType := ghcv2messages.OmittableMovingExpenseType(*movingExpense.MovingExpenseType)
		payload.MovingExpenseType = &movingExpenseType
	}

	if movingExpense.Status != nil {
		status := ghcv2messages.OmittablePPMDocumentStatus(*movingExpense.Status)
		payload.Status = &status
	}

	if movingExpense.Reason != nil {
		reason := ghcv2messages.PPMDocumentStatusReason(*movingExpense.Reason)
		payload.Reason = &reason
	}

	if movingExpense.SITStartDate != nil {
		payload.SitStartDate = handlers.FmtDatePtr(movingExpense.SITStartDate)
	}

	if movingExpense.SITEndDate != nil {
		payload.SitEndDate = handlers.FmtDatePtr(movingExpense.SITEndDate)
	}

	return payload
}

func MovingExpenses(storer storage.FileStorer, movingExpenses models.MovingExpenses) []*ghcv2messages.MovingExpense {
	payload := make([]*ghcv2messages.MovingExpense, len(movingExpenses))
	for i, movingExpense := range movingExpenses {
		copyOfMovingExpense := movingExpense
		payload[i] = MovingExpense(storer, &copyOfMovingExpense)
	}
	return payload
}

func WeightTickets(storer storage.FileStorer, weightTickets models.WeightTickets) []*ghcv2messages.WeightTicket {
	payload := make([]*ghcv2messages.WeightTicket, len(weightTickets))
	for i, weightTicket := range weightTickets {
		copyOfWeightTicket := weightTicket
		weightTicketPayload := WeightTicket(storer, &copyOfWeightTicket)
		payload[i] = weightTicketPayload
	}
	return payload
}

// WeightTicket payload
func WeightTicket(storer storage.FileStorer, weightTicket *models.WeightTicket) *ghcv2messages.WeightTicket {
	ppmShipment := strfmt.UUID(weightTicket.PPMShipmentID.String())

	emptyDocument, err := PayloadForDocumentModel(storer, weightTicket.EmptyDocument)
	if err != nil {
		return nil
	}

	fullDocument, err := PayloadForDocumentModel(storer, weightTicket.FullDocument)
	if err != nil {
		return nil
	}

	proofOfTrailerOwnershipDocument, err := PayloadForDocumentModel(storer, weightTicket.ProofOfTrailerOwnershipDocument)
	if err != nil {
		return nil
	}

	payload := &ghcv2messages.WeightTicket{
		ID:                                strfmt.UUID(weightTicket.ID.String()),
		PpmShipmentID:                     ppmShipment,
		CreatedAt:                         *handlers.FmtDateTime(weightTicket.CreatedAt),
		UpdatedAt:                         *handlers.FmtDateTime(weightTicket.UpdatedAt),
		VehicleDescription:                weightTicket.VehicleDescription,
		EmptyWeight:                       handlers.FmtPoundPtr(weightTicket.EmptyWeight),
		MissingEmptyWeightTicket:          weightTicket.MissingEmptyWeightTicket,
		EmptyDocumentID:                   *handlers.FmtUUID(weightTicket.EmptyDocumentID),
		EmptyDocument:                     emptyDocument,
		FullWeight:                        handlers.FmtPoundPtr(weightTicket.FullWeight),
		MissingFullWeightTicket:           weightTicket.MissingFullWeightTicket,
		FullDocumentID:                    *handlers.FmtUUID(weightTicket.FullDocumentID),
		FullDocument:                      fullDocument,
		OwnsTrailer:                       weightTicket.OwnsTrailer,
		TrailerMeetsCriteria:              weightTicket.TrailerMeetsCriteria,
		ProofOfTrailerOwnershipDocumentID: *handlers.FmtUUID(weightTicket.ProofOfTrailerOwnershipDocumentID),
		ProofOfTrailerOwnershipDocument:   proofOfTrailerOwnershipDocument,
		AdjustedNetWeight:                 handlers.FmtPoundPtr(weightTicket.AdjustedNetWeight),
		NetWeightRemarks:                  weightTicket.NetWeightRemarks,
		ETag:                              etag.GenerateEtag(weightTicket.UpdatedAt),
	}

	if weightTicket.Status != nil {
		status := ghcv2messages.OmittablePPMDocumentStatus(*weightTicket.Status)
		payload.Status = &status
	}

	if weightTicket.Reason != nil {
		reason := ghcv2messages.PPMDocumentStatusReason(*weightTicket.Reason)
		payload.Reason = &reason
	}

	return payload
}

// ShipmentAddressUpdate payload
func ShipmentAddressUpdate(shipmentAddressUpdate *models.ShipmentAddressUpdate) *ghcv2messages.ShipmentAddressUpdate {
	if shipmentAddressUpdate == nil || shipmentAddressUpdate.ID.IsNil() {
		return nil
	}

	payload := &ghcv2messages.ShipmentAddressUpdate{
		ID:                strfmt.UUID(shipmentAddressUpdate.ID.String()),
		ShipmentID:        strfmt.UUID(shipmentAddressUpdate.ShipmentID.String()),
		NewAddress:        Address(&shipmentAddressUpdate.NewAddress),
		OriginalAddress:   Address(&shipmentAddressUpdate.OriginalAddress),
		ContractorRemarks: shipmentAddressUpdate.ContractorRemarks,
		OfficeRemarks:     shipmentAddressUpdate.OfficeRemarks,
		Status:            ghcv2messages.ShipmentAddressUpdateStatus(shipmentAddressUpdate.Status),
	}

	return payload
}

// MTOShipment payload
func MTOShipment(storer storage.FileStorer, mtoShipment *models.MTOShipment, sitStatusPayload *ghcv2messages.SITStatus) *ghcv2messages.MTOShipment {

	payload := &ghcv2messages.MTOShipment{
		ID:                          strfmt.UUID(mtoShipment.ID.String()),
		MoveTaskOrderID:             strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:                ghcv2messages.MTOShipmentType(mtoShipment.ShipmentType),
		Status:                      ghcv2messages.MTOShipmentStatus(mtoShipment.Status),
		CounselorRemarks:            mtoShipment.CounselorRemarks,
		CustomerRemarks:             mtoShipment.CustomerRemarks,
		RejectionReason:             mtoShipment.RejectionReason,
		PickupAddress:               Address(mtoShipment.PickupAddress),
		SecondaryDeliveryAddress:    Address(mtoShipment.SecondaryDeliveryAddress),
		SecondaryPickupAddress:      Address(mtoShipment.SecondaryPickupAddress),
		DestinationAddress:          Address(mtoShipment.DestinationAddress),
		HasSecondaryDeliveryAddress: mtoShipment.HasSecondaryDeliveryAddress,
		HasSecondaryPickupAddress:   mtoShipment.HasSecondaryPickupAddress,
		PrimeEstimatedWeight:        handlers.FmtPoundPtr(mtoShipment.PrimeEstimatedWeight),
		PrimeActualWeight:           handlers.FmtPoundPtr(mtoShipment.PrimeActualWeight),
		NtsRecordedWeight:           handlers.FmtPoundPtr(mtoShipment.NTSRecordedWeight),
		MtoAgents:                   *MTOAgents(&mtoShipment.MTOAgents),
		MtoServiceItems:             MTOServiceItemModels(mtoShipment.MTOServiceItems, storer),
		Diversion:                   mtoShipment.Diversion,
		Reweigh:                     Reweigh(mtoShipment.Reweigh, sitStatusPayload),
		CreatedAt:                   strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                   strfmt.DateTime(mtoShipment.UpdatedAt),
		ETag:                        etag.GenerateEtag(mtoShipment.UpdatedAt),
		DeletedAt:                   handlers.FmtDateTimePtr(mtoShipment.DeletedAt),
		ApprovedDate:                handlers.FmtDateTimePtr(mtoShipment.ApprovedDate),
		SitDaysAllowance:            handlers.FmtIntPtrToInt64(mtoShipment.SITDaysAllowance),
		SitExtensions:               *SITDurationUpdates(&mtoShipment.SITDurationUpdates),
		BillableWeightCap:           handlers.FmtPoundPtr(mtoShipment.BillableWeightCap),
		BillableWeightJustification: mtoShipment.BillableWeightJustification,
		UsesExternalVendor:          mtoShipment.UsesExternalVendor,
		ServiceOrderNumber:          mtoShipment.ServiceOrderNumber,
		StorageFacility:             StorageFacility(mtoShipment.StorageFacility),
		PpmShipment:                 PPMShipment(storer, mtoShipment.PPMShipment),
		DeliveryAddressUpdate:       ShipmentAddressUpdate(mtoShipment.DeliveryAddressUpdate),
	}

	if sitStatusPayload != nil {
		// If we have a sitStatusPayload, overwrite SitDaysAllowance from the shipment model.
		totalSITAllowance := 0
		if sitStatusPayload.TotalDaysRemaining != nil {
			totalSITAllowance += int(*sitStatusPayload.TotalDaysRemaining)
		}
		if sitStatusPayload.TotalSITDaysUsed != nil {
			totalSITAllowance += int(*sitStatusPayload.TotalSITDaysUsed)
		}
		payload.SitDaysAllowance = handlers.FmtIntPtrToInt64(&totalSITAllowance)
	}

	if mtoShipment.SITDurationUpdates != nil && len(mtoShipment.SITDurationUpdates) > 0 {
		payload.SitExtensions = *SITDurationUpdates(&mtoShipment.SITDurationUpdates)
	}

	if mtoShipment.RequestedPickupDate != nil && !mtoShipment.RequestedPickupDate.IsZero() {
		payload.RequestedPickupDate = handlers.FmtDatePtr(mtoShipment.RequestedPickupDate)
	}

	if mtoShipment.ActualPickupDate != nil && !mtoShipment.ActualPickupDate.IsZero() {
		payload.ActualPickupDate = handlers.FmtDatePtr(mtoShipment.ActualPickupDate)
	}

	if mtoShipment.ActualDeliveryDate != nil && !mtoShipment.ActualDeliveryDate.IsZero() {
		payload.ActualDeliveryDate = handlers.FmtDatePtr(mtoShipment.ActualDeliveryDate)
	}

	if mtoShipment.RequestedDeliveryDate != nil && !mtoShipment.RequestedDeliveryDate.IsZero() {
		payload.RequestedDeliveryDate = handlers.FmtDatePtr(mtoShipment.RequestedDeliveryDate)
	}

	if mtoShipment.RequiredDeliveryDate != nil && !mtoShipment.RequiredDeliveryDate.IsZero() {
		payload.RequiredDeliveryDate = handlers.FmtDatePtr(mtoShipment.RequiredDeliveryDate)
	}

	if mtoShipment.ScheduledPickupDate != nil {
		payload.ScheduledPickupDate = handlers.FmtDatePtr(mtoShipment.ScheduledPickupDate)
	}

	if mtoShipment.ScheduledDeliveryDate != nil {
		payload.ScheduledDeliveryDate = handlers.FmtDatePtr(mtoShipment.ScheduledDeliveryDate)
	}

	if mtoShipment.DestinationType != nil {
		destinationType := ghcv2messages.DestinationType(*mtoShipment.DestinationType)
		payload.DestinationType = &destinationType
	}

	if sitStatusPayload != nil {
		payload.SitStatus = sitStatusPayload
	}

	if mtoShipment.TACType != nil {
		tt := ghcv2messages.LOAType(*mtoShipment.TACType)
		payload.TacType = &tt
	}

	if mtoShipment.SACType != nil {
		st := ghcv2messages.LOAType(*mtoShipment.SACType)
		payload.SacType = &st
	}

	weightsCalculator := mtoshipment.NewShipmentBillableWeightCalculator()
	calculatedWeights := weightsCalculator.CalculateShipmentBillableWeight(mtoShipment)

	// CalculatedBillableWeight is intentionally not a part of the mto_shipments model
	// because we don't want to store a derived value in the database
	payload.CalculatedBillableWeight = handlers.FmtPoundPtr(calculatedWeights.CalculatedBillableWeight)

	return payload
}

// MTOAgent payload
func MTOAgent(mtoAgent *models.MTOAgent) *ghcv2messages.MTOAgent {
	payload := &ghcv2messages.MTOAgent{
		ID:            strfmt.UUID(mtoAgent.ID.String()),
		MtoShipmentID: strfmt.UUID(mtoAgent.MTOShipmentID.String()),
		CreatedAt:     strfmt.DateTime(mtoAgent.CreatedAt),
		UpdatedAt:     strfmt.DateTime(mtoAgent.UpdatedAt),
		FirstName:     mtoAgent.FirstName,
		LastName:      mtoAgent.LastName,
		AgentType:     string(mtoAgent.MTOAgentType),
		Email:         mtoAgent.Email,
		Phone:         mtoAgent.Phone,
		ETag:          etag.GenerateEtag(mtoAgent.UpdatedAt),
	}
	return payload
}

// MTOAgents payload
func MTOAgents(mtoAgents *models.MTOAgents) *ghcv2messages.MTOAgents {
	payload := make(ghcv2messages.MTOAgents, len(*mtoAgents))
	for i, m := range *mtoAgents {
		copyOfMtoAgent := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOAgent(&copyOfMtoAgent)
	}
	return &payload
}

func ServiceRequestDoc(serviceRequest models.ServiceRequestDocument, storer storage.FileStorer) (*ghcv2messages.ServiceRequestDocument, error) {

	uploads := make([]*ghcv2messages.Upload, len(serviceRequest.ServiceRequestDocumentUploads))

	if serviceRequest.ServiceRequestDocumentUploads != nil && len(serviceRequest.ServiceRequestDocumentUploads) > 0 {
		for i, serviceRequestUpload := range serviceRequest.ServiceRequestDocumentUploads {
			url, err := storer.PresignedURL(serviceRequestUpload.Upload.StorageKey, serviceRequestUpload.Upload.ContentType)
			if err != nil {
				return nil, err
			}
			uploads[i] = Upload(storer, serviceRequestUpload.Upload, url)
		}
	}

	return &ghcv2messages.ServiceRequestDocument{
		Uploads: uploads,
	}, nil

}

// MTOServiceItemModel payload
func MTOServiceItemModel(s *models.MTOServiceItem, storer storage.FileStorer) *ghcv2messages.MTOServiceItem {
	if s == nil {
		return nil
	}

	serviceRequestDocs := make(ghcv2messages.ServiceRequestDocuments, len(s.ServiceRequestDocuments))

	if s.ServiceRequestDocuments != nil && len(s.ServiceRequestDocuments) > 0 {
		for i, serviceRequest := range s.ServiceRequestDocuments {
			payload, err := ServiceRequestDoc(serviceRequest, storer)
			if err != nil {
				return nil
			}
			serviceRequestDocs[i] = payload
		}
	}

	return &ghcv2messages.MTOServiceItem{
		ID:                            handlers.FmtUUID(s.ID),
		MoveTaskOrderID:               handlers.FmtUUID(s.MoveTaskOrderID),
		MtoShipmentID:                 handlers.FmtUUIDPtr(s.MTOShipmentID),
		ReServiceID:                   handlers.FmtUUID(s.ReServiceID),
		ReServiceCode:                 handlers.FmtString(string(s.ReService.Code)),
		ReServiceName:                 handlers.FmtStringPtr(&s.ReService.Name),
		Reason:                        handlers.FmtStringPtr(s.Reason),
		RejectionReason:               handlers.FmtStringPtr(s.RejectionReason),
		PickupPostalCode:              handlers.FmtStringPtr(s.PickupPostalCode),
		SITPostalCode:                 handlers.FmtStringPtr(s.SITPostalCode),
		SitEntryDate:                  handlers.FmtDateTimePtr(s.SITEntryDate),
		SitDepartureDate:              handlers.FmtDateTimePtr(s.SITDepartureDate),
		SitCustomerContacted:          handlers.FmtDatePtr(s.SITCustomerContacted),
		SitRequestedDelivery:          handlers.FmtDatePtr(s.SITRequestedDelivery),
		Status:                        ghcv2messages.MTOServiceItemStatus(s.Status),
		Description:                   handlers.FmtStringPtr(s.Description),
		Dimensions:                    MTOServiceItemDimensions(s.Dimensions),
		CustomerContacts:              MTOServiceItemCustomerContacts(s.CustomerContacts),
		SitAddressUpdates:             SITAddressUpdates(s.SITAddressUpdates),
		SitDestinationOriginalAddress: Address(s.SITDestinationOriginalAddress),
		SitDestinationFinalAddress:    Address(s.SITDestinationFinalAddress),
		EstimatedWeight:               handlers.FmtPoundPtr(s.EstimatedWeight),
		CreatedAt:                     strfmt.DateTime(s.CreatedAt),
		ApprovedAt:                    handlers.FmtDateTimePtr(s.ApprovedAt),
		RejectedAt:                    handlers.FmtDateTimePtr(s.RejectedAt),
		ETag:                          etag.GenerateEtag(s.UpdatedAt),
		ServiceRequestDocuments:       serviceRequestDocs,
		ConvertToCustomerExpense:      *handlers.FmtBool(s.CustomerExpense),
		CustomerExpenseReason:         handlers.FmtStringPtr(s.CustomerExpenseReason),
	}
}

// MTOServiceItemModels payload
func MTOServiceItemModels(s models.MTOServiceItems, storer storage.FileStorer) ghcv2messages.MTOServiceItems {
	serviceItems := ghcv2messages.MTOServiceItems{}
	for _, item := range s {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		serviceItems = append(serviceItems, MTOServiceItemModel(&copyOfServiceItem, storer))
	}

	return serviceItems
}

// MTOServiceItemDimension payload
func MTOServiceItemDimension(d *models.MTOServiceItemDimension) *ghcv2messages.MTOServiceItemDimension {
	return &ghcv2messages.MTOServiceItemDimension{
		ID:     *handlers.FmtUUID(d.ID),
		Type:   ghcv2messages.DimensionType(d.Type),
		Length: *d.Length.Int32Ptr(),
		Height: *d.Height.Int32Ptr(),
		Width:  *d.Width.Int32Ptr(),
	}
}

// MTOServiceItemDimensions payload
func MTOServiceItemDimensions(d models.MTOServiceItemDimensions) ghcv2messages.MTOServiceItemDimensions {
	payload := make(ghcv2messages.MTOServiceItemDimensions, len(d))
	for i, item := range d {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOServiceItemDimension(&copyOfServiceItem)
	}
	return payload
}

// MTOServiceItemCustomerContact payload
func MTOServiceItemCustomerContact(c *models.MTOServiceItemCustomerContact) *ghcv2messages.MTOServiceItemCustomerContact {
	return &ghcv2messages.MTOServiceItemCustomerContact{
		Type:                       ghcv2messages.CustomerContactType(c.Type),
		DateOfContact:              *handlers.FmtDate(c.DateOfContact),
		TimeMilitary:               c.TimeMilitary,
		FirstAvailableDeliveryDate: *handlers.FmtDate(c.FirstAvailableDeliveryDate),
	}
}

// MTOServiceItemCustomerContacts payload
func MTOServiceItemCustomerContacts(c models.MTOServiceItemCustomerContacts) ghcv2messages.MTOServiceItemCustomerContacts {
	payload := make(ghcv2messages.MTOServiceItemCustomerContacts, len(c))
	for i, item := range c {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOServiceItemCustomerContact(&copyOfServiceItem)
	}
	return payload
}

// SITAddressUpdate payload
func SITAddressUpdate(u models.SITAddressUpdate) *ghcv2messages.SITAddressUpdate {
	return &ghcv2messages.SITAddressUpdate{
		ID:                *handlers.FmtUUID(u.ID),
		MtoServiceItemID:  *handlers.FmtUUID(u.MTOServiceItemID),
		Distance:          handlers.FmtInt64(int64(u.Distance)),
		ContractorRemarks: u.ContractorRemarks,
		OfficeRemarks:     u.OfficeRemarks,
		Status:            u.Status,
		OldAddress:        Address(&u.OldAddress),
		NewAddress:        Address(&u.NewAddress),
		CreatedAt:         strfmt.DateTime(u.CreatedAt),
		UpdatedAt:         strfmt.DateTime(u.UpdatedAt),
		ETag:              etag.GenerateEtag(u.UpdatedAt)}
}

// SITAddressUpdates payload
func SITAddressUpdates(u models.SITAddressUpdates) ghcv2messages.SITAddressUpdates {
	payload := make(ghcv2messages.SITAddressUpdates, len(u))
	for i, item := range u {
		payload[i] = SITAddressUpdate(item)
	}
	return payload
}

// Upload payload
func Upload(storer storage.FileStorer, upload models.Upload, url string) *ghcv2messages.Upload {
	uploadPayload := &ghcv2messages.Upload{
		ID:          handlers.FmtUUIDValue(upload.ID),
		Filename:    upload.Filename,
		ContentType: upload.ContentType,
		URL:         strfmt.URI(url),
		Bytes:       upload.Bytes,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload
}

// Upload payload for when a Proof of Service doc is designated as a weight ticket
// This adds an isWeightTicket key to the payload for the UI to use
func WeightTicketUpload(storer storage.FileStorer, upload models.Upload, url string, isWeightTicket bool) *ghcv2messages.Upload {
	uploadPayload := &ghcv2messages.Upload{
		ID:             handlers.FmtUUIDValue(upload.ID),
		Filename:       upload.Filename,
		ContentType:    upload.ContentType,
		URL:            strfmt.URI(url),
		Bytes:          upload.Bytes,
		CreatedAt:      strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:      strfmt.DateTime(upload.UpdatedAt),
		IsWeightTicket: isWeightTicket,
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload
}

func PayloadForUploadModel(
	storer storage.FileStorer,
	upload models.Upload,
	url string,
) *ghcv2messages.Upload {
	uploadPayload := &ghcv2messages.Upload{
		ID:          handlers.FmtUUIDValue(upload.ID),
		Filename:    upload.Filename,
		ContentType: upload.ContentType,
		URL:         strfmt.URI(url),
		Bytes:       upload.Bytes,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload
}

func PayloadForDocumentModel(storer storage.FileStorer, document models.Document) (*ghcv2messages.Document, error) {
	uploads := make([]*ghcv2messages.Upload, len(document.UserUploads))
	for i, userUpload := range document.UserUploads {
		if userUpload.Upload.ID == uuid.Nil {
			return nil, errors.New("no uploads for user")
		}
		url, err := storer.PresignedURL(userUpload.Upload.StorageKey, userUpload.Upload.ContentType)
		if err != nil {
			return nil, err
		}

		uploadPayload := PayloadForUploadModel(storer, userUpload.Upload, url)
		uploads[i] = uploadPayload
	}

	documentPayload := &ghcv2messages.Document{
		ID:              handlers.FmtUUID(document.ID),
		ServiceMemberID: handlers.FmtUUID(document.ServiceMemberID),
		Uploads:         uploads,
	}
	return documentPayload, nil
}

var (
	// QueuePaymentRequestPaymentRequested status payment requested
	QueuePaymentRequestPaymentRequested = "Payment requested"
	// QueuePaymentRequestReviewed status Payment request reviewed
	QueuePaymentRequestReviewed = "Reviewed"
	// QueuePaymentRequestRejected status Payment request rejected
	QueuePaymentRequestRejected = "Rejected"
	// QueuePaymentRequestPaid status PaymentRequest paid
	QueuePaymentRequestPaid = "Paid"
	// QueuePaymentRequestDeprecated status PaymentRequest deprecated
	QueuePaymentRequestDeprecated = "Deprecated"
	// QueuePaymentRequestError status PaymentRequest error
	QueuePaymentRequestError = "Error"
)

// Reweigh payload
func Reweigh(reweigh *models.Reweigh, _ *ghcv2messages.SITStatus) *ghcv2messages.Reweigh {
	if reweigh == nil || reweigh.ID == uuid.Nil {
		return nil
	}
	payload := &ghcv2messages.Reweigh{
		ID:                     strfmt.UUID(reweigh.ID.String()),
		RequestedAt:            strfmt.DateTime(reweigh.RequestedAt),
		RequestedBy:            ghcv2messages.ReweighRequester(reweigh.RequestedBy),
		VerificationReason:     reweigh.VerificationReason,
		Weight:                 handlers.FmtPoundPtr(reweigh.Weight),
		VerificationProvidedAt: handlers.FmtDateTimePtr(reweigh.VerificationProvidedAt),
		ShipmentID:             strfmt.UUID(reweigh.ShipmentID.String()),
	}

	return payload
}
