package ghcapi

import (
	"log"
	"net/http"

	moveorder "github.com/transcom/mymove/pkg/services/move_order"

	"github.com/transcom/mymove/pkg/services/office_user/customer"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"

	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

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
	ghcAPI.PaymentRequestsGetPaymentRequestHandler = ShowPaymentRequestHandler{
		context,
		paymentrequest.NewPaymentRequestFetcher(context.DB()),
	}

	ghcAPI.PaymentRequestsListPaymentRequestsHandler = ListPaymentRequestsHandler{
		context,
		paymentrequest.NewPaymentRequestListFetcher(context.DB()),
	}
	ghcAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandler{
		context,
		movetaskorder.NewMoveTaskOrderFetcher(context.DB()),
	}
	ghcAPI.CustomerGetCustomerHandler = GetCustomerHandler{
		context,
		customer.NewCustomerFetcher(context.DB()),
	}
	ghcAPI.MoveOrderGetMoveOrderHandler = GetMoveOrdersHandler{
		context,
		moveorder.NewMoveOrderFetcher(context.DB()),
	}

	return ghcAPI.Serve(nil)
}
