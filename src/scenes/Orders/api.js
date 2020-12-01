import { getClient, checkResponse } from 'shared/Swagger/api';
import { formatPayload } from 'shared/utils';

export async function GetOrders(ordersId) {
  const client = await getClient();
  const response = await client.apis.orders.showOrders({
    ordersId,
  });
  checkResponse(response, 'failed to get orders due to server error');
  return response.body;
}

export async function UpdateOrders(ordersId, ordersPayload) {
  const client = await getClient();
  const payloadDef = client.spec.definitions.CreateUpdateOrders;
  const response = await client.apis.orders.updateOrders({
    ordersId,
    updateOrders: formatPayload(ordersPayload, payloadDef),
  });
  checkResponse(response, 'failed to update orders due to server error');
  return response.body;
}
