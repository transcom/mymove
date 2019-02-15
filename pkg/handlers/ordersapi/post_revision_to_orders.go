package ordersapi

import (
	"reflect"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	beeline "github.com/honeycombio/beeline-go"
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
	if clientCert == nil || !clientCert.AllowOrdersAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return ordersoperations.NewPostRevisionToOrdersUnauthorized()
	}

	// TODO check enumerated values, or are these already checked for me by swagger?

	id, err := uuid.FromString(params.UUID.String())
	if err != nil {
		h.Logger().Info("Not a valid UUID")
		return ordersoperations.NewPostRevisionToOrdersBadRequest()
	}

	orders, err := models.FetchElectronicOrderByID(h.DB(), id)
	if err == models.ErrFetchNotFound {
		return ordersoperations.NewPostRevisionToOrdersNotFound()
	}

	for _, r := range orders.Revisions {
		// SeqNum collision
		if r.SeqNum == int(*params.Revision.SeqNum) {
			return ordersoperations.NewPostRevisionToOrdersConflict()
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
