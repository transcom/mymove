package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/storage"
	tdgs "github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/uploader"
)

func stringSliceContains(stringSlice []string, value string) bool {
	for _, x := range stringSlice {
		if value == x {
			return true
		}
	}
	return false
}

const (
	scenarioFlag      string = "scenario"
	namedScenarioFlag string = "named-scenario"
)

type errInvalidScenario struct {
	Scenario int
}

func (e *errInvalidScenario) Error() string {
	return fmt.Sprintf("invalid scenario %d", e.Scenario)
}

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	scenario := v.GetInt(scenarioFlag)
	if scenario < 0 || scenario > 7 {
		return errors.Wrap(&errInvalidScenario{Scenario: scenario}, fmt.Sprintf("%s is invalid, expected value between 0 and 7 not %d", scenarioFlag, scenario))
	}

	namedScenarios := []string{
		tdgs.E2eBasicScenario.Name,
		tdgs.DevSeedScenario.Name,
		tdgs.BandwidthScenario.Name,
	}
	namedScenario := v.GetString(namedScenarioFlag)
	if !stringSliceContains(namedScenarios, namedScenario) {
		return errors.Wrap(&errInvalidScenario{Scenario: scenario}, fmt.Sprintf("%s is invalid, expected a value from %v", namedScenarioFlag, namedScenarios))
	}

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// Scenario config
	flag.Int(scenarioFlag, 0, "Specify which scenario you'd like to run. Current options: 1, 2, 3, 4, 5, 6, 7.")
	flag.String(namedScenarioFlag, "", "It's like a scenario, but more descriptive.")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Storage
	cli.InitStorageFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

// Hey, refactoring self: you can pull the UUIDs from the objects rather than
// querying the db for them again.
func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	parseErr := flag.Parse(os.Args[1:])
	if parseErr != nil {
		log.Fatalf("Could not parse flags: %v\n", parseErr)
	}

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		log.Fatal("failed to bind flags", zap.Error(bindErr))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	scenario := v.GetInt(scenarioFlag)
	namedScenario := v.GetString(namedScenarioFlag)

	if scenario == 4 {
		err = tdgs.RunPPMSITEstimateScenario1(dbConnection)
	} else if scenario == 5 {
		err = tdgs.RunRateEngineScenario1(dbConnection)
	} else if scenario == 6 {
		query := `DELETE FROM transportation_service_provider_performances;
				  DELETE FROM transportation_service_providers;
				  DELETE FROM traffic_distribution_lists;
				  DELETE FROM tariff400ng_zip3s;
				  DELETE FROM tariff400ng_zip5_rate_areas;
				  DELETE FROM tariff400ng_service_areas;
				  DELETE FROM tariff400ng_linehaul_rates;
				  DELETE FROM tariff400ng_shorthaul_rates;
				  DELETE FROM tariff400ng_full_pack_rates;
				  DELETE FROM tariff400ng_full_unpack_rates;`

		err = dbConnection.RawQuery(query).Exec()
		if err != nil {
			logger.Fatal("Failed to run raw query", zap.Error(err))
		}
		err = tdgs.RunRateEngineScenario2(dbConnection)
	} else if namedScenario != "" {
		// Initialize logger
		logger, newDevelopmentErr := zap.NewDevelopment()
		if newDevelopmentErr != nil {
			logger.Fatal("Problem with zap NewDevelopment", zap.Error(newDevelopmentErr))
		}

		// Initialize storage and uploader
		var session *awssession.Session
		storageBackend := v.GetString(cli.StorageBackendFlag)
		if storageBackend == "s3" {
			c := &aws.Config{
				Region: aws.String(v.GetString(cli.AWSRegionFlag)),
			}
			s, errorSession := awssession.NewSession(c)

			if errorSession != nil {
				logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
			}

			session = s
		}
		storer := storage.InitStorage(v, session, logger)

		userUploader, uploaderErr := uploader.NewUserUploader(dbConnection, logger, storer, 25*uploader.MB)
		if uploaderErr != nil {
			logger.Fatal("could not instantiate user uploader", zap.Error(err))
		}
		primeUploader, uploaderErr := uploader.NewPrimeUploader(dbConnection, logger, storer, 25*uploader.MB)
		if uploaderErr != nil {
			logger.Fatal("could not instantiate prime uploader", zap.Error(err))
		}

		if namedScenario == tdgs.E2eBasicScenario.Name {
			tdgs.E2eBasicScenario.Run(dbConnection, userUploader, primeUploader, logger)
		} else if namedScenario == tdgs.DevSeedScenario.Name {
			// Something is different about our cert config in CI so only running this
			// for the devseed scenario not e2e_basic for Cypress
			certificates, rootCAs, certErr := certs.InitDoDCertificates(v, logger)
			if certificates == nil || rootCAs == nil || certErr != nil {
				logger.Fatal("Failed to initialize DOD certificates", zap.Error(certErr))
			}

			// Create a secondary planner specifically for GHC.
			routeTLSConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}
			routePlanner, plannerErr := route.InitGHCRoutePlanner(v, logger, dbConnection, routeTLSConfig)

			if plannerErr != nil {
				logger.Fatal("Failed to initialize GHC route planner")
			}

			tdgs.DevSeedScenario.Run(dbConnection, userUploader, primeUploader, routePlanner, logger)
		} else if namedScenario == tdgs.BandwidthScenario.Name {
			tdgs.BandwidthScenario.Run(dbConnection, userUploader, primeUploader)
		}

		logger.Info("Success! Created e2e test data.")
	} else {
		flag.PrintDefaults()
	}
	if err != nil {
		log.Fatal("Failed to load scenario", zap.Error(err))
	}
}
