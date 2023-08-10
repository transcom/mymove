package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/storage"
	tdgs "github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/uploader"
)

const (
	scenarioFlag         string = "scenario"
	namedScenarioFlag    string = "named-scenario"
	namedSubScenarioFlag string = "named-sub-scenario" // name of the sub scenario in the main scenario
)

type errInvalidScenario struct {
	Name string
}

func (e *errInvalidScenario) Error() string {
	return fmt.Sprintf("invalid scenario: %s", e.Name)
}

func checkConfig(v *viper.Viper, logger *zap.Logger) error {

	logger.Debug("checking config")

	scenario := v.GetInt(scenarioFlag)
	if scenario > 0 {
		return errors.Wrap(&errInvalidScenario{Name: "0"}, "Numeric scenarios not supported")
	}

	namedScenario := v.GetString(namedScenarioFlag)
	_, err := findNamedScenarioByName(namedScenario)
	if err != nil {
		return err
	}

	err = cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func checkConfigNamedSubScenarioFlag(v *viper.Viper, namedScenarioStruct tdgs.NamedScenario) error {
	namedSubScenario := v.GetString(namedSubScenarioFlag)
	// optional flag, ok if value is empty
	// ok if there are not any named sub scenarios
	if namedSubScenario == "" || len(namedScenarioStruct.SubScenarios) == 0 {
		return nil
	}

	// continue, check if named sub scenarios matches expectations
	if _, ok := namedScenarioStruct.SubScenarios[namedSubScenario]; !ok {
		// to get the list of names
		var namedSubScenarioStringList []string
		for key := range namedScenarioStruct.SubScenarios {
			namedSubScenarioStringList = append(namedSubScenarioStringList, key)
		}

		return fmt.Errorf("%s is an invalid sub-scenario, expected "+
			"a value from %v or empty value", namedSubScenario, namedSubScenarioStringList)
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// Scenario config
	flag.Int(scenarioFlag, 0, "Specify which scenario you'd like to run. Current options: 1, 2, 3, 4, 5, 6, 7.")
	flag.String(namedScenarioFlag, "", "It's like a scenario, but more descriptive.")
	flag.String(namedSubScenarioFlag, "", "Specify a named-sub-scenario after specifying a named-scenario. "+
		"This is meant to run specific seed data setup in the main named-scenario without having to seed everything.")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Storage
	cli.InitStorageFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func findNamedScenarioByName(name string) (*tdgs.NamedScenario, error) {
	for _, scenario := range namedScenarios {
		result := scenario
		if name == scenario.Name {
			return &result, nil
		}
	}

	// to get the list of names
	var namedScenarioStringList []string
	for _, namedScenario := range namedScenarios {
		namedScenarioStringList = append(namedScenarioStringList, namedScenario.Name)
	}

	return nil, errors.Wrap(&errInvalidScenario{Name: name}, fmt.Sprintf("%s is invalid, expected "+
		"a value from %v", name, namedScenarioStringList))
}

var namedScenarios = []tdgs.NamedScenario{
	tdgs.NamedScenario(tdgs.E2eBasicScenario),
	tdgs.NamedScenario(tdgs.DevSeedScenario),
	tdgs.NamedScenario(tdgs.BandwidthScenario),
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

	logger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	// run inside a transaction as some of the testdatagen needs it
	err = appcontext.NewAppContext(dbConnection, logger, nil).NewTransaction(
		func(appCtx appcontext.AppContext) error {

			namedScenario := v.GetString(namedScenarioFlag)
			namedSubScenario := v.GetString(namedSubScenarioFlag)

			if namedScenario != "" {

				storer := storage.InitStorage(v, logger)

				userUploader, uploaderErr := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
				if uploaderErr != nil {
					logger.Fatal("could not instantiate user uploader", zap.Error(err))
				}
				primeUploader, uploaderErr := uploader.NewPrimeUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
				if uploaderErr != nil {
					logger.Fatal("could not instantiate prime uploader", zap.Error(err))
				}

				if namedScenario == tdgs.E2eBasicScenario.Name {
					tdgs.E2eBasicScenario.Run(appCtx, userUploader, primeUploader)
				} else if namedScenario == tdgs.DevSeedScenario.Name {
					// Something is different about our cert config in CI so only running this
					// for the devseed scenario not e2e_basic for Cypress
					certificates, rootCAs, certErr := certs.InitDoDCertificates(v, logger)
					if certificates == nil || rootCAs == nil || certErr != nil {
						logger.Fatal("Failed to initialize DOD certificates", zap.Error(certErr))
					}

					// Initialize setup
					tdgs.DevSeedScenario.Setup(appCtx, userUploader, primeUploader)

					// Sub-scenarios are generated at run time
					// Check config
					// optional flag
					if serr := checkConfigNamedSubScenarioFlag(v, tdgs.NamedScenario(tdgs.DevSeedScenario)); serr != nil {
						logger.Fatal("invalid configuration", zap.Error(serr))
					}

					// Run seed
					tdgs.DevSeedScenario.Run(appCtx, namedSubScenario)
				} else if namedScenario == tdgs.BandwidthScenario.Name {
					tdgs.BandwidthScenario.Run(appCtx, userUploader, primeUploader)
				}

				logger.Info("Success! Created e2e test data.")
			} else {
				flag.PrintDefaults()
			}
			return nil
		})
	if err != nil {
		log.Fatal("Failed to load scenario", zap.Error(err))
	}
}
