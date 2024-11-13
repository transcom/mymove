package payloads

import (
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
)

// MoveTaskOrder payload
func MoveTaskOrder(moveTaskOrder *models.Move) *primemessages.MoveTaskOrder {
	if moveTaskOrder == nil {
		return nil
	}
	paymentRequests := PaymentRequests(&moveTaskOrder.PaymentRequests)
	mtoServiceItems := MTOServiceItems(&moveTaskOrder.MTOServiceItems)
	mtoShipments := MTOShipmentsWithoutServiceItems(&moveTaskOrder.MTOShipments)

	payload := &primemessages.MoveTaskOrder{
		ID:                         strfmt.UUID(moveTaskOrder.ID.String()),
		MoveCode:                   moveTaskOrder.Locator,
		CreatedAt:                  strfmt.DateTime(moveTaskOrder.CreatedAt),
		AvailableToPrimeAt:         handlers.FmtDateTimePtr(moveTaskOrder.AvailableToPrimeAt),
		ApprovedAt:                 handlers.FmtDateTimePtr(moveTaskOrder.ApprovedAt),
		PrimeCounselingCompletedAt: handlers.FmtDateTimePtr(moveTaskOrder.PrimeCounselingCompletedAt),
		ExcessWeightQualifiedAt:    handlers.FmtDateTimePtr(moveTaskOrder.ExcessWeightQualifiedAt),
		ExcessWeightAcknowledgedAt: handlers.FmtDateTimePtr(moveTaskOrder.ExcessWeightAcknowledgedAt),
		ExcessWeightUploadID:       handlers.FmtUUIDPtr(moveTaskOrder.ExcessWeightUploadID),
		OrderID:                    strfmt.UUID(moveTaskOrder.OrdersID.String()),
		Order:                      Order(&moveTaskOrder.Orders),
		ReferenceID:                *moveTaskOrder.ReferenceID,
		PaymentRequests:            *paymentRequests,
		MtoShipments:               *mtoShipments,
		UpdatedAt:                  strfmt.DateTime(moveTaskOrder.UpdatedAt),
		ETag:                       etag.GenerateEtag(moveTaskOrder.UpdatedAt),
	}

	if moveTaskOrder.PPMType != nil {
		payload.PpmType = *moveTaskOrder.PPMType
	}

	// mto service item references a polymorphic type which auto-generates an interface and getters and setters
	payload.SetMtoServiceItems(*mtoServiceItems)

	// update originDutyLocationGBLOC to match TOO's gbloc and not service counselors's gbloc
	if len(moveTaskOrder.ShipmentGBLOC) > 0 && moveTaskOrder.ShipmentGBLOC[0].GBLOC != nil {
		payload.Order.OriginDutyLocationGBLOC = swag.StringValue(moveTaskOrder.ShipmentGBLOC[0].GBLOC)
	}

	return payload
}

// ListMove payload
func ListMove(move *models.Move, moveOrderAmendmentsCount *services.MoveOrderAmendmentAvailableSinceCount) *primemessages.ListMove {
	if move == nil {
		return nil
	}

	payload := &primemessages.ListMove{
		ID:                 strfmt.UUID(move.ID.String()),
		MoveCode:           move.Locator,
		CreatedAt:          strfmt.DateTime(move.CreatedAt),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(move.AvailableToPrimeAt),
		ApprovedAt:         handlers.FmtDateTimePtr(move.ApprovedAt),
		OrderID:            strfmt.UUID(move.OrdersID.String()),
		ReferenceID:        *move.ReferenceID,
		UpdatedAt:          strfmt.DateTime(move.UpdatedAt),
		ETag:               etag.GenerateEtag(move.UpdatedAt),
		Amendments: &primemessages.Amendments{
			Total:          handlers.FmtInt64(0),
			AvailableSince: handlers.FmtInt64(0),
		},
	}

	if move.PPMType != nil {
		payload.PpmType = *move.PPMType
	}

	if moveOrderAmendmentsCount != nil {
		payload.Amendments.Total = handlers.FmtInt64(int64(moveOrderAmendmentsCount.Total))
		payload.Amendments.AvailableSince = handlers.FmtInt64(int64(moveOrderAmendmentsCount.AvailableSinceTotal))
	}

	return payload
}

