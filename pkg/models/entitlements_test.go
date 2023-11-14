package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGetEntitlementWithValidValues() {
	E1 := models.ServiceMemberRankE1

	suite.Run("E1 with dependents", func() {
		E1FullLoad, err := models.GetEntitlement(E1)
		suite.NoError(err)
		suite.Assertions.Equal(8000, E1FullLoad.TotalWeightSelfPlusDependents)
	})

	suite.Run("E1 without dependents", func() {
		E1Solo, err := models.GetEntitlement(E1)
		suite.NoError(err)
		suite.Assertions.Equal(5000, E1Solo.TotalWeightSelf)
	})

	suite.Run("E1 Pro Gear", func() {
		E1ProGear, err := models.GetEntitlement(E1)
		suite.NoError(err)
		suite.Assertions.Equal(2000, E1ProGear.ProGearWeight)
	})

	suite.Run("E1 Pro Gear Spouse", func() {
		E1ProGearSpouse, err := models.GetEntitlement(E1)
		suite.NoError(err)
		suite.Assertions.Equal(500, E1ProGearSpouse.ProGearWeightSpouse)
	})
}
