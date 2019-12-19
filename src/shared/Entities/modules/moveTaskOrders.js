import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { get } from 'lodash';

const updateMoveTaskOrders = 'moveTaskOrder.updateMoveTaskOrderStatus';
export function updateMoveTaskOrderStatus(moveTaskOrderID, isAvailableToPrime, label = updateMoveTaskOrders) {
  const swaggerTag = 'moveTaskOrder.updateMoveTaskOrderStatus';
  return swaggerRequest(getGHCClient, swaggerTag, { moveTaskOrderID }, { updateMoveTaskOrders }, { label });
}

const getMoveOrderLabel = 'moveOrder.getMoveOrder';
export function getMoveOrder(moveOrderID, label = getMoveOrderLabel) {
  const swaggerTag = 'moveOrder.getMoveOrder';
  return swaggerRequest(getGHCClient, swaggerTag, { moveOrderID }, { label });
}

export function selectMoveOrder(state, moveOrderId) {
  return get(state, `entities.moveOrder.${moveOrderId}`, {});
}

const getMoveTaskOrderLabel = 'moveTaskOrder.getMoveTaskOrder';
export function getMoveTaskOrder(moveTaskOrderID, label = getMoveTaskOrderLabel) {
  const swaggerTag = 'moveTaskOrder.getMoveTaskOrder';
  return swaggerRequest(getGHCClient, swaggerTag, { moveTaskOrderID }, { label });
}

const getAllMoveTaskOrdersLabel = 'moveOrder.listMoveTaskOrders';
export function getAllMoveTaskOrders(moveOrderID, label = getAllMoveTaskOrdersLabel) {
  const swaggerTag = 'moveOrder.listMoveTaskOrders';
  return swaggerRequest(getGHCClient, swaggerTag, { moveOrderID }, { label });
}

const getCustomerOperation = 'customer.getCustomer';
export function getCustomer(customerID, label = getCustomerOperation) {
  return swaggerRequest(getGHCClient, getCustomerOperation, { customerID }, { label });
}

export function selectCustomer(state, customerId) {
  return get(state, `entities.customer.${customerId}`, {});
}
