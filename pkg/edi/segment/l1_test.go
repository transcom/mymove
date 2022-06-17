package edisegment

func (suite *SegmentSuite) TestValidateL1() {
	freightRate := float64(0)
	validL1 := L1{
		LadingLineItemNumber: 1,
		FreightRate:          &freightRate,
		RateValueQualifier:   "LB",
		Charge:               10000,
	}

	altValidL1 := L1{
		LadingLineItemNumber: 12,
		Charge:               10000,
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validL1)
		suite.NoError(err)
		suite.Equal([]string{"L1", "1", "0.00", "LB", "10000"}, validL1.StringArray())
	})

	suite.Run("validate alt success", func() {
		err := suite.validator.Struct(altValidL1)
		suite.NoError(err)
		suite.Equal([]string{"L1", "12", "", "", "10000"}, altValidL1.StringArray())
	})

	suite.Run("validate failure 1", func() {
		freightRate := float64(-1)
		l1 := L1{
			// LadingLineItemNumber:          // required
			FreightRate:        &freightRate, // min
			RateValueQualifier: "XX",         // oneof
			// Charge:                        // required
		}

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "required")
		suite.ValidateError(err, "FreightRate", "min")
		suite.ValidateError(err, "RateValueQualifier", "oneof")
		suite.ValidateError(err, "Charge", "required")
		suite.ValidateErrorLen(err, 4)
	})

	suite.Run("validate failure 2", func() {
		l1 := validL1
		l1.LadingLineItemNumber = 1000 // max
		l1.Charge = 1000000000000      // max

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "max")
		suite.ValidateError(err, "Charge", "max")
		suite.ValidateErrorLen(err, 2)
	})

	suite.Run("validate failure 3", func() {
		l1 := validL1
		l1.LadingLineItemNumber = -3 // min
		l1.RateValueQualifier = ""   // required_with
		l1.Charge = -1000000000000   // min

		err := suite.validator.Struct(l1)
		suite.ValidateError(err, "LadingLineItemNumber", "min")
		suite.ValidateError(err, "RateValueQualifier", "required_with")
		suite.ValidateError(err, "Charge", "min")
		suite.ValidateErrorLen(err, 3)
	})
}
