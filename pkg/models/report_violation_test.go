package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestReportViolation() {
	suite.Run("Create and query a reportViolation successfully", func() {
		reportViolations := models.ReportViolations{}

		_, err := testdatagen.MakeReportViolation(suite.DB(), testdatagen.Assertions{})
		suite.NoError(err)

		err = suite.DB().All(&reportViolations)
		suite.NoError(err)
	})

}
