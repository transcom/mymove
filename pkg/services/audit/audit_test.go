package audit

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

func TestCapture(t *testing.T) {
	uuidString := "77c9922f-58c7-45cd-8c10-48f2a52bb55d"
	uuid, _ := uuid.FromString(uuidString)
	model := models.OfficeUser{
		ID:        uuid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	logger := zap.NewNop()
	session := auth.Session{}

	t.Run("failure when a non-pointer is passed in", func(t *testing.T) {
		err := Capture(model, logger, &session, "create_office_user")

		if err == nil {
			t.Error("Expected pointer error")
		}
	})
}
