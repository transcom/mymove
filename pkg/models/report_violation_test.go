package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestReportViolation() {
	suite.Run("Create and query a reportViolation successfully", func() {
		reportViolations := models.ReportViolations{}

		usprc, err := models.FindByZipCodeAndCity(suite.DB(), "90210", "Beverly Hills")
		suite.NoError(err)

		testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				UsPostRegionCityID: &usprc.ID,
			},
		})
		err = suite.DB().All(&reportViolations)
		suite.NoError(err)
	})

}
