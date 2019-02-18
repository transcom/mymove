import { denormalize } from 'normalizr';

import { shipments } from '../schema';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient, getPublicClient } from 'shared/Swagger/api';

const approveShipmentLabel = 'Shipments.approveShipment';
const completeShipmentLabel = 'Shipments.completeShipment';

export function createOrUpdateShipment(label, moveId, shipment, id) {
  if (id) {
    return updateShipment(label, id, shipment);
  } else {
    return createShipment(label, moveId, shipment);
  }
}

export function getShipment(label, shipmentId) {
  return swaggerRequest(getClient, 'shipments.getShipment', { shipmentId }, { label });
}

export function getPublicShipment(label, shipmentId) {
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

export function updatePublicShipment(
  shipmentId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
  label = 'shipments.updateShipment',
) {
  return swaggerRequest(getPublicClient, 'shipments.patchShipment', { shipmentId, update: shipment }, { label });
}

export function approveShipment(shipmentId) {
  const label = approveShipmentLabel;
  const swaggerTag = 'shipments.approveHHG';
  return swaggerRequest(getClient, swaggerTag, { shipmentId }, { label });
}

export function completeShipment(shipmentId) {
  const label = completeShipmentLabel;
  const swaggerTag = 'shipments.completeHHG';
  return swaggerRequest(getClient, swaggerTag, { shipmentId }, { label });
}

export function selectShipment(state, id) {
  if (!id) {
    return {};
  }
  return denormalize([id], shipments, state.entities)[0] || {};
}

export function selectShipmentStatus(state, id) {
  const shipment = selectShipment(state, id);
  return shipment.status;
}

export function selectShipmentForMove(state, moveId) {
  const shipment = Object.values(state.entities.shipments).find(shipment => shipment.move_id === moveId);
  return shipment || {};
}
