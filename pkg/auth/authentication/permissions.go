package authentication

import (
	"github.com/spf13/cast"
	"go.uber.org/zap"

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
	RoleType: roles.RoleTypeTOO,
	Permissions: []string{"update.move", "create.serviceItem",
		"update.shipment", "update.financialReviewFlag", "update.orders", "update.allowances"},
}

var TIO = RolePermissions{
	RoleType:    roles.RoleTypeTIO,
	Permissions: []string{"create.serviceItem", "update.shipment", "update.financialReviewFlag", "update.orders", "update.allowances"},
}

var ServicesCounselor = RolePermissions{
	RoleType:    roles.RoleTypeServicesCounselor,
	Permissions: []string{"update.financialReviewFlag", "update.shipment", "update.orders", "update.allowances"},
}

var QAECSR = RolePermissions{
	RoleType:    roles.RoleTypeQaeCsr,
	Permissions: []string{"read.move"},
}

var AllRolesPermissions = []RolePermissions{TOO, TIO, ServicesCounselor, QAECSR}

// check if a [user.role] has permissions on a given object
func checkUserPermission(appCtx appcontext.AppContext, session *auth.Session, permission string) (bool, error) {

	logger := appCtx.Logger()
	userPermissions, err := getPermissionsForUser(appCtx, session.UserID)

	if err != nil {
		logger.Error("Error while looking up permissions: ", zap.String("permission error", err.Error()))
		return false, err
	}

	for _, perm := range userPermissions {
		if permission == perm {
			logger.Info("PERMISSION GRANTED: ", zap.String("permission", permission))
			return true, nil
		}
	}

	logger.Warn("Permission not granted for user, ", zap.String("permission denied to user with session IDToken: ", session.IDToken))
	return false, nil
}

// for a given user return the permissions associated with their roles
func getPermissionsForUser(appCtx appcontext.AppContext, userID uuid.UUID) ([]string, error) {
	logger := appCtx.Logger()
	var userPermissions []string

	//check the users roles
	userRoles, err := getRolesForUser(appCtx, userID)
	if err != nil {
		logger.Error("Error while looking up user roles: ", zap.String("permission error", err.Error()))
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
	logger := appCtx.Logger()
	var userRoleTypes []roles.RoleType

	err := appCtx.DB().RawQuery(`SELECT roles.role_type
		FROM roles
			LEFT JOIN users_roles ur
			    ON roles.id = ur.role_id
			WHERE ur.deleted_at IS NULL AND ur.user_id = ?`, userID).All(&userRoleTypes)

	if err != nil {
		logger.Warn("Error while looking up user roles: ", zap.String("user role lookup error: ", err.Error()))
		return nil, err
	}

	logger.Info("User has the following roles: ", zap.String("user roles", cast.ToString(userRoleTypes)))

	return userRoleTypes, err
}
