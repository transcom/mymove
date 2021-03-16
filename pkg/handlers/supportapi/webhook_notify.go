package supportapi

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	webhookops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/webhook"
	supportmessages "github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/event"

	"github.com/transcom/mymove/pkg/handlers"
)

// ReceiveWebhookNotificationHandler passes through a message
type ReceiveWebhookNotificationHandler struct {
	handlers.HandlerContext
}

// Handle receipt of message
func (h ReceiveWebhookNotificationHandler) Handle(params webhookops.ReceiveWebhookNotificationParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	notif := params.Body

	// This is a test endpoint, it receives a notification, logs it and simply responds with a 200
	payload := &supportmessages.WebhookNotification{
		ID:        params.Body.ID,
		EventKey:  params.Body.EventKey,
		CreatedAt: params.Body.CreatedAt,
		Object:    params.Body.Object,
	}
	objectString := "<empty>"
	if notif.Object != nil {
		objectString = *notif.Object
	}

	logger.Info("Received Webhook Notification: ",
		zap.String("ID", notif.ID.String()),
		zap.String("EventKey", notif.EventKey),
		zap.String("createdAt", notif.CreatedAt.String()),
		zap.String("object", objectString))
	return webhookops.NewReceiveWebhookNotificationOK().WithPayload(payload)
}

// CreateWebhookNotificationHandler is the interface to handle the createWebhookNotification
type CreateWebhookNotificationHandler struct {
	handlers.HandlerContext
}

// WebhookNotificationM2P converts model to payload
func WebhookNotificationM2P(model *models.WebhookNotification) *supportmessages.WebhookNotification {
	payload := supportmessages.WebhookNotification{
		ID:               *handlers.FmtUUID(model.ID),
		EventKey:         model.EventKey,
		Object:           swag.String(model.Payload),
		CreatedAt:        *handlers.FmtDateTime(model.CreatedAt),
		UpdatedAt:        *handlers.FmtDateTime(model.UpdatedAt),
		FirstAttemptedAt: handlers.FmtDateTimePtr(model.FirstAttemptedAt),
		ObjectID:         handlers.FmtUUIDPtr(model.ObjectID),
		TraceID:          *handlers.FmtUUIDPtr(model.TraceID),
		MoveTaskOrderID:  handlers.FmtUUIDPtr(model.MoveTaskOrderID),
	}
	return &payload
}

// WebhookNotificatonP2M converts payload to model
func WebhookNotificatonP2M(payload *supportmessages.WebhookNotification, traceID uuid.UUID) (*models.WebhookNotification, *validate.Errors) {
	verrs := validate.NewErrors()
	var notification *models.WebhookNotification
	if payload == nil {
		// create default notification
		message := "{ \"message\": \"This is a test notification\" }"
		notification = &models.WebhookNotification{
			EventKey: string(event.TestCreateEventKey),
			TraceID:  &traceID,
			Payload:  message,
			Status:   models.WebhookNotificationPending,
		}
	} else {
		if !event.ExistsEventKey(payload.EventKey) {
			verrs.Add("eventKey", "must be a registered event key")
			return nil, verrs
		}
		notification = &models.WebhookNotification{
			// ID is managed by pop
			EventKey:        payload.EventKey,
			TraceID:         &traceID,
			MoveTaskOrderID: handlers.FmtUUIDPtrToPopPtr(payload.MoveTaskOrderID),
			ObjectID:        handlers.FmtUUIDPtrToPopPtr(payload.ObjectID),
			// Payload updated below
			Status: models.WebhookNotificationPending,
			// CreatedAt is managed by pop
			// UpdatedAt is managed by pop
			// FirstAttemptedAt is never provided by user
		}
		if payload.Object != nil {
			notification.Payload = *payload.Object
		}
	}
	return notification, nil
}

// Handle handles the endpoint request to the createWebhookNotification handler
func (h CreateWebhookNotificationHandler) Handle(params webhookops.CreateWebhookNotificationParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	payload := params.Body

	var err error
	notification, verrs := WebhookNotificatonP2M(payload, h.GetTraceID())
	if verrs == nil {
		verrs, err = h.DB().ValidateAndCreate(notification)
	}
	if verrs != nil && verrs.HasAny() {
		logger.Error("Error validating WebhookNotification: ", zap.Error(verrs))

		return webhookops.NewCreateWebhookNotificationUnprocessableEntity().WithPayload(payloads.ValidationError(
			"The notification definition is invalid.", h.GetTraceID(), verrs))
	}
	if err != nil {
		logger.Error("Error creating WebhookNotification: ", zap.Error(err))
		return webhookops.NewCreateWebhookNotificationInternalServerError().WithPayload(payloads.InternalServerError(swag.String("that thing happened"), h.GetTraceID()))
	}
	fmt.Println("got here 3")
	payload = WebhookNotificationM2P(notification)
	fmt.Println("got here 4")
	return webhookops.NewCreateWebhookNotificationCreated().WithPayload(payload)
}
