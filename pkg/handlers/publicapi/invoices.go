package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/server"
	"go.uber.org/zap"

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
	session := server.SessionFromRequestContext(params.HTTPRequest)

	if session == nil {
		return accessorialop.NewGetInvoiceUnauthorized()
	}

	// Fetch invoice
	invoiceID, _ := uuid.FromString(params.InvoiceID.String())
	invoice, err := models.FetchInvoice(h.DB(), session, invoiceID)
	if err != nil {
		if err == models.ErrFetchNotFound {
			h.Logger().Warn("Invoice not found", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		} else if err == models.ErrFetchForbidden {
			h.Logger().Error("User not permitted to access invoice", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		} else if err == models.ErrUserUnauthorized {
			h.Logger().Error("User not authorized to access invoice", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
		h.Logger().Error("Error fetching invoice", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	payload := payloadForInvoiceModel(invoice)
	return accessorialop.NewGetInvoiceOK().WithPayload(payload)
}
