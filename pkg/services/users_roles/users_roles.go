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

//UpdateUserRoles associates a given user with a set of roles
func (u usersRolesCreator) UpdateUserRoles(userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
	_, err := u.addUserRoles(userID, rs)
	if err != nil {
		return []models.UsersRoles{}, err
	}
	_, err = u.removeUserRoles(userID, rs)
	if err != nil {
		return []models.UsersRoles{}, err
	}
	var usersRoles []models.UsersRoles
	// fetch + return updated roles
	err = u.db.Where("user_id = ?", userID).All(&usersRoles)
	if err != nil {
		return []models.UsersRoles{}, err
	}
	return usersRoles, nil
}

func (u usersRolesCreator) addUserRoles(userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
	//Having to use somewhat convoluted right join syntax b/c FROM clause in pop is derived from the model
	//and for the RawQuery was having trouble passing in array into the in clause with additional params
	//ideally would just be the query below
	//SELECT r.id                              AS role_id,
	//	'3b9360a3-3304-4c60-90f4-83d687884079' AS user_id,
	//	ur.user_id,
	//	role_type
	//FROM roles r
	//		LEFT JOIN users_roles ur ON r.id = ur.role_id
	//	AND ur.user_id = '3b9360a3-3304-4c60-90f4-83d687884079'
	//WHERE role_type IN ('transportation_ordering_officer', 'contracting_officer', 'customer')
	//	AND ur.user_id ISNULL;
	var userRolesToAdd []models.UsersRoles
	err := u.db.Select("r.id as role_id, ? as user_id").
		RightJoin("roles r", "r.id=users_roles.role_id AND users_roles.user_id = ?", userID, userID).
		Where("role_type IN (?) AND users_roles.user_id IS NULL", rs).
		All(&userRolesToAdd)
	if err != nil {
		return []models.UsersRoles{}, err

	}
	err = u.db.Create(userRolesToAdd)
	if err != nil {
		return []models.UsersRoles{}, err

	}
	return userRolesToAdd, nil
}

func (u usersRolesCreator) removeUserRoles(userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
	//Having to use somewhat convoluted right join syntax b/c FROM clause in pop is derived from the model
	//and for the RawQuery was having trouble passing in array into the in clause with additional params
	//ideally would just be the query below
	//SELECT r.id                              AS role_id,
	//	'3b9360a3-3304-4c60-90f4-83d687884079' AS user_id,
	//	ur.user_id,
	//	role_type
	//FROM roles r
	//		LEFT JOIN users_roles ur ON r.id = ur.role_id
	//	AND ur.user_id = '3b9360a3-3304-4c60-90f4-83d687884079'
	//WHERE role_type NOT IN ('transportation_ordering_officer', 'contracting_officer')
	//	AND ur.user_id IS NOT NULL;
	var userRolesToDelete []models.UsersRoles
	err := u.db.Select("users_roles.id, r.id as role_id, ? as user_id").
		RightJoin("roles r", "r.id=users_roles.role_id AND users_roles.user_id = ?", userID, userID).
		Where("role_type NOT IN (?) AND users_roles.id IS NOT NULL", rs).
		All(&userRolesToDelete)
	if err != nil {
		return []models.UsersRoles{}, err
	}
	err = u.db.Destroy(userRolesToDelete)
	if err != nil {
		return []models.UsersRoles{}, err
	}
	return userRolesToDelete, nil
}
