package edisegment

func (suite *SegmentSuite) TestValidateLX() {
	validLX := LX{
		AssignedNumber: 1,
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validLX)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		lx := LX{
			AssignedNumber: 0, // min
		}

		err := suite.validator.Struct(lx)
		suite.ValidateError(err, "AssignedNumber", "min")
		suite.ValidateErrorLen(err, 1)
	})

	suite.Run("validate failure 2", func() {
		lx := LX{
			AssignedNumber: 1000000, // max
		}

		err := suite.validator.Struct(lx)
		suite.ValidateError(err, "AssignedNumber", "max")
		suite.ValidateErrorLen(err, 1)
	})
}
