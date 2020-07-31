package supportapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	supportoperations "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/webhook"
	"github.com/transcom/mymove/pkg/handlers"
)

// PostNotificationHandler passes through a message
type PostNotificationHandler struct {
	handlers.HandlerContext
}

// Handle posts message
func (h PostNotificationHandler) Handle(params supportoperations.PostNotificationParams) middleware.Responder {

	//ctx := params.HTTPRequest.Context()
	//logger := h.LoggerFromContext(ctx)

	newNotification := &supportoperations.PostNotificationOKBody{ // replace
		Message: params.Message.Message,
	}

	notificationPayload := supportoperations.PostNotificationOKBody(*newNotification)
	fmt.Printf("Do we get a nil here??? 🐙")
	fmt.Print(notificationPayload)
	// if err != nil {
	// 	return handlers.ResponseForError(logger, err)
	// }
	return supportoperations.NewPostNotificationOK().WithPayload(&notificationPayload)

}
