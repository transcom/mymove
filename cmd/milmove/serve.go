package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/transcom/mymove/pkg/models"
	"golang.org/x/crypto/ocsp"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gobuffalo/pop/v6"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/ecs"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/handlers/routing"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/telemetry"
)

// initServeFlags - Order matters!
func initServeFlags(flag *pflag.FlagSet) {

	// Environment
	cli.InitEnvironmentFlags(flag)

	// Build Files
	cli.InitBuildFlags(flag)

	// Webserver
	cli.InitWebserverFlags(flag)

	// Hosts
	cli.InitHostFlags(flag)

	// Initialize Swagger
	cli.InitSwaggerFlags(flag)

	// Certs
	cli.InitCertFlags(flag)

	// Ports to listen to
	cli.InitPortFlags(flag)

	// Enable listeners
	cli.InitListenerFlags(flag)

	// Login.Gov Auth config
	cli.InitAuthFlags(flag)

	// Devlocal Auth config
	cli.InitDevlocalFlags(flag)

	// HERE Route Config
	cli.InitRouteFlags(flag)

	// EDI Invoice Config
	cli.InitGEXFlags(flag)

	// Storage
	cli.InitStorageFlags(flag)

	// Email
	cli.InitEmailFlags(flag)

	// IWS
	cli.InitIWSFlags(flag)

	// DB Config
	cli.InitDatabaseFlags(flag)

	// CSRF Protection
	cli.InitCSRFFlags(flag)

	// Middleware
	cli.InitMiddlewareFlags(flag)

	// Logging
	cli.InitLoggingFlags(flag)

	// pprof flags
	cli.InitDebugFlags(flag)

	// Service Flags
	cli.InitServiceFlags(flag)

	// Redis Flags
	cli.InitRedisFlags(flag)

	// SessionFlags
	cli.InitSessionFlags(flag)

	// Telemetry flag config
	cli.InitTelemetryFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

func checkServeConfig(v *viper.Viper, logger *zap.Logger) error {

	logger.Info("checking webserver config")

	if err := cli.CheckEnvironment(v); err != nil {
		logger.Info(fmt.Sprintf("Environment check failed: %v", err.Error()))
	}

	if err := cli.CheckBuild(v); err != nil {
		return err
	}

	if err := cli.CheckWebserver(v); err != nil {
		return err
	}

	if err := cli.CheckHosts(v); err != nil {
		return err
	}

	if err := cli.CheckSwagger(v); err != nil {
		return err
	}

	if err := cli.CheckCert(v); err != nil {
		return err
	}

	if err := cli.CheckListeners(v); err != nil {
		return err
	}

	if err := cli.CheckPorts(v); err != nil {
		return err
	}

	if err := cli.CheckAuth(v); err != nil {
		return err
	}

	if err := cli.CheckDevlocal(v); err != nil {
		return err
	}

	if err := cli.CheckRoute(v); err != nil {
		return err
	}

	if err := cli.CheckGEX(v); err != nil {
		return err
	}

	if err := cli.CheckStorage(v); err != nil {
		return err
	}

	if err := cli.CheckEmail(v); err != nil {
		return err
	}

	if err := cli.CheckIWS(v); err != nil {
		return err
	}

	if err := cli.CheckDatabase(v, logger); err != nil {
		return err
	}

	if err := cli.CheckCSRF(v); err != nil {
		return err
	}

	if err := cli.CheckMiddleWare(v); err != nil {
		return err
	}

	if err := cli.CheckLogging(v); err != nil {
		return err
	}

	if err := cli.CheckDebugFlags(v); err != nil {
		return err
	}

	if err := cli.CheckServices(v); err != nil {
		return err
	}

	if err := cli.CheckRedis(v); err != nil {
		return err
	}

	return cli.CheckSession(v)
}

func startListener(srv *server.NamedServer, logger *zap.Logger, useTLS bool) {
	logger.Info("Starting listener",
		zap.String("name", srv.Name),
		zap.Duration("idle-timeout", srv.IdleTimeout),
		zap.Any("listen-address", srv.Addr),
		zap.Int("max-header-bytes", srv.MaxHeaderBytes),
		zap.Int("port", srv.Port()),
		zap.Bool("tls", useTLS),
	)
	var err error
	if useTLS {
		err = srv.ListenAndServeTLS()
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		logger.Fatal("server error", zap.String("name", srv.Name), zap.Error(err))
	}
}

func initializeViper(cmd *cobra.Command, args []string) (*viper.Viper, error) {
	err := cmd.ParseFlags(args)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse flags")
	}

	flag := cmd.Flags()

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		return nil, errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	return v, nil
}

