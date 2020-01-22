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
		IsCanceled:         &moveTaskOrder.IsCanceled,
		MoveOrderID:        strfmt.UUID(moveTaskOrder.MoveOrderID.String()),
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
		Agency:             customer.Agency,
		CurrentAddress:     Address(&customer.CurrentAddress),
		DestinationAddress: Address(&customer.DestinationAddress),
		DodID:              customer.DODID,
		Email:              customer.Email,
		FirstName:          customer.FirstName,
		ID:                 strfmt.UUID(customer.ID.String()),
		LastName:           customer.LastName,
		Phone:              customer.PhoneNumber,
		UserID:             strfmt.UUID(customer.UserID.String()),
	}
	return &payload
}

func MoveOrder(moveOrder *models.MoveOrder) *ghcmessages.MoveOrder {
	if moveOrder == nil {
		return nil
	}
	destinationDutyStation := DutyStation(&moveOrder.DestinationDutyStation)
	originDutyStation := DutyStation(&moveOrder.OriginDutyStation)
	moveOrder.Entitlement.SetWeightAllotment(moveOrder.Grade)
	entitlements := Entitlement(&moveOrder.Entitlement)
	payload := ghcmessages.MoveOrder{
		Agency:                 moveOrder.Customer.Agency,
		CustomerID:             strfmt.UUID(moveOrder.CustomerID.String()),
		FirstName:              moveOrder.Customer.FirstName,
		LastName:               moveOrder.Customer.LastName,
		ConfirmationNumber:     moveOrder.ConfirmationNumber,
		DestinationDutyStation: destinationDutyStation,
		Entitlement:            entitlements,
		Grade:                  moveOrder.Grade,
		ID:                     strfmt.UUID(moveOrder.ID.String()),
		OriginDutyStation:      originDutyStation,
	}
	return &payload
}

func Entitlement(entitlement *models.Entitlement) *ghcmessages.Entitlements {
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
	return &ghcmessages.Entitlements{
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

func MTOShipment(mtoShipment *models.MTOShipment) *ghcmessages.MTOShipment {
	return &ghcmessages.MTOShipment{
		ID:                  strfmt.UUID(mtoShipment.ID.String()),
		MoveTaskOrderID:     strfmt.UUID(mtoShipment.MoveTaskOrderID.String()),
		ShipmentType:        "HHG",
		Status:              string(mtoShipment.Status),
		CustomerRemarks:     *mtoShipment.CustomerRemarks,
		RequestedPickupDate: strfmt.Date(*mtoShipment.RequestedPickupDate),
		CreatedAt:           strfmt.Date(mtoShipment.CreatedAt),
		UpdatedAt:           strfmt.Date(mtoShipment.UpdatedAt),
	}
}

func MTOShipments(mtoShipments *models.MTOShipments) *ghcmessages.MTOShipments {
	payload := make(ghcmessages.MTOShipments, len(*mtoShipments))

	for i, m := range *mtoShipments {
		payload[i] = MTOShipment(&m)
	}
	return &payload
}
