package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/pprof"
	"os"
	"os/signal"
	"path"
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
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/trussworks/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/ecs"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/adminapi"
	"github.com/transcom/mymove/pkg/handlers/dpsapi"
	"github.com/transcom/mymove/pkg/handlers/ghcapi"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"github.com/transcom/mymove/pkg/handlers/ordersapi"
	"github.com/transcom/mymove/pkg/handlers/primeapi"
	"github.com/transcom/mymove/pkg/handlers/supportapi"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/middleware"
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

	// SDDC + DPS Auth config
	cli.InitDPSFlags(flag)

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

	if err := cli.CheckDPS(v); err != nil {
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

// fileHandler serves up a single file
func fileHandler(entrypoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
}

// indexHandler returns a handler that will serve the resulting content
func indexHandler(buildDir string, globalLogger *zap.Logger) http.HandlerFunc {

	indexPath := path.Join(buildDir, "index.html")
	indexHTML, err := ioutil.ReadFile(filepath.Clean(indexPath))
	if err != nil {
		globalLogger.Fatal("could not read index.html template: run make client_build", zap.Error(err))
	}

	stat, err := os.Stat(indexPath)
	if err != nil {
		globalLogger.Fatal("could not stat index.html template", zap.Error(err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "index.html", stat.ModTime(), bytes.NewReader(indexHTML))
	}
}

func redisHealthCheck(pool *redis.Pool, logger *zap.Logger, data map[string]interface{}) map[string]interface{} {
	conn := pool.Get()

	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			logger.Error("Failed to close redis connection", zap.Error(closeErr))
		}
	}()

	pong, err := redis.String(conn.Do("PING"))
	if err != nil {
		logger.Error("Failed to ping Redis during health check", zap.Error(err))
	}
	logger.Info("Health check Redis ping", zap.String("ping_response", pong))

	data["redis"] = err == nil

	return data
}

