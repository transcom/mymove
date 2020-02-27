package primeapi

import (
	"log"
	"net/http"

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

	primeSpec, err := loads.Analyzed(primeapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPI := primeops.NewMymoveAPI(primeSpec)

	primeAPI.MoveTaskOrderFetchMTOUpdatesHandler = FetchMTOUpdatesHandler{
		context,
	}

	primeAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		context,
		mtoshipment.NewMTOShipmentUpdater(context.DB(), builder, fetcher),
	}

	primeAPI.PaymentRequestsCreatePaymentRequestHandler = CreatePaymentRequestHandler{
		context,
		paymentrequest.NewPaymentRequestCreator(context.DB()),
	}

	return primeAPI.Serve(nil)
}
