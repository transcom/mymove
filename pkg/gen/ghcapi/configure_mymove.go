// This file is safe to edit. Once it exists it will not be overwritten

package ghcapi

import (
	"crypto/tls"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/addresses"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/application_parameters"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/calendar"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer_support_remarks"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/evaluation_reports"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ghc_documents"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/lines_of_accounting"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/move_task_order"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_agent"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/office_users"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/order"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/orders"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/pws_violations"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/re_service_items"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/report_violations"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/shipment"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/tac"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/transportation_office"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/uploads"
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
	api.MultipartformConsumer = runtime.DiscardConsumer

	api.BinProducer = runtime.ByteStreamProducer()
	api.JSONProducer = runtime.JSONProducer()
	api.TextEventStreamProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("textEventStream producer has not yet been implemented")
	})

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// ppm.CreatePPMUploadMaxParseMemory = 32 << 20
	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// uploads.CreateUploadMaxParseMemory = 32 << 20
	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// move.UploadAdditionalDocumentsMaxParseMemory = 32 << 20
	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// order.UploadAmendedOrdersMaxParseMemory = 32 << 20

	if api.OrderAcknowledgeExcessUnaccompaniedBaggageWeightRiskHandler == nil {
		api.OrderAcknowledgeExcessUnaccompaniedBaggageWeightRiskHandler = order.AcknowledgeExcessUnaccompaniedBaggageWeightRiskHandlerFunc(func(params order.AcknowledgeExcessUnaccompaniedBaggageWeightRiskParams) middleware.Responder {
			return middleware.NotImplemented("operation order.AcknowledgeExcessUnaccompaniedBaggageWeightRisk has not yet been implemented")
		})
	}
	if api.OrderAcknowledgeExcessWeightRiskHandler == nil {
		api.OrderAcknowledgeExcessWeightRiskHandler = order.AcknowledgeExcessWeightRiskHandlerFunc(func(params order.AcknowledgeExcessWeightRiskParams) middleware.Responder {
			return middleware.NotImplemented("operation order.AcknowledgeExcessWeightRisk has not yet been implemented")
		})
	}
	if api.EvaluationReportsAddAppealToSeriousIncidentHandler == nil {
		api.EvaluationReportsAddAppealToSeriousIncidentHandler = evaluation_reports.AddAppealToSeriousIncidentHandlerFunc(func(params evaluation_reports.AddAppealToSeriousIncidentParams) middleware.Responder {
			return middleware.NotImplemented("operation evaluation_reports.AddAppealToSeriousIncident has not yet been implemented")
		})
	}
	if api.EvaluationReportsAddAppealToViolationHandler == nil {
		api.EvaluationReportsAddAppealToViolationHandler = evaluation_reports.AddAppealToViolationHandlerFunc(func(params evaluation_reports.AddAppealToViolationParams) middleware.Responder {
			return middleware.NotImplemented("operation evaluation_reports.AddAppealToViolation has not yet been implemented")
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
	if api.ShipmentApproveShipmentsHandler == nil {
		api.ShipmentApproveShipmentsHandler = shipment.ApproveShipmentsHandlerFunc(func(params shipment.ApproveShipmentsParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.ApproveShipments has not yet been implemented")
		})
	}
	if api.ReportViolationsAssociateReportViolationsHandler == nil {
		api.ReportViolationsAssociateReportViolationsHandler = report_violations.AssociateReportViolationsHandlerFunc(func(params report_violations.AssociateReportViolationsParams) middleware.Responder {
			return middleware.NotImplemented("operation report_violations.AssociateReportViolations has not yet been implemented")
		})
	}
	if api.PaymentRequestsBulkDownloadHandler == nil {
		api.PaymentRequestsBulkDownloadHandler = payment_requests.BulkDownloadHandlerFunc(func(params payment_requests.BulkDownloadParams) middleware.Responder {
			return middleware.NotImplemented("operation payment_requests.BulkDownload has not yet been implemented")
		})
	}
	if api.MoveCheckForLockedMovesAndUnlockHandler == nil {
		api.MoveCheckForLockedMovesAndUnlockHandler = move.CheckForLockedMovesAndUnlockHandlerFunc(func(params move.CheckForLockedMovesAndUnlockParams) middleware.Responder {
			return middleware.NotImplemented("operation move.CheckForLockedMovesAndUnlock has not yet been implemented")
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
	if api.ShipmentCreateApprovedSITDurationUpdateHandler == nil {
		api.ShipmentCreateApprovedSITDurationUpdateHandler = shipment.CreateApprovedSITDurationUpdateHandlerFunc(func(params shipment.CreateApprovedSITDurationUpdateParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.CreateApprovedSITDurationUpdate has not yet been implemented")
		})
	}
	if api.CustomerSupportRemarksCreateCustomerSupportRemarkForMoveHandler == nil {
		api.CustomerSupportRemarksCreateCustomerSupportRemarkForMoveHandler = customer_support_remarks.CreateCustomerSupportRemarkForMoveHandlerFunc(func(params customer_support_remarks.CreateCustomerSupportRemarkForMoveParams) middleware.Responder {
			return middleware.NotImplemented("operation customer_support_remarks.CreateCustomerSupportRemarkForMove has not yet been implemented")
		})
	}
	if api.CustomerCreateCustomerWithOktaOptionHandler == nil {
		api.CustomerCreateCustomerWithOktaOptionHandler = customer.CreateCustomerWithOktaOptionHandlerFunc(func(params customer.CreateCustomerWithOktaOptionParams) middleware.Responder {
			return middleware.NotImplemented("operation customer.CreateCustomerWithOktaOption has not yet been implemented")
		})
	}
	if api.GhcDocumentsCreateDocumentHandler == nil {
		api.GhcDocumentsCreateDocumentHandler = ghc_documents.CreateDocumentHandlerFunc(func(params ghc_documents.CreateDocumentParams) middleware.Responder {
			return middleware.NotImplemented("operation ghc_documents.CreateDocument has not yet been implemented")
		})
	}
	if api.EvaluationReportsCreateEvaluationReportHandler == nil {
		api.EvaluationReportsCreateEvaluationReportHandler = evaluation_reports.CreateEvaluationReportHandlerFunc(func(params evaluation_reports.CreateEvaluationReportParams) middleware.Responder {
			return middleware.NotImplemented("operation evaluation_reports.CreateEvaluationReport has not yet been implemented")
		})
	}
	if api.MtoShipmentCreateMTOShipmentHandler == nil {
		api.MtoShipmentCreateMTOShipmentHandler = mto_shipment.CreateMTOShipmentHandlerFunc(func(params mto_shipment.CreateMTOShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_shipment.CreateMTOShipment has not yet been implemented")
		})
	}
	if api.PpmCreateMovingExpenseHandler == nil {
		api.PpmCreateMovingExpenseHandler = ppm.CreateMovingExpenseHandlerFunc(func(params ppm.CreateMovingExpenseParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.CreateMovingExpense has not yet been implemented")
		})
	}
	if api.OrderCreateOrderHandler == nil {
		api.OrderCreateOrderHandler = order.CreateOrderHandlerFunc(func(params order.CreateOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation order.CreateOrder has not yet been implemented")
		})
	}
	if api.PpmCreatePPMUploadHandler == nil {
		api.PpmCreatePPMUploadHandler = ppm.CreatePPMUploadHandlerFunc(func(params ppm.CreatePPMUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.CreatePPMUpload has not yet been implemented")
		})
	}
	if api.PpmCreateProGearWeightTicketHandler == nil {
		api.PpmCreateProGearWeightTicketHandler = ppm.CreateProGearWeightTicketHandlerFunc(func(params ppm.CreateProGearWeightTicketParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.CreateProGearWeightTicket has not yet been implemented")
		})
	}
	if api.OfficeUsersCreateRequestedOfficeUserHandler == nil {
		api.OfficeUsersCreateRequestedOfficeUserHandler = office_users.CreateRequestedOfficeUserHandlerFunc(func(params office_users.CreateRequestedOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office_users.CreateRequestedOfficeUser has not yet been implemented")
		})
	}
	if api.ShipmentCreateTerminationHandler == nil {
		api.ShipmentCreateTerminationHandler = shipment.CreateTerminationHandlerFunc(func(params shipment.CreateTerminationParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.CreateTermination has not yet been implemented")
		})
	}
	if api.UploadsCreateUploadHandler == nil {
		api.UploadsCreateUploadHandler = uploads.CreateUploadHandlerFunc(func(params uploads.CreateUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.CreateUpload has not yet been implemented")
		})
	}
	if api.PpmCreateWeightTicketHandler == nil {
		api.PpmCreateWeightTicketHandler = ppm.CreateWeightTicketHandlerFunc(func(params ppm.CreateWeightTicketParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.CreateWeightTicket has not yet been implemented")
		})
	}
	if api.MoveDeleteAssignedOfficeUserHandler == nil {
		api.MoveDeleteAssignedOfficeUserHandler = move.DeleteAssignedOfficeUserHandlerFunc(func(params move.DeleteAssignedOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation move.DeleteAssignedOfficeUser has not yet been implemented")
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
	if api.PpmDeleteMovingExpenseHandler == nil {
		api.PpmDeleteMovingExpenseHandler = ppm.DeleteMovingExpenseHandlerFunc(func(params ppm.DeleteMovingExpenseParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.DeleteMovingExpense has not yet been implemented")
		})
	}
	if api.PpmDeleteProGearWeightTicketHandler == nil {
		api.PpmDeleteProGearWeightTicketHandler = ppm.DeleteProGearWeightTicketHandlerFunc(func(params ppm.DeleteProGearWeightTicketParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.DeleteProGearWeightTicket has not yet been implemented")
		})
	}
	if api.ShipmentDeleteShipmentHandler == nil {
		api.ShipmentDeleteShipmentHandler = shipment.DeleteShipmentHandlerFunc(func(params shipment.DeleteShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.DeleteShipment has not yet been implemented")
		})
	}
	if api.UploadsDeleteUploadHandler == nil {
		api.UploadsDeleteUploadHandler = uploads.DeleteUploadHandlerFunc(func(params uploads.DeleteUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.DeleteUpload has not yet been implemented")
		})
	}
	if api.PpmDeleteWeightTicketHandler == nil {
		api.PpmDeleteWeightTicketHandler = ppm.DeleteWeightTicketHandlerFunc(func(params ppm.DeleteWeightTicketParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.DeleteWeightTicket has not yet been implemented")
		})
	}
	if api.ShipmentDenySITExtensionHandler == nil {
		api.ShipmentDenySITExtensionHandler = shipment.DenySITExtensionHandlerFunc(func(params shipment.DenySITExtensionParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.DenySITExtension has not yet been implemented")
		})
	}
	if api.EvaluationReportsDownloadEvaluationReportHandler == nil {
		api.EvaluationReportsDownloadEvaluationReportHandler = evaluation_reports.DownloadEvaluationReportHandlerFunc(func(params evaluation_reports.DownloadEvaluationReportParams) middleware.Responder {
			return middleware.NotImplemented("operation evaluation_reports.DownloadEvaluationReport has not yet been implemented")
		})
	}
	if api.MtoAgentFetchMTOAgentListHandler == nil {
		api.MtoAgentFetchMTOAgentListHandler = mto_agent.FetchMTOAgentListHandlerFunc(func(params mto_agent.FetchMTOAgentListParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_agent.FetchMTOAgentList has not yet been implemented")
		})
	}
	if api.PpmFinishDocumentReviewHandler == nil {
		api.PpmFinishDocumentReviewHandler = ppm.FinishDocumentReviewHandlerFunc(func(params ppm.FinishDocumentReviewParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.FinishDocumentReview has not yet been implemented")
		})
	}
	if api.ReServiceItemsGetAllReServiceItemsHandler == nil {
		api.ReServiceItemsGetAllReServiceItemsHandler = re_service_items.GetAllReServiceItemsHandlerFunc(func(params re_service_items.GetAllReServiceItemsParams) middleware.Responder {
			return middleware.NotImplemented("operation re_service_items.GetAllReServiceItems has not yet been implemented")
		})
	}
	if api.QueuesGetBulkAssignmentDataHandler == nil {
		api.QueuesGetBulkAssignmentDataHandler = queues.GetBulkAssignmentDataHandlerFunc(func(params queues.GetBulkAssignmentDataParams) middleware.Responder {
			return middleware.NotImplemented("operation queues.GetBulkAssignmentData has not yet been implemented")
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
	if api.QueuesGetDestinationRequestsQueueHandler == nil {
		api.QueuesGetDestinationRequestsQueueHandler = queues.GetDestinationRequestsQueueHandlerFunc(func(params queues.GetDestinationRequestsQueueParams) middleware.Responder {
			return middleware.NotImplemented("operation queues.GetDestinationRequestsQueue has not yet been implemented")
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
	if api.AddressesGetLocationByZipCityStateHandler == nil {
		api.AddressesGetLocationByZipCityStateHandler = addresses.GetLocationByZipCityStateHandlerFunc(func(params addresses.GetLocationByZipCityStateParams) middleware.Responder {
			return middleware.NotImplemented("operation addresses.GetLocationByZipCityState has not yet been implemented")
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
	if api.PpmGetPPMActualWeightHandler == nil {
		api.PpmGetPPMActualWeightHandler = ppm.GetPPMActualWeightHandlerFunc(func(params ppm.GetPPMActualWeightParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.GetPPMActualWeight has not yet been implemented")
		})
	}
	if api.PpmGetPPMCloseoutHandler == nil {
		api.PpmGetPPMCloseoutHandler = ppm.GetPPMCloseoutHandlerFunc(func(params ppm.GetPPMCloseoutParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.GetPPMCloseout has not yet been implemented")
		})
	}
	if api.PpmGetPPMDocumentsHandler == nil {
		api.PpmGetPPMDocumentsHandler = ppm.GetPPMDocumentsHandlerFunc(func(params ppm.GetPPMDocumentsParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.GetPPMDocuments has not yet been implemented")
		})
	}
	if api.PpmGetPPMSITEstimatedCostHandler == nil {
		api.PpmGetPPMSITEstimatedCostHandler = ppm.GetPPMSITEstimatedCostHandlerFunc(func(params ppm.GetPPMSITEstimatedCostParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.GetPPMSITEstimatedCost has not yet been implemented")
		})
	}
	if api.PwsViolationsGetPWSViolationsHandler == nil {
		api.PwsViolationsGetPWSViolationsHandler = pws_violations.GetPWSViolationsHandlerFunc(func(params pws_violations.GetPWSViolationsParams) middleware.Responder {
			return middleware.NotImplemented("operation pws_violations.GetPWSViolations has not yet been implemented")
		})
	}
	if api.ApplicationParametersGetParamHandler == nil {
		api.ApplicationParametersGetParamHandler = application_parameters.GetParamHandlerFunc(func(params application_parameters.GetParamParams) middleware.Responder {
			return middleware.NotImplemented("operation application_parameters.GetParam has not yet been implemented")
		})
	}
	if api.OrdersGetPayGradesHandler == nil {
		api.OrdersGetPayGradesHandler = orders.GetPayGradesHandlerFunc(func(params orders.GetPayGradesParams) middleware.Responder {
			return middleware.NotImplemented("operation orders.GetPayGrades has not yet been implemented")
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
	if api.ReportViolationsGetReportViolationsByReportIDHandler == nil {
		api.ReportViolationsGetReportViolationsByReportIDHandler = report_violations.GetReportViolationsByReportIDHandlerFunc(func(params report_violations.GetReportViolationsByReportIDParams) middleware.Responder {
			return middleware.NotImplemented("operation report_violations.GetReportViolationsByReportID has not yet been implemented")
		})
	}
	if api.QueuesGetServicesCounselingOriginListHandler == nil {
		api.QueuesGetServicesCounselingOriginListHandler = queues.GetServicesCounselingOriginListHandlerFunc(func(params queues.GetServicesCounselingOriginListParams) middleware.Responder {
			return middleware.NotImplemented("operation queues.GetServicesCounselingOriginList has not yet been implemented")
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
	if api.TransportationOfficeGetTransportationOfficesHandler == nil {
		api.TransportationOfficeGetTransportationOfficesHandler = transportation_office.GetTransportationOfficesHandlerFunc(func(params transportation_office.GetTransportationOfficesParams) middleware.Responder {
			return middleware.NotImplemented("operation transportation_office.GetTransportationOffices has not yet been implemented")
		})
	}
	if api.TransportationOfficeGetTransportationOfficesGBLOCsHandler == nil {
		api.TransportationOfficeGetTransportationOfficesGBLOCsHandler = transportation_office.GetTransportationOfficesGBLOCsHandlerFunc(func(params transportation_office.GetTransportationOfficesGBLOCsParams) middleware.Responder {
			return middleware.NotImplemented("operation transportation_office.GetTransportationOfficesGBLOCs has not yet been implemented")
		})
	}
	if api.TransportationOfficeGetTransportationOfficesOpenHandler == nil {
		api.TransportationOfficeGetTransportationOfficesOpenHandler = transportation_office.GetTransportationOfficesOpenHandlerFunc(func(params transportation_office.GetTransportationOfficesOpenParams) middleware.Responder {
			return middleware.NotImplemented("operation transportation_office.GetTransportationOfficesOpen has not yet been implemented")
		})
	}
	if api.UploadsGetUploadHandler == nil {
		api.UploadsGetUploadHandler = uploads.GetUploadHandlerFunc(func(params uploads.GetUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.GetUpload has not yet been implemented")
		})
	}
	if api.UploadsGetUploadStatusHandler == nil {
		api.UploadsGetUploadStatusHandler = uploads.GetUploadStatusHandlerFunc(func(params uploads.GetUploadStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.GetUploadStatus has not yet been implemented")
		})
	}
	if api.CalendarIsDateWeekendHolidayHandler == nil {
		api.CalendarIsDateWeekendHolidayHandler = calendar.IsDateWeekendHolidayHandlerFunc(func(params calendar.IsDateWeekendHolidayParams) middleware.Responder {
			return middleware.NotImplemented("operation calendar.IsDateWeekendHoliday has not yet been implemented")
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
	if api.QueuesListPrimeMovesHandler == nil {
		api.QueuesListPrimeMovesHandler = queues.ListPrimeMovesHandlerFunc(func(params queues.ListPrimeMovesParams) middleware.Responder {
			return middleware.NotImplemented("operation queues.ListPrimeMoves has not yet been implemented")
		})
	}
	if api.MoveMoveCancelerHandler == nil {
		api.MoveMoveCancelerHandler = move.MoveCancelerHandlerFunc(func(params move.MoveCancelerParams) middleware.Responder {
			return middleware.NotImplemented("operation move.MoveCanceler has not yet been implemented")
		})
	}
	if api.ShipmentRejectShipmentHandler == nil {
		api.ShipmentRejectShipmentHandler = shipment.RejectShipmentHandlerFunc(func(params shipment.RejectShipmentParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.RejectShipment has not yet been implemented")
		})
	}
	if api.LinesOfAccountingRequestLineOfAccountingHandler == nil {
		api.LinesOfAccountingRequestLineOfAccountingHandler = lines_of_accounting.RequestLineOfAccountingHandlerFunc(func(params lines_of_accounting.RequestLineOfAccountingParams) middleware.Responder {
			return middleware.NotImplemented("operation lines_of_accounting.RequestLineOfAccounting has not yet been implemented")
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
	if api.ShipmentReviewShipmentAddressUpdateHandler == nil {
		api.ShipmentReviewShipmentAddressUpdateHandler = shipment.ReviewShipmentAddressUpdateHandlerFunc(func(params shipment.ReviewShipmentAddressUpdateParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.ReviewShipmentAddressUpdate has not yet been implemented")
		})
	}
	if api.QueuesSaveBulkAssignmentDataHandler == nil {
		api.QueuesSaveBulkAssignmentDataHandler = queues.SaveBulkAssignmentDataHandlerFunc(func(params queues.SaveBulkAssignmentDataParams) middleware.Responder {
			return middleware.NotImplemented("operation queues.SaveBulkAssignmentData has not yet been implemented")
		})
	}
	if api.EvaluationReportsSaveEvaluationReportHandler == nil {
		api.EvaluationReportsSaveEvaluationReportHandler = evaluation_reports.SaveEvaluationReportHandlerFunc(func(params evaluation_reports.SaveEvaluationReportParams) middleware.Responder {
			return middleware.NotImplemented("operation evaluation_reports.SaveEvaluationReport has not yet been implemented")
		})
	}
	if api.AddressesSearchCountriesHandler == nil {
		api.AddressesSearchCountriesHandler = addresses.SearchCountriesHandlerFunc(func(params addresses.SearchCountriesParams) middleware.Responder {
			return middleware.NotImplemented("operation addresses.SearchCountries has not yet been implemented")
		})
	}
	if api.CustomerSearchCustomersHandler == nil {
		api.CustomerSearchCustomersHandler = customer.SearchCustomersHandlerFunc(func(params customer.SearchCustomersParams) middleware.Responder {
			return middleware.NotImplemented("operation customer.SearchCustomers has not yet been implemented")
		})
	}
	if api.MoveSearchMovesHandler == nil {
		api.MoveSearchMovesHandler = move.SearchMovesHandlerFunc(func(params move.SearchMovesParams) middleware.Responder {
			return middleware.NotImplemented("operation move.SearchMoves has not yet been implemented")
		})
	}
	if api.PpmSendPPMToCustomerHandler == nil {
		api.PpmSendPPMToCustomerHandler = ppm.SendPPMToCustomerHandlerFunc(func(params ppm.SendPPMToCustomerParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.SendPPMToCustomer has not yet been implemented")
		})
	}
	if api.MoveSetFinancialReviewFlagHandler == nil {
		api.MoveSetFinancialReviewFlagHandler = move.SetFinancialReviewFlagHandlerFunc(func(params move.SetFinancialReviewFlagParams) middleware.Responder {
			return middleware.NotImplemented("operation move.SetFinancialReviewFlag has not yet been implemented")
		})
	}
	if api.PpmShowAOAPacketHandler == nil {
		api.PpmShowAOAPacketHandler = ppm.ShowAOAPacketHandlerFunc(func(params ppm.ShowAOAPacketParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.ShowAOAPacket has not yet been implemented")
		})
	}
	if api.TransportationOfficeShowCounselingOfficesHandler == nil {
		api.TransportationOfficeShowCounselingOfficesHandler = transportation_office.ShowCounselingOfficesHandlerFunc(func(params transportation_office.ShowCounselingOfficesParams) middleware.Responder {
			return middleware.NotImplemented("operation transportation_office.ShowCounselingOffices has not yet been implemented")
		})
	}
	if api.PpmShowPaymentPacketHandler == nil {
		api.PpmShowPaymentPacketHandler = ppm.ShowPaymentPacketHandlerFunc(func(params ppm.ShowPaymentPacketParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.ShowPaymentPacket has not yet been implemented")
		})
	}
	if api.EvaluationReportsSubmitEvaluationReportHandler == nil {
		api.EvaluationReportsSubmitEvaluationReportHandler = evaluation_reports.SubmitEvaluationReportHandlerFunc(func(params evaluation_reports.SubmitEvaluationReportParams) middleware.Responder {
			return middleware.NotImplemented("operation evaluation_reports.SubmitEvaluationReport has not yet been implemented")
		})
	}
	if api.PpmSubmitPPMShipmentDocumentationHandler == nil {
		api.PpmSubmitPPMShipmentDocumentationHandler = ppm.SubmitPPMShipmentDocumentationHandlerFunc(func(params ppm.SubmitPPMShipmentDocumentationParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.SubmitPPMShipmentDocumentation has not yet been implemented")
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
	if api.MoveUpdateAssignedOfficeUserHandler == nil {
		api.MoveUpdateAssignedOfficeUserHandler = move.UpdateAssignedOfficeUserHandlerFunc(func(params move.UpdateAssignedOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation move.UpdateAssignedOfficeUser has not yet been implemented")
		})
	}
	if api.OrderUpdateBillableWeightHandler == nil {
		api.OrderUpdateBillableWeightHandler = order.UpdateBillableWeightHandlerFunc(func(params order.UpdateBillableWeightParams) middleware.Responder {
			return middleware.NotImplemented("operation order.UpdateBillableWeight has not yet been implemented")
		})
	}
	if api.MoveUpdateCloseoutOfficeHandler == nil {
		api.MoveUpdateCloseoutOfficeHandler = move.UpdateCloseoutOfficeHandlerFunc(func(params move.UpdateCloseoutOfficeParams) middleware.Responder {
			return middleware.NotImplemented("operation move.UpdateCloseoutOffice has not yet been implemented")
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
	if api.MoveTaskOrderUpdateMTOReviewedBillableWeightsAtHandler == nil {
		api.MoveTaskOrderUpdateMTOReviewedBillableWeightsAtHandler = move_task_order.UpdateMTOReviewedBillableWeightsAtHandlerFunc(func(params move_task_order.UpdateMTOReviewedBillableWeightsAtParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMTOReviewedBillableWeightsAt has not yet been implemented")
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
	if api.MoveTaskOrderUpdateMoveTaskOrderStatusHandler == nil {
		api.MoveTaskOrderUpdateMoveTaskOrderStatusHandler = move_task_order.UpdateMoveTaskOrderStatusHandlerFunc(func(params move_task_order.UpdateMoveTaskOrderStatusParams) middleware.Responder {
			return middleware.NotImplemented("operation move_task_order.UpdateMoveTaskOrderStatus has not yet been implemented")
		})
	}
	if api.PpmUpdateMovingExpenseHandler == nil {
		api.PpmUpdateMovingExpenseHandler = ppm.UpdateMovingExpenseHandlerFunc(func(params ppm.UpdateMovingExpenseParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.UpdateMovingExpense has not yet been implemented")
		})
	}
	if api.OfficeUsersUpdateOfficeUserHandler == nil {
		api.OfficeUsersUpdateOfficeUserHandler = office_users.UpdateOfficeUserHandlerFunc(func(params office_users.UpdateOfficeUserParams) middleware.Responder {
			return middleware.NotImplemented("operation office_users.UpdateOfficeUser has not yet been implemented")
		})
	}
	if api.OrderUpdateOrderHandler == nil {
		api.OrderUpdateOrderHandler = order.UpdateOrderHandlerFunc(func(params order.UpdateOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation order.UpdateOrder has not yet been implemented")
		})
	}
	if api.PpmUpdatePPMSITHandler == nil {
		api.PpmUpdatePPMSITHandler = ppm.UpdatePPMSITHandlerFunc(func(params ppm.UpdatePPMSITParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.UpdatePPMSIT has not yet been implemented")
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
	if api.PpmUpdateProGearWeightTicketHandler == nil {
		api.PpmUpdateProGearWeightTicketHandler = ppm.UpdateProGearWeightTicketHandlerFunc(func(params ppm.UpdateProGearWeightTicketParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.UpdateProGearWeightTicket has not yet been implemented")
		})
	}
	if api.ShipmentUpdateSITServiceItemCustomerExpenseHandler == nil {
		api.ShipmentUpdateSITServiceItemCustomerExpenseHandler = shipment.UpdateSITServiceItemCustomerExpenseHandlerFunc(func(params shipment.UpdateSITServiceItemCustomerExpenseParams) middleware.Responder {
			return middleware.NotImplemented("operation shipment.UpdateSITServiceItemCustomerExpense has not yet been implemented")
		})
	}
	if api.MtoServiceItemUpdateServiceItemSitEntryDateHandler == nil {
		api.MtoServiceItemUpdateServiceItemSitEntryDateHandler = mto_service_item.UpdateServiceItemSitEntryDateHandlerFunc(func(params mto_service_item.UpdateServiceItemSitEntryDateParams) middleware.Responder {
			return middleware.NotImplemented("operation mto_service_item.UpdateServiceItemSitEntryDate has not yet been implemented")
		})
	}
	if api.UploadsUpdateUploadHandler == nil {
		api.UploadsUpdateUploadHandler = uploads.UpdateUploadHandlerFunc(func(params uploads.UpdateUploadParams) middleware.Responder {
			return middleware.NotImplemented("operation uploads.UpdateUpload has not yet been implemented")
		})
	}
	if api.PpmUpdateWeightTicketHandler == nil {
		api.PpmUpdateWeightTicketHandler = ppm.UpdateWeightTicketHandlerFunc(func(params ppm.UpdateWeightTicketParams) middleware.Responder {
			return middleware.NotImplemented("operation ppm.UpdateWeightTicket has not yet been implemented")
		})
	}
	if api.MoveUploadAdditionalDocumentsHandler == nil {
		api.MoveUploadAdditionalDocumentsHandler = move.UploadAdditionalDocumentsHandlerFunc(func(params move.UploadAdditionalDocumentsParams) middleware.Responder {
			return middleware.NotImplemented("operation move.UploadAdditionalDocuments has not yet been implemented")
		})
	}
	if api.OrderUploadAmendedOrdersHandler == nil {
		api.OrderUploadAmendedOrdersHandler = order.UploadAmendedOrdersHandlerFunc(func(params order.UploadAmendedOrdersParams) middleware.Responder {
			return middleware.NotImplemented("operation order.UploadAmendedOrders has not yet been implemented")
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
