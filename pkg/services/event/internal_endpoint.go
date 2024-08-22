package event

// -------------------- API NAMES --------------------

// InternalAPIName is a const string to use the EndpointTypes
const InternalAPIName string = "internalapi"

// -------------------- ENDPOINT KEYS --------------------

// InternalShowLoggedInUserEndpointKey is the key for the showLoggedInUser endpoint in internal
const InternalShowLoggedInUserEndpointKey = "Internal.ShowLoggedInUser"

// InternalIsLoggedInUserEndpointKey is the key for the isLoggedInUser endpoint in internal
const InternalIsLoggedInUserEndpointKey = "Internal.IsLoggedInUser"

// InternalCreateOrdersEndpointKey is the key for the createOrders endpoint in internal
const InternalCreateOrdersEndpointKey = "Internal.CreateOrders"

// InternalUpdateOrdersEndpointKey is the key for the updateOrders endpoint in internal
const InternalUpdateOrdersEndpointKey = "Internal.UpdateOrders"

// InternalShowOrdersEndpointKey is the key for the showOrders endpoint in internal
const InternalShowOrdersEndpointKey = "Internal.ShowOrders"

// InternalUploadAmendedOrdersEndpointKey is the key for the UploadAmendedOrders endpoint in internal
const InternalUploadAmendedOrdersEndpointKey = "Internal.UploadAmendedOrders"

// InternalPatchMoveEndpointKey is the key for the patchMove endpoint in internal
const InternalPatchMoveEndpointKey = "Internal.PatchMove"

// InternalShowMoveEndpointKey is the key for the showMove endpoint in internal
const InternalShowMoveEndpointKey = "Internal.ShowMove"

// InternalCreateSignedCertificationEndpointKey is the key for the createSignedCertification endpoint in internal
const InternalCreateSignedCertificationEndpointKey = "Internal.CreateSignedCertification"

// InternalIndexSignedCertificationEndpointKey is the key for the indexSignedCertification endpoint in internal
const InternalIndexSignedCertificationEndpointKey = "Internal.IndexSignedCertification"

// InternalApproveReimbursementEndpointKey is the key for the approveReimbursement endpoint in internal
const InternalApproveReimbursementEndpointKey = "Internal.ApproveReimbursement"

// InternalShowOfficeOrdersEndpointKey is the key for the showOfficeOrders endpoint in internal
const InternalShowOfficeOrdersEndpointKey = "Internal.ShowOfficeOrders"

// InternalIndexMoveDocumentsEndpointKey is the key for the indexMoveDocuments endpoint in internal
const InternalIndexMoveDocumentsEndpointKey = "Internal.IndexMoveDocuments"

// InternalCreateGenericMoveDocumentEndpointKey is the key for the createGenericMoveDocument endpoint in internal
const InternalCreateGenericMoveDocumentEndpointKey = "Internal.CreateGenericMoveDocument"

// InternalUpdateMoveDocumentEndpointKey is the key for the updateMoveDocument endpoint in internal
const InternalUpdateMoveDocumentEndpointKey = "Internal.UpdateMoveDocument"

// InternalCreateMovingExpenseDocumentEndpointKey is the key for the createMovingExpenseDocument endpoint in internal
const InternalCreateMovingExpenseDocumentEndpointKey = "Internal.CreateMovingExpenseDocument"

// InternalApproveMoveEndpointKey is the key for the approveMove endpoint in internal
const InternalApproveMoveEndpointKey = "Internal.ApproveMove"

// InternalSubmitMoveForApprovalEndpointKey is the key for the submitMoveForApproval endpoint in internal
const InternalSubmitMoveForApprovalEndpointKey = "Internal.SubmitMoveForApproval"

// InternalCancelMoveEndpointKey is the key for the cancelMove endpoint in internal
const InternalCancelMoveEndpointKey = "Internal.CancelMove"

// InternalShowMoveDatesSummaryEndpointKey is the key for the showMoveDatesSummary endpoint in internal
const InternalShowMoveDatesSummaryEndpointKey = "Internal.ShowMoveDatesSummary"

// InternalShowShipmentSummaryWorksheetEndpointKey is the key for the showShipmentSummaryWorksheet endpoint in internal
const InternalShowShipmentSummaryWorksheetEndpointKey = "Internal.ShowShipmentSummaryWorksheet"

