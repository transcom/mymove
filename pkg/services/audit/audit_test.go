package audit

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
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
	dummyRequest := http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	logger := zap.NewNop()

	t.Run("success", func(t *testing.T) {
		uuidString := "88c9922f-58c7-45cd-8c10-48f2a52bbabc"
		adminUserID, _ := uuid.FromString(uuidString)

		session := auth.Session{
			AdminUserID: adminUserID,
		}

		appCtx := appcontext.NewAppContext(nil, logger, &session)

		req := &http.Request{
			URL: &url.URL{
				Path: "/admin/v1/admin-users",
			},
			Method: "POST",
		}

		zapFields, _ := Capture(appCtx, &model, nil, req)
		var eventType string
		for _, field := range zapFields {
			if field.Key == "event_type" {
				eventType = field.String
			}
		}

		if assert.NotEmpty(t, zapFields) {
			assert.Equal(t, "event_type", zapFields[0].Key)
			assert.Equal(t, "audit_post_admin_users", eventType)
		}
	})

	t.Run("success with optional patch payload", func(t *testing.T) {
		uuidString := "88c9922f-58c7-45cd-8c10-48f2a52bbabc"
		adminUserID, _ := uuid.FromString(uuidString)

		type fakePatchPayload struct {
			Active         bool    `json:"active,omitempty"`
			FirstName      string  `json:"first_name,omitempty"`
			LastName       string  `json:"last_name,omitempty"`
			MiddleInitials *string `json:"middle_initials,omitempty"`
			Telephone      string  `json:"telephone,omitempty"`
		}

		payload := fakePatchPayload{
			Active:    true,
			FirstName: "Leo",
			LastName:  "Spaceman",
			Telephone: "800-588-2300",
		}

		session := auth.Session{
			AdminUserID: adminUserID,
		}

		appCtx := appcontext.NewAppContext(nil, logger, &session)

		req := &http.Request{
			URL: &url.URL{
				Path: "/admin/v1/admin-users/778acee1-bb04-4ccf-80bf-eae3c66e8c22",
			},
			Method: "PATCH",
		}

		zapFields, _ := Capture(appCtx, &model, &payload, req)

		var fieldsChanged string
		var eventType string
		for _, field := range zapFields {
			if field.Key == "fields_changed" {
				fieldsChanged = field.String
			}

			if field.Key == "event_type" {
				eventType = field.String
			}
		}

		if assert.NotEmpty(t, zapFields) {
			assert.Equal(t, "active,first_name,last_name,telephone", fieldsChanged)
			assert.Equal(t, "audit_patch_admin_users", eventType)
		}
	})

	t.Run("service member session should not include names", func(t *testing.T) {
		uuidString := "88c9922f-58c7-45cd-8c10-48f2a52bbabc"
		serviceMemberID, _ := uuid.FromString(uuidString)

		session := auth.Session{
			ServiceMemberID: serviceMemberID,
		}

		appCtx := appcontext.NewAppContext(nil, logger, &session)

		zapFields, _ := Capture(appCtx, &model, nil, &dummyRequest)

		if assert.NotEmpty(t, zapFields) {
			var keys []string
			for _, field := range zapFields {
				keys = append(keys, field.Key)
			}

			assert.NotContains(t, "responsible_user_name", keys)
		}
	})

	t.Run("success when a non-pointer is passed in", func(t *testing.T) {
		session := auth.Session{}

		appCtx := appcontext.NewAppContext(nil, logger, &session)

		zapFields, err := Capture(appCtx, model, nil, &dummyRequest)

		var eventType string
		for _, field := range zapFields {
			if field.Key == "event_type" {
				eventType = field.String
			}
		}
		assert.Nil(t, err)
		if assert.NotEmpty(t, zapFields) {
			assert.Equal(t, "event_type", zapFields[0].Key)
			assert.Equal(t, "audit__", eventType)
		}
	})

	t.Run("success when a non-struct is passed in", func(t *testing.T) {
		session := auth.Session{}

		appCtx := appcontext.NewAppContext(nil, logger, &session)

		invalidArg := 5
		zapFields, err := Capture(appCtx, &invalidArg, nil, &dummyRequest)

		var eventType string
		for _, field := range zapFields {
			if field.Key == "event_type" {
				eventType = field.String
			}
		}
		assert.Nil(t, err)
		if assert.NotEmpty(t, zapFields) {
			assert.Equal(t, "event_type", zapFields[0].Key)
			assert.Equal(t, "audit__", eventType)
		}
	})
}

