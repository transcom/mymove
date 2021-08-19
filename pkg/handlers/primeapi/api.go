package primeapi

import (
	"context"
	"log"
	"net/http"

	mtoagent "github.com/transcom/mymove/pkg/services/mto_agent"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	move "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/transcom/mymove/pkg/gen/primeapi"
	primeops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"
)

// NewPrimeAPIHandler returns a handler for the Prime API
func NewPrimeAPIHandler(ctx handlers.HandlerContext) http.Handler {
	builder := query.NewQueryBuilder(ctx.DB())
	fetcher := fetch.NewFetcher(builder)

	primeSpec, err := loads.Analyzed(primeapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPI := primeops.NewMymoveAPI(primeSpec)
	queryBuilder := query.NewQueryBuilder(ctx.DB())
	moveRouter := move.NewMoveRouter(ctx.DB(), ctx.Logger())
	moveWeights := move.NewMoveWeights()

	primeAPI.ServeError = handlers.ServeCustomError

	primeAPI.MoveTaskOrderFetchMTOUpdatesHandler = FetchMTOUpdatesHandler{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(ctx.DB()),
	}

	primeAPI.MoveTaskOrderListMovesHandler = ListMovesHandler{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(ctx.DB()),
	}

	primeAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandlerFunc{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(ctx.DB()),
	}

	primeAPI.MtoServiceItemCreateMTOServiceItemHandler = CreateMTOServiceItemHandler{
		ctx,
		mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter),
		movetaskorder.NewMoveTaskOrderChecker(ctx.DB()),
	}

	primeAPI.MtoServiceItemUpdateMTOServiceItemHandler = UpdateMTOServiceItemHandler{
		ctx,
		mtoserviceitem.NewMTOServiceItemUpdater(builder, moveRouter),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		ctx,
		mtoshipment.NewMTOShipmentUpdater(ctx.DB(), builder, fetcher, ctx.Planner(), moveRouter, moveWeights),
	}

	primeAPI.PaymentRequestCreatePaymentRequestHandler = CreatePaymentRequestHandler{
		ctx,
		paymentrequest.NewPaymentRequestCreator(
			ctx.DB(),
			ctx.GHCPlanner(),
			ghcrateengine.NewServiceItemPricer(ctx.DB()),
		),
	}

	primeAPI.PaymentRequestCreateUploadHandler = CreateUploadHandler{
		ctx,
		paymentrequest.NewPaymentRequestUploadCreator(ctx.DB(), ctx.LoggerFromContext(context.Background()), ctx.FileStorer()),
	}

	primeAPI.MoveTaskOrderUpdateMTOPostCounselingInformationHandler = UpdateMTOPostCounselingInformationHandler{
		ctx,
		fetch.NewFetcher(queryBuilder),
		movetaskorder.NewMoveTaskOrderUpdater(
			ctx.DB(),
			queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
			moveRouter,
		),
		movetaskorder.NewMoveTaskOrderChecker(ctx.DB()),
	}

	primeAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		ctx,
		mtoshipment.NewMTOShipmentCreator(ctx.DB(), builder, fetcher, moveRouter),
		movetaskorder.NewMoveTaskOrderChecker(ctx.DB()),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentAddressHandler = UpdateMTOShipmentAddressHandler{
		ctx,
		mtoshipment.NewMTOShipmentAddressUpdater(ctx.DB()),
	}

	primeAPI.MtoShipmentCreateMTOAgentHandler = CreateMTOAgentHandler{
		ctx,
		mtoagent.NewMTOAgentCreator(ctx.DB(), movetaskorder.NewMoveTaskOrderChecker(ctx.DB())),
	}

	primeAPI.MtoShipmentUpdateMTOAgentHandler = UpdateMTOAgentHandler{
		ctx,
		mtoagent.NewMTOAgentUpdater(ctx.DB(), movetaskorder.NewMoveTaskOrderChecker(ctx.DB())),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentStatusHandler = UpdateMTOShipmentStatusHandler{
		ctx,
		mtoshipment.NewMTOShipmentUpdater(ctx.DB(), builder, fetcher, ctx.Planner(), moveRouter, moveWeights),
		mtoshipment.NewMTOShipmentStatusUpdater(ctx.DB(), queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter), ctx.Planner()),
	}

	return primeAPI.Serve(nil)
}
