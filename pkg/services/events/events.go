package events

import (
	"net/http"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// EventKeyType is a string representing the event
// An event is generally a Object.Action event
// You would use the key in an Event object to trigger an event
type EventKeyType string

// eventModel is stored in the map of key, values
// It contains info like the model type of the object associated with this event
type eventModel struct {
	EventKey      EventKeyType
	ModelInstance interface{}
}

// Event holds a single event
// It is passed to EventRecord to trigger an event
type Event struct {
	EventKey        EventKeyType    // Pick from a select list of predefined events (PaymentRequest.Create)
	Request         *http.Request   // We expect to get this from the handler
	MtoID           uuid.UUID       // This is the ID of the MTO that the object is associated with
	UpdatedObjectID uuid.UUID       // This is the ID of the object itself (PaymentRequest.ID)
	EndpointKey     EndpointKeyType // Pick from a select list of endpoints
	logger          handlers.Logger
	session         *auth.Session
	clientCert      *models.ClientCert
	db              *pop.Connection
	hctx            handlers.HandlerContext
}

// PaymentRequestCreateEventKey is a key containing PaymentRequest.Create
const PaymentRequestCreateEventKey EventKeyType = "PaymentRequest.Create"

// PaymentRequestUpdateEventKey is a key containing PaymentRequest.Update
const PaymentRequestUpdateEventKey EventKeyType = "PaymentRequest.Update"

var eventModels map[EventKeyType]eventModel = map[EventKeyType]eventModel{
	PaymentRequestCreateEventKey: {PaymentRequestCreateEventKey, models.PaymentRequest{}},
	PaymentRequestUpdateEventKey: {PaymentRequestUpdateEventKey, models.PaymentRequest{}},
}

// IsCreateEvent returns true if this event is a create event
func IsCreateEvent(e EventKeyType) (bool, error) {
	s := strings.Split(string(e), ".")
	//TODO return error
	if s[1] == "Create" {
		return true, nil
	}
	return false, nil
}

// GetModelFromEvent returns a model instance associated with this event
func GetModelFromEvent(e EventKeyType) (interface{}, error) {
	// TODO return error
	return eventModels[e].ModelInstance, nil
}

// Auditor holds on to contextual information we need to create an AuditRecording
// type Event struct {
// 	EventType  type
// 	hctx       handlers.HandlerContext
// 	logger     Logger
// 	session    *auth.Session
// 	clientCert *models.ClientCert
// 	request    *http.Request
// 	model      interface{}
// 	payload    interface{}
// }

// EventHandlerFunc is a type of func that can be registered as an event handler
// to be called by the eventing system
type EventHandlerFunc func(event *Event) error

// registeredEventHandlers are the handlers that will be run on each event
var registeredEventHandlers = []EventHandlerFunc{
	EventNotificationsHandler,
}

//EventGenerator is the service object to generate events
type EventGenerator struct {
	db   *pop.Connection
	hctx handlers.HandlerContext
}

// NewEventGenerator instantiates a new EventGenerator
func NewEventGenerator(db *pop.Connection, hctx handlers.HandlerContext) EventGenerator {
	return EventGenerator{
		db:   db,
		hctx: hctx,
	}
}

// EventRecord creates an event recording
func (e *EventGenerator) EventRecord(event Event) (*Event, error) {

	event.clientCert = authentication.ClientCertFromRequestContext(event.Request)

	// Get logger info
	if event.clientCert != nil {
		event.logger = e.hctx.LoggerFromRequest(event.Request)
	} else {
		event.session, event.logger = e.hctx.SessionAndLoggerFromRequest(event.Request)
		//add back session
	}

	event.db = e.db
	event.hctx = e.hctx

	// Call each registered event handler with the event info and context
	for i := 0; i < len(registeredEventHandlers); i++ {
		registeredEventHandlers[i](&event)
	}
	return &event, nil
}