// InternalApprovePPMEndpointKey is the key for the approvePPM endpoint in internal
const InternalApprovePPMEndpointKey = "Internal.ApprovePPM"

// InternalCreateDocumentEndpointKey is the key for the createDocument endpoint in internal
const InternalCreateDocumentEndpointKey = "Internal.CreateDocument"

// InternalShowDocumentEndpointKey is the key for the showDocument endpoint in internal
const InternalShowDocumentEndpointKey = "Internal.ShowDocument"

// InternalCreateUploadEndpointKey is the key for the createUpload endpoint in internal
const InternalCreateUploadEndpointKey = "Internal.CreateUpload"

// InternalCreateServiceMemberEndpointKey is the key for the createServiceMember endpoint in internal
const InternalCreateServiceMemberEndpointKey = "Internal.CreateServiceMember"

// InternalShowServiceMemberEndpointKey is the key for the showServiceMember endpoint in internal
const InternalShowServiceMemberEndpointKey = "Internal.ShowServiceMember"

// InternalPatchServiceMemberEndpointKey is the key for the patchServiceMember endpoint in internal
const InternalPatchServiceMemberEndpointKey = "Internal.PatchServiceMember"

// InternalShowServiceMemberOrdersEndpointKey is the key for the showServiceMemberOrders endpoint in internal
const InternalShowServiceMemberOrdersEndpointKey = "Internal.ShowServiceMemberOrders"

// InternalCreateServiceMemberBackupContactEndpointKey is the key for the createServiceMemberBackupContact endpoint in internal
const InternalCreateServiceMemberBackupContactEndpointKey = "Internal.CreateServiceMemberBackupContact"

// InternalIndexServiceMemberBackupContactsEndpointKey is the key for the indexServiceMemberBackupContacts endpoint in internal
const InternalIndexServiceMemberBackupContactsEndpointKey = "Internal.IndexServiceMemberBackupContacts"

// InternalShowServiceMemberBackupContactEndpointKey is the key for the showServiceMemberBackupContact endpoint in internal
const InternalShowServiceMemberBackupContactEndpointKey = "Internal.ShowServiceMemberBackupContact"

// InternalUpdateServiceMemberBackupContactEndpointKey is the key for the updateServiceMemberBackupContact endpoint in internal
const InternalUpdateServiceMemberBackupContactEndpointKey = "Internal.UpdateServiceMemberBackupContact"

// InternalSearchDutyLocationsEndpointKey is the key for the searchDutyLocations endpoint in internal
const InternalSearchDutyLocationsEndpointKey = "Internal.SearchDutyLocations"

// InternalShowCounselingOfficesEndpointKey is the key for the showCounselingOffices endpoint in internal
const InternalShowCounselingOfficesEndpointKey = "Internal.ShowCounselingOffices"

// InternalShowDutyLocationsTransportationOfficeEndpointKey is the key for the showDutyLocationTransportationOffice endpoint in internal
const InternalShowDutyLocationTransportationOfficeEndpointKey = "Internal.ShowDutyLocationTransportationOffice"

// InternalShowQueueEndpointKey is the key for the showQueue endpoint in internal
const InternalShowQueueEndpointKey = "Internal.ShowQueue"

// InternalIndexEntitlementsEndpointKey is the key for the indexEntitlements endpoint in internal
const InternalIndexEntitlementsEndpointKey = "Internal.IndexEntitlements"

// InternalValidateEntitlementEndpointKey is the key for the validateEntitlement endpoint in internal
const InternalValidateEntitlementEndpointKey = "Internal.ValidateEntitlement"

// InternalShowAvailableMoveDatesEndpointKey is the key for the showAvailableMoveDates endpoint in internal
const InternalShowAvailableMoveDatesEndpointKey = "Internal.ShowAvailableMoveDates"

// InternalGetCookieURLEndpointKey is the key for the getCookieURL endpoint in internal
const InternalGetCookieURLEndpointKey = "Internal.GetCookieURL"

// InternalCreateWeightTicketDocumentEndpointKey is the key for the createWeightTicketDocument endpoint in internal
const InternalCreateWeightTicketDocumentEndpointKey = "Internal.CreateWeightTicketDocument"

// InternalValidatePostalCodeWithRateDataEndpointKey is the key for the validatePostalCodeWithRateData endpoint in internal
const InternalValidatePostalCodeWithRateDataEndpointKey = "Internal.ValidatePostalCodeWithRateData"

