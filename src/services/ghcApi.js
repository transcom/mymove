import Swagger from 'swagger-client';

import { makeSwaggerRequest, requestInterceptor, responseInterceptor } from './swaggerRequest';

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

export async function getPaymentRequest(key, paymentRequestID) {
  return makeGHCRequest('paymentRequests.getPaymentRequest', { paymentRequestID });
}

export async function getMove(key, locator) {
  return makeGHCRequest('move.getMove', { locator }, { normalize: false });
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

export async function updateCustomerSupportRemarkForMove({ body, locator }) {
  return makeGHCRequest('customerSupportRemarks.updateCustomerSupportRemarkForMove', { body, locator });
}

export async function deleteCustomerSupportRemark({ customerSupportRemarkID }) {
  return makeGHCRequest('customerSupportRemarks.deleteCustomerSupportRemark', { customerSupportRemarkID });
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

export async function searchMoves(key, locator, dodID, customerName) {
  return makeGHCRequest(
    'move.searchMoves',
    { body: { locator, dodID, customerName } },
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

export async function updateOrder({ orderID, ifMatchETag, body }) {
  const operationPath = 'order.updateOrder';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag, body });
}

export async function counselingUpdateOrder({ orderID, ifMatchETag, body }) {
  const operationPath = 'order.counselingUpdateOrder';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag, body });
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

export async function updateCustomerInfo({ customerId, ifMatchETag, body }) {
  const operationPath = 'customer.updateCustomer';
  return makeGHCRequest(operationPath, { customerID: customerId, 'If-Match': ifMatchETag, body });
}

export async function updateMTOReviewedBillableWeights({ moveTaskOrderID, ifMatchETag }) {
  const operationPath = 'moveTaskOrder.UpdateMTOReviewedBillableWeightsAt';
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

export function updateMTOShipmentStatus({
  shipmentID,
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
  const operationPath = 'shipment.createSITExtensionAsTOO';
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

export async function getMovesQueue(key, { sort, order, filters = [], currentPage = 1, currentPageSize = 20 }) {
  const operationPath = 'queues.getMovesQueue';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  return makeGHCRequest(
    operationPath,
    { sort, order, page: currentPage, perPage: currentPageSize, ...paramFilters },
    { schemaKey: 'queueMovesResult', normalize: false },
  );
}

export async function getServicesCounselingQueue(
  key,
  { sort, order, filters = [], currentPage = 1, currentPageSize = 20 },
) {
  const operationPath = 'queues.getServicesCounselingQueue';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  return makeGHCRequest(
    operationPath,
    { sort, order, page: currentPage, perPage: currentPageSize, ...paramFilters },
    { schemaKey: 'queueMovesResult', normalize: false },
  );
}

export async function getPaymentRequestsQueue(
  key,
  { sort, order, filters = [], currentPage = 1, currentPageSize = 20 },
) {
  const operationPath = 'queues.getPaymentRequestsQueue';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  return makeGHCRequest(
    operationPath,
    { sort, order, page: currentPage, perPage: currentPageSize, ...paramFilters },
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
