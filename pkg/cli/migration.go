package cli

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MigrationManifestFlag is the migration manifest flag
	MigrationManifestFlag    string = "migration-manifest" //deprecated
	DMLMigrationManifestFlag string = "dml-migration-manifest"
	// MigrationWaitFlag is the migration wait flag
	MigrationWaitFlag string = "migration-wait"

	DDLTablesMigrationPathFlag     = "ddl-tables-migration-path"
	DDLTablesMigrationManifestFlag = "ddl-tables-migration-manifest"

	DDLTypesMigrationPathFlag     = "ddl-types-migration-path"
	DDLTypesMigrationManifestFlag = "ddl-types-migration-manifest"

	DDLViewsMigrationPathFlag     = "ddl-views-migration-path"
	DDLViewsMigrationManifestFlag = "ddl-views-migration-manifest"

	DDLFunctionsMigrationPathFlag     = "ddl-functions-migration-path"
	DDLFunctionsMigrationManifestFlag = "ddl-functions-migration-manifest"
)

var (
	errMissingMigrationManifest = errors.New("missing migration manifest, expected to be set")
)

// InitMigrationFlags initializes the Migration command line flags
func InitMigrationFlags(flag *pflag.FlagSet) {
	flag.StringP(MigrationManifestFlag, "m", "migrations/app/migrations_manifest.txt", "Path to the manifest")
	flag.StringP(DMLMigrationManifestFlag, "d", "migrations/app/dml_migrations_manifest.txt", "Path to the manifest")
	flag.DurationP(MigrationWaitFlag, "w", time.Millisecond*10, "duration to wait when polling for new data from migration file")
	flag.String(DDLTablesMigrationPathFlag, "", "Path to DDL tables migrations directory")
	flag.String(DDLTablesMigrationManifestFlag, "", "Path to DDL tables migrations manifest")
	flag.String(DDLTypesMigrationPathFlag, "", "Path to DDL types migrations directory")
	flag.String(DDLTypesMigrationManifestFlag, "", "Path to DDL types migrations manifest")
	flag.String(DDLViewsMigrationPathFlag, "", "Path to DDL views migrations directory")
	flag.String(DDLViewsMigrationManifestFlag, "", "Path to DDL views migrations manifest")
	flag.String(DDLFunctionsMigrationPathFlag, "", "Path to DDL functions migrations directory")
	flag.String(DDLFunctionsMigrationManifestFlag, "", "Path to DDL functions migrations manifest")

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
	if len(DMLMigrationManifestFlag) == 0 {
		return errMissingMigrationManifest
	}
	if len(DDLTypesMigrationManifestFlag) == 0 {
		return errMissingMigrationManifest
	}
	if len(DDLTablesMigrationManifestFlag) == 0 {
		return errMissingMigrationManifest
	}
	if len(DDLViewsMigrationManifestFlag) == 0 {
		return errMissingMigrationManifest
	}
	if len(DDLFunctionsMigrationManifestFlag) == 0 {
		return errMissingMigrationManifest
	}
	return nil
}
