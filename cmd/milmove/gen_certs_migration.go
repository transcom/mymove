package main

import (
	"crypto/sha256"
	"encoding/hex"
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
	// FingerprintFlag is the Certificate Fingerprint flag
	FingerprintFlag string = "fingerprint"
	// SubjectFlag is the Certificate Subject flag
	SubjectFlag string = "subject"
	// UpdateFlag is the flag to enable Update mode
	UpdateFlag string = "update"
	// ClientCertIDFlag is the ID flag for Client Cert
	ClientCertIDFlag string = "certid"

	// template for adding client certificates
	createCertsMigration string = `
-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
INSERT INTO public.client_certs (
	id,
	sha256_digest,
	subject,
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
	allow_navy_orders_write,
	allow_prime)
VALUES (
	'{{.ID}}',
	'{{.Fingerprint}}',
	'{{.Subject}}',
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
	true,
	true);
`
	// template for updating client certificates
	updateCertsMigration string = `
-- This migration allows a CAC cert to have read/write access to all orders and the prime API.
-- The Orders API and the Prime API use client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
-- Updating sha256 for new cac
UPDATE client_certs
SET sha256_digest = '{{.Fingerprint}}',
	subject = '{{.Subject}}',
	updated_at = NOW()
WHERE id = '{{.ID}}';
`
)

// CertsTemplate is a struct that stores the context from which to generate the migration
type CertsTemplate struct {
	ID          string
	Fingerprint string
	Subject     string
}

// InitCertsMigrationFlags initializes certs migration command line flags
func InitCertsMigrationFlags(flag *pflag.FlagSet) {
	flag.StringP(FingerprintFlag, "f", "", "Certificate fingerprint in SHA 256 form")
	flag.StringP(SubjectFlag, "s", "", "Certificate subject")
	flag.Bool(UpdateFlag, false, "Create an update migration")
	flag.String(ClientCertIDFlag, "", "Previous ID of client cert entry")

}

// CheckCertsMigration validates command line flags
func CheckCertsMigration(v *viper.Viper) error {
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

		subject := v.GetString(SubjectFlag)
		if len(subject) == 0 {
			return errors.Errorf("%s is missing", SubjectFlag)
		}
	}

	if v.GetBool(UpdateFlag) {
		certID := v.GetString(ClientCertIDFlag)
		if len(certID) == 0 {
			return fmt.Errorf("%s required for generating UPDATE migration", ClientCertIDFlag)
		}
		uuid.Must(uuid.FromString(certID))
	}

	return nil
}

func initGenCertsMigrationFlags(flag *pflag.FlagSet) {
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

func genCertsMigration(cmd *cobra.Command, args []string) error {
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

	err = CheckCertsMigration(v)
	if err != nil {
		return err
	}
	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)

	var fingerprint, subject string
	if v.GetBool(cli.CACFlag) {

		store, errStore := cli.GetCACStore(v)
		if errStore != nil {
			return fmt.Errorf("Ensure CAC reader and card inserted: %w", errStore)
		}

		defer func() {
			if closeErr := store.Close(); closeErr != nil {
				fmt.Println("Failed to close CAC store")
			}
		}()

		cert, errTLSCert := store.TLSCertificate()
		if errTLSCert != nil {
			return errTLSCert
		}

		// Get the fingerprint
		hash := sha256.Sum256(cert.Certificate[0])
		fingerprint = hex.EncodeToString(hash[:])

		// Get the subject in RFC2253 format
		subject = cert.Leaf.Subject.String()
	} else {
		fingerprint = v.GetString(FingerprintFlag)
		subject = v.GetString(SubjectFlag)
	}

	certsTemplate := CertsTemplate{
		Fingerprint: fingerprint,
		Subject:     subject,
	}

	var t1 *template.Template
	updateMode := v.GetBool(UpdateFlag)

	if updateMode {
		certsTemplate.ID = v.GetString(ClientCertIDFlag)
		t1 = template.Must(template.New("certs_migration").Parse(updateCertsMigration))
	} else {
		certsTemplate.ID = uuid.Must(uuid.NewV4()).String()
		t1 = template.Must(template.New("certs_migration").Parse(createCertsMigration))
	}

	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", migrationVersion, migrationName)

	err = createMigration(tempMigrationPath, secureMigrationName, t1, certsTemplate)
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
