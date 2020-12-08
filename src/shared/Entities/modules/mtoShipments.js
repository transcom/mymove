import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient, getClient } from 'shared/Swagger/api';
import { selectMoveTaskOrders } from 'shared/Entities/modules/moveTaskOrders';
import { filter } from 'lodash';
import { denormalize } from 'normalizr';
import { mtoShipments } from '../schema';

/** REMAINING EXPORTS ARE USED BY OFFICE */

const mtoShipmentsSchemaKey = 'mtoShipments';
const getMTOShipmentsOperation = 'mtoShipment.listMTOShipments';
export function getMTOShipments(moveTaskOrderID, label = getMTOShipmentsOperation, schemaKey = mtoShipmentsSchemaKey) {
  return swaggerRequest(getGHCClient, getMTOShipmentsOperation, { moveTaskOrderID }, { label, schemaKey });
}

const mtoShipmentSchemaKey = 'mtoShipment';
const patchMTOShipmentStatusOperation = 'mtoShipment.patchMTOShipmentStatus';
export function patchMTOShipmentStatus(
  moveTaskOrderID,
  shipmentID,
  shipmentStatus,
  ifMatchETag,
  rejectionReason,
  label = patchMTOShipmentStatusOperation,
  schemaKey = mtoShipmentSchemaKey,
) {
  return swaggerRequest(
    getGHCClient,
    patchMTOShipmentStatusOperation,
    {
      moveTaskOrderID,
      shipmentID,
      'If-Match': ifMatchETag,
      body: { status: shipmentStatus, rejectionReason },
    },
    { label, schemaKey },
  );
}

const createMTOShipmentOperation = 'mtoShipment.createMTOShipment';
export function createMTOShipment(mtoShipment, label = createMTOShipmentOperation, schemaKey = mtoShipmentSchemaKey) {
  return swaggerRequest(getClient, createMTOShipmentOperation, { body: mtoShipment }, { label, schemaKey });
}

const updateMTOShipmentOperation = 'mtoShipment.updateMTOShipment';
export function updateMTOShipment(
  mtoShipmentId,
  mtoShipment,
  ifMatchETag,
  label = updateMTOShipmentOperation,
  schemaKey = mtoShipmentSchemaKey,
) {
  return swaggerRequest(
    getClient,
    updateMTOShipmentOperation,
    { mtoShipmentId, 'If-Match': ifMatchETag, body: mtoShipment },
    { label, schemaKey },
  );
}

const loadMTOShipmentsOperation = 'mtoShipment.listMTOShipments';
export function loadMTOShipments(
  moveTaskOrderID,
  label = loadMTOShipmentsOperation,
  schemaKey = mtoShipmentsSchemaKey,
) {
  return swaggerRequest(getClient, loadMTOShipmentsOperation, { moveTaskOrderID }, { label, schemaKey });
}

export function selectMTOShipments(state, moveOrderId) {
  const moveTaskOrders = selectMoveTaskOrders(state, moveOrderId);
  return filter(state.entities.mtoShipments, (item) => moveTaskOrders.find((mto) => mto.id === item.moveTaskOrderID));
}

export function selectMTOShipmentsByMoveId(state, moveId) {
  const mtoShipments = filter(state.entities.mtoShipments, (mtoShipment) => mtoShipment.moveTaskOrderID === moveId);
  return mtoShipments;
}

export function selectMTOShipmentForMTO(state, moveTaskOrderId) {
  const mtoShipment = Object.values(state.entities.mtoShipments).find(
    (mtoShipment) => mtoShipment.moveTaskOrderID === moveTaskOrderId,
  );
  return mtoShipment || {};
}

export function selectMTOShipmentById(state, id) {
  const emptyShipment = {};
  if (!id) return emptyShipment;
  return denormalize([id], mtoShipments, state.entities)[0] || emptyShipment;
}
