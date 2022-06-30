package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestReport() {
	suite.Run("blargh", func() {
		reports := models.EvaluationReports{}
		testdatagen.MakeEvaluationReport(suite.DB(), testdatagen.Assertions{})
		err := suite.DB().Q().All(&reports)
		suite.NoError(err)
	})
}
