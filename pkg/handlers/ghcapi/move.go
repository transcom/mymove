package ghcapi

import (
	"errors"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
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
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
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

			moveOrders, err := models.FetchOrder(appCtx.DB(), move.OrdersID)
			if err != nil {
				appCtx.Logger().Error("Error retreiving user privileges", zap.Error(err))
			}

			if moveOrders.OrdersType == "SAFETY" && !privileges.HasPrivilege(models.PrivilegeTypeSafety) {
				appCtx.Logger().Error("Invalid permissions")
				errMsg := "Page is inaccessible"
				return moveop.NewGetMoveNotFound().WithPayload(&ghcmessages.Error{Message: &errMsg}), apperror.NewNotFoundError(uuid.Nil, "Page is inaccessible")
			} else {
				payload, err := payloads.Move(move, h.FileStorer())
				if err != nil {
					return nil, err
				}
				return moveop.NewGetMoveOK().WithPayload(payload), nil
			}
		})
}

type SearchMovesHandler struct {
	handlers.HandlerConfig
	services.MoveSearcher
	services.MoveUnlocker
}

func (h SearchMovesHandler) Handle(params moveop.SearchMovesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			searchMovesParams := services.SearchMovesParams{
				Branch:                params.Body.Branch,
				Locator:               params.Body.Locator,
				DodID:                 params.Body.DodID,
				Emplid:                params.Body.Emplid,
				CustomerName:          params.Body.CustomerName,
				PaymentRequestCode:    params.Body.PaymentRequestCode,
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

			// if the search move office user is accessing the queue, we need to unlock move/moves they have locked
			if appCtx.Session().IsOfficeUser() {
				officeUserID := appCtx.Session().OfficeUserID
				for i, move := range moves {
					lockedOfficeUserID := move.LockedByOfficeUserID
					if lockedOfficeUserID != nil && *lockedOfficeUserID == officeUserID {
						copyOfMove := move
						unlockedMove, err := h.UnlockMove(appCtx, &copyOfMove, officeUserID)
						if err != nil {
							return moveop.NewSearchMovesInternalServerError(), err
						}
						moves[i] = *unlockedMove
					}
				}
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

			payload, err := payloads.Move(move, h.FileStorer())
			if err != nil {
				return nil, err
			}
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

			payload, err := payloads.Move(move, h.FileStorer())
			if err != nil {
				return nil, err
			}
			return moveop.NewUpdateCloseoutOfficeOK().WithPayload(payload), nil
		})
}

type UploadAdditionalDocumentsHandler struct {
	handlers.HandlerConfig
	uploader services.MoveAdditionalDocumentsUploader
}

func (h UploadAdditionalDocumentsHandler) Handle(params moveop.UploadAdditionalDocumentsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			file, ok := params.File.(*runtime.File)
			if !ok {
				errMsg := "This should always be a runtime.File, something has changed in go-swagger."

				appCtx.Logger().Error(errMsg)

				return moveop.NewUploadAdditionalDocumentsInternalServerError(), nil
			}

			appCtx.Logger().Info(
				"File uploader and size",
				zap.String("userID", appCtx.Session().UserID.String()),
				zap.String("serviceMemberID", appCtx.Session().ServiceMemberID.String()),
				zap.String("officeUserID", appCtx.Session().OfficeUserID.String()),
				zap.String("AdminUserID", appCtx.Session().AdminUserID.String()),
				zap.Int64("size", file.Header.Size),
			)

			moveID, err := uuid.FromString(params.MoveID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			upload, url, verrs, err := h.uploader.CreateAdditionalDocumentsUpload(appCtx, appCtx.Session().UserID, moveID, file.Data, file.Header.Filename, h.FileStorer(), models.UploadTypeOFFICE)

			if verrs.HasAny() || err != nil {
				switch err.(type) {
				case uploader.ErrTooLarge:
					return moveop.NewUploadAdditionalDocumentsRequestEntityTooLarge(), err
				case uploader.ErrFile:
					return moveop.NewUploadAdditionalDocumentsInternalServerError(), err
				case uploader.ErrFailedToInitUploader:
					return moveop.NewUploadAdditionalDocumentsInternalServerError(), err
				case apperror.NotFoundError:
					return moveop.NewUploadAdditionalDocumentsNotFound(), err
				default:
					return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
				}
			}

			uploadPayload, err := payloadForUploadModelFromAdditionalDocumentsUpload(h.FileStorer(), upload, url)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return moveop.NewUploadAdditionalDocumentsCreated().WithPayload(uploadPayload), nil
		})
}

type MoveCancelerHandler struct {
	handlers.HandlerConfig
	services.MoveCanceler
}

