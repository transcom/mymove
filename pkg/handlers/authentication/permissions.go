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
		"read.shipmentsPaymentSITBalance",
		"read.paymentServiceItemStatus",
		"update.move",
		"update.shipment",
		"update.financialReviewFlag",
		"update.orders",
		"update.allowances",
		"update.billableWeight",
		"update.SITExtension",
		"update.MTOServiceItem",
		"update.excessWeightRisk",
		"update.customer",
		"view.closeoutOffice",
		"update.MTOPage",
	},
}

// TODO: FOR NOW HQ WILL USE SAME PERMISSIONS AS TOO, B-20014 WILL CHANGE THESE TO READ ONLY
var HQ = RolePermissions{
	RoleType: roles.RoleTypeHQ,
	Permissions: []string{
		"create.serviceItem",
		"create.shipmentDiversionRequest",
		"create.reweighRequest",
		"create.shipmentCancellation",
		"create.SITExtension",
		"read.paymentRequest",
		"read.shipmentsPaymentSITBalance",
		"read.paymentServiceItemStatus",
		"update.move",
		"update.shipment",
		"update.financialReviewFlag",
		"update.orders",
		"update.allowances",
		"update.billableWeight",
		"update.SITExtension",
		"update.MTOServiceItem",
		"update.excessWeightRisk",
		"update.customer",
		"view.closeoutOffice",
		"update.MTOPage",
	},
}

var TIO = RolePermissions{
	RoleType: roles.RoleTypeTIO,
	Permissions: []string{
		"create.serviceItem",
		"read.paymentRequest",
		"read.shipmentsPaymentSITBalance",
		"update.financialReviewFlag",
		"update.orders",
		"update.billableWeight",
		"update.maxBillableWeight",
		"update.paymentRequest",
		"update.paymentServiceItemStatus",
		"update.MTOPage",
		"update.customer",
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
		"update.closeoutOffice",
		"view.closeoutOffice",
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
		"view.closeoutOffice",
		"read.shipmentsPaymentSITBalance",
	},
}

var AllRolesPermissions = []RolePermissions{TOO, TIO, ServicesCounselor, QAECSR, HQ}

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
	userRoles, err := roles.FetchRolesForUser(appCtx.DB(), userID)

	var userRoleTypes []roles.RoleType
	for i := range userRoles {
		userRoleTypes = append(userRoleTypes, userRoles[i].RoleType)
	}

	if err != nil {
		logger.Warn("Error while looking up user roles: ", zap.String("user role lookup error: ", err.Error()))
		return nil, err
	}

	logger.Info("User has the following roles: ", zap.Any("user roles", userRoleTypes))

	return userRoleTypes, nil
}
