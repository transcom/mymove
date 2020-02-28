package primeapi

import (
	"time"

	"github.com/transcom/mymove/pkg/services"
	movetaskorderservice "github.com/transcom/mymove/pkg/services/move_task_order"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

// FetchMTOUpdatesHandler lists move task orders with the option to filter since a particular date
type FetchMTOUpdatesHandler struct {
	handlers.HandlerContext
}

// Handle fetches all move task orders with the option to filter since a particular date
func (h FetchMTOUpdatesHandler) Handle(params movetaskorderops.FetchMTOUpdatesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	var mtos models.MoveTaskOrders

	query := h.DB().Where("is_available_to_prime = ?", true).Eager(
		"PaymentRequests",
		"MTOServiceItems",
		"MTOServiceItems.ReService",
		"MTOShipments",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MoveOrder",
		"MoveOrder.Customer",
		"MoveOrder.Entitlement")
	if params.Since != nil {
		since := time.Unix(*params.Since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err := query.All(&mtos)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return movetaskorderops.NewFetchMTOUpdatesInternalServerError()
	}

	payload := payloads.MoveTaskOrders(&mtos)

	return movetaskorderops.NewFetchMTOUpdatesOK().WithPayload(payload)
}

// UpdateMTOPostCounselingInformationHandler updates the move task order with post-counseling information
type UpdateMTOPostCounselingInformationHandler struct {
	handlers.HandlerContext
	services.Fetcher
	services.MoveTaskOrderUpdater
}

// Handle updates to move task order post-counseling
func (h UpdateMTOPostCounselingInformationHandler) Handle(params movetaskorderops.UpdateMTOPostCounselingInformationParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	mtoID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	eTag := params.IfMatch
	mto, err := h.MoveTaskOrderUpdater.UpdatePostCounselingInfo(mtoID, params.Body, eTag)
	if err != nil {
		logger.Error("primeapi.UpdateMTOPostCounselingInformation error YOOOOOOOOOOO", zap.Error(err))
		switch err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationNotFound()
		case movetaskorderservice.PreconditionFailedError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationPreconditionFailed()
		case movetaskorderservice.ValidationError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationUnprocessableEntity()
		default:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationInternalServerError()
		}
	}
	mtoPayload := payloads.MoveTaskOrderWithEtag(mto)
	return movetaskorderops.NewUpdateMTOPostCounselingInformationOK().WithPayload(mtoPayload)
}
