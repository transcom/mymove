package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	// DisableUserMigrationFilenameFlag sql file containing the migration
	DisableUserMigrationFilenameFlag string = "migration-filename"
)

const (
	// template for adding office users
	disableUser string = `UPDATE admin_users
SET disabled=true
WHERE email='{{.EmailPrefix}}@truss.works';

UPDATE office_users
SET disabled=true
WHERE email='{{.EmailPrefix}}@truss.works';

UPDATE tsp_users
SET disabled=true
WHERE email='{{.EmailPrefix}}+pyvl@truss.works'
	OR email='{{.EmailPrefix}}+dlxm@truss.works'
	OR email='{{.EmailPrefix}}+ssow@truss.works';
`
)

// UserTemplate is a struct that stores the EmailPrefix from which to generate the migration
type UserTemplate struct {
	EmailPrefix string
}

func genDisableUserMigration(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
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

	migrationsPath := v.GetString(cli.MigrationPathFlag)
	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationFileName := v.GetString(DisableUserMigrationFilenameFlag)

	user := UserTemplate{EmailPrefix: "NAME"}

	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", time.Now().Format(VersionTimeFormat), migrationFileName)
	t1 := template.Must(template.New("disable_user").Parse(disableUser))
	err = createMigration("./tmp", secureMigrationName, t1, user)
	if err != nil {
		return err
	}
	localMigrationPath := filepath.Join("local_migrations", secureMigrationName)
	localMigrationFile, err := os.Create(localMigrationPath)
	defer closeFile(localMigrationFile)
	if err != nil {
		return errors.Wrapf(err, "error creating %s", localMigrationPath)
	}
	log.Printf("new migration file created at:  %q\n", localMigrationPath)

	migrationName := fmt.Sprintf("%s_%s.up.fizz", time.Now().Format(VersionTimeFormat), migrationFileName)
	t2 := template.Must(template.New("migration").Parse(migration))
	err = createMigration(migrationsPath, migrationName, t2, secureMigrationName)
	if err != nil {
		return err
	}

	err = addMigrationToManifest(migrationManifest, migrationName)
	if err != nil {
		return err
	}
	return nil
}
