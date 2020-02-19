import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient, getPrimeClient } from 'shared/Swagger/api';

const getAllMoveOrdersOperation = 'moveOrder.listMoveOrders';
export function getAllMoveOrders(label = getAllMoveOrdersOperation) {
  return swaggerRequest(getGHCClient, getAllMoveOrdersOperation, {}, { label });
}

const getAllMTOsOperation = 'moveTaskOrder.fetchMTOUpdates';
export function getAllMTOs(label = getAllMTOsOperation) {
  return swaggerRequest(getPrimeClient, getAllMTOsOperation, {}, { label });
}
