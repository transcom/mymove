package event

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
)

// KeyType is a string representing the event
// An event is generally a Object.Action event
// You would use the key in an Event object to trigger an event
type KeyType string

// eventModel is stored in the map of key, values
// It contains info like the model type of the object associated with this event
type eventModel struct {
	EventKey      KeyType
	ModelInstance interface{}
}

// Event holds a single event
// It is passed to EventTrigger to trigger an event
type Event struct {
	EventKey        KeyType               // Pick from a select list of predefined events (PaymentRequest.Create)
	MtoID           uuid.UUID             // This is the ID of the MTO that the object is associated with
	UpdatedObjectID uuid.UUID             // This is the ID of the object itself (PaymentRequest.ID)
	EndpointKey     EndpointKeyType       // Pick from a select list of endpoints
	AppContext      appcontext.AppContext // The AppContext from the Request
	TraceID         uuid.UUID             // the TraceID from the Request
}

// OrderUpdateEventKey is a key containing Order.Update
const OrderUpdateEventKey KeyType = "Order.Update"

// MoveTaskOrderCreateEventKey is a key containing MoveTaskOrder.Create
const MoveTaskOrderCreateEventKey KeyType = "MoveTaskOrder.Create"

// MoveTaskOrderUpdateEventKey is a key containing MoveTaskOrder.Update
const MoveTaskOrderUpdateEventKey KeyType = "MoveTaskOrder.Update"

// MTOShipmentCreateEventKey is a key containing MTOShipment.Create
const MTOShipmentCreateEventKey KeyType = "MTOShipment.Create"

// MTOShipmentUpdateEventKey is a key containing MTOShipment.Update
const MTOShipmentUpdateEventKey KeyType = "MTOShipment.Update"

// ShipmentDeleteEventKey is a key containing Shipment.Delete
const ShipmentDeleteEventKey KeyType = "Shipment.Delete"

// ShipmentApproveEventKey is a key containing Shipment.Approve
const ShipmentApproveEventKey KeyType = "Shipment.Approve"

// ShipmentRequestDiversionEventKey is a key containing Shipment.RequestDiversion
const ShipmentRequestDiversionEventKey KeyType = "Shipment.RequestDiversion"

// ShipmentApproveDiversionEventKey is a key containing Shipment.ApproveDiversion
const ShipmentApproveDiversionEventKey KeyType = "Shipment.ApproveDiversion"

// ShipmentRejectEventKey is a key containing Shipment.Reject
const ShipmentRejectEventKey KeyType = "Shipment.Reject"

// ShipmentRequestCancellationEventKey is a key containing Shipment.RequestCancellation
const ShipmentRequestCancellationEventKey KeyType = "Shipment.RequestCancellation"

// ShipmentRequestReweighEventKey is a key containing Shipment.RequestReweigh
const ShipmentRequestReweighEventKey KeyType = "Shipment.RequestReweigh"

// ApproveSITExtensionEventKey is a key containing Shipment.ApproveSIT
const ApproveSITExtensionEventKey KeyType = "Shipment.ApproveSITExtension"

// DenySITExtensionEventKey is a key containing Shipment.DenySIT
const DenySITExtensionEventKey KeyType = "Shipment.DenySITExtension"

// MTOServiceItemCreateEventKey is a key containing MTOServiceItem.Create
const MTOServiceItemCreateEventKey KeyType = "MTOServiceItem.Create"

// MTOServiceItemUpdateEventKey is a key containing MTOServiceItem.Update
const MTOServiceItemUpdateEventKey KeyType = "MTOServiceItem.Update"

// PaymentRequestCreateEventKey is a key containing PaymentRequest.Create
const PaymentRequestCreateEventKey KeyType = "PaymentRequest.Create"

// PaymentRequestUpdateEventKey is a key containing PaymentRequest.Update
const PaymentRequestUpdateEventKey KeyType = "PaymentRequest.Update"

// CustomerSupportRemarkDeleteEventKey is a key containing CustomerSupportRemark.Delete
const CustomerSupportRemarkDeleteEventKey KeyType = "CustomerSupportRemark.Delete"

// TestCreateEventKey is a key containing Test.Create
const TestCreateEventKey KeyType = "Test.Create"

// TestUpdateEventKey is a key containing Test.Update
const TestUpdateEventKey KeyType = "Test.Update"

// TestDeleteEventKey is a key containing Test.Delete
const TestDeleteEventKey KeyType = "Test.Delete"

