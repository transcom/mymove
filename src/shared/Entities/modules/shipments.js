import { denormalize } from 'normalizr';
import { shipments } from '../schema';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient, getPublicClient } from 'shared/Swagger/api';
import { isNull } from 'lodash';
import { getEntitlements } from 'shared/entitlements.js';

const approveShipmentLabel = 'Shipments.approveShipment';
export const getShipmentLabel = 'Shipments.getShipment';
export const getPublicShipmentLabel = 'Shipments.getPublicShipment';
const createShipmentLabel = 'Shipments.createShipment';
const updateShipmentLabel = 'shipments.updateShipment';
const updatePublicShipmentLabel = 'shipments.updatePublicShipment';
export const acceptPublicShipmentLabel = 'shipments.acceptShipment';
const transportPublicShipmentLabel = 'shipments.transportShipment';
const deliverPublicShipmentLabel = 'shipments.deliverShipment';
const completePmSurveyLabel = 'shipments.completePmSurvey';

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

export function approveShipment(shipmentId, shipmentApproveDate, label = approveShipmentLabel) {
  const swaggerTag = 'shipments.approveHHG';
  return swaggerRequest(
    getClient,
    swaggerTag,
    {
      shipmentId,
      approveShipmentPayload: {
        approve_date: shipmentApproveDate,
      },
    },
    { label },
  );
}

export function acceptShipment(shipmentId, label = acceptPublicShipmentLabel) {
  const swaggerTag = 'shipments.acceptShipment';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId }, { label });
}

export function transportShipment(shipmentId, payload, label = transportPublicShipmentLabel) {
  const swaggerTag = 'shipments.transportShipment';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId, payload }, label);
}

export function deliverShipment(shipmentId, payload, label = deliverPublicShipmentLabel) {
  const swaggerTag = 'shipments.deliverShipment';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId, payload }, label);
}

export function completePmSurvey(shipmentId, label = completePmSurveyLabel) {
  const swaggerTag = 'shipments.completePmSurvey';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId }, label);
}

export function calculateEntitlementsForShipment(state, shipmentId) {
  const shipment = selectShipment(state, shipmentId);
  const move = shipment.move || {};
  const serviceMember = shipment.service_member || {};
  const hasDependents = move.has_dependents;
  const spouseHasProGear = move.spouse_has_pro_gear;
  const rank = serviceMember.rank;

  if (isNull(hasDependents) || isNull(spouseHasProGear) || isNull(rank)) {
    return null;
  }
  return getEntitlements(rank, hasDependents, spouseHasProGear);
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
