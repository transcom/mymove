package cli

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MigrationManifestFlag is the migration manifest flag
	MigrationManifestFlag string = "migration-manifest"
	// MigrationWaitFlag is the migration wait flag
	MigrationWaitFlag string = "migration-wait"
	// MigrationSchemaPathFlag contains the path to the schema file
	MigrationSchemaPathFlag string = "migration-schema-path"
	// MigrationLoadDevSeedFlag is whether to load dev seed data
	MigrationLoadDevSeedFlag string = "migration-load-dev-seed"
	// MigrationPrintStatusFlag is whether to only check the status of migrations
	MigrationPrintStatusFlag string = "migration-print-status"
	// MigrationCheckAppliedFlag exits with a zero status unless all
	// migrations are applied"
	MigrationCheckAppliedFlag string = "migration-check-applied"
)

var (
	errMissingMigrationManifest = errors.New("missing migration manifest, expected to be set")
)

// InitMigrationFlags initializes the Migration command line flags
func InitMigrationFlags(flag *pflag.FlagSet) {
	flag.StringP(MigrationManifestFlag, "m", "migrations/app/migrations_manifest.txt", "Path to the manifest")
	flag.DurationP(MigrationWaitFlag, "w", time.Millisecond*10, "duration to wait when polling for new data from migration file")
	flag.String(MigrationSchemaPathFlag, "", "Path to full schema file")
	flag.Bool(MigrationLoadDevSeedFlag, false, "Load the dev seed data")
	flag.Bool(MigrationPrintStatusFlag, false, "Print migration status only")
	flag.Bool(MigrationCheckAppliedFlag, false, "Check migration status, exit with nonzero if pending migrations need to be applied")
}

// CheckMigration validates migration command line flags
func CheckMigration(v *viper.Viper) error {
	migrationManifest := v.GetString(MigrationManifestFlag)
	if len(migrationManifest) == 0 {
		return errMissingMigrationManifest
	}
	if len(MigrationManifestFlag) == 0 {
		return errMissingMigrationManifest
	}
	return nil
}
