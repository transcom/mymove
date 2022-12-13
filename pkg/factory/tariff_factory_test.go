package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildTariff() {

	suite.Run("Successful creation of default Tariff Zip", func() {
		// Under test: BuildTariff400ngZip3
		// Mocked: none
		// Set up:      Create a Tariff400ngZip3 with no customizations or traits
		// Expected outcome: Tariff400ngZip3 should be created with default values
		zip3Tariff := BuildTariff400ngZip3(suite.DB(), nil, nil)

		defaultZip3 := models.Tariff400ngZip3{
			Zip3:          DefaultZip3,
			BasepointCity: "Beverly Hills",
			State:         "CA",
			ServiceArea:   "56",
			RateArea:      "US88",
			Region:        "2",
		}
		suite.Equal(defaultZip3.Zip3, zip3Tariff.Zip3)
		suite.Equal(defaultZip3.BasepointCity, zip3Tariff.BasepointCity)
		suite.Equal(defaultZip3.State, zip3Tariff.State)
		suite.Equal(defaultZip3.ServiceArea, zip3Tariff.ServiceArea)
		suite.Equal(defaultZip3.RateArea, zip3Tariff.RateArea)
		suite.Equal(defaultZip3.Region, zip3Tariff.Region)
	})

	suite.Run("Successful creation of Tariff with customization", func() {
		// Under test:      BuildTariff400ngZip3
		// Set up:          Create a Tariff400ngZip3 with a customized state and no trait
		// Expected outcome:Tariff400ngZip3 should be created with state
		customTariff := models.Tariff400ngZip3{
			Zip3:          "921",
			State:         "NV",
			BasepointCity: "San Diego",
			ServiceArea:   "27",
			RateArea:      "US22",
			Region:        "5",
		}
		zip3Tariff := BuildTariff400ngZip3(suite.DB(), []Customization{
			{Model: customTariff},
		}, nil)

		suite.Equal(customTariff.Zip3, zip3Tariff.Zip3)
		suite.Equal(customTariff.BasepointCity, zip3Tariff.BasepointCity)
		suite.Equal(customTariff.State, zip3Tariff.State)
		suite.Equal(customTariff.ServiceArea, zip3Tariff.ServiceArea)
		suite.Equal(customTariff.RateArea, zip3Tariff.RateArea)
		suite.Equal(customTariff.Region, zip3Tariff.Region)
	})

	suite.Run("Two tariffs should not be created", func() {
		// Under test:      FetchOrBuildTariff400ngZip3
		// Set up:          Create a Tariff400ngZip3 with no customized state or traits
		// Expected outcome:Only 1 Tariff400ngZip3 should be created
		count, potentialErr := suite.DB().Where("zip3 = ?", DefaultZip3).Count(&models.Tariff400ngZip3s{})
		suite.NoError(potentialErr)
		suite.Zero(count)

		zip3Tariff := FetchOrBuildTariff400ngZip3(suite.DB(), nil, nil)

		zip3Tariff1 := FetchOrBuildTariff400ngZip3(suite.DB(), nil, nil)

		suite.Equal(zip3Tariff.ID, zip3Tariff1.ID)

		existingZip3sCount, err := suite.DB().Where("zip3 = ?", DefaultZip3).Count(&models.Tariff400ngZip3s{})
		suite.NoError(err)
		suite.Equal(1, existingZip3sCount)
	})

	suite.Run("No tariff created if linkOnly passed in", func() {
		// Under test:      FetchOrBuildTariff400ngZip3
		// Set up:          Pass in a Tariff400ngZip3 with linkOnly
		// Expected outcome:No Tariff400ngZip3 should be created

		// Check zero tariffs
		count, err := suite.DB().Where("zip3 = ?", DefaultZip3).Count(&models.Tariff400ngZip3s{})
		suite.NoError(err)
		suite.Zero(count)

		// Create a new tariff (not in db)
		zip3 := models.Tariff400ngZip3{
			ID:            uuid.Must(uuid.NewV4()),
			Zip3:          DefaultZip3,
			BasepointCity: "Sacramento",
			State:         "CA",
		}

		// Pass in as linkOnly
		zip3Tariff := BuildTariff400ngZip3(suite.DB(), []Customization{
			{
				Model:    zip3,
				LinkOnly: true,
			},
		}, nil)

		// Check zero tariffs in DB still
		count, err = suite.DB().Where("zip3 = ?", DefaultZip3).Count(&models.Tariff400ngZip3s{})
		suite.NoError(err)
		suite.Zero(count)
		suite.Equal(zip3.BasepointCity, zip3Tariff.BasepointCity)

	})
}
