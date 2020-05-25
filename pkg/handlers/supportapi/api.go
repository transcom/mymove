package supportapi

import (
	"log"
	"net/http"

	"github.com/transcom/mymove/pkg/services/support"

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

	supportAPI.MoveTaskOrderUpdateMoveTaskOrderStatusHandler = UpdateMoveTaskOrderStatusHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderUpdater(context.DB(), queryBuilder),
	}

	supportAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB())}

	supportAPI.MoveTaskOrderCreateMoveTaskOrderHandler = CreateMoveTaskOrderHandler{
		context,
		support.NewInternalMoveTaskOrderCreator(context.DB()),
	}

	supportAPI.PaymentRequestsUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              context,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(queryBuilder),
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
