package transportationofficeassignments

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type TransportationOfficeAssignmentsServiceSuite struct {
	*testingsuite.PopTestSuite
	toaFetcher services.TransportaionOfficeAssignmentFetcher
}

func TestTransportationOfficeAssignmentsServiceSuite(t *testing.T) {
	ts := &TransportationOfficeAssignmentsServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *TransportationOfficeAssignmentsServiceSuite) Test_FetchTransportaionOfficeAssignmentsByOfficeUserID() {
	suite.toaFetcher = NewTransportaionOfficeAssignmentFetcher()
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	assignments, err := suite.toaFetcher.FetchTransportaionOfficeAssignmentsByOfficeUserID(suite.AppContextForTest(), officeUser.ID)

	suite.NoError(err)
	suite.Equal(1, len(assignments))
	suite.Equal(officeUser.ID, assignments[0].ID)
}
