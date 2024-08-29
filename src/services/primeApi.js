import Swagger from 'swagger-client';

import { makeSwaggerRequest, makeSwaggerRequestRaw, requestInterceptor, responseInterceptor } from './swaggerRequest';

let primeSimulatorClient = null;
let primeSimulatorClientV2 = null;
let primeSimulatorClientV3 = null;

// setting up the same config from Swagger/api.js
export async function getPrimeSimulatorClient() {
  if (!primeSimulatorClient) {
    primeSimulatorClient = await Swagger({
      url: '/prime/v1/swagger.yaml',
      requestInterceptor,
      responseInterceptor,
    });
  }
  return primeSimulatorClient;
}

// setting up the same config from Swagger/api.js
export async function getPrimeSimulatorClientV2() {
  if (!primeSimulatorClientV2) {
    primeSimulatorClientV2 = await Swagger({
      url: '/prime/v2/swagger.yaml',
      requestInterceptor,
      responseInterceptor,
    });
  }
  return primeSimulatorClientV2;
}

// setting up the same config from Swagger/api.js
export async function getPrimeSimulatorClientV3() {
  if (!primeSimulatorClientV3) {
    primeSimulatorClientV3 = await Swagger({
      url: '/prime/v3/swagger.yaml',
      requestInterceptor,
      responseInterceptor,
    });
  }
  return primeSimulatorClientV3;
}

export async function makePrimeSimulatorRequest(operationPath, params = {}, options = {}) {
  const client = await getPrimeSimulatorClient();
  return makeSwaggerRequest(client, operationPath, params, options);
}

export async function makePrimeSimulatorRequestV2(operationPath, params = {}, options = {}) {
  const client = await getPrimeSimulatorClientV2();
  return makeSwaggerRequest(client, operationPath, params, options);
}

export async function makePrimeSimulatorRequestV3(operationPath, params = {}, options = {}) {
  const client = await getPrimeSimulatorClientV3();
  return makeSwaggerRequest(client, operationPath, params, options);
}

export async function getPrimeSimulatorAvailableMoves() {
  const operationPath = 'moveTaskOrder.listMoves';
  return makePrimeSimulatorRequest(operationPath, {}, { schemaKey: 'listMoves', normalize: false });
}

export async function getPrimeSimulatorMove(key, locator) {
  return makePrimeSimulatorRequestV3('moveTaskOrder.getMoveTaskOrder', { moveID: locator }, { normalize: false });
}

export async function createPaymentRequest({ moveTaskOrderID, serviceItems }) {
  return makePrimeSimulatorRequest(
    'paymentRequest.createPaymentRequest',
    { body: { moveTaskOrderID, serviceItems } },
    { normalize: false },
  );
}

export async function completeCounseling({ moveTaskOrderID, ifMatchETag }) {
  return makePrimeSimulatorRequest(
    'moveTaskOrder.updateMTOPostCounselingInformation',
    { moveTaskOrderID, 'If-Match': ifMatchETag },
    { normalize: false },
  );
}

export async function deleteShipment({ mtoShipmentID }) {
  return makePrimeSimulatorRequest('mtoShipment.deleteMTOShipment', { mtoShipmentID }, { normalize: false });
}

export async function createUpload({ paymentRequestID, file, isWeightTicket }) {
  return makePrimeSimulatorRequest(
    'paymentRequest.createUpload',
    { paymentRequestID, file, isWeightTicket },
    { normalize: false },
  );
}

export async function createServiceRequestDocumentUpload({ mtoServiceItemID, file }) {
  return makePrimeSimulatorRequest(
    'mtoServiceItem.createServiceRequestDocumentUpload',
    { mtoServiceItemID, file },
    { normalize: false },
  );
}

export function createPrimeMTOShipmentV2({ normalize = false, schemaKey = 'mtoShipment', body }) {
  const operationPath = 'mtoShipment.createMTOShipment';
  return makePrimeSimulatorRequestV2(
    operationPath,
    {
      body,
    },
    { schemaKey, normalize },
  );
}

export function createPrimeMTOShipmentV3({ normalize = false, schemaKey = 'mtoShipment', body }) {
  const operationPath = 'mtoShipment.createMTOShipment';
  return makePrimeSimulatorRequestV3(
    operationPath,
    {
      body,
    },
    { schemaKey, normalize },
  );
}

export function updatePrimeMTOShipmentV2({
  mtoShipmentID,
  ifMatchETag,
  normalize = true,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'mtoShipment.updateMTOShipment';
  return makePrimeSimulatorRequestV2(
    operationPath,
    {
      mtoShipmentID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function updatePrimeMTOShipmentV3({
  mtoShipmentID,
  ifMatchETag,
  normalize = true,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'mtoShipment.updateMTOShipment';
  return makePrimeSimulatorRequestV3(
    operationPath,
    {
      mtoShipmentID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function createServiceItem({ body }) {
  return makePrimeSimulatorRequest('mtoServiceItem.createMTOServiceItem', { body: { ...body } }, { normalize: false });
}

export function updateMTOServiceItem({ mtoServiceItemID, eTag, body }) {
  return makePrimeSimulatorRequest(
    'mtoServiceItem.updateMTOServiceItem',
    { mtoServiceItemID, 'If-Match': eTag, body },
    { normalize: false },
  );
}

export function updateShipmentDestinationAddress({
  mtoShipmentID,
  ifMatchETag,
  body,
  schemaKey = 'mtoShipment',
  normalize = true,
}) {
  const operationPath = 'mtoShipment.updateShipmentDestinationAddress';
  return makePrimeSimulatorRequest(
    operationPath,
    {
      mtoShipmentID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function updatePrimeMTOShipmentAddress({
  mtoShipmentID,
  ifMatchETag,
  addressID,
  normalize = false,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'mtoShipment.updateMTOShipmentAddress';
  return makePrimeSimulatorRequest(
    operationPath,
    {
      mtoShipmentID,
      addressID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function updatePrimeMTOShipmentReweigh({
  mtoShipmentID,
  reweighID,
  ifMatchETag,
  normalize = false,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'mtoShipment.updateReweigh';
  return makePrimeSimulatorRequest(
    operationPath,
    {
      mtoShipmentID,
      reweighID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

// updatePrimeMTOShipmentStatus This function is used to by the Prime Simulator
// to send a cancellation request to the Prime API.
export function updatePrimeMTOShipmentStatus({
  mtoShipmentID,
  ifMatchETag,
  normalize = false,
  schemaKey = 'mtoShipment',
}) {
  const operationPath = 'mtoShipment.updateMTOShipmentStatus';
  // Default body is defined here as we can only send a status of CANCELED at
  // this time. See documentation here:
  // https://transcom.github.io/mymove-docs/api/prime#tag/mtoShipment/operation/updateMTOShipmentStatus
  const body = {
    status: 'CANCELED',
  };
  return makePrimeSimulatorRequest(
    operationPath,
    {
      mtoShipmentID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

// Sends api request for SIT Extension Request from Prime Sim
export function createSITExtensionRequest({ mtoShipmentID, normalize = false, schemaKey = 'mtoShipment', body }) {
  const operationPath = 'mtoShipment.createSITExtension';
  return makePrimeSimulatorRequest(
    operationPath,
    {
      mtoShipmentID,
      body,
    },
    { schemaKey, normalize },
  );
}

export async function downloadMoveOrder({ locator, type }) {
  const client = await getPrimeSimulatorClient();
  return makeSwaggerRequestRaw(client, 'moveTaskOrder.downloadMoveOrder', { locator, type });
}
