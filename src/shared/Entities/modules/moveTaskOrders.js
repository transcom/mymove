import { swaggerRequest } from 'shared/Swagger/request';
import { getGHCClient } from 'shared/Swagger/api';

const getEntitlementsLabel = 'Entitlements.getEntitlements';

export function getEntitlements(moveTaskOrderID, label = getEntitlementsLabel) {
  const swaggerTag = 'Entitlements.getEntitlements';
  return swaggerRequest(getGHCClient, swaggerTag, { moveTaskOrderID }, { label });
}
