package event

// -------------------- API NAMES --------------------

// InternalAPIName is a const string to use the EndpointTypes
const InternalAPIName string = "internalapi"

// -------------------- ENDPOINT KEYS --------------------

// InternalShowPPMEstimateEndpointKey is the key for the showPPMEstimate endpoint in internal
const InternalShowPPMEstimateEndpointKey = "Internal.ShowPPMEstimate"

// InternalShowPPMSitEstimateEndpointKey is the key for the showPPMSitEstimate endpoint in internal
const InternalShowPPMSitEstimateEndpointKey = "Internal.ShowPPMSitEstimate"

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

// InternalCreatePersonallyProcuredMoveEndpointKey is the key for the createPersonallyProcuredMove endpoint in internal
const InternalCreatePersonallyProcuredMoveEndpointKey = "Internal.CreatePersonallyProcuredMove"

// InternalIndexPersonallyProcuredMovesEndpointKey is the key for the indexPersonallyProcuredMoves endpoint in internal
const InternalIndexPersonallyProcuredMovesEndpointKey = "Internal.IndexPersonallyProcuredMoves"

// InternalUpdatePersonallyProcuredMoveEstimateEndpointKey is the key for the updatePersonallyProcuredMoveEstimate endpoint in internal
const InternalUpdatePersonallyProcuredMoveEstimateEndpointKey = "Internal.UpdatePersonallyProcuredMoveEstimate"

// InternalUpdatePersonallyProcuredMoveEndpointKey is the key for the updatePersonallyProcuredMove endpoint in internal
const InternalUpdatePersonallyProcuredMoveEndpointKey = "Internal.UpdatePersonallyProcuredMove"

// InternalPatchPersonallyProcuredMoveEndpointKey is the key for the patchPersonallyProcuredMove endpoint in internal
const InternalPatchPersonallyProcuredMoveEndpointKey = "Internal.PatchPersonallyProcuredMove"

// InternalShowPersonallyProcuredMoveEndpointKey is the key for the showPersonallyProcuredMove endpoint in internal
const InternalShowPersonallyProcuredMoveEndpointKey = "Internal.ShowPersonallyProcuredMove"

// InternalSubmitPersonallyProcuredMoveEndpointKey is the key for the submitPersonallyProcuredMove endpoint in internal
const InternalSubmitPersonallyProcuredMoveEndpointKey = "Internal.SubmitPersonallyProcuredMove"

// InternalRequestPPMExpenseSummaryEndpointKey is the key for the requestPPMExpenseSummary endpoint in internal
const InternalRequestPPMExpenseSummaryEndpointKey = "Internal.RequestPPMExpenseSummary"

// InternalRequestPPMPaymentEndpointKey is the key for the requestPPMPayment endpoint in internal
const InternalRequestPPMPaymentEndpointKey = "Internal.RequestPPMPayment"

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

// InternalCreatePPMAttachmentsEndpointKey is the key for the createPPMAttachments endpoint in internal
const InternalCreatePPMAttachmentsEndpointKey = "Internal.CreatePPMAttachments"

// InternalApprovePPMEndpointKey is the key for the approvePPM endpoint in internal
const InternalApprovePPMEndpointKey = "Internal.ApprovePPM"

// InternalShowPPMIncentiveEndpointKey is the key for the showPPMIncentive endpoint in internal
const InternalShowPPMIncentiveEndpointKey = "Internal.ShowPPMIncentive"

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

// InternalFetchAccessCodeEndpointKey is the key for the fetchAccessCode endpoint in internal
const InternalFetchAccessCodeEndpointKey = "Internal.FetchAccessCode"

// InternalValidateAccessCodeEndpointKey is the key for the validateAccessCode endpoint in internal
const InternalValidateAccessCodeEndpointKey = "Internal.ValidateAccessCode"

// InternalClaimAccessCodeEndpointKey is the key for the claimAccessCode endpoint in internal
const InternalClaimAccessCodeEndpointKey = "Internal.ClaimAccessCode"

// InternalShowAddressEndpointKey is the key for the showAddress endpoint in internal
const InternalShowAddressEndpointKey = "Internal.ShowAddress"

// InternalCreateMTOShipmentEndpointKey is the key for the createMTOShipment endpoint in internal
const InternalCreateMTOShipmentEndpointKey = "Internal.CreateMTOShipment"

// -------------------- ENDPOINT MAP ENTRIES --------------------
var internalEndpoints = EndpointMapType{
	InternalShowPPMEstimateEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showPPMEstimate",
	},
	InternalShowPPMSitEstimateEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showPPMSitEstimate",
	},
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
	InternalCreatePersonallyProcuredMoveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createPersonallyProcuredMove",
	},
	InternalIndexPersonallyProcuredMovesEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "indexPersonallyProcuredMoves",
	},
	InternalUpdatePersonallyProcuredMoveEstimateEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "updatePersonallyProcuredMoveEstimate",
	},
	InternalUpdatePersonallyProcuredMoveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "updatePersonallyProcuredMove",
	},
	InternalPatchPersonallyProcuredMoveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "patchPersonallyProcuredMove",
	},
	InternalShowPersonallyProcuredMoveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showPersonallyProcuredMove",
	},
	InternalSubmitPersonallyProcuredMoveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "submitPersonallyProcuredMove",
	},
	InternalRequestPPMExpenseSummaryEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "requestPPMExpenseSummary",
	},
	InternalRequestPPMPaymentEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "requestPPMPayment",
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
	InternalCreatePPMAttachmentsEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "createPPMAttachments",
	},
	InternalApprovePPMEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "approvePPM",
	},
	InternalShowPPMIncentiveEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "showPPMIncentive",
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
	InternalValidatePostalCodeWithRateDataEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "validatePostalCodeWithRateData",
	},
	InternalFetchAccessCodeEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "fetchAccessCode",
	},
	InternalValidateAccessCodeEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "validateAccessCode",
	},
	InternalClaimAccessCodeEndpointKey: {
		APIName:     InternalAPIName,
		OperationID: "claimAccessCode",
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
