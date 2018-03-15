package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-openapi/runtime/middleware"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

type Document struct {
	ID uuid.UUID
}

func payloadForDocumentModel(document Document) internalmessages.DocumentPayload {
	documentPayload := internalmessages.DocumentPayload{
		ID: fmtUUID(document.ID),
	}
	return documentPayload
}

// CreateDocumentHandler creates a new document via POST /issue
type CreateDocumentHandler HandlerContext

// Handle creates a new Document from a request payload
func (h CreateDocumentHandler) Handle(params documentop.CreateDocumentParams) middleware.Responder {
	// newDocument := Document{}

	file := params.File

	fmt.Printf("%s has a length of %d bytes.\n", file.Header.Filename, file.Header.Size)

	cwd, err := os.Getwd()
	if err != nil {
		h.logger.Error("Could not get cwd", zap.Error(err))
	}

	uploadsDir := filepath.Join(cwd, "uploads")
	if err = os.Mkdir(uploadsDir, 0777); err != nil {
		h.logger.Error("Could not make directory", zap.Error(err))
	}

	destinationPath := filepath.Join(uploadsDir, file.Header.Filename)
	destination, err := os.Create(destinationPath)
	defer destination.Close()

	if err != nil {
		h.logger.Error("Could on open file", zap.Error(err))
	}

	io.Copy(destination, file.Data)

	var response middleware.Responder
	// if _, err := h.db.ValidateAndCreate(&newDocument); err != nil {
	// 	h.logger.Error("DB Insertion", zap.Error(err))
	// 	response = documentop.NewCreateDocumentBadRequest()
	// } else {
	// 	documentPayload := payloadForDocumentModel(newDocument)
	// 	response = documentop.NewCreateDocumentCreated().WithPayload(&documentPayload)
	// }
	return response
}
