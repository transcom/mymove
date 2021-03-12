package payloads

import (
	"math"
	"time"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// Contractor payload
func Contractor(contractor *models.Contractor) *ghcmessages.Contractor {
	if contractor == nil {
		return nil
	}

	payload := &ghcmessages.Contractor{
		ID:             strfmt.UUID(contractor.ID.String()),
		ContractNumber: contractor.ContractNumber,
		Name:           contractor.Name,
		Type:           contractor.Type,
	}

	return payload
}

// Move payload
func Move(move *models.Move) *ghcmessages.Move {
	if move == nil {
		return nil
	}

	payload := &ghcmessages.Move{
		ID:                 strfmt.UUID(move.ID.String()),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(move.AvailableToPrimeAt),
		ContractorID:       handlers.FmtUUIDPtr(move.ContractorID),
		Contractor:         Contractor(move.Contractor),
		Locator:            move.Locator,
		OrdersID:           strfmt.UUID(move.OrdersID.String()),
		Orders:             Order(&move.Orders),
		ReferenceID:        handlers.FmtStringPtr(move.ReferenceID),
		Status:             ghcmessages.MoveStatus(move.Status),
		CreatedAt:          strfmt.DateTime(move.CreatedAt),
		SubmittedAt:        handlers.FmtDateTimePtr(move.SubmittedAt),
		UpdatedAt:          strfmt.DateTime(move.UpdatedAt),
		ETag:               etag.GenerateEtag(move.UpdatedAt),
	}

	return payload
}

// MoveTaskOrder payload
func MoveTaskOrder(moveTaskOrder *models.Move) *ghcmessages.MoveTaskOrder {
	if moveTaskOrder == nil {
		return nil
	}

	payload := &ghcmessages.MoveTaskOrder{
		ID:                 strfmt.UUID(moveTaskOrder.ID.String()),
		CreatedAt:          strfmt.DateTime(moveTaskOrder.CreatedAt),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(moveTaskOrder.AvailableToPrimeAt),
		IsCanceled:         moveTaskOrder.IsCanceled(),
		OrderID:            strfmt.UUID(moveTaskOrder.OrdersID.String()),
		ReferenceID:        *moveTaskOrder.ReferenceID,
		UpdatedAt:          strfmt.DateTime(moveTaskOrder.UpdatedAt),
		ETag:               etag.GenerateEtag(moveTaskOrder.UpdatedAt),
		Locator:            moveTaskOrder.Locator,
	}
	return payload
}

// Customer payload
func Customer(customer *models.ServiceMember) *ghcmessages.Customer {
	if customer == nil {
		return nil
	}

	payload := ghcmessages.Customer{
		Agency:         swag.StringValue((*string)(customer.Affiliation)),
		CurrentAddress: Address(customer.ResidentialAddress),
		DodID:          swag.StringValue(customer.Edipi),
		Email:          customer.PersonalEmail,
		FirstName:      swag.StringValue(customer.FirstName),
		ID:             strfmt.UUID(customer.ID.String()),
		LastName:       swag.StringValue(customer.LastName),
		Phone:          customer.Telephone,
		UserID:         strfmt.UUID(customer.UserID.String()),
		ETag:           etag.GenerateEtag(customer.UpdatedAt),
		BackupContact:  BackupContact(customer.BackupContacts),
	}
	return &payload
}

// Order payload
func Order(order *models.Order) *ghcmessages.Order {
	if order == nil {
		return nil
	}
	if order.ID == uuid.Nil {
		return nil
	}

	destinationDutyStation := DutyStation(&order.NewDutyStation)
	originDutyStation := DutyStation(order.OriginDutyStation)
	if order.Grade != nil && order.Entitlement != nil {
		order.Entitlement.SetWeightAllotment(*order.Grade)
	}
	entitlements := Entitlement(order.Entitlement)

	var deptIndicator ghcmessages.DeptIndicator
	if order.DepartmentIndicator != nil {
		deptIndicator = ghcmessages.DeptIndicator(*order.DepartmentIndicator)
	}

	var ordersTypeDetail ghcmessages.OrdersTypeDetail
	if order.OrdersTypeDetail != nil {
		ordersTypeDetail = ghcmessages.OrdersTypeDetail(*order.OrdersTypeDetail)
	}

	var grade ghcmessages.Grade
	if order.Grade != nil {
		grade = ghcmessages.Grade(*order.Grade)
	}
	//
	var branch ghcmessages.Branch
	if order.ServiceMember.Affiliation != nil {
		branch = ghcmessages.Branch(*order.ServiceMember.Affiliation)
	}

	var moveCode string
	if order.Moves != nil && len(order.Moves) > 0 {
		moveCode = order.Moves[0].Locator
	}

	payload := ghcmessages.Order{
		DestinationDutyStation: destinationDutyStation,
		Entitlement:            entitlements,
		Grade:                  &grade,
		OrderNumber:            order.OrdersNumber,
		OrderTypeDetail:        &ordersTypeDetail,
		ID:                     strfmt.UUID(order.ID.String()),
		OriginDutyStation:      originDutyStation,
		ETag:                   etag.GenerateEtag(order.UpdatedAt),
		Agency:                 branch,
		CustomerID:             strfmt.UUID(order.ServiceMemberID.String()),
		Customer:               Customer(&order.ServiceMember),
		FirstName:              swag.StringValue(order.ServiceMember.FirstName),
		LastName:               swag.StringValue(order.ServiceMember.LastName),
		ReportByDate:           strfmt.Date(order.ReportByDate),
		DateIssued:             strfmt.Date(order.IssueDate),
		OrderType:              ghcmessages.OrdersType(order.OrdersType),
		DepartmentIndicator:    &deptIndicator,
		Tac:                    handlers.FmtStringPtr(order.TAC),
		Sac:                    handlers.FmtStringPtr(order.SAC),
		UploadedOrderID:        strfmt.UUID(order.UploadedOrdersID.String()),
		MoveCode:               moveCode,
	}

	return &payload
}

// Entitlement payload
func Entitlement(entitlement *models.Entitlement) *ghcmessages.Entitlements {
	if entitlement == nil {
		return nil
	}
	var proGearWeight, proGearWeightSpouse, totalWeight int64
	if weightAllotment := entitlement.WeightAllotment(); weightAllotment != nil {
		proGearWeight = int64(weightAllotment.ProGearWeight)
		proGearWeightSpouse = int64(weightAllotment.ProGearWeightSpouse)
		if *entitlement.DependentsAuthorized {
			totalWeight = int64(weightAllotment.TotalWeightSelfPlusDependents)
		} else {
			totalWeight = int64(weightAllotment.TotalWeightSelf)
		}
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
	return &ghcmessages.Entitlements{
		ID:                    strfmt.UUID(entitlement.ID.String()),
		AuthorizedWeight:      authorizedWeight,
		DependentsAuthorized:  entitlement.DependentsAuthorized,
		NonTemporaryStorage:   entitlement.NonTemporaryStorage,
		PrivatelyOwnedVehicle: entitlement.PrivatelyOwnedVehicle,
		ProGearWeight:         proGearWeight,
		ProGearWeightSpouse:   proGearWeightSpouse,
		StorageInTransit:      &sit,
		TotalDependents:       totalDependents,
		TotalWeight:           totalWeight,
		ETag:                  etag.GenerateEtag(entitlement.UpdatedAt),
	}
}

// DutyStation payload
func DutyStation(dutyStation *models.DutyStation) *ghcmessages.DutyStation {
	if dutyStation == nil {
		return nil
	}
	address := Address(&dutyStation.Address)
	payload := ghcmessages.DutyStation{
		Address:   address,
		AddressID: address.ID,
		ID:        strfmt.UUID(dutyStation.ID.String()),
		Name:      dutyStation.Name,
		ETag:      etag.GenerateEtag(dutyStation.UpdatedAt),
	}
	return &payload
}

// Address payload
func Address(address *models.Address) *ghcmessages.Address {
	if address == nil {
		return nil
	}
	return &ghcmessages.Address{
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

// BackupContact payload
func BackupContact(contacts models.BackupContacts) *ghcmessages.BackupContact {
	var name, email, phone string

	if len(contacts) != 0 {
		contact := contacts[0]
		name = contact.Name
		email = contact.Email
		phone = ""
		contactPhone := contact.Phone
		if contactPhone != nil {
			phone = *contactPhone
		}
	}

	return &ghcmessages.BackupContact{
		Name:  &name,
		Email: &email,
		Phone: &phone,
	}
}

// MTOShipment payload
func MTOShipment(mtoShipment *models.MTOShipment) *ghcmessages.MTOShipment {
	strfmt.MarshalFormat = strfmt.RFC3339Micro

	payload := &ghcmessages.MTOShipment{
		ID:                       strfmt.UUID(mtoShipment.ID.String()),
		MoveTaskOrderID:          strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:             mtoShipment.ShipmentType,
		Status:                   ghcmessages.MTOShipmentStatus(mtoShipment.Status),
		CustomerRemarks:          mtoShipment.CustomerRemarks,
		RejectionReason:          mtoShipment.RejectionReason,
		PickupAddress:            Address(mtoShipment.PickupAddress),
		SecondaryDeliveryAddress: Address(mtoShipment.SecondaryDeliveryAddress),
		SecondaryPickupAddress:   Address(mtoShipment.SecondaryPickupAddress),
		DestinationAddress:       Address(mtoShipment.DestinationAddress),
		PrimeEstimatedWeight:     handlers.FmtPoundPtr(mtoShipment.PrimeEstimatedWeight),
		PrimeActualWeight:        handlers.FmtPoundPtr(mtoShipment.PrimeActualWeight),
		MtoAgents:                *MTOAgents(&mtoShipment.MTOAgents),
		MtoServiceItems:          MTOServiceItemModels(mtoShipment.MTOServiceItems),
		CreatedAt:                strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                strfmt.DateTime(mtoShipment.UpdatedAt),
		ETag:                     etag.GenerateEtag(mtoShipment.UpdatedAt),
	}

	if mtoShipment.RequestedPickupDate != nil && !mtoShipment.RequestedPickupDate.IsZero() {
		payload.RequestedPickupDate = *handlers.FmtDatePtr(mtoShipment.RequestedPickupDate)
	}

	if mtoShipment.ApprovedDate != nil {
		payload.ApprovedDate = strfmt.Date(*mtoShipment.ApprovedDate)
	}

	if mtoShipment.ScheduledPickupDate != nil {
		payload.ScheduledPickupDate = strfmt.Date(*mtoShipment.ScheduledPickupDate)
	}

	return payload
}

// MTOShipments payload
func MTOShipments(mtoShipments *models.MTOShipments) *ghcmessages.MTOShipments {
	payload := make(ghcmessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		copyOfMtoShipment := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOShipment(&copyOfMtoShipment)
	}
	return &payload
}

// MTOAgent payload
func MTOAgent(mtoAgent *models.MTOAgent) *ghcmessages.MTOAgent {
	payload := &ghcmessages.MTOAgent{
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
func MTOAgents(mtoAgents *models.MTOAgents) *ghcmessages.MTOAgents {
	payload := make(ghcmessages.MTOAgents, len(*mtoAgents))
	for i, m := range *mtoAgents {
		copyOfMtoAgent := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOAgent(&copyOfMtoAgent)
	}
	return &payload
}

// PaymentRequests payload
func PaymentRequests(prs *models.PaymentRequests, storer storage.FileStorer) (*ghcmessages.PaymentRequests, error) {
	payload := make(ghcmessages.PaymentRequests, len(*prs))

	for i, p := range *prs {
		paymentRequest := p
		pr, err := PaymentRequest(&paymentRequest, storer)
		if err != nil {
			return nil, err
		}
		payload[i] = pr
	}
	return &payload, nil
}

// PaymentRequest payload
func PaymentRequest(pr *models.PaymentRequest, storer storage.FileStorer) (*ghcmessages.PaymentRequest, error) {
	serviceDocs := make(ghcmessages.ProofOfServiceDocs, len(pr.ProofOfServiceDocs))

	if pr.ProofOfServiceDocs != nil && len(pr.ProofOfServiceDocs) > 0 {
		for i, proofOfService := range pr.ProofOfServiceDocs {
			payload, err := ProofOfServiceDoc(proofOfService, storer)
			if err != nil {
				return nil, err
			}
			serviceDocs[i] = payload
		}
	}

	return &ghcmessages.PaymentRequest{
		ID:                   *handlers.FmtUUID(pr.ID),
		IsFinal:              &pr.IsFinal,
		MoveTaskOrderID:      *handlers.FmtUUID(pr.MoveTaskOrderID),
		MoveTaskOrder:        Move(&pr.MoveTaskOrder),
		PaymentRequestNumber: pr.PaymentRequestNumber,
		RejectionReason:      pr.RejectionReason,
		Status:               ghcmessages.PaymentRequestStatus(pr.Status),
		ETag:                 etag.GenerateEtag(pr.UpdatedAt),
		ServiceItems:         *PaymentServiceItems(&pr.PaymentServiceItems),
		ReviewedAt:           handlers.FmtDateTimePtr(pr.ReviewedAt),
		ProofOfServiceDocs:   serviceDocs,
		CreatedAt:            strfmt.DateTime(pr.CreatedAt),
	}, nil
}

// PaymentServiceItem payload
func PaymentServiceItem(ps *models.PaymentServiceItem) *ghcmessages.PaymentServiceItem {
	if ps == nil {
		return nil
	}
	paymentServiceItemParams := PaymentServiceItemParams(&ps.PaymentServiceItemParams)

	return &ghcmessages.PaymentServiceItem{
		ID:                       *handlers.FmtUUID(ps.ID),
		MtoServiceItemID:         *handlers.FmtUUID(ps.MTOServiceItemID),
		MtoServiceItemName:       ps.MTOServiceItem.ReService.Name,
		MtoShipmentType:          ghcmessages.MTOShipmentType(ps.MTOServiceItem.MTOShipment.ShipmentType),
		MtoShipmentID:            handlers.FmtUUIDPtr(ps.MTOServiceItem.MTOShipmentID),
		CreatedAt:                strfmt.DateTime(ps.CreatedAt),
		PriceCents:               handlers.FmtCost(ps.PriceCents),
		RejectionReason:          ps.RejectionReason,
		Status:                   ghcmessages.PaymentServiceItemStatus(ps.Status),
		ReferenceID:              ps.ReferenceID,
		ETag:                     etag.GenerateEtag(ps.UpdatedAt),
		PaymentServiceItemParams: *paymentServiceItemParams,
	}
}

// PaymentServiceItems payload
func PaymentServiceItems(paymentServiceItems *models.PaymentServiceItems) *ghcmessages.PaymentServiceItems {
	payload := make(ghcmessages.PaymentServiceItems, len(*paymentServiceItems))
	for i, m := range *paymentServiceItems {
		copyOfPaymentServiceItem := m // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = PaymentServiceItem(&copyOfPaymentServiceItem)
	}
	return &payload
}

// PaymentServiceItemParam payload
func PaymentServiceItemParam(paymentServiceItemParam models.PaymentServiceItemParam) *ghcmessages.PaymentServiceItemParam {
	return &ghcmessages.PaymentServiceItemParam{
		ID:                   strfmt.UUID(paymentServiceItemParam.ID.String()),
		PaymentServiceItemID: strfmt.UUID(paymentServiceItemParam.PaymentServiceItemID.String()),
		Key:                  ghcmessages.ServiceItemParamName(paymentServiceItemParam.ServiceItemParamKey.Key),
		Value:                paymentServiceItemParam.Value,
		Type:                 ghcmessages.ServiceItemParamType(paymentServiceItemParam.ServiceItemParamKey.Type),
		Origin:               ghcmessages.ServiceItemParamOrigin(paymentServiceItemParam.ServiceItemParamKey.Origin),
		ETag:                 etag.GenerateEtag(paymentServiceItemParam.UpdatedAt),
	}
}

// PaymentServiceItemParams payload
func PaymentServiceItemParams(paymentServiceItemParams *models.PaymentServiceItemParams) *ghcmessages.PaymentServiceItemParams {
	if paymentServiceItemParams == nil {
		return nil
	}

	payload := make(ghcmessages.PaymentServiceItemParams, len(*paymentServiceItemParams))

	for i, p := range *paymentServiceItemParams {
		payload[i] = PaymentServiceItemParam(p)
	}
	return &payload
}

// MTOServiceItemModel payload
func MTOServiceItemModel(s *models.MTOServiceItem) *ghcmessages.MTOServiceItem {
	if s == nil {
		return nil
	}

	return &ghcmessages.MTOServiceItem{
		ID:               handlers.FmtUUID(s.ID),
		MoveTaskOrderID:  handlers.FmtUUID(s.MoveTaskOrderID),
		MtoShipmentID:    handlers.FmtUUIDPtr(s.MTOShipmentID),
		ReServiceID:      handlers.FmtUUID(s.ReServiceID),
		ReServiceCode:    handlers.FmtString(string(s.ReService.Code)),
		ReServiceName:    handlers.FmtStringPtr(&s.ReService.Name),
		Reason:           handlers.FmtStringPtr(s.Reason),
		RejectionReason:  handlers.FmtStringPtr(s.RejectionReason),
		PickupPostalCode: handlers.FmtStringPtr(s.PickupPostalCode),
		Status:           ghcmessages.MTOServiceItemStatus(s.Status),
		Description:      handlers.FmtStringPtr(s.Description),
		Dimensions:       MTOServiceItemDimensions(s.Dimensions),
		CustomerContacts: MTOServiceItemCustomerContacts(s.CustomerContacts),
		CreatedAt:        strfmt.DateTime(s.CreatedAt),
		ApprovedAt:       handlers.FmtDateTimePtr(s.ApprovedAt),
		RejectedAt:       handlers.FmtDateTimePtr(s.RejectedAt),
		ETag:             etag.GenerateEtag(s.UpdatedAt),
	}
}

// MTOServiceItemModels payload
func MTOServiceItemModels(s models.MTOServiceItems) ghcmessages.MTOServiceItems {
	serviceItems := ghcmessages.MTOServiceItems{}
	for _, item := range s {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		serviceItems = append(serviceItems, MTOServiceItemModel(&copyOfServiceItem))
	}

	return serviceItems
}

// MTOServiceItemDimension payload
func MTOServiceItemDimension(d *models.MTOServiceItemDimension) *ghcmessages.MTOServiceItemDimension {
	return &ghcmessages.MTOServiceItemDimension{
		ID:     *handlers.FmtUUID(d.ID),
		Type:   ghcmessages.DimensionType(d.Type),
		Length: *d.Length.Int32Ptr(),
		Height: *d.Height.Int32Ptr(),
		Width:  *d.Width.Int32Ptr(),
	}
}

// MTOServiceItemDimensions payload
func MTOServiceItemDimensions(d models.MTOServiceItemDimensions) ghcmessages.MTOServiceItemDimensions {
	payload := make(ghcmessages.MTOServiceItemDimensions, len(d))
	for i, item := range d {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOServiceItemDimension(&copyOfServiceItem)
	}
	return payload
}

// MTOServiceItemCustomerContact payload
func MTOServiceItemCustomerContact(c *models.MTOServiceItemCustomerContact) *ghcmessages.MTOServiceItemCustomerContact {
	return &ghcmessages.MTOServiceItemCustomerContact{
		Type:                       ghcmessages.CustomerContactType(c.Type),
		TimeMilitary:               c.TimeMilitary,
		FirstAvailableDeliveryDate: *handlers.FmtDate(c.FirstAvailableDeliveryDate),
	}
}

// MTOServiceItemCustomerContacts payload
func MTOServiceItemCustomerContacts(c models.MTOServiceItemCustomerContacts) ghcmessages.MTOServiceItemCustomerContacts {
	payload := make(ghcmessages.MTOServiceItemCustomerContacts, len(c))
	for i, item := range c {
		copyOfServiceItem := item // Make copy to avoid implicit memory aliasing of items from a range statement.
		payload[i] = MTOServiceItemCustomerContact(&copyOfServiceItem)
	}
	return payload
}

// Upload payload
func Upload(storer storage.FileStorer, upload models.Upload, url string) *ghcmessages.Upload {
	uploadPayload := &ghcmessages.Upload{
		ID:          handlers.FmtUUID(upload.ID),
		Filename:    swag.String(upload.Filename),
		ContentType: swag.String(upload.ContentType),
		URL:         handlers.FmtURI(url),
		Bytes:       &upload.Bytes,
		CreatedAt:   handlers.FmtDateTime(upload.CreatedAt),
		UpdatedAt:   handlers.FmtDateTime(upload.UpdatedAt),
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload
}

// ProofOfServiceDoc payload from model
func ProofOfServiceDoc(proofOfService models.ProofOfServiceDoc, storer storage.FileStorer) (*ghcmessages.ProofOfServiceDoc, error) {

	uploads := make([]*ghcmessages.Upload, len(proofOfService.PrimeUploads))
	if proofOfService.PrimeUploads != nil && len(proofOfService.PrimeUploads) > 0 {
		for i, primeUpload := range proofOfService.PrimeUploads {
			url, err := storer.PresignedURL(primeUpload.Upload.StorageKey, primeUpload.Upload.ContentType)
			if err != nil {
				return nil, err
			}
			uploads[i] = Upload(storer, primeUpload.Upload, url)
		}
	}

	return &ghcmessages.ProofOfServiceDoc{
		Uploads: uploads,
	}, nil
}

// QueueMoves payload
func QueueMoves(moves []models.Move) *ghcmessages.QueueMoves {
	queueMoves := make(ghcmessages.QueueMoves, len(moves))
	for i, move := range moves {
		customer := move.Orders.ServiceMember

		var validMTOShipments []models.MTOShipment
		for _, shipment := range move.MTOShipments {
			if shipment.Status == models.MTOShipmentStatusSubmitted || shipment.Status == models.MTOShipmentStatusApproved {
				validMTOShipments = append(validMTOShipments, shipment)
			}
		}

		var deptIndicator ghcmessages.DeptIndicator
		if move.Orders.DepartmentIndicator != nil {
			deptIndicator = ghcmessages.DeptIndicator(*move.Orders.DepartmentIndicator)
		}

		queueMoves[i] = &ghcmessages.QueueMove{
			Customer:               Customer(&customer),
			Status:                 ghcmessages.QueueMoveStatus(move.Status),
			ID:                     *handlers.FmtUUID(move.Orders.ID),
			Locator:                move.Locator,
			DepartmentIndicator:    &deptIndicator,
			ShipmentsCount:         int64(len(validMTOShipments)),
			DestinationDutyStation: DutyStation(&move.Orders.NewDutyStation),
			OriginGBLOC:            ghcmessages.GBLOC(move.Orders.OriginDutyStation.TransportationOffice.Gbloc),
		}
	}
	return &queueMoves
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
)

// This is a helper function to calculate the inferred status needed for QueuePaymentRequest payload
func queuePaymentRequestStatus(paymentRequest models.PaymentRequest) string {
	// If a payment request is in the PENDING state, let's use the term 'payment requested'
	if paymentRequest.Status == models.PaymentRequestStatusPending {
		return QueuePaymentRequestPaymentRequested
	}

	// If a payment request is either reviewed, sent_to_gex or recieved_by_gex then we'll use 'reviewed'
	if paymentRequest.Status == models.PaymentRequestStatusSentToGex ||
		paymentRequest.Status == models.PaymentRequestStatusReceivedByGex ||
		paymentRequest.Status == models.PaymentRequestStatusReviewed {
		return QueuePaymentRequestReviewed
	}

	if paymentRequest.Status == models.PaymentRequestStatusReviewedAllRejected {
		return QueuePaymentRequestRejected
	}

	return QueuePaymentRequestPaid
}

// QueuePaymentRequests payload
func QueuePaymentRequests(paymentRequests *models.PaymentRequests) *ghcmessages.QueuePaymentRequests {
	queuePaymentRequests := make(ghcmessages.QueuePaymentRequests, len(*paymentRequests))

	for i, paymentRequest := range *paymentRequests {
		moveTaskOrder := paymentRequest.MoveTaskOrder
		orders := moveTaskOrder.Orders

		queuePaymentRequests[i] = &ghcmessages.QueuePaymentRequest{
			ID:          *handlers.FmtUUID(paymentRequest.ID),
			MoveID:      *handlers.FmtUUID(moveTaskOrder.ID),
			Customer:    Customer(&orders.ServiceMember),
			Status:      ghcmessages.PaymentRequestStatus(queuePaymentRequestStatus(paymentRequest)),
			Age:         int64(math.Ceil(time.Since(paymentRequest.CreatedAt).Hours() / 24.0)),
			SubmittedAt: *handlers.FmtDateTime(paymentRequest.CreatedAt), // RequestedAt does not seem to be populated
			Locator:     moveTaskOrder.Locator,
			OriginGBLOC: ghcmessages.GBLOC(orders.OriginDutyStation.TransportationOffice.Gbloc),
		}

		if orders.DepartmentIndicator != nil {
			deptIndicator := ghcmessages.DeptIndicator(*orders.DepartmentIndicator)
			queuePaymentRequests[i].DepartmentIndicator = &deptIndicator
		}
	}

	return &queuePaymentRequests
}
