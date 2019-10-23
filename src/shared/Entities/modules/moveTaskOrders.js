import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';

const getEntitlementsLabel = 'Entitlements.getEntitlements';
const getCustomerInfoOperation = 'Customer.getCustomerInfo';

export function getEntitlements(moveTaskOrderID, label = getEntitlementsLabel) {
  const swaggerTag = 'Entitlements.getEntitlements';
  return swaggerRequest(getGHCClient, swaggerTag, { moveTaskOrderID }, { label });
}

export function getCustomerInfo(customerID, label = getCustomerInfoOperation) {
  return swaggerRequest(getGHCClient, getCustomerInfoOperation, { customerID }, { label });
}
