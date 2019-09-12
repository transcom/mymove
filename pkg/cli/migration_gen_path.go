package cli

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MigrationGenPathFlag is the migration path (for generated migrations) flag
	MigrationGenPathFlag string = "migration-gen-path"
)

var (
	errMissingMigrationGenPath = errors.New("missing migration path, expected to be set")
)

type errInvalidMigrationGenPath struct {
	Path string
}

func (e *errInvalidMigrationGenPath) Error() string {
	return fmt.Sprintf("invalid migration file path %q", e.Path)
}

// InitMigrationFlags initializes the Migration command line flags
func InitMigrationGenPathFlags(flag *pflag.FlagSet) {
	flag.StringP(MigrationGenPathFlag, "p", "./migrations", "Path to the migrations folder")
}

// CheckMigration validates migration command line flags
func CheckMigrationGenPath(v *viper.Viper) error {
	migrationGenPath := v.GetString(MigrationGenPathFlag)
	if len(migrationGenPath) == 0 {
		return errMissingMigrationGenPath
	}
	if _, err := os.Stat(migrationGenPath); os.IsNotExist(err) {
		return errors.Wrapf(&errInvalidMigrationGenPath{Path: migrationGenPath}, "Expected %s to be a path in the filesystem", migrationGenPath)
	}
	return nil
}
