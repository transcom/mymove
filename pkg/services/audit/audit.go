package audit

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/gobuffalo/flect"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
)

func Capture(model interface{}, logger Logger, session *auth.Session, eventType string) error {
	msg := flect.Titleize(eventType)
	t := reflect.TypeOf(model)

	if t.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer to a record")
	}

	recordType := parseRecordType(t.Elem().String())
	elem := reflect.ValueOf(model).Elem()
	createdAt := elem.FieldByName("CreatedAt").Interface().(time.Time).String()
	updatedAt := elem.FieldByName("UpdatedAt").Interface().(time.Time).String()
	uuid := elem.FieldByName("ID").Interface().(uuid.UUID).String()

	logger.Info(msg,
		zap.String("record_id", uuid),
		zap.String("record_type", recordType),
		zap.String("record_created_at", createdAt),
		zap.String("record_updated_at", updatedAt),
		zap.String("responsible_user_id", session.UserID.String()),
		zap.String("responsible_user_email", session.Email),
		zap.String("responsible_user_name", fullName(session.FirstName, session.LastName)),
		zap.String("event_type", eventType))

	return nil
}

func parseRecordType(rt string) string {
	parts := strings.Split(rt, ".")

	return parts[1]
}

func fullName(first, last string) string {
	return first + " " + last
}
