package audit

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
)

func Capture(model interface{}, payload interface{}, logger Logger, session *auth.Session, eventType string) ([]zap.Field, error) {
	msg := flect.Titleize(eventType)

	t, err := validateInterface(model)
	if err != nil {
		return nil, err
	}

	recordType := parseRecordType(t.String())
	elem := reflect.ValueOf(model).Elem()
	createdAt := elem.FieldByName("CreatedAt").Interface().(time.Time).String()
	updatedAt := elem.FieldByName("UpdatedAt").Interface().(time.Time).String()
	uuid := elem.FieldByName("ID").Interface().(uuid.UUID).String()

	logItems := []zap.Field{
		zap.String("record_id", uuid),
		zap.String("record_type", recordType),
		zap.String("record_created_at", createdAt),
		zap.String("record_updated_at", updatedAt),
		zap.String("responsible_user_id", session.UserID.String()),
		zap.String("responsible_user_email", session.Email),
		zap.String("event_type", eventType),
	}

	if session.IsAdminUser() || session.IsOfficeUser() {
		logItems = append(logItems,
			zap.String("responsible_user_name", fullName(session.FirstName, session.LastName)),
		)
	}

	if !(payload == nil) {
		_, err := validateInterface(payload)
		if err != nil {
			return nil, err
		}

		patchPayload, _ := json.Marshal(payload)

		logItems = append(logItems, zap.String("patch_payload", string(patchPayload)))
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
