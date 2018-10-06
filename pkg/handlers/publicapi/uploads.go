package publicapi

import (
	"github.com/go-openapi/runtime/middleware"

	internaluploadop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/uploads"
	uploadop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
)

// TODO: These handlers just call the handlers for uploads in the internal API. This is not good.
// Ideally, the api should only live in one place and we should have the ability to call the public and the private
// api from our apps. Additionally, we are starting to investigate the model of using Services
// to separate logic from the handlers. This should be replaced at some point with those fixes.

// CreateUploadHandler creates a new upload via POST /documents/{documentID}/uploads
type CreateUploadHandler struct {
	handlers.HandlerContext
}

// Handle creates a new Upload from a request payload
func (h CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {

	internalUploadParams := internaluploadop.CreateUploadParams{
		HTTPRequest: params.HTTPRequest,
		DocumentID:  params.DocumentID,
		File:        params.File,
	}

	internalHandler := internalapi.CreateUploadHandler{HandlerContext: h.HandlerContext}
	return internalHandler.Handle(internalUploadParams)
}

// DeleteUploadHandler deletes an upload
type DeleteUploadHandler struct {
	handlers.HandlerContext
}

// Handle deletes an upload
func (h DeleteUploadHandler) Handle(params uploadop.DeleteUploadParams) middleware.Responder {
	internalDeleteParams := internaluploadop.DeleteUploadParams{
		HTTPRequest: params.HTTPRequest,
		UploadID:    params.UploadID,
	}

	internalHandler := internalapi.DeleteUploadHandler{HandlerContext: h.HandlerContext}
	return internalHandler.Handle(internalDeleteParams)
}

// DeleteUploadsHandler deletes a collection of uploads
type DeleteUploadsHandler struct {
	handlers.HandlerContext
}

// Handle deletes uploads
func (h DeleteUploadsHandler) Handle(params uploadop.DeleteUploadsParams) middleware.Responder {
	internalDeleteParams := internaluploadop.DeleteUploadsParams{
		HTTPRequest: params.HTTPRequest,
		UploadIds:   params.UploadIds,
	}

	internalHandler := internalapi.DeleteUploadsHandler{HandlerContext: h.HandlerContext}
	return internalHandler.Handle(internalDeleteParams)
}
