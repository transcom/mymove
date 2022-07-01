package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestReport() {
	suite.Run("Create and query a report successfully", func() {
		reports := models.EvaluationReports{}
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		err := suite.DB().All(&reports)
		suite.NoError(err)
	})
}
