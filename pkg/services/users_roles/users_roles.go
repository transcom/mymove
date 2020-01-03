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
	usersRoles, err := u.fetchUnassociatedRoles(userID, rs)
	if err != nil {
		return usersRoles, err
	}
	err = u.db.Create(usersRoles)
	if err != nil {
		return usersRoles, err
	}
	return usersRoles, err
}

func (u usersRolesCreator) fetchUnassociatedRoles(userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
	// select all roles not already associated with this user
	var userRoles []models.UsersRoles
	rss := make([]interface{}, len(rs))
	for i := 1; i < len(rss); {
		rss[i] = rs[i]
		i++
	}
	err := u.db.RawQuery(
		`SELECT $1::uuid as user_id,
					  roles.id as role_id
	FROM roles
	WHERE role_type NOT IN (
		SELECT role_type
		FROM roles
				 JOIN users_roles ur ON roles.id = ur.role_id
		WHERE user_id = $1)`, userID).All(&userRoles)
	if err != nil {
		return userRoles, err
	}
	return userRoles, nil
}
