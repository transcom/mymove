import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient, getClient } from 'shared/Swagger/api';
import { selectMoveTaskOrders } from 'shared/Entities/modules/moveTaskOrders';
import { filter } from 'lodash';

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

// TODO - deprecate this selector when we refactor the wizard flow
export function selectMTOShipmentForMTO(state, moveTaskOrderId) {
  const mtoShipment = Object.values(state.entities.mtoShipments).find(
    (mtoShipment) => mtoShipment.moveTaskOrderID === moveTaskOrderId,
  );
  return mtoShipment || {};
}
