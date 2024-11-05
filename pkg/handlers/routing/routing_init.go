package routing

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/csrf"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/adminapi"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/handlers/ghcapi"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"github.com/transcom/mymove/pkg/handlers/pptasapi"
	"github.com/transcom/mymove/pkg/handlers/primeapi"
	"github.com/transcom/mymove/pkg/handlers/primeapiv2"
	"github.com/transcom/mymove/pkg/handlers/primeapiv3"
	"github.com/transcom/mymove/pkg/handlers/supportapi"
	"github.com/transcom/mymove/pkg/handlers/testharnessapi"
	"github.com/transcom/mymove/pkg/logging"
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

	// BuildRoot is where the client build is located (e.g. "build")
	BuildRoot string

	// If running in local development mode, where should uploaded
	// files be stored? LocalStorageRoot and LocalStorageWebRoot
	// configure that
	LocalStorageRoot    string
	LocalStorageWebRoot string

	// What is the maximum body size that should be accepted?
	MaxBodySize int64

	// Should serve client collector endpoint?
	ServeClientCollector bool

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
	// The path to the prime V2 api swagger definition
	PrimeV2SwaggerPath string
	// The path to the prime V3 api swagger definition
	PrimeV3SwaggerPath string

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

	// Should the pptas api be served?
	ServePPTAS bool
	// The path to the pptas api swagger definition
	PPTASSwaggerPath string

	// Should devlocal auth be enabled? Definitely never enabled in
	// production
	ServeDevlocalAuth bool

	// The git branch and commit used when building this server
	GitBranch string
	GitCommit string

	// To prevent CSRF, configure a CSRF Middlware
	// Configuring it here lets us re-use it / override it effectively
	// in tests
	CSRFMiddleware func(http.Handler) http.Handler
}

func InitCSRFMiddlware(csrfAuthKey []byte, secure bool, path string, cookieName string) func(http.Handler) http.Handler {
	return csrf.Protect(csrfAuthKey,
		csrf.Secure(secure),
		csrf.Path(path),
		csrf.CookieName(cookieName),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reason := csrf.FailureReason(r)
			logger := logging.FromContext(r.Context())
			logger.Info("Request forbidden by CSRF Middleware", zap.String("reason", reason.Error()))
			http.Error(w, fmt.Sprintf("%s - %s",
				http.StatusText(http.StatusForbidden), reason),
				http.StatusForbidden)
		})),
	)
}

// a custom host router that ignores the port
type HostRouter struct {
	routes map[string]chi.Router
}

// make sure HostRoutes implements chi.Routes
var _ chi.Routes = &HostRouter{}

func NewHostRouter() *HostRouter {
	return &HostRouter{
		routes: make(map[string]chi.Router),
	}
}

// Map adds a chi.Router for a hostname
func (hr *HostRouter) Map(host string, r chi.Router) {
	hr.routes[strings.ToLower(host)] = r
}

func (hr *HostRouter) Match(_ *chi.Context, _ string, _ string) bool {
	// the chi.Context does not contain information about which host
	// was requested and the host router is not distinguishing based
	// on method or path, so the host router always matches every
	// route. The per host dispatch happens below in ServeHTTP
	return true
}

