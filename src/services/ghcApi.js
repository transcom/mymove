import Swagger from 'swagger-client';

import { makeSwaggerRequest, requestInterceptor, responseInterceptor, makeSwaggerRequestRaw } from './swaggerRequest';

let ghcClient = null;

// setting up the same config from Swagger/api.js
export async function getGHCClient() {
  if (!ghcClient) {
    ghcClient = await Swagger({
      url: '/ghc/v1/swagger.yaml',
      requestInterceptor,
      responseInterceptor,
    });
  }
  return ghcClient;
}

export async function makeGHCRequest(operationPath, params = {}, options = {}) {
  const client = await getGHCClient();
  return makeSwaggerRequest(client, operationPath, params, options);
}

export async function makeGHCRequestRaw(operationPath, params = {}) {
  const client = await getGHCClient();
  return makeSwaggerRequestRaw(client, operationPath, params);
}

export async function getPaymentRequest(key, paymentRequestID) {
  return makeGHCRequest('paymentRequests.getPaymentRequest', { paymentRequestID });
}

export async function getPPMDocuments(key, shipmentID) {
  return makeGHCRequest('ppm.getPPMDocuments', { shipmentID }, { normalize: false });
}

export async function patchWeightTicket({ ppmShipmentId, weightTicketId, payload, eTag }) {
  return makeGHCRequest(
    'ppm.updateWeightTicket',
    {
      ppmShipmentId,
      weightTicketId,
      'If-Match': eTag,
      updateWeightTicketPayload: payload,
    },
    {
      normalize: false,
    },
  );
}

export async function patchExpense({ ppmShipmentId, movingExpenseId, payload, eTag }) {
  return makeGHCRequest(
    'ppm.updateMovingExpense',
    {
      ppmShipmentId,
      movingExpenseId,
      'If-Match': eTag,
      updateMovingExpense: payload,
    },
    {
      normalize: false,
    },
  );
}

export async function patchProGearWeightTicket({ ppmShipmentId, proGearWeightTicketId, payload, eTag }) {
  return makeGHCRequest(
    'ppm.updateProGearWeightTicket',
    {
      ppmShipmentId,
      proGearWeightTicketId,
      'If-Match': eTag,
      updateProGearWeightTicket: payload,
    },
    {
      normalize: false,
    },
  );
}

export async function getPPMCloseout(key, ppmShipmentId) {
  return makeGHCRequest('ppm.getPPMCloseout', { ppmShipmentId }, { normalize: false });
}

export async function getPPMSITEstimatedCost(
  key,
  ppmShipmentId,
  sitLocation,
  sitEntryDate,
  sitDepartureDate,
  weightStored,
) {
  return makeGHCRequest(
    'ppm.getPPMSITEstimatedCost',
    {
      ppmShipmentId,
      sitLocation,
      sitEntryDate,
      sitDepartureDate,
      weightStored,
    },
    { normalize: false },
  );
}

export async function getPPMActualWeight(key, ppmShipmentId) {
  return makeGHCRequest('ppm.getPPMActualWeight', { ppmShipmentId }, { normalize: false });
}

export async function patchPPMDocumentsSetStatus({ ppmShipmentId, eTag }) {
  return makeGHCRequest(
    'ppm.finishDocumentReview',
    {
      ppmShipmentId,
      'If-Match': eTag,
    },
    {
      normalize: false,
    },
  );
}

export async function getMove(key, locator) {
  return makeGHCRequest('move.getMove', { locator }, { normalize: false });
}

export async function getPrimeSimulatorAvailableMoves(key, { filters = [], currentPage = 1, currentPageSize = 20 }) {
  const operationPath = 'queues.listPrimeMoves';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  return makeGHCRequest(
    operationPath,
    { page: currentPage, perPage: currentPageSize, ...paramFilters },
    { schemaKey: 'listMoves', normalize: false },
  );
}

export async function getCustomerSupportRemarksForMove(key, locator) {
  return makeGHCRequest('customerSupportRemarks.getCustomerSupportRemarksForMove', { locator }, { normalize: false });
}

export async function createCustomerSupportRemarkForMove({ body, locator }) {
  return makeGHCRequest('customerSupportRemarks.createCustomerSupportRemarkForMove', {
    body,
    locator,
  });
}