// InternalShowAddressEndpointKey is the key for the showAddress endpoint in internal
const InternalShowAddressEndpointKey = "Internal.ShowAddress"

// InternalCreateMTOShipmentEndpointKey is the key for the createMTOShipment endpoint in internal
const InternalCreateMTOShipmentEndpointKey = "Internal.CreateMTOShipment"

// InternalCreateWeightTicketEndpointKey is the key for the createWeightTicket endpoint in internal
const InternalCreateWeightTicketEndpointKey = "Internal.CreateWeightTicket"

// InternalUpdateWeightTicketEndpointKey is the key for the updateWeightTicket endpoint in internal
const InternalUpdateWeightTicketEndpointKey = "Internal.UpdateWeightTicket"

// -------------------- ENDPOINT MAP ENTRIES --------------------
var internalEndpoints = EndpointMapType{
	InternalShowLoggedInUserEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showLoggedInUser",
	},
	InternalIsLoggedInUserEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "isLoggedInUser",
	},
	InternalCreateOrdersEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createOrders",
	},
	InternalUpdateOrdersEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "updateOrders",
	},
	InternalShowOrdersEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showOrders",
	},
	InternalUploadAmendedOrdersEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "uploadAmendedOrders",
	},
	InternalPatchMoveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "patchMove",
	},
	InternalShowMoveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showMove",
	},
	InternalCreateSignedCertificationEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createSignedCertification",
	},
	InternalIndexSignedCertificationEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "indexSignedCertification",
	},
	InternalApproveReimbursementEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "approveReimbursement",
	},
	InternalShowOfficeOrdersEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showOfficeOrders",
	},
	InternalIndexMoveDocumentsEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "indexMoveDocuments",
	},
	InternalCreateGenericMoveDocumentEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createGenericMoveDocument",
	},
	InternalUpdateMoveDocumentEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "updateMoveDocument",
	},
	InternalCreateMovingExpenseDocumentEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createMovingExpenseDocument",
	},
	InternalApproveMoveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "approveMove",
	},
	InternalSubmitMoveForApprovalEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "submitMoveForApproval",
	},
	InternalCancelMoveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "cancelMove",
	},
	InternalShowMoveDatesSummaryEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showMoveDatesSummary",
	},
	InternalShowShipmentSummaryWorksheetEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showShipmentSummaryWorksheet",
	},
	InternalApprovePPMEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "approvePPM",
	},
	InternalCreateDocumentEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createDocument",
	},
	InternalShowDocumentEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showDocument",
	},
	InternalCreateUploadEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createUpload",
	},
	InternalCreateServiceMemberEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createServiceMember",
	},
	InternalShowServiceMemberEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showServiceMember",
	},
	InternalPatchServiceMemberEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "patchServiceMember",
	},
	InternalShowServiceMemberOrdersEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showServiceMemberOrders",
	},
	InternalCreateServiceMemberBackupContactEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createServiceMemberBackupContact",
	},
	InternalIndexServiceMemberBackupContactsEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "indexServiceMemberBackupContacts",
	},
	InternalShowServiceMemberBackupContactEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showServiceMemberBackupContact",
	},
	InternalUpdateServiceMemberBackupContactEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "updateServiceMemberBackupContact",
	},
	InternalSearchDutyLocationsEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "searchDutyLocations",
	},
	InternalShowCounselingOfficesEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showCounselingOffices",
	},
	InternalShowDutyLocationTransportationOfficeEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showDutyLocationTransportationOffice",
	},
	InternalShowQueueEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showQueue",
	},
	InternalIndexEntitlementsEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "indexEntitlements",
	},
	InternalValidateEntitlementEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "validateEntitlement",
	},
	InternalShowAvailableMoveDatesEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showAvailableMoveDates",
	},
	InternalGetCookieURLEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "getCookieURL",
	},
	InternalCreateWeightTicketDocumentEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createWeightTicketDocument",
	},
	InternalCreateWeightTicketEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createWeightTicket",
	},
	InternalUpdateWeightTicketEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "updateWeightTicket",
	},
	InternalValidatePostalCodeWithRateDataEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "validatePostalCodeWithRateData",
	},
	InternalShowAddressEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showAddress",
	},
	InternalCreateMTOShipmentEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createMTOShipment",
	},
}
