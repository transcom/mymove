package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	documentop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ghc_documents"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForDocumentModel(storer storage.FileStorer, document models.Document) (*ghcmessages.DocumentPayload, error) {
	uploads := make([]*ghcmessages.Upload, len(document.UserUploads))
	for i, userUpload := range document.UserUploads {
		if userUpload.Upload.ID == uuid.Nil {
			return nil, errors.New("No uploads for user")
		}
		url, err := storer.PresignedURL(userUpload.Upload.StorageKey, userUpload.Upload.ContentType)
		if err != nil {
			return nil, err
		}

		uploadPayload := payloads.Upload(storer, userUpload.Upload, url)
		uploads[i] = uploadPayload
	}

	documentPayload := &ghcmessages.DocumentPayload{
		ID:              handlers.FmtUUID(document.ID),
		ServiceMemberID: handlers.FmtUUID(document.ServiceMemberID),
		Uploads:         uploads,
	}
	return documentPayload, nil
}

// GetDocumentHandler shows a document via GETT /documents/:document_id
type GetDocumentHandler struct {
	handlers.HandlerConfig
}

// Handle creates a new Document from a request payload
func (h GetDocumentHandler) Handle(params documentop.GetDocumentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			documentID, err := uuid.FromString(params.DocumentID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			document, err := models.FetchDocument(appCtx.DB(), appCtx.Session(), documentID, false)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			documentPayload, err := payloadForDocumentModel(h.FileStorer(), document)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			return documentop.NewGetDocumentOK().WithPayload(documentPayload), nil
		})
}