export async function updateCustomerSupportRemarkForMove({ body, customerSupportRemarkID }) {
  return makeGHCRequest('customerSupportRemarks.updateCustomerSupportRemarkForMove', {
    body,
    customerSupportRemarkID,
  });
}

export async function deleteCustomerSupportRemark({ customerSupportRemarkID }) {
  return makeGHCRequest(
    'customerSupportRemarks.deleteCustomerSupportRemark',
    { customerSupportRemarkID },
    { normalize: false },
  );
}

export async function createCounselingEvaluationReport({ moveCode }) {
  return makeGHCRequest('evaluationReports.createEvaluationReport', { locator: moveCode }, { normalize: false });
}

export async function createShipmentEvaluationReport({ body, moveCode }) {
  return makeGHCRequest('evaluationReports.createEvaluationReport', { locator: moveCode, body }, { normalize: false });
}

export async function deleteEvaluationReport(reportID) {
  return makeGHCRequest('evaluationReports.deleteEvaluationReport', { reportID }, { normalize: false });
}

export async function saveEvaluationReport({ reportID, ifMatchETag, body }) {
  return makeGHCRequest(
    'evaluationReports.saveEvaluationReport',
    { reportID, 'If-Match': ifMatchETag, body },
    { normalize: false },
  );
}

export async function submitEvaluationReport({ reportID, ifMatchETag }) {
  return makeGHCRequest(
    'evaluationReports.submitEvaluationReport',
    { reportID, 'If-Match': ifMatchETag },
    { normalize: false },
  );
}

export async function addSeriousIncidentAppeal({ reportID, body }) {
  return makeGHCRequest('evaluationReports.addAppealToSeriousIncident', { reportID, body }, { normalize: false });
}

export async function addViolationAppeal({ reportID, reportViolationID, body }) {
  return makeGHCRequest(
    'evaluationReports.addAppealToViolation',
    { reportID, reportViolationID, body },
    { normalize: false },
  );
}

export async function associateReportViolations({ reportID, body }) {
  return makeGHCRequest('reportViolations.associateReportViolations', { reportID, body }, { normalize: false });
}

export async function getReportViolationsByReportID(key, reportID) {
  return makeGHCRequest('reportViolations.getReportViolationsByReportID', { reportID }, { normalize: false });
}

export async function getMoveHistory(key, { moveCode, currentPage = 1, currentPageSize = 20 }) {
  return makeGHCRequest(
    'move.getMoveHistory',
    { locator: moveCode, page: currentPage, perPage: currentPageSize },
    { schemaKey: 'MoveHistoryResult', normalize: false },
  );
}

export async function getOrder(key, orderID) {
  return makeGHCRequest('order.getOrder', { orderID });
}

export async function getMovePaymentRequests(key, locator) {
  return makeGHCRequest(
    'paymentRequests.getPaymentRequestsForMove',
    { locator },
    { schemaKey: 'paymentRequests', normalize: false },
  );
}

export async function getMTOShipments(key, moveTaskOrderID, normalize = true) {
  return makeGHCRequest('mtoShipment.listMTOShipments', { moveTaskOrderID }, { schemaKey: 'mtoShipments', normalize });
}

export async function getShipmentEvaluationReports(key, moveID) {
  return makeGHCRequest(
    'move.getMoveShipmentEvaluationReportsList',
    { moveID },
    { schemaKey: 'evaluationReports', normalize: false },
  );
}

export async function getEvaluationReportByID(key, reportID) {
  return makeGHCRequest('evaluationReports.getEvaluationReport', { reportID }, { normalize: false });
}

export async function getCounselingEvaluationReports(key, moveID) {
  return makeGHCRequest(
    'move.getMoveCounselingEvaluationReportsList',
    { moveID },
    { schemaKey: 'evaluationReports', normalize: false },
  );
}

export async function getPWSViolations() {
  return makeGHCRequest('pwsViolations.getPWSViolations', {}, { normalize: false });
}

export async function getMTOServiceItems(key, moveTaskOrderID, normalize = true) {
  return makeGHCRequest(
    'mtoServiceItem.listMTOServiceItems',
    { moveTaskOrderID },
    { schemaKey: 'mtoServiceItems', normalize },
  );
}

