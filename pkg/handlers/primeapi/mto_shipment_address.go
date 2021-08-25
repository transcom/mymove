package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMTOShipmentAddressHandler is the handler to update an address
type UpdateMTOShipmentAddressHandler struct {
	handlers.HandlerContext
	MTOShipmentAddressUpdater services.MTOShipmentAddressUpdater
}

// Handle updates an address on a shipment
func (h UpdateMTOShipmentAddressHandler) Handle(params mtoshipmentops.UpdateMTOShipmentAddressParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	// Get the params and payload
	payload := params.Body
	eTag := params.IfMatch
	mtoShipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())
	addressID := uuid.FromStringOrNil(params.AddressID.String())

	// Get the new address model
	newAddress := payloads.AddressModel(payload)
	newAddress.ID = addressID

	// Call the service object
	updatedAddress, err := h.MTOShipmentAddressUpdater.UpdateMTOShipmentAddress(appCtx, newAddress, mtoShipmentID, eTag, true)

	// Convert the errors into error responses to return to caller
	if err != nil {
		logger.Error("primeapi.UpdateMTOShipmentAddressHandler", zap.Error(err))

		switch e := err.(type) {
		case services.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOShipmentAddressPreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		// Not Found Error -> Not Found Response
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentAddressNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		// InvalidInputError -> Unprocessable Entity Response
		case services.InvalidInputError:
			return mtoshipmentops.NewUpdateMTOShipmentAddressUnprocessableEntity().WithPayload(
				payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceID(), e.ValidationErrors))
		// ConflictError -> ConflictError Response
		case services.ConflictError:
			return mtoshipmentops.NewUpdateMTOShipmentAddressConflict().WithPayload(
				payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID()))
		// QueryError -> Internal Server Error
		case services.QueryError:
			if e.Unwrap() != nil {
				logger.Error("primeapi.UpdateMTOShipmentAddressHandler error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewUpdateMTOShipmentAddressInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		// Unknown -> Internal Server Error
		default:
			return mtoshipmentops.NewUpdateMTOShipmentAddressInternalServerError().
				WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}

	}

	// If no error, create a successful payload to return
	mtoShipmentAddressPayload := payloads.Address(updatedAddress)
	return mtoshipmentops.NewUpdateMTOShipmentAddressOK().WithPayload(mtoShipmentAddressPayload)

}
