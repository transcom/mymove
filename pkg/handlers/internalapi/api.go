package internalapi

import (
	"github.com/pkg/errors"
	"go.uber.org/dig"
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
	handlers.HandlerContext
	SliuHandler userops.ShowLoggedInUserHandler
}

// NewInternalAPIHandler returns a handler for the internal API
func NewInternalAPIHandler(p HandlerParams) (Handler, error) {

	internalSpec, err := loads.Analyzed(internalapi.SwaggerJSON, "")
	if err != nil {
		return nil, errors.Wrap(err, "Failed parsing internal Swagger")
	}
	internalAPI := internalops.NewMymoveAPI(internalSpec)

	internalAPI.UsersShowLoggedInUserHandler = p.SliuHandler

	internalAPI.IssuesCreateIssueHandler = CreateIssueHandler{p.HandlerContext}
	internalAPI.IssuesIndexIssuesHandler = IndexIssuesHandler{p.HandlerContext}

	internalAPI.CertificationCreateSignedCertificationHandler = CreateSignedCertificationHandler{p.HandlerContext}
	internalAPI.PpmCreatePersonallyProcuredMoveHandler = CreatePersonallyProcuredMoveHandler{p.HandlerContext}
	internalAPI.PpmIndexPersonallyProcuredMovesHandler = IndexPersonallyProcuredMovesHandler{p.HandlerContext}
	internalAPI.PpmPatchPersonallyProcuredMoveHandler = PatchPersonallyProcuredMoveHandler{p.HandlerContext}
	internalAPI.PpmSubmitPersonallyProcuredMoveHandler = SubmitPersonallyProcuredMoveHandler{p.HandlerContext}
	internalAPI.PpmShowPPMEstimateHandler = ShowPPMEstimateHandler{p.HandlerContext}
	internalAPI.PpmShowPPMSitEstimateHandler = ShowPPMSitEstimateHandler{p.HandlerContext}
	internalAPI.PpmShowPPMIncentiveHandler = ShowPPMIncentiveHandler{p.HandlerContext}
	internalAPI.PpmRequestPPMPaymentHandler = RequestPPMPaymentHandler{p.HandlerContext}
	internalAPI.PpmCreatePPMAttachmentsHandler = CreatePersonallyProcuredMoveAttachmentsHandler{p.HandlerContext}
	internalAPI.PpmRequestPPMExpenseSummaryHandler = RequestPPMExpenseSummaryHandler{p.HandlerContext}

	internalAPI.DutyStationsSearchDutyStationsHandler = SearchDutyStationsHandler{p.HandlerContext}

	internalAPI.TransportationOfficesShowDutyStationTransportationOfficeHandler = ShowDutyStationTransportationOfficeHandler{p.HandlerContext}

	internalAPI.OrdersCreateOrdersHandler = CreateOrdersHandler{p.HandlerContext}
	internalAPI.OrdersUpdateOrdersHandler = UpdateOrdersHandler{p.HandlerContext}
	internalAPI.OrdersShowOrdersHandler = ShowOrdersHandler{p.HandlerContext}

	internalAPI.MovesCreateMoveHandler = CreateMoveHandler{p.HandlerContext}
	internalAPI.MovesPatchMoveHandler = PatchMoveHandler{p.HandlerContext}
	internalAPI.MovesShowMoveHandler = ShowMoveHandler{p.HandlerContext}
	internalAPI.MovesSubmitMoveForApprovalHandler = SubmitMoveHandler{p.HandlerContext}
	internalAPI.MovesShowMoveDatesSummaryHandler = ShowMoveDatesSummaryHandler{p.HandlerContext}

	internalAPI.MoveDocsCreateGenericMoveDocumentHandler = CreateGenericMoveDocumentHandler{p.HandlerContext}
	internalAPI.MoveDocsUpdateMoveDocumentHandler = UpdateMoveDocumentHandler{p.HandlerContext}
	internalAPI.MoveDocsIndexMoveDocumentsHandler = IndexMoveDocumentsHandler{p.HandlerContext}

	internalAPI.MoveDocsCreateMovingExpenseDocumentHandler = CreateMovingExpenseDocumentHandler{p.HandlerContext}

	internalAPI.ServiceMembersCreateServiceMemberHandler = CreateServiceMemberHandler{p.HandlerContext}
	internalAPI.ServiceMembersPatchServiceMemberHandler = PatchServiceMemberHandler{p.HandlerContext}
	internalAPI.ServiceMembersShowServiceMemberHandler = ShowServiceMemberHandler{p.HandlerContext}
	internalAPI.ServiceMembersShowServiceMemberOrdersHandler = ShowServiceMemberOrdersHandler{p.HandlerContext}

	internalAPI.DocumentsCreateDocumentHandler = CreateDocumentHandler{p.HandlerContext}
	internalAPI.DocumentsShowDocumentHandler = ShowDocumentHandler{p.HandlerContext}
	internalAPI.UploadsCreateUploadHandler = CreateUploadHandler{p.HandlerContext}
	internalAPI.UploadsDeleteUploadHandler = DeleteUploadHandler{p.HandlerContext}
	internalAPI.UploadsDeleteUploadsHandler = DeleteUploadsHandler{p.HandlerContext}

	internalAPI.QueuesShowQueueHandler = ShowQueueHandler{p.HandlerContext}

	internalAPI.ShipmentsCreateShipmentHandler = CreateShipmentHandler{p.HandlerContext}
	internalAPI.ShipmentsPatchShipmentHandler = PatchShipmentHandler{p.HandlerContext}
	internalAPI.ShipmentsGetShipmentHandler = GetShipmentHandler{p.HandlerContext}
	internalAPI.ShipmentsApproveHHGHandler = ApproveHHGHandler{p.HandlerContext}
	internalAPI.ShipmentsCompleteHHGHandler = CompleteHHGHandler{p.HandlerContext}
	internalAPI.ShipmentsSendHHGInvoiceHandler = ShipmentInvoiceHandler{p.HandlerContext}

	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler{p.HandlerContext}
	internalAPI.OfficeApprovePPMHandler = ApprovePPMHandler{p.HandlerContext}
	internalAPI.OfficeApproveReimbursementHandler = ApproveReimbursementHandler{p.HandlerContext}
	internalAPI.OfficeCancelMoveHandler = CancelMoveHandler{p.HandlerContext}

	internalAPI.EntitlementsValidateEntitlementHandler = ValidateEntitlementHandler{p.HandlerContext}

	internalAPI.GexSendGexRequestHandler = SendGexRequestHandler{p.HandlerContext}

	internalAPI.CalendarShowAvailableMoveDatesHandler = ShowAvailableMoveDatesHandler{p.HandlerContext}
	internalAPI.DpsAuthGetCookieURLHandler = DPSAuthGetCookieURLHandler{p.HandlerContext}
	return internalAPI.Serve(nil), nil
}
