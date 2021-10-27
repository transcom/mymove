package ordersapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// GetOrdersHandler returns Orders by uuid
type GetOrdersHandler struct {
	handlers.HandlerContext
}

// Handle (GetOrdersHandler) responds to GET /orders/{uuid}
func (h GetOrdersHandler) Handle(params ordersoperations.GetOrdersParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	logger := h.LoggerFromContext(ctx)

	clientCert := authentication.ClientCertFromContext(ctx)
	if clientCert == nil {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to access this API"))
	}

	id, err := uuid.FromString(params.UUID.String())
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	orders, err := models.FetchElectronicOrderByID(h.DB(), id)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	if !verifyOrdersReadAccess(orders.Issuer, clientCert) {
		return handlers.ResponseForError(logger, errors.WithMessage(models.ErrFetchForbidden, "Not permitted to read Orders from this issuer"))
	}

	ordersPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	return ordersoperations.NewGetOrdersOK().WithPayload(ordersPayload)
}
