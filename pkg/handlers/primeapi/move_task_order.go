package primeapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
	movetaskorderservice "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/handlers"
)

//ListMoveTaskOrdersHandler handler for updating MoveTaskOrder Destination Address
type ListMoveTaskOrdersHandler struct {
	handlers.HandlerContext
}

//Handle handles requests to ListMoveTaskOrdersHandler
func (h ListMoveTaskOrdersHandler) Handle(params movetaskorderops.ListMoveTaskOrdersParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	var mtos models.MoveTaskOrders

	query := h.DB().Q()
	if params.Since != nil {
		since := time.Unix(*params.Since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err := query.All(&mtos)

	if err != nil {
		logger.Error("Unable to fetch records:", zap.Error(err))
		return movetaskorderops.NewListMoveTaskOrdersInternalServerError()
	}

	payload := make(primemessages.MoveTaskOrders, len(mtos))

	for i, m := range mtos {
		payload[i] = payloads.MoveTaskOrder(m)
	}

	return movetaskorderops.NewListMoveTaskOrdersOK().WithPayload(payload)
}

//UpdateMoveTaskOrderEstimatedWeightHandler handler for updating MoveTaskOrder Destination Address
type UpdateMoveTaskOrderEstimatedWeightHandler struct {
	handlers.HandlerContext
	moveTaskOrderPrimeEstimatedWeightUpdater services.MoveTaskOrderPrimeEstimatedWeightUpdater
}

//Handle handles requests to UpdateMoveTaskOrderEstimatedWeightHandler
func (h UpdateMoveTaskOrderEstimatedWeightHandler) Handle(params movetaskorderops.UpdateMoveTaskOrderEstimatedWeightParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	primeEstimatedWeight := unit.Pound(params.Body.PrimeEstimatedWeight)
	mto, err := h.moveTaskOrderPrimeEstimatedWeightUpdater.UpdatePrimeEstimatedWeight(moveTaskOrderID, primeEstimatedWeight, time.Now())
	if err != nil {
		logger.Error("ghciap.UpdateMoveTaskOrderEstimatedWeightHandler error", zap.Error(err))
		switch e := err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewUpdateMoveTaskOrderEstimatedWeightNotFound()
		case movetaskorderservice.ErrInvalidInput:
			payload := &primemessages.ValidationError{
				InvalidFields: e.InvalidFields(),
				ClientError: primemessages.ClientError{
					Title:    handlers.FmtString(handlers.ValidationErrMessage),
					Detail:   handlers.FmtString(e.Error()),
					Instance: handlers.FmtUUID(h.GetTraceID()),
				},
			}
			return movetaskorderops.NewUpdateMoveTaskOrderEstimatedWeightUnprocessableEntity().WithPayload(payload)
		default:
			return movetaskorderops.NewListMoveTaskOrdersInternalServerError()
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(*mto)
	return movetaskorderops.NewUpdateMoveTaskOrderEstimatedWeightOK().WithPayload(moveTaskOrderPayload)
}

//UpdateMoveTaskOrderPostCounselingInformationHandler handler for updating MoveTaskOrder Destination Address
type UpdateMoveTaskOrderPostCounselingInformationHandler struct {
	handlers.HandlerContext
	moveTaskOrderPostCounselingInformationUpdater services.MoveTaskOrderPrimePostCounselingUpdater
}

//Handle handles requests to UpdateMoveTaskOrderPostCounselingInformationHandler
func (h UpdateMoveTaskOrderPostCounselingInformationHandler) Handle(params movetaskorderops.UpdateMoveTaskOrderPostCounselingInformationParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	scheduledMoveDate := time.Time(params.Body.ScheduledMoveDate)
	mto, err := h.moveTaskOrderPostCounselingInformationUpdater.UpdateMoveTaskOrderPostCounselingInformation(moveTaskOrderID,
		services.PostCounselingInformation{
			PPMIsIncluded:            params.Body.PpmIsIncluded,
			ScheduledMoveDate:        scheduledMoveDate,
			SecondaryDeliveryAddress: payloads.AddressModel(params.Body.SecondaryDeliveryAddress),
			SecondaryPickupAddress:   payloads.AddressModel(params.Body.SecondaryPickupAddress),
		},
	)
	if err != nil {
		logger.Error("ghciap.UpdateMoveTaskOrderPostCounselingInformationHandler error", zap.Error(err))
		switch e := err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewUpdateMoveTaskOrderPostCounselingInformationNotFound()
		case movetaskorderservice.ErrInvalidInput:
			payload := &primemessages.ValidationError{
				InvalidFields: e.InvalidFields(),
				ClientError: primemessages.ClientError{
					Title:    handlers.FmtString(handlers.ValidationErrMessage),
					Detail:   handlers.FmtString(e.Error()),
					Instance: handlers.FmtUUID(h.GetTraceID()),
				},
			}
			return movetaskorderops.NewUpdateMoveTaskOrderPostCounselingInformationUnprocessableEntity().WithPayload(payload)
		default:
			return movetaskorderops.NewUpdateMoveTaskOrderPostCounselingInformationInternalServerError()
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(*mto)
	return movetaskorderops.NewUpdateMoveTaskOrderPostCounselingInformationOK().WithPayload(moveTaskOrderPayload)
}

//UpdateMoveTaskOrderDestinationAddressHandler handler for updating MoveTaskOrder Destination Address
type UpdateMoveTaskOrderDestinationAddressHandler struct {
	handlers.HandlerContext
	moveTaskOrderDestinationAddressUpdater services.MoveTaskOrderDestinationAddressUpdater
}

//Handle handles requests to UpdateMoveTaskOrderDestinationAddressHandler
func (h UpdateMoveTaskOrderDestinationAddressHandler) Handle(params movetaskorderops.UpdateMoveTaskOrderDestinationAddressParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	moveTaskOrderID := uuid.FromStringOrNil(params.MoveTaskOrderID)
	addressModel := payloads.AddressModel(params.DestinationAddress)
	mto, err := h.moveTaskOrderDestinationAddressUpdater.UpdateMoveTaskOrderDestinationAddress(moveTaskOrderID, addressModel)
	if err != nil {
		logger.Error("ghciap.UpdateMoveTaskOrderPostCounselingInformationHandler error", zap.Error(err))
		switch e := err.(type) {
		case movetaskorderservice.ErrNotFound:
			return movetaskorderops.NewUpdateMoveTaskOrderDestinationAddressNotFound()
		case movetaskorderservice.ErrInvalidInput:
			payload := &primemessages.ValidationError{
				InvalidFields: e.InvalidFields(),
				ClientError: primemessages.ClientError{
					Title:    handlers.FmtString(handlers.ValidationErrMessage),
					Detail:   handlers.FmtString(e.Error()),
					Instance: handlers.FmtUUID(h.GetTraceID()),
				},
			}
			return movetaskorderops.NewUpdateMoveTaskOrderDestinationAddressUnprocessableEntity().WithPayload(payload)
		default:
			return movetaskorderops.NewUpdateMoveTaskOrderDestinationAddressInternalServerError()
		}
	}
	moveTaskOrderPayload := payloads.MoveTaskOrder(*mto)
	return movetaskorderops.NewUpdateMoveTaskOrderDestinationAddressOK().WithPayload(moveTaskOrderPayload)
}
