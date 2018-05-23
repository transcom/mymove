package handlers

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalapi"
	internalops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/storage"
)

// FileStorer is the set of methods needed to store and retrieve objects.
type FileStorer interface {
	Store(string, io.ReadSeeker, string) (*storage.StoreResult, error)
	Delete(string) error
	Key(...string) string
	PresignedURL(string, string) (string, error)
}

// HandlerContext contains dependencies that are shared between all handlers.
// Each individual handler is declared as a type alias for HandlerContext so that the Handle() method
// can be declared on it. When wiring up a handler, you can create a HandlerContext and cast it to the type you want.
type HandlerContext struct {
	db      *pop.Connection
	logger  *zap.Logger
	planner route.Planner
	storage FileStorer
}

// NewHandlerContext returns a new HandlerContext with its required private fields set.
func NewHandlerContext(db *pop.Connection, logger *zap.Logger) HandlerContext {
	return HandlerContext{
		db:     db,
		logger: logger,
	}
}

// SetFileStorer is a simple setter for storage private field
func (context *HandlerContext) SetFileStorer(storer FileStorer) {
	context.storage = storer
}

// SetPlanner is a simple setter for the route.Planner private field
func (context *HandlerContext) SetPlanner(planner route.Planner) {
	context.planner = planner
}

// NewPublicAPIHandler returns a handler for the public API
func NewPublicAPIHandler(context HandlerContext) http.Handler {

	// Wire up the handlers to the publicAPIMux
	apiSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	publicAPI := publicops.NewMymoveAPI(apiSpec)
	publicAPI.IndexTSPsHandler = TSPIndexHandler(context)
	publicAPI.TspShipmentsHandler = TSPShipmentsHandler(context)
	return publicAPI.Serve(nil)
}

// NewInternalAPIHandler returns a handler for the public API
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

	internalAPI.DutyStationsSearchDutyStationsHandler = SearchDutyStationsHandler(context)

	internalAPI.ShipmentsIndexShipmentsHandler = IndexShipmentsHandler(context)

	internalAPI.OrdersCreateOrdersHandler = CreateOrdersHandler(context)
	internalAPI.OrdersUpdateOrdersHandler = UpdateOrdersHandler(context)
	internalAPI.OrdersShowOrdersHandler = ShowOrdersHandler(context)

	internalAPI.MovesCreateMoveHandler = CreateMoveHandler(context)
	internalAPI.MovesPatchMoveHandler = PatchMoveHandler(context)
	internalAPI.MovesShowMoveHandler = ShowMoveHandler(context)
	internalAPI.MovesSubmitMoveForApprovalHandler = SubmitMoveHandler(context)

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

	internalAPI.OfficeApproveMoveHandler = ApproveMoveHandler(context)
	internalAPI.OfficeApprovePPMHandler = ApprovePPMHandler(context)

	return internalAPI.Serve(nil)
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
