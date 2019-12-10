package audit

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

func TestCapture(t *testing.T) {
	uuidString := "77c9922f-58c7-45cd-8c10-48f2a52bb55d"
	officeUserID, _ := uuid.FromString(uuidString)
	model := models.OfficeUser{
		ID:        officeUserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	logger := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		uuidString := "88c9922f-58c7-45cd-8c10-48f2a52bbabc"
		adminUserID, _ := uuid.FromString(uuidString)

		session := auth.Session{
			AdminUserID: adminUserID,
		}

		zapFields, _ := Capture(&model, logger, &session, "create_office_user")

		if assert.NotEmpty(t, zapFields) {
			assert.Equal(t, "record_id", zapFields[0].Key)
		}
	})

	t.Run("service member session should not include names", func(t *testing.T) {
		uuidString := "88c9922f-58c7-45cd-8c10-48f2a52bbabc"
		serviceMemberID, _ := uuid.FromString(uuidString)

		session := auth.Session{
			ServiceMemberID: serviceMemberID,
		}

		zapFields, _ := Capture(&model, logger, &session, "create_office_user")

		if assert.NotEmpty(t, zapFields) {
			var keys []string
			for _, field := range zapFields {
				keys = append(keys, field.Key)
			}

			assert.NotContains(t, "responsible_user_name", keys)
		}
	})

	t.Run("failure when a non-pointer is passed in", func(t *testing.T) {
		session := auth.Session{}
		_, err := Capture(model, logger, &session, "create_office_user")

		if err == nil {
			t.Error("Expected pointer error")
		}
	})
}
