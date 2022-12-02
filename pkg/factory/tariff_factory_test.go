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
}

func (suite *FactorySuite) TestFetchOrBuildDefaultTariff400ngZip3() {
	defaultState := "CA"
	suite.Run("Successful creation of default Tariff400ngZip3", func() {
		// Under test:      FetchOrBuildDefaultTariff400ngZip3
		// Mocked:          None
		// Set up:          Use helper function FetchOrBuildDefaultTariff400ngZip3
		// Expected outcome:Tariff should be created with default values

		zip3Tariff := FetchOrBuildDefaultTariff400ngZip3(suite.DB())
		suite.Equal(defaultState, zip3Tariff.State)
	})
}
