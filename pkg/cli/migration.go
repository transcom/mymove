package cli

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MigrationPathFlag is the migration path flag
	MigrationPathFlag string = "path"
)

type errInvalidMigrationPath struct {
	Path string
}

func (e *errInvalidMigrationPath) Error() string {
	return fmt.Sprintf("invalid migration path '%s'", e.Path)
}

// InitMigrationFlags initializes the Migration command line flags
func InitMigrationFlags(flag *pflag.FlagSet) {
	flag.StringP(MigrationPathFlag, "p", "./migrations", "Path to the migrations folder")
}

// CheckMigration validates migration command line flags
func CheckMigration(v *viper.Viper) error {
	if migrationPath := v.GetString(MigrationPathFlag); len(migrationPath) == 0 {
		return errors.Wrap(&errInvalidMigrationPath{Path: migrationPath}, "Expected a migration path to be set")
	}
	return nil
}
