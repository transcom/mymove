package authentication

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models/roles"
)

// TODO: placeholder until we figure out where these should be stored
type RolePermissions struct {
	RoleType    roles.RoleType
	Permissions []string
}

var TOO = RolePermissions{
	RoleType: roles.RoleTypeTOO,
	Permissions: []string{
		"create.serviceItem",
		"create.shipmentDiversionRequest",
		"create.reweighRequest",
		"create.shipmentCancellation",
		"create.SITExtension",
		"read.paymentRequest",
		"update.move",
		"update.shipment",
		"update.financialReviewFlag",
		"update.orders",
		"update.allowances",
		"update.billableWeight",
		"update.SITExtension",
		"update.MTOServiceItem",
		"update.paymentServiceItemStatus",
		"update.excessWeightRisk",
	},
}

var TIO = RolePermissions{
	RoleType: roles.RoleTypeTIO,
	Permissions: []string{
		"create.ShipmentCancellation",
		"create.shipmentDiversionRequest",
		"create.reweighRequest",
		"create.serviceItem",
		"read.paymentRequest",
		"read.shipmentsPaymentSITBalance",
		"update.shipment",
		"update.financialReviewFlag",
		"update.orders",
		"update.allowances",
		"update.billableWeight",
		"update.maxBillableWeight",
		"update.paymentRequest",
		"update.paymentServiceItemStatus",
		"update.MTOServiceItem",
	},
}

var ServicesCounselor = RolePermissions{
	RoleType: roles.RoleTypeServicesCounselor,
	Permissions: []string{
		"create.shipmentDiversionRequest",
		"create.reweighRequest",
		"update.financialReviewFlag",
		"update.shipment",
		"update.orders",
		"update.allowances",
		"update.billableWeight",
		"update.MTOServiceItem",
		"update.customer",
	},
}

var QAECSR = RolePermissions{
	RoleType: roles.RoleTypeQaeCsr,
	Permissions: []string{
		"create.reportViolation",
		"create.evaluationReport",
		"read.paymentRequest",
		"update.evaluationReport",
		"delete.evaluationReport",
	},
}

var AllRolesPermissions = []RolePermissions{TOO, TIO, ServicesCounselor, QAECSR}

// check if a [user.role] has permissions on a given object
func checkUserPermission(appCtx appcontext.AppContext, session *auth.Session, permission string) (bool, error) {

	logger := appCtx.Logger()
	userPermissions := getPermissionsForUser(appCtx, session.UserID)

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
func getPermissionsForUser(appCtx appcontext.AppContext, userID uuid.UUID) []string {
	var userPermissions []string

	// check the users roles
	userRoles, err := getRolesForUser(appCtx, userID)
	// if there's an error looking up roles return an empty permission array
	if err != nil {
		return userPermissions
	}

	for _, ur := range userRoles {
		for _, rp := range AllRolesPermissions {

			if ur == rp.RoleType {
				userPermissions = append(userPermissions, rp.Permissions...)
			}
		}
	}

	return userPermissions
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

	logger.Info("User has the following roles: ", zap.Any("user roles", userRoleTypes))

	return userRoleTypes, nil
}
