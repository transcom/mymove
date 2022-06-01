package authentication

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"

	"github.com/transcom/mymove/pkg/auth"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models/roles"
)

// TODO: placeholder until we figure out where these should be stored
type RolePermissions struct {
	RoleType    roles.RoleType
	Permissions []string
}

var TOO = RolePermissions{
	RoleType:    roles.RoleTypeTOO,
	Permissions: []string{"update.financial_review_flag", "update.orders"},
}

var TIO = RolePermissions{
	RoleType:    roles.RoleTypeTIO,
	Permissions: []string{"update.financial_review_flag", "update.orders"},
}

var ServicesCounselor = RolePermissions{
	RoleType:    roles.RoleTypeServicesCounselor,
	Permissions: []string{"update.financial_review_flag", "update.orders"},
}

var QAECSR = RolePermissions{
	RoleType:    roles.RoleTypeQaeCsr,
	Permissions: []string{},
}

var AllRolesPermissions = []RolePermissions{TOO, TIO, ServicesCounselor, QAECSR}

// check if a [user.role] has permissions on a given object
func checkUserPermission(appCtx appcontext.AppContext, session *auth.Session, permission string) (bool, error) {

	userPermissions, err := getPermissionsForUser(appCtx, session.UserID)

	if err != nil {
		return false, err
	}

	for _, perm := range userPermissions {
		if permission == perm {
			fmt.Println("PERMISSION GRANTED: ", permission)
			return true, nil
		}
	}

	return false, nil
}

// for a given user return the permissions associated with their roles
func getPermissionsForUser(appCtx appcontext.AppContext, userID uuid.UUID) ([]string, error) {
	var userPermissions []string

	//check the users roles
	userRoles, err := getRolesForUser(appCtx, userID)
	if err != nil {
		return nil, err
	}

	for _, ur := range userRoles {
		for _, rp := range AllRolesPermissions {

			if ur == rp.RoleType {
				userPermissions = append(userPermissions, rp.Permissions...)
			}
		}
	}

	return userPermissions, nil
}

// load the [user.role] given a valid user ID
// what we care about here is the string, so we can look it up for permissions --> roles.role_type
func getRolesForUser(appCtx appcontext.AppContext, userID uuid.UUID) ([]roles.RoleType, error) {
	var userRoleTypes []roles.RoleType

	err := appCtx.DB().RawQuery(`SELECT roles.role_type
		FROM roles
			LEFT JOIN users_roles ur
			    ON roles.id = ur.role_id
			WHERE ur.deleted_at IS NULL AND ur.user_id = ?`, userID).All(&userRoleTypes)

	fmt.Printf("USER ROLESS: %+v\n", userRoleTypes)

	return userRoleTypes, err
}
