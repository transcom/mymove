package auth

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models/roles"
)

type authSessionKey string

const sessionContextKey authSessionKey = "session"

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

// SetupSessionManagers configures the session manager for each app: mil, admin,
// and office. It's necessary to have separate session managers to allow users
// to be signed in on multiple apps at the same time.
func SetupSessionManagers(redisEnabled bool, sessionStore scs.Store, useSecureCookie bool, idleTimeout time.Duration, lifetime time.Duration) [3]*scs.SessionManager {
	if !redisEnabled {
		return [3]*scs.SessionManager{}
	}
	var milSession, adminSession, officeSession *scs.SessionManager
	gob.Register(Session{})

	milSession = scs.New()
	milSession.Store = sessionStore
	milSession.Cookie.Name = "mil_session_token"

	adminSession = scs.New()
	adminSession.Store = sessionStore
	adminSession.Cookie.Name = "admin_session_token"

	officeSession = scs.New()
	officeSession.Store = sessionStore
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

	return [3]*scs.SessionManager{milSession, adminSession, officeSession}
}

// IsOfficeApp returns true iff the request is for the office.move.mil host
func (s *Session) IsOfficeApp() bool {
	return s.ApplicationName == OfficeApp
}

// IsMilApp returns true iff the request is for the my.move.mil host
func (s *Session) IsMilApp() bool {
	return s.ApplicationName == MilApp
}

// IsAdminApp returns true iff the request is for the admin.move.mil host
func (s *Session) IsAdminApp() bool {
	return s.ApplicationName == AdminApp
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
	return context.WithValue(r.Context(), sessionContextKey, session)
}

// SetSessionInContext modifies the context to add the session data.
func SetSessionInContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
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
