import { getClient, checkResponse } from 'shared/Swagger/api';

export async function GetShipment(moveId, shipmentId) {
  const client = await getClient();
  const response = await client.apis.shipment.Shipment({
    moveId,
    shipmentId,
  });
  checkResponse(response, 'failed to get hhg shipment due to server error');
  return response.body;
}

export async function GetMoveDatesSummary(moveDate) {
  const client = await getClient();
  const response = await client.apis.moves.showMoveDatesSummary({
    moveDate,
  });
  checkResponse(response, 'failed to get hhg shipment due to server error');
  return response.body;
}
