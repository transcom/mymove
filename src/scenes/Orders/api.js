import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function GetOrders(serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.service_members.showServiceMemberOrders({
    serviceMemberId,
  });
  checkResponse(response, 'failed to get service member due to server error');
  return response.body;
}

export async function UpdateOrders(ordersId, ordersPayload) {
  const client = await getClient();
  const response = await client.apis.service_members.updateOrders({
    ordersId,
    updateOrdersPayload: ordersPayload,
  });
  checkResponse(response, 'failed to update orders due to server error');
  return response.body;
}
