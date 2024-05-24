package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/primev3messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// MoveTaskOrder payload
func MoveTaskOrder(moveTaskOrder *models.Move) *primev3messages.MoveTaskOrder {
	if moveTaskOrder == nil {
		return nil
	}
	paymentRequests := PaymentRequests(&moveTaskOrder.PaymentRequests)
	mtoServiceItems := MTOServiceItems(&moveTaskOrder.MTOServiceItems)
	mtoShipments := MTOShipmentsWithoutServiceItems(&moveTaskOrder.MTOShipments)

	payload := &primev3messages.MoveTaskOrder{
		ID:                         strfmt.UUID(moveTaskOrder.ID.String()),
		MoveCode:                   moveTaskOrder.Locator,
		CreatedAt:                  strfmt.DateTime(moveTaskOrder.CreatedAt),
		AvailableToPrimeAt:         handlers.FmtDateTimePtr(moveTaskOrder.AvailableToPrimeAt),
		PrimeCounselingCompletedAt: handlers.FmtDateTimePtr(moveTaskOrder.PrimeCounselingCompletedAt),
		ExcessWeightQualifiedAt:    handlers.FmtDateTimePtr(moveTaskOrder.ExcessWeightQualifiedAt),
		ExcessWeightAcknowledgedAt: handlers.FmtDateTimePtr(moveTaskOrder.ExcessWeightAcknowledgedAt),
		ExcessWeightUploadID:       handlers.FmtUUIDPtr(moveTaskOrder.ExcessWeightUploadID),
		OrderID:                    strfmt.UUID(moveTaskOrder.OrdersID.String()),
		Order:                      Order(&moveTaskOrder.Orders),
		ReferenceID:                *moveTaskOrder.ReferenceID,
		PaymentRequests:            *paymentRequests,
		MtoShipments:               *mtoShipments,
		ContractNumber:             moveTaskOrder.Contractor.ContractNumber,
		UpdatedAt:                  strfmt.DateTime(moveTaskOrder.UpdatedAt),
		ETag:                       etag.GenerateEtag(moveTaskOrder.UpdatedAt),
	}

	if moveTaskOrder.PPMType != nil {
		payload.PpmType = *moveTaskOrder.PPMType
	}

	// mto service item references a polymorphic type which auto-generates an interface and getters and setters
	payload.SetMtoServiceItems(*mtoServiceItems)

	return payload
}

// Customer payload
func Customer(customer *models.ServiceMember) *primev3messages.Customer {
	if customer == nil {
		return nil
	}
	payload := primev3messages.Customer{
		FirstName:      swag.StringValue(customer.FirstName),
		LastName:       swag.StringValue(customer.LastName),
		DodID:          swag.StringValue(customer.Edipi),
		ID:             strfmt.UUID(customer.ID.String()),
		UserID:         strfmt.UUID(customer.UserID.String()),
		CurrentAddress: Address(customer.ResidentialAddress),
		ETag:           etag.GenerateEtag(customer.UpdatedAt),
		Branch:         swag.StringValue((*string)(customer.Affiliation)),
	}

	if customer.Telephone != nil {
		payload.Phone = *customer.Telephone
	}

	if customer.PersonalEmail != nil {
		payload.Email = *customer.PersonalEmail
	}
	return &payload
}

// Order payload
func Order(order *models.Order) *primev3messages.Order {
	if order == nil {
		return nil
	}
	destinationDutyLocation := DutyLocation(&order.NewDutyLocation)
	originDutyLocation := DutyLocation(order.OriginDutyLocation)
	if order.Grade != nil && order.Entitlement != nil {
		order.Entitlement.SetWeightAllotment(string(*order.Grade))
	}

	var grade string
	if order.Grade != nil {
		grade = string(*order.Grade)
	}

	payload := primev3messages.Order{
		CustomerID:                     strfmt.UUID(order.ServiceMemberID.String()),
		Customer:                       Customer(&order.ServiceMember),
		DestinationDutyLocation:        destinationDutyLocation,
		Entitlement:                    Entitlement(order.Entitlement),
		ID:                             strfmt.UUID(order.ID.String()),
		OriginDutyLocation:             originDutyLocation,
		OriginDutyLocationGBLOC:        swag.StringValue(order.OriginDutyLocationGBLOC),
		OrderNumber:                    order.OrdersNumber,
		LinesOfAccounting:              order.TAC,
		Rank:                           &grade, // Convert prime API "Rank" into our internal tracking of "Grade"
		ETag:                           etag.GenerateEtag(order.UpdatedAt),
		ReportByDate:                   strfmt.Date(order.ReportByDate),
		OrdersType:                     primev3messages.OrdersType(order.OrdersType),
		SupplyAndServicesCostEstimate:  order.SupplyAndServicesCostEstimate,
		PackingAndShippingInstructions: order.PackingAndShippingInstructions,
		MethodOfPayment:                order.MethodOfPayment,
		Naics:                          order.NAICS,
	}

	if payload.Customer.Branch == "MARINES" {
		payload.OriginDutyLocationGBLOC = "USMC"
	}

	return &payload
}

