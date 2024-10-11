package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	pp "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/payment_request_syncada_file"
	pre "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/payment_request_syncada_files"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type IndexPaymentRequestSyncadaFilesHandler struct {
	handlers.HandlerConfig
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

func (h IndexPaymentRequestSyncadaFilesHandler) Handle(params pre.IndexPaymentRequestSyncadaFilesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, paymentRequestNumberFilter)
			ordering := query.NewQueryOrder(params.Sort, params.Order)
			pagination := h.NewPagination(params.Page, params.PerPage)
			var paymentRequestEdiFiles models.PaymentRequestEdiFiles
			err := h.ListFetcher.FetchRecordList(appCtx, &paymentRequestEdiFiles, queryFilters, nil, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			totalPaymentRequestSyncadaFilesCount, err := h.ListFetcher.FetchRecordCount(appCtx, &paymentRequestEdiFiles, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			queriedPaymentRequestEdiFilesCount := len(paymentRequestEdiFiles)
			payload := make([]*adminmessages.PaymentRequestSyncadaFile, queriedPaymentRequestEdiFilesCount)
			for i, paymentRequestEdiFile := range paymentRequestEdiFiles {
				payload[i] = payloadForPaymentRequestEdiFile(paymentRequestEdiFile)
			}
			return pre.NewIndexPaymentRequestSyncadaFilesOK().WithContentRange(fmt.Sprintf("payment-request-syncada-files %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedPaymentRequestEdiFilesCount, totalPaymentRequestSyncadaFilesCount)).WithPayload(payload), nil
		})
}

func payloadForPaymentRequestEdiFile(paymentRequestEdiFile models.PaymentRequestEdiFile) *adminmessages.PaymentRequestSyncadaFile {
	paymentRequestSyncadaFilePayload := &adminmessages.PaymentRequestSyncadaFile{
		ID:                   *handlers.FmtUUID(paymentRequestEdiFile.ID),
		PaymentRequestNumber: paymentRequestEdiFile.PaymentRequestNumber,
		FileName:             paymentRequestEdiFile.Filename,
		EdiString:            paymentRequestEdiFile.EdiString,
		CreatedAt:            *handlers.FmtDateTime(paymentRequestEdiFile.CreatedAt),
	}
	return paymentRequestSyncadaFilePayload
}

type GetPaymentRequestSyncadaFileHandler struct {
	handlers.HandlerConfig
	services.PaymentRequestSyncadaFileFetcher
	services.NewQueryFilter
}

// Handle implements payment_request_syncada_file.PaymentRequestSyncadaFileHandler.
func (g GetPaymentRequestSyncadaFileHandler) Handle(params pp.PaymentRequestSyncadaFileParams) middleware.Responder {
	return g.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			paymentRequestSyncadaFileID := uuid.FromStringOrNil(params.PaymentRequestSyncadaFileID.String())
			queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", paymentRequestSyncadaFileID)}

			paymentRequestSyncadaFile, err := g.PaymentRequestSyncadaFileFetcher.FetchPaymentRequestSyncadaFile(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			payload := payloadForPaymentRequestEdiFile(paymentRequestSyncadaFile)
			return pp.NewPaymentRequestSyncadaFileOK().WithPayload(payload), nil
		})
}

var paymentRequestNumberFilter = map[string]func(string) []services.QueryFilter{
	"paymentRequestNumber": func(content string) []services.QueryFilter {
		return []services.QueryFilter{
			query.NewQueryFilter("payment_request_number", "ILIKE", fmt.Sprintf("%%%s%%", content))}
	},
}
