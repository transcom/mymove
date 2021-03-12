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
func NewGhcAPIHandler(context handlers.HandlerContext) *ghcops.MymoveAPI {
	ghcSpec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcAPI := ghcops.NewMymoveAPI(ghcSpec)
	queryBuilder := query.NewQueryBuilder(context.DB())
	ghcAPI.ServeError = handlers.ServeCustomError

	ghcAPI.MoveGetMoveHandler = GetMoveHandler{
		HandlerContext: context,
		MoveFetcher:    move.NewMoveFetcher(context.DB()),
	}

	ghcAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{
		HandlerContext:        context,
		MTOServiceItemUpdater: mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder),
		Fetcher:               fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.MtoServiceItemListMTOServiceItemsHandler = ListMTOServiceItemsHandler{
		context,
		fetch.NewListFetcher(queryBuilder),
		fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.PaymentRequestsGetPaymentRequestHandler = GetPaymentRequestHandler{
		context,
		paymentrequest.NewPaymentRequestFetcher(context.DB()),
	}

	ghcAPI.PaymentRequestsGetPaymentRequestsForMoveHandler = GetPaymentRequestForMoveHandler{
		HandlerContext:            context,
		PaymentRequestListFetcher: paymentrequest.NewPaymentRequestListFetcher(context.DB()),
	}

	ghcAPI.PaymentRequestsUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              context,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(context.DB()),
	}

	ghcAPI.PaymentServiceItemUpdatePaymentServiceItemStatusHandler = UpdatePaymentServiceItemStatusHandler{
		HandlerContext: context,
		Fetcher:        fetch.NewFetcher(queryBuilder),
		Builder:        *queryBuilder,
	}

	ghcAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandler{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB()),
	}

	ghcAPI.CustomerGetCustomerHandler = GetCustomerHandler{
		context,
		customer.NewCustomerFetcher(context.DB()),
	}
	ghcAPI.OrderGetOrderHandler = GetOrdersHandler{
		context,
		order.NewOrderFetcher(context.DB()),
	}
	ghcAPI.OrderUpdateOrderHandler = UpdateOrderHandler{
		context,
		order.NewOrderUpdater(context.DB()),
	}
	ghcAPI.OrderListMoveTaskOrdersHandler = ListMoveTaskOrdersHandler{context, movetaskorder.NewMoveTaskOrderFetcher(context.DB())}

	ghcAPI.MoveTaskOrderUpdateMoveTaskOrderStatusHandler = UpdateMoveTaskOrderStatusHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderUpdater(context.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)),
	}

	ghcAPI.MtoShipmentListMTOShipmentsHandler = ListMTOShipmentsHandler{
		context,
		fetch.NewListFetcher(queryBuilder),
		fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.MtoShipmentPatchMTOShipmentStatusHandler = PatchShipmentHandler{
		context,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentStatusUpdater(context.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder), context.Planner()),
	}

	ghcAPI.MtoAgentFetchMTOAgentListHandler = ListMTOAgentsHandler{
		HandlerContext: context,
		ListFetcher:    fetch.NewListFetcher(queryBuilder),
	}

	ghcAPI.GhcDocumentsGetDocumentHandler = GetDocumentHandler{context}

	ghcAPI.QueuesGetMovesQueueHandler = GetMovesQueueHandler{
		context,
		order.NewOrderFetcher(context.DB()),
	}

	ghcAPI.QueuesGetPaymentRequestsQueueHandler = GetPaymentRequestsQueueHandler{
		context,
		paymentrequest.NewPaymentRequestListFetcher(context.DB()),
	}

	ghcAPI.TacTacValidationHandler = TacValidationHandler{
		context,
	}

	return ghcAPI
}
