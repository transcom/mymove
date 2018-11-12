package document

import (
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
)

type fetchDocumentService struct {
	docDB              models.DocumentDB
	fetchServiceMember services.FetchServiceMember
}

// NewFetchDocumentService is the DI provider to create a FetchDocument service object
func NewFetchDocumentService(docDB models.DocumentDB, fetchServiceMember services.FetchServiceMember) services.FetchDocument {
	return &fetchDocumentService{
		docDB,
		fetchServiceMember,
	}
}

// Execute fetches a document with appropriate authorization checks for the current session
func (s *fetchDocumentService) Execute(session *server.Session, id uuid.UUID) (models.Document, error) {
	// Load the document
	document, err := s.docDB.Fetch(id)
	if err != nil {
		return models.Document{}, err
	}

	// Now see if we have permissions to access the associated ServiceMember
	_, smErr := s.fetchServiceMember.Execute(document.ServiceMemberID, session)
	if smErr != nil {
		return models.Document{}, smErr
	}
	return *document, nil
}
