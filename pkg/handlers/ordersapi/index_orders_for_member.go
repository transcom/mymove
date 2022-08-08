package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/models"
)

// IndexOrdersForMemberHandler returns a list of Orders matching the provided search parameters
type IndexOrdersForMemberHandler struct {
	handlers.HandlerConfig
}

// Handle (IndexOrdersForMemberHandler) responds to GET /edipis/{edipi}/orders
func (h IndexOrdersForMemberHandler) Handle(params ordersoperations.IndexOrdersForMemberParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	clientCert := authentication.ClientCertFromContext(ctx)
	if clientCert == nil {
		return handlers.ResponseForError(appCtx.Logger(), errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(appCtx.Logger(), errors.WithMessage(models.ErrFetchForbidden, "Not permitted to access this API"))
	}
	allowedIssuers := clientCert.GetAllowedOrdersIssuersRead()
	if len(allowedIssuers) == 0 {
		return handlers.ResponseForError(appCtx.Logger(), errors.WithMessage(models.ErrFetchForbidden, "Not permitted to read any Orders"))
	}

	orders, err := models.FetchElectronicOrdersByEdipiAndIssuers(appCtx.DB(), params.Edipi, allowedIssuers)
	if err == models.ErrFetchNotFound {
		return ordersoperations.NewIndexOrdersForMemberOK().WithPayload([]*ordersmessages.Orders{})
	} else if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	ordersPayloads := make([]*ordersmessages.Orders, len(orders))
	for i, o := range orders {
		payload, err := payloadForElectronicOrderModel(o)
		if err != nil {
			return handlers.ResponseForError(appCtx.Logger(), err)
		}
		ordersPayloads[i] = payload
	}

	return ordersoperations.NewIndexOrdersForMemberOK().WithPayload(ordersPayloads)
}
