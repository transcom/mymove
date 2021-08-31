package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateReweighHandler is the handler to update a reweigh
type UpdateReweighHandler struct {
	handlers.HandlerContext
	ReweighUpdater services.ReweighUpdater
}

// Handle updates on a reweigh
func (h UpdateReweighHandler) Handle(params mtoshipmentops.UpdateReweighParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	// Get the etag and payload
	payload := params.Body
	eTag := params.IfMatch

	// Get the new reweigh model
	newReweigh := payloads.ReweighModelFromUpdate(payload, params.ReweighID, params.MtoShipmentID)

	// Call the service object
	updatedReweigh, err := h.ReweighUpdater.UpdateReweigh(appCtx, newReweigh, eTag)

	// Convert the errors into error responses to return to caller
	if err != nil {
		logger.Error("primeapi.UpdateReweighHandler", zap.Error(err))

		switch e := err.(type) {
		case services.PreconditionFailedError:
			return mtoshipmentops.NewUpdateReweighPreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		// Not Found Error -> Not Found Response
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateReweighNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		// InvalidInputError -> Unprocessable Entity Response
		case services.InvalidInputError:
			return mtoshipmentops.NewUpdateReweighUnprocessableEntity().WithPayload(
				payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))
		// ConflictError -> ConflictError Response
		case services.ConflictError:
			return mtoshipmentops.NewUpdateReweighConflict().WithPayload(
				payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		// QueryError -> Internal Server Error
		case services.QueryError:
			if e.Unwrap() != nil {
				logger.Error("primeapi.UpdateReweighHandler error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewUpdateReweighInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		// Unknown -> Internal Server Error
		default:
			return mtoshipmentops.NewUpdateReweighInternalServerError().
				WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}

	}

	// If no error, create a successful payload to return
	reweighPayload := payloads.Reweigh(updatedReweigh)
	return mtoshipmentops.NewUpdateReweighOK().WithPayload(reweighPayload)

}
