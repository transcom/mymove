import Swagger from 'swagger-client';
import * as Cookies from 'js-cookie';

import { makeSwaggerRequest } from './swaggerRequest';

// setting up the same config from Swagger/api.js
const requestInterceptor = (req) => {
  if (!req.loadSpec) {
    const token = Cookies.get('masked_gorilla_csrf');
    if (token) {
      req.headers['X-CSRF-Token'] = token;
    } else {
      // eslint-disable-next-line no-console
      console.warn('Unable to retrieve CSRF Token from cookie');
    }
  }
  return req;
};

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

export async function getPaymentRequestList() {
  return makeGHCRequest('paymentRequests.listPaymentRequests');
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
