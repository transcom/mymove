package event

// -------------------- API NAMES --------------------

// GhcAPIName is a const string to use the EndpointTypes
const GhcAPIName string = "ghcapi"

// -------------------- ENDPOINT KEYS --------------------

// GhcGetCustomerEndpointKey is the key for the getCustomer endpoint in ghc
const GhcGetCustomerEndpointKey = "Ghc.GetCustomer"

// GhcGetMoveEndpointKey is the key for the getMove endpoint in ghc
const GhcGetMoveEndpointKey = "Ghc.GetMove"

// GhcGetMovesQueueEndpointKey is the key for the getMovesQueue endpoint in ghc
const GhcGetMovesQueueEndpointKey = "Ghc.GetMovesQueue"

// GhcGetOrderEndpointKey is the key for the getOrder endpoint in ghc
const GhcGetOrderEndpointKey = "Ghc.GetOrder"

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

// GhcUpdateMTOShipmentEndpointKey is the key for the updateMTOShipment endpoint in ghc
const GhcUpdateMTOShipmentEndpointKey = "Ghc.UpdateMTOShipment"

// GhcDeleteShipmentEndpointKey is the key for the deleteShipment endpoint in ghc
const GhcDeleteShipmentEndpointKey = "Ghc.DeleteShipment"

// GhcApproveShipmentEndpointKey is the key for the approveShipment endpoint in ghc
const GhcApproveShipmentEndpointKey = "Ghc.ApproveShipment"

// GhcRequestShipmentDiversionEndpointKey is the key for the requestShipmentDiversion endpoint in ghc
const GhcRequestShipmentDiversionEndpointKey = "Ghc.RequestShipmentDiversion"

// GhcApproveShipmentDiversionEndpointKey is the key for the approveShipmentDiversion endpoint in ghc
const GhcApproveShipmentDiversionEndpointKey = "Ghc.ApproveShipmentDiversion"

// GhcRejectShipmentEndpointKey is the key for the rejectShipment endpoint in ghc
const GhcRejectShipmentEndpointKey = "Ghc.RejectShipment"

// GhcRequestShipmentCancellationEndpointKey is the key for the requestShipmentCancellation endpoint in ghc
const GhcRequestShipmentCancellationEndpointKey = "Ghc.RequestShipmentCancellation"

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

// GhcUpdateMTOStatusServiceCounselingCompletedEndpointKey is the key for the updateMTOStatusServiceCounselingCompleted endpoint in ghc
const GhcUpdateMTOStatusServiceCounselingCompletedEndpointKey = "Ghc.UpdateMTOStatusServiceCounselingCompleted"

// GhcUpdatePaymentServiceItemStatusEndpointKey is the key for the updatePaymentServiceItemStatus endpoint in ghc
const GhcUpdatePaymentServiceItemStatusEndpointKey = "Ghc.UpdatePaymentServiceItemStatus"

// GhcGetEntitlementsEndpointKey is the key for the getEntitlements endpoint in ghc
const GhcGetEntitlementsEndpointKey = "Ghc.GetEntitlements"

// GhcGetPaymentRequestsQueueEndpointKey is the key for the getPaymentRequestsQueue endpoint in ghc
const GhcGetPaymentRequestsQueueEndpointKey = "Ghc.GetPaymentRequestsQueue"

// GhcGetPaymentRequestEndpointKey is the key for the getPaymentRequest endpoint in ghc
const GhcGetPaymentRequestEndpointKey = "Ghc.GetPaymentRequest"

// GhcUpdatePaymentRequestStatusEndpointKey is the key for the updatePaymentRequestStatus endpoint in ghc
const GhcUpdatePaymentRequestStatusEndpointKey = "Ghc.UpdatePaymentRequestStatus"

// GhcUpdateOrderEndpointKey is the key for the updateOrder endpoint in ghc
const GhcUpdateOrderEndpointKey = "Ghc.UpdateOrder"

// GhcCounselingUpdateOrderEndpointKey is the key for the counselingUpdateOrder endpoint in ghc
const GhcCounselingUpdateOrderEndpointKey = "Ghc.CounselingUpdateOrder"

// GhcUpdateAllowanceEndpointKey is the key for the updateAllowance endpoint in ghc
const GhcUpdateAllowanceEndpointKey = "Ghc.UpdateAllowance"

// GhcCounselingUpdateAllowanceEndpointKey is the key for the counselingUpdateAllowance endpoint in ghc
const GhcCounselingUpdateAllowanceEndpointKey = "Ghc.CounselingUpdateAllowance"

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
	GhcGetMovesQueueEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getMovesQueue",
	},
	GhcGetOrderEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getOrder",
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
	GhcUpdateMTOShipmentEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updateMTOShipment",
	},
	GhcDeleteShipmentEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "deleteShipment",
	},
	GhcApproveShipmentEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "approveShipment",
	},
	GhcRequestShipmentDiversionEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "requestShipmentDiversion",
	},
	GhcApproveShipmentDiversionEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "approveShipmentDiversion",
	},
	GhcRejectShipmentEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "rejectShipment",
	},
	GhcRequestShipmentCancellationEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "requestShipmentCancellation",
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
	GhcUpdateMTOStatusServiceCounselingCompletedEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updateMTOStatusServiceCounselingCompleted",
	},
	GhcUpdatePaymentServiceItemStatusEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updatePaymentServiceItemStatus",
	},
	GhcGetEntitlementsEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getEntitlements",
	},
	GhcGetPaymentRequestsQueueEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getPaymentRequestsQueue",
	},
	GhcGetPaymentRequestEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "getPaymentRequest",
	},
	GhcUpdatePaymentRequestStatusEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updatePaymentRequestStatus",
	},
	GhcUpdateOrderEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updateOrder",
	},
	GhcCounselingUpdateOrderEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "counselingUpdateOrder",
	},
	GhcUpdateAllowanceEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "updateAllowance",
	},
	GhcCounselingUpdateAllowanceEndpointKey: {
		APIName:     GhcAPIName,
		OperationID: "counselingUpdateAllowance",
	},
}
