import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';

const getEntitlementsLabel = 'Entitlements.getEntitlements';

export function getEntitlements(moveTaskOrderId, label = getEntitlementsLabel) {
  const swaggerTag = 'entitlements.getEntitlements';
  return swaggerRequest(getGHCClient, swaggerTag, { moveTaskOrderId }, { label });
}
