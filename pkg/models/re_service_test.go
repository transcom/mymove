package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReServiceValidation() {
	suite.Run("test valid ReService", func() {
		validReService := models.ReService{
			Code: "123abc",
			Name: "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReService, expErrors)
	})

	suite.Run("test empty ReService", func() {
		emptyReService := models.ReService{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReService, expErrors)
	})
}

func (suite *ModelSuite) TestFetchReServiceBycode() {
	suite.Run("success - receive ReService when code is provided", func() {
		reService, err := models.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIHPK)
		suite.NoError(err)
		suite.NotNil(reService)
	})

	suite.Run("failure - receive error when code is not provided", func() {
		var blankReServiceCode models.ReServiceCode
		reService, err := models.FetchReServiceByCode(suite.DB(), blankReServiceCode)
		suite.Error(err)
		suite.Nil(reService)
		suite.Contains(err.Error(), "error fetching from re_services - required code not provided")
	})
}

func (suite *ModelSuite) TestIsDestinationRequest() {
	suite.Run("returns true when a service item is a destination request", func() {
		destinationSit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
				},
			},
		}, nil)

		destinationSITBool := models.IsDestinationRequest(destinationSit.ReService.Code)
		suite.True(destinationSITBool)
	})
	suite.Run("returns false when a service item is not a destination request", func() {
		originSit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		originSITBool := models.IsDestinationRequest(originSit.ReService.Code)
		suite.False(originSITBool)
	})
}
