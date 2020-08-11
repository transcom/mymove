package event

// -------------------- API NAMES --------------------

// primeAPIName is a const string to use the EndpointTypes
const supportAPIName string = "supportapi"

// -------------------- ENDPOINT KEYS --------------------

// SupportListMTOsEndpointKey is the key for the listMTOs endpoint in support
const SupportListMTOsEndpointKey = "Support.ListMTOs"

// SupportCreateMoveTaskOrderEndpointKey is the key for the createMoveTaskOrder endpoint in support
const SupportCreateMoveTaskOrderEndpointKey = "Support.CreateMoveTaskOrder"

// SupportGetMoveTaskOrderEndpointKey is the key for the getMoveTaskOrder endpoint in support
const SupportGetMoveTaskOrderEndpointKey = "Support.GetMoveTaskOrder"

// SupportMakeMoveTaskOrderAvailableEndpointKey is the key for the makeMoveTaskOrderAvailable endpoint in support
const SupportMakeMoveTaskOrderAvailableEndpointKey = "Support.MakeMoveTaskOrderAvailable"

// SupportListMTOPaymentRequestsEndpointKey is the key for the listMTOPaymentRequests endpoint in support
const SupportListMTOPaymentRequestsEndpointKey = "Support.ListMTOPaymentRequests"

// SupportUpdatePaymentRequestStatusEndpointKey is the key for the updatePaymentRequestStatus endpoint in support
const SupportUpdatePaymentRequestStatusEndpointKey = "Support.UpdatePaymentRequestStatus"

// SupportUpdateMTOServiceItemStatusEndpointKey is the key for the updateMTOServiceItemStatus endpoint in support
const SupportUpdateMTOServiceItemStatusEndpointKey = "Support.UpdateMTOServiceItemStatus"

// SupportUpdateMTOShipmentStatusEndpointKey is the key for the updateMTOShipmentStatus endpoint in support
const SupportUpdateMTOShipmentStatusEndpointKey = "Support.UpdateMTOShipmentStatus"

// -------------------- ENDPOINT MAP --------------------
var supportEndpoints = EndpointMapType{
	SupportListMTOsEndpointKey: {
		APIName:     supportAPIName,
		OperationID: "listMTOs",
	},
	SupportCreateMoveTaskOrderEndpointKey: {
		APIName:     supportAPIName,
		OperationID: "createMoveTaskOrder",
	},
	SupportGetMoveTaskOrderEndpointKey: {
		APIName:     supportAPIName,
		OperationID: "getMoveTaskOrder",
	},
	SupportMakeMoveTaskOrderAvailableEndpointKey: {
		APIName:     supportAPIName,
		OperationID: "makeMoveTaskOrderAvailable",
	},
	SupportListMTOPaymentRequestsEndpointKey: {
		APIName:     supportAPIName,
		OperationID: "listMTOPaymentRequests",
	},
	SupportUpdatePaymentRequestStatusEndpointKey: {
		APIName:     supportAPIName,
		OperationID: "updatePaymentRequestStatus",
	},
	SupportUpdateMTOServiceItemStatusEndpointKey: {
		APIName:     supportAPIName,
		OperationID: "updateMTOServiceItemStatus",
	},
	SupportUpdateMTOShipmentStatusEndpointKey: {
		APIName:     supportAPIName,
		OperationID: "updateMTOShipmentStatus",
	},
}
