package edisegment

func (suite *SegmentSuite) TestValidateL10() {
	validL10 := L10{
		Weight:          100.0,
		WeightQualifier: "B",
		WeightUnitCode:  "L",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validL10)
		suite.NoError(err)
	})

	suite.Run("validate failure", func() {
		l10 := L10{
			// Weight required
			WeightQualifier: "X",
			WeightUnitCode:  "X",
		}

		err := suite.validator.Struct(l10)
		suite.ValidateError(err, "Weight", "required")
		suite.ValidateError(err, "WeightQualifier", "eq")
		suite.ValidateError(err, "WeightUnitCode", "eq")
		suite.ValidateErrorLen(err, 3)
	})
}
