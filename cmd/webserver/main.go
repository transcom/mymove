package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/pop"
	"github.com/gorilla/csrf"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	goji "goji.io"
	"goji.io/pat"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/dpsapi"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"github.com/transcom/mymove/pkg/handlers/ordersapi"
	"github.com/transcom/mymove/pkg/handlers/publicapi"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/storage"
)

// GitCommit is empty unless set as a build flag
// See https://blog.alexellis.io/inject-build-time-vars-golang/
var gitBranch string
var gitCommit string

// max request body size is 20 mb
const maxBodySize int64 = 200 * 1000 * 1000

type errInvalidProtocol struct {
	Protocol string
}

func (e *errInvalidProtocol) Error() string {
	return fmt.Sprintf("invalid protocol %s, must be http or https", e.Protocol)
}

type errInvalidPort struct {
	Port int
}

func (e *errInvalidPort) Error() string {
	return fmt.Sprintf("invalid port %d, must be > 0 and <= 65535", e.Port)
}

type errInvalidHost struct {
	Host string
}

func (e *errInvalidHost) Error() string {
	return fmt.Sprintf("invalid host %s, must not contain whitespace, :, /, or \\", e.Host)
}

type errInvalidRegion struct {
	Region string
}

func (e *errInvalidRegion) Error() string {
	return fmt.Sprintf("invalid region %s", e.Region)
}

type errInvalidPKCS7 struct {
	Path string
}

func (e *errInvalidPKCS7) Error() string {
	return fmt.Sprintf("invalid DER encoded PKCS7 package: %s", e.Path)
}

func stringSliceContains(stringSlice []string, value string) bool {
	for _, x := range stringSlice {
		if value == x {
			return true
		}
	}
	return false
}

