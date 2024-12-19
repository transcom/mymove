package primeapiv3

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primev3api/primev3operations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapiv3/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// GetMoveTaskOrderHandler returns the details for a particular move
type GetMoveTaskOrderHandler struct {
	handlers.HandlerConfig
	moveTaskOrderFetcher   services.MoveTaskOrderFetcher
	shipmentRateAreaFinder services.ShipmentRateAreaFinder
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
			const featureFlagName = "boat"
			isBoatFeatureOn := false
			flag, err := h.FeatureFlagFetcher().GetBooleanFlag(params.HTTPRequest.Context(), appCtx.Logger(), "", featureFlagName, map[string]string{})
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
			flagMH, err := h.FeatureFlagFetcher().GetBooleanFlag(params.HTTPRequest.Context(), appCtx.Logger(), "", featureFlagNameMH, map[string]string{})
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

			// Add oconus rate area information to payload
			shipmentPostalCodeRateArea, err := h.shipmentRateAreaFinder.GetPrimeMoveShipmentOconusRateArea(appCtx, *mto)
			if err != nil {
				appCtx.Logger().Error("primeapi.GetMoveTaskOrderHandler error", zap.Error(err))
				return movetaskorderops.NewGetMoveTaskOrderInternalServerError().WithPayload(
					payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			appCtx.Logger().Info(fmt.Sprintf("getMoveTaskOrder shipmentPostalCodeRateArea.length=%v", len(*shipmentPostalCodeRateArea)))
			moveTaskOrderPayload := payloads.MoveTaskOrderWithShipmentOconusRateArea(mto, shipmentPostalCodeRateArea)

			return movetaskorderops.NewGetMoveTaskOrderOK().WithPayload(moveTaskOrderPayload), nil
		})
}
