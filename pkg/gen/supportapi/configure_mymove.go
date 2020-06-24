// This file is safe to edit. Once it exists it will not be overwritten

package supportapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/supportapi/supportoperations"
	"github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/payment_requests"
)

//go:generate swagger generate server --target ../../gen --name Mymove --spec ../../../swagger/support.yaml --api-package supportoperations --model-package supportmessages --server-package supportapi --exclude-main

func configureFlags(api *supportoperations.MymoveAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *supportoperations.MymoveAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.MoveTaskOrderCreateMoveTaskOrderHandler == nil {
		api.MoveTaskOrderCreateMoveTaskOrderHandler = move_task_order.CreateMoveTaskOrderHandlerFunc(func(params move_task_order.CreateMoveTaskOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.CreateMoveTaskOrder has not yet been implemented")
		})
	}
	if api.MoveTaskOrderGetMoveTaskOrderHandler == nil {
		api.MoveTaskOrderGetMoveTaskOrderHandler = move_task_order.GetMoveTaskOrderHandlerFunc(func(params move_task_order.GetMoveTaskOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.GetMoveTaskOrder has not yet been implemented")
		})
	}
	if api.PaymentRequestsListMTOPaymentRequestsHandler == nil {
		api.PaymentRequestsListMTOPaymentRequestsHandler = payment_requests.ListMTOPaymentRequestsHandlerFunc(func(params payment_requests.ListMTOPaymentRequestsParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.ListMTOPaymentRequests has not yet been implemented")
		})
	}
	if api.MoveTaskOrderListMTOsHandler == nil {
		api.MoveTaskOrderListMTOsHandler = move_task_order.ListMTOsHandlerFunc(func(params move_task_order.ListMTOsParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.ListMTOs has not yet been implemented")
		})
	}
	if api.MoveTaskOrderMakeMoveTaskOrderAvailableHandler == nil {
		api.MoveTaskOrderMakeMoveTaskOrderAvailableHandler = move_task_order.MakeMoveTaskOrderAvailableHandlerFunc(func(params move_task_order.MakeMoveTaskOrderAvailableParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.MakeMoveTaskOrderAvailable has not yet been implemented")
		})
	}
	if api.MtoServiceItemUpdateMTOServiceItemStatusHandler == nil {
		api.MtoServiceItemUpdateMTOServiceItemStatusHandler = mto_service_item.UpdateMTOServiceItemStatusHandlerFunc(func(params mto_service_item.UpdateMTOServiceItemStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.UpdateMTOServiceItemStatus has not yet been implemented")
		})
	}
	if api.MtoShipmentUpdateMTOShipmentStatusHandler == nil {
		api.MtoShipmentUpdateMTOShipmentStatusHandler = mto_shipment.UpdateMTOShipmentStatusHandlerFunc(func(params mto_shipment.UpdateMTOShipmentStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.UpdateMTOShipmentStatus has not yet been implemented")
		})
	}
	if api.PaymentRequestsUpdatePaymentRequestStatusHandler == nil {
		api.PaymentRequestsUpdatePaymentRequestStatusHandler = payment_requests.UpdatePaymentRequestStatusHandlerFunc(func(params payment_requests.UpdatePaymentRequestStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.UpdatePaymentRequestStatus has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
