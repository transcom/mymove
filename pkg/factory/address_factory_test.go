package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildAddress() {
	defaultAddress1 := "123 Any Street"
	defaultCity := "Beverly Hills"
	defaultState := "CA"
	defaultPostalCode := "90210"

	customAddress1 := "101 This is Awesome Street"
	customAddress2 := models.StringPointer("Unit 2525")
	customAddress3 := models.StringPointer("c/o Another Person")
	customCity := "Modesto"
	customState := "ID"
	customPostalCode := "83725"
	suite.Run("Successful creation of default address", func() {
		// Under test:      BuildAddress
		// Mocked:          None
		// Set up:          Create an Address with no customizations or traits
		// Expected outcome:Address should be created with default values

		address := BuildAddress(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultAddress1, address.StreetAddress1)
		suite.Equal("P.O. Box 12345", *address.StreetAddress2)
		suite.Equal("c/o Some Person", *address.StreetAddress3)
		suite.Equal(defaultCity, address.City)
		suite.Equal(defaultState, address.State)
		suite.Equal(defaultPostalCode, address.PostalCode)
		suite.Equal("US", *address.Country)
	})

	suite.Run("Successful creation of an address with customization", func() {
		// Under test:      BuildAddress
		// Set up:          Create an Address with a customized StreetAddress1 and no trait
		// Expected outcome:Address should be created with custom street address
		address := BuildAddress(suite.DB(), []Customization{
			{
				Model: models.Address{
					StreetAddress1: customAddress1,
					StreetAddress2: customAddress2,
					StreetAddress3: customAddress3,
					City:           customCity,
					State:          customState,
					PostalCode:     customPostalCode,
				},
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customAddress1, address.StreetAddress1)
		suite.Equal(customAddress2, address.StreetAddress2)
		suite.Equal(customAddress3, address.StreetAddress3)
		suite.Equal(customCity, address.City)
		suite.Equal(customState, address.State)
		suite.Equal(customPostalCode, address.PostalCode)
		suite.Equal(models.StringPointer("US"), address.Country)
	})

	suite.Run("Successful creation of an address with trait", func() {
		// Under test:      BuildAddress
		// Set up:          Create an Address with a trait
		// Expected outcome:Address should be created with custom StreetAddress1 and active status
		address := BuildAddress(suite.DB(), nil,
			[]Trait{
				GetTraitAddress2,
			})

		// VALIDATE RESULTS
		suite.Equal("987 Any Avenue", address.StreetAddress1)
		suite.Equal("P.O. Box 9876", *address.StreetAddress2)
		suite.Equal("c/o Some Person", *address.StreetAddress3)
		suite.Equal("Fairfield", address.City)
		suite.Equal("CA", address.State)
		suite.Equal("94535", address.PostalCode)
		suite.Equal("US", *address.Country)
	})

	suite.Run("Successful creation of address with both", func() {
		// Under test:      BuildAddress
		// Set up:          Create an Address with a customized StreetAddress1 and address trait
		// Expected outcome:Address should be created with email
		address := BuildAddress(suite.DB(), []Customization{
			{
				Model: models.Address{
					StreetAddress1: customAddress1,
					StreetAddress2: customAddress2,
				},
			},
		}, []Trait{
			GetTraitAddress3,
		})

		// VALIDATE RESULTS
		suite.Equal(customAddress1, address.StreetAddress1)
		suite.Equal(customAddress2, address.StreetAddress2)
		suite.Equal("c/o Another Person", *address.StreetAddress3)
		suite.Equal("Des Moines", address.City)
		suite.Equal("IA", address.State)
		suite.Equal("50309", address.PostalCode)
		suite.Equal("US", *address.Country)
	})

	suite.Run("Successful creation of stubbed address", func() {
		// Under test:      BuildAddress
		// Set up:          Create a customized address, but don't pass in a db
		// Expected outcome:Address should be created with email
		//                  No address should be created in database
		precount, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)

		address := BuildAddress(nil, []Customization{
			{
				Model: models.Address{
					StreetAddress1: customAddress1,
				},
			},
		}, []Trait{
			GetTraitAddress4,
		})

		// VALIDATE RESULTS
		suite.Equal(customAddress1, address.StreetAddress1)
		suite.Equal("P.O. Box 1234", *address.StreetAddress2)
		suite.Equal("c/o Another Person", *address.StreetAddress3)
		suite.Equal("Houston", address.City)
		suite.Equal("TX", address.State)
		suite.Equal("77083", address.PostalCode)
		suite.Equal("US", *address.Country)

		// Count how many addresses are in the DB, no new addresses should have been created
		count, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful creation of address with linked address", func() {
		// Under test:       BuildAddress
		// Set up:           Create an address and pass in a linkOnly address
		// Expected outcome: No new address should be created.

		// Check num addresses
		precount, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)

		address := BuildAddress(suite.DB(), []Customization{
			{
				Model: models.Address{
					ID:             uuid.Must(uuid.NewV4()),
					StreetAddress1: customAddress1,
					StreetAddress2: customAddress2,
					StreetAddress3: customAddress3,
					City:           customCity,
					State:          customState,
					PostalCode:     customPostalCode,
					Country:        models.StringPointer("Canada"),
				},
				LinkOnly: true,
			},
		}, nil)

		count, err := suite.DB().Count(&models.Address{})
		suite.Equal(precount, count)
		suite.NoError(err)

		// VALIDATE RESULTS
		suite.Equal(customAddress1, address.StreetAddress1)
		suite.Equal(customAddress2, address.StreetAddress2)
		suite.Equal(customAddress3, address.StreetAddress3)
		suite.Equal(customCity, address.City)
		suite.Equal(customState, address.State)
		suite.Equal(customPostalCode, address.PostalCode)
		suite.Equal("Canada", *address.Country)
	})
}
