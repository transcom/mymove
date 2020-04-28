package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// FilenameFlag is the name of the file being passed in
	FilenameFlag string = "filename"
	// ETagFlag is the etag for the mto shipment being updated
	ETagFlag string = "etag"
)

func containsDash(args []string) bool {
	for _, arg := range args {
		if arg == "-" {
			return true
		}
	}
	return false
}

// decodeJSONFileToPayload takes a filename, or if not available
func decodeJSONFileToPayload(filename string, isStdin bool, payload interface{}) error {
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
