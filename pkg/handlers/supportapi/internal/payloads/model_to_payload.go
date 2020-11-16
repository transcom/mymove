package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// MoveTaskOrders payload
func MoveTaskOrders(moveTaskOrders *models.Moves) []*supportmessages.MoveTaskOrder {
	payload := make(supportmessages.MoveTaskOrders, len(*moveTaskOrders))

	for i, m := range *moveTaskOrders {
		//  G601 TODO needs review
		payload[i] = MoveTaskOrder(&m)
	}
	return payload
}

// MoveTaskOrder payload
func MoveTaskOrder(moveTaskOrder *models.Move) *supportmessages.MoveTaskOrder {
	if moveTaskOrder == nil {
		return nil
	}
	mtoShipments := MTOShipments(&moveTaskOrder.MTOShipments)
	payload := &supportmessages.MoveTaskOrder{
		ID:                 strfmt.UUID(moveTaskOrder.ID.String()),
		CreatedAt:          strfmt.DateTime(moveTaskOrder.CreatedAt),
		AvailableToPrimeAt: handlers.FmtDateTimePtr(moveTaskOrder.AvailableToPrimeAt),
		IsCanceled:         moveTaskOrder.IsCanceled(),
		MoveOrder:          MoveOrder(&moveTaskOrder.Orders),
		ReferenceID:        *moveTaskOrder.ReferenceID,
		ContractorID:       strfmt.UUID(moveTaskOrder.ContractorID.String()),
		MtoShipments:       *mtoShipments,
		UpdatedAt:          strfmt.DateTime(moveTaskOrder.UpdatedAt),
		ETag:               etag.GenerateEtag(moveTaskOrder.UpdatedAt),
		Status:             (supportmessages.MoveStatus)(moveTaskOrder.Status),
		Locator:            moveTaskOrder.Locator,
	}

	if moveTaskOrder.PPMEstimatedWeight != nil {
		payload.PpmEstimatedWeight = int64(*moveTaskOrder.PPMEstimatedWeight)
	}

	if moveTaskOrder.PPMType != nil {
		payload.PpmType = *moveTaskOrder.PPMType
	}

	return payload
}

// Customer payload
func Customer(customer *models.ServiceMember) *supportmessages.Customer {
	if customer == nil {
		return nil
	}
	payload := supportmessages.Customer{
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
	}
	return &payload
}

// MoveOrder payload
func MoveOrder(moveOrder *models.Order) *supportmessages.MoveOrder {
	if moveOrder == nil {
		return nil
	}
	destinationDutyStation := DutyStation(&moveOrder.NewDutyStation)
	originDutyStation := DutyStation(moveOrder.OriginDutyStation)
	uploadedOrders := Document(&moveOrder.UploadedOrders)
	if moveOrder.Grade != nil && moveOrder.Entitlement != nil {
		moveOrder.Entitlement.SetWeightAllotment(*moveOrder.Grade)
	}

	reportByDate := strfmt.Date(moveOrder.ReportByDate)
	issueDate := strfmt.Date(moveOrder.IssueDate)

	payload := supportmessages.MoveOrder{
		DestinationDutyStation:   destinationDutyStation,
		DestinationDutyStationID: destinationDutyStation.ID,
		Entitlement:              Entitlement(moveOrder.Entitlement),
		Customer:                 Customer(&moveOrder.ServiceMember),
		OrderNumber:              moveOrder.OrdersNumber,
		OrdersType:               supportmessages.OrdersType(moveOrder.OrdersType),
		ID:                       strfmt.UUID(moveOrder.ID.String()),
		OriginDutyStation:        originDutyStation,
		ETag:                     etag.GenerateEtag(moveOrder.UpdatedAt),
		Status:                   supportmessages.OrdersStatus(moveOrder.Status),
		UploadedOrders:           uploadedOrders,
		UploadedOrdersID:         strfmt.UUID(uploadedOrders.ID.String()),
		ReportByDate:             &reportByDate,
		IssueDate:                &issueDate,
	}

	if moveOrder.Grade != nil {
		payload.Rank = moveOrder.Grade
	}
	return &payload
}

