package models_test

import (
	"github.com/gofrs/uuid"

	m "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicHHGAllowanceInstantiation() {

	newHHGAllowance := &m.HHGAllowance{
		PayGradeID:                    uuid.Must(uuid.NewV4()),
		TotalWeightSelf:               5000,
		TotalWeightSelfPlusDependents: 8000,
		ProGearWeight:                 2000,
		ProGearWeightSpouse:           500,
	}

	verrs, err := newHHGAllowance.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) TestEmptyHHGAllowanceInstantiation() {
	newHHGAllowance := m.HHGAllowance{}

	expErrors := map[string][]string{
		"pay_grade_id": {"PayGradeID can not be blank."},
	}

	suite.verifyValidationErrors(&newHHGAllowance, expErrors)
}

// Test validation fields that pass when empty but fail with faulty values
func (suite *ModelSuite) TestFaultyHHGAllowanceInstantiation() {
	newHHGAllowance := m.HHGAllowance{
		PayGradeID:                    uuid.Must(uuid.NewV4()),
		TotalWeightSelf:               -1,
		TotalWeightSelfPlusDependents: -1,
		ProGearWeight:                 -1,
		ProGearWeightSpouse:           -1,
	}

	expErrors := map[string][]string{
		"total_weight_self":                 {"-1 is not greater than -1."},
		"total_weight_self_plus_dependents": {"-1 is not greater than -1."},
		"pro_gear_weight":                   {"-1 is not greater than -1."},
		"pro_gear_weight_spouse":            {"-1 is not greater than -1."},
	}

	suite.verifyValidationErrors(&newHHGAllowance, expErrors)
}
