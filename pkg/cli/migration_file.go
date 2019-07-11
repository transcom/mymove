package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// MigrationVersionFlag is the migration version flag
	MigrationVersionFlag string = "version"
	// MigrationNameFlag is the migration path flag
	MigrationNameFlag string = "name"
	// MigrationTypeFlag is the migration manifest flag
	MigrationTypeFlag string = "type"
	// VersionTimeFormat is the Go time format for creating a version number.
	VersionTimeFormat string = "20060102150405"
)

var (
	errMissingMigrationVersion = errors.New("missing migration version, expected to be set")
	errMissingMigrationName    = errors.New("missing migration name, expected to be set")
	errMissingMigrationType    = errors.New("missing migration type, expected to be set")
)

type errInvalidMigrationName struct {
	Value string
}

func (e *errInvalidMigrationName) Error() string {
	return fmt.Sprintf("invalid migration name %q, expecting no spaces", e.Value)
}

type errInvalidMigrationVersion struct {
	Value string
}

func (e *errInvalidMigrationVersion) Error() string {
	return fmt.Sprintf("invalid migration version %q, expecting an integer", e.Value)
}

type errInvalidMigrationType struct {
	Value string
}

func (e *errInvalidMigrationType) Error() string {
	return fmt.Sprintf("invalid migration type %q, expecting sql or fizz.", e.Value)
}

// InitMigrationFlags initializes the Migration command line flags
func InitMigrationFileFlags(flag *pflag.FlagSet) {
	flag.String(MigrationVersionFlag, time.Now().Format(VersionTimeFormat), "migration version: integer representation of datetime, default is current time using Go format "+VersionTimeFormat)
	flag.StringP(MigrationNameFlag, "n", "", "migration name: alphanumeric, no spaces, underscores and dashes allowed")
	flag.StringP(MigrationTypeFlag, "t", "fizz", "migration type: fizz or sql.")
}

// CheckMigration validates migration command line flags
func CheckMigrationFile(v *viper.Viper) error {
	migrationVersion := v.GetString(MigrationVersionFlag)
	if len(migrationVersion) == 0 {
		return errMissingMigrationVersion
	}
	if _, err := strconv.Atoi(migrationVersion); err != nil {
		return &errInvalidMigrationVersion{Value: migrationVersion}
	}
	migrationName := v.GetString(MigrationNameFlag)
	if len(migrationName) == 0 {
		return errMissingMigrationName
	}
	if strings.Contains(migrationName, " ") {
		return &errInvalidMigrationName{Value: migrationName}
	}
	migrationType := v.GetString(MigrationTypeFlag)
	if len(migrationType) == 0 {
		return errMissingMigrationType
	}
	if migrationType != "sql" && migrationType != "fizz" {
		return &errInvalidMigrationType{Value: migrationType}
	}
	return nil
}