func (hr *HostRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// host without the port
	hostOnly := strings.Split(r.Host, ":")[0]
	if router, ok := hr.routes[strings.ToLower(hostOnly)]; ok {
		router.ServeHTTP(w, r)
		return
	}
	// wildcard
	if router, ok := hr.routes["*"]; ok {
		router.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

// The HostRouter does not use chi.Route as mentioned in Match, and so
// the list of routes is always empty
func (hr *HostRouter) Routes() []chi.Route {
	return []chi.Route{}
}

// The HostRouter does not support setting up middleware for all
// hosts, it should be done on a per host basis using the chi.Router
// added via Map. Thus the host router always has empty middleware
func (hr *HostRouter) Middlewares() chi.Middlewares {
	return chi.Middlewares{}
}

func newBaseRouter(appCtx appcontext.AppContext, routingConfig *Config, telemetryConfig *telemetry.Config, serverName string) chi.Router {
	router := chi.NewRouter()
	// Add middleware: they are evaluated in the reverse order in which they
	// are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	router.Use(telemetry.NewOtelHTTPMiddleware(telemetryConfig, serverName, appCtx.Logger()))
	router.Use(auth.SessionIDMiddleware(routingConfig.HandlerConfig.AppNames(), routingConfig.HandlerConfig.SessionManagers()))
	router.Use(middleware.Trace(telemetryConfig)) // injects trace id into the context
	router.Use(middleware.ContextLogger("milmove_trace_id", appCtx.Logger()))
	router.Use(middleware.Recovery(appCtx.Logger()))
	router.Use(middleware.SecurityHeaders(appCtx.Logger()))

	if routingConfig.MaxBodySize > 0 {
		router.Use(middleware.LimitBodySize(routingConfig.MaxBodySize, appCtx.Logger()))
	}

	return router
}

func mountHealthRoute(appCtx appcontext.AppContext, redisPool *redis.Pool,
	routingConfig *Config, site chi.Router) {
	requestLoggerMiddleware := middleware.RequestLogger()
	healthHandler := handlers.NewHealthHandler(appCtx,
		redisPool, routingConfig.GitBranch, routingConfig.GitCommit)
	site.Method("GET", "/health", requestLoggerMiddleware(healthHandler))
}

func mountLocalStorageRoute(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.LocalStorageRoot != "" && routingConfig.LocalStorageWebRoot != "" {
		localStorageHandlerFunc := storage.NewFilesystemHandler(
			routingConfig.FileSystem, routingConfig.LocalStorageRoot)

		// path.Join removes trailing slashes, but we want it
		storageHandlerPath := path.Join("/", routingConfig.LocalStorageWebRoot) + "/"
		appCtx.Logger().Info("Registering storage handler",
			zap.Any("storageHandlerPath", storageHandlerPath))
		site.Route(storageHandlerPath, func(r chi.Router) {
			r.Use(middleware.RequestLogger())
			r.Get("/*", localStorageHandlerFunc)
			r.Head("/*", localStorageHandlerFunc)
		})
	}
}

func mountStaticRoutes(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	// Serves files out of build folder
	cfs := handlers.NewCustomFileSystem(
		// use afero HttpFS so we can wrap the existing FileSystem
		// Super useful for testing
		afero.NewHttpFs(routingConfig.FileSystem).Dir(routingConfig.BuildRoot),
		"index.html",
		appCtx.Logger(),
	)

	clientHandler := handlers.NewSpaHandler(
		routingConfig.BuildRoot,
		"index.html",
		cfs,
	)

	site.Route("/static/", func(r chi.Router) {
		r.Use(middleware.RequestLogger())
		r.Method("GET", "/*", clientHandler)
		r.Method("HEAD", "/*", clientHandler)
	})

	site.Route("/downloads/", func(r chi.Router) {
		r.Use(middleware.RequestLogger())
		r.Method("GET", "/*", clientHandler)
		r.Method("HEAD", "/*", clientHandler)
	})

	site.Method("GET", "/favicon.ico", clientHandler)

	// Explicitly disable swagger.json route
	site.Method("GET", "/swagger.json", http.NotFoundHandler())

	// Not sure this is right. We should serve swagger-ui per
	// swagger endpoint, not globally - ahobson 2023-07-06
	// if routingConfig.ServeSwaggerUI {
	// 	appCtx.Logger().Info("Swagger UI static file serving is enabled")
	// 	site.Method("GET", "/swagger-ui/", clientHandler)
	// } else {
	// 	site.Method("GET", "/swagger-ui/", http.NotFoundHandler())
	// }
}

func mountCollectorRoutes(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.ServeClientCollector {
		appCtx.Logger().Info("client collecting service is enabled")
		clientLogHandler := handlers.NewClientLogHandler(appCtx)
		site.Route("/client", func(r chi.Router) {
			r.Use(middleware.RequestLogger())
			r.Post("/log", clientLogHandler)
		})
	}
}

func mountPrimeAPI(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.ServePrime {
		clientCertMiddleware := authentication.ClientCertMiddleware(appCtx)

		// Setup shared middleware
		site.Route("/prime", func(primeRouter chi.Router) {
			if routingConfig.ServeDevlocalAuth {
				devlocalClientCertMiddleware := authentication.DevlocalClientCertMiddleware(appCtx)
				primeRouter.Use(devlocalClientCertMiddleware)
			} else {
				primeRouter.Use(clientCertMiddleware)
			}
			primeRouter.Use(authentication.PrimeAuthorizationMiddleware(appCtx.Logger()))
			primeRouter.Use(middleware.NoCache())
			primeRouter.Use(middleware.RequestLogger())

			// Setup version specific info for v1
			primeRouter.Route("/v1", func(r chi.Router) {
				r.Method("GET", "/swagger.yaml",
					handlers.NewFileHandler(routingConfig.FileSystem,
						routingConfig.PrimeSwaggerPath))
				if routingConfig.ServeSwaggerUI {
					appCtx.Logger().Info("Prime API Swagger UI serving is enabled")
					r.Method("GET", "/docs",
						handlers.NewFileHandler(routingConfig.FileSystem,
							path.Join(routingConfig.BuildRoot, "swagger-ui", "prime.html")))
				} else {
					r.Method("GET", "/docs", http.NotFoundHandler())
				}
				api := primeapi.NewPrimeAPI(routingConfig.HandlerConfig)
				tracingMiddleware := middleware.OpenAPITracing(api)
				r.Mount("/", api.Serve(tracingMiddleware))
			})
			// Setup version specific info for v2
			primeRouter.Route("/v2", func(r chi.Router) {
				r.Method("GET", "/swagger.yaml",
					handlers.NewFileHandler(routingConfig.FileSystem,
						routingConfig.PrimeV2SwaggerPath))
				if routingConfig.ServeSwaggerUI {
					r.Method("GET", "/docs",
						handlers.NewFileHandler(routingConfig.FileSystem,
							path.Join(routingConfig.BuildRoot, "swagger-ui", "prime_v2.html")))
				} else {
					r.Method("GET", "/docs", http.NotFoundHandler())
				}
				api := primeapiv2.NewPrimeAPI(routingConfig.HandlerConfig)
				tracingMiddleware := middleware.OpenAPITracing(api)
				r.Mount("/", api.Serve(tracingMiddleware))
			})
			// Setup version specific info for v3
			primeRouter.Route("/v3", func(r chi.Router) {
				r.Method("GET", "/swagger.yaml",
					handlers.NewFileHandler(routingConfig.FileSystem,
						routingConfig.PrimeV3SwaggerPath))
				if routingConfig.ServeSwaggerUI {
					r.Method("GET", "/docs",
						handlers.NewFileHandler(routingConfig.FileSystem,
							path.Join(routingConfig.BuildRoot, "swagger-ui", "prime_v3.html")))
				} else {
					r.Method("GET", "/docs", http.NotFoundHandler())
				}
				api := primeapiv3.NewPrimeAPI(routingConfig.HandlerConfig)
				tracingMiddleware := middleware.OpenAPITracing(api)
				r.Mount("/", api.Serve(tracingMiddleware))
			})
		})
	}
}

// PPTAS API to serve under the mTLS "Api" / "Prime" API server to support Navy requests
func mountPPTASAPI(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.ServePPTAS {
		clientCertMiddleware := authentication.ClientCertMiddleware(appCtx)
		site.Route("/pptas/v1", func(r chi.Router) {
			if routingConfig.ServeDevlocalAuth {
				devlocalClientCertMiddleware := authentication.DevlocalClientCertMiddleware(appCtx)
				r.Use(devlocalClientCertMiddleware)
			} else {
				r.Use(clientCertMiddleware)
			}
			r.Use(authentication.PPTASAuthorizationMiddleware(appCtx.Logger()))
			r.Use(middleware.NoCache())
			r.Use(middleware.RequestLogger())
			r.Method(
				"GET",
				"/swagger.yaml",
				handlers.NewFileHandler(routingConfig.FileSystem,
					routingConfig.PPTASSwaggerPath))
			if routingConfig.ServeSwaggerUI {
				r.Method("GET", "/docs",
					handlers.NewFileHandler(routingConfig.FileSystem,
						path.Join(routingConfig.BuildRoot, "swagger-ui", "pptas.html")))
			} else {
				r.Method("GET", "/docs", http.NotFoundHandler())
			}
			api := pptasapi.NewPPTASAPI(routingConfig.HandlerConfig)
			tracingMiddleware := middleware.OpenAPITracing(api)
			r.Mount("/", api.Serve(tracingMiddleware))
		})
	}
}

// Remember that the support api is to assist inside of dev/stg for endpoints such as
// manually invoking the EDI858 generator for a given payment request. It should never
// be utilized in production
func mountSupportAPI(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.ServeSupport {
		clientCertMiddleware := authentication.ClientCertMiddleware(appCtx)

		site.Route("/support/v1/", func(r chi.Router) {
			r.Use(clientCertMiddleware)
			r.Use(authentication.PrimeAuthorizationMiddleware(appCtx.Logger()))
			r.Use(middleware.NoCache())
			r.Use(middleware.RequestLogger())
			r.Method("GET", "/swagger.yaml",
				handlers.NewFileHandler(routingConfig.FileSystem,
					routingConfig.SupportSwaggerPath))
			if routingConfig.ServeSwaggerUI {
				appCtx.Logger().Info("Support API Swagger UI serving is enabled")
				r.Method("GET", "/docs",
					handlers.NewFileHandler(routingConfig.FileSystem,
						path.Join(routingConfig.BuildRoot, "swagger-ui", "support.html")))
			} else {
				r.Method("GET", "/docs", http.NotFoundHandler())
			}
			r.Mount("/", supportapi.NewSupportAPIHandler(routingConfig.HandlerConfig))
		})
	}
}

func mountTestharnessAPI(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	// only enable the test harness if support and devlocal auth
	// is enabled, and do it before CSRF and other middleware
	if routingConfig.ServeDevlocalAuth {
		addAuditUserToRequestContextMiddleware := authentication.AddAuditUserIDToRequestContextMiddleware(appCtx)
		appCtx.Logger().Info("Enabling testharness")
		site.Route("/testharness", func(r chi.Router) {
			r.Use(middleware.RequestLogger())
			r.Use(addAuditUserToRequestContextMiddleware)
			r.Method("POST", "/build/{action}",
				testharnessapi.NewDefaultBuilder(routingConfig.HandlerConfig))
			r.Method("GET", "/list",
				testharnessapi.NewBuilderList(routingConfig.HandlerConfig))

		})
	}
}

func mountDebugRoutes(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.ServeDebugPProf {
		appCtx.Logger().Info("Enabling pprof routes")

		site.Route("/debug/pprof/", func(r chi.Router) {
			r.Use(middleware.RequestLogger())
			r.Get("/", pprof.Index)
			r.Method("GET", "/allocs", pprof.Handler("allocs"))
			r.Method("GET", "/block", pprof.Handler("block"))
			r.Get("/cmdline", pprof.Cmdline)
			r.Method("GET", "/goroutine", pprof.Handler("goroutine"))
			r.Method("GET", "/heap", pprof.Handler("heap"))
			r.Method("GET", "/mutex", pprof.Handler("mutex"))
			r.Get("/profile", pprof.Profile)
			r.Get("/trace", pprof.Trace)
			r.Method("GET", "/threadcreate", pprof.Handler("threadcreate"))
			r.Get("/symbol", pprof.Symbol)
		})
	}
}

func mountInternalAPI(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.ServeAPIInternal {
		isLoggedInMiddleware := authentication.IsLoggedInMiddleware(appCtx.Logger())
		userAuthMiddleware := authentication.UserAuthMiddleware(appCtx.Logger())
		addAuditUserToRequestContextMiddleware := authentication.AddAuditUserIDToRequestContextMiddleware(appCtx)
		site.Route("/internal", func(r chi.Router) {
			r.Method("GET", "/swagger.yaml",
				handlers.NewFileHandler(routingConfig.FileSystem,
					routingConfig.APIInternalSwaggerPath))
			if routingConfig.ServeSwaggerUI {
				appCtx.Logger().Info("Internal API Swagger UI serving is enabled")
				r.Method("GET", "/docs",
					handlers.NewFileHandler(routingConfig.FileSystem,
						path.Join(routingConfig.BuildRoot, "swagger-ui", "internal.html")))
			} else {
				r.Method("GET", "/docs", http.NotFoundHandler())
			}
			r.Method("GET", "/users/is_logged_in", isLoggedInMiddleware)
			// Mux for internal API that enforces auth
			r.Route("/", func(rAuth chi.Router) {
				rAuth.Use(userAuthMiddleware)
				rAuth.Use(addAuditUserToRequestContextMiddleware)
				rAuth.Use(middleware.NoCache())
				api := internalapi.NewInternalAPI(routingConfig.HandlerConfig)
				// This middleware enables stricter checks for most of the internal api endpoints
				customerAPIAuthMiddleware := authentication.CustomerAPIAuthMiddleware(appCtx, api)
				rAuth.Use(customerAPIAuthMiddleware)
				tracingMiddleware := middleware.OpenAPITracing(api)
				rAuth.Mount("/", api.Serve(tracingMiddleware))
			})
		})
	}
}

func mountAdminAPI(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.ServeAdmin {
		userAuthMiddleware := authentication.UserAuthMiddleware(appCtx.Logger())
		addAuditUserToRequestContextMiddleware := authentication.AddAuditUserIDToRequestContextMiddleware(appCtx)
		site.Route("/admin/v1", func(r chi.Router) {
			r.Method("GET", "/swagger.yaml",
				handlers.NewFileHandler(routingConfig.FileSystem,
					routingConfig.AdminSwaggerPath))
			if routingConfig.ServeSwaggerUI {
				appCtx.Logger().Info("Admin API Swagger UI serving is enabled")
				r.Method("GET", "/docs",
					handlers.NewFileHandler(routingConfig.FileSystem,
						path.Join(routingConfig.BuildRoot, "swagger-ui", "admin.html")))
			} else {
				r.Method("GET", "/docs", http.NotFoundHandler())
			}

			// Mux for admin API that enforces auth
			r.Route("/", func(rAuth chi.Router) {
				rAuth.Use(userAuthMiddleware)
				rAuth.Use(addAuditUserToRequestContextMiddleware)
				rAuth.Use(authentication.AdminAuthMiddleware(appCtx.Logger()))
				rAuth.Use(middleware.NoCache())
				api := adminapi.NewAdminAPI(routingConfig.HandlerConfig)
				tracingMiddleware := middleware.OpenAPITracing(api)
				rAuth.Mount("/", api.Serve(tracingMiddleware))
			})
		})
	}
}

func mountPrimeSimulatorAPI(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.ServePrimeSimulator {
		userAuthMiddleware := authentication.UserAuthMiddleware(appCtx.Logger())
		addAuditUserToRequestContextMiddleware := authentication.AddAuditUserIDToRequestContextMiddleware(appCtx)
		site.Route("/prime/v1", func(r chi.Router) {
			r.Method("GET", "/swagger.yaml",
				handlers.NewFileHandler(routingConfig.FileSystem,
					routingConfig.PrimeSwaggerPath))
			if routingConfig.ServeSwaggerUI {
				appCtx.Logger().Info("Prime Simulator API Swagger UI serving is enabled")
				r.Method("GET", "/docs",
					handlers.NewFileHandler(routingConfig.FileSystem,
						path.Join(routingConfig.BuildRoot, "swagger-ui", "prime.html")))
			} else {
				r.Method("GET", "/docs", http.NotFoundHandler())
			}

			// Mux for prime simulator API that enforces auth
			r.Route("/", func(rAuth chi.Router) {
				rAuth.Use(userAuthMiddleware)
				rAuth.Use(addAuditUserToRequestContextMiddleware)
				rAuth.Use(authentication.PrimeSimulatorAuthorizationMiddleware(appCtx.Logger()))
				rAuth.Use(middleware.NoCache())
				api := primeapi.NewPrimeAPI(routingConfig.HandlerConfig)
				tracingMiddleware := middleware.OpenAPITracing(api)
				rAuth.Mount("/", api.Serve(tracingMiddleware))
			})
		})
		site.Route("/prime/v2", func(r chi.Router) {
			r.Method("GET", "/swagger.yaml",
				handlers.NewFileHandler(routingConfig.FileSystem,
					routingConfig.PrimeV2SwaggerPath))
			if routingConfig.ServeSwaggerUI {
				appCtx.Logger().Info("Prime Simulator API Swagger UI serving is enabled")
				r.Method("GET", "/docs",
					handlers.NewFileHandler(routingConfig.FileSystem,
						path.Join(routingConfig.BuildRoot, "swagger-ui", "prime.html")))
			} else {
				r.Method("GET", "/docs", http.NotFoundHandler())
			}

			// Mux for prime simulator API that enforces auth
			r.Route("/", func(rAuth chi.Router) {
				rAuth.Use(userAuthMiddleware)
				rAuth.Use(addAuditUserToRequestContextMiddleware)
				rAuth.Use(authentication.PrimeSimulatorAuthorizationMiddleware(appCtx.Logger()))
				rAuth.Use(middleware.NoCache())
				api := primeapiv2.NewPrimeAPI(routingConfig.HandlerConfig)
				tracingMiddleware := middleware.OpenAPITracing(api)
				rAuth.Mount("/", api.Serve(tracingMiddleware))
			})
		})
		site.Route("/prime/v3", func(r chi.Router) {
			r.Method("GET", "/swagger.yaml",
				handlers.NewFileHandler(routingConfig.FileSystem,
					routingConfig.PrimeV3SwaggerPath))
			if routingConfig.ServeSwaggerUI {
				appCtx.Logger().Info("Prime Simulator API Swagger UI serving is enabled")
				r.Method("GET", "/docs",
					handlers.NewFileHandler(routingConfig.FileSystem,
						path.Join(routingConfig.BuildRoot, "swagger-ui", "prime.html")))
			} else {
				r.Method("GET", "/docs", http.NotFoundHandler())
			}

			// Mux for prime simulator API that enforces auth
			r.Route("/", func(rAuth chi.Router) {
				rAuth.Use(userAuthMiddleware)
				rAuth.Use(addAuditUserToRequestContextMiddleware)
				rAuth.Use(authentication.PrimeSimulatorAuthorizationMiddleware(appCtx.Logger()))
				rAuth.Use(middleware.NoCache())
				api := primeapiv3.NewPrimeAPI(routingConfig.HandlerConfig)
				tracingMiddleware := middleware.OpenAPITracing(api)
				rAuth.Mount("/", api.Serve(tracingMiddleware))
			})
		})
		site.Route("/pptas/v1", func(r chi.Router) {
			r.Method("GET", "/swagger.yaml",
				handlers.NewFileHandler(routingConfig.FileSystem,
					routingConfig.PPTASSwaggerPath))
			if routingConfig.ServeSwaggerUI {
				appCtx.Logger().Info("PPTAS Simulator API Swagger UI serving is enabled")
				r.Method("GET", "/docs",
					handlers.NewFileHandler(routingConfig.FileSystem,
						path.Join(routingConfig.BuildRoot, "swagger-ui", "pptas.html")))
			} else {
				r.Method("GET", "/docs", http.NotFoundHandler())
			}

			// Mux for prime simulator API that enforces auth
			r.Route("/", func(rAuth chi.Router) {
				rAuth.Use(userAuthMiddleware)
				rAuth.Use(addAuditUserToRequestContextMiddleware)
				rAuth.Use(authentication.PrimeSimulatorAuthorizationMiddleware(appCtx.Logger()))
				rAuth.Use(middleware.NoCache())
				api := pptasapi.NewPPTASAPI(routingConfig.HandlerConfig)
				tracingMiddleware := middleware.OpenAPITracing(api)
				rAuth.Mount("/", api.Serve(tracingMiddleware))
			})
		})
		// Support API serves to support Prime API testing outside of production environments, hence why it is
		// mounted inside the Prime sim API without client cert middleware
		if routingConfig.ServeSupport {
			site.Route("/support/v1", func(r chi.Router) {
				r.Method("GET", "/swagger.yaml",
					handlers.NewFileHandler(routingConfig.FileSystem,
						routingConfig.SupportSwaggerPath))
				if routingConfig.ServeSwaggerUI {
					appCtx.Logger().Info("Support API Swagger UI serving is enabled")
					r.Method("GET", "/docs",
						handlers.NewFileHandler(routingConfig.FileSystem,
							path.Join(routingConfig.BuildRoot, "swagger-ui", "support.html")))
				} else {
					r.Method("GET", "/docs", http.NotFoundHandler())
				}

				// Mux for support API that enforces auth
				r.Route("/", func(rAuth chi.Router) {
					rAuth.Use(userAuthMiddleware)
					rAuth.Use(addAuditUserToRequestContextMiddleware)
					rAuth.Use(authentication.PrimeSimulatorAuthorizationMiddleware(appCtx.Logger()))
					rAuth.Use(middleware.NoCache())
					rAuth.Use(middleware.RequestLogger())
					rAuth.Mount("/", supportapi.NewSupportAPIHandler(routingConfig.HandlerConfig))
				})
			})
		}
	}
}

func mountGHCAPI(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	if routingConfig.ServeGHC {
		userAuthMiddleware := authentication.UserAuthMiddleware(appCtx.Logger())
		addAuditUserToRequestContextMiddleware := authentication.AddAuditUserIDToRequestContextMiddleware(appCtx)
		site.Route("/ghc/v1", func(r chi.Router) {
			r.Method("GET", "/swagger.yaml",
				handlers.NewFileHandler(routingConfig.FileSystem,
					routingConfig.GHCSwaggerPath))
			if routingConfig.ServeSwaggerUI {
				appCtx.Logger().Info("GHC API Swagger UI serving is enabled")
				r.Method("GET", "/docs",
					handlers.NewFileHandler(routingConfig.FileSystem,
						path.Join(routingConfig.BuildRoot, "swagger-ui", "ghc.html")))
			} else {
				r.Method("GET", "/docs", http.NotFoundHandler())
			}

			api := ghcapi.NewGhcAPIHandler(routingConfig.HandlerConfig)
			tracingMiddleware := middleware.OpenAPITracing(api)

			// Mux for GHC API open routes
			r.Route("/open", func(rOpen chi.Router) {
				rOpen.Mount("/", api.Serve(tracingMiddleware))
			})
			// Mux for GHC API that enforces auth
			r.Route("/", func(rAuth chi.Router) {
				rAuth.Use(userAuthMiddleware)
				rAuth.Use(addAuditUserToRequestContextMiddleware)
				rAuth.Use(middleware.NoCache())
				permissionsMiddleware := authentication.PermissionsMiddleware(appCtx, api)
				rAuth.Use(permissionsMiddleware)
				rAuth.Mount("/", api.Serve(tracingMiddleware))
			})
		})
	}
}

func mountAuthRoutes(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	site.Route("/auth/", func(r chi.Router) {
		r.Use(middleware.NoCache())
		r.Method("GET", "/okta", authentication.NewRedirectHandler(routingConfig.AuthContext, routingConfig.HandlerConfig, routingConfig.HandlerConfig.UseSecureCookie()))
		r.Method("GET", "/okta/callback", authentication.NewCallbackHandler(routingConfig.AuthContext, routingConfig.HandlerConfig, routingConfig.HandlerConfig.NotificationSender()))
		r.Method("POST", "/logout", authentication.NewLogoutHandler(routingConfig.AuthContext, routingConfig.HandlerConfig))
		r.Method("POST", "/logoutOktaRedirect", authentication.NewLogoutOktaRedirectHandler(routingConfig.AuthContext, routingConfig.HandlerConfig))
	})

	if routingConfig.ServeDevlocalAuth {
		appCtx.Logger().Info("Enabling devlocal auth")
		site.Route("/devlocal-auth/", func(r chi.Router) {
			r.Use(middleware.NoCache())
			r.Method("GET", "/login", authentication.NewUserListHandler(routingConfig.AuthContext, routingConfig.HandlerConfig))
			r.Method("POST", "/login", authentication.NewAssignUserHandler(routingConfig.AuthContext, routingConfig.HandlerConfig))
			r.Method("POST", "/new", authentication.NewCreateAndLoginUserHandler(routingConfig.AuthContext, routingConfig.HandlerConfig))
			r.Method("POST", "/create", authentication.NewCreateUserHandler(routingConfig.AuthContext, routingConfig.HandlerConfig))
		})
	}
}

func mountDefaultStaticRoute(appCtx appcontext.AppContext, routingConfig *Config, site chi.Router) {
	// Serve index.html to all requests that haven't matches a
	// previous route,

	defaultHandler := indexHandler(routingConfig, appCtx.Logger())
	site.NotFound(defaultHandler)
}

func mountSessionRoutes(appCtx appcontext.AppContext, routingConfig *Config, baseSite chi.Router, sessionManager auth.SessionManager, fn func(r chi.Router)) {
	sessionCookieMiddleware := auth.SessionCookieMiddleware(appCtx.Logger(), routingConfig.HandlerConfig.AppNames(), routingConfig.HandlerConfig.SessionManagers())
	maskedCSRFMiddleware := auth.MaskedCSRFMiddleware(routingConfig.HandlerConfig.UseSecureCookie())

	baseSite.Route("/", func(r chi.Router) {
		// need to load and save in the session manager before any other
		r.Use(sessionManager.LoadAndSave)
		r.Use(sessionCookieMiddleware)
		r.Use(middleware.RequestLogger())

		appCtx.Logger().Info("Enabling CSRF protection")
		r.Use(routingConfig.CSRFMiddleware)
		r.Use(maskedCSRFMiddleware)

		mountAuthRoutes(appCtx, routingConfig, r)
		// invoke callback to mount session routes
		fn(r)
	})

}

func newMilRouter(appCtx appcontext.AppContext, redisPool *redis.Pool,
	routingConfig *Config, telemetryConfig *telemetry.Config, serverName string) chi.Router {

	site := newBaseRouter(appCtx, routingConfig, telemetryConfig, serverName)

	mountHealthRoute(appCtx, redisPool, routingConfig, site)
	mountLocalStorageRoute(appCtx, routingConfig, site)
	mountStaticRoutes(appCtx, routingConfig, site)
	mountTestharnessAPI(appCtx, routingConfig, site)
	mountCollectorRoutes(appCtx, routingConfig, site)

	milSessionManager := routingConfig.HandlerConfig.SessionManagers().Mil
	mountSessionRoutes(appCtx, routingConfig, site, milSessionManager,
		func(sessionRoute chi.Router) {
			mountInternalAPI(appCtx, routingConfig, sessionRoute)
		},
	)

	mountDefaultStaticRoute(appCtx, routingConfig, site)

	return site
}

func newOfficeRouter(appCtx appcontext.AppContext, redisPool *redis.Pool,
	routingConfig *Config, telemetryConfig *telemetry.Config, serverName string) chi.Router {

	site := newBaseRouter(appCtx, routingConfig, telemetryConfig, serverName)

	mountHealthRoute(appCtx, redisPool, routingConfig, site)
	mountLocalStorageRoute(appCtx, routingConfig, site)
	mountStaticRoutes(appCtx, routingConfig, site)
	mountTestharnessAPI(appCtx, routingConfig, site)
	mountCollectorRoutes(appCtx, routingConfig, site)

	officeSessionManager := routingConfig.HandlerConfig.SessionManagers().Office
	mountSessionRoutes(appCtx, routingConfig, site, officeSessionManager,
		func(sessionRoute chi.Router) {
			mountInternalAPI(appCtx, routingConfig, sessionRoute)
			mountPrimeSimulatorAPI(appCtx, routingConfig, sessionRoute)
			mountGHCAPI(appCtx, routingConfig, sessionRoute)
		},
	)

	mountDefaultStaticRoute(appCtx, routingConfig, site)

	return site
}

func newAdminRouter(appCtx appcontext.AppContext, redisPool *redis.Pool,
	routingConfig *Config, telemetryConfig *telemetry.Config, serverName string) chi.Router {

	site := newBaseRouter(appCtx, routingConfig, telemetryConfig, serverName)

	mountHealthRoute(appCtx, redisPool, routingConfig, site)
	mountStaticRoutes(appCtx, routingConfig, site)
	mountTestharnessAPI(appCtx, routingConfig, site)
	mountCollectorRoutes(appCtx, routingConfig, site)

	adminSessionManager := routingConfig.HandlerConfig.SessionManagers().Admin
	mountSessionRoutes(appCtx, routingConfig, site, adminSessionManager,
		func(sessionRoute chi.Router) {
			mountInternalAPI(appCtx, routingConfig, sessionRoute)
			mountAdminAPI(appCtx, routingConfig, sessionRoute)
			// debug routes can go anywhere, but the admin site seems most appropriate
			mountDebugRoutes(appCtx, routingConfig, sessionRoute)
		},
	)

	mountDefaultStaticRoute(appCtx, routingConfig, site)

	return site
}

// This "Prime" router is really just the "API" router for MilMove.
// It was initially just named under "Prime" as it was the only use of the router
func newPrimeRouter(appCtx appcontext.AppContext, redisPool *redis.Pool,
	routingConfig *Config, telemetryConfig *telemetry.Config, serverName string) chi.Router {

	site := newBaseRouter(appCtx, routingConfig, telemetryConfig, serverName)

	mountHealthRoute(appCtx, redisPool, routingConfig, site)
	mountPrimeAPI(appCtx, routingConfig, site)
	mountSupportAPI(appCtx, routingConfig, site)
	mountTestharnessAPI(appCtx, routingConfig, site)
	mountPPTASAPI(appCtx, routingConfig, site)
	return site
}

// InitRouting sets up the routing for all the hosts (mil, office,
// admin, api)
func InitRouting(serverName string, appCtx appcontext.AppContext, redisPool *redis.Pool,
	routingConfig *Config, telemetryConfig *telemetry.Config) (http.Handler, error) {

	// check for missing CSRF middleware ASAP
	if routingConfig.CSRFMiddleware == nil {
		return nil, errors.New("missing CSRF Middleware")
	}

	// With chi, we have to register all middleware before setting up
	// routes
	//
	// Because we want different middleware for different hosts, we
	// need to set up a router per host

	hostRouter := NewHostRouter()

	milServerName := routingConfig.HandlerConfig.AppNames().MilServername
	milRouter := newMilRouter(appCtx, redisPool, routingConfig, telemetryConfig, serverName)
	hostRouter.Map(milServerName, milRouter)

	officeServerName := routingConfig.HandlerConfig.AppNames().OfficeServername
	officeRouter := newOfficeRouter(appCtx, redisPool, routingConfig, telemetryConfig, serverName)
	hostRouter.Map(officeServerName, officeRouter)

	adminServerName := routingConfig.HandlerConfig.AppNames().AdminServername
	adminRouter := newAdminRouter(appCtx, redisPool, routingConfig, telemetryConfig, serverName)
	hostRouter.Map(adminServerName, adminRouter)

	primeServerName := routingConfig.HandlerConfig.AppNames().PrimeServername
	primeRouter := newPrimeRouter(appCtx, redisPool, routingConfig, telemetryConfig, serverName)
	hostRouter.Map(primeServerName, primeRouter)

	// need a wildcard health router as the ELB makes requests to the
	// IP, not the hostname
	healthRouter := chi.NewRouter()
	mountHealthRoute(appCtx, redisPool, routingConfig, healthRouter)
	hostRouter.Map("*", healthRouter)

	return hostRouter, nil
}

// InitHealthRouting sets up the routing for the internal health
// server used by the ECS health check
func InitHealthRouting(serverName string, appCtx appcontext.AppContext, redisPool *redis.Pool,
	routingConfig *Config, telemetryConfig *telemetry.Config) (http.Handler, error) {

	site := chi.NewRouter()
	site.Use(telemetry.NewOtelHTTPMiddleware(telemetryConfig, serverName, appCtx.Logger()))
	mountHealthRoute(appCtx, redisPool, routingConfig, site)

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
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			http.ServeContent(w, r, "index.html", stat.ModTime(), reader)
		} else {
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}
}
