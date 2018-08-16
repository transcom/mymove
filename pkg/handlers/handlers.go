package handlers

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/internalapi"
	internalops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	"github.com/transcom/mymove/pkg/gen/ordersapi"
	ordersops "github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/handlers/internal"
	"github.com/transcom/mymove/pkg/handlers/orders"
	"github.com/transcom/mymove/pkg/handlers/public"
	"github.com/transcom/mymove/pkg/handlers/utils"
)

// NewPublicAPIHandler returns a handler for the public API
func NewPublicAPIHandler(context utils.HandlerContext) http.Handler {

	// Wire up the handlers to the publicAPIMux
	apiSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	publicAPI := publicops.NewMymoveAPI(apiSpec)

	// Blackouts

	// Documents

	// Shipments
	publicAPI.ShipmentsIndexShipmentsHandler = public.PublicIndexShipmentsHandler(context)
	publicAPI.ShipmentsGetShipmentHandler = public.PublicGetShipmentHandler(context)
	publicAPI.ShipmentsCreateShipmentAcceptHandler = public.PublicCreateShipmentAcceptHandler(context)
	publicAPI.ShipmentsCreateShipmentRejectHandler = public.PublicCreateShipmentRejectHandler(context)

	// TSPs
	publicAPI.TspsIndexTSPsHandler = public.PublicTspsIndexTSPsHandler(context)
	publicAPI.TspsGetTspShipmentsHandler = public.PublicTspsGetTspShipmentsHandler(context)

	return publicAPI.Serve(nil)
}

