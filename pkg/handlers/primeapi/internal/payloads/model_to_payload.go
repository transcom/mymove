package payloads

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
)

func MoveTaskOrder(moveTaskOrder *models.MoveTaskOrder) *primemessages.MoveTaskOrder {
	if moveTaskOrder == nil {
		return nil
	}
	paymentRequests := PaymentRequests(&moveTaskOrder.PaymentRequests)
	payload := &primemessages.MoveTaskOrder{
		ID:                 strfmt.UUID(moveTaskOrder.ID.String()),
		CreatedAt:          strfmt.Date(moveTaskOrder.CreatedAt),
		IsAvailableToPrime: &moveTaskOrder.IsAvailableToPrime,
		IsCanceled:         &moveTaskOrder.IsCanceled,
		MoveOrderID:        strfmt.UUID(moveTaskOrder.MoveOrderID.String()),
		ReferenceID:        moveTaskOrder.ReferenceID,
		PaymentRequests:    paymentRequests,
		UpdatedAt:          strfmt.Date(moveTaskOrder.UpdatedAt),
	}
	return payload
}

func MoveTaskOrders(moveTaskOrders *models.MoveTaskOrders) []*primemessages.MoveTaskOrder {
	payload := make(primemessages.MoveTaskOrders, len(*moveTaskOrders))

	for i, m := range *moveTaskOrders {
		payload[i] = MoveTaskOrder(&m)
	}
	return payload
}

func Customer(customer *models.Customer) *primemessages.Customer {
	if customer == nil {
		return nil
	}
	payload := primemessages.Customer{
		DodID:  customer.DODID,
		ID:     strfmt.UUID(customer.ID.String()),
		UserID: strfmt.UUID(customer.UserID.String()),
	}
	return &payload
}

func MoveOrder(moveOrders *models.MoveOrder) *primemessages.MoveOrder {
	if moveOrders == nil {
		return nil
	}
	destinationDutyStation := DutyStation(&moveOrders.DestinationDutyStation)
	originDutyStation := DutyStation(&moveOrders.OriginDutyStation)
	entitlements := Entitlement(&moveOrders.Entitlement)
	payload := primemessages.MoveOrder{
		CustomerID:             strfmt.UUID(moveOrders.CustomerID.String()),
		DestinationDutyStation: destinationDutyStation,
		Entitlement:            entitlements,
		ID:                     strfmt.UUID(moveOrders.ID.String()),
		OriginDutyStation:      originDutyStation,
	}
	return &payload
}

func Entitlement(entitlement *models.Entitlement) *primemessages.Entitlements {
	if entitlement == nil {
		return nil
	}
	var proGearWeight int64
	if entitlement.ProGearWeight != nil {
		proGearWeight = int64(*entitlement.ProGearWeight)
	}
	var proGearWeightSpouse int64
	if entitlement.ProGearWeightSpouse != nil {
		proGearWeightSpouse = int64(*entitlement.ProGearWeightSpouse)
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
		ID:                    strfmt.UUID(entitlement.ID.String()),
		DependentsAuthorized:  entitlement.DependentsAuthorized,
		NonTemporaryStorage:   entitlement.NonTemporaryStorage,
		PrivatelyOwnedVehicle: entitlement.PrivatelyOwnedVehicle,
		ProGearWeight:         proGearWeight,
		ProGearWeightSpouse:   proGearWeightSpouse,
		StorageInTransit:      sit,
		TotalDependents:       totalDependents,
	}
}

func DutyStation(dutyStation *models.DutyStation) *primemessages.DutyStation {
	if dutyStation == nil {
		return nil
	}
	address := Address(&dutyStation.Address)
	payload := primemessages.DutyStation{
		Address:   address,
		AddressID: address.ID,
		ID:        strfmt.UUID(dutyStation.ID.String()),
		Name:      dutyStation.Name,
	}
	return &payload
}

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
	}
}

func PaymentRequest(paymentRequest *models.PaymentRequest) *primemessages.PaymentRequest {
	return &primemessages.PaymentRequest{
		ID:              strfmt.UUID(paymentRequest.ID.String()),
		Status:          primemessages.PaymentRequestStatus(paymentRequest.Status),
		IsFinal:         &paymentRequest.IsFinal,
		MoveTaskOrderID: strfmt.UUID(paymentRequest.MoveTaskOrderID.String()),
		RejectionReason: paymentRequest.RejectionReason,
	}
}

func PaymentRequests(paymentRequests *[]models.PaymentRequest) []*primemessages.PaymentRequest {
	payload := make(primemessages.PaymentRequests, len(*paymentRequests))

	for i, p := range *paymentRequests {
		payload[i] = PaymentRequest(&p)
	}
	return payload
}