func limitBodySizeMiddleware(inner http.Handler) http.Handler {
	zap.L().Debug("limitBodySizeMiddleware installed")
	mw := func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		inner.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

func noCacheMiddleware(inner http.Handler) http.Handler {
	zap.L().Debug("noCacheMiddleware installed")
	mw := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		inner.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

func httpsComplianceMiddleware(inner http.Handler) http.Handler {
	zap.L().Debug("httpsComplianceMiddleware installed")
	mw := func(w http.ResponseWriter, r *http.Request) {
		// set the HSTS header using values recommended by OWASP
		// https://www.owasp.org/index.php/HTTP_Strict_Transport_Security_Cheat_Sheet#Examples
		w.Header().Set("strict-transport-security", "max-age=31536000; includeSubdomains; preload")
		inner.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

func securityHeadersMiddleware(inner http.Handler) http.Handler {
	zap.L().Debug("securityHeadersMiddleware installed")
	mw := func(w http.ResponseWriter, r *http.Request) {
		// Sets headers to prevent rendering our page in an iframe, prevents clickjacking
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
		w.Header().Set("X-Frame-Options", "deny")
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/frame-ancestors
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		inner.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

func initFlags(flag *pflag.FlagSet) {

	flag.String("build", "build", "the directory to serve static files from.")
	flag.String("config-dir", "config", "The location of server config files")
	flag.String("env", "development", "The environment to run in, which configures the database.")
	flag.String("interface", "", "The interface spec to listen for connections on. Default is all.")
	flag.String("service-name", "app", "The service name identifies the application for instrumentation.")

	flag.String("http-my-server-name", "milmovelocal", "Hostname according to environment.")
	flag.String("http-office-server-name", "officelocal", "Hostname according to environment.")
	flag.String("http-tsp-server-name", "tsplocal", "Hostname according to environment.")
	flag.String("http-orders-server-name", "orderslocal", "Hostname according to environment.")
	flag.String("http-dps-server-name", "dpslocal", "Hostname according to environment.")

	// SDDC + DPS Auth config
	flag.String("http-sddc-server-name", "sddclocal", "Hostname according to envrionment.")
	flag.String("http-sddc-protocol", "https", "Protocol for sddc")
	flag.String("http-sddc-port", "", "The port for sddc")
	flag.String("dps-auth-secret-key", "", "DPS auth JWT secret key")
	flag.String("dps-redirect-url", "", "DPS url to redirect to")
	flag.String("dps-cookie-name", "", "Name of the DPS cookie")
	flag.String("dps-cookie-domain", "sddclocal", "Domain of the DPS cookie")
	flag.String("dps-auth-cookie-secret-key", "", "DPS auth cookie secret key, 32 byte long")
	flag.Int("dps-cookie-expires-in-minutes", 240, "DPS cookie expiration in minutes")

	// Initialize Swagger
	flag.String("swagger", "swagger/api.yaml", "The location of the public API swagger definition")
	flag.String("internal-swagger", "swagger/internal.yaml", "The location of the internal API swagger definition")
	flag.String("orders-swagger", "swagger/orders.yaml", "The location of the Orders API swagger definition")
	flag.String("dps-swagger", "swagger/dps.yaml", "The location of the DPS API swagger definition")

	flag.Bool("debug-logging", false, "log messages at the debug level.")
	flag.String("client-auth-secret-key", "", "Client auth secret JWT key.")
	flag.Bool("no-session-timeout", false, "whether user sessions should timeout.")

	flag.String("devlocal-ca", "", "Path to PEM-encoded devlocal CA certificate, enabled in development and test builds")
	flag.String("dod-ca-package", "", "Path to PKCS#7 package containing certificates of all DoD root and intermediate CAs")
	flag.String("move-mil-dod-ca-cert", "", "The DoD CA certificate used to sign the move.mil TLS certificate.")
	flag.String("move-mil-dod-tls-cert", "", "The DoD-signed TLS certificate for various move.mil services.")
	flag.String("move-mil-dod-tls-key", "", "The private key for the DoD-signed TLS certificate for various move.mil services.")

	// Ports to listen to
	flag.Int("mutual-tls-port", 9443, "The `port` for the mutual TLS listener.")
	flag.Int("tls-port", 8443, "the `port` for the server side TLS listener.")
	flag.Int("no-tls-port", 8080, "the `port` for the listener not requiring any TLS.")

	// Login.Gov config
	flag.String("login-gov-callback-protocol", "https", "Protocol for non local environments.")
	flag.Int("login-gov-callback-port", 443, "The port for callback urls.")
	flag.String("login-gov-secret-key", "", "Login.gov auth secret JWT key.")
	flag.String("login-gov-my-client-id", "", "Client ID registered with login gov.")
	flag.String("login-gov-office-client-id", "", "Client ID registered with login gov.")
	flag.String("login-gov-tsp-client-id", "", "Client ID registered with login gov.")
	flag.String("login-gov-hostname", "", "Hostname for communicating with login gov.")

	/* For bing Maps use the following
	bingMapsEndpoint := flag.String("bing_maps_endpoint", "", "URL for the Bing Maps Truck endpoint to use")
	bingMapsKey := flag.String("bing_maps_key", "", "Authentication key to use for the Bing Maps endpoint")
	*/

	// HERE Maps Config
	flag.String("here-maps-geocode-endpoint", "", "URL for the HERE maps geocoder endpoint")
	flag.String("here-maps-routing-endpoint", "", "URL for the HERE maps routing endpoint")
	flag.String("here-maps-app-id", "", "HERE maps App ID for this application")
	flag.String("here-maps-app-code", "", "HERE maps App API code")

	// EDI Invoice Config
	flag.String("gex-basic-auth-username", "", "GEX api auth username")
	flag.String("gex-basic-auth-password", "", "GEX api auth password")
	flag.Bool("send-prod-invoice", false, "Flag (bool) for EDI Invoices to signify if they should be sent with Production or Test indicator")
	flag.String("gex-url", "", "URL for sending an HTTP POST request to GEX")

	flag.String("storage-backend", "local", "Storage backend to use, either local, memory or s3.")
	flag.String("local-storage-root", "tmp", "Local storage root directory. Default is tmp.")
	flag.String("local-storage-web-root", "storage", "Local storage web root directory. Default is storage.")
	flag.String("email-backend", "local", "Email backend to use, either SES or local")
	flag.String("aws-s3-bucket-name", "", "S3 bucket used for file storage")
	flag.String("aws-s3-region", "", "AWS region used for S3 file storage")
	flag.String("aws-s3-key-namespace", "", "Key prefix for all objects written to S3")
	flag.String("aws-ses-region", "", "AWS region used for SES")
	flag.String("aws-ses-domain", "", "Domain used for SES")

	// Honeycomb Config
	flag.Bool("honeycomb-enabled", false, "Honeycomb enabled")
	flag.String("honeycomb-api-host", "https://api.honeycomb.io/", "API Host for Honeycomb")
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

	// CSRF Protection
	flag.String("csrf-auth-key", "", "CSRF Auth Key, 32 byte long")
}

func initDODCertificates(v *viper.Viper, logger *webserverLogger) ([]tls.Certificate, *x509.CertPool, error) {

	// https://tools.ietf.org/html/rfc7468#section-2
	//	- https://stackoverflow.com/questions/20173472/does-go-regexps-any-charcter-match-newline
	re := regexp.MustCompile("(?s)([-]{5}BEGIN CERTIFICATE[-]{5})(\\s*)(.+?)(\\s*)([-]{5}END CERTIFICATE[-]{5})")

	certFormat := "-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----"

	tlsCert := v.GetString("move-mil-dod-tls-cert")
	if len(tlsCert) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing", "move-mil-dod-tls-cert")
	}

	tlsCertMatches := re.FindAllStringSubmatch(tlsCert, -1)
	if len(tlsCertMatches) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing certificate PEM block", "move-mil-dod-tls-cert")
	}
	if len(tlsCertMatches) > 1 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s has too many certificate PEM blocks", "move-mil-dod-tls-cert")
	}

	tlsCerts := make([]string, 0, len(tlsCertMatches))
	for _, m := range tlsCertMatches {
		// each match will include a slice of strings starting with
		// (0) the full match, then
		// (1) "-----BEGIN CERTIFICATE-----",
		// (2) whitespace if any,
		// (3) base64-encoded certificate data,
		// (4) whitespace if any, and then
		// (5) -----END CERTIFICATE-----
		tlsCerts = append(tlsCerts, fmt.Sprintf(certFormat, m[3]))
	}

	logger.Info("certitficate chain from move-mil-dod-tls-cert parsed", zap.Any("count", len(tlsCerts)))

	caCert := v.GetString("move-mil-dod-ca-cert")
	if len(caCert) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing", "move-mil-dod-ca-cert")
	}

	caCertMatches := re.FindAllStringSubmatch(caCert, -1)
	if len(caCertMatches) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing certificate PEM block", "move-mil-dod-tls-cert")
	}

	caCerts := make([]string, 0, len(caCertMatches))
	for _, m := range caCertMatches {
		// each match will include a slice of strings starting with
		// (0) the full match, then
		// (1) "-----BEGIN CERTIFICATE-----",
		// (2) whitespace if any,
		// (3) base64-encoded certificate data,
		// (4) whitespace if any, and then
		// (5) -----END CERTIFICATE-----
		caCerts = append(caCerts, fmt.Sprintf(certFormat, m[3]))
	}

	logger.Info("certitficate chain from move-mil-dod-ca-cert parsed", zap.Any("count", len(caCerts)))

	//Append move.mil cert with intermediate CA to create a validate certificate chain
	cert := strings.Join(append(append(make([]string, 0), tlsCerts...), caCerts...), "\n")

	key := v.GetString("move-mil-dod-tls-key")
	if len(key) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Errorf("%s is missing", "move-mil-dod-tls-key")
	}

	keyPair, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, "failed to parse DOD x509 keypair for server")
	}

	logger.Info("DOD keypair", zap.Any("certificates", len(keyPair.Certificate)))

	pathToPackage := v.GetString("dod-ca-package")
	if len(pathToPackage) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Wrap(&errInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is missing", "dod-ca-package"))
	}

	pkcs7Package, err := ioutil.ReadFile(pathToPackage) // #nosec
	if err != nil {
		return make([]tls.Certificate, 0), nil, errors.Wrap(err, fmt.Sprintf("%s is invalid", "dod-ca-package"))
	}

	if len(pkcs7Package) == 0 {
		return make([]tls.Certificate, 0), nil, errors.Wrap(&errInvalidPKCS7{Path: pathToPackage}, fmt.Sprintf("%s is an empty file", "dod-ca-package"))
	}

	dodCACertPool, err := server.LoadCertPoolFromPkcs7Package(pkcs7Package)
	if err != nil {
		return make([]tls.Certificate, 0), dodCACertPool, errors.Wrap(err, "Failed to parse DoD CA certificate package")
	}

	return []tls.Certificate{keyPair}, dodCACertPool, nil

}

