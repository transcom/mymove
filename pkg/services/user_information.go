package services

import (
	"github.com/gofrs/uuid"
)

// UserInformation defines the user fields that are displayed
// when looking up a particular user on the admin site
type UserInformation struct {
	UserID                 uuid.UUID `db:"user_id"`
	LoginGovEmail          *string   `db:"login_gov_email"`
	CurrentAdminSessionID  *string   `db:"current_admin_session_id"`
	CurrentOfficeSessionID *string   `db:"current_office_session_id"`
	CurrentMilSessionID    *string   `db:"current_mil_session_id"`
}

// UserInformationFetcher is the service object interface for FetchUserInformation
//go:generate mockery -name UserInformationFetcher
type UserInformationFetcher interface {
	FetchUserInformation(uuid uuid.UUID) (UserInformation, error)
}
