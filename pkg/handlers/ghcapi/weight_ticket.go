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

// GetWeightTicketsHandler is the handler that fetches all of the weight tickets for a PPM shipment for the office api
type GetWeightTicketsHandler struct {
	handlers.HandlerConfig
	services.WeightTicketFetcher
}

// Handle retrieves all weight tickets for a PPM shipment
func (h GetWeightTicketsHandler) Handle(params weightticketops.GetWeightTicketsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))

			errPayload := &ghcmessages.Error{Message: &errInstance}

			if !appCtx.Session().IsOfficeApp() {
				return weightticketops.NewGetWeightTicketsForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmShipmentID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			weightTickets, err := h.WeightTicketFetcher.ListWeightTickets(appCtx, ppmShipmentID)

			if err != nil {
				appCtx.Logger().Error("ghcapi.GetWeightTicketsHandler error", zap.Error(err))

				switch e := err.(type) {
				case apperror.ForbiddenError:
					return weightticketops.NewGetWeightTicketsForbidden().WithPayload(errPayload), nil
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the error (usually a pq error) for better debugging
						appCtx.Logger().Error(
							"ghcapi.GetWeightTicketsHandler error",
							zap.Error(e.Unwrap()),
						)
					}

					return weightticketops.NewGetWeightTicketsInternalServerError().WithPayload(errPayload), nil
				default:
					return weightticketops.NewGetWeightTicketsInternalServerError().WithPayload(errPayload), nil
				}
			}

			returnPayload := payloads.WeightTickets(h.FileStorer(), weightTickets)

			return weightticketops.NewGetWeightTicketsOK().WithPayload(returnPayload), nil
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
