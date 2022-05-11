package supportapi

import (
	"log"
	"net/http"

	"github.com/benbjohnson/clock"
	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/supportapi"
	supportops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/invoice"
	move "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"
	internalmovetaskorder "github.com/transcom/mymove/pkg/services/support/move_task_order"
)

// NewSupportAPIHandler returns a handler for the Prime API
func NewSupportAPIHandler(ctx handlers.HandlerContext) http.Handler {
	queryBuilder := query.NewQueryBuilder()
	moveRouter := move.NewMoveRouter()
	supportSpec, err := loads.Analyzed(supportapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	supportAPI := supportops.NewMymoveAPI(supportSpec)

	supportAPI.ServeError = handlers.ServeCustomError

	supportAPI.MoveTaskOrderListMTOsHandler = ListMTOsHandler{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	supportAPI.MoveTaskOrderMakeMoveTaskOrderAvailableHandler = MakeMoveTaskOrderAvailableHandlerFunc{
		ctx,
		movetaskorder.NewMoveTaskOrderUpdater(
			queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
			moveRouter,
		),
	}

	supportAPI.MoveTaskOrderHideNonFakeMoveTaskOrdersHandler = HideNonFakeMoveTaskOrdersHandlerFunc{
		ctx,
		movetaskorder.NewMoveTaskOrderHider(),
	}

	supportAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandlerFunc{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher()}

	supportAPI.MoveTaskOrderCreateMoveTaskOrderHandler = CreateMoveTaskOrderHandler{
		ctx,
		internalmovetaskorder.NewInternalMoveTaskOrderCreator(),
	}

	supportAPI.PaymentRequestUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              ctx,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(),
	}

	supportAPI.PaymentRequestListMTOPaymentRequestsHandler = ListMTOPaymentRequestsHandler{
		ctx,
	}

	supportAPI.MtoShipmentUpdateMTOShipmentStatusHandler = UpdateMTOShipmentStatusHandlerFunc{
		ctx,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter), ctx.Planner()),
	}

	supportAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{ctx, mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter)}
	supportAPI.WebhookReceiveWebhookNotificationHandler = ReceiveWebhookNotificationHandler{ctx}

	supportAPI.PaymentRequestGetPaymentRequestEDIHandler = GetPaymentRequestEDIHandler{
		HandlerContext:                    ctx,
		PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(),
		GHCPaymentRequestInvoiceGenerator: invoice.NewGHCPaymentRequestInvoiceGenerator(ctx.ICNSequencer(), clock.New()),
	}

	supportAPI.PaymentRequestProcessReviewedPaymentRequestsHandler = ProcessReviewedPaymentRequestsHandler{
		HandlerContext:                ctx,
		PaymentRequestFetcher:         paymentrequest.NewPaymentRequestFetcher(),
		PaymentRequestStatusUpdater:   paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestReviewedFetcher: paymentrequest.NewPaymentRequestReviewedFetcher(),
		// Unable to get logger to pass in for the instantiation of
		// paymentrequest.InitNewPaymentRequestReviewedProcessor(appCtx.DB(), appCtx.Logger(), true),
		// This limitation has come up a few times
		// - https://dp3.atlassian.net/browse/MB-2352 (story to address issue)
		// - https://ustcdp3.slack.com/archives/CP6F568DC/p1592508325118600
		// - https://github.com/transcom/mymove/blob/c42adf61735be8ee8e5e83f41a656206f1e59b9d/pkg/handlers/primeapi/api.go
		// As a temporary workaround paymentrequest.InitNewPaymentRequestReviewedProcessor
		// is called directly in the handler
	}

	supportAPI.PaymentRequestRecalculatePaymentRequestHandler = RecalculatePaymentRequestHandler{
		HandlerContext: ctx,
		PaymentRequestRecalculator: paymentrequest.NewPaymentRequestRecalculator(
			paymentrequest.NewPaymentRequestCreator(
				ctx.GHCPlanner(),
				ghcrateengine.NewServiceItemPricer(),
			),
			paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		),
	}

	supportAPI.WebhookCreateWebhookNotificationHandler = CreateWebhookNotificationHandler{
		HandlerContext: ctx,
	}

	return supportAPI.Serve(nil)
}
