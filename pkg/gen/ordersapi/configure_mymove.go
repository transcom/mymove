// This file is safe to edit. Once it exists it will not be overwritten

package ordersapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
)

//go:generate swagger generate server --target ../../gen --name Mymove --spec ../../../swagger/orders.yaml --api-package ordersoperations --model-package ordersmessages --server-package ordersapi --exclude-main

func configureFlags(api *ordersoperations.MymoveAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *ordersoperations.MymoveAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.GetOrdersHandler == nil {
		api.GetOrdersHandler = ordersoperations.GetOrdersHandlerFunc(func(params ordersoperations.GetOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation ordersoperations.GetOrders has not yet been implemented")
		})
	}
	if api.GetOrdersByIssuerAndOrdersNumHandler == nil {
		api.GetOrdersByIssuerAndOrdersNumHandler = ordersoperations.GetOrdersByIssuerAndOrdersNumHandlerFunc(func(params ordersoperations.GetOrdersByIssuerAndOrdersNumParams) middleware.Responder {
			return middleware.NotImplemented("operation ordersoperations.GetOrdersByIssuerAndOrdersNum has not yet been implemented")
		})
	}
	if api.IndexOrdersForMemberHandler == nil {
		api.IndexOrdersForMemberHandler = ordersoperations.IndexOrdersForMemberHandlerFunc(func(params ordersoperations.IndexOrdersForMemberParams) middleware.Responder {
			return middleware.NotImplemented("operation ordersoperations.IndexOrdersForMember has not yet been implemented")
		})
	}
	if api.PostRevisionHandler == nil {
		api.PostRevisionHandler = ordersoperations.PostRevisionHandlerFunc(func(params ordersoperations.PostRevisionParams) middleware.Responder {
			return middleware.NotImplemented("operation ordersoperations.PostRevision has not yet been implemented")
		})
	}
	if api.PostRevisionToOrdersHandler == nil {
		api.PostRevisionToOrdersHandler = ordersoperations.PostRevisionToOrdersHandlerFunc(func(params ordersoperations.PostRevisionToOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation ordersoperations.PostRevisionToOrders has not yet been implemented")
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
