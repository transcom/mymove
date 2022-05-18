package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/services"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
)

// CreateSITExtensionHandler is the handler to create a sit extension
type CreateSITExtensionHandler struct {
	handlers.HandlerConfig
	SITExtensionCreator services.SITExtensionCreator
}

// Handle created a sit extension for a shipment
func (h CreateSITExtensionHandler) Handle(params mtoshipmentops.CreateSITExtensionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// Get the new extension model
			SITExtension := payloads.SITExtensionModel(params.Body, params.MtoShipmentID)

			// Call the service object
			createdExtension, err := h.SITExtensionCreator.CreateSITExtension(appCtx, SITExtension)

			// Convert the errors into error responses to return to caller
			if err != nil {
				appCtx.Logger().Error("primeapi.CreateSITExtensionHandler", zap.Error(err))

				switch e := err.(type) {
				// NotFoundError -> Not Found Response
				case apperror.NotFoundError:
					return mtoshipmentops.NewCreateSITExtensionNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
					// ConflictError -> Conflict Response
				case apperror.ConflictError:
					return mtoshipmentops.NewCreateSITExtensionConflict().WithPayload(
						payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// InvalidInputError -> Unprocessable Entity Response
				case apperror.InvalidInputError:
					return mtoshipmentops.NewCreateSITExtensionUnprocessableEntity().WithPayload(
						payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				// QueryError -> Internal Server Error
				case apperror.QueryError:
					if e.Unwrap() != nil {
						appCtx.Logger().Error("primeapi.CreateSITExtensionHandler error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewCreateSITExtensionInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// Unknown -> Internal Server Error
				default:
					return mtoshipmentops.NewCreateSITExtensionInternalServerError().
						WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}

			}
			// If no error, create a successful payload to return
			payload := payloads.SITExtension(createdExtension)
			return mtoshipmentops.NewCreateSITExtensionCreated().WithPayload(payload), nil
		})
}
