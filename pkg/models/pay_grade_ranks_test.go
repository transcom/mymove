package models_test

import (
	m "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGetPayGradeRankDropdownOptions() {
	suite.Run("Fetch a affiliations Pay Grade/Ranks", func() {
		options, err := m.GetPayGradeRankDropdownOptions(suite.DB(), "ARMY")
		suite.NoError(err)

		suite.NotNil(options)
		suite.Equal(33, len(options))
	})
}
