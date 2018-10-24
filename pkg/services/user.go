package services

import (
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/server"
)

/*
FetchServiceMember is the interface for a service object(SO) to loads a Service memeber from the database, applying
appropriate authorization checks for the session

*/
type FetchServiceMember interface {
	// Execute ensures that the session passed in is authorized to access the details of the ServiceMember identified by id
	Execute(session *server.Session, id uuid.UUID) (*models.ServiceMember, error)
}
