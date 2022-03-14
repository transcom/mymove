package event

import (
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
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
	appCtx := event.AppContext
	payloadString := string(*payload)
	newNotification := models.WebhookNotification{
		EventKey:        string(event.EventKey),
		MoveTaskOrderID: &event.MtoID,
		ObjectID:        &event.UpdatedObjectID,
		Payload:         payloadString,
		Status:          models.WebhookNotificationPending,
	}

	if !event.TraceID.IsNil() {
		newNotification.TraceID = &event.TraceID
	}

	// Creates the notification entry in the DB
	verrs, err := appCtx.DB().ValidateAndCreate(&newNotification)
	if verrs.Count() > 0 {
		appCtx.Logger().Error("event.notificationSave error", zap.Error(verrs))
		return apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "")
	} else if err != nil {
		appCtx.Logger().Error("event.notificationSave error", zap.Error(err))
		e := apperror.NewQueryError("Notification", err, "Unable to save Notification.")
		return e
	}
	return nil

}

// checkAvailabilityToPrime returns true if the MTO is
// available to Prime. If there is a query error, it returns an
// error as well.
func checkAvailabilityToPrime(event *Event) (bool, error) {
	appCtx := event.AppContext
	mtoChecker := movetaskorder.NewMoveTaskOrderChecker()
	availableToPrime, err := mtoChecker.MTOAvailableToPrime(appCtx, event.MtoID)
	if err != nil {
		unknownErr := apperror.NewEventError("Unknown error checking prime availability", err)
		return false, unknownErr
	}
	// No need to store notification to send if not available to prime
	if !availableToPrime {
		return false, nil
	}
	return true, nil

}

// assembleMTOShipmentPayload assembles the MTOShipment Payload and returns the JSON in bytes and a bool
// representing whether this notification should continue (we don't want to notify when the shipment
// is handled by an external vendor, for instance).
func assembleMTOShipmentPayload(appCtx appcontext.AppContext, updatedObjectID uuid.UUID) ([]byte, bool, error) {
	// First, get the MTOShipment and ensure we need to notify before loading other relationships.
	var mtoShipment models.MTOShipment
	err := appCtx.DB().Find(&mtoShipment, updatedObjectID.String())
	if err != nil {
		notFoundError := apperror.NewNotFoundError(updatedObjectID, "looking for MTOShipment")
		notFoundError.Wrap(err)
		return nil, false, notFoundError
	}

	// If this shipment is being handled by an external vendor, don't notify the prime.
	if mtoShipment.UsesExternalVendor {
		return nil, false, nil
	}

	// Now load any additional required relationships since we now know we intend to send this notification.
	err = appCtx.DB().Load(&mtoShipment, "PickupAddress", "DestinationAddress", "SecondaryPickupAddress",
		"SecondaryDeliveryAddress", "MTOAgents", "StorageFacility")
	if err != nil {
		return nil, false, apperror.NewQueryError("MTOShipment", err, "Unable to load MTOShipment relationships")
	}

	if mtoShipment.StorageFacility != nil && uuid.Nil != mtoShipment.StorageFacility.AddressID {
		err = appCtx.DB().Load(mtoShipment.StorageFacility, "Address")
		if err != nil {
			notFoundError := apperror.NewNotFoundError(updatedObjectID, "looking for MTOShipment.StorageFacility.Address")
			notFoundError.Wrap(err)
			return nil, false, notFoundError
		}
	}

	payload := payloads.MTOShipment(&mtoShipment)
	payloadArray, err := json.Marshal(payload)
	if err != nil {
		return nil, false, apperror.NewEventError("Unknown error creating MTOShipment payload.", err)
	}

	return payloadArray, true, nil
}

// assembleMTOPayload assembles the MoveTaskOrder Payload and returns the JSON in bytes
func assembleMTOPayload(appCtx appcontext.AppContext, updatedObjectID uuid.UUID) ([]byte, error) {
	model := models.Move{}
	// If using eager, important to be specific about which addl associations to load to reduce DB hits
	err := appCtx.DB().Find(&model, updatedObjectID)

	if err != nil {
		notFoundError := apperror.NewNotFoundError(updatedObjectID, "looking for MoveTaskOrder")
		notFoundError.Wrap(err)
		return nil, notFoundError
	}

	payload := MoveTaskOrderModelToPayload(&model)
	payloadArray, err := json.Marshal(payload)
	if err != nil {
		unknownErr := apperror.NewEventError("Unknown error creating MoveTaskOrder payload", err)
		return nil, unknownErr
	}

	return payloadArray, nil
}

// assembleMTOServiceItemPayload assembles the MTOServiceItem Payload and returns the JSON in bytes
func assembleMTOServiceItemPayload(appCtx appcontext.AppContext, updatedObjectID uuid.UUID) ([]byte, error) {
	model := models.MTOServiceItem{}
	// Important to be specific about which addl associations to load to reduce DB hits
	err := appCtx.DB().Eager("ReService", "Dimensions", "CustomerContacts").Find(&model, updatedObjectID)

	if err != nil {
		notFoundError := apperror.NewNotFoundError(updatedObjectID, "looking for MTOServiceItem")
		notFoundError.Wrap(err)
		return nil, notFoundError
	}

	payload := payloads.MTOServiceItem(&model)
	payloadArray, err := json.Marshal(payload)
	if err != nil {
		unknownErr := apperror.NewEventError("Unknown error creating MTOServiceItem payload", err)
		return nil, unknownErr
	}

	return payloadArray, nil

}

