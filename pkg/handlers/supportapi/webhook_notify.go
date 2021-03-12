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

// PostWebhookNotifyHandler passes through a message
type PostWebhookNotifyHandler struct {
	handlers.HandlerContext
}

// Handle posts message
func (h PostWebhookNotifyHandler) Handle(params webhookoperations.PostWebhookNotifyParams) middleware.Responder {

	payload := &supportmessages.WebhookNotification{
		ID:          params.Body.ID,
		EventName:   params.Body.EventName,
		TriggeredAt: params.Body.TriggeredAt,
		ObjectType:  params.Body.ObjectType,
		Object:      params.Body.Object,
	}

	return webhookoperations.NewPostWebhookNotifyOK().WithPayload(payload)
}

type CreateWebhookNotificationHandler struct {
	handlers.HandlerContext
}

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
		EventName: notification.EventKey,
		Object:    *notification.Payload,
	}
	return webhookoperations.NewCreateWebhookNotificationCreated().WithPayload(&payload)
}
