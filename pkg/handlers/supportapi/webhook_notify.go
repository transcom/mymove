package supportapi

import (
	"github.com/go-openapi/runtime/middleware"

	webhookoperations "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/webhook"
	"github.com/transcom/mymove/pkg/handlers"
)

// PostWebhookNotifyHandler passes through a message
type PostWebhookNotifyHandler struct {
	handlers.HandlerContext
}

// Handle posts message
func (h PostWebhookNotifyHandler) Handle(params webhookoperations.PostWebhookNotifyParams) middleware.Responder {

	// Leaving logger commented out here because we'll probably want to use it
	// as we build this out

	//logger := h.LoggerFromContext(ctx)

	payload := &webhookoperations.PostWebhookNotifyOKBody{
		Message: params.Message.Message,
	}

	return webhookoperations.NewPostWebhookNotifyOK().WithPayload(payload)
}
