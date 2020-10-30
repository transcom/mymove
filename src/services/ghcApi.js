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

export async function getMoveOrder(key, moveOrderID) {
  return makeGHCRequest('moveOrder.getMoveOrder', { moveOrderID });
}

export async function getMoveTaskOrderList(key, moveOrderID) {
  return makeGHCRequest('moveOrder.listMoveTaskOrders', { moveOrderID });
}

export async function getMTOShipments(key, moveTaskOrderID) {
  return makeGHCRequest('mtoShipment.listMTOShipments', { moveTaskOrderID }, { schemaKey: 'mtoShipments' });
}

export async function getMTOServiceItems(key, moveTaskOrderID) {
  return makeGHCRequest('mtoServiceItem.listMTOServiceItems', { moveTaskOrderID }, { schemaKey: 'mtoServiceItems' });
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

export async function updateMoveOrder({ moveOrderID, ifMatchETag, body }) {
  const operationPath = 'moveOrder.updateMoveOrder';
  return makeGHCRequest(operationPath, { moveOrderID, 'If-Match': ifMatchETag, body });
}

export async function getMovesQueue(key, { filters = [] }) {
  const operationPath = 'queues.getMovesQueue';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  return makeGHCRequest(operationPath, { ...paramFilters }, { schemaKey: 'queueMovesResult', normalize: false });
}

export async function getPaymentRequestsQueue() {
  const operationPath = 'queues.getPaymentRequestsQueue';
  return makeGHCRequest(operationPath, {}, { schemaKey: 'queuePaymentRequestsResult', normalize: false });
}
