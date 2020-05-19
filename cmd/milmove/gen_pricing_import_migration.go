package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

//type PricingTemplate struct {}
const (
	//ImportFlag string = "import"
	DbHost   = "localhost"
	DbUser   = "postgres"
	DbName   = "dev_db"
	Filename = "pricing_data_dump"
)

//create a pg dump

func pgDump() {
	fileName := Filename
	// GO script for pg_dump -h localhost -U postgres -d dev_db -t re_* --data-only -T re_services* --data-only > pricing_data_dump.sql
	cmd := exec.Command(
		"pg_dump",
		"-h"+DbHost,
		"-U"+DbUser,
		DbName,
		"-t re_* --data-only -T re_services --data-only",
	)

	//Open the output file
	outfile, err := os.Create(fileName + ".sql")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	//	Send stdout to the outfile
	cmd.Stdout = outfile

	// Start the command
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Print("Waiting for command to finish…")

	// Wait for the command to finish.
	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func initGenPricingImportMigrationFlags(flag *pflag.FlagSet) {
	// DB Config
	cli.InitDatabaseFlags(flag)

	// Migration Config
	cli.InitMigrationFlags(flag)

	// Migration File Config
	cli.InitMigrationFileFlags(flag)

	// Migration Gen Path Config
	cli.InitMigrationGenPathFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

// generate migration
func genPricingImportMigration(cmd *cobra.Command, args []string) error {

	v := viper.New()

	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)

	// will i need a template?

	//store the results of the pgDump in to a variable
	pgDump()
	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", migrationVersion, migrationName)

	//pricingTemplate := PricingTemplate{}

	t1 := template.New("pricing_import_migration")
	err := createMigration(tempMigrationPath, secureMigrationName, t1, nil)
	if err != nil {
		return err
	}

	t2 := template.Must(template.New("migrations/app/secure").Parse(localMigrationTemplate))
	err = createMigration("./migrations/app/secure", secureMigrationName, t2, nil)
	if err != nil {
		return err
	}

	err = addMigrationToManifest(migrationManifest, secureMigrationName)
	if err != nil {
		return err
	}
	return nil

}
