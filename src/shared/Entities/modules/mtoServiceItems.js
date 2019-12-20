import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { selectMoveTaskOrders } from 'shared/Entities/modules/moveTaskOrders';
import { get, filter } from 'lodash';

const getMTOServiceItemsOperation = 'mtoServiceItem.listMTOServiceItems';
const mtoServiceItemsSchemaKey = 'mtoServiceItems';
export function getMTOServiceItems(
  moveTaskOrderID,
  label = getMTOServiceItemsOperation,
  schemaKey = mtoServiceItemsSchemaKey,
) {
  return swaggerRequest(getGHCClient, getMTOServiceItemsOperation, { moveTaskOrderID }, { label, schemaKey });
}

export function selectMTOServiceItem(state, mtoServiceItemId) {
  return get(state, `entities.mtoServiceItem.${mtoServiceItemId}`, {});
}

export function selectMTOServiceItems(state, moveOrderId) {
  const moveTaskOrders = selectMoveTaskOrders(state, moveOrderId);

  return filter(state.entities.mtoServiceItem, item => moveTaskOrders.find(mto => mto.id === item.moveTaskOrderID));
}
