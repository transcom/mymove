package models_test

import (
	"testing"

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
	suite.T().Run("with no dependents authorized, TotalWeightSelf is AuthorizedWeight", func(t *testing.T) {
		entitlement := models.Entitlement{}
		entitlement.SetWeightAllotment("E_1")

		suite.Equal(entitlement.WeightAllotment().TotalWeightSelf, *entitlement.AuthorizedWeight())
	})

	suite.T().Run("with dependents authorized, TotalWeightSelfPlusDependents is AuthorizedWeight", func(t *testing.T) {
		dependentsAuthorized := true
		entitlement := models.Entitlement{DependentsAuthorized: &dependentsAuthorized}
		entitlement.SetWeightAllotment("E_1")

		suite.Equal(entitlement.WeightAllotment().TotalWeightSelfPlusDependents, *entitlement.AuthorizedWeight())
	})
}
