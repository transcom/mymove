package route

func (suite *PlannerSuite) TestZip5ToZip3LatLong() {
	// 	"028": {41.8230, -71.4187},  //providence,ri
	zip02807LatLong := LatLong{
		Latitude:  41.8230,
		Longitude: -71.4187,
	}

	// With leading 0
	ll, err := Zip5ToZip3LatLong("02807")
	suite.Assertions.Nil(err, "Should not get error from Zip5")
	suite.Assertions.Equal(zip02807LatLong, ll, "Lat long for zip with leading zero")

	// With delivery route
	ll, err = Zip5ToZip3LatLong("02807-9999")
	suite.Assertions.Nil(err, "Should not get error from Zip5 with route")
	suite.Assertions.Equal(zip02807LatLong, ll, "Lat long for zip with route")

	// Without leading 0
	ll, err = Zip5ToZip3LatLong("2807")
	suite.Assertions.Nil(err, "Should not get error from Zip5 no leading 0")
	suite.Assertions.Equal(zip02807LatLong, ll, "Lat long for zip with no leading zero")

	// Greater than 65636
	// 	"941": {37.7562, -122.4430}, //san francisco,ca
	zip94103LatLong := LatLong{
		Latitude:  37.7562,
		Longitude: -122.4430,
	}

	ll, err = Zip5ToZip3LatLong("94103")
	suite.Assertions.Nil(err, "Should not get error from Zip5 >64k")
	suite.Assertions.Equal(zip94103LatLong, ll, "Lat long for zip with no leading zero")

	// Not a number
	_, err = Zip5ToZip3LatLong("charleston")
	suite.Assertions.NotNil(err, "Should get error from Zip5 not number")

	// Not a valid zip
	_, err = Zip5ToZip3LatLong("00001")
	suite.Assertions.NotNil(err, "Should get error from Zip5 not valid")

	// With more than 5 numbers
	_, err = Zip5ToZip3LatLong("0280799")
	suite.Assertions.NotNil(err, "Should get error from Zip5 with more than 5 numbers")

}