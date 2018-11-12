package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGetEntitlementWithValidValues() {
	E1 := models.ServiceMemberRankE1

	// When: E1 has dependents and spouse gear
	E1FullLoad, err := models.GetEntitlement(E1, true, true)
	suite.Nil(err)
	suite.Assertions.Equal(10500, E1FullLoad)
	// When: E1 doesn't have dependents or spouse gear
	E1Solo, err := models.GetEntitlement(E1, false, false)
	suite.Nil(err)
	suite.Assertions.Equal(7000, E1Solo)
	// When: E1 doesn't have dependents but has spouse gear - impossible state
	E1FakeSpouse, err := models.GetEntitlement(E1, false, true)
	suite.Nil(err)
	suite.Assertions.Equal(7000, E1FakeSpouse)
	// When: E1 has dependents but no spouse gear
	E1DivorcedWithKids, err := models.GetEntitlement(E1, true, false)
	suite.Nil(err)
	suite.Assertions.Equal(10000, E1DivorcedWithKids)
}