export async function getDocument(key, documentId) {
  return makeGHCRequest('ghcDocuments.getDocument', { documentId }, { schemaKey: 'document' });
}
export async function getCustomer(key, customerID) {
  return makeGHCRequest('customer.getCustomer', { customerID });
}

export async function searchMoves(key, { sort, order, filters = [], currentPage = 1, currentPageSize = 20 }) {
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  if (paramFilters.status) {
    paramFilters.status = paramFilters.status.split(',');
  }
  if (paramFilters.shipmentsCount) {
    paramFilters.shipmentsCount = Number(paramFilters.shipmentsCount);
  }
  return makeGHCRequest(
    'move.searchMoves',
    {
      body: {
        sort,
        order,
        page: currentPage,
        perPage: currentPageSize,
        ...paramFilters,
      },
    },
    { schemaKey: 'searchMovesResult', normalize: false },
  );
}

export async function patchMTOServiceItemStatus({ moveId, mtoServiceItemID, ifMatchEtag, status, rejectionReason }) {
  return makeGHCRequest(
    'mtoServiceItem.updateMTOServiceItemStatus',
    {
      moveTaskOrderID: moveId,
      mtoServiceItemID,
      'If-Match': ifMatchEtag,
      body: { status, rejectionReason },
    },
    { schemaKey: 'mtoServiceItem' },
  );
}

export async function patchPaymentRequest({ paymentRequestID, status, ifMatchETag, rejectionReason }) {
  return makeGHCRequest('paymentRequests.updatePaymentRequestStatus', {
    paymentRequestID,
    'If-Match': ifMatchETag,
    body: { status, rejectionReason },
  });
}

export async function patchPaymentServiceItemStatus({
  moveTaskOrderID,
  paymentServiceItemID,
  status,
  ifMatchEtag,
  rejectionReason,
}) {
  const operationPath = 'paymentServiceItem.updatePaymentServiceItemStatus';
  return makeGHCRequest(
    operationPath,
    {
      moveTaskOrderID,
      paymentServiceItemID,
      'If-Match': ifMatchEtag,
      body: { status, rejectionReason },
    },
    { label: operationPath, schemaKey: 'paymentServiceItem' },
  );
}

export async function getTacValid({ tac }) {
  const operationPath = 'order.tacValidation';
  return makeGHCRequest(operationPath, { tac }, { normalize: false });
}

// Retrieves the line of accounting based on a given TAC,
// effective date, and service member affiliation
export async function getLoa({ tacCode, effectiveDate, departmentIndicator }) {
  const operationPath = 'linesOfAccounting.requestLineOfAccounting';
  return makeGHCRequest(operationPath, { body: { tacCode, effectiveDate, departmentIndicator } }, { normalize: false });
}

export async function updateOrder({ orderID, ifMatchETag, body }) {
  const operationPath = 'order.updateOrder';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag, body });
}

export async function counselingUpdateOrder({ orderID, ifMatchETag, body }) {
  const operationPath = 'order.counselingUpdateOrder';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag, body });
}

export async function counselingCreateOrder({ body }) {
  const operationPath = 'order.createOrder';
  return makeGHCRequest(operationPath, { createOrders: body }, { normalize: true });
}

export async function updateUpload({ uploadID, body }) {
  const operationPath = 'uploads.updateUpload';
  return makeGHCRequest(operationPath, { uploadID, body });
}

export async function updateAllowance({ orderID, ifMatchETag, body }) {
  const operationPath = 'order.updateAllowance';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag, body });
}

export async function counselingUpdateAllowance({ orderID, ifMatchETag, body }) {
  const operationPath = 'order.counselingUpdateAllowance';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag, body });
}

export async function updateBillableWeight({ orderID, ifMatchETag, body }) {
  const operationPath = 'order.updateBillableWeight';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag, body });
}

export async function updateMaxBillableWeightAsTIO({ orderID, ifMatchETag, body }) {
  const operationPath = 'order.updateMaxBillableWeightAsTIO';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag, body });
}

export async function acknowledgeExcessWeightRisk({ orderID, ifMatchETag }) {
  const operationPath = 'order.acknowledgeExcessWeightRisk';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag });
}

export async function createCustomerWithOktaOption({ body }) {
  const operationPath = 'customer.createCustomerWithOktaOption';
  return makeGHCRequest(operationPath, { body });
}

