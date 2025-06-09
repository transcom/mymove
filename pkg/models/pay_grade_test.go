package models_test

import (
	m "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicPayGradeInstantiation() {
	newPayGrade := &m.PayGrade{
		Grade: "NewGrade",
	}

	verrs, err := newPayGrade.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) TestEmptyPayGradeInstantiation() {
	newPayGrade := m.PayGrade{}

	expErrors := map[string][]string{
		"grade": {"Grade can not be blank."},
	}
	suite.verifyValidationErrors(&newPayGrade, expErrors, nil)
}

func (suite *ModelSuite) TestGetPayGradesForAffiliation() {
	payGrades, err := m.GetPayGradesForAffiliation(suite.DB(), m.AffiliationAIRFORCE.String())
	suite.NoError(err)

	suite.Equal(26, len(payGrades))
}
