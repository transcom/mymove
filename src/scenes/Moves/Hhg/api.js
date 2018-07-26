import { getClient, checkResponse } from 'shared/api';

export async function GetShipment(moveId) {
  const client = await getClient();
  const response = await client.apis.shipment.Shipment({
    moveId,
  });
  checkResponse(response, 'failed to get hhg due to server error');
  return response.body;
}
