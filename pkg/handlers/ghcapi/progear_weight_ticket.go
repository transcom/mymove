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
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return progearops.NewCreateProGearWeightTicketUnauthorized(), noSessionErr
			}
			if !appCtx.Session().IsOfficeApp() {
				return progearops.NewCreateProGearWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
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
				appCtx.Logger().Error("Returned Payload is empty", zap.Error(err))
				return progearops.NewCreateProGearWeightTicketInternalServerError(), nil
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

// DeleteProGearWeightTicketHandler
type DeleteProGearWeightTicketHandler struct {
	handlers.HandlerConfig
	progearDeleter services.ProgearWeightTicketDeleter
}

// Handle deletes a pro-gear weight ticket
func (h DeleteProGearWeightTicketHandler) Handle(params progearops.DeleteProGearWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return progearops.NewDeleteProGearWeightTicketUnauthorized(), noSessionErr
			}
			if !appCtx.Session().IsOfficeApp() {
				return progearops.NewDeleteProGearWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			progearWeightTicketID := uuid.FromStringOrNil(params.ProGearWeightTicketID.String())
			err := h.progearDeleter.DeleteProgearWeightTicket(appCtx, ppmID, progearWeightTicketID)
			if err != nil {
				appCtx.Logger().Error("ghcapi.DeleteProgearWeightTicketHandler", zap.Error(err))

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
