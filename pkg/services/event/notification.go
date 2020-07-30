package event

import (
	"go.uber.org/zap"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// isSourcePrime returns true if the source of the event is prime
func isSourcePrime(event *Event) bool {
	// We can assume the source is prime if the endpoint that caused the event was
	// a prime endpoint. Support is considered non-prime endpoint.
	// Generally we should err on the side of more notifications rather than fewer
	// since we would not like them to miss an update.
	apiName := GetEndpointAPI(event.EndpointKey)
	if apiName == PrimeAPIName {
		return true
	}
	return false
}

// notificationSave saves the record in the webhook_notification table
// If it fails, it will return an error.
func notificationSave(event *Event, objectType string, payload *[]byte) error {
	payloadString := string(*payload)
	newNotification := models.WebhookNotification{
		EventKey: string(event.EventKey),
		Payload:  swag.String(payloadString),
		Status:   models.WebhookNotificationPending,
	}

// PaymentRequestModelToPayload This is an example to log the data
// Will need to move to payloads to model file.
func PaymentRequestModelToPayload(paymentRequest *models.PaymentRequest) *primemessages.PaymentRequest {
	if paymentRequest == nil {
		return nil
	}

	paymentServiceItems := PaymentServiceItemsModelToPayload(&paymentRequest.PaymentServiceItems)
	return &primemessages.PaymentRequest{
		ID:                   strfmt.UUID(paymentRequest.ID.String()),
		IsFinal:              &paymentRequest.IsFinal,
		MoveTaskOrderID:      strfmt.UUID(paymentRequest.MoveTaskOrderID.String()),
		PaymentRequestNumber: paymentRequest.PaymentRequestNumber,
		RejectionReason:      paymentRequest.RejectionReason,
		Status:               primemessages.PaymentRequestStatus(paymentRequest.Status),
		PaymentServiceItems:  *paymentServiceItems,
		ETag:                 etag.GenerateEtag(paymentRequest.UpdatedAt),
	}
}

// PaymentServiceItemModelToPayload payload
func PaymentServiceItemModelToPayload(paymentServiceItem *models.PaymentServiceItem) *primemessages.PaymentServiceItem {
	if paymentServiceItem == nil {
		return nil
	}

	payload := &primemessages.PaymentServiceItem{
		ID:               strfmt.UUID(paymentServiceItem.ID.String()),
		PaymentRequestID: strfmt.UUID(paymentServiceItem.PaymentRequestID.String()),
		MtoServiceItemID: strfmt.UUID(paymentServiceItem.MTOServiceItemID.String()),
		Status:           primemessages.PaymentServiceItemStatus(paymentServiceItem.Status),
		RejectionReason:  paymentServiceItem.RejectionReason,
		ETag:             etag.GenerateEtag(paymentServiceItem.UpdatedAt),
	}

	if paymentServiceItem.PriceCents != nil {
		payload.PriceCents = swag.Int64(int64(*paymentServiceItem.PriceCents))
	}

	return payload
}

// PaymentServiceItemsModelToPayload payload
func PaymentServiceItemsModelToPayload(paymentServiceItems *models.PaymentServiceItems) *primemessages.PaymentServiceItems {
	if paymentServiceItems == nil {
		return nil
	}

	payload := make(primemessages.PaymentServiceItems, len(*paymentServiceItems))

	for i, p := range *paymentServiceItems {
		payload[i] = PaymentServiceItemModelToPayload(&p)
	}
	return &payload
}

// NotificationEventHandler receives notifications from the events package
func NotificationEventHandler(event *Event) error {

	var logger = event.logger
	var db = event.DBConnection

	// Currently it logs information about the event. Eventually it will create an entry
	// in the notification table in the database
	logger.Info("Event handler ran:",
		zap.String("endpoint", string(event.EndpointKey)),
		zap.String("mtoID", event.MtoID.String()),
		zap.String("objectID", event.UpdatedObjectID.String()))

	//Get the type of model which is stored in the eventType
	modelBeingUpdated, err := GetModelFromEvent(event.EventKey)
	if err != nil {
		return err
	}

	// Based on which model was updated, construct the proper payload
	switch modelBeingUpdated.(type) {
	case models.PaymentRequest:
		model := models.PaymentRequest{}
		//Todo May need to check who made the call
		//Todo May need to check whether mto was available
		// Important to be specific about which addl associations to load to reduce DB hits
		err := db.Eager("PaymentServiceItems", "PaymentServiceItems.PaymentServiceItemParams").Find(&model, event.UpdatedObjectID.String())
		if err != nil {
			notFoundError := services.NewNotFoundError(event.UpdatedObjectID, "looking for PaymentRequest")
			notFoundError.Wrap(err)
			return notFoundError
		}
		payload := PaymentRequestModelToPayload(&model)
		logger.Info("Notification payload:", zap.Any("payload", *payload))
	}

	return nil
}
