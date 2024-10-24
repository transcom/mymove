package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// TransportaionOfficeAssignmentUpdater is the service object interface for UpdateTransportaionOfficeAssignments
//
//go:generate mockery --name TransportaionOfficeAssignmentUpdater
type TransportaionOfficeAssignmentUpdater interface {
	UpdateTransportaionOfficeAssignments(appCtx appcontext.AppContext, officeUserId uuid.UUID, transportationOfficeAssignments models.TransportationOfficeAssignments) (models.TransportationOfficeAssignments, error)
}
