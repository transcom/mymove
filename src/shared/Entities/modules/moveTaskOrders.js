import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';
import { get } from 'lodash';

const getEntitlementsLabel = 'moveTaskOrder.getEntitlements';

export function getEntitlements(moveTaskOrderID, label = getEntitlementsLabel) {
  const swaggerTag = 'moveTaskOrder.getEntitlements';
  return swaggerRequest(getGHCClient, swaggerTag, { moveTaskOrderID }, { label });
}

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

export function selectMoveTaskOrder(state, moveTaskOrderId) {
  return get(state, `entities.moveTaskOrders.${moveTaskOrderId}`, {});
}

const getCustomerInfoOperation = 'Customer.getCustomerInfo';
export function getCustomerInfo(customerID, label = getCustomerInfoOperation) {
  return swaggerRequest(getGHCClient, getCustomerInfoOperation, { customerID }, { label });
}

const getAllCustomerMovesOperation = 'Customer.getAllCustomerMoves';
export function getAllCustomerMoves(label = getAllCustomerMovesOperation) {
  return swaggerRequest(getGHCClient, getAllCustomerMovesOperation, {}, { label });
}
