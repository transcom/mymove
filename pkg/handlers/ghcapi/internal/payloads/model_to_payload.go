package payloads

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
)

func MoveTaskOrder(moveTaskOrder *models.MoveTaskOrder) *ghcmessages.MoveTaskOrder {
	if moveTaskOrder == nil {
		return nil
	}
	payload := &ghcmessages.MoveTaskOrder{
		ID:                 strfmt.UUID(moveTaskOrder.ID.String()),
		CreatedAt:          strfmt.Date(moveTaskOrder.CreatedAt),
		IsAvailableToPrime: &moveTaskOrder.IsAvailableToPrime,
		IsCanceled:         &moveTaskOrder.IsCancelled,
		MoveOrdersID:       strfmt.UUID(moveTaskOrder.MoveOrderID.String()),
		ReferenceID:        moveTaskOrder.ReferenceID,
		UpdatedAt:          strfmt.Date(moveTaskOrder.UpdatedAt),
	}
	return payload
}

func Customer(customer *models.Customer) *ghcmessages.Customer {
	if customer == nil {
		return nil
	}
	payload := ghcmessages.Customer{
		DodID:  customer.DODID,
		ID:     strfmt.UUID(customer.ID.String()),
		UserID: strfmt.UUID(customer.UserID.String()),
	}
	return &payload
}

func MoveOrders(moveOrders *models.MoveOrder) *ghcmessages.MoveOrder {
	if moveOrders == nil {
		return nil
	}
	destinationDutyStation := DutyStation(&moveOrders.DestinationDutyStation)
	originDutyStation := DutyStation(&moveOrders.OriginDutyStation)
	entitlements := Entitlements(&moveOrders.Entitlement)
	payload := ghcmessages.MoveOrder{
		CustomerID:             strfmt.UUID(moveOrders.CustomerID.String()),
		DestinationDutyStation: destinationDutyStation,
		Entitlement:            entitlements,
		ID:                     strfmt.UUID(moveOrders.ID.String()),
		OriginDutyStation:      originDutyStation,
	}
	return &payload
}

func Entitlements(entitlement *models.Entitlement) *ghcmessages.Entitlements {
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
	return &ghcmessages.Entitlements{
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
	}
	return &payload
}

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
	}
}
