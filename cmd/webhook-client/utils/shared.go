package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Logger type exports the logger for use in the command files
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

// Root flags that may be used by any command
const (
	// CertPathFlag is the path to the client mTLS certificate
	CertPathFlag string = "certpath"
	// KeyPathFlag is the path to the key mTLS certificate
	KeyPathFlag string = "keypath"
	// HostnameFlag is the hostname to connect to
	HostnameFlag string = "hostname"
	// PortFlag is the port to connect to
	PortFlag string = "port"
	// InsecureFlag indicates that TLS verification and validation can be skipped
	InsecureFlag string = "insecure"
)

// ParseFlags parses the flags and binds both cobra and viper
func ParseFlags(cmd *cobra.Command, v *viper.Viper, args []string) error {

	errParseFlags := cmd.ParseFlags(args)
	if errParseFlags != nil {
		return fmt.Errorf("Could not parse args: %w", errParseFlags)
	}
	flags := cmd.Flags()
	errBindPFlags := v.BindPFlags(flags)
	if errBindPFlags != nil {
		return fmt.Errorf("Could not bind flags: %w", errBindPFlags)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	return nil
}

// DecodeJSONFileToPayload takes a filename, or stdin and decodes the file into
// the supplied json payload.
// If the filename is not supplied, the isStdin bool should be set to true to use stdin.
// If the file contains parameters that do not exist in the payload struct, it will fail with an error
// Otherwise it will populate the payload
func DecodeJSONFileToPayload(filename string, isStdin bool, payload interface{}) error {
	var reader *bufio.Reader
	if filename != "" {
		file, err := os.Open(filepath.Clean(filename))
		if err != nil {
			return fmt.Errorf("File open failed: %w", err)
		}
		reader = bufio.NewReader(file)
	} else if isStdin { // Uses std in if "-"" is provided instead
		reader = bufio.NewReader(os.Stdin)
	} else {
		return errors.New("no file input was found")
	}

	jsonDecoder := json.NewDecoder(reader)
	jsonDecoder.DisallowUnknownFields()

	// Read the json into the mto payload
	err := jsonDecoder.Decode(payload)
	if err != nil {
		return fmt.Errorf("File decode failed: %w", err)
	}

	return nil
}

// ContainsDash returns true if the original command included an empty dash
func ContainsDash(args []string) bool {
	for _, arg := range args {
		if arg == "-" {
			return true
		}
	}
	return false
}
