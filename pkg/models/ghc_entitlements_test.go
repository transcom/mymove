package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestAuthorizedWeightWhenExistsInDB() {
	aw := 3000
	entitlement := models.Entitlement{DBAuthorizedWeight: &aw}

	err := suite.DB().Create(&entitlement)
	suite.NoError(err)

	suite.Equal(entitlement.DBAuthorizedWeight, entitlement.AuthorizedWeight())
}

func (suite *ModelSuite) TestAuthorizedWeightWhenNotInDBAndHaveWeightAllotment() {
	entitlement := models.Entitlement{}
	entitlement.SetWeightAllotment("E_1")

	err := suite.DB().Create(&entitlement)
	suite.NoError(err)

	suite.Equal(entitlement.WeightAllotment().TotalWeightSelf, *entitlement.AuthorizedWeight())
}
