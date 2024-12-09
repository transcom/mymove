package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildPrimaryTransportationOfficeAssignment() {
	suite.Run("Successful creation of default Primary TransportationOfficeAssignment", func() {
		// Under test:      BuildPrimaryTransportationOfficeAssignment
		// Mocked:          None
		// Set up:          Create a transportation office assignment with no customizations or traits
		// Expected outcome:transportationOfficeAssignment should be created with default values

		// SETUP
		defaultOffice := models.TransportationOffice{
			Name:      "JPPSO Testy McTest",
			Gbloc:     "KKFA",
			Latitude:  1.23445,
			Longitude: -23.34455,
		}
		defaultAddress := models.Address{
			StreetAddress1: "123 Any Street",
		}
		defaultOfficeUser := models.OfficeUser{
			FirstName: "Leo",
			LastName:  "Spaceman",
			Telephone: "415-555-1212",
		}

		// CALL FUNCTION UNDER TEST
		// transportationOffice := BuildTransportationOffice(suite.DB(), nil, nil)
		// officeUser := BuildOfficeUserWithoutTransportationOfficeAssignment(suite.DB(), nil, nil)
		primaryTransportationOfficeAssignment := BuildPrimaryTransportationOfficeAssignment(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(models.BoolPointer(true), primaryTransportationOfficeAssignment.PrimaryOffice)

		// Associated Transportation Office is the default transportation office
		suite.Equal(defaultOffice.Name, primaryTransportationOfficeAssignment.TransportationOffice.Name)
		suite.Equal(defaultOffice.Gbloc, primaryTransportationOfficeAssignment.TransportationOffice.Gbloc)
		suite.Equal(defaultOffice.Latitude, primaryTransportationOfficeAssignment.TransportationOffice.Latitude)
		suite.Equal(defaultOffice.Longitude, primaryTransportationOfficeAssignment.TransportationOffice.Longitude)
		suite.Equal((*string)(nil), primaryTransportationOfficeAssignment.TransportationOffice.Hours)
		suite.Equal((*string)(nil), primaryTransportationOfficeAssignment.TransportationOffice.Services)
		suite.Equal((*string)(nil), primaryTransportationOfficeAssignment.TransportationOffice.Note)
		suite.Equal(defaultAddress.StreetAddress1, primaryTransportationOfficeAssignment.TransportationOffice.Address.StreetAddress1)

		// Associated Office User is the default Office User
		var matchingOfficeUser models.OfficeUser
		err := suite.DB().Find(&matchingOfficeUser, primaryTransportationOfficeAssignment.ID)
		suite.NoError(err)
		suite.Equal(defaultOfficeUser.FirstName, matchingOfficeUser.FirstName)
		suite.Equal(defaultOfficeUser.LastName, matchingOfficeUser.LastName)
		suite.Equal(defaultOfficeUser.Telephone, matchingOfficeUser.Telephone)
	})

	suite.Run("Successful creation of default Alternate TransportationOfficeAssignment", func() {
		// Under test:      BuildAlternateTransportationOfficeAssignment
		// Mocked:          None
		// Set up:          Create a transportation office assignment with no customizations or traits
		// Expected outcome:transportationOfficeAssignment should be created with default values

		// SETUP
		defaultOffice := models.TransportationOffice{
			Name:      "JPPSO Testy McTest",
			Gbloc:     "KKFA",
			Latitude:  1.23445,
			Longitude: -23.34455,
		}
		defaultAddress := models.Address{
			StreetAddress1: "123 Any Street",
		}
		defaultOfficeUser := models.OfficeUser{
			FirstName: "Leo",
			LastName:  "Spaceman",
			Telephone: "415-555-1212",
		}

		// CALL FUNCTION UNDER TEST
		primaryTransportationOfficeAssignment := BuildAlternateTransportationOfficeAssignment(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(models.BoolPointer(false), primaryTransportationOfficeAssignment.PrimaryOffice)

		// Associated Transportation Office is the default transportation office
		suite.Equal(defaultOffice.Name, primaryTransportationOfficeAssignment.TransportationOffice.Name)
		suite.Equal(defaultOffice.Gbloc, primaryTransportationOfficeAssignment.TransportationOffice.Gbloc)
		suite.Equal(defaultOffice.Latitude, primaryTransportationOfficeAssignment.TransportationOffice.Latitude)
		suite.Equal(defaultOffice.Longitude, primaryTransportationOfficeAssignment.TransportationOffice.Longitude)
		suite.Equal((*string)(nil), primaryTransportationOfficeAssignment.TransportationOffice.Hours)
		suite.Equal((*string)(nil), primaryTransportationOfficeAssignment.TransportationOffice.Services)
		suite.Equal((*string)(nil), primaryTransportationOfficeAssignment.TransportationOffice.Note)
		suite.Equal(defaultAddress.StreetAddress1, primaryTransportationOfficeAssignment.TransportationOffice.Address.StreetAddress1)

		// Associated Office User is the default Office User
		var matchingOfficeUser models.OfficeUser
		err := suite.DB().Find(&matchingOfficeUser, primaryTransportationOfficeAssignment.ID)
		suite.NoError(err)
		suite.Equal(defaultOfficeUser.FirstName, matchingOfficeUser.FirstName)
		suite.Equal(defaultOfficeUser.LastName, matchingOfficeUser.LastName)
		suite.Equal(defaultOfficeUser.Telephone, matchingOfficeUser.Telephone)
	})

	suite.Run("Successful creation of customized TransportationOfficeAssignments", func() {
		// Under test:      BuildPrimaryTransportationOfficeAssignment & BuildAlternateTransportationOfficeAssignment
		// Set up:          Create a Transportation Office Assignments and pass custom fields
		// Expected outcome:transportationOfficeAssignments should be created with custom fields
		// SETUP
		customOffice := models.TransportationOffice{
			ID:        uuid.Must(uuid.NewV4()),
			Name:      "JPPSO Coronado",
			Gbloc:     "CCRD",
			Latitude:  32.6806,
			Longitude: -117.1779,
		}

		secondaryCustomOffice := models.TransportationOffice{
			ID:        uuid.Must(uuid.NewV4()),
			Name:      "JPPSO Coronadon't",
			Gbloc:     "CCRD",
			Latitude:  33.6806,
			Longitude: -118.1779,
		}

		transportationOffice := BuildTransportationOfficeWithPhoneLine(suite.DB(), []Customization{
			{Model: customOffice},
		}, nil)

		secondaryTransportationOffice := BuildTransportationOffice(suite.DB(), []Customization{
			{Model: secondaryCustomOffice},
		}, nil)

		customOfficeUser := models.OfficeUser{
			ID:                     uuid.Must(uuid.NewV4()),
			FirstName:              "TOA",
			LastName:               "Tester",
			Email:                  "toa.tester@mail.com",
			TransportationOfficeID: transportationOffice.ID,
		}

		officeUser := BuildOfficeUserWithoutTransportationOfficeAssignment(suite.DB(), []Customization{
			{Model: customOfficeUser},
		}, nil)

		transportationOfficeAssignment := BuildPrimaryTransportationOfficeAssignment(suite.DB(), []Customization{
			{Model: officeUser, LinkOnly: true},
			{Model: transportationOffice, LinkOnly: true},
		}, nil)

		secondaryTransportationOfficeAssignment := BuildAlternateTransportationOfficeAssignment(suite.DB(), []Customization{
			{Model: officeUser, LinkOnly: true},
			{Model: secondaryTransportationOffice, LinkOnly: true},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customOffice.ID, transportationOfficeAssignment.TransportationOffice.ID)
		suite.Equal(customOffice.Name, transportationOfficeAssignment.TransportationOffice.Name)
		suite.Equal(customOffice.Gbloc, transportationOfficeAssignment.TransportationOffice.Gbloc)
		suite.Equal(customOfficeUser.ID, transportationOfficeAssignment.ID)
		suite.Equal(models.BoolPointer(true), transportationOfficeAssignment.PrimaryOffice)

		suite.Equal(secondaryCustomOffice.ID, secondaryTransportationOfficeAssignment.TransportationOffice.ID)
		suite.Equal(secondaryCustomOffice.Name, secondaryTransportationOfficeAssignment.TransportationOffice.Name)
		suite.Equal(secondaryCustomOffice.Gbloc, secondaryTransportationOfficeAssignment.TransportationOffice.Gbloc)
		suite.Equal(customOfficeUser.ID, secondaryTransportationOfficeAssignment.ID)
		suite.Equal(models.BoolPointer(false), secondaryTransportationOfficeAssignment.PrimaryOffice)

	})
}
