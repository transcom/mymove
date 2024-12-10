package internalapi

import (
	"io"
	"log"

	"github.com/benbjohnson/clock"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/gen/internalapi"
	internalops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	"github.com/transcom/mymove/pkg/handlers"
	paperworkgenerator "github.com/transcom/mymove/pkg/paperwork"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	boatshipment "github.com/transcom/mymove/pkg/services/boat_shipment"
	dateservice "github.com/transcom/mymove/pkg/services/calendar"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mobilehomeshipment "github.com/transcom/mymove/pkg/services/mobile_home_shipment"
	move "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	movingexpense "github.com/transcom/mymove/pkg/services/moving_expense"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	"github.com/transcom/mymove/pkg/services/order"
	"github.com/transcom/mymove/pkg/services/paperwork"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	postalcodeservice "github.com/transcom/mymove/pkg/services/postal_codes"
	ppmcloseout "github.com/transcom/mymove/pkg/services/ppm_closeout"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	progear "github.com/transcom/mymove/pkg/services/progear_weight_ticket"
	"github.com/transcom/mymove/pkg/services/query"
	shipmentsummaryworksheet "github.com/transcom/mymove/pkg/services/shipment_summary_worksheet"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	"github.com/transcom/mymove/pkg/services/upload"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
	weightticketparser "github.com/transcom/mymove/pkg/services/weight_ticket_parser"
	"github.com/transcom/mymove/pkg/uploader"
)

