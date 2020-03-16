package primeapi

import (
	"log"
	"net/http"

	"github.com/transcom/mymove/pkg/services/audit"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/go-openapi/loads"

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
	auditor := audit.NewAuditor(builder, context)

	primeSpec, err := loads.Analyzed(primeapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPI := primeops.NewMymoveAPI(primeSpec)
	queryBuilder := query.NewQueryBuilder(context.DB())

	primeAPI.MoveTaskOrderFetchMTOUpdatesHandler = FetchMTOUpdatesHandler{
		context,
		&auditor,
	}

	primeAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		context,
		mtoshipment.NewMTOShipmentUpdater(context.DB(), builder, fetcher),
	}

	primeAPI.PaymentRequestsCreatePaymentRequestHandler = CreatePaymentRequestHandler{
		context,
		paymentrequest.NewPaymentRequestCreator(context.DB()),
	}

	primeAPI.MoveTaskOrderUpdateMTOPostCounselingInformationHandler = UpdateMTOPostCounselingInformationHandler{
		context,
		fetch.NewFetcher(queryBuilder),
		movetaskorder.NewMoveTaskOrderUpdater(context.DB(), queryBuilder),
	}

	return primeAPI.Serve(nil)
}
