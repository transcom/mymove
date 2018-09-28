package internalapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"
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

	internalAPI.IssuesCreateIssueHandler = CreateIssueHandler{context}
	internalAPI.IssuesIndexIssuesHandler = IndexIssuesHandler{context}

	internalAPI.CertificationCreateSignedCertificationHandler = CreateSignedCertificationHandler{context}
	internalAPI.CertificationIndexSignedCertificationsHandler = IndexSignedCertificationsHandler{context}

	internalAPI.PpmCreatePersonallyProcuredMoveHandler = CreatePersonallyProcuredMoveHandler{context}
	internalAPI.PpmIndexPersonallyProcuredMovesHandler = IndexPersonallyProcuredMovesHandler{context}
	internalAPI.PpmPatchPersonallyProcuredMoveHandler = PatchPersonallyProcuredMoveHandler{context}
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

	internalAPI.MovesCreateMoveHandler = CreateMoveHandler{context}
	internalAPI.MovesPatchMoveHandler = PatchMoveHandler{context}
	internalAPI.MovesShowMoveHandler = ShowMoveHandler{context}
	internalAPI.MovesSubmitMoveForApprovalHandler = SubmitMoveHandler{context}

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
	internalAPI.ShipmentsSendHHGInvoiceHandler = ShipmentInvoiceHandler{context}

	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler{context}
	internalAPI.OfficeApprovePPMHandler = ApprovePPMHandler{context}
	internalAPI.OfficeApproveReimbursementHandler = ApproveReimbursementHandler{context}
	internalAPI.OfficeCancelMoveHandler = CancelMoveHandler{context}

	internalAPI.EntitlementsValidateEntitlementHandler = ValidateEntitlementHandler{context}

	internalAPI.GexSendGexRequestHandler = SendGexRequestHandler{context}

	internalAPI.CalendarShowUnavailableMoveDatesHandler = ShowUnavailableMoveDatesHandler{context}
	internalAPI.CalendarShowMoveDatesSummaryHandler = ShowMoveDatesSummaryHandler{context}

	return internalAPI.Serve(nil)
}
