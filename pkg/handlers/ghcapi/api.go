package ghcapi

import (
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewGhcAPIHandler returns a handler for the GHC API
func NewGhcAPIHandler(context handlers.HandlerContext) http.Handler {

	ghcSpec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcAPI := ghcops.NewMymoveAPI(ghcSpec)

	ghcAPI.PaymentRequestsListPaymentRequestsHandler = ListPaymentRequestsHandler{
		context,
		paymentrequest.NewPaymentRequestListFetcher(context.DB()),
	}
	return ghcAPI.Serve(nil)
}
