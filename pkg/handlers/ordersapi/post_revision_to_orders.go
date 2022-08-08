package ordersapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/models"
)

// PostRevisionToOrdersHandler adds a Revision to Orders by uuid
type PostRevisionToOrdersHandler struct {
	handlers.HandlerConfig
}

// Handle (params ordersoperations.PostRevisionToOrdersParams) responds to POST /orders/{uuid}
func (h PostRevisionToOrdersHandler) Handle(params ordersoperations.PostRevisionToOrdersParams) middleware.Responder {

	ctx := params.HTTPRequest.Context()

	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	clientCert := authentication.ClientCertFromContext(ctx)
	if clientCert == nil {
		return handlers.ResponseForError(appCtx.Logger(), errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(appCtx.Logger(), errors.WithMessage(models.ErrWriteForbidden, "Not permitted to access this API"))
	}

	id, err := uuid.FromString(params.UUID.String())
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	orders, err := models.FetchElectronicOrderByID(appCtx.DB(), id)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	if !verifyOrdersWriteAccess(orders.Issuer, clientCert) {
		return handlers.ResponseForError(appCtx.Logger(), errors.WithMessage(models.ErrWriteForbidden, "Not permitted to write Orders from this issuer"))
	}

	for _, r := range orders.Revisions {
		// SeqNum collision
		if r.SeqNum == int(*params.Revision.SeqNum) {
			return handlers.ResponseForError(
				appCtx.Logger(),
				errors.WithMessage(
					models.ErrWriteConflict,
					fmt.Sprintf("Cannot POST Revision with SeqNum %d to Orders %s: a Revision with that SeqNum already exists in those Orders", r.SeqNum, params.UUID)))
		}
	}

	newRevision := toElectronicOrdersRevision(orders, params.Revision)
	verrs, err := appCtx.DB().ValidateAndCreate(newRevision)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err)
	}

	orders.Revisions = append(orders.Revisions, *newRevision)

	orderPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}
	return ordersoperations.NewPostRevisionToOrdersCreated().WithPayload(orderPayload)
}
