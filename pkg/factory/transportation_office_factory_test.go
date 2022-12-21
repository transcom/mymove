package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildTransportationOffice() {
	suite.Run("Successful creation of default TransportationOffice", func() {
		// Under test:      BuildTransportationOffice
		// Mocked:          None
		// Set up:          Create a transportation office with no customizations or traits
		// Expected outcome:transportationOffice should be created with default values

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

		// CALL FUNCTION UNDER TEST
		transportationOffice := BuildTransportationOffice(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultOffice.Name, transportationOffice.Name)
		suite.Equal(defaultOffice.Gbloc, transportationOffice.Gbloc)
		suite.Equal(defaultOffice.Latitude, transportationOffice.Latitude)
		suite.Equal(defaultOffice.Longitude, transportationOffice.Longitude)
		suite.Equal((*string)(nil), transportationOffice.Hours)
		suite.Equal((*string)(nil), transportationOffice.Services)
		suite.Equal((*string)(nil), transportationOffice.Note)
		// Check that address was hooked in
		suite.Equal(defaultAddress.StreetAddress1, transportationOffice.Address.StreetAddress1)

	})

	suite.Run("Successful creation of customized TransportationOffice", func() {
		// Under test:      BuildTransportationOffice
		// Set up:          Create a Transportation Office and pass custom fields
		// Expected outcome:transportationOffice should be created with custom fields
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
		customPhoneLine := models.OfficePhoneLine{
			Number: "555-775-8829",
		}

		// CALL FUNCTION UNDER TEST
		transportationOffice := BuildTransportationOfficeWithPhoneLine(suite.DB(), []Customization{
			{Model: customOffice},
			{Model: customAddress},
			{Model: customPhoneLine},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customOffice.ID, transportationOffice.ID)
		suite.Equal(customOffice.Name, transportationOffice.Name)
		suite.Equal(customOffice.Gbloc, transportationOffice.Gbloc)
		suite.Equal(customOffice.Latitude, transportationOffice.Latitude)
		suite.Equal(customOffice.Longitude, transportationOffice.Longitude)
		suite.Equal(*customOffice.Note, *transportationOffice.Note)
		suite.Equal(*customOffice.Hours, *transportationOffice.Hours)
		suite.Equal(*customOffice.Services, *transportationOffice.Services)
		// MYTODO Check that address was customized .
		suite.Equal(customPhoneLine.Number, transportationOffice.PhoneLines[0].Number)
	})

	suite.Run("Successful return of linkOnly TransportationOffice", func() {
		// Under test:       BuildTransportationOffice
		// Set up:           Pass in a linkOnly transportationOffice
		// Expected outcome: No new TransportationOffice should be created.

		// Check num TransportationOffice records
		precount, err := suite.DB().Count(&models.TransportationOffice{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		office := BuildTransportationOffice(suite.DB(), []Customization{
			{
				Model: models.TransportationOffice{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.TransportationOffice{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, office.ID)

	})
	suite.Run("Successful return of stubbed TransportationOffice", func() {
		// Under test:       BuildTransportationOffice
		// Set up:           Pass in a linkOnly transportationOffice
		// Expected outcome: No new TransportationOffice should be created.

		// Check num TransportationOffice records
		precount, err := suite.DB().Count(&models.TransportationOffice{})
		suite.NoError(err)

		// Nil passed in as db
		office := BuildTransportationOffice(nil, []Customization{
			{
				Model: models.TransportationOffice{
					Name: "JPSSO Coronado",
				},
			},
		}, nil)
		count, err := suite.DB().Count(&models.TransportationOffice{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal("JPSSO Coronado", office.Name)
		suite.Equal("KKFA", office.Gbloc)

	})
}
