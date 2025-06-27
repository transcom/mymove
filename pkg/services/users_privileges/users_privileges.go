package usersprivileges

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

type usersPrivilegesCreator struct {
}

// NewUsersPrivilegesCreator creates a new struct with the service dependencies
func NewUsersPrivilegesCreator() services.UserPrivilegeAssociator {
	return usersPrivilegesCreator{}
}

// UpdateUserPrivileges associates a given user with a set of privileges
func (u usersPrivilegesCreator) UpdateUserPrivileges(appCtx appcontext.AppContext, userID uuid.UUID, rs []roles.PrivilegeType) ([]models.UsersPrivileges, error) {
	_, err := u.addUserPrivileges(appCtx, userID, rs)
	if err != nil {
		return []models.UsersPrivileges{}, err
	}
	_, err = u.removeUserPrivileges(appCtx, userID, rs)
	if err != nil {
		return []models.UsersPrivileges{}, err
	}
	var usersPrivileges []models.UsersPrivileges
	// fetch + return updated privileges
	err = appCtx.DB().Where("user_id = ?", userID).All(&usersPrivileges)
	if err != nil {
		return []models.UsersPrivileges{}, err
	}
	return usersPrivileges, nil
}

func (u usersPrivilegesCreator) addUserPrivileges(appCtx appcontext.AppContext, userID uuid.UUID, rs []roles.PrivilegeType) ([]models.UsersPrivileges, error) {
	//Having to use somewhat convoluted right join syntax b/c FROM clause in pop is derived from the model
	//and for the RawQuery was having trouble passing in array into the in clause with additional params
	//ideally would just be the query below
	//SELECT r.id                              AS privilege_id,
	//	'3b9360a3-3304-4c60-90f4-83d687884079' AS user_id,
	//	ur.user_id,
	//	privilege_type
	//FROM privileges r
	//		LEFT JOIN users_privileges ur ON r.id = ur.privilege_id
	//	AND ur.user_id = '3b9360a3-3304-4c60-90f4-83d687884079'
	//WHERE privilege_type IN ('supervisor')
	//	AND ur.user_id ISNULL;
	var userPrivilegesToAdd []models.UsersPrivileges
	if len(rs) > 0 {
		err := appCtx.DB().Select("r.id as privilege_id, ? as user_id").
			RightJoin("privileges r", "r.id=users_privileges.privilege_id AND users_privileges.user_id = ? AND users_privileges.deleted_at IS NULL", userID, userID).
			Where("privilege_type IN (?) AND (users_privileges.user_id IS NULL)", rs).
			All(&userPrivilegesToAdd)
		if err != nil {
			return []models.UsersPrivileges{}, err

		}
	}
	err := appCtx.DB().Create(userPrivilegesToAdd)
	if err != nil {
		return []models.UsersPrivileges{}, err

	}
	return userPrivilegesToAdd, nil
}

func (u usersPrivilegesCreator) removeUserPrivileges(appCtx appcontext.AppContext, userID uuid.UUID, rs []roles.PrivilegeType) ([]models.UsersPrivileges, error) {
	//Having to use somewhat convoluted right join syntax b/c FROM clause in pop is derived from the model
	//and for the RawQuery was having trouble passing in array into the in clause with additional params
	//ideally would just be the query below
	//SELECT r.id                              AS privilege_id,
	//	'3b9360a3-3304-4c60-90f4-83d687884079' AS user_id,
	//	ur.user_id,
	//	privilege_type
	//FROM privileges r
	//		LEFT JOIN users_privileges ur ON r.id = ur.privilege_id
	//	AND ur.user_id = '3b9360a3-3304-4c60-90f4-83d687884079'
	//WHERE privilege_type NOT IN ('supervisor')
	//	AND ur.user_id IS NOT NULL;
	var userPrivilegesToDelete []models.UsersPrivileges
	if len(rs) > 0 {
		err := appCtx.DB().Select("users_privileges.id, r.id as privilege_id, ? as user_id, users_privileges.deleted_at").
			RightJoin("privileges r", "r.id=users_privileges.privilege_id AND users_privileges.user_id = ? AND users_privileges.deleted_at IS NULL", userID, userID).
			Where("privilege_type NOT IN (?) AND users_privileges.id IS NOT NULL", rs).
			All(&userPrivilegesToDelete)
		if err != nil {
			return []models.UsersPrivileges{}, err
		}
	}
	// query above wont work if nothing in rs array (i.e this user should have no privileges)
	// b/c of how pop expands empty array rs, below just removes all privileges in this situation
	if len(rs) == 0 {
		err := appCtx.DB().Where("user_id = ?", userID).
			All(&userPrivilegesToDelete)
		if err != nil {
			return []models.UsersPrivileges{}, err
		}
	}
	for _, privilegeToDelete := range userPrivilegesToDelete {
		copyOfPrivilegeToDelete := privilegeToDelete // Make copy to avoid implicit memory aliasing of items from a range statement.
		err := utilities.SoftDestroy(appCtx.DB(), &copyOfPrivilegeToDelete)
		if err != nil {
			return []models.UsersPrivileges{}, err
		}
	}
	return userPrivilegesToDelete, nil
}

func (u usersPrivilegesCreator) VerifyUserPrivilegeAllowed(appCtx appcontext.AppContext, roles []*adminmessages.OfficeUserRole, privileges []*adminmessages.OfficeUserPrivilege) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	for _, privilege := range privileges {
		for _, role := range roles {

			var allowed bool
			err := appCtx.DB().RawQuery("SELECT * FROM is_role_privilege_allowed($1, $2)", role.RoleType, privilege.PrivilegeType).First(&allowed)

			if err != nil {
				return nil, err
			}

			if !allowed {
				err = apperror.NewBadDataError(fmt.Sprintf("%s is not an authorized role for %s privileges", *role.Name, *privilege.Name))
				appCtx.Logger().Error(err.Error())
				verrs.Add("Validation Error", err.Error())
			}
		}
	}

	if verrs.HasAny() {
		return verrs, nil

	}

	return nil, nil
}

func (f usersPrivilegesCreator) FetchPrivilegesForUser(appCtx appcontext.AppContext, userID uuid.UUID) (roles.Privileges, error) {
	var privileges roles.Privileges
	err := appCtx.DB().Q().Join("users_privileges", "users_privileges.privilege_id = privileges.id").
		Where("users_privileges.deleted_at IS NULL AND users_privileges.user_id = ?", (userID)).
		All(&privileges)
	return privileges, err
}
