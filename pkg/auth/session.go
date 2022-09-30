package auth

import (
	"context"
	"encoding/gob"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gofrs/uuid"
	"github.com/gomodule/redigo/redis"

	"github.com/transcom/mymove/pkg/models/roles"
)

type authSessionKey string

const sessionContextKey authSessionKey = "session"

type sessionIDKey string

const sessionIDContextKey sessionIDKey = "sessionID"

// Application describes the application name
type Application string

const (
	// OfficeApp indicates office.move.mil
	OfficeApp Application = "office"
	// MilApp indicates my.move.mil (DNS still points to my.move.mil and not mil.move.mil)
	MilApp Application = "mil"
	// AdminApp indicates admin.move.mil
	AdminApp Application = "admin"
)

// IsOfficeApp returns true if the application is the office app
func (a Application) IsOfficeApp() bool {
	return a == OfficeApp
}

// IsMilApp returns true if the application is the mil app
func (a Application) IsMilApp() bool {
	return a == MilApp
}

// IsAdminApp returns true if the application is the admin app
func (a Application) IsAdminApp() bool {
	return a == AdminApp
}

type SessionManager interface {
	Get(context.Context, string) interface{}
	Put(context.Context, string, interface{})
	Destroy(context.Context) error
	RenewToken(context.Context) error
	Commit(context.Context) (string, time.Time, error)
	Load(context.Context, string) (context.Context, error)
	LoadAndSave(http.Handler) http.Handler
	Store() scs.Store
}

type ScsSessionManagerWrapper struct {
	ScsSessionManager *scs.SessionManager
}

func (s ScsSessionManagerWrapper) Get(ctx context.Context, key string) interface{} {
	return s.ScsSessionManager.Get(ctx, key)
}

func (s ScsSessionManagerWrapper) Put(ctx context.Context, key string, val interface{}) {
	s.ScsSessionManager.Put(ctx, key, val)
}

func (s ScsSessionManagerWrapper) Destroy(ctx context.Context) error {
	return s.ScsSessionManager.Destroy(ctx)
}

func (s ScsSessionManagerWrapper) RenewToken(ctx context.Context) error {
	return s.ScsSessionManager.RenewToken(ctx)
}

func (s ScsSessionManagerWrapper) Commit(ctx context.Context) (string, time.Time, error) {
	return s.ScsSessionManager.Commit(ctx)
}

func (s ScsSessionManagerWrapper) Load(ctx context.Context, token string) (context.Context, error) {
	return s.ScsSessionManager.Load(ctx, token)
}

func (s ScsSessionManagerWrapper) LoadAndSave(next http.Handler) http.Handler {
	return s.ScsSessionManager.LoadAndSave(next)
}

func (s ScsSessionManagerWrapper) Store() scs.Store {
	return s.ScsSessionManager.Store
}

type AppSessionManagers struct {
	Mil    SessionManager
	Office SessionManager
	Admin  SessionManager
}

// sessionManagerForSession returns the appropriate session manager
// for the session
func (a AppSessionManagers) SessionManagerForApplication(app Application) SessionManager {
	if app.IsMilApp() {
		return a.Mil
	} else if app.IsAdminApp() {
		return a.Admin
	} else if app.IsOfficeApp() {
		return a.Office
	}

	return nil
}

// SetupSessionManagers configures the session manager for each app: mil, admin,
// and office. It's necessary to have separate session managers to allow users
// to be signed in on multiple apps at the same time.
func SetupSessionManagers(redisPool *redis.Pool, useSecureCookie bool, idleTimeout time.Duration, lifetime time.Duration) AppSessionManagers {
	var milSession, adminSession, officeSession *scs.SessionManager
	gob.Register(Session{})

	// we need to ensure each session manager has its own store so
	// that sessions don't leak between apps. If a redisPool is
	// provided, we can use a prefix to ensure sessions are separated.
	// If redis is not configured, this would be local testing of some
	// kind in which case we can create a separate memory store for
	// each app
	newSessionStoreFn := func(prefix string) scs.Store {
		if redisPool != nil {
			return redisstore.NewWithPrefix(redisPool, prefix)
		}
		// For local testing, we don't need a background thread
		// cleaning up session
		return memstore.NewWithCleanupInterval(time.Duration(0))
	}

	milSession = scs.New()
	milSession.Store = newSessionStoreFn("mil")
	milSession.Cookie.Name = "mil_session_token"

	adminSession = scs.New()
	adminSession.Store = newSessionStoreFn("admin")
	adminSession.Cookie.Name = "admin_session_token"

	officeSession = scs.New()
	officeSession.Store = newSessionStoreFn("office")
	officeSession.Cookie.Name = "office_session_token"

	// IdleTimeout controls the maximum length of time a session can be inactive
	// before it expires. The default is 15 minutes. To disable idle timeout in
	// a non-production environment, set SESSION_IDLE_TIMEOUT_IN_MINUTES to 0.
	milSession.IdleTimeout = idleTimeout
	adminSession.IdleTimeout = idleTimeout
	officeSession.IdleTimeout = idleTimeout

	// Lifetime controls the maximum length of time that a session is valid for
	// before it expires. The lifetime is an 'absolute expiry' which is set when
	// the session is first created or renewed (such as when a user signs in)
	// and does not change. The default value is 24 hours.
	milSession.Lifetime = lifetime
	adminSession.Lifetime = lifetime
	officeSession.Lifetime = lifetime

	milSession.Cookie.Path = "/"
	adminSession.Cookie.Path = "/"
	officeSession.Cookie.Path = "/"

	// A value of false means the session cookie will be deleted when the
	// browser is closed.
	milSession.Cookie.Persist = false
	adminSession.Cookie.Persist = false
	officeSession.Cookie.Persist = false

	if useSecureCookie {
		milSession.Cookie.Secure = true
		adminSession.Cookie.Secure = true
		officeSession.Cookie.Secure = true
	}

	return AppSessionManagers{
		Mil: ScsSessionManagerWrapper{
			ScsSessionManager: milSession,
		},
		Office: ScsSessionManagerWrapper{
			ScsSessionManager: officeSession,
		},
		Admin: ScsSessionManagerWrapper{
			ScsSessionManager: adminSession,
		},
	}
}

