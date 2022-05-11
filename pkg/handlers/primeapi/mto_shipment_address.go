package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMTOShipmentAddressHandler is the handler to update an address
type UpdateMTOShipmentAddressHandler struct {
	handlers.HandlerContext
	MTOShipmentAddressUpdater services.MTOShipmentAddressUpdater
}

// Handle updates an address on a shipment
func (h UpdateMTOShipmentAddressHandler) Handle(params mtoshipmentops.UpdateMTOShipmentAddressParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

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
				appCtx.Logger().Error("primeapi.UpdateMTOShipmentAddressHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.PreconditionFailedError:
					return mtoshipmentops.NewUpdateMTOShipmentAddressPreconditionFailed().WithPayload(
						payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// Not Found Error -> Not Found Response
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentAddressNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// InvalidInputError -> Unprocessable Entity Response
				case apperror.InvalidInputError:
					return mtoshipmentops.NewUpdateMTOShipmentAddressUnprocessableEntity().WithPayload(
						payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				// ConflictError -> ConflictError Response
				case apperror.ConflictError:
					return mtoshipmentops.NewUpdateMTOShipmentAddressConflict().WithPayload(
						payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// QueryError -> Internal Server Error
				case apperror.QueryError:
					if e.Unwrap() != nil {
						appCtx.Logger().Error("primeapi.UpdateMTOShipmentAddressHandler error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewUpdateMTOShipmentAddressInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// Unknown -> Internal Server Error
				default:
					return mtoshipmentops.NewUpdateMTOShipmentAddressInternalServerError().
						WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}

			}

			// If no error, create a successful payload to return
			mtoShipmentAddressPayload := payloads.Address(updatedAddress)
			return mtoshipmentops.NewUpdateMTOShipmentAddressOK().WithPayload(mtoShipmentAddressPayload), nil
		})
}
