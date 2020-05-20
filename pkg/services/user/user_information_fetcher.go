package user

import (
	"database/sql"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"
)

type userInformationFetcher struct {
	db *pop.Connection
}

// NewUserInformationFetcher return an implementation of the UserInformationFetcher interface
func NewUserInformationFetcher(db *pop.Connection) services.UserInformationFetcher {
	return &userInformationFetcher{db}
}

// FetchUserInformation fetches user information
func (uif *userInformationFetcher) FetchUserInformation(userID uuid.UUID) (services.UserInformation, error) {
	pop.Debug = true
	q := `
SELECT users.id as user_id,
			 users.login_gov_email,
			 users.current_admin_session_id,
			 users.current_office_session_id,
			 users.current_mil_session_id
FROM users
where users.id = $1`
	ui := services.UserInformation{}
	err := uif.db.RawQuery(q, userID).First(&ui)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return services.UserInformation{}, services.NewNotFoundError(userID, "")
		default:
			return services.UserInformation{}, err
		}
	}
	pop.Debug = false
	return ui, nil
}
