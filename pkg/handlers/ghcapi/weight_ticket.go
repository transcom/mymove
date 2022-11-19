package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	weightticketops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateWeightTicketHandler
type UpdateWeightTicketHandler struct {
	handlers.HandlerConfig
	weighTicketUpdater services.WeightTicketUpdater
}

func (h UpdateWeightTicketHandler) Handle(params weightticketops.UpdateWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return weightticketops.NewUpdateWeightTicketUnauthorized(), noSessionErr
			}

			payload := params.UpdateWeightTicketPayload
			if payload == nil {
				appCtx.Logger().Error("Invalid Weight Ticket: params Body is nil")
				emptyBodyError := apperror.NewBadDataError("The request body cannot be empty.")
				payload := payloadForValidationError(
					"Empty body error",
					emptyBodyError.Error(),
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors(),
				)

				return weightticketops.NewUpdateWeightTicketUnprocessableEntity().WithPayload(payload), emptyBodyError
			}

			ppmshipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())
			//oldPPMShipment, err := mtoshipment.FindShipment(appCtx, ppmshipmentID)
			// Can't find original weight ticket
			//if err != nil {
			//	appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
			//	switch err.(type) {
			//	case apperror.NotFoundError:
			//		return weightticketops.NewUpdateWeightTicketNotFound(), err
			//	default:
			//		msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
			//
			//		return weightticketops.NewUpdateWeightTicketInternalServerError().WithPayload(
			//			&ghcmessages.Error{Message: &msg},
			//		), err
			//	}
			//}

			weightTicket := payloads.WeightTicketModelFromUpdate(payload)
			weightTicket.ID = ppmshipmentID

			//handleError := func(err error) (middleware.Responder, error) {
			//	appCtx.Logger().Error("ghcapi.UpdateWeightTicketHandler", zap.Error(err))
			//
			//	switch e := err.(type) {
			//	case apperror.NotFoundError:
			//		return weightticketops.NewUpdateWeightTicketNotFound(), err
			//	case apperror.ForbiddenError:
			//		msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
			//		return weightticketops.NewUpdateWeightTicketForbidden().WithPayload(
			//			&ghcmessages.Error{Message: &msg},
			//		), err
			//	case apperror.InvalidInputError:
			//		return weightticketops.NewUpdateWeightTicketUnprocessableEntity().WithPayload(
			//			payloadForValidationError(
			//				handlers.ValidationErrMessage,
			//				err.Error(),
			//				h.GetTraceIDFromRequest(params.HTTPRequest),
			//				e.ValidationErrors,
			//			),
			//		), err
			//	case apperror.PreconditionFailedError:
			//		msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
			//		return weightticketops.NewUpdateWeightTicketPreconditionFailed().WithPayload(
			//			&ghcmessages.Error{Message: &msg},
			//		), err
			//	default:
			//		msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
			//
			//		return weightticketops.NewUpdateWeightTicketInternalServerError().WithPayload(
			//			&ghcmessages.Error{Message: &msg},
			//		), err
			//	}
			//}

			updatedWeightTicket, _ := h.weighTicketUpdater.UpdateWeightTicket(appCtx, *weightTicket, params.IfMatch)
			//if err != nil {
			//	return handleError(err)
			//}
			returnPayload := payloads.WeightTicket(h.FileStorer(), updatedWeightTicket)
			return weightticketops.NewUpdateWeightTicketOK().WithPayload(returnPayload), nil
		})
}
