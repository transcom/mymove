package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	edierrorsop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/e_d_i_errors"
	singleedierrorop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/single_e_d_i_error"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForEdiErrorModel(e models.EdiError) *adminmessages.EdiError {
	return &adminmessages.EdiError{
		ID:                   handlers.FmtUUID(e.ID),
		PaymentRequestID:     handlers.FmtUUID(e.PaymentRequestID),
		PaymentRequestNumber: e.PaymentRequest.PaymentRequestNumber,
		Code:                 e.Code,
		Description:          e.Description,
		EdiType:              (*string)(&e.EDIType),
		CreatedAt:            strfmt.DateTime(e.CreatedAt),
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

// GetEdiErrorHandler returns a single EDI error by ID via GET /edi-errors/{id}
type GetEdiErrorHandler struct {
	handlers.HandlerConfig
	ediErrorFetcher services.EDIErrorFetcher
}

// Handle retrieves a specific EDI error
func (h GetEdiErrorHandler) Handle(params singleedierrorop.GetEdiErrorParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest, func(appCtx appcontext.AppContext) (middleware.Responder, error) {
		ediErrorID, err := uuid.FromString(params.EdiErrorID.String())
		if err != nil {
			return handlers.ResponseForError(appCtx.Logger(), err), err
		}

		ediError, err := h.ediErrorFetcher.FetchEdiErrorByID(appCtx, ediErrorID)
		if err != nil {
			return handlers.ResponseForError(appCtx.Logger(), err), err
		}

		return singleedierrorop.NewGetEdiErrorOK().WithPayload(payloadForEdiErrorModel(ediError)), nil
	})
}
