package edisegment

func (suite *SegmentSuite) TestValidateAK1() {
	validAK1 := AK1{
		FunctionalIdentifierCode: "SI",
		GroupControlNumber:       1234567,
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validAK1)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		ak1 := AK1{
			FunctionalIdentifierCode: "XX", // eq
			GroupControlNumber:       0,    // min
		}

		err := suite.validator.Struct(ak1)
		suite.ValidateError(err, "FunctionalIdentifierCode", "eq")
		suite.ValidateError(err, "GroupControlNumber", "min")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure 2", func() {
		ak1 := AK1{
			FunctionalIdentifierCode: "SI",
			GroupControlNumber:       999999999999999, // max
		}

		err := suite.validator.Struct(ak1)
		suite.ValidateError(err, "GroupControlNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
