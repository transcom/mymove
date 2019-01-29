package internalapi

import (
	"io"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/gen/internalapi"
	internalops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	"github.com/transcom/mymove/pkg/handlers"
)

// NewInternalAPIHandler returns a handler for the internal API
func NewInternalAPIHandler(context handlers.HandlerContext) http.Handler {

	internalSpec, err := loads.Analyzed(internalapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	internalAPI := internalops.NewMymoveAPI(internalSpec)

	internalAPI.UsersShowLoggedInUserHandler = ShowLoggedInUserHandler{context}

	internalAPI.CertificationCreateSignedCertificationHandler = CreateSignedCertificationHandler{context}

	internalAPI.PpmCreatePersonallyProcuredMoveHandler = CreatePersonallyProcuredMoveHandler{context}
	internalAPI.PpmIndexPersonallyProcuredMovesHandler = IndexPersonallyProcuredMovesHandler{context}
	internalAPI.PpmPatchPersonallyProcuredMoveHandler = PatchPersonallyProcuredMoveHandler{context}
	internalAPI.PpmSubmitPersonallyProcuredMoveHandler = SubmitPersonallyProcuredMoveHandler{context}
	internalAPI.PpmShowPPMEstimateHandler = ShowPPMEstimateHandler{context}
	internalAPI.PpmShowPPMSitEstimateHandler = ShowPPMSitEstimateHandler{context}
	internalAPI.PpmShowPPMIncentiveHandler = ShowPPMIncentiveHandler{context}
	internalAPI.PpmRequestPPMPaymentHandler = RequestPPMPaymentHandler{context}
	internalAPI.PpmCreatePPMAttachmentsHandler = CreatePersonallyProcuredMoveAttachmentsHandler{context}
	internalAPI.PpmRequestPPMExpenseSummaryHandler = RequestPPMExpenseSummaryHandler{context}

	internalAPI.DutyStationsSearchDutyStationsHandler = SearchDutyStationsHandler{context}

	internalAPI.TransportationOfficesShowDutyStationTransportationOfficeHandler = ShowDutyStationTransportationOfficeHandler{context}

	internalAPI.OrdersCreateOrdersHandler = CreateOrdersHandler{context}
	internalAPI.OrdersUpdateOrdersHandler = UpdateOrdersHandler{context}
	internalAPI.OrdersShowOrdersHandler = ShowOrdersHandler{context}

	internalAPI.MovesPatchMoveHandler = PatchMoveHandler{context}
	internalAPI.MovesShowMoveHandler = ShowMoveHandler{context}
	internalAPI.MovesSubmitMoveForApprovalHandler = SubmitMoveHandler{context}
	internalAPI.MovesShowMoveDatesSummaryHandler = ShowMoveDatesSummaryHandler{context}

	internalAPI.MoveDocsCreateGenericMoveDocumentHandler = CreateGenericMoveDocumentHandler{context}
	internalAPI.MoveDocsUpdateMoveDocumentHandler = UpdateMoveDocumentHandler{context}
	internalAPI.MoveDocsIndexMoveDocumentsHandler = IndexMoveDocumentsHandler{context}

	internalAPI.MoveDocsCreateMovingExpenseDocumentHandler = CreateMovingExpenseDocumentHandler{context}

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

	internalAPI.ShipmentsCreateShipmentHandler = CreateShipmentHandler{context}
	internalAPI.ShipmentsPatchShipmentHandler = PatchShipmentHandler{context}
	internalAPI.ShipmentsGetShipmentHandler = GetShipmentHandler{context}
	internalAPI.ShipmentsApproveHHGHandler = ApproveHHGHandler{context}
	internalAPI.ShipmentsCompleteHHGHandler = CompleteHHGHandler{context}
	internalAPI.ShipmentsCreateAndSendHHGInvoiceHandler = ShipmentInvoiceHandler{context}

	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler{context}
	internalAPI.OfficeApprovePPMHandler = ApprovePPMHandler{context}
	internalAPI.OfficeApproveReimbursementHandler = ApproveReimbursementHandler{context}
	internalAPI.OfficeCancelMoveHandler = CancelMoveHandler{context}

	internalAPI.EntitlementsValidateEntitlementHandler = ValidateEntitlementHandler{context}

	internalAPI.CalendarShowAvailableMoveDatesHandler = ShowAvailableMoveDatesHandler{context}

	internalAPI.DpsAuthGetCookieURLHandler = DPSAuthGetCookieURLHandler{context}

	internalAPI.MovesShowShipmentSummaryWorksheetHandler = ShowShipmentSummaryWorksheetHandler{context}

	internalAPI.ApplicationPdfProducer = PDFProducer()

	return internalAPI.Serve(nil)
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
