package supportapi

import (
	"log"
	"net/http"

	"github.com/transcom/mymove/pkg/services/fetch"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"

	"github.com/go-openapi/loads"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/gen/supportapi"
	supportops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations"
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

	supportAPI.PaymentRequestsUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              context,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(queryBuilder),
	}
	supportAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandlerFunc{context, movetaskorder.NewMoveTaskOrderFetcher(context.DB())}

	supportAPI.MtoShipmentPatchMTOShipmentStatusHandler = PatchMTOShipmentStatusHandlerFunc{
		context,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentStatusUpdater(context.DB(), queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder), context.Planner()),
	}

	return supportAPI.Serve(nil)
}
