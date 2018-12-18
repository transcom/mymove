package auth

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
)

type authSessionKey string

const sessionContextKey authSessionKey = "session"

// Session stores information about the currently logged in session
type Session struct {
	ApplicationName application
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
	Features        []Feature
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

// CanAccessFeature checks whether or not the authenticated user can access
// a specific feature
func (s *Session) CanAccessFeature(feature Feature) bool {
	for _, f := range s.Features {
		if f == feature {
			return true
		}
	}
	return false
}
