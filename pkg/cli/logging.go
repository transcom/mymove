package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// LoggingEnvFlag is the logging environment flag
	LoggingEnvFlag string = "logging-env"
	// LogTaskMetadataFlag is the Log Task Metadata Flag
	LogTaskMetadataFlag string = "log-task-metadata"

	// LoggingEnvProduction is the production logging environment
	LoggingEnvProduction string = "production"
	// LoggingEnvDevelopment is the development logging environment
	LoggingEnvDevelopment string = "development"
)

var (
	allLoggingEnvs = []string{
		LoggingEnvProduction,
		LoggingEnvDevelopment,
	}
)

type errInvalidLoggingEnv struct {
	Value       string
	LoggingEnvs []string
}

func (e *errInvalidLoggingEnv) Error() string {
	return fmt.Sprintf("invalid logging env %s, must be one of: ", e.Value) + strings.Join(e.LoggingEnvs, ", ")
}

// InitLoggingFlags initializes the logging command line flags
func InitLoggingFlags(flag *pflag.FlagSet) {
	flag.String(LoggingEnvFlag, LoggingEnvDevelopment, "logging environment: "+strings.Join(allLoggingEnvs, ", "))
	flag.Bool(LogTaskMetadataFlag, false, "Fetch AWS Task Metadata and add to log lines.")
}

// CheckLogging validates logging command line flags
func CheckLogging(v *viper.Viper) error {
	if str := v.GetString(LoggingEnvFlag); !stringSliceContains(allLoggingEnvs, str) {
		return &errInvalidLoggingEnv{Value: str, LoggingEnvs: allLoggingEnvs}
	}
	return nil
}
