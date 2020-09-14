package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gofrs/uuid"

	"go.uber.org/zap"

	uploadop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/uploads"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
)

func payloadForPaymentRequestUploadModel(u models.Upload) *primemessages.Upload {
	return &primemessages.Upload{
		Bytes:       &u.Bytes,
		ContentType: &u.ContentType,
		Filename:    &u.Filename,
		CreatedAt:   (*strfmt.DateTime)(&u.CreatedAt),
		UpdatedAt:   (*strfmt.DateTime)(&u.UpdatedAt),
	}
}

// CreateUploadHandler is the create upload handler
type CreateUploadHandler struct {
	handlers.HandlerContext
	// To be fixed under this story: https://github.com/transcom/mymove/pull/3775/files#r397219200
	// unable to get logger to pass in for instantiation
	//services.PaymentRequestUploadCreator
}

// Handle creates uploads
func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {
	_, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	var contractorID uuid.UUID
	contractor, err := models.FetchGHCPrimeTestContractor(h.DB())
	if err != nil {
		logger.Error("error getting TEST GHC Prime Contractor", zap.Error(err))
		// Setting a custom message so we don't reveal the SQL error:
		return uploadop.NewCreateUploadBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
			"Unable to get the TEST GHC Prime Contractor.", h.GetTraceID()))
	}
	if contractor != nil {
		contractorID = contractor.ID
	} else {
		logger.Error("error with TEST GHC Prime Contractor value is nil")
		// Same message as before (same base issue):
		return uploadop.NewCreateUploadBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
			"Unable to get the TEST GHC Prime Contractor.", h.GetTraceID()))
	}

	paymentRequestID, err := uuid.FromString(params.PaymentRequestID)
	if err != nil {
		logger.Error("error creating uuid from string", zap.Error(err))
		return uploadop.NewCreateUploadUnprocessableEntity().WithPayload(payloads.ValidationError(
			"The payment request ID must be a valid UUID.", h.GetTraceID(), nil))
	}

	uploadCreator := paymentrequest.NewPaymentRequestUploadCreator(h.DB(), logger, h.FileStorer())
	createdUpload, err := uploadCreator.CreateUpload(params.File, paymentRequestID, contractorID)
	if err != nil {
		logger.Error("primeapi.CreateUploadHandler error", zap.Error(err))
		switch e := err.(type) {
		case *services.BadDataError:
			return uploadop.NewCreateUploadBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, err.Error(), h.GetTraceID()))
		case services.NotFoundError:
			return uploadop.NewCreateUploadNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return uploadop.NewCreateUploadUnprocessableEntity().WithPayload(payloads.ValidationError(err.Error(), h.GetTraceID(), e.ValidationErrors))
		default:
			return uploadop.NewCreateUploadInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}

	returnPayload := payloadForPaymentRequestUploadModel(*createdUpload)
	return uploadop.NewCreateUploadCreated().WithPayload(returnPayload)
}
