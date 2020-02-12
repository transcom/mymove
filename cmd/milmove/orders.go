package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gobuffalo/pop"
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
	"github.com/transcom/mymove/pkg/ecs"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ordersapi"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/middleware"
	"github.com/transcom/mymove/pkg/server"
)

// initOrdersFlags - Order matters!
func initOrdersFlags(flag *pflag.FlagSet) {

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

	// IWS
	cli.InitIWSFlags(flag)

	// DB Config
	cli.InitDatabaseFlags(flag)

	// aws-vault
	cli.InitVaultFlags(flag)

	// Logging
	cli.InitLoggingFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// pprof flags
	cli.InitDebugFlags(flag)

	// Service Flags
	cli.InitServiceFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

func checkOrdersConfig(v *viper.Viper, logger logger) error {

	logger.Info("checking webserver config")

	if err := cli.CheckEnvironment(v); err != nil {
		return err
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

	if err := cli.CheckIWS(v); err != nil {
		return err
	}

	if err := cli.CheckDatabase(v, logger); err != nil {
		return err
	}

	if err := cli.CheckVault(v); err != nil {
		return err
	}

	if err := cli.CheckLogging(v); err != nil {
		return err
	}

	if err := cli.CheckVerbose(v); err != nil {
		return err
	}

	if err := cli.CheckDebugFlags(v); err != nil {
		return err
	}

	if err := cli.CheckServices(v); err != nil {
		return err
	}

	return nil
}

func ordersFunction(cmd *cobra.Command, args []string) error {

	var logger *zap.Logger
	var dbConnection *pop.Connection
	dbClose := &sync.Once{}

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
			logger.Sync()
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

	logger, err = logging.Config(v.GetString(cli.LoggingEnvFlag), v.GetBool(cli.VerboseFlag))
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

	err = checkOrdersConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	if v.GetBool(cli.DbDebugFlag) {
		pop.Debug = true
	}

	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) {
		c, errorConfig := cli.GetAWSConfig(v, v.GetBool(cli.VerboseFlag))
		if errorConfig != nil {
			logger.Fatal(errors.Wrap(errorConfig, "error creating aws config").Error())
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

	// Session management and authentication middleware
	clientCertMiddleware := authentication.ClientCertMiddleware(logger, dbConnection)

	handlerContext := handlers.NewHandlerContext(dbConnection, logger)

	build := v.GetString(cli.BuildRootFlag)

	certificates, rootCAs, err := certs.InitDoDCertificates(v, logger)
	if certificates == nil || rootCAs == nil || err != nil {
		logger.Fatal("Failed to initialize DOD certificates", zap.Error(err))
	}

	logger.Debug("Server DOD Key Pair Loaded")
	logger.Debug("Trusted Certificate Authorities", zap.Any("subjects", rootCAs.Subjects()))

	rbs, err := iws.InitRBSPersonLookup(v, logger)
	if err != nil {
		logger.Fatal("Could not instantiate IWS RBS", zap.Error(err))
	}
	handlerContext.SetIWSPersonLookup(*rbs)

	// bare is the base muxer. Not intended to have any middleware attached.
	bare := goji.NewMux()

	// Base routes
	site := goji.SubMux()
	bare.Handle(pat.New("/*"), site)
	// Add middleware: they are evaluated in the reverse order in which they
	// are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
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
			dbErr := dbConnection.RawQuery("SELECT 1;").Exec()
			if dbErr != nil {
				logger.Error("Failed database health check", zap.Error(dbErr))
			}
			data["database"] = dbErr == nil
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

	// Explicitly disable swagger.json route
	site.Handle(pat.Get("/swagger.json"), http.NotFoundHandler())
	if v.GetBool(cli.ServeSwaggerUIFlag) {
		logger.Info("Swagger UI static file serving is enabled")
		site.Handle(pat.Get("/swagger-ui/*"), staticMux)
	} else {
		site.Handle(pat.Get("/swagger-ui/*"), http.NotFoundHandler())
	}

	if v.GetBool(cli.ServeOrdersFlag) {
		OrdersServername := v.GetString(cli.HTTPOrdersServerNameFlag)

		ordersMux := goji.SubMux()
		ordersDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, OrdersServername)
		ordersMux.Use(ordersDetectionMiddleware)
		ordersMux.Use(middleware.NoCache(logger))
		ordersMux.Use(clientCertMiddleware)
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

	// Handlers under mutual TLS need to go before this section that sets up middleware that shouldn't be enabled for mutual TLS (such as CSRF)
	root := goji.NewMux()
	root.Use(middleware.Recovery(logger))
	root.Use(middleware.Trace(logger, &handlerContext))            // injects http request trace id
	root.Use(middleware.ContextLogger("milmove_trace_id", logger)) // injects http request logger
	root.Use(middleware.RequestLogger(logger))

	site.Handle(pat.New("/*"), root)

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
	logger.Sync()

	// Create a buffered channel that accepts 1 signal at a time.
	quit := make(chan os.Signal, 1)

	// Only send the SIGINT and SIGTERM signals to the quit channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait until the quit channel receieves a signal
	sig := <-quit

	logger.Info("received signal for graceful shutdown of server", zap.Any("signal", sig))

	// flush message that we received signal
	logger.Sync()

	gracefulShutdownTimeout := v.GetDuration(cli.GracefulShutdownTimeoutFlag)

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	logger.Info("Waiting for listeners to be shutdown", zap.Duration("timeout", gracefulShutdownTimeout))

	// flush message that we are waiting on listeners
	logger.Sync()

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
	logger.Sync()

	var dbCloseErr error
	dbClose.Do(func() {
		logger.Info("closing database connections")
		dbCloseErr = dbConnection.Close()
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

	logger.Sync()

	if shutdownError {
		os.Exit(1)
	}

	return nil
}
