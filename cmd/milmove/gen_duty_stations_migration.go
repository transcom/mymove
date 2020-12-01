package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	dutyStations "github.com/transcom/mymove/pkg/services/duty_stations"
)

const (
	// DutyStationsFilenameFlag filename containing the details for new duty stations
	DutyStationsFilenameFlag string = "duty-stations-filename"
)

// MigrationInfo carries the filename of the migration
type MigrationInfo struct {
	Filename string
}

const (
	// DutyStationMigration is the duty station migration template
	DutyStationMigration string = `
-- Migration generated using: cmd/milmove/gen_duty_stations_migration.go
-- Duty stations file: {{.Filename}}`
)

// InitAddDutyStationsFlags initializes add_duty_stations command line flags
func InitAddDutyStationsFlags(flag *pflag.FlagSet) {
	flag.StringP(DutyStationsFilenameFlag, "f", "", "File name of csv file containing the new duty stations users")
}

// CheckAddDutyStations validates add_duty_stations command line flags
func CheckAddDutyStations(v *viper.Viper, logger logger) error {
	if err := cli.CheckDatabase(v, logger); err != nil {
		return err
	}

	if err := cli.CheckMigration(v); err != nil {
		return err
	}

	if err := cli.CheckMigrationFile(v); err != nil {
		return err
	}

	if err := cli.CheckMigrationGenPath(v); err != nil {
		return err
	}

	DutyStationsFilename := v.GetString(DutyStationsFilenameFlag)
	if DutyStationsFilename == "" {
		return fmt.Errorf("--duty-stations-filename is required")
	}
	return nil
}

func initGenDutyStationsMigrationFlags(flag *pflag.FlagSet) {
	// DB Config
	cli.InitDatabaseFlags(flag)

	// Migration Config
	cli.InitMigrationFlags(flag)

	// Migration File Config
	cli.InitMigrationFileFlags(flag)

	// Migration Gen Path Config
	cli.InitMigrationGenPathFlags(flag)

	// Add Duty Stations
	InitAddDutyStationsFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

func createDutyStationMigration(path string, filename string, ds []dutyStations.DutyStationMigration) error {
	migrationPath := filepath.Join(path, filename)
	migrationFile, err := os.Create(migrationPath)
	defer closeFile(migrationFile)
	if err != nil {
		return errors.Wrapf(err, "error creating %s", migrationPath)
	}

	t1 := template.Must(template.New("temp1").Parse(DutyStationMigration))
	err = t1.Execute(migrationFile, MigrationInfo{filename})
	if err != nil {
		log.Println("error executing template 1: ", err)
	}
	t2 := template.Must(template.New("temp2").Parse(dutyStations.InsertTemplate))
	err = t2.Execute(migrationFile, ds)
	if err != nil {
		log.Println("error executing template 2: ", err)
	}
	log.Printf("new migration file created at:  %q\n", migrationPath)
	return nil
}

func genDutyStationsMigration(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "could not ParseFlags on args")
	}

	flag := cmd.Flags()
	err = flag.Parse(os.Args[1:])
	if err != nil {
		return errors.Wrap(err, "could not parse flags")
	}

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		return errors.Wrap(err, "could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger, err := logging.Config(v.GetString(cli.LoggingEnvFlag), v.GetString(cli.LoggingLevelFlag))
	if err != nil {
		return err
	}

	err = CheckAddDutyStations(v, logger)
	if err != nil {
		return err
	}

	migrationPath := v.GetString(cli.MigrationGenPathFlag)
	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)
	dutyStationsFilename := v.GetString(DutyStationsFilenameFlag)

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		if dbConnection == nil {
			// No connection object means that the configuraton failed to validate and we should kill server startup
			logger.Fatal("Invalid DB Configuration", zap.Error(err))
		} else {
			// A valid connection object that still has an error indicates that the DB is not up and
			// thus is not ready for migrations
			logger.Fatal("DB is not ready for connections", zap.Error(err))
		}
	}

	builder := dutyStations.NewMigrationBuilder(dbConnection, logger)
	insertions, err := builder.Build(dutyStationsFilename)
	if err != nil {
		logger.Panic("Error while building migration", zap.Error(err))
	}

	migrationFilename := fmt.Sprintf("%s_%s.up.sql", migrationVersion, migrationName)
	err = createDutyStationMigration(migrationPath, migrationFilename, insertions)
	if err != nil {
		return err
	}

	err = addMigrationToManifest(migrationManifest, migrationFilename)
	if err != nil {
		return err
	}
	return nil
}
