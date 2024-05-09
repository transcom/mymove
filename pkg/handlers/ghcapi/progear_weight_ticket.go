package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	progearops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
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
				return progearops.NewCreateWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsOfficeApp() {
				return progearops.NewUpdateMovingExpenseForbidden(), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				appCtx.Logger().Error("missing PPM Shipment ID", zap.Error(err))
				return progearops.NewCreateProGearWeightTicketBadRequest(), nil
			}

			progear, err := h.progearCreator.CreateProgearWeightTicket(appCtx, ppmShipmentID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.CreateProgearWeightTicketHandler", zap.Error(err))
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

// UpdateProgearWeightTicketHandler
type UpdateProgearWeightTicketHandler struct {
	handlers.HandlerConfig
	progearUpdater services.ProgearWeightTicketUpdater
}

func (h UpdateProgearWeightTicketHandler) Handle(params progearops.UpdateProGearWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.UpdateProGearWeightTicket
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

			progearWeightTicket := payloads.ProgearWeightTicketModelFromUpdate(payload)

			if !appCtx.Session().IsOfficeApp() {
				return progearops.NewUpdateProGearWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			progearWeightTicket.ID = uuid.FromStringOrNil(params.ProGearWeightTicketID.String())

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.UpdateWeightTicketHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return progearops.NewUpdateProGearWeightTicketNotFound(), err
				case apperror.InvalidInputError:
					return progearops.NewUpdateProGearWeightTicketUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					), err
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return progearops.NewUpdateProGearWeightTicketPreconditionFailed().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the error (usually a pq error) for better debugging
						appCtx.Logger().Error(
							"ghcapi.GetWeightTicketsHandler error",
							zap.Error(e.Unwrap()),
						)
					}
					return progearops.NewUpdateProGearWeightTicketInternalServerError(), err
				default:
					return progearops.NewUpdateProGearWeightTicketInternalServerError(), err
				}
			}

			updatedProgearWeightTicket, err := h.progearUpdater.UpdateProgearWeightTicket(appCtx, *progearWeightTicket, params.IfMatch)

			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.ProGearWeightTicket(h.FileStorer(), updatedProgearWeightTicket)
			return progearops.NewUpdateProGearWeightTicketOK().WithPayload(returnPayload), nil
		})
}

// DeleteProgearWeightTicketHandler
type DeleteProgearWeightTicketHandler struct {
	handlers.HandlerConfig
	progearDeleter services.ProgearWeightTicketDeleter
}

func (h DeleteProgearWeightTicketHandler) Handle(params progearops.DeleteProGearWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return progearops.NewDeleteProGearWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			progearWeightTicketID := uuid.FromStringOrNil(string(params.ProGearWeightTicketID))
			ppmID := uuid.FromStringOrNil(string(params.PpmShipmentID))

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.DeleteProgearWeightTicketHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return progearops.NewDeleteProGearWeightTicketNotFound(), err
				case apperror.InvalidInputError:
					return progearops.NewDeleteProGearWeightTicketUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					), err
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return progearops.NewDeleteProGearWeightTicketPreconditionFailed().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.QueryError:
					return progearops.NewDeleteProGearWeightTicketInternalServerError(), err
				default:
					return progearops.NewDeleteProGearWeightTicketInternalServerError(), err
				}
			}

			err := h.progearDeleter.DeleteProgearWeightTicket(appCtx, ppmID, progearWeightTicketID)

			if err != nil {
				return handleError(err)
			}

			return progearops.NewDeleteProGearWeightTicketNoContent(), nil
		})
}
