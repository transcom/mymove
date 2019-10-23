package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	serviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

type ListServiceItemsHandler struct {
	handlers.HandlerContext
	services.ServiceItemListFetcher
	services.NewQueryFilter
}

func (h ListServiceItemsHandler) Handle(params serviceitemop.ListServiceItemsParams) middleware.Responder {
	id, _ := uuid.NewV4()
	serviceItem := &ghcmessages.ServiceItem{
		ID: handlers.FmtUUID(id),
	}

	var payload ghcmessages.ServiceItems
	payload = append(payload, serviceItem)
	return serviceitemop.NewListServiceItemsOK().WithPayload(payload)
}
