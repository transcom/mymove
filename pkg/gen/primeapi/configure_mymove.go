// This file is safe to edit. Once it exists it will not be overwritten

package primeapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
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

	if api.MtoShipmentCreateMTOAgentHandler == nil {
		api.MtoShipmentCreateMTOAgentHandler = mto_shipment.CreateMTOAgentHandlerFunc(func(params mto_shipment.CreateMTOAgentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.CreateMTOAgent has not yet been implemented")
		})
	}
	if api.MtoServiceItemCreateMTOServiceItemHandler == nil {
		api.MtoServiceItemCreateMTOServiceItemHandler = mto_service_item.CreateMTOServiceItemHandlerFunc(func(params mto_service_item.CreateMTOServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.CreateMTOServiceItem has not yet been implemented")
		})
	}
	if api.MtoShipmentCreateMTOShipmentHandler == nil {
		api.MtoShipmentCreateMTOShipmentHandler = mto_shipment.CreateMTOShipmentHandlerFunc(func(params mto_shipment.CreateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.CreateMTOShipment has not yet been implemented")
		})
	}
	if api.PaymentRequestCreatePaymentRequestHandler == nil {
		api.PaymentRequestCreatePaymentRequestHandler = payment_request.CreatePaymentRequestHandlerFunc(func(params payment_request.CreatePaymentRequestParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_request.CreatePaymentRequest has not yet been implemented")
		})
	}
	if api.PaymentRequestCreateUploadHandler == nil {
		api.PaymentRequestCreateUploadHandler = payment_request.CreateUploadHandlerFunc(func(params payment_request.CreateUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_request.CreateUpload has not yet been implemented")
		})
	}
	if api.MoveTaskOrderFetchMTOUpdatesHandler == nil {
		api.MoveTaskOrderFetchMTOUpdatesHandler = move_task_order.FetchMTOUpdatesHandlerFunc(func(params move_task_order.FetchMTOUpdatesParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.FetchMTOUpdates has not yet been implemented")
		})
	}
	if api.MtoShipmentUpdateMTOAgentHandler == nil {
		api.MtoShipmentUpdateMTOAgentHandler = mto_shipment.UpdateMTOAgentHandlerFunc(func(params mto_shipment.UpdateMTOAgentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.UpdateMTOAgent has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMTOPostCounselingInformationHandler == nil {
		api.MoveTaskOrderUpdateMTOPostCounselingInformationHandler = move_task_order.UpdateMTOPostCounselingInformationHandlerFunc(func(params move_task_order.UpdateMTOPostCounselingInformationParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMTOPostCounselingInformation has not yet been implemented")
		})
	}
	if api.MtoServiceItemUpdateMTOServiceItemHandler == nil {
		api.MtoServiceItemUpdateMTOServiceItemHandler = mto_service_item.UpdateMTOServiceItemHandlerFunc(func(params mto_service_item.UpdateMTOServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.UpdateMTOServiceItem has not yet been implemented")
		})
	}
	if api.MtoShipmentUpdateMTOShipmentHandler == nil {
		api.MtoShipmentUpdateMTOShipmentHandler = mto_shipment.UpdateMTOShipmentHandlerFunc(func(params mto_shipment.UpdateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.UpdateMTOShipment has not yet been implemented")
		})
	}
	if api.MtoShipmentUpdateMTOShipmentAddressHandler == nil {
		api.MtoShipmentUpdateMTOShipmentAddressHandler = mto_shipment.UpdateMTOShipmentAddressHandlerFunc(func(params mto_shipment.UpdateMTOShipmentAddressParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.UpdateMTOShipmentAddress has not yet been implemented")
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
