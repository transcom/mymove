package transportationofficeassignments

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type transportationOfficeAssignmentUpdater struct {
}

// NewTransportationOfficeAssignmentUpdater creates a new struct with the service dependencies
func NewTransportationOfficeAssignmentUpdater() services.TransportationOfficeAssignmentUpdater {
	return transportationOfficeAssignmentUpdater{}
}

func (updater transportationOfficeAssignmentUpdater) UpdateTransportationOfficeAssignments(
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

	var assignments models.TransportationOfficeAssignments
	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err = appCtx.DB().Destroy(existingAssignments)
		if err != nil {
			return err
		}

		err = appCtx.DB().Create(assignmentsToCreate)
		if err != nil {
			return err
		}

		err = appCtx.DB().Where("id = ?", officeUserId).All(&assignments)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return models.TransportationOfficeAssignments{}, txErr
	}

	return assignments, nil
}
