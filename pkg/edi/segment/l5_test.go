package edisegment

import (
	"testing"
)

func (suite *SegmentSuite) TestValidateL5() {
	validL5 := L5{
		LadingLineItemNumber:   1,
		LadingDescription:      "DLH - Domestic Line Haul",
		CommodityCode:          "CCode",
		CommodityCodeQualifier: "D",
	}

	suite.T().Run("validate success", func(t *testing.T) {
		err := suite.validator.Struct(validL5)
		suite.NoError(err)
	})

	suite.T().Run("validate failure of lading description", func(t *testing.T) {
		l5 := L5{
			LadingLineItemNumber: 1,
		}

		err := suite.validator.Struct(l5)
		suite.ValidateError(err, "LadingDescription", "required")
		suite.ValidateErrorLen(err, 1)
	})

	suite.T().Run("validate failure of missing CommodityCodeQualifier", func(t *testing.T) {
		l5 := L5{
			LadingLineItemNumber: 1,
			LadingDescription:    "DLH - Domestic Line Haul",
			CommodityCode:        "CCode",
		}

		err := suite.validator.Struct(l5)
		suite.ValidateError(err, "CommodityCodeQualifier", "required_with")
		suite.ValidateErrorLen(err, 1)
	})

	suite.T().Run("validate failure of missing CommodityCode ", func(t *testing.T) {
		l5 := L5{
			LadingLineItemNumber:   1,
			LadingDescription:      "DLH - Domestic Line Haul",
			CommodityCodeQualifier: "D",
		}

		err := suite.validator.Struct(l5)
		suite.ValidateError(err, "CommodityCode", "required_with")
		suite.ValidateErrorLen(err, 1)
	})
}
