import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { shipmentAccessorials } from '../schema';
import { denormalize } from 'normalizr';

export function createShipmentAccessorial(label, shipmentId, payload) {
  return swaggerRequest(getPublicClient, 'accessorials.createShipmentAccessorial', { shipmentId, payload }, { label });
}

export function updateShipmentAccessorial(label, shipmentAccessorialId, payload) {
  return swaggerRequest(
    getPublicClient,
    'accessorials.updateShipmentAccessorial',
    { shipmentAccessorialId, payload },
    { label },
  );
}

export function deleteShipmentAccessorial(label, shipmentAccessorialId) {
  return swaggerRequest(
    getPublicClient,
    'accessorials.deleteShipmentAccessorial',
    { shipmentAccessorialId },
    { label },
  );
}

export function approveShipmentAccessorial(label, shipmentAccessorialId) {
  return swaggerRequest(
    getPublicClient,
    'accessorials.approveShipmentAccessorial',
    { shipmentAccessorialId },
    { label },
  );
}

export function getAllShipmentAccessorials(label, shipmentId) {
  return swaggerRequest(getPublicClient, 'accessorials.getShipmentAccessorials', { shipmentId }, { label });
}

export const selectShipmentAccessorials = state => Object.values(state.entities.shipmentAccessorials);

export const getShipmentAccessorialsLabel = 'ShipmentAccessorials.getAllShipmentAccessorials';
export const createShipmentAccessorialLabel = 'ShipmentAccessorials.createShipmentAccessorial';

export const selectShipmentAccessorial = (state, id) => denormalize([id], shipmentAccessorials, state.entities)[0];
