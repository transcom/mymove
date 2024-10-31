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
	defaultCounty := models.StringPointer("LOS ANGELES")

	customAddress1 := "101 This is Awesome Street"
	customAddress2 := models.StringPointer("Unit 2525")
	customAddress3 := models.StringPointer("c/o Another Person")
	customCity := "Modesto"
	customState := "ID"
	customPostalCode := "83725"
	customCounty := models.StringPointer("ADA")
	suite.Run("Successful creation of default address", func() {
		// Under test:      BuildAddress
		// Mocked:          None
		// Set up:          Create an Address with no customizations or traits
		// Expected outcome:Address should be created with default values

		address := BuildAddress(suite.DB(), nil, nil)

		country, err := models.FetchCountryByID(suite.DB(), *address.CountryId)
		suite.NoError(err)

		// VALIDATE RESULTS
		suite.Equal(defaultAddress1, address.StreetAddress1)
		suite.Equal("P.O. Box 12345", *address.StreetAddress2)
		suite.Equal("c/o Some Person", *address.StreetAddress3)
		suite.Equal(defaultCity, address.City)
		suite.Equal(defaultState, address.State)
		suite.Equal(defaultPostalCode, address.PostalCode)
		suite.Equal(country.ID, *address.CountryId)
		suite.Equal(defaultCounty, address.County)
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
					County:         customCounty,
					IsOconus:       models.BoolPointer(false),
				},
			},
			{
				Model: models.Country{
					Country:     "US",
					CountryName: "UNITED STATES",
				},
			},
		}, nil)

		country, err := models.FetchCountryByID(suite.DB(), *address.CountryId)
		suite.NoError(err)

		// VALIDATE RESULTS
		suite.Equal(customAddress1, address.StreetAddress1)
		suite.Equal(customAddress2, address.StreetAddress2)
		suite.Equal(customAddress3, address.StreetAddress3)
		suite.Equal(customCity, address.City)
		suite.Equal(customState, address.State)
		suite.Equal(customPostalCode, address.PostalCode)
		suite.Equal(country.ID, *address.CountryId)
		suite.Equal(customCounty, address.County)
	})

	suite.Run("Successful creation of an address with trait", func() {
		// Under test:      BuildAddress
		// Set up:          Create an Address with a trait
		// Expected outcome:Address should be created with custom StreetAddress1 and active status
		address := BuildAddress(suite.DB(), nil,
			[]Trait{
				GetTraitAddress2,
			})

		country, err := models.FetchCountryByID(suite.DB(), *address.CountryId)
		suite.NoError(err)
		// VALIDATE RESULTS
		suite.Equal("987 Any Avenue", address.StreetAddress1)
		suite.Equal("P.O. Box 9876", *address.StreetAddress2)
		suite.Equal("c/o Some Person", *address.StreetAddress3)
		suite.Equal("Fairfield", address.City)
		suite.Equal("CA", address.State)
		suite.Equal("94535", address.PostalCode)
		suite.Equal(country.ID, *address.CountryId)
		suite.Equal(models.StringPointer("SOLANO"), address.County)
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
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, []Trait{
			GetTraitAddress3,
		})

		country, err := models.FetchCountryByID(suite.DB(), *address.CountryId)
		suite.NoError(err)

		// VALIDATE RESULTS
		suite.Equal(customAddress1, address.StreetAddress1)
		suite.Equal(customAddress2, address.StreetAddress2)
		suite.Equal("c/o Another Person", *address.StreetAddress3)
		suite.Equal("Des Moines", address.City)
		suite.Equal("IA", address.State)
		suite.Equal("50309", address.PostalCode)
		suite.Equal(country.ID, *address.CountryId)
		suite.Equal(models.StringPointer("POLK"), address.County)
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
					IsOconus:       models.BoolPointer(false),
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
		suite.Equal(models.StringPointer("db nil when created"), address.County)

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
					County:         models.StringPointer("County"),
					IsOconus:       models.BoolPointer(false),
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
		suite.Equal(models.StringPointer("County"), address.County)
	})
}

