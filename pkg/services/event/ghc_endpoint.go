package event

// -------------------- API NAMES --------------------

// GhcAPIName is a const string to use the EndpointTypes
const GhcAPIName string = "ghcapi"

// -------------------- ENDPOINT KEYS --------------------

// GhcGetCustomerEndpointKey is the key for the getCustomer endpoint in ghc
const GhcGetCustomerEndpointKey = "Ghc.GetCustomer"

// GhcGetMoveEndpointKey is the key for the getMove endpoint in ghc
const GhcGetMoveEndpointKey = "Ghc.GetMove"

// GhcListMoveOrdersEndpointKey is the key for the listMoveOrders endpoint in ghc
const GhcListMoveOrdersEndpointKey = "Ghc.ListMoveOrders"

// GhcGetMoveOrderEndpointKey is the key for the getMoveOrder endpoint in ghc
const GhcGetMoveOrderEndpointKey = "Ghc.GetMoveOrder"

// GhcListMoveTaskOrdersEndpointKey is the key for the listMoveTaskOrders endpoint in ghc
const GhcListMoveTaskOrdersEndpointKey = "Ghc.ListMoveTaskOrders"

// GhcGetMoveTaskOrderEndpointKey is the key for the getMoveTaskOrder endpoint in ghc
const GhcGetMoveTaskOrderEndpointKey = "Ghc.GetMoveTaskOrder"

// GhcUpdateMoveTaskOrderEndpointKey is the key for the updateMoveTaskOrder endpoint in ghc
const GhcUpdateMoveTaskOrderEndpointKey = "Ghc.UpdateMoveTaskOrder"

// GhcListMTOServiceItemsEndpointKey is the key for the listMTOServiceItems endpoint in ghc
const GhcListMTOServiceItemsEndpointKey = "Ghc.ListMTOServiceItems"

// GhcCreateMTOServiceItemEndpointKey is the key for the createMTOServiceItem endpoint in ghc
const GhcCreateMTOServiceItemEndpointKey = "Ghc.CreateMTOServiceItem"

// GhcListMTOShipmentsEndpointKey is the key for the listMTOShipments endpoint in ghc
const GhcListMTOShipmentsEndpointKey = "Ghc.ListMTOShipments"

// GhcPatchMTOShipmentStatusEndpointKey is the key for the patchMTOShipmentStatus endpoint in ghc
const GhcPatchMTOShipmentStatusEndpointKey = "Ghc.PatchMTOShipmentStatus"

// GhcFetchMTOAgentListEndpointKey is the key for the fetchMTOAgentList endpoint in ghc
const GhcFetchMTOAgentListEndpointKey = "Ghc.FetchMTOAgentList"

// GhcGetMTOServiceItemEndpointKey is the key for the getMTOServiceItem endpoint in ghc
const GhcGetMTOServiceItemEndpointKey = "Ghc.GetMTOServiceItem"

// GhcUpdateMTOServiceItemEndpointKey is the key for the updateMTOServiceItem endpoint in ghc
const GhcUpdateMTOServiceItemEndpointKey = "Ghc.UpdateMTOServiceItem"

// GhcUpdateMTOServiceItemStatusEndpointKey is the key for the updateMTOServiceItemStatus endpoint in ghc
const GhcUpdateMTOServiceItemStatusEndpointKey = "Ghc.UpdateMTOServiceItemStatus"

// GhcUpdateMoveTaskOrderStatusEndpointKey is the key for the updateMoveTaskOrderStatus endpoint in ghc
const GhcUpdateMoveTaskOrderStatusEndpointKey = "Ghc.UpdateMoveTaskOrderStatus"

// GhcUpdatePaymentServiceItemStatusEndpointKey is the key for the updatePaymentServiceItemStatus endpoint in ghc
const GhcUpdatePaymentServiceItemStatusEndpointKey = "Ghc.UpdatePaymentServiceItemStatus"

// GhcGetEntitlementsEndpointKey is the key for the getEntitlements endpoint in ghc
const GhcGetEntitlementsEndpointKey = "Ghc.GetEntitlements"

// GhcListPaymentRequestsEndpointKey is the key for the listPaymentRequests endpoint in ghc
const GhcListPaymentRequestsEndpointKey = "Ghc.ListPaymentRequests"

// GhcGetPaymentRequestEndpointKey is the key for the getPaymentRequest endpoint in ghc
const GhcGetPaymentRequestEndpointKey = "Ghc.GetPaymentRequest"

// GhcUpdatePaymentRequestStatusEndpointKey is the key for the updatePaymentRequestStatus endpoint in ghc
const GhcUpdatePaymentRequestStatusEndpointKey = "Ghc.UpdatePaymentRequestStatus"

// -------------------- ENDPOINT MAP ENTRIES --------------------
var ghcEndpoints = EndpointMapType{
	GhcGetCustomerEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getCustomer",
	},
	GhcGetMoveEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getMove",
	},
	GhcListMoveOrdersEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "listMoveOrders",
	},
	GhcGetMoveOrderEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getMoveOrder",
	},
	GhcListMoveTaskOrdersEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "listMoveTaskOrders",
	},
	GhcGetMoveTaskOrderEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getMoveTaskOrder",
	},
	GhcUpdateMoveTaskOrderEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updateMoveTaskOrder",
	},
	GhcListMTOServiceItemsEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "listMTOServiceItems",
	},
	GhcCreateMTOServiceItemEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "createMTOServiceItem",
	},
	GhcListMTOShipmentsEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "listMTOShipments",
	},
	GhcPatchMTOShipmentStatusEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "patchMTOShipmentStatus",
	},
	GhcFetchMTOAgentListEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "fetchMTOAgentList",
	},
	GhcGetMTOServiceItemEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getMTOServiceItem",
	},
	GhcUpdateMTOServiceItemEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updateMTOServiceItem",
	},
	GhcUpdateMTOServiceItemStatusEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updateMTOServiceItemStatus",
	},
	GhcUpdateMoveTaskOrderStatusEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updateMoveTaskOrderStatus",
	},
	GhcUpdatePaymentServiceItemStatusEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updatePaymentServiceItemStatus",
	},
	GhcGetEntitlementsEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getEntitlements",
	},
	GhcListPaymentRequestsEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "listPaymentRequests",
	},
	GhcGetPaymentRequestEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getPaymentRequest",
	},
	GhcUpdatePaymentRequestStatusEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updatePaymentRequestStatus",
	},
}
