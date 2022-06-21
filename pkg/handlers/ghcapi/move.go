package ghcapi

import (
	"errors"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"

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
	handlers.HandlerConfig
	services.MoveFetcher
}

// Handle handles the getMove by locator request
func (h GetMoveHandler) Handle(params moveop.GetMoveParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			locator := params.Locator
			if locator == "" {
				return moveop.NewGetMoveBadRequest(), apperror.NewBadDataError("missing required parameter: locator")
			}

			move, err := h.FetchMove(appCtx, locator, nil)

			if err != nil {
				appCtx.Logger().Error("Error retrieving move by locator", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return moveop.NewGetMoveNotFound(), err
				default:
					return moveop.NewGetMoveInternalServerError(), err
				}
			}

			payload := payloads.Move(move)
			return moveop.NewGetMoveOK().WithPayload(payload), nil
		})
}

type SearchMovesHandler struct {
	handlers.HandlerConfig
	services.MoveSearcher
}

func (h SearchMovesHandler) Handle(params moveop.SearchMovesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moves, err := h.MoveSearcher.SearchMoves(appCtx, params.Body.Locator, params.Body.DodID, params.Body.CustomerName)

			if err != nil {
				appCtx.Logger().Error("Error retrieving move by locator", zap.Error(err))
				return moveop.NewSearchMovesInternalServerError(), err
			}

			searchMoves := payloads.SearchMoves(moves)
			payload := &ghcmessages.SearchMovesResult{
				Page:        1,
				PerPage:     100,
				TotalCount:  int64(len(moves)),
				SearchMoves: *searchMoves,
			}
			return moveop.NewSearchMovesOK().WithPayload(payload), nil
		})
}

type SetFinancialReviewFlagHandler struct {
	handlers.HandlerConfig
	services.MoveFinancialReviewFlagSetter
}

// Handle flags a move for financial review
func (h SetFinancialReviewFlagHandler) Handle(params moveop.SetFinancialReviewFlagParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			moveID := uuid.FromStringOrNil(params.MoveID.String())

			remarks := params.Body.Remarks
			flagForReview := params.Body.FlagForReview
			if flagForReview == nil {
				badDataError := apperror.NewBadDataError("missing FlagForReview field")
				payload := payloadForValidationError("Unable to flag move for financial review", badDataError.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
				return moveop.NewSetFinancialReviewFlagUnprocessableEntity().WithPayload(payload), badDataError
			}
			// We require remarks when the move is going to be flagged for review.
			if *flagForReview && remarks == nil {
				badDataError := apperror.NewBadDataError("missing remarks field")
				payload := payloadForValidationError("Unable to flag move for financial review", badDataError.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
				return moveop.NewSetFinancialReviewFlagUnprocessableEntity().WithPayload(payload), badDataError
			}

			move, err := h.SetFinancialReviewFlag(appCtx, moveID, *params.IfMatch, *flagForReview, remarks)

			if err != nil {
				appCtx.Logger().Error("Error flagging move for financial review", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return moveop.NewSetFinancialReviewFlagNotFound(), err
				case apperror.PreconditionFailedError:
					return moveop.NewSetFinancialReviewFlagPreconditionFailed(), err
				case apperror.InvalidInputError:
					var e *apperror.InvalidInputError
					_ = errors.As(err, &e)
					payload := payloadForValidationError("Unable to flag move for financial review", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return moveop.NewSetFinancialReviewFlagUnprocessableEntity().WithPayload(payload), err
				default:
					return moveop.NewSetFinancialReviewFlagInternalServerError(), err
				}
			}

			payload := payloads.Move(move)
			return moveop.NewSetFinancialReviewFlagOK().WithPayload(payload), nil
		})
}
