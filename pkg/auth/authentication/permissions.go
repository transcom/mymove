package authentication

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
)

// TODO: placeholder until we figure out where these should be stored
type RolePermissions struct {
	RoleType    string
	Permissions []string
}

var TOO = RolePermissions{
	RoleType: "transportation_ordering_officer",
	Permissions: []string{"update.move", "create.serviceItem",
		"update.shipment"},
}

var TIO = RolePermissions{
	RoleType:    "transportation_invoicing_officer",
	Permissions: []string{"create.serviceItem", "update.shipment"},
}

var AllRolesPermissions = []RolePermissions{TOO, TIO}

// check if a [user.role] has permissions on a given object
func checkUserPermission(appCtx appcontext.AppContext, permission string) (bool, error) {
	userID := appCtx.Session().UserID

	userPermissions, err := getPermissionsForUser(appCtx, userID)

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
func getRolesForUser(appCtx appcontext.AppContext, userID uuid.UUID) ([]string, error) {

	var roleTypes []string

	err := appCtx.DB().RawQuery(`SELECT DISTINCT role_type
		FROM roles
			RIGHT JOIN users_roles ur
			    ON ur.role_id = roles.id
			    	AND ur.deleted_at IS NULL
			        AND ur.user_id = ? `, userID).All(&roleTypes)

	fmt.Printf("USER ROLESS: %+v\n", roleTypes)

	return roleTypes, err

}
