package internalapi

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	progearops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
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

			if returnPayload == nil {
				return progearops.NewCreateProGearWeightTicketInternalServerError(), err
			}
			return progearops.NewCreateProGearWeightTicketCreated().WithPayload(returnPayload), nil
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
			return progearops.NewUpdateProGearWeightTicketCreated().WithPayload(returnPayload), nil
		})
}

// DeleteProGearWeightTicketHandler
type DeleteProGearWeightTicketHandler struct {
	handlers.HandlerConfig
	progearDeleter services.ProgearWeightTicketDeleter
}

// Handle deletes a pro-gear weight ticket
func (h DeleteProGearWeightTicketHandler) Handle(params progearops.DeleteProGearWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				appCtx.Logger().Error("internalapi.DeleteProgearWeightTicketHandler", zap.Error(noSessionErr))
				return progearops.NewDeleteProGearWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() || appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				appCtx.Logger().Error("internalapi.DeleteProgearWeightTicketHandler", zap.Error(noServiceMemberIDErr))
				return progearops.NewDeleteProGearWeightTicketForbidden(), noServiceMemberIDErr
			}

			// Make sure the service member is not modifying another service member's PPM
			ppmID := uuid.FromStringOrNil(params.PpmShipmentID.String())
			var ppmShipment models.PPMShipment
			err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
				EagerPreload(
					"Shipment.MoveTaskOrder.Orders",
					"ProgearWeightTickets",
				).
				Find(&ppmShipment, ppmID)
			if err != nil {
				if err == sql.ErrNoRows {
					return progearops.NewDeleteWeightTicketNotFound(), err
				}
				return progearops.NewDeleteProGearWeightTicketInternalServerError(), err
			}
			if ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID != appCtx.Session().ServiceMemberID {
				wrongServiceMemberIDErr := apperror.NewSessionError("Attempted delete by wrong service member")
				appCtx.Logger().Error("internalapi.DeleteProgearWeightTicketHandler", zap.Error(wrongServiceMemberIDErr))
				return progearops.NewDeleteProGearWeightTicketForbidden(), wrongServiceMemberIDErr
			}
			progearWeightTicketID := uuid.FromStringOrNil(params.ProGearWeightTicketID.String())
			found := false
			for _, lineItem := range ppmShipment.ProgearWeightTickets {
				if lineItem.ID == progearWeightTicketID {
					found = true
					break
				}
			}
			if !found {
				mismatchedPPMShipmentAndProgearWeightTicketIDErr := apperror.NewSessionError("Pro-gear weight ticket does not exist on ppm shipment")
				appCtx.Logger().Error("internalapi.DeleteProGearWeightTicketHandler", zap.Error(mismatchedPPMShipmentAndProgearWeightTicketIDErr))
				return progearops.NewDeleteProGearWeightTicketNotFound(), mismatchedPPMShipmentAndProgearWeightTicketIDErr
			}

			err = h.progearDeleter.DeleteProgearWeightTicket(appCtx, progearWeightTicketID)
			if err != nil {
				appCtx.Logger().Error("internalapi.DeleteProgearWeightTicketHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return progearops.NewDeleteProGearWeightTicketNotFound(), err
				case apperror.ConflictError:
					return progearops.NewDeleteProGearWeightTicketConflict(), err
				case apperror.ForbiddenError:
					return progearops.NewDeleteProGearWeightTicketForbidden(), err
				case apperror.UnprocessableEntityError:
					return progearops.NewDeleteProGearWeightTicketUnprocessableEntity(), err
				default:
					return progearops.NewDeleteProGearWeightTicketInternalServerError(), err
				}
			}

			return progearops.NewDeleteProGearWeightTicketNoContent(), nil
		})
}
