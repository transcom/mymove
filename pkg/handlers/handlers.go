package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"go.uber.org/zap"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalapi"
	internalops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/ordersapi"
	ordersops "github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/storage"
)

// HandlerContext contains dependencies that are shared between all handlers.
// Each individual handler is declared as a type alias for HandlerContext so that the Handle() method
// can be declared on it. When wiring up a handler, you can create a HandlerContext and cast it to the type you want.
type HandlerContext struct {
	db                 *pop.Connection
	logger             *zap.Logger
	cookieSecret       string
	noSessionTimeout   bool
	planner            route.Planner
	storage            storage.FileStorer
	notificationSender notifications.NotificationSender
}

// NewHandlerContext returns a new HandlerContext with its required private fields set.
func NewHandlerContext(db *pop.Connection, logger *zap.Logger) HandlerContext {
	return HandlerContext{
		db:     db,
		logger: logger,
	}
}

// SetFileStorer is a simple setter for storage private field
func (context *HandlerContext) SetFileStorer(storer storage.FileStorer) {
	context.storage = storer
}

// SetNotificationSender is a simple setter for AWS SES private field
func (context *HandlerContext) SetNotificationSender(sender notifications.NotificationSender) {
	context.notificationSender = sender
}

// SetPlanner is a simple setter for the route.Planner private field
func (context *HandlerContext) SetPlanner(planner route.Planner) {
	context.planner = planner
}

// SetCookieSecret is a simple setter for the cookieSeecret private Field
func (context *HandlerContext) SetCookieSecret(cookieSecret string) {
	context.cookieSecret = cookieSecret
}

// SetNoSessionTimeout is a simple setter for the noSessionTimeout private Field
func (context *HandlerContext) SetNoSessionTimeout() {
	context.noSessionTimeout = true
}

// CookieUpdateResponder wraps a swagger middleware.Responder in code which sets the session_cookie
// See: https://github.com/go-swagger/go-swagger/issues/748
type CookieUpdateResponder struct {
	session          *auth.Session
	cookieSecret     string
	noSessionTimeout bool
	logger           *zap.Logger
	responder        middleware.Responder
}

// NewCookieUpdateResponder constructs a wrapper for the responder which will update cookies
func NewCookieUpdateResponder(request *http.Request, secret string, noSessionTimeout bool, logger *zap.Logger, responder middleware.Responder) middleware.Responder {
	return &CookieUpdateResponder{
		session:          auth.SessionFromRequestContext(request),
		cookieSecret:     secret,
		noSessionTimeout: noSessionTimeout,
		logger:           logger,
		responder:        responder,
	}
}

// WriteResponse updates the session cookie before writing out the details of the response
func (cur *CookieUpdateResponder) WriteResponse(rw http.ResponseWriter, p runtime.Producer) {
	auth.WriteSessionCookie(rw, cur.session, cur.cookieSecret, cur.noSessionTimeout, cur.logger)
	cur.responder.WriteResponse(rw, p)
}

// NewPublicAPIHandler returns a handler for the public API
func NewPublicAPIHandler(context HandlerContext) http.Handler {

	// Wire up the handlers to the publicAPIMux
	apiSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	publicAPI := publicops.NewMymoveAPI(apiSpec)

	// Blackouts

	// Documents

	// Shipments
	publicAPI.ShipmentsIndexShipmentsHandler = PublicIndexShipmentsHandler(context)
	publicAPI.ShipmentsGetShipmentHandler = PublicGetShipmentHandler(context)
	publicAPI.ShipmentsCreateShipmentAcceptHandler = PublicCreateShipmentAcceptHandler(context)
	publicAPI.ShipmentsCreateShipmentRejectHandler = PublicCreateShipmentRejectHandler(context)

	// TSPs
	publicAPI.TspsIndexTSPsHandler = PublicTspsIndexTSPsHandler(context)
	publicAPI.TspsGetTspShipmentsHandler = PublicTspsGetTspShipmentsHandler(context)

	return publicAPI.Serve(nil)
}

