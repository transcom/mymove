package internalapi

import (
	"io"
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalapi"
	internalops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	move "github.com/transcom/mymove/pkg/services/move"
	movingexpense "github.com/transcom/mymove/pkg/services/moving_expense"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	"github.com/transcom/mymove/pkg/services/order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	postalcodeservice "github.com/transcom/mymove/pkg/services/postal_codes"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	progear "github.com/transcom/mymove/pkg/services/progear_weight_ticket"
	"github.com/transcom/mymove/pkg/services/query"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	weightticket "github.com/transcom/mymove/pkg/services/weight_ticket"
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
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	mtoShipmentRouter := mtoshipment.NewShipmentRouter()
	ppmShipmentRouter := ppmshipment.NewPPMShipmentRouter(mtoShipmentRouter)
	transportationOfficeFetcher := transportationoffice.NewTransportationOfficesFetcher()

	internalAPI.UsersShowLoggedInUserHandler = ShowLoggedInUserHandler{handlerConfig, officeuser.NewOfficeUserFetcherPop()}
	internalAPI.CertificationCreateSignedCertificationHandler = CreateSignedCertificationHandler{handlerConfig}
	internalAPI.CertificationIndexSignedCertificationHandler = IndexSignedCertificationsHandler{handlerConfig}

	internalAPI.PpmCreatePersonallyProcuredMoveHandler = CreatePersonallyProcuredMoveHandler{handlerConfig}
	internalAPI.PpmIndexPersonallyProcuredMovesHandler = IndexPersonallyProcuredMovesHandler{handlerConfig}
	internalAPI.PpmPatchPersonallyProcuredMoveHandler = PatchPersonallyProcuredMoveHandler{handlerConfig}
	internalAPI.PpmSubmitPersonallyProcuredMoveHandler = SubmitPersonallyProcuredMoveHandler{handlerConfig}
	internalAPI.PpmShowPPMIncentiveHandler = ShowPPMIncentiveHandler{handlerConfig}
	internalAPI.PpmRequestPPMPaymentHandler = RequestPPMPaymentHandler{handlerConfig}

	internalAPI.DutyLocationsSearchDutyLocationsHandler = SearchDutyLocationsHandler{handlerConfig}

	internalAPI.AddressesShowAddressHandler = ShowAddressHandler{handlerConfig}

	internalAPI.TransportationOfficesShowDutyLocationTransportationOfficeHandler = ShowDutyLocationTransportationOfficeHandler{handlerConfig}

	internalAPI.OrdersCreateOrdersHandler = CreateOrdersHandler{handlerConfig}
	internalAPI.OrdersUpdateOrdersHandler = UpdateOrdersHandler{handlerConfig}
	internalAPI.OrdersShowOrdersHandler = ShowOrdersHandler{handlerConfig}
	internalAPI.OrdersUploadAmendedOrdersHandler = UploadAmendedOrdersHandler{
		handlerConfig,
		order.NewOrderUpdater(moveRouter),
	}

	internalAPI.MovesPatchMoveHandler = PatchMoveHandler{handlerConfig, transportationOfficeFetcher}
	internalAPI.MovesShowMoveHandler = ShowMoveHandler{handlerConfig}
	internalAPI.MovesSubmitMoveForApprovalHandler = SubmitMoveHandler{
		handlerConfig,
		moveRouter,
	}
	internalAPI.MovesSubmitAmendedOrdersHandler = SubmitAmendedOrdersHandler{
		handlerConfig,
		moveRouter,
	}
	internalAPI.MovesShowMoveDatesSummaryHandler = ShowMoveDatesSummaryHandler{handlerConfig}

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
	internalAPI.UploadsDeleteUploadHandler = DeleteUploadHandler{handlerConfig}
	internalAPI.UploadsDeleteUploadsHandler = DeleteUploadsHandler{handlerConfig}

	internalAPI.QueuesShowQueueHandler = ShowQueueHandler{handlerConfig}
	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler{handlerConfig, moveRouter}
	internalAPI.OfficeApprovePPMHandler = ApprovePPMHandler{handlerConfig}
	internalAPI.OfficeApproveReimbursementHandler = ApproveReimbursementHandler{handlerConfig}
	internalAPI.OfficeCancelMoveHandler = CancelMoveHandler{handlerConfig, moveRouter}

	internalAPI.EntitlementsIndexEntitlementsHandler = IndexEntitlementsHandler{handlerConfig}

	internalAPI.CalendarShowAvailableMoveDatesHandler = ShowAvailableMoveDatesHandler{handlerConfig}

	internalAPI.MovesShowShipmentSummaryWorksheetHandler = ShowShipmentSummaryWorksheetHandler{handlerConfig}

	internalAPI.RegisterProducer("application/pdf", PDFProducer())

	internalAPI.PostalCodesValidatePostalCodeWithRateDataHandler = ValidatePostalCodeWithRateDataHandler{
		handlerConfig,
		postalcodeservice.NewPostalCodeValidator(),
	}

	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
	shipmentRouter := mtoshipment.NewShipmentRouter()
	shipmentCreator := shipment.NewShipmentCreator(mtoShipmentCreator, ppmshipment.NewPPMShipmentCreator(ppmEstimator), shipmentRouter)

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

	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()
	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(ppmEstimator, addressCreator, addressUpdater)
	shipmentUpdater := shipment.NewShipmentUpdater(
		mtoshipment.NewCustomerMTOShipmentUpdater(
			builder,
			fetcher,
			handlerConfig.Planner(),
			moveRouter,
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
			handlerConfig.NotificationSender(),
			paymentRequestShipmentRecalculator,
		),
		ppmShipmentUpdater,
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
		mtoshipment.NewShipmentDeleter(),
	}

	internalAPI.PpmCreateMovingExpenseHandler = CreateMovingExpenseHandler{handlerConfig, movingexpense.NewMovingExpenseCreator()}
	internalAPI.PpmUpdateMovingExpenseHandler = UpdateMovingExpenseHandler{handlerConfig, movingexpense.NewMovingExpenseUpdater()}

	internalAPI.PpmCreateWeightTicketHandler = CreateWeightTicketHandler{handlerConfig, weightticket.NewCustomerWeightTicketCreator()}

	weightTicketFetcher := weightticket.NewWeightTicketFetcher()
	internalAPI.PpmUpdateWeightTicketHandler = UpdateWeightTicketHandler{handlerConfig, weightticket.NewCustomerWeightTicketUpdater(weightTicketFetcher, ppmShipmentUpdater)}

	internalAPI.PpmCreateProGearWeightTicketHandler = CreateProGearWeightTicketHandler{handlerConfig, progear.NewCustomerProgearWeightTicketCreator()}
	internalAPI.PpmUpdateProGearWeightTicketHandler = UpdateProGearWeightTicketHandler{handlerConfig, progear.NewCustomerProgearWeightTicketUpdater()}

	internalAPI.PpmCreatePPMUploadHandler = CreatePPMUploadHandler{handlerConfig}

	ppmShipmentNewSubmitter := ppmshipment.NewPPMShipmentNewSubmitter(signedCertificationCreator, ppmShipmentRouter)

	internalAPI.PpmSubmitPPMShipmentDocumentationHandler = SubmitPPMShipmentDocumentationHandler{handlerConfig, ppmShipmentNewSubmitter}

	ppmShipmentUpdatedSubmitter := ppmshipment.NewPPMShipmentUpdatedSubmitter(signedCertificationUpdater, ppmShipmentRouter)

	internalAPI.PpmResubmitPPMShipmentDocumentationHandler = ResubmitPPMShipmentDocumentationHandler{handlerConfig, ppmShipmentUpdatedSubmitter}

	internalAPI.TransportationOfficesGetTransportationOfficesHandler = GetTransportationOfficesHandler{
		handlerConfig,
		transportationOfficeFetcher,
	}

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
