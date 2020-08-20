package main

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
)

// ParseFlags parses the root command line flags
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
