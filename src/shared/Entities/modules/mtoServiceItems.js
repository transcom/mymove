import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { filter } from 'lodash';

const getMTOServiceItemsOperation = 'mtoServiceItem.listMTOServiceItems';
const mtoServiceItemsSchemaKey = 'mtoServiceItems';
export function getMTOServiceItems(
  moveTaskOrderID,
  label = getMTOServiceItemsOperation,
  schemaKey = mtoServiceItemsSchemaKey,
) {
  return swaggerRequest(getGHCClient, getMTOServiceItemsOperation, { moveTaskOrderID }, { label, schemaKey });
}

export function selectMTOServiceItems(state, moveTaskOrderId) {
  return filter(state.entities.mtoServiceItem, item => {
    return item.moveTaskOrderID === moveTaskOrderId;
  });
}
