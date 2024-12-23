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
	toaFetcher services.TransportaionOfficeAssignmentFetcher
	toaUpdater services.TransportaionOfficeAssignmentUpdater
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

func (suite *TransportationOfficeAssignmentsUpdaterServiceSuite) Test_UpdateTransportaionOfficeAssignments_AddInitialAssignment() {
	suite.toaUpdater = NewTransportaionOfficeAssignmentUpdater()

	officeUser := factory.BuildOfficeUserWithoutTransportationOfficeAssignment(suite.DB(), nil, nil)
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

	transportaionOfficeAssignmentToAdd := models.TransportationOfficeAssignment{
		ID:                     officeUser.ID,
		TransportationOfficeID: transportationOffice.ID,
		PrimaryOffice:          models.BoolPointer(true),
	}
	toasToAdd := models.TransportationOfficeAssignments{transportaionOfficeAssignmentToAdd}

	assignments, err := suite.toaUpdater.UpdateTransportaionOfficeAssignments(suite.AppContextForTest(), officeUser.ID, toasToAdd)

	suite.NoError(err)
	suite.Equal(1, len(assignments))
	suite.Equal(officeUser.ID, assignments[0].ID)
	suite.Equal(true, *assignments[0].PrimaryOffice)
}

func (suite *TransportationOfficeAssignmentsUpdaterServiceSuite) Test_UpdateTransportaionOfficeAssignments_AddAdditionalAssignment() {
	suite.toaFetcher = NewTransportaionOfficeAssignmentFetcher()
	suite.toaUpdater = NewTransportaionOfficeAssignmentUpdater()

	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

	existingAssignments, _ := suite.toaFetcher.FetchTransportaionOfficeAssignmentsByOfficeUserID(suite.AppContextForTest(), officeUser.ID)

	transportaionOfficeAssignmentToAdd := models.TransportationOfficeAssignment{
		ID:                     officeUser.ID,
		TransportationOfficeID: transportationOffice.ID,
		PrimaryOffice:          models.BoolPointer(false),
	}

	toasToAdd := models.TransportationOfficeAssignments{existingAssignments[0], transportaionOfficeAssignmentToAdd}

	assignments, err := suite.toaUpdater.UpdateTransportaionOfficeAssignments(suite.AppContextForTest(), officeUser.ID, toasToAdd)

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

func (suite *TransportationOfficeAssignmentsUpdaterServiceSuite) Test_UpdateTransportaionOfficeAssignments_SwapPrimaryOffice() {
	suite.toaFetcher = NewTransportaionOfficeAssignmentFetcher()
	suite.toaUpdater = NewTransportaionOfficeAssignmentUpdater()

	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	secondaryTransportationOfficeAssignment := factory.BuildAlternateTransportationOfficeAssignment(suite.DB(), []factory.Customization{
		{Model: officeUser, LinkOnly: true},
	}, nil)

	existingAssignments, err := suite.toaFetcher.FetchTransportaionOfficeAssignmentsByOfficeUserID(suite.AppContextForTest(), officeUser.ID)

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
	assignments, err := suite.toaUpdater.UpdateTransportaionOfficeAssignments(suite.AppContextForTest(), officeUser.ID, toasToAdd)

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
