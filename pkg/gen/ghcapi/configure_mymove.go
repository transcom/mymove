// This file is safe to edit. Once it exists it will not be overwritten

package ghcapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ghc_documents"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_agent"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/order"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
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
	if api.MtoAgentFetchMTOAgentListHandler == nil {
		api.MtoAgentFetchMTOAgentListHandler = mto_agent.FetchMTOAgentListHandlerFunc(func(params mto_agent.FetchMTOAgentListParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_agent.FetchMTOAgentList has not yet been implemented")
		})
	}
	if api.CustomerGetCustomerHandler == nil {
		api.CustomerGetCustomerHandler = customer.GetCustomerHandlerFunc(func(params customer.GetCustomerParams) middleware.Responder {
			return middleware.NotImplemented("operation customer.GetCustomer has not yet been implemented")
		})
	}
	if api.GhcDocumentsGetDocumentHandler == nil {
		api.GhcDocumentsGetDocumentHandler = ghc_documents.GetDocumentHandlerFunc(func(params ghc_documents.GetDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation ghc_documents.GetDocument has not yet been implemented")
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
	if api.MoveGetMoveHandler == nil {
		api.MoveGetMoveHandler = move.GetMoveHandlerFunc(func(params move.GetMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation move.GetMove has not yet been implemented")
		})
	}
	if api.MoveTaskOrderGetMoveTaskOrderHandler == nil {
		api.MoveTaskOrderGetMoveTaskOrderHandler = move_task_order.GetMoveTaskOrderHandlerFunc(func(params move_task_order.GetMoveTaskOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.GetMoveTaskOrder has not yet been implemented")
		})
	}
	if api.QueuesGetMovesQueueHandler == nil {
		api.QueuesGetMovesQueueHandler = queues.GetMovesQueueHandlerFunc(func(params queues.GetMovesQueueParams) middleware.Responder {
			return middleware.NotImplemented("operation queues.GetMovesQueue has not yet been implemented")
		})
	}
	if api.OrderGetOrderHandler == nil {
		api.OrderGetOrderHandler = order.GetOrderHandlerFunc(func(params order.GetOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation order.GetOrder has not yet been implemented")
		})
	}
	if api.PaymentRequestsGetPaymentRequestHandler == nil {
		api.PaymentRequestsGetPaymentRequestHandler = payment_requests.GetPaymentRequestHandlerFunc(func(params payment_requests.GetPaymentRequestParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.GetPaymentRequest has not yet been implemented")
		})
	}
	if api.PaymentRequestsGetPaymentRequestsForMoveHandler == nil {
		api.PaymentRequestsGetPaymentRequestsForMoveHandler = payment_requests.GetPaymentRequestsForMoveHandlerFunc(func(params payment_requests.GetPaymentRequestsForMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.GetPaymentRequestsForMove has not yet been implemented")
		})
	}
	if api.QueuesGetPaymentRequestsQueueHandler == nil {
		api.QueuesGetPaymentRequestsQueueHandler = queues.GetPaymentRequestsQueueHandlerFunc(func(params queues.GetPaymentRequestsQueueParams) middleware.Responder {
			return middleware.NotImplemented("operation queues.GetPaymentRequestsQueue has not yet been implemented")
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
	if api.OrderListMoveTaskOrdersHandler == nil {
		api.OrderListMoveTaskOrdersHandler = order.ListMoveTaskOrdersHandlerFunc(func(params order.ListMoveTaskOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation order.ListMoveTaskOrders has not yet been implemented")
		})
	}
	if api.MtoShipmentPatchMTOShipmentStatusHandler == nil {
		api.MtoShipmentPatchMTOShipmentStatusHandler = mto_shipment.PatchMTOShipmentStatusHandlerFunc(func(params mto_shipment.PatchMTOShipmentStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.PatchMTOShipmentStatus has not yet been implemented")
		})
	}
	if api.TacTacValidationHandler == nil {
		api.TacTacValidationHandler = tac.TacValidationHandlerFunc(func(params tac.TacValidationParams) middleware.Responder {
			return middleware.NotImplemented("operation tac.TacValidation has not yet been implemented")
		})
	}
	if api.MtoServiceItemUpdateMTOServiceItemHandler == nil {
		api.MtoServiceItemUpdateMTOServiceItemHandler = mto_service_item.UpdateMTOServiceItemHandlerFunc(func(params mto_service_item.UpdateMTOServiceItemParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.UpdateMTOServiceItem has not yet been implemented")
		})
	}
	if api.MtoServiceItemUpdateMTOServiceItemStatusHandler == nil {
		api.MtoServiceItemUpdateMTOServiceItemStatusHandler = mto_service_item.UpdateMTOServiceItemStatusHandlerFunc(func(params mto_service_item.UpdateMTOServiceItemStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.UpdateMTOServiceItemStatus has not yet been implemented")
		})
	}
	if api.OrderUpdateMoveOrderHandler == nil {
		api.OrderUpdateMoveOrderHandler = order.UpdateMoveOrderHandlerFunc(func(params order.UpdateMoveOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation order.UpdateMoveOrder has not yet been implemented")
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
	if api.PaymentRequestsUpdatePaymentRequestStatusHandler == nil {
		api.PaymentRequestsUpdatePaymentRequestStatusHandler = payment_requests.UpdatePaymentRequestStatusHandlerFunc(func(params payment_requests.UpdatePaymentRequestStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.UpdatePaymentRequestStatus has not yet been implemented")
		})
	}
	if api.PaymentServiceItemUpdatePaymentServiceItemStatusHandler == nil {
		api.PaymentServiceItemUpdatePaymentServiceItemStatusHandler = payment_service_item.UpdatePaymentServiceItemStatusHandlerFunc(func(params payment_service_item.UpdatePaymentServiceItemStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_service_item.UpdatePaymentServiceItemStatus has not yet been implemented")
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
