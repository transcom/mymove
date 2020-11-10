package supportapi

import (
	"log"
	"net/http"

	"github.com/transcom/mymove/pkg/services/invoice"
	internalmovetaskorder "github.com/transcom/mymove/pkg/services/support/move_task_order"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/gen/supportapi"
	supportops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations"
	"github.com/transcom/mymove/pkg/handlers"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
)

// NewSupportAPIHandler returns a handler for the Prime API
func NewSupportAPIHandler(context handlers.HandlerContext) http.Handler {
	queryBuilder := query.NewQueryBuilder(context.DB())
	var logger paymentrequest.Logger
	// logger = log.Printf

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
		movetaskorder.NewMoveTaskOrderUpdater(context.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)),
	}

	supportAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandlerFunc{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB())}

	supportAPI.MoveTaskOrderCreateMoveTaskOrderHandler = CreateMoveTaskOrderHandler{
		context,
		internalmovetaskorder.NewInternalMoveTaskOrderCreator(context.DB()),
	}

	supportAPI.PaymentRequestUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              context,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(context.DB()),
	}

	supportAPI.PaymentRequestListMTOPaymentRequestsHandler = ListMTOPaymentRequestsHandler{
		context,
	}

	supportAPI.MtoShipmentUpdateMTOShipmentStatusHandler = UpdateMTOShipmentStatusHandlerFunc{
		context,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentStatusUpdater(context.DB(), queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder), context.Planner()),
	}

	supportAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{context, mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder)}
	supportAPI.WebhookPostWebhookNotifyHandler = PostWebhookNotifyHandler{context}

	supportAPI.PaymentRequestGetPaymentRequestEDIHandler = GetPaymentRequestEDIHandler{
		HandlerContext:                    context,
		PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(context.DB()),
		GHCPaymentRequestInvoiceGenerator: invoice.NewGHCPaymentRequestInvoiceGenerator(context.DB()),
	}

	supportAPI.PaymentRequestProcessReviewedPaymentRequestsHandler = ProcessReviewedPaymentRequestsHandler{
		HandlerContext:                  context,
		PaymentRequestFetcher:           paymentrequest.NewPaymentRequestFetcher(context.DB()),
		PaymentRequestStatusUpdater:     paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestReviewedFetcher:   paymentrequest.NewPaymentRequestReviewedFetcher(context.DB()),
		PaymentRequestReviewedProcessor: paymentrequest.InitNewPaymentRequestReviewedProcessor(context.DB(), logger, true),
	}

	return supportAPI.Serve(nil)
}
