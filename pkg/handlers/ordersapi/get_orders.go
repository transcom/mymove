package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/models"
)

// GetOrdersHandler returns Orders by uuid
type GetOrdersHandler struct {
	handlers.HandlerConfig
}

// Handle (GetOrdersHandler) responds to GET /orders/{uuid}
func (h GetOrdersHandler) Handle(params ordersoperations.GetOrdersParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	clientCert := authentication.ClientCertFromContext(ctx)
	if clientCert == nil {
		return handlers.ResponseForError(appCtx.Logger(), errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(appCtx.Logger(), errors.WithMessage(models.ErrFetchForbidden, "Not permitted to access this API"))
	}

	id, err := uuid.FromString(params.UUID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	orders, err := models.FetchElectronicOrderByID(appCtx.DB(), id)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	if !verifyOrdersReadAccess(orders.Issuer, clientCert) {
		return handlers.ResponseForError(appCtx.Logger(), errors.WithMessage(models.ErrFetchForbidden, "Not permitted to read Orders from this issuer"))
	}

	ordersPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	return ordersoperations.NewGetOrdersOK().WithPayload(ordersPayload)
}
