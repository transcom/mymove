package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	progearops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// CreateProgearHandler
type CreateProgearHandler struct {
	handlers.HandlerConfig
	progearCreator services.ProgearCreator
}

// Handle creating a progear weight ticket
func (h CreateProgearHandler) Handle(params progearops.CreateProGearWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return progearops.NewCreateProGearWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return progearops.NewCreateProGearWeightTicketForbidden(), noServiceMemberIDErr
			}

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				appCtx.Logger().Error("missing PPM Shipment ID", zap.Error(err))
				return progearops.NewCreateProGearWeightTicketBadRequest(), nil
			}

			progear, err := h.progearCreator.CreateProgear(appCtx, ppmShipmentID)

			if err != nil {
				appCtx.Logger().Error("internalapi.CreateProgearHandler", zap.Error(err))
				switch err.(type) {
				case apperror.InvalidInputError:
					return progearops.NewCreateProGearWeightTicketUnprocessableEntity(), err
				case apperror.ForbiddenError:
					return progearops.NewCreateProGearWeightTicketForbidden(), err
				case apperror.NotFoundError:
					return progearops.NewCreateProGearWeightTicketNotFound(), err
				default:
					return progearops.NewCreateProGearWeightTicketInternalServerError(), err
				}
			}
			returnPayload := payloads.ProGearWeightTicket(h.FileStorer(), progear)
			return progearops.NewCreateProGearWeightTicketOK().WithPayload(returnPayload), nil
		})
}

// // UpdateProgearHandler
// type UpdateProgearHandler struct {
// 	handlers.HandlerConfig
// 	weightTicketUpdater services.ProgearUpdater
// }

// func (h UpdateProgearHandler) Handle(params weightticketops.UpdateProgearParams) middleware.Responder {
// 	// track every request with middleware:
// 	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
// 		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

// 			if appCtx.Session() == nil {
// 				noSessionErr := apperror.NewSessionError("No user session")
// 				return weightticketops.NewUpdateProgearUnauthorized(), noSessionErr
// 			}

// 			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
// 				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
// 				return weightticketops.NewUpdateProgearForbidden(), noServiceMemberIDErr
// 			}

// 			payload := params.UpdateProgearPayload
// 			if payload == nil {
// 				noBodyErr := apperror.NewBadDataError("Invalid weight ticket: params UpdateProgearPayload is nil")
// 				appCtx.Logger().Error(noBodyErr.Error())
// 				return weightticketops.NewUpdateProgearBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
// 					"The Weight Ticket request payload cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest))), noBodyErr
// 			}

// 			weightTicket := payloads.ProgearModelFromUpdate(payload)
// 			weightTicket.ID = uuid.FromStringOrNil(params.ProgearID.String())

// 			updateProgear, err := h.weightTicketUpdater.UpdateProgear(appCtx, *weightTicket, params.IfMatch)

// 			if err != nil {
// 				appCtx.Logger().Error("internalapi.UpdateProgearHandler", zap.Error(err))
// 				switch e := err.(type) {
// 				case apperror.InvalidInputError:
// 					return weightticketops.NewUpdateProgearUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
// 				case apperror.PreconditionFailedError:
// 					return weightticketops.NewUpdateProgearPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
// 				case apperror.NotFoundError:
// 					return weightticketops.NewUpdateProgearNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
// 				case apperror.QueryError:
// 					if e.Unwrap() != nil {
// 						// If you can unwrap, log the internal error (usually a pq error) for better debugging
// 						appCtx.
// 							Logger().
// 							Error(
// 								"internalapi.UpdateProgearHandler error",
// 								zap.Error(e.Unwrap()),
// 							)
// 					}
// 					return weightticketops.
// 						NewUpdateProgearInternalServerError().
// 						WithPayload(
// 							payloads.InternalServerError(
// 								nil,
// 								h.GetTraceIDFromRequest(params.HTTPRequest),
// 							),
// 						), err
// 				default:
// 					return weightticketops.
// 						NewUpdateProgearInternalServerError().
// 						WithPayload(
// 							payloads.InternalServerError(
// 								nil,
// 								h.GetTraceIDFromRequest(params.HTTPRequest),
// 							),
// 						), err
// 				}

// 			}
// 			returnPayload := payloads.Progear(h.FileStorer(), updateProgear)
// 			return weightticketops.NewUpdateProgearOK().WithPayload(returnPayload), nil
// 		})
// }
