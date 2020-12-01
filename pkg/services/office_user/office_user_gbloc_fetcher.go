package officeuser

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserGblocFetcher struct {
	db *pop.Connection
}

// FetchGblocForOfficeUser fetches the GBLOC for the office user in the current session
func (f *officeUserGblocFetcher) FetchGblocForOfficeUser(officeUserID uuid.UUID) (string, error) {
	var transportationOffice models.TransportationOffice

	err := f.db.Q().
		Join("office_users", "transportation_offices.id = office_users.transportation_office_id").
		Where("office_users.id = ?", officeUserID).
		First(&transportationOffice)

	if err != nil {
		return "", fmt.Errorf("error fetching transportationOffice for officeUserID: %s with error %w", officeUserID, err)
	}

	gbloc := transportationOffice.Gbloc

	return gbloc, nil
}

// NewOfficeUserGblocFetcher returns an implementation of the OfficeUserGblocFetcher interface
func NewOfficeUserGblocFetcher(db *pop.Connection) services.OfficeUserGblocFetcher {
	return &officeUserGblocFetcher{db}
}
