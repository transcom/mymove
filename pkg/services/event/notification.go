package event

import (
	"go.uber.org/zap"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
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
func notificationSave(event *Event, objectType string, payload *[]byte) error {
	payloadString := string(*payload)
	newNotification := models.WebhookNotification{
		EventKey:   string(event.EventKey),
		ObjectType: swag.String(objectType),
		Payload:    swag.String(payloadString),
		Status:     models.WebhookNotificationPending,
	}

	trace := event.HandlerContext.GetTraceID()
	if trace != uuid.Nil {
		newNotification.TraceID = &trace
	}

	// Creates the notification entry in the DB
	verrs, err := event.DBConnection.ValidateAndCreate(&newNotification)
	if verrs.Count() > 0 {
		event.logger.Error("event.notificationSave error", zap.Error(verrs))
		return services.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	} else if err != nil {
		event.logger.Error("event.notificationSave error", zap.Error(err))
		e := services.NewQueryError("Notification", err, "Unable to save Notification.")
		return e
	}
	return nil

}
func checkAvailabilityToPrime(event *Event) (bool, error) {
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker(event.DBConnection)
	availableToPrime, err := mtoChecker.MTOAvailableToPrime(event.MtoID)
	if err != nil {
		unknownErr := services.NewEventError("Unknown error checking prime availability", err)
		return false, unknownErr
	}
	// No need to store notification to send if not available to prime
	if !availableToPrime {
		return false, nil
	}
	return true, nil

}

func paymentRequestEventHandler(event *Event) (bool, error) {
	// CHECK SOURCE
	// Continue only if source of event is not Prime
	if isSourcePrime(event) {
		return false, nil
	}

	// CHECK FOR AVAILABILITY TO PRIME
	// Continue only if MTO is available to Prime
	if isAvailableToPrime, err := checkAvailabilityToPrime(event); !isAvailableToPrime {
		return false, err
	}

	// ASSEMBLE PAYLOAD
	var db = event.DBConnection
	model := models.PaymentRequest{}

	// Important to be specific about which addl associations to load to reduce DB hits
	err := db.Eager("PaymentServiceItems", "PaymentServiceItems.PaymentServiceItemParams").Find(&model, event.UpdatedObjectID.String())
	if err != nil {
		notFoundError := services.NewNotFoundError(event.UpdatedObjectID, "looking for PaymentRequest")
		notFoundError.Wrap(err)
		return false, notFoundError
	}
	payload := PaymentRequestModelToPayload(&model)
	payloadArray, err := payload.MarshalBinary()
	if err != nil {
		unknownErr := services.NewEventError("Unknown error creating payload", err)
		return false, unknownErr
	}

	// STORE NOTIFICATION IN DB
	err = notificationSave(event, "PaymentRequest", &payloadArray)
	if err != nil {
		unknownErr := services.NewEventError("Unknown error storing notification", err)
		return false, unknownErr
	}

	return true, nil
}

// NotificationEventHandler receives notifications from the events package
// For alerting ALL errors should be logged here.
func NotificationEventHandler(event *Event) error {

	var logger = event.logger

	// Currently it logs information about the event. Eventually it will create an entry
	// in the notification table in the database
	logger.Info("event.NotificationEventHandler: Running with event:",
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
		stored, err := paymentRequestEventHandler(event)
		if err != nil {
			event.logger.Error("event.NotificationEventHandler: ", zap.Error(err))
			return err
		} else if !stored {
			event.logger.Info("event.NotificationEventHandler: No notification created.")
		}
		event.logger.Info("SUCCESS!")
	}

	return nil
}
