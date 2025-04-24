// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
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
