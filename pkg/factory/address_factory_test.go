package factory

import "github.com/transcom/mymove/pkg/models"

func (suite *FactorySuite) TestBuildAddress() {
	defaultAddress1 := "123 Any Street"
	defaultPostalCode := "90210"
	customAddress1 := "101 This is Awesome Street"
	suite.Run("Successful creation of default address", func() {
		// Under test:      BuildAddress
		// Mocked:          None
		// Set up:          Create an Address with no customizations or traits
		// Expected outcome:Address should be created with default values

		address := BuildAddress(suite.DB(), nil, nil)
		suite.Equal(defaultAddress1, address.StreetAddress1)
		suite.Equal(defaultPostalCode, address.PostalCode)
	})

	suite.Run("Successful creation of an address with customization", func() {
		// Under test:      BuildAddress
		// Set up:          Create an Address with a customized StreetAddress1 and no trait
		// Expected outcome:Address should be created with custom street address
		address := BuildAddress(suite.DB(), []Customization{
			{
				Model: models.Address{
					StreetAddress1: customAddress1,
				},
			},
		}, nil)
		suite.Equal(customAddress1, address.StreetAddress1)
	})

	suite.Run("Successful creation of an address with trait", func() {
		address := BuildAddress(suite.DB(), nil,
			[]Trait{
				GetTraitAddress2,
			})
		suite.Equal("987 Any Avenue", address.StreetAddress1)
		suite.Equal("94535", address.PostalCode)
	})

	suite.Run("Successful creation of address with both", func() {
		address := BuildAddress(suite.DB(), []Customization{
			{
				Model: models.Address{
					StreetAddress1: customAddress1,
				},
			},
		}, []Trait{
			GetTraitAddress3,
		})

		suite.Equal(customAddress1, address.StreetAddress1)
		suite.Equal("50309", address.PostalCode)
	})

	suite.Run("Successful creation of stubbed address", func() {
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

		suite.Equal(customAddress1, address.StreetAddress1)
		suite.Equal("77083", address.PostalCode)
		// Count how many addresses are in the DB, no new addresses should have been created
		count, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})
}
