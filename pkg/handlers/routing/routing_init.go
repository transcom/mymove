package routing

import (
	"encoding/hex"
	"net/http"
	"net/http/pprof"
	"path"
	"path/filepath"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/spf13/afero"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/adminapi"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/handlers/ghcapi"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"github.com/transcom/mymove/pkg/handlers/ordersapi"
	"github.com/transcom/mymove/pkg/handlers/primeapi"
	"github.com/transcom/mymove/pkg/handlers/supportapi"
	"github.com/transcom/mymove/pkg/middleware"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/telemetry"
)

type Config struct {
	HandlerConfig handlers.HandlerConfig

	AuthContext authentication.Context

	// Use the afero filesystem interface to allow for replacement
	// during testing
	FileSystem afero.Fs

	// routing config

	// BuildRoot is where the client build is located (e.g. "build")
	BuildRoot string

	// If running in local development mode, where should uploaded
	// files be stored? LocalStorageRoot and LocalStorageWebRoot
	// configure that
	LocalStorageRoot    string
	LocalStorageWebRoot string

	// What is the maximum body size that should be accepted?
	MaxBodySize int64

	// To prevent CSRF, configure a authentication key
	CSRFAuthKey string

	// Should the swagger ui be served? Generally only enabled in development
	ServeSwaggerUI bool

	// Should the orders api be served? This is deprecated now
	ServeOrders bool
	// The path to the orders api swagger definition
	OrdersSwaggerPath string

	// Should the prime api be served?
	ServePrime bool
	// The path to the prime api swagger definition
	PrimeSwaggerPath string

	// Should the support api be served? Mostly only used in dev environments
	ServeSupport bool
	// The path to the support api swagger definition
	SupportSwaggerPath string

	// Should the API endpoint for profiling be enabled. Mostly only
	// used in dev environments
	ServeDebugPProf bool

	// Should the internal api be served?
	ServeAPIInternal bool
	// The path to the internal api swagger definition
	APIInternalSwaggerPath string

	// Should the admin api be served?
	ServeAdmin bool
	// The path to the admin api swagger definition
	AdminSwaggerPath string

	// Should the prime simulator be enabled? Definitely never enabled
	// in production
	ServePrimeSimulator bool

	// Should the ghc api be served?
	ServeGHC bool
	// The path to the ghc api swagger definition
	GHCSwaggerPath string

	// Should devlocal auth be enabled? Definitely never enabled in
	// production
	ServeDevlocalAuth bool

	// The git branch and commit used when building this server
	GitBranch string
	GitCommit string
}

