package ordersapi

import (
	"fmt"
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// PostRevisionToOrdersHandler adds a Revision to Orders by uuid
type PostRevisionToOrdersHandler struct {
	handlers.HandlerContext
}

// Handle (params ordersoperations.PostRevisionToOrdersParams) responds to POST /orders/{uuid}
func (h PostRevisionToOrdersHandler) Handle(params ordersoperations.PostRevisionToOrdersParams) middleware.Responder {
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil {
		return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrUserUnauthorized, "No client certificate provided"))
	}
	if !clientCert.AllowOrdersAPI {
		return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrWriteForbidden, "Not permitted to access this API"))
	}

	id, err := uuid.FromString(params.UUID.String())
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	orders, err := models.FetchElectronicOrderByID(h.DB(), id)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	if !verifyOrdersWriteAccess(orders.Issuer, clientCert) {
		return handlers.ResponseForError(h.Logger(), errors.WithMessage(models.ErrWriteForbidden, "Not permitted to write Orders from this issuer"))
	}

	for _, r := range orders.Revisions {
		// SeqNum collision
		if r.SeqNum == int(*params.Revision.SeqNum) {
			return handlers.ResponseForError(
				h.Logger(),
				errors.WithMessage(
					models.ErrWriteConflict,
					fmt.Sprintf("Cannot POST Revision with SeqNum %d to Orders %s: a Revision with that SeqNum already exists in those Orders", r.SeqNum, params.UUID)))
		}
	}

	newRevision := toElectronicOrdersRevision(orders, params.Revision)
	verrs, err := models.CreateElectronicOrdersRevision(ctx, h.DB(), newRevision)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	orders.Revisions = append(orders.Revisions, *newRevision)

	orderPayload, err := payloadForElectronicOrderModel(orders)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	return ordersoperations.NewPostRevisionToOrdersCreated().WithPayload(orderPayload)
}
