import { getClient, checkResponse } from 'shared/api';

export async function ShipmentsIndex() {
  const client = await getClient();
  const response = await client.apis.shipments.indexShipments();
  checkResponse(response, 'failed to load shipments index due to server error');
  return response.body;
}

export async function CreateShipment(
  moveId,
  payload /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  const client = await getClient();
  const response = await client.apis.shipments.createShipment({
    moveId,
    createShipment: payload,
  });
  checkResponse(response, 'failed to create shipment due to server error');
  return response.body;
}

export async function UpdateShipment(
  moveId,
  shipmentId,
  payload /*shape: {pickup_address, requested_pickup_date, weight_estimate}*/,
) {
  const client = await getClient();
  const response = await client.apis.shipments.patchShipment({
    moveId,
    shipmentId,
    patchShipmentPayload: payload,
  });
  checkResponse(response, 'failed to update shipment due to server error');
  return response.body;
}
