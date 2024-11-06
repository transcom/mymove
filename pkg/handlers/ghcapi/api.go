package ghcapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/handlers"
	paperwork "github.com/transcom/mymove/pkg/paperwork"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services/address"
	boatshipment "github.com/transcom/mymove/pkg/services/boat_shipment"
	dateservice "github.com/transcom/mymove/pkg/services/calendar"
	customerserviceremarks "github.com/transcom/mymove/pkg/services/customer_support_remarks"
	evaluationreport "github.com/transcom/mymove/pkg/services/evaluation_report"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	lineofaccounting "github.com/transcom/mymove/pkg/services/line_of_accounting"
	movelocker "github.com/transcom/mymove/pkg/services/lock_move"
	mobileHomeShipment "github.com/transcom/mymove/pkg/services/mobile_home_shipment"
	"github.com/transcom/mymove/pkg/services/move"
	movehistory "github.com/transcom/mymove/pkg/services/move_history"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	movingexpense "github.com/transcom/mymove/pkg/services/moving_expense"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	officeusercreator "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/office_user/customer"
	"github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	order "github.com/transcom/mymove/pkg/services/order"
	paperwork_service "github.com/transcom/mymove/pkg/services/paperwork"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	paymentserviceitem "github.com/transcom/mymove/pkg/services/payment_service_item"
	ppmcloseout "github.com/transcom/mymove/pkg/services/ppm_closeout"
	ppmshipment "github.com/transcom/mymove/pkg/services/ppmshipment"
	progear "github.com/transcom/mymove/pkg/services/progear_weight_ticket"
	pwsviolation "github.com/transcom/mymove/pkg/services/pws_violation"
	"github.com/transcom/mymove/pkg/services/query"
	reportviolation "github.com/transcom/mymove/pkg/services/report_violation"
	"github.com/transcom/mymove/pkg/services/roles"
	shipmentaddressupdate "github.com/transcom/mymove/pkg/services/shipment_address_update"
	shipmentsummaryworksheet "github.com/transcom/mymove/pkg/services/shipment_summary_worksheet"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
	sitentrydateupdate "github.com/transcom/mymove/pkg/services/sit_entry_date_update"
	sitextension "github.com/transcom/mymove/pkg/services/sit_extension"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
	transportationaccountingcode "github.com/transcom/mymove/pkg/services/transportation_accounting_code"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	transportaionofficeassignments "github.com/transcom/mymove/pkg/services/transportation_office_assignments"
	"github.com/transcom/mymove/pkg/services/upload"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
	"github.com/transcom/mymove/pkg/uploader"
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
	moveLocker := movelocker.NewMoveLocker()
	addressCreator := address.NewAddressCreator()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	officerUserCreator := officeusercreator.NewOfficeUserCreator(
		queryBuilder,
		handlerConfig.NotificationSender(),
	)
	newQueryFilter := query.NewQueryFilter
	newUserRolesCreator := usersroles.NewUsersRolesCreator()
	newRolesFetcher := roles.NewRolesFetcher()
	newTransportaionOfficeAssignmentUpdater := transportaionofficeassignments.NewTransportaionOfficeAssignmentUpdater()
	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
		moveRouter, signedCertificationCreator, signedCertificationUpdater,
	)

	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	ppmCloseoutFetcher := ppmcloseout.NewPPMCloseoutFetcher(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{}, ppmEstimator)
	SSWPPMComputer := shipmentsummaryworksheet.NewSSWPPMComputer(ppmCloseoutFetcher)
	uploadCreator := upload.NewUploadCreator(handlerConfig.FileStorer())

	userUploader, err := uploader.NewUserUploader(handlerConfig.FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		log.Fatalln(err)
	}

	pdfGenerator, err := paperwork.NewGenerator(userUploader.Uploader())
	if err != nil {
		log.Fatalln(err)
	}

	SSWPPMGenerator, err := shipmentsummaryworksheet.NewSSWPPMGenerator(pdfGenerator)
	if err != nil {
		log.Fatalln(err)
	}

	transportationOfficeFetcher := transportationoffice.NewTransportationOfficesFetcher()
	closeoutOfficeUpdater := move.NewCloseoutOfficeUpdater(move.NewMoveFetcher(), transportationOfficeFetcher)
	assignedOfficeUserUpdater := move.NewAssignedOfficeUserUpdater(move.NewMoveFetcher())

	shipmentSITStatus := sitstatus.NewShipmentSITStatus()

	ghcAPI.ServeError = handlers.ServeCustomError

	ghcAPI.MoveGetMoveHandler = GetMoveHandler{
		HandlerConfig: handlerConfig,
		MoveFetcher:   move.NewMoveFetcher(),
		MoveLocker:    moveLocker,
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

	ghcAPI.EvaluationReportsDownloadEvaluationReportHandler = DownloadEvaluationReportHandler{
		HandlerConfig:           handlerConfig,
		EvaluationReportFetcher: evaluationreport.NewEvaluationReportFetcher(),
		MTOShipmentFetcher:      mtoshipment.NewMTOShipmentFetcher(),
		OrderFetcher:            order.NewOrderFetcher(),
		ReportViolationFetcher:  reportviolation.NewReportViolationFetcher(),
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

	ghcAPI.MtoServiceItemGetMTOServiceItemHandler = GetMTOServiceItemHandler{
		HandlerConfig:         handlerConfig,
		mtoServiceItemFetcher: mtoserviceitem.NewMTOServiceItemFetcher(),
	}

	ghcAPI.MoveUploadAdditionalDocumentsHandler = UploadAdditionalDocumentsHandler{
		HandlerConfig: handlerConfig,
		uploader:      move.NewMoveAdditionalDocumentsUploader(uploadCreator),
	}
	paymentRequestRecalculator := paymentrequest.NewPaymentRequestRecalculator(
		paymentrequest.NewPaymentRequestCreator(
			handlerConfig.HHGPlanner(),
			ghcrateengine.NewServiceItemPricer(),
		),
		paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
	)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(paymentRequestRecalculator)
	addressUpdater := address.NewAddressUpdater()
	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(ppmEstimator, addressCreator, addressUpdater)
	boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
	mobileHomeShipmentUpdater := mobileHomeShipment.NewMobileHomeShipmentUpdater()

	noCheckUpdater := mtoshipment.NewMTOShipmentUpdater(queryBuilder,
		fetch.NewFetcher(queryBuilder),
		handlerConfig.HHGPlanner(),
		moveRouter,
		move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
		handlerConfig.NotificationSender(),
		paymentRequestShipmentRecalculator,
		addressUpdater,
		addressCreator)
	sitExtensionShipmentUpdater := shipment.NewShipmentUpdater(noCheckUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)

	ghcAPI.MtoServiceItemUpdateServiceItemSitEntryDateHandler = UpdateServiceItemSitEntryDateHandler{
		HandlerConfig:       handlerConfig,
		sitEntryDateUpdater: sitentrydateupdate.NewSitEntryDateUpdater(),
		ShipmentSITStatus:   sitstatus.NewShipmentSITStatus(),
		MTOShipmentFetcher:  mtoshipment.NewMTOShipmentFetcher(),
		ShipmentUpdater:     sitExtensionShipmentUpdater,
	}

	ghcAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{
		HandlerConfig:         handlerConfig,
		MTOServiceItemUpdater: mtoserviceitem.NewMTOServiceItemUpdater(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, shipmentFetcher, addressCreator),
		Fetcher:               fetch.NewFetcher(queryBuilder),
		ShipmentSITStatus:     sitstatus.NewShipmentSITStatus(),
		MTOShipmentFetcher:    mtoshipment.NewMTOShipmentFetcher(),
		ShipmentUpdater:       sitExtensionShipmentUpdater,
	}

	ghcAPI.MtoServiceItemListMTOServiceItemsHandler = ListMTOServiceItemsHandler{
		handlerConfig,
		fetch.NewListFetcher(queryBuilder),
		fetch.NewFetcher(queryBuilder),
		ghcrateengine.NewCounselingServicesPricer(),
		ghcrateengine.NewManagementServicesPricer(),
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
	ghcAPI.CustomerCreateCustomerWithOktaOptionHandler = CreateCustomerWithOktaOptionHandler{
		handlerConfig,
	}
	ghcAPI.OrderGetOrderHandler = GetOrdersHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}
	ghcAPI.OrderCounselingUpdateOrderHandler = CounselingUpdateOrderHandler{
		handlerConfig,
		order.NewOrderUpdater(moveRouter),
	}
	ghcAPI.OrderCreateOrderHandler = CreateOrderHandler{
		handlerConfig,
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

	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreatorV1(
		queryBuilder,
		fetch.NewFetcher(queryBuilder),
		moveRouter,
		addressCreator,
	)

	primeDownloadMoveUploadPDFGenerator, err := paperwork_service.NewMoveUserUploadToPDFDownloader(pdfGenerator)
	if err != nil {
		log.Fatalln(err)
	}

	AOAPacketCreator := ppmshipment.NewAOAPacketCreator(SSWPPMGenerator, SSWPPMComputer, primeDownloadMoveUploadPDFGenerator, userUploader, pdfGenerator)
	if err != nil {
		log.Fatalln(err)
	}

	ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(ppmEstimator, addressCreator)
	boatShipmentCreator := boatshipment.NewBoatShipmentCreator()
	mobileHomeShipmentCreator := mobileHomeShipment.NewMobileHomeShipmentCreator()
	ghcAPI.PpmShowAOAPacketHandler = showAOAPacketHandler{handlerConfig, SSWPPMComputer, SSWPPMGenerator, AOAPacketCreator}

	shipmentRouter := mtoshipment.NewShipmentRouter()
	shipmentCreator := shipment.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator, boatShipmentCreator, mobileHomeShipmentCreator, shipmentRouter, moveTaskOrderUpdater)
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
		mtoshipment.NewShipmentDeleter(moveTaskOrderUpdater, moveRouter),
	}

	ghcAPI.ShipmentApproveShipmentHandler = ApproveShipmentHandler{
		handlerConfig,
		mtoshipment.NewShipmentApprover(
			mtoshipment.NewShipmentRouter(),
			mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			handlerConfig.HHGPlanner(),
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
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
			moveRouter,
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
			moveRouter,
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentRequestShipmentReweighHandler = RequestShipmentReweighHandler{
		handlerConfig,
		mtoshipment.NewShipmentReweighRequester(),
		shipmentSITStatus,
		mtoshipment.NewOfficeMTOShipmentUpdater(
			queryBuilder,
			fetch.NewFetcher(queryBuilder),
			handlerConfig.HHGPlanner(),
			moveRouter,
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
			handlerConfig.NotificationSender(),
			paymentRequestShipmentRecalculator,
			addressUpdater,
			addressCreator,
		),
	}
	mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(
		queryBuilder,
		fetch.NewFetcher(queryBuilder),
		handlerConfig.HHGPlanner(),
		moveRouter,
		move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
		handlerConfig.NotificationSender(),
		paymentRequestShipmentRecalculator,
		addressUpdater,
		addressCreator,
	)

	shipmentUpdater := shipment.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)

	ghcAPI.MoveSearchMovesHandler = SearchMovesHandler{
		HandlerConfig: handlerConfig,
		MoveSearcher:  move.NewMoveSearcher(),
	}

	ghcAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateShipmentHandler{
		handlerConfig,
		shipmentUpdater,
		shipmentSITStatus,
	}

	shipmentAddressUpdater := shipmentaddressupdate.NewShipmentAddressUpdateRequester(handlerConfig.HHGPlanner(), addressCreator, moveRouter)

	ghcAPI.ShipmentReviewShipmentAddressUpdateHandler = ReviewShipmentAddressUpdateHandler{
		handlerConfig,
		shipmentAddressUpdater,
	}

	ghcAPI.MtoAgentFetchMTOAgentListHandler = ListMTOAgentsHandler{
		HandlerConfig: handlerConfig,
		ListFetcher:   fetch.NewListFetcher(queryBuilder),
	}

	ghcAPI.ShipmentApproveSITExtensionHandler = ApproveSITExtensionHandler{
		handlerConfig,
		sitextension.NewSITExtensionApprover(moveRouter),
		shipmentSITStatus,
		sitExtensionShipmentUpdater,
	}

	ghcAPI.ShipmentDenySITExtensionHandler = DenySITExtensionHandler{
		handlerConfig,
		sitextension.NewSITExtensionDenier(moveRouter),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentUpdateSITServiceItemCustomerExpenseHandler = UpdateSITServiceItemCustomerExpenseHandler{
		handlerConfig,
		mtoserviceitem.NewMTOServiceItemUpdater(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, shipmentFetcher, addressCreator),
		mtoshipment.NewMTOShipmentFetcher(),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentCreateApprovedSITDurationUpdateHandler = CreateApprovedSITDurationUpdateHandler{
		handlerConfig,
		sitextension.NewApprovedSITDurationUpdateCreator(),
		shipmentSITStatus,
		sitExtensionShipmentUpdater,
	}

	ghcAPI.GhcDocumentsGetDocumentHandler = GetDocumentHandler{handlerConfig}
	ghcAPI.GhcDocumentsCreateDocumentHandler = CreateDocumentHandler{handlerConfig}

	ghcAPI.QueuesGetMovesQueueHandler = GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		movelocker.NewMoveUnlocker(),
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	ghcAPI.QueuesListPrimeMovesHandler = ListPrimeMovesHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	ghcAPI.QueuesGetPaymentRequestsQueueHandler = GetPaymentRequestsQueueHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestListFetcher(),
		movelocker.NewMoveUnlocker(),
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	ghcAPI.QueuesGetServicesCounselingQueueHandler = GetServicesCounselingQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		movelocker.NewMoveUnlocker(),
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	ghcAPI.QueuesGetServicesCounselingOriginListHandler = GetServicesCounselingOriginListHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	ghcAPI.TacTacValidationHandler = TacValidationHandler{
		handlerConfig,
	}

	ghcAPI.PaymentRequestsGetShipmentsPaymentSITBalanceHandler = ShipmentsSITBalanceHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestShipmentsSITBalance(),
	}

	ghcAPI.PpmUpdateProGearWeightTicketHandler = UpdateProgearWeightTicketHandler{
		handlerConfig,
		progear.NewOfficeProgearWeightTicketUpdater(),
	}

	ghcAPI.PpmFinishDocumentReviewHandler = FinishDocumentReviewHandler{
		handlerConfig,
		ppmshipment.NewPPMShipmentReviewDocuments(
			ppmshipment.NewPPMShipmentRouter(mtoshipment.NewShipmentRouter()),
			signedCertificationCreator, signedCertificationUpdater,
		),
	}

	ppmDocumentsFetcher := ppmshipment.NewPPMDocumentFetcher()

	ghcAPI.PpmGetPPMDocumentsHandler = GetPPMDocumentsHandler{
		handlerConfig,
		ppmDocumentsFetcher,
	}

	ghcAPI.PpmGetPPMCloseoutHandler = GetPPMCloseoutHandler{
		handlerConfig,
		ppmCloseoutFetcher,
	}

	ghcAPI.PpmGetPPMActualWeightHandler = GetPPMActualWeightHandler{
		handlerConfig,
		ppmCloseoutFetcher,
		ppmshipment.NewPPMShipmentFetcher(),
	}

	weightTicketFetcher := weightticket.NewWeightTicketFetcher()

	ghcAPI.PpmUpdateWeightTicketHandler = UpdateWeightTicketHandler{
		handlerConfig,
		weightticket.NewOfficeWeightTicketUpdater(weightTicketFetcher, ppmShipmentUpdater),
	}

	ghcAPI.PpmUpdateMovingExpenseHandler = UpdateMovingExpenseHandler{
		handlerConfig,
		movingexpense.NewOfficeMovingExpenseUpdater(ppmEstimator),
	}

	ghcAPI.PwsViolationsGetPWSViolationsHandler = GetPWSViolationsHandler{
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

	ghcAPI.TransportationOfficeGetTransportationOfficesHandler = GetTransportationOfficesHandler{
		handlerConfig,
		transportationOfficeFetcher,
	}

	ghcAPI.TransportationOfficeGetTransportationOfficesOpenHandler = GetTransportationOfficesOpenHandler{
		handlerConfig,
		transportationOfficeFetcher,
	}

	ghcAPI.TransportationOfficeGetTransportationOfficesGBLOCsHandler = GetTransportationOfficesGBLOCsHandler{
		handlerConfig,
		transportationOfficeFetcher,
	}

	ghcAPI.MoveUpdateCloseoutOfficeHandler = UpdateMoveCloseoutOfficeHandler{
		handlerConfig,
		closeoutOfficeUpdater,
	}

	ghcAPI.OfficeUsersCreateRequestedOfficeUserHandler = RequestOfficeUserHandler{
		handlerConfig,
		officerUserCreator,
		newQueryFilter,
		newUserRolesCreator,
		newRolesFetcher,
		newTransportaionOfficeAssignmentUpdater,
	}
	ppmShipmentFetcher := ppmshipment.NewPPMShipmentFetcher()
	paymentPacketCreator := ppmshipment.NewPaymentPacketCreator(ppmShipmentFetcher, pdfGenerator, AOAPacketCreator)
	ghcAPI.PpmShowPaymentPacketHandler = ShowPaymentPacketHandler{handlerConfig, paymentPacketCreator}

	ghcAPI.UploadsCreateUploadHandler = CreateUploadHandler{handlerConfig}
	ghcAPI.UploadsUpdateUploadHandler = UpdateUploadHandler{handlerConfig, upload.NewUploadInformationFetcher()}
	ghcAPI.UploadsDeleteUploadHandler = DeleteUploadHandler{handlerConfig, upload.NewUploadInformationFetcher()}

	ghcAPI.CustomerSearchCustomersHandler = SearchCustomersHandler{
		HandlerConfig:    handlerConfig,
		CustomerSearcher: customer.NewCustomerSearcher(),
	}

	// Create TAC and LOA services
	tacFetcher := transportationaccountingcode.NewTransportationAccountingCodeFetcher()
	loaFetcher := lineofaccounting.NewLinesOfAccountingFetcher(tacFetcher)

	ghcAPI.LinesOfAccountingRequestLineOfAccountingHandler = LinesOfAccountingRequestLineOfAccountingHandler{
		HandlerConfig:           handlerConfig,
		LineOfAccountingFetcher: loaFetcher,
	}

	ghcAPI.ApplicationParametersGetParamHandler = ApplicationParametersParamHandler{handlerConfig}
	ghcAPI.PpmUpdatePPMSITHandler = UpdatePPMSITHandler{handlerConfig, ppmShipmentUpdater, ppmShipmentFetcher}
	ghcAPI.PpmGetPPMSITEstimatedCostHandler = GetPPMSITEstimatedCostHandler{handlerConfig, ppmEstimator, ppmShipmentFetcher}

	ghcAPI.OrderUploadAmendedOrdersHandler = UploadAmendedOrdersHandler{
		handlerConfig,
		order.NewOrderUpdater(moveRouter),
	}

	ghcAPI.MoveMoveCancelerHandler = MoveCancelerHandler{
		handlerConfig,
		move.NewMoveCanceler(),
	}

	paymentRequestBulkDownloadCreator := paymentrequest.NewPaymentRequestBulkDownloadCreator(pdfGenerator)
	ghcAPI.PaymentRequestsBulkDownloadHandler = PaymentRequestBulkDownloadHandler{
		handlerConfig,
		paymentRequestBulkDownloadCreator,
	}

	dateSelectionChecker := dateservice.NewDateSelectionChecker()
	ghcAPI.CalendarIsDateWeekendHolidayHandler = IsDateWeekendHolidayHandler{handlerConfig, dateSelectionChecker}

	ghcAPI.MoveUpdateAssignedOfficeUserHandler = UpdateAssignedOfficeUserHandler{
		handlerConfig,
		assignedOfficeUserUpdater,
		officeusercreator.NewOfficeUserFetcherPop(),
	}
	ghcAPI.MoveDeleteAssignedOfficeUserHandler = DeleteAssignedOfficeUserHandler{
		handlerConfig,
		assignedOfficeUserUpdater,
	}

	return ghcAPI
}
