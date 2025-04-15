package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/appcontext"
	edierrorsop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/e_d_i_errors"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForEdiErrorModel(e models.EdiError) *adminmessages.EdiError {
	return &adminmessages.EdiError{
		ID:               handlers.FmtUUID(e.ID),
		PaymentRequestID: handlers.FmtUUID(e.PaymentRequestID),
		Code:             e.Code,
		Description:      e.Description,
		EdiType:          (*string)(&e.EDIType),
		CreatedAt:        strfmt.DateTime(e.CreatedAt),
		// UpdatedAt:        handlers.FmtDateTimePtr(&e.UpdatedAt),
	}
}

type FetchEdiErrorsHandler struct {
	handlers.HandlerConfig
	ediErrorFetcher services.EDIErrorFetcher
}

// Handle retrieves a list of edi errors
func (h FetchEdiErrorsHandler) Handle(params edierrorsop.FetchEdiErrorsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			ediErrorPaymentRequests, err := h.ediErrorFetcher.FetchEdiErrors(appCtx)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			payload := make(adminmessages.EdiErrors, len(ediErrorPaymentRequests))
			for i, r := range ediErrorPaymentRequests {
				payload[i] = payloadForEdiErrorModel(r)
			}

			return edierrorsop.NewFetchEdiErrorsOK().
				WithContentRange(fmt.Sprintf("edi_errors %d-%d/%d", 0, len(ediErrorPaymentRequests), len(ediErrorPaymentRequests))).
				WithPayload(payload), nil
		})
}
