package payloads

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// CustomerModel converts payload to model - currently does not tackle addresses
func CustomerModel(customer *supportmessages.Customer) *models.ServiceMember {
	if customer == nil {
		return nil
	}
	return &models.ServiceMember{
		ID:            uuid.FromStringOrNil(customer.ID.String()),
		Affiliation:   (*models.ServiceMemberAffiliation)(customer.Agency),
		Edipi:         customer.DodID,
		FirstName:     customer.FirstName,
		LastName:      customer.LastName,
		PersonalEmail: customer.Email,
		Telephone:     customer.Phone,
	}
}

// MoveOrderModel converts payload to model - it does not convert nested
// duty stations but will preserve the ID if provided.
// It will create nested customer and entitlement models
// if those are provided in the payload
func MoveOrderModel(moveOrderPayload *supportmessages.MoveOrder) *models.Order {
	if moveOrderPayload == nil {
		return nil
	}
	model := &models.Order{
		ID:            uuid.FromStringOrNil(moveOrderPayload.ID.String()),
		Grade:         swag.String((string)(moveOrderPayload.Rank)),
		OrdersNumber:  moveOrderPayload.OrderNumber,
		ServiceMember: *CustomerModel(moveOrderPayload.Customer),
		Entitlement:   EntitlementModel(moveOrderPayload.Entitlement),
		TAC:           moveOrderPayload.Tac,
	}

	customerID := uuid.FromStringOrNil(moveOrderPayload.CustomerID.String())
	model.ServiceMemberID = customerID

	destinationDutyStationID := uuid.FromStringOrNil(moveOrderPayload.DestinationDutyStationID.String())
	model.NewDutyStationID = destinationDutyStationID

	originDutyStationID := uuid.FromStringOrNil(moveOrderPayload.OriginDutyStationID.String())
	model.OriginDutyStationID = &originDutyStationID

	reportByDate := time.Time(*moveOrderPayload.ReportByDate)
	if !reportByDate.IsZero() {
		model.ReportByDate = reportByDate
	}
	return model
}

// EntitlementModel converts the payload to model
func EntitlementModel(entitlementPayload *supportmessages.Entitlement) *models.Entitlement {
	if entitlementPayload == nil {
		return nil
	}

	// proGearWeight and ProGearWeightSpouse currently not handled as
	// they are not in the entitlement record in the db.
	model := &models.Entitlement{
		ID:                    uuid.FromStringOrNil(entitlementPayload.ID.String()),
		DependentsAuthorized:  entitlementPayload.DependentsAuthorized,
		NonTemporaryStorage:   entitlementPayload.NonTemporaryStorage,
		PrivatelyOwnedVehicle: entitlementPayload.PrivatelyOwnedVehicle,
	}

	if entitlementPayload.AuthorizedWeight != nil {
		model.DBAuthorizedWeight = swag.Int(int(*entitlementPayload.AuthorizedWeight))
	}

	totalDependents := int(entitlementPayload.TotalDependents)
	model.TotalDependents = &totalDependents

	storageInTransit := int(entitlementPayload.StorageInTransit)
	model.StorageInTransit = &storageInTransit

	return model
}

// MoveTaskOrderModel return an MTO model constructed from the payload.
// Does not create nested mtoServiceItems, mtoShipments, or paymentRequests
func MoveTaskOrderModel(mtoPayload *supportmessages.MoveTaskOrder) *models.Move {
	if mtoPayload == nil {
		return nil
	}
	ppmEstimatedWeight := unit.Pound(mtoPayload.PpmEstimatedWeight)
	contractorID := uuid.FromStringOrNil(mtoPayload.ContractorID.String())
	model := &models.Move{
		ReferenceID:        &mtoPayload.ReferenceID,
		PPMEstimatedWeight: &ppmEstimatedWeight,
		PPMType:            &mtoPayload.PpmType,
		ContractorID:       &contractorID,
	}

	if mtoPayload.AvailableToPrimeAt != nil {
		availableToPrimeAt := time.Time(*mtoPayload.AvailableToPrimeAt)
		model.AvailableToPrimeAt = &availableToPrimeAt
	}

	if mtoPayload.IsCanceled != nil && *mtoPayload.IsCanceled == true {
		model.Status = models.MoveStatusCANCELED
	}

	return model
}
