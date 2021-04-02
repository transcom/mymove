package supportapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/services/event"

	webhookops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/webhook"
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

	objectID := "<empty>"
	if notif.ObjectID != nil {
		objectID = notif.ObjectID.String()
	}
	mtoID := "<empty>"
	if notif.MoveTaskOrderID != nil {
		mtoID = notif.MoveTaskOrderID.String()
	}

	logger.Info("Received Webhook Notification: ",
		zap.String("id", notif.ID.String()),
		zap.String("eventKey", notif.EventKey),
		zap.String("createdAt", notif.CreatedAt.String()),
		zap.String("traceID", notif.TraceID.String()),
		zap.String("moveID", mtoID),
		zap.String("objectID", objectID))
	return webhookops.NewReceiveWebhookNotificationOK().WithPayload(notif)
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
	if payload == nil {
		// create a default notification payload
		message := "{ \"message\": \"This is a test notification\" }"
		payload = &supportmessages.WebhookNotification{
			EventKey: string(event.TestCreateEventKey),
			TraceID:  *handlers.FmtUUID(h.GetTraceID()),
			Object:   swag.String(message),
			Status:   supportmessages.WebhookNotificationStatusPENDING,
		}
	}
	// Convert to model and create in DB
	notification, verrs := payloads.WebhookNotificationModel(payload, h.GetTraceID())
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
		return webhookops.NewCreateWebhookNotificationInternalServerError().WithPayload(payloads.InternalServerError(swag.String(err.Error()), h.GetTraceID()))
	}

	payload = payloads.WebhookNotification(notification)
	return webhookops.NewCreateWebhookNotificationCreated().WithPayload(payload)
}
