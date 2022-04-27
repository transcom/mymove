package payloads

import (
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/transcom/mymove/pkg/storage"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// MoveTaskOrder payload
func MoveTaskOrder(moveTaskOrder *models.Move) *primemessages.MoveTaskOrder {
	if moveTaskOrder == nil {
		return nil
	}
	paymentRequests := PaymentRequests(&moveTaskOrder.PaymentRequests)
	mtoServiceItems := MTOServiceItems(&moveTaskOrder.MTOServiceItems)
	mtoShipments := MTOShipments(&moveTaskOrder.MTOShipments)

	payload := &primemessages.MoveTaskOrder{
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
		UpdatedAt:                  strfmt.DateTime(moveTaskOrder.UpdatedAt),
		ETag:                       etag.GenerateEtag(moveTaskOrder.UpdatedAt),
	}

	if moveTaskOrder.PPMEstimatedWeight != nil {
		payload.PpmEstimatedWeight = int64(*moveTaskOrder.PPMEstimatedWeight)
	}

	if moveTaskOrder.PPMType != nil {
		payload.PpmType = *moveTaskOrder.PPMType
	}

	// mto service item references a polymorphic type which auto-generates an interface and getters and setters
	payload.SetMtoServiceItems(*mtoServiceItems)

	return payload
}

// ListMove payload
func ListMove(move *models.Move) *primemessages.ListMove {
	if move == nil {
		return nil
	}
	payload := &primemessages.ListMove{
		ID:                 strfmt.UUID(move.ID.String()),
		MoveCode:           move.Locator,
		CreatedAt:          strfmt.DateTime(move.CreatedAt),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(move.AvailableToPrimeAt),
		OrderID:            strfmt.UUID(move.OrdersID.String()),
		ReferenceID:        *move.ReferenceID,
		UpdatedAt:          strfmt.DateTime(move.UpdatedAt),
		ETag:               etag.GenerateEtag(move.UpdatedAt),
	}

	if move.PPMEstimatedWeight != nil {
		payload.PpmEstimatedWeight = int64(*move.PPMEstimatedWeight)
	}

	if move.PPMType != nil {
		payload.PpmType = *move.PPMType
	}

	return payload
}

// ListMoves payload
func ListMoves(moves *models.Moves) []*primemessages.ListMove {
	payload := make(primemessages.ListMoves, len(*moves))

	for i, m := range *moves {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = ListMove(&copyOfM)
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
		order.Entitlement.SetWeightAllotment(*order.Grade)
	}

	payload := primemessages.Order{
		CustomerID:              strfmt.UUID(order.ServiceMemberID.String()),
		Customer:                Customer(&order.ServiceMember),
		DestinationDutyLocation: destinationDutyLocation,
		Entitlement:             Entitlement(order.Entitlement),
		ID:                      strfmt.UUID(order.ID.String()),
		OriginDutyLocation:      originDutyLocation,
		OrderNumber:             order.OrdersNumber,
		LinesOfAccounting:       order.TAC,
		Rank:                    order.Grade,
		ETag:                    etag.GenerateEtag(order.UpdatedAt),
		ReportByDate:            strfmt.Date(order.ReportByDate),
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
	return &primemessages.Entitlements{
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
		Country:        address.Country,
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

// PaymentRequest payload
func PaymentRequest(paymentRequest *models.PaymentRequest) *primemessages.PaymentRequest {
	if paymentRequest == nil {
		return nil
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
		payload.PriceCents = swag.Int64(int64(*paymentServiceItem.PriceCents))
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

// PPMShipment payload
func PPMShipment(ppmShipment *models.PPMShipment) *primemessages.PPMShipment {
	if ppmShipment == nil || ppmShipment.ID.IsNil() {
		return nil
	}

	payloadPPMShipment := &primemessages.PPMShipment{
		ID:                             *handlers.FmtUUID(ppmShipment.ID),
		ShipmentID:                     *handlers.FmtUUID(ppmShipment.ShipmentID),
		CreatedAt:                      strfmt.DateTime(ppmShipment.CreatedAt),
		UpdatedAt:                      strfmt.DateTime(ppmShipment.UpdatedAt),
		Status:                         primemessages.PPMShipmentStatus(ppmShipment.Status),
		ExpectedDepartureDate:          handlers.FmtDate(ppmShipment.ExpectedDepartureDate),
		ActualMoveDate:                 handlers.FmtDatePtr(ppmShipment.ActualMoveDate),
		SubmittedAt:                    handlers.FmtDateTimePtr(ppmShipment.SubmittedAt),
		ReviewedAt:                     handlers.FmtDateTimePtr(ppmShipment.ReviewedAt),
		ApprovedAt:                     handlers.FmtDateTimePtr(ppmShipment.ApprovedAt),
		PickupPostalCode:               &ppmShipment.PickupPostalCode,
		SecondaryPickupPostalCode:      ppmShipment.SecondaryPickupPostalCode,
		DestinationPostalCode:          &ppmShipment.DestinationPostalCode,
		SecondaryDestinationPostalCode: ppmShipment.SecondaryDestinationPostalCode,
		SitExpected:                    ppmShipment.SitExpected,
		EstimatedWeight:                handlers.FmtPoundPtr(ppmShipment.EstimatedWeight),
		EstimatedIncentive:             handlers.FmtCost(ppmShipment.EstimatedIncentive),
		NetWeight:                      handlers.FmtPoundPtr(ppmShipment.NetWeight),
		HasProGear:                     ppmShipment.HasProGear,
		ProGearWeight:                  handlers.FmtPoundPtr(ppmShipment.ProGearWeight),
		SpouseProGearWeight:            handlers.FmtPoundPtr(ppmShipment.SpouseProGearWeight),
		Advance:                        handlers.FmtCost(ppmShipment.Advance),
		AdvanceRequested:               ppmShipment.AdvanceRequested,
		ETag:                           etag.GenerateEtag(ppmShipment.UpdatedAt),
	}

	return payloadPPMShipment
}

// MTOShipment converts MTOShipment model to payload
func MTOShipment(mtoShipment *models.MTOShipment) *primemessages.MTOShipment {
	payload := &primemessages.MTOShipment{
		ID:                               strfmt.UUID(mtoShipment.ID.String()),
		ActualPickupDate:                 handlers.FmtDatePtr(mtoShipment.ActualPickupDate),
		ApprovedDate:                     handlers.FmtDatePtr(mtoShipment.ApprovedDate),
		FirstAvailableDeliveryDate:       handlers.FmtDatePtr(mtoShipment.FirstAvailableDeliveryDate),
		PrimeEstimatedWeightRecordedDate: handlers.FmtDatePtr(mtoShipment.PrimeEstimatedWeightRecordedDate),
		RequestedPickupDate:              handlers.FmtDatePtr(mtoShipment.RequestedPickupDate),
		RequiredDeliveryDate:             handlers.FmtDatePtr(mtoShipment.RequiredDeliveryDate),
		ScheduledPickupDate:              handlers.FmtDatePtr(mtoShipment.ScheduledPickupDate),
		Agents:                           *MTOAgents(&mtoShipment.MTOAgents),
		SitExtensions:                    *SITExtensions(&mtoShipment.SITExtensions),
		Reweigh:                          Reweigh(mtoShipment.Reweigh),
		MoveTaskOrderID:                  strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:                     primemessages.MTOShipmentType(mtoShipment.ShipmentType),
		CustomerRemarks:                  mtoShipment.CustomerRemarks,
		CounselorRemarks:                 mtoShipment.CounselorRemarks,
		Status:                           string(mtoShipment.Status),
		Diversion:                        bool(mtoShipment.Diversion),
		CreatedAt:                        strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                        strfmt.DateTime(mtoShipment.UpdatedAt),
		PpmShipment:                      PPMShipment(mtoShipment.PPMShipment),
		ETag:                             etag.GenerateEtag(mtoShipment.UpdatedAt),
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

	if mtoShipment.MTOServiceItems != nil {
		payload.SetMtoServiceItems(*MTOServiceItems(&mtoShipment.MTOServiceItems))
	} else {
		payload.SetMtoServiceItems([]primemessages.MTOServiceItem{})
	}

	if mtoShipment.PrimeEstimatedWeight != nil {
		payload.PrimeEstimatedWeight = int64(*mtoShipment.PrimeEstimatedWeight)
	}

	if mtoShipment.PrimeActualWeight != nil {
		payload.PrimeActualWeight = int64(*mtoShipment.PrimeActualWeight)
	}

	if mtoShipment.NTSRecordedWeight != nil {
		payload.NtsRecordedWeight = handlers.FmtInt64(mtoShipment.NTSRecordedWeight.Int64())
	}

	return payload
}

// MTOShipments converts an array of MTOShipment models to a payload
func MTOShipments(mtoShipments *models.MTOShipments) *primemessages.MTOShipments {
	payload := make(primemessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOShipment(&copyOfM)
	}
	return &payload
}

// MTOServiceItem payload
func MTOServiceItem(mtoServiceItem *models.MTOServiceItem) primemessages.MTOServiceItem {
	var payload primemessages.MTOServiceItem
	// here we determine which payload model to use based on the re service code
	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDOPSIT:
		var sitDepartureDate time.Time
		if mtoServiceItem.SITDepartureDate != nil {
			sitDepartureDate = *mtoServiceItem.SITDepartureDate
		}
		payload = &primemessages.MTOServiceItemOriginSIT{
			ReServiceCode:      handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Reason:             mtoServiceItem.Reason,
			SitDepartureDate:   handlers.FmtDate(sitDepartureDate),
			SitEntryDate:       handlers.FmtDatePtr(mtoServiceItem.SITEntryDate),
			SitPostalCode:      mtoServiceItem.SITPostalCode,
			SitHHGActualOrigin: Address(mtoServiceItem.SITOriginHHGActualAddress),
		}
	case models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT, models.ReServiceCodeDDDSIT:
		var sitDepartureDate, firstAvailableDeliveryDate1, firstAvailableDeliveryDate2 time.Time
		var timeMilitary1, timeMilitary2 *string

		if mtoServiceItem.SITDepartureDate != nil {
			sitDepartureDate = *mtoServiceItem.SITDepartureDate
		}

		firstContact := GetCustomerContact(mtoServiceItem.CustomerContacts, models.CustomerContactTypeFirst)
		secondContact := GetCustomerContact(mtoServiceItem.CustomerContacts, models.CustomerContactTypeSecond)
		timeMilitary1 = &firstContact.TimeMilitary
		timeMilitary2 = &secondContact.TimeMilitary

		if !firstContact.FirstAvailableDeliveryDate.IsZero() {
			firstAvailableDeliveryDate1 = firstContact.FirstAvailableDeliveryDate
		}

		if !secondContact.FirstAvailableDeliveryDate.IsZero() {
			firstAvailableDeliveryDate2 = firstContact.FirstAvailableDeliveryDate
		}

		payload = &primemessages.MTOServiceItemDestSIT{
			ReServiceCode:               handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			TimeMilitary1:               handlers.FmtStringPtrNonEmpty(timeMilitary1),
			FirstAvailableDeliveryDate1: handlers.FmtDate(firstAvailableDeliveryDate1),
			TimeMilitary2:               handlers.FmtStringPtrNonEmpty(timeMilitary2),
			FirstAvailableDeliveryDate2: handlers.FmtDate(firstAvailableDeliveryDate2),
			SitDepartureDate:            handlers.FmtDate(sitDepartureDate),
			SitEntryDate:                handlers.FmtDatePtr(mtoServiceItem.SITEntryDate),
			SitDestinationFinalAddress:  Address(mtoServiceItem.SITDestinationFinalAddress),
		}

	case models.ReServiceCodeDCRT, models.ReServiceCodeDUCRT:
		item := GetDimension(mtoServiceItem.Dimensions, models.DimensionTypeItem)
		crate := GetDimension(mtoServiceItem.Dimensions, models.DimensionTypeCrate)
		cratingSI := primemessages.MTOServiceItemDomesticCrating{
			ReServiceCode: handlers.FmtString(string(mtoServiceItem.ReService.Code)),
			Description:   mtoServiceItem.Description,
			Reason:        mtoServiceItem.Reason,
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

	one := mtoServiceItem.ID.String()
	two := strfmt.UUID(one)
	payload.SetID(two)
	payload.SetMoveTaskOrderID(handlers.FmtUUID(mtoServiceItem.MoveTaskOrderID))
	payload.SetMtoShipmentID(strfmt.UUID(shipmentIDStr))
	payload.SetReServiceName(mtoServiceItem.ReService.Name)
	payload.SetStatus(primemessages.MTOServiceItemStatus(mtoServiceItem.Status))
	payload.SetETag(etag.GenerateEtag(mtoServiceItem.UpdatedAt))
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
		payload.Upload = *upload
	}

	return payload
}

// Upload returns the data for an uploaded file.
func Upload(appCtx appcontext.AppContext, storer storage.FileStorer, upload *models.Upload) *primemessages.Upload {
	if upload == nil || upload.ID == uuid.Nil {
		return nil
	}

	payload := &primemessages.Upload{
		ID:          strfmt.UUID(upload.ID.String()),
		Bytes:       &upload.Bytes,
		ContentType: &upload.ContentType,
		Filename:    &upload.Filename,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
	}

	url, err := storer.PresignedURL(upload.StorageKey, upload.ContentType)
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

// SITExtension payload
func SITExtension(sitExtension *models.SITExtension) *primemessages.SITExtension {
	if sitExtension == nil {
		return nil
	}
	payload := &primemessages.SITExtension{
		ID:                strfmt.UUID(sitExtension.ID.String()),
		ETag:              etag.GenerateEtag(sitExtension.UpdatedAt),
		MtoShipmentID:     strfmt.UUID(sitExtension.MTOShipmentID.String()),
		RequestReason:     string(sitExtension.RequestReason),
		RequestedDays:     int64(sitExtension.RequestedDays),
		Status:            string(sitExtension.Status),
		CreatedAt:         strfmt.DateTime(sitExtension.CreatedAt),
		UpdatedAt:         strfmt.DateTime(sitExtension.UpdatedAt),
		ApprovedDays:      handlers.FmtIntPtrToInt64(sitExtension.ApprovedDays),
		ContractorRemarks: handlers.FmtStringPtr(sitExtension.ContractorRemarks),
		DecisionDate:      handlers.FmtDateTimePtr(sitExtension.DecisionDate),
		OfficeRemarks:     handlers.FmtStringPtr(sitExtension.OfficeRemarks),
	}

	return payload
}

// SITExtensions payload\
func SITExtensions(sitExtensions *models.SITExtensions) *primemessages.SITExtensions {
	if sitExtensions == nil {
		return nil
	}

	payload := make(primemessages.SITExtensions, len(*sitExtensions))

	for i, m := range *sitExtensions {
		copyOfM := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = SITExtension(&copyOfM)
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