// NewInternalAPI returns the internal API
func NewInternalAPI(handlerConfig handlers.HandlerConfig) *internalops.MymoveAPI {

	internalSpec, err := loads.Analyzed(internalapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	internalAPI := internalops.NewMymoveAPI(internalSpec)

	internalAPI.ServeError = handlers.ServeCustomError
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := move.NewMoveRouter()
	uploadCreator := upload.NewUploadCreator(handlerConfig.FileStorer())
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	ppmCloseoutFetcher := ppmcloseout.NewPPMCloseoutFetcher(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{}, ppmEstimator)
	SSWPPMComputer := shipmentsummaryworksheet.NewSSWPPMComputer(ppmCloseoutFetcher)

	userUploader, err := uploader.NewUserUploader(handlerConfig.FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		log.Fatalln(err)
	}

	pdfGenerator, err := paperworkgenerator.NewGenerator(userUploader.Uploader())
	if err != nil {
		log.Fatalln(err)
	}

	SSWPPMGenerator, err := shipmentsummaryworksheet.NewSSWPPMGenerator(pdfGenerator)
	if err != nil {
		log.Fatalln(err)
	}

	parserComputer := weightticketparser.NewWeightTicketComputer()
	weightGenerator, err := weightticketparser.NewWeightTicketParserGenerator(pdfGenerator)

	if err != nil {
		log.Fatalln(err)
	}

	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	mtoShipmentRouter := mtoshipment.NewShipmentRouter()
	ppmShipmentRouter := ppmshipment.NewPPMShipmentRouter(mtoShipmentRouter)
	transportationOfficeFetcher := transportationoffice.NewTransportationOfficesFetcher()
	closeoutOfficeUpdater := move.NewCloseoutOfficeUpdater(move.NewMoveFetcher(), transportationOfficeFetcher)
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()

	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(ppmEstimator, addressCreator, addressUpdater)
	boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
	mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

	primeDownloadMoveUploadPDFGenerator, err := paperwork.NewMoveUserUploadToPDFDownloader(pdfGenerator)
	if err != nil {
		log.Fatalln(err)
	}
	ppmShipmentFetcher := ppmshipment.NewPPMShipmentFetcher()

	AOAPacketCreator := ppmshipment.NewAOAPacketCreator(SSWPPMGenerator, SSWPPMComputer, primeDownloadMoveUploadPDFGenerator, userUploader, pdfGenerator)
	if err != nil {
		log.Fatalln(err)
	}
	internalAPI.FeatureFlagsBooleanFeatureFlagForUserHandler = BooleanFeatureFlagsForUserHandler{handlerConfig}
	internalAPI.FeatureFlagsVariantFeatureFlagForUserHandler = VariantFeatureFlagsForUserHandler{handlerConfig}

	internalAPI.UsersShowLoggedInUserHandler = ShowLoggedInUserHandler{handlerConfig, officeuser.NewOfficeUserFetcherPop()}
	internalAPI.CertificationCreateSignedCertificationHandler = CreateSignedCertificationHandler{handlerConfig}
	internalAPI.CertificationIndexSignedCertificationHandler = IndexSignedCertificationsHandler{handlerConfig}

	internalAPI.DutyLocationsSearchDutyLocationsHandler = SearchDutyLocationsHandler{handlerConfig}
	internalAPI.TransportationOfficesShowCounselingOfficesHandler = ShowCounselingOfficesHandler{
		handlerConfig,
		transportationOfficeFetcher,
	}

	internalAPI.AddressesShowAddressHandler = ShowAddressHandler{handlerConfig}

	internalAPI.TransportationOfficesShowDutyLocationTransportationOfficeHandler = ShowDutyLocationTransportationOfficeHandler{handlerConfig}

	internalAPI.OrdersCreateOrdersHandler = CreateOrdersHandler{handlerConfig}
	internalAPI.OrdersUpdateOrdersHandler = UpdateOrdersHandler{handlerConfig}
	internalAPI.OrdersShowOrdersHandler = ShowOrdersHandler{handlerConfig}
	internalAPI.OrdersUploadAmendedOrdersHandler = UploadAmendedOrdersHandler{
		handlerConfig,
		order.NewOrderUpdater(moveRouter),
	}

	internalAPI.MovesPatchMoveHandler = PatchMoveHandler{handlerConfig, closeoutOfficeUpdater}
	internalAPI.MovesGetAllMovesHandler = GetAllMovesHandler{handlerConfig}

	internalAPI.ApplicationParametersValidateHandler = ApplicationParametersValidateHandler{handlerConfig}

	internalAPI.MovesShowMoveHandler = ShowMoveHandler{handlerConfig}
	internalAPI.MovesSubmitMoveForApprovalHandler = SubmitMoveHandler{
		handlerConfig,
		moveRouter,
	}
	internalAPI.MovesSubmitAmendedOrdersHandler = SubmitAmendedOrdersHandler{
		handlerConfig,
		moveRouter,
	}
	internalAPI.MovesUploadAdditionalDocumentsHandler = UploadAdditionalDocumentsHandler{
		handlerConfig,
		move.NewMoveAdditionalDocumentsUploader(uploadCreator),
	}

	internalAPI.OktaProfileShowOktaInfoHandler = GetOktaProfileHandler{handlerConfig}
	internalAPI.OktaProfileUpdateOktaInfoHandler = UpdateOktaProfileHandler{handlerConfig}

	internalAPI.ServiceMembersCreateServiceMemberHandler = CreateServiceMemberHandler{handlerConfig}
	internalAPI.ServiceMembersPatchServiceMemberHandler = PatchServiceMemberHandler{handlerConfig}
	internalAPI.ServiceMembersShowServiceMemberHandler = ShowServiceMemberHandler{handlerConfig}
	internalAPI.ServiceMembersShowServiceMemberOrdersHandler = ShowServiceMemberOrdersHandler{handlerConfig}

	internalAPI.BackupContactsIndexServiceMemberBackupContactsHandler = IndexBackupContactsHandler{handlerConfig}
	internalAPI.BackupContactsCreateServiceMemberBackupContactHandler = CreateBackupContactHandler{handlerConfig}
	internalAPI.BackupContactsUpdateServiceMemberBackupContactHandler = UpdateBackupContactHandler{handlerConfig}
	internalAPI.BackupContactsShowServiceMemberBackupContactHandler = ShowBackupContactHandler{handlerConfig}

	internalAPI.DocumentsCreateDocumentHandler = CreateDocumentHandler{handlerConfig}
	internalAPI.DocumentsShowDocumentHandler = ShowDocumentHandler{handlerConfig}
	internalAPI.UploadsCreateUploadHandler = CreateUploadHandler{handlerConfig}
	internalAPI.UploadsDeleteUploadHandler = DeleteUploadHandler{handlerConfig, upload.NewUploadInformationFetcher()}
	internalAPI.UploadsDeleteUploadsHandler = DeleteUploadsHandler{handlerConfig}

	internalAPI.QueuesShowQueueHandler = ShowQueueHandler{handlerConfig}
	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler{handlerConfig, moveRouter}
	internalAPI.OfficeApproveReimbursementHandler = ApproveReimbursementHandler{handlerConfig}
	internalAPI.OfficeCancelMoveHandler = CancelMoveHandler{handlerConfig, moveRouter}

	internalAPI.EntitlementsIndexEntitlementsHandler = IndexEntitlementsHandler{handlerConfig}

	internalAPI.CalendarShowAvailableMoveDatesHandler = ShowAvailableMoveDatesHandler{handlerConfig}

	internalAPI.PpmShowAOAPacketHandler = showAOAPacketHandler{handlerConfig, SSWPPMComputer, SSWPPMGenerator, AOAPacketCreator}

	internalAPI.RegisterProducer(uploader.FileTypePDF, PDFProducer())

	internalAPI.PostalCodesValidatePostalCodeWithRateDataHandler = ValidatePostalCodeWithRateDataHandler{
		handlerConfig,
		postalcodeservice.NewPostalCodeValidator(clock.New()),
	}

	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveRouter, addressCreator)
	shipmentRouter := mtoshipment.NewShipmentRouter()
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		builder,
		mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
		moveRouter, signedCertificationCreator, signedCertificationUpdater, ppmEstimator,
	)
	boatShipmentCreator := boatshipment.NewBoatShipmentCreator()
	mobileHomeShipmentCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
	shipmentCreator := shipment.NewShipmentCreator(mtoShipmentCreator, ppmshipment.NewPPMShipmentCreator(ppmEstimator, addressCreator), boatShipmentCreator, mobileHomeShipmentCreator, shipmentRouter, moveTaskOrderUpdater)

	internalAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		handlerConfig,
		shipmentCreator,
	}

	paymentRequestRecalculator := paymentrequest.NewPaymentRequestRecalculator(
		paymentrequest.NewPaymentRequestCreator(
			handlerConfig.HHGPlanner(),
			ghcrateengine.NewServiceItemPricer(),
		),
		paymentrequest.NewPaymentRequestStatusUpdater(builder),
	)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(paymentRequestRecalculator)

	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)

	shipmentUpdater := shipment.NewShipmentUpdater(
		mtoshipment.NewCustomerMTOShipmentUpdater(
			builder,
			fetcher,
			handlerConfig.DTODPlanner(),
			moveRouter,
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
			handlerConfig.NotificationSender(),
			paymentRequestShipmentRecalculator,
			addressUpdater,
			addressCreator,
		),
		ppmShipmentUpdater,
		boatShipmentUpdater,
		mobileHomeShipmentUpdater,
	)

	internalAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		handlerConfig,
		shipmentUpdater,
	}

	internalAPI.MtoShipmentListMTOShipmentsHandler = ListMTOShipmentsHandler{
		handlerConfig,
		mtoshipment.NewMTOShipmentFetcher(),
	}

	internalAPI.MtoShipmentDeleteShipmentHandler = DeleteShipmentHandler{
		handlerConfig,
		mtoshipment.NewShipmentDeleter(moveTaskOrderUpdater, moveRouter),
	}

	dateSelectionChecker := dateservice.NewDateSelectionChecker()
	internalAPI.CalendarIsDateWeekendHolidayHandler = IsDateWeekendHolidayHandler{handlerConfig, dateSelectionChecker}

	internalAPI.PpmCreateMovingExpenseHandler = CreateMovingExpenseHandler{handlerConfig, movingexpense.NewMovingExpenseCreator()}
	internalAPI.PpmUpdateMovingExpenseHandler = UpdateMovingExpenseHandler{handlerConfig, movingexpense.NewCustomerMovingExpenseUpdater(ppmEstimator)}
	internalAPI.PpmDeleteMovingExpenseHandler = DeleteMovingExpenseHandler{handlerConfig, movingexpense.NewMovingExpenseDeleter()}

	internalAPI.PpmCreateWeightTicketHandler = CreateWeightTicketHandler{handlerConfig, weightticket.NewCustomerWeightTicketCreator()}

	weightTicketFetcher := weightticket.NewWeightTicketFetcher()
	internalAPI.PpmUpdateWeightTicketHandler = UpdateWeightTicketHandler{handlerConfig, weightticket.NewCustomerWeightTicketUpdater(weightTicketFetcher, ppmShipmentUpdater)}
	internalAPI.PpmDeleteWeightTicketHandler = DeleteWeightTicketHandler{handlerConfig, weightticket.NewWeightTicketDeleter(weightTicketFetcher, ppmEstimator)}

	internalAPI.PpmCreateProGearWeightTicketHandler = CreateProGearWeightTicketHandler{handlerConfig, progear.NewCustomerProgearWeightTicketCreator()}
	internalAPI.PpmUpdateProGearWeightTicketHandler = UpdateProGearWeightTicketHandler{handlerConfig, progear.NewCustomerProgearWeightTicketUpdater()}
	internalAPI.PpmDeleteProGearWeightTicketHandler = DeleteProGearWeightTicketHandler{handlerConfig, progear.NewProgearWeightTicketDeleter()}

	internalAPI.PpmCreatePPMUploadHandler = CreatePPMUploadHandler{handlerConfig, weightGenerator, parserComputer, userUploader}

	ppmShipmentNewSubmitter := ppmshipment.NewPPMShipmentNewSubmitter(ppmShipmentFetcher, signedCertificationCreator, ppmShipmentRouter)

	internalAPI.PpmSubmitPPMShipmentDocumentationHandler = SubmitPPMShipmentDocumentationHandler{handlerConfig, ppmShipmentNewSubmitter}

	ppmShipmentUpdatedSubmitter := ppmshipment.NewPPMShipmentUpdatedSubmitter(signedCertificationUpdater, ppmShipmentRouter)

	internalAPI.PpmResubmitPPMShipmentDocumentationHandler = ResubmitPPMShipmentDocumentationHandler{handlerConfig, ppmShipmentUpdatedSubmitter}

	internalAPI.TransportationOfficesGetTransportationOfficesHandler = GetTransportationOfficesHandler{
		handlerConfig,
		transportationOfficeFetcher,
	}

	paymentPacketCreator := ppmshipment.NewPaymentPacketCreator(ppmShipmentFetcher, pdfGenerator, AOAPacketCreator)
	internalAPI.PpmShowPaymentPacketHandler = ShowPaymentPacketHandler{handlerConfig, paymentPacketCreator}

	return internalAPI
}

// PDFProducer creates a new PDF producer
func PDFProducer() runtime.Producer {
	return runtime.ProducerFunc(func(writer io.Writer, data interface{}) error {
		rw, ok := data.(io.ReadCloser)
		if !ok {
			return errors.Errorf("could not convert %+v to io.ReadCloser", data)
		}
		_, err := io.Copy(writer, rw)
		return err
	})
}
