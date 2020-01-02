package usersroles

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type usersRolesCreator struct {
	db *pop.Connection
}

// NewNewUsersRolesCreator creates a new struct with the service dependencies
func NewUsersRolesCreator(db *pop.Connection) services.UserRoleAssociator {
	return usersRolesCreator{db}
}

//AssociateUserRoles associates a given user with a set of roles
func (u usersRolesCreator) AssociateUserRoles(userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
	var usersRoles []models.UsersRoles
	user := models.User{}
	err := u.db.Find(&user, userID)
	if err != nil {
		return usersRoles, err
	}
	var roleIDs []uuid.UUID
	for _, roleType := range rs {
		var roleID uuid.UUID
		err := u.db.RawQuery("select id from roles where role_type = $1", roleType).First(&roleID)
		if err == nil {
			roleIDs = append(roleIDs, roleID)
		}
	}

	var allRoles []models.UsersRoles
	for _, r := range roleIDs {
		ur := models.UsersRoles{
			UserID: user.ID,
			RoleID: r,
		}
		allRoles = append(allRoles, ur)
	}
	err = u.db.Create(allRoles)
	if err != nil {
		return usersRoles, err
	}
	return allRoles, err
}
