package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// HTTPSDDCServerNameFlag is the HTTP SDDC Server Name Flag
	HTTPSDDCServerNameFlag string = "http-sddc-server-name"
	// HTTPSDDCProtocolFlag is the HTTP SDDC Protocol Flag
	HTTPSDDCProtocolFlag string = "http-sddc-protocol"
	// HTTPSDDCPortFlag is the HTTP SDDC Port Flag
	HTTPSDDCPortFlag string = "http-sddc-port"
	// HTTPDPSServerNameFlag is the HTTP DPS Server Name Flag
	HTTPDPSServerNameFlag string = "http-dps-server-name"
	// DPSAuthSecretKeyFlag is the DPS Auth Secret Key Flag  G101
	DPSAuthSecretKeyFlag string = "dps-auth-secret-key"
	// DPSRedirectURLFlag is the DPS Redirect URL Flag
	DPSRedirectURLFlag string = "dps-redirect-url"
	// DPSCookieNameFlag is the DPS Cookie Name Flag
	DPSCookieNameFlag string = "dps-cookie-name"
	// DPSCookieDomainFlag is the DPS Cookie Domain Flag Flag
	DPSCookieDomainFlag string = "dps-cookie-domain"
	// DPSAuthCookieSecretKeyFlag is the DPS Auth Cookie Scret Key Flag  G101
	DPSAuthCookieSecretKeyFlag string = "dps-auth-cookie-secret-key"
	// DPSCookieExpiresInMinutesFlag is the DPS Cookie Expires In Minutes Flag
	DPSCookieExpiresInMinutesFlag string = "dps-cookie-expires-in-minutes"
)

// InitDPSFlags initializes the DPS command line flags
func InitDPSFlags(flag *pflag.FlagSet) {
	flag.String(HTTPSDDCServerNameFlag, "sddclocal", "Hostname according to envrionment.")
	flag.String(HTTPSDDCProtocolFlag, "https", "Protocol for sddc")
	flag.Int(HTTPSDDCPortFlag, 443, "The port for sddc")

	flag.String(HTTPDPSServerNameFlag, "dpslocal", "Hostname according to environment.")
	flag.String(DPSAuthSecretKeyFlag, "", "DPS auth JWT secret key")
	flag.String(DPSRedirectURLFlag, "", "DPS url to redirect to")
	flag.String(DPSCookieNameFlag, "", "Name of the DPS cookie")
	flag.String(DPSCookieDomainFlag, "sddclocal", "Domain of the DPS cookie")
	flag.String(DPSAuthCookieSecretKeyFlag, "", "DPS auth cookie secret key, 32 byte long")
	flag.Int(DPSCookieExpiresInMinutesFlag, 240, "DPS cookie expiration in minutes")
}

// CheckDPS validates DPS command line flags
func CheckDPS(v *viper.Viper) error {

	if err := ValidateProtocol(v, HTTPSDDCProtocolFlag); err != nil {
		return err
	}

	hostVars := []string{
		HTTPSDDCServerNameFlag,
		HTTPDPSServerNameFlag,
		DPSCookieDomainFlag,
	}

	for _, c := range hostVars {
		err := ValidateHost(v, c)
		if err != nil {
			return err
		}
	}

	if err := ValidatePort(v, HTTPSDDCPortFlag); err != nil {
		return err
	}

	dpsCookieSecret := []byte(v.GetString(DPSAuthCookieSecretKeyFlag))
	if len(dpsCookieSecret) != 32 {
		return errors.Errorf("DPS Cookie Secret Key is not 32 bytes. Cookie Secret Key length: %d", len(dpsCookieSecret))
	}

	return nil
}
