import { getClient, checkResponse } from 'shared/Swagger/api';

export async function ShipmentsIndex() {
  const client = await getClient();
  const response = await client.apis.shipments.indexShipments();
  checkResponse(response, 'failed to load shipments index due to server error');
  return response.body;
}
