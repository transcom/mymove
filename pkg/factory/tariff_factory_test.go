package factory

import (
	"database/sql"

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
		var modelZip3s models.Tariff400ngZip3s
		defaultZip3 := DefaultZip3
		potentialErr := suite.DB().Where("zip3 = ?", defaultZip3).All(&modelZip3s)
		if potentialErr != nil && potentialErr != sql.ErrNoRows {
			suite.Fail("Failed to gather Tariff400ngZip3 rows from DB", potentialErr)
		}
		suite.Equal(0, len(modelZip3s))

		zip3Tariff := FetchOrBuildTariff400ngZip3(suite.DB(), nil, nil)

		zip3Tariff1 := FetchOrBuildTariff400ngZip3(suite.DB(), nil, nil)

		suite.Equal(zip3Tariff.ID, zip3Tariff1.ID)

		var existingZip3s models.Tariff400ngZip3s
		err := suite.DB().Where("zip3 = ?", defaultZip3).All(&existingZip3s)
		if err != nil && err != sql.ErrNoRows {
			suite.Fail("Failed to gather Tariff400ngZip3 rows from DB", err)
		}
		suite.Equal(1, len(existingZip3s))
	})
}
