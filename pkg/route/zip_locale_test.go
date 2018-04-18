package route

func (suite *PlannerSuite) TestZip5ToLatLong() {
	zip02807LatLong := LatLong{
		Latitude:  41.176815,
		Longitude: -71.577085,
	}

	// With leading 0
	ll, err := Zip5ToLatLong("02807")
	suite.Assertions.Nil(err, "Should not get error from Zip5")
	suite.Assertions.Equal(zip02807LatLong, ll, "Lat long for zip with leading zero")

	// Without leading 0
	ll, err = Zip5ToLatLong("2807")
	suite.Assertions.Nil(err, "Should not get error from Zip5 no leading 0")
	suite.Assertions.Equal(zip02807LatLong, ll, "Lat long for zip with no leading zero")

	// Not a number
	ll, err = Zip5ToLatLong("charleston")
	suite.Assertions.NotNil(err, "Should get error from Zip5 not number")

	// Not a valid zip
	ll, err = Zip5ToLatLong("12345")
	suite.Assertions.NotNil(err, "Should get error from Zip5 not valid")
}
