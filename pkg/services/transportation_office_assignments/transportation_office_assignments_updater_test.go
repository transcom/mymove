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

type TransportationOfficeAssignmentsUpdaterServiceSuite struct {
	*testingsuite.PopTestSuite
	toaFetcher services.TransportationOfficeAssignmentFetcher
	toaUpdater services.TransportationOfficeAssignmentUpdater
}

func TestTransportationOfficeAssignmentsUpdaterServiceSuite(t *testing.T) {
	ts := &TransportationOfficeAssignmentsUpdaterServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *TransportationOfficeAssignmentsUpdaterServiceSuite) Test_UpdateTransportationOfficeAssignments_AddInitialAssignment() {
	suite.toaUpdater = NewTransportationOfficeAssignmentUpdater()

	officeUser := factory.BuildOfficeUserWithoutTransportationOfficeAssignment(suite.DB(), nil, nil)
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

	transportationOfficeAssignmentToAdd := models.TransportationOfficeAssignment{
		ID:                     officeUser.ID,
		TransportationOfficeID: transportationOffice.ID,
		PrimaryOffice:          models.BoolPointer(true),
	}
	toasToAdd := models.TransportationOfficeAssignments{transportationOfficeAssignmentToAdd}

	assignments, err := suite.toaUpdater.UpdateTransportationOfficeAssignments(suite.AppContextForTest(), officeUser.ID, toasToAdd)

	suite.NoError(err)
	suite.Equal(1, len(assignments))
	suite.Equal(officeUser.ID, assignments[0].ID)
	suite.Equal(true, *assignments[0].PrimaryOffice)
}

func (suite *TransportationOfficeAssignmentsUpdaterServiceSuite) Test_UpdateTransportationOfficeAssignments_AddAdditionalAssignment() {
	suite.toaFetcher = NewTransportationOfficeAssignmentFetcher()
	suite.toaUpdater = NewTransportationOfficeAssignmentUpdater()

	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

	existingAssignments, _ := suite.toaFetcher.FetchTransportationOfficeAssignmentsByOfficeUserID(suite.AppContextForTest(), officeUser.ID)

	transportationOfficeAssignmentToAdd := models.TransportationOfficeAssignment{
		ID:                     officeUser.ID,
		TransportationOfficeID: transportationOffice.ID,
		PrimaryOffice:          models.BoolPointer(false),
	}

	toasToAdd := models.TransportationOfficeAssignments{existingAssignments[0], transportationOfficeAssignmentToAdd}

	assignments, err := suite.toaUpdater.UpdateTransportationOfficeAssignments(suite.AppContextForTest(), officeUser.ID, toasToAdd)

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
	suite.Equal(transportationOffice.ID, assignments[secondaryAssignmentIndex].TransportationOfficeID)
}

func (suite *TransportationOfficeAssignmentsUpdaterServiceSuite) Test_UpdateTransportationOfficeAssignments_SwapPrimaryOffice() {
	suite.toaFetcher = NewTransportationOfficeAssignmentFetcher()
	suite.toaUpdater = NewTransportationOfficeAssignmentUpdater()

	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	secondaryTransportationOfficeAssignment := factory.BuildAlternateTransportationOfficeAssignment(suite.DB(), []factory.Customization{
		{Model: officeUser, LinkOnly: true},
	}, nil)

	existingAssignments, err := suite.toaFetcher.FetchTransportationOfficeAssignmentsByOfficeUserID(suite.AppContextForTest(), officeUser.ID)

	primaryAssignmentIndex := slices.IndexFunc(existingAssignments, func(toa models.TransportationOfficeAssignment) bool {
		return *toa.PrimaryOffice
	})
	secondaryAssignmentIndex := slices.IndexFunc(existingAssignments, func(toa models.TransportationOfficeAssignment) bool {
		return *toa.PrimaryOffice == false
	})

	suite.NoError(err)
	suite.Equal(2, len(existingAssignments))
	suite.Equal(officeUser.TransportationOfficeID, existingAssignments[primaryAssignmentIndex].TransportationOfficeID)
	suite.Equal(secondaryTransportationOfficeAssignment.TransportationOfficeID, existingAssignments[secondaryAssignmentIndex].TransportationOfficeID)

	existingAssignments[primaryAssignmentIndex].PrimaryOffice = models.BoolPointer(false)
	existingAssignments[secondaryAssignmentIndex].PrimaryOffice = models.BoolPointer(true)

	toasToAdd := models.TransportationOfficeAssignments{existingAssignments[secondaryAssignmentIndex], existingAssignments[primaryAssignmentIndex]}
	assignments, err := suite.toaUpdater.UpdateTransportationOfficeAssignments(suite.AppContextForTest(), officeUser.ID, toasToAdd)

	suite.NoError(err)
	suite.Equal(2, len(assignments))

	newPrimaryAssignmentIndex := slices.IndexFunc(assignments, func(toa models.TransportationOfficeAssignment) bool {
		return *toa.PrimaryOffice
	})
	newSecondaryAssignmentIndex := slices.IndexFunc(assignments, func(toa models.TransportationOfficeAssignment) bool {
		return *toa.PrimaryOffice == false
	})

	suite.Equal(existingAssignments[secondaryAssignmentIndex].TransportationOfficeID, assignments[newPrimaryAssignmentIndex].TransportationOfficeID)
	suite.Equal(existingAssignments[primaryAssignmentIndex].TransportationOfficeID, assignments[newSecondaryAssignmentIndex].TransportationOfficeID)
}
