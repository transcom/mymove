package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestFetchOconusRateAreaByCityId() {
	usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99801")
	suite.NotNil(usprc)
	suite.FatalNoError(err)
	oconusRateArea, err := models.FetchOconusRateAreaByCityId(suite.DB(), usprc.ID.String())
	suite.NotNil(oconusRateArea)
	suite.NoError(err)
}
