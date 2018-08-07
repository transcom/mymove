import { getPublicClient, checkResponse } from 'shared/api';

// SHIPMENT QUEUE
export async function RetrieveShipmentsForTSP(queueType) {
  const queueToStatus = {
    new: ['AWARDED'],
    all: [],
  };
  /* eslint-disable security/detect-object-injection */
  const status =
    (queueType &&
      queueToStatus[queueType] &&
      queueToStatus[queueType].join(',')) ||
    '';
  /* eslint-enable security/detect-object-injection */
  const client = await getPublicClient();
  const response = await client.apis.shipments.indexShipments({
    status,
    limit: 25,
    offset: 1,
  });
  checkResponse(response, 'failed to retrieve moves due to server error');
  return response.body;
}

// SHIPMENT
export async function LoadShipment(shipmentId) {
  const client = await getPublicClient();
  const response = await client.apis.shipments.getShipment({
    shipment_uuid: shipmentId,
  });
  checkResponse(response, 'failed to load shipment due to server error');
  return response.body;
}

export async function AcceptShipment(shipmentId, originShippingAgent, destinationShippingAgent) {
  const client = await getPublicClient();
  const response = await client.apis.shipments.createShipmentAccept({
    shipment_uuid: shipmentId,
    originShippingAgent,
    destinationShippingAgent,
  });
  return response.body;
}

export async function RejectShipment(shipmentId, rejectReason) {
  const client = await getPublicClient();
  const response = await client.apis.shipments.createShipmentReject({
    shipment_uuid: shipmentId,
    rejectReason,
  });
  return response.body;
}
