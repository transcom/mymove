package ghcapi

import (
	"log"
	"net/http"

	"github.com/transcom/mymove/pkg/services/query"

	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewGhcAPIHandler returns a handler for the GHC API
func NewGhcAPIHandler(context handlers.HandlerContext) http.Handler {

	queryBuilder := query.NewQueryBuilder(context.DB())

	ghcSpec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcAPI := ghcops.NewMymoveAPI(ghcSpec)
	ghcAPI.PaymentRequestsGetPaymentRequestHandler = GetPaymentRequestHandler{
		context,
		paymentrequest.NewPaymentRequestFetcher(queryBuilder),
	}

	ghcAPI.PaymentRequestsListPaymentRequestsHandler = ListPaymentRequestsHandler{
		context,
		paymentrequest.NewPaymentRequestListFetcher(context.DB()),
	}
	return ghcAPI.Serve(nil)
}
