import { denormalize } from 'normalizr';

import { shipments } from '../schema';
import { swaggerRequest } from 'shared/Swagger/request';

export function createOrUpdateShipment(label, moveId, shipment, id) {
  if (id) {
    return updateShipment(label, id, shipment);
  } else {
    return createShipment(label, moveId, shipment);
  }
}

export function getShipment(label, shipmentId, moveId) {
  return swaggerRequest(
    'shipments.getShipment',
    { moveId, shipmentId },
    { label },
  );
}

export function createShipment(
  label,
  moveId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  return swaggerRequest(
    'shipments.createShipment',
    { moveId, shipment },
    { label },
  );
}

export function updateShipment(
  label,
  shipmentId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  return swaggerRequest(
    'shipments.patchShipment',
    { shipmentId, shipment },
    { label },
  );
}

export function selectShipment(state, id) {
  if (!id) {
    return null;
  }
  return denormalize([id], shipments, state.entities)[0];
}
