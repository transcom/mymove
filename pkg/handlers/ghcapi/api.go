package ghcapi

import (
	"log"

	"github.com/transcom/mymove/pkg/services/fetch"
	moveorder "github.com/transcom/mymove/pkg/services/move_order"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/services/office_user/customer"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/handlers"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
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
		Fetcher:        fetch.NewFetcher(queryBuilder),
		NewQueryFilter: query.NewQueryFilter,
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

	ghcAPI.PaymentRequestsListPaymentRequestsHandler = ListPaymentRequestsHandler{
		context,
		paymentrequest.NewPaymentRequestListFetcher(context.DB()),
	}

	ghcAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandler{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB()),
	}
	ghcAPI.CustomerGetCustomerHandler = GetCustomerHandler{
		context,
		customer.NewCustomerFetcher(context.DB()),
	}
	ghcAPI.MoveOrderListMoveOrdersHandler = ListMoveOrdersHandler{context, moveorder.NewMoveOrderFetcher(context.DB())}
	ghcAPI.MoveOrderGetMoveOrderHandler = GetMoveOrdersHandler{
		context,
		moveorder.NewMoveOrderFetcher(context.DB()),
	}
	ghcAPI.MoveOrderUpdateMoveOrderHandler = UpdateMoveOrderHandler{
		context,
		moveorder.NewMoveOrderUpdater(context.DB(), queryBuilder),
	}
	ghcAPI.MoveOrderListMoveTaskOrdersHandler = ListMoveTaskOrdersHandler{context, movetaskorder.NewMoveTaskOrderFetcher(context.DB())}

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
		officeuser.NewOfficeUserFetcher(queryBuilder),
		moveorder.NewMoveOrderFetcher(context.DB()),
	}

	ghcAPI.QueuesGetPaymentRequestsQueueHandler = GetPaymentRequestsQueueHandler{
		context,
		officeuser.NewOfficeUserFetcher(queryBuilder),
		paymentrequest.NewPaymentRequestListFetcher(context.DB()),
	}

	return ghcAPI
}
