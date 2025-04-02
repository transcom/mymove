package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReServiceValidation() {
	suite.Run("test valid ReService", func() {
		validReService := models.ReService{
			Code: "123abc",
			Name: "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReService, expErrors, nil)
	})

	suite.Run("test empty ReService", func() {
		emptyReService := models.ReService{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReService, expErrors, nil)
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
