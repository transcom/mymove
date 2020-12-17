package dtod

import (
	"crypto/tls"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
)

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

func (suite *DTODPlannerServiceSuite) TestDTODZip5Distance() {

	/*
		// For local testing, run `go test ./pkg/services/dtod_planner -v` to trigger the call
		// to the real DTOD test SOAP service. We DO NOT WANT TO RUN in regular UT/CircleCI cycle
		suite.T().Run("real call to DTOD uncomment locally to test", func(t *testing.T) {
			v  := suite.initFlags()
			tlsConfig := suite.getTLSConfig(v)

			dtodUsername, dtodPassword, dtodUrl, dtodWsdl, err := GetDTODFlags(v)
			suite.NoError(err)

			dtod := NewDTODZip5Distance(suite.DB(), suite.logger, tlsConfig, dtodUsername, dtodPassword, dtodUrl, dtodWsdl)
			dtod.DTODZip5Distance("05030", "05091") // actual distance is 22.195 miles
		})
	*/

}
