package services

import (
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/server"
)

/*
FetchDocument is the interface for a service object(SO) to loads a Document from the database, applying
appropriate authorization check for the current session.

*/
type FetchDocument interface {
	// Execute ensures that the session passed in is authorized to access the details of the ServiceMember identified by id
	Execute(session *server.Session, id uuid.UUID) (models.Document, error)
}

/*
FetchUpload is the interface for a service object(SO) to loads an Upload from the database, applying
appropriate authorization check for the current session.
*/
type FetchUpload interface {
	// Execute ensures that the session passed in is authorized to access the details of the ServiceMember identified by id
	Execute(session *server.Session, id uuid.UUID) (models.Upload, error)
}
