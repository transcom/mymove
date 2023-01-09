package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildOfficePhoneLine() {
	suite.Run("Successful creation of default OfficePhoneLine", func() {
		// Under test:      BuildOfficePhoneLine
		// Mocked:          None
		// Set up:          Create a line with no customizations or traits
		// Expected outcome:officePhoneLine should be created with default values

		// SETUP
		defaultPhone := models.OfficePhoneLine{
			Number:      "(510) 555-5555",
			IsDsnNumber: false,
			Type:        "voice",
		}
		defaultOffice := models.TransportationOffice{
			Name: "JPPSO Testy McTest",
		}

		// CALL FUNCTION UNDER TEST
		officePhoneLine := BuildOfficePhoneLine(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultPhone.Number, officePhoneLine.Number)
		suite.Equal(defaultPhone.IsDsnNumber, officePhoneLine.IsDsnNumber)
		suite.Equal(defaultPhone.Type, officePhoneLine.Type)
		suite.Equal((*string)(nil), officePhoneLine.Label)
		suite.Equal("voice", officePhoneLine.Type)
		// Check that office was hooked in
		suite.Equal(defaultOffice.Name, officePhoneLine.TransportationOffice.Name)

	})

	suite.Run("Successful creation of customized OfficePhoneLine", func() {
		// Under test:       BuildOfficePhoneLine
		// Set up:           Create a phoneLine Office and pass custom fields,
		//                   including a transportationOffice
		// Expected outcome: officePhoneLine and transportationOffice
		//                   should be created with custom fields
		// SETUP
		customPhone := models.OfficePhoneLine{
			ID:          uuid.Must(uuid.NewV4()),
			Number:      "555-7758829",
			Label:       models.StringPointer("Newman"),
			IsDsnNumber: true,
			Type:        "fax",
		}
		customOffice := models.TransportationOffice{
			Name: "NEX National City",
		}

		// CALL FUNCTION UNDER TEST
		officePhoneLine := BuildOfficePhoneLine(suite.DB(), []Customization{
			{Model: customOffice},
			{Model: customPhone},
		}, nil)

		suite.Equal(customPhone.Number, officePhoneLine.Number)
		suite.Equal(customPhone.IsDsnNumber, officePhoneLine.IsDsnNumber)
		suite.Equal(customPhone.Type, officePhoneLine.Type)
		suite.Equal(*customPhone.Label, *officePhoneLine.Label)
		suite.Equal("fax", officePhoneLine.Type)
		// Check that office was hooked in
		suite.Equal(customOffice.Name, officePhoneLine.TransportationOffice.Name)
		// Check that the transportationOffice was customized

	})

	suite.Run("Successful return of linkOnly OfficePhoneLine", func() {
		// Under test:       BuildOfficePhoneLine
		// Set up:           Pass in a linkOnly officePhoneLine
		// Expected outcome: No new OfficePhoneLine should be created.

		// Check num OfficePhoneLine records
		precount, err := suite.DB().Count(&models.OfficePhoneLine{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		office := BuildOfficePhoneLine(suite.DB(), []Customization{
			{
				Model: models.OfficePhoneLine{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.OfficePhoneLine{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, office.ID)

	})
	suite.Run("Successful return of stubbed OfficePhoneLine", func() {
		// Under test:       BuildOfficePhoneLine
		// Set up:           Pass in a linkOnly officePhoneLine
		// Expected outcome: No new OfficePhoneLine should be created.

		// Check num OfficePhoneLine records
		precount, err := suite.DB().Count(&models.OfficePhoneLine{})
		suite.NoError(err)

		// Nil passed in as db
		phone := BuildOfficePhoneLine(nil, []Customization{
			{
				Model: models.OfficePhoneLine{
					Label: models.StringPointer("JPSSO Coronado"),
				},
			},
		}, nil)
		count, err := suite.DB().Count(&models.OfficePhoneLine{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal("JPSSO Coronado", *phone.Label)
		suite.Equal("voice", phone.Type)

	})
}
