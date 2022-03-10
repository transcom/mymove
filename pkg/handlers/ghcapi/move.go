package ghcapi

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"go.uber.org/zap"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"

	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// GetMoveHandler gets a move by locator
type GetMoveHandler struct {
	handlers.HandlerContext
	services.MoveFetcher
}

// Handle handles the getMove by locator request
func (h GetMoveHandler) Handle(params moveop.GetMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			locator := params.Locator
			if locator == "" {
				return moveop.NewGetMoveBadRequest()
			}

			move, err := h.FetchMove(appCtx, locator, nil)

			if err != nil {
				appCtx.Logger().Error("Error retrieving move by locator", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return moveop.NewGetMoveNotFound()
				default:
					return moveop.NewGetMoveInternalServerError()
				}
			}

			payload := payloads.Move(move)
			return moveop.NewGetMoveOK().WithPayload(payload)
		})
}

type SetFinancialReviewFlagHandler struct {
	handlers.HandlerContext
	services.MoveFinancialReviewFlagSetter
}

// Handle flags a move for financial review
func (h SetFinancialReviewFlagHandler) Handle(params moveop.SetFinancialReviewFlagParams) middleware.Responder {
	return h.AuditableAppContextFromRequest(params.HTTPRequest,
		func(appCtx appcontext.AppContext) middleware.Responder {

			moveID := uuid.FromStringOrNil(params.MoveID.String())

			remarks := params.Body.Remarks
			flagForReview := params.Body.FlagForReview
			if flagForReview == nil {
				payload := payloadForValidationError("Unable to flag move for financial review", "missing FlagForReview field", h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
				return moveop.NewSetFinancialReviewFlagUnprocessableEntity().WithPayload(payload)
			}
			// We require remarks when the move is going to be flagged for review.
			if *flagForReview && remarks == nil {
				payload := payloadForValidationError("Unable to flag move for financial review", "missing remarks field", h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
				return moveop.NewSetFinancialReviewFlagUnprocessableEntity().WithPayload(payload)
			}

			move, err := h.SetFinancialReviewFlag(appCtx, moveID, *params.IfMatch, *flagForReview, remarks)

			if err != nil {
				appCtx.Logger().Error("Error flagging move for financial review", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return moveop.NewSetFinancialReviewFlagNotFound()
				case apperror.PreconditionFailedError:
					return moveop.NewSetFinancialReviewFlagPreconditionFailed()
				case apperror.InvalidInputError:
					var e *apperror.InvalidInputError
					_ = errors.As(err, &e)
					payload := payloadForValidationError("Unable to flag move for financial review", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return moveop.NewSetFinancialReviewFlagUnprocessableEntity().WithPayload(payload)
				default:
					return moveop.NewSetFinancialReviewFlagInternalServerError()
				}
			}

			payload := payloads.Move(move)
			return moveop.NewSetFinancialReviewFlagOK().WithPayload(payload)
		})
}
