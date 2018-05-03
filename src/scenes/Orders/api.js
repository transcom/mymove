import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreateOrders(orders) {
  const client = await getClient();
  const response = await client.apis.orders.createOrders({
    createOrdersPayload: orders,
  });
  checkResponse(response, 'failed to create a orders due to server error');
  return response.body;
}

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
  const response = await client.apis.orders.updateOrders({
    ordersId,
    updateOrdersPayload: ordersPayload,
  });
  checkResponse(response, 'failed to update orders due to server error');
  return response.body;
}

export async function ShowCurrentOrdersAPI(serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.service_members.showServiceMemberOrders({
    serviceMemberId,
  });
  checkResponse(response, 'failed to get current orders due to server error');
  return response.body;
}
