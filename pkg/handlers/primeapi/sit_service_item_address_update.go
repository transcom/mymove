package primeapi

import (
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	sitaddressupdateops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/sit_address_update"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// CreateSITAddressUpdateRequestHandler is the handler to create a address update request
type CreateSITAddressUpdateRequestHandler struct {
	handlers.HandlerConfig
	SITAddressUpdateRequestCreator services.SITAddressUpdateRequestCreator
}

// Handle creates the address update request
func (h CreateSITAddressUpdateRequestHandler) Handle(params sitaddressupdateops.CreateSITAddressUpdateRequestParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// Get the new address update model
			SITAddressUpdateRequest := payloads.SITAddressUpdateModel(params.Body)

			// Call the service object
			createdSITAddressUpdateRequest, err := h.SITAddressUpdateRequestCreator.CreateSITAddressUpdateRequest(appCtx, SITAddressUpdateRequest)

			//Convert the errors into error responses to return to caller
			if err != nil {
				appCtx.Logger().Error("primeapi.CreateSITAddressUpdateRequestHandler", zap.Error(err))

				switch e := err.(type) {
				// NotFoundError -> Not Found Response
				case apperror.NotFoundError:
					return sitaddressupdateops.NewCreateSITAddressUpdateRequestNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// ConflictError -> Conflict Response
				case apperror.ConflictError:
					return sitaddressupdateops.NewCreateSITAddressUpdateRequestConflict().WithPayload(
						payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// InvalidInputError -> Unprocessable Entity Response
				case apperror.InvalidInputError:
					return sitaddressupdateops.NewCreateSITAddressUpdateRequestUnprocessableEntity().WithPayload(
						payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				// QueryError -> Internal Server Error
				case apperror.QueryError:
					if e.Unwrap() != nil {
						appCtx.Logger().Error("primeapi.CreateSITAddressUpdateRequestHandler error", zap.Error(e.Unwrap()))
					}
					return sitaddressupdateops.NewCreateSITAddressUpdateRequestInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				// Unknown -> Internal Server Error
				default:
					return sitaddressupdateops.NewCreateSITAddressUpdateRequestInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			// If no error, create a succesful payload to return
			payload := payloads.SITAddressUpdate(*createdSITAddressUpdateRequest)
			return sitaddressupdateops.NewCreateSITAddressUpdateRequestCreated().WithPayload(payload), nil
		})
}
