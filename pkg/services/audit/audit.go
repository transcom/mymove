package audit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/pop/slices"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

// Auditor holds on to contextual information we need to create an AuditRecording
type Auditor struct {
	builder    *query.Builder
	hctx       handlers.HandlerContext
	logger     Logger
	session    *auth.Session
	clientCert *models.ClientCert
	request    *http.Request
}

// NewAuditor instantiates a new Auditor
func NewAuditor(builder *query.Builder, hctx handlers.HandlerContext) Auditor {
	return Auditor{
		builder: builder,
		hctx:    hctx,
	}
}

// SetRequestContext adds the request to the Auditor struct for later use
func (a *Auditor) SetRequestContext(request *http.Request) {
	a.request = request

	clientCert := authentication.ClientCertFromRequestContext(request)

	// Request is coming from the Prime
	if clientCert != nil {
		a.clientCert = clientCert
		a.logger = a.hctx.LoggerFromRequest(request)
	} else {
		a.session, a.logger = a.hctx.SessionAndLoggerFromRequest(request)
	}

}

// Record creates an audit recording
func (a *Auditor) Record(name string, model, payload interface{}) (*models.AuditRecording, error) {
	modelMap := slices.Map{}
	payloadMap := slices.Map{}
	val := reflect.ValueOf(model)

	if model != nil {
		var modelValue reflect.Value
		if val.Kind() == reflect.Ptr {
			modelValue = reflect.ValueOf(model).Elem()
		}

		if val.Kind() == reflect.Struct {
			modelValue = val
		}

		for i := 0; i < modelValue.NumField(); i++ {
			fieldFromType := modelValue.Type().Field(i)
			fieldFromValue := modelValue.Field(i)
			fieldName := flect.Underscore(fieldFromType.Name)
			_, ok := fieldFromType.Tag.Lookup("db")

			if !ok || fieldFromValue.IsZero() {
				continue
			}

			modelMap[fieldName] = fieldFromValue.Interface()
		}
	}

	if payload != nil {
		payloadVal := reflect.ValueOf(payload).Elem()
		for i := 0; i < payloadVal.NumField(); i++ {
			fieldFromType := payloadVal.Type().Field(i)
			fieldFromValue := payloadVal.Field(i)
			fieldName := flect.Underscore(fieldFromType.Name)

			if !fieldFromValue.IsZero() {
				payloadMap[fieldName] = fieldFromValue.Interface()
			}
		}

	}

	metadata := slices.Map{
		"milmove_trace_id": a.hctx.GetTraceID(),
	}

	// Tie to MTO if there is a MTO ID on the model
	// Add friendly shipment identifier if model is shipment

	auditRecording := models.AuditRecording{
		EventName:  "MYNAME",
		RecordType: "type",
		RecordData: modelMap,
		Payload:    payloadMap,
		Metadata:   metadata,
	}

	if a.session == nil {
		if a.clientCert != nil {
			auditRecording.ClientCertID = &a.clientCert.ID
		}
	} else {
		auditRecording.UserID = &a.session.UserID
		auditRecording.FirstName = &a.session.FirstName
		auditRecording.LastName = &a.session.LastName
		auditRecording.Email = &a.session.Email

	}

	_, err := a.builder.CreateOne(&auditRecording)

	if err != nil {
		return nil, err
	}

	// if this is a security event
	//   log to cloudwatch with the logic from the Capture function below

	return &auditRecording, nil
}

// Capture captures an audit record
func Capture(model interface{}, payload interface{}, logger Logger, session *auth.Session, request *http.Request) ([]zap.Field, error) {
	// metadataColumns := map[string]map[string]string{
	// 	"Customer requested shipments + pick up dates": { meta: "requested_pickup_date", identifier: "pretty_shipment_id" },
	// }

	var logItems []zap.Field
	eventType := extractEventType(request)
	msg := flect.Titleize(eventType)

	logItems = append(logItems,
		zap.String("event_type", eventType),
		zap.String("responsible_user_id", session.UserID.String()),
		zap.String("responsible_user_email", session.Email),
	)

	if session.IsAdminUser() || session.IsOfficeUser() {
		logItems = append(logItems,
			zap.String("responsible_user_name", fullName(session.FirstName, session.LastName)),
		)
	}

	t, err := validateInterface(model)
	if err == nil && reflect.ValueOf(model).IsValid() == true && reflect.ValueOf(model).IsNil() == false && reflect.ValueOf(model).IsZero() == false {
		recordType := parseRecordType(t.String())
		elem := reflect.ValueOf(model).Elem()

		var createdAt string
		if elem.FieldByName("CreatedAt").IsValid() == true {
			createdAt = elem.FieldByName("CreatedAt").Interface().(time.Time).String()
		} else {
			createdAt = time.Now().String()
		}

		var updatedAt string
		if elem.FieldByName("updatedAt").IsValid() == true {
			updatedAt = elem.FieldByName("updatedAt").Interface().(time.Time).String()
		} else {
			updatedAt = time.Now().String()
		}

		var id string
		if elem.FieldByName("ID").IsValid() == true {
			id = elem.FieldByName("ID").Interface().(uuid.UUID).String()
		} else {
			id = ""
		}

		logItems = append(logItems,
			zap.String("record_id", id),
			zap.String("record_type", recordType),
			zap.String("record_created_at", createdAt),
			zap.String("record_updated_at", updatedAt),
		)

		if payload != nil {
			_, err = validateInterface(payload)
			if err != nil {
				return nil, err
			}

			var payloadFields []string
			payloadValue := reflect.ValueOf(payload).Elem()
			for i := 0; i < payloadValue.NumField(); i++ {
				fieldFromType := payloadValue.Type().Field(i)
				fieldFromValue := payloadValue.Field(i)
				fieldName := flect.Underscore(fieldFromType.Name)

				if !fieldFromValue.IsZero() {
					payloadFields = append(payloadFields, fieldName)
				}
			}

			logItems = append(logItems, zap.String("fields_changed", strings.Join(payloadFields, ",")))

			var payloadJSON []byte
			payloadJSON, err = json.Marshal(payload)

			if err != nil {
				return nil, err
			}

			logger.Debug("Audit patch payload", zap.String("patch_payload", string(payloadJSON)))
		}
	} else {
		msg += " invalid or zero or nil model interface received from request handler"
		logItems = append(logItems,
			zap.Error(err),
		)
	}

	logger.Info(msg, logItems...)

	return logItems, nil
}

func parseRecordType(rt string) string {
	parts := strings.Split(rt, ".")

	return parts[1]
}

func fullName(first, last string) string {
	return first + " " + last
}

func validateInterface(thing interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(thing)

	if t != nil {
		if t.Kind() != reflect.Ptr {
			return nil, errors.New("must pass a pointer to a struct")
		}

		t = t.Elem()
		if t.Kind() != reflect.Struct {
			return nil, errors.New("must pass a pointer to a struct")
		}
	}

	return t, nil
}

func extractEventType(request *http.Request) string {
	path := request.URL.Path
	apiRegex := regexp.MustCompile("\\/[a-zA-Z]+\\/v1")
	uuidRegex := regexp.MustCompile("/([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}){1}") // https://adamscheller.com/regular-expressions/uuid-regex/
	cleanPath := uuidRegex.ReplaceAllString(apiRegex.ReplaceAllString(path, ""), "")
	return fmt.Sprintf("audit_%s_%s", strings.ToLower(request.Method), flect.Underscore(cleanPath))
}
