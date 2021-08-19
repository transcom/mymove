package ghcapi

import (
	"log"

	"github.com/transcom/mymove/pkg/services/fetch"
	order "github.com/transcom/mymove/pkg/services/order"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/services/office_user/customer"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

// NewGhcAPIHandler returns a handler for the GHC API
func NewGhcAPIHandler(ctx handlers.HandlerContext) *ghcops.MymoveAPI {
	ghcSpec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcAPI := ghcops.NewMymoveAPI(ghcSpec)
	queryBuilder := query.NewQueryBuilder(ctx.DB())
	moveRouter := move.NewMoveRouter(ctx.DB(), ctx.Logger())
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		ctx.DB(),
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
		moveRouter,
	)

	ghcAPI.ServeError = handlers.ServeCustomError

	ghcAPI.MoveGetMoveHandler = GetMoveHandler{
		HandlerContext: ctx,
		MoveFetcher:    move.NewMoveFetcher(ctx.DB()),
	}

	ghcAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{
		HandlerContext:        ctx,
		MTOServiceItemUpdater: mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter),
		Fetcher:               fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.MtoServiceItemListMTOServiceItemsHandler = ListMTOServiceItemsHandler{
		ctx,
		fetch.NewListFetcher(queryBuilder),
		fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.PaymentRequestsGetPaymentRequestHandler = GetPaymentRequestHandler{
		ctx,
		paymentrequest.NewPaymentRequestFetcher(ctx.DB()),
	}

	ghcAPI.PaymentRequestsGetPaymentRequestsForMoveHandler = GetPaymentRequestForMoveHandler{
		HandlerContext:            ctx,
		PaymentRequestListFetcher: paymentrequest.NewPaymentRequestListFetcher(ctx.DB()),
	}

	ghcAPI.PaymentRequestsUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              ctx,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(ctx.DB()),
	}

	ghcAPI.PaymentServiceItemUpdatePaymentServiceItemStatusHandler = UpdatePaymentServiceItemStatusHandler{
		HandlerContext: ctx,
		Fetcher:        fetch.NewFetcher(queryBuilder),
		Builder:        *queryBuilder,
	}

	ghcAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandler{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(ctx.DB()),
	}

	ghcAPI.CustomerGetCustomerHandler = GetCustomerHandler{
		ctx,
		customer.NewCustomerFetcher(ctx.DB()),
	}
	ghcAPI.CustomerUpdateCustomerHandler = UpdateCustomerHandler{
		ctx,
		customer.NewCustomerUpdater(ctx.DB()),
	}
	ghcAPI.OrderGetOrderHandler = GetOrdersHandler{
		ctx,
		order.NewOrderFetcher(ctx.DB()),
	}
	ghcAPI.OrderCounselingUpdateOrderHandler = CounselingUpdateOrderHandler{
		ctx,
		order.NewOrderUpdater(ctx.DB()),
	}

	ghcAPI.OrderUpdateOrderHandler = UpdateOrderHandler{
		ctx,
		order.NewOrderUpdater(ctx.DB()),
		moveTaskOrderUpdater,
	}

	ghcAPI.OrderUpdateAllowanceHandler = UpdateAllowanceHandler{
		ctx,
		order.NewOrderUpdater(ctx.DB()),
	}
	ghcAPI.OrderCounselingUpdateAllowanceHandler = CounselingUpdateAllowanceHandler{
		ctx,
		order.NewOrderUpdater(ctx.DB()),
	}

	ghcAPI.MoveTaskOrderUpdateMoveTaskOrderStatusHandler = UpdateMoveTaskOrderStatusHandlerFunc{
		ctx,
		moveTaskOrderUpdater,
	}

	ghcAPI.MoveTaskOrderUpdateMTOStatusServiceCounselingCompletedHandler = UpdateMTOStatusServiceCounselingCompletedHandlerFunc{
		ctx,
		moveTaskOrderUpdater,
	}

	ghcAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		ctx,
		mtoshipment.NewMTOShipmentCreator(
			ctx.DB(),
			queryBuilder,
			fetch.NewFetcher(queryBuilder),
			moveRouter,
		),
	}

	ghcAPI.MtoShipmentListMTOShipmentsHandler = ListMTOShipmentsHandler{
		ctx,
		fetch.NewListFetcher(queryBuilder),
		fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.ShipmentDeleteShipmentHandler = DeleteShipmentHandler{
		ctx,
		mtoshipment.NewShipmentDeleter(ctx.DB()),
	}

	ghcAPI.ShipmentApproveShipmentHandler = ApproveShipmentHandler{
		ctx,
		mtoshipment.NewShipmentApprover(
			ctx.DB(),
			mtoshipment.NewShipmentRouter(ctx.DB()),
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
			ctx.Planner(),
		),
	}

	ghcAPI.ShipmentRequestShipmentDiversionHandler = RequestShipmentDiversionHandler{
		ctx,
		mtoshipment.NewShipmentDiversionRequester(
			ctx.DB(),
			mtoshipment.NewShipmentRouter(ctx.DB()),
		),
	}

	ghcAPI.ShipmentApproveShipmentDiversionHandler = ApproveShipmentDiversionHandler{
		ctx,
		mtoshipment.NewShipmentDiversionApprover(
			ctx.DB(),
			mtoshipment.NewShipmentRouter(ctx.DB()),
		),
	}

	ghcAPI.ShipmentRejectShipmentHandler = RejectShipmentHandler{
		ctx,
		mtoshipment.NewShipmentRejecter(
			ctx.DB(),
			mtoshipment.NewShipmentRouter(ctx.DB()),
		),
	}

	ghcAPI.ShipmentRequestShipmentCancellationHandler = RequestShipmentCancellationHandler{
		ctx,
		mtoshipment.NewShipmentCancellationRequester(
			ctx.DB(),
			mtoshipment.NewShipmentRouter(ctx.DB()),
		),
	}

	ghcAPI.ShipmentRequestShipmentReweighHandler = RequestShipmentReweighHandler{
		ctx,
		mtoshipment.NewShipmentReweighRequester(
			ctx.DB(),
		),
	}

	ghcAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateShipmentHandler{
		ctx,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentUpdater(
			ctx.DB(),
			queryBuilder,
			fetch.NewFetcher(queryBuilder),
			ctx.Planner(),
			moveRouter,
		),
	}

	ghcAPI.MtoAgentFetchMTOAgentListHandler = ListMTOAgentsHandler{
		HandlerContext: ctx,
		ListFetcher:    fetch.NewListFetcher(queryBuilder),
	}

	ghcAPI.GhcDocumentsGetDocumentHandler = GetDocumentHandler{ctx}

	ghcAPI.QueuesGetMovesQueueHandler = GetMovesQueueHandler{
		ctx,
		order.NewOrderFetcher(ctx.DB()),
	}

	ghcAPI.QueuesGetPaymentRequestsQueueHandler = GetPaymentRequestsQueueHandler{
		ctx,
		paymentrequest.NewPaymentRequestListFetcher(ctx.DB()),
	}

	ghcAPI.QueuesGetServicesCounselingQueueHandler = GetServicesCounselingQueueHandler{
		ctx,
		order.NewOrderFetcher(ctx.DB()),
	}

	ghcAPI.TacTacValidationHandler = TacValidationHandler{
		ctx,
	}

	return ghcAPI
}
