package dtod

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/services/dtod_planner/mocks"
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

func (suite *DTODPlannerServiceSuite) initFlags() *viper.Viper {
	flag := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	// Init TLS Flags
	cli.InitCertFlags(flag)

	// Init DTOD Flags
	InitDTODFlags(flag)

	flagSet := []string{}
	flag.Parse(flagSet)

	/*
		err := flag.Parse(os.Args[1:])
		if err != nil {
			suite.logger.Fatal("could not parse flags", zap.Error(err))
		}
	*/

	v := viper.New()
	err := v.BindPFlags(flag)
	if err != nil {
		suite.logger.Fatal("could not bind flags", zap.Error(err))
	}

	pflagsErr := v.BindPFlags(flag)
	if pflagsErr != nil {
		suite.logger.Fatal("invalid configuration", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return v
}

func (suite *DTODPlannerServiceSuite) getTLSConfig(v *viper.Viper) *tls.Config {
	certificates, rootCAs, err := certs.InitDoDCertificates(v, suite.logger)
	if certificates == nil || rootCAs == nil || err != nil {
		log.Fatal("Error in getting tls certs", err)
	}

	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}

	return tlsConfig
}

/*
func (suite *DTODPlannerServiceSuite) TestDTODZip5DistanceReal() {
	// For local testing, run `go test ./pkg/services/dtod_planner -v` to trigger the call
	// to the real DTOD test SOAP service. We DO NOT WANT TO RUN in regular UT/CircleCI cycle

	suite.T().Run("real call to DTOD uncomment locally to test", func(t *testing.T) {
		v := suite.initFlags()
		tlsConfig := suite.getTLSConfig(v)

		dtodUsername, dtodPassword, dtodURL, dtodWSDL, err := GetDTODFlags(v)
		suite.NoError(err)

		tr := &http.Transport{TLSClientConfig: tlsConfig}
		httpClient := &http.Client{Transport: tr, Timeout: time.Duration(30) * time.Second}

		// Use SoapClientWithConfig instead if you want to see the request and response
		soap, err := gosoap.SoapClient(dtodWSDL, httpClient)
		// soap, err := gosoap.SoapClientWithConfig(dtodWSDL, httpClient, &gosoap.Config{Dump: true})
		suite.NoError(err)
		soap.URL = dtodURL

		dtod := NewDTODZip5Distance(suite.logger, dtodUsername, dtodPassword, soap)
		distance, err := dtod.DTODZip5Distance("05030", "05091") // actual distance is 23.664 miles
		suite.NoError(err)
		suite.Equal(24, distance)
	})
}
*/

func (suite *DTODPlannerServiceSuite) TestDTODZip5DistanceFake() {
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
		{"error from call method", "0", true, 0, true},
		{"distance round down", "72.323", false, 72, false},
	}

	for _, test := range tests {
		suite.T().Run("fake call to DTOD: "+test.name, func(t *testing.T) {
			var soapError error
			if test.responseError {
				soapError = errors.New("some error")
			}

			testSoapClient := &mocks.SoapCaller{}
			testSoapClient.On("Call",
				mock.Anything,
				mock.Anything,
			).Return(soapResponseForDistance(test.responseDistance), soapError)

			dtod := NewDTODZip5Distance(suite.logger, "fake_username", "fake_password", testSoapClient)
			distance, err := dtod.DTODZip5Distance("05030", "05091")

			if test.shouldError {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
			suite.Equal(test.expectedDistance, distance)
		})
	}
}
