package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	// DisableUserEmailFlag is the email of the user to disable
	DisableUserEmailFlag string = "migration-email"

	// template for adding office users
	disableUser string = `UPDATE users
SET active=false
WHERE okta_email='{{.EmailPrefix}}@{{.EmailDomain}}';
`
)

// UserTemplate is a struct that stores the EmailPrefix from which to generate the migration
type UserTemplate struct {
	EmailPrefix string
	EmailDomain string
}

// InitDisableUserFlags initializes command line flags
func InitDisableUserFlags(flag *pflag.FlagSet) {
	flag.StringP(DisableUserEmailFlag, "e", "", "Email address")
}

func initDisableUserMigrationFlags(flag *pflag.FlagSet) {
	// Migration Config
	cli.InitMigrationFlags(flag)

	// Migration File Config
	cli.InitMigrationFileFlags(flag)

	// Disable User
	InitDisableUserFlags(flag)

	// Don't sort command line flags
	flag.SortFlags = false
}

// CheckDisableUserFlags validates add_office_users command line flags
func CheckDisableUserFlags(v *viper.Viper) error {
	if err := cli.CheckMigration(v); err != nil {
		return err
	}

	if err := cli.CheckMigrationFile(v); err != nil {
		return err
	}

	email := v.GetString(DisableUserEmailFlag)
	if len(email) == 0 {
		return errors.Errorf("%s is missing", DisableUserEmailFlag)
	}
	return nil
}

func genDisableUserMigration(cmd *cobra.Command, args []string) error {
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

	err = CheckDisableUserFlags(v)
	if err != nil {
		return err
	}

	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)
	migrationEmail := strings.Split(v.GetString(DisableUserEmailFlag), "@")

	emailPrefix := migrationEmail[0]
	emailDomain := migrationEmail[1]

	user := UserTemplate{EmailPrefix: emailPrefix, EmailDomain: emailDomain}

	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", migrationVersion, migrationName)
	t1 := template.Must(template.New("disable_user").Parse(disableUser))
	err = createMigration(tempMigrationPath, secureMigrationName, t1, user)
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
