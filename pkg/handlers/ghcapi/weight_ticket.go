package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	weightticketops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// CreateWeightTicketHandler
type CreateWeightTicketHandler struct {
	handlers.HandlerConfig
	weightTicketCreator services.WeightTicketCreator
}

// Handle creates a weight ticket
// Depending on the SO, may need to change the document params to weight ticket params
func (h CreateWeightTicketHandler) Handle(params weightticketops.CreateWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return weightticketops.NewCreateWeightTicketUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsOfficeApp() {
				return weightticketops.NewUpdateMovingExpenseForbidden(), apperror.NewSessionError("Request should come from the office app.")
			}

			// NO NEED FOR payload_to_model, will need for Update
			ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
			if err != nil {
				appCtx.Logger().Error("missing PPM Shipment ID", zap.Error(err))
				return weightticketops.NewCreateWeightTicketBadRequest(), nil
			}

			weightTicket, err := h.weightTicketCreator.CreateWeightTicket(appCtx, ppmShipmentID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.CreateWeightTicketHandler", zap.Error(err))
				// Can get a status error
				// Can get an DB error - does the weight ticket, doc create?
				// Can get an error for whether the PPM exist
				switch err.(type) {
				case apperror.InvalidInputError:
					return weightticketops.NewCreateWeightTicketUnprocessableEntity(), err
				case apperror.ForbiddenError:
					return weightticketops.NewCreateWeightTicketForbidden(), err
				case apperror.NotFoundError:
					return weightticketops.NewCreateWeightTicketNotFound(), err
				default:
					return weightticketops.NewCreateWeightTicketInternalServerError(), err
				}
			}
			returnPayload := payloads.WeightTicket(h.FileStorer(), weightTicket)
			return weightticketops.NewCreateWeightTicketOK().WithPayload(returnPayload), nil
		})
}

// UpdateWeightTicketHandler
type UpdateWeightTicketHandler struct {
	handlers.HandlerConfig
	weighTicketUpdater services.WeightTicketUpdater
}

func (h UpdateWeightTicketHandler) Handle(params weightticketops.UpdateWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.UpdateWeightTicketPayload
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}
			weightTicket := payloads.WeightTicketModelFromUpdate(payload)

			if !appCtx.Session().IsOfficeApp() {
				return weightticketops.NewUpdateWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			weightTicket.ID = uuid.FromStringOrNil(params.WeightTicketID.String())

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.UpdateWeightTicketHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return weightticketops.NewUpdateWeightTicketNotFound(), err
				case apperror.InvalidInputError:
					return weightticketops.NewUpdateWeightTicketUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					), err
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return weightticketops.NewUpdateWeightTicketPreconditionFailed().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.QueryError:
					return weightticketops.NewUpdateWeightTicketInternalServerError(), err
				default:
					return weightticketops.NewUpdateWeightTicketInternalServerError(), err
				}
			}

			updatedWeightTicket, err := h.weighTicketUpdater.UpdateWeightTicket(appCtx, *weightTicket, params.IfMatch)

			if err != nil {
				return handleError(err)
			}

			returnPayload := payloads.WeightTicket(h.FileStorer(), updatedWeightTicket)

			return weightticketops.NewUpdateWeightTicketOK().WithPayload(returnPayload), nil
		})
}

// DeleteWeightTicketHandler
type DeleteWeightTicketHandler struct {
	handlers.HandlerConfig
	weightTicketDeleter services.WeightTicketDeleter
}

func (h DeleteWeightTicketHandler) Handle(params weightticketops.DeleteWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return weightticketops.NewDeleteWeightTicketForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			weightTicketID := uuid.FromStringOrNil(params.WeightTicketID.String())
			ppmID := uuid.FromStringOrNil(string(params.PpmShipmentID.String()))

			handleError := func(err error) (middleware.Responder, error) {
				appCtx.Logger().Error("ghcapi.DeleteWeightTicketHandler", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return weightticketops.NewDeleteWeightTicketNotFound(), err
				case apperror.InvalidInputError:
					return weightticketops.NewDeleteWeightTicketUnprocessableEntity().WithPayload(
						payloadForValidationError(
							handlers.ValidationErrMessage,
							err.Error(),
							h.GetTraceIDFromRequest(params.HTTPRequest),
							e.ValidationErrors,
						),
					), err
				case apperror.PreconditionFailedError:
					msg := fmt.Sprintf("%v | Instance: %v", handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))
					return weightticketops.NewDeleteWeightTicketPreconditionFailed().WithPayload(
						&ghcmessages.Error{Message: &msg},
					), err
				case apperror.QueryError:
					return weightticketops.NewDeleteWeightTicketInternalServerError(), err
				default:
					return weightticketops.NewDeleteWeightTicketInternalServerError(), err
				}
			}

			err := h.weightTicketDeleter.DeleteWeightTicket(appCtx, ppmID, weightTicketID)

			if err != nil {
				return handleError(err)
			}

			return weightticketops.NewDeleteWeightTicketNoContent(), nil
		})
}
