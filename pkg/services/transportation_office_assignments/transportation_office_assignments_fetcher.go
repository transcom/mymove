package transportationofficeassignments

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type transportaionOfficeAssignmentFetcher struct {
}

// NewTransportaionOfficeAssignmentUpdater creates a new struct with the service dependencies
func NewTransportaionOfficeAssignmentFetcher() services.TransportaionOfficeAssignmentFetcher {
	return transportaionOfficeAssignmentFetcher{}
}

func (fetcher transportaionOfficeAssignmentFetcher) FetchTransportaionOfficeAssignmentsByOfficeUserID(
	appCtx appcontext.AppContext,
	officeUserId uuid.UUID,
) (models.TransportationOfficeAssignments, error) {

	var transportationOfficeAssignments models.TransportationOfficeAssignments

	err := appCtx.DB().Q().EagerPreload("TransportationOffice").
		Join("transportation_offices", "transportation_office_assignments.transportation_office_id = transportation_offices.id").
		Where("transportation_office_assignments.id = ?", (officeUserId)).
		All(&transportationOfficeAssignments)

	if err != nil {
		return models.TransportationOfficeAssignments{}, err
	}

	return transportationOfficeAssignments, nil
}
