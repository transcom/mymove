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
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

// event name constants

// Capture (name, slice, logger, session, request) (models.AuditRecording, error)
// Capture captures an audit record
func Capture(model interface{}, payload interface{}, logger Logger, session *auth.Session, request *http.Request, db *pop.Connection) ([]zap.Field, error) {
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
		modelMap := slices.Map{}
		payloadMap := slices.Map{}
		metadata := slices.Map{
			"weight": 90,
		}

		modelValue := elem
		for i := 0; i < modelValue.NumField(); i++ {
			fieldFromType := modelValue.Type().Field(i)
			fieldFromValue := modelValue.Field(i)
			fieldName := flect.Underscore(fieldFromType.Name)

			if !fieldFromValue.IsZero() {
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

		fmt.Println("=============================")
		fmt.Println("=============================")
		fmt.Println("=============================")
		fmt.Println("=============================")
		fmt.Printf("%#v\n", modelMap)
		fmt.Printf("%#v\n", payloadMap)
		fmt.Println("=============================")
		fmt.Println("=============================")
		fmt.Println("=============================")
		fmt.Println("=============================")

		auditRecording := models.AuditRecording{
			Name:       "MYNAME",
			RecordType: "type",
			RecordData: modelMap,
			Payload:    payloadMap,
			Metadata:   metadata,
			UserID:     &session.UserID,
		}

		builder := query.NewQueryBuilder(db)
		_, err := builder.CreateOne(&auditRecording)

		if err != nil {
			panic(err)
		}

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
	apiRegex := regexp.MustCompile("\\/[a-zA-Z]+\\/v1")
	uuidRegex := regexp.MustCompile("/([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}){1}") // https://adamscheller.com/regular-expressions/uuid-regex/
	cleanPath := uuidRegex.ReplaceAllString(apiRegex.ReplaceAllString(path, ""), "")
	return fmt.Sprintf("audit_%s_%s", strings.ToLower(request.Method), flect.Underscore(cleanPath))
}
