package transportationofficeassignments

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type transportaionOfficeAssignmentUpdater struct {
}

// NewTransportaionOfficeAssignmentUpdater creates a new struct with the service dependencies
func NewTransportaionOfficeAssignmentUpdater() services.TransportaionOfficeAssignmentUpdater {
	return transportaionOfficeAssignmentUpdater{}
}

func (updater transportaionOfficeAssignmentUpdater) UpdateTransportaionOfficeAssignments(
	appCtx appcontext.AppContext,
	officeUserId uuid.UUID,
	newAssignments models.TransportationOfficeAssignments,
) (models.TransportationOfficeAssignments, error) {

	var existingAssignments models.TransportationOfficeAssignments
	err := appCtx.DB().Where("id = ?", officeUserId).All(&existingAssignments)

	if err != nil {
		return models.TransportationOfficeAssignments{}, err
	}

	var assignmentsToCreate models.TransportationOfficeAssignments

	for _, newAssignment := range newAssignments {
		newAssignment.ID = officeUserId
		assignmentsToCreate = append(assignmentsToCreate, newAssignment)
	}

	err = appCtx.DB().Destroy(existingAssignments)
	if err != nil {
		return models.TransportationOfficeAssignments{}, err
	}

	err = appCtx.DB().Create(assignmentsToCreate)
	if err != nil {
		return models.TransportationOfficeAssignments{}, err
	}

	var assignments models.TransportationOfficeAssignments
	err = appCtx.DB().Where("id = ?", officeUserId).All(&assignments)
	if err != nil {
		return models.TransportationOfficeAssignments{}, err
	}

	return assignments, nil
}
