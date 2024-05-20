package ghcapi

import (
	"errors"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	moveop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// GetMoveHandler gets a move by locator
type GetMoveHandler struct {
	handlers.HandlerConfig
	services.MoveFetcher
	services.MoveLocker
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

			privileges, err := models.FetchPrivilegesForUser(appCtx.DB(), appCtx.Session().UserID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}

			// if this user is accessing the move record, we need to lock it so others can't edit it
			// to allow for locking a move, we need to look at these things
			// 1. Is the user an office user?
			// 2. Are the columns empty (lock_expires_at & locked_by) in the db?
			// 3. Is the lock_expires_at after right now?
			// 4. Is the current user the one that locked it? This will reset the locked_at time.
			// if all of those questions have the answer "yes", then we will proceed with locking the move by the current user
			officeUserID := appCtx.Session().OfficeUserID
			lockedOfficeUserID := move.LockedByOfficeUserID
			lockExpiresAt := move.LockExpiresAt
			now := time.Now()
			if appCtx.Session().IsOfficeUser() {
				if move.LockedByOfficeUserID == nil && move.LockExpiresAt == nil || (lockExpiresAt != nil && now.After(*lockExpiresAt)) || (*lockedOfficeUserID == officeUserID && lockedOfficeUserID != nil) {
					move, err = h.LockMove(appCtx, move, officeUserID)
					if err != nil {
						return moveop.NewGetMoveInternalServerError(), err
					}
				}
			}

			if move.Orders.OrdersType == "SAFETY" && !privileges.HasPrivilege(models.PrivilegeTypeSafety) {
				appCtx.Logger().Error("Invalid permissions")
				return moveop.NewGetMoveNotFound(), nil
			} else {
				payload := payloads.Move(move)
				return moveop.NewGetMoveOK().WithPayload(payload), nil
			}
		})
}

type SearchMovesHandler struct {
	handlers.HandlerConfig
	services.MoveSearcher
}

func (h SearchMovesHandler) Handle(params moveop.SearchMovesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			searchMovesParams := services.SearchMovesParams{
				Branch:                params.Body.Branch,
				Locator:               params.Body.Locator,
				DodID:                 params.Body.DodID,
				CustomerName:          params.Body.CustomerName,
				DestinationPostalCode: params.Body.DestinationPostalCode,
				OriginPostalCode:      params.Body.OriginPostalCode,
				Status:                params.Body.Status,
				ShipmentsCount:        params.Body.ShipmentsCount,
				Page:                  params.Body.Page,
				PerPage:               params.Body.PerPage,
				Sort:                  params.Body.Sort,
				Order:                 params.Body.Order,
				PickupDate:            handlers.FmtDateTimePtrToPopPtr(params.Body.PickupDate),
				DeliveryDate:          handlers.FmtDateTimePtrToPopPtr(params.Body.DeliveryDate),
			}

			moves, totalCount, err := h.MoveSearcher.SearchMoves(appCtx, &searchMovesParams)

			if err != nil {
				appCtx.Logger().Error("Error searching for move", zap.Error(err))
				return moveop.NewSearchMovesInternalServerError(), err
			}
			searchMoves := payloads.SearchMoves(appCtx, moves)
			payload := &ghcmessages.SearchMovesResult{
				Page:        searchMovesParams.Page,
				PerPage:     searchMovesParams.PerPage,
				TotalCount:  int64(totalCount),
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

type UpdateMoveCloseoutOfficeHandler struct {
	handlers.HandlerConfig
	services.MoveCloseoutOfficeUpdater
}

func (h UpdateMoveCloseoutOfficeHandler) Handle(params moveop.UpdateCloseoutOfficeParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			closeoutOfficeID := uuid.FromStringOrNil(params.Body.CloseoutOfficeID.String())

			move, err := h.MoveCloseoutOfficeUpdater.UpdateCloseoutOffice(appCtx, params.Locator, closeoutOfficeID, params.IfMatch)
			if err != nil {
				appCtx.Logger().Error("UpdateMoveCloseoutOfficeHandler error", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return moveop.NewUpdateCloseoutOfficeNotFound(), err
				case apperror.PreconditionFailedError:
					return moveop.NewUpdateCloseoutOfficePreconditionFailed(), err
				case apperror.InvalidInputError:
					return moveop.NewUpdateCloseoutOfficeUnprocessableEntity(), err
				default:
					return moveop.NewUpdateCloseoutOfficeInternalServerError(), err
				}
			}

			return moveop.NewUpdateCloseoutOfficeOK().WithPayload(payloads.Move(move)), nil
		})
}
