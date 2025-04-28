package ghcapi

import (
	"database/sql"
	"fmt"

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

// CreateMovingExpenseHandler

type CreateMovingExpenseHandler struct {
	handlers.HandlerConfig
	movingExpenseCreator services.MovingExpenseCreator
}

// Handle creates a moving expense
func (h CreateMovingExpenseHandler) Handle(params movingexpenseops.CreateMovingExpenseParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest, func(appCtx appcontext.AppContext) (middleware.Responder, error) {

		/** Feature Flag - COMPLETE_PPM_CLOSEOUT_FOR_CUSTOMER **/
		const featureFlagNameCloseoutForCustomer = "complete_ppm_closeout_for_customer"
		isCloseoutForCustomerFeatureOn := false
		flag, ffErr := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagNameCloseoutForCustomer, map[string]string{})

		if ffErr != nil {
			appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagNameCloseoutForCustomer), zap.Error(ffErr))
		} else {
			isCloseoutForCustomerFeatureOn = flag.Match
		}

		if !isCloseoutForCustomerFeatureOn {
			return movingexpenseops.NewCreateMovingExpenseUnprocessableEntity().WithPayload(payloadForValidationError(
				"Unable to create a moving expense", "Moving expenses cannot be created unless the complete_ppm_closeout_for_customer feature flag is enabled.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), nil
		}

		if appCtx.Session() == nil {
			noSessionErr := apperror.NewSessionError("No user session")
			return movingexpenseops.NewCreateMovingExpenseUnauthorized(), noSessionErr
		}

		// No need for payload_to_model for Create
		ppmShipmentID, err := uuid.FromString(params.PpmShipmentID.String())
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, apperror.NewNotFoundError(ppmShipmentID, "Incorrect PPMShipmentID")
			default:
				appCtx.Logger().Error("missing PPM Shipment ID", zap.Error(err))
				return movingexpenseops.NewCreateMovingExpenseBadRequest(), nil
			}
		}

		movingExpense, err := h.movingExpenseCreator.CreateMovingExpense(appCtx, ppmShipmentID)

		if err != nil {
			appCtx.Logger().Error("ghcapi.CreateMovingExpenseHandler", zap.Error(err))
			switch e := err.(type) {
			case apperror.InvalidInputError:
				return movingexpenseops.NewCreateMovingExpenseUnprocessableEntity().WithPayload(
					payloadForValidationError(
						handlers.ValidationErrMessage,
						err.Error(),
						h.GetTraceIDFromRequest(params.HTTPRequest),
						e.ValidationErrors,
					),
				), err
			case apperror.NotFoundError:
				return movingexpenseops.NewCreateMovingExpenseNotFound(), err
			case apperror.QueryError:
				if e.Unwrap() != nil {
					// If you can unwrap, log the internal error (usually a pq error) for better debugging
					appCtx.Logger().Error("ghcapi.CreateMovingExpenseHandler error", zap.Error(e.Unwrap()))
				}
				return movingexpenseops.
					NewCreateMovingExpenseInternalServerError().
					WithPayload(
						payloads.InternalServerError(
							nil,
							h.GetTraceIDFromRequest(params.HTTPRequest),
						),
					), err
			default:
				return movingexpenseops.NewCreateMovingExpenseInternalServerError().WithPayload(
					payloads.InternalServerError(
						nil,
						h.GetTraceIDFromRequest(params.HTTPRequest),
					),
				), err
			}
		}

		// Add to payload
		returnPayload := payloads.MovingExpense(h.FileStorer(), movingExpense)
		return movingexpenseops.NewCreateMovingExpenseCreated().WithPayload(returnPayload), nil
	})
}

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

// DeleteMovingExpenseHandler
type DeleteMovingExpenseHandler struct {
	handlers.HandlerConfig
	movingExpenseDeleter services.MovingExpenseDeleter
}

// Handle deletes a moving expense
func (h DeleteMovingExpenseHandler) Handle(params movingexpenseops.DeleteMovingExpenseParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			/** Feature Flag - COMPLETE_PPM_CLOSEOUT_FOR_CUSTOMER **/
			const featureFlagNameCloseoutForCustomer = "complete_ppm_closeout_for_customer"
			isCloseoutForCustomerFeatureOn := false
			flag, ffErr := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagNameCloseoutForCustomer, map[string]string{})

			if ffErr != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagNameCloseoutForCustomer), zap.Error(ffErr))
			} else {
				isCloseoutForCustomerFeatureOn = flag.Match
			}

			if !isCloseoutForCustomerFeatureOn {
				return movingexpenseops.NewDeleteMovingExpenseUnprocessableEntity().WithPayload(payloadForValidationError(
					"Unable to delete a moving expense", "Moving expenses cannot be deleted unless the complete_ppm_closeout_for_customer feature flag is enabled.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), nil
			}

			errInstance := fmt.Sprintf("Instance: %s", h.GetTraceIDFromRequest(params.HTTPRequest))
			errPayload := &ghcmessages.Error{Message: &errInstance}
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				appCtx.Logger().Error("ghcapi.DeleteMovingExpenseHandler", zap.Error(noSessionErr))
				return movingexpenseops.NewDeleteMovingExpenseUnauthorized(), noSessionErr
			}
			if !appCtx.Session().IsOfficeApp() {
				return movingexpenseops.NewDeleteMovingExpenseForbidden().WithPayload(errPayload), apperror.NewSessionError("Request should come from the office app.")
			}

			ppmID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			movingExpenseID := uuid.FromStringOrNil(params.MovingExpenseID.String())

			err := h.movingExpenseDeleter.DeleteMovingExpense(appCtx, ppmID, movingExpenseID)
			if err != nil {
				appCtx.Logger().Error("ghcapi.DeleteMovingExpenseHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return movingexpenseops.NewDeleteMovingExpenseNotFound(), err
				case apperror.ConflictError:
					return movingexpenseops.NewDeleteMovingExpenseConflict(), err
				case apperror.ForbiddenError:
					return movingexpenseops.NewDeleteMovingExpenseForbidden(), err
				case apperror.UnprocessableEntityError:
					return movingexpenseops.NewDeleteMovingExpenseUnprocessableEntity(), err
				default:
					return movingexpenseops.NewDeleteMovingExpenseInternalServerError(), err
				}
			}

			return movingexpenseops.NewDeleteMovingExpenseNoContent(), nil
		})
}
