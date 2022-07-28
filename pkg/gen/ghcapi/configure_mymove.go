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
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer_support_remarks"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/evaluation_reports"
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
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/shipment"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
)

//go:generate swagger generate server --target ../../gen --name Mymove --spec ../../../swagger/ghc.yaml --api-package ghcoperations --model-package ghcmessages --server-package ghcapi --principal interface{} --exclude-main

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

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.MoveTaskOrderUpdateMTOReviewedBillableWeightsAtHandler == nil {
		api.MoveTaskOrderUpdateMTOReviewedBillableWeightsAtHandler = move_task_order.UpdateMTOReviewedBillableWeightsAtHandlerFunc(func(params move_task_order.UpdateMTOReviewedBillableWeightsAtParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMTOReviewedBillableWeightsAt has not yet been implemented")
		})
	}
	if api.OrderAcknowledgeExcessWeightRiskHandler == nil {
		api.OrderAcknowledgeExcessWeightRiskHandler = order.AcknowledgeExcessWeightRiskHandlerFunc(func(params order.AcknowledgeExcessWeightRiskParams) middleware.Responder {
			return middleware.NotImplemented("operation order.AcknowledgeExcessWeightRisk has not yet been implemented")
		})
	}
	if api.ShipmentApproveSITExtensionHandler == nil {
		api.ShipmentApproveSITExtensionHandler = shipment.ApproveSITExtensionHandlerFunc(func(params shipment.ApproveSITExtensionParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.ApproveSITExtension has not yet been implemented")
		})
	}
	if api.ShipmentApproveShipmentHandler == nil {
		api.ShipmentApproveShipmentHandler = shipment.ApproveShipmentHandlerFunc(func(params shipment.ApproveShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.ApproveShipment has not yet been implemented")
		})
	}
	if api.ShipmentApproveShipmentDiversionHandler == nil {
		api.ShipmentApproveShipmentDiversionHandler = shipment.ApproveShipmentDiversionHandlerFunc(func(params shipment.ApproveShipmentDiversionParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.ApproveShipmentDiversion has not yet been implemented")
		})
	}
	if api.OrderCounselingUpdateAllowanceHandler == nil {
		api.OrderCounselingUpdateAllowanceHandler = order.CounselingUpdateAllowanceHandlerFunc(func(params order.CounselingUpdateAllowanceParams) middleware.Responder {
			return middleware.NotImplemented("operation order.CounselingUpdateAllowance has not yet been implemented")
		})
	}
	if api.OrderCounselingUpdateOrderHandler == nil {
		api.OrderCounselingUpdateOrderHandler = order.CounselingUpdateOrderHandlerFunc(func(params order.CounselingUpdateOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation order.CounselingUpdateOrder has not yet been implemented")
		})
	}
	if api.CustomerSupportRemarksCreateCustomerSupportRemarkForMoveHandler == nil {
		api.CustomerSupportRemarksCreateCustomerSupportRemarkForMoveHandler = customer_support_remarks.CreateCustomerSupportRemarkForMoveHandlerFunc(func(params customer_support_remarks.CreateCustomerSupportRemarkForMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation customer_support_remarks.CreateCustomerSupportRemarkForMove has not yet been implemented")
		})
	}
	if api.EvaluationReportsCreateEvaluationReportForShipmentHandler == nil {
		api.EvaluationReportsCreateEvaluationReportForShipmentHandler = evaluation_reports.CreateEvaluationReportForShipmentHandlerFunc(func(params evaluation_reports.CreateEvaluationReportForShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation evaluation_reports.CreateEvaluationReportForShipment has not yet been implemented")
		})
	}
	if api.MtoShipmentCreateMTOShipmentHandler == nil {
		api.MtoShipmentCreateMTOShipmentHandler = mto_shipment.CreateMTOShipmentHandlerFunc(func(params mto_shipment.CreateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.CreateMTOShipment has not yet been implemented")
		})
	}
	if api.ShipmentCreateSITExtensionAsTOOHandler == nil {
		api.ShipmentCreateSITExtensionAsTOOHandler = shipment.CreateSITExtensionAsTOOHandlerFunc(func(params shipment.CreateSITExtensionAsTOOParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.CreateSITExtensionAsTOO has not yet been implemented")
		})
	}
	if api.CustomerSupportRemarksDeleteCustomerSupportRemarkHandler == nil {
		api.CustomerSupportRemarksDeleteCustomerSupportRemarkHandler = customer_support_remarks.DeleteCustomerSupportRemarkHandlerFunc(func(params customer_support_remarks.DeleteCustomerSupportRemarkParams) middleware.Responder {
			return middleware.NotImplemented("operation customer_support_remarks.DeleteCustomerSupportRemark has not yet been implemented")
		})
	}
	if api.EvaluationReportsDeleteEvaluationReportHandler == nil {
		api.EvaluationReportsDeleteEvaluationReportHandler = evaluation_reports.DeleteEvaluationReportHandlerFunc(func(params evaluation_reports.DeleteEvaluationReportParams) middleware.Responder {
			return middleware.NotImplemented("operation evaluation_reports.DeleteEvaluationReport has not yet been implemented")
		})
	}
	if api.ShipmentDeleteShipmentHandler == nil {
		api.ShipmentDeleteShipmentHandler = shipment.DeleteShipmentHandlerFunc(func(params shipment.DeleteShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.DeleteShipment has not yet been implemented")
		})
	}
	if api.ShipmentDenySITExtensionHandler == nil {
		api.ShipmentDenySITExtensionHandler = shipment.DenySITExtensionHandlerFunc(func(params shipment.DenySITExtensionParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.DenySITExtension has not yet been implemented")
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
	if api.CustomerSupportRemarksGetCustomerSupportRemarksForMoveHandler == nil {
		api.CustomerSupportRemarksGetCustomerSupportRemarksForMoveHandler = customer_support_remarks.GetCustomerSupportRemarksForMoveHandlerFunc(func(params customer_support_remarks.GetCustomerSupportRemarksForMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation customer_support_remarks.GetCustomerSupportRemarksForMove has not yet been implemented")
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
	if api.EvaluationReportsGetEvaluationReportHandler == nil {
		api.EvaluationReportsGetEvaluationReportHandler = evaluation_reports.GetEvaluationReportHandlerFunc(func(params evaluation_reports.GetEvaluationReportParams) middleware.Responder {
			return middleware.NotImplemented("operation evaluation_reports.GetEvaluationReport has not yet been implemented")
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
	if api.MoveGetMoveCounselingEvaluationReportsListHandler == nil {
		api.MoveGetMoveCounselingEvaluationReportsListHandler = move.GetMoveCounselingEvaluationReportsListHandlerFunc(func(params move.GetMoveCounselingEvaluationReportsListParams) middleware.Responder {
			return middleware.NotImplemented("operation move.GetMoveCounselingEvaluationReportsList has not yet been implemented")
		})
	}
	if api.MoveGetMoveHistoryHandler == nil {
		api.MoveGetMoveHistoryHandler = move.GetMoveHistoryHandlerFunc(func(params move.GetMoveHistoryParams) middleware.Responder {
			return middleware.NotImplemented("operation move.GetMoveHistory has not yet been implemented")
		})
	}
	if api.MoveGetMoveShipmentEvaluationReportsListHandler == nil {
		api.MoveGetMoveShipmentEvaluationReportsListHandler = move.GetMoveShipmentEvaluationReportsListHandlerFunc(func(params move.GetMoveShipmentEvaluationReportsListParams) middleware.Responder {
			return middleware.NotImplemented("operation move.GetMoveShipmentEvaluationReportsList has not yet been implemented")
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
	if api.QueuesGetServicesCounselingQueueHandler == nil {
		api.QueuesGetServicesCounselingQueueHandler = queues.GetServicesCounselingQueueHandlerFunc(func(params queues.GetServicesCounselingQueueParams) middleware.Responder {
			return middleware.NotImplemented("operation queues.GetServicesCounselingQueue has not yet been implemented")
		})
	}
	if api.MtoShipmentGetShipmentHandler == nil {
		api.MtoShipmentGetShipmentHandler = mto_shipment.GetShipmentHandlerFunc(func(params mto_shipment.GetShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.GetShipment has not yet been implemented")
		})
	}
	if api.PaymentRequestsGetShipmentsPaymentSITBalanceHandler == nil {
		api.PaymentRequestsGetShipmentsPaymentSITBalanceHandler = payment_requests.GetShipmentsPaymentSITBalanceHandlerFunc(func(params payment_requests.GetShipmentsPaymentSITBalanceParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.GetShipmentsPaymentSITBalance has not yet been implemented")
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
	if api.ShipmentRejectShipmentHandler == nil {
		api.ShipmentRejectShipmentHandler = shipment.RejectShipmentHandlerFunc(func(params shipment.RejectShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.RejectShipment has not yet been implemented")
		})
	}
	if api.ShipmentRequestShipmentCancellationHandler == nil {
		api.ShipmentRequestShipmentCancellationHandler = shipment.RequestShipmentCancellationHandlerFunc(func(params shipment.RequestShipmentCancellationParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.RequestShipmentCancellation has not yet been implemented")
		})
	}
	if api.ShipmentRequestShipmentDiversionHandler == nil {
		api.ShipmentRequestShipmentDiversionHandler = shipment.RequestShipmentDiversionHandlerFunc(func(params shipment.RequestShipmentDiversionParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.RequestShipmentDiversion has not yet been implemented")
		})
	}
	if api.ShipmentRequestShipmentReweighHandler == nil {
		api.ShipmentRequestShipmentReweighHandler = shipment.RequestShipmentReweighHandlerFunc(func(params shipment.RequestShipmentReweighParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.RequestShipmentReweigh has not yet been implemented")
		})
	}
	if api.MoveSearchMovesHandler == nil {
		api.MoveSearchMovesHandler = move.SearchMovesHandlerFunc(func(params move.SearchMovesParams) middleware.Responder {
			return middleware.NotImplemented("operation move.SearchMoves has not yet been implemented")
		})
	}
	if api.MoveSetFinancialReviewFlagHandler == nil {
		api.MoveSetFinancialReviewFlagHandler = move.SetFinancialReviewFlagHandlerFunc(func(params move.SetFinancialReviewFlagParams) middleware.Responder {
			return middleware.NotImplemented("operation move.SetFinancialReviewFlag has not yet been implemented")
		})
	}
	if api.TacTacValidationHandler == nil {
		api.TacTacValidationHandler = tac.TacValidationHandlerFunc(func(params tac.TacValidationParams) middleware.Responder {
			return middleware.NotImplemented("operation tac.TacValidation has not yet been implemented")
		})
	}
	if api.OrderUpdateAllowanceHandler == nil {
		api.OrderUpdateAllowanceHandler = order.UpdateAllowanceHandlerFunc(func(params order.UpdateAllowanceParams) middleware.Responder {
			return middleware.NotImplemented("operation order.UpdateAllowance has not yet been implemented")
		})
	}
	if api.OrderUpdateBillableWeightHandler == nil {
		api.OrderUpdateBillableWeightHandler = order.UpdateBillableWeightHandlerFunc(func(params order.UpdateBillableWeightParams) middleware.Responder {
			return middleware.NotImplemented("operation order.UpdateBillableWeight has not yet been implemented")
		})
	}
	if api.CustomerUpdateCustomerHandler == nil {
		api.CustomerUpdateCustomerHandler = customer.UpdateCustomerHandlerFunc(func(params customer.UpdateCustomerParams) middleware.Responder {
			return middleware.NotImplemented("operation customer.UpdateCustomer has not yet been implemented")
		})
	}
	if api.CustomerSupportRemarksUpdateCustomerSupportRemarkForMoveHandler == nil {
		api.CustomerSupportRemarksUpdateCustomerSupportRemarkForMoveHandler = customer_support_remarks.UpdateCustomerSupportRemarkForMoveHandlerFunc(func(params customer_support_remarks.UpdateCustomerSupportRemarkForMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation customer_support_remarks.UpdateCustomerSupportRemarkForMove has not yet been implemented")
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
	if api.MtoShipmentUpdateMTOShipmentHandler == nil {
		api.MtoShipmentUpdateMTOShipmentHandler = mto_shipment.UpdateMTOShipmentHandlerFunc(func(params mto_shipment.UpdateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.UpdateMTOShipment has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMTOStatusServiceCounselingCompletedHandler == nil {
		api.MoveTaskOrderUpdateMTOStatusServiceCounselingCompletedHandler = move_task_order.UpdateMTOStatusServiceCounselingCompletedHandlerFunc(func(params move_task_order.UpdateMTOStatusServiceCounselingCompletedParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMTOStatusServiceCounselingCompleted has not yet been implemented")
		})
	}
	if api.OrderUpdateMaxBillableWeightAsTIOHandler == nil {
		api.OrderUpdateMaxBillableWeightAsTIOHandler = order.UpdateMaxBillableWeightAsTIOHandlerFunc(func(params order.UpdateMaxBillableWeightAsTIOParams) middleware.Responder {
			return middleware.NotImplemented("operation order.UpdateMaxBillableWeightAsTIO has not yet been implemented")
		})
	}
	if api.MoveTaskOrderUpdateMoveTIORemarksHandler == nil {
		api.MoveTaskOrderUpdateMoveTIORemarksHandler = move_task_order.UpdateMoveTIORemarksHandlerFunc(func(params move_task_order.UpdateMoveTIORemarksParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMoveTIORemarks has not yet been implemented")
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
	if api.OrderUpdateOrderHandler == nil {
		api.OrderUpdateOrderHandler = order.UpdateOrderHandlerFunc(func(params order.UpdateOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation order.UpdateOrder has not yet been implemented")
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
