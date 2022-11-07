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

// CreateProGearWeightTicketHandler
type CreateProGearWeightTicketHandler struct {
	handlers.HandlerConfig
	progearCreator services.ProgearWeightTicketCreator
}

// Handle creating a progear weight ticket
func (h CreateProGearWeightTicketHandler) Handle(params progearops.CreateProGearWeightTicketParams) middleware.Responder {
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

			progear, err := h.progearCreator.CreateProgearWeightTicket(appCtx, ppmShipmentID)

			if err != nil {
				appCtx.Logger().Error("internalapi.CreateProgearWeightTicketHandler", zap.Error(err))
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

// UpdateProGearWeightTicketHandler
type UpdateProGearWeightTicketHandler struct {
	handlers.HandlerConfig
	progearUpdater services.ProgearWeightTicketUpdater
}

func (h UpdateProGearWeightTicketHandler) Handle(params progearops.UpdateProGearWeightTicketParams) middleware.Responder {
	// track every request with middleware:
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return progearops.NewUpdateProGearWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return progearops.NewUpdateProGearWeightTicketForbidden(), noServiceMemberIDErr
			}

			payload := params.UpdateProGearWeightTicket
			if payload == nil {
				noBodyErr := apperror.NewBadDataError("Invalid weight ticket: params UpdateProgearPayload is nil")
				appCtx.Logger().Error(noBodyErr.Error())
				return progearops.NewUpdateProGearWeightTicketBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The Weight Ticket request payload cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest))), noBodyErr
			}

			weightTicket := payloads.ProgearWeightTicketModelFromUpdate(payload)
			weightTicket.ID = uuid.FromStringOrNil(params.ProGearWeightTicketID.String())

			updateProgear, err := h.progearUpdater.UpdateProgearWeightTicket(appCtx, *weightTicket, params.IfMatch)

			if err != nil {
				appCtx.Logger().Error("internalapi.UpdateProGearWeightTicketHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.InvalidInputError:
					return progearops.NewUpdateProGearWeightTicketUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.PreconditionFailedError:
					return progearops.NewUpdateProGearWeightTicketPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.NotFoundError:
					return progearops.NewUpdateProGearWeightTicketNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.
							Logger().
							Error(
								"internalapi.UpdateProGearWeightTicketHandler error",
								zap.Error(e.Unwrap()),
							)
					}
					return progearops.
						NewUpdateProGearWeightTicketInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				default:
					return progearops.
						NewUpdateProGearWeightTicketInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				}

			}
			returnPayload := payloads.ProGearWeightTicket(h.FileStorer(), updateProgear)
			return progearops.NewUpdateProGearWeightTicketOK().WithPayload(returnPayload), nil
		})
}
