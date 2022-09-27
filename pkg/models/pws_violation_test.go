package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestPWSViolation() {
	suite.Run("Create and query a PWS violation successfully", func() {
		violations := models.PWSViolations{}
		testdatagen.MakePWSViolation(suite.DB(), testdatagen.Assertions{})

		err := suite.DB().All(&violations)
		suite.NoError(err)
	})
}
