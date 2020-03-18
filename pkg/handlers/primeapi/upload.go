package primeapi

import (
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	uploadop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
)

func payloadForPaymentRequestUploadModel(u models.Upload) *primemessages.Upload {
	return &primemessages.Upload{
		Bytes:       &u.Bytes,
		ContentType: &u.ContentType,
		Filename:    &u.Filename,
	}
}

// CreateUploadHandler is the create upload handler
type CreateUploadHandler struct {
	handlers.HandlerContext
	services.PaymentRequestUploadCreator
}

// Handle creates uploads
func (h *CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	// TODO https://dp3.atlassian.net/browse/MB-1969
	var contractorID uuid.UUID // TODO not populated. Do not know how get from MTO to Contractor ID
	contractor, err := models.FetchGHCPrimeTestContractor(h.DB())
	if err != nil {
		logger.Error("error getting TEST GHC Prime Contractor", zap.Error(err))
	}
	if contractor != nil {
		contractorID = contractor.ID
	} else {
		logger.Error("error with TEST GHC Prime Contractor value is nil")
	}

	paymentRequestID, err := uuid.FromString(params.PaymentRequestID)
	if err != nil {
		logger.Error("error creating uuid from string", zap.Error(err))
	}

	uploadCreator := paymentrequest.NewPaymentRequestUploadCreator(h.DB(), logger, h.FileStorer())
	createdUpload, err := uploadCreator.CreateUpload(params.File, paymentRequestID, contractorID)
	if err != nil {
		logger.Error("cannot create payment request upload", zap.Error(err))
		return uploadop.NewCreateUploadBadRequest()
	}

	returnPayload := payloadForPaymentRequestUploadModel(*createdUpload)
	return uploadop.NewCreateUploadCreated().WithPayload(returnPayload)
}
