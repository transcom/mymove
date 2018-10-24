package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/di"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/storage"
	"go.uber.org/dig"
	"log"
)

// HoneycombConfig contains is configuration for connecting to Honeycomb service
type HoneycombConfig struct {
	Debug        *bool
	Enabled      *bool
	APIKey       *string
	DataSet      *string
	UseHoneycomb bool
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
}

// NewRelicConfig contains the App ID and Key for New Relic
type NewRelicConfig struct {
	AppID string
	Key   string
}

// ListenerConfig contains configuration for the various HTTP(S) listeners
type ListenerConfig struct {
	NoTLSPort     string // Port with no TLS
	TLSPort       string // Port for regular TLS access
	MutualTLSPort string // Port for TLS with client certs
	DoDCACert     string // The DoD CA certificate used to sign the move.mil TLS certificates
	DoDTLSCert    string // The DoD signed tls certificate for various move.mil services
	DoDTLSKey     string // The DoD signed tls key for various move.mil services
}

// WebServerConfig rolls up the various bits of config, so parseConfig provider has a sensible return value
type WebServerConfig struct {
	dig.Out
	Logger         *di.Config
	Honeycomb      *HoneycombConfig
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
}

/*
ParseConfig parses the config for the MyMoveMil web server.
TODO. This code should really live with the main webserver in /cmd/webserver but I couldn't work out how to get packages to work there
*/
func ParseConfig() WebServerConfig {

	// FOR NOW. PatrickD's viper proposal should hopefully simplify this
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	build := flag.String("build", "build", "the directory to serve static files from.")
	config := flag.String("config-dir", "config", "The location of server config files")

	debugLogging := flag.Bool("debug_logging", false, "log messages at the debug level.")

	listenInterface := flag.String("interface", "", "The interface spec to listen for connections on. Default is all.")
	myHostname := flag.String("http_my_server_name", "localhost", "Hostname according to environment.")
	officeHostname := flag.String("http_office_server_name", "officelocal", "Hostname according to environment.")
	tspHostname := flag.String("http_tsp_server_name", "tsplocal", "Hostname according to environment.")
	ordersHostname := flag.String("http_orders_server_name", "orderslocal", "Hostname according to environment.")

	clientAuthSecretKey := flag.String("client_auth_secret_key", "", "Client auth secret JWT key.")
	noSessionTimeout := flag.Bool("no_session_timeout", false, "whether user sessions should timeout.")

	internalSwagger := flag.String("internal-swagger", "swagger/internal.yaml", "The location of the internal API swagger definition")
	apiSwagger := flag.String("swagger", "swagger/api.yaml", "The location of the public API swagger definition")
	ordersSwagger := flag.String("orders-swagger", "swagger/orders.yaml", "The location of the Orders API swagger definition")

	moveMilDODCACert := flag.String("move_mil_dod_ca_cert", "", "The DoD CA certificate used to sign the move.mil TLS certificates.")
	moveMilDODTLSCert := flag.String("move_mil_dod_tls_cert", "", "the DoD signed tls certificate for various move.mil services.")
	moveMilDODTLSKey := flag.String("move_mil_dod_tls_key", "", "the DoD signed tls key for various move.mil services.")

	mutualTLSPort := flag.String("mutual_tls_port", "9443", "The `port` for the mutual TLS listener.")
	tlsPort := flag.String("tls_port", "8443", "the `port` for the server side TLS listener.")
	noTLSPort := flag.String("no_tls_port", "8080", "the `port` for the listener not requiring any TLS.")

	loginGovCallbackProtocol := flag.String("login_gov_callback_protocol", "https://", "Protocol for non local environments.")
	loginGovCallbackPort := flag.String("login_gov_callback_port", "443", "The port for callback urls.")
	loginGovSecretKey := flag.String("login_gov_secret_key", "", "Login.gov auth secret JWT key.")
	loginGovMyClientID := flag.String("login_gov_my_client_id", "", "Client ID registered with login gov.")
	loginGovOfficeClientID := flag.String("login_gov_office_client_id", "", "Client ID registered with login gov.")
	loginGovTSPClientID := flag.String("login_gov_tsp_client_id", "", "Client ID registered with login gov.")
	loginGovHostname := flag.String("login_gov_hostname", "", "Hostname for communicating with login gov.")

	/* For bing Maps use the following
	bingMapsEndpoint := flag.String("bing_maps_endpoint", "", "URL for the Bing Maps Truck endpoint to use")
	bingMapsKey := flag.String("bing_maps_key", "", "Authentication key to use for the Bing Maps endpoint")
	*/
	hereGeoEndpoint := flag.String("here_maps_geocode_endpoint", "", "URL for the HERE maps geocoder endpoint")
	hereRouteEndpoint := flag.String("here_maps_routing_endpoint", "", "URL for the HERE maps routing endpoint")
	hereAppID := flag.String("here_maps_app_id", "", "HERE maps App ID for this application")
	hereAppCode := flag.String("here_maps_app_code", "", "HERE maps App API code")

	storageBackend := flag.String("storage_backend", "filesystem", "Storage backend to use, either filesystem or s3.")

	s3Bucket := flag.String("aws_s3_bucket_name", "", "S3 bucket used for file storage")
	s3Region := flag.String("aws_s3_region", "", "AWS region used for S3 file storage")
	s3KeyNamespace := flag.String("aws_s3_key_namespace", "", "Key prefix for all objects written to S3")

	awsSesRegion := flag.String("aws_ses_region", "", "AWS region used for SES")
	emailBackend := flag.String("email_backend", "local", "Email backend to use, either SES or local")

	newRelicApplicationID := flag.String("new_relic_application_id", "", "App ID for New Relic Browser")
	newRelicLicenseKey := flag.String("new_relic_license_key", "", "License key for New Relic Browser")

	honeyConfig := HoneycombConfig{
		flag.Bool("honeycomb_debug", false, "Debug Honeycomb using stdout."),
		flag.Bool("honeycomb_enabled", false, "Honeycomb enabled"),
		flag.String("honeycomb_api_key", "", "API Key for Honeycomb"),
		flag.String("honeycomb_dataset", "", "Dataset for Honeycomb"),
		false,
	}

	flag.Parse()

	if *loginGovHostname == "" {
		log.Fatal("Must provide the Login.gov hostname parameter, exiting")
	}

	var sesConfig *notifications.SESNotificationConfig
	if *emailBackend == "ses" {
		sesConfig = &notifications.SESNotificationConfig{Config: aws.Config{Region: aws.String(*awsSesRegion)}}
	}

	var s3Config *storage.S3StorerConfig
	if *storageBackend == "s3" {
		s3Config = &storage.S3StorerConfig{Bucket: *s3Bucket, Region: *s3Region, KeyNamespace: *s3KeyNamespace}
	}
	return WebServerConfig{
		Out: dig.Out{},
		Logger: &di.Config{
			DebugLogging: *debugLogging},
		Honeycomb: &honeyConfig,
		DB: &DatabaseConfig{
			*config,
			*env},
		Hosts: &server.HostsConfig{
			ListenInterface: *listenInterface,
			MyName:          *myHostname,
			OfficeName:      *officeHostname,
			TspName:         *tspHostname,
			OrdersName:      *ordersHostname},
		Cookie: &server.SessionCookieConfig{
			Secret:    *clientAuthSecretKey,
			NoTimeout: *noSessionTimeout,
		},
		Swagger: &SwaggerConfig{
			*internalSwagger,
			*apiSwagger,
			*ordersSwagger,
		},
		Here: &route.HEREConfig{
			RouteEndpoint:   *hereRouteEndpoint,
			GeocodeEndpoint: *hereGeoEndpoint,
			AppCode:         *hereAppCode,
			AppID:           *hereAppID,
		},
		SesSender: sesConfig,
		S3Config:  s3Config,
		EnvConfig: &server.LocalEnvConfig{
			Environment: *env,
			SiteDir:     *build,
			ConfigDir:   *config,
		},
		NewRelicConfig: &NewRelicConfig{
			AppID: *newRelicApplicationID,
			Key:   *newRelicLicenseKey,
		},
		LoginGovConfig: &authentication.LoginGovConfig{
			Host:             *loginGovHostname,
			CallbackProtocol: *loginGovCallbackProtocol,
			CallbackPort:     *loginGovCallbackPort,
			MyClientID:       *loginGovMyClientID,
			OfficeClientID:   *loginGovOfficeClientID,
			TspClientID:      *loginGovTSPClientID,
			Secret:           *loginGovSecretKey,
		},
		TLSConfig: &ListenerConfig{
			NoTLSPort:     *noTLSPort,
			TLSPort:       *tlsPort,
			MutualTLSPort: *mutualTLSPort,
			DoDCACert:     *moveMilDODCACert,
			DoDTLSCert:    *moveMilDODTLSCert,
			DoDTLSKey:     *moveMilDODTLSKey,
		},
	}
}
