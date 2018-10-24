package document

import (
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
)

type fetchUploadService struct {
	docDB         models.DocumentDB
	fetchDocument services.FetchDocument
}

// NewFetchUploadService is the DI provider to create a FetchDocument service object
func NewFetchUploadService(docDB models.DocumentDB, fetchDocument services.FetchDocument) services.FetchUpload {
	return &fetchUploadService{
		docDB,
		fetchDocument,
	}
}

// FetchUpload returns an Upload if the user has access to that upload
func (s *fetchUploadService) Execute(session *server.Session, id uuid.UUID) (models.Upload, error) {
	upload, err := s.docDB.FetchUpload(id)
	if err != nil {
		return models.Upload{}, err
	}

	// If there's a document, check permissions. Otherwise user must have been the uploader
	if upload.DocumentID != nil {
		_, docErr := s.fetchDocument.Execute(session, *upload.DocumentID)
		if docErr != nil {
			return models.Upload{}, docErr
		}
	} else if upload.UploaderID != session.UserID {
		return models.Upload{}, errors.Wrap(services.ErrFetchForbidden, "user ID doesn't match uploader ID")
	}
	return *upload, nil
}
