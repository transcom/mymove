package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	reserviceitemsop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/re_service_items"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// GetReServiceItemsHandler returns a list of service items
type GetReServiceItemsHandler struct {
	handlers.HandlerConfig
	services.ServiceItemListFetcher
}

func (h GetReServiceItemsHandler) Handle(params reserviceitemsop.GetAllReServiceItemsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			serviceItemList, err := h.FetchServiceItemList(appCtx)
			if err != nil {
				appCtx.Logger().Error("Error fetching Service Item List", zap.Error(err))
				return reserviceitemsop.NewGetAllReServiceItemsInternalServerError(), err
			}
			returnPayload := payloads.ReServiceItems(*serviceItemList)
			return reserviceitemsop.NewGetAllReServiceItemsOK().WithPayload(returnPayload), nil
		})
}