func initializeLogger(v *viper.Viper) (*zap.Logger, func()) {
	logger, loggerSync, err := logging.Config(
		logging.WithEnvironment(v.GetString(cli.LoggingEnvFlag)),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)

	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	fields := make([]zap.Field, 0)
	if len(gitBranch) > 0 {
		fields = append(fields, zap.String("git_branch", gitBranch))
	}
	if len(gitCommit) > 0 {
		fields = append(fields, zap.String("git_commit", gitCommit))
	}
	logger = logger.With(fields...)

	if v.GetBool(cli.LogTaskMetadataFlag) {
		// according to
		// https://docs.aws.amazon.com/AmazonECS/latest/userguide/task-metadata-endpoint-v4-fargate.html
		//
		//     Beginning with Fargate platform version 1.4.0, an
		//     environment variable named
		//     ECS_CONTAINER_METADATA_URI_V4 is injected into each
		//     container in a task
		metadataURL := os.Getenv("ECS_CONTAINER_METADATA_URI_V4")
		if metadataURL != "" {
			var ecsTaskMetadataV4 ecs.TaskMetadataV4
			r, gerr := http.Get(metadataURL + "/task")
			if gerr != nil {
				logger.Error("Cannot fetch v4 task metadata", zap.Error(gerr))
			} else {
				derr := json.NewDecoder(r.Body).Decode(&ecsTaskMetadataV4)
				if derr != nil {
					logger.Error("Cannot decode v4 task metadata", zap.Error(derr))
				} else {
					logger.Info("V4 Task", zap.Any("metadata", ecsTaskMetadataV4))
					logger = logger.With(
						zap.String("ecs_cluster", ecsTaskMetadataV4.Cluster),
						zap.String("ecs_task_def_family", ecsTaskMetadataV4.Family),
						zap.String("ecs_task_def_revision", ecsTaskMetadataV4.Revision),
					)
				}
			}
		}
	}

	zap.ReplaceGlobals(logger)

	return logger, loggerSync
}

func initializeAwsSession(v *viper.Viper, logger *zap.Logger) *awssession.Session {
	if v.GetBool(cli.DbIamFlag) || (v.GetString(cli.EmailBackendFlag) == "ses") || (v.GetString(cli.StorageBackendFlag) == "s3") {
		c := &aws.Config{
			Region: aws.String(v.GetString(cli.AWSRegionFlag)),
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
		}
		return s
	}
	return nil
}

func initializeDB(v *viper.Viper, logger *zap.Logger,
	awsSession *awssession.Session) *pop.Connection {

	if v.GetBool(cli.DbDebugFlag) {
		pop.Debug = true
	}

	var dbCreds *credentials.Credentials
	if v.GetBool(cli.DbIamFlag) {
		if awsSession != nil {
			// We want to get the credentials from the logged in AWS session rather than create directly,
			// because the session conflates the environment, shared, and container metdata config
			// within NewSession.  With stscreds, we use the Secure Token Service,
			// to assume the given role (that has rds db connect permissions).
			dbIamRole := v.GetString(cli.DbIamRoleFlag)
			logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))
			dbCreds = stscreds.NewCredentials(awsSession, dbIamRole)
			stsService := sts.New(awsSession)
			callerIdentity, callerIdentityErr := stsService.GetCallerIdentity(&sts.GetCallerIdentityInput{})
			if callerIdentityErr != nil {
				logger.Error(errors.Wrap(callerIdentityErr, "error getting aws sts caller identity").Error())
			} else {
				logger.Info(fmt.Sprintf("STS Caller Identity - Account: %s, ARN: %s, UserId: %s", *callerIdentity.Account, *callerIdentity.Arn, *callerIdentity.UserId))
			}
		}
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, dbCreds, logger)
	if err != nil {
		logger.Fatal("Invalid DB Configuration", zap.Error(err))
	}

	err = cli.PingPopConnection(dbConnection, logger)
	if err != nil {
		// if the db is not up yet, the server can still start. This
		// prevents a failure loop when deploying containers
		logger.Warn("DB is not ready for connections", zap.Error(err))
	}

	return dbConnection
}

