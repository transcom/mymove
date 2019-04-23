package cli

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type errInvalidRegion struct {
	Region string
}

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %s", e.Region)
}

type errInvalidProtocol struct {
	Protocol string
}

func (e *errInvalidProtocol) Error() string {
	return fmt.Sprintf("invalid protocol %s, must be http or https", e.Protocol)
}

type errInvalidURL struct {
	URL string
}

func (e *errInvalidURL) Error() string {
	return fmt.Sprintf("invalid url %s", e.URL)
}

func stringSliceContains(stringSlice []string, value string) bool {
	for _, x := range stringSlice {
		if value == x {
			return true
		}
	}
	return false
}

// ValidateProtocol validates a Protocol passed in from the command line
func ValidateProtocol(v *viper.Viper, flagname string) error {
	if p := v.GetString(flagname); p != "http" && p != "https" {
		return errors.Wrap(&errInvalidProtocol{Protocol: p}, fmt.Sprintf("%s is invalid", flagname))
	}
	return nil
}

// ValidateURL validates a URL passed in from the command line
func ValidateURL(v *viper.Viper, flagname string) error {
	endpoint := v.GetString(flagname)
	if _, err := url.Parse(endpoint); endpoint != "" && err != nil {
		return errors.Wrapf(err, "Unable to parse %s endpoint %s", flagname, endpoint)
	}
	return nil
}
