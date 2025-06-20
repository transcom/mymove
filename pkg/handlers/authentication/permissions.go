package authentication

import (
	"slices"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models/roles"
)

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
		"create.supportingDocuments",
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
		"update.closeoutOffice",
		"update.MTOPage",
		"create.TXOShipment",
		"update.cancelMoveFlag",
	},
}

var HQ = RolePermissions{
	RoleType: roles.RoleTypeHQ,
	Permissions: []string{
		"read.paymentRequest",
		"read.shipmentsPaymentSITBalance",
		"read.paymentServiceItemStatus",
		"view.closeoutOffice",
	},
}

var TIO = RolePermissions{
	RoleType: roles.RoleTypeTIO,
	Permissions: []string{
		"create.serviceItem",
		"create.supportingDocuments",
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
		"create.supportingDocuments",
		"update.financialReviewFlag",
		"update.shipment",
		"update.orders",
		"update.allowances",
		"update.billableWeight",
		"update.MTOServiceItem",
		"update.customer",
		"update.closeoutOffice",
		"view.closeoutOffice",
		"update.cancelMoveFlag",
	},
}

var QAE = RolePermissions{
	RoleType: roles.RoleTypeQae,
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

var ContractingOfficer = RolePermissions{
	RoleType: roles.RoleTypeContractingOfficer,
	Permissions: []string{
		"create.shipmentTermination",
		"read.paymentRequest",
		"view.closeoutOffice",
		"read.shipmentsPaymentSITBalance",
	},
}

var CustomerServiceRepresentative = RolePermissions{
	RoleType: roles.RoleTypeCustomerServiceRepresentative,
	Permissions: []string{
		"read.paymentRequest",
		"view.closeoutOffice",
		"read.shipmentsPaymentSITBalance",
	},
}

var GSR = RolePermissions{
	RoleType: roles.RoleTypeGSR,
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

var AllRolesPermissions = []RolePermissions{TOO, TIO, ServicesCounselor, QAE, CustomerServiceRepresentative, HQ, GSR, ContractingOfficer}

// check if a [user.role] has permissions on a given object
func checkUserPermission(session auth.Session, permission string) bool {
	return slices.Contains(session.Permissions, permission)
}

// for a given user return the permissions associated with their roles given the current session role
func getPermissionsForUser(appCtx appcontext.AppContext) []string {
	var userPermissions []string

	session := appCtx.Session()
	if session != nil {
		return GetPermissionsForRole(session.ActiveRole.RoleType)
	}

	return userPermissions
}

func GetPermissionsForRole(roleType roles.RoleType) []string {
	for _, rp := range AllRolesPermissions {
		if rp.RoleType == roleType {
			return rp.Permissions
		}
	}
	return []string{}
}
