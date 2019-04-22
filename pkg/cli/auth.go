package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
)

const (
	// ClientAuthSecretKeyFlag is the Client Auth Secret Key Flag
	ClientAuthSecretKeyFlag string = "client-auth-secret-key"
	// NoSessionTimeoutFlag is the No Session Timeout Flag
	NoSessionTimeoutFlag string = "no-session-timeout"

	// LoginGovCallbackProtocolFlag is the Login.gov Callback Protocol Flag
	LoginGovCallbackProtocolFlag string = "login-gov-callback-protocol"
	// LoginGovCallbackPortFlag is the Login.gov Callback Port Flag
	LoginGovCallbackPortFlag string = "login-gov-callback-port"
	// LoginGovSecretKeyFlag is the Login.gov Secret Key Flag
	LoginGovSecretKeyFlag string = "login-gov-secret-key"
	// LoginGovMyClientIDFlag is the Login.gov My Client ID Flag
	LoginGovMyClientIDFlag string = "login-gov-my-client-id"
	// LoginGovOfficeClientIDFlag is the Login.gov Office Client ID Flag
	LoginGovOfficeClientIDFlag string = "login-gov-office-client-id"
	// LoginGovTSPClientIDFlag is the Login.gov TSP Client ID Flag
	LoginGovTSPClientIDFlag string = "login-gov-tsp-client-id"
	// LoginGovHostnameFlag is the Login.gov Hostname Flag
	LoginGovHostnameFlag string = "login-gov-hostname"
)

// InitAuthFlags initializes Auth command line flags
func InitAuthFlags(flag *pflag.FlagSet) {
	flag.String(ClientAuthSecretKeyFlag, "", "Client auth secret JWT key.")
	flag.Bool(NoSessionTimeoutFlag, false, "whether user sessions should timeout.")

	flag.String(LoginGovCallbackProtocolFlag, "https", "Protocol for non local environments.")
	flag.Int(LoginGovCallbackPortFlag, 443, "The port for callback urls.")
	flag.String(LoginGovSecretKeyFlag, "", "Login.gov auth secret JWT key.")
	flag.String(LoginGovMyClientIDFlag, "", "Client ID registered with login gov.")
	flag.String(LoginGovOfficeClientIDFlag, "", "Client ID registered with login gov.")
	flag.String(LoginGovTSPClientIDFlag, "", "Client ID registered with login gov.")
	flag.String(LoginGovHostnameFlag, "", "Hostname for communicating with login gov.")
}

// InitAuth initializes the Login.gov provider
func InitAuth(v *viper.Viper, logger logger, appnames auth.ApplicationServername) (authentication.LoginGovProvider, error) {
	loginGovCallbackProtocol := v.GetString(LoginGovCallbackProtocolFlag)
	loginGovCallbackPort := v.GetInt(LoginGovCallbackPortFlag)
	loginGovSecretKey := v.GetString(LoginGovSecretKeyFlag)
	loginGovHostname := v.GetString(LoginGovHostnameFlag)

	loginGovProvider := authentication.NewLoginGovProvider(loginGovHostname, loginGovSecretKey, logger)
	err := loginGovProvider.RegisterProvider(
		appnames.MilServername,
		v.GetString(LoginGovMyClientIDFlag),
		appnames.OfficeServername,
		v.GetString(LoginGovOfficeClientIDFlag),
		appnames.TspServername,
		v.GetString(LoginGovTSPClientIDFlag),
		loginGovCallbackProtocol,
		loginGovCallbackPort)
	return loginGovProvider, err
}
