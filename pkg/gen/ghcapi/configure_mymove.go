// This file is safe to edit. Once it exists it will not be overwritten

package ghcapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/entitlements"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/service_item"
)

//go:generate swagger generate server --target ../../gen --name Mymove --spec ../../../swagger/ghc.yaml --api-package ghcoperations --model-package ghcmessages --server-package ghcapi --exclude-main

func configureFlags(api *ghcoperations.MymoveAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *ghcoperations.MymoveAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.PaymentRequestsCreatePaymentRequestHandler == nil {
		api.PaymentRequestsCreatePaymentRequestHandler = payment_requests.CreatePaymentRequestHandlerFunc(func(params payment_requests.CreatePaymentRequestParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.CreatePaymentRequest has not yet been implemented")
		})
	}
	if api.ServiceItemCreateServiceItemHandler == nil {
		api.ServiceItemCreateServiceItemHandler = service_item.CreateServiceItemHandlerFunc(func(params service_item.CreateServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation service_item.CreateServiceItem has not yet been implemented")
		})
	}
	if api.MoveTaskOrderDeleteMoveTaskOrderHandler == nil {
		api.MoveTaskOrderDeleteMoveTaskOrderHandler = move_task_order.DeleteMoveTaskOrderHandlerFunc(func(params move_task_order.DeleteMoveTaskOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.DeleteMoveTaskOrder has not yet been implemented")
		})
	}
	if api.ServiceItemDeleteServiceItemHandler == nil {
		api.ServiceItemDeleteServiceItemHandler = service_item.DeleteServiceItemHandlerFunc(func(params service_item.DeleteServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation service_item.DeleteServiceItem has not yet been implemented")
		})
	}
	if api.CustomerGetCustomerInfoHandler == nil {
		api.CustomerGetCustomerInfoHandler = customer.GetCustomerInfoHandlerFunc(func(params customer.GetCustomerInfoParams) middleware.Responder {
			return middleware.NotImplemented("operation customer.GetCustomerInfo has not yet been implemented")
		})
	}
	if api.EntitlementsGetEntitlementsHandler == nil {
		api.EntitlementsGetEntitlementsHandler = entitlements.GetEntitlementsHandlerFunc(func(params entitlements.GetEntitlementsParams) middleware.Responder {
			return middleware.NotImplemented("operation entitlements.GetEntitlements has not yet been implemented")
		})
	}
	if api.MoveTaskOrderGetMoveTaskOrderHandler == nil {
		api.MoveTaskOrderGetMoveTaskOrderHandler = move_task_order.GetMoveTaskOrderHandlerFunc(func(params move_task_order.GetMoveTaskOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.GetMoveTaskOrder has not yet been implemented")
		})
	}
	if api.PaymentRequestsGetPaymentRequestHandler == nil {
		api.PaymentRequestsGetPaymentRequestHandler = payment_requests.GetPaymentRequestHandlerFunc(func(params payment_requests.GetPaymentRequestParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.GetPaymentRequest has not yet been implemented")
		})
	}
	if api.ServiceItemGetServiceItemHandler == nil {
		api.ServiceItemGetServiceItemHandler = service_item.GetServiceItemHandlerFunc(func(params service_item.GetServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation service_item.GetServiceItem has not yet been implemented")
		})
	}
	if api.MoveTaskOrderListMoveTaskOrdersHandler == nil {
		api.MoveTaskOrderListMoveTaskOrdersHandler = move_task_order.ListMoveTaskOrdersHandlerFunc(func(params move_task_order.ListMoveTaskOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.ListMoveTaskOrders has not yet been implemented")
		})
	}
	if api.PaymentRequestsListPaymentRequestsHandler == nil {
		api.PaymentRequestsListPaymentRequestsHandler = payment_requests.ListPaymentRequestsHandlerFunc(func(params payment_requests.ListPaymentRequestsParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.ListPaymentRequests has not yet been implemented")
		})
	}
	if api.ServiceItemListServiceItemsHandler == nil {
		api.ServiceItemListServiceItemsHandler = service_item.ListServiceItemsHandlerFunc(func(params service_item.ListServiceItemsParams) middleware.Responder {
			return middleware.NotImplemented("operation service_item.ListServiceItems has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMoveTaskOrderHandler == nil {
		api.MoveTaskOrderUpdateMoveTaskOrderHandler = move_task_order.UpdateMoveTaskOrderHandlerFunc(func(params move_task_order.UpdateMoveTaskOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMoveTaskOrder has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMoveTaskOrderActualWeightHandler == nil {
		api.MoveTaskOrderUpdateMoveTaskOrderActualWeightHandler = move_task_order.UpdateMoveTaskOrderActualWeightHandlerFunc(func(params move_task_order.UpdateMoveTaskOrderActualWeightParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMoveTaskOrderActualWeight has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMoveTaskOrderStatusHandler == nil {
		api.MoveTaskOrderUpdateMoveTaskOrderStatusHandler = move_task_order.UpdateMoveTaskOrderStatusHandlerFunc(func(params move_task_order.UpdateMoveTaskOrderStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMoveTaskOrderStatus has not yet been implemented")
		})
	}
	if api.PaymentRequestsUpdatePaymentRequestHandler == nil {
		api.PaymentRequestsUpdatePaymentRequestHandler = payment_requests.UpdatePaymentRequestHandlerFunc(func(params payment_requests.UpdatePaymentRequestParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.UpdatePaymentRequest has not yet been implemented")
		})
	}
	if api.PaymentRequestsUpdatePaymentRequestStatusHandler == nil {
		api.PaymentRequestsUpdatePaymentRequestStatusHandler = payment_requests.UpdatePaymentRequestStatusHandlerFunc(func(params payment_requests.UpdatePaymentRequestStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.UpdatePaymentRequestStatus has not yet been implemented")
		})
	}
	if api.ServiceItemUpdateServiceItemHandler == nil {
		api.ServiceItemUpdateServiceItemHandler = service_item.UpdateServiceItemHandlerFunc(func(params service_item.UpdateServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation service_item.UpdateServiceItem has not yet been implemented")
		})
	}
	if api.ServiceItemUpdateServiceItemStatusHandler == nil {
		api.ServiceItemUpdateServiceItemStatusHandler = service_item.UpdateServiceItemStatusHandlerFunc(func(params service_item.UpdateServiceItemStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation service_item.UpdateServiceItemStatus has not yet been implemented")
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
