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
)

// IsTspApp returns true iff the request is for the office.move.mil host
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

// Session stores information about the currently logged in session
type Session struct {
	ApplicationName Application
	Hostname        string
	IDToken         string
	Disabled        bool
	UserID          uuid.UUID
	Email           string
	FirstName       string
	Middle          string
	LastName        string
	ServiceMemberID uuid.UUID
	OfficeUserID    uuid.UUID
	TspUserID       uuid.UUID
	DpsUserID       uuid.UUID
}

// SetSessionInRequestContext modifies the request's Context() to add the session data
func SetSessionInRequestContext(r *http.Request, session *Session) context.Context {
	return context.WithValue(r.Context(), sessionContextKey, session)
}

// SessionFromRequestContext gets the reference to the Session stored in the request.Context()
func SessionFromRequestContext(r *http.Request) *Session {
	if session, ok := r.Context().Value(sessionContextKey).(*Session); ok {
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

// IsDpsUser checks whether the authenticated user is a DpsUser
func (s *Session) IsDpsUser() bool {
	return s.DpsUserID != uuid.Nil
}
