package edisegment

func (suite *SegmentSuite) TestValidateG62() {
	validG62ActualPickupDateTime := G62{
		DateQualifier: 86,
		Date:          "20200909",
		TimeQualifier: 8,
		Time:          "1617",
	}
	validG62RequestedPickupDateTime := G62{
		DateQualifier: 10,
		Date:          "20200909",
		TimeQualifier: 5,
		Time:          "1617",
	}
	validG62ScheduledPickupDate := G62{
		DateQualifier: 76,
		Date:          "20200909",
	}

	suite.Run("validate success", func() {
		err := suite.validator.Struct(validG62ActualPickupDateTime)
		suite.NoError(err)
		err = suite.validator.Struct(validG62RequestedPickupDateTime)
		suite.NoError(err)
		err = suite.validator.Struct(validG62ScheduledPickupDate)
		suite.NoError(err)
	})

	suite.Run("validate failure 1", func() {
		g62 := G62{
			DateQualifier: 42,         // oneof
			Date:          "20190945", // datetime
			TimeQualifier: 42,         // oneof
			Time:          "2517",     // datetime
		}

		err := suite.validator.Struct(g62)
		suite.ValidateError(err, "DateQualifier", "oneof")
		suite.ValidateError(err, "Date", "datetime")
		suite.ValidateError(err, "TimeQualifier", "oneof")
		suite.ValidateError(err, "Time", "datetime")
		suite.ValidateErrorLen(err, 4)
	})
}
