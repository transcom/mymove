package usersroles

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type usersRolesCreator struct {
	checks []usersRolesValidator
}

// NewUsersRolesCreator creates a new struct with the service dependencies
func NewUsersRolesCreator() services.UserRoleAssociator {
	return usersRolesCreator{
		checks: []usersRolesValidator{
			checkTransportationOfficerPolicyViolation(),
		},
	}
}

// UpdateUserRoles associates a given user with a set of roles
func (u usersRolesCreator) UpdateUserRoles(appCtx appcontext.AppContext, userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
	var usersRoles []models.UsersRoles

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		_, err := u.addUserRoles(appCtx, userID, rs)
		if err != nil {
			return err
		}
		_, err = u.removeUserRoles(appCtx, userID, rs)
		if err != nil {
			return err
		}
		// fetch + return updated roles
		err = appCtx.DB().Where("user_id = ?", userID).All(&usersRoles)
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		return []models.UsersRoles{}, txErr
	}

	return usersRoles, nil
}

func (u usersRolesCreator) addUserRoles(appCtx appcontext.AppContext, userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
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

	// Retrieve existing active roles for the user
	var existingUserRoles []models.UsersRoles
	err := appCtx.DB().
		Select("users_roles.*").
		RightJoin("roles r", "r.id = users_roles.role_id").
		Where("users_roles.user_id = ? AND users_roles.deleted_at IS NULL", userID).
		All(&existingUserRoles)
	if err != nil {
		return []models.UsersRoles{}, err
	}

	// Identify which roles need to be added
	var userRolesToAdd []models.UsersRoles
	if len(rs) > 0 {
		err := appCtx.DB().Select("r.id as role_id, ? as user_id").
			RightJoin("roles r", "r.id=users_roles.role_id AND users_roles.user_id = ? AND users_roles.deleted_at IS NULL", userID, userID).
			Where("role_type IN (?) AND (users_roles.user_id IS NULL)", rs).
			All(&userRolesToAdd)
		if err != nil {
			return []models.UsersRoles{}, err

		}
	}
	err = validateUsersRoles(appCtx, &userRolesToAdd, &existingUserRoles, u.checks...)
	if err != nil {
		return nil, err
	}

	err = appCtx.DB().Create(userRolesToAdd)
	if err != nil {
		return []models.UsersRoles{}, err

	}
	return userRolesToAdd, nil
}

func (u usersRolesCreator) removeUserRoles(appCtx appcontext.AppContext, userID uuid.UUID, rs []roles.RoleType) ([]models.UsersRoles, error) {
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
	if len(rs) > 0 {
		err := appCtx.DB().Select("users_roles.id, r.id as role_id, ? as user_id, users_roles.deleted_at").
			RightJoin("roles r", "r.id=users_roles.role_id AND users_roles.user_id = ? AND users_roles.deleted_at IS NULL", userID, userID).
			Where("role_type NOT IN (?) AND users_roles.id IS NOT NULL", rs).
			All(&userRolesToDelete)
		if err != nil {
			return []models.UsersRoles{}, err
		}
	}
	// query above wont work if nothing in rs array (i.e this user should have no roles)
	// b/c of how pop expands empty array rs, below just removes all roles in this situation
	if len(rs) == 0 {
		err := appCtx.DB().Where("user_id = ?", userID).
			All(&userRolesToDelete)
		if err != nil {
			return []models.UsersRoles{}, err
		}
	}
	for _, roleToDelete := range userRolesToDelete {
		copyOfRoleToDelete := roleToDelete // Make copy to avoid implicit memory aliasing of items from a range statement.
		err := utilities.SoftDestroy(appCtx.DB(), &copyOfRoleToDelete)
		if err != nil {
			return []models.UsersRoles{}, err
		}
	}
	return userRolesToDelete, nil
}
