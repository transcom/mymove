import { ensureClientIsLoaded, checkResponse } from 'shared/api';

export async function ShipmentsIndex() {
  const client = await ensureClientIsLoaded();
  const response = await client.apis.shipments.indexShipments();
  checkResponse(response, 'failed to load shipments index due to server error');
  return response.body;
}
