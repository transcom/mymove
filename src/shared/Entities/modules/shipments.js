import { denormalize } from 'normalizr';

import { shipments } from '../schema';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient, getPublicClient } from 'shared/Swagger/api';

const approveShipmentLabel = 'Shipments.approveShipment';
const completeShipmentLabel = 'Shipments.completeShipment';
export const getShipmentLabel = 'Shipments.getShipment';
const getPublicShipmentLabel = 'Shipments.getPublicShipment';
const createShipmentLabel = 'Shipments.createShipment';
const updateShipmentLabel = 'shipments.updateShipment';
const updatePublicShipmentLabel = 'shipments.updatePublicShipment';

export function createOrUpdateShipment(moveId, shipment, id, label) {
  if (id) {
    return updateShipment(id, shipment, label);
  } else {
    return createShipment(moveId, shipment, label);
  }
}

export function getShipment(shipmentId, label = getShipmentLabel) {
  return swaggerRequest(getClient, 'shipments.getShipment', { shipmentId }, { label });
}

export function getPublicShipment(shipmentId, label = getPublicShipmentLabel) {
  return swaggerRequest(getPublicClient, 'shipments.getShipment', { shipmentId }, { label });
}

export function createShipment(
  moveId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
  label = createShipmentLabel,
) {
  return swaggerRequest(getClient, 'shipments.createShipment', { moveId, shipment }, { label });
}

export function updateShipment(
  shipmentId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
  label = updateShipmentLabel,
) {
  return swaggerRequest(getClient, 'shipments.patchShipment', { shipmentId, shipment }, { label });
}

export function updatePublicShipment(
  shipmentId,
  shipment /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
  label = updatePublicShipmentLabel,
) {
  return swaggerRequest(getPublicClient, 'shipments.patchShipment', { shipmentId, update: shipment }, { label });
}

export function approveShipment(shipmentId, label = approveShipmentLabel) {
  const swaggerTag = 'shipments.approveHHG';
  return swaggerRequest(getClient, swaggerTag, { shipmentId }, { label });
}

export function completeShipment(shipmentId, label = completeShipmentLabel) {
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
