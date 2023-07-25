package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *ModelSuite) Test_AllFieldsOptionalCanSave() {
	loa := factory.BuildDefaultLineOfAccounting(suite.DB())

	verrs, err := suite.DB().ValidateAndSave(&loa)
	suite.NoVerrs(verrs)
	suite.NoError(err)
}

func (suite *ModelSuite) Test_AllFieldsPresentCanSave() {

	// TODO: use Factory
	loa := factory.BuildFullLineOfAccounting(suite.DB())

	verrs, err := suite.DB().ValidateAndSave(&loa)

	suite.NoVerrs(verrs)
	suite.NoError(err)
}
