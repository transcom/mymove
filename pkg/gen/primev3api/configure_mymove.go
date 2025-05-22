// This file is safe to edit. Once it exists it will not be overwritten

package primev3api

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/primev3api/primev3operations"
	"github.com/transcom/mymove/pkg/gen/primev3api/primev3operations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/primev3api/primev3operations/mto_shipment"
)

//go:generate swagger generate server --target ../../gen --name Mymove --spec ../../../swagger/prime_v3.yaml --api-package primev3operations --model-package primev3messages --server-package primev3api --principal interface{} --exclude-main

func configureFlags(api *primev3operations.MymoveAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *primev3operations.MymoveAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.MtoShipmentCreateMTOShipmentHandler == nil {
		api.MtoShipmentCreateMTOShipmentHandler = mto_shipment.CreateMTOShipmentHandlerFunc(func(params mto_shipment.CreateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.CreateMTOShipment has not yet been implemented")
		})
	}
	if api.MoveTaskOrderGetMoveTaskOrderHandler == nil {
		api.MoveTaskOrderGetMoveTaskOrderHandler = move_task_order.GetMoveTaskOrderHandlerFunc(func(params move_task_order.GetMoveTaskOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.GetMoveTaskOrder has not yet been implemented")
		})
	}
	if api.MtoShipmentUpdateMTOShipmentHandler == nil {
		api.MtoShipmentUpdateMTOShipmentHandler = mto_shipment.UpdateMTOShipmentHandlerFunc(func(params mto_shipment.UpdateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.UpdateMTOShipment has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

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
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