func (h MoveCancelerHandler) Handle(params moveop.MoveCancelerParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID := uuid.FromStringOrNil(params.MoveID.String())

			move, err := h.MoveCanceler.CancelMove(appCtx, moveID)
			if err != nil {
				appCtx.Logger().Error("MoveCancelerHandler error", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return moveop.NewMoveCancelerNotFound(), err
				case apperror.PreconditionFailedError:
					return moveop.NewMoveCancelerPreconditionFailed(), err
				case apperror.InvalidInputError:
					return moveop.NewMoveCancelerUnprocessableEntity(), err
				case apperror.ConflictError:
					return moveop.NewMoveCancelerConflict(), err
				default:
					return moveop.NewMoveCancelerInternalServerError(), err
				}
			}

			payload, err := payloads.Move(move, h.FileStorer())
			if err != nil {
				return nil, err
			}
			return moveop.NewMoveCancelerOK().WithPayload(payload), nil
		})
}

type CheckForLockedMovesAndUnlockHandler struct {
	handlers.HandlerConfig
	services.MoveUnlocker
}

func (h CheckForLockedMovesAndUnlockHandler) Handle(params moveop.CheckForLockedMovesAndUnlockParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// if the search move office user is accessing the queue, we need to unlock move/moves they have locked
			err := h.CheckForLockedMovesAndUnlock(appCtx, uuid.FromStringOrNil(params.OfficeUserID.String()))
			if err != nil {
				return moveop.NewCheckForLockedMovesAndUnlockInternalServerError(), err
			}
			var payload moveop.CheckForLockedMovesAndUnlockOK
			payload.SetPayload(&moveop.CheckForLockedMovesAndUnlockOKBody{
				SuccessMessage: "Successfully unlocked all move(s) for current office user"})

			return &payload, nil
		})
}

type DeleteAssignedOfficeUserHandler struct {
	handlers.HandlerConfig
	services.MoveAssignedOfficeUserUpdater
}

func (h DeleteAssignedOfficeUserHandler) Handle(params moveop.DeleteAssignedOfficeUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID := uuid.FromStringOrNil(params.MoveID.String())

			role := getRole(*params.Body.RoleType)

			move, err := h.MoveAssignedOfficeUserUpdater.DeleteAssignedOfficeUser(appCtx, moveID, role)
			if err != nil {
				appCtx.Logger().Error("Error updating move", zap.Error(err))
				return moveop.NewDeleteAssignedOfficeUserInternalServerError(), err
			}

			payload, err := payloads.Move(move, h.FileStorer())
			if err != nil {
				return nil, err
			}

			return moveop.NewDeleteAssignedOfficeUserOK().WithPayload(payload), nil
		})
}

type UpdateAssignedOfficeUserHandler struct {
	handlers.HandlerConfig
	services.MoveAssignedOfficeUserUpdater
	officeUserFetcherPop services.OfficeUserFetcherPop
}

func (h UpdateAssignedOfficeUserHandler) Handle(params moveop.UpdateAssignedOfficeUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID := uuid.FromStringOrNil(params.MoveID.String())

			role := getRole(*params.Body.RoleType)

			officeUserID := uuid.FromStringOrNil(params.Body.OfficeUserID.String())

			officeUser, err := h.officeUserFetcherPop.FetchOfficeUserByID(appCtx, officeUserID)
			if err != nil {
				appCtx.Logger().Error("Error fetching office user", zap.Error(err))
				return moveop.NewUpdateAssignedOfficeUserNotFound(), err
			}

			move, err := h.MoveAssignedOfficeUserUpdater.UpdateAssignedOfficeUser(appCtx, moveID, &officeUser, role)
			if err != nil {
				appCtx.Logger().Error("Error updating move with an assigned office user", zap.Error(err))
				return moveop.NewUpdateAssignedOfficeUserInternalServerError(), err
			}

			payload, err := payloads.Move(move, h.FileStorer())
			if err != nil {
				return nil, err
			}

			return moveop.NewUpdateAssignedOfficeUserOK().WithPayload(payload), nil
		})
}

func payloadForUploadModelFromAdditionalDocumentsUpload(storer storage.FileStorer, upload models.Upload, url string) (*ghcmessages.Upload, error) {
	uploadPayload := &ghcmessages.Upload{
		ID:          handlers.FmtUUIDValue(upload.ID),
		Filename:    upload.Filename,
		ContentType: upload.ContentType,
		URL:         strfmt.URI(url),
		Bytes:       upload.Bytes,
		CreatedAt:   strfmt.DateTime(upload.CreatedAt),
		UpdatedAt:   strfmt.DateTime(upload.UpdatedAt),
	}
	tags, err := storer.Tags(upload.StorageKey)
	if err != nil || len(tags) == 0 {
		uploadPayload.Status = "PROCESSING"
	} else {
		uploadPayload.Status = tags["av-status"]
	}
	return uploadPayload, nil
}

func getRole(role string) roles.RoleType {
	var roleType roles.RoleType
	switch role {
	case "services_counselor":
		roleType = roles.RoleTypeServicesCounselor
	case "task_ordering_officer":
		roleType = roles.RoleTypeTOO
	case "task_invoicing_officer":
		roleType = roles.RoleTypeTIO
	}

	return roleType
}