// NewInternalAPIHandler returns a handler for the internal API
func NewInternalAPIHandler(context utils.HandlerContext) http.Handler {

	internalSpec, err := loads.Analyzed(internalapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	internalAPI := internalops.NewMymoveAPI(internalSpec)

	internalAPI.UsersShowLoggedInUserHandler = internal.ShowLoggedInUserHandler(context)

	internalAPI.IssuesCreateIssueHandler = internal.CreateIssueHandler(context)
	internalAPI.IssuesIndexIssuesHandler = internal.IndexIssuesHandler(context)

	internalAPI.CertificationCreateSignedCertificationHandler = internal.CreateSignedCertificationHandler(context)
	internalAPI.CertificationIndexSignedCertificationsHandler = internal.IndexSignedCertificationsHandler(context)

	internalAPI.PpmCreatePersonallyProcuredMoveHandler = internal.CreatePersonallyProcuredMoveHandler(context)
	internalAPI.PpmIndexPersonallyProcuredMovesHandler = internal.IndexPersonallyProcuredMovesHandler(context)
	internalAPI.PpmPatchPersonallyProcuredMoveHandler = internal.PatchPersonallyProcuredMoveHandler(context)
	internalAPI.PpmShowPPMEstimateHandler = internal.ShowPPMEstimateHandler(context)
	internalAPI.PpmShowPPMSitEstimateHandler = internal.ShowPPMSitEstimateHandler(context)
	internalAPI.PpmShowPPMIncentiveHandler = internal.ShowPPMIncentiveHandler(context)
	internalAPI.PpmRequestPPMPaymentHandler = internal.RequestPPMPaymentHandler(context)
	internalAPI.PpmCreatePPMAttachmentsHandler = internal.CreatePersonallyProcuredMoveAttachmentsHandler(context)
	internalAPI.PpmRequestPPMExpenseSummaryHandler = internal.RequestPPMExpenseSummaryHandler(context)

	internalAPI.DutyStationsSearchDutyStationsHandler = internal.SearchDutyStationsHandler(context)

	internalAPI.TransportationOfficesShowDutyStationTransportationOfficeHandler = internal.ShowDutyStationTransportationOfficeHandler(context)

	internalAPI.OrdersCreateOrdersHandler = internal.CreateOrdersHandler(context)
	internalAPI.OrdersUpdateOrdersHandler = internal.UpdateOrdersHandler(context)
	internalAPI.OrdersShowOrdersHandler = internal.ShowOrdersHandler(context)

	internalAPI.MovesCreateMoveHandler = internal.CreateMoveHandler(context)
	internalAPI.MovesPatchMoveHandler = internal.PatchMoveHandler(context)
	internalAPI.MovesShowMoveHandler = internal.ShowMoveHandler(context)
	internalAPI.MovesSubmitMoveForApprovalHandler = internal.SubmitMoveHandler(context)

	internalAPI.MoveDocsCreateGenericMoveDocumentHandler = internal.CreateGenericMoveDocumentHandler(context)
	internalAPI.MoveDocsUpdateMoveDocumentHandler = internal.UpdateMoveDocumentHandler(context)
	internalAPI.MoveDocsIndexMoveDocumentsHandler = internal.IndexMoveDocumentsHandler(context)

	internalAPI.MoveDocsCreateMovingExpenseDocumentHandler = internal.CreateMovingExpenseDocumentHandler(context)

	internalAPI.ServiceMembersCreateServiceMemberHandler = internal.CreateServiceMemberHandler(context)
	internalAPI.ServiceMembersPatchServiceMemberHandler = internal.PatchServiceMemberHandler(context)
	internalAPI.ServiceMembersShowServiceMemberHandler = internal.ShowServiceMemberHandler(context)
	internalAPI.ServiceMembersShowServiceMemberOrdersHandler = internal.ShowServiceMemberOrdersHandler(context)

	internalAPI.BackupContactsIndexServiceMemberBackupContactsHandler = internal.IndexBackupContactsHandler(context)
	internalAPI.BackupContactsCreateServiceMemberBackupContactHandler = internal.CreateBackupContactHandler(context)
	internalAPI.BackupContactsUpdateServiceMemberBackupContactHandler = internal.UpdateBackupContactHandler(context)
	internalAPI.BackupContactsShowServiceMemberBackupContactHandler = internal.ShowBackupContactHandler(context)

	internalAPI.DocumentsCreateDocumentHandler = internal.CreateDocumentHandler(context)
	internalAPI.DocumentsShowDocumentHandler = internal.ShowDocumentHandler(context)
	internalAPI.UploadsCreateUploadHandler = internal.CreateUploadHandler(context)
	internalAPI.UploadsDeleteUploadHandler = internal.DeleteUploadHandler(context)
	internalAPI.UploadsDeleteUploadsHandler = internal.DeleteUploadsHandler(context)

	internalAPI.QueuesShowQueueHandler = internal.ShowQueueHandler(context)

	internalAPI.ShipmentsCreateShipmentHandler = internal.CreateShipmentHandler(context)
	internalAPI.ShipmentsPatchShipmentHandler = internal.PatchShipmentHandler(context)
	internalAPI.ShipmentsGetShipmentHandler = internal.GetShipmentHandler(context)

	internalAPI.OfficeApproveMoveHandler = internal.ApproveMoveHandler(context)
	internalAPI.OfficeApprovePPMHandler = internal.ApprovePPMHandler(context)
	internalAPI.OfficeApproveReimbursementHandler = internal.ApproveReimbursementHandler(context)
	internalAPI.OfficeCancelMoveHandler = internal.CancelMoveHandler(context)

	internalAPI.EntitlementsValidateEntitlementHandler = internal.ValidateEntitlementHandler(context)

	return internalAPI.Serve(nil)
}

// NewOrdersAPIHandler returns a handler for the Orders API
func NewOrdersAPIHandler(context utils.HandlerContext) http.Handler {

	// Wire up the handlers to the ordersAPIMux
	ordersSpec, err := loads.Analyzed(ordersapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	ordersAPI := ordersops.NewMymoveAPI(ordersSpec)
	ordersAPI.GetOrdersHandler = orders.GetOrdersHandler(context)
	ordersAPI.IndexOrdersHandler = orders.IndexOrdersHandler(context)
	ordersAPI.PostRevisionHandler = orders.PostRevisionHandler(context)
	ordersAPI.PostRevisionToOrdersHandler = orders.PostRevisionToOrdersHandler(context)
	return ordersAPI.Serve(nil)
}
