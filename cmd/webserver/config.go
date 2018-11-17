package main

import (
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"os"
	"strings"

	"go.uber.org/dig"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/transcom/mymove/pkg/authentication"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/storage"
)

func initFlags(flag *pflag.FlagSet) {

	flag.String("build", "build", "the directory to serve static files from.")
	flag.String("config-dir", "config", "The location of server config files")
	flag.String("env", "development", "The environment to run in, which configures the database.")
	flag.String("interface", "", "The interface spec to listen for connections on. Default is all.")
	flag.String("service-name", "app", "The service name identifies the application for instrumentation.")

	flag.String("http-my-server-name", "localhost", "Hostname according to environment.")
	flag.String("http-office-server-name", "officelocal", "Hostname according to environment.")
	flag.String("http-tsp-server-name", "tsplocal", "Hostname according to environment.")
	flag.String("http-orders-server-name", "orderslocal", "Hostname according to environment.")
	flag.String("http-dps-server-name", "dpslocal", "Hostname according to environment.")

	// Initialize Swagger
	flag.String("swagger", "swagger/api.yaml", "The location of the public API swagger definition")
	flag.String("internal-swagger", "swagger/internal.yaml", "The location of the internal API swagger definition")
	flag.String("orders-swagger", "swagger/orders.yaml", "The location of the Orders API swagger definition")
	flag.String("dps-swagger", "swagger/dps.yaml", "The location of the DPS API swagger definition")

	// SDDC + DPS Auth config
	flag.String("http-sddc-server-name", "sddclocal", "Hostname according to envrionment.")
	flag.String("http-sddc-protocol", "https", "Protocol for sddc")
	flag.String("http-sddc-port", "", "The port for sddc")
	flag.String("dps-auth-secret-key", "", "DPS auth JWT secret key")
	flag.String("dps-redirect-url", "", "DPS url to redirect to")
	flag.String("dps-cookie-name", "", "Name of the DPS cookie")
	flag.String("dps-cookie-domain", "sddclocal", "Domain of the DPS cookie")

	flag.Bool("debug-logging", false, "log messages at the debug level.")
	flag.String("client-auth-secret-key", "", "Client auth secret JWT key.")
	flag.Bool("no-session-timeout", false, "whether user sessions should timeout.")

	flag.String("dod-ca-package", "", "Path to PKCS#7 package containing certificates of all DoD root and intermediate CAs")
	flag.String("move-mil-dod-ca-cert", "", "The DoD CA certificate used to sign the move.mil TLS certificate.")
	flag.String("move-mil-dod-tls-cert", "", "The DoD-signed TLS certificate for various move.mil services.")
	flag.String("move-mil-dod-tls-key", "", "The private key for the DoD-signed TLS certificate for various move.mil services.")

	// Ports to listen to
	flag.Int("mutual-tls-port", 9443, "The `port` for the mutual TLS listener.")
	flag.Int("tls-port", 8443, "the `port` for the server side TLS listener.")
	flag.Int("no-tls-port", 8080, "the `port` for the listener not requiring any TLS.")

	// Login.Gov config
	flag.String("login-gov-callback-protocol", "https://", "Protocol for non local environments.")
	flag.Int("login-gov-callback-port", 443, "The port for callback urls.")
	flag.String("login-gov-secret-key", "", "Login.gov auth secret JWT key.")
	flag.String("login-gov-my-client-id", "", "Client ID registered with login gov.")
	flag.String("login-gov-office-client-id", "", "Client ID registered with login gov.")
	flag.String("login-gov-tsp-client-id", "", "Client ID registered with login gov.")
	flag.String("login-gov-hostname", "", "Hostname for communicating with login gov.")

	flag.String("bing_maps_endpoint", "", "URL for the Bing Maps Truck endpoint to use")
	flag.String("bing_maps_key", "", "Authentication key to use for the Bing Maps endpoint")

	// HERE Maps Config
	flag.String("here-maps-geocode-endpoint", "", "URL for the HERE maps geocoder endpoint")
	flag.String("here-maps-routing-endpoint", "", "URL for the HERE maps routing endpoint")
	flag.String("here-maps-app-id", "", "HERE maps App ID for this application")
	flag.String("here-maps-app-code", "", "HERE maps App API code")

	flag.String("storage-backend", "filesystem", "Storage backend to use, either filesystem or s3.")
	flag.String("email-backend", "local", "Email backend to use, either SES or local")
	flag.String("aws-s3-bucket-name", "", "S3 bucket used for file storage")
	flag.String("aws-s3-region", "", "AWS region used for S3 file storage")
	flag.String("aws-s3-key-namespace", "", "Key prefix for all objects written to S3")
	flag.String("aws-ses-region", "", "AWS region used for SES")

	// New Relic Config
	flag.String("new-relic-application-id", "", "App ID for New Relic Browser")
	flag.String("new-relic-license-key", "", "License key for New Relic Browser")

	// Honeycomb Config
	flag.Bool("honeycomb-enabled", false, "Honeycomb enabled")
	flag.String("honeycomb-api-key", "", "API Key for Honeycomb")
	flag.String("honeycomb-dataset", "", "Dataset for Honeycomb")
	flag.Bool("honeycomb-debug", false, "Debug honeycomb using stdout.")

	// IWS
	flag.String("iws-rbs-host", "", "Hostname for the IWS RBS")

	// DB Config
	flag.String("db-name", "dev_db", "Database Name")
	flag.String("db-host", "localhost", "Database Hostname")
	flag.Int("db-port", 5432, "Database Port")
	flag.String("db-user", "postgres", "Database Username")
	flag.String("db-password", "", "Database Password")
}

