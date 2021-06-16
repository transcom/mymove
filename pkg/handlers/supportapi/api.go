package supportapi

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/logging"

	"github.com/transcom/mymove/pkg/services/invoice"
	internalmovetaskorder "github.com/transcom/mymove/pkg/services/support/move_task_order"

	"github.com/benbjohnson/clock"
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
func NewSupportAPIHandler(ctx handlers.HandlerContext) http.Handler {
	queryBuilder := query.NewQueryBuilder(ctx.DB())
	supportSpec, err := loads.Analyzed(supportapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	gexURL := v.GetString(cli.GEXURLFlag)
	dbEnv := v.GetString(cli.DbEnvFlag)

	// Set the ICNSequencer in the handler: if we are in dev/test mode and sending to a real
	// GEX URL, then we should use a random ICN number within a defined range to avoid duplicate
	// test ICNs in Syncada.
	var icnSequencer sequence.Sequencer
	// ICNs are 9-digit numbers; reserve the ones in an upper range for development/testing.
	icnSequencer, err = sequence.NewRandomSequencer(ediinvoice.ICNRandomMin, ediinvoice.ICNRandomMax)
	if err != nil {
		log.Fatalln("Could not create random sequencer for ICN", zap.Error(err))
	}
	certLogger, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		log.Fatalln("Failed to initialize Zap logging", zap.Error(err))
	}
	certificates, rootCAs, err := certs.InitDoDEntrustCertificates(v, certLogger)
	if certificates == nil || rootCAs == nil || err != nil {
		log.Fatalln("Error in getting tls certs", zap.Error(err))
	}
	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}

	gexSender := invoice.NewGexSenderHTTP(
		gexURL,
		cli.GEXChannelInvoice,
		true,
		tlsConfig,
		v.GetString(cli.GEXBasicAuthUsernameFlag),
		v.GetString(cli.GEXBasicAuthPasswordFlag))
	reviewedPaymentRequestProcessor, err := paymentrequest.InitNewPaymentRequestReviewedProcessor(ctx.DB(), ctx.LoggerFromContext(context.Background()), true, icnSequencer, gexSender)
	if err != nil {
		msg := "failed to initialize InitNewPaymentRequestReviewedProcessor"
		log.Fatalln(msg, zap.Error(err))
	}

	supportAPI := supportops.NewMymoveAPI(supportSpec)

	supportAPI.ServeError = handlers.ServeCustomError

	supportAPI.MoveTaskOrderListMTOsHandler = ListMTOsHandler{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(ctx.DB()),
	}

	supportAPI.MoveTaskOrderMakeMoveTaskOrderAvailableHandler = MakeMoveTaskOrderAvailableHandlerFunc{
		ctx,
		movetaskorder.NewMoveTaskOrderUpdater(ctx.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)),
	}

	supportAPI.MoveTaskOrderHideNonFakeMoveTaskOrdersHandler = HideNonFakeMoveTaskOrdersHandlerFunc{
		ctx,
		movetaskorder.NewMoveTaskOrderHider(ctx.DB()),
	}

	supportAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandlerFunc{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(ctx.DB())}

	supportAPI.MoveTaskOrderCreateMoveTaskOrderHandler = CreateMoveTaskOrderHandler{
		ctx,
		internalmovetaskorder.NewInternalMoveTaskOrderCreator(ctx.DB()),
	}

	supportAPI.PaymentRequestUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerContext:              ctx,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(ctx.DB()),
	}

	supportAPI.PaymentRequestListMTOPaymentRequestsHandler = ListMTOPaymentRequestsHandler{
		ctx,
	}

	supportAPI.MtoShipmentUpdateMTOShipmentStatusHandler = UpdateMTOShipmentStatusHandlerFunc{
		ctx,
		fetch.NewFetcher(queryBuilder),
		mtoshipment.NewMTOShipmentStatusUpdater(ctx.DB(), queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder), ctx.Planner()),
	}

	supportAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{ctx, mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder)}
	supportAPI.WebhookReceiveWebhookNotificationHandler = ReceiveWebhookNotificationHandler{ctx}

	supportAPI.PaymentRequestGetPaymentRequestEDIHandler = GetPaymentRequestEDIHandler{
		HandlerContext:                    ctx,
		PaymentRequestFetcher:             paymentrequest.NewPaymentRequestFetcher(ctx.DB()),
		GHCPaymentRequestInvoiceGenerator: invoice.NewGHCPaymentRequestInvoiceGenerator(ctx.ICNSequencer(), clock.New()),
	}

	supportAPI.PaymentRequestProcessReviewedPaymentRequestsHandler = ProcessReviewedPaymentRequestsHandler{
		HandlerContext:                  ctx,
		PaymentRequestFetcher:           paymentrequest.NewPaymentRequestFetcher(ctx.DB()),
		PaymentRequestStatusUpdater:     paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestReviewedFetcher:   paymentrequest.NewPaymentRequestReviewedFetcher(ctx.DB()),
		PaymentRequestReviewedProcessor: reviewedPaymentRequestProcessor,
	}

	supportAPI.WebhookCreateWebhookNotificationHandler = CreateWebhookNotificationHandler{
		HandlerContext: ctx,
	}

	return supportAPI.Serve(nil)
}
