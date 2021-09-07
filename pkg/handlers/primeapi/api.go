package primeapi

import (
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
	reweigh "github.com/transcom/mymove/pkg/services/reweigh"

	"github.com/transcom/mymove/pkg/gen/primeapi"
	primeops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"
)

// NewPrimeAPIHandler returns a handler for the Prime API
func NewPrimeAPIHandler(ctx handlers.HandlerContext) http.Handler {
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)

	primeSpec, err := loads.Analyzed(primeapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPI := primeops.NewMymoveAPI(primeSpec)
	queryBuilder := query.NewQueryBuilder()
	moveRouter := move.NewMoveRouter()
	moveWeights := move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

	primeAPI.ServeError = handlers.ServeCustomError

	primeAPI.MoveTaskOrderListMovesHandler = ListMovesHandler{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	primeAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandlerFunc{
		ctx,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	primeAPI.MtoServiceItemCreateMTOServiceItemHandler = CreateMTOServiceItemHandler{
		ctx,
		mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter),
		movetaskorder.NewMoveTaskOrderChecker(),
	}

	primeAPI.MtoServiceItemUpdateMTOServiceItemHandler = UpdateMTOServiceItemHandler{
		ctx,
		mtoserviceitem.NewMTOServiceItemUpdater(builder, moveRouter),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		ctx,
		mtoshipment.NewMTOShipmentUpdater(builder, fetcher, ctx.Planner(), moveRouter, moveWeights),
	}

	primeAPI.PaymentRequestCreatePaymentRequestHandler = CreatePaymentRequestHandler{
		ctx,
		paymentrequest.NewPaymentRequestCreator(
			ctx.GHCPlanner(),
			ghcrateengine.NewServiceItemPricer(),
		),
	}

	primeAPI.PaymentRequestCreateUploadHandler = CreateUploadHandler{
		ctx,
		paymentrequest.NewPaymentRequestUploadCreator(ctx.FileStorer()),
	}

	primeAPI.MoveTaskOrderUpdateMTOPostCounselingInformationHandler = UpdateMTOPostCounselingInformationHandler{
		ctx,
		fetch.NewFetcher(queryBuilder),
		movetaskorder.NewMoveTaskOrderUpdater(
			queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
			moveRouter,
		),
		movetaskorder.NewMoveTaskOrderChecker(),
	}

	primeAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		ctx,
		mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter),
		movetaskorder.NewMoveTaskOrderChecker(),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentAddressHandler = UpdateMTOShipmentAddressHandler{
		ctx,
		mtoshipment.NewMTOShipmentAddressUpdater(),
	}

	primeAPI.MtoShipmentCreateMTOAgentHandler = CreateMTOAgentHandler{
		ctx,
		mtoagent.NewMTOAgentCreator(movetaskorder.NewMoveTaskOrderChecker()),
	}

	primeAPI.MtoShipmentUpdateMTOAgentHandler = UpdateMTOAgentHandler{
		ctx,
		mtoagent.NewMTOAgentUpdater(movetaskorder.NewMoveTaskOrderChecker()),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentStatusHandler = UpdateMTOShipmentStatusHandler{
		ctx,
		mtoshipment.NewMTOShipmentUpdater(builder, fetcher, ctx.Planner(), moveRouter, moveWeights),
		mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter), ctx.Planner()),
	}

	primeAPI.MtoShipmentUpdateReweighHandler = UpdateReweighHandler{
		ctx,
		reweigh.NewReweighUpdater(),
	}

	return primeAPI.Serve(nil)
}
