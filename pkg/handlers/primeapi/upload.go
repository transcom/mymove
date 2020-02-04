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

type CreateUploadHandler struct {
	handlers.HandlerContext
	services.PaymentRequestUploadCreator
}

func (h *CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	userID := session.UserID // TODO: restrict to prime user when prime auth is implemented
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID)
	if err != nil {
		logger.Error("error creating uuid from string", zap.Error(err))
	}

	uploadCreator := paymentrequest.NewPaymentRequestUploadCreator(h.DB(), logger, h.FileStorer())
	createdUpload, err := uploadCreator.CreateUpload(params.File, paymentRequestID, userID)
	if err != nil {
		logger.Error("cannot create payment request upload", zap.Error(err))
		return uploadop.NewCreateUploadBadRequest()
	}

	returnPayload := payloadForPaymentRequestUploadModel(*createdUpload)
	return uploadop.NewCreateUploadCreated().WithPayload(returnPayload)
}