var eventModels = map[KeyType]eventModel{
	OrderUpdateEventKey:                 {OrderUpdateEventKey, models.Order{}},
	MoveTaskOrderCreateEventKey:         {MoveTaskOrderCreateEventKey, models.Move{}},
	MoveTaskOrderUpdateEventKey:         {MoveTaskOrderUpdateEventKey, models.Move{}},
	MTOShipmentCreateEventKey:           {MTOShipmentCreateEventKey, models.MTOShipment{}},
	MTOShipmentUpdateEventKey:           {MTOShipmentUpdateEventKey, models.MTOShipment{}},
	ShipmentDeleteEventKey:              {ShipmentDeleteEventKey, models.MTOShipment{}},
	ShipmentApproveEventKey:             {ShipmentApproveEventKey, models.MTOShipment{}},
	ShipmentRequestDiversionEventKey:    {ShipmentRequestDiversionEventKey, models.MTOShipment{}},
	ShipmentApproveDiversionEventKey:    {ShipmentApproveDiversionEventKey, models.MTOShipment{}},
	ShipmentRejectEventKey:              {ShipmentRejectEventKey, models.MTOShipment{}},
	ShipmentRequestCancellationEventKey: {ShipmentRequestCancellationEventKey, models.MTOShipment{}},
	ShipmentRequestReweighEventKey:      {ShipmentRequestReweighEventKey, models.MTOShipment{}},
	ApproveSITExtensionEventKey:         {ApproveSITExtensionEventKey, models.MTOShipment{}},
	DenySITExtensionEventKey:            {DenySITExtensionEventKey, models.MTOShipment{}},
	MTOServiceItemCreateEventKey:        {MTOServiceItemCreateEventKey, models.MTOServiceItem{}},
	MTOServiceItemUpdateEventKey:        {MTOServiceItemUpdateEventKey, models.MTOServiceItem{}},
	PaymentRequestCreateEventKey:        {PaymentRequestCreateEventKey, models.PaymentRequest{}},
	PaymentRequestUpdateEventKey:        {PaymentRequestUpdateEventKey, models.PaymentRequest{}},
	CustomerSupportRemarkDeleteEventKey: {CustomerSupportRemarkDeleteEventKey, models.CustomerSupportRemark{}},
	TestCreateEventKey:                  {TestCreateEventKey, nil},
	TestUpdateEventKey:                  {TestUpdateEventKey, nil},
	TestDeleteEventKey:                  {TestDeleteEventKey, nil}}

// IsCreateEvent returns true if this event is a create event
func IsCreateEvent(e KeyType) (bool, error) {
	s := strings.Split(string(e), ".")
	if len(s) != 2 {
		err := apperror.NewEventError(fmt.Sprintf("Event Key %s is malformed. Should be of form Object.Action.", e), nil)
		return false, err
	}
	if s[1] == "Create" {
		return true, nil
	}
	return false, nil
}

// GetModelFromEvent returns a model instance associated with this event
func GetModelFromEvent(e KeyType) (interface{}, error) {
	eventModel, success := eventModels[e]

	if !success {
		err := apperror.NewEventError(fmt.Sprintf("Event Key %s was not found in eventModels. Must use known event key.", e), nil)
		return nil, err
	}
	return eventModel.ModelInstance, nil
}

// ExistsEventKey returns true if the event key exists
func ExistsEventKey(e string) bool {
	_, ok := eventModels[KeyType(e)]
	return ok
}

// RegisteredEventHandlerFunc is a type of func that can be registered as an event handler
// to be called by the eventing system
type RegisteredEventHandlerFunc func(event *Event) error

// registeredEventHandlers are the handlers that will be run on each event
var registeredEventHandlers = []RegisteredEventHandlerFunc{
	NotificationEventHandler,
}

func consolidateError(errorList []error) string {
	switch len(errorList) {
	case 0:
		return "no errors"
	default:
		errMessage := ""
		for _, e := range errorList {
			errMessage += e.Error() + ". "
		}
		return errMessage
	}
}

// TriggerEvent triggers an event to send to various handlers
func TriggerEvent(event Event) (*Event, error) {

	var errorList []error
	// Check eventKey
	_, success := eventModels[event.EventKey]
	if !success {
		err := apperror.NewEventError(fmt.Sprintf("Event Key %s was not found in eventModels. Must use known event key.", event.EventKey), nil)
		return nil, err
	}
	// Check that Request and context were passed in
	if event.AppContext == nil {
		err := apperror.NewEventError("The AppContext must be passed to TriggerEvent.", nil)
		return nil, err
	}
	// Check endpointKey if exists
	if event.EndpointKey != "" {
		result := GetEndpointAPI(event.EndpointKey)
		if result == nil {
			err := apperror.NewEventError(fmt.Sprintf("Endpoint Key %s was not found in endpoints. Must use known endpoint key.", event.EndpointKey), nil)
			return nil, err
		}
	}

	// Call each registered event handler with the event info and context
	// Collect errors, this is to avoid one registered handler failure to
	// affect another.
	for i := 0; i < len(registeredEventHandlers); i++ {
		err := registeredEventHandlers[i](&event)
		if err != nil {
			errorList = append(errorList, err)
		}
	}
	if len(errorList) > 0 {
		err := apperror.NewEventError(consolidateError(errorList), nil)
		return &event, err
	}
	return &event, nil
}
