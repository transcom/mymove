package auth

import (
	"context"
	"github.com/gobuffalo/uuid"
	"net/http"
)

type authSessionKey string

const sessionContextKey authSessionKey = "session"

// Session stores information about the currently logged in session
// EntityID is the associated entity for a user
//   - office_user -> transportation_office_id
//   - tsp_user -> transportation_service_provider_id
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
	EntityID        uuid.UUID
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
