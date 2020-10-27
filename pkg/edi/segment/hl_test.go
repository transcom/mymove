package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateHL() {
	validHL := HL{
		HierarchicalIDNumber:  "303",
		HierarchicalLevelCode: "SS",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validHL)
		suite.NoError(err)
	})

	suite.T().Run("validate failure", func(t *testing.T) {
		hl := HL{
			HierarchicalIDNumber:       "300", // oneof
			HierarchicalParentIDNumber: "1",   // isdefault
			HierarchicalLevelCode:      "XX",  // eq
		}

		err := suite.validator.Struct(hl)
		suite.ValidateError(err, "HierarchicalIDNumber", "oneof")
		suite.ValidateError(err, "HierarchicalParentIDNumber", "isdefault")
		suite.ValidateError(err, "HierarchicalLevelCode", "oneof")
		suite.ValidateErrorLen(err, 3)
	})
}