func initializeTLSConfig(appCtx appcontext.AppContext, v *viper.Viper) *tls.Config {
	certificates, rootCAs, err := certs.InitDoDCertificates(v, appCtx.Logger()) //initialize DOD config
	if certificates == nil || rootCAs == nil || err != nil {
		appCtx.Logger().Fatal("Failed to initialize DOD certificates", zap.Error(err))
	}

	useDevlocalAuthCA := stringSliceContains([]string{cli.EnvironmentTest, cli.EnvironmentDevelopment, cli.EnvironmentReview, cli.EnvironmentLoadtest}, v.GetString(cli.EnvironmentFlag))
	devlocalCAPath := v.GetString(cli.DevlocalCAFlag)
	if useDevlocalAuthCA && devlocalCAPath != "" {
		appCtx.Logger().Info("Adding devlocal CA to root CAs")
		devlocalCa, readFileErr := os.ReadFile(filepath.Clean(devlocalCAPath))
		if readFileErr != nil {
			appCtx.Logger().Error(fmt.Sprintf("Unable to read devlocal CA from path %s", devlocalCAPath), zap.Error(readFileErr))
		} else {
			rootCAs.AppendCertsFromPEM(devlocalCa)
		}
	}
	// RA Summary: staticcheck - SA1019 - Using a deprecated function, variable, constant or field
	// RA: Linter is flagging: rootCAs.Subjects is deprecated: if s was returned by SystemCertPool, Subjects will not include the system roots.
	// RA: Why code valuable: It allows us to log the root CA subjects that are being trusted.
	// RA: Mitigation: The deprecation notes this is a problem when reading SystemCertPool, but we do not use this here and are building our own cert pool instead.
	// RA Developer Status: Mitigated
	// RA Validator Status: Mitigated
	// RA Validator: leodis.f.scott.civ@mail.mil
	// RA Modified Severity: CAT III
	// nolint:staticcheck
	subjects := rootCAs.Subjects()
	appCtx.Logger().Info("Trusted CAs", zap.Any("num", len(subjects)), zap.Any("subjects", subjects))

	return &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}
}

// HTTP client
// Hard failure - failure to check revocation list status due to network or other failures.
var failToCheckRevList bool = false

func certRevokedCheck(appCtx appcontext.AppContext, clientCert *x509.Certificate) (error) {
	ocspResponse, err := sendOCSPRequestAndGetResponse(clientCert.Subject.CommonName, clientCert.Issuer, clientCert.OCSPServer[0])
	if err != nil {
		return err
	}
	switch ocspResponse.Status {
		case ocsp.Good:
			fmt.Printf("[+] Certificate status is Good\n")
		case ocsp.Revoked:
			fmt.Printf("[-] Certificate status is Revoked\n")
			return fmt.Errorf("The certificate was revoked!")
		case ocsp.Unknown:
			fmt.Printf("[-] Certificate status is Unknown\n")
			return fmt.Errorf("The certificate is unknown to OCSP server!")
		}

		fmt.Printf("Server certificate was allowed\n")
	return nil
	}
	//if revoked, ok, err := certRevokedOCSP(clientCert, failToCheckRevList); !ok {
	//	appCtx.Logger().Info("error checking revocation via OCSP")
	//	if failToCheckRevList {
	//		return true, false, err
	//	}
	//	return false, false, err
	//} else if revoked {
	//	appCtx.Logger().Info("certificate is revoked via OCSP")
	//	return true, true, err
	//}

	return false, true, nil
}

func getCRL(url string) (rawCerts [][]byte, verifiedChains [][]*x509.Certificate) {
	//userCert := verifiedChains[0][0]
	//issuerCert := verifiedChains[0][1]
	//responseBody := [][]byte //TODO: read response body here
	//return x509.ParseCRL(responseBody) // NOTE: ParseCRL is deprecated, use something else to parse
	return nil, nil
}

func getIssuer(issuerCert *x509.Certificate) *x509.Certificate {
	return issuer
}

// Check cert against a specific CRL
func certRevokedCRL(cert *x509.Certificate, url string) (revoked, ok bool, err error) {
	// Add check that cer hasn't expired
	return false, true, err
}

func certRevokedOCSP(cert *x509.Certificate, strict bool) (revoked, ok bool, err error) {
	return revoked, ok, err
}