func initRoutePlanner(v *viper.Viper, logger *zap.Logger) route.Planner {
	return route.NewHEREPlanner(
		logger,
		v.GetString("here-maps-geocode-endpoint"),
		v.GetString("here-maps-routing-endpoint"),
		v.GetString("here-maps-app-id"),
		v.GetString("here-maps-app-code"))
}

func initHoneycomb(v *viper.Viper, logger *webserverLogger) bool {

	honeycombAPIHost := v.GetString("honeycomb-api-host")
	honeycombAPIKey := v.GetString("honeycomb-api-key")
	honeycombDataset := v.GetString("honeycomb-dataset")
	honeycombServiceName := v.GetString("service-name")

	if v.GetBool("honeycomb-enabled") && len(honeycombAPIKey) > 0 && len(honeycombDataset) > 0 && len(honeycombServiceName) > 0 {
		logger.Debug("Honeycomb Integration enabled",
			zap.String("honeycomb-api-host", honeycombAPIHost),
			zap.String("honeycomb-dataset", honeycombDataset))
		beeline.Init(beeline.Config{
			APIHost:     honeycombAPIHost,
			WriteKey:    honeycombAPIKey,
			Dataset:     honeycombDataset,
			Debug:       v.GetBool("honeycomb-debug"),
			ServiceName: honeycombServiceName,
		})
		return true
	}

	logger.Debug("Honeycomb Integration disabled")
	return false
}

func initRBSPersonLookup(v *viper.Viper, logger *webserverLogger) (*iws.RBSPersonLookup, error) {
	return iws.NewRBSPersonLookup(
		v.GetString("iws-rbs-host"),
		v.GetString("dod-ca-package"),
		v.GetString("move-mil-dod-tls-cert"),
		v.GetString("move-mil-dod-tls-key"))
}

func initDatabase(v *viper.Viper, logger *webserverLogger) (*pop.Connection, error) {

	env := v.GetString("env")
	dbName := v.GetString("db-name")
	dbHost := v.GetString("db-host")
	dbPort := strconv.Itoa(v.GetInt("db-port"))
	dbUser := v.GetString("db-user")
	dbPassword := v.GetString("db-password")

	// Modify DB options by environment
	dbOptions := map[string]string{"sslmode": "disable"}
	if env == "test" {
		// Leave the test database name hardcoded, since we run tests in the same
		// environment as development, and it's extra confusing to have to swap env
		// variables before running tests.
		dbName = "test_db"
	} else if env == "container" {
		// Require sslmode for containers
		dbOptions["sslmode"] = "require"
	}

	// Construct a safe URL and log it
	s := "postgres://%s:%s@%s:%s/%s?sslmode=%s"
	dbURL := fmt.Sprintf(s, dbUser, "*****", dbHost, dbPort, dbName, dbOptions["sslmode"])
	logger.Debug("Connecting to the database", zap.String("url", dbURL))

	// Configure DB connection details
	dbConnectionDetails := pop.ConnectionDetails{
		Dialect:  "postgres",
		Database: dbName,
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		Options:  dbOptions,
	}
	err := dbConnectionDetails.Finalize()
	if err != nil {
		logger.Error("Failed to finalize DB connection details", zap.Error(err))
		return nil, err
	}

	// Set up the connection
	connection, err := pop.NewConnection(&dbConnectionDetails)
	if err != nil {
		logger.Error("Failed create DB connection", zap.Error(err))
		return nil, err
	}

	// Open the connection
	err = connection.Open()
	if err != nil {
		logger.Error("Failed to open DB connection", zap.Error(err))
		return nil, err
	}

	// Check the connection
	db, err := sqlx.Open(connection.Dialect.Details().Dialect, connection.Dialect.URL())
	err = db.Ping()
	if err != nil {
		logger.Warn("Failed to ping DB connection", zap.Error(err))
		return connection, err
	}

	// Return the open connection
	return connection, nil
}

