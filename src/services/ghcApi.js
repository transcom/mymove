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

export async function getPaymentRequestList() {
  const operationPath = 'paymentRequests.listPaymentRequests';
  const client = await getGHCClient();

  return makeSwaggerRequest(client, operationPath);
}

export async function getPaymentRequest(key, paymentRequestID) {
  const operationPath = 'paymentRequests.getPaymentRequest';
  const client = await getGHCClient();

  return makeSwaggerRequest(client, operationPath, { paymentRequestID });
}

export async function getMTOShipments(key, moveTaskOrderID) {
  const operationPath = 'mtoShipment.listMTOShipments';
  const client = await getGHCClient();

  return makeSwaggerRequest(client, operationPath, { moveTaskOrderID }, { schemaKey: 'mtoShipments' });
}

export async function getMTOServiceItems(key, moveTaskOrderID) {
  const operationPath = 'mtoServiceItem.listMTOServiceItems';
  const client = await getGHCClient();

  return makeSwaggerRequest(client, operationPath, { moveTaskOrderID }, { schemaKey: 'mtoServiceItems' });
}

export async function patchPaymentServiceItemStatus({
  moveTaskOrderID,
  paymentServiceItemID,
  status,
  ifMatchEtag,
  rejectionReason,
}) {
  const operationPath = 'paymentServiceItem.updatePaymentServiceItemStatus';
  const client = await getGHCClient();

  return makeSwaggerRequest(
    client,
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
