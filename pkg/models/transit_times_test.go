package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_TransitDaysLookup() {
	days, err := TransitDays(unit.Pound(2500), 1100)
	suite.Nil(err)
	suite.Equal(11, days, "wrong number of days")

	days, err = TransitDays(unit.Pound(4300), 6100)
	suite.Nil(err)
	suite.Equal(30, days, "wrong number of days")
}

func (suite *ModelSuite) Test_TransitDaysLookupFail() {
	// Too much weight
	_, err := TransitDays(unit.Pound(100000), 2000)
	suite.Error(err)

	// Too many miles
	_, err = TransitDays(unit.Pound(2000), 8001)
	suite.Error(err)
}
