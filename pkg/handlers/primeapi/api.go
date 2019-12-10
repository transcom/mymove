package primeapi

import (
	"log"
	"net/http"

	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/primeapi"
	primeops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewPrimeAPIHandler returns a handler for the Prime API
func NewPrimeAPIHandler(context handlers.HandlerContext) http.Handler {

	primeSpec, err := loads.Analyzed(primeapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPI := primeops.NewMymoveAPI(primeSpec)

	primeAPI.PaymentRequestsCreatePaymentRequestHandler = CreatePaymentRequestHandler{
		context,
		paymentrequest.NewPaymentRequestCreator(context.DB()),
	}

	return primeAPI.Serve(nil)
}
