import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { get } from 'lodash';

const updateMoveTaskOrders = 'moveTaskOrder.updateMoveTaskOrderStatus';
export function updateMoveTaskOrderStatus(moveTaskOrderID, status, label = updateMoveTaskOrders) {
  const swaggerTag = 'moveTaskOrder.updateMoveTaskOrderStatus';
  return swaggerRequest(
    getGHCClient,
    swaggerTag,
    { moveTaskOrderID, body: { status } },
    { updateMoveTaskOrders },
    { label },
  );
}

const getMoveOrderLabel = 'moveOrder.getMoveOrder';
export function getMoveOrder(moveOrderID, label = getMoveOrderLabel) {
  console.log(moveOrderID);
  const swaggerTag = 'moveOrder.getMoveOrder';
  return swaggerRequest(getGHCClient, swaggerTag, { moveOrderID }, { getMoveOrder }, { label });
}

export function selectMoveOrder(state, moveOrderId) {
  return get(state, `entities.moveOrders.${moveOrderId}`, {});
}

const getMoveTaskOrderLabel = 'moveTaskOrder.updateMoveTaskOrderStatus';
export function getMoveTaskOrder(moveTaskOrderID, label = getMoveTaskOrderLabel) {
  const swaggerTag = 'moveTaskOrder.getMoveTaskOrder';
  return swaggerRequest(getGHCClient, swaggerTag, { moveTaskOrderID }, { getMoveTaskOrder }, { label });
}

export function selectMoveTaskOrder(state, moveTaskOrderId) {
  return get(state, `entities.moveTaskOrders.${moveTaskOrderId}`, {});
}

const getCustomerOperation = 'Customer.getCustomer';
export function getCustomer(customerID, label = getCustomerOperation) {
  return swaggerRequest(getGHCClient, getCustomerOperation, { customerID }, { label });
}

const getAllCustomerMovesOperation = 'Customer.getAllCustomerMoves';
export function getAllCustomerMoves(label = getAllCustomerMovesOperation) {
  return swaggerRequest(getGHCClient, getAllCustomerMovesOperation, {}, { label });
}
