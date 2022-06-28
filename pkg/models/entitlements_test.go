package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGetEntitlementWithValidValues() {
	E1 := models.ServiceMemberRankE1

	suite.Run("E1 with dependents", func() {
		E1FullLoad, err := models.GetEntitlement(E1, true)
		suite.NoError(err)
		suite.Assertions.Equal(8000, E1FullLoad)
	})

	suite.Run("E1 without dependents", func() {
		E1Solo, err := models.GetEntitlement(E1, false)
		suite.NoError(err)
		suite.Assertions.Equal(5000, E1Solo)
	})
}
