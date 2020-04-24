package payloads

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"
)

// CustomerModel model
func CustomerModel(customer *supportmessages.Customer) *models.Customer {
	if customer == nil {
		return nil
	}
	return &models.Customer{
		ID:          uuid.FromStringOrNil(customer.ID.String()),
		Agency:      &customer.Agency,
		FirstName:   &customer.FirstName,
		LastName:    &customer.LastName,
		DODID:       &customer.DodID,
		Email:       customer.Email,
		PhoneNumber: customer.Phone,
	}
}

// MoveOrderModel returns a moveOrder model contstructed from the moveOrder message
func MoveOrderModel(moveOrderPayload *supportmessages.MoveOrder) *models.MoveOrder {
	if moveOrderPayload == nil {
		return nil
	}
	model := &models.MoveOrder{
		ID:              uuid.FromStringOrNil(moveOrderPayload.ID.String()),
		Grade:           &moveOrderPayload.Rank,
		OrderNumber:     moveOrderPayload.OrderNumber,
		OrderType:       moveOrderPayload.OrderType,
		OrderTypeDetail: moveOrderPayload.OrderTypeDetail,
	}

	destinationDutyStationID := uuid.FromStringOrNil(moveOrderPayload.DestinationDutyStationID.String())
	model.DestinationDutyStationID = &destinationDutyStationID

	originDutyStationID := uuid.FromStringOrNil(moveOrderPayload.OriginDutyStationID.String())
	model.OriginDutyStationID = &originDutyStationID

	reportByDate := time.Time(moveOrderPayload.ReportByDate)
	if !reportByDate.IsZero() {
		model.ReportByDate = &reportByDate
	}
	return model
}

// EntitlementModel returns a entitlement model constructed from the entitlement in the payload
func EntitlementModel(entitlementPayload *supportmessages.Entitlement) *models.Entitlement {
	if entitlementPayload == nil {
		return nil
	}
	model := &models.Entitlement{
		DependentsAuthorized:  entitlementPayload.DependentsAuthorized,
		NonTemporaryStorage:   entitlementPayload.NonTemporaryStorage,
		PrivatelyOwnedVehicle: entitlementPayload.PrivatelyOwnedVehicle,
	}

	totalDependents := int(entitlementPayload.TotalDependents)
	model.TotalDependents = &totalDependents

	storageInTransit := int(entitlementPayload.StorageInTransit)
	model.StorageInTransit = &storageInTransit

	return model
}
