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

	payload := &webhookoperations.PostWebhookNotifyOKBody{
		ID:          params.Body.ID,
		EventName:   params.Body.EventName,
		TriggeredAt: params.Body.TriggeredAt,
		ObjectType:  params.Body.ObjectType,
		Object:      params.Body.Object,
	}

	return webhookoperations.NewPostWebhookNotifyOK().WithPayload(payload)
}
