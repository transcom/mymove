// This file is safe to edit. Once it exists it will not be overwritten

package adminapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/electronic_order"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/upload"
)

//go:generate swagger generate server --target ../../gen --name Mymove --spec ../../../swagger/admin.yaml --api-package adminoperations --model-package adminmessages --server-package adminapi --exclude-main

func configureFlags(api *adminoperations.MymoveAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *adminoperations.MymoveAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.OfficeCreateOfficeUserHandler == nil {
		api.OfficeCreateOfficeUserHandler = office.CreateOfficeUserHandlerFunc(func(params office.CreateOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office.CreateOfficeUser has not yet been implemented")
		})
	}
	if api.ElectronicOrderGetElectronicOrdersTotalsHandler == nil {
		api.ElectronicOrderGetElectronicOrdersTotalsHandler = electronic_order.GetElectronicOrdersTotalsHandlerFunc(func(params electronic_order.GetElectronicOrdersTotalsParams) middleware.Responder {
			return middleware.NotImplemented("operation electronic_order.GetElectronicOrdersTotals has not yet been implemented")
		})
	}
	if api.OfficeGetOfficeUserHandler == nil {
		api.OfficeGetOfficeUserHandler = office.GetOfficeUserHandlerFunc(func(params office.GetOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office.GetOfficeUser has not yet been implemented")
		})
	}
	if api.UploadGetUploadHandler == nil {
		api.UploadGetUploadHandler = upload.GetUploadHandlerFunc(func(params upload.GetUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation upload.GetUpload has not yet been implemented")
		})
	}
	if api.OfficeIndexAccessCodesHandler == nil {
		api.OfficeIndexAccessCodesHandler = office.IndexAccessCodesHandlerFunc(func(params office.IndexAccessCodesParams) middleware.Responder {
			return middleware.NotImplemented("operation office.IndexAccessCodes has not yet been implemented")
		})
	}
	if api.ElectronicOrderIndexElectronicOrdersHandler == nil {
		api.ElectronicOrderIndexElectronicOrdersHandler = electronic_order.IndexElectronicOrdersHandlerFunc(func(params electronic_order.IndexElectronicOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation electronic_order.IndexElectronicOrders has not yet been implemented")
		})
	}
	if api.OfficeIndexOfficeUsersHandler == nil {
		api.OfficeIndexOfficeUsersHandler = office.IndexOfficeUsersHandlerFunc(func(params office.IndexOfficeUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation office.IndexOfficeUsers has not yet been implemented")
		})
	}
	if api.OfficeIndexOfficesHandler == nil {
		api.OfficeIndexOfficesHandler = office.IndexOfficesHandlerFunc(func(params office.IndexOfficesParams) middleware.Responder {
			return middleware.NotImplemented("operation office.IndexOffices has not yet been implemented")
		})
	}
	if api.OfficeUpdateOfficeUserHandler == nil {
		api.OfficeUpdateOfficeUserHandler = office.UpdateOfficeUserHandlerFunc(func(params office.UpdateOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office.UpdateOfficeUser has not yet been implemented")
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
