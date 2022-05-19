package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	documentop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/documents"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

func payloadForDocumentModel(storer storage.FileStorer, document models.Document) (*internalmessages.DocumentPayload, error) {
	uploads := make([]*internalmessages.UploadPayload, len(document.UserUploads))
	for i, userUpload := range document.UserUploads {
		if userUpload.Upload.ID == uuid.Nil {
			return nil, errors.New("No uploads for user")
		}
		url, err := storer.PresignedURL(userUpload.Upload.StorageKey, userUpload.Upload.ContentType)
		if err != nil {
			return nil, err
		}

		uploadPayload := payloadForUploadModel(storer, userUpload.Upload, url)
		uploads[i] = uploadPayload
	}

	documentPayload := &internalmessages.DocumentPayload{
		ID:              handlers.FmtUUID(document.ID),
		ServiceMemberID: handlers.FmtUUID(document.ServiceMemberID),
		Uploads:         uploads,
	}
	return documentPayload, nil
}

// CreateDocumentHandler creates a new document via POST /documents/
type CreateDocumentHandler struct {
	handlers.HandlerConfig
}

// Handle creates a new Document from a request payload
func (h CreateDocumentHandler) Handle(params documentop.CreateDocumentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			serviceMemberID, err := uuid.FromString(params.DocumentPayload.ServiceMemberID.String())
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			// Fetch to check auth
			serviceMember, err := models.FetchServiceMemberForUser(appCtx.DB(), appCtx.Session(), serviceMemberID)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			newDocument := models.Document{
				ServiceMemberID: serviceMember.ID,
			}

			verrs, err := appCtx.DB().ValidateAndCreate(&newDocument)
			if err != nil {
				dbInsertErr := apperror.NewInternalServerError("DB Insertion")
				appCtx.Logger().Info(dbInsertErr.Error(), zap.Error(err))
				return documentop.NewCreateDocumentInternalServerError(), dbInsertErr
			} else if verrs.HasAny() {
				noSaveErr := apperror.NewInternalServerError("Could not save document")
				appCtx.Logger().Error(noSaveErr.Error(), zap.String("errors", verrs.Error()))
				return documentop.NewCreateDocumentBadRequest(), noSaveErr
			}

			appCtx.Logger().Info("created a document with id", zap.Any("new_document_id", newDocument.ID))
			documentPayload, err := payloadForDocumentModel(h.FileStorer(), newDocument)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			return documentop.NewCreateDocumentCreated().WithPayload(documentPayload), nil
		})
}

// ShowDocumentHandler shows a document via GETT /documents/:document_id
type ShowDocumentHandler struct {
	handlers.HandlerConfig
}

// Handle creates a new Document from a request payload
func (h ShowDocumentHandler) Handle(params documentop.ShowDocumentParams) middleware.Responder {
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

			return documentop.NewShowDocumentOK().WithPayload(documentPayload), nil
		})
}
