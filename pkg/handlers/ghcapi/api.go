package ghcapi

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations"
	"github.com/transcom/mymove/pkg/handlers"
	paperwork "github.com/transcom/mymove/pkg/paperwork"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services/address"
	boatshipment "github.com/transcom/mymove/pkg/services/boat_shipment"
	dateservice "github.com/transcom/mymove/pkg/services/calendar"
	customerserviceremarks "github.com/transcom/mymove/pkg/services/customer_support_remarks"
	"github.com/transcom/mymove/pkg/services/entitlements"
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
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/office_user/customer"
	"github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	order "github.com/transcom/mymove/pkg/services/order"
	paperwork_service "github.com/transcom/mymove/pkg/services/paperwork"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	paymentserviceitem "github.com/transcom/mymove/pkg/services/payment_service_item"
	portlocation "github.com/transcom/mymove/pkg/services/port_location"
	ppmcloseout "github.com/transcom/mymove/pkg/services/ppm_closeout"
	ppmshipment "github.com/transcom/mymove/pkg/services/ppmshipment"
	progear "github.com/transcom/mymove/pkg/services/progear_weight_ticket"
	pwsviolation "github.com/transcom/mymove/pkg/services/pws_violation"
	"github.com/transcom/mymove/pkg/services/query"
	reportviolation "github.com/transcom/mymove/pkg/services/report_violation"
	"github.com/transcom/mymove/pkg/services/roles"
	serviceitem "github.com/transcom/mymove/pkg/services/service_item"
	shipmentaddressupdate "github.com/transcom/mymove/pkg/services/shipment_address_update"
	shipmentsummaryworksheet "github.com/transcom/mymove/pkg/services/shipment_summary_worksheet"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
	sitentrydateupdate "github.com/transcom/mymove/pkg/services/sit_entry_date_update"
	sitextension "github.com/transcom/mymove/pkg/services/sit_extension"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
	transportationaccountingcode "github.com/transcom/mymove/pkg/services/transportation_accounting_code"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	transportationofficeassignments "github.com/transcom/mymove/pkg/services/transportation_office_assignments"
	"github.com/transcom/mymove/pkg/services/upload"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
	weightticketparser "github.com/transcom/mymove/pkg/services/weight_ticket_parser"
	"github.com/transcom/mymove/pkg/uploader"
)

