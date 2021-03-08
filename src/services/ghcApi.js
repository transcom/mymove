import Swagger from 'swagger-client';

import { makeSwaggerRequest, requestInterceptor } from './swaggerRequest';

let ghcClient = null;

// setting up the same config from Swagger/api.js
export async function getGHCClient() {
  if (!ghcClient) {
    ghcClient = await Swagger({
      url: '/ghc/v1/swagger.yaml',
      requestInterceptor,
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

export async function getMoveTaskOrderList(key, orderID) {
  return makeGHCRequest('order.listMoveTaskOrders', { orderID });
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

export async function patchMTOServiceItemStatus({
  moveTaskOrderId,
  mtoServiceItemID,
  ifMatchEtag,
  status,
  rejectionReason,
}) {
  return makeGHCRequest(
    'mtoServiceItem.updateMTOServiceItemStatus',
    {
      moveTaskOrderID: moveTaskOrderId,
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

export async function updateMoveOrder({ orderID, ifMatchETag, body }) {
  const operationPath = 'order.updateMoveOrder';
  return makeGHCRequest(operationPath, { orderID, 'If-Match': ifMatchETag, body });
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

export function updateMTOShipmentStatus({
  moveTaskOrderID,
  shipmentID,
  shipmentStatus,
  ifMatchETag,
  rejectionReason,
  normalize = true,
  schemaKey = 'mtoShipment',
}) {
  const operationPath = 'mtoShipment.patchMTOShipmentStatus';
  return makeGHCRequest(
    operationPath,
    {
      moveTaskOrderID,
      shipmentID,
      'If-Match': ifMatchETag,
      body: { status: shipmentStatus, rejectionReason },
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