// assemblePaymentRequestPayload assembles the payload and returns the JSON in bytes
func assemblePaymentRequestPayload(appCtx appcontext.AppContext, updatedObjectID uuid.UUID) ([]byte, error) {
	// ASSEMBLE PAYLOAD
	model := models.PaymentRequest{}

	// Important to be specific about which addl associations to load to reduce DB hits
	err := appCtx.DB().Eager("PaymentServiceItems", "PaymentServiceItems.PaymentServiceItemParams").Find(&model, updatedObjectID.String())
	if err != nil {
		notFoundError := apperror.NewNotFoundError(updatedObjectID, "looking for PaymentRequest")
		notFoundError.Wrap(err)
		return nil, notFoundError
	}
	payload := PaymentRequestModelToPayload(&model)
	payloadArray, err := payload.MarshalBinary()
	if err != nil {
		unknownErr := apperror.NewEventError("Unknown error creating payload", err)
		return nil, unknownErr
	}
	return payloadArray, nil

}

// assembleOrderPayload assembles the Order Payload and returns the JSON in bytes
func assembleOrderPayload(appCtx appcontext.AppContext, updatedObjectID uuid.UUID) ([]byte, error) {
	model := models.Order{}
	// Important to be specific about which addl associations to load to reduce DB hits
	err := appCtx.DB().Eager(
		"ServiceMember", "Entitlement", "OriginDutyLocation", "NewDutyLocation.Address").Find(&model, updatedObjectID)

	// Due to a bug in pop (https://github.com/gobuffalo/pop/issues/578), we
	// cannot eager load the address as "OriginDutyLocation.Address" because
	// OriginDutyLocation is a pointer.
	if model.OriginDutyLocation != nil {
		err = appCtx.DB().Load(model.OriginDutyLocation, "Address")
	}

	if err != nil {
		notFoundError := apperror.NewNotFoundError(updatedObjectID, "looking for Order")
		notFoundError.Wrap(err)
		return nil, notFoundError
	}

	payload := payloads.Order(&model)
	payloadArray, err := json.Marshal(payload)
	if err != nil {
		unknownErr := apperror.NewEventError("Unknown error creating payload", err)
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
	appCtx := event.AppContext

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
		payloadArray, err = assemblePaymentRequestPayload(appCtx, event.UpdatedObjectID)
	case models.MTOShipment:
		var shouldNotify bool
		payloadArray, shouldNotify, err = assembleMTOShipmentPayload(appCtx, event.UpdatedObjectID)
		if !shouldNotify {
			return false, err
		}
	case models.MTOServiceItem:
		payloadArray, err = assembleMTOServiceItemPayload(appCtx, event.UpdatedObjectID)
	case models.Move:
		payloadArray, err = assembleMTOPayload(appCtx, event.UpdatedObjectID)
	default:
		appCtx.Logger().Error("event.NotificationEventHandler: Unknown logical object being updated.")
		err = apperror.NewEventError(fmt.Sprintf("No notification handler for event %s", event.EventKey), nil)
	}
	if err != nil {
		return false, err
	}

	// STORE NOTIFICATION IN DB
	err = notificationSave(event, &payloadArray)
	if err != nil {
		unknownErr := apperror.NewEventError("Unknown error storing notification", err)
		return false, unknownErr
	}

	return true, nil
}

// The purpose of this function is to handle order specific events.

func orderEventHandler(event *Event, modelBeingUpdated interface{}) (bool, error) {
	appCtx := event.AppContext
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
	payloadArray, _ = assembleOrderPayload(appCtx, event.UpdatedObjectID)

	// STORE NOTIFICATION IN DB
	err = notificationSave(event, &payloadArray)
	if err != nil {
		unknownErr := apperror.NewEventError("Unknown error storing notification", err)
		return false, unknownErr
	}

	return true, nil
}

// NotificationEventHandler receives notifications from the events package
// For alerting ALL errors should be logged here.
func NotificationEventHandler(event *Event) error {
	appCtx := event.AppContext
	//Get the model type of the updated logical object
	modelBeingUpdated, err := GetModelFromEvent(event.EventKey)

	if err != nil {
		appCtx.Logger().Error("event.NotificationEventHandler: EventKey does not exist.")
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
		appCtx.Logger().Error("event.NotificationEventHandler: ", zap.Error(err))
		return err
	} else if !stored {
		appCtx.Logger().Info("event.NotificationEventHandler: No notification needed to be created.")
	} else {
		msg := fmt.Sprintf("event.NotificationEventHandler: Notification stored for %s event triggered by %s endpoint.", event.EventKey, event.EndpointKey)
		appCtx.Logger().Info(msg)
	}

	return nil
}
