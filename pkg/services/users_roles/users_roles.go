package usersroles

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type usersRolesCreator struct {
	//logger Logger
	db *pop.Connection
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewUsersRolesCreator(db *pop.Connection) services.UserRoleAssociator {
	return usersRolesCreator{db}
}

func (u usersRolesCreator) AssociateUserRoles(userID uuid.UUID, rs roles.Roles) ([]models.UsersRoles, error) {
	var usersRoles []models.UsersRoles
	user := models.User{}
	err := u.db.Find(&user, userID)
	if err != nil {
		//logger.Error("Error saving user", zap.Error(err))
		return usersRoles, err
	}
	var allRoles []models.UsersRoles
	for _, r := range rs {
		ur := models.UsersRoles{
			UserID: user.ID,
			RoleID: r.ID,
		}
		allRoles = append(allRoles, ur)
	}
	err = u.db.Create(allRoles)
	if err != nil {
		//logger.Error("Error saving role", zap.Error(err))
		return usersRoles, err
	}
	return allRoles, err
}
