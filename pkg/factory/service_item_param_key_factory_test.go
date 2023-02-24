package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildServiceItemParamKey() {
	defaultKey := models.ServiceItemParamNameWeightEstimated
	defaultDescription := "test name weight estimated description"
	defaultType := models.ServiceItemParamTypeInteger
	defaultOrigin := models.ServiceItemParamOriginPrime

	suite.Run("Successful creation of default service item param key", func() {
		// Under test:      BuildServiceItemParamKey
		// Mocked:          None
		// Set up:          Create a Service Item Param Key with no customizations or traits
		// Expected outcome:Service Item Param Key should be created with default values

		// CALL FUNCTION UNDER TEST
		serviceItemParamKey := BuildServiceItemParamKey(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultKey, serviceItemParamKey.Key)
		suite.Equal(defaultDescription, serviceItemParamKey.Description)
		suite.Equal(defaultType, serviceItemParamKey.Type)
		suite.Equal(defaultOrigin, serviceItemParamKey.Origin)
	})

	suite.Run("Successful creation of customized ServiceItemParamKey", func() {
		// Under test:      BuildServiceItemParamKey
		// Set up:          Create a Service Item Param Key and pass custom fields
		// Expected outcome:serviceItemParamKey should be created with custom fields
		// SETUP
		customKey := models.ServiceItemParamNameContractYearName
		customType := models.ServiceItemParamTypeString
		customServiceItemParamKey := models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameContractYearName,
			Description: "test contract year name",
			Type:        models.ServiceItemParamTypeString,
		}

		// CALL FUNCTION UNDER TEST
		serviceItemParamKey := BuildServiceItemParamKey(suite.DB(), []Customization{
			{Model: customServiceItemParamKey},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customKey, serviceItemParamKey.Key)
		suite.Equal("test contract year name", serviceItemParamKey.Description)
		suite.Equal(customType, serviceItemParamKey.Type)
		suite.Equal(defaultOrigin, serviceItemParamKey.Origin)
	})

	suite.Run("Successful return of linkOnly ServiceItemParamKey", func() {
		// Under test:       BuildServiceItemParamKey
		// Set up:           Pass in a linkOnly serviceItemParamKey
		// Expected outcome: No new ServiceItemParamKey should be created.

		// Check num of ServiceItemParamKey records
		precount, err := suite.DB().Count(&models.ServiceItemParamKey{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		serviceItemParamKey := BuildServiceItemParamKey(suite.DB(), []Customization{
			{
				Model: models.ServiceItemParamKey{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)

		count, err := suite.DB().Count(&models.ServiceItemParamKey{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, serviceItemParamKey.ID)
	})

	suite.Run("Successful return of stubbed ServiceItemParamKey", func() {
		// Check num of ServiceItemParamKey records
		precount, err := suite.DB().Count(&models.ServiceItemParamKey{})
		suite.NoError(err)

		customKey := models.ServiceItemParamNameActualPickupDate

		// Nil passed in as db
		serviceItemParamKey := BuildServiceItemParamKey(nil, []Customization{
			{
				Model: models.ServiceItemParamKey{
					Key: models.ServiceItemParamNameActualPickupDate,
				},
			},
		}, nil)

		count, err := suite.DB().Count(&models.ServiceItemParamKey{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customKey, serviceItemParamKey.Key)
		suite.Equal(defaultDescription, serviceItemParamKey.Description)
		suite.Equal(defaultType, serviceItemParamKey.Type)
		suite.Equal(defaultOrigin, serviceItemParamKey.Origin)
	})

	suite.Run("Two service item param keys should not be created", func() {
		// Under test:      FetchOrBuildServiceItemParamKey
		// Set up:          Create a service item param key with no customized state or traits
		// Expected outcome:Only 1 service item param key should be created
		defaultKey := models.ServiceItemParamNameWeightEstimated
		count, potentialErr := suite.DB().Where(`key=$1`, defaultKey).Count(&models.ServiceItemParamKeys{})
		suite.NoError(potentialErr)
		suite.Zero(count)

		firstServiceItemParamKey := FetchOrBuildServiceItemParamKey(suite.DB(), nil, nil)

		secondServiceItemParamKey := FetchOrBuildServiceItemParamKey(suite.DB(), nil, nil)

		suite.Equal(firstServiceItemParamKey.ID, secondServiceItemParamKey.ID)

		existingServiceItemParamKeyCount, err := suite.DB().Where(`key=$1`, defaultKey).Count(&models.ServiceItemParamKeys{})
		suite.NoError(err)
		suite.Equal(1, existingServiceItemParamKeyCount)
	})
}
