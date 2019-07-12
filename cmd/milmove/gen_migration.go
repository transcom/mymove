package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

	return nil
}

func addMigrationToManifest(migrationManifest string, filename string) error {
	mmf, err := os.OpenFile(migrationManifest, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "could not open migration manifest")
	}
	defer mmf.Close()

	_, err = mmf.WriteString(filename + "\n")
	if err != nil {
		return errors.Wrap(err, "could not append to the migration manifest")
	}

	log.Println(fmt.Sprintf("new migration appended to manifest at %q", migrationManifest))
	return nil
}

func genMigrationFunction(cmd *cobra.Command, args []string) error {

	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "Could not parse flags")
	}

	flag := cmd.Flags()

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		return errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	err = checkGenMigrationConfig(v)
	if err != nil {
		return err
	}

	migrationPath := v.GetString(cli.MigrationPathFlag)
	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationType := v.GetString(cli.MigrationTypeFlag)

	filename := fmt.Sprintf("%s_%s.up.%s", migrationVersion, migrationName, migrationType)
	p := filepath.Join(migrationPath, filename)

	err = ioutil.WriteFile(p, []byte{}, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write new migration file")
	}

	fmt.Println(fmt.Sprintf("new migration file created at %q", p))

	err = addMigrationToManifest(migrationManifest, filename)
	if err != nil {
		return err
	}
	return nil
}
