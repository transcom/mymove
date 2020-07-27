package supportapi

import (
	"github.com/go-openapi/runtime/middleware"

	supportoperations "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/supportapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
)

// PostNotificationHandler passes through a message
type PostNotificationHandler struct {
	handlers.HandlerContext
}

// Handle posts message
func (h PostNotificationHandler) Handle(params supportoperations.PostNotificationParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	logger := h.LoggerFromContext(ctx)

	newNotification := &models.ServerNotification{
		Message: params.Body.Message,
	}

	// notificationPayload, err := params.Body.Message
	notificationPayload, err := payloads.Notification(newNotification)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	return supportoperations.NewPostNotificationOK().WithPayload(notificationPayload)

}
