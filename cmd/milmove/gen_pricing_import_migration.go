package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

//type PricingTemplate struct {}
//Include instructions to start with a clean table/db because everything is going to be dumped
// ideally these are coming from the flags
const (
	//ImportFlag string = "import"
	DbHost              = "localhost"
	DbUser              = "postgres"
	DbName              = "dev_db"
	FilenameFlag string = "price_import_data"
)

//create a pg dump
func pgDump(fileName string) {
	fileName = FilenameFlag + ".sql"

	// GO script to run the following:
	// pg_dump -h localhost -U postgres -d dev_db -t re_* --data-only -T re_services* --data-only > pricing_data_dump.sql
	cmd := exec.Command(
		"pg_dump",
		"--compress=9",
		"-h "+DbHost,
		"-U "+DbUser,
		DbName,
		"-t re_* --data-only -T re_services* --data-only ",
	)

	log.Print(cmd)

	//Open the output file
	outfile, err := os.Create(fileName)
	if err != nil {
		log.Print("error in the os.Create")
		log.Fatal(err)
	}
	defer outfile.Close()

	cmd.Stdout = outfile

	//stdOut, err :=  cmd.StdoutPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//	Send stdout to the outfile
	//_, err = io.Copy(outfile, stdOut)

	// Start the command
	if err = cmd.Start(); err != nil {
		log.Print("error in the cmd.Start")
		log.Fatal(err)
	}

	log.Print("Waiting for command to finish…")

	// Wait for the command to finish.
	if err = cmd.Wait(); err != nil {
		log.Print("error in the cmd.Wait")
		log.Fatal(err)
	}
}

func initPricingImportMigrationFlags(flag *pflag.FlagSet) {
	flag.StringP(FilenameFlag, "f", "", "name file")
}

func checkPricingImportMigration(v *viper.Viper) error {
	migrationName := v.GetString(cli.MigrationNameFlag)
	if len(migrationName) == 0 {
		return fmt.Errorf("%s is missing", cli.MigrationNameFlag)
	}
	return nil
}

func initGenPricingImportMigrationFlags(flag *pflag.FlagSet) {
	// Flag for filename
	initPricingImportMigrationFlags(flag)

	// DB Config
	cli.InitDatabaseFlags(flag)

	//// Migration Config
	cli.InitMigrationFlags(flag)
	//
	//// Migration File Config
	cli.InitMigrationFileFlags(flag)
	//
	//// Migration Gen Path Config
	cli.InitMigrationGenPathFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

// generate migration
func genPricingImportMigration(cmd *cobra.Command, args []string) error {
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

	err = checkPricingImportMigration(v)
	if err != nil {
		return err
	}

	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)
	fileName := v.GetString(FilenameFlag)

	pgDump(fileName)

	// prompt migration_name preceded by -n
	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", migrationVersion, migrationName)

	err = createMigration(tempMigrationPath, secureMigrationName, nil, nil)
	if err != nil {
		log.Print("error in the createMigration1")

		return err
	}
	err = createMigration("./migrations/app/secure", secureMigrationName, nil, nil)
	if err != nil {
		log.Print("error in the createMigration")

		return err
	}

	err = addMigrationToManifest(migrationManifest, secureMigrationName)
	if err != nil {
		return err
	}
	return nil

}
