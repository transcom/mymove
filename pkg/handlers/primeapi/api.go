package primeapi

import (
	"log"
	"net/http"

	mtoagent "github.com/transcom/mymove/pkg/services/mto_agent"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/services/ghcrateengine"
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
func NewPrimeAPIHandler(context handlers.HandlerContext) http.Handler {
	builder := query.NewQueryBuilder(context.DB())
	fetcher := fetch.NewFetcher(builder)

	primeSpec, err := loads.Analyzed(primeapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPI := primeops.NewMymoveAPI(primeSpec)
	queryBuilder := query.NewQueryBuilder(context.DB())

	primeAPI.ServeError = handlers.ServeCustomError

	primeAPI.MoveTaskOrderFetchMTOUpdatesHandler = FetchMTOUpdatesHandler{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB()),
	}

	primeAPI.MtoServiceItemCreateMTOServiceItemHandler = CreateMTOServiceItemHandler{
		context,
		mtoserviceitem.NewMTOServiceItemCreator(builder),
		movetaskorder.NewMoveTaskOrderChecker(context.DB()),
	}

	primeAPI.MtoServiceItemUpdateMTOServiceItemHandler = UpdateMTOServiceItemHandler{
		context,
		mtoserviceitem.NewMTOServiceItemUpdater(builder),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		context,
		mtoshipment.NewMTOShipmentUpdater(context.DB(), builder, fetcher, context.Planner()),
	}

	primeAPI.PaymentRequestCreatePaymentRequestHandler = CreatePaymentRequestHandler{
		context,
		paymentrequest.NewPaymentRequestCreator(
			context.DB(),
			context.GHCPlanner(),
			ghcrateengine.NewServiceItemPricer(context.DB()),
		),
	}

	primeAPI.PaymentRequestCreateUploadHandler = CreateUploadHandler{
		context,
		// To be fixed under this story: https://github.com/transcom/mymove/pull/3775/files#r397219200
		// unable to get logger to pass in for instantiation
		//paymentrequest.NewPaymentRequestUploadCreator(context.DB(), &logger,
		//	context.FileStorer()),
	}

	primeAPI.MoveTaskOrderUpdateMTOPostCounselingInformationHandler = UpdateMTOPostCounselingInformationHandler{
		context,
		fetch.NewFetcher(queryBuilder),
		movetaskorder.NewMoveTaskOrderUpdater(context.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)),
		movetaskorder.NewMoveTaskOrderChecker(context.DB()),
	}

	primeAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		context,
		mtoshipment.NewMTOShipmentCreator(context.DB(), builder, fetcher),
		movetaskorder.NewMoveTaskOrderChecker(context.DB()),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentAddressHandler = UpdateMTOShipmentAddressHandler{
		context,
		mtoshipment.NewMTOShipmentAddressUpdater(context.DB()),
	}

	primeAPI.MtoShipmentCreateMTOAgentHandler = CreateMTOAgentHandler{
		context,
		mtoagent.NewMTOAgentCreator(context.DB(), movetaskorder.NewMoveTaskOrderChecker(context.DB())),
	}

	primeAPI.MtoShipmentUpdateMTOAgentHandler = UpdateMTOAgentHandler{
		context,
		mtoagent.NewMTOAgentUpdater(context.DB()),
	}

	return primeAPI.Serve(nil)
}
