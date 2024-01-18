package primeapi

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// ListMovesHandler lists moves with the option to filter since a particular date. Optimized ver.
type ListMovesHandler struct {
	handlers.HandlerConfig
	services.MoveTaskOrderFetcher
}

// Handle fetches all moves with the option to filter since a particular date. Optimized version.
func (h ListMovesHandler) Handle(params movetaskorderops.ListMovesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			var searchParams services.MoveTaskOrderFetcherParams
			if params.Since != nil {
				since := handlers.FmtDateTimePtrToPop(params.Since)
				searchParams.Since = &since
			}

			mtos, err := h.MoveTaskOrderFetcher.ListPrimeMoveTaskOrders(appCtx, &searchParams)

			if err != nil {
				appCtx.Logger().Error("Unexpected error while fetching moves:", zap.Error(err))
				return movetaskorderops.NewListMovesInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			payload := payloads.ListMoves(&mtos)

			return movetaskorderops.NewListMovesOK().WithPayload(payload), nil
		})
}

// GetMoveTaskOrderHandler returns the details for a particular move
type GetMoveTaskOrderHandler struct {
	handlers.HandlerConfig
	moveTaskOrderFetcher services.MoveTaskOrderFetcher
}

// Handle fetches a move from the database using its UUID or move code
func (h GetMoveTaskOrderHandler) Handle(params movetaskorderops.GetMoveTaskOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			searchParams := services.MoveTaskOrderFetcherParams{
				IsAvailableToPrime:       true,
				ExcludeExternalShipments: true,
			}

			// Add either ID or Locator to search params
			moveTaskOrderID := uuid.FromStringOrNil(params.MoveID)
			if moveTaskOrderID != uuid.Nil {
				searchParams.MoveTaskOrderID = moveTaskOrderID
			} else {
				searchParams.Locator = params.MoveID
			}

			mto, err := h.moveTaskOrderFetcher.FetchMoveTaskOrder(appCtx, &searchParams)
			if err != nil {
				appCtx.Logger().Error("primeapi.GetMoveTaskOrderHandler error", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewGetMoveTaskOrderNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, *handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return movetaskorderops.NewGetMoveTaskOrderInternalServerError().WithPayload(
						payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			moveTaskOrderPayload := payloads.MoveTaskOrder(mto)

			return movetaskorderops.NewGetMoveTaskOrderOK().WithPayload(moveTaskOrderPayload), nil
		})
}

// CreateExcessWeightRecordHandler uploads an excess weight record file
type CreateExcessWeightRecordHandler struct {
	handlers.HandlerConfig
	uploader services.MoveExcessWeightUploader
}

// Handle uploads the file passed into the request and updates the move
func (h CreateExcessWeightRecordHandler) Handle(params movetaskorderops.CreateExcessWeightRecordParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			moveID := uuid.FromStringOrNil(params.MoveTaskOrderID.String())

			file, ok := params.File.(*runtime.File)
			if !ok {
				err := apperror.NewInternalServerError("This should always be a runtime.File, something has changed in go-swagger.")
				appCtx.Logger().Error(err.Error())
				return movetaskorderops.NewCreateExcessWeightRecordInternalServerError().WithPayload(
					payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			excessWeightRecord, err := h.uploader.CreateExcessWeightUpload(
				appCtx, moveID, file.Data, file.Header.Filename, models.UploadTypePRIME)
			if err != nil {
				appCtx.Logger().Error("primeapi.CreateExcessWeightRecord error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewCreateExcessWeightRecordNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return movetaskorderops.NewCreateExcessWeightRecordUnprocessableEntity().WithPayload(
						payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.InvalidCreateInputError:
					return movetaskorderops.NewCreateExcessWeightRecordUnprocessableEntity().WithPayload(
						payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						appCtx.Logger().Error("primeapi.CreateExcessWeightRecord QueryError", zap.Error(e.Unwrap()))
					}
					return movetaskorderops.NewCreateExcessWeightRecordInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return movetaskorderops.NewCreateExcessWeightRecordInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			payload := payloads.ExcessWeightRecord(appCtx, h.FileStorer(), excessWeightRecord)
			return movetaskorderops.NewCreateExcessWeightRecordCreated().WithPayload(payload), nil
		})
}

// UpdateMTOPostCounselingInformationHandler updates the move with post-counseling information
type UpdateMTOPostCounselingInformationHandler struct {
	handlers.HandlerConfig
	services.Fetcher
	services.MoveTaskOrderUpdater
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// Handle updates to move post-counseling
func (h UpdateMTOPostCounselingInformationHandler) Handle(params movetaskorderops.UpdateMTOPostCounselingInformationParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			mtoID := uuid.FromStringOrNil(params.MoveTaskOrderID)
			eTag := params.IfMatch

			mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(appCtx, mtoID)

			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateMTOPostCounselingInformation error", zap.Error(err))
				return movetaskorderops.NewUpdateMTOPostCounselingInformationUnprocessableEntity().WithPayload(
					payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), nil)), err
			}

			if !mtoAvailableToPrime {
				err = apperror.NewInternalServerError("primeapi.UpdateMTOPostCounselingInformationHandler error - MTO is not available to Prime")
				appCtx.Logger().Error(err.Error())
				return movetaskorderops.NewUpdateMTOPostCounselingInformationNotFound().WithPayload(payloads.ClientError(
					handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", mtoID), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			mto, err := h.MoveTaskOrderUpdater.UpdatePostCounselingInfo(appCtx, mtoID, eTag)
			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateMTOPostCounselingInformation error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return movetaskorderops.NewUpdateMTOPostCounselingInformationNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.PreconditionFailedError:
					return movetaskorderops.NewUpdateMTOPostCounselingInformationPreconditionFailed().WithPayload(
						payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.ConflictError:
					return movetaskorderops.NewUpdateMTOPostCounselingInformationConflict().WithPayload(
						payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return movetaskorderops.NewUpdateMTOPostCounselingInformationUnprocessableEntity().WithPayload(
						payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				default:
					return movetaskorderops.NewUpdateMTOPostCounselingInformationInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			mtoPayload := payloads.MoveTaskOrder(mto)
			return movetaskorderops.NewUpdateMTOPostCounselingInformationOK().WithPayload(mtoPayload), nil
		})
}

// DownloadMoveOrderHandler is the struct to download all move orders by locator as a PDF
type DownloadMoveOrderHandler struct {
	handlers.HandlerConfig
	services.MoveSearcher
	services.OrderFetcher
}

// Handler for downloading move order by locator as a PDF
func (h DownloadMoveOrderHandler) Handle(params movetaskorderops.DownloadMoveOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			locator := strings.TrimSpace(params.Locator)

			if len(locator) == 0 {
				err := apperror.NewBadDataError("missing/empty required URI parameter: locator")
				appCtx.Logger().Error(err.Error())
				return movetaskorderops.NewDownloadMoveOrderBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			searchMovesParams := services.SearchMovesParams{
				Locator: &locator,
			}
			moves, totalCount, err := h.MoveSearcher.SearchMoves(appCtx, &searchMovesParams)
			if err != nil {
				appCtx.Logger().Error("Unexpected server error", zap.Error(err))
				return movetaskorderops.NewDownloadMoveOrderInternalServerError().WithPayload(
					payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			if totalCount == 0 {
				errMessage := fmt.Sprintf("Move not found, locator: %s.", locator)
				err := apperror.NewNotFoundError(uuid.Nil, errMessage)
				appCtx.Logger().Error(err.Error())
				return movetaskorderops.NewDownloadMoveOrderNotFound().WithPayload(
					payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			for _, move := range moves {
				var errMessage string
				// Check if move has requested counseling
				if move.Status != models.MoveStatusNeedsServiceCounseling {
					errMessage = fmt.Sprintf("Move is not in 'needs counseling state', locator: %s ", locator)
				}

				// Note: OriginDutyLocation.ProvidesServicesCounseling == True means location has government based counseling.
				// FALSE indicates the location requires PRIME/GHC counseling.
				if move.Orders.OriginDutyLocation.ProvidesServicesCounseling {
					errMessage = fmt.Sprintf("Duty location of client's move currently does not have Prime counseling enabled, locator: %s", locator)
				}

				if len(errMessage) > 0 {
					unprocessableErr := apperror.NewUnprocessableEntityError(errMessage)
					appCtx.Logger().Info(unprocessableErr.Error())
					payload := payloads.ValidationError(unprocessableErr.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), nil)
					return movetaskorderops.NewDownloadMoveOrderUnprocessableEntity().
						WithPayload(payload), unprocessableErr
				}
			}

			// For now return mock empty PDF file for 200 response.
			// TODO: (B-18027) - https://www13.v1host.com/USTRANSCOM38/story.mvc/Summary?oidToken=Story%3A870406
			// - Retrieve all uploaded move order docs
			// - Create new PDF service layer to merge all uploaded docs in one payload
			// - Wire up PDF service to generate response PDF payload
			buf := new(bytes.Buffer)
			payload := io.NopCloser(buf)
			filename := fmt.Sprintf("inline; filename=\"%s QA-%s %s.pdf\"", "MOCK", locator, time.Now().Format("01-02-2006"))

			return movetaskorderops.NewDownloadMoveOrderOK().WithContentDisposition(filename).WithPayload(payload), nil
		})
}
