package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalapi"
	internalops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations"
	"github.com/transcom/mymove/pkg/gen/restapi"
	publicops "github.com/transcom/mymove/pkg/gen/restapi/apioperations"
	"github.com/transcom/mymove/pkg/storage"
)

// HandlerContext contains dependencies that are shared between all handlers.
// Each individual handler is declared as a type alias for HandlerContext so that the Handle() method
// can be declared on it. When wiring up a handler, you can create a HandlerContext and cast it to the type you want.
type HandlerContext struct {
	db     *pop.Connection
	logger *zap.Logger
}

// NewHandlerContext returns a new HandlerContext with its private fields set.
func NewHandlerContext(db *pop.Connection, logger *zap.Logger) HandlerContext {
	return HandlerContext{
		db:     db,
		logger: logger,
	}
}

type fileStorer interface {
	Store(string, io.ReadSeeker, string) (*storage.StoreResult, error)
	Key(...string) string
	PresignedURL(string) (string, error)
}

// FileHandlerContext wraps a HandlerContext with an additional dependency for file
// manipulation
type FileHandlerContext struct {
	HandlerContext
	storage fileStorer
}

// NewFileHandlerContext returns a new FileHandlerContext with its private fields set.
func NewFileHandlerContext(context HandlerContext, storage fileStorer) FileHandlerContext {
	return FileHandlerContext{
		HandlerContext: context,
		storage:        storage,
	}
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

	internalAPI.IssuesCreateIssueHandler = CreateIssueHandler(context)
	internalAPI.IssuesIndexIssuesHandler = IndexIssuesHandler(context)

	internalAPI.Form1299sCreateForm1299Handler = CreateForm1299Handler(context)
	internalAPI.Form1299sIndexForm1299sHandler = IndexForm1299sHandler(context)
	internalAPI.Form1299sShowForm1299Handler = ShowForm1299Handler(context)

	internalAPI.CertificationCreateSignedCertificationHandler = CreateSignedCertificationHandler(context)

	internalAPI.PpmCreatePersonallyProcuredMoveHandler = CreatePersonallyProcuredMoveHandler(context)
	internalAPI.PpmIndexPersonallyProcuredMovesHandler = IndexPersonallyProcuredMovesHandler(context)
	internalAPI.PpmPatchPersonallyProcuredMoveHandler = PatchPersonallyProcuredMoveHandler(context)

	internalAPI.ShipmentsIndexShipmentsHandler = IndexShipmentsHandler(context)

	internalAPI.MovesCreateMoveHandler = CreateMoveHandler(context)
	internalAPI.MovesIndexMovesHandler = IndexMovesHandler(context)
	internalAPI.MovesPatchMoveHandler = PatchMoveHandler(context)
	return internalAPI.Serve(nil)
}
