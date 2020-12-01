import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { get } from 'lodash';

export function selectMoveOrderList(state) {
  return Object.values(get(state, 'entities.moveOrders', {}));
}

const getAllMoveOrdersOperation = 'moveOrder.listMoveOrders';
export function getAllMoveOrders(label = getAllMoveOrdersOperation) {
  return swaggerRequest(getGHCClient, getAllMoveOrdersOperation, {}, { label });
}
