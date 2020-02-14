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
)

// NewGhcAPIHandler returns a handler for the GHC API
func NewGhcAPIHandler(context handlers.HandlerContext) *ghcops.MymoveAPI {
	ghcSpec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcAPI := ghcops.NewMymoveAPI(ghcSpec)
	queryBuilder := query.NewQueryBuilder(context.DB())

	ghcAPI.MtoServiceItemCreateMTOServiceItemHandler = CreateMTOServiceItemHandler{
		context,
		mtoserviceitem.NewMTOServiceItemCreator(queryBuilder),
	}

	ghcAPI.MtoServiceItemListMTOServiceItemsHandler = ListMTOServiceItemsHandler{
		context,
		fetch.NewListFetcher(queryBuilder),
		fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.PaymentRequestsGetPaymentRequestHandler = GetPaymentRequestHandler{
		context,
		paymentrequest.NewPaymentRequestFetcher(queryBuilder),
	}

	ghcAPI.PaymentRequestsUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              context,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(queryBuilder),
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
	ghcAPI.MoveOrderListMoveTaskOrdersHandler = ListMoveTaskOrdersHandler{context, movetaskorder.NewMoveTaskOrderFetcher(context.DB())}

	ghcAPI.MoveTaskOrderUpdateMoveTaskOrderStatusHandler = UpdateMoveTaskOrderStatusHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderUpdater(context.DB()),
	}

	ghcAPI.MtoShipmentListMTOShipmentsHandler = ListMTOShipmentsHandler{
		context,
		fetch.NewListFetcher(queryBuilder),
		fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.MtoShipmentPatchMTOShipmentStatusHandler = PatchShipmentHandler{
		context,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentStatusUpdater(context.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)),
	}

	return ghcAPI
}
