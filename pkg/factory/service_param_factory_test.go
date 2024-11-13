package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildServiceParam() {

	suite.Run("Successful creation of default service item param key", func() {
		// Under test:      BuildServiceParam
		// Mocked:          None
		// Set up:          Create a Service Paramwith no customizations or traits
		// Expected outcome:Service Paramshould be created with default values

		// CALL FUNCTION UNDER TEST
		serviceParam := BuildServiceParam(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.True(serviceParam.IsOptional)
		suite.False(serviceParam.ServiceID.IsNil())
		suite.False(serviceParam.ServiceItemParamKeyID.IsNil())
	})

	suite.Run("Successful creation of customized ServiceParam", func() {
		// Under test:      BuildServiceParam
		// Set up:          Create a Service Param and pass custom fields
		// Expected outcome:serviceParam should be created with custom fields
		// SETUP
		customServiceParam := models.ServiceParam{
			IsOptional: false,
		}

		reService := FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)

		serviceItemParamKey := BuildServiceItemParamKey(suite.DB(), []Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameIsPeak,
					Description: "custom desc",
				},
			},
		}, nil)

		// CALL FUNCTION UNDER TEST
		serviceParam := BuildServiceParam(suite.DB(), []Customization{
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    serviceItemParamKey,
				LinkOnly: true,
			},
			{
				Model: customServiceParam,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(reService.ID, serviceParam.ServiceID)
		suite.Equal(serviceItemParamKey.ID, serviceParam.ServiceItemParamKeyID)
		suite.Equal(customServiceParam.IsOptional, serviceParam.IsOptional)
	})

	suite.Run("Successful return of linkOnly ServiceParam", func() {
		// Under test:       BuildServiceParam
		// Set up:           Pass in a linkOnly serviceParam
		// Expected outcome: No new ServiceParam should be created.

		// Check num of ServiceParam records
		precount, err := suite.DB().Count(&models.ServiceParam{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		serviceParam := BuildServiceParam(suite.DB(), []Customization{
			{
				Model: models.ServiceParam{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)

		count, err := suite.DB().Count(&models.ServiceParam{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, serviceParam.ID)
	})

	suite.Run("Successful return of stubbed ServiceParam", func() {
		// Check num of ServiceParam records
		precount, err := suite.DB().Count(&models.ServiceParam{})
		suite.NoError(err)

		customServiceParam := models.ServiceParam{
			ServiceID:             uuid.Must(uuid.NewV4()),
			ServiceItemParamKeyID: uuid.Must(uuid.NewV4()),
			IsOptional:            false,
		}

		// Nil passed in as db
		serviceParam := BuildServiceParam(nil, []Customization{
			{
				Model: customServiceParam,
			},
		}, nil)

		count, err := suite.DB().Count(&models.ServiceParam{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.True(serviceParam.ID.IsNil())
		suite.Equal(customServiceParam.ServiceID, serviceParam.ServiceID)
		suite.Equal(customServiceParam.ServiceItemParamKeyID, serviceParam.ServiceItemParamKeyID)
	})

	suite.Run("Two service params should not be created", func() {
		// Under test:      FetchOrBuildServiceParam
		// Set up:          Create a service paramwith no customized state or traits
		// Expected outcome:Only 1 service item param key should be created
		count, potentialErr := suite.DB().Count(&models.ServiceParams{})
		suite.NoError(potentialErr)
		suite.Zero(count)

		firstServiceParam := FetchOrBuildServiceParam(suite.DB(), nil, nil)

		secondServiceParam := FetchOrBuildServiceParam(suite.DB(), nil, nil)

		suite.Equal(firstServiceParam.ID, secondServiceParam.ID)

		existingServiceParamCount, err := suite.DB().Count(&models.ServiceParams{})
		suite.NoError(err)
		suite.Equal(1, existingServiceParamCount)
	})
}
