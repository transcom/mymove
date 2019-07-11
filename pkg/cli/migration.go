package cli

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MigrationPathFlag is the migration path flag
	MigrationPathFlag string = "migration-path"
	// MigrationManifestFlag is the migration manifest flag
	MigrationManifestFlag string = "migration-manifest"
)

var (
	errMissingMigrationPath     = errors.New("missing migration path, expected to be set")
	errMissingMigrationManifest = errors.New("missing migration manifest, expected to be set")
)

type errInvalidMigrationPath struct {
	Path string
}

func (e *errInvalidMigrationPath) Error() string {
	return fmt.Sprintf("invalid migration path %q", e.Path)
}

// InitMigrationFlags initializes the Migration command line flags
func InitMigrationFlags(flag *pflag.FlagSet) {
	flag.StringP(MigrationPathFlag, "p", "./migrations", "Path to the migrations folder")
	flag.StringP(MigrationManifestFlag, "m", "./migrations_manifest.txt", "Path to the manifest")
}

// CheckMigration validates migration command line flags
func CheckMigration(v *viper.Viper) error {
	migrationPath := v.GetString(MigrationPathFlag)
	if len(migrationPath) == 0 {
		return errMissingMigrationPath
	}
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		return errors.Wrapf(&errInvalidMigrationPath{Path: migrationPath}, "Expected %s to exist", migrationPath)
	}
	if len(MigrationManifestFlag) == 0 {
		return errMissingMigrationManifest
	}
	return nil
}
