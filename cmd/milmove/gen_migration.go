package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

// initGenFlags - Order matters!
func initGenMigrationFlags(flag *pflag.FlagSet) {

	// Migration Config
	cli.InitMigrationFlags(flag)

	// Migration File Config
	cli.InitMigrationFileFlags(flag)

	// Migration Gen Path Config
	cli.InitMigrationGenPathFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

func checkGenMigrationConfig(v *viper.Viper) error {

	if err := cli.CheckMigration(v); err != nil {
		return err
	}

	if err := cli.CheckMigrationFile(v); err != nil {
		return err
	}

	return cli.CheckMigrationGenPath(v)
}

func genMigrationFunction(cmd *cobra.Command, args []string) error {

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

	err = checkGenMigrationConfig(v)
	if err != nil {
		return err
	}

	migrationPath := v.GetString(cli.MigrationGenPathFlag)
	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationType := v.GetString(cli.MigrationTypeFlag)

	filename := fmt.Sprintf("%s_%s.up.%s", migrationVersion, migrationName, migrationType)
	err = writeEmptyFile(migrationPath, filename)
	if err != nil {
		return err
	}

	err = addMigrationToManifest(migrationManifest, filename)
	if err != nil {
		return err
	}
	return nil
}
