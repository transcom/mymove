package internalapi

import (
	"io"
	"log"

	"github.com/transcom/mymove/pkg/services/move"
	officeuser "github.com/transcom/mymove/pkg/services/office_user"

	"github.com/transcom/mymove/pkg/services/fetch"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"

	accesscodeservice "github.com/transcom/mymove/pkg/services/accesscode"
	movedocument "github.com/transcom/mymove/pkg/services/move_documents"
	postalcodeservice "github.com/transcom/mymove/pkg/services/postal_codes"
	"github.com/transcom/mymove/pkg/services/ppmservices"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/gen/internalapi"
	internalops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewInternalAPI returns the internal API
func NewInternalAPI(context handlers.HandlerContext) *internalops.MymoveAPI {

	internalSpec, err := loads.Analyzed(internalapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	internalAPI := internalops.NewMymoveAPI(internalSpec)

	internalAPI.ServeError = handlers.ServeCustomError
	builder := query.NewQueryBuilder(context.DB())
	fetcher := fetch.NewFetcher(builder)

	internalAPI.UsersShowLoggedInUserHandler = ShowLoggedInUserHandler{context, officeuser.NewOfficeUserFetcherPop(context.DB())}
	internalAPI.CertificationCreateSignedCertificationHandler = CreateSignedCertificationHandler{context}
	internalAPI.CertificationIndexSignedCertificationHandler = IndexSignedCertificationsHandler{context}

	internalAPI.PpmCreatePersonallyProcuredMoveHandler = CreatePersonallyProcuredMoveHandler{context}
	internalAPI.PpmIndexPersonallyProcuredMovesHandler = IndexPersonallyProcuredMovesHandler{context}
	internalAPI.PpmPatchPersonallyProcuredMoveHandler = PatchPersonallyProcuredMoveHandler{context}
	internalAPI.PpmUpdatePersonallyProcuredMoveEstimateHandler = UpdatePersonallyProcuredMoveEstimateHandler{context, ppmservices.NewEstimateCalculator(context.DB(), context.Planner())}
	internalAPI.PpmSubmitPersonallyProcuredMoveHandler = SubmitPersonallyProcuredMoveHandler{context}
	internalAPI.PpmShowPPMEstimateHandler = ShowPPMEstimateHandler{context}
	internalAPI.PpmShowPPMSitEstimateHandler = ShowPPMSitEstimateHandler{context, ppmservices.NewEstimateCalculator(context.DB(), context.Planner())}
	internalAPI.PpmShowPPMIncentiveHandler = ShowPPMIncentiveHandler{context}
	internalAPI.PpmRequestPPMPaymentHandler = RequestPPMPaymentHandler{context}
	internalAPI.PpmCreatePPMAttachmentsHandler = CreatePersonallyProcuredMoveAttachmentsHandler{context}
	internalAPI.PpmRequestPPMExpenseSummaryHandler = RequestPPMExpenseSummaryHandler{context}

	internalAPI.DutyStationsSearchDutyStationsHandler = SearchDutyStationsHandler{context}

	internalAPI.AddressesShowAddressHandler = ShowAddressHandler{context}

	internalAPI.TransportationOfficesShowDutyStationTransportationOfficeHandler = ShowDutyStationTransportationOfficeHandler{context}

	internalAPI.OrdersCreateOrdersHandler = CreateOrdersHandler{context}
	internalAPI.OrdersUpdateOrdersHandler = UpdateOrdersHandler{context}
	internalAPI.OrdersShowOrdersHandler = ShowOrdersHandler{context}

	internalAPI.MovesPatchMoveHandler = PatchMoveHandler{context}
	internalAPI.MovesShowMoveHandler = ShowMoveHandler{context}
	internalAPI.MovesSubmitMoveForApprovalHandler = SubmitMoveHandler{
		context,
		move.NewMoveStatusRouter(context.DB()),
	}
	internalAPI.MovesShowMoveDatesSummaryHandler = ShowMoveDatesSummaryHandler{context}

	internalAPI.MoveDocsCreateGenericMoveDocumentHandler = CreateGenericMoveDocumentHandler{context}
	internalAPI.MoveDocsUpdateMoveDocumentHandler = UpdateMoveDocumentHandler{context,
		movedocument.NewMoveDocumentUpdater(context.DB()),
	}
	internalAPI.MoveDocsIndexMoveDocumentsHandler = IndexMoveDocumentsHandler{context}
	internalAPI.MoveDocsDeleteMoveDocumentHandler = DeleteMoveDocumentHandler{context}

	internalAPI.MoveDocsCreateMovingExpenseDocumentHandler = CreateMovingExpenseDocumentHandler{context}

	internalAPI.MoveDocsCreateWeightTicketDocumentHandler = CreateWeightTicketSetDocumentHandler{context}

	internalAPI.ServiceMembersCreateServiceMemberHandler = CreateServiceMemberHandler{context}
	internalAPI.ServiceMembersPatchServiceMemberHandler = PatchServiceMemberHandler{context}
	internalAPI.ServiceMembersShowServiceMemberHandler = ShowServiceMemberHandler{context}
	internalAPI.ServiceMembersShowServiceMemberOrdersHandler = ShowServiceMemberOrdersHandler{context}

	internalAPI.BackupContactsIndexServiceMemberBackupContactsHandler = IndexBackupContactsHandler{context}
	internalAPI.BackupContactsCreateServiceMemberBackupContactHandler = CreateBackupContactHandler{context}
	internalAPI.BackupContactsUpdateServiceMemberBackupContactHandler = UpdateBackupContactHandler{context}
	internalAPI.BackupContactsShowServiceMemberBackupContactHandler = ShowBackupContactHandler{context}

	internalAPI.DocumentsCreateDocumentHandler = CreateDocumentHandler{context}
	internalAPI.DocumentsShowDocumentHandler = ShowDocumentHandler{context}
	internalAPI.UploadsCreateUploadHandler = CreateUploadHandler{context}
	internalAPI.UploadsDeleteUploadHandler = DeleteUploadHandler{context}
	internalAPI.UploadsDeleteUploadsHandler = DeleteUploadsHandler{context}

	internalAPI.QueuesShowQueueHandler = ShowQueueHandler{context}

	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler{context}
	internalAPI.OfficeApprovePPMHandler = ApprovePPMHandler{context}
	internalAPI.OfficeApproveReimbursementHandler = ApproveReimbursementHandler{context}
	internalAPI.OfficeCancelMoveHandler = CancelMoveHandler{context}

	internalAPI.EntitlementsIndexEntitlementsHandler = IndexEntitlementsHandler{context}
	internalAPI.EntitlementsValidateEntitlementHandler = ValidateEntitlementHandler{context}

	internalAPI.CalendarShowAvailableMoveDatesHandler = ShowAvailableMoveDatesHandler{context}

	internalAPI.DpsAuthGetCookieURLHandler = DPSAuthGetCookieURLHandler{context}

	internalAPI.MovesShowShipmentSummaryWorksheetHandler = ShowShipmentSummaryWorksheetHandler{context}

	internalAPI.RegisterProducer("application/pdf", PDFProducer())

	internalAPI.PostalCodesValidatePostalCodeWithRateDataHandler = ValidatePostalCodeWithRateDataHandler{
		context,
		postalcodeservice.NewPostalCodeValidator(context.DB()),
	}

	// Access Codes
	internalAPI.AccesscodeFetchAccessCodeHandler = FetchAccessCodeHandler{context, accesscodeservice.NewAccessCodeFetcher(context.DB())}
	internalAPI.AccesscodeValidateAccessCodeHandler = ValidateAccessCodeHandler{context, accesscodeservice.NewAccessCodeValidator(context.DB())}
	internalAPI.AccesscodeClaimAccessCodeHandler = ClaimAccessCodeHandler{context, accesscodeservice.NewAccessCodeClaimer(context.DB())}

	// GHC Endpoint

	internalAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		context,
		mtoshipment.NewMTOShipmentCreator(context.DB(), builder, fetcher),
	}

	internalAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		context,
		mtoshipment.NewMTOShipmentUpdater(context.DB(), builder, fetcher, context.Planner()),
	}

	internalAPI.MtoShipmentListMTOShipmentsHandler = ListMTOShipmentsHandler{
		context,
		fetch.NewListFetcher(builder),
		fetch.NewFetcher(builder),
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
