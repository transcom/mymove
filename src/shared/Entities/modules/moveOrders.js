import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { get } from 'lodash';

const getAllMoveOrdersOperation = 'moveOrder.listMoveOrders';
export function getAllMoveOrders(label = getAllMoveOrdersOperation) {
  return swaggerRequest(getGHCClient, getAllMoveOrdersOperation, {}, { label });
}

export const selectMoveOrder = (state, moveOrderId) => {
  return get(state, `entities.moveOrder.${moveOrderId}`, {});
};