// Entitlement payload
func Entitlement(entitlement *models.Entitlement) *primev3messages.Entitlements {
	if entitlement == nil {
		return nil
	}
	var totalWeight int64
	if weightAllowance := entitlement.WeightAllowance(); weightAllowance != nil {
		totalWeight = int64(*weightAllowance)
	}
	var authorizedWeight *int64
	if entitlement.AuthorizedWeight() != nil {
		aw := int64(*entitlement.AuthorizedWeight())
		authorizedWeight = &aw
	}
	var sit int64
	if entitlement.StorageInTransit != nil {
		sit = int64(*entitlement.StorageInTransit)
	}
	var totalDependents int64
	if entitlement.TotalDependents != nil {
		totalDependents = int64(*entitlement.TotalDependents)
	}
	return &primev3messages.Entitlements{
		ID:                             strfmt.UUID(entitlement.ID.String()),
		AuthorizedWeight:               authorizedWeight,
		DependentsAuthorized:           entitlement.DependentsAuthorized,
		NonTemporaryStorage:            entitlement.NonTemporaryStorage,
		PrivatelyOwnedVehicle:          entitlement.PrivatelyOwnedVehicle,
		ProGearWeight:                  int64(entitlement.ProGearWeight),
		ProGearWeightSpouse:            int64(entitlement.ProGearWeightSpouse),
		RequiredMedicalEquipmentWeight: int64(entitlement.RequiredMedicalEquipmentWeight),
		OrganizationalClothingAndIndividualEquipment: entitlement.OrganizationalClothingAndIndividualEquipment,
		StorageInTransit: sit,
		TotalDependents:  totalDependents,
		TotalWeight:      totalWeight,
		ETag:             etag.GenerateEtag(entitlement.UpdatedAt),
	}
}

// DutyLocation payload
func DutyLocation(dutyLocation *models.DutyLocation) *primev3messages.DutyLocation {
	if dutyLocation == nil {
		return nil
	}
	address := Address(&dutyLocation.Address)
	payload := primev3messages.DutyLocation{
		Address:   address,
		AddressID: address.ID,
		ID:        strfmt.UUID(dutyLocation.ID.String()),
		Name:      dutyLocation.Name,
	}
	return &payload
}

