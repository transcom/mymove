package public

import (
	"github.com/go-openapi/runtime/middleware"

	publicdocumentsop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/documents"
	"github.com/transcom/mymove/pkg/handlers/utils"
)

/* NOTE - The code above is for the INTERNAL API. The code below is for the public API. These will, obviously,
need to be reconciled. This will be done when the NotImplemented code below is Implemented
*/

// CreateDocumentUploadHandler creates a new document upload via POST /document/{document_uuid}/uploads
type CreateDocumentUploadHandler utils.HandlerContext

// Handle creates a new DocumentUpload from a request payload
func (h CreateDocumentUploadHandler) Handle(params publicdocumentsop.CreateDocumentUploadParams) middleware.Responder {
	return middleware.NotImplemented("operation .createDocumentUpload has not yet been implemented")
}
