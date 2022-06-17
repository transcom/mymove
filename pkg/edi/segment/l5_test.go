package edisegment

func (suite *SegmentSuite) TestValidateL5() {
	validL5 := L5{
		LadingLineItemNumber:   1,
		LadingDescription:      "DLH - Domestic Line Haul",
		CommodityCode:          "CCode",
		CommodityCodeQualifier: "D",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validL5)
		suite.NoError(err)
	})

	suite.Run("validate failure of lading description", func() {
		l5 := L5{
			LadingLineItemNumber: 1,
		}

		err := suite.validator.Struct(l5)
		suite.ValidateError(err, "LadingDescription", "required")
		suite.ValidateErrorLen(err, 1)
	})

	suite.Run("validate failure of missing CommodityCodeQualifier", func() {
		l5 := L5{
			LadingLineItemNumber: 1,
			LadingDescription:    "DLH - Domestic Line Haul",
			CommodityCode:        "CCode",
		}

		err := suite.validator.Struct(l5)
		suite.ValidateError(err, "CommodityCodeQualifier", "required_with")
		suite.ValidateErrorLen(err, 1)
	})

	suite.Run("validate failure of missing CommodityCode ", func() {
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
