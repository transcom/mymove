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
	queryBuilder := query.NewQueryBuilder()
	moveRouter := move.NewMoveRouter()
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
		moveRouter,
	)
	shipmentSITStatus := mtoshipment.NewShipmentSITStatus()

	ghcAPI.ServeError = handlers.ServeCustomError

	ghcAPI.MoveGetMoveHandler = GetMoveHandler{
		HandlerContext: ctx,
		MoveFetcher:    move.NewMoveFetcher(),
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
		paymentrequest.NewPaymentRequestFetcher(),
	}

	ghcAPI.PaymentRequestsGetPaymentRequestsForMoveHandler = GetPaymentRequestForMoveHandler{
		HandlerContext:            ctx,
		PaymentRequestListFetcher: paymentrequest.NewPaymentRequestListFetcher(),
	}

	ghcAPI.PaymentRequestsUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              ctx,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(),
	}

	ghcAPI.PaymentServiceItemUpdatePaymentServiceItemStatusHandler = UpdatePaymentServiceItemStatusHandler{
		HandlerContext: ctx,
		Fetcher:        fetch.NewFetcher(queryBuilder),
		Builder:        *queryBuilder,
	}

	ghcAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandler{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	ghcAPI.CustomerGetCustomerHandler = GetCustomerHandler{
		ctx,
		customer.NewCustomerFetcher(),
	}
	ghcAPI.CustomerUpdateCustomerHandler = UpdateCustomerHandler{
		ctx,
		customer.NewCustomerUpdater(),
	}
	ghcAPI.OrderGetOrderHandler = GetOrdersHandler{
		ctx,
		order.NewOrderFetcher(),
	}
	ghcAPI.OrderCounselingUpdateOrderHandler = CounselingUpdateOrderHandler{
		ctx,
		order.NewOrderUpdater(),
	}

	ghcAPI.OrderUpdateOrderHandler = UpdateOrderHandler{
		ctx,
		order.NewOrderUpdater(),
		moveTaskOrderUpdater,
	}

	ghcAPI.OrderUpdateAllowanceHandler = UpdateAllowanceHandler{
		ctx,
		order.NewOrderUpdater(),
	}
	ghcAPI.OrderCounselingUpdateAllowanceHandler = CounselingUpdateAllowanceHandler{
		ctx,
		order.NewOrderUpdater(),
	}
	ghcAPI.OrderUpdateBillableWeightHandler = UpdateBillableWeightHandler{
		ctx,
		order.NewExcessWeightRiskManager(moveRouter),
	}

	ghcAPI.OrderAcknowledgeExcessWeightRiskHandler = AcknowledgeExcessWeightRiskHandler{
		ctx,
		order.NewExcessWeightRiskManager(moveRouter),
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
			queryBuilder,
			fetch.NewFetcher(queryBuilder),
			moveRouter,
		),
	}

	ghcAPI.MtoShipmentListMTOShipmentsHandler = ListMTOShipmentsHandler{
		ctx,
		mtoshipment.NewMTOShipmentFetcher(),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentDeleteShipmentHandler = DeleteShipmentHandler{
		ctx,
		mtoshipment.NewShipmentDeleter(),
	}

	ghcAPI.ShipmentApproveShipmentHandler = ApproveShipmentHandler{
		ctx,
		mtoshipment.NewShipmentApprover(
			mtoshipment.NewShipmentRouter(),
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
			ctx.Planner(),
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentRequestShipmentDiversionHandler = RequestShipmentDiversionHandler{
		ctx,
		mtoshipment.NewShipmentDiversionRequester(
			mtoshipment.NewShipmentRouter(),
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentApproveShipmentDiversionHandler = ApproveShipmentDiversionHandler{
		ctx,
		mtoshipment.NewShipmentDiversionApprover(
			mtoshipment.NewShipmentRouter(),
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentRejectShipmentHandler = RejectShipmentHandler{
		ctx,
		mtoshipment.NewShipmentRejecter(
			mtoshipment.NewShipmentRouter(),
		),
	}

	ghcAPI.ShipmentRequestShipmentCancellationHandler = RequestShipmentCancellationHandler{
		ctx,
		mtoshipment.NewShipmentCancellationRequester(
			mtoshipment.NewShipmentRouter(),
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentRequestShipmentReweighHandler = RequestShipmentReweighHandler{
		ctx,
		mtoshipment.NewShipmentReweighRequester(),
		shipmentSITStatus,
	}

	ghcAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateShipmentHandler{
		ctx,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentUpdater(
			queryBuilder,
			fetch.NewFetcher(queryBuilder),
			ctx.Planner(),
			moveRouter,
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
		),
		shipmentSITStatus,
	}

	ghcAPI.MtoAgentFetchMTOAgentListHandler = ListMTOAgentsHandler{
		HandlerContext: ctx,
		ListFetcher:    fetch.NewListFetcher(queryBuilder),
	}

	ghcAPI.ShipmentApproveSitExtensionHandler = ApproveSITExtensionHandler{
		ctx,
		mtoshipment.NewSITExtensionApprover(),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentDenySitExtensionHandler = DenySITExtensionHandler{
		ctx,
		mtoshipment.NewSITExtensionDenier(),
		shipmentSITStatus,
	}

	ghcAPI.GhcDocumentsGetDocumentHandler = GetDocumentHandler{ctx}

	ghcAPI.QueuesGetMovesQueueHandler = GetMovesQueueHandler{
		ctx,
		order.NewOrderFetcher(),
	}

	ghcAPI.QueuesGetPaymentRequestsQueueHandler = GetPaymentRequestsQueueHandler{
		ctx,
		paymentrequest.NewPaymentRequestListFetcher(),
	}

	ghcAPI.QueuesGetServicesCounselingQueueHandler = GetServicesCounselingQueueHandler{
		ctx,
		order.NewOrderFetcher(),
	}

	ghcAPI.TacTacValidationHandler = TacValidationHandler{
		ctx,
	}

	return ghcAPI
}
