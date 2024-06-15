package internalapi

import (
	"database/sql"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
)

// CreateMovingExpenseHandler

type CreateMovingExpenseHandler struct {
	handlers.HandlerConfig
	movingExpenseCreator services.MovingExpenseCreator
}

// Handle creates a moving expense
func (h CreateMovingExpenseHandler) Handle(params movingexpenseops.CreateMovingExpenseParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest, func(appCtx appcontext.AppContext) (middleware.Responder, error) {
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
			appCtx.Logger().Error("internalapi.CreateMovingExpenseHandler", zap.Error(err))
			switch e := err.(type) {
			case apperror.InvalidInputError:
				return movingexpenseops.NewCreateMovingExpenseUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage,
					h.GetTraceIDFromRequest(params.HTTPRequest),
					e.ValidationErrors)), err
			case apperror.NotFoundError:
				return movingexpenseops.NewCreateMovingExpenseNotFound(), err
			case apperror.QueryError:
				if e.Unwrap() != nil {
					// If you can unwrap, log the internal error (usually a pq error) for better debugging
					appCtx.Logger().Error("internalapi.CreateMovingExpenseHandler error", zap.Error(e.Unwrap()))
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

// UpdateMovingExpenseHandler struct
type UpdateMovingExpenseHandler struct {
	handlers.HandlerConfig
	movingExpenseUpdater services.MovingExpenseUpdater
}

func (h UpdateMovingExpenseHandler) Handle(params movingexpenseops.UpdateMovingExpenseParams) middleware.Responder {
	// track every request with middleware:
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				return movingexpenseops.NewUpdateMovingExpenseUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() && appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				return movingexpenseops.NewUpdateMovingExpenseForbidden(), noServiceMemberIDErr
			}

			payload := params.UpdateMovingExpense
			if payload == nil {
				noBodyErr := apperror.NewBadDataError("Invalid moving expense: params UpdateMovingExpense is nil")
				appCtx.Logger().Error(noBodyErr.Error())
				return movingexpenseops.NewUpdateMovingExpenseBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					"The moving expense request payload cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest))), noBodyErr
			}

			movingExpense := payloads.MovingExpenseModelFromUpdate(payload)
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
			updateMovingExpense, err := h.movingExpenseUpdater.UpdateMovingExpense(appCtx, *movingExpense, params.IfMatch)

			if err != nil {
				appCtx.Logger().Error("internalapi.UpdateMovingExpenseHandler", zap.Error(err))
				switch e := err.(type) {
				case apperror.InvalidInputError:
					return movingexpenseops.NewUpdateMovingExpenseUnprocessableEntity().WithPayload(payloads.ValidationError(handlers.ValidationErrMessage, h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.PreconditionFailedError:
					return movingexpenseops.NewUpdateMovingExpensePreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.NotFoundError:
					return movingexpenseops.NewUpdateMovingExpenseNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.
							Logger().
							Error(
								"internalapi.UpdateMovingExpenseHandler error",
								zap.Error(e.Unwrap()),
							)
					}
					return movingexpenseops.
						NewUpdateMovingExpenseInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				default:
					return movingexpenseops.
						NewUpdateMovingExpenseInternalServerError().
						WithPayload(
							payloads.InternalServerError(
								nil,
								h.GetTraceIDFromRequest(params.HTTPRequest),
							),
						), err
				}

			}
			returnPayload := payloads.MovingExpense(h.FileStorer(), updateMovingExpense)
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
			if appCtx.Session() == nil {
				noSessionErr := apperror.NewSessionError("No user session")
				appCtx.Logger().Error("internalapi.DeleteMovingExpenseHandler", zap.Error(noSessionErr))
				return movingexpenseops.NewDeleteMovingExpenseUnauthorized(), noSessionErr
			}

			if !appCtx.Session().IsMilApp() || appCtx.Session().ServiceMemberID == uuid.Nil {
				noServiceMemberIDErr := apperror.NewSessionError("No service member ID")
				appCtx.Logger().Error("internalapi.DeleteMovingExpenseHandler", zap.Error(noServiceMemberIDErr))
				return movingexpenseops.NewDeleteMovingExpenseForbidden(), noServiceMemberIDErr
			}

			// Make sure the service member is not modifying another service member's PPM
			ppmID := uuid.FromStringOrNil(params.PpmShipmentID.String())

			movingExpenseID := uuid.FromStringOrNil(params.MovingExpenseID.String())

			err := h.movingExpenseDeleter.DeleteMovingExpense(appCtx, ppmID, movingExpenseID)
			if err != nil {
				appCtx.Logger().Error("internalapi.DeleteMovingExpenseHandler", zap.Error(err))

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
