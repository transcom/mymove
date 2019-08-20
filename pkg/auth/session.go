package auth

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
)

type authSessionKey string

const sessionContextKey authSessionKey = "session"

// Application describes the application name
type Application string

const (
	// TspApp indicates tsp.move.mil
	TspApp Application = "tsp"
	// OfficeApp indicates office.move.mil
	OfficeApp Application = "office"
	// MilApp indicates my.move.mil (DNS still points to my.move.mil and not mil.move.mil)
	MilApp Application = "mil"
	// AdminApp indicates admin.move.mil
	AdminApp Application = "admin"
)

// IsTspApp returns true iff the request is for the tsp.move.mil host
func (s *Session) IsTspApp() bool {
	return s.ApplicationName == TspApp
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
	TspUserID       uuid.UUID
	AdminUserID     uuid.UUID
	AdminUserRole   string
	DpsUserID       uuid.UUID
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

// IsTspUser checks whether the authenticated user is a TspUser
func (s *Session) IsTspUser() bool {
	return s.TspUserID != uuid.Nil
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

// IsDpsUser checks whether the authenticated user is a DpsUser
func (s *Session) IsDpsUser() bool {
	return s.DpsUserID != uuid.Nil
}
