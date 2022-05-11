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
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	move "github.com/transcom/mymove/pkg/services/move"
	movedocument "github.com/transcom/mymove/pkg/services/move_documents"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"
	"github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	"github.com/transcom/mymove/pkg/services/order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	postalcodeservice "github.com/transcom/mymove/pkg/services/postal_codes"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
)

// NewInternalAPI returns the internal API
func NewInternalAPI(ctx handlers.HandlerContext) *internalops.MymoveAPI {

	internalSpec, err := loads.Analyzed(internalapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	internalAPI := internalops.NewMymoveAPI(internalSpec)

	internalAPI.ServeError = handlers.ServeCustomError
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := move.NewMoveRouter()
	ppmEstimator := ppmshipment.NewEstimatePPM(ctx.GHCPlanner(), &paymentrequesthelper.RequestPaymentHelper{})

	internalAPI.UsersShowLoggedInUserHandler = ShowLoggedInUserHandler{ctx, officeuser.NewOfficeUserFetcherPop()}
	internalAPI.CertificationCreateSignedCertificationHandler = CreateSignedCertificationHandler{ctx}
	internalAPI.CertificationIndexSignedCertificationHandler = IndexSignedCertificationsHandler{ctx}

	internalAPI.PpmCreatePersonallyProcuredMoveHandler = CreatePersonallyProcuredMoveHandler{ctx}
	internalAPI.PpmIndexPersonallyProcuredMovesHandler = IndexPersonallyProcuredMovesHandler{ctx}
	internalAPI.PpmPatchPersonallyProcuredMoveHandler = PatchPersonallyProcuredMoveHandler{ctx}
	internalAPI.PpmSubmitPersonallyProcuredMoveHandler = SubmitPersonallyProcuredMoveHandler{ctx}
	internalAPI.PpmShowPPMIncentiveHandler = ShowPPMIncentiveHandler{ctx}
	internalAPI.PpmRequestPPMPaymentHandler = RequestPPMPaymentHandler{ctx}
	internalAPI.PpmCreatePPMAttachmentsHandler = CreatePersonallyProcuredMoveAttachmentsHandler{ctx}
	internalAPI.PpmRequestPPMExpenseSummaryHandler = RequestPPMExpenseSummaryHandler{ctx}

	internalAPI.DutyLocationsSearchDutyLocationsHandler = SearchDutyLocationsHandler{ctx}

	internalAPI.AddressesShowAddressHandler = ShowAddressHandler{ctx}

	internalAPI.TransportationOfficesShowDutyLocationTransportationOfficeHandler = ShowDutyLocationTransportationOfficeHandler{ctx}

	internalAPI.OrdersCreateOrdersHandler = CreateOrdersHandler{ctx}
	internalAPI.OrdersUpdateOrdersHandler = UpdateOrdersHandler{ctx}
	internalAPI.OrdersShowOrdersHandler = ShowOrdersHandler{ctx}
	internalAPI.OrdersUploadAmendedOrdersHandler = UploadAmendedOrdersHandler{
		ctx,
		order.NewOrderUpdater(moveRouter),
	}

	internalAPI.MovesPatchMoveHandler = PatchMoveHandler{ctx}
	internalAPI.MovesShowMoveHandler = ShowMoveHandler{ctx}
	internalAPI.MovesSubmitMoveForApprovalHandler = SubmitMoveHandler{
		ctx,
		moveRouter,
	}
	internalAPI.MovesSubmitAmendedOrdersHandler = SubmitAmendedOrdersHandler{
		ctx,
		moveRouter,
	}
	internalAPI.MovesShowMoveDatesSummaryHandler = ShowMoveDatesSummaryHandler{ctx}

	internalAPI.MoveDocsCreateGenericMoveDocumentHandler = CreateGenericMoveDocumentHandler{ctx}
	internalAPI.MoveDocsUpdateMoveDocumentHandler = UpdateMoveDocumentHandler{ctx,
		movedocument.NewMoveDocumentUpdater(),
	}
	internalAPI.MoveDocsIndexMoveDocumentsHandler = IndexMoveDocumentsHandler{ctx}
	internalAPI.MoveDocsDeleteMoveDocumentHandler = DeleteMoveDocumentHandler{ctx}

	internalAPI.MoveDocsCreateMovingExpenseDocumentHandler = CreateMovingExpenseDocumentHandler{ctx}

	internalAPI.MoveDocsCreateWeightTicketDocumentHandler = CreateWeightTicketSetDocumentHandler{ctx}

	internalAPI.ServiceMembersCreateServiceMemberHandler = CreateServiceMemberHandler{ctx}
	internalAPI.ServiceMembersPatchServiceMemberHandler = PatchServiceMemberHandler{ctx}
	internalAPI.ServiceMembersShowServiceMemberHandler = ShowServiceMemberHandler{ctx}
	internalAPI.ServiceMembersShowServiceMemberOrdersHandler = ShowServiceMemberOrdersHandler{ctx}

	internalAPI.BackupContactsIndexServiceMemberBackupContactsHandler = IndexBackupContactsHandler{ctx}
	internalAPI.BackupContactsCreateServiceMemberBackupContactHandler = CreateBackupContactHandler{ctx}
	internalAPI.BackupContactsUpdateServiceMemberBackupContactHandler = UpdateBackupContactHandler{ctx}
	internalAPI.BackupContactsShowServiceMemberBackupContactHandler = ShowBackupContactHandler{ctx}

	internalAPI.DocumentsCreateDocumentHandler = CreateDocumentHandler{ctx}
	internalAPI.DocumentsShowDocumentHandler = ShowDocumentHandler{ctx}
	internalAPI.UploadsCreateUploadHandler = CreateUploadHandler{ctx}
	internalAPI.UploadsDeleteUploadHandler = DeleteUploadHandler{ctx}
	internalAPI.UploadsDeleteUploadsHandler = DeleteUploadsHandler{ctx}

	internalAPI.QueuesShowQueueHandler = ShowQueueHandler{ctx}
	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler{ctx, moveRouter}
	internalAPI.OfficeApprovePPMHandler = ApprovePPMHandler{ctx}
	internalAPI.OfficeApproveReimbursementHandler = ApproveReimbursementHandler{ctx}
	internalAPI.OfficeCancelMoveHandler = CancelMoveHandler{ctx, moveRouter}

	internalAPI.EntitlementsIndexEntitlementsHandler = IndexEntitlementsHandler{ctx}

	internalAPI.CalendarShowAvailableMoveDatesHandler = ShowAvailableMoveDatesHandler{ctx}

	internalAPI.DpsAuthGetCookieURLHandler = DPSAuthGetCookieURLHandler{ctx}

	internalAPI.MovesShowShipmentSummaryWorksheetHandler = ShowShipmentSummaryWorksheetHandler{ctx}

	internalAPI.RegisterProducer("application/pdf", PDFProducer())

	internalAPI.PostalCodesValidatePostalCodeWithRateDataHandler = ValidatePostalCodeWithRateDataHandler{
		ctx,
		postalcodeservice.NewPostalCodeValidator(),
	}

	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
	ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(ppmEstimator)
	shipmentCreator := shipment.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator)

	internalAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		ctx,
		shipmentCreator,
	}

	paymentRequestRecalculator := paymentrequest.NewPaymentRequestRecalculator(
		paymentrequest.NewPaymentRequestCreator(
			ctx.GHCPlanner(),
			ghcrateengine.NewServiceItemPricer(),
		),
		paymentrequest.NewPaymentRequestStatusUpdater(builder),
	)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(paymentRequestRecalculator)

	shipmentUpdater := shipment.NewShipmentUpdater(
		mtoshipment.NewMTOShipmentUpdater(
			builder,
			fetcher,
			ctx.Planner(),
			moveRouter,
			move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
			ctx.NotificationSender(),
			paymentRequestShipmentRecalculator,
		),
		ppmshipment.NewPPMShipmentUpdater(ppmEstimator),
	)

	internalAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		ctx,
		shipmentUpdater,
	}

	internalAPI.MtoShipmentListMTOShipmentsHandler = ListMTOShipmentsHandler{
		ctx,
		fetch.NewListFetcher(builder),
		fetch.NewFetcher(builder),
	}

	internalAPI.MtoShipmentDeleteShipmentHandler = DeleteShipmentHandler{
		ctx,
		mtoshipment.NewShipmentDeleter(),
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