// IsOfficeApp returns true iff the request is for the office.move.mil host
func (s *Session) IsOfficeApp() bool {
	return s.ApplicationName.IsOfficeApp()
}

// IsMilApp returns true iff the request is for the my.move.mil host
func (s *Session) IsMilApp() bool {
	return s.ApplicationName.IsMilApp()
}

// IsAdminApp returns true iff the request is for the admin.move.mil host
func (s *Session) IsAdminApp() bool {
	return s.ApplicationName.IsAdminApp()
}

// Session stores information about the currently logged in session
type Session struct {
	ApplicationName Application
	Hostname        string
	IDToken         string
	UserID          uuid.UUID
	Email           string
	FirstName       string
	Middle          string
	LastName        string
	ServiceMemberID uuid.UUID
	OfficeUserID    uuid.UUID
	AdminUserID     uuid.UUID
	AdminUserRole   string
	Roles           roles.Roles
	Permissions     []string
}

// SetSessionInRequestContext modifies the request's Context() to add the session data
func SetSessionInRequestContext(r *http.Request, session *Session) context.Context {
	return SetSessionInContext(r.Context(), session)
}

// SetSessionInContext modifies the context to add the session data.
func SetSessionInContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}

func setSessionIDInContext(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, sessionIDContextKey, sessionID)
}

func SessionIDFromContext(ctx context.Context) string {
	if sessionID, ok := ctx.Value(sessionIDContextKey).(string); ok {
		return sessionID
	}
	return ""
}

// SessionFromRequestContext gets the reference to the Session stored in the request.Context()
func SessionFromRequestContext(r *http.Request) *Session {
	return SessionFromContext(r.Context())
}

// SessionFromContext gets the reference to the Session stored in the request.Context()
func SessionFromContext(ctx context.Context) *Session {
	if session, ok := ctx.Value(sessionContextKey).(*Session); ok {
		return session
	}
	return nil
}

// IsServiceMember checks whether the authenticated user is a ServiceMember
func (s *Session) IsServiceMember() bool {
	return s.ServiceMemberID != uuid.Nil
}

// IsOfficeUser checks whether the authenticated user is an OfficeUser
func (s *Session) IsOfficeUser() bool {
	return s.OfficeUserID != uuid.Nil
}

// IsAdminUser checks whether the authenticated user is an AdminUser
func (s *Session) IsAdminUser() bool {
	return s.AdminUserID != uuid.Nil
}

// IsSystemAdmin checks whether the authenticated admin user is a system admin
func (s *Session) IsSystemAdmin() bool {
	role := "SYSTEM_ADMIN"
	return s.IsAdminUser() && s.AdminUserRole == role
}

// IsProgramAdmin checks whether the authenticated admin user is a program admin
func (s *Session) IsProgramAdmin() bool {
	role := "PROGRAM_ADMIN"
	return s.IsAdminUser() && s.AdminUserRole == role
}

func SessionIDMiddleware(appnames ApplicationServername, sessionManagers AppSessionManagers) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// Split the hostname from the port
			hostname := strings.Split(r.Host, ":")[0]
			app, err := ApplicationName(hostname, appnames)
			// the hostname may not be recognized if this is e.g.
			// a prime session. We won't have a sessionManager in
			// that case, so we won't add a sessionID to the context
			if err == nil {
				sessionManager := sessionManagers.SessionManagerForApplication(app)
				if sessionManager != nil {
					wrapper, ok := sessionManager.(ScsSessionManagerWrapper)
					if ok {
						cookie, err := r.Cookie(wrapper.ScsSessionManager.Cookie.Name)
						if err == nil {
							ctx = setSessionIDInContext(ctx, cookie.Value)
						}
					}
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
