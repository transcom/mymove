package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildDutyLocation() {
	defaultAffiliation := internalmessages.AffiliationAIRFORCE

	suite.Run("test fetch system", func() {
		dutyLocation := FetchOrBuildCurrentDutyLocation(suite.DB())

		suite.Equal(dutyLocation.ID, dutyLocation.ID)
	})
	suite.Run("Successful creation of default duty location", func() {
		// Under test:      BuildDutyLocation
		// Mocked:          None
		// Set up:          Create a Duty Location with no customizations or traits
		// Expected outcome:Duty Location should be created with default values

		defaultOffice := models.TransportationOffice{
			Name:      "JPPSO Testy McTest",
			Gbloc:     "KKFA",
			Latitude:  1.23445,
			Longitude: -23.34455,
		}
		defaultAddress := models.Address{
			StreetAddress1: "987 Other Avenue",
		}

		// CALL FUNCTION UNDER TEST
		dutyLocation := BuildDutyLocation(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultAffiliation, *dutyLocation.Affiliation)
		// Check that address was hooked in
		suite.Equal(defaultAddress.StreetAddress1, dutyLocation.Address.StreetAddress1)
		// Check that transportation office was hooked in
		suite.Equal(defaultOffice.Name, dutyLocation.TransportationOffice.Name)
		suite.Equal(defaultOffice.Gbloc, dutyLocation.TransportationOffice.Gbloc)
		suite.Equal(defaultOffice.Latitude, dutyLocation.TransportationOffice.Latitude)
		suite.Equal(defaultOffice.Longitude, dutyLocation.TransportationOffice.Longitude)
	})

	suite.Run("Successful creation of customized DutyLocation", func() {
		// Under test:      BuiltDutyLocation
		// Set up:          Create a Duty Location and pass custom fields
		// Expected outcome:dutyLocation should be created with custom fields
		// SETUP
		customOffice := models.TransportationOffice{
			ID:        uuid.Must(uuid.NewV4()),
			Name:      "JPPSO Coronado",
			Gbloc:     "CCRD",
			Latitude:  32.6806,
			Longitude: -117.1779,
			Note:      models.StringPointer("Accessible to Public"),
			Hours:     models.StringPointer("9am-9pm"),
			Services:  models.StringPointer("CAC creation"),
		}

		customAddress := models.Address{
			StreetAddress1: "123 Any Street",
		}

		customAffiliation := internalmessages.AffiliationNAVY

		customDutyLocation := models.DutyLocation{
			ID:          uuid.Must(uuid.NewV4()),
			Affiliation: &customAffiliation,
		}

		// CALL FUNCTION UNDER TEST
		dutyLocation := BuildDutyLocation(suite.DB(), []Customization{
			{Model: customDutyLocation},
			{Model: customAddress},
			{Model: customOffice},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customDutyLocation.ID, dutyLocation.ID)
		// Check that the transportation office was customized
		suite.Equal(customOffice.ID, dutyLocation.TransportationOffice.ID)
		suite.Equal(customOffice.Name, dutyLocation.TransportationOffice.Name)
		suite.Equal(customOffice.Gbloc, dutyLocation.TransportationOffice.Gbloc)
		suite.Equal(customOffice.Latitude, dutyLocation.TransportationOffice.Latitude)
		suite.Equal(customOffice.Longitude, dutyLocation.TransportationOffice.Longitude)
		suite.Equal(*customOffice.Note, *dutyLocation.TransportationOffice.Note)
		suite.Equal(*customOffice.Hours, *dutyLocation.TransportationOffice.Hours)
		suite.Equal(*customOffice.Services, *dutyLocation.TransportationOffice.Services)
		// Check that the address was customized
		suite.Equal(customAddress.StreetAddress1, dutyLocation.Address.StreetAddress1)
	})

	suite.Run("Successful creation of duty location with customized addresses", func() {
		// Under test:      BuiltDutyLocation
		// Set up:          Create a Duty Location with unique dutylocation address and transportation office
		// Expected outcome:dutyLocation should be created with custom address different from address attached for TO

		// SETUP
		customDutyLocationAddress := models.Address{
			StreetAddress1: "123 Any Street",
		}

		customTransportationOfficeAddress := models.Address{
			StreetAddress1: "456 Something Street",
		}

		customAffiliation := internalmessages.AffiliationNAVY

		customDutyLocation := models.DutyLocation{
			ID:          uuid.Must(uuid.NewV4()),
			Affiliation: &customAffiliation,
		}

		// CALL FUNCTION UNDER TEST
		dutyLocation := BuildDutyLocation(suite.DB(), []Customization{
			{Model: customDutyLocation},
			{Model: customDutyLocationAddress, Type: &Addresses.DutyLocationAddress},
			{Model: customTransportationOfficeAddress, Type: &Addresses.DutyLocationTOAddress},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customDutyLocation.ID, dutyLocation.ID)
		suite.Equal(customAffiliation, *dutyLocation.Affiliation)
		// Check that the address was customized
		suite.Equal(customDutyLocationAddress.StreetAddress1, dutyLocation.Address.StreetAddress1)
		// Check that Transportation Office Address is different
		suite.NotEqual(dutyLocation.Address.StreetAddress1, dutyLocation.TransportationOffice.Address.StreetAddress1)
	})

	suite.Run("Successful return of linkOnly DutyLocation", func() {
		// Under test:       BuildDutyLocation
		// Set up:           Pass in a linkOnly dutyLocation
		// Expected outcome: No new DutyLocation should be created.

		// Check num of DutyLocation records
		precount, err := suite.DB().Count(&models.DutyLocation{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		dutyLocation := BuildDutyLocation(suite.DB(), []Customization{
			{
				Model: models.DutyLocation{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)

		count, err := suite.DB().Count(&models.DutyLocation{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, dutyLocation.ID)
	})

	suite.Run("Successful return of stubbed DutyLocation", func() {
		// Check num of DutyLocation records
		precount, err := suite.DB().Count(&models.DutyLocation{})
		suite.NoError(err)

		affiliation := internalmessages.AffiliationNAVY

		// Nil passed in as db
		dutyLocation := BuildDutyLocation(nil, []Customization{
			{
				Model: models.DutyLocation{
					Affiliation: &affiliation,
				},
			},
		}, nil)

		count, err := suite.DB().Count(&models.DutyLocation{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(affiliation, *dutyLocation.Affiliation)
	})

	suite.Run("Successful creation of Fort Gordon DutyLocation", func() {
		// Under test:       FetchOrMakeDefaultNewOrdersDutyLocation
		// Set up:           Count how many dutylocations and transportation office records we have
		//                   Call the function twice.
		// Expected outcome: The first time should create the appropriate new records, the second time should not.

		// Check num of DutyLocation records
		precountDL, err := suite.DB().Count(&models.DutyLocation{})
		suite.NoError(err)
		precountTO, err := suite.DB().Count(&models.TransportationOffice{})
		suite.NoError(err)

		// FUNCTION UNDER TEST
		dutyLocation := FetchOrBuildOrdersDutyLocation(suite.DB())

		// We should have one more TO and DL than we had before
		countDL, err := suite.DB().Count(&models.DutyLocation{})
		suite.NoError(err)
		suite.Equal(precountDL+1, countDL)
		countTO, err := suite.DB().Count(&models.TransportationOffice{})
		suite.NoError(err)
		suite.Equal(precountTO+1, countTO)
		suite.Equal("Fort Gordon", dutyLocation.Name)

		// Calling the function again
		dutyLocation2 := FetchOrBuildOrdersDutyLocation(suite.DB())

		// We should get the same record we created before
		suite.Equal("Fort Gordon", dutyLocation2.Name)
		suite.Equal(dutyLocation.ID, dutyLocation2.ID)
		// There should still only be one more of each record than we had at the start of the test.
		countDL, err = suite.DB().Count(&models.DutyLocation{})
		suite.NoError(err)
		suite.Equal(precountDL+1, countDL)
		countTO, err = suite.DB().Count(&models.TransportationOffice{})
		suite.NoError(err)
		suite.Equal(precountTO+1, countTO)

	})
}