// Address payload
func Address(address *models.Address) *primev3messages.Address {
	if address == nil {
		return nil
	}
	return &primev3messages.Address{
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
func StorageFacility(storage *models.StorageFacility) *primev3messages.StorageFacility {
	if storage == nil {
		return nil
	}

	return &primev3messages.StorageFacility{
		ID:           strfmt.UUID(storage.ID.String()),
		Address:      Address(&storage.Address),
		ETag:         etag.GenerateEtag(storage.UpdatedAt),
		Email:        storage.Email,
		FacilityName: storage.FacilityName,
		LotNumber:    storage.LotNumber,
		Phone:        storage.Phone,
	}
}

// MTOAgent payload
func MTOAgent(mtoAgent *models.MTOAgent) *primev3messages.MTOAgent {
	if mtoAgent == nil {
		return nil
	}

	return &primev3messages.MTOAgent{
		AgentType:     primev3messages.MTOAgentType(mtoAgent.MTOAgentType),
		FirstName:     mtoAgent.FirstName,
		LastName:      mtoAgent.LastName,
		Phone:         mtoAgent.Phone,
		Email:         mtoAgent.Email,
		ID:            strfmt.UUID(mtoAgent.ID.String()),
		MtoShipmentID: strfmt.UUID(mtoAgent.MTOShipmentID.String()),
		CreatedAt:     strfmt.DateTime(mtoAgent.CreatedAt),
		UpdatedAt:     strfmt.DateTime(mtoAgent.UpdatedAt),
		ETag:          etag.GenerateEtag(mtoAgent.UpdatedAt),
	}
}

// MTOAgents payload
func MTOAgents(mtoAgents *models.MTOAgents) *primev3messages.MTOAgents {
	if mtoAgents == nil {
		return nil
	}

	agents := make(primev3messages.MTOAgents, len(*mtoAgents))

	for i, m := range *mtoAgents {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		agents[i] = MTOAgent(&copyOfM)
	}

	return &agents
}

func ProofOfServiceDoc(proofOfServiceDoc models.ProofOfServiceDoc) *primev3messages.ProofOfServiceDoc {
	uploads := make([]*primev3messages.UploadWithOmissions, len(proofOfServiceDoc.PrimeUploads))
	if proofOfServiceDoc.PrimeUploads != nil && len(proofOfServiceDoc.PrimeUploads) > 0 {
		for i, primeUpload := range proofOfServiceDoc.PrimeUploads { //#nosec G601
			uploads[i] = basicUpload(&primeUpload.Upload) //#nosec G601
		}
	}

	return &primev3messages.ProofOfServiceDoc{
		Uploads: uploads,
	}
}

// PaymentRequest payload
func PaymentRequest(paymentRequest *models.PaymentRequest) *primev3messages.PaymentRequest {
	if paymentRequest == nil {
		return nil
	}

	serviceDocs := make(primev3messages.ProofOfServiceDocs, len(paymentRequest.ProofOfServiceDocs))

	if paymentRequest.ProofOfServiceDocs != nil && len(paymentRequest.ProofOfServiceDocs) > 0 {
		for i, proofOfService := range paymentRequest.ProofOfServiceDocs {
			serviceDocs[i] = ProofOfServiceDoc(proofOfService)
		}
	}

	paymentServiceItems := PaymentServiceItems(&paymentRequest.PaymentServiceItems)
	return &primev3messages.PaymentRequest{
		ID:                              strfmt.UUID(paymentRequest.ID.String()),
		IsFinal:                         &paymentRequest.IsFinal,
		MoveTaskOrderID:                 strfmt.UUID(paymentRequest.MoveTaskOrderID.String()),
		PaymentRequestNumber:            paymentRequest.PaymentRequestNumber,
		RecalculationOfPaymentRequestID: handlers.FmtUUIDPtr(paymentRequest.RecalculationOfPaymentRequestID),
		RejectionReason:                 paymentRequest.RejectionReason,
		Status:                          primev3messages.PaymentRequestStatus(paymentRequest.Status),
		PaymentServiceItems:             *paymentServiceItems,
		ProofOfServiceDocs:              serviceDocs,
		ETag:                            etag.GenerateEtag(paymentRequest.UpdatedAt),
	}
}

// PaymentRequests payload
func PaymentRequests(paymentRequests *models.PaymentRequests) *primev3messages.PaymentRequests {
	if paymentRequests == nil {
		return nil
	}

	payload := make(primev3messages.PaymentRequests, len(*paymentRequests))

	for i, p := range *paymentRequests {
		copyOfP := p // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentRequest(&copyOfP)
	}
	return &payload
}

// PaymentServiceItem payload
func PaymentServiceItem(paymentServiceItem *models.PaymentServiceItem) *primev3messages.PaymentServiceItem {
	if paymentServiceItem == nil {
		return nil
	}

	paymentServiceItemParams := PaymentServiceItemParams(&paymentServiceItem.PaymentServiceItemParams)

	payload := &primev3messages.PaymentServiceItem{
		ID:                       strfmt.UUID(paymentServiceItem.ID.String()),
		PaymentRequestID:         strfmt.UUID(paymentServiceItem.PaymentRequestID.String()),
		MtoServiceItemID:         strfmt.UUID(paymentServiceItem.MTOServiceItemID.String()),
		Status:                   primev3messages.PaymentServiceItemStatus(paymentServiceItem.Status),
		RejectionReason:          paymentServiceItem.RejectionReason,
		ReferenceID:              paymentServiceItem.ReferenceID,
		PaymentServiceItemParams: *paymentServiceItemParams,
		ETag:                     etag.GenerateEtag(paymentServiceItem.UpdatedAt),
	}

	if paymentServiceItem.PriceCents != nil {
		payload.PriceCents = models.Int64Pointer(int64(*paymentServiceItem.PriceCents))
	}

	return payload
}

// PaymentServiceItems payload
func PaymentServiceItems(paymentServiceItems *models.PaymentServiceItems) *primev3messages.PaymentServiceItems {
	if paymentServiceItems == nil {
		return nil
	}

	payload := make(primev3messages.PaymentServiceItems, len(*paymentServiceItems))

	for i, p := range *paymentServiceItems {
		copyOfP := p // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentServiceItem(&copyOfP)
	}
	return &payload
}

// PaymentServiceItemParam payload
func PaymentServiceItemParam(paymentServiceItemParam *models.PaymentServiceItemParam) *primev3messages.PaymentServiceItemParam {
	if paymentServiceItemParam == nil {
		return nil
	}

	return &primev3messages.PaymentServiceItemParam{
		ID:                   strfmt.UUID(paymentServiceItemParam.ID.String()),
		PaymentServiceItemID: strfmt.UUID(paymentServiceItemParam.PaymentServiceItemID.String()),
		Key:                  primev3messages.ServiceItemParamName(paymentServiceItemParam.ServiceItemParamKey.Key),
		Value:                paymentServiceItemParam.Value,
		Type:                 primev3messages.ServiceItemParamType(paymentServiceItemParam.ServiceItemParamKey.Type),
		Origin:               primev3messages.ServiceItemParamOrigin(paymentServiceItemParam.ServiceItemParamKey.Origin),
		ETag:                 etag.GenerateEtag(paymentServiceItemParam.UpdatedAt),
	}
}

// PaymentServiceItemParams payload
func PaymentServiceItemParams(paymentServiceItemParams *models.PaymentServiceItemParams) *primev3messages.PaymentServiceItemParams {
	if paymentServiceItemParams == nil {
		return nil
	}

	payload := make(primev3messages.PaymentServiceItemParams, len(*paymentServiceItemParams))

	for i, p := range *paymentServiceItemParams {
		copyOfP := p // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentServiceItemParam(&copyOfP)
	}
	return &payload
}

func ServiceRequestDocument(serviceRequestDocument models.ServiceRequestDocument) *primev3messages.ServiceRequestDocument {
	uploads := make([]*primev3messages.UploadWithOmissions, len(serviceRequestDocument.ServiceRequestDocumentUploads))
	if serviceRequestDocument.ServiceRequestDocumentUploads != nil && len(serviceRequestDocument.ServiceRequestDocumentUploads) > 0 {
		for i, proofOfServiceDocumentUpload := range serviceRequestDocument.ServiceRequestDocumentUploads {
			uploads[i] = basicUpload(&proofOfServiceDocumentUpload.Upload) //#nosec G601
		}
	}

	return &primev3messages.ServiceRequestDocument{
		Uploads: uploads,
	}
}

// PPMShipment payload
func PPMShipment(ppmShipment *models.PPMShipment) *primev3messages.PPMShipment {
	if ppmShipment == nil || ppmShipment.ID.IsNil() {
		return nil
	}

	payloadPPMShipment := &primev3messages.PPMShipment{
		ID:                             *handlers.FmtUUID(ppmShipment.ID),
		ShipmentID:                     *handlers.FmtUUID(ppmShipment.ShipmentID),
		CreatedAt:                      strfmt.DateTime(ppmShipment.CreatedAt),
		UpdatedAt:                      strfmt.DateTime(ppmShipment.UpdatedAt),
		Status:                         primev3messages.PPMShipmentStatus(ppmShipment.Status),
		ExpectedDepartureDate:          handlers.FmtDate(ppmShipment.ExpectedDepartureDate),
		ActualMoveDate:                 handlers.FmtDatePtr(ppmShipment.ActualMoveDate),
		SubmittedAt:                    handlers.FmtDateTimePtr(ppmShipment.SubmittedAt),
		ReviewedAt:                     handlers.FmtDateTimePtr(ppmShipment.ReviewedAt),
		ApprovedAt:                     handlers.FmtDateTimePtr(ppmShipment.ApprovedAt),
		HasSecondaryPickupAddress:      ppmShipment.HasSecondaryPickupAddress,
		HasSecondaryDestinationAddress: ppmShipment.HasSecondaryDestinationAddress,
		ActualPickupPostalCode:         ppmShipment.ActualPickupPostalCode,
		ActualDestinationPostalCode:    ppmShipment.ActualDestinationPostalCode,
		SitExpected:                    ppmShipment.SITExpected,
		SitEstimatedWeight:             handlers.FmtPoundPtr(ppmShipment.SITEstimatedWeight),
		SitEstimatedEntryDate:          handlers.FmtDatePtr(ppmShipment.SITEstimatedEntryDate),
		SitEstimatedDepartureDate:      handlers.FmtDatePtr(ppmShipment.SITEstimatedDepartureDate),
		SitEstimatedCost:               handlers.FmtCost(ppmShipment.SITEstimatedCost),
		EstimatedWeight:                handlers.FmtPoundPtr(ppmShipment.EstimatedWeight),
		EstimatedIncentive:             handlers.FmtCost(ppmShipment.EstimatedIncentive),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.FmtPoundPtr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.FmtPoundPtr(ppmShipment.SpouseProGearWeight),
		HasRequestedAdvance:            ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:         handlers.FmtCost(ppmShipment.AdvanceAmountRequested),
		HasReceivedAdvance:             ppmShipment.HasReceivedAdvance,
		AdvanceAmountReceived:          handlers.FmtCost(ppmShipment.AdvanceAmountReceived),
		ETag:                           etag.GenerateEtag(ppmShipment.UpdatedAt),
	}

	if ppmShipment.SITLocation != nil {
		sitLocation := primev3messages.SITLocationType(*ppmShipment.SITLocation)
		payloadPPMShipment.SitLocation = &sitLocation
	}

	// Set up address payloads
	if ppmShipment.PickupAddress != nil {
		payloadPPMShipment.PickupAddress = Address(ppmShipment.PickupAddress)
	}
	if ppmShipment.DestinationAddress != nil {
		payloadPPMShipment.DestinationAddress = Address(ppmShipment.DestinationAddress)
	}
	if ppmShipment.SecondaryPickupAddress != nil {
		payloadPPMShipment.SecondaryPickupAddress = Address(ppmShipment.SecondaryPickupAddress)
	}
	if ppmShipment.SecondaryDestinationAddress != nil {
		payloadPPMShipment.SecondaryDestinationAddress = Address(ppmShipment.SecondaryDestinationAddress)
	}

	return payloadPPMShipment
}

func MTOShipmentWithoutServiceItems(mtoShipment *models.MTOShipment) *primev3messages.MTOShipmentWithoutServiceItems {
	payload := &primev3messages.MTOShipmentWithoutServiceItems{
		ID:                               strfmt.UUID(mtoShipment.ID.String()),
		ActualPickupDate:                 handlers.FmtDatePtr(mtoShipment.ActualPickupDate),
		ApprovedDate:                     handlers.FmtDatePtr(mtoShipment.ApprovedDate),
		FirstAvailableDeliveryDate:       handlers.FmtDatePtr(mtoShipment.FirstAvailableDeliveryDate),
		PrimeEstimatedWeightRecordedDate: handlers.FmtDatePtr(mtoShipment.PrimeEstimatedWeightRecordedDate),
		RequestedPickupDate:              handlers.FmtDatePtr(mtoShipment.RequestedPickupDate),
		RequestedDeliveryDate:            handlers.FmtDatePtr(mtoShipment.RequestedDeliveryDate),
		RequiredDeliveryDate:             handlers.FmtDatePtr(mtoShipment.RequiredDeliveryDate),
		ScheduledPickupDate:              handlers.FmtDatePtr(mtoShipment.ScheduledPickupDate),
		ScheduledDeliveryDate:            handlers.FmtDatePtr(mtoShipment.ScheduledDeliveryDate),
		ActualDeliveryDate:               handlers.FmtDatePtr(mtoShipment.ActualDeliveryDate),
		Agents:                           *MTOAgents(&mtoShipment.MTOAgents),
		SitExtensions:                    *SITDurationUpdates(&mtoShipment.SITDurationUpdates),
		Reweigh:                          Reweigh(mtoShipment.Reweigh),
		MoveTaskOrderID:                  strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:                     primev3messages.MTOShipmentType(mtoShipment.ShipmentType),
		CustomerRemarks:                  mtoShipment.CustomerRemarks,
		CounselorRemarks:                 mtoShipment.CounselorRemarks,
		Status:                           string(mtoShipment.Status),
		Diversion:                        bool(mtoShipment.Diversion),
		DeliveryAddressUpdate:            ShipmentAddressUpdate(mtoShipment.DeliveryAddressUpdate),
		CreatedAt:                        strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                        strfmt.DateTime(mtoShipment.UpdatedAt),
		PpmShipment:                      PPMShipment(mtoShipment.PPMShipment),
		ETag:                             etag.GenerateEtag(mtoShipment.UpdatedAt),
		OriginSitAuthEndDate:             (*strfmt.Date)(mtoShipment.OriginSITAuthEndDate),
		DestinationSitAuthEndDate:        (*strfmt.Date)(mtoShipment.DestinationSITAuthEndDate),
	}

	// Set up address payloads
	if mtoShipment.PickupAddress != nil {
		payload.PickupAddress.Address = *Address(mtoShipment.PickupAddress)
	}
	if mtoShipment.DestinationAddress != nil {
		payload.DestinationAddress.Address = *Address(mtoShipment.DestinationAddress)
	}
	if mtoShipment.DestinationType != nil {
		destinationType := primev3messages.DestinationType(*mtoShipment.DestinationType)
		payload.DestinationType = &destinationType
	}
	if mtoShipment.SecondaryPickupAddress != nil {
		payload.SecondaryPickupAddress.Address = *Address(mtoShipment.SecondaryPickupAddress)
	}
	if mtoShipment.SecondaryDeliveryAddress != nil {
		payload.SecondaryDeliveryAddress.Address = *Address(mtoShipment.SecondaryDeliveryAddress)
	}

	if mtoShipment.StorageFacility != nil {
		payload.StorageFacility = StorageFacility(mtoShipment.StorageFacility)
	}

	if mtoShipment.PrimeEstimatedWeight != nil {
		payload.PrimeEstimatedWeight = handlers.FmtInt64(mtoShipment.PrimeEstimatedWeight.Int64())
	}

	if mtoShipment.PrimeActualWeight != nil {
		payload.PrimeActualWeight = handlers.FmtInt64(mtoShipment.PrimeActualWeight.Int64())
	}

	if mtoShipment.NTSRecordedWeight != nil {
		payload.NtsRecordedWeight = handlers.FmtInt64(mtoShipment.NTSRecordedWeight.Int64())
	}

	if mtoShipment.ShipmentType == models.MTOShipmentTypePPM {
		if mtoShipment.PPMShipment.PickupAddress != nil {
			payload.PpmShipment.PickupAddress = Address(mtoShipment.PPMShipment.PickupAddress)
		}
		if mtoShipment.PPMShipment.SecondaryPickupAddress != nil {
			payload.PpmShipment.SecondaryPickupAddress = Address(mtoShipment.PPMShipment.SecondaryPickupAddress)
		}
		if mtoShipment.PPMShipment.DestinationAddress != nil {
			payload.PpmShipment.DestinationAddress = Address(mtoShipment.PPMShipment.DestinationAddress)
		}
		if mtoShipment.PPMShipment.SecondaryDestinationAddress != nil {
			payload.PpmShipment.SecondaryDestinationAddress = Address(mtoShipment.PPMShipment.SecondaryDestinationAddress)
		}
		payload.PpmShipment.HasSecondaryPickupAddress = mtoShipment.PPMShipment.HasSecondaryPickupAddress
		payload.PpmShipment.HasSecondaryDestinationAddress = mtoShipment.PPMShipment.HasSecondaryDestinationAddress
	}

	return payload
}

func MTOShipmentsWithoutServiceItems(mtoShipments *models.MTOShipments) *primev3messages.MTOShipmentsWithoutServiceObjects {
	payload := make(primev3messages.MTOShipmentsWithoutServiceObjects, len(*mtoShipments))

	for i, m := range *mtoShipments {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOShipmentWithoutServiceItems(&copyOfM)
	}
	return &payload
}

// MTOServiceItem payload
func MTOServiceItem(mtoServiceItem *models.MTOServiceItem) primev3messages.MTOServiceItem {
	var payload primev3messages.MTOServiceItem
	// here we determine which payload model to use based on the re service code
	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDOPSIT, models.ReServiceCodeDOSFSC:
		var sitDepartureDate time.Time
		if mtoServiceItem.SITDepartureDate != nil {
			sitDepartureDate = *mtoServiceItem.SITDepartureDate
		}
		payload = &primev3messages.MTOServiceItemOriginSIT{
			ReServiceCode:        handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Reason:               mtoServiceItem.Reason,
			SitDepartureDate:     handlers.FmtDate(sitDepartureDate),
			SitEntryDate:         handlers.FmtDatePtr(mtoServiceItem.SITEntryDate),
			SitPostalCode:        mtoServiceItem.SITPostalCode,
			SitHHGActualOrigin:   Address(mtoServiceItem.SITOriginHHGActualAddress),
			SitHHGOriginalOrigin: Address(mtoServiceItem.SITOriginHHGOriginalAddress),
		}
	case models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT, models.ReServiceCodeDDSFSC:
		var sitDepartureDate, firstAvailableDeliveryDate1, firstAvailableDeliveryDate2, dateOfContact1, dateOfContact2 time.Time
		var timeMilitary1, timeMilitary2 *string

		if mtoServiceItem.SITDepartureDate != nil {
			sitDepartureDate = *mtoServiceItem.SITDepartureDate
		}

		firstContact := GetCustomerContact(mtoServiceItem.CustomerContacts, models.CustomerContactTypeFirst)
		secondContact := GetCustomerContact(mtoServiceItem.CustomerContacts, models.CustomerContactTypeSecond)
		timeMilitary1 = &firstContact.TimeMilitary
		timeMilitary2 = &secondContact.TimeMilitary

		if !firstContact.DateOfContact.IsZero() {
			dateOfContact1 = firstContact.DateOfContact
		}

		if !secondContact.DateOfContact.IsZero() {
			dateOfContact2 = secondContact.DateOfContact
		}

		if !firstContact.FirstAvailableDeliveryDate.IsZero() {
			firstAvailableDeliveryDate1 = firstContact.FirstAvailableDeliveryDate
		}

		if !secondContact.FirstAvailableDeliveryDate.IsZero() {
			firstAvailableDeliveryDate2 = secondContact.FirstAvailableDeliveryDate
		}

		payload = &primev3messages.MTOServiceItemDestSIT{
			ReServiceCode:               handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Reason:                      mtoServiceItem.Reason,
			DateOfContact1:              handlers.FmtDate(dateOfContact1),
			TimeMilitary1:               handlers.FmtStringPtrNonEmpty(timeMilitary1),
			FirstAvailableDeliveryDate1: handlers.FmtDate(firstAvailableDeliveryDate1),
			DateOfContact2:              handlers.FmtDate(dateOfContact2),
			TimeMilitary2:               handlers.FmtStringPtrNonEmpty(timeMilitary2),
			FirstAvailableDeliveryDate2: handlers.FmtDate(firstAvailableDeliveryDate2),
			SitDepartureDate:            handlers.FmtDate(sitDepartureDate),
			SitEntryDate:                handlers.FmtDatePtr(mtoServiceItem.SITEntryDate),
			SitDestinationFinalAddress:  Address(mtoServiceItem.SITDestinationFinalAddress),
			SitAddressUpdates:           SITAddressUpdates(mtoServiceItem.SITAddressUpdates),
			SitCustomerContacted:        handlers.FmtDatePtr(mtoServiceItem.SITCustomerContacted),
			SitRequestedDelivery:        handlers.FmtDatePtr(mtoServiceItem.SITRequestedDelivery),
		}

	case models.ReServiceCodeDCRT, models.ReServiceCodeDUCRT:
		item := GetDimension(mtoServiceItem.Dimensions, models.DimensionTypeItem)
		crate := GetDimension(mtoServiceItem.Dimensions, models.DimensionTypeCrate)
		cratingSI := primev3messages.MTOServiceItemDomesticCrating{
			ReServiceCode: handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Description:   mtoServiceItem.Description,
			Reason:        mtoServiceItem.Reason,
		}
		cratingSI.Item.MTOServiceItemDimension = primev3messages.MTOServiceItemDimension{
			ID:     strfmt.UUID(item.ID.String()),
			Height: item.Height.Int32Ptr(),
			Length: item.Length.Int32Ptr(),
			Width:  item.Width.Int32Ptr(),
		}
		cratingSI.Crate.MTOServiceItemDimension = primev3messages.MTOServiceItemDimension{
			ID:     strfmt.UUID(crate.ID.String()),
			Height: crate.Height.Int32Ptr(),
			Length: crate.Length.Int32Ptr(),
			Width:  crate.Width.Int32Ptr(),
		}
		payload = &cratingSI
	case models.ReServiceCodeDDSHUT, models.ReServiceCodeDOSHUT:
		payload = &primev3messages.MTOServiceItemShuttle{
			ReServiceCode:   handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Reason:          mtoServiceItem.Reason,
			EstimatedWeight: handlers.FmtPoundPtr(mtoServiceItem.EstimatedWeight),
			ActualWeight:    handlers.FmtPoundPtr(mtoServiceItem.ActualWeight),
		}
	default:
		// otherwise, basic service item
		payload = &primev3messages.MTOServiceItemBasic{
			ReServiceCode: primev3messages.NewReServiceCode(primev3messages.ReServiceCode(mtoServiceItem.ReService.Code)),
		}
	}

	// set all relevant fields that apply to all service items
	var shipmentIDStr string
	if mtoServiceItem.MTOShipmentID != nil {
		shipmentIDStr = mtoServiceItem.MTOShipmentID.String()
	}

	serviceRequestDocuments := make(primev3messages.ServiceRequestDocuments, len(mtoServiceItem.ServiceRequestDocuments))

	if mtoServiceItem.ServiceRequestDocuments != nil && len(mtoServiceItem.ServiceRequestDocuments) > 0 {
		for i, serviceRequestDocument := range mtoServiceItem.ServiceRequestDocuments {
			serviceRequestDocuments[i] = ServiceRequestDocument(serviceRequestDocument)
		}
	}

	one := mtoServiceItem.ID.String()
	two := strfmt.UUID(one)
	payload.SetID(two)
	payload.SetMoveTaskOrderID(handlers.FmtUUID(mtoServiceItem.MoveTaskOrderID))
	payload.SetMtoShipmentID(strfmt.UUID(shipmentIDStr))
	payload.SetReServiceName(mtoServiceItem.ReService.Name)
	payload.SetStatus(primev3messages.MTOServiceItemStatus(mtoServiceItem.Status))
	payload.SetRejectionReason(mtoServiceItem.RejectionReason)
	payload.SetETag(etag.GenerateEtag(mtoServiceItem.UpdatedAt))
	payload.SetServiceRequestDocuments(serviceRequestDocuments)
	return payload
}

// MTOServiceItems payload
func MTOServiceItems(mtoServiceItems *models.MTOServiceItems) *[]primev3messages.MTOServiceItem {
	payload := []primev3messages.MTOServiceItem{}

	for _, p := range *mtoServiceItems {
		copyOfP := p // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload = append(payload, MTOServiceItem(&copyOfP))
	}
	return &payload
}

// Reweigh returns the reweigh payload
func Reweigh(reweigh *models.Reweigh) *primev3messages.Reweigh {
	if reweigh == nil || reweigh.ID == uuid.Nil {
		return nil
	}

	payload := &primev3messages.Reweigh{
		ID:                     strfmt.UUID(reweigh.ID.String()),
		ShipmentID:             strfmt.UUID(reweigh.ShipmentID.String()),
		RequestedAt:            strfmt.DateTime(reweigh.RequestedAt),
		RequestedBy:            primev3messages.ReweighRequester(reweigh.RequestedBy),
		CreatedAt:              strfmt.DateTime(reweigh.CreatedAt),
		UpdatedAt:              strfmt.DateTime(reweigh.UpdatedAt),
		ETag:                   etag.GenerateEtag(reweigh.UpdatedAt),
		Weight:                 handlers.FmtPoundPtr(reweigh.Weight),
		VerificationReason:     handlers.FmtStringPtr(reweigh.VerificationReason),
		VerificationProvidedAt: handlers.FmtDateTimePtr(reweigh.VerificationProvidedAt),
	}

	return payload
}

func basicUpload(upload *models.Upload) *primev3messages.UploadWithOmissions {
	if upload == nil || upload.ID == uuid.Nil {
		return nil
	}

	payload := &primev3messages.UploadWithOmissions{
		ID:          strfmt.UUID(upload.ID.String()),
		Bytes:       &upload.Bytes,
		ContentType: &upload.ContentType,
		Filename:    &upload.Filename,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
	}

	return payload
}

// SITDurationUpdate payload
func SITDurationUpdate(sitDurationUpdate *models.SITDurationUpdate) *primev3messages.SITExtension {
	if sitDurationUpdate == nil {
		return nil
	}
	payload := &primev3messages.SITExtension{
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
func SITDurationUpdates(sitDurationUpdates *models.SITDurationUpdates) *primev3messages.SITExtensions {
	if sitDurationUpdates == nil {
		return nil
	}

	payload := make(primev3messages.SITExtensions, len(*sitDurationUpdates))

	for i, m := range *sitDurationUpdates {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = SITDurationUpdate(&copyOfM)
	}

	return &payload
}

// GetDimension will get the first dimension of the passed in type.
func GetDimension(dimensions models.MTOServiceItemDimensions, dimensionType models.DimensionType) models.MTOServiceItemDimension {
	if len(dimensions) == 0 {
		return models.MTOServiceItemDimension{}
	}

	for _, dimension := range dimensions {
		if dimension.Type == dimensionType {
			return dimension
		}
	}

	return models.MTOServiceItemDimension{}
}

// GetCustomerContact will get the first customer contact for destination 1st day SIT based on type.
func GetCustomerContact(customerContacts models.MTOServiceItemCustomerContacts, customerContactType models.CustomerContactType) models.MTOServiceItemCustomerContact {
	if len(customerContacts) == 0 {
		return models.MTOServiceItemCustomerContact{}
	}

	for _, customerContact := range customerContacts {
		if customerContact.Type == customerContactType {
			return customerContact
		}
	}

	return models.MTOServiceItemCustomerContact{}
}

// ShipmentAddressUpdate payload
func ShipmentAddressUpdate(shipmentAddressUpdate *models.ShipmentAddressUpdate) *primev3messages.ShipmentAddressUpdate {
	if shipmentAddressUpdate == nil || shipmentAddressUpdate.ID.IsNil() {
		return nil
	}

	payload := &primev3messages.ShipmentAddressUpdate{
		ID:                strfmt.UUID(shipmentAddressUpdate.ID.String()),
		ShipmentID:        strfmt.UUID(shipmentAddressUpdate.ShipmentID.String()),
		NewAddress:        Address(&shipmentAddressUpdate.NewAddress),
		OriginalAddress:   Address(&shipmentAddressUpdate.OriginalAddress),
		ContractorRemarks: shipmentAddressUpdate.ContractorRemarks,
		OfficeRemarks:     shipmentAddressUpdate.OfficeRemarks,
		Status:            primev3messages.ShipmentAddressUpdateStatus(shipmentAddressUpdate.Status),
	}

	return payload
}

// SITAddressUpdates payload
func SITAddressUpdates(u models.SITAddressUpdates) primev3messages.SitAddressUpdates {
	payload := make(primev3messages.SitAddressUpdates, len(u))
	for i, item := range u {
		copyOfItem := item
		payload[i] = SITAddressUpdate(&copyOfItem)
	}
	return payload
}

// SITAddressUpdate payload
func SITAddressUpdate(sitAddressUpdate *models.SITAddressUpdate) *primev3messages.SitAddressUpdate {
	if sitAddressUpdate == nil {
		return nil
	}

	payload := &primev3messages.SitAddressUpdate{
		ID:                strfmt.UUID(sitAddressUpdate.ID.String()),
		ETag:              etag.GenerateEtag(sitAddressUpdate.UpdatedAt),
		MtoServiceItemID:  strfmt.UUID(sitAddressUpdate.MTOServiceItemID.String()),
		NewAddressID:      strfmt.UUID(sitAddressUpdate.NewAddressID.String()),
		NewAddress:        Address(&sitAddressUpdate.NewAddress),
		ContractorRemarks: handlers.FmtStringPtr(sitAddressUpdate.ContractorRemarks),
		OfficeRemarks:     handlers.FmtStringPtr(sitAddressUpdate.OfficeRemarks),
		OldAddressID:      strfmt.UUID(sitAddressUpdate.OldAddressID.String()),
		OldAddress:        Address(&sitAddressUpdate.OldAddress),
		Status:            primev3messages.SitAddressUpdateStatus(sitAddressUpdate.Status),
		CreatedAt:         strfmt.DateTime(sitAddressUpdate.CreatedAt),
		UpdatedAt:         strfmt.DateTime(sitAddressUpdate.UpdatedAt),
	}

	return payload
}

// ClientError describes errors in a standard structure to be returned in the payload
func ClientError(title string, detail string, instance uuid.UUID) *primev3messages.ClientError {
	return &primev3messages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
}

// InternalServerError describes errors in a standard structure to be returned in the payload.
// If detail is nil, string defaults to "An internal server error has occurred."
func InternalServerError(detail *string, traceID uuid.UUID) *primev3messages.Error {
	payload := primev3messages.Error{
		Title:    handlers.FmtString(handlers.InternalServerErrMessage),
		Detail:   handlers.FmtString(handlers.InternalServerErrDetail),
		Instance: strfmt.UUID(traceID.String()),
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}

// ValidationError describes validation errors from the model or properties
func ValidationError(detail string, instance uuid.UUID, validationErrors *validate.Errors) *primev3messages.ValidationError {
	payload := &primev3messages.ValidationError{
		ClientError: *ClientError(handlers.ValidationErrMessage, detail, instance),
	}
	if validationErrors != nil {
		payload.InvalidFields = handlers.NewValidationErrorListResponse(validationErrors).Errors
	}
	return payload
}

// MTOShipment converts MTOShipment model to payload
func MTOShipment(mtoShipment *models.MTOShipment) *primev3messages.MTOShipment {
	payload := &primev3messages.MTOShipment{
		MTOShipmentWithoutServiceItems: *MTOShipmentWithoutServiceItems(mtoShipment),
	}

	if mtoShipment.MTOServiceItems != nil {
		payload.SetMtoServiceItems(*MTOServiceItems(&mtoShipment.MTOServiceItems))
	} else {
		payload.SetMtoServiceItems([]primev3messages.MTOServiceItem{})
	}

	return payload
}
