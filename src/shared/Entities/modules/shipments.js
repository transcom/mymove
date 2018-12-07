import { denormalize } from 'normalizr';

import { shipments } from '../schema';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient, getPublicClient } from 'shared/Swagger/api';

export function createOrUpdateShipment(label, moveId, shipment, id) {
  if (id) {
    return updateShipment(label, id, shipment);
  } else {
    return createShipment(label, moveId, shipment);
  }
}

export function getShipment(label, shipmentId) {
  return swaggerRequest(getPublicClient, 'shipments.getShipment', { shipmentId }, { label });
}

export function createShipment(
  label,
  moveId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  return swaggerRequest(getClient, 'shipments.createShipment', { moveId, shipment }, { label });
}

export function updateShipment(
  label,
  shipmentId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  return swaggerRequest(getClient, 'shipments.patchShipment', { shipmentId, shipment }, { label });
}

export function selectShipment(state, id) {
  if (!id) {
    return null;
  }
  return denormalize([id], shipments, state.entities)[0];
}