func serveFunction(cmd *cobra.Command, args []string) error {

	var logger *zap.Logger
	var loggerSync func()
	var dbConnection *pop.Connection
	dbClose := &sync.Once{}
	var redisPool *redis.Pool
	redisClose := &sync.Once{}

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

	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "Could not parse flags")
	}

	flag := cmd.Flags()

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		return errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger, loggerSync, err = logging.Config(
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

	logger.Info("webserver starting up")

	err = checkServeConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}
	telemetryConfig, err := cli.CheckTelemetry(v)
	if err != nil {
		logger.Fatal("invalid trace config", zap.Error(err))
	}

	telemetryShutdownFn := telemetry.Init(logger, telemetryConfig)
	defer telemetryShutdownFn()

	dbEnv := v.GetString(cli.DbEnvFlag)

	isDevOrTest := dbEnv == "development" || dbEnv == "test"
	if isDevOrTest {
		logger.Info(fmt.Sprintf("Starting in %s mode, which enables additional features", dbEnv))
	}

	clientAuthSecretKey := v.GetString(cli.ClientAuthSecretKeyFlag)
	loginGovCallbackProtocol := v.GetString(cli.LoginGovCallbackProtocolFlag)
	loginGovCallbackPort := v.GetInt(cli.LoginGovCallbackPortFlag)
	loginGovSecretKey := v.GetString(cli.LoginGovSecretKeyFlag)
	loginGovHostname := v.GetString(cli.LoginGovHostnameFlag)

	// Assert that our secret keys can be parsed into actual private keys
	// TODO: Store the parsed key in handlers/AppContext instead of parsing every time
	if _, parseRSAPrivateKeyFromPEMErr := jwt.ParseRSAPrivateKeyFromPEM([]byte(loginGovSecretKey)); parseRSAPrivateKeyFromPEMErr != nil {
		logger.Fatal("Login.gov private key", zap.Error(parseRSAPrivateKeyFromPEMErr))
	}
	if _, parseRSAPrivateKeyFromPEMErr := jwt.ParseRSAPrivateKeyFromPEM([]byte(clientAuthSecretKey)); parseRSAPrivateKeyFromPEMErr != nil {
		logger.Fatal("Client auth private key", zap.Error(parseRSAPrivateKeyFromPEMErr))
	}
	if len(loginGovHostname) == 0 {
		logger.Fatal("Must provide the Login.gov hostname parameter, exiting")
	}

	if v.GetBool(cli.DbDebugFlag) {
		pop.Debug = true
	}

	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) || (v.GetString(cli.EmailBackendFlag) == "ses") || (v.GetString(cli.StorageBackendFlag) == "s3") {
		c := &aws.Config{
			Region: aws.String(v.GetString(cli.AWSRegionFlag)),
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
		}
		session = s
	}

	var dbCreds *credentials.Credentials
	if v.GetBool(cli.DbIamFlag) {
		if session != nil {
			// We want to get the credentials from the logged in AWS session rather than create directly,
			// because the session conflates the environment, shared, and container metdata config
			// within NewSession.  With stscreds, we use the Secure Token Service,
			// to assume the given role (that has rds db connect permissions).
			dbIamRole := v.GetString(cli.DbIamRoleFlag)
			logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))
			dbCreds = stscreds.NewCredentials(session, dbIamRole)
			stsService := sts.New(session)
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

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)

	telemetry.RegisterDBStatsObserver(appCtx, telemetryConfig)
	telemetry.RegisterRuntimeObserver(appCtx, telemetryConfig)

	// Create a connection to Redis
	redisPool, errRedisConnection := cli.InitRedis(appCtx, v)
	if errRedisConnection != nil {
		logger.Fatal("Invalid Redis Configuration", zap.Error(errRedisConnection))
	}

	// Collect the servernames into a handy struct
	appnames := auth.ApplicationServername{
		MilServername:    v.GetString(cli.HTTPMyServerNameFlag),
		OfficeServername: v.GetString(cli.HTTPOfficeServerNameFlag),
		AdminServername:  v.GetString(cli.HTTPAdminServerNameFlag),
		OrdersServername: v.GetString(cli.HTTPOrdersServerNameFlag),
		DpsServername:    v.GetString(cli.HTTPDPSServerNameFlag),
		SddcServername:   v.GetString(cli.HTTPSDDCServerNameFlag),
		PrimeServername:  v.GetString(cli.HTTPPrimeServerNameFlag),
	}

	// Register Login.gov authentication provider for My.(move.mil)
	loginGovProvider, err := authentication.InitAuth(v, logger, appnames)
	if err != nil {
		logger.Fatal("Registering login provider", zap.Error(err))
	}

	useSecureCookie := !isDevOrTest
	redisEnabled := v.GetBool(cli.RedisEnabledFlag)
	sessionStore := redisstore.New(redisPool)
	idleTimeout := time.Duration(v.GetInt(cli.SessionIdleTimeoutInMinutesFlag)) * time.Minute
	lifetime := time.Duration(v.GetInt(cli.SessionLifetimeInHoursFlag)) * time.Hour
	sessionManagers := auth.SetupSessionManagers(redisEnabled, sessionStore, useSecureCookie, idleTimeout, lifetime)
	milSession := sessionManagers[0]
	adminSession := sessionManagers[1]
	officeSession := sessionManagers[2]

	// Session management and authentication middleware
	sessionCookieMiddleware := auth.SessionCookieMiddleware(logger, appnames, sessionManagers)
	maskedCSRFMiddleware := auth.MaskedCSRFMiddleware(logger, useSecureCookie)
	userAuthMiddleware := authentication.UserAuthMiddleware(logger)
	isLoggedInMiddleware := authentication.IsLoggedInMiddleware(logger)
	clientCertMiddleware := authentication.ClientCertMiddleware(appCtx)

	handlerConfig := handlers.NewHandlerConfig(dbConnection, logger)
	handlerConfig.SetSessionManagers(sessionManagers)
	handlerConfig.SetCookieSecret(clientAuthSecretKey)
	handlerConfig.SetUseSecureCookie(useSecureCookie)
	handlerConfig.SetAppNames(appnames)

	// Email
	notificationSender, notificationSenderErr := notifications.InitEmail(v, session, logger)
	if notificationSenderErr != nil {
		logger.Fatal("notification sender sending not enabled", zap.Error(notificationSenderErr))
	}
	handlerConfig.SetNotificationSender(notificationSender)

	build := v.GetString(cli.BuildRootFlag)

	// Serves files out of build folder
	clientHandler := spaHandler{
		staticPath: build,
		indexPath:  "index.html",
	}
	// Set SendProductionInvoice for ediinvoice
	handlerConfig.SetSendProductionInvoice(v.GetBool(cli.GEXSendProdInvoiceFlag))

	// Storage
	storer := storage.InitStorage(v, session, logger)
	handlerConfig.SetFileStorer(storer)

	certificates, rootCAs, err := certs.InitDoDCertificates(v, logger)
	if certificates == nil || rootCAs == nil || err != nil {
		logger.Fatal("Failed to initialize DOD certificates", zap.Error(err))
	}

	logger.Debug("Server DOD Key Pair Loaded")
	// RA Summary: staticcheck - SA1019 - Using a deprecated function, variable, constant or field
	// RA: Linter is flagging: rootCAs.Subjects is deprecated: if s was returned by SystemCertPool, Subjects will not include the system roots.
	// RA: Why code valuable: It allows us to log the root CA subjects that are being trusted.
	// RA: Mitigation: The deprecation notes this is a problem when reading SystemCertPool, but we do not use this here and are building our own cert pool instead.
	// RA Developer Status: Mitigated
	// RA Validator Status: Mitigated
	// RA Validator: leodis.f.scott.civ@mail.mil
	// RA Modified Severity: CAT III
	// nolint:staticcheck
	logger.Debug("Trusted Certificate Authorities", zap.Any("subjects", rootCAs.Subjects()))

	// Get route planner for handlers to calculate transit distances
	// routePlanner := route.NewBingPlanner(logger, bingMapsEndpoint, bingMapsKey)
	routePlanner := route.InitRoutePlanner(v)
	handlerConfig.SetPlanner(routePlanner)

	// Create a secondary planner specifically for HHG.
	routeTLSConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}
	hhgRoutePlanner, initRouteErr := route.InitHHGRoutePlanner(v, routeTLSConfig)
	if initRouteErr != nil {
		logger.Fatal("Could not instantiate HHG route planner", zap.Error(initRouteErr))
	}
	handlerConfig.SetHHGPlanner(hhgRoutePlanner)

	// Create a secondary planner specifically for DTOD.
	dtodRoutePlanner, initRouteErr := route.InitDtodRoutePlanner(v, routeTLSConfig)
	if initRouteErr != nil {
		logger.Fatal("Could not instantiate dtod route planner", zap.Error(initRouteErr))
	}
	handlerConfig.SetDtodPlanner(dtodRoutePlanner)

	// Set the GexSender() and GexSender fields
	gexTLSConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}
	var gexRequester services.GexSender
	gexURL := v.GetString(cli.GEXURLFlag)
	if len(gexURL) == 0 {
		// this spins up a local test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		gexRequester = invoice.NewGexSenderHTTP(
			server.URL,
			false,
			&tls.Config{MinVersion: tls.VersionTLS12},
			"",
			"",
		)
	} else {
		gexRequester = invoice.NewGexSenderHTTP(
			gexURL,
			true,
			gexTLSConfig,
			v.GetString(cli.GEXBasicAuthUsernameFlag),
			v.GetString(cli.GEXBasicAuthPasswordFlag),
		)
	}
	handlerConfig.SetGexSender(gexRequester)

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
		icnSequencer = sequence.NewDatabaseSequencer(ediinvoice.ICNSequenceName)
	}
	handlerConfig.SetICNSequencer(icnSequencer)

	rbs, err := iws.InitRBSPersonLookup(appCtx, v)
	if err != nil {
		logger.Fatal("Could not instantiate IWS RBS", zap.Error(err))
	}
	handlerConfig.SetIWSPersonLookup(rbs)

	dpsAuthSecretKey := v.GetString(cli.DPSAuthSecretKeyFlag)
	dpsCookieDomain := v.GetString(cli.DPSCookieDomainFlag)
	dpsCookieSecret := []byte(v.GetString(cli.DPSAuthCookieSecretKeyFlag))
	dpsCookieExpires := v.GetInt(cli.DPSCookieExpiresInMinutesFlag)

	dpsAuthParams := dpsauth.InitDPSAuthParams(v, appnames)
	handlerConfig.SetDPSAuthParams(dpsAuthParams)

	// site is the base
	site := mux.NewRouter()
	storageBackend := v.GetString(cli.StorageBackendFlag)
	if storageBackend == "local" {
		localStorageRoot := v.GetString(cli.LocalStorageRootFlag)
		localStorageWebRoot := v.GetString(cli.LocalStorageWebRootFlag)
		//Add a file handler to provide access to files uploaded in development
		fs := storage.NewFilesystemHandler(localStorageRoot)

		site.HandleFunc(path.Join("/", localStorageWebRoot), fs)
	}
	// Add middleware: they are evaluated in the reverse order in which they
	// are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	site.Use(middleware.Trace(logger)) // injects trace id into the context
	site.Use(middleware.ContextLogger("milmove_trace_id", logger))
	site.Use(middleware.Recovery(logger))
	site.Use(middleware.SecurityHeaders(logger))

	if maxBodySize := v.GetInt64(cli.MaxBodySizeFlag); maxBodySize > 0 {
		site.Use(middleware.LimitBodySize(maxBodySize, logger))
	}

	// Stub health check
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
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
			logger.Info("Health check connecting to the DB")
			dbErr := dbConnection.RawQuery("SELECT 1;").Exec()
			if dbErr != nil {
				logger.Error("Failed database health check", zap.Error(dbErr))
			}
			data["database"] = dbErr == nil
			if redisEnabled {
				data = redisHealthCheck(redisPool, logger, data)
			}
		}
		newEncoderErr := json.NewEncoder(w).Encode(data)
		if newEncoderErr != nil {
			logger.Error("Failed encoding health check response", zap.Error(newEncoderErr))
		}

		// We are not using request middleware here so logging directly in the check
		var protocol string
		if r.TLS == nil {
			protocol = "http"
		} else {
			protocol = "https"
		}

		fields := []zap.Field{
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
		}

		// Append x- headers, e.g., x-forwarded-for.
		for name, values := range r.Header {
			if nameLowerCase := strings.ToLower(name); strings.HasPrefix(nameLowerCase, "x-") {
				if len(values) > 0 {
					fields = append(fields, zap.String(nameLowerCase, values[0]))
				}
			}
		}

		// Log the number of headers, which can be used for finding abnormal requests
		fields = append(fields, zap.Int("headers", len(r.Header)))

		logger.Info("Request health", fields...)

	}
	site.HandleFunc("/health", healthHandler).Methods("GET")

	staticMux := site.PathPrefix("/static/").Subrouter()
	staticMux.Use(middleware.ValidMethodsStatic(logger))
	staticMux.Use(middleware.RequestLogger(logger))
	if telemetryConfig.Enabled {
		staticMux.Use(otelmux.Middleware("static"))
	}
	staticMux.PathPrefix("/").Handler(clientHandler).Methods("GET", "HEAD")

	downloadMux := site.PathPrefix("/downloads/").Subrouter()
	downloadMux.Use(middleware.ValidMethodsStatic(logger))
	downloadMux.Use(middleware.RequestLogger(logger))
	if telemetryConfig.Enabled {
		downloadMux.Use(otelmux.Middleware("download"))
	}
	downloadMux.PathPrefix("/").Handler(clientHandler).Methods("GET", "HEAD")

	site.Handle("/favicon.ico", clientHandler)

	// Explicitly disable swagger.json route
	site.Handle("/swagger.json", http.NotFoundHandler()).Methods("GET")
	if v.GetBool(cli.ServeSwaggerUIFlag) {
		logger.Info("Swagger UI static file serving is enabled")
		site.PathPrefix("/swagger-ui/").Handler(clientHandler).Methods("GET")
	} else {
		site.PathPrefix("/swagger-ui/").Handler(http.NotFoundHandler()).Methods("GET")
	}

	if v.GetBool(cli.ServeOrdersFlag) {
		ordersMux := site.Host(appnames.OrdersServername).PathPrefix("/orders/v1/").Subrouter()
		ordersDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.OrdersServername)
		ordersMux.Use(ordersDetectionMiddleware)
		ordersMux.Use(middleware.NoCache(logger))
		ordersMux.Use(clientCertMiddleware)
		ordersMux.Use(middleware.RequestLogger(logger))
		ordersMux.HandleFunc("/swagger.yaml", fileHandler(v.GetString(cli.OrdersSwaggerFlag))).Methods("GET")
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Orders API Swagger UI serving is enabled")
			ordersMux.HandleFunc("/docs", fileHandler(path.Join(build, "swagger-ui", "orders.html"))).Methods("GET")
		} else {
			ordersMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}
		api := ordersapi.NewOrdersAPI(handlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		ordersMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}
	if v.GetBool(cli.ServeDPSFlag) {
		dpsMux := site.Host(appnames.DpsServername).PathPrefix("/dps/v0/").Subrouter()
		dpsDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.DpsServername)
		dpsMux.Use(dpsDetectionMiddleware)
		dpsMux.Use(middleware.NoCache(logger))
		dpsMux.Use(clientCertMiddleware)
		dpsMux.Use(middleware.RequestLogger(logger))
		dpsMux.HandleFunc("/swagger.yaml", fileHandler(v.GetString(cli.DPSSwaggerFlag))).Methods("GET")
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("DPS API Swagger UI serving is enabled")
			dpsMux.HandleFunc("/docs", fileHandler(path.Join(build, "swagger-ui", "dps.html"))).Methods("GET")
		} else {
			dpsMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}
		api := dpsapi.NewDPSAPI(handlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		dpsMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if v.GetBool(cli.ServeSDDCFlag) {
		sddcDPSMux := site.Host(appnames.SddcServername).PathPrefix("/dps_auth/").Subrouter()

		sddcDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.SddcServername)
		sddcDPSMux.Use(sddcDetectionMiddleware)
		sddcDPSMux.Use(middleware.NoCache(logger))
		sddcDPSMux.Use(middleware.RequestLogger(logger))
		if telemetryConfig.Enabled {
			sddcDPSMux.Use(otelmux.Middleware("sddc"))
		}
		sddcDPSMux.Handle("/set_cookie",
			dpsauth.NewSetCookieHandler(dpsAuthSecretKey,
				dpsCookieDomain,
				dpsCookieSecret,
				dpsCookieExpires))
	}

	if v.GetBool(cli.ServePrimeFlag) {
		primeMux := site.Host(appnames.PrimeServername).PathPrefix("/prime/v1/").Subrouter()

		primeDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.PrimeServername)
		primeMux.Use(primeDetectionMiddleware)
		if v.GetBool(cli.DevlocalAuthFlag) {
			devlocalClientCertMiddleware := authentication.DevlocalClientCertMiddleware(appCtx)
			primeMux.Use(devlocalClientCertMiddleware)
		} else {
			primeMux.Use(clientCertMiddleware)
		}
		primeMux.Use(authentication.PrimeAuthorizationMiddleware(logger))
		primeMux.Use(middleware.NoCache(logger))
		primeMux.Use(middleware.RequestLogger(logger))
		primeMux.HandleFunc("/swagger.yaml", fileHandler(v.GetString(cli.PrimeSwaggerFlag))).Methods("GET")
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Prime API Swagger UI serving is enabled")
			primeMux.HandleFunc("/docs", fileHandler(path.Join(build, "swagger-ui", "prime.html"))).Methods("GET")
		} else {
			primeMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}
		api := primeapi.NewPrimeAPI(handlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		primeMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if v.GetBool(cli.ServeSupportFlag) {
		supportMux := site.Host(appnames.PrimeServername).PathPrefix("/support/v1/").Subrouter()

		supportDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.PrimeServername)
		supportMux.Use(supportDetectionMiddleware)
		supportMux.Use(clientCertMiddleware)
		supportMux.Use(authentication.PrimeAuthorizationMiddleware(logger))
		supportMux.Use(middleware.NoCache(logger))
		supportMux.Use(middleware.RequestLogger(logger))
		supportMux.HandleFunc("/swagger.yaml", fileHandler(v.GetString(cli.SupportSwaggerFlag))).Methods("GET")
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Support API Swagger UI serving is enabled")
			supportMux.HandleFunc("/docs", fileHandler(path.Join(build, "swagger-ui", "support.html"))).Methods("GET")
		} else {
			supportMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}
		supportMux.PathPrefix("/").Handler(supportapi.NewSupportAPIHandler(handlerConfig))
	}

	// Handlers under mutual TLS need to go before this section that sets up middleware that shouldn't be enabled for mutual TLS (such as CSRF)
	root := mux.NewRouter()
	root.Use(sessionCookieMiddleware)
	root.Use(middleware.RequestLogger(logger))

	debug := root.PathPrefix("/debug/pprof/").Subrouter()
	debug.Use(userAuthMiddleware)
	if v.GetBool(cli.DebugPProfFlag) {
		logger.Info("Enabling pprof routes")
		debug.HandleFunc("/", pprof.Index).Methods("GET")
		debug.Handle("/allocs", pprof.Handler("allocs")).Methods("GET")
		debug.Handle("/block", pprof.Handler("block")).Methods("GET")
		debug.HandleFunc("/cmdline", pprof.Cmdline).Methods("GET")
		debug.Handle("/goroutine", pprof.Handler("goroutine")).Methods("GET")
		debug.Handle("/heap", pprof.Handler("heap")).Methods("GET")
		debug.Handle("/mutex", pprof.Handler("mutex")).Methods("GET")
		debug.HandleFunc("/profile", pprof.Profile).Methods("GET")
		debug.HandleFunc("/trace", pprof.Trace).Methods("GET")
		debug.Handle("/threadcreate", pprof.Handler("threadcreate")).Methods("GET")
		debug.HandleFunc("/symbol", pprof.Symbol).Methods("GET")
	} else {
		debug.HandleFunc("/", http.NotFound).Methods("GET")
	}

	// CSRF path is set specifically at the root to avoid duplicate tokens from different paths
	csrfAuthKey, err := hex.DecodeString(v.GetString(cli.CSRFAuthKeyFlag))
	if err != nil {
		logger.Fatal("Failed to decode csrf auth key", zap.Error(err))
	}
	logger.Info("Enabling CSRF protection")
	root.Use(csrf.Protect(csrfAuthKey, csrf.Secure(!isDevOrTest), csrf.Path("/"), csrf.CookieName(auth.GorillaCSRFToken)))
	root.Use(maskedCSRFMiddleware)

	site.Host(appnames.MilServername).PathPrefix("/").Handler(milSession.LoadAndSave(root))
	site.Host(appnames.AdminServername).PathPrefix("/").Handler(adminSession.LoadAndSave(root))
	site.Host(appnames.OfficeServername).PathPrefix("/").Handler(officeSession.LoadAndSave(root))

	if v.GetBool(cli.ServeAPIInternalFlag) {
		internalMux := root.PathPrefix("/internal/").Subrouter()
		internalMux.HandleFunc("/swagger.yaml", fileHandler(v.GetString(cli.InternalSwaggerFlag))).Methods("GET")
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Internal API Swagger UI serving is enabled")
			internalMux.HandleFunc("/docs", fileHandler(path.Join(build, "swagger-ui", "internal.html"))).Methods("GET")
		} else {
			internalMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}
		internalMux.Use(middleware.RequestLogger(logger))
		internalMux.HandleFunc("/users/is_logged_in", isLoggedInMiddleware).Methods("GET")
		// Mux for internal API that enforces auth
		internalAPIMux := internalMux.PathPrefix("/").Subrouter()
		internalAPIMux.Use(userAuthMiddleware)
		internalAPIMux.Use(middleware.NoCache(logger))
		api := internalapi.NewInternalAPI(handlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		internalAPIMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if v.GetBool(cli.ServeAdminFlag) {
		adminMux := root.PathPrefix("/admin/v1/").Subrouter()

		adminMux.HandleFunc("/swagger.yaml", fileHandler(v.GetString(cli.AdminSwaggerFlag))).Methods("GET")
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Admin API Swagger UI serving is enabled")
			adminMux.HandleFunc("/docs", fileHandler(path.Join(build, "swagger-ui", "admin.html"))).Methods("GET")
		} else {
			adminMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}

		// Mux for admin API that enforces auth
		adminAPIMux := adminMux.PathPrefix("/").Subrouter()
		adminAPIMux.Use(userAuthMiddleware)
		adminAPIMux.Use(authentication.AdminAuthMiddleware(logger))
		adminAPIMux.Use(middleware.NoCache(logger))
		api := adminapi.NewAdminAPI(handlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		adminAPIMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if v.GetBool(cli.ServePrimeSimulatorFlag) {
		// attach prime simulator API to root so cookies are handled
		primeSimulatorMux := root.Host(appnames.OfficeServername).PathPrefix("/prime/v1/").Subrouter()
		primeSimulatorMux.HandleFunc("/swagger.yaml", fileHandler(v.GetString(cli.PrimeSwaggerFlag))).Methods("GET")
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Prime Simulator API Swagger UI serving is enabled")
			primeSimulatorMux.HandleFunc("/docs", fileHandler(path.Join(build, "swagger-ui", "prime.html"))).Methods("GET")
		} else {
			primeSimulatorMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}

		// Mux for prime simulator API that enforces auth
		primeSimulatorAPIMux := primeSimulatorMux.PathPrefix("/").Subrouter()
		primeSimulatorAPIMux.Use(userAuthMiddleware)
		primeSimulatorAPIMux.Use(authentication.PrimeSimulatorAuthorizationMiddleware(logger))
		primeSimulatorAPIMux.Use(middleware.NoCache(logger))
		api := primeapi.NewPrimeAPI(handlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		primeSimulatorAPIMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if v.GetBool(cli.ServeGHCFlag) {
		ghcMux := root.PathPrefix("/ghc/v1/").Subrouter()
		ghcMux.HandleFunc("/swagger.yaml", fileHandler(v.GetString(cli.GHCSwaggerFlag))).Methods("GET")
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("GHC API Swagger UI serving is enabled")
			ghcMux.HandleFunc("/docs", fileHandler(path.Join(build, "swagger-ui", "ghc.html"))).Methods("GET")
		} else {
			ghcMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}

		// Mux for GHC API that enforces auth
		ghcAPIMux := ghcMux.PathPrefix("/").Subrouter()
		ghcAPIMux.Use(userAuthMiddleware)
		ghcAPIMux.Use(middleware.NoCache(logger))
		api := ghcapi.NewGhcAPIHandler(handlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		ghcAPIMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	authContext := authentication.NewAuthContext(logger, loginGovProvider, loginGovCallbackProtocol, loginGovCallbackPort, sessionManagers)
	authMux := root.PathPrefix("/auth/").Subrouter()
	authMux.Use(middleware.NoCache(logger))
	authMux.Use(otelmux.Middleware("auth"))
	authMux.Handle("/login-gov", authentication.NewRedirectHandler(authContext, handlerConfig, useSecureCookie)).Methods("GET")
	authMux.Handle("/login-gov/callback", authentication.NewCallbackHandler(authContext, handlerConfig, notificationSender)).Methods("GET")
	authMux.Handle("/logout", authentication.NewLogoutHandler(authContext, handlerConfig)).Methods("POST")

	if v.GetBool(cli.DevlocalAuthFlag) {
		logger.Info("Enabling devlocal auth")
		localAuthMux := root.PathPrefix("/devlocal-auth/").Subrouter()
		localAuthMux.Use(middleware.NoCache(logger))
		localAuthMux.Use(otelmux.Middleware("devlocal"))
		localAuthMux.Handle("/login", authentication.NewUserListHandler(authContext, handlerConfig)).Methods("GET")
		localAuthMux.Handle("/login", authentication.NewAssignUserHandler(authContext, handlerConfig, appnames)).Methods("POST")
		localAuthMux.Handle("/new", authentication.NewCreateAndLoginUserHandler(authContext, handlerConfig, appnames)).Methods("POST")
		localAuthMux.Handle("/create", authentication.NewCreateUserHandler(authContext, handlerConfig, appnames)).Methods("POST")

		if stringSliceContains([]string{cli.EnvironmentTest, cli.EnvironmentDevelopment, cli.EnvironmentReview, cli.EnvironmentLoadtest}, v.GetString(cli.EnvironmentFlag)) {
			logger.Info("Adding devlocal CA to root CAs")
			devlocalCAPath := v.GetString(cli.DevlocalCAFlag)
			devlocalCa, readFileErr := ioutil.ReadFile(filepath.Clean(devlocalCAPath))
			if readFileErr != nil {
				logger.Error(fmt.Sprintf("Unable to read devlocal CA from path %s", devlocalCAPath), zap.Error(readFileErr))
			} else {
				rootCAs.AppendCertsFromPEM(devlocalCa)
			}
		}
	}

	// Serve index.html to all requests that haven't matches a previous route,
	root.PathPrefix("/").Handler(indexHandler(build, logger)).Methods("GET", "HEAD")

	otelHTTPOptions := []otelhttp.Option{}
	if telemetryConfig.ReadEvents {
		otelHTTPOptions = append(otelHTTPOptions, otelhttp.WithMessageEvents(otelhttp.ReadEvents))
	}
	if telemetryConfig.WriteEvents {
		otelHTTPOptions = append(otelHTTPOptions, otelhttp.WithMessageEvents(otelhttp.WriteEvents))
	}
	listenInterface := v.GetString(cli.InterfaceFlag)

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
			Certificates: certificates,
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
			Certificates: certificates,
			ClientCAs:    rootCAs,
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
