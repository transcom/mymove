package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MigrationManifestFlag is the migration manifest flag
	MigrationManifestFlag string = "migration-manifest"
	// MigrationPathFlag is the migration path flag
	MigrationPathFlag string = "migration-path"
	// MigrationWaitFlag is the migration wait flag
	MigrationWaitFlag string = "migration-wait"
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
	flag.StringP(MigrationManifestFlag, "m", "", "Path to the migrations manifest")
	flag.StringP(MigrationPathFlag, "p", "", "Semicolon-separated path to the migrations directories")
	flag.DurationP(MigrationWaitFlag, "w", time.Millisecond*10, "duration to wait when polling for new data from migration file")
}

// CheckMigration validates migration command line flags
func CheckMigration(v *viper.Viper) error {
	migrationPath := v.GetString(MigrationPathFlag)
	if len(migrationPath) == 0 {
		return errMissingMigrationPath
	}
	for _, p := range strings.Split(migrationPath, ";") {
		if strings.HasPrefix(p, "file://") {
			if _, err := os.Stat(p[len("file://"):]); os.IsNotExist(err) {
				return errors.Wrapf(&errInvalidMigrationPath{Path: p}, "Expected %s to exist", p)
			}
		}
		if strings.HasSuffix(p, "/") {
			return errors.Wrapf(&errInvalidMigrationPath{Path: p}, "Path %s Cannot end in slash", p)
		}
	}
	migrationManifest := v.GetString(MigrationManifestFlag)
	if len(migrationManifest) == 0 {
		return errMissingMigrationManifest
	}
	return nil
}
