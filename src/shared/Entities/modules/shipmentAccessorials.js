import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { shipmentAccessorials } from '../schema';
import { denormalize } from 'normalizr';
import { orderBy } from 'lodash';
import { createSelector } from 'reselect';

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

export const selectShipmentAccessorials = createSelector(
  [
    state => {
      const sortedList = orderBy(
        Object.values(state.entities.shipmentAccessorials),
        ['status', 'approved_date', 'submitted_date'],
        ['asc', 'desc', 'desc'],
      );
      return sortedList;
    },
  ],
  value => value,
);

export const getShipmentAccessorialsLabel = 'ShipmentAccessorials.getAllShipmentAccessorials';
export const createShipmentAccessorialLabel = 'ShipmentAccessorials.createShipmentAccessorial';
export const deleteShipmentAccessorialLabel = 'ShipmentAccessorials.deleteShipmentAccessorial';
export const approveShipmentAccessorialLabel = 'ShipmentAccessorials.approveShipmentAccessorial';
export const updateShipmentAccessorialLabel = 'ShipmentAccessorials.updateShipmentAccessorial';

export const selectShipmentAccessorial = (state, id) => denormalize([id], shipmentAccessorials, state.entities)[0];
