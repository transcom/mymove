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
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_order"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
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

	if api.MtoServiceItemCreateMTOServiceItemHandler == nil {
		api.MtoServiceItemCreateMTOServiceItemHandler = mto_service_item.CreateMTOServiceItemHandlerFunc(func(params mto_service_item.CreateMTOServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.CreateMTOServiceItem has not yet been implemented")
		})
	}
	if api.MtoServiceItemDeleteMTOServiceItemHandler == nil {
		api.MtoServiceItemDeleteMTOServiceItemHandler = mto_service_item.DeleteMTOServiceItemHandlerFunc(func(params mto_service_item.DeleteMTOServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.DeleteMTOServiceItem has not yet been implemented")
		})
	}
	if api.MoveTaskOrderDeleteMoveTaskOrderHandler == nil {
		api.MoveTaskOrderDeleteMoveTaskOrderHandler = move_task_order.DeleteMoveTaskOrderHandlerFunc(func(params move_task_order.DeleteMoveTaskOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.DeleteMoveTaskOrder has not yet been implemented")
		})
	}
	if api.CustomerGetAllCustomerMovesHandler == nil {
		api.CustomerGetAllCustomerMovesHandler = customer.GetAllCustomerMovesHandlerFunc(func(params customer.GetAllCustomerMovesParams) middleware.Responder {
			return middleware.NotImplemented("operation customer.GetAllCustomerMoves has not yet been implemented")
		})
	}
	if api.CustomerGetCustomerHandler == nil {
		api.CustomerGetCustomerHandler = customer.GetCustomerHandlerFunc(func(params customer.GetCustomerParams) middleware.Responder {
			return middleware.NotImplemented("operation customer.GetCustomer has not yet been implemented")
		})
	}
	if api.MoveTaskOrderGetEntitlementsHandler == nil {
		api.MoveTaskOrderGetEntitlementsHandler = move_task_order.GetEntitlementsHandlerFunc(func(params move_task_order.GetEntitlementsParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.GetEntitlements has not yet been implemented")
		})
	}
	if api.MtoServiceItemGetMTOServiceItemHandler == nil {
		api.MtoServiceItemGetMTOServiceItemHandler = mto_service_item.GetMTOServiceItemHandlerFunc(func(params mto_service_item.GetMTOServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.GetMTOServiceItem has not yet been implemented")
		})
	}
	if api.MoveOrderGetMoveOrderHandler == nil {
		api.MoveOrderGetMoveOrderHandler = move_order.GetMoveOrderHandlerFunc(func(params move_order.GetMoveOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_order.GetMoveOrder has not yet been implemented")
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
	if api.MtoServiceItemListMTOServiceItemsHandler == nil {
		api.MtoServiceItemListMTOServiceItemsHandler = mto_service_item.ListMTOServiceItemsHandlerFunc(func(params mto_service_item.ListMTOServiceItemsParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.ListMTOServiceItems has not yet been implemented")
		})
	}
	if api.MtoShipmentListMTOShipmentsHandler == nil {
		api.MtoShipmentListMTOShipmentsHandler = mto_shipment.ListMTOShipmentsHandlerFunc(func(params mto_shipment.ListMTOShipmentsParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.ListMTOShipments has not yet been implemented")
		})
	}
	if api.MoveOrderListMoveOrdersHandler == nil {
		api.MoveOrderListMoveOrdersHandler = move_order.ListMoveOrdersHandlerFunc(func(params move_order.ListMoveOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation move_order.ListMoveOrders has not yet been implemented")
		})
	}
	if api.MoveOrderListMoveTaskOrdersHandler == nil {
		api.MoveOrderListMoveTaskOrdersHandler = move_order.ListMoveTaskOrdersHandlerFunc(func(params move_order.ListMoveTaskOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation move_order.ListMoveTaskOrders has not yet been implemented")
		})
	}
	if api.PaymentRequestsListPaymentRequestsHandler == nil {
		api.PaymentRequestsListPaymentRequestsHandler = payment_requests.ListPaymentRequestsHandlerFunc(func(params payment_requests.ListPaymentRequestsParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.ListPaymentRequests has not yet been implemented")
		})
	}
	if api.MtoServiceItemUpdateMTOServiceItemHandler == nil {
		api.MtoServiceItemUpdateMTOServiceItemHandler = mto_service_item.UpdateMTOServiceItemHandlerFunc(func(params mto_service_item.UpdateMTOServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.UpdateMTOServiceItem has not yet been implemented")
		})
	}
	if api.MtoServiceItemUpdateMTOServiceItemstatusHandler == nil {
		api.MtoServiceItemUpdateMTOServiceItemstatusHandler = mto_service_item.UpdateMTOServiceItemstatusHandlerFunc(func(params mto_service_item.UpdateMTOServiceItemstatusParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.UpdateMTOServiceItemstatus has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMoveTaskOrderHandler == nil {
		api.MoveTaskOrderUpdateMoveTaskOrderHandler = move_task_order.UpdateMoveTaskOrderHandlerFunc(func(params move_task_order.UpdateMoveTaskOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMoveTaskOrder has not yet been implemented")
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
