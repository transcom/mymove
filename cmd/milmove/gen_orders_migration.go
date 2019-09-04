package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	// OrdersFingerprintFlag is the Orders Certificate Fingerprint flag
	OrdersFingerprintFlag string = "fingerprint"
	// OrdersSubjectFlag is the Orders Certificate Subject flag
	OrdersSubjectFlag string = "subject"

	// template for adding orders certificates
	createOrdersMigration string = `
-- Until the admin UI is in place and has visibility on the electronic orders table,
-- we need certificates that can look at the Orders that have been uploaded.
-- This migration allows a CAC cert to have read/write access to all orders.
-- The Orders API uses client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to upload Orders and review Orders. Eventually
-- this CAC certificate should be removed.
INSERT INTO public.client_certs (
	id,
	sha256_digest,
	subject,
	allow_dps_auth_api,
	allow_orders_api,
	created_at,
	updated_at,
	allow_air_force_orders_read,
	allow_air_force_orders_write,
	allow_army_orders_read,
	allow_army_orders_write,
	allow_coast_guard_orders_read,
	allow_coast_guard_orders_write,
	allow_marine_corps_orders_read,
	allow_marine_corps_orders_write,
	allow_navy_orders_read,
	allow_navy_orders_write)
VALUES (
	'{{.ID}}',
	'{{.Fingerprint}}',
	'{{.Subject}}',
	false,
	true,
	now(),
	now(),
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	true);
`
)

// OrdersTemplate is a struct that stores the context from which to generate the migration
type OrdersTemplate struct {
	ID          string
	Fingerprint string
	Subject     string
}

// InitOrdersMigrationFlags initializes orders migration command line flags
func InitOrdersMigrationFlags(flag *pflag.FlagSet) {
	flag.StringP(OrdersFingerprintFlag, "f", "", "Certificate fingerprint in SHA 256 form")
	flag.StringP(OrdersSubjectFlag, "s", "", "Certificate subject")
}

// CheckOrdersMigration validates add_office_users command line flags
func CheckOrdersMigration(v *viper.Viper) error {
	if err := cli.CheckMigration(v); err != nil {
		return err
	}

	if err := cli.CheckMigrationFile(v); err != nil {
		return err
	}

	fingerprint := v.GetString(OrdersFingerprintFlag)
	if len(fingerprint) == 0 {
		return fmt.Errorf("%s is missing", OrdersFingerprintFlag)
	}
	sha256Pattern := "^[a-f0-9]{64}$"
	_, err := regexp.MatchString(sha256Pattern, fingerprint)
	if err != nil {
		return errors.Errorf("Fingerprint must be a valid SHA 256 hash")
	}

	subject := v.GetString(OrdersSubjectFlag)
	if len(subject) == 0 {
		return errors.Errorf("%s is missing", OrdersSubjectFlag)
	}

	return nil
}

func initGenOrdersMigrationFlags(flag *pflag.FlagSet) {
	// Migration Config
	cli.InitMigrationFlags(flag)

	// Migration File Config
	cli.InitMigrationFileFlags(flag)

	// Init Orders Migration Flags
	InitOrdersMigrationFlags(flag)

	// Don't sort command line flags
	flag.SortFlags = false
}

func genOrdersMigration(cmd *cobra.Command, args []string) error {
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

	err = CheckOrdersMigration(v)
	if err != nil {
		return err
	}
	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)

	ordersTemplate := OrdersTemplate{
		ID:          uuid.Must(uuid.NewV4()).String(),
		Fingerprint: v.GetString(OrdersFingerprintFlag),
		Subject:     v.GetString(OrdersSubjectFlag),
	}

	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", migrationVersion, migrationName)
	t1 := template.Must(template.New("orders_migration").Parse(createOrdersMigration))
	err = createMigration(tempMigrationPath, secureMigrationName, t1, ordersTemplate)
	if err != nil {
		return err
	}

	t2 := template.Must(template.New("local_migrations").Parse(localMigrationTemplate))
	err = createMigration("./local_migrations", secureMigrationName, t2, nil)
	if err != nil {
		return err
	}

	err = addMigrationToManifest(migrationManifest, secureMigrationName)
	if err != nil {
		return err
	}
	return nil
}