// InitRouting sets up the routing
func InitRouting(appCtx appcontext.AppContext, redisPool *redis.Pool,
	routingConfig *Config, telemetryConfig *telemetry.Config) (http.Handler, error) {

	// site is the base
	site := mux.NewRouter()

	if routingConfig.LocalStorageRoot != "" && routingConfig.LocalStorageWebRoot != "" {
		localStorageHandlerFunc := storage.NewFilesystemHandler(routingConfig.LocalStorageRoot)

		site.HandleFunc(path.Join("/", routingConfig.LocalStorageWebRoot),
			localStorageHandlerFunc)
	}
	// Add middleware: they are evaluated in the reverse order in which they
	// are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	site.Use(middleware.Trace(appCtx.Logger())) // injects trace id into the context
	site.Use(middleware.ContextLogger("milmove_trace_id", appCtx.Logger()))
	site.Use(middleware.Recovery(appCtx.Logger()))
	site.Use(middleware.SecurityHeaders(appCtx.Logger()))

	if routingConfig.MaxBodySize > 0 {
		site.Use(middleware.LimitBodySize(routingConfig.MaxBodySize, appCtx.Logger()))
	}

	// Session management and authentication middleware
	sessionCookieMiddleware := auth.SessionCookieMiddleware(appCtx.Logger(), routingConfig.HandlerConfig.AppNames(), routingConfig.HandlerConfig.GetSessionManagers())
	maskedCSRFMiddleware := auth.MaskedCSRFMiddleware(appCtx.Logger(), routingConfig.HandlerConfig.UseSecureCookie())
	userAuthMiddleware := authentication.UserAuthMiddleware(appCtx.Logger())
	isLoggedInMiddleware := authentication.IsLoggedInMiddleware(appCtx.Logger())
	clientCertMiddleware := authentication.ClientCertMiddleware(appCtx)

	// Serves files out of build folder
	clientHandler := handlers.NewSpaHandler(
		routingConfig.BuildRoot,
		"index.html",
	)

	// Stub health check
	healthHandler := handlers.NewHealthHandler(appCtx, redisPool,
		routingConfig.GitBranch, routingConfig.GitCommit)
	site.HandleFunc("/health", healthHandler).Methods("GET")

	staticMux := site.PathPrefix("/static/").Subrouter()
	staticMux.Use(middleware.ValidMethodsStatic(appCtx.Logger()))
	staticMux.Use(middleware.RequestLogger(appCtx.Logger()))
	if telemetryConfig.Enabled {
		staticMux.Use(otelmux.Middleware("static"))
	}
	staticMux.PathPrefix("/").Handler(clientHandler).Methods("GET", "HEAD")

	downloadMux := site.PathPrefix("/downloads/").Subrouter()
	downloadMux.Use(middleware.ValidMethodsStatic(appCtx.Logger()))
	downloadMux.Use(middleware.RequestLogger(appCtx.Logger()))
	if telemetryConfig.Enabled {
		downloadMux.Use(otelmux.Middleware("download"))
	}
	downloadMux.PathPrefix("/").Handler(clientHandler).Methods("GET", "HEAD")

	site.Handle("/favicon.ico", clientHandler)

	// Explicitly disable swagger.json route
	site.Handle("/swagger.json", http.NotFoundHandler()).Methods("GET")
	if routingConfig.ServeSwaggerUI {
		appCtx.Logger().Info("Swagger UI static file serving is enabled")
		site.PathPrefix("/swagger-ui/").Handler(clientHandler).Methods("GET")
	} else {
		site.PathPrefix("/swagger-ui/").Handler(http.NotFoundHandler()).Methods("GET")
	}

	if routingConfig.ServeOrders {
		ordersServerName := routingConfig.HandlerConfig.AppNames().OrdersServername
		ordersMux := site.Host(ordersServerName).PathPrefix("/orders/v1/").Subrouter()
		ordersDetectionMiddleware := auth.HostnameDetectorMiddleware(appCtx.Logger(), ordersServerName)
		ordersMux.Use(ordersDetectionMiddleware)
		ordersMux.Use(middleware.NoCache(appCtx.Logger()))
		ordersMux.Use(clientCertMiddleware)
		ordersMux.Use(middleware.RequestLogger(appCtx.Logger()))
		ordersMux.HandleFunc("/swagger.yaml", handlers.NewFileHandler(routingConfig.OrdersSwaggerPath)).Methods("GET")
		if routingConfig.ServeSwaggerUI {
			appCtx.Logger().Info("Orders API Swagger UI serving is enabled")
			ordersMux.HandleFunc("/docs", handlers.NewFileHandler(path.Join(routingConfig.BuildRoot, "swagger-ui", "orders.html"))).Methods("GET")
		} else {
			ordersMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}
		api := ordersapi.NewOrdersAPI(routingConfig.HandlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		ordersMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if routingConfig.ServePrime {
		primeServerName := routingConfig.HandlerConfig.AppNames().PrimeServername
		primeMux := site.Host(primeServerName).PathPrefix("/prime/v1/").Subrouter()

		primeDetectionMiddleware := auth.HostnameDetectorMiddleware(appCtx.Logger(), primeServerName)
		primeMux.Use(primeDetectionMiddleware)
		if routingConfig.ServeDevlocalAuth {
			devlocalClientCertMiddleware := authentication.DevlocalClientCertMiddleware(appCtx)
			primeMux.Use(devlocalClientCertMiddleware)
		} else {
			primeMux.Use(clientCertMiddleware)
		}
		primeMux.Use(authentication.PrimeAuthorizationMiddleware(appCtx.Logger()))
		primeMux.Use(middleware.NoCache(appCtx.Logger()))
		primeMux.Use(middleware.RequestLogger(appCtx.Logger()))
		primeMux.HandleFunc("/swagger.yaml", handlers.NewFileHandler(routingConfig.PrimeSwaggerPath)).Methods("GET")
		if routingConfig.ServeSwaggerUI {
			appCtx.Logger().Info("Prime API Swagger UI serving is enabled")
			primeMux.HandleFunc("/docs", handlers.NewFileHandler(path.Join(routingConfig.BuildRoot, "swagger-ui", "prime.html"))).Methods("GET")
		} else {
			primeMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}
		api := primeapi.NewPrimeAPI(routingConfig.HandlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		primeMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if routingConfig.ServeSupport {
		primeServerName := routingConfig.HandlerConfig.AppNames().PrimeServername
		supportMux := site.Host(primeServerName).PathPrefix("/support/v1/").Subrouter()

		supportDetectionMiddleware := auth.HostnameDetectorMiddleware(appCtx.Logger(), primeServerName)
		supportMux.Use(supportDetectionMiddleware)
		supportMux.Use(clientCertMiddleware)
		supportMux.Use(authentication.PrimeAuthorizationMiddleware(appCtx.Logger()))
		supportMux.Use(middleware.NoCache(appCtx.Logger()))
		supportMux.Use(middleware.RequestLogger(appCtx.Logger()))
		supportMux.HandleFunc("/swagger.yaml", handlers.NewFileHandler(routingConfig.SupportSwaggerPath)).Methods("GET")
		if routingConfig.ServeSwaggerUI {
			appCtx.Logger().Info("Support API Swagger UI serving is enabled")
			supportMux.HandleFunc("/docs", handlers.NewFileHandler(path.Join(routingConfig.BuildRoot, "swagger-ui", "support.html"))).Methods("GET")
		} else {
			supportMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}
		supportMux.PathPrefix("/").Handler(supportapi.NewSupportAPIHandler(routingConfig.HandlerConfig))
	}

	// Handlers under mutual TLS need to go before this section that sets up middleware that shouldn't be enabled for mutual TLS (such as CSRF)
	root := mux.NewRouter()
	root.Use(sessionCookieMiddleware)
	root.Use(middleware.RequestLogger(appCtx.Logger()))

	debug := root.PathPrefix("/debug/pprof/").Subrouter()
	debug.Use(userAuthMiddleware)
	if routingConfig.ServeDebugPProf {
		appCtx.Logger().Info("Enabling pprof routes")
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
	csrfAuthKey, err := hex.DecodeString(routingConfig.CSRFAuthKey)
	if err != nil {
		appCtx.Logger().Fatal("Failed to decode csrf auth key", zap.Error(err))
	}
	appCtx.Logger().Info("Enabling CSRF protection")
	root.Use(csrf.Protect(csrfAuthKey, csrf.Secure(routingConfig.HandlerConfig.UseSecureCookie()), csrf.Path("/"), csrf.CookieName(auth.GorillaCSRFToken)))
	root.Use(maskedCSRFMiddleware)

	site.Host(routingConfig.HandlerConfig.AppNames().MilServername).PathPrefix("/").Handler(routingConfig.HandlerConfig.GetMilSessionManager().LoadAndSave(root))
	site.Host(routingConfig.HandlerConfig.AppNames().AdminServername).PathPrefix("/").Handler(routingConfig.HandlerConfig.GetAdminSessionManager().LoadAndSave(root))
	site.Host(routingConfig.HandlerConfig.AppNames().OfficeServername).PathPrefix("/").Handler(routingConfig.HandlerConfig.GetOfficeSessionManager().LoadAndSave(root))

	if routingConfig.ServeAPIInternal {
		internalMux := root.PathPrefix("/internal/").Subrouter()
		internalMux.HandleFunc("/swagger.yaml", handlers.NewFileHandler(routingConfig.APIInternalSwaggerPath)).Methods("GET")
		if routingConfig.ServeSwaggerUI {
			appCtx.Logger().Info("Internal API Swagger UI serving is enabled")
			internalMux.HandleFunc("/docs", handlers.NewFileHandler(path.Join(routingConfig.BuildRoot, "swagger-ui", "internal.html"))).Methods("GET")
		} else {
			internalMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}
		internalMux.HandleFunc("/users/is_logged_in", isLoggedInMiddleware).Methods("GET")
		// Mux for internal API that enforces auth
		internalAPIMux := internalMux.PathPrefix("/").Subrouter()
		internalAPIMux.Use(userAuthMiddleware)
		internalAPIMux.Use(middleware.NoCache(appCtx.Logger()))
		api := internalapi.NewInternalAPI(routingConfig.HandlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		internalAPIMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if routingConfig.ServeAdmin {
		adminMux := root.PathPrefix("/admin/v1/").Subrouter()

		adminMux.HandleFunc("/swagger.yaml", handlers.NewFileHandler(routingConfig.AdminSwaggerPath)).Methods("GET")
		if routingConfig.ServeSwaggerUI {
			appCtx.Logger().Info("Admin API Swagger UI serving is enabled")
			adminMux.HandleFunc("/docs", handlers.NewFileHandler(path.Join(routingConfig.BuildRoot, "swagger-ui", "admin.html"))).Methods("GET")
		} else {
			adminMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}

		// Mux for admin API that enforces auth
		adminAPIMux := adminMux.PathPrefix("/").Subrouter()
		adminAPIMux.Use(userAuthMiddleware)
		adminAPIMux.Use(authentication.AdminAuthMiddleware(appCtx.Logger()))
		adminAPIMux.Use(middleware.NoCache(appCtx.Logger()))
		api := adminapi.NewAdminAPI(routingConfig.HandlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		adminAPIMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if routingConfig.ServePrimeSimulator {
		// attach prime simulator API to root so cookies are handled
		officeServerName := routingConfig.HandlerConfig.AppNames().OfficeServername
		primeSimulatorMux := root.Host(officeServerName).PathPrefix("/prime/v1/").Subrouter()
		primeSimulatorMux.HandleFunc("/swagger.yaml", handlers.NewFileHandler(routingConfig.PrimeSwaggerPath)).Methods("GET")
		if routingConfig.ServeSwaggerUI {
			appCtx.Logger().Info("Prime Simulator API Swagger UI serving is enabled")
			primeSimulatorMux.HandleFunc("/docs", handlers.NewFileHandler(path.Join(routingConfig.BuildRoot, "swagger-ui", "prime.html"))).Methods("GET")
		} else {
			primeSimulatorMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}

		// Mux for prime simulator API that enforces auth
		primeSimulatorAPIMux := primeSimulatorMux.PathPrefix("/").Subrouter()
		primeSimulatorAPIMux.Use(userAuthMiddleware)
		primeSimulatorAPIMux.Use(authentication.PrimeSimulatorAuthorizationMiddleware(appCtx.Logger()))
		primeSimulatorAPIMux.Use(middleware.NoCache(appCtx.Logger()))
		api := primeapi.NewPrimeAPI(routingConfig.HandlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		primeSimulatorAPIMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	if routingConfig.ServeGHC {
		ghcMux := root.PathPrefix("/ghc/v1/").Subrouter()
		ghcMux.HandleFunc("/swagger.yaml", handlers.NewFileHandler(routingConfig.GHCSwaggerPath)).Methods("GET")
		if routingConfig.ServeSwaggerUI {
			appCtx.Logger().Info("GHC API Swagger UI serving is enabled")
			ghcMux.HandleFunc("/docs", handlers.NewFileHandler(path.Join(routingConfig.BuildRoot, "swagger-ui", "ghc.html"))).Methods("GET")
		} else {
			ghcMux.Handle("/docs", http.NotFoundHandler()).Methods("GET")
		}

		// Mux for GHC API that enforces auth
		ghcAPIMux := ghcMux.PathPrefix("/").Subrouter()
		ghcAPIMux.Use(userAuthMiddleware)
		ghcAPIMux.Use(middleware.NoCache(appCtx.Logger()))
		api := ghcapi.NewGhcAPIHandler(routingConfig.HandlerConfig)
		tracingMiddleware := middleware.OpenAPITracing(api)
		ghcAPIMux.PathPrefix("/").Handler(api.Serve(tracingMiddleware))
	}

	authMux := root.PathPrefix("/auth/").Subrouter()
	authMux.Use(middleware.NoCache(appCtx.Logger()))
	authMux.Use(otelmux.Middleware("auth"))
	authMux.Handle("/login-gov", authentication.NewRedirectHandler(routingConfig.AuthContext, routingConfig.HandlerConfig, routingConfig.HandlerConfig.UseSecureCookie())).Methods("GET")
	authMux.Handle("/login-gov/callback", authentication.NewCallbackHandler(routingConfig.AuthContext, routingConfig.HandlerConfig, routingConfig.HandlerConfig.NotificationSender())).Methods("GET")
	authMux.Handle("/logout", authentication.NewLogoutHandler(routingConfig.AuthContext, routingConfig.HandlerConfig)).Methods("POST")

	if routingConfig.ServeDevlocalAuth {
		appCtx.Logger().Info("Enabling devlocal auth")
		localAuthMux := root.PathPrefix("/devlocal-auth/").Subrouter()
		localAuthMux.Use(middleware.NoCache(appCtx.Logger()))
		localAuthMux.Use(otelmux.Middleware("devlocal"))
		localAuthMux.Handle("/login", authentication.NewUserListHandler(routingConfig.AuthContext, routingConfig.HandlerConfig)).Methods("GET")
		localAuthMux.Handle("/login", authentication.NewAssignUserHandler(routingConfig.AuthContext, routingConfig.HandlerConfig, routingConfig.HandlerConfig.AppNames())).Methods("POST")
		localAuthMux.Handle("/new", authentication.NewCreateAndLoginUserHandler(routingConfig.AuthContext, routingConfig.HandlerConfig, routingConfig.HandlerConfig.AppNames())).Methods("POST")
		localAuthMux.Handle("/create", authentication.NewCreateUserHandler(routingConfig.AuthContext, routingConfig.HandlerConfig, routingConfig.HandlerConfig.AppNames())).Methods("POST")

	}

	// Serve index.html to all requests that haven't matches a previous route,
	root.PathPrefix("/").Handler(indexHandler(routingConfig, appCtx.Logger())).Methods("GET", "HEAD")

	return site, nil
}

// indexHandler returns a handler that will serve the resulting content
func indexHandler(routingConfig *Config, globalLogger *zap.Logger) http.HandlerFunc {

	indexPath := path.Join(routingConfig.BuildRoot, "index.html")
	reader, err := routingConfig.FileSystem.Open(filepath.Clean(indexPath))
	if err != nil {
		globalLogger.Fatal("could not read index.html template: run make client_build", zap.Error(err))
	}

	stat, err := routingConfig.FileSystem.Stat(indexPath)
	if err != nil {
		globalLogger.Fatal("could not stat index.html template", zap.Error(err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "index.html", stat.ModTime(), reader)
	}
}
