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
			HierarchicalIDNumber:       "A-123", // alphanum
			HierarchicalParentIDNumber: "1",     // isdefault
			HierarchicalLevelCode:      "XX",    // eq
		}

		err := suite.validator.Struct(hl)
		suite.ValidateError(err, "HierarchicalIDNumber", "alphanum")
		suite.ValidateError(err, "HierarchicalParentIDNumber", "isdefault")
		suite.ValidateError(err, "HierarchicalLevelCode", "oneof")
		suite.ValidateErrorLen(err, 3)
	})

	suite.T().Run("validate failure 2", func(t *testing.T) {
		hl := validHL
		hl.HierarchicalIDNumber = "" // alphanum takes precidence over min

		err := suite.validator.Struct(hl)
		suite.ValidateError(err, "HierarchicalIDNumber", "alphanum")
		suite.ValidateErrorLen(err, 1)
	})

	suite.T().Run("validate failure 3", func(t *testing.T) {
		hl := validHL
		hl.HierarchicalIDNumber = "0123456789ABCDF" // max

		err := suite.validator.Struct(hl)
		suite.ValidateError(err, "HierarchicalIDNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
