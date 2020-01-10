import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { selectMoveTaskOrders } from 'shared/Entities/modules/moveTaskOrders';
import { filter } from 'lodash';

const getMTOShipmentsOperation = 'mtoServiceItem.listMTOShipments';
const mtoShipmentsSchemaKey = 'mtoShipments';
export function getMTOShipments(moveTaskOrderID, label = getMTOShipmentsOperation, schemaKey = mtoShipmentsSchemaKey) {
  return swaggerRequest(getGHCClient, getMTOShipmentsOperation, { moveTaskOrderID }, { label, schemaKey });
}

export function selectMTOShipments(state, moveOrderId) {
  const moveTaskOrders = selectMoveTaskOrders(state, moveOrderId);

  return filter(state.entities.mtoShipments, item => moveTaskOrders.find(mto => mto.id === item.moveTaskOrderID));
}
