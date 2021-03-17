package supportapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	webhookops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/webhook"
	supportmessages "github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
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

// Handle handles the endpoint request to the createWebhookNotification handler
func (h CreateWebhookNotificationHandler) Handle(params webhookops.CreateWebhookNotificationParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	payload := params.Body

	var err error
	notification, verrs := payloads.WebhookNotificatonModel(payload, h.GetTraceID())
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

	payload = payloads.WebhookNotification(notification)
	return webhookops.NewCreateWebhookNotificationCreated().WithPayload(payload)
}
