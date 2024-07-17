package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
	ppmshipment "github.com/transcom/mymove/pkg/services/ppmshipment"
)

// UpdateMovingExpenseHandler
type UpdateMovingExpenseHandler struct {
	handlers.HandlerConfig
	movingExpenseUpdater services.MovingExpenseUpdater
}

func (h UpdateMovingExpenseHandler) Handle(params movingexpenseops.UpdateMovingExpenseParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest, func(appCtx appcontext.AppContext) (middleware.Responder, error) {
		if !appCtx.Session().IsOfficeApp() {
			return movingexpenseops.NewUpdateMovingExpenseForbidden(), apperror.NewSessionError("Request should come from the office app.")
		}

		movingExpense := payloads.MovingExpenseModelFromUpdate(params.UpdateMovingExpense)

		movingExpense.ID = uuid.FromStringOrNil(params.MovingExpenseID.String())
		movingExpense.PPMShipmentID = uuid.FromStringOrNil(params.PpmShipmentID.String())
		ppmEagerAssociations := []string{"PickupAddress",
			"DestinationAddress",
			"SecondaryPickupAddress",
			"SecondaryDestinationAddress",
		}
		ppmShipmentFetcher := ppmshipment.NewPPMShipmentFetcher()
		ppmShipment, ppmShipmentErr := ppmShipmentFetcher.GetPPMShipment(appCtx, movingExpense.PPMShipmentID, ppmEagerAssociations, nil)

		if ppmShipmentErr != nil {
			return nil, ppmShipmentErr
		}

		movingExpense.PPMShipment = *ppmShipment
		updatedMovingExpense, err := h.movingExpenseUpdater.UpdateMovingExpense(appCtx, *movingExpense, params.IfMatch)

		if err != nil {
			appCtx.Logger().Error("ghcapi.UpdateMovingExpenseHandler error", zap.Error(err))

			switch e := err.(type) {
			case apperror.NotFoundError:
				return movingexpenseops.NewUpdateMovingExpenseNotFound(), nil
			case apperror.QueryError:
				if e.Unwrap() != nil {
					// If you can unwrap, log the error (usually a pq error) for better debugging
					appCtx.Logger().Error(
						"ghcapi.UpdateMovingExpenseHandler error",
						zap.Error(e.Unwrap()),
					)
				}

				return movingexpenseops.NewUpdateMovingExpenseInternalServerError(), nil
			case apperror.PreconditionFailedError:
				return movingexpenseops.NewUpdateMovingExpensePreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), nil
			case apperror.InvalidInputError:
				return movingexpenseops.NewUpdateMovingExpenseUnprocessableEntity().WithPayload(
					payloadForValidationError(
						handlers.ValidationErrMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors,
					),
				), nil
			default:
				return movingexpenseops.NewUpdateMovingExpenseInternalServerError(), nil
			}
		}

		returnPayload := payloads.MovingExpense(h.FileStorer(), updatedMovingExpense)

		return movingexpenseops.NewUpdateMovingExpenseOK().WithPayload(returnPayload), nil
	})
}
