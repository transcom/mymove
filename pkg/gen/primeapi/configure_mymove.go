// This file is safe to edit. Once it exists it will not be overwritten

package primeapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/uploads"
)

//go:generate swagger generate server --target ../../gen --name Mymove --spec ../../../swagger/prime.yaml --api-package primeoperations --model-package primemessages --server-package primeapi --exclude-main

func configureFlags(api *primeoperations.MymoveAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *primeoperations.MymoveAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	if api.PaymentRequestsCreatePaymentRequestHandler == nil {
		api.PaymentRequestsCreatePaymentRequestHandler = payment_requests.CreatePaymentRequestHandlerFunc(func(params payment_requests.CreatePaymentRequestParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.CreatePaymentRequest has not yet been implemented")
		})
	}
	if api.UploadsCreateUploadHandler == nil {
		api.UploadsCreateUploadHandler = uploads.CreateUploadHandlerFunc(func(params uploads.CreateUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.CreateUpload has not yet been implemented")
		})
	}
	if api.MoveTaskOrderFetchMTOUpdatesHandler == nil {
		api.MoveTaskOrderFetchMTOUpdatesHandler = move_task_order.FetchMTOUpdatesHandlerFunc(func(params move_task_order.FetchMTOUpdatesParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.FetchMTOUpdates has not yet been implemented")
		})
	}
	if api.MoveTaskOrderGetMoveTaskOrderCustomerHandler == nil {
		api.MoveTaskOrderGetMoveTaskOrderCustomerHandler = move_task_order.GetMoveTaskOrderCustomerHandlerFunc(func(params move_task_order.GetMoveTaskOrderCustomerParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.GetMoveTaskOrderCustomer has not yet been implemented")
		})
	}
	if api.MoveTaskOrderGetPrimeEntitlementsHandler == nil {
		api.MoveTaskOrderGetPrimeEntitlementsHandler = move_task_order.GetPrimeEntitlementsHandlerFunc(func(params move_task_order.GetPrimeEntitlementsParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.GetPrimeEntitlements has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMTOPostCounselingInformationHandler == nil {
		api.MoveTaskOrderUpdateMTOPostCounselingInformationHandler = move_task_order.UpdateMTOPostCounselingInformationHandlerFunc(func(params move_task_order.UpdateMTOPostCounselingInformationParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMTOPostCounselingInformation has not yet been implemented")
		})
	}
	if api.MtoShipmentUpdateMTOShipmentHandler == nil {
		api.MtoShipmentUpdateMTOShipmentHandler = mto_shipment.UpdateMTOShipmentHandlerFunc(func(params mto_shipment.UpdateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.UpdateMTOShipment has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMoveTaskOrderActualWeightHandler == nil {
		api.MoveTaskOrderUpdateMoveTaskOrderActualWeightHandler = move_task_order.UpdateMoveTaskOrderActualWeightHandlerFunc(func(params move_task_order.UpdateMoveTaskOrderActualWeightParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMoveTaskOrderActualWeight has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMoveTaskOrderDestinationAddressHandler == nil {
		api.MoveTaskOrderUpdateMoveTaskOrderDestinationAddressHandler = move_task_order.UpdateMoveTaskOrderDestinationAddressHandlerFunc(func(params move_task_order.UpdateMoveTaskOrderDestinationAddressParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMoveTaskOrderDestinationAddress has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMoveTaskOrderEstimatedWeightHandler == nil {
		api.MoveTaskOrderUpdateMoveTaskOrderEstimatedWeightHandler = move_task_order.UpdateMoveTaskOrderEstimatedWeightHandlerFunc(func(params move_task_order.UpdateMoveTaskOrderEstimatedWeightParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMoveTaskOrderEstimatedWeight has not yet been implemented")
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
