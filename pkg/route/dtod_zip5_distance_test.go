package route

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"

	"github.com/transcom/mymove/pkg/route/ghcmocks"
)

const distanceResponseTemplate = `<ProcessRequestResponse xmlns="https://dtod.sddc.army.mil/service/">
  <ProcessRequestResult>
	<Date>2020-12-22T16:09:41.7847017+00:00</Date>
	<Version>Current</Version>
	<Function>Distance</Function>
	<Region>Unassigned</Region>
	<RouteType>CommercialPersonalProperty</RouteType>
	<Units>Miles</Units>
	<Origin>
	  <City>Origin City</City>
	  <StateCountry>VT</StateCountry>
	  <County>Origin County</County>
	  <SplCode/>
	  <ZipCode>05030</ZipCode>
	  <Latitude>42.123456</Latitude>
	  <Longitude>-71.456789</Longitude>
	  <IsSplc>false</IsSplc>
	  <IsLatLong>false</IsLatLong>
	  <Region>NorthAmerica</Region>
	  <NameString>05030 Origin City, VT, Origin County</NameString>
	  <AbbreviationFormat>FIPS</AbbreviationFormat>
	</Origin>
	<Destination>
	  <City>Destination City</City>
	  <StateCountry>VT</StateCountry>
	  <County>Destination County</County>
	  <SplCode/>
	  <ZipCode>05091</ZipCode>
	  <Latitude>42.829292</Latitude>
	  <Longitude>-71.089378</Longitude>
	  <IsSplc>false</IsSplc>
	  <IsLatLong>false</IsLatLong>
	  <Region>NorthAmerica</Region>
	  <NameString>05091 Destination City, VT, Destination County</NameString>
	  <AbbreviationFormat>FIPS</AbbreviationFormat>
	</Destination>
	<Distance>%v</Distance>
	<Time>35</Time>
	<Matches/>
	<Directions/>
	<Coordinates/>
  </ProcessRequestResult>
</ProcessRequestResponse>`

func soapResponseForDistance(distance string) *gosoap.Response {
	// Note: Passing distance as a string so we can put bad data into the response
	return &gosoap.Response{
		Body: []byte(fmt.Sprintf(distanceResponseTemplate, distance)),
	}
}

func (suite *GHCTestSuite) TestDTODZip5DistanceFake() {
	tests := []struct {
		name             string
		responseDistance string
		responseError    bool
		expectedDistance int
		shouldError      bool
	}{
		{"distance round down", "72.323", false, 72, false},
		{"distance round up", "72.5", false, 73, false},
		{"negative distance", "-1", false, 0, true},
		{"error from call method", "25", true, 0, true},
		{"distance round down", "72.323", false, 72, false},
	}

	for _, test := range tests {
		suite.Run("fake call to DTOD: "+test.name, func() {
			var soapError error
			if test.responseError {
				soapError = errors.New("some error")
			}

			testSoapClient := &ghcmocks.SoapCaller{}
			testSoapClient.On("Call",
				mock.Anything,
				mock.Anything,
			).Return(soapResponseForDistance(test.responseDistance), soapError)

			dtod := NewDTODZip5Distance(fakeUsername, fakePassword, testSoapClient)
			distance, err := dtod.DTODZip5Distance(suite.AppContextForTest(), "05030", "05091")

			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
			suite.Equal(test.expectedDistance, distance)
		})
	}
}
