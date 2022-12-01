package audit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/audit"
	"github.com/transcom/mymove/pkg/auth"
)

// Capture captures an audit record
func Capture(appCtx appcontext.AppContext, model interface{}, payload interface{}, request *http.Request) ([]zap.Field, error) {
	var logItems []zap.Field
	eventType := extractEventType(request)
	msg := flect.Titleize(eventType)

	logItems = append(logItems, zap.String("event_type", eventType))
	logItems = append(logItems, extractAuditUser(request)...)

	item, err := validateInterface(model)
	if err == nil && reflect.ValueOf(model).IsValid() && !reflect.ValueOf(model).IsNil() && !reflect.ValueOf(model).IsZero() {
		logItems = append(logItems, extractRecordInformation(item, model)...)

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

			appCtx.Logger().Info("Audit patch payload", zap.String("patch_payload", string(payloadJSON)))
		}
	} else {
		msg += " invalid or zero or nil model interface received from request handler"
		logItems = append(logItems,
			zap.Error(err),
		)
	}

	appCtx.Logger().Info(msg, logItems...)

	return logItems, nil
}

// CaptureAccountStatus captures an audit record when a user account is enabled or disabled
func CaptureAccountStatus(appCtx appcontext.AppContext, model interface{}, activeValue bool, request *http.Request) ([]zap.Field, error) {
	var logItems []zap.Field
	eventType := extractEventType(request) + "_active_status_changed"
	msg := flect.Titleize(eventType)
	logItems = append(logItems, zap.String("event_type", eventType))
	logItems = append(logItems, extractResponsibleUser(appCtx.Session())...)

	item, err := validateInterface(model)
	if err == nil && reflect.ValueOf(model).IsValid() && !reflect.ValueOf(model).IsNil() && !reflect.ValueOf(model).IsZero() {
		logItems = append(logItems, extractRecordInformation(item, model)...)

		// Create log message and view value of active
		activeMessage := "disabled"
		if activeValue {
			activeMessage = "enabled"
		}

		logItems = append(logItems, zap.String("active_value", strconv.FormatBool(activeValue)))
		msg += fmt.Sprintf(" - account %s", activeMessage)
	} else {
		msg += " invalid or zero or nil model interface received from request handler"
		logItems = append(logItems,
			zap.Error(err),
		)
	}

	appCtx.Logger().Info(msg, logItems...)
	return logItems, nil
}

func extractAuditUser(request *http.Request) []zap.Field {
	var logItems []zap.Field
	auditUser := audit.RetrieveAuditUserFromContext(request.Context())
	logItems = append(logItems,
		zap.String("audit_user_id", auditUser.ID.String()),
		zap.String("audit_user_email", auditUser.LoginGovEmail),
		zap.String("BLAHBLAHBLAH", "BLAH"),
	)
	return logItems
}

func extractResponsibleUser(session *auth.Session) []zap.Field {
	var logItems []zap.Field
	logItems = append(logItems,
		zap.String("responsible_user_id", session.UserID.String()),
		zap.String("responsible_user_email", session.Email),
	)

	if session.IsAdminUser() || session.IsOfficeUser() {
		logItems = append(logItems,
			zap.String("responsible_user_name", fullName(session.FirstName, session.LastName)),
		)
	}
	return logItems
}

func extractRecordInformation(item reflect.Type, model interface{}) []zap.Field {
	var logItems []zap.Field
	recordType := parseRecordType(item.String())
	elem := reflect.ValueOf(model).Elem()

	var createdAt string
	if elem.FieldByName("CreatedAt").IsValid() {
		createdAt = elem.FieldByName("CreatedAt").Interface().(time.Time).String()
	} else {
		createdAt = time.Now().String()
	}

	var updatedAt string
	if elem.FieldByName("UpdatedAt").IsValid() {
		updatedAt = elem.FieldByName("UpdatedAt").Interface().(time.Time).String()
	} else {
		updatedAt = time.Now().String()
	}

	var id string
	if elem.FieldByName("ID").IsValid() {
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

	return logItems
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
	if t.Kind() != reflect.Ptr {
		return nil, errors.New("must pass a pointer to a struct")
	}

	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil, errors.New("must pass a pointer to a struct")
	}

	return t, nil
}

func extractEventType(request *http.Request) string {
	path := request.URL.Path
	apiRegex := regexp.MustCompile(`\/[a-zA-Z]+\/v1`)
	uuidRegex := regexp.MustCompile("/([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}){1}") // https://adamscheller.com/regular-expressions/uuid-regex/
	cleanPath := uuidRegex.ReplaceAllString(apiRegex.ReplaceAllString(path, ""), "")
	return fmt.Sprintf("audit_%s_%s", strings.ToLower(request.Method), flect.Underscore(cleanPath))
}
