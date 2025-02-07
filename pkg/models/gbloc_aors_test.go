package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestFetchGblocAorsByJppsoCodeRateAreaDept() {
	usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99801")
	suite.NotNil(usprc)
	suite.FatalNoError(err)
	oconusRateArea, err := models.FetchOconusRateAreaByCityId(suite.DB(), usprc.ID.String())
	suite.NotNil(oconusRateArea)
	suite.NoError(err)

	jppsoRegion, err := models.FetchJppsoRegionByCode(suite.DB(), "MAPK")
	suite.NotNil(jppsoRegion)
	suite.NoError(err)

	gblocAors, err := models.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, models.DepartmentIndicatorARMY.String())
	suite.NotNil(gblocAors)
	suite.NoError(err)
}
