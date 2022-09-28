package ghcapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	customerserviceremarks "github.com/transcom/mymove/pkg/services/customer_support_remarks"
	evaluationreport "github.com/transcom/mymove/pkg/services/evaluation_report"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/move"
	movehistory "github.com/transcom/mymove/pkg/services/move_history"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/office_user/customer"
	"github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	order "github.com/transcom/mymove/pkg/services/order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	paymentserviceitem "github.com/transcom/mymove/pkg/services/payment_service_item"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	pwsviolation "github.com/transcom/mymove/pkg/services/pws_violation"
	"github.com/transcom/mymove/pkg/services/query"
	reportviolation "github.com/transcom/mymove/pkg/services/report_violation"
)

// NewGhcAPIHandler returns a handler for the GHC API
func NewGhcAPIHandler(handlerConfig handlers.HandlerConfig) *ghcops.MymoveAPI {
	ghcSpec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcAPI := ghcops.NewMymoveAPI(ghcSpec)
	queryBuilder := query.NewQueryBuilder()
	moveRouter := move.NewMoveRouter()
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
		moveRouter,
	)
	shipmentSITStatus := mtoshipment.NewShipmentSITStatus()

	ghcAPI.ServeError = handlers.ServeCustomError

	ghcAPI.MoveGetMoveHandler = GetMoveHandler{
		HandlerConfig: handlerConfig,
		MoveFetcher:   move.NewMoveFetcher(),
	}

	ghcAPI.MoveGetMoveHistoryHandler = GetMoveHistoryHandler{
		HandlerConfig:      handlerConfig,
		MoveHistoryFetcher: movehistory.NewMoveHistoryFetcher(),
	}

	ghcAPI.CustomerSupportRemarksGetCustomerSupportRemarksForMoveHandler = ListCustomerSupportRemarksHandler{
		HandlerConfig:                 handlerConfig,
		CustomerSupportRemarksFetcher: customerserviceremarks.NewCustomerSupportRemarks(),
	}

	ghcAPI.CustomerSupportRemarksCreateCustomerSupportRemarkForMoveHandler = CreateCustomerSupportRemarksHandler{
		HandlerConfig:                 handlerConfig,
		CustomerSupportRemarksCreator: customerserviceremarks.NewCustomerSupportRemarksCreator(),
	}

	ghcAPI.CustomerSupportRemarksUpdateCustomerSupportRemarkForMoveHandler = UpdateCustomerSupportRemarkHandler{
		HandlerConfig:                handlerConfig,
		CustomerSupportRemarkUpdater: customerserviceremarks.NewCustomerSupportRemarkUpdater(),
	}

	ghcAPI.CustomerSupportRemarksDeleteCustomerSupportRemarkHandler = DeleteCustomerSupportRemarkHandler{
		HandlerConfig:                handlerConfig,
		CustomerSupportRemarkDeleter: customerserviceremarks.NewCustomerSupportRemarkDeleter(),
	}

	ghcAPI.EvaluationReportsCreateEvaluationReportHandler = CreateEvaluationReportHandler{
		HandlerConfig:           handlerConfig,
		EvaluationReportCreator: evaluationreport.NewEvaluationReportCreator(),
	}

	ghcAPI.MoveGetMoveCounselingEvaluationReportsListHandler = GetCounselingEvaluationReportsHandler{
		HandlerConfig:           handlerConfig,
		EvaluationReportFetcher: evaluationreport.NewEvaluationReportFetcher(),
	}

	ghcAPI.MoveGetMoveShipmentEvaluationReportsListHandler = GetShipmentEvaluationReportsHandler{
		HandlerConfig:           handlerConfig,
		EvaluationReportFetcher: evaluationreport.NewEvaluationReportFetcher(),
	}

	ghcAPI.EvaluationReportsGetEvaluationReportHandler = GetEvaluationReportHandler{
		HandlerConfig:           handlerConfig,
		EvaluationReportFetcher: evaluationreport.NewEvaluationReportFetcher(),
	}

	ghcAPI.EvaluationReportsDeleteEvaluationReportHandler = DeleteEvaluationReportHandler{
		HandlerConfig:           handlerConfig,
		EvaluationReportDeleter: evaluationreport.NewEvaluationReportDeleter(),
	}

	ghcAPI.EvaluationReportsSaveEvaluationReportHandler = SaveEvaluationReportHandler{
		HandlerConfig:           handlerConfig,
		EvaluationReportUpdater: evaluationreport.NewEvaluationReportUpdater(),
	}

	ghcAPI.EvaluationReportsSubmitEvaluationReportHandler = SubmitEvaluationReportHandler{
		HandlerConfig:           handlerConfig,
		EvaluationReportUpdater: evaluationreport.NewEvaluationReportUpdater(),
	}

	ghcAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{
		HandlerConfig:         handlerConfig,
		MTOServiceItemUpdater: mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter),
		Fetcher:               fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.MtoServiceItemListMTOServiceItemsHandler = ListMTOServiceItemsHandler{
		handlerConfig,
		fetch.NewListFetcher(queryBuilder),
		fetch.NewFetcher(queryBuilder),
	}

	ghcAPI.PaymentRequestsGetPaymentRequestHandler = GetPaymentRequestHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestFetcher(),
	}

	ghcAPI.PaymentRequestsGetPaymentRequestsForMoveHandler = GetPaymentRequestForMoveHandler{
		HandlerConfig:             handlerConfig,
		PaymentRequestListFetcher: paymentrequest.NewPaymentRequestListFetcher(),
	}

	ghcAPI.PaymentRequestsUpdatePaymentRequestStatusHandler = UpdatePaymentRequestStatusHandler{
		HandlerConfig:               handlerConfig,
		PaymentRequestStatusUpdater: paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
		PaymentRequestFetcher:       paymentrequest.NewPaymentRequestFetcher(),
	}

	ghcAPI.PaymentServiceItemUpdatePaymentServiceItemStatusHandler = UpdatePaymentServiceItemStatusHandler{
		HandlerConfig:                   handlerConfig,
		PaymentServiceItemStatusUpdater: paymentserviceitem.NewPaymentServiceItemStatusUpdater(),
	}

	ghcAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}
	ghcAPI.MoveSetFinancialReviewFlagHandler = SetFinancialReviewFlagHandler{
		handlerConfig,
		move.NewFinancialReviewFlagSetter(),
	}

	ghcAPI.CustomerGetCustomerHandler = GetCustomerHandler{
		handlerConfig,
		customer.NewCustomerFetcher(),
	}
	ghcAPI.CustomerUpdateCustomerHandler = UpdateCustomerHandler{
		handlerConfig,
		customer.NewCustomerUpdater(),
	}
	ghcAPI.OrderGetOrderHandler = GetOrdersHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}
	ghcAPI.OrderCounselingUpdateOrderHandler = CounselingUpdateOrderHandler{
		handlerConfig,
		order.NewOrderUpdater(moveRouter),
	}

	ghcAPI.OrderUpdateOrderHandler = UpdateOrderHandler{
		handlerConfig,
		order.NewOrderUpdater(moveRouter),
		moveTaskOrderUpdater,
	}

	ghcAPI.OrderUpdateAllowanceHandler = UpdateAllowanceHandler{
		handlerConfig,
		order.NewOrderUpdater(moveRouter),
	}
	ghcAPI.OrderCounselingUpdateAllowanceHandler = CounselingUpdateAllowanceHandler{
		handlerConfig,
		order.NewOrderUpdater(moveRouter),
	}
	ghcAPI.OrderUpdateBillableWeightHandler = UpdateBillableWeightHandler{
		handlerConfig,
		order.NewExcessWeightRiskManager(moveRouter),
	}

	ghcAPI.OrderUpdateMaxBillableWeightAsTIOHandler = UpdateMaxBillableWeightAsTIOHandler{
		handlerConfig,
		order.NewExcessWeightRiskManager(moveRouter),
	}

	ghcAPI.OrderAcknowledgeExcessWeightRiskHandler = AcknowledgeExcessWeightRiskHandler{
		handlerConfig,
		order.NewExcessWeightRiskManager(moveRouter),
	}

	ghcAPI.MoveTaskOrderUpdateMoveTaskOrderStatusHandler = UpdateMoveTaskOrderStatusHandlerFunc{
		handlerConfig,
		moveTaskOrderUpdater,
	}

	ghcAPI.MoveTaskOrderUpdateMTOStatusServiceCounselingCompletedHandler = UpdateMTOStatusServiceCounselingCompletedHandlerFunc{
		handlerConfig,
		moveTaskOrderUpdater,
	}

	ghcAPI.MoveTaskOrderUpdateMTOReviewedBillableWeightsAtHandler = UpdateMTOReviewedBillableWeightsAtHandlerFunc{
		handlerConfig,
		moveTaskOrderUpdater,
	}

	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(
		queryBuilder,
		fetch.NewFetcher(queryBuilder),
		moveRouter,
	)
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(ppmEstimator)
	shipmentCreator := shipment.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator)
	ghcAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		handlerConfig,
		shipmentCreator,
		shipmentSITStatus,
	}

	ghcAPI.MtoShipmentListMTOShipmentsHandler = ListMTOShipmentsHandler{
		handlerConfig,
		mtoshipment.NewMTOShipmentFetcher(),
		shipmentSITStatus,
	}

	ghcAPI.MtoShipmentGetShipmentHandler = GetMTOShipmentHandler{
		HandlerConfig:      handlerConfig,
		mtoShipmentFetcher: mtoshipment.NewMTOShipmentFetcher(),
	}

	ghcAPI.ShipmentDeleteShipmentHandler = DeleteShipmentHandler{
		handlerConfig,
		mtoshipment.NewShipmentDeleter(),
	}

	ghcAPI.ShipmentApproveShipmentHandler = ApproveShipmentHandler{
		handlerConfig,
		mtoshipment.NewShipmentApprover(
			mtoshipment.NewShipmentRouter(),
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
			handlerConfig.Planner(),
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentRequestShipmentDiversionHandler = RequestShipmentDiversionHandler{
		handlerConfig,
		mtoshipment.NewShipmentDiversionRequester(
			mtoshipment.NewShipmentRouter(),
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentApproveShipmentDiversionHandler = ApproveShipmentDiversionHandler{
		handlerConfig,
		mtoshipment.NewShipmentDiversionApprover(
			mtoshipment.NewShipmentRouter(),
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentRejectShipmentHandler = RejectShipmentHandler{
		handlerConfig,
		mtoshipment.NewShipmentRejecter(
			mtoshipment.NewShipmentRouter(),
		),
	}

	ghcAPI.ShipmentRequestShipmentCancellationHandler = RequestShipmentCancellationHandler{
		handlerConfig,
		mtoshipment.NewShipmentCancellationRequester(
			mtoshipment.NewShipmentRouter(),
		),
		shipmentSITStatus,
	}

	paymentRequestRecalculator := paymentrequest.NewPaymentRequestRecalculator(
		paymentrequest.NewPaymentRequestCreator(
			handlerConfig.HHGPlanner(),
			ghcrateengine.NewServiceItemPricer(),
		),
		paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
	)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(paymentRequestRecalculator)

	ghcAPI.ShipmentRequestShipmentReweighHandler = RequestShipmentReweighHandler{
		handlerConfig,
		mtoshipment.NewShipmentReweighRequester(),
		shipmentSITStatus,
		mtoshipment.NewOfficeMTOShipmentUpdater(
			queryBuilder,
			fetch.NewFetcher(queryBuilder),
			handlerConfig.Planner(),
			moveRouter,
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
			handlerConfig.NotificationSender(),
			paymentRequestShipmentRecalculator,
		),
	}
	mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(
		queryBuilder,
		fetch.NewFetcher(queryBuilder),
		handlerConfig.Planner(),
		moveRouter,
		move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
		handlerConfig.NotificationSender(),
		paymentRequestShipmentRecalculator,
	)

	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(ppmEstimator)
	shipmentUpdater := shipment.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)

	ghcAPI.MoveSearchMovesHandler = SearchMovesHandler{
		HandlerConfig: handlerConfig,
		MoveSearcher:  move.NewMoveSearcher(),
	}

	ghcAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateShipmentHandler{
		handlerConfig,
		shipmentUpdater,
		shipmentSITStatus,
	}

	ghcAPI.MtoAgentFetchMTOAgentListHandler = ListMTOAgentsHandler{
		HandlerConfig: handlerConfig,
		ListFetcher:   fetch.NewListFetcher(queryBuilder),
	}

	ghcAPI.ShipmentApproveSITExtensionHandler = ApproveSITExtensionHandler{
		handlerConfig,
		mtoshipment.NewSITExtensionApprover(moveRouter),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentDenySITExtensionHandler = DenySITExtensionHandler{
		handlerConfig,
		mtoshipment.NewSITExtensionDenier(moveRouter),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentCreateSITExtensionAsTOOHandler = CreateSITExtensionAsTOOHandler{
		handlerConfig,
		mtoshipment.NewCreateSITExtensionAsTOO(),
		shipmentSITStatus,
	}

	ghcAPI.GhcDocumentsGetDocumentHandler = GetDocumentHandler{handlerConfig}

	ghcAPI.QueuesGetMovesQueueHandler = GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}

	ghcAPI.QueuesGetPaymentRequestsQueueHandler = GetPaymentRequestsQueueHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestListFetcher(),
	}

	ghcAPI.QueuesGetServicesCounselingQueueHandler = GetServicesCounselingQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}

	ghcAPI.TacTacValidationHandler = TacValidationHandler{
		handlerConfig,
	}

	ghcAPI.PaymentRequestsGetShipmentsPaymentSITBalanceHandler = ShipmentsSITBalanceHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestShipmentsSITBalance(),
	}

	ghcAPI.EvaluationReportsGetPWSViolationsHandler = GetPWSViolationsHandler{
		handlerConfig,
		pwsviolation.NewPWSViolationsFetcher(),
	}

	ghcAPI.ReportViolationsAssociateReportViolationsHandler = AssociateReportViolationsHandler{
		handlerConfig,
		reportviolation.NewReportViolationCreator(),
	}

	ghcAPI.ReportViolationsGetReportViolationsByReportIDHandler = GetReportViolationsHandler{
		handlerConfig,
		reportviolation.NewReportViolationFetcher(),
	}

	return ghcAPI
}