// NewInternalAPIHandler returns a handler for the internal API
func NewInternalAPIHandler(context HandlerContext) http.Handler {

	internalSpec, err := loads.Analyzed(internalapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	internalAPI := internalops.NewMymoveAPI(internalSpec)

	internalAPI.UsersShowLoggedInUserHandler = ShowLoggedInUserHandler(context)

	internalAPI.IssuesCreateIssueHandler = CreateIssueHandler(context)
	internalAPI.IssuesIndexIssuesHandler = IndexIssuesHandler(context)

	internalAPI.CertificationCreateSignedCertificationHandler = CreateSignedCertificationHandler(context)
	internalAPI.CertificationIndexSignedCertificationsHandler = IndexSignedCertificationsHandler(context)

	internalAPI.PpmCreatePersonallyProcuredMoveHandler = CreatePersonallyProcuredMoveHandler(context)
	internalAPI.PpmIndexPersonallyProcuredMovesHandler = IndexPersonallyProcuredMovesHandler(context)
	internalAPI.PpmPatchPersonallyProcuredMoveHandler = PatchPersonallyProcuredMoveHandler(context)
	internalAPI.PpmShowPPMEstimateHandler = ShowPPMEstimateHandler(context)
	internalAPI.PpmShowPPMSitEstimateHandler = ShowPPMSitEstimateHandler(context)
	internalAPI.PpmShowPPMIncentiveHandler = ShowPPMIncentiveHandler(context)
	internalAPI.PpmRequestPPMPaymentHandler = RequestPPMPaymentHandler(context)

	internalAPI.DutyStationsSearchDutyStationsHandler = SearchDutyStationsHandler(context)

	internalAPI.TransportationOfficesShowDutyStationTransportationOfficeHandler = ShowDutyStationTransportationOfficeHandler(context)

	internalAPI.OrdersCreateOrdersHandler = CreateOrdersHandler(context)
	internalAPI.OrdersUpdateOrdersHandler = UpdateOrdersHandler(context)
	internalAPI.OrdersShowOrdersHandler = ShowOrdersHandler(context)

	internalAPI.MovesCreateMoveHandler = CreateMoveHandler(context)
	internalAPI.MovesPatchMoveHandler = PatchMoveHandler(context)
	internalAPI.MovesShowMoveHandler = ShowMoveHandler(context)
	internalAPI.MovesSubmitMoveForApprovalHandler = SubmitMoveHandler(context)

	internalAPI.MoveDocsCreateGenericMoveDocumentHandler = CreateGenericMoveDocumentHandler(context)
	internalAPI.MoveDocsUpdateMoveDocumentHandler = UpdateMoveDocumentHandler(context)
	internalAPI.MoveDocsIndexMoveDocumentsHandler = IndexMoveDocumentsHandler(context)

	internalAPI.MoveDocsCreateMovingExpenseDocumentHandler = CreateMovingExpenseDocumentHandler(context)

	internalAPI.ServiceMembersCreateServiceMemberHandler = CreateServiceMemberHandler(context)
	internalAPI.ServiceMembersPatchServiceMemberHandler = PatchServiceMemberHandler(context)
	internalAPI.ServiceMembersShowServiceMemberHandler = ShowServiceMemberHandler(context)
	internalAPI.ServiceMembersShowServiceMemberOrdersHandler = ShowServiceMemberOrdersHandler(context)

	internalAPI.BackupContactsIndexServiceMemberBackupContactsHandler = IndexBackupContactsHandler(context)
	internalAPI.BackupContactsCreateServiceMemberBackupContactHandler = CreateBackupContactHandler(context)
	internalAPI.BackupContactsUpdateServiceMemberBackupContactHandler = UpdateBackupContactHandler(context)
	internalAPI.BackupContactsShowServiceMemberBackupContactHandler = ShowBackupContactHandler(context)

	internalAPI.DocumentsCreateDocumentHandler = CreateDocumentHandler(context)
	internalAPI.DocumentsShowDocumentHandler = ShowDocumentHandler(context)
	internalAPI.UploadsCreateUploadHandler = CreateUploadHandler(context)
	internalAPI.UploadsDeleteUploadHandler = DeleteUploadHandler(context)
	internalAPI.UploadsDeleteUploadsHandler = DeleteUploadsHandler(context)

	internalAPI.QueuesShowQueueHandler = ShowQueueHandler(context)

	internalAPI.ShipmentsCreateShipmentHandler = CreateShipmentHandler(context)
	internalAPI.ShipmentsPatchShipmentHandler = PatchShipmentHandler(context)
	internalAPI.ShipmentsGetShipmentHandler = GetShipmentHandler(context)

	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler(context)
	internalAPI.OfficeApprovePPMHandler = ApprovePPMHandler(context)
	internalAPI.OfficeApproveReimbursementHandler = ApproveReimbursementHandler(context)
	internalAPI.OfficeCancelMoveHandler = CancelMoveHandler(context)

	internalAPI.EntitlementsValidateEntitlementHandler = ValidateEntitlementHandler(context)

	return internalAPI.Serve(nil)
}

// NewOrdersAPIHandler returns a handler for the Orders API
func NewOrdersAPIHandler(context HandlerContext) http.Handler {

	// Wire up the handlers to the ordersAPIMux
	ordersSpec, err := loads.Analyzed(ordersapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	ordersAPI := ordersops.NewMymoveAPI(ordersSpec)
	ordersAPI.GetOrdersHandler = GetOrdersHandler(context)
	ordersAPI.IndexOrdersHandler = IndexOrdersHandler(context)
	ordersAPI.PostRevisionHandler = PostRevisionHandler(context)
	ordersAPI.PostRevisionToOrdersHandler = PostRevisionToOrdersHandler(context)
	return ordersAPI.Serve(nil)
}

// Converts the value returned by Pop's ValidateAnd* methods into a payload that can
// be returned to clients. This payload contains an object with a key,  `errors`, the
// value of which is a name -> validation error object.
func createFailedValidationPayload(verrs *validate.Errors) *internalmessages.InvalidRequestResponsePayload {
	errs := make(map[string]string)
	for _, key := range verrs.Keys() {
		errs[key] = strings.Join(verrs.Get(key), " ")
	}
	return &internalmessages.InvalidRequestResponsePayload{
		Errors: errs,
	}
}
