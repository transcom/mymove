package primeapi

import (
	"time"

	"github.com/transcom/mymove/pkg/services"

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

	query := h.DB().Where("available_to_prime_at IS NOT NULL").Eager(
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.CustomerContacts",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"MoveOrder.Customer",
		"MoveOrder.Entitlement")
	if params.Since != nil {
		since := time.Unix(*params.Since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err := query.All(&mtos)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return movetaskorderops.NewFetchMTOUpdatesInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
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
	logger.Info("primeapi.UpdateMTOPostCounselingInformationHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

	mto, err := h.MoveTaskOrderUpdater.UpdatePostCounselingInfo(mtoID, params.Body, eTag)
	if err != nil {
		logger.Error("primeapi.UpdateMTOPostCounselingInformation error", zap.Error(err))
		switch e := err.(type) {
		case services.NotFoundError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationNotFound().WithPayload(
				payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID()))
		case services.PreconditionFailedError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationPreconditionFailed().WithPayload(
				payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceID()))
		case services.InvalidInputError:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationUnprocessableEntity().WithPayload(
				payloads.ValidationError(err.Error(), h.GetTraceID(), e.ValidationErrors))
		default:
			return movetaskorderops.NewUpdateMTOPostCounselingInformationInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}
	mtoPayload := payloads.MoveTaskOrder(mto)
	return movetaskorderops.NewUpdateMTOPostCounselingInformationOK().WithPayload(mtoPayload)
}