func checkConfig(v *viper.Viper) error {

	err := checkProtocols(v)
	if err != nil {
		return err
	}

	err = checkHosts(v)
	if err != nil {
		return err
	}

	err = checkPorts(v)
	if err != nil {
		return err
	}

	err = checkDPS(v)
	if err != nil {
		return err
	}

	err = checkCSRF(v)
	if err != nil {
		return err
	}

	err = checkEmail(v)
	if err != nil {
		return err
	}

	err = checkStorage(v)
	if err != nil {
		return err
	}

	err = checkGEX(v)
	if err != nil {
		return err
	}

	return nil
}

func checkProtocols(v *viper.Viper) error {

	protocolVars := []string{
		"login-gov-callback-protocol",
		"http-sddc-protocol",
	}

	for _, c := range protocolVars {
		if p := v.GetString(c); p != "http" && p != "https" {
			return errors.Wrap(&errInvalidProtocol{Protocol: p}, fmt.Sprintf("%s is invalid", c))
		}
	}

	return nil
}

func checkHosts(v *viper.Viper) error {
	invalidChars := ":/\\ \t\n\v\f\r"

	hostVars := []string{
		"http-my-server-name",
		"http-office-server-name",
		"http-tsp-server-name",
		"http-orders-server-name",
		"http-dps-server-name",
		"http-sddc-server-name",
		"dps-cookie-domain",
		"login-gov-hostname",
		"iws-rbs-host",
		"db-host",
	}

	for _, c := range hostVars {
		if h := v.GetString(c); len(h) == 0 || strings.ContainsAny(h, invalidChars) {
			return errors.Wrap(&errInvalidHost{Host: h}, fmt.Sprintf("%s is invalid", c))
		}
	}

	return nil
}

func checkPorts(v *viper.Viper) error {
	portVars := []string{
		"mutual-tls-port",
		"tls-port",
		"no-tls-port",
		"login-gov-callback-port",
		"db-port",
	}

	for _, c := range portVars {
		if p := v.GetInt(c); p <= 0 || p > 65535 {
			return errors.Wrap(&errInvalidPort{Port: p}, fmt.Sprintf("%s is invalid", c))
		}
	}

	return nil
}

func checkDPS(v *viper.Viper) error {

	dpsCookieSecret := []byte(v.GetString("dps-auth-cookie-secret-key"))
	if len(dpsCookieSecret) != 32 {
		return errors.New("DPS Cookie Secret Key is not 32 bytes. Cookie Secret Key length: " + strconv.Itoa(len(dpsCookieSecret)))
	}

	return nil
}

func checkCSRF(v *viper.Viper) error {

	csrfAuthKey, err := hex.DecodeString(v.GetString("csrf-auth-key"))
	if err != nil {
		return errors.Wrap(err, "Error decoding CSRF Auth Key")
	}
	if len(csrfAuthKey) != 32 {
		return errors.New("CSRF Auth Key is not 32 bytes. Auth Key length: " + strconv.Itoa(len(csrfAuthKey)))
	}

	return nil
}

func checkEmail(v *viper.Viper) error {
	emailBackend := v.GetString("email-backend")
	if !stringSliceContains([]string{"local", "ses"}, emailBackend) {
		return fmt.Errorf("invalid email-backend %s, expecting local or ses", emailBackend)
	}

	if emailBackend == "ses" {
		// SES is only available in 3 regions: us-east-1, us-west-2, and eu-west-1
		// - see https://docs.aws.amazon.com/ses/latest/DeveloperGuide/regions.html#region-endpoints
		if r := v.GetString("aws-ses-region"); len(r) == 0 || !stringSliceContains([]string{"us-east-1", "us-west-2", "eu-west-1"}, r) {
			return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-ses-region"))
		}
		if h := v.GetString("aws-ses-domain"); len(h) == 0 {
			return errors.Wrap(&errInvalidHost{Host: h}, fmt.Sprintf("%s is invalid", "aws-ses-domain"))
		}
	}

	return nil
}

func checkGEX(v *viper.Viper) error {
	gexURL := v.GetString("gex-url")
	if len(gexURL) > 0 && gexURL != "https://gexweba.daas.dla.mil/msg_data/submit/" {
		return fmt.Errorf("invalid gexUrl %s, expecting "+
			"https://gexweba.daas.dla.mil/msg_data/submit/ or an empty string", gexURL)
	}

	if len(gexURL) > 0 {
		if len(v.GetString("gex-basic-auth-username")) == 0 {
			return fmt.Errorf("GEX_BASIC_AUTH_USERNAME is missing")
		}
		if len(v.GetString("gex-basic-auth-password")) == 0 {
			return fmt.Errorf("GEX_BASIC_AUTH_PASSWORD is missing")
		}
	}

	return nil
}

