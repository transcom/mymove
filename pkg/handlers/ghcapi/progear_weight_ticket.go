package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	progearops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateProgearWeightTicketHandler
type UpdateProgearWeightTicketHandler struct {
	handlers.HandlerConfig
	progearUpdater services.ProgearWeightTicketUpdater
}

func (h UpdateProgearWeightTicketHandler) Handle(params progearops.UpdateProGearWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return progearops.NewUpdateProGearWeightTicketUnauthorized(), noSessionErr
			}

			payload := params.UpdateProGearWeightTicket
			if payload == nil {
				appCtx.Logger().Error("Invalid Progear Weight Ticket: params Body is nil")
				emptyBodyError := apperror.NewBadDataError("The request body cannot be empty.")
				payload := payloadForValidationError(
					"Empty body error",
					emptyBodyError.Error(),
					h.GetTraceIDFromRequest(params.HTTPRequest),
					validate.NewErrors(),
				)

				return progearops.NewUpdateProGearWeightTicketUnprocessableEntity().WithPayload(payload), emptyBodyError
			}

			//ppmshipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())
			//oldPPMShipment, err := mtoshipment.FindShipment(appCtx, ppmshipmentID)
			// Can't find original progear weight ticket
			//if err != nil {
			//	appCtx.Logger().Error("ghcapi.UpdateShipmentHandler", zap.Error(err))
			//	switch err.(type) {
			//	case apperror.NotFoundError:
			//		return progearops.NewUpdateProGearWeightTicketNotFound(), err
			//	default:
			//		msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
			//
			//		return progearops.NewUpdateProGearWeightTicketInternalServerError().WithPayload(
			//			&ghcmessages.Error{Message: &msg},
			//		), err
			//	}
			//}

			progearWeightTicket := payloads.ProgearWeightTicketModelFromUpdate(payload)
			progearWeightTicket.ID = uuid.FromStringOrNil(params.ProGearWeightTicketID.String())

			//handleError := func(err error) (middleware.Responder, error) {
			//	appCtx.Logger().Error("ghcapi.UpdateProgearWeightTicketHandler", zap.Error(err))
			//
			//	switch e := err.(type) {
			//	case apperror.NotFoundError:
			//		return progearops.NewUpdateProGearWeightTicketNotFound(), err
			//	case apperror.ForbiddenError:
			//		msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
			//		return progearops.NewUpdateProGearWeightTicketForbidden().WithPayload(
			//			&ghcmessages.Error{Message: &msg},
			//		), err
			//	case apperror.InvalidInputError:
			//		return progearops.NewUpdateProGearWeightTicketUnprocessableEntity().WithPayload(
			//			payloadForValidationError(
			//				handlers.ValidationErrMessage,
			//				err.Error(),
			//				h.GetTraceIDFromRequest(params.HTTPRequest),
			//				e.ValidationErrors,
			//			),
			//		), err
			//	case apperror.PreconditionFailedError:
			//		msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
			//		return progearops.NewUpdateProGearWeightTicketPreconditionFailed().WithPayload(
			//			&ghcmessages.Error{Message: &msg},
			//		), err
			//	default:
			//		msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
			//
			//		return progearops.NewUpdateProGearWeightTicketInternalServerError().WithPayload(
			//			&ghcmessages.Error{Message: &msg},
			//		), err
			//	}
			//}

			updatedProgearWeightTicket, _ := h.progearUpdater.UpdateProgearWeightTicket(appCtx, *progearWeightTicket, params.IfMatch)
			//if err != nil {
			//	return handleError(err)
			//}
			returnPayload := payloads.ProGearWeightTicket(h.FileStorer(), updatedProgearWeightTicket)
			return progearops.NewUpdateProGearWeightTicketOK().WithPayload(returnPayload), nil
		})
}
