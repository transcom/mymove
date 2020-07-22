package events

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/handlers"
)

// https://raw.githubusercontent.com/transcom/mymove/ba76aeae09b007219fe60e5a7002d8c0f7c14b1d/pkg/services/audit/audit.go

// EventType is the name of the audit events we care about
type EventType string

// Endpoint is the name of the api + endpoint
type Endpoint string

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
	eventType EventType
	request   *http.Request
	mtoID     uuid.UUID
	objectID  uuid.UUID
	endpoint  Endpoint
}

//EventGenerator is the service object to generate events
type EventGenerator struct {
	//	builder *query.Builder
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

// LogIt logs an event to the console
func LogIt(event *Event, logger handlers.Logger) {
	fmt.Println("\n\n Event handler ran: ", event.endpoint, event.mtoID, event.objectID)
	logger.Info("Event handler ran:", zap.String("endpoint", string(event.endpoint)))
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
func (e *EventGenerator) EventRecord(endpoint Endpoint, r *http.Request) (*Event, error) {

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

	event := Event{
		eventType: "new event",
		endpoint:  endpoint,
		request:   r,
	}

	clientCert := authentication.ClientCertFromRequestContext(event.request)

	// Request is coming from the Prime
	var logger handlers.Logger
	if clientCert != nil {
		logger = e.hctx.LoggerFromRequest(event.request)
	} else {
		_, logger = e.hctx.SessionAndLoggerFromRequest(event.request)
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
	LogIt(&event, logger)
	return &event, nil
}

// // Capture captures an audit record
// func Capture(model interface{}, payload interface{}, logger Logger, session *auth.Session, request *http.Request) ([]zap.Field, error) {
// 	// metadataColumns := map[string]map[string]string{
// 	// 	"Customer requested shipments + pick up dates": { meta: "requested_pickup_date", identifier: "pretty_shipment_id" },
// 	// }

// 	var logItems []zap.Field
// 	eventType := extractEventType(request)
// 	msg := flect.Titleize(eventType)

// 	logItems = append(logItems,
// 		zap.String("event_type", eventType),
// 		zap.String("responsible_user_id", session.UserID.String()),
// 		zap.String("responsible_user_email", session.Email),
// 	)

// 	if session.IsAdminUser() || session.IsOfficeUser() {
// 		logItems = append(logItems,
// 			zap.String("responsible_user_name", fullName(session.FirstName, session.LastName)),
// 		)
// 	}

// 	t, err := validateInterface(model)
// 	if err == nil && reflect.ValueOf(model).IsValid() == true && reflect.ValueOf(model).IsNil() == false && reflect.ValueOf(model).IsZero() == false {
// 		recordType := parseRecordType(t.String())
// 		elem := reflect.ValueOf(model).Elem()

// 		var createdAt string
// 		if elem.FieldByName("CreatedAt").IsValid() == true {
// 			createdAt = elem.FieldByName("CreatedAt").Interface().(time.Time).String()
// 		} else {
// 			createdAt = time.Now().String()
// 		}

// 		var updatedAt string
// 		if elem.FieldByName("updatedAt").IsValid() == true {
// 			updatedAt = elem.FieldByName("updatedAt").Interface().(time.Time).String()
// 		} else {
// 			updatedAt = time.Now().String()
// 		}

// 		var id string
// 		if elem.FieldByName("ID").IsValid() == true {
// 			id = elem.FieldByName("ID").Interface().(uuid.UUID).String()
// 		} else {
// 			id = ""
// 		}

// 		logItems = append(logItems,
// 			zap.String("record_id", id),
// 			zap.String("record_type", recordType),
// 			zap.String("record_created_at", createdAt),
// 			zap.String("record_updated_at", updatedAt),
// 		)

// 		if payload != nil {
// 			_, err = validateInterface(payload)
// 			if err != nil {
// 				return nil, err
// 			}

// 			var payloadFields []string
// 			payloadValue := reflect.ValueOf(payload).Elem()
// 			for i := 0; i < payloadValue.NumField(); i++ {
// 				fieldFromType := payloadValue.Type().Field(i)
// 				fieldFromValue := payloadValue.Field(i)
// 				fieldName := flect.Underscore(fieldFromType.Name)

// 				if !fieldFromValue.IsZero() {
// 					payloadFields = append(payloadFields, fieldName)
// 				}
// 			}

// 			logItems = append(logItems, zap.String("fields_changed", strings.Join(payloadFields, ",")))

// 			var payloadJSON []byte
// 			payloadJSON, err = json.Marshal(payload)

// 			if err != nil {
// 				return nil, err
// 			}

// 			logger.Debug("Audit patch payload", zap.String("patch_payload", string(payloadJSON)))
// 		}
// 	} else {
// 		msg += " invalid or zero or nil model interface received from request handler"
// 		logItems = append(logItems,
// 			zap.Error(err),
// 		)
// 	}

// 	logger.Info(msg, logItems...)

// 	return logItems, nil
// }

// func parseRecordType(rt string) string {
// 	parts := strings.Split(rt, ".")

// 	return parts[1]
// }

// func fullName(first, last string) string {
// 	return first + " " + last
// }

// func validateInterface(thing interface{}) (reflect.Type, error) {
// 	t := reflect.TypeOf(thing)

// 	if t != nil {
// 		if t.Kind() != reflect.Ptr {
// 			return nil, errors.New("must pass a pointer to a struct")
// 		}

// 		t = t.Elem()
// 		if t.Kind() != reflect.Struct {
// 			return nil, errors.New("must pass a pointer to a struct")
// 		}
// 	}

// 	return t, nil
// }

// func extractEventType(request *http.Request) string {
// 	path := request.URL.Path
// 	apiRegex := regexp.MustCompile("\\/[a-zA-Z]+\\/v1")
// 	uuidRegex := regexp.MustCompile("/([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}){1}") // https://adamscheller.com/regular-expressions/uuid-regex/
// 	cleanPath := uuidRegex.ReplaceAllString(apiRegex.ReplaceAllString(path, ""), "")
// 	return fmt.Sprintf("audit_%s_%s", strings.ToLower(request.Method), flect.Underscore(cleanPath))
// }
