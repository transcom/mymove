package event

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
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
	if apiName != nil && *apiName == PrimeAPIName {
		return true
	}
	return false
}

// notificationSave saves the record in the webhook_notification table
// If it fails, it will return an error.
func notificationSave(event *Event, payload *[]byte) error {
	payloadString := string(*payload)
	newNotification := models.WebhookNotification{
		EventKey:        string(event.EventKey),
		MoveTaskOrderID: &event.MtoID,
		ObjectID:        &event.UpdatedObjectID,
		Payload:         payloadString,
		Status:          models.WebhookNotificationPending,
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

// assembleMTOShipmentPayload assembles the MTOShipment Payload and returns the JSON in bytes
func assembleMTOShipmentPayload(db *pop.Connection, updatedObjectID uuid.UUID) ([]byte, error) {
	model := models.MTOShipment{}

	// Important to be specific about which addl associations to load to reduce DB hits
	err := db.Eager("PickupAddress", "DestinationAddress",
		"SecondaryPickupAddress", "SecondaryDeliveryAddress",
		"MTOAgents").Find(&model, updatedObjectID.String())

	if err != nil {
		notFoundError := services.NewNotFoundError(updatedObjectID, "looking for MTOShipment")
		notFoundError.Wrap(err)
		return nil, notFoundError
	}

	payload := payloads.MTOShipment(&model)
	payloadArray, err := json.Marshal(payload)
	if err != nil {
		unknownErr := services.NewEventError("Unknown error creating MTOShipment payload.", err)
		return nil, unknownErr
	}
	return payloadArray, nil

}

// assembleMTOPayload assembles the MoveTaskOrder Payload and returns the JSON in bytes
func assembleMTOPayload(db *pop.Connection, updatedObjectID uuid.UUID) ([]byte, error) {
	model := models.Move{}
	// If using eager, important to be specific about which addl associations to load to reduce DB hits
	err := db.Find(&model, updatedObjectID)

	if err != nil {
		notFoundError := services.NewNotFoundError(updatedObjectID, "looking for MoveTaskOrder")
		notFoundError.Wrap(err)
		return nil, notFoundError
	}

	payload := MoveTaskOrderModelToPayload(&model)
	payloadArray, err := json.Marshal(payload)
	if err != nil {
		unknownErr := services.NewEventError("Unknown error creating MoveTaskOrder payload", err)
		return nil, unknownErr
	}

	return payloadArray, nil
}

// assembleMTOServiceItemPayload assembles the MTOServiceItem Payload and returns the JSON in bytes
func assembleMTOServiceItemPayload(db *pop.Connection, updatedObjectID uuid.UUID) ([]byte, error) {
	model := models.MTOServiceItem{}
	// Important to be specific about which addl associations to load to reduce DB hits
	err := db.Eager("ReService", "Dimensions", "CustomerContacts").Find(&model, updatedObjectID)

	if err != nil {
		notFoundError := services.NewNotFoundError(updatedObjectID, "looking for MTOServiceItem")
		notFoundError.Wrap(err)
		return nil, notFoundError
	}

	payload := payloads.MTOServiceItem(&model)
	payloadArray, err := json.Marshal(payload)
	if err != nil {
		unknownErr := services.NewEventError("Unknown error creating MTOServiceItem payload", err)
		return nil, unknownErr
	}

	return payloadArray, nil

}

// assemblePaymentRequestPayload assembles the payload and returns the JSON in bytes
func assemblePaymentRequestPayload(db *pop.Connection, updatedObjectID uuid.UUID) ([]byte, error) {
	// ASSEMBLE PAYLOAD
	model := models.PaymentRequest{}

	// Important to be specific about which addl associations to load to reduce DB hits
	err := db.Eager("PaymentServiceItems", "PaymentServiceItems.PaymentServiceItemParams").Find(&model, updatedObjectID.String())
	if err != nil {
		notFoundError := services.NewNotFoundError(updatedObjectID, "looking for PaymentRequest")
		notFoundError.Wrap(err)
		return nil, notFoundError
	}
	payload := PaymentRequestModelToPayload(&model)
	payloadArray, err := payload.MarshalBinary()
	if err != nil {
		unknownErr := services.NewEventError("Unknown error creating payload", err)
		return nil, unknownErr
	}
	return payloadArray, nil

}

// assembleOrderPayload assembles the Order Payload and returns the JSON in bytes
func assembleOrderPayload(db *pop.Connection, updatedObjectID uuid.UUID) ([]byte, error) {
	model := models.Order{}
	// Important to be specific about which addl associations to load to reduce DB hits
	err := db.Eager(
		"ServiceMember", "Entitlement", "OriginDutyStation", "NewDutyStation.Address").Find(&model, updatedObjectID)

	// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
	// cannot eager load the address as "OriginDutyStation.Address" because
	// OriginDutyStation is a pointer.
	if model.OriginDutyStation != nil {
		err = db.Load(model.OriginDutyStation, "Address")
	}

	if err != nil {
		notFoundError := services.NewNotFoundError(updatedObjectID, "looking for Order")
		notFoundError.Wrap(err)
		return nil, notFoundError
	}

	payload := payloads.Order(&model)
	payloadArray, err := json.Marshal(payload)
	if err != nil {
		unknownErr := services.NewEventError("Unknown error creating payload", err)
		return nil, unknownErr
	}

	return payloadArray, nil
}

// objectEventHandler is the default handler. It checks the source of the event and
// whether the event is available to Prime. If it determines this is a notification that we should
// send prime, it calls the appropriate function to assemble the json payload
// and stores the notification in the db.
// Returns bool indicating whether notification was stored, and error if there was one
// encountered.
func objectEventHandler(event *Event, modelBeingUpdated interface{}) (bool, error) {
	db := event.DBConnection

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

	// Based on the type of the event, call the appropriate handler to assemble the json payload
	var payloadArray []byte
	var err error

	switch modelBeingUpdated.(type) {
	case models.PaymentRequest:
		payloadArray, err = assemblePaymentRequestPayload(db, event.UpdatedObjectID)
	case models.MTOShipment:
		payloadArray, err = assembleMTOShipmentPayload(db, event.UpdatedObjectID)
	case models.MTOServiceItem:
		payloadArray, err = assembleMTOServiceItemPayload(db, event.UpdatedObjectID)
	case models.Move:
		payloadArray, err = assembleMTOPayload(db, event.UpdatedObjectID)
	default:
		event.logger.Error("event.NotificationEventHandler: Unknown logical object being updated.")
		err = services.NewEventError(fmt.Sprintf("No notification handler for event %s", event.EventKey), nil)
	}
	if err != nil {
		return false, err
	}

	// STORE NOTIFICATION IN DB
	err = notificationSave(event, &payloadArray)
	if err != nil {
		unknownErr := services.NewEventError("Unknown error storing notification", err)
		return false, unknownErr
	}

	return true, nil
}

// The purpose of this function is to handle order specific events.

func orderEventHandler(event *Event, modelBeingUpdated interface{}) (bool, error) {
	db := event.DBConnection
	// CHECK SOURCE
	// Continue only if source of event is not Prime
	if isSourcePrime(event) {
		return false, nil
	}

	// CHECK IF MOVE ID IS NIL
	// If moveID (mto ID) is nil, then return false, nil
	if event.MtoID == uuid.Nil {
		return false, nil
	}

	// CHECK FOR AVAILABILITY TO PRIME
	// Continue only if MTO is available to Prime
	if isAvailableToPrime, _ := checkAvailabilityToPrime(event); !isAvailableToPrime {
		return false, nil
	}

	// case models.Order:
	var payloadArray []byte
	var err error
	payloadArray, _ = assembleOrderPayload(db, event.UpdatedObjectID)

	// STORE NOTIFICATION IN DB
	err = notificationSave(event, &payloadArray)
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

	// Call the object specific handler if it exists else call the default
	stored := false
	switch modelBeingUpdated.(type) {
	case models.Order:
		stored, err = orderEventHandler(event, modelBeingUpdated)
	default:
		stored, err = objectEventHandler(event, modelBeingUpdated)
	}

	// Log what happened.
	if err != nil {
		event.logger.Error("event.NotificationEventHandler: ", zap.Error(err))
		return err
	} else if !stored {
		event.logger.Info("event.NotificationEventHandler: No notification needed to be created.")
	} else {
		msg := fmt.Sprintf("event.NotificationEventHandler: Notification stored for %s event triggered by %s endpoint.", event.EventKey, event.EndpointKey)
		event.logger.Info(msg)
	}

	return nil
}
