package publicapi

import (
	"github.com/go-openapi/runtime/middleware"

	documentsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/documents"
	"github.com/transcom/mymove/pkg/handlers"
)

// CreateDocumentUploadHandler creates a new document upload
type CreateDocumentUploadHandler struct {
	handlers.HandlerContext
}

// Handle creates a new DocumentUpload from a request payload
func (h CreateDocumentUploadHandler) Handle(params documentsop.CreateDocumentUploadParams) middleware.Responder {
	return middleware.NotImplemented("operation .createDocumentUpload has not yet been implemented")
}
