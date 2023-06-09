package primeapi

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type CreateServiceRequestDocumentUploadHandler struct {
	handlers.HandlerConfig
	services.ServiceRequestDocumentUploadCreator
}

func (h CreateServiceRequestDocumentUploadHandler) Handle(params mtoserviceitemops.CreateServiceRequestDocumentUploadParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest, func(appCtx appcontext.AppContext) (middleware.Responder, error) {
		var contractorID uuid.UUID
		contractor, err := models.FetchGHCPrimeContractor(appCtx.DB())
		if err != nil {
			appCtx.Logger().Error("error getting GHC Prime Contractor", zap.Error(err))
			return mtoserviceitemops.NewCreateServiceRequestDocumentUploadBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, "Unable to get the GHC Prime Contractor.", h.GetTraceIDFromRequest(params.HTTPRequest))), err
		}

		if contractor != nil {
			contractorID = contractor.ID
		} else {
			err = apperror.NewBadDataError("error with GHC Prime Contractor value is nil")
			appCtx.Logger().Error(err.Error())
			// Same message as before (same base issue):
			return mtoserviceitemops.NewCreateServiceRequestDocumentUploadBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
				"Unable to get the GHC Prime Contractor.", h.GetTraceIDFromRequest(params.HTTPRequest))), err
		}

		serviceItemID, err := uuid.FromString(params.MtoServiceItemID)
		if err != nil {
			appCtx.Logger().Error("error creating uuid from string", zap.Error(err))
			return mtoserviceitemops.NewCreateServiceRequestDocumentUploadUnprocessableEntity().WithPayload(payloads.ValidationError("The service item ID must be a valid UUID.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), err
		}

		file, ok := params.File.(*runtime.File)
		if !ok {
			err = apperror.NewInternalServerError("this should always be a runtime.File, something has changed in go-swagger")
			appCtx.Logger().Error(err.Error())
			return mtoserviceitemops.NewCreateServiceRequestDocumentUploadInternalServerError(), err
		}

		createServiceRequestDocumentUpload, err := h.ServiceRequestDocumentUploadCreator.CreateUpload(appCtx, file.Data, serviceItemID, contractorID, file.Header.Filename)
		if err != nil {
			appCtx.Logger().Error("primeapi.CreateUploadHandler error", zap.Error(err))
			switch e := err.(type) {
			case *apperror.BadDataError:
				return mtoserviceitemops.NewCreateServiceRequestDocumentUploadBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			case apperror.NotFoundError:
				return mtoserviceitemops.NewCreateServiceRequestDocumentUploadNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			case apperror.InvalidInputError:
				return mtoserviceitemops.NewCreateServiceRequestDocumentUploadUnprocessableEntity().WithPayload(payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
			default:
				return mtoserviceitemops.NewCreateServiceRequestDocumentUploadInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}
		}

		returnPayload := payloads.ServiceRequestDocumentUploadModel(*createServiceRequestDocumentUpload)
		return mtoserviceitemops.NewCreateServiceRequestDocumentUploadCreated().WithPayload(returnPayload), nil
	})
}
