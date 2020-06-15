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
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/pop/v5"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	goji "goji.io"
	"goji.io/pat"

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

	// Feature Flags
	cli.InitFeatureFlags(flag)

	// pprof flags
	cli.InitDebugFlags(flag)

	// Service Flags
	cli.InitServiceFlags(flag)

	// Redis Flags
	cli.InitRedisFlags(flag)

	// SessionFlags
	cli.InitSessionFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

func checkServeConfig(v *viper.Viper, logger logger) error {

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

	if err := cli.CheckFeatureFlag(v); err != nil {
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

func startListener(srv *server.NamedServer, logger logger, useTLS bool) {
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
func indexHandler(buildDir string, logger logger) http.HandlerFunc {

	indexPath := path.Join(buildDir, "index.html")
	// #nosec - indexPath does not come from user input
	indexHTML, err := ioutil.ReadFile(filepath.Clean(indexPath))
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

			loggerSyncErr := logger.Sync()
			if loggerSyncErr != nil {
				logger.Error("Failed to sync logger", zap.Error(loggerSyncErr))
			}
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

	logger, err = logging.Config(
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

	// Create a connection to Redis
	redisPool, errRedisConnection := cli.InitRedis(v, logger)
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
	clientCertMiddleware := authentication.ClientCertMiddleware(logger, dbConnection)

	handlerContext := handlers.NewHandlerContext(dbConnection, logger)
	handlerContext.SetSessionManagers(sessionManagers)
	handlerContext.SetCookieSecret(clientAuthSecretKey)
	handlerContext.SetUseSecureCookie(useSecureCookie)
	handlerContext.SetAppNames(appnames)

	// Email
	notificationSender, notificationSenderErr := notifications.InitEmail(v, session, logger)
	if notificationSenderErr != nil {
		logger.Fatal("notification sender sending not enabled", zap.Error(notificationSenderErr))
	}
	handlerContext.SetNotificationSender(notificationSender)

	build := v.GetString(cli.BuildRootFlag)

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(build))

	// Set SendProductionInvoice for ediinvoice
	handlerContext.SetSendProductionInvoice(v.GetBool(cli.GEXSendProdInvoiceFlag))

	// Storage
	storer := storage.InitStorage(v, session, logger)
	handlerContext.SetFileStorer(storer)

	certificates, rootCAs, err := certs.InitDoDCertificates(v, logger)
	if certificates == nil || rootCAs == nil || err != nil {
		logger.Fatal("Failed to initialize DOD certificates", zap.Error(err))
	}

	logger.Debug("Server DOD Key Pair Loaded")
	logger.Debug("Trusted Certificate Authorities", zap.Any("subjects", rootCAs.Subjects()))

	// Get route planner for handlers to calculate transit distances
	// routePlanner := route.NewBingPlanner(logger, bingMapsEndpoint, bingMapsKey)
	routePlanner := route.InitRoutePlanner(v, logger)
	handlerContext.SetPlanner(routePlanner)

	// Create a secondary planner specifically for GHC.
	routeTLSConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}
	ghcRoutePlanner, initRouteErr := route.InitGHCRoutePlanner(v, logger, dbConnection, routeTLSConfig)
	if initRouteErr != nil {
		logger.Fatal("Could not instantiate GHC route planner", zap.Error(initRouteErr))
	}
	handlerContext.SetGHCPlanner(ghcRoutePlanner)

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
	handlerContext.SetGexSender(gexRequester)

	// Set feature flags
	handlerContext.SetFeatureFlag(
		handlers.FeatureFlag{Name: cli.FeatureFlagAccessCode, Active: v.GetBool(cli.FeatureFlagAccessCode)},
	)

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

	rbs, err := iws.InitRBSPersonLookup(v, logger)
	if err != nil {
		logger.Fatal("Could not instantiate IWS RBS", zap.Error(err))
	}
	handlerContext.SetIWSPersonLookup(rbs)

	dpsAuthSecretKey := v.GetString(cli.DPSAuthSecretKeyFlag)
	dpsCookieDomain := v.GetString(cli.DPSCookieDomainFlag)
	dpsCookieSecret := []byte(v.GetString(cli.DPSAuthCookieSecretKeyFlag))
	dpsCookieExpires := v.GetInt(cli.DPSCookieExpiresInMinutesFlag)

	dpsAuthParams := dpsauth.InitDPSAuthParams(v, appnames)
	handlerContext.SetDPSAuthParams(dpsAuthParams)

	// bare is the base muxer. Not intended to have any middleware attached.
	bare := goji.NewMux()
	storageBackend := v.GetString(cli.StorageBackendFlag)
	if storageBackend == "local" {
		localStorageRoot := v.GetString(cli.LocalStorageRootFlag)
		localStorageWebRoot := v.GetString(cli.LocalStorageWebRootFlag)
		//Add a file handler to provide access to files uploaded in development
		fs := storage.NewFilesystemHandler(localStorageRoot)
		bare.Handle(pat.Get(path.Join("/", localStorageWebRoot, "/*")), fs)
	}
	// Base routes
	site := goji.SubMux()
	bare.Handle(pat.New("/*"), site)
	// Add middleware: they are evaluated in the reverse order in which they
	// are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	site.Use(middleware.Trace(logger, &handlerContext)) // injects trace id into the context
	site.Use(middleware.ContextLogger("milmove_trace_id", logger))
	site.Use(middleware.Recovery(logger))
	site.Use(middleware.SecurityHeaders(logger))

	if maxBodySize := v.GetInt64(cli.MaxBodySizeFlag); maxBodySize > 0 {
		site.Use(middleware.LimitBodySize(maxBodySize, logger))
	}

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

		logger.Info("Request", fields...)

	})

	staticMux := goji.SubMux()
	staticMux.Use(middleware.ValidMethodsStatic(logger))
	staticMux.Use(middleware.RequestLogger(logger))
	staticMux.Handle(pat.Get("/*"), clientHandler)
	// Needed to serve static paths (like favicon)
	staticMux.Handle(pat.Get(""), clientHandler)

	// Allow public content through without any auth or app checks
	site.Handle(pat.New("/static/*"), staticMux)
	site.Handle(pat.New("/downloads/*"), staticMux)
	site.Handle(pat.New("/favicon.ico"), staticMux)

	// Explicitly disable swagger.json route
	site.Handle(pat.Get("/swagger.json"), http.NotFoundHandler())
	if v.GetBool(cli.ServeSwaggerUIFlag) {
		logger.Info("Swagger UI static file serving is enabled")
		site.Handle(pat.Get("/swagger-ui/*"), staticMux)
	} else {
		site.Handle(pat.Get("/swagger-ui/*"), http.NotFoundHandler())
	}

	if v.GetBool(cli.ServeOrdersFlag) {
		ordersMux := goji.SubMux()
		ordersDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.OrdersServername)
		ordersMux.Use(ordersDetectionMiddleware)
		ordersMux.Use(middleware.NoCache(logger))
		ordersMux.Use(clientCertMiddleware)
		ordersMux.Use(middleware.RequestLogger(logger))
		ordersMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString(cli.OrdersSwaggerFlag)))
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Orders API Swagger UI serving is enabled")
			ordersMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "orders.html")))
		} else {
			ordersMux.Handle(pat.Get("/docs"), http.NotFoundHandler())
		}
		ordersMux.Handle(pat.New("/*"), ordersapi.NewOrdersAPIHandler(handlerContext))
		site.Handle(pat.New("/orders/v1/*"), ordersMux)
	}
	if v.GetBool(cli.ServeDPSFlag) {
		dpsMux := goji.SubMux()
		dpsDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.DpsServername)
		dpsMux.Use(dpsDetectionMiddleware)
		dpsMux.Use(middleware.NoCache(logger))
		dpsMux.Use(clientCertMiddleware)
		dpsMux.Use(middleware.RequestLogger(logger))
		dpsMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString(cli.DPSSwaggerFlag)))
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("DPS API Swagger UI serving is enabled")
			dpsMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "dps.html")))
		} else {
			dpsMux.Handle(pat.Get("/docs"), http.NotFoundHandler())
		}
		dpsMux.Handle(pat.New("/*"), dpsapi.NewDPSAPIHandler(handlerContext))
		site.Handle(pat.New("/dps/v0/*"), dpsMux)
	}

	if v.GetBool(cli.ServeSDDCFlag) {
		sddcDPSMux := goji.SubMux()
		sddcDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.SddcServername)
		sddcDPSMux.Use(sddcDetectionMiddleware)
		sddcDPSMux.Use(middleware.NoCache(logger))
		sddcDPSMux.Use(middleware.RequestLogger(logger))
		site.Handle(pat.New("/dps_auth/*"), sddcDPSMux)
		sddcDPSMux.Handle(pat.Get("/set_cookie"),
			dpsauth.NewSetCookieHandler(logger,
				dpsAuthSecretKey,
				dpsCookieDomain,
				dpsCookieSecret,
				dpsCookieExpires))
	}

	if v.GetBool(cli.ServePrimeFlag) {
		primeMux := goji.SubMux()
		primeDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.PrimeServername)
		primeMux.Use(primeDetectionMiddleware)
		if v.GetBool(cli.DevlocalAuthFlag) {
			devlocalClientCertMiddleware := authentication.DevlocalClientCertMiddleware(logger, dbConnection)
			primeMux.Use(devlocalClientCertMiddleware)
		} else {
			primeMux.Use(clientCertMiddleware)
		}
		primeMux.Use(authentication.PrimeAuthorizationMiddleware(logger))
		primeMux.Use(middleware.NoCache(logger))
		primeMux.Use(middleware.RequestLogger(logger))
		primeMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString(cli.PrimeSwaggerFlag)))
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Prime API Swagger UI serving is enabled")
			primeMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "prime.html")))
		} else {
			primeMux.Handle(pat.Get("/docs"), http.NotFoundHandler())
		}
		primeMux.Handle(pat.New("/*"), primeapi.NewPrimeAPIHandler(handlerContext))
		site.Handle(pat.New("/prime/v1/*"), primeMux)
	}

	if v.GetBool(cli.ServeSupportFlag) {
		supportMux := goji.SubMux()
		supportDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.PrimeServername)
		supportMux.Use(supportDetectionMiddleware)
		supportMux.Use(clientCertMiddleware)
		supportMux.Use(authentication.PrimeAuthorizationMiddleware(logger))
		supportMux.Use(middleware.NoCache(logger))
		supportMux.Use(middleware.RequestLogger(logger))
		supportMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString(cli.SupportSwaggerFlag)))
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Support API Swagger UI serving is enabled")
			supportMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "support.html")))
		} else {
			supportMux.Handle(pat.Get("/docs"), http.NotFoundHandler())
		}
		supportMux.Handle(pat.New("/*"), supportapi.NewSupportAPIHandler(handlerContext))
		site.Handle(pat.New("/support/v1/*"), supportMux)
	}

	// Handlers under mutual TLS need to go before this section that sets up middleware that shouldn't be enabled for mutual TLS (such as CSRF)
	root := goji.NewMux()
	root.Use(sessionCookieMiddleware)
	root.Use(middleware.RequestLogger(logger))

	debug := goji.SubMux()
	debug.Use(userAuthMiddleware)
	root.Handle(pat.New("/debug/pprof/*"), debug)
	if v.GetBool(cli.DebugPProfFlag) {
		logger.Info("Enabling pprof routes")
		debug.HandleFunc(pat.Get("/"), pprof.Index)
		debug.Handle(pat.Get("/allocs"), pprof.Handler("allocs"))
		debug.Handle(pat.Get("/block"), pprof.Handler("block"))
		debug.HandleFunc(pat.Get("/cmdline"), pprof.Cmdline)
		debug.Handle(pat.Get("/goroutine"), pprof.Handler("goroutine"))
		debug.Handle(pat.Get("/heap"), pprof.Handler("heap"))
		debug.Handle(pat.Get("/mutex"), pprof.Handler("mutex"))
		debug.HandleFunc(pat.Get("/profile"), pprof.Profile)
		debug.HandleFunc(pat.Get("/trace"), pprof.Trace)
		debug.Handle(pat.Get("/threadcreate"), pprof.Handler("threadcreate"))
		debug.HandleFunc(pat.Get("/symbol"), pprof.Symbol)
	} else {
		debug.HandleFunc(pat.Get("/*"), http.NotFound)
	}

	// CSRF path is set specifically at the root to avoid duplicate tokens from different paths
	csrfAuthKey, err := hex.DecodeString(v.GetString(cli.CSRFAuthKeyFlag))
	if err != nil {
		logger.Fatal("Failed to decode csrf auth key", zap.Error(err))
	}
	logger.Info("Enabling CSRF protection")
	root.Use(csrf.Protect(csrfAuthKey, csrf.Secure(!isDevOrTest), csrf.Path("/"), csrf.CookieName(auth.GorillaCSRFToken)))
	root.Use(maskedCSRFMiddleware)

	site.Handle(
		pat.New("/*"),
		milSession.LoadAndSave(adminSession.LoadAndSave(officeSession.LoadAndSave(root))))

	if v.GetBool(cli.ServeAPIInternalFlag) {
		internalMux := goji.SubMux()
		root.Handle(pat.New("/internal/*"), internalMux)
		internalMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString(cli.InternalSwaggerFlag)))
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Internal API Swagger UI serving is enabled")
			internalMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "internal.html")))
		} else {
			internalMux.Handle(pat.Get("/docs"), http.NotFoundHandler())
		}
		// Mux for internal API that enforces auth
		internalAPIMux := goji.SubMux()
		internalMux.HandleFunc(pat.Get("/users/is_logged_in"), isLoggedInMiddleware)
		internalMux.Handle(pat.New("/*"), internalAPIMux)
		internalAPIMux.Use(userAuthMiddleware)
		internalAPIMux.Use(middleware.NoCache(logger))
		api := internalapi.NewInternalAPI(handlerContext)
		internalAPIMux.Handle(pat.New("/*"), api.Serve(nil))
	}

	if v.GetBool(cli.ServeAdminFlag) {
		adminMux := goji.SubMux()
		root.Handle(pat.New("/admin/v1/*"), adminMux)
		adminMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString(cli.AdminSwaggerFlag)))
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("Admin API Swagger UI serving is enabled")
			adminMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "admin.html")))
		} else {
			adminMux.Handle(pat.Get("/docs"), http.NotFoundHandler())
		}

		// Mux for admin API that enforces auth
		adminAPIMux := goji.SubMux()
		adminMux.Handle(pat.New("/*"), adminAPIMux)
		adminAPIMux.Use(userAuthMiddleware)
		adminAPIMux.Use(authentication.AdminAuthMiddleware(logger))
		adminAPIMux.Use(middleware.NoCache(logger))
		adminAPIMux.Handle(pat.New("/*"), adminapi.NewAdminAPIHandler(handlerContext))
	}

	if v.GetBool(cli.ServeGHCFlag) {
		ghcMux := goji.SubMux()
		root.Handle(pat.New("/ghc/v1/*"), ghcMux)
		ghcMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString(cli.GHCSwaggerFlag)))
		if v.GetBool(cli.ServeSwaggerUIFlag) {
			logger.Info("GHC API Swagger UI serving is enabled")
			ghcMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "ghc.html")))
		} else {
			ghcMux.Handle(pat.Get("/docs"), http.NotFoundHandler())
		}

		// Mux for GHC API that enforces auth
		ghcAPIMux := goji.SubMux()
		ghcMux.Handle(pat.New("/*"), ghcAPIMux)
		ghcAPIMux.Use(userAuthMiddleware)
		ghcAPIMux.Use(middleware.NoCache(logger))
		api := ghcapi.NewGhcAPIHandler(handlerContext)
		ghcAPIMux.Handle(pat.New("/*"), api.Serve(nil))
	}

	authContext := authentication.NewAuthContext(logger, loginGovProvider, loginGovCallbackProtocol, loginGovCallbackPort, sessionManagers)
	authContext.SetFeatureFlag(
		authentication.FeatureFlag{
			Name:   cli.FeatureFlagAccessCode,
			Active: v.GetBool(cli.FeatureFlagAccessCode),
		},
	)
	authMux := goji.SubMux()
	root.Handle(pat.New("/auth/*"), authMux)
	authMux.Handle(pat.Get("/login-gov"), authentication.RedirectHandler{Context: authContext})
	authMux.Handle(pat.Get("/login-gov/callback"), authentication.NewCallbackHandler(authContext, dbConnection))
	authMux.Handle(pat.Post("/logout"), authentication.NewLogoutHandler(authContext, dbConnection))

	if v.GetBool(cli.DevlocalAuthFlag) {
		logger.Info("Enabling devlocal auth")
		localAuthMux := goji.SubMux()
		root.Handle(pat.New("/devlocal-auth/*"), localAuthMux)
		localAuthMux.Handle(pat.Get("/login"), authentication.NewUserListHandler(authContext, dbConnection))
		localAuthMux.Handle(pat.Post("/login"), authentication.NewAssignUserHandler(authContext, dbConnection, appnames))
		localAuthMux.Handle(pat.Post("/new"), authentication.NewCreateAndLoginUserHandler(authContext, dbConnection, appnames))
		localAuthMux.Handle(pat.Post("/create"), authentication.NewCreateUserHandler(authContext, dbConnection, appnames))

		if stringSliceContains([]string{cli.EnvironmentTest, cli.EnvironmentDevelopment, cli.EnvironmentReview}, v.GetString(cli.EnvironmentFlag)) {
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
	root.HandleFunc(pat.Get("/*"), indexHandler(build, logger))

	listenInterface := v.GetString(cli.InterfaceFlag)

	noTLSEnabled := v.GetBool(cli.NoTLSListenerFlag)
	var noTLSServer *server.NamedServer
	if noTLSEnabled {
		noTLSServer, err = server.CreateNamedServer(&server.CreateNamedServerInput{
			Name:        "no-tls",
			Host:        listenInterface,
			Port:        v.GetInt(cli.NoTLSPortFlag),
			Logger:      logger,
			HTTPHandler: bare,
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
			HTTPHandler:  bare,
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
			HTTPHandler:  bare,
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
	loggerSyncErr := logger.Sync()
	if loggerSyncErr != nil {
		logger.Error("Failed to sync logger", zap.Error(loggerSyncErr))
	}

	// Create a buffered channel that accepts 1 signal at a time.
	quit := make(chan os.Signal, 1)

	// Only send the SIGINT and SIGTERM signals to the quit channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait until the quit channel receieves a signal
	sig := <-quit

	logger.Info("received signal for graceful shutdown of server", zap.Any("signal", sig))

	// flush message that we received signal
	loggerSyncErr = logger.Sync()
	if loggerSyncErr != nil {
		logger.Error("Failed to sync logger", zap.Error(loggerSyncErr))
	}

	gracefulShutdownTimeout := v.GetDuration(cli.GracefulShutdownTimeoutFlag)

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	logger.Info("Waiting for listeners to be shutdown", zap.Duration("timeout", gracefulShutdownTimeout))

	// flush message that we are waiting on listeners
	loggerSyncErr = logger.Sync()
	if loggerSyncErr != nil {
		logger.Error("Failed to sync logger", zap.Error(loggerSyncErr))
	}

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
	loggerSyncErr = logger.Sync()
	if loggerSyncErr != nil {
		logger.Error("Failed to sync logger", zap.Error(loggerSyncErr))
	}

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

	loggerSyncErr = logger.Sync()
	if loggerSyncErr != nil {
		logger.Error("Failed to sync logger", zap.Error(loggerSyncErr))
	}

	if shutdownError {
		os.Exit(1)
	}

	return nil
}