// ListMoves payload
func ListMoves(moves *models.Moves, moveOrderAmendmentAvailableSinceCounts services.MoveOrderAmendmentAvailableSinceCounts) []*primemessages.ListMove {
	payload := make(primemessages.ListMoves, len(*moves))

	moveOrderAmendmentsFilterCountMap := make(map[uuid.UUID]services.MoveOrderAmendmentAvailableSinceCount, len(*moves))
	for _, info := range moveOrderAmendmentAvailableSinceCounts {
		moveOrderAmendmentsFilterCountMap[info.MoveID] = info
	}

	for i, m := range *moves {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		if value, ok := moveOrderAmendmentsFilterCountMap[m.ID]; ok {
			payload[i] = ListMove(&copyOfM, &value)
		} else {
			payload[i] = ListMove(&copyOfM, nil)
		}
	}

	return payload
}

// Customer payload
func Customer(customer *models.ServiceMember) *primemessages.Customer {
	if customer == nil {
		return nil
	}
	payload := primemessages.Customer{
		FirstName:      swag.StringValue(customer.FirstName),
		LastName:       swag.StringValue(customer.LastName),
		DodID:          swag.StringValue(customer.Edipi),
		Emplid:         swag.StringValue(customer.Emplid),
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
func Order(order *models.Order) *primemessages.Order {
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

	payload := primemessages.Order{
		CustomerID:                   strfmt.UUID(order.ServiceMemberID.String()),
		Customer:                     Customer(&order.ServiceMember),
		DestinationDutyLocation:      destinationDutyLocation,
		DestinationDutyLocationGBLOC: swag.StringValue(order.DestinationGBLOC),
		Entitlement:                  Entitlement(order.Entitlement),
		ID:                           strfmt.UUID(order.ID.String()),
		OriginDutyLocation:           originDutyLocation,
		OriginDutyLocationGBLOC:      swag.StringValue(order.OriginDutyLocationGBLOC),
		OrderNumber:                  order.OrdersNumber,
		LinesOfAccounting:            order.TAC,
		Rank:                         &grade, // Convert prime API "Rank" into our internal tracking of "Grade"
		ETag:                         etag.GenerateEtag(order.UpdatedAt),
		ReportByDate:                 strfmt.Date(order.ReportByDate),
		OrdersType:                   primemessages.OrdersType(order.OrdersType),
	}

	if strings.ToLower(payload.Customer.Branch) == "marines" {
		payload.OriginDutyLocationGBLOC = "USMC"
		payload.DestinationDutyLocationGBLOC = "USMC"
	}

	return &payload
}

// Entitlement payload
func Entitlement(entitlement *models.Entitlement) *primemessages.Entitlements {
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
	var ubAllowance int64
	if entitlement.UBAllowance != nil {
		ubAllowance = int64(*entitlement.UBAllowance)
	}
	return &primemessages.Entitlements{
		ID:                             strfmt.UUID(entitlement.ID.String()),
		AuthorizedWeight:               authorizedWeight,
		UnaccompaniedBaggageAllowance:  &ubAllowance,
		DependentsAuthorized:           entitlement.DependentsAuthorized,
		GunSafe:                        entitlement.GunSafe,
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
func DutyLocation(dutyLocation *models.DutyLocation) *primemessages.DutyLocation {
	if dutyLocation == nil {
		return nil
	}
	address := Address(&dutyLocation.Address)
	payload := primemessages.DutyLocation{
		Address:   address,
		AddressID: address.ID,
		ID:        strfmt.UUID(dutyLocation.ID.String()),
		Name:      dutyLocation.Name,
	}
	return &payload
}

// Country payload
func Country(country *models.Country) *string {
	if country == nil {
		return nil
	}
	return &country.Country
}

// Address payload
func Address(address *models.Address) *primemessages.Address {
	if address == nil {
		return nil
	}
	return &primemessages.Address{
		ID:             strfmt.UUID(address.ID.String()),
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
		City:           &address.City,
		State:          &address.State,
		PostalCode:     &address.PostalCode,
		Country:        Country(address.Country),
		County:         &address.County,
		ETag:           etag.GenerateEtag(address.UpdatedAt),
	}
}

// StorageFacility payload
func StorageFacility(storage *models.StorageFacility) *primemessages.StorageFacility {
	if storage == nil {
		return nil
	}

	return &primemessages.StorageFacility{
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
func MTOAgent(mtoAgent *models.MTOAgent) *primemessages.MTOAgent {
	if mtoAgent == nil {
		return nil
	}

	return &primemessages.MTOAgent{
		AgentType:     primemessages.MTOAgentType(mtoAgent.MTOAgentType),
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
func MTOAgents(mtoAgents *models.MTOAgents) *primemessages.MTOAgents {
	if mtoAgents == nil {
		return nil
	}

	agents := make(primemessages.MTOAgents, len(*mtoAgents))

	for i, m := range *mtoAgents {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		agents[i] = MTOAgent(&copyOfM)
	}

	return &agents
}

func ProofOfServiceDoc(proofOfServiceDoc models.ProofOfServiceDoc) *primemessages.ProofOfServiceDoc {
	uploads := make([]*primemessages.UploadWithOmissions, len(proofOfServiceDoc.PrimeUploads))
	if len(proofOfServiceDoc.PrimeUploads) > 0 {
		for i, primeUpload := range proofOfServiceDoc.PrimeUploads {
			uploads[i] = basicUpload(&primeUpload.Upload) //#nosec G601
		}
	}

	return &primemessages.ProofOfServiceDoc{
		Uploads: uploads,
	}
}

// PaymentRequest payload
func PaymentRequest(paymentRequest *models.PaymentRequest) *primemessages.PaymentRequest {
	if paymentRequest == nil {
		return nil
	}

	serviceDocs := make(primemessages.ProofOfServiceDocs, len(paymentRequest.ProofOfServiceDocs))

	if len(paymentRequest.ProofOfServiceDocs) > 0 {
		for i, proofOfService := range paymentRequest.ProofOfServiceDocs {
			serviceDocs[i] = ProofOfServiceDoc(proofOfService)
		}
	}

	paymentServiceItems := PaymentServiceItems(&paymentRequest.PaymentServiceItems)
	return &primemessages.PaymentRequest{
		ID:                              strfmt.UUID(paymentRequest.ID.String()),
		IsFinal:                         &paymentRequest.IsFinal,
		MoveTaskOrderID:                 strfmt.UUID(paymentRequest.MoveTaskOrderID.String()),
		PaymentRequestNumber:            paymentRequest.PaymentRequestNumber,
		RecalculationOfPaymentRequestID: handlers.FmtUUIDPtr(paymentRequest.RecalculationOfPaymentRequestID),
		RejectionReason:                 paymentRequest.RejectionReason,
		Status:                          primemessages.PaymentRequestStatus(paymentRequest.Status),
		PaymentServiceItems:             *paymentServiceItems,
		ProofOfServiceDocs:              serviceDocs,
		ETag:                            etag.GenerateEtag(paymentRequest.UpdatedAt),
	}
}

// PaymentRequests payload
func PaymentRequests(paymentRequests *models.PaymentRequests) *primemessages.PaymentRequests {
	if paymentRequests == nil {
		return nil
	}

	payload := make(primemessages.PaymentRequests, len(*paymentRequests))

	for i, p := range *paymentRequests {
		copyOfP := p // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentRequest(&copyOfP)
	}
	return &payload
}

// PaymentServiceItem payload
func PaymentServiceItem(paymentServiceItem *models.PaymentServiceItem) *primemessages.PaymentServiceItem {
	if paymentServiceItem == nil {
		return nil
	}

	paymentServiceItemParams := PaymentServiceItemParams(&paymentServiceItem.PaymentServiceItemParams)

	payload := &primemessages.PaymentServiceItem{
		ID:                       strfmt.UUID(paymentServiceItem.ID.String()),
		PaymentRequestID:         strfmt.UUID(paymentServiceItem.PaymentRequestID.String()),
		MtoServiceItemID:         strfmt.UUID(paymentServiceItem.MTOServiceItemID.String()),
		Status:                   primemessages.PaymentServiceItemStatus(paymentServiceItem.Status),
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
func PaymentServiceItems(paymentServiceItems *models.PaymentServiceItems) *primemessages.PaymentServiceItems {
	if paymentServiceItems == nil {
		return nil
	}

	payload := make(primemessages.PaymentServiceItems, len(*paymentServiceItems))

	for i, p := range *paymentServiceItems {
		copyOfP := p // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentServiceItem(&copyOfP)
	}
	return &payload
}

// PaymentServiceItemParam payload
func PaymentServiceItemParam(paymentServiceItemParam *models.PaymentServiceItemParam) *primemessages.PaymentServiceItemParam {
	if paymentServiceItemParam == nil {
		return nil
	}

	return &primemessages.PaymentServiceItemParam{
		ID:                   strfmt.UUID(paymentServiceItemParam.ID.String()),
		PaymentServiceItemID: strfmt.UUID(paymentServiceItemParam.PaymentServiceItemID.String()),
		Key:                  primemessages.ServiceItemParamName(paymentServiceItemParam.ServiceItemParamKey.Key),
		Value:                paymentServiceItemParam.Value,
		Type:                 primemessages.ServiceItemParamType(paymentServiceItemParam.ServiceItemParamKey.Type),
		Origin:               primemessages.ServiceItemParamOrigin(paymentServiceItemParam.ServiceItemParamKey.Origin),
		ETag:                 etag.GenerateEtag(paymentServiceItemParam.UpdatedAt),
	}
}

// PaymentServiceItemParams payload
func PaymentServiceItemParams(paymentServiceItemParams *models.PaymentServiceItemParams) *primemessages.PaymentServiceItemParams {
	if paymentServiceItemParams == nil {
		return nil
	}

	payload := make(primemessages.PaymentServiceItemParams, len(*paymentServiceItemParams))

	for i, p := range *paymentServiceItemParams {
		copyOfP := p // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentServiceItemParam(&copyOfP)
	}
	return &payload
}

//nolint:gosec //G601
func ServiceRequestDocument(serviceRequestDocument models.ServiceRequestDocument) *primemessages.ServiceRequestDocument {
	uploads := make([]*primemessages.UploadWithOmissions, len(serviceRequestDocument.ServiceRequestDocumentUploads))
	if len(serviceRequestDocument.ServiceRequestDocumentUploads) > 0 {
		for i, proofOfServiceDocumentUpload := range serviceRequestDocument.ServiceRequestDocumentUploads {
			uploads[i] = basicUpload(&proofOfServiceDocumentUpload.Upload)
		}
	}

	return &primemessages.ServiceRequestDocument{
		Uploads: uploads,
	}
}

// PPMShipment payload
func PPMShipment(ppmShipment *models.PPMShipment) *primemessages.PPMShipment {
	if ppmShipment == nil || ppmShipment.ID.IsNil() {
		return nil
	}

	payloadPPMShipment := &primemessages.PPMShipment{
		ID:                           *handlers.FmtUUID(ppmShipment.ID),
		ShipmentID:                   *handlers.FmtUUID(ppmShipment.ShipmentID),
		CreatedAt:                    strfmt.DateTime(ppmShipment.CreatedAt),
		UpdatedAt:                    strfmt.DateTime(ppmShipment.UpdatedAt),
		Status:                       primemessages.PPMShipmentStatus(ppmShipment.Status),
		ExpectedDepartureDate:        handlers.FmtDate(ppmShipment.ExpectedDepartureDate),
		ActualMoveDate:               handlers.FmtDatePtr(ppmShipment.ActualMoveDate),
		SubmittedAt:                  handlers.FmtDateTimePtr(ppmShipment.SubmittedAt),
		ReviewedAt:                   handlers.FmtDateTimePtr(ppmShipment.ReviewedAt),
		ApprovedAt:                   handlers.FmtDateTimePtr(ppmShipment.ApprovedAt),
		ActualPickupPostalCode:       ppmShipment.ActualPickupPostalCode,
		ActualDestinationPostalCode:  ppmShipment.ActualDestinationPostalCode,
		SitExpected:                  ppmShipment.SITExpected,
		SitEstimatedWeight:           handlers.FmtPoundPtr(ppmShipment.SITEstimatedWeight),
		SitEstimatedEntryDate:        handlers.FmtDatePtr(ppmShipment.SITEstimatedEntryDate),
		SitEstimatedDepartureDate:    handlers.FmtDatePtr(ppmShipment.SITEstimatedDepartureDate),
		SitEstimatedCost:             handlers.FmtCost(ppmShipment.SITEstimatedCost),
		EstimatedWeight:              handlers.FmtPoundPtr(ppmShipment.EstimatedWeight),
		EstimatedIncentive:           handlers.FmtCost(ppmShipment.EstimatedIncentive),
		HasProGear:                   ppmShipment.HasProGear,
		ProGearWeight:                handlers.FmtPoundPtr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:          handlers.FmtPoundPtr(ppmShipment.SpouseProGearWeight),
		HasRequestedAdvance:          ppmShipment.HasRequestedAdvance,
		AdvanceAmountRequested:       handlers.FmtCost(ppmShipment.AdvanceAmountRequested),
		HasReceivedAdvance:           ppmShipment.HasReceivedAdvance,
		AdvanceAmountReceived:        handlers.FmtCost(ppmShipment.AdvanceAmountReceived),
		IsActualExpenseReimbursement: ppmShipment.IsActualExpenseReimbursement,
		ETag:                         etag.GenerateEtag(ppmShipment.UpdatedAt),
	}

	if ppmShipment.SITLocation != nil {
		sitLocation := primemessages.SITLocationType(*ppmShipment.SITLocation)
		payloadPPMShipment.SitLocation = &sitLocation
	}

	if ppmShipment.IsActualExpenseReimbursement != nil {
		payloadPPMShipment.IsActualExpenseReimbursement = ppmShipment.IsActualExpenseReimbursement
	}

	return payloadPPMShipment
}

// MarketCode payload
func MarketCode(marketCode *models.MarketCode) string {
	if marketCode == nil {
		return "" // Or a default string value
	}
	return string(*marketCode)
}

func MTOShipmentWithoutServiceItems(mtoShipment *models.MTOShipment) *primemessages.MTOShipmentWithoutServiceItems {
	payload := &primemessages.MTOShipmentWithoutServiceItems{
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
		ShipmentType:                     primemessages.MTOShipmentType(mtoShipment.ShipmentType),
		CustomerRemarks:                  mtoShipment.CustomerRemarks,
		CounselorRemarks:                 mtoShipment.CounselorRemarks,
		Status:                           string(mtoShipment.Status),
		Diversion:                        bool(mtoShipment.Diversion),
		DiversionReason:                  mtoShipment.DiversionReason,
		DeliveryAddressUpdate:            ShipmentAddressUpdate(mtoShipment.DeliveryAddressUpdate),
		CreatedAt:                        strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                        strfmt.DateTime(mtoShipment.UpdatedAt),
		PpmShipment:                      PPMShipment(mtoShipment.PPMShipment),
		ETag:                             etag.GenerateEtag(mtoShipment.UpdatedAt),
		OriginSitAuthEndDate:             (*strfmt.Date)(mtoShipment.OriginSITAuthEndDate),
		DestinationSitAuthEndDate:        (*strfmt.Date)(mtoShipment.DestinationSITAuthEndDate),
		MarketCode:                       MarketCode(&mtoShipment.MarketCode),
	}

	// Set up address payloads
	if mtoShipment.PickupAddress != nil {
		payload.PickupAddress.Address = *Address(mtoShipment.PickupAddress)
	}
	if mtoShipment.DestinationAddress != nil {
		payload.DestinationAddress.Address = *Address(mtoShipment.DestinationAddress)
	}
	if mtoShipment.DestinationType != nil {
		destinationType := primemessages.DestinationType(*mtoShipment.DestinationType)
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

	return payload
}

func MTOShipmentsWithoutServiceItems(mtoShipments *models.MTOShipments) *primemessages.MTOShipmentsWithoutServiceObjects {
	payload := make(primemessages.MTOShipmentsWithoutServiceObjects, len(*mtoShipments))

	for i, m := range *mtoShipments {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOShipmentWithoutServiceItems(&copyOfM)
	}
	return &payload
}

// MTOShipment converts MTOShipment model to payload
func MTOShipment(mtoShipment *models.MTOShipment) *primemessages.MTOShipment {
	payload := &primemessages.MTOShipment{
		MTOShipmentWithoutServiceItems: *MTOShipmentWithoutServiceItems(mtoShipment),
	}

	if mtoShipment.MTOServiceItems != nil {
		payload.SetMtoServiceItems(*MTOServiceItems(&mtoShipment.MTOServiceItems))
	} else {
		payload.SetMtoServiceItems([]primemessages.MTOServiceItem{})
	}

	return payload
}

// MTOServiceItem payload
func MTOServiceItem(mtoServiceItem *models.MTOServiceItem) primemessages.MTOServiceItem {
	var payload primemessages.MTOServiceItem
	// here we determine which payload model to use based on the re service code
	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDOPSIT, models.ReServiceCodeDOSFSC:
		var sitDepartureDate time.Time
		if mtoServiceItem.SITDepartureDate != nil {
			sitDepartureDate = *mtoServiceItem.SITDepartureDate
		}
		payload = &primemessages.MTOServiceItemOriginSIT{
			ReServiceCode:                   handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Reason:                          mtoServiceItem.Reason,
			SitDepartureDate:                handlers.FmtDate(sitDepartureDate),
			SitEntryDate:                    handlers.FmtDatePtr(mtoServiceItem.SITEntryDate),
			SitPostalCode:                   mtoServiceItem.SITPostalCode,
			SitHHGActualOrigin:              Address(mtoServiceItem.SITOriginHHGActualAddress),
			SitHHGOriginalOrigin:            Address(mtoServiceItem.SITOriginHHGOriginalAddress),
			RequestApprovalsRequestedStatus: *mtoServiceItem.RequestedApprovalsRequestedStatus,
			SitCustomerContacted:            handlers.FmtDatePtr(mtoServiceItem.SITCustomerContacted),
			SitRequestedDelivery:            handlers.FmtDatePtr(mtoServiceItem.SITRequestedDelivery),
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

		payload = &primemessages.MTOServiceItemDestSIT{
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
			SitCustomerContacted:        handlers.FmtDatePtr(mtoServiceItem.SITCustomerContacted),
			SitRequestedDelivery:        handlers.FmtDatePtr(mtoServiceItem.SITRequestedDelivery),
		}

	case models.ReServiceCodeDCRT, models.ReServiceCodeDUCRT:
		item := GetDimension(mtoServiceItem.Dimensions, models.DimensionTypeItem)
		crate := GetDimension(mtoServiceItem.Dimensions, models.DimensionTypeCrate)
		cratingSI := primemessages.MTOServiceItemDomesticCrating{
			ReServiceCode:   handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Description:     mtoServiceItem.Description,
			Reason:          mtoServiceItem.Reason,
			StandaloneCrate: mtoServiceItem.StandaloneCrate,
		}
		cratingSI.Item.MTOServiceItemDimension = primemessages.MTOServiceItemDimension{
			ID:     strfmt.UUID(item.ID.String()),
			Height: item.Height.Int32Ptr(),
			Length: item.Length.Int32Ptr(),
			Width:  item.Width.Int32Ptr(),
		}
		cratingSI.Crate.MTOServiceItemDimension = primemessages.MTOServiceItemDimension{
			ID:     strfmt.UUID(crate.ID.String()),
			Height: crate.Height.Int32Ptr(),
			Length: crate.Length.Int32Ptr(),
			Width:  crate.Width.Int32Ptr(),
		}
		payload = &cratingSI
	case models.ReServiceCodeDDSHUT, models.ReServiceCodeDOSHUT:
		payload = &primemessages.MTOServiceItemShuttle{
			ReServiceCode:   handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Reason:          mtoServiceItem.Reason,
			EstimatedWeight: handlers.FmtPoundPtr(mtoServiceItem.EstimatedWeight),
			ActualWeight:    handlers.FmtPoundPtr(mtoServiceItem.ActualWeight),
		}
	default:
		// otherwise, basic service item
		payload = &primemessages.MTOServiceItemBasic{
			ReServiceCode: primemessages.NewReServiceCode(primemessages.ReServiceCode(mtoServiceItem.ReService.Code)),
		}
	}

	// set all relevant fields that apply to all service items
	var shipmentIDStr string
	if mtoServiceItem.MTOShipmentID != nil {
		shipmentIDStr = mtoServiceItem.MTOShipmentID.String()
	}

	serviceRequestDocuments := make(primemessages.ServiceRequestDocuments, len(mtoServiceItem.ServiceRequestDocuments))

	if len(mtoServiceItem.ServiceRequestDocuments) > 0 {
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
	payload.SetStatus(primemessages.MTOServiceItemStatus(mtoServiceItem.Status))
	payload.SetRejectionReason(mtoServiceItem.RejectionReason)
	payload.SetETag(etag.GenerateEtag(mtoServiceItem.UpdatedAt))
	payload.SetServiceRequestDocuments(serviceRequestDocuments)
	return payload
}

// MTOServiceItems payload
func MTOServiceItems(mtoServiceItems *models.MTOServiceItems) *[]primemessages.MTOServiceItem {
	payload := []primemessages.MTOServiceItem{}

	for _, p := range *mtoServiceItems {
		copyOfP := p // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload = append(payload, MTOServiceItem(&copyOfP))
	}
	return &payload
}

// Reweigh returns the reweigh payload
func Reweigh(reweigh *models.Reweigh) *primemessages.Reweigh {
	if reweigh == nil || reweigh.ID == uuid.Nil {
		return nil
	}

	payload := &primemessages.Reweigh{
		ID:                     strfmt.UUID(reweigh.ID.String()),
		ShipmentID:             strfmt.UUID(reweigh.ShipmentID.String()),
		RequestedAt:            strfmt.DateTime(reweigh.RequestedAt),
		RequestedBy:            primemessages.ReweighRequester(reweigh.RequestedBy),
		CreatedAt:              strfmt.DateTime(reweigh.CreatedAt),
		UpdatedAt:              strfmt.DateTime(reweigh.UpdatedAt),
		ETag:                   etag.GenerateEtag(reweigh.UpdatedAt),
		Weight:                 handlers.FmtPoundPtr(reweigh.Weight),
		VerificationReason:     handlers.FmtStringPtr(reweigh.VerificationReason),
		VerificationProvidedAt: handlers.FmtDateTimePtr(reweigh.VerificationProvidedAt),
	}

	return payload
}

// ExcessWeightRecord returns the fields on the move related to excess weights,
// and returns the uploaded document set as the ExcessWeightUpload on the move.
func ExcessWeightRecord(appCtx appcontext.AppContext, storer storage.FileStorer, move *models.Move) *primemessages.ExcessWeightRecord {
	if move == nil || move.ID == uuid.Nil {
		return nil
	}

	payload := &primemessages.ExcessWeightRecord{
		MoveID:                         handlers.FmtUUIDPtr(&move.ID),
		MoveExcessWeightQualifiedAt:    handlers.FmtDateTimePtr(move.ExcessWeightQualifiedAt),
		MoveExcessWeightAcknowledgedAt: handlers.FmtDateTimePtr(move.ExcessWeightAcknowledgedAt),
	}

	upload := Upload(appCtx, storer, move.ExcessWeightUpload)
	if upload != nil {
		payload.UploadWithOmissions = *upload
	}

	return payload
}

// Upload returns the data for an uploaded file.
func Upload(appCtx appcontext.AppContext, storer storage.FileStorer, upload *models.Upload) *primemessages.UploadWithOmissions {
	if upload == nil || upload.ID == uuid.Nil {
		return nil
	}

	payload := &primemessages.UploadWithOmissions{
		ID:          strfmt.UUID(upload.ID.String()),
		Bytes:       &upload.Bytes,
		ContentType: &upload.ContentType,
		Filename:    &upload.Filename,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
	}

	url, err := storer.PresignedURL(upload.StorageKey, upload.ContentType, upload.Filename)
	if err == nil {
		payload.URL = *handlers.FmtURI(url)
	} else {
		appCtx.Logger().Error("primeapi error with getting url for Upload payload", zap.Error(err))
	}

	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || tags == nil {
		payload.Status = "PROCESSING"
	} else {
		status, ok := tags["av-status"]
		if !ok {
			status = "PROCESSING"
		}
		payload.Status = status
	}

	return payload
}

func basicUpload(upload *models.Upload) *primemessages.UploadWithOmissions {
	if upload == nil || upload.ID == uuid.Nil {
		return nil
	}

	payload := &primemessages.UploadWithOmissions{
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
func SITDurationUpdate(sitDurationUpdate *models.SITDurationUpdate) *primemessages.SITExtension {
	if sitDurationUpdate == nil {
		return nil
	}
	payload := &primemessages.SITExtension{
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

// ShipmentAddressUpdate payload
func ShipmentAddressUpdate(shipmentAddressUpdate *models.ShipmentAddressUpdate) *primemessages.ShipmentAddressUpdate {
	if shipmentAddressUpdate == nil || shipmentAddressUpdate.ID.IsNil() {
		return nil
	}

	payload := &primemessages.ShipmentAddressUpdate{
		ID:                strfmt.UUID(shipmentAddressUpdate.ID.String()),
		ShipmentID:        strfmt.UUID(shipmentAddressUpdate.ShipmentID.String()),
		NewAddress:        Address(&shipmentAddressUpdate.NewAddress),
		OriginalAddress:   Address(&shipmentAddressUpdate.OriginalAddress),
		ContractorRemarks: shipmentAddressUpdate.ContractorRemarks,
		OfficeRemarks:     shipmentAddressUpdate.OfficeRemarks,
		Status:            primemessages.ShipmentAddressUpdateStatus(shipmentAddressUpdate.Status),
	}

	return payload
}

// SITDurationUpdates payload
func SITDurationUpdates(sitDurationUpdates *models.SITDurationUpdates) *primemessages.SITExtensions {
	if sitDurationUpdates == nil {
		return nil
	}

	payload := make(primemessages.SITExtensions, len(*sitDurationUpdates))

	for i, m := range *sitDurationUpdates {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = SITDurationUpdate(&copyOfM)
	}

	return &payload
}

// InternalServerError describes errors in a standard structure to be returned in the payload.
// If detail is nil, string defaults to "An internal server error has occurred."
func InternalServerError(detail *string, traceID uuid.UUID) *primemessages.Error {
	payload := primemessages.Error{
		Title:    handlers.FmtString(handlers.InternalServerErrMessage),
		Detail:   handlers.FmtString(handlers.InternalServerErrDetail),
		Instance: strfmt.UUID(traceID.String()),
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}

// NotImplementedError describes errors for endpoints and functions that haven't been fully developed yet.
// If detail is nil, string defaults to "This feature is in development"
func NotImplementedError(detail *string, traceID uuid.UUID) *primemessages.Error {
	payload := primemessages.Error{
		Title:    handlers.FmtString(handlers.NotImplementedErrMessage),
		Detail:   handlers.FmtString(handlers.NotImplementedErrDetail),
		Instance: strfmt.UUID(traceID.String()),
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}

// ValidationError describes validation errors from the model or properties
func ValidationError(detail string, instance uuid.UUID, validationErrors *validate.Errors) *primemessages.ValidationError {
	payload := &primemessages.ValidationError{
		ClientError: *ClientError(handlers.ValidationErrMessage, detail, instance),
	}
	if validationErrors != nil {
		payload.InvalidFields = handlers.NewValidationErrorListResponse(validationErrors).Errors
	}
	return payload
}

// ClientError describes errors in a standard structure to be returned in the payload
func ClientError(title string, detail string, instance uuid.UUID) *primemessages.ClientError {
	return &primemessages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
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
