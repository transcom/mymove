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
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/csrf"
	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	goji "goji.io"
	"goji.io/pat"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/ecs"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/adminapi"
	"github.com/transcom/mymove/pkg/handlers/dpsapi"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"github.com/transcom/mymove/pkg/handlers/ordersapi"
	"github.com/transcom/mymove/pkg/handlers/publicapi"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/middleware"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/invoice"
	"github.com/transcom/mymove/pkg/storage"
)

// GitCommit is empty unless set as a build flag
// See https://blog.alexellis.io/inject-build-time-vars-golang/
var gitBranch string
var gitCommit string

type errInvalidHost struct {
	Host string
}

func (e *errInvalidHost) Error() string {
	return fmt.Sprintf("invalid host %s, must not contain whitespace, :, /, or \\", e.Host)
}

// initServeFlags - Order matters!
func initServeFlags(flag *pflag.FlagSet) {

	// Build Server
	cli.InitBuildFlags(flag)

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

	// Login.Gov Auth config
	cli.InitAuthFlags(flag)

	// HERE Route Config
	cli.InitRouteFlags(flag)

	// EDI Invoice Config
	cli.InitGEXFlags(flag)

	// Storage
	cli.InitStorageFlags(flag)

	// Email
	cli.InitEmailFlags(flag)

	// Honeycomb Config
	cli.InitHoneycombFlags(flag)

	// IWS
	cli.InitIWSFlags(flag)

	// DB Config
	cli.InitDatabaseFlags(flag)

	// CSRF Protection
	cli.InitCSRFFlags(flag)

	// Middleware
	cli.InitMiddlewareFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Info("checking webserver config")

	if err := cli.CheckBuild(v); err != nil {
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

	if err := cli.CheckPorts(v); err != nil {
		return err
	}

	if err := cli.CheckAuth(v); err != nil {
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

	if err := cli.CheckHoneycomb(v); err != nil {
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

	if err := cli.CheckVerbose(v); err != nil {
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

func versionFunction(cmd *cobra.Command, args []string) error {
	str, err := json.Marshal(map[string]interface{}{
		"gitBranch": gitBranch,
		"gitCommit": gitCommit,
	})
	if err != nil {
		return err
	}
	fmt.Println(string(str))
	return nil
}

func serveFunction(cmd *cobra.Command, args []string) error {

	err := cmd.ParseFlags(os.Args[1:])
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

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
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

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	isDevOrTest := dbEnv == "development" || dbEnv == "test"
	if isDevOrTest {
		logger.Info(fmt.Sprintf("Starting in %s mode, which enables additional features", dbEnv))
	}

	// Honeycomb initialization also initializes beeline, so keep this near the top of the stack
	useHoneycomb := cli.InitHoneycomb(v, logger)

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

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, logger)
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

	// Collect the servernames into a handy struct
	appnames := auth.ApplicationServername{
		MilServername:    v.GetString(cli.HTTPMyServerNameFlag),
		OfficeServername: v.GetString(cli.HTTPOfficeServerNameFlag),
		TspServername:    v.GetString(cli.HTTPTSPServerNameFlag),
		AdminServername:  v.GetString(cli.HTTPAdminServerNameFlag),
		OrdersServername: v.GetString(cli.HTTPOrdersServerNameFlag),
		DpsServername:    v.GetString(cli.HTTPDPSServerNameFlag),
		SddcServername:   v.GetString(cli.HTTPSDDCServerNameFlag),
	}

	// Register Login.gov authentication provider for My.(move.mil)
	loginGovProvider, err := cli.InitAuth(v, logger, appnames)
	if err != nil {
		logger.Fatal("Registering login provider", zap.Error(err))
	}

	useSecureCookie := !isDevOrTest
	// Session management and authentication middleware
	noSessionTimeout := v.GetBool(cli.NoSessionTimeoutFlag)
	sessionCookieMiddleware := auth.SessionCookieMiddleware(logger, clientAuthSecretKey, noSessionTimeout, appnames, useSecureCookie)
	maskedCSRFMiddleware := auth.MaskedCSRFMiddleware(logger, useSecureCookie)
	userAuthMiddleware := authentication.UserAuthMiddleware(logger)
	clientCertMiddleware := authentication.ClientCertMiddleware(logger, dbConnection)

	handlerContext := handlers.NewHandlerContext(dbConnection, logger)
	handlerContext.SetCookieSecret(clientAuthSecretKey)
	handlerContext.SetUseSecureCookie(useSecureCookie)
	if noSessionTimeout {
		handlerContext.SetNoSessionTimeout()
	}

	// Email
	notificationSender := cli.InitEmail(v, logger)
	handlerContext.SetNotificationSender(notificationSender)

	build := v.GetString(cli.BuildFlag)

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(build))

	// Get route planner for handlers to calculate transit distances
	// routePlanner := route.NewBingPlanner(logger, bingMapsEndpoint, bingMapsKey)
	routePlanner := cli.InitRoutePlanner(v, logger)
	handlerContext.SetPlanner(routePlanner)

	// Set SendProductionInvoice for ediinvoice
	handlerContext.SetSendProductionInvoice(v.GetBool(cli.GEXSendProdInvoiceFlag))

	// Storage
	storer := cli.InitStorage(v, logger)
	handlerContext.SetFileStorer(storer)

	certificates, rootCAs, err := cli.InitDoDCertificates(v, logger)
	if certificates == nil || rootCAs == nil || err != nil {
		logger.Fatal("Failed to initialize DOD certificates", zap.Error(err))
	}

	logger.Debug("Server DOD Key Pair Loaded")
	logger.Debug("Trusted Certificate Authorities", zap.Any("subjects", rootCAs.Subjects()))

	// Set the GexSender() and GexSender fields
	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs}
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
			&tls.Config{},
			"",
			"",
		)
	} else {
		gexRequester = invoice.NewGexSenderHTTP(
			gexURL,
			true,
			tlsConfig,
			v.GetString(cli.GEXBasicAuthUsernameFlag),
			v.GetString(cli.GEXBasicAuthPasswordFlag),
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

	rbs, err := cli.InitRBSPersonLookup(v, logger)
	if err != nil {
		logger.Fatal("Could not instantiate IWS RBS", zap.Error(err))
	}
	handlerContext.SetIWSPersonLookup(*rbs)

	dpsAuthSecretKey := v.GetString(cli.DPSAuthSecretKeyFlag)
	dpsCookieDomain := v.GetString(cli.DPSCookieDomainFlag)
	dpsCookieSecret := []byte(v.GetString(cli.DPSAuthCookieSecretKeyFlag))
	dpsCookieExpires := v.GetInt(cli.DPSCookieExpiresInMinutesFlag)

	dpsAuthParams := cli.InitDPSAuthParams(v, appnames)
	handlerContext.SetDPSAuthParams(dpsAuthParams)

	// Base routes
	site := goji.NewMux()
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

		logger.Info("Request", fields...)

	})

	staticMux := goji.SubMux()
	staticMux.Use(middleware.ValidMethodsStatic(logger))
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

	ordersMux := goji.SubMux()
	ordersDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.OrdersServername)
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

	dpsMux := goji.SubMux()
	dpsDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.DpsServername)
	dpsMux.Use(dpsDetectionMiddleware)
	dpsMux.Use(middleware.NoCache(logger))
	dpsMux.Use(clientCertMiddleware)
	dpsMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString(cli.DPSSwaggerFlag)))
	if v.GetBool(cli.ServeSwaggerUIFlag) {
		logger.Info("DPS API Swagger UI serving is enabled")
		dpsMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "dps.html")))
	} else {
		dpsMux.Handle(pat.Get("/docs"), http.NotFoundHandler())
	}
	dpsMux.Handle(pat.New("/*"), dpsapi.NewDPSAPIHandler(handlerContext))
	site.Handle(pat.New("/dps/v0/*"), dpsMux)

	sddcDPSMux := goji.SubMux()
	sddcDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, appnames.SddcServername)
	sddcDPSMux.Use(sddcDetectionMiddleware)
	sddcDPSMux.Use(middleware.NoCache(logger))
	site.Handle(pat.New("/dps_auth/*"), sddcDPSMux)
	sddcDPSMux.Handle(pat.Get("/set_cookie"),
		dpsauth.NewSetCookieHandler(logger,
			dpsAuthSecretKey,
			dpsCookieDomain,
			dpsCookieSecret,
			dpsCookieExpires))

	root := goji.NewMux()
	root.Use(middleware.Recovery(logger))
	root.Use(sessionCookieMiddleware)
	root.Use(middleware.RequestLogger(logger))

	// CSRF path is set specifically at the root to avoid duplicate tokens from different paths
	csrfAuthKey, err := hex.DecodeString(v.GetString(cli.CSRFAuthKeyFlag))
	if err != nil {
		logger.Fatal("Failed to decode csrf auth key", zap.Error(err))
	}
	logger.Info("Enabling CSRF protection")
	root.Use(csrf.Protect(csrfAuthKey, csrf.Secure(!isDevOrTest), csrf.Path("/"), csrf.CookieName(auth.GorillaCSRFToken)))
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
	apiMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString(cli.SwaggerFlag)))
	if v.GetBool(cli.ServeSwaggerUIFlag) {
		logger.Info("Public API Swagger UI serving is enabled")
		apiMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "api.html")))
	} else {
		apiMux.Handle(pat.Get("/docs"), http.NotFoundHandler())
	}
	externalAPIMux := goji.SubMux()
	apiMux.Handle(pat.New("/*"), externalAPIMux)
	externalAPIMux.Use(middleware.NoCache(logger))
	externalAPIMux.Use(userAuthMiddleware)
	externalAPIMux.Handle(pat.New("/*"), publicapi.NewPublicAPIHandler(handlerContext))

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
	internalMux.Handle(pat.New("/*"), internalAPIMux)
	internalAPIMux.Use(userAuthMiddleware)
	internalAPIMux.Use(middleware.NoCache(logger))
	internalAPIMux.Handle(pat.New("/*"), internalapi.NewInternalAPIHandler(handlerContext))

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
	adminAPIMux.Use(middleware.NoCache(logger))
	adminAPIMux.Handle(pat.New("/*"), adminapi.NewAdminAPIHandler(handlerContext))

	authContext := authentication.NewAuthContext(logger, loginGovProvider, loginGovCallbackProtocol, loginGovCallbackPort)
	authMux := goji.SubMux()
	root.Handle(pat.New("/auth/*"), authMux)
	authMux.Handle(pat.Get("/login-gov"), authentication.RedirectHandler{Context: authContext})
	authMux.Handle(pat.Get("/login-gov/callback"), authentication.NewCallbackHandler(authContext, dbConnection, clientAuthSecretKey, noSessionTimeout, useSecureCookie))
	authMux.Handle(pat.Post("/logout"), authentication.NewLogoutHandler(authContext, clientAuthSecretKey, noSessionTimeout, useSecureCookie))

	if isDevOrTest {
		logger.Info("Enabling devlocal auth")
		localAuthMux := goji.SubMux()
		root.Handle(pat.New("/devlocal-auth/*"), localAuthMux)
		localAuthMux.Handle(pat.Get("/login"), authentication.NewUserListHandler(authContext, dbConnection))
		localAuthMux.Handle(pat.Post("/login"), authentication.NewAssignUserHandler(authContext, dbConnection, appnames, clientAuthSecretKey, noSessionTimeout, useSecureCookie))
		localAuthMux.Handle(pat.Post("/new"), authentication.NewCreateAndLoginUserHandler(authContext, dbConnection, appnames, clientAuthSecretKey, noSessionTimeout, useSecureCookie))
		localAuthMux.Handle(pat.Post("/create"), authentication.NewCreateUserHandler(authContext, dbConnection, appnames, clientAuthSecretKey, noSessionTimeout, useSecureCookie))

		devlocalCAPath := v.GetString(cli.DevlocalCAFlag)
		devlocalCa, readFileErr := ioutil.ReadFile(devlocalCAPath) // #nosec
		if readFileErr != nil {
			logger.Error(fmt.Sprintf("Unable to read devlocal CA from path %s", devlocalCAPath), zap.Error(readFileErr))
		} else {
			rootCAs.AppendCertsFromPEM(devlocalCa)
		}
	}

	storageBackend := v.GetString(cli.StorageBackendFlag)
	if storageBackend == "local" {
		localStorageRoot := v.GetString(cli.LocalStorageRootFlag)
		localStorageWebRoot := v.GetString(cli.LocalStorageWebRootFlag)

		// Add a file handler to provide access to files uploaded in development
		fs := storage.NewFilesystemHandler(localStorageRoot)
		root.Handle(pat.Get(path.Join("/", localStorageWebRoot, "/*")), fs)
	}

	// Serve index.html to all requests that haven't matches a previous route,
	root.HandleFunc(pat.Get("/*"), indexHandler(build, logger))

	var httpHandler http.Handler
	if useHoneycomb {
		httpHandler = hnynethttp.WrapHandler(site)
	} else {
		httpHandler = site
	}

	listenInterface := v.GetString(cli.InterfaceFlag)

	noTLSServer, err := server.CreateNamedServer(&server.CreateNamedServerInput{
		Name:        "no-tls",
		Host:        listenInterface,
		Port:        v.GetInt(cli.NoTLSPortFlag),
		Logger:      logger,
		HTTPHandler: httpHandler,
	})
	if err != nil {
		logger.Fatal("error creating no-tls server", zap.Error(err))
	}
	go startListener(noTLSServer, logger, false)

	tlsServer, err := server.CreateNamedServer(&server.CreateNamedServerInput{
		Name:         "tls",
		Host:         listenInterface,
		Port:         v.GetInt(cli.TLSPortFlag),
		Logger:       logger,
		HTTPHandler:  httpHandler,
		ClientAuth:   tls.NoClientCert,
		Certificates: certificates,
	})
	if err != nil {
		logger.Fatal("error creating tls server", zap.Error(err))
	}
	go startListener(tlsServer, logger, true)

	mutualTLSServer, err := server.CreateNamedServer(&server.CreateNamedServerInput{
		Name:         "mutual-tls",
		Host:         listenInterface,
		Port:         v.GetInt(cli.MutualTLSPortFlag),
		Logger:       logger,
		HTTPHandler:  httpHandler,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: certificates,
		ClientCAs:    rootCAs,
	})
	if err != nil {
		logger.Fatal("error creating mutual-tls server", zap.Error(err))
	}
	go startListener(mutualTLSServer, logger, true)

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

	wg.Add(1)
	go func() {
		shutdownErrors.Store(noTLSServer, noTLSServer.Shutdown(ctx))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		shutdownErrors.Store(tlsServer, tlsServer.Shutdown(ctx))
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		shutdownErrors.Store(mutualTLSServer, mutualTLSServer.Shutdown(ctx))
		wg.Done()
	}()

	wg.Wait()
	logger.Info("All listeners are shutdown")
	logger.Sync()

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
	logger.Sync()

	if shutdownError {
		os.Exit(1)
	}

	return nil
}

func main() {

	root := cobra.Command{
		Use:   "milmove [flags]",
		Short: "Webserver for MilMove",
		Long:  "Webserver for MilMove",
	}

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information to stdout",
		Long:  "Print version information to stdout",
		RunE:  versionFunction,
	})

	serveCommand := &cobra.Command{
		Use:   "serve",
		Short: "Runs MilMove webserver",
		Long:  "Runs MilMove webserver",
		RunE:  serveFunction,
	}
	initServeFlags(serveCommand.Flags())
	root.AddCommand(serveCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\nmilmove completion > /usr/local/etc/bash_completion.d/milmove",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
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
