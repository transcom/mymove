import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';

const getAllMoveOrdersOperation = 'moveOrder.listMoveOrders';
export function getAllMoveOrders(label = getAllMoveOrdersOperation) {
  return swaggerRequest(getGHCClient, getAllMoveOrdersOperation, {}, { label });
}
