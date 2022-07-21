package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	"github.com/alexedwards/scs/redisstore"
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
	"github.com/trussworks/otelhttp"
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

	if err := cli.CheckSession(v); err != nil {
		return err
	}

	return nil
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
		resp, httpGetErr := http.Get("http://169.254.170.2/v2/metadata")
		if httpGetErr != nil {
			logger.Error(errors.Wrap(httpGetErr, "could not fetch task metadata").Error())
		} else {
			body, readAllErr := ioutil.ReadAll(resp.Body)
			if readAllErr != nil {
				logger.Error(errors.Wrap(readAllErr, "could not read task metadata").Error())
			} else {
				taskMetadata := &ecs.TaskMetadata{}
				unmarshallErr := json.Unmarshal(body, taskMetadata)
				if unmarshallErr != nil {
					logger.Error(errors.Wrap(unmarshallErr, "could not parse task metadata").Error())
				} else {
					logger = logger.With(
						zap.String("ecs_cluster", taskMetadata.Cluster),
						zap.String("ecs_task_def_family", taskMetadata.Family),
						zap.String("ecs_task_def_revision", taskMetadata.Revision),
					)
				}
			}
			err = resp.Body.Close()
			if err != nil {
				logger.Error(errors.Wrap(err, "could not close task metadata response").Error())
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
	dbConnection, errDbConnection := cli.InitDatabase(v, dbCreds, logger)
	if errDbConnection != nil {
		if dbConnection == nil {
			// No connection object means that the configuraton failed to validate and we should kill server startup
			logger.Fatal("Invalid DB Configuration", zap.Error(errDbConnection))
		} else {
			// A valid connection object that still has an error indicates that the DB is not up but we
			// can proceed (this avoids a failure loop when deploying containers).
			logger.Warn("DB is not ready for connections", zap.Error(errDbConnection))
		}
	}

	return dbConnection
}

func initializeTLSConfig(appCtx appcontext.AppContext, v *viper.Viper) *tls.Config {
	certificates, rootCAs, err := certs.InitDoDCertificates(v, appCtx.Logger())
	if certificates == nil || rootCAs == nil || err != nil {
		appCtx.Logger().Fatal("Failed to initialize DOD certificates", zap.Error(err))
	}
	appCtx.Logger().Debug("Server DOD Key Pair Loaded")
	appCtx.Logger().Debug("Trusted Certificate Authorities", zap.Any("subjects", rootCAs.Subjects()))

	useDevlocalAuthCA := stringSliceContains([]string{cli.EnvironmentTest, cli.EnvironmentDevelopment, cli.EnvironmentReview, cli.EnvironmentLoadtest}, v.GetString(cli.EnvironmentFlag))
	if useDevlocalAuthCA {
		appCtx.Logger().Info("Adding devlocal CA to root CAs")
		devlocalCAPath := v.GetString(cli.DevlocalCAFlag)
		devlocalCa, readFileErr := ioutil.ReadFile(filepath.Clean(devlocalCAPath))
		if readFileErr != nil {
			appCtx.Logger().Error(fmt.Sprintf("Unable to read devlocal CA from path %s", devlocalCAPath), zap.Error(readFileErr))
		} else {
			rootCAs.AppendCertsFromPEM(devlocalCa)
		}
	}

	return &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}
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
	routingConfig.CSRFAuthKey = v.GetString(cli.CSRFAuthKeyFlag)

	routingConfig.GitBranch = gitBranch
	routingConfig.GitCommit = gitCommit
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

	redisEnabled := v.GetBool(cli.RedisEnabledFlag)
	var sessionStore *redisstore.RedisStore
	if redisEnabled {
		sessionStore = redisstore.New(redisPool)
	}
	sessionIdleTimeout := time.Duration(v.GetInt(cli.SessionIdleTimeoutInMinutesFlag)) * time.Minute
	sessionLifetime := time.Duration(v.GetInt(cli.SessionLifetimeInHoursFlag)) * time.Hour

	useSecureCookie := !isDevOrTest
	sessionManagers := auth.SetupSessionManagers(redisEnabled,
		sessionStore, useSecureCookie,
		sessionIdleTimeout, sessionLifetime)
	routingConfig.AuthContext = authentication.NewAuthContext(appCtx.Logger(), loginGovProvider, loginGovCallbackProtocol, loginGovCallbackPort, sessionManagers)

	// Email
	notificationSender, err := notifications.InitEmail(v, awsSession, appCtx.Logger())
	if err != nil {
		appCtx.Logger().Fatal("notification sender sending not enabled", zap.Error(err))
	}

	routingConfig.BuildRoot = v.GetString(cli.BuildRootFlag)
	sendProductionInvoice := v.GetBool(cli.GEXSendProdInvoiceFlag)

	// Storage
	fileStorer := storage.InitStorage(v, awsSession, appCtx.Logger())

	// Get route planner for handlers to calculate transit distances
	// routePlanner := route.NewBingPlanner(logger, bingMapsEndpoint, bingMapsKey)
	routePlanner := route.InitRoutePlanner(v)

	// Create a secondary planner specifically for HHG.
	hhgRoutePlanner, err := route.InitHHGRoutePlanner(v, tlsConfig)
	if err != nil {
		appCtx.Logger().Fatal("Could not instantiate HHG route planner", zap.Error(err))
	}

	// Create a secondary planner specifically for DTOD.
	dtodRoutePlanner, err := route.InitDtodRoutePlanner(v, tlsConfig)
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

	routingConfig.HandlerConfig = handlers.NewHandlerConfig(
		appCtx.DB(),
		appCtx.Logger(),
		clientAuthSecretKey,
		routePlanner,
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
		[]handlers.FeatureFlag{},
		sessionManagers,
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
	telemetry.RegisterDBStatsObserver(appCtx, telemetryConfig)
	telemetry.RegisterRuntimeObserver(appCtx, telemetryConfig)

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

	// initialize the router
	site, err := routing.InitRouting(appCtx, redisPool, routingConfig, telemetryConfig)
	if err != nil {
		return err
	}

	// set up telemetry options for the server
	otelHTTPOptions := []otelhttp.Option{}
	if telemetryConfig.ReadEvents {
		otelHTTPOptions = append(otelHTTPOptions, otelhttp.WithMessageEvents(otelhttp.ReadEvents))
	}
	if telemetryConfig.WriteEvents {
		otelHTTPOptions = append(otelHTTPOptions, otelhttp.WithMessageEvents(otelhttp.WriteEvents))
	}
	listenInterface := v.GetString(cli.InterfaceFlag)

	// start each server:
	//
	// * noTLSServer that is not listening on TLS. This server is
	//   generally only run in local development environments
	// * tlsServer that does listen using TLS. This server is run in
	//   production and generally handles all non prime API requests
	// * mutualTLSServer that requires mutual TLS authentication. This
	//   server handles prime API requests
	//
	// However, you will note that for historical reasons, each server
	// has the entire routing setup for all servers. It's thus the
	// responsibility of routing config and middleware to prevent the
	// server from responding to the wrong requests. The ideal would
	// be for each server to have separate routing config to limit the
	// options.

	noTLSEnabled := v.GetBool(cli.NoTLSListenerFlag)
	var noTLSServer *server.NamedServer
	if noTLSEnabled {
		noTLSServer, err = server.CreateNamedServer(&server.CreateNamedServerInput{
			Name:        "no-tls",
			Host:        listenInterface,
			Port:        v.GetInt(cli.NoTLSPortFlag),
			Logger:      logger,
			HTTPHandler: otelhttp.NewHandler(site, "server-no-tls", otelHTTPOptions...),
		})
		if err != nil {
			logger.Fatal("error creating no-tls server", zap.Error(err))
		}
		go startListener(noTLSServer, logger, false)
	}

	tlsEnabled := v.GetBool(cli.TLSListenerFlag)
	var tlsServer *server.NamedServer
	if tlsEnabled {
		tlsServer, err = server.CreateNamedServer(&server.CreateNamedServerInput{
			Name:         "tls",
			Host:         listenInterface,
			Port:         v.GetInt(cli.TLSPortFlag),
			Logger:       logger,
			HTTPHandler:  otelhttp.NewHandler(site, "server-tls", otelHTTPOptions...),
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
		mutualTLSServer, err = server.CreateNamedServer(&server.CreateNamedServerInput{
			Name:         "mutual-tls",
			Host:         listenInterface,
			Port:         v.GetInt(cli.MutualTLSPortFlag),
			Logger:       logger,
			HTTPHandler:  otelhttp.NewHandler(site, "server-mtls", otelHTTPOptions...),
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
