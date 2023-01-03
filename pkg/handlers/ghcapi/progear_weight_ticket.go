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
