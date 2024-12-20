package supportapi

import (
	"log"
	"net/http"

	"github.com/benbjohnson/clock"
	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/supportapi"
	supportops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/invoice"
	lineofaccounting "github.com/transcom/mymove/pkg/services/line_of_accounting"
	move "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	portlocation "github.com/transcom/mymove/pkg/services/port_location"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
	internalmovetaskorder "github.com/transcom/mymove/pkg/services/support/move_task_order"
	transportationaccountingcode "github.com/transcom/mymove/pkg/services/transportation_accounting_code"
)

// NewSupportAPIHandler returns a handler for the Prime API
func NewSupportAPIHandler(handlerConfig handlers.HandlerConfig) http.Handler {
	queryBuilder := query.NewQueryBuilder()
	moveRouter := move.NewMoveRouter()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	addressCreator := address.NewAddressCreator()
	portLocationFetcher := portlocation.NewPortLocationFetcher()
	supportSpec, err := loads.Analyzed(supportapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	supportAPI := supportops.NewMymoveAPI(supportSpec)

	supportAPI.ServeError = handlers.ServeCustomError

	supportAPI.MoveTaskOrderListMTOsHandler = ListMTOsHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	supportAPI.MoveTaskOrderMakeMoveTaskOrderAvailableHandler = MakeMoveTaskOrderAvailableHandlerFunc{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderUpdater(
			queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveRouter, signedCertificationCreator, signedCertificationUpdater, ppmEstimator,
		),
	}

	supportAPI.MoveTaskOrderHideNonFakeMoveTaskOrdersHandler = HideNonFakeMoveTaskOrdersHandlerFunc{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderHider(),
	}

	supportAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandlerFunc{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher()}

	supportAPI.MoveTaskOrderCreateMoveTaskOrderHandler = CreateMoveTaskOrderHandler{
		handlerConfig,
		internalmovetaskorder.NewInternalMoveTaskOrderCreator(),
	}

	supportAPI.PaymentRequestUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerConfig:               handlerConfig,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(),
	}

	supportAPI.PaymentRequestListMTOPaymentRequestsHandler = ListMTOPaymentRequestsHandler{
		handlerConfig,
	}

	supportAPI.MtoShipmentUpdateMTOShipmentStatusHandler = UpdateMTOShipmentStatusHandlerFunc{
		handlerConfig,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()), handlerConfig.HHGPlanner()),
	}

	supportAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{handlerConfig, mtoserviceitem.NewMTOServiceItemUpdater(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer(), portLocationFetcher)}
	supportAPI.WebhookReceiveWebhookNotificationHandler = ReceiveWebhookNotificationHandler{handlerConfig}

	// Create TAC and LOA services
	tacFetcher := transportationaccountingcode.NewTransportationAccountingCodeFetcher()
	loaFetcher := lineofaccounting.NewLinesOfAccountingFetcher(tacFetcher)
	supportAPI.PaymentRequestGetPaymentRequestEDIHandler = GetPaymentRequestEDIHandler{
		HandlerConfig:                     handlerConfig,
		PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(),
		GHCPaymentRequestInvoiceGenerator: invoice.NewGHCPaymentRequestInvoiceGenerator(handlerConfig.ICNSequencer(), clock.New(), loaFetcher),
	}

	supportAPI.PaymentRequestProcessReviewedPaymentRequestsHandler = ProcessReviewedPaymentRequestsHandler{
		HandlerConfig:                 handlerConfig,
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
		HandlerConfig: handlerConfig,
		PaymentRequestRecalculator: paymentrequest.NewPaymentRequestRecalculator(
			paymentrequest.NewPaymentRequestCreator(
				handlerConfig.HHGPlanner(),
				ghcrateengine.NewServiceItemPricer(),
			),
			paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		),
	}

	supportAPI.WebhookCreateWebhookNotificationHandler = CreateWebhookNotificationHandler{
		HandlerConfig: handlerConfig,
	}

	return supportAPI.Serve(nil)
}