func checkStorage(v *viper.Viper) error {

	storageBackend := v.GetString("storage-backend")
	if !stringSliceContains([]string{"local", "memory", "s3"}, storageBackend) {
		return fmt.Errorf("invalid storage-backend %s, expecting local, memory or s3", storageBackend)
	}

	if storageBackend == "s3" {
		regions, ok := endpoints.RegionsForService(endpoints.DefaultPartitions(), endpoints.AwsPartitionID, endpoints.S3ServiceID)
		if !ok {
			return fmt.Errorf("could not find regions for service %s", endpoints.S3ServiceID)
		}

		r := v.GetString("aws-s3-region")
		if len(r) == 0 {
			return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-s3-region"))
		}

		if _, ok := regions[r]; !ok {
			return errors.Wrap(&errInvalidRegion{Region: r}, fmt.Sprintf("%s is invalid", "aws-s3-region"))
		}
	} else if storageBackend == "local" {
		localStorageRoot := v.GetString("local-storage-root")
		if _, err := filepath.Abs(localStorageRoot); err != nil {
			return fmt.Errorf("could not get absolute path for %s", localStorageRoot)
		}
	}

	return nil
}

func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	env := v.GetString("env")
	isDevOrTest := env == "development" || env == "test"

	zapLogger, err := logging.Config(env, v.GetBool("debug-logging"))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(zapLogger)

	logger := &webserverLogger{zapLogger}

	logger.Debug("webserver starting up")

	err = checkConfig(v)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Honeycomb
	useHoneycomb := initHoneycomb(v, logger)

	clientAuthSecretKey := v.GetString("client-auth-secret-key")

	loginGovCallbackProtocol := v.GetString("login-gov-callback-protocol")
	loginGovCallbackPort := v.GetInt("login-gov-callback-port")
	loginGovSecretKey := v.GetString("login-gov-secret-key")
	loginGovHostname := v.GetString("login-gov-hostname")

	// Assert that our secret keys can be parsed into actual private keys
	// TODO: Store the parsed key in handlers/AppContext instead of parsing every time
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(loginGovSecretKey)); err != nil {
		logger.Fatal("Login.gov private key", zap.Error(err))
	}
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(clientAuthSecretKey)); err != nil {
		logger.Fatal("Client auth private key", zap.Error(err))
	}
	if len(loginGovHostname) == 0 {
		logger.Fatal("Must provide the Login.gov hostname parameter, exiting")
	}

	// Create a connection to the DB
	dbConnection, err := initDatabase(v, logger)
	if err != nil {
		if dbConnection == nil {
			// No connection object means that the configuraton failed to validate and we should kill server startup
			logger.Fatal("Connecting to DB", zap.Error(err))
		} else {
			// A valid connection object that still has an error indicates that the DB is not up but we
			// can proceed (this avoids a failure loop when deploying containers).
			logger.Warn("Starting server without DB connection")
		}
	}

	myHostname := v.GetString("http-my-server-name")
	officeHostname := v.GetString("http-office-server-name")
	tspHostname := v.GetString("http-tsp-server-name")

	// Register Login.gov authentication provider for My.(move.mil)
	loginGovProvider := authentication.NewLoginGovProvider(loginGovHostname, loginGovSecretKey, zapLogger)
	err = loginGovProvider.RegisterProvider(
		myHostname,
		v.GetString("login-gov-my-client-id"),
		officeHostname,
		v.GetString("login-gov-office-client-id"),
		tspHostname,
		v.GetString("login-gov-tsp-client-id"),
		loginGovCallbackProtocol,
		loginGovCallbackPort)
	if err != nil {
		logger.Fatal("Registering login provider", zap.Error(err))
	}

	// Session management and authentication middleware
	noSessionTimeout := v.GetBool("no-session-timeout")
	sessionCookieMiddleware := auth.SessionCookieMiddleware(zapLogger, clientAuthSecretKey, noSessionTimeout)
	maskedCSRFMiddleware := auth.MaskedCSRFMiddleware(zapLogger, noSessionTimeout)
	appDetectionMiddleware := auth.DetectorMiddleware(zapLogger, myHostname, officeHostname, tspHostname)
	userAuthMiddleware := authentication.UserAuthMiddleware(zapLogger)
	clientCertMiddleware := authentication.ClientCertMiddleware(zapLogger, dbConnection)

	handlerContext := handlers.NewHandlerContext(dbConnection, zapLogger)
	handlerContext.SetCookieSecret(clientAuthSecretKey)
	if noSessionTimeout {
		handlerContext.SetNoSessionTimeout()
	}

	if v.GetString("email-backend") == "ses" {
		// Setup Amazon SES (email) service
		// TODO: This might be able to be combined with the AWS Session that we're using for S3 down
		// below.
		awsSESRegion := v.GetString("aws-ses-region")
		awsSESDomain := v.GetString("aws-ses-domain")
		logger.Info("Using ses email backend",
			zap.String("region", awsSESRegion),
			zap.String("domain", awsSESDomain))
		sesSession, err := awssession.NewSession(&aws.Config{
			Region: aws.String(awsSESRegion),
		})
		if err != nil {
			logger.Fatal("Failed to create a new AWS client config provider", zap.Error(err))
		}
		sesService := ses.New(sesSession)
		handlerContext.SetNotificationSender(notifications.NewNotificationSender(sesService, awsSESDomain, zapLogger))
	} else {
		domain := "milmovelocal"
		logger.Info("Using local email backend", zap.String("domain", domain))
		handlerContext.SetNotificationSender(notifications.NewStubNotificationSender(domain, zapLogger))
	}

	build := v.GetString("build")

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(build))

	// Get route planner for handlers to calculate transit distances
	// routePlanner := route.NewBingPlanner(logger, bingMapsEndpoint, bingMapsKey)
	routePlanner := initRoutePlanner(v, zapLogger)
	handlerContext.SetPlanner(routePlanner)

	// Set SendProductionInvoice for ediinvoice
	handlerContext.SetSendProductionInvoice(v.GetBool("send-prod-invoice"))

	storageBackend := v.GetString("storage-backend")
	localStorageRoot := v.GetString("local-storage-root")
	localStorageWebRoot := v.GetString("local-storage-web-root")

	var storer storage.FileStorer
	if storageBackend == "s3" {
		awsS3Bucket := v.GetString("aws-s3-bucket-name")
		awsS3Region := v.GetString("aws-s3-region")
		awsS3KeyNamespace := v.GetString("aws-s3-key-namespace")
		logger.Info("Using s3 storage backend",
			zap.String("bucket", awsS3Bucket),
			zap.String("region", awsS3Region),
			zap.String("key", awsS3KeyNamespace))
		if len(awsS3Bucket) == 0 {
			logger.Fatal("must provide aws-s3-bucket-name parameter, exiting")
		}
		if len(awsS3Region) == 0 {
			logger.Fatal("Must provide aws-s3-region parameter, exiting")
		}
		if len(awsS3KeyNamespace) == 0 {
			logger.Fatal("Must provide aws_s3_key_namespace parameter, exiting")
		}
		aws := awssession.Must(awssession.NewSession(&aws.Config{
			Region: aws.String(awsS3Region),
		}))
		storer = storage.NewS3(awsS3Bucket, awsS3KeyNamespace, zapLogger, aws)
	} else if storageBackend == "memory" {
		logger.Info("Using memory storage backend",
			zap.String("root", path.Join(localStorageRoot, localStorageWebRoot)),
			zap.String("web root", localStorageWebRoot))
		fsParams := storage.NewMemoryParams(localStorageRoot, localStorageWebRoot, zapLogger)
		storer = storage.NewMemory(fsParams)
	} else {
		logger.Info("Using local storage backend",
			zap.String("root", path.Join(localStorageRoot, localStorageWebRoot)),
			zap.String("web root", localStorageWebRoot))
		fsParams := storage.NewFilesystemParams(localStorageRoot, localStorageWebRoot, zapLogger)
		storer = storage.NewFilesystem(fsParams)
	}
	handlerContext.SetFileStorer(storer)

	certificates, rootCAs, err := initDODCertificates(v, logger)
	if certificates == nil || rootCAs == nil || err != nil {
		logger.Fatal("Failed to initialize DOD certificates", zap.Error(err))
	}

	logger.Debug("Server DOD Key Pair Loaded")
	logger.Debug("Trusted Certificate Authorities", zap.Any("subjects", rootCAs.Subjects()))

	// Set the GexSender() and GexSender fields
	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs}
	var gexRequester services.GexSender
	gexURL := v.GetString("gex-url")
	if len(gexURL) == 0 {
		// this spins up a local test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		gexRequester = invoice.NewGexSenderHTTP(
			server.URL,
			false,
			&tls.Config{},
			"",
			"",
		)
	} else {
		gexRequester = invoice.NewGexSenderHTTP(
			v.GetString("gex-url"),
			true,
			tlsConfig,
			v.GetString("gex-basic-auth-username"),
			v.GetString("gex-basic-auth-password"),
		)
	}
	handlerContext.SetGexSender(gexRequester)

	// Set the ICNSequencer in the handler: if we are in dev/test mode and sending to a real
	// GEX URL, then we should use a random ICN number within a defined range to avoid duplicate
	// test ICNs in Syncada.
	var icnSequencer sequence.Sequencer
	if isDevOrTest && len(gexURL) > 0 {
		// ICNs are 9-digit numbers; reserve the ones in an upper range for development/testing.
		icnSequencer, err = sequence.NewRandomSequencer(ediinvoice.ICNRandomMin, ediinvoice.ICNRandomMax)
		if err != nil {
			logger.Fatal("Could not create random sequencer for ICN", zap.Error(err))
		}
	} else {
		icnSequencer = sequence.NewDatabaseSequencer(dbConnection, ediinvoice.ICNSequenceName)
	}
	handlerContext.SetICNSequencer(icnSequencer)

	rbs, err := initRBSPersonLookup(v, logger)
	if err != nil {
		logger.Fatal("Could not instantiate IWS RBS", zap.Error(err))
	}
	handlerContext.SetIWSPersonLookup(*rbs)

	sddcHostname := v.GetString("http-sddc-server-name")
	dpsAuthSecretKey := v.GetString("dps-auth-secret-key")
	dpsCookieDomain := v.GetString("dps-cookie-domain")
	dpsCookieSecret := []byte(v.GetString("dps-auth-cookie-secret-key"))
	dpsCookieExpires := v.GetInt("dps-cookie-expires-in-minutes")
	handlerContext.SetDPSAuthParams(
		dpsauth.Params{
			SDDCProtocol:   v.GetString("http-sddc-protocol"),
			SDDCHostname:   sddcHostname,
			SDDCPort:       v.GetString("http-sddc-port"),
			SecretKey:      dpsAuthSecretKey,
			DPSRedirectURL: v.GetString("dps-redirect-url"),
			CookieName:     v.GetString("dps-cookie-name"),
			CookieDomain:   dpsCookieDomain,
			CookieSecret:   dpsCookieSecret,
			CookieExpires:  dpsCookieExpires,
		},
	)

	// Base routes
	site := goji.NewMux()
	// Add middleware: they are evaluated in the reverse order in which they
	// are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	site.Use(httpsComplianceMiddleware)
	site.Use(securityHeadersMiddleware)
	site.Use(limitBodySizeMiddleware)

	// Stub health check
	site.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, r *http.Request) {

		data := map[string]interface{}{
			"gitBranch": gitBranch,
			"gitCommit": gitCommit,
		}

		// Check and see if we should disable DB query with '?database=false'
		// Disabling the DB is useful for Route53 health checks which require the TLS
		// handshake be less than 4 seconds and the status code return in less than
		// two seconds. https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/dns-failover-determining-health-of-endpoints.html
		showDB, ok := r.URL.Query()["database"]

		// Always show DB unless key set to "false"
		if !ok || (ok && showDB[0] != "false") {
			dbErr := dbConnection.RawQuery("SELECT 1;").Exec()
			if dbErr != nil {
				logger.Error("Failed database health check", zap.Error(dbErr))
			}
			data["database"] = dbErr == nil
		}

		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			logger.Error("Failed encoding health check response", zap.Error(err))
		}

		// We are not using request middleware here so logging directly in the check
		var protocol string
		if r.TLS == nil {
			protocol = "http"
		} else {
			protocol = "https"
		}
		zapLogger.Info("Request",
			zap.String("git-branch", gitBranch),
			zap.String("git-commit", gitCommit),
			zap.String("accepted-language", r.Header.Get("accepted-language")),
			zap.Int64("content-length", r.ContentLength),
			zap.String("host", r.Host),
			zap.String("method", r.Method),
			zap.String("protocol", protocol),
			zap.String("protocol-version", r.Proto),
			zap.String("referer", r.Header.Get("referer")),
			zap.String("source", r.RemoteAddr),
			zap.String("url", r.URL.String()),
			zap.String("user-agent", r.UserAgent()),
			zap.String("x-amzn-trace-id", r.Header.Get("x-amzn-trace-id")),
			zap.String("x-forwarded-for", r.Header.Get("x-forwarded-for")),
			zap.String("x-forwarded-host", r.Header.Get("x-forwarded-host")),
			zap.String("x-forwarded-proto", r.Header.Get("x-forwarded-proto")),
		)
	})

	// Allow public content through without any auth or app checks
	site.Handle(pat.Get("/static/*"), clientHandler)
	site.Handle(pat.Get("/swagger-ui/*"), clientHandler)
	site.Handle(pat.Get("/downloads/*"), clientHandler)
	site.Handle(pat.Get("/favicon.ico"), clientHandler)

	ordersMux := goji.SubMux()
	ordersDetectionMiddleware := auth.HostnameDetectorMiddleware(zapLogger, v.GetString("http-orders-server-name"))
	ordersMux.Use(ordersDetectionMiddleware)
	ordersMux.Use(noCacheMiddleware)
	ordersMux.Use(clientCertMiddleware)
	ordersMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString("orders-swagger")))
	ordersMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "orders.html")))
	ordersMux.Handle(pat.New("/*"), ordersapi.NewOrdersAPIHandler(handlerContext))
	site.Handle(pat.Get("/orders/v0/*"), ordersMux)

	dpsMux := goji.SubMux()
	dpsDetectionMiddleware := auth.HostnameDetectorMiddleware(zapLogger, v.GetString("http-dps-server-name"))
	dpsMux.Use(dpsDetectionMiddleware)
	dpsMux.Use(noCacheMiddleware)
	dpsMux.Use(clientCertMiddleware)
	dpsMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString("dps-swagger")))
	dpsMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "dps.html")))
	dpsMux.Handle(pat.New("/*"), dpsapi.NewDPSAPIHandler(handlerContext))
	site.Handle(pat.New("/dps/v0/*"), dpsMux)

	sddcDPSMux := goji.SubMux()
	sddcDetectionMiddleware := auth.HostnameDetectorMiddleware(zapLogger, sddcHostname)
	sddcDPSMux.Use(sddcDetectionMiddleware)
	sddcDPSMux.Use(noCacheMiddleware)
	site.Handle(pat.New("/dps_auth/*"), sddcDPSMux)
	sddcDPSMux.Handle(pat.Get("/set_cookie"),
		dpsauth.NewSetCookieHandler(zapLogger,
			dpsAuthSecretKey,
			dpsCookieDomain,
			dpsCookieSecret,
			dpsCookieExpires))

	root := goji.NewMux()
	root.Use(sessionCookieMiddleware)
	root.Use(appDetectionMiddleware) // Comes after the sessionCookieMiddleware as it sets session state
	root.Use(logging.LogRequestMiddleware(gitBranch, gitCommit))

	// CSRF path is set specifically at the root to avoid duplicate tokens from different paths
	csrfAuthKey, err := hex.DecodeString(v.GetString("csrf-auth-key"))
	if err != nil {
		logger.Fatal("Failed to decode csrf auth key", zap.Error(err))
	}
	logger.Info("Enabling CSRF protection")
	root.Use(csrf.Protect(csrfAuthKey, csrf.Secure(!isDevOrTest), csrf.Path("/")))
	root.Use(maskedCSRFMiddleware)

	// Sends build variables to honeycomb
	if len(gitBranch) > 0 && len(gitCommit) > 0 {
		root.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx, span := beeline.StartSpan(r.Context(), "BuildVariablesMiddleware")
				defer span.Send()
				span.AddTraceField("git.branch", gitBranch)
				span.AddTraceField("git.commit", gitCommit)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
	}
	site.Handle(pat.New("/*"), root)

	apiMux := goji.SubMux()
	root.Handle(pat.New("/api/v1/*"), apiMux)
	apiMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString("swagger")))
	apiMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "api.html")))

	externalAPIMux := goji.SubMux()
	apiMux.Handle(pat.New("/*"), externalAPIMux)
	externalAPIMux.Use(noCacheMiddleware)
	externalAPIMux.Use(userAuthMiddleware)
	externalAPIMux.Handle(pat.New("/*"), publicapi.NewPublicAPIHandler(handlerContext))

	internalMux := goji.SubMux()
	root.Handle(pat.New("/internal/*"), internalMux)
	internalMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString("internal-swagger")))
	internalMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "internal.html")))

	// Mux for internal API that enforces auth
	internalAPIMux := goji.SubMux()
	internalMux.Handle(pat.New("/*"), internalAPIMux)
	internalAPIMux.Use(userAuthMiddleware)
	internalAPIMux.Use(noCacheMiddleware)
	internalAPIMux.Handle(pat.New("/*"), internalapi.NewInternalAPIHandler(handlerContext))

	authContext := authentication.NewAuthContext(zapLogger, loginGovProvider, loginGovCallbackProtocol, loginGovCallbackPort)
	authMux := goji.SubMux()
	root.Handle(pat.New("/auth/*"), authMux)
	authMux.Handle(pat.Get("/login-gov"), authentication.RedirectHandler{Context: authContext})
	authMux.Handle(pat.Get("/login-gov/callback"), authentication.NewCallbackHandler(authContext, dbConnection, clientAuthSecretKey, noSessionTimeout))
	authMux.Handle(pat.Get("/logout"), authentication.NewLogoutHandler(authContext, clientAuthSecretKey, noSessionTimeout))

	if isDevOrTest {
		logger.Info("Enabling devlocal auth")
		localAuthMux := goji.SubMux()
		root.Handle(pat.New("/devlocal-auth/*"), localAuthMux)
		localAuthMux.Handle(pat.Get("/login"), authentication.NewUserListHandler(authContext, dbConnection))
		localAuthMux.Handle(pat.Post("/login"), authentication.NewAssignUserHandler(authContext, dbConnection, clientAuthSecretKey, noSessionTimeout))
		localAuthMux.Handle(pat.Post("/new"), authentication.NewCreateAndLoginUserHandler(authContext, dbConnection, clientAuthSecretKey, noSessionTimeout))
		localAuthMux.Handle(pat.Post("/create"), authentication.NewCreateUserHandler(authContext, dbConnection, clientAuthSecretKey, noSessionTimeout))

		devlocalCa, err := ioutil.ReadFile(v.GetString("devlocal-ca")) // #nosec
		if err != nil {
			logger.Error("No devlocal CA path defined")
		} else {
			rootCAs.AppendCertsFromPEM(devlocalCa)
		}
	}

	if storageBackend == "local" {
		// Add a file handler to provide access to files uploaded in development
		fs := storage.NewFilesystemHandler(localStorageRoot)
		root.Handle(pat.Get(path.Join("/", localStorageWebRoot, "/*")), fs)
	}

	// Serve index.html to all requests that haven't matches a previous route,
	root.HandleFunc(pat.Get("/*"), indexHandler(build, zapLogger))

	var httpHandler http.Handler
	if useHoneycomb {
		httpHandler = hnynethttp.WrapHandler(site)
	} else {
		httpHandler = site
	}

	errChan := make(chan error)

	listenInterface := v.GetString("interface")

	go func() {
		noTLSServer := server.Server{
			ListenAddress: listenInterface,
			HTTPHandler:   httpHandler,
			Logger:        zapLogger,
			Port:          v.GetInt("no-tls-port"),
		}
		errChan <- noTLSServer.ListenAndServe()
	}()

	go func() {
		tlsServer := server.Server{
			ClientAuthType: tls.NoClientCert,
			ListenAddress:  listenInterface,
			HTTPHandler:    httpHandler,
			Logger:         zapLogger,
			Port:           v.GetInt("tls-port"),
			TLSCerts:       certificates,
		}
		errChan <- tlsServer.ListenAndServeTLS()
	}()

	go func() {
		mutualTLSServer := server.Server{
			// Ensure that any DoD-signed client certificate can be validated,
			// using the package of DoD root and intermediate CAs provided by DISA
			CaCertPool:     rootCAs,
			ClientAuthType: tls.RequireAndVerifyClientCert,
			ListenAddress:  listenInterface,
			HTTPHandler:    httpHandler,
			Logger:         zapLogger,
			Port:           v.GetInt("mutual-tls-port"),
			TLSCerts:       certificates,
		}
		errChan <- mutualTLSServer.ListenAndServeTLS()
	}()

	logger.Fatal("listener error", zap.Error(<-errChan))
}

// fileHandler serves up a single file
func fileHandler(entrypoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
}

// indexHandler returns a handler that will serve the resulting content
func indexHandler(buildDir string, logger *zap.Logger) http.HandlerFunc {

	indexPath := path.Join(buildDir, "index.html")
	// #nosec - indexPath does not come from user input
	indexHTML, err := ioutil.ReadFile(indexPath)
	if err != nil {
		logger.Fatal("could not read index.html template: run make client_build", zap.Error(err))
	}

	stat, err := os.Stat(indexPath)
	if err != nil {
		logger.Fatal("could not stat index.html template", zap.Error(err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "index.html", stat.ModTime(), bytes.NewReader(indexHTML))
	}
}
