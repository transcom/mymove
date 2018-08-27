import { denormalize } from 'normalizr';
import { get } from 'lodash';

import { shipments } from '../schema';
import { swaggerRequest } from 'shared/api';

export function createOrUpdateShipment(label, moveId, shipment, id) {
  if (id) {
    return updateShipment(label, moveId, id, shipment);
  } else {
    return createShipment(label, moveId, shipment);
  }
}

export function getShipment(label, shipmentId, moveId) {
  return swaggerRequest(
    'shipments.getShipment',
    { moveId, shipmentId },
    {
      label,
      schemaKey: 'shipment',
    },
  );
}

export function createShipment(
  label,
  moveId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  return swaggerRequest(
    'shipments.createShipment',
    {
      moveId,
      shipment,
    },
    {
      label,
      schemaKey: 'shipment',
    },
  );
}

export function updateShipment(
  label,
  moveId,
  shipmentId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  return swaggerRequest(
    'shipments.patchShipment',
    {
      moveId,
      shipmentId,
      shipment,
    },
    {
      label,
      schemaKey: 'shipment',
    },
  );
}

export function selectShipment(state, id) {
  if (!id) {
    return null;
  }
  return denormalize([id], shipments, state.entities)[0];
}