export async function updateCustomerInfo({ customerId, ifMatchETag, body }) {
  const operationPath = 'customer.updateCustomer';
  return makeGHCRequest(operationPath, { customerID: customerId, 'If-Match': ifMatchETag, body });
}

export async function updateMTOReviewedBillableWeights({ moveTaskOrderID, ifMatchETag }) {
  const operationPath = 'moveTaskOrder.updateMTOReviewedBillableWeightsAt';
  return makeGHCRequest(operationPath, { moveTaskOrderID, 'If-Match': ifMatchETag });
}

export async function updateTIORemarks({ moveTaskOrderID, ifMatchETag, body }) {
  const operationPath = 'moveTaskOrder.updateMoveTIORemarks';
  return makeGHCRequest(operationPath, { moveTaskOrderID, 'If-Match': ifMatchETag, body });
}

export function updateMoveStatus({ moveTaskOrderID, ifMatchETag, mtoApprovalServiceItemCodes, normalize = true }) {
  const operationPath = 'moveTaskOrder.updateMoveTaskOrderStatus';
  return makeGHCRequest(
    operationPath,
    {
      moveTaskOrderID,
      'If-Match': ifMatchETag,
      serviceItemCodes: mtoApprovalServiceItemCodes,
    },
    { normalize },
  );
}

export function updateMoveStatusServiceCounselingCompleted({ moveTaskOrderID, ifMatchETag, normalize = false }) {
  const operationPath = 'moveTaskOrder.updateMTOStatusServiceCounselingCompleted';
  return makeGHCRequest(
    operationPath,
    {
      moveTaskOrderID,
      'If-Match': ifMatchETag,
    },
    { normalize },
  );
}

export function cancelMove({ moveID, normalize = false }) {
  const operationPath = 'move.moveCanceler';
  return makeGHCRequest(
    operationPath,
    {
      moveID,
    },
    { normalize },
  );
}

export function updateMTOShipmentStatus({
  shipmentID,
  diversionReason,
  operationPath,
  ifMatchETag,
  normalize = true,
  schemaKey = 'mtoShipment',
}) {
  return makeGHCRequest(
    operationPath,
    {
      shipmentID,
      'If-Match': ifMatchETag,
      body: {
        diversionReason,
      },
    },
    { schemaKey, normalize },
  );
}

export function updateMTOShipmentRequestReweigh({
  shipmentID,
  ifMatchETag,
  normalize = false,
  schemaKey = 'mtoShipment',
}) {
  const operationPath = 'shipment.requestShipmentReweigh';
  return makeGHCRequest(
    operationPath,
    {
      shipmentID,
      'If-Match': ifMatchETag,
    },
    { schemaKey, normalize },
  );
}

export function createMTOShipment({ body, normalize = true, schemaKey = 'mtoShipment' }) {
  const operationPath = 'mtoShipment.createMTOShipment';
  return makeGHCRequest(operationPath, { body }, { schemaKey, normalize });
}

export async function getMTOShipmentByID(key, shipmentID) {
  return makeGHCRequest('mtoShipment.getShipment', { shipmentID }, { schemaKey: 'mtoShipment', normalize: false });
}

