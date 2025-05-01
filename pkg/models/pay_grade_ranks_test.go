package models_test

import (
	m "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGetPayGradeRankDropdownOptions() {
	testCases := [6]m.ServiceMemberAffiliation{
		m.ServiceMemberAffiliation(m.AffiliationARMY),
		m.ServiceMemberAffiliation(m.AffiliationNAVY),
		m.ServiceMemberAffiliation(m.AffiliationMARINES),
		m.ServiceMemberAffiliation(m.AffiliationAIRFORCE),
		m.ServiceMemberAffiliation(m.AffiliationCOASTGUARD),
		m.ServiceMemberAffiliation(m.AffiliationSPACEFORCE),
	}
	for _, testCase := range testCases {
		suite.Run("No errors for all affiliations", func() {
			options, err := m.GetPayGradeRankDropdownOptions(suite.DB(), string(testCase))
			suite.NoError(err)

			suite.Greater(len(options), 0)
		})
	}
	suite.Run("Fetch a affiliations Pay Grade/Ranks", func() {
		options, err := m.GetPayGradeRankDropdownOptions(suite.DB(), string(m.AffiliationARMY))
		suite.NoError(err)

		suite.NotNil(options)
		suite.Equal(33, len(options))
	})
}
