package payloads

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func MoveTaskOrder(moveTaskOrder models.MoveTaskOrder) *primemessages.MoveTaskOrder {
	destinationAddress := Address(&moveTaskOrder.DestinationAddress)
	pickupAddress := Address(&moveTaskOrder.PickupAddress)
	entitlements := Entitlements(&moveTaskOrder.Entitlements)
	var primeEstimatedWeight *int64
	if moveTaskOrder.PrimeEstimatedWeight != nil {
		wt := moveTaskOrder.PrimeEstimatedWeight.Int64()
		primeEstimatedWeight = &wt
	}
	var primeEstimatedWeightRecordedDate *strfmt.Date
	if moveTaskOrder.PrimeEstimatedWeight != nil {
		primeEstimatedWeightRecordedDate = handlers.FmtDatePtr(moveTaskOrder.PrimeEstimatedWeightRecordedDate)
	}
	payload := &primemessages.MoveTaskOrder{
		CustomerID:                       strfmt.UUID(moveTaskOrder.CustomerID.String()),
		DestinationAddress:               destinationAddress,
		DestinationDutyStation:           strfmt.UUID(moveTaskOrder.DestinationDutyStation.ID.String()),
		Entitlements:                     entitlements,
		ID:                               strfmt.UUID(moveTaskOrder.ID.String()),
		MoveDate:                         strfmt.Date(moveTaskOrder.RequestedPickupDate),
		MoveID:                           strfmt.UUID(moveTaskOrder.MoveID.String()),
		OriginDutyStation:                strfmt.UUID(moveTaskOrder.OriginDutyStationID.String()),
		PickupAddress:                    pickupAddress,
		PrimeEstimatedWeight:             primeEstimatedWeight,
		PrimeEstimatedWeightRecordedDate: primeEstimatedWeightRecordedDate,
		ReferenceID:                      moveTaskOrder.ReferenceID,
		Remarks:                          moveTaskOrder.CustomerRemarks,
		RequestedPickupDate:              strfmt.Date(moveTaskOrder.RequestedPickupDate),
		Status:                           string(moveTaskOrder.Status),
		UpdatedAt:                        strfmt.Date(moveTaskOrder.UpdatedAt),
	}
	return payload
}

func Address(a *models.Address) *primemessages.Address {
	if a == nil {
		return nil
	}
	return &primemessages.Address{
		ID:             strfmt.UUID(a.ID.String()),
		StreetAddress1: swag.String(a.StreetAddress1),
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           swag.String(a.City),
		State:          swag.String(a.State),
		PostalCode:     swag.String(a.PostalCode),
		Country:        a.Country,
	}
}

func Entitlements(entitlement *models.GHCEntitlement) *primemessages.Entitlements {
	if entitlement == nil {
		return nil
	}
	return &primemessages.Entitlements{
		DependentsAuthorized:  entitlement.DependentsAuthorized,
		NonTemporaryStorage:   handlers.FmtBool(entitlement.NonTemporaryStorage),
		PrivatelyOwnedVehicle: handlers.FmtBool(entitlement.PrivatelyOwnedVehicle),
		ProGearWeight:         int64(entitlement.ProGearWeight),
		ProGearWeightSpouse:   int64(entitlement.ProGearWeightSpouse),
		StorageInTransit:      int64(entitlement.StorageInTransit),
		TotalDependents:       int64(entitlement.TotalDependents),
	}
}

func Customer(serviceMember *models.ServiceMember) *primemessages.Customer {
	if serviceMember == nil {
		return nil
	}
	var agency *string
	if serviceMember.Affiliation != nil {
		agency = handlers.FmtString(string(*serviceMember.Affiliation))
	}
	var rank *string
	if serviceMember.Rank != nil {
		rank = handlers.FmtString(string(*serviceMember.Rank))
	}

	return &primemessages.Customer{
		ID:            strfmt.UUID(serviceMember.ID.String()),
		Agency:        agency,
		Email:         serviceMember.PersonalEmail,
		FirstName:     serviceMember.FirstName,
		Grade:         rank,
		LastName:      serviceMember.LastName,
		MiddleName:    serviceMember.MiddleName,
		PickupAddress: Address(serviceMember.ResidentialAddress),
		Suffix:        serviceMember.Suffix,
		Telephone:     serviceMember.Telephone,
	}
}

// TODO maybe remove
func CustomerWithMTO(moveTaskOrder *models.MoveTaskOrder) *primemessages.Customer {
	if moveTaskOrder == nil {
		return nil
	}
	customer := Customer(&moveTaskOrder.Customer)
	return &primemessages.Customer{
		ID:                     strfmt.UUID(customer.ID.String()),
		Agency:                 customer.Agency,
		DestinationAddress:     Address(&moveTaskOrder.DestinationAddress),
		DestinationDutyStation: &moveTaskOrder.DestinationDutyStation.Name,
		Email:                  customer.Email,
		FirstName:              customer.FirstName,
		Grade:                  customer.Grade,
		LastName:               customer.LastName,
		MiddleName:             customer.MiddleName,
		OriginDutyStation:      &moveTaskOrder.OriginDutyStation.Name,
		PickupAddress:          Address(&moveTaskOrder.PickupAddress),
		Remarks:                moveTaskOrder.CustomerRemarks,
		RequestedPickupDate:    strfmt.Date(moveTaskOrder.RequestedPickupDate),
		Suffix:                 customer.Suffix,
		Telephone:              customer.Telephone,
	}
}
