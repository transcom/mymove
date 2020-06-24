package supportapi

import (
	"log"
	"net/http"

	internalmovetaskorder "github.com/transcom/mymove/pkg/services/support/move_task_order"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"

	supportops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/transcom/mymove/pkg/gen/supportapi"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewSupportAPIHandler returns a handler for the Prime API
func NewSupportAPIHandler(context handlers.HandlerContext) http.Handler {
	queryBuilder := query.NewQueryBuilder(context.DB())

	supportSpec, err := loads.Analyzed(supportapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	supportAPI := supportops.NewMymoveAPI(supportSpec)

	supportAPI.ServeError = handlers.ServeCustomError

	supportAPI.MoveTaskOrderListMTOsHandler = ListMTOsHandler{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB()),
	}

	supportAPI.MoveTaskOrderMakeMoveTaskOrderAvailableHandler = MakeMoveTaskOrderAvailableHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderUpdater(context.DB(), queryBuilder),
	}

	supportAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB())}

	supportAPI.MoveTaskOrderCreateMoveTaskOrderHandler = CreateMoveTaskOrderHandler{
		context,
		internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
	}

	supportAPI.PaymentRequestsUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              context,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(queryBuilder),
	}

	supportAPI.PaymentRequestsListMTOPaymentRequestsHandler = ListMTOPaymentRequestsHandler{
		context,
	}

	supportAPI.MtoShipmentUpdateMTOShipmentStatusHandler = UpdateMTOShipmentStatusHandlerFunc{
		context,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentStatusUpdater(context.DB(), queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder), context.Planner()),
	}

	supportAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{context, mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder)}
	return supportAPI.Serve(nil)
}
