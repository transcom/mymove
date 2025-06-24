package transportationofficeassignments

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type TransportationOfficeAssignmentsFetcherServiceSuite struct {
	*testingsuite.PopTestSuite
	toaFetcher services.TransportationOfficeAssignmentFetcher
}

func TestTransportationOfficeAssignmentsFetcherServiceSuite(t *testing.T) {
	ts := &TransportationOfficeAssignmentsFetcherServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *TransportationOfficeAssignmentsFetcherServiceSuite) Test_FetchTransportationOfficeAssignmentByOfficeUserID() {
	suite.toaFetcher = NewTransportationOfficeAssignmentFetcher()

	// Creating an office user requires creating a transportation office assignment and we will need the office user's ID
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	assignments, err := suite.toaFetcher.FetchTransportationOfficeAssignmentsByOfficeUserID(suite.AppContextForTest(), officeUser.ID)

	suite.NoError(err)
	suite.Equal(1, len(assignments))
	suite.Equal(officeUser.ID, assignments[0].ID)
	suite.Equal(true, *assignments[0].PrimaryOffice)
}

func (suite *TransportationOfficeAssignmentsFetcherServiceSuite) Test_FetchTransportationOfficeAssignmentsByOfficeUserID() {
	suite.toaFetcher = NewTransportationOfficeAssignmentFetcher()

	// Creating an office user requires creating a transportation office assignment and we will need the office user's ID
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	secondaryTransportationOfficeAssignment := factory.BuildAlternateTransportationOfficeAssignment(suite.DB(), []factory.Customization{
		{Model: officeUser, LinkOnly: true},
	}, nil)
	assignments, err := suite.toaFetcher.FetchTransportationOfficeAssignmentsByOfficeUserID(suite.AppContextForTest(), officeUser.ID)

	suite.NoError(err)
	suite.Equal(2, len(assignments))
	suite.Equal(officeUser.ID, assignments[0].ID)
	suite.Equal(officeUser.ID, assignments[1].ID)

	primaryAssignmentIndex := slices.IndexFunc(assignments, func(toa models.TransportationOfficeAssignment) bool {
		return *toa.PrimaryOffice
	})
	secondaryAssignmentIndex := slices.IndexFunc(assignments, func(toa models.TransportationOfficeAssignment) bool {
		return *toa.PrimaryOffice == false
	})

	suite.Equal(officeUser.TransportationOfficeID, assignments[primaryAssignmentIndex].TransportationOfficeID)
	suite.Equal(secondaryTransportationOfficeAssignment.TransportationOfficeID, assignments[secondaryAssignmentIndex].TransportationOfficeID)
}
