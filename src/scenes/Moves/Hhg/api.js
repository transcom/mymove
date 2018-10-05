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

export async function GetMoveDatesSummary(moveId, moveDate) {
  const client = await getClient();
  const response = await client.apis.moves.showMoveDatesSummary({
    moveId,
    move_date: formatDateString(moveDate),
  });
  checkResponse(response, 'failed to get hhg shipment due to server error');
  return response.body;
}
