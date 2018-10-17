import { getClient, checkResponse } from 'shared/Swagger/api';
import { formatDateString } from 'shared/utils';

export async function GetShipment(moveId, shipmentId) {
  const client = await getClient();
  const response = await client.apis.shipment.Shipment({
    moveId,
    shipmentId,
  });
  checkResponse(response, 'failed to get hhg shipment due to server error');
  return response.body;
}
