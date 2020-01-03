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
	roleMap, err := u.fetchUnassociatedRoles(userID)
	if err != nil {
		return usersRoles, err
	}
	for _, r := range rs {
		if rID, ok := roleMap[r]; ok {
			ur := models.UsersRoles{
				UserID: user.ID,
				RoleID: rID,
			}
			usersRoles = append(usersRoles, ur)
		}
	}
	err = u.db.Create(usersRoles)
	if err != nil {
		return usersRoles, err
	}
	return usersRoles, err
}

func (u usersRolesCreator) fetchUnassociatedRoles(userID uuid.UUID) (map[roles.RoleType]uuid.UUID, error) {
	var allRoles roles.Roles
	var roleMap = make(map[roles.RoleType]uuid.UUID)
	err := u.db.RawQuery(
		`SELECT *
	FROM roles
	WHERE role_type NOT IN (
    	SELECT role_type
    	FROM roles
        	JOIN users_roles ur ON roles.id = ur.role_id
    	WHERE user_id = $1);`, userID).
		All(&allRoles)
	if err != nil {
		return roleMap, err
	}
	for _, role := range allRoles {
		roleMap[role.RoleType] = role.ID
	}
	return roleMap, nil
}