func (suite *FactorySuite) TestBuildMinimalAddress() {
	defaultStreet := "N/A"
	defaultCity := "Fort Gorden"
	defaultState := "GA"
	defaultPostalCode := "30813"

	customStreet := "101 Custom Street"
	customPostalCode := "98765"

	suite.Run("Successful creation of default minimal address", func() {
		// Under test:      BuildMinimalAddress
		// Set up:          No customizations or traits provided
		// Expected outcome: Address should be created with default values

		address := BuildMinimalAddress(suite.DB(), nil, nil)

		country, err := models.FetchCountryByID(suite.DB(), *address.CountryId)
		suite.NoError(err)

		suite.Equal(defaultStreet, address.StreetAddress1)
		suite.Equal(defaultCity, address.City)
		suite.Equal(defaultState, address.State)
		suite.Equal(defaultPostalCode, address.PostalCode)
		suite.Equal(country.ID, *address.CountryId)
	})

	suite.Run("Successful creation of minimal address with customization", func() {
		// Under test:      BuildMinimalAddress
		// Set up:          Create a minimal address with a customized StreetAddress1 and PostalCode
		// Expected outcome: Address should be created with custom values

		address := BuildMinimalAddress(suite.DB(), []Customization{
			{
				Model: models.Address{
					StreetAddress1: customStreet,
					PostalCode:     customPostalCode,
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		country, err := models.FetchCountryByID(suite.DB(), *address.CountryId)
		suite.NoError(err)

		suite.Equal(customStreet, address.StreetAddress1)
		suite.Equal(defaultCity, address.City)
		suite.Equal(defaultState, address.State)
		suite.Equal(customPostalCode, address.PostalCode)
		suite.Equal(country.ID, *address.CountryId)
	})

	suite.Run("Successful creation of minimal address with trait", func() {
		// Under test:      BuildMinimalAddress
		// Set up:          Create a minimal address with a trait
		// Expected outcome: Address should be created with trait values

		address := BuildMinimalAddress(suite.DB(), nil, []Trait{
			GetTraitAddress2,
		})

		country, err := models.FetchCountryByID(suite.DB(), *address.CountryId)
		suite.NoError(err)

		suite.Equal("987 Any Avenue", address.StreetAddress1)
		suite.Equal("Fairfield", address.City)
		suite.Equal("CA", address.State)
		suite.Equal("94535", address.PostalCode)
		suite.Equal(country.ID, *address.CountryId)
	})

	suite.Run("Successful creation of stubbed address with customization", func() {
		// Under test:      BuildMinimalAddress
		// Set up:          Create a customized address, but don't pass in a db (stub)
		// Expected outcome: Address should be created, but not saved in the database

		precount, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)

		address := BuildMinimalAddress(nil, []Customization{
			{
				Model: models.Address{
					StreetAddress1: customStreet,
					PostalCode:     customPostalCode,
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		suite.Equal(customStreet, address.StreetAddress1)
		suite.Equal(customPostalCode, address.PostalCode)

		// Count should remain the same as the address is not saved in the DB
		count, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful creation of minimal address with LinkOnly customization", func() {
		// Under test:      BuildMinimalAddress
		// Set up:          Create an address with a link-only customization
		// Expected outcome: No new address should be created in the database.

		precount, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)

		address := BuildMinimalAddress(suite.DB(), []Customization{
			{
				Model: models.Address{
					ID:             uuid.Must(uuid.NewV4()),
					StreetAddress1: customStreet,
					PostalCode:     customPostalCode,
					IsOconus:       models.BoolPointer(false),
				},
				LinkOnly: true,
			},
		}, nil)

		count, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)
		suite.Equal(precount, count)

		suite.Equal(customStreet, address.StreetAddress1)
		suite.Equal(customPostalCode, address.PostalCode)
	})
}