export function updateMTOShipment({
  moveTaskOrderID,
  shipmentID,
  ifMatchETag,
  normalize = true,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'mtoShipment.updateMTOShipment';
  return makeGHCRequest(
    operationPath,
    {
      moveTaskOrderID,
      shipmentID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function approveSITExtension({
  shipmentID,
  sitExtensionID,
  ifMatchETag,
  normalize = true,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'shipment.approveSITExtension';
  return makeGHCRequest(
    operationPath,
    {
      shipmentID,
      sitExtensionID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function denySITExtension({
  shipmentID,
  sitExtensionID,
  ifMatchETag,
  normalize = true,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'shipment.denySITExtension';
  return makeGHCRequest(
    operationPath,
    {
      shipmentID,
      sitExtensionID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function submitSITExtension({ shipmentID, ifMatchETag, normalize = true, schemaKey = 'mtoShipment', body }) {
  const operationPath = 'shipment.createApprovedSITDurationUpdate';
  return makeGHCRequest(
    operationPath,
    {
      shipmentID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function updateSITServiceItemCustomerExpense({
  shipmentID,
  ifMatchETag,
  normalize = true,
  convertToCustomerExpense,
  customerExpenseReason,
}) {
  return makeGHCRequest(
    'shipment.updateSITServiceItemCustomerExpense',
    {
      shipmentID,
      'If-Match': ifMatchETag,
      body: { convertToCustomerExpense, customerExpenseReason },
    },
    { schemaKey: 'mtoShipment', normalize },
  );
}

export function deleteShipment({ shipmentID, normalize = false, schemaKey = 'shipment' }) {
  const operationPath = 'shipment.deleteShipment';
  return makeGHCRequest(
    operationPath,
    {
      shipmentID,
    },
    { schemaKey, normalize },
  );
}

export async function getMovesQueue(
  key,
  { sort, order, filters = [], currentPage = 1, currentPageSize = 20, viewAsGBLOC },
) {
  const operationPath = 'queues.getMovesQueue';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  return makeGHCRequest(
    operationPath,
    { sort, order, page: currentPage, perPage: currentPageSize, viewAsGBLOC, ...paramFilters },
    { schemaKey: 'queueMovesResult', normalize: false },
  );
}

export async function getServicesCounselingQueue(
  key,
  { sort, order, filters = [], currentPage = 1, currentPageSize = 20, needsPPMCloseout = false, viewAsGBLOC },
) {
  const operationPath = 'queues.getServicesCounselingQueue';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });

  return makeGHCRequest(
    operationPath,
    {
      sort,
      order,
      page: currentPage,
      perPage: currentPageSize,
      needsPPMCloseout,
      viewAsGBLOC,
      ...paramFilters,
    },

    { schemaKey: 'queueMovesResult', normalize: false },
  );
}

export async function getServicesCounselingOriginLocations(needsPPMCloseout) {
  const operationPath = 'queues.getServicesCounselingOriginList';

  return makeGHCRequest(
    operationPath,
    {
      needsPPMCloseout,
    },

    { schemaKey: 'Locations', normalize: false },
  );
}

export async function getServicesCounselingPPMQueue(
  key,
  { sort, order, filters = [], currentPage = 1, currentPageSize = 20, needsPPMCloseout = true, viewAsGBLOC },
) {
  const operationPath = 'queues.getServicesCounselingQueue';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });

  return makeGHCRequest(
    operationPath,
    { sort, order, page: currentPage, perPage: currentPageSize, needsPPMCloseout, viewAsGBLOC, ...paramFilters },
    { schemaKey: 'queueMovesResult', normalize: false },
  );
}

export async function getPaymentRequestsQueue(
  key,
  { sort, order, filters = [], currentPage = 1, currentPageSize = 20, viewAsGBLOC },
) {
  const operationPath = 'queues.getPaymentRequestsQueue';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  return makeGHCRequest(
    operationPath,
    { sort, order, page: currentPage, perPage: currentPageSize, viewAsGBLOC, ...paramFilters },
    { schemaKey: 'queuePaymentRequestsResult', normalize: false },
  );
}

export async function getShipmentsPaymentSITBalance(key, paymentRequestID) {
  return makeGHCRequest('paymentRequests.getShipmentsPaymentSITBalance', { paymentRequestID });
}

export function updateFinancialFlag({ moveID, ifMatchETag, body }) {
  const operationPath = 'move.setFinancialReviewFlag';
  // What is the schemakey and normalize for?
  return makeGHCRequest(
    operationPath,
    {
      moveID,
      'If-Match': ifMatchETag,
      body,
    },
    { normalize: false },
  );
}

export function updateMoveCloseoutOffice({ locator, ifMatchETag, body }) {
  const operationPath = 'move.updateCloseoutOffice';
  return makeGHCRequest(
    operationPath,
    {
      locator,
      'If-Match': ifMatchETag,
      body,
    },
    { normalize: false },
  );
}

export function updateServiceItemSITEntryDate({ mtoServiceItemID, body }) {
  const operationPath = 'mtoServiceItem.updateServiceItemSitEntryDate';
  return makeGHCRequest(
    operationPath,
    {
      mtoServiceItemID,
      body,
    },
    { normalize: false },
  );
}

export async function searchTransportationOffices(search) {
  const operationPath = 'transportationOffice.getTransportationOffices';
  return makeGHCRequest(operationPath, { search }, { normalize: false });
}

export async function searchTransportationOfficesOpen(search) {
  const operationPath = 'transportationOffice.getTransportationOfficesOpen';
  return makeGHCRequest(operationPath, { search }, { normalize: false });
}

export async function getGBLOCs() {
  const operationPath = 'transportationOffice.getTransportationOfficesGBLOCs';
  return makeGHCRequest(operationPath, {}, { normalize: false });
}

export const reviewShipmentAddressUpdate = async ({ shipmentID, ifMatchETag, body }) => {
  const operationPath = 'shipment.reviewShipmentAddressUpdate';
  const schemaKey = 'ShipmentAddressUpdate';
  const normalize = false;

  return makeGHCRequest(
    operationPath,
    {
      shipmentID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
};

export async function downloadPPMAOAPacket(ppmShipmentId) {
  return makeGHCRequestRaw('ppm.showAOAPacket', { ppmShipmentId });
}

export async function downloadPPMPaymentPacket(ppmShipmentId) {
  return makeGHCRequestRaw('ppm.showPaymentPacket', { ppmShipmentId });
}

export async function createOfficeAccountRequest({ body }) {
  return makeGHCRequest('officeUsers.createRequestedOfficeUser', { officeUser: body }, { normalize: false });
}

export async function createUploadForDocument(file, documentId) {
  return makeGHCRequest(
    'uploads.createUpload',
    {
      documentId,
      file,
    },
    {
      normalize: false,
    },
  );
}

export async function createUploadForAmdendedOrders(file, orderID) {
  return makeGHCRequest(
    'order.uploadAmendedOrders',
    {
      orderID,
      file,
    },
    {
      normalize: false,
    },
  );
}

export async function createUploadForSupportingDocuments(file, moveID) {
  return makeGHCRequest(
    'move.uploadAdditionalDocuments',
    {
      moveID,
      file,
    },
    {
      normalize: false,
    },
  );
}

export async function deleteUploadForDocument(uploadID, orderID) {
  return makeGHCRequest(
    'uploads.deleteUpload',
    {
      uploadID,
      orderID,
    },
    {
      normalize: false,
    },
  );
}

export async function searchCustomers(key, { sort, order, filters = [], currentPage = 1, currentPageSize = 20 }) {
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  return makeGHCRequest(
    'customer.searchCustomers',
    {
      body: {
        sort,
        order,
        page: currentPage,
        perPage: currentPageSize,
        ...paramFilters,
      },
    },
    { schemaKey: 'searchMovesResult', normalize: false },
  );
}

export async function patchPPMSIT({ ppmShipmentId, payload, eTag }) {
  return makeGHCRequest(
    'ppm.updatePPMSIT',
    {
      ppmShipmentId,
      'If-Match': eTag,
      body: payload,
    },
    {
      normalize: false,
    },
  );
}

export async function bulkDownloadPaymentRequest(paymentRequestID) {
  return makeGHCRequestRaw('paymentRequests.bulkDownload', { paymentRequestID });
}

export async function searchLocationByZipCityState(search) {
  return makeGHCRequest('addresses.getLocationByZipCityState', { search }, { normalize: false });
}

export async function dateSelectionIsWeekendHoliday(countryCode, date) {
  return makeGHCRequestRaw(
    'calendar.isDateWeekendHoliday',
    {
      countryCode,
      date,
    },
    { normalize: false },
  );
}

export async function updateAssignedOfficeUserForMove({ moveID, officeUserId, roleType }) {
  return makeGHCRequest('move.updateAssignedOfficeUser', {
    moveID,
    body: { officeUserId, roleType },
  });
}

export async function checkForLockedMovesAndUnlock(key, officeUserID) {
  return makeGHCRequestRaw('move.checkForLockedMovesAndUnlock', {
    officeUserID,
  });
}

export async function deleteAssignedOfficeUserForMove({ moveID, roleType }) {
  return makeGHCRequest('move.deleteAssignedOfficeUser', {
    moveID,
    body: { roleType },
  });
}

export async function getAllReServiceItems() {
  return makeGHCRequestRaw('reServiceItems.getAllReServiceItems', {}, { normalize: false });
}
