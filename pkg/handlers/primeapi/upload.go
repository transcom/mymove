package primeapi

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gofrs/uuid"

	"go.uber.org/zap"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForPaymentRequestUploadModel(u models.Upload) *primemessages.Upload {
	return &primemessages.Upload{
		Bytes:       &u.Bytes,
		ContentType: &u.ContentType,
		Filename:    &u.Filename,
		CreatedAt:   (strfmt.DateTime)(u.CreatedAt),
		UpdatedAt:   (strfmt.DateTime)(u.UpdatedAt),
	}
}

// CreateUploadHandler is the create upload handler
type CreateUploadHandler struct {
	handlers.HandlerConfig
	services.PaymentRequestUploadCreator
}

// Handle creates uploads
func (h CreateUploadHandler) Handle(params paymentrequestop.CreateUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			var contractorID uuid.UUID
			contractor, err := models.FetchGHCPrimeTestContractor(appCtx.DB())
			if err != nil {
				appCtx.Logger().Error("error getting TEST GHC Prime Contractor", zap.Error(err))
				// Setting a custom message so we don't reveal the SQL error:
				return paymentrequestop.NewCreateUploadBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"Unable to get the TEST GHC Prime Contractor.", h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}
			if contractor != nil {
				contractorID = contractor.ID
			} else {
				err = apperror.NewBadDataError("error with TEST GHC Prime Contractor value is nil")
				appCtx.Logger().Error(err.Error())
				// Same message as before (same base issue):
				return paymentrequestop.NewCreateUploadBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"Unable to get the TEST GHC Prime Contractor.", h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			paymentRequestID, err := uuid.FromString(params.PaymentRequestID)
			if err != nil {
				appCtx.Logger().Error("error creating uuid from string", zap.Error(err))
				return paymentrequestop.NewCreateUploadUnprocessableEntity().WithPayload(payloads.ValidationError(
					"The payment request ID must be a valid UUID.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), err
			}

			file, ok := params.File.(*runtime.File)
			if !ok {
				err = apperror.NewInternalServerError("this should always be a runtime.File, something has changed in go-swagger")
				appCtx.Logger().Error(err.Error())
				return paymentrequestop.NewCreateUploadInternalServerError(), err
			}

			createdUpload, err := h.PaymentRequestUploadCreator.CreateUpload(appCtx, file.Data, paymentRequestID, contractorID, file.Header.Filename)
			if err != nil {
				appCtx.Logger().Error("primeapi.CreateUploadHandler error", zap.Error(err))
				switch e := err.(type) {
				case *apperror.BadDataError:
					return paymentrequestop.NewCreateUploadBadRequest().WithPayload(
						payloads.ClientError(handlers.BadRequestErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.NotFoundError:
					return paymentrequestop.NewCreateUploadNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return paymentrequestop.NewCreateUploadUnprocessableEntity().WithPayload(
						payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				default:
					return paymentrequestop.NewCreateUploadInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			returnPayload := payloadForPaymentRequestUploadModel(*createdUpload)
			return paymentrequestop.NewCreateUploadCreated().WithPayload(returnPayload), nil
		})
}