// Entitlement payload
func Entitlement(entitlement *models.Entitlement) *supportmessages.Entitlement {
	if entitlement == nil {
		return nil
	}
	var proGearWeight, proGearWeightSpouse, totalWeight int64
	if entitlement.WeightAllotment() != nil {
		proGearWeight = int64(entitlement.WeightAllotment().ProGearWeight)
		proGearWeightSpouse = int64(entitlement.WeightAllotment().ProGearWeightSpouse)
		totalWeight = int64(entitlement.WeightAllotment().TotalWeightSelf)
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
	return &supportmessages.Entitlement{
		ID:                    strfmt.UUID(entitlement.ID.String()),
		AuthorizedWeight:      authorizedWeight,
		DependentsAuthorized:  entitlement.DependentsAuthorized,
		NonTemporaryStorage:   entitlement.NonTemporaryStorage,
		PrivatelyOwnedVehicle: entitlement.PrivatelyOwnedVehicle,
		ProGearWeight:         proGearWeight,
		ProGearWeightSpouse:   proGearWeightSpouse,
		StorageInTransit:      sit,
		TotalDependents:       totalDependents,
		TotalWeight:           totalWeight,
		ETag:                  etag.GenerateEtag(entitlement.UpdatedAt),
	}
}

// DutyStation payload
func DutyStation(dutyStation *models.DutyStation) *supportmessages.DutyStation {
	if dutyStation == nil {
		return nil
	}
	address := Address(&dutyStation.Address)
	payload := supportmessages.DutyStation{
		Address:   address,
		AddressID: address.ID,
		ID:        strfmt.UUID(dutyStation.ID.String()),
		Name:      dutyStation.Name,
		ETag:      etag.GenerateEtag(dutyStation.UpdatedAt),
	}
	return &payload
}

// Document payload
func Document(document *models.Document) *supportmessages.Document {
	if document == nil {
		return nil
	}
	formattedID := strfmt.UUID(document.ID.String())
	formattedServiceMemberID := strfmt.UUID(document.ServiceMemberID.String())
	payload := supportmessages.Document{
		ID:              &formattedID,
		ServiceMemberID: &formattedServiceMemberID,
	}
	return &payload
}

// Address payload
func Address(address *models.Address) *supportmessages.Address {
	if address == nil {
		return nil
	}
	return &supportmessages.Address{
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

// MTOShipment payload
func MTOShipment(mtoShipment *models.MTOShipment) *supportmessages.MTOShipment {
	strfmt.MarshalFormat = strfmt.RFC3339Micro
	var primeActualWeight int64
	if mtoShipment.PrimeActualWeight != nil {
		primeActualWeight = int64(*mtoShipment.PrimeActualWeight)
	}
	payload := &supportmessages.MTOShipment{
		ID:                       strfmt.UUID(mtoShipment.ID.String()),
		MoveTaskOrderID:          strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:             mtoShipment.ShipmentType,
		Status:                   string(mtoShipment.Status),
		CustomerRemarks:          mtoShipment.CustomerRemarks,
		RejectionReason:          mtoShipment.RejectionReason,
		PickupAddress:            Address(mtoShipment.PickupAddress),
		PrimeActualWeight:        primeActualWeight,
		SecondaryDeliveryAddress: Address(mtoShipment.SecondaryDeliveryAddress),
		SecondaryPickupAddress:   Address(mtoShipment.SecondaryPickupAddress),
		DestinationAddress:       Address(mtoShipment.DestinationAddress),
		CreatedAt:                strfmt.DateTime(mtoShipment.CreatedAt),
		UpdatedAt:                strfmt.DateTime(mtoShipment.UpdatedAt),
		ETag:                     etag.GenerateEtag(mtoShipment.UpdatedAt),
	}

	if mtoShipment.RequestedPickupDate != nil {
		payload.RequestedPickupDate = strfmt.Date(*mtoShipment.RequestedPickupDate)
	}

	if mtoShipment.ApprovedDate != nil {
		payload.ApprovedDate = strfmt.Date(*mtoShipment.ApprovedDate)
	}

	return payload
}

// MTOServiceItem payload
func MTOServiceItem(mtoServiceItem *models.MTOServiceItem) *supportmessages.UpdateMTOServiceItemStatus {
	strfmt.MarshalFormat = strfmt.RFC3339Micro
	payload := &supportmessages.UpdateMTOServiceItemStatus{
		ETag:            etag.GenerateEtag(mtoServiceItem.UpdatedAt),
		ID:              strfmt.UUID(mtoServiceItem.ID.String()),
		MoveTaskOrderID: strfmt.UUID(mtoServiceItem.MoveTaskOrderID.String()),
		MtoShipmentID:   strfmt.UUID(mtoServiceItem.MTOShipmentID.String()),
		Status:          supportmessages.MTOServiceItemStatus(mtoServiceItem.Status),
		RejectionReason: mtoServiceItem.RejectionReason,
	}

	return payload
}

// MTOShipments payload
func MTOShipments(mtoShipments *models.MTOShipments) *supportmessages.MTOShipments {
	payload := make(supportmessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		//  G601 TODO needs review
		payload[i] = MTOShipment(&m)
	}
	return &payload
}

// MTOAgent payload
func MTOAgent(mtoAgent *models.MTOAgent) *supportmessages.MTOAgent {
	payload := &supportmessages.MTOAgent{
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
func MTOAgents(mtoAgents *models.MTOAgents) *supportmessages.MTOAgents {
	payload := make(supportmessages.MTOAgents, len(*mtoAgents))
	for i, m := range *mtoAgents {
		//  G601 TODO needs review
		payload[i] = MTOAgent(&m)
	}
	return &payload
}

// PaymentRequest payload
func PaymentRequest(pr *models.PaymentRequest) *supportmessages.PaymentRequest {
	return &supportmessages.PaymentRequest{
		ID:                   *handlers.FmtUUID(pr.ID),
		IsFinal:              &pr.IsFinal,
		MoveTaskOrderID:      *handlers.FmtUUID(pr.MoveTaskOrderID),
		PaymentRequestNumber: pr.PaymentRequestNumber,
		RejectionReason:      pr.RejectionReason,
		Status:               supportmessages.PaymentRequestStatus(pr.Status),
		ETag:                 etag.GenerateEtag(pr.UpdatedAt),
	}
}

// PaymentRequests payload
func PaymentRequests(paymentRequests *models.PaymentRequests) *supportmessages.PaymentRequests {
	payload := make(supportmessages.PaymentRequests, len(*paymentRequests))

	for i, pr := range *paymentRequests {
		//  G601 TODO needs review
		payload[i] = PaymentRequest(&pr)
	}
	return &payload
}

// InternalServerError describes errors in a standard structure to be returned in the payload.
// If detail is nil, string defaults to "An internal server error has occurred."
func InternalServerError(detail *string, traceID uuid.UUID) *supportmessages.Error {
	payload := supportmessages.Error{
		Title:    handlers.FmtString(handlers.InternalServerErrMessage),
		Detail:   handlers.FmtString(handlers.InternalServerErrDetail),
		Instance: strfmt.UUID(traceID.String()),
	}
	if detail != nil {
		payload.Detail = detail
	}
	return &payload
}

// ValidationError payload describes validation errors from the model or properties
func ValidationError(detail string, instance uuid.UUID, validationErrors *validate.Errors) *supportmessages.ValidationError {
	payload := &supportmessages.ValidationError{
		ClientError: *ClientError(handlers.ValidationErrMessage, detail, instance),
	}
	if validationErrors != nil {
		payload.InvalidFields = handlers.NewValidationErrorListResponse(validationErrors).Errors
	}
	return payload
}

// ClientError payload contains the default information we send to the client on errors
func ClientError(title string, detail string, instance uuid.UUID) *supportmessages.ClientError {
	return &supportmessages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
}
