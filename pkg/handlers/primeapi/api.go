package primeapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

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

	primeAPI.MoveTaskOrderFetchMTOUpdatesHandler = FetchMTOUpdatesHandler{
		context,
	}

	primeAPI.MtoServiceItemCreateMTOServiceItemHandler = CreateMTOServiceItemHandler{
		context,
		mtoserviceitem.NewMTOServiceItemCreator(builder),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		context,
		mtoshipment.NewMTOShipmentUpdater(context.DB(), builder, fetcher),
	}

	primeAPI.PaymentRequestsCreatePaymentRequestHandler = CreatePaymentRequestHandler{
		context,
		paymentrequest.NewPaymentRequestCreator(context.DB()),
	}

	primeAPI.UploadsCreateUploadHandler = CreateUploadHandler{
		context,
		//paymentrequest.NewPaymentRequestUploadCreator(context.DB(), &logger,
		//	context.FileStorer()),
	}

	primeAPI.MoveTaskOrderUpdateMTOPostCounselingInformationHandler = UpdateMTOPostCounselingInformationHandler{
		context,
		fetch.NewFetcher(queryBuilder),
		movetaskorder.NewMoveTaskOrderUpdater(context.DB(), queryBuilder),
	}

	return primeAPI.Serve(nil)
}