// NewGhcAPIHandler returns a handler for the GHC API
func NewGhcAPIHandler(handlerConfig handlers.HandlerConfig) *ghcops.MymoveAPI {
	waf := entitlements.NewWeightAllotmentFetcher()
	ghcSpec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcAPI := ghcops.NewMymoveAPI(ghcSpec)
	queryBuilder := query.NewQueryBuilder()
	moveRouter := move.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher())
	moveLocker := movelocker.NewMoveLocker()
	addressCreator := address.NewAddressCreator()
	portLocationFetcher := portlocation.NewPortLocationFetcher()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	officerUserCreator := officeuser.NewOfficeUserCreator(
		queryBuilder,
		handlerConfig.NotificationSender(),
	)
	officeUserUpdater := officeuser.NewOfficeUserUpdater(queryBuilder)
	newQueryFilter := query.NewQueryFilter
	newUserRolesCreator := usersroles.NewUsersRolesCreator()
	newRolesFetcher := roles.NewRolesFetcher()
	newTransportationOfficeAssignmentUpdater := transportationofficeassignments.NewTransportationOfficeAssignmentUpdater()
	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})

	mtoServiceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(
		handlerConfig.HHGPlanner(),
		queryBuilder,
		moveRouter,
		ghcrateengine.NewDomesticUnpackPricer(),
		ghcrateengine.NewDomesticPackPricer(),
		ghcrateengine.NewDomesticLinehaulPricer(),
		ghcrateengine.NewDomesticShorthaulPricer(),
		ghcrateengine.NewDomesticOriginPricer(),
		ghcrateengine.NewDomesticDestinationPricer(),
		ghcrateengine.NewFuelSurchargePricer())

	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoServiceItemCreator,
		moveRouter, signedCertificationCreator, signedCertificationUpdater, ppmEstimator,
	)

	ppmCloseoutFetcher := ppmcloseout.NewPPMCloseoutFetcher(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{}, ppmEstimator)
	SSWPPMComputer := shipmentsummaryworksheet.NewSSWPPMComputer(ppmCloseoutFetcher)
	uploadCreator := upload.NewUploadCreator(handlerConfig.FileStorer())

	serviceItemFetcher := serviceitem.NewServiceItemFetcher()

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

	weightTicketFetcher := weightticket.NewWeightTicketFetcher()
	parserComputer := weightticketparser.NewWeightTicketComputer()
	weightGenerator, err := weightticketparser.NewWeightTicketParserGenerator(pdfGenerator)
	if err != nil {
		log.Fatalln(err)
	}

	transportationOfficeFetcher := transportationoffice.NewTransportationOfficesFetcher()
	closeoutOfficeUpdater := move.NewCloseoutOfficeUpdater(move.NewMoveFetcher(), transportationOfficeFetcher)
	assignedOfficeUserUpdater := move.NewAssignedOfficeUserUpdater(move.NewMoveFetcher())
	vLocation := address.NewVLocation()

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
		OrderFetcher:            order.NewOrderFetcher(waf),
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

	ghcAPI.EvaluationReportsAddAppealToViolationHandler = AddAppealToViolationHandler{
		HandlerConfig:             handlerConfig,
		ReportViolationsAddAppeal: reportviolation.NewReportViolationsAddAppeal(),
	}

	ghcAPI.EvaluationReportsAddAppealToSeriousIncidentHandler = AddAppealToSeriousIncidentHandler{
		HandlerConfig:            handlerConfig,
		SeriousIncidentAddAppeal: evaluationreport.NewEvaluationReportSeriousIncidentAddAppeal(),
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
		move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()), waf),
		handlerConfig.NotificationSender(),
		paymentRequestShipmentRecalculator,
		addressUpdater,
		addressCreator)
	sitExtensionShipmentUpdater := shipment.NewShipmentUpdater(noCheckUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)

	ghcAPI.MtoServiceItemUpdateServiceItemSitEntryDateHandler = UpdateServiceItemSitEntryDateHandler{
		HandlerConfig:       handlerConfig,
		sitEntryDateUpdater: sitentrydateupdate.NewSitEntryDateUpdater(),
		ShipmentSITStatus:   sitstatus.NewShipmentSITStatus(),
		MTOShipmentFetcher:  mtoshipment.NewMTOShipmentFetcher(),
		ShipmentUpdater:     sitExtensionShipmentUpdater,
	}

	ghcAPI.MtoServiceItemUpdateMTOServiceItemStatusHandler = UpdateMTOServiceItemStatusHandler{
		HandlerConfig:         handlerConfig,
		MTOServiceItemUpdater: mtoserviceitem.NewMTOServiceItemUpdater(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, shipmentFetcher, addressCreator, portLocationFetcher, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
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
		movetaskorder.NewMoveTaskOrderFetcher(waf),
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
		order.NewOrderFetcher(waf),
	}
	ghcAPI.OrderCounselingUpdateOrderHandler = CounselingUpdateOrderHandler{
		handlerConfig,
		order.NewOrderUpdater(moveRouter),
	}
	ghcAPI.OrderCreateOrderHandler = CreateOrderHandler{
		handlerConfig,
		waf,
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

	ghcAPI.OrderAcknowledgeExcessUnaccompaniedBaggageWeightRiskHandler = AcknowledgeExcessUnaccompaniedBaggageWeightRiskHandler{
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
	shipmentCreator := shipment.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator, boatShipmentCreator, mobileHomeShipmentCreator, shipmentRouter, moveTaskOrderUpdater, move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()), waf))
	ghcAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		handlerConfig,
		shipmentCreator,
		shipmentSITStatus,
		closeoutOfficeUpdater,
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

	ghcAPI.ShipmentCreateTerminationHandler = TerminateShipmentHandler{
		handlerConfig,
		mtoshipment.NewShipmentTermination(),
	}

	ghcAPI.ShipmentApproveShipmentsHandler = ApproveShipmentsHandler{
		handlerConfig,
		mtoshipment.NewShipmentApprover(
			shipmentRouter,
			mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			handlerConfig.HHGPlanner(),
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()), waf),
			moveTaskOrderUpdater,
			moveRouter,
		),
		shipmentSITStatus,
		moveTaskOrderUpdater,
		move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()), waf),
		mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()),
	}

	ghcAPI.ShipmentApproveShipmentHandler = ApproveShipmentHandler{
		handlerConfig,
		mtoshipment.NewShipmentApprover(
			shipmentRouter,
			mtoServiceItemCreator,
			handlerConfig.HHGPlanner(),
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()), waf),
			moveTaskOrderUpdater,
			moveRouter,
		),
		shipmentSITStatus,
		moveTaskOrderUpdater,
		move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()), waf),
		mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()),
	}

	ghcAPI.ShipmentRequestShipmentDiversionHandler = RequestShipmentDiversionHandler{
		handlerConfig,
		mtoshipment.NewShipmentDiversionRequester(
			shipmentRouter,
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentApproveShipmentDiversionHandler = ApproveShipmentDiversionHandler{
		handlerConfig,
		mtoshipment.NewShipmentDiversionApprover(
			shipmentRouter,
			moveRouter,
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentRejectShipmentHandler = RejectShipmentHandler{
		handlerConfig,
		mtoshipment.NewShipmentRejecter(
			shipmentRouter,
		),
	}

	ghcAPI.ShipmentRequestShipmentCancellationHandler = RequestShipmentCancellationHandler{
		handlerConfig,
		mtoshipment.NewShipmentCancellationRequester(
			shipmentRouter,
			moveRouter,
		),
		shipmentSITStatus,
	}

	ghcAPI.ShipmentRequestShipmentReweighHandler = RequestShipmentReweighHandler{
		handlerConfig,
		mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()),
		shipmentSITStatus,
		mtoshipment.NewOfficeMTOShipmentUpdater(
			queryBuilder,
			fetch.NewFetcher(queryBuilder),
			handlerConfig.HHGPlanner(),
			moveRouter,
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()), waf),
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
		move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester(handlerConfig.NotificationSender()), waf),
		handlerConfig.NotificationSender(),
		paymentRequestShipmentRecalculator,
		addressUpdater,
		addressCreator,
	)

	shipmentUpdater := shipment.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)

	ghcAPI.MoveSearchMovesHandler = SearchMovesHandler{
		HandlerConfig: handlerConfig,
		MoveSearcher:  move.NewMoveSearcher(),
		MoveUnlocker:  movelocker.NewMoveUnlocker(),
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
		mtoserviceitem.NewMTOServiceItemUpdater(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, shipmentFetcher, addressCreator, portLocationFetcher, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
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

	ghcAPI.QueuesGetBulkAssignmentDataHandler = GetBulkAssignmentDataHandler{
		handlerConfig,
		officeuser.NewOfficeUserFetcherPop(),
		move.NewMoveFetcherBulkAssignment(),
		moveLocker,
	}

	ghcAPI.QueuesSaveBulkAssignmentDataHandler = SaveBulkAssignmentDataHandler{
		handlerConfig,
		officeuser.NewOfficeUserFetcherPop(),
		move.NewMoveFetcher(),
		move.NewMoveAssignerBulkAssignment(),
		movelocker.NewMoveUnlocker(),
	}

	ghcAPI.QueuesGetMovesQueueHandler = GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(waf),
		movelocker.NewMoveUnlocker(),
		officeuser.NewOfficeUserFetcherPop(),
	}

	ghcAPI.QueuesGetDestinationRequestsQueueHandler = GetDestinationRequestsQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(waf),
		movelocker.NewMoveUnlocker(),
		officeuser.NewOfficeUserFetcherPop(),
	}

	ghcAPI.QueuesListPrimeMovesHandler = ListPrimeMovesHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(waf),
	}

	ghcAPI.QueuesGetPaymentRequestsQueueHandler = GetPaymentRequestsQueueHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestListFetcher(),
		movelocker.NewMoveUnlocker(),
		officeuser.NewOfficeUserFetcherPop(),
	}

	ghcAPI.QueuesGetServicesCounselingQueueHandler = GetServicesCounselingQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(waf),
		movelocker.NewMoveUnlocker(),
		officeuser.NewOfficeUserFetcherPop(),
	}

	ghcAPI.QueuesGetServicesCounselingOriginListHandler = GetServicesCounselingOriginListHandler{
		handlerConfig,
		order.NewOrderFetcher(waf),
		officeuser.NewOfficeUserFetcherPop(),
	}

	ghcAPI.TacTacValidationHandler = TacValidationHandler{
		handlerConfig,
	}

	ghcAPI.PaymentRequestsGetShipmentsPaymentSITBalanceHandler = ShipmentsSITBalanceHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestShipmentsSITBalance(),
	}

	ghcAPI.PpmCreateProGearWeightTicketHandler = CreateProGearWeightTicketHandler{handlerConfig, progear.NewOfficeProgearWeightTicketCreator()}
	ghcAPI.PpmDeleteProGearWeightTicketHandler = DeleteProGearWeightTicketHandler{handlerConfig, progear.NewProgearWeightTicketDeleter()}

	ghcAPI.PpmUpdateProGearWeightTicketHandler = UpdateProgearWeightTicketHandler{
		handlerConfig,
		progear.NewOfficeProgearWeightTicketUpdater(),
	}

	ppmShipmentFetcher := ppmshipment.NewPPMShipmentFetcher()
	if err != nil {
		log.Fatalln(err)
	}

	ppmShipmentRouter := ppmshipment.NewPPMShipmentRouter(shipmentRouter)
	if err != nil {
		log.Fatalln(err)
	}

	ghcAPI.PpmSendPPMToCustomerHandler = SendPPMToCustomerHandler{
		handlerConfig,
		ppmShipmentFetcher,
		moveTaskOrderUpdater,
	}

	ghcAPI.PpmFinishDocumentReviewHandler = FinishDocumentReviewHandler{
		handlerConfig,
		ppmshipment.NewPPMShipmentReviewDocuments(
			ppmShipmentRouter,
			signedCertificationCreator, signedCertificationUpdater, SSWPPMComputer,
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

	ghcAPI.PpmCreatePPMUploadHandler = CreatePPMUploadHandler{handlerConfig, weightGenerator, parserComputer, userUploader}
	ghcAPI.PpmCreateWeightTicketHandler = CreateWeightTicketHandler{handlerConfig, weightticket.NewOfficeWeightTicketCreator()}
	ghcAPI.PpmDeleteWeightTicketHandler = DeleteWeightTicketHandler{handlerConfig, weightticket.NewWeightTicketDeleter(weightTicketFetcher, ppmEstimator)}

	ghcAPI.PpmUpdateWeightTicketHandler = UpdateWeightTicketHandler{
		handlerConfig,
		weightticket.NewOfficeWeightTicketUpdater(weightTicketFetcher, ppmShipmentUpdater),
	}

	ghcAPI.PpmCreateMovingExpenseHandler = CreateMovingExpenseHandler{
		handlerConfig,
		movingexpense.NewMovingExpenseCreator(),
	}

	ghcAPI.PpmUpdateMovingExpenseHandler = UpdateMovingExpenseHandler{
		handlerConfig,
		movingexpense.NewOfficeMovingExpenseUpdater(ppmEstimator),
	}

	ghcAPI.PpmDeleteMovingExpenseHandler = DeleteMovingExpenseHandler{
		handlerConfig,
		movingexpense.NewMovingExpenseDeleter(),
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

	ghcAPI.TransportationOfficeShowCounselingOfficesHandler = ShowCounselingOfficesHandler{
		handlerConfig,
		transportationOfficeFetcher,
	}

	ghcAPI.MoveUpdateCloseoutOfficeHandler = UpdateMoveCloseoutOfficeHandler{
		handlerConfig,
		closeoutOfficeUpdater,
	}

	ghcAPI.AddressesGetLocationByZipCityStateHandler = GetLocationByZipCityStateHandler{
		handlerConfig,
		vLocation,
	}

	ghcAPI.OfficeUsersCreateRequestedOfficeUserHandler = RequestOfficeUserHandler{
		handlerConfig,
		officerUserCreator,
		newQueryFilter,
		newUserRolesCreator,
		newRolesFetcher,
		newTransportationOfficeAssignmentUpdater,
	}
	ghcAPI.OfficeUsersUpdateOfficeUserHandler = UpdateOfficeUserHandler{
		handlerConfig,
		officeUserUpdater,
	}
	paymentPacketCreator := ppmshipment.NewPaymentPacketCreator(ppmShipmentFetcher, pdfGenerator, AOAPacketCreator)
	ghcAPI.PpmShowPaymentPacketHandler = ShowPaymentPacketHandler{handlerConfig, paymentPacketCreator}

	ghcAPI.UploadsCreateUploadHandler = CreateUploadHandler{handlerConfig}
	ghcAPI.UploadsUpdateUploadHandler = UpdateUploadHandler{handlerConfig, upload.NewUploadInformationFetcher()}
	ghcAPI.UploadsDeleteUploadHandler = DeleteUploadHandler{handlerConfig, upload.NewUploadInformationFetcher()}
	ghcAPI.UploadsGetUploadStatusHandler = GetUploadStatusHandler{handlerConfig, upload.NewUploadInformationFetcher()}
	ghcAPI.TextEventStreamProducer = runtime.ByteStreamProducer() // GetUploadStatus produces Event Stream

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
		officeuser.NewOfficeUserFetcherPop(),
	}
	ghcAPI.MoveDeleteAssignedOfficeUserHandler = DeleteAssignedOfficeUserHandler{
		handlerConfig,
		assignedOfficeUserUpdater,
	}

	ghcAPI.MoveCheckForLockedMovesAndUnlockHandler = CheckForLockedMovesAndUnlockHandler{
		HandlerConfig: handlerConfig,
		MoveUnlocker:  movelocker.NewMoveUnlocker(),
	}

	ghcAPI.ReServiceItemsGetAllReServiceItemsHandler = GetReServiceItemsHandler{
		handlerConfig,
		serviceItemFetcher,
	}

	ppmShipmentNewSubmitter := ppmshipment.NewPPMShipmentNewSubmitter(ppmShipmentFetcher, signedCertificationCreator, ppmShipmentRouter)
	ghcAPI.PpmSubmitPPMShipmentDocumentationHandler = SubmitPPMShipmentDocumentationHandler{handlerConfig, ppmShipmentNewSubmitter}

	return ghcAPI
}
