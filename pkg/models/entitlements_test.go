package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGetEntitlementWithValidValues() {
	E1 := models.ServiceMemberRankE1

	// When: E1 has dependents and spouse gear
	E1FullLoad, _ := models.GetEntitlement(E1, true, true)
	suite.Assertions.Equal(10500, E1FullLoad)
	// When: E1 doesn't have dependents or spouse gear
	E1Solo, _ := models.GetEntitlement(E1, false, false)
	suite.Assertions.Equal(7000, E1Solo)
	// When: E1 doesn't have dependents but has spouse gear - impossible state
	E1FakeSpouse, _ := models.GetEntitlement(E1, false, true)
	suite.Assertions.Equal(7000, E1FakeSpouse)
	// When: E1 has dependents but no spouse gear
	E1DivorcedWithKids, _ := models.GetEntitlement(E1, true, false)
	suite.Assertions.Equal(10000, E1DivorcedWithKids)
}
