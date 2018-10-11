package internalapi

import (
	"go.uber.org/dig"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/transcom/mymove/pkg/gen/internalapi"
	internalops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	userops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/users"
	"github.com/transcom/mymove/pkg/handlers"
)

// Handler is a package specific type for DI checking
type Handler http.Handler

// HandlerParams bundles up the dependencies of NewInternalApiHandler
type HandlerParams struct {
	dig.In
	sliuHandler userops.ShowLoggedInUserHandler
}

// NewInternalAPIHandler returns a handler for the internal API
func NewInternalAPIHandler(params HandlerParams, context handlers.HandlerContext) Handler {

	internalSpec, err := loads.Analyzed(internalapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	internalAPI := internalops.NewMymoveAPI(internalSpec)

	internalAPI.UsersShowLoggedInUserHandler = params.sliuHandler

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
	internalAPI.ShipmentsSendHHGInvoiceHandler = ShipmentInvoiceHandler{context}

	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler{context}
	internalAPI.OfficeApprovePPMHandler = ApprovePPMHandler{context}
	internalAPI.OfficeApproveReimbursementHandler = ApproveReimbursementHandler{context}
	internalAPI.OfficeCancelMoveHandler = CancelMoveHandler{context}

	internalAPI.EntitlementsValidateEntitlementHandler = ValidateEntitlementHandler{context}

	internalAPI.GexSendGexRequestHandler = SendGexRequestHandler{context}

	internalAPI.CalendarShowAvailableMoveDatesHandler = ShowAvailableMoveDatesHandler{context}

	return internalAPI.Serve(nil)
}
