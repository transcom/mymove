package event

import (
	"go.uber.org/zap"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

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

// checkAvailabilityToPrime returns true if the MTO is
// available to Prime. If there is a query error, it returns an
// error as well.
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

// paymentRequestEventHandler handles all events pertaining to paymentRequests
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

	//Get the model type of the updated logical object
	modelBeingUpdated, err := GetModelFromEvent(event.EventKey)
	if err != nil {
		event.logger.Error("event.NotificationEventHandler: EventKey does not exist.")
		return err
	}

	// Based on the model type, call the appropriate handler
	var stored bool
	switch modelBeingUpdated.(type) {
	case models.PaymentRequest:
		stored, err = paymentRequestEventHandler(event)
	default:
		event.logger.Error("event.NotificationEventHandler: Unknown logical object being updated.")
	}

	// Log what happened.
	if err != nil {
		event.logger.Error("event.NotificationEventHandler: ", zap.Error(err))
		return err
	} else if !stored {
		event.logger.Info("event.NotificationEventHandler: No notification needed to be created.")
	} else {
		event.logger.Info("event.NotificationEventHandler: Notification created and stored.")
	}

	return nil
}
