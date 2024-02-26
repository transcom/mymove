package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	// FingerprintFlag is the Certificate Fingerprint flag
	DisableFingerprintFlag string = "fingerprint"

	// template for adding client certificates
	createDisableCertsMigration string = `
-- This migration removes an expired cert which had read/write access to all orders and the prime API.
-- OKTA Email is always the SHA-256 digest of the client SSL certificate @api.move.mil
-- Delete the public cert and the user role to prevent login privilege utilization.
-- For data integrity not delete the user account but set it to inactive.

DELETE FROM public.client_certs
where
user_id = (select id from users where okta_email = '{{.Fingerprint}}@api.move.mil');

DELETE FROM users_roles
where
user_id = (select id from users where okta_email = '{{.Fingerprint}}@api.move.mil');

UPDATE users
SET active=false
WHERE
okta_email='{{.Fingerprint}}@api.move.mil';

`
)

// CertsTemplate is a struct that stores the context from which to generate the migration
type DisableCertsTemplate struct {
	Fingerprint string
}

// InitCertsMigrationFlags initializes certs migration command line flags
func InitDisableCertsMigrationFlags(flag *pflag.FlagSet) {
	flag.StringP(FingerprintFlag, "f", "", "Certificate fingerprint in SHA 256 form")
}

// CheckCertsMigration validates command line flags
func CheckDisableCertsMigration(v *viper.Viper) error {
	if err := cli.CheckMigration(v); err != nil {
		return err
	}

	if err := cli.CheckMigrationFile(v); err != nil {
		return err
	}

	if err := cli.CheckCAC(v); err != nil {
		return err
	}

	if !v.GetBool(cli.CACFlag) {
		fingerprint := v.GetString(FingerprintFlag)
		if len(fingerprint) == 0 {
			return fmt.Errorf("%s is missing", FingerprintFlag)
		}
		sha256Pattern := "^[a-f0-9]{64}$"
		_, err := regexp.MatchString(sha256Pattern, fingerprint)
		if err != nil {
			return errors.Errorf("Fingerprint must be a valid SHA 256 hash")
		}
	}
	return nil
}

func initGenDisableCertsMigrationFlags(flag *pflag.FlagSet) {
	// Migration Config
	cli.InitMigrationFlags(flag)

	// Migration File Config
	cli.InitMigrationFileFlags(flag)

	// CAC Config
	cli.InitCACFlags(flag)

	// Init Certs Migration Flags
	InitCertsMigrationFlags(flag)

	// Don't sort command line flags
	flag.SortFlags = false
}

func genDisableCertsMigration(cmd *cobra.Command, args []string) error {
	// println("Function genDisableCertsMigration")
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

	err = CheckDisableCertsMigration(v)
	if err != nil {
		return err
	}
	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)

	var fingerprint string
	{
		fingerprint = v.GetString(FingerprintFlag)
	}

	certsDisableTemplate := DisableCertsTemplate{
		Fingerprint: fingerprint,
	}

	var t1 = template.Must(template.New("disable_certs_migration").Parse(createDisableCertsMigration))

	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", migrationVersion, migrationName)

	err = createMigration(tempMigrationPath, secureMigrationName, t1, certsDisableTemplate)
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
