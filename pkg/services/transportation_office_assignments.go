package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// TransportationOfficeAssignmentFetcher is the service object interface for FetchTransportationOfficeAssignmentsByOfficeUserID
//
//go:generate mockery --name TransportationOfficeAssignmentFetcher
type TransportationOfficeAssignmentFetcher interface {
	FetchTransportationOfficeAssignmentsByOfficeUserID(appCtx appcontext.AppContext, officeUserId uuid.UUID) (models.TransportationOfficeAssignments, error)
}

// TransportationOfficeAssignmentUpdater is the service object interface for UpdateTransportationOfficeAssignments
//
//go:generate mockery --name TransportationOfficeAssignmentUpdater
type TransportationOfficeAssignmentUpdater interface {
	UpdateTransportationOfficeAssignments(appCtx appcontext.AppContext, officeUserId uuid.UUID, transportationOfficeAssignments models.TransportationOfficeAssignments) (models.TransportationOfficeAssignments, error)
}
