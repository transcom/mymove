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
	// DDLMigrationManifestFlag is the ddl migration manifest flag
	DDLMigrationManifestFlag = "ddl-migration-manifest"
	// DDLMigrationPathFlag is the ddl migration path flag
	DDLMigrationPathFlag = "ddl-migration-path"
)

var (
	errMissingMigrationManifest = errors.New("missing migration manifest, expected to be set")
)

// InitMigrationFlags initializes the Migration command line flags
func InitMigrationFlags(flag *pflag.FlagSet) {
	flag.StringP(MigrationManifestFlag, "m", "migrations/app/migrations_manifest.txt", "Path to the manifest")
	flag.DurationP(MigrationWaitFlag, "w", time.Millisecond*10, "duration to wait when polling for new data from migration file")
	flag.String(DDLMigrationManifestFlag, "", "Path to DDL migrations manifest")
	flag.String(DDLMigrationPathFlag, "", "Path to DDL migrations directory")
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
