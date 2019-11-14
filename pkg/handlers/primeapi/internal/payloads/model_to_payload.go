package payloads

import (
	"github.com/go-openapi/strfmt"

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
	var scheduledMoveDate strfmt.Date
	if moveTaskOrder.ScheduledMoveDate != nil {
		scheduledMoveDate = *handlers.FmtDate(*moveTaskOrder.ScheduledMoveDate)
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
		PpmIsIncluded:                    moveTaskOrder.PpmIsIncluded,
		PickupAddress:                    pickupAddress,
		PrimeEstimatedWeight:             primeEstimatedWeight,
		PrimeEstimatedWeightRecordedDate: primeEstimatedWeightRecordedDate,
		Remarks:                          moveTaskOrder.CustomerRemarks,
		RequestedPickupDate:              strfmt.Date(moveTaskOrder.RequestedPickupDate),
		ScheduledMoveDate:                scheduledMoveDate,
		SecondaryPickupAddress:           *moveTaskOrder.SecondaryPickupAddress,
		SecondaryDeliveryAddress:         *moveTaskOrder.SecondaryDeliveryAddress,
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
		StreetAddress1: &a.StreetAddress1,
		StreetAddress2: a.StreetAddress2,
		StreetAddress3: a.StreetAddress3,
		City:           &a.City,
		State:          &a.State,
		PostalCode:     &a.PostalCode,
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
