package events

import (
	"net/http"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// EventType is the name of the audit events we care about
type EventType struct {
	Object        string
	Action        string
	ModelInstance interface{}
}

// Todo switch to an event key type as well

// PaymentRequestCreateEvent is an event
var PaymentRequestCreateEvent = EventType{"paymentRequest", "create", models.PaymentRequest{}}

// PaymentRequestUpdateEvent is an event
var PaymentRequestUpdateEvent = EventType{"paymentRequest", "update", models.PaymentRequest{}}

// MoveTaskOrderCreateEvent is an event
var MoveTaskOrderCreateEvent = EventType{"moveTaskOrder", "update", models.MoveTaskOrder{}}

// String is a function to convert to string
func (e EventType) String() string {
	return e.Object + "." + e.Action
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

// Event is the type that holds a single event
type Event struct {
	EventType       EventType
	Request         *http.Request
	MtoID           uuid.UUID
	UpdatedObjectID uuid.UUID
	EndpointKey     EndpointKeyType
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

// // Payload returns a function that sets the Auditor's payload field
// func Payload(payload interface{}) func(*Auditor) error {
// 	return func(a *Auditor) error {
// 		a.payload = payload
// 		return nil
// 	}
// }

// // Model returns a function that sets the Auditor's model field
// func Model(model interface{}) func(*Auditor) error {
// 	return func(a *Auditor) error {
// 		a.model = model
// 		return nil
// 	}
// }

// EventHandlerFunc is a type of func that can be registered as an event handler
// to be called by the eventing system
type EventHandlerFunc func(event *Event, db *pop.Connection, logger handlers.Logger) error

var registeredEventHandlers = []EventHandlerFunc{
	EventNotificationsHandler,
}

// // setRequestContext adds the request to the Auditor struct for later use
// func (a *Auditor) setRequestContext(request *http.Request) error {
// 	a.request = request

// 	clientCert := authentication.ClientCertFromRequestContext(request)

// 	// Request is coming from the Prime
// 	if clientCert != nil {
// 		a.clientCert = clientCert
// 		a.logger = a.hctx.LoggerFromRequest(request)
// 	} else {
// 		a.session, a.logger = a.hctx.SessionAndLoggerFromRequest(request)
// 	}

// 	return nil
// }

// EventRecord creates an event recording
func (e *EventGenerator) EventRecord(event Event) (*Event, error) {

	// a.setRequestContext(request)
	// for _, option := range options {
	// 	err := option(a)

	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// modelMap := slices.Map{}
	// // payloadMap := slices.Map{}
	// model := a.model
	// // payload := a.payload

	// val := reflect.ValueOf(model)

	// if model != nil {
	// 	var modelValue reflect.Value
	// 	if val.Kind() == reflect.Ptr {
	// 		modelValue = reflect.ValueOf(model).Elem()
	// 	}

	// 	if val.Kind() == reflect.Struct {
	// 		modelValue = val
	// 	}

	// 	for i := 0; i < modelValue.NumField(); i++ {
	// 		fieldFromType := modelValue.Type().Field(i)
	// 		fieldFromValue := modelValue.Field(i)
	// 		fieldName := flect.Underscore(fieldFromType.Name)
	// 		_, ok := fieldFromType.Tag.Lookup("db")

	// 		if !ok || fieldFromValue.IsZero() {
	// 			continue
	// 		}

	// 		modelMap[fieldName] = fieldFromValue.Interface()
	// 	}
	// }

	// metadata := slices.Map{
	// 	"milmove_trace_id": a.hctx.GetTraceID(),
	// }

	// // Tie to MTO if there is a MTO ID on the model
	// // Add friendly shipment identifier if model is shipment

	clientCert := authentication.ClientCertFromRequestContext(event.Request)

	// Request is coming from the Prime
	var logger handlers.Logger
	if clientCert != nil {
		logger = e.hctx.LoggerFromRequest(event.Request)
	} else {
		_, logger = e.hctx.SessionAndLoggerFromRequest(event.Request)
		//add back session
	}

	// auditRecording := models.AuditRecording{
	// 	EventName:  "MYNAME",
	// 	RecordType: "type",
	// 	RecordData: modelMap,
	// 	Metadata:   metadata,
	// }

	// if a.session == nil {
	// 	if a.clientCert != nil {
	// 		auditRecording.ClientCertID = &a.clientCert.ID
	// 	}
	// } else {
	// 	auditRecording.UserID = &a.session.UserID
	// 	auditRecording.FirstName = &a.session.FirstName
	// 	auditRecording.LastName = &a.session.LastName
	// 	auditRecording.Email = &a.session.Email
	// }

	// _, err := a.builder.CreateOne(&auditRecording)

	// if err != nil {
	// 	return nil, err
	// }

	// if this is a security event
	//   log to cloudwatch with the logic from the Capture function below
	//LogIt(&event, logger)
	for i := 0; i < len(registeredEventHandlers); i++ {
		registeredEventHandlers[i](&event, e.db, logger)
	}
	return &event, nil
}
