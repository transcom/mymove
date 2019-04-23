package cli

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dpsauth"
)

const (
	// HTTPSDDCServerNameFlag is the HTTP SDDC Server Name Flag
	HTTPSDDCServerNameFlag string = "http-sddc-server-name"
	// HTTPSDDCProtocolFlag is the HTTP SDDC Protocol Flag
	HTTPSDDCProtocolFlag string = "http-sddc-protocol"
	// HTTPSDDCPortFlag is the HTTP SDDC Port Flag
	HTTPSDDCPortFlag string = "http-sddc-port"
	// DPSAuthSecretKeyFlag is the DPS Auth Secret Key Flag
	DPSAuthSecretKeyFlag string = "dps-auth-secret-key"
	// DPSRedirectURLFlag is the DPS Redirect URL Flag
	DPSRedirectURLFlag string = "dps-redirect-url"
	// DPSCookieNameFlag is the DPS Cookie Name Flag
	DPSCookieNameFlag string = "dps-cookie-name"
	// DPSCookieDomainFlag is the DPS Cookie Domain Flag Flag
	DPSCookieDomainFlag string = "dps-cookie-domain"
	// DPSAuthCookieSecretKeyFlag is the DPS Auth Cookie Scret Key Flag
	DPSAuthCookieSecretKeyFlag string = "dps-auth-cookie-secret-key"
	// DPSCookieExpiresInMinutesFlag is the DPS Cookie Expires In Minutes Flag
	DPSCookieExpiresInMinutesFlag string = "dps-cookie-expires-in-minutes"
)

// InitDPSFlags initializes the DPS command line flags
func InitDPSFlags(flag *pflag.FlagSet) {
	flag.String(HTTPSDDCServerNameFlag, "sddclocal", "Hostname according to envrionment.")
	flag.String(HTTPSDDCProtocolFlag, "https", "Protocol for sddc")
	flag.Int(HTTPSDDCPortFlag, 443, "The port for sddc")
	flag.String(DPSAuthSecretKeyFlag, "", "DPS auth JWT secret key")
	flag.String(DPSRedirectURLFlag, "", "DPS url to redirect to")
	flag.String(DPSCookieNameFlag, "", "Name of the DPS cookie")
	flag.String(DPSCookieDomainFlag, "sddclocal", "Domain of the DPS cookie")
	flag.String(DPSAuthCookieSecretKeyFlag, "", "DPS auth cookie secret key, 32 byte long")
	flag.Int(DPSCookieExpiresInMinutesFlag, 240, "DPS cookie expiration in minutes")
}

// InitDPSAuthParams initializes the DPS Auth Params
func InitDPSAuthParams(v *viper.Viper, appnames auth.ApplicationServername) dpsauth.Params {
	dpsAuthSecretKey := v.GetString(DPSAuthSecretKeyFlag)
	dpsCookieDomain := v.GetString(DPSCookieDomainFlag)
	dpsCookieSecret := []byte(v.GetString(DPSAuthCookieSecretKeyFlag))
	dpsCookieExpires := v.GetInt(DPSCookieExpiresInMinutesFlag)
	return dpsauth.Params{
		SDDCProtocol:   v.GetString(HTTPSDDCProtocolFlag),
		SDDCHostname:   appnames.SddcServername,
		SDDCPort:       v.GetInt(HTTPSDDCPortFlag),
		SecretKey:      dpsAuthSecretKey,
		DPSRedirectURL: v.GetString(DPSRedirectURLFlag),
		CookieName:     v.GetString(DPSCookieNameFlag),
		CookieDomain:   dpsCookieDomain,
		CookieSecret:   dpsCookieSecret,
		CookieExpires:  dpsCookieExpires,
	}
}

// CheckDPS validates DPS command line flags
func CheckDPS(v *viper.Viper) error {

	dpsCookieSecret := []byte(v.GetString("dps-auth-cookie-secret-key"))
	if len(dpsCookieSecret) != 32 {
		return errors.New("DPS Cookie Secret Key is not 32 bytes. Cookie Secret Key length: " + strconv.Itoa(len(dpsCookieSecret)))
	}

	return nil
}
