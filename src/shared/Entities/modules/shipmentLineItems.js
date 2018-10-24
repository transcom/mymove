import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';
import { shipmentLineItems } from '../schema';
import { denormalize } from 'normalizr';

export function createShipmentLineItem(label, shipmentId, payload) {
  return swaggerRequest(getPublicClient, 'accessorials.createShipmentLineItem', { shipmentId, payload }, { label });
}

export function updateShipmentLineItem(label, shipmentLineItemId, payload) {
  return swaggerRequest(
    getPublicClient,
    'accessorials.updateShipmentLineItem',
    { shipmentLineItemId, payload },
    { label },
  );
}

export function deleteShipmentLineItem(label, shipmentLineItemId) {
  return swaggerRequest(getPublicClient, 'accessorials.deleteShipmentLineItem', { shipmentLineItemId }, { label });
}

export function approveShipmentLineItem(label, shipmentLineItemId) {
  return swaggerRequest(getPublicClient, 'accessorials.approveShipmentLineItem', { shipmentLineItemId }, { label });
}

export function getAllShipmentLineItems(label, shipmentId) {
  return swaggerRequest(getPublicClient, 'accessorials.getShipmentLineItems', { shipmentId }, { label });
}

export const selectShipmentLineItems = state => Object.values(state.entities.shipmentLineItems);

export const getShipmentLineItemsLabel = 'ShipmentLineItems.getAllShipmentLineItems';
export const createShipmentLineItemLabel = 'ShipmentLineItems.createShipmentLineItem';
export const deleteShipmentLineItemLabel = 'ShipmentLineItems.deleteShipmentLineItem';
export const approveShipmentLineItemLabel = 'ShipmentLineItems.approveShipmentLineItem';
export const updateShipmentLineItemLabel = 'ShipmentLineItems.updateShipmentLineItem';

export const selectShipmentLineItem = (state, id) => denormalize([id], shipmentLineItems, state.entities)[0];