// Request OCSP response from server.
// Returns error if the server can't get the certificate
func sendOCSPRequestAndGetResponse(appCtx appcontext.AppContext, ocspServer string, request *http.Request, responseWriter http.ResponseWriter, clientCert, issuerCert *x509.Certificate) (*ocsp.Response, error) {
	var ocspRead = io.ReadAll
	// NOTE: http requests must be made with TLS
	newAppCtx := appcontext.NewAppContextFromContext(request.Context(), appCtx)
	hash := sha256.Sum256(request.TLS.PeerCertificates[0].Raw)
	hashString := hex.EncodeToString(hash[:])

	clientCert, err := models.FetchClientCert(newAppCtx.DB(), hashString)
	if err != nil {
		// This is not a known client certificate at all
		newAppCtx.Logger().Info(
			"Unknown / unregistered client certificate",
			zap.String("SHA256_hash_string", hashString),
		)
		http.Error(responseWriter, http.StatusText(401), http.StatusUnauthorized)
		return nil, nil
	}

	ocspRequestOpts := &ocsp.RequestOptions{Hash: crypto.SHA256}
	buffer, err := ocsp.CreateRequest(clientCert, issuerCert, ocspRequestOpts)
	if err != nil {
		return nil, err
	}
	httpRequest, err := http.NewRequest(http.MethodPost, ocspServer, bytes.NewBuffer(buffer))
	if err != nil {
		return nil, err
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	body, err := ocspRead(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	ocspResponse, err := ocsp.ParseResponseForCert(body, clientCert, issuerCert)
	//return ocsp.ParseResponseForCert(body, leafCertificate, issuerCert)
	return ocspResponse, err
}

//func getOCSPResponse(commonName string, clientCert, issuerCert *x509.Certificate, ocspServerURL string) (*ocsp.Response, error) {
// // OCSP Request
//}

func fetchRemoteCert(url string) (*x509.Certificate, error) {
	response, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
}

func initializeRouteOptions(v *viper.Viper, routingConfig *routing.Config) {
	routingConfig.MaxBodySize = v.GetInt64(cli.MaxBodySizeFlag)
	routingConfig.ServeSwaggerUI = v.GetBool(cli.ServeSwaggerUIFlag)
	routingConfig.ServeOrders = v.GetBool(cli.ServeOrdersFlag)
	if routingConfig.ServeOrders {
		routingConfig.OrdersSwaggerPath = v.GetString(cli.OrdersSwaggerFlag)
	}
	routingConfig.ServePrime = v.GetBool(cli.ServePrimeFlag)
	routingConfig.ServePrimeSimulator = v.GetBool(cli.ServePrimeSimulatorFlag)
	if routingConfig.ServePrime || routingConfig.ServePrimeSimulator {
		routingConfig.PrimeSwaggerPath = v.GetString(cli.PrimeSwaggerFlag)
		routingConfig.PrimeV2SwaggerPath = v.GetString(cli.PrimeV2SwaggerFlag)
	}

	routingConfig.ServeSupport = v.GetBool(cli.ServeSupportFlag)
	if routingConfig.ServeSupport {
		routingConfig.SupportSwaggerPath = v.GetString(cli.SupportSwaggerFlag)
	}
	routingConfig.ServeDebugPProf = v.GetBool(cli.DebugPProfFlag)
	routingConfig.ServeAPIInternal = v.GetBool(cli.ServeAPIInternalFlag)
	if routingConfig.ServeAPIInternal {
		routingConfig.APIInternalSwaggerPath = v.GetString(cli.InternalSwaggerFlag)
	}
	routingConfig.ServeAdmin = v.GetBool(cli.ServeAdminFlag)
	if routingConfig.ServeAdmin {
		routingConfig.AdminSwaggerPath = v.GetString(cli.AdminSwaggerFlag)
	}
	routingConfig.ServeGHC = v.GetBool(cli.ServeGHCFlag)
	if routingConfig.ServeGHC {
		routingConfig.GHCSwaggerPath = v.GetString(cli.GHCSwaggerFlag)
	}
	routingConfig.ServeDevlocalAuth = v.GetBool(cli.DevlocalAuthFlag)

	routingConfig.GitBranch = gitBranch
	routingConfig.GitCommit = gitCommit

	csrfAuthKey, err := hex.DecodeString(v.GetString(cli.CSRFAuthKeyFlag))
	if err == nil {
		routingConfig.CSRFMiddleware = routing.InitCSRFMiddlware(csrfAuthKey, routingConfig.HandlerConfig.UseSecureCookie(), "/", auth.GorillaCSRFToken)
	}
}

func buildRoutingConfig(appCtx appcontext.AppContext, v *viper.Viper, redisPool *redis.Pool, awsSession *awssession.Session, isDevOrTest bool, tlsConfig *tls.Config) *routing.Config {
	routingConfig := &routing.Config{}

	// always use the OS Filesystem when serving for real
	routingConfig.FileSystem = afero.NewOsFs()

	// Collect the servernames into a handy struct
	appNames := auth.ApplicationServername{
		MilServername:    v.GetString(cli.HTTPMyServerNameFlag),
		OfficeServername: v.GetString(cli.HTTPOfficeServerNameFlag),
		AdminServername:  v.GetString(cli.HTTPAdminServerNameFlag),
		OrdersServername: v.GetString(cli.HTTPOrdersServerNameFlag),
		PrimeServername:  v.GetString(cli.HTTPPrimeServerNameFlag),
	}

	clientAuthSecretKey := v.GetString(cli.ClientAuthSecretKeyFlag)
	loginGovCallbackProtocol := v.GetString(cli.LoginGovCallbackProtocolFlag)
	loginGovCallbackPort := v.GetInt(cli.LoginGovCallbackPortFlag)
	loginGovSecretKey := v.GetString(cli.LoginGovSecretKeyFlag)
	loginGovHostname := v.GetString(cli.LoginGovHostnameFlag)

	// Assert that our secret keys can be parsed into actual private keys
	// TODO: Store the parsed key in handlers/AppContext instead of parsing every time
	if _, parseRSAPrivateKeyFromPEMErr := jwt.ParseRSAPrivateKeyFromPEM([]byte(loginGovSecretKey)); parseRSAPrivateKeyFromPEMErr != nil {
		appCtx.Logger().Fatal("Login.gov private key", zap.Error(parseRSAPrivateKeyFromPEMErr))
	}
	if _, parseRSAPrivateKeyFromPEMErr := jwt.ParseRSAPrivateKeyFromPEM([]byte(clientAuthSecretKey)); parseRSAPrivateKeyFromPEMErr != nil {
		appCtx.Logger().Fatal("Client auth private key", zap.Error(parseRSAPrivateKeyFromPEMErr))
	}
	if len(loginGovHostname) == 0 {
		appCtx.Logger().Fatal("Must provide the Login.gov hostname parameter, exiting")
	}

	// Register Login.gov authentication provider for My.(move.mil)
	loginGovProvider, err := authentication.InitAuth(v, appCtx.Logger(),
		appNames)
	if err != nil {
		appCtx.Logger().Fatal("Registering login provider", zap.Error(err))
	}

	routingConfig.AuthContext = authentication.NewAuthContext(appCtx.Logger(), loginGovProvider, loginGovCallbackProtocol, loginGovCallbackPort)

	// Email
	notificationSender, err := notifications.InitEmail(v, awsSession, appCtx.Logger())
	if err != nil {
		appCtx.Logger().Fatal("notification sender sending not enabled", zap.Error(err))
	}

	routingConfig.BuildRoot = v.GetString(cli.BuildRootFlag)
	sendProductionInvoice := v.GetBool(cli.GEXSendProdInvoiceFlag)

	// Storage
	fileStorer := storage.InitStorage(v, awsSession, appCtx.Logger())

	// Create a secondary planner specifically for HHG.
	hhgRoutePlanner, err := route.InitHHGRoutePlanner(appCtx, v, tlsConfig)
	if err != nil {
		appCtx.Logger().Fatal("Could not instantiate HHG route planner", zap.Error(err))
	}

	// Create a secondary planner specifically for DTOD.
	dtodRoutePlanner, err := route.InitDTODRoutePlanner(appCtx, v, tlsConfig)
	if err != nil {
		appCtx.Logger().Fatal("Could not instantiate dtod route planner", zap.Error(err))
	}

	// Set the GexSender() and GexSender fields
	gexURL := v.GetString(cli.GEXURLFlag)
	var gexSender services.GexSender
	if len(gexURL) != 0 {
		gexSender = invoice.NewGexSenderHTTP(
			gexURL,
			true,
			tlsConfig,
			v.GetString(cli.GEXBasicAuthUsernameFlag),
			v.GetString(cli.GEXBasicAuthPasswordFlag),
		)
	} else {
		// this spins up a local test server
		fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		gexSender = invoice.NewGexSenderHTTP(
			fakeServer.URL,
			false,
			&tls.Config{MinVersion: tls.VersionTLS12},
			"",
			"",
		)
	}

	// Set the ICNSequencer in the handler: if we are in dev/test mode and sending to a real
	// GEX URL, then we should use a random ICN number within a defined range to avoid duplicate
	// test ICNs in Syncada.
	var icnSequencer sequence.Sequencer
	if isDevOrTest && len(gexURL) > 0 {
		// ICNs are 9-digit numbers; reserve the ones in an upper range for development/testing.
		icnSequencer, err = sequence.NewRandomSequencer(ediinvoice.ICNRandomMin, ediinvoice.ICNRandomMax)
		if err != nil {
			appCtx.Logger().Fatal("Could not create random sequencer for ICN", zap.Error(err))
		}
	} else {
		icnSequencer = sequence.NewDatabaseSequencer(ediinvoice.ICNSequenceName)
	}

	iwsPersonLookup, err := iws.InitRBSPersonLookup(appCtx, v)
	if err != nil {
		appCtx.Logger().Fatal("Could not instantiate IWS RBS", zap.Error(err))
	}

	storageBackend := v.GetString(cli.StorageBackendFlag)
	if storageBackend == "local" {
		routingConfig.LocalStorageRoot = v.GetString(cli.LocalStorageRootFlag)
		//Add a file handler to provide access to files uploaded in development
		routingConfig.LocalStorageWebRoot = v.GetString(cli.LocalStorageWebRootFlag)
	}

	sessionIdleTimeout := time.Duration(v.GetInt(cli.SessionIdleTimeoutInMinutesFlag)) * time.Minute
	sessionLifetime := time.Duration(v.GetInt(cli.SessionLifetimeInHoursFlag)) * time.Hour

	useSecureCookie := !isDevOrTest
	sessionManagers := auth.SetupSessionManagers(redisPool, useSecureCookie,
		sessionIdleTimeout, sessionLifetime)

	routingConfig.HandlerConfig = handlers.NewHandlerConfig(
		appCtx.DB(),
		appCtx.Logger(),
		clientAuthSecretKey,
		hhgRoutePlanner,
		dtodRoutePlanner,
		fileStorer,
		notificationSender,
		iwsPersonLookup,
		sendProductionInvoice,
		gexSender,
		icnSequencer,
		useSecureCookie,
		appNames,
		sessionManagers,
		[]handlers.FeatureFlag{},
	)

	initializeRouteOptions(v, routingConfig)

	return routingConfig
}

func serveFunction(cmd *cobra.Command, args []string) error {

	// variables that are initialized in this function and needed
	// during cleanup after the function exits
	var logger *zap.Logger
	var loggerSync func()
	var dbConnection *pop.Connection
	dbClose := &sync.Once{}
	var redisPool *redis.Pool
	redisClose := &sync.Once{}

	// cleanup that runs when this function ends ensuring we close
	// the database connection, the redis connection, and flush the
	// logger if needed
	defer func() {
		if logger != nil {
			if r := recover(); r != nil {
				logger.Error("server recovered from panic", zap.Any("recover", r))
			}
			if dbConnection != nil {
				dbClose.Do(func() {
					logger.Info("closing database connections")
					if err := dbConnection.Close(); err != nil {
						logger.Error("error closing database connections", zap.Error(err))
					}
				})
			}
			if redisPool != nil {
				redisClose.Do(func() {
					logger.Info("closing redis connections")
					if err := redisPool.Close(); err != nil {
						logger.Error("error closing redis connections", zap.Error(err))
					}
				})
			}

			loggerSync()
		}
	}()

	// Prepare to parse command line options / environment variables
	// using the viper library
	v, err := initializeViper(cmd, args)
	if err != nil {
		return err
	}

	// set up the logger and a function to flush the logger as needed
	logger, loggerSync = initializeLogger(v)
	logger.Info("webserver starting up")

	// ensure that the provided configuration options make sense
	// before starting the server
	err = checkServeConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Telemetry is used for reporting stats. Ensure the telemetry
	// config makes sense before starting
	telemetryConfig, err := cli.CheckTelemetry(v)
	if err != nil {
		logger.Fatal("invalid trace config", zap.Error(err))
	}

	// initialize the telemetry system and ensure it is shut down when
	// the server finishes
	telemetryShutdownFn := telemetry.Init(logger, telemetryConfig)
	defer telemetryShutdownFn()

	dbEnv := v.GetString(cli.DbEnvFlag)
	isDevOrTest := dbEnv == "development" || dbEnv == "test"
	if isDevOrTest {
		logger.Info(fmt.Sprintf("Starting in %s mode, which enables additional features", dbEnv))
	}

	// set up AWS (as needed)
	session := initializeAwsSession(v, logger)

	// connect to the db
	dbConnection = initializeDB(v, logger, session)

	// set up appcontext
	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	// now that we have the appcontext, register telemetry observers
	err = telemetry.RegisterDBStatsObserver(appCtx, telemetryConfig)
	if err != nil {
		logger.Fatal("Failed to register db stats observer", zap.Error(err))
	}
	err = telemetry.RegisterRuntimeObserver(appCtx, telemetryConfig)
	if err != nil {
		logger.Fatal("Failed to register runtime observer", zap.Error(err))
	}
	err = telemetry.RegisterMilmoveDataObserver(appCtx, telemetryConfig)
	if err != nil {
		logger.Fatal("Failed to register runtime observer", zap.Error(err))
	}

	// Create a connection to Redis
	redisPool, errRedisConnection := cli.InitRedis(appCtx, v)
	if errRedisConnection != nil {
		logger.Fatal("Invalid Redis Configuration", zap.Error(errRedisConnection))
	}

	// set up the tls configuration
	tlsConfig := initializeTLSConfig(appCtx, v)

	// build the routing configuration
	routingConfig := buildRoutingConfig(appCtx, v, redisPool, session,
		isDevOrTest, tlsConfig)

	listenInterface := v.GetString(cli.InterfaceFlag)

	// start each server:
	//
	// * healthServer that is only for health checks
	// * noTLSServer that is not listening on TLS. This server is
	//   generally only run in local development environments
	// * tlsServer that does listen using TLS. This server is run in
	//   production and generally handles all non prime API requests
	// * mutualTLSServer that requires mutual TLS authentication. This
	//   server handles prime API requests
	//
	// However, you will note that for historical reasons, each server
	// (other than the healthServer) has the entire routing setup for
	// all servers. It's thus the responsibility of routing config and
	// middleware to prevent the server from responding to the wrong
	// requests. The ideal would be for each server to have separate
	// routing config to limit the options.

	// see cmd/milmove/health.go for more rationale about why a
	// separate thread for a health listener was chosen
	healthEnabled := v.GetBool(cli.HealthListenerFlag)
	var healthServer *server.NamedServer
	if healthEnabled {
		serverName := "health"
		healthPort := v.GetInt(cli.HealthPortFlag)
		healthSite, err := routing.InitHealthRouting(serverName, appCtx, redisPool,
			routingConfig, telemetryConfig)
		if err != nil {
			return err
		}

		healthServer, err = server.CreateNamedServer(&server.CreateNamedServerInput{
			Name:        "health",
			Host:        "127.0.0.1", // health server is always localhost only
			Port:        healthPort,
			Logger:      logger,
			HTTPHandler: healthSite,
		})
		if err != nil {
			logger.Fatal("error creating health server", zap.Error(err))
		}
		go startListener(healthServer, logger, false)
	}

	noTLSEnabled := v.GetBool(cli.NoTLSListenerFlag)
	var noTLSServer *server.NamedServer
	if noTLSEnabled {
		serverName := "no-tls"
		noTLSPort := v.GetInt(cli.NoTLSPortFlag)
		// initialize the router
		site, err := routing.InitRouting(serverName, appCtx, redisPool,
			routingConfig, telemetryConfig)
		if err != nil {
			return err
		}

		noTLSServer, err = server.CreateNamedServer(&server.CreateNamedServerInput{
			Name:        serverName,
			Host:        listenInterface,
			Port:        noTLSPort,
			Logger:      logger,
			HTTPHandler: site,
		})
		if err != nil {
			logger.Fatal("error creating no-tls server", zap.Error(err))
		}
		go startListener(noTLSServer, logger, false)
	}

	tlsEnabled := v.GetBool(cli.TLSListenerFlag)
	var tlsServer *server.NamedServer
	if tlsEnabled {
		serverName := "tls"
		tlsPort := v.GetInt(cli.TLSPortFlag)
		// initialize the router
		site, err := routing.InitRouting(serverName, appCtx, redisPool,
			routingConfig, telemetryConfig)
		if err != nil {
			return err
		}
		tlsServer, err = server.CreateNamedServer(&server.CreateNamedServerInput{
			Name:         serverName,
			Host:         listenInterface,
			Port:         tlsPort,
			Logger:       logger,
			HTTPHandler:  site,
			ClientAuth:   tls.NoClientCert,
			Certificates: tlsConfig.Certificates,
		})
		if err != nil {
			logger.Fatal("error creating tls server", zap.Error(err))
		}
		go startListener(tlsServer, logger, true)
	}

	mutualTLSEnabled := v.GetBool(cli.MutualTLSListenerFlag)
	var mutualTLSServer *server.NamedServer
	if mutualTLSEnabled {
		serverName := "mutual-tls"
		mtlsPort := v.GetInt(cli.MutualTLSPortFlag)
		// initialize the router
		site, err := routing.InitRouting(serverName, appCtx, redisPool,
			routingConfig, telemetryConfig)
		if err != nil {
			return err
		}

		mutualTLSServer, err = server.CreateNamedServer(&server.CreateNamedServerInput{
			Name:         serverName,
			Host:         listenInterface,
			Port:         mtlsPort,
			Logger:       logger,
			HTTPHandler:  site,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: tlsConfig.Certificates,
			ClientCAs:    tlsConfig.RootCAs,
		})
		if err != nil {
			logger.Fatal("error creating mutual-tls server", zap.Error(err))
		}
		go startListener(mutualTLSServer, logger, true)
	}

	// make sure we flush any pending startup messages
	loggerSync()

	// Create a buffered channel that accepts 1 signal at a time.
	quit := make(chan os.Signal, 1)

	// Only send the SIGINT and SIGTERM signals to the quit channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait until the quit channel receieves a signal
	sig := <-quit

	logger.Info("received signal for graceful shutdown of server", zap.Any("signal", sig))

	loggerSync()

	gracefulShutdownTimeout := v.GetDuration(cli.GracefulShutdownTimeoutFlag)

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	logger.Info("Waiting for listeners to be shutdown", zap.Duration("timeout", gracefulShutdownTimeout))

	loggerSync()

	wg := &sync.WaitGroup{}
	var shutdownErrors sync.Map

	if noTLSEnabled {
		wg.Add(1)
		go func() {
			shutdownErrors.Store(noTLSServer, noTLSServer.Shutdown(ctx))
			wg.Done()
		}()
	}

	if tlsEnabled {
		wg.Add(1)
		go func() {
			shutdownErrors.Store(tlsServer, tlsServer.Shutdown(ctx))
			wg.Done()
		}()
	}

	if mutualTLSEnabled {
		wg.Add(1)
		go func() {
			shutdownErrors.Store(mutualTLSServer, mutualTLSServer.Shutdown(ctx))
			wg.Done()
		}()
	}

	if healthEnabled {
		wg.Add(1)
		go func() {
			shutdownErrors.Store(healthServer, healthServer.Shutdown(ctx))
			wg.Done()
		}()
	}

	wg.Wait()
	logger.Info("All listeners are shutdown")
	loggerSync()

	var dbCloseErr error
	dbClose.Do(func() {
		logger.Info("closing database connections")
		dbCloseErr = dbConnection.Close()
	})

	var redisCloseErr error
	redisClose.Do(func() {
		logger.Info("closing redis connections")
		redisCloseErr = redisPool.Close()
	})

	shutdownError := false
	shutdownErrors.Range(func(key, value interface{}) bool {
		if srv, ok := key.(*server.NamedServer); ok {
			if err, ok := value.(error); ok {
				logger.Error("shutdown error", zap.String("name", srv.Name), zap.String("addr", srv.Addr), zap.Int("port", srv.Port()), zap.Error(err))
				shutdownError = true
			} else {
				logger.Info("shutdown server", zap.String("name", srv.Name), zap.String("addr", srv.Addr), zap.Int("port", srv.Port()))
			}
		}
		return true
	})

	if dbCloseErr != nil {
		logger.Error("error closing database connections", zap.Error(dbCloseErr))
	}

	if redisCloseErr != nil {
		logger.Error("error closing redis connections", zap.Error(redisCloseErr))
	}

	loggerSync()

	if shutdownError {
		os.Exit(1)
	}

	return nil
}
