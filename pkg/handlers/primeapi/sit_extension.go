package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/transcom/mymove/pkg/services"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
)

// CreateSITExtensionHandler is the handler to create a sit extension
type CreateSITExtensionHandler struct {
	handlers.HandlerContext
	SITExtensionCreator services.SITExtensionCreator
}

// Handle created a sit extension for a shipment
func (h CreateSITExtensionHandler) Handle(params mtoshipmentops.CreateSITExtensionParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	// Get the new extension model
	SITExtension := payloads.SITExtensionModel(params.Body, params.MtoShipmentID)

	// Call the service object
	createdextension, err := h.SITExtensionCreator.CreateSITExtension(appCtx, SITExtension)

	// Convert the errors into error responses to return to caller
	if err != nil {
		logger.Error("primeapi.CreateSITExtensionHandler", zap.Error(err))

		switch e := err.(type) {
		// NotFoundError -> Not Found Response
		case services.NotFoundError:
			return mtoshipmentops.NewCreateSITExtensionNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
			// ConflictError -> Conflict Response
		case services.ConflictError:
			return mtoshipmentops.NewCreateSITExtensionConflict().WithPayload(
				payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		// InvalidInputError -> Unprocessable Entity Response
		case services.InvalidInputError:
			return mtoshipmentops.NewCreateSITExtensionUnprocessableEntity().WithPayload(
				payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))
		// QueryError -> Internal Server Error
		case services.QueryError:
			if e.Unwrap() != nil {
				logger.Error("primeapi.CreateSITExtensionHandler error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewCreateSITExtensionInternalServerError().WithPayload(
				payloads.InternalServerError(nil, h.GetTraceID()))
		// Unknown -> Internal Server Error
		default:
			return mtoshipmentops.NewCreateSITExtensionInternalServerError().
				WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}

	}
	// If no error, create a successful payload to return
	payload := payloads.SITExtension(createdextension)
	return mtoshipmentops.NewCreateSITExtensionCreated().WithPayload(payload)
}
