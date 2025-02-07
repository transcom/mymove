package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestFetchJppsoRegionByCode() {
	jppsoRegion, err := models.FetchJppsoRegionByCode(suite.DB(), "MAPK")
	suite.NotNil(jppsoRegion)
	suite.NoError(err)
	suite.Equal("USCG Base Ketchikan", jppsoRegion.Name)
}