func parseConfig() *viper.Viper {
	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	return v
}

// DatabaseConfig contains where to find per environment configs and the environment name
type DatabaseConfig struct {
	ConfigDir   string
	Environment string
}

// SwaggerConfig contains names of the various swagger yaml files
type SwaggerConfig struct {
	Internal string
	API      string
	Orders   string
	DPS      string
}

// NewRelicConfig contains the App ID and Key for New Relic
type NewRelicConfig struct {
	AppID string
	Key   string
}

// ListenerConfig contains configuration for the various HTTP(S) listeners
type ListenerConfig struct {
	NoTLSPort        int    // Port with no TLS
	TLSPort          int    // Port for regular TLS access
	MutualTLSPort    int    // Port for TLS with client certs
	DoDCACert        string // The DoD CA certificate used to sign the move.mil TLS certificates
	DoDTLSCert       string // The DoD signed tls certificate for various move.mil services
	DoDTLSKey        string // The DoD signed tls key for various move.mil services
	DoDCACertPackage string // Path to PKCS#7 package containing certificates of all DoD root and intermediate CAs
}

// WebServerConfig rolls up the various bits of config, so parseConfig provider has a sensible return value
type WebServerConfig struct {
	dig.Out
	DB             *DatabaseConfig
	Hosts          *server.HostsConfig
	Cookie         *server.SessionCookieConfig
	Swagger        *SwaggerConfig
	Here           *route.HEREConfig
	SesSender      *notifications.SESNotificationConfig
	S3Config       *storage.S3StorerConfig
	EnvConfig      *server.LocalEnvConfig
	NewRelicConfig *NewRelicConfig
	LoginGovConfig *authentication.LoginGovConfig
	TLSConfig      *ListenerConfig
	DPSAuthParams  *dpsauth.Params
	handlers.SendProdInvoice
}

