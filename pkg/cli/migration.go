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
	// MigrationPathFlag is the migration path flag
	MigrationPathFlag string = "migration-path"
	// MigrationManifestFlag is the migration manifest flag
	MigrationManifestFlag string = "migration-manifest"
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
	flag.StringP(MigrationPathFlag, "p", "./migrations", "Path to the migrations folder")
	flag.StringP(MigrationManifestFlag, "m", "./migrations_manifest.txt", "Path to the manifest")
	flag.DurationP(MigrationWaitFlag, "w", time.Millisecond*10, "duration to wait when polling for new data from migration file")
}

// CheckMigration validates migration command line flags
func CheckMigration(v *viper.Viper) error {
	migrationPath := v.GetString(MigrationPathFlag)
	if len(migrationPath) == 0 {
		return errMissingMigrationPath
	}
	for _, p := range strings.Split(migrationPath, ";") {
		if len(p) == 0 {
			continue
		}
		if strings.HasPrefix(p, "file://") {
			filesystemPath := p[len("file://"):]
			if _, err := os.Stat(filesystemPath); os.IsNotExist(err) {
				return errors.Wrapf(&errInvalidMigrationPath{Path: filesystemPath}, "Expected %s to exist", filesystemPath)
			}
		} else if !strings.HasPrefix(p, "s3://") {
			return errors.Wrapf(&errInvalidMigrationPath{Path: p}, "Expected %s to have prefix file:// or s3://", p)
		}
		if strings.HasSuffix(p, "/") {
			return errors.Wrapf(&errInvalidMigrationPath{Path: p}, "Path %s Cannot end in slash", p)
		}
	}
	migrationManifest := v.GetString(MigrationManifestFlag)
	if len(migrationManifest) == 0 {
		return errMissingMigrationManifest
	}
	if len(MigrationManifestFlag) == 0 {
		return errMissingMigrationManifest
	}
	return nil
}
