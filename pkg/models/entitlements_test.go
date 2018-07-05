package models_test

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGetEntitlementWithValidValues() {
	E1 := internalmessages.ServiceMemberRankE1

	// When: E1 has dependents and spouse gear
	suite.Assertions.Equal(10500, models.GetEntitlement(E1, true, true))
	// When: E1 doesn't have dependents or spouse gear
	suite.Assertions.Equal(7000, models.GetEntitlement(E1, false, false))
	// When: E1 doesn't have dependents but has spouse gear - impossible state
	suite.Assertions.Equal(7000, models.GetEntitlement(E1, false, true))
	// When: E1 has dependents but no spouse gear
	suite.Assertions.Equal(10000, models.GetEntitlement(E1, true, false))
}
