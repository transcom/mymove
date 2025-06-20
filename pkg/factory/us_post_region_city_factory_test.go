package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestFetchOrBuildUsPostRegionCity() {
	suite.Run("Successful creation of default UsPostRegionCity", func() {
		// Under test:      FetchOrBuildUsPostRegionCity
		// Mocked:          None
		// Set up:          Create a us post region city with no customizations or traits
		// Expected outcome:UsPostRegionCity should be created with default values

		// SETUP
		defaultUsPostRegionCity := models.UsPostRegionCity{
			UsprZipID:          "90210",
			USPostRegionCityNm: "BEVERLY HILLS",
			UsprcCountyNm:      "LOS ANGELES",
			State:              "CA",
			CtryGencDgphCd:     "US",
			UsPostRegionId:     uuid.FromStringOrNil("5a6c650f-f4a9-428a-ae9d-20a251769dc5"),
			CityId:             uuid.FromStringOrNil("d684959a-f59c-4c05-b7c8-0a16df6718aa"),
		}

		// CALL FUNCTION UNDER TEST
		usPostRegionCity := FetchOrBuildUsPostRegionCity(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultUsPostRegionCity.UsprZipID, usPostRegionCity.UsprZipID)
		suite.Equal(defaultUsPostRegionCity.USPostRegionCityNm, usPostRegionCity.USPostRegionCityNm)
		suite.Equal(defaultUsPostRegionCity.UsprcCountyNm, usPostRegionCity.UsprcCountyNm)
		suite.Equal(defaultUsPostRegionCity.State, usPostRegionCity.State)
		suite.Equal(defaultUsPostRegionCity.CtryGencDgphCd, usPostRegionCity.CtryGencDgphCd)
		suite.Equal(defaultUsPostRegionCity.UsPostRegionId, usPostRegionCity.UsPostRegionId)
		suite.Equal(defaultUsPostRegionCity.CityId, usPostRegionCity.CityId)

	})

	suite.Run("Successful creation of customized UsPostRegionCity", func() {
		// Under test:      FetchOrBuildUsPostRegionCity
		// Set up:          Create or fetch a UsPostRegionCity and pass custom fields
		// Expected outcome:UsPostRegionCity should be created with custom fields
		// SETUP
		customUsPostRegionCity := models.UsPostRegionCity{
			ID:                 uuid.Must(uuid.NewV4()),
			UsprZipID:          "29229",
			USPostRegionCityNm: "New City",
			UsprcCountyNm:      "New County",
			State:              "SC",
			CtryGencDgphCd:     "US",
			UsPostRegionId:     uuid.FromStringOrNil("5a6c650f-f4a9-428a-ae9d-20a251769dc5"),
			CityId:             uuid.FromStringOrNil("d684959a-f59c-4c05-b7c8-0a16df6718aa"),
		}

		// CALL FUNCTION UNDER TEST
		usPostRegionCity := FetchOrBuildUsPostRegionCity(suite.DB(), []Customization{
			{Model: customUsPostRegionCity},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customUsPostRegionCity.ID, usPostRegionCity.ID)
		suite.Equal(customUsPostRegionCity.UsprZipID, usPostRegionCity.UsprZipID)
		suite.Equal(customUsPostRegionCity.USPostRegionCityNm, usPostRegionCity.USPostRegionCityNm)
		suite.Equal(customUsPostRegionCity.UsprcCountyNm, usPostRegionCity.UsprcCountyNm)
		suite.Equal(customUsPostRegionCity.State, usPostRegionCity.State)
		suite.Equal(customUsPostRegionCity.CtryGencDgphCd, usPostRegionCity.CtryGencDgphCd)
		suite.Equal(customUsPostRegionCity.UsPostRegionId, usPostRegionCity.UsPostRegionId)
		suite.Equal(customUsPostRegionCity.CityId, usPostRegionCity.CityId)
	})

	suite.Run("Successful return of linkOnly UsPostRegionCity", func() {
		// Under test:       FetchOrBuildUsPostRegionCity
		// Set up:           Pass in a linkOnly UsPostRegionCity
		// Expected outcome: No new UsPostRegionCity should be created.

		// Check num UsPostRegionCity records
		precount, err := suite.DB().Count(&models.UsPostRegionCity{})
		suite.NoError(err)

		expectedUsPostRegionCity := models.UsPostRegionCity{
			ID:                 uuid.Must(uuid.NewV4()),
			UsprZipID:          "29229",
			USPostRegionCityNm: "New City",
			UsprcCountyNm:      "New County",
			State:              "SC",
			CtryGencDgphCd:     "US",
			UsPostRegionId:     uuid.FromStringOrNil("5a6c650f-f4a9-428a-ae9d-20a251769dc5"),
			CityId:             uuid.FromStringOrNil("d684959a-f59c-4c05-b7c8-0a16df6718aa"),
		}
		usPostRegionCity := FetchOrBuildUsPostRegionCity(suite.DB(), []Customization{
			{
				Model:    expectedUsPostRegionCity,
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.UsPostRegionCity{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(expectedUsPostRegionCity.ID, usPostRegionCity.ID)
		suite.Equal(expectedUsPostRegionCity.UsprZipID, usPostRegionCity.UsprZipID)
		suite.Equal(expectedUsPostRegionCity.USPostRegionCityNm, usPostRegionCity.USPostRegionCityNm)
		suite.Equal(expectedUsPostRegionCity.UsprcCountyNm, usPostRegionCity.UsprcCountyNm)
		suite.Equal(expectedUsPostRegionCity.State, usPostRegionCity.State)
		suite.Equal(expectedUsPostRegionCity.CtryGencDgphCd, usPostRegionCity.CtryGencDgphCd)
		suite.Equal(expectedUsPostRegionCity.UsPostRegionId, usPostRegionCity.UsPostRegionId)
		suite.Equal(expectedUsPostRegionCity.CityId, usPostRegionCity.CityId)

	})

	suite.Run("Successful return of stubbed UsPostRegionCity", func() {
		// Under test:       FetchOrBuildUsPostRegionCity
		// Set up:           Pass in a linkOnly UsPostRegionCity
		// Expected outcome: No new UsPostRegionCity should be created.

		// Check num UsPostRegionCity records
		precount, err := suite.DB().Count(&models.UsPostRegionCity{})
		suite.NoError(err)

		expectedUsPostRegionCity := models.UsPostRegionCity{
			UsprZipID:          "29229",
			USPostRegionCityNm: "New City",
			UsprcCountyNm:      "New County",
			State:              "SC",
			CtryGencDgphCd:     "US",
			UsPostRegionId:     uuid.FromStringOrNil("5a6c650f-f4a9-428a-ae9d-20a251769dc5"),
			CityId:             uuid.FromStringOrNil("d684959a-f59c-4c05-b7c8-0a16df6718aa"),
		}

		// Nil passed in as db
		usPostRegionCity := FetchOrBuildUsPostRegionCity(nil, []Customization{
			{
				Model: expectedUsPostRegionCity,
			},
		}, nil)

		count, err := suite.DB().Count(&models.UsPostRegionCity{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(expectedUsPostRegionCity.UsprZipID, usPostRegionCity.UsprZipID)
		suite.Equal(expectedUsPostRegionCity.USPostRegionCityNm, usPostRegionCity.USPostRegionCityNm)
		suite.Equal(expectedUsPostRegionCity.UsprcCountyNm, usPostRegionCity.UsprcCountyNm)
		suite.Equal(expectedUsPostRegionCity.State, usPostRegionCity.State)
		suite.Equal(expectedUsPostRegionCity.CtryGencDgphCd, usPostRegionCity.CtryGencDgphCd)
		suite.Equal(expectedUsPostRegionCity.UsPostRegionId, usPostRegionCity.UsPostRegionId)
		suite.Equal(expectedUsPostRegionCity.CityId, usPostRegionCity.CityId)

	})
}