/*
parseConfig parses the config for the MyMoveMil web server.
*/
func serverConfig(cfg *viper.Viper) WebServerConfig {

	var sesConfig *notifications.SESNotificationConfig
	if cfg.GetString("email-backend") == "ses" {
		sesConfig = &notifications.SESNotificationConfig{Region: cfg.GetString("aws-ses-region")}
	}

	var s3Config *storage.S3StorerConfig
	if cfg.GetString("storage-backend") == "s3" {
		s3Config = &storage.S3StorerConfig{
			Bucket:       cfg.GetString("aws-s3-bucket-name"),
			Region:       cfg.GetString("aws-s3-region"),
			KeyNamespace: cfg.GetString("aws-s3-key-namespace"),
		}
	}

	return WebServerConfig{
		Out: dig.Out{},
		DB: &DatabaseConfig{
			ConfigDir:   cfg.GetString("config-dir"),
			Environment: cfg.GetString("env"),
		},
		Hosts: &server.HostsConfig{
			ListenInterface: cfg.GetString("interface"),
			MyName:          cfg.GetString("http-my-server-name"),
			OfficeName:      cfg.GetString("http-office-server-name"),
			TspName:         cfg.GetString("http-tsp-server-name"),
			OrdersName:      cfg.GetString("http-orders-server-name"),
			DPSName:         cfg.GetString("http-dps-server-name"),
		},
		Cookie: &server.SessionCookieConfig{
			Secret:    cfg.GetString("client-auth-secret-key"),
			NoTimeout: cfg.GetBool("no-session-timeout"),
		},
		Swagger: &SwaggerConfig{
			Internal: cfg.GetString("internal-swagger"),
			API:      cfg.GetString("swagger"),
			Orders:   cfg.GetString("orders-swagger"),
			DPS:      cfg.GetString("dps-swagger"),
		},
		Here: &route.HEREConfig{
			RouteEndpoint:   cfg.GetString("here-maps-routing-endpoint"),
			GeocodeEndpoint: cfg.GetString("here-maps-geocode-endpoint"),
			AppCode:         cfg.GetString("here-maps-app-code"),
			AppID:           cfg.GetString("here-maps-app-id"),
		},
		SesSender: sesConfig,
		S3Config:  s3Config,
		EnvConfig: &server.LocalEnvConfig{
			Environment: cfg.GetString("env"),
			SiteDir:     cfg.GetString("build"),
			ConfigDir:   cfg.GetString("config-dir"),
		},
		NewRelicConfig: &NewRelicConfig{
			AppID: cfg.GetString("new-relic-application-id"),
			Key:   cfg.GetString("new-relic-license-key"),
		},
		LoginGovConfig: &authentication.LoginGovConfig{
			Host:             cfg.GetString("login-gov-hostname"),
			CallbackProtocol: cfg.GetString("login-gov-callback-protocol"),
			CallbackPort:     cfg.GetInt("login-gov-callback-port"),
			MyClientID:       cfg.GetString("login-gov-my-client-id"),
			OfficeClientID:   cfg.GetString("login-gov-office-client-id"),
			TspClientID:      cfg.GetString("login-gov-tsp-client-id"),
			Secret:           cfg.GetString("login-gov-secret-key"),
		},
		TLSConfig: &ListenerConfig{
			NoTLSPort:        cfg.GetInt("no-tls-port"),
			TLSPort:          cfg.GetInt("tls-port"),
			MutualTLSPort:    cfg.GetInt("mutual-tls-port"),
			DoDCACert:        cfg.GetString("move-mil-dod-ca-cert"),
			DoDTLSCert:       cfg.GetString("move-mil-dod-tls-cert"),
			DoDTLSKey:        cfg.GetString("move-mil-dod-tls-key"),
			DoDCACertPackage: cfg.GetString("dod-ca-package"),
		},
		DPSAuthParams: &dpsauth.Params{
			SDDCProtocol:   cfg.GetString("http-sddc-protocol"),
			SDDCHostname:   cfg.GetString("http-sddc-server-name"),
			SDDCPort:       cfg.GetString("http-sddc-port"),
			SecretKey:      cfg.GetString("dps-auth-secret-key"),
			DPSRedirectURL: cfg.GetString("dps-redirect-url"),
			CookieName:     cfg.GetString("dps-cookie-name"),
		},
		SendProdInvoice: handlers.SendProdInvoice(cfg.GetBool("send-prod-invoice")),
	}
}
