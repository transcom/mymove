package payloads

import (
	"time"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/event"

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

// OrderModel converts payload to model - it does not convert nested
// duty stations but will preserve the ID if provided.
// It will create nested customer and entitlement models
// if those are provided in the payload
func OrderModel(orderPayload *supportmessages.Order) *models.Order {
	if orderPayload == nil {
		return nil
	}
	model := &models.Order{
		ID:            uuid.FromStringOrNil(orderPayload.ID.String()),
		Grade:         swag.String((string)(orderPayload.Rank)),
		OrdersNumber:  orderPayload.OrderNumber,
		ServiceMember: *CustomerModel(orderPayload.Customer),
		Entitlement:   EntitlementModel(orderPayload.Entitlement),
		TAC:           orderPayload.Tac,
	}

	customerID := uuid.FromStringOrNil(orderPayload.CustomerID.String())
	model.ServiceMemberID = customerID

	destinationDutyStationID := uuid.FromStringOrNil(orderPayload.DestinationDutyStationID.String())
	model.NewDutyStationID = destinationDutyStationID

	originDutyStationID := uuid.FromStringOrNil(orderPayload.OriginDutyStationID.String())
	model.OriginDutyStationID = &originDutyStationID

	reportByDate := time.Time(*orderPayload.ReportByDate)
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

	return model
}

// WebhookNotificationModel converts payload to model
func WebhookNotificationModel(payload *supportmessages.WebhookNotification, traceID uuid.UUID) (*models.WebhookNotification, *validate.Errors) {
	verrs := validate.NewErrors()

	if !event.ExistsEventKey(payload.EventKey) {
		verrs.Add("eventKey", "must be a registered event key")
		return nil, verrs
	}
	notification := &models.WebhookNotification{
		// ID is managed by pop, readonly
		EventKey:        payload.EventKey,
		TraceID:         &traceID,
		MoveTaskOrderID: handlers.FmtUUIDPtrToPopPtr(payload.MoveTaskOrderID),
		ObjectID:        handlers.FmtUUIDPtrToPopPtr(payload.ObjectID),
		Status:          models.WebhookNotificationPending,
		// CreatedAt is managed by pop, readonly
		// UpdatedAt is managed by pop, readonly
		// FirstAttemptedAt is never provided by user, readonly
	}
	if payload.Object != nil {
		notification.Payload = *payload.Object
	}
	return notification, nil
}
