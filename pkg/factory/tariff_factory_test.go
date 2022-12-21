package factory

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildTariff() {
	defaultState := "CA"
	customState := "NV"
	suite.Run("Successful creation of default Tariff Zip", func() {
		// Under test: BuildTariff400ngZip3
		// Mocked: none
		// Set up:      Create a Tariff400ngZip3 with no customizations or traits
		// Expected outcome: Tariff400ngZip3 should be created with default values

		zip3Tariff := BuildTariff400ngZip3(suite.DB(), nil, nil)
		suite.Equal(defaultState, zip3Tariff.State)
	})

	suite.Run("Successful creation of Tariff with customization", func() {
		// Under test:      BuildTariff400ngZip3
		// Set up:          Create a Tariff400ngZip3 with a customized state and no trait
		// Expected outcome:Tariff400ngZip3 should be created with state
		zip3Tariff := BuildTariff400ngZip3(suite.DB(), []Customization{
			{
				Model: models.Tariff400ngZip3{
					State: customState,
				},
			},
		}, nil)

		suite.Equal(customState, zip3Tariff.State)
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
}
