package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForInvoiceModel(a *models.Invoice) *apimessages.Invoice {
	if a == nil {
		return nil
	}

	return &apimessages.Invoice{
		ID: *handlers.FmtUUID(a.ID),

		Status:    apimessages.InvoiceStatus(a.Status),
		CreatedAt: *handlers.FmtDateTime(a.CreatedAt),
		UpdatedAt: *handlers.FmtDateTime(a.UpdatedAt),
	}
}

// GetInvoiceHandler returns an invoice
type GetInvoiceHandler struct {
	handlers.HandlerContext
}

// Handle returns a specified invoice
func (h GetInvoiceHandler) Handle(params accessorialop.GetInvoiceParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if session == nil {
		return accessorialop.NewGetInvoiceUnauthorized()
	}
	invoiceID, _ := uuid.FromString(params.InvoiceID.String())
	// TODO: check that if TSP app, user is authorized to view shipment

	// Fetch invoice
	invoice, err := models.FetchInvoice(h.DB(), session, invoiceID)
	if err != nil {
		h.Logger().Error("Error fetching invoice", zap.Error(err))
		return accessorialop.NewGetInvoiceInternalServerError()
	}
	payload := payloadForInvoiceModel(invoice)
	return accessorialop.NewGetInvoiceOK().WithPayload(payload)
}
