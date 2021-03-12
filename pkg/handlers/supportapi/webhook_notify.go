package supportapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	webhookoperations "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/webhook"
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
func (h ReceiveWebhookNotificationHandler) Handle(params webhookoperations.ReceiveWebhookNotificationParams) middleware.Responder {
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
	return webhookoperations.NewReceiveWebhookNotificationOK().WithPayload(payload)
}

// CreateWebhookNotificationHandler is the interface to handle the createWebhookNotification
type CreateWebhookNotificationHandler struct {
	handlers.HandlerContext
}

// Handle handles the endpoint request to the createWebhookNotification handler
func (h CreateWebhookNotificationHandler) Handle(params webhookoperations.CreateWebhookNotificationParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	message := "{ \"message\": \"This is a test notification\" }"
	traceID := h.GetTraceID()
	notification := models.WebhookNotification{
		EventKey: string(event.MoveTaskOrderUpdateEventKey),
		TraceID:  &traceID,
		Payload:  &message,
		Status:   models.WebhookNotificationPending,
	}
	verrs, err := h.DB().ValidateAndCreate(&notification)
	if verrs != nil && verrs.HasAny() {
		logger.Error("Could not store ", zap.Error(verrs))
		return webhookoperations.NewCreateWebhookNotificationInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
	}
	if err != nil {
		return webhookoperations.NewCreateWebhookNotificationInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
	}

	payload := supportmessages.WebhookNotification{
		EventKey: notification.EventKey,
		Object:   notification.Payload,
	}
	return webhookoperations.NewCreateWebhookNotificationCreated().WithPayload(&payload)
}
