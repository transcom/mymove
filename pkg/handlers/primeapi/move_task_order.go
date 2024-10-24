package primeapi

import (
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
	"github.com/transcom/mymove/pkg/notifications"
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

			mtos, amendmentCountInfo, err := h.MoveTaskOrderFetcher.ListPrimeMoveTaskOrdersAmendments(appCtx, &searchParams)

			if err != nil {
				appCtx.Logger().Error("Unexpected error while fetching moves:", zap.Error(err))
				return movetaskorderops.NewListMovesInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			payload := payloads.ListMoves(&mtos, amendmentCountInfo)

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

			/** Feature Flag - Boat Shipment **/
			isBoatFeatureOn := false
			const featureFlagName = "boat"
			flag, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
			} else {
				isBoatFeatureOn = flag.Match
			}

			// Remove Boat shipments if Boat FF is off
			if !isBoatFeatureOn {
				var filteredShipments models.MTOShipments
				if mto.MTOShipments != nil {
					filteredShipments = models.MTOShipments{}
				}
				for i, shipment := range mto.MTOShipments {
					if shipment.ShipmentType == models.MTOShipmentTypeBoatHaulAway || shipment.ShipmentType == models.MTOShipmentTypeBoatTowAway {
						continue
					}

					filteredShipments = append(filteredShipments, mto.MTOShipments[i])
				}
				mto.MTOShipments = filteredShipments
			}
			/** End of Feature Flag **/

			/** Feature Flag - Mobile Home Shipment **/
			isMobileHomeFeatureOn := false
			const featureFlagNameMH = "mobile_home"
			flagMH, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagNameMH, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureFlagNameMH), zap.Error(err))
			} else {
				isMobileHomeFeatureOn = flagMH.Match
			}

			// Remove MobileHome shipments if MobileHome FF is off
			if !isMobileHomeFeatureOn {
				var filteredShipments models.MTOShipments
				if mto.MTOShipments != nil {
					filteredShipments = models.MTOShipments{}
				}
				for i, shipment := range mto.MTOShipments {
					if shipment.ShipmentType == models.MTOShipmentTypeMobileHome {
						continue
					}

					filteredShipments = append(filteredShipments, mto.MTOShipments[i])
				}
				mto.MTOShipments = filteredShipments
			}
			/** End of Feature Flag **/

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

			/* Don't send prime related emails on BLUEBARK moves */
			if mto.Orders.CanSendEmailWithOrdersType() {
				err = h.NotificationSender().SendNotification(appCtx,
					notifications.NewPrimeCounselingComplete(*mtoPayload),
				)
				if err != nil {
					appCtx.Logger().Error(err.Error())
					return movetaskorderops.NewUpdateMTOPostCounselingInformationInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			return movetaskorderops.NewUpdateMTOPostCounselingInformationOK().WithPayload(mtoPayload), nil
		})
}

// DownloadMoveOrderHandler is the struct to download all move orders by locator as a PDF
type DownloadMoveOrderHandler struct {
	handlers.HandlerConfig
	services.MoveSearcher
	services.OrderFetcher
	services.PrimeDownloadMoveUploadPDFGenerator
}

// Handler for downloading move order by locator as a PDF
func (h DownloadMoveOrderHandler) Handle(params movetaskorderops.DownloadMoveOrderParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			locator := strings.TrimSpace(params.Locator)
			docType := strings.TrimSpace(*params.Type)

			if len(locator) == 0 {
				err := apperror.NewBadDataError("primeapi.DownloadMoveOrder: missing/empty required URI parameter: locator")
				appCtx.Logger().Error(err.Error())
				return movetaskorderops.NewDownloadMoveOrderBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			searchMovesParams := services.SearchMovesParams{
				Locator: &locator,
			}
			moves, totalCount, err := h.MoveSearcher.SearchMoves(appCtx, &searchMovesParams)
			if err != nil {
				appCtx.Logger().Error("primeapi.DownloadMoveOrder error", zap.Error(err))
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
				// Note: OriginDutyLocation.ProvidesServicesCounseling == True means location has government based counseling.
				// FALSE indicates the location requires PRIME/GHC counseling.
				if move.Orders.OriginDutyLocation.ProvidesServicesCounseling {
					unprocessableErr := apperror.NewUnprocessableEntityError(
						fmt.Sprintf("primeapi.DownloadMoveOrder: Duty location of client's move currently does not have Prime counseling enabled, locator: %s", locator))
					appCtx.Logger().Warn(unprocessableErr.Error())
					payload := payloads.ValidationError(unprocessableErr.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), nil)
					return movetaskorderops.NewDownloadMoveOrderUnprocessableEntity().
						WithPayload(payload), unprocessableErr
				}
			}

			move := moves[len(moves)-1]

			var moveOrderUploadType = services.MoveOrderUploadAll
			if docType == "ORDERS" {
				moveOrderUploadType = services.MoveOrderUpload
			} else if docType == "AMENDMENTS" {
				moveOrderUploadType = services.MoveOrderAmendmentUpload
			}

			outputFile, err := h.PrimeDownloadMoveUploadPDFGenerator.GenerateDownloadMoveUserUploadPDF(appCtx, moveOrderUploadType, move, true)

			if err != nil {
				switch e := err.(type) {
				case apperror.UnprocessableEntityError:
					appCtx.Logger().Warn("primeapi.DownloadMoveOrder warn", zap.Error(err))
					payload := payloads.ValidationError(e.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), nil)
					return movetaskorderops.NewDownloadMoveOrderUnprocessableEntity().WithPayload(payload), err
				default:
					appCtx.Logger().Error("primeapi.DownloadMoveOrder error", zap.Error(err))
					return movetaskorderops.NewDownloadMoveOrderInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			payload := io.NopCloser(outputFile)

			// Build fileName in format: Customer-{type}-for-MTO-{locator}-{TIMESTAMP}.pdf
			// example:
			// Customer-ORDERS,AMENDMENTS-for-MTO-PPMSIT-2024-01-11T17-02.pdf   (all)
			// Customer-ORDERS-for-MTO-PPMSIT-2024-01-11T17-02.pdf
			// Customer-AMENDMENTS-for-MTO-PPMSIT-2024-01-11T17-02.pdf
			var fileNamePrefix = "Customer"
			if docType == "ALL" {
				fileNamePrefix += "-ORDERS,AMENDMENTS"
			} else {
				fileNamePrefix += "-" + docType
			}
			contentDisposition := fmt.Sprintf("inline; filename=\"%s-for-MTO-%s-%s.pdf\"", fileNamePrefix, locator, time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))

			return movetaskorderops.NewDownloadMoveOrderOK().WithContentDisposition(contentDisposition).WithPayload(payload), nil
		})
}
