import { getClient, checkResponse } from 'shared/Swagger/api';

export async function ValidateEntitlement(moveId) {
  const client = await getClient();
  const response = await client.apis.entitlements.validateEntitlement({
    moveId,
  });
  checkResponse(
    response,
    'failed to validate entitlement due to server error.',
  );
  return response.body;
}