func TestCaptureAccountStatus(t *testing.T) {
	uuidStringOffice := "1127bdbd-0610-4e52-9f10-1fa3c063bad3"
	officeUserID, _ := uuid.FromString(uuidStringOffice)
	model := models.OfficeUser{
		ID:        officeUserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	logger := zap.NewNop()

	uuidStringAdmin := "4ad12fe7-1514-4b6b-a35d-ce68e6c5b1fc"
	adminUserID, _ := uuid.FromString(uuidStringAdmin)

	session := auth.Session{
		AdminUserID: adminUserID,
	}

	appCtx := appcontext.NewAppContext(nil, logger, &session)

	req := &http.Request{
		URL: &url.URL{
			Path: "/admin/v1/admin-users",
		},
		Method: "POST",
	}

	t.Run("Sucessfully logs account enabled", func(t *testing.T) {
		zapFields, _ := CaptureAccountStatus(appCtx, &model, true, req)

		if assert.NotEmpty(t, zapFields) {
			fieldsMap := map[string]string{}
			for _, f := range zapFields {
				fieldsMap[f.Key] = f.String
			}

			assert.Equal(t, "audit_post_admin_users_active_status_changed", fieldsMap["event_type"])
			assert.Equal(t, "true", fieldsMap["active_value"])
		}
	})

	t.Run("Sucessfully logs account disabled", func(t *testing.T) {
		zapFields, _ := CaptureAccountStatus(appCtx, &model, false, req)

		if assert.NotEmpty(t, zapFields) {
			fieldsMap := map[string]string{}
			for _, f := range zapFields {
				fieldsMap[f.Key] = f.String
			}

			assert.Equal(t, "audit_post_admin_users_active_status_changed", fieldsMap["event_type"])
			assert.Equal(t, "false", fieldsMap["active_value"])
		}
	})
}

func TestExtractResponsibleUser(t *testing.T) {
	uuidStringAdmin := "4ad12fe7-1514-4b6b-a35d-ce68e6c5b1fc"
	adminUserID, _ := uuid.FromString(uuidStringAdmin)
	uuidStringUser := "4ad12fe7-1514-4b6b-a35d-ce68e6c5b1fb"
	userID, _ := uuid.FromString(uuidStringUser)

	userEmail := "test@fake.com"

	session := auth.Session{
		AdminUserID: adminUserID,
		UserID:      userID,
		Email:       userEmail,
		FirstName:   "John",
		LastName:    "Doe",
	}

	var zapFields []zap.Field

	t.Run("Returns the require fields", func(t *testing.T) {
		zapFields = extractResponsibleUser(&session)

		if assert.NotEmpty(t, zapFields) {
			fieldsMap := map[string]string{}
			for _, f := range zapFields {
				fieldsMap[f.Key] = f.String
			}

			assert.Equal(t, uuidStringUser, fieldsMap["responsible_user_id"])
			assert.Equal(t, userEmail, fieldsMap["responsible_user_email"])
			assert.Equal(t, "John Doe", fieldsMap["responsible_user_name"])
		}
	})
}

func TestExtractRecordInformation(t *testing.T) {
	uuidStringUser := "4ad12fe7-1514-4b6b-a35d-ce68e6c5b1fb"
	userID, _ := uuid.FromString(uuidStringUser)

	model := &models.OfficeUser{
		CreatedAt: time.Now(),
		ID:        userID,
		UpdatedAt: time.Now(),
	}

	item, _ := validateInterface(model)
	var zapFields []zap.Field
	t.Run("Returns the require fields", func(t *testing.T) {
		zapFields = extractRecordInformation(item, model)

		if assert.NotEmpty(t, zapFields) {
			fieldsMap := map[string]string{}
			for _, f := range zapFields {
				fieldsMap[f.Key] = f.String
			}

			assert.Equal(t, uuidStringUser, fieldsMap["record_id"])
			assert.Equal(t, "OfficeUser", fieldsMap["record_type"])
			assert.Equal(t, model.CreatedAt.String(), fieldsMap["record_created_at"])
			assert.Equal(t, model.UpdatedAt.String(), fieldsMap["record_updated_at"])
		}
	})
}
