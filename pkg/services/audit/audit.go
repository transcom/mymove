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

	"github.com/transcom/mymove/pkg/auth"
)

// Capture captures an audit record
func Capture(model interface{}, payload interface{}, logger Logger, session *auth.Session, request *http.Request) ([]zap.Field, error) {
	var logItems []zap.Field
	eventType := extractEventType(request)
	msg := flect.Titleize(eventType)

	logItems = append(logItems, zap.String("event_type", eventType))
	logItems = extractResponsibleUser(logItems, session)

	item, err := validateInterface(model)
	if err == nil && reflect.ValueOf(model).IsValid() && !reflect.ValueOf(model).IsNil() && !reflect.ValueOf(model).IsZero() {
		logItems = extractRecordInformation(item, model, logItems)

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

// CaptureAccountStatus captures an audit record when a user account is enabled or disabled
func CaptureAccountStatus(model interface{}, payload interface{}, logger Logger, session *auth.Session, request *http.Request) ([]zap.Field, error) {
	var logItems []zap.Field
	eventType := extractEventType(request)
	msg := flect.Titleize(eventType)
	logItems = append(logItems, zap.String("event_type", eventType))

	logItems = extractResponsibleUser(logItems, session)

	item, err := validateInterface(model)
	if err == nil && reflect.ValueOf(model).IsValid() && !reflect.ValueOf(model).IsNil() && !reflect.ValueOf(model).IsZero() {
		logItems = extractRecordInformation(item, model, logItems)

		// Create log message and view value of active
		elem := reflect.ValueOf(model).Elem()
		var activeValue bool
		if elem.FieldByName("Active").IsValid() {
			activeValue = elem.FieldByName("Active").Bool()
		} else {
			msg += " invalid model interface received from request handler - active not in model"
			logItems = append(logItems,
				zap.Error(err),
			)
		}

		activeMessage := "disabled"
		if activeValue {
			activeMessage = "enabled"
		}

		logItems = append(logItems, zap.String("active_value", strconv.FormatBool(activeValue)))

		msg += fmt.Sprintf(" - account %s 🎉🍑", activeMessage)
	} else {
		msg += " invalid or zero or nil model interface received from request handler"
		logItems = append(logItems,
			zap.Error(err),
		)
	}

	logger.Info(msg, logItems...)

	return logItems, nil
}

func extractResponsibleUser(logItems []zap.Field, session *auth.Session) []zap.Field {
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

func extractRecordInformation(item reflect.Type, model interface{}, logItems []zap.Field) []zap.Field {
	recordType := parseRecordType(item.String())
	elem := reflect.ValueOf(model).Elem()

	var createdAt string
	if elem.FieldByName("CreatedAt").IsValid() {
		createdAt = elem.FieldByName("CreatedAt").Interface().(time.Time).String()
	} else {
		createdAt = time.Now().String()
	}

	var updatedAt string
	if elem.FieldByName("updatedAt").IsValid() {
		updatedAt = elem.FieldByName("updatedAt").Interface().(time.Time).String()
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
